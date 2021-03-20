package Controller

import (
	"BackEnd/Middleware"
	"BackEnd/Model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func Auth(c *gin.Context) {
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "无效的参数",
		})
	}
	if user.Username == "123" && user.Password == "456" {
		TokenString, _ := Middleware.GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"data": gin.H{"token": TokenString},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "failed",
	})
	return
}

func HoeHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
		"data": gin.H{
			"username": username,
		},
	})

}

func CheckedImage(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		studentID := c.PostForm("student_id")
		homeworkID := c.PostForm("homework_id")
		questionID := c.PostForm("question_id")
		score := c.PostForm("score")
		fmt.Println(studentID, homeworkID, questionID, score)
		scoreNum, err := strconv.Atoi(score)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		file.Filename = studentID + "-checked"
		dst := fmt.Sprintf("D:/FinProject/BackEnd/images/%s/%s/%s", homeworkID, questionID, file.Filename)
		c.SaveUploadedFile(file, dst)
		db.Model(&Model.HomeworkUploadRecord{}).Where("homework_id = ? AND question_id = ? AND student_id = ?", homeworkID, questionID, studentID).Updates(&Model.HomeworkUploadRecord{
			Score:    scoreNum,
			IsScored: true,
		})
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	}
}

func ImageSend(c *gin.Context) {
	homeworkid := c.Query("homework")
	questionid := c.Query("question")
	studentid := c.Query("student")
	c.File("./images/" + homeworkid + "/" + questionid + "/" + studentid)
}

func ImageSendTest(c *gin.Context) {
	homeworkid := c.Query("homework")
	questionid := c.Query("question")
	studentid := c.Query("student")
	c.File("./images/" + homeworkid + "/" + questionid + "/" + studentid + "-checked")
}

func HomeworkJudgeList(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var HomeworkJudgeLists Model.HomeworkUploadRecordsForSelects
		var checked int64
		homework := c.Query("homework_id")
		db.Raw("SELECT name,homework_upload_records.question_id,student_id,score,is_scored,content,question_max_score from homework_upload_records RIGHT JOIN questions ON homework_upload_records.question_id = questions.id LEFT JOIN students On students.id=homework_upload_records.student_id RIGHT JOIN homeworks on homeworks.id = homework_upload_records.homework_id AND homeworks.question_id = homework_upload_records.question_id where homework_id = ? AND homework_upload_records.deleted_at is null AND questions.deleted_at is null AND students.deleted_at is null AND homeworks.deleted_at is null", homework).Scan(&HomeworkJudgeLists)
		db.Model(&Model.HomeworkUploadRecord{}).Where("homework_id = ? AND is_scored = 1", homework).Count(&checked)
		c.JSON(http.StatusOK, gin.H{
			"lists":   HomeworkJudgeLists,
			"checked": checked,
		})
	}
}

func GetCourseList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var teacherId = c.Query("teacher_id")
		var courses Model.CourseForSelects
		db.Model(&Model.Course{}).Select("courses.id,course_name,credit,teacher_id,classroom_id,max_choose_num,selected_num,start_time,end_time,name").Joins("right join teachers on courses.teacher_id = teachers.id").Where("teacher_id = ?", teacherId).Group("course_name").Scan(&courses)
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
		})
	}
}

func GetHomeworkList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var courseId = c.Query("course_id")
		var homeworks Model.HomeworkForSelects
		db.Model(&[]Model.Homework{}).Where("course_id = ?", courseId).Group("id").Scan(&homeworks)
		//db.Raw("Select * from homeworks where course_id = ? AND deleted_at is NULL GROUP BY id", courseId).Scan(&homeworks)
		c.JSON(http.StatusOK, gin.H{
			"homeworks": homeworks,
			"time":      time.Now(),
		})
	}
}

func DeleteHomework(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var id = c.Query("id")
		db.Where("id = ?", id).Delete(&Model.Homework{})
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}

func PostQuestion(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var content = c.PostForm("content")
		var question Model.Question
		var newQuestion Model.Question
		db.Last(&question)
		fmt.Println(question)
		newQuestion.ID = question.ID + 1
		newQuestion.Content = content
		db.Create(&newQuestion)
		c.JSON(http.StatusOK, gin.H{
			"id": newQuestion.ID,
		})
	}
}

func GetContentList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var questions Model.QuestionForSelects
		db.Raw("select * from questions").Scan(&questions)
		c.JSON(http.StatusOK, gin.H{
			"questions": questions,
		})
	}
}

func GetNewHomeworkID(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var homework Model.Homework
		db.Unscoped().Order("id desc").Last(&homework)
		homework.ID++
		db.Create(&homework)
		c.JSON(http.StatusOK, gin.H{
			"id": homework.ID,
		})
	}
}

func CreateHomework(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var questions Model.Homeworks
		if err := c.ShouldBindJSON(&questions); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Unscoped().Where("id = ?", questions[0].ID).Delete(&Model.Homework{})
		id := utils.ToString(questions[0].ID)
		for _, item := range questions {
			err := os.MkdirAll("./images/"+id+"/"+utils.ToString(item.QuestionID), 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		db.Create(&questions)
	}
}

func GetQuestionList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var questions Model.QuestionForStudents
		var records Model.HomeworkUploadRecordsForSelects
		homeworkID := c.Query("homework_id")
		studentId := c.Query("student_id")
		db.Raw("select dead_line,homework_title,question_max_score,content,question_id from homeworks right join questions on homeworks.question_id = questions.id where homeworks.id = ?", homeworkID).Scan(&questions)
		if studentId != "" {
			db.Raw("select * from homework_upload_records where homework_id = ? AND student_id = ?", homeworkID, studentId).Scan(&records)
			for _, value := range records {
				for i, v := range questions {
					fmt.Println(v.QuestionId == value.QuestionID)
					if v.QuestionId == value.QuestionID {
						questions[i].Uploaded = true
						break
					}
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"questions": questions,
		})
	}
}

func PostHomework(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		var record Model.HomeworkUploadRecord
		studentID := c.PostForm("student_id")
		homeworkID := c.PostForm("homework_id")
		questionID := c.PostForm("question_id")
		record.StudentID, err = strconv.Atoi(studentID)
		record.HomeworkID, err = strconv.Atoi(homeworkID)
		record.QuestionID, err = strconv.Atoi(questionID)
		record.Score = 0
		record.IsUpload = true
		fmt.Println(record.StudentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		file.Filename = studentID
		dst := fmt.Sprintf("./images/%s/%s/%s", homeworkID, questionID, file.Filename)
		_ = c.SaveUploadedFile(file, dst)
		db.Debug().Create(&record)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	}
}

func GetReviewList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reviwes Model.RecordForReviews
		homeworkID := c.Query("homework_id")
		studentId := c.Query("student_id")
		db.Raw("select * from homework_upload_records right join homeworks h on h.id = homework_upload_records.homework_id AND h.question_id = homework_upload_records.question_id right join questions q on q.id = h.question_id where homework_id =? AND student_id = ?", homeworkID, studentId).Scan(&reviwes)
		c.JSON(http.StatusOK, gin.H{
			"reviews": reviwes,
		})
	}
}

func GetClassesList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		level := c.Query("level")
		classes := &Model.CourseForSelects{}
		teachers := &Model.TeacherForSelects{}
		faculties := Model.Faculties{}
		directions := Model.DirectionForSelects{}
		classrooms := &Model.ClassroomForSelects{}
		sourceMap := make(map[int]Model.DirectionForSelects)
		fmt.Println(level)
		if level == "2" {
			tId := c.Query("teacher_id")
			db.Find("teacher_id = ?", tId, &Model.Course{}).Scan(&classes)
			c.JSON(http.StatusOK, gin.H{
				"classes": classes,
			})
		} else {
			db.Debug().Find(&[]Model.Teacher{}).Scan(&teachers)
			db.Find(&[]Model.Course{}).Order("id asc").Scan(&classes)
			db.Find(&[]Model.Faculty{}).Scan(&faculties)
			db.Find(&[]Model.Direction{}).Scan(&directions)
			db.Find(&[]Model.Classroom{}).Scan(&classrooms)
			for _, v := range directions {
				sourceMap[v.FacultyID] = append(sourceMap[v.FacultyID], v)
			}
			fmt.Println(teachers)
			c.JSON(http.StatusOK, gin.H{
				"teachers":   teachers,
				"classes":    classes,
				"faculties":  faculties,
				"directions": sourceMap,
				"classrooms": classrooms,
			})
		}
	}
}

func PostNewClass(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		course := &Model.Course{}
		courseTemp := &Model.Course{}
		course.FacultyID, _ = strconv.Atoi(c.PostForm("faculty_id"))
		course.CourseName = c.PostForm("name")
		course.Credit = c.PostForm("credit")
		course.MaxChooseNum, _ = strconv.Atoi(c.PostForm("max_choose_num"))
		course.WeekTime, _ = strconv.Atoi(c.PostForm("week_time"))
		course.StartTime, _ = strconv.Atoi(c.PostForm("start_time"))
		course.EndTime, _ = strconv.Atoi(c.PostForm("end_time"))
		course.TeacherID, _ = strconv.Atoi(c.PostForm("teacher_id"))
		course.StartWeek, _ = strconv.Atoi(c.PostForm("start_week"))
		course.DirectionID, _ = strconv.Atoi(c.PostForm("direction_id"))
		course.EndWeek, _ = strconv.Atoi(c.PostForm("end_week"))
		course.ClassroomID, _ = strconv.Atoi(c.PostForm("classroom_id"))
		flag := c.PostForm("flag")
		if flag == "copy" {
			course.CopyFlag = true
			db.Select("id").Where("course_name = ?", course.CourseName).Last(&courseTemp)
		} else {
			course.CopyFlag = false
			db.Debug().Unscoped().Select("id").Order("id desc").Limit(1).Find(&courseTemp)
			courseTemp.ID++
		}
		db.Debug().Unscoped().Select("record_id").Order("record_id desc").Limit(1).Find(&course)
		course.RecordID++
		course.ID = courseTemp.ID
		fmt.Println(course)
		db.Create(&course)
		c.JSON(http.StatusOK, gin.H{
			"snackbar": "setSuccess",
			"msg":      "添加课程成功",
		})

		//db.Create(&course)
	}
}

func DeleteClass(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var id = c.Query("course_id")
		result := db.Where("id = ?", id).Delete(&Model.Course{})
		fmt.Println(result.RowsAffected)
		c.JSON(http.StatusOK, gin.H{
			"snackbar": "setSuccess",
			"msg":      "删除课程成功！",
			"id":       id,
		})
	}
}

func UploadClass(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		classes := Model.CourseForSelects{}
		err := c.ShouldBindJSON(&classes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"snackbar": "setError",
				"msg":      err,
			})
			return
		}
		for _, v := range classes {
			fmt.Println(v)
			db.Debug().Updates(&Model.Course{
				ID:           v.ID,
				CourseName:   v.CourseName,
				Credit:       v.Credit,
				TeacherID:    v.TeacherID,
				ClassroomID:  v.ClassroomID,
				MaxChooseNum: v.MaxChooseNum,
				SelectedNum:  v.SelectedNum,
				StartTime:    v.StartTime,
				EndTime:      v.EndTime,
				WeekTime:     v.WeekTime,
				FacultyID:    v.FacultyID,
				DirectionID:  v.DirectionID,
				StartWeek:    v.StartWeek,
				EndWeek:      v.EndWeek,
			})
		}
	}
}

func PostApply(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		apply := &Model.ApplyForCourseChange{}
		c.ShouldBindJSON(&apply)
		db.Create(&apply)
		c.JSON(http.StatusOK, gin.H{
			"snackbar": "setSuccess",
			"msg":      "提交请求成功！",
		})
	}
}

func GetApply(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		level, _ := strconv.Atoi(c.Query("level"))
		classrooms := &Model.ClassroomForSelects{}
		db.Find(&[]Model.Classroom{}).Scan(&classrooms)
		applies := &Model.Applies{}
		if level < 3 {
			teacher_id := c.Query("teacher_id")
			db.Model(&Model.ApplyForCourseChange{}).Select("apply_for_course_changes.*,t.name,c2.course_name").Joins("left join teachers t on t.id = apply_for_course_changes.teacher_id left join courses c2 on apply_for_course_changes.course_id = c2.id ").Where("c2.copy_flag = 0 AND c.teacher_id = ?", teacher_id).Scan(&applies)
		} else {
			db.Model(&Model.ApplyForCourseChange{}).Select("apply_for_course_changes.*,t.name,c2.course_name").Joins("left join teachers t on t.id = apply_for_course_changes.teacher_id left join courses c2 on apply_for_course_changes.course_id = c2.id ").Where("copy_flag = 0").Scan(&applies)
		}
		c.JSON(http.StatusOK, gin.H{
			"applies":    applies,
			"classrooms": classrooms,
		})
	}
}

func OperateApply(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		result, _ := strconv.Atoi(c.PostForm("result"))
		id := c.PostForm("id")
		db.Debug().Model(&Model.ApplyForCourseChange{}).Where("id = ?", id).Update("result", result)
		if result == 1 {
			apply := &Model.ApplyForCourseChange{}
			course := &Model.Course{}
			tempCourse :=&Model.Course{}
			db.Where("id = ?", id).Find(&apply)
			db.Where("id = ? AND week_time = ? AND ( IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?) )", apply.CourseID, apply.BeforeWeekTime, apply.BeforeStartTime, apply.BeforeStartTime, apply.BeforeEndTime, apply.BeforeEndTime).Find(&course)
			if course.StartWeek == apply.BeforeStartWeek || course.EndWeek == apply.BeforeEndWeek {
				if course.StartWeek == apply.BeforeStartWeek {
					db.Debug().Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("start_week", apply.BeforeEndWeek+1)
				} else {
					db.Debug().Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("end_week", apply.BeforeStartWeek-1)
				}
			} else {
				fmt.Println(course)
				courseOBJ1 := course
				courseOBJ1.StartWeek = apply.BeforeEndWeek + 1
				courseOBJ1.EndWeek = course.EndWeek
				courseOBJ1.CopyFlag = true
				db.Debug().Unscoped().Select("record_id").Order("record_id desc").Limit(1).Find(&tempCourse)
				db.Debug().Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("end_week", apply.BeforeStartWeek-1)
				courseOBJ1.RecordID = tempCourse.RecordID + 1
				db.Create(&courseOBJ1)
			}
			courseOBJ2 := course
			courseOBJ2.StartWeek = apply.AfterStartWeek
			courseOBJ2.EndWeek = apply.AfterEndWeek
			courseOBJ2.CopyFlag = true
			courseOBJ2.WeekTime = apply.AfterWeekTime
			courseOBJ2.StartTime = apply.AfterStartTime
			courseOBJ2.EndTime = apply.AfterEndTime
			course = &Model.Course{}
			db.Debug().Unscoped().Select("record_id").Order("record_id desc").Limit(1).Find(&course)
			courseOBJ2.RecordID = course.RecordID + 1
			db.Create(courseOBJ2)
		}
		c.JSON(http.StatusOK, 0)
	}
}

func ValidClassrooms(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		course := &Model.Course{}
		tempClassrooms := &Model.ClassroomForSelects{}
		classrooms := &Model.ClassroomForSelects{}
		course.ID,_ = strconv.Atoi(c.PostForm("course_id"))
		course.WeekTime, _ = strconv.Atoi(c.PostForm("week_time"))
		course.StartTime, _ = strconv.Atoi(c.PostForm("start_time"))
		course.EndTime, _ = strconv.Atoi(c.PostForm("end_time"))
		course.StartWeek, _ = strconv.Atoi(c.PostForm("start_week"))
		course.EndWeek, _ = strconv.Atoi(c.PostForm("end_week"))
		course.ClassroomID, _ = strconv.Atoi(c.PostForm("classroom_id"))
		timeNode := c.PostForm("time")
		if timeNode == "after" {
			db.Debug().Model(&Model.Course{}).Select("courses.classroom_id as id,c.name").Joins("left join classrooms c on c.id = courses.classroom_id").Where("(IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ?",course.StartTime,course.StartTime,course.EndTime,course.EndTime,course.StartWeek,course.StartWeek,course.EndWeek,course.EndWeek,course.WeekTime).Group("c.name").Scan(&tempClassrooms)
			db.Debug().Model(&Model.Classroom{}).Not(tempClassrooms).Scan(&classrooms)
		}else{
			db.Debug().Model(&Model.Course{}).Select("courses.classroom_id as id,c.name").Joins("left join classrooms c on c.id = courses.classroom_id").Where(" courses.id = ? AND start_time = ? and end_time = ? AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ?",course.ID,course.StartTime,course.EndTime,course.StartWeek,course.StartWeek,course.EndWeek,course.EndWeek,course.WeekTime).Group("c.name").Scan(&classrooms)
		}
		c.JSON(http.StatusOK,gin.H{
			"classrooms":classrooms,
		})
	}
}

func GetClassSheet (db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var classSheet [7][14]string
		courses:= Model.ClassSheets{}
		teacher_id := c.Query("teacher_id")
		week := c.Query("week")
		db.Debug().Model(&Model.Course{}).Select("courses.*,t.name as teacher_name,c.name as classroom_name").Joins("left join teachers t on t.id = courses.teacher_id").Joins("left join classrooms c on c.id = courses.classroom_id ").Where("teacher_id = ? AND start_week <= ? AND end_week >= ?",teacher_id,week,week).Scan(&courses)
		for _,v:=range courses{
			for i:= v.StartTime;i<=v.EndTime;i++ {
				week_time := v.WeekTime
				course_name := v.CourseName
				course_classroom := v.ClassroomName
				teacher_name := v.TeacherName
				start_week := utils.ToString(v.StartWeek)
				end_week := utils.ToString(v.EndWeek)
				classSheet[week_time-1][i-1] = course_name+"\n教师："+teacher_name+"\n教室："+course_classroom+"\n第"+start_week+"周-第"+end_week+"周"
			}
		}
		c.JSON(http.StatusOK,gin.H{
			"classSheet":classSheet,
		})
	}
}

func ReturnOK(c *gin.Context) {
	c.JSON(http.StatusOK, 0)
}
