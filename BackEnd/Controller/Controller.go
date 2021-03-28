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
		var studentId = c.Query("student_id")
		var courses Model.CourseForSelects
		if studentId == "" {
			db.Debug().Model(&Model.Course{}).Select("courses.id,course_name,credit,teacher_id,classroom_id,max_choose_num,selected_num,start_time,end_time,name,courses.record_id").Joins("right join teachers on courses.teacher_id = teachers.id").Where("teacher_id = ? AND copy_flag = 0", teacherId).Group("record_id").Scan(&courses)
		}else {
			db.Debug().Model(&Model.Student2Course{}).Select("courses.id,course_name,credit,teacher_id,classroom_id,max_choose_num,selected_num,start_time,end_time,name,courses.record_id").Joins("left join courses on courses.record_id = student2_courses.record_id right join teachers on courses.teacher_id = teachers.id").Where("student_id = ?",studentId).Scan(&courses)
		}
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
		})
	}
}

func GetHomeworkList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var courseId = c.Query("course_id")
		var recordId = c.Query("record_id")
		var homeworks Model.HomeworkForSelects
		db.Model(&[]Model.Homework{}).Where("course_id = ? AND record_id = ?", courseId, recordId).Group("id").Scan(&homeworks)
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
		db.Raw("select dead_line,homework_title,question_max_score,content,question_id,record_id from homeworks right join questions on homeworks.question_id = questions.id where homeworks.id = ?", homeworkID).Scan(&questions)
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
		recordID := c.PostForm("record_id")
		fmt.Println(recordID)
		record.StudentID, err = strconv.Atoi(studentID)
		record.HomeworkID, err = strconv.Atoi(homeworkID)
		record.QuestionID, err = strconv.Atoi(questionID)
		record.RecordID,err = strconv.Atoi(recordID)
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
		specialties := Model.DirectionToSpecialtyForSelects{}
		directions := Model.DirectionToSpecialtyForSelects{}
		classrooms := &Model.ClassroomForSelects{}
		facultyToSpecialty := make(map[int]Model.DirectionToSpecialtyForSelects)
		SpecialtyToDirection := make(map[int]Model.DirectionToSpecialtyForSelects)
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
			db.Find(&[]Model.DirectionToSpecialty{}).Scan(&specialties)
			db.Find(&[]Model.DirectionToSpecialty{}).Scan(&directions)
			db.Find(&[]Model.Classroom{}).Scan(&classrooms)
			for _, v := range specialties {
				facultyToSpecialty[v.FacultyID] = append(facultyToSpecialty[v.FacultyID], v)
			}
			for _, v := range directions {
				SpecialtyToDirection[v.SpecialtyID] = append(SpecialtyToDirection[v.SpecialtyID], v)
			}
			fmt.Println(teachers)
			c.JSON(http.StatusOK, gin.H{
				"teachers":    teachers,
				"classes":     classes,
				"faculties":   faculties,
				"specialties": facultyToSpecialty,
				"directions":  SpecialtyToDirection,
				"classrooms":  classrooms,
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
		course.SpecialtyID, _ = strconv.Atoi(c.PostForm("specialty_id"))
		course.EndWeek, _ = strconv.Atoi(c.PostForm("end_week"))
		course.ClassroomID, _ = strconv.Atoi(c.PostForm("classroom_id"))
		course.CopyFlag, _ = strconv.Atoi(c.PostForm("copy_flag"))
		flag := c.PostForm("flag")
		if flag == "copy" || flag == "course" {
			db.Select("id").Where("course_name = ?", course.CourseName).Last(&courseTemp)
		} else {
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
		var id = c.Query("record_id")
		result := db.Where("record_id = ? OR copy_flag = ?", id, id).Delete(&Model.Course{})
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
			db.Debug().Where("record_id = ?", v.RecordID).Updates(&Model.Course{
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
				SpecialtyID:  v.SpecialtyID,
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
			teacherId := c.Query("teacher_id")
			db.Model(&Model.ApplyForCourseChange{}).Select("apply_for_course_changes.*,t.name,c2.course_name").Joins("left join teachers t on t.id = apply_for_course_changes.teacher_id left join courses c2 on apply_for_course_changes.course_id = c2.id ").Where("c2.copy_flag = 0 AND c.teacher_id = ?", teacherId).Scan(&applies)
		} else {
			db.Model(&Model.ApplyForCourseChange{}).Select("apply_for_course_changes.*,t.name,c2.course_name").Joins("left join teachers t on t.id = apply_for_course_changes.teacher_id left join courses c2 on apply_for_course_changes.course_id = c2.id ").Where("copy_flag = 0").Group("apply_for_course_changes.id").Scan(&applies)
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
		id := c.PostForm("id") //course_id
		db.Debug().Model(&Model.ApplyForCourseChange{}).Where("id = ?", id).Update("result", result)
		if result == 1 {
			apply := &Model.ApplyForCourseChange{}
			course := &Model.Course{}
			tempCourse := &Model.Course{}
			db.Where("id = ?", id).Find(&apply)
			db.Where("id = ? AND week_time = ? AND ( IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?) )", apply.CourseID, apply.BeforeWeekTime, apply.BeforeStartTime, apply.BeforeStartTime, apply.BeforeEndTime, apply.BeforeEndTime).Find(&course)

			if course.StartWeek == apply.BeforeStartWeek && course.EndWeek == apply.BeforeEndWeek {
				db.Debug().Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("start_time", apply.AfterStartTime).Update("end_time", apply.AfterEndWeek)
				c.JSON(http.StatusOK, 0)
				return
			} else if course.StartWeek == apply.BeforeStartWeek || course.EndWeek == apply.BeforeEndWeek {
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
				courseOBJ1.CopyFlag = course.RecordID
				db.Debug().Unscoped().Select("record_id").Order("record_id desc").Limit(1).Find(&tempCourse)
				db.Debug().Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("end_week", apply.BeforeStartWeek-1)
				courseOBJ1.RecordID = tempCourse.RecordID + 1
				db.Create(&courseOBJ1)
			}
			courseOBJ2 := course
			courseOBJ2.StartWeek = apply.AfterStartWeek
			courseOBJ2.EndWeek = apply.AfterEndWeek
			courseOBJ2.CopyFlag = course.RecordID
			courseOBJ2.WeekTime = apply.AfterWeekTime
			courseOBJ2.StartTime = apply.AfterStartTime
			courseOBJ2.EndTime = apply.AfterEndTime
			courseOBJ2.CopyFlag = course.RecordID
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
		tempCourse := &Model.CourseForSelect{}
		tempClassrooms1 := &Model.ClassroomForSelects{}
		tempClassrooms2 := &Model.ClassroomForSelects{}
		classrooms := &Model.ClassroomForSelects{}
		course.RecordID, _ = strconv.Atoi(c.PostForm("record_id"))
		course.ID, _ = strconv.Atoi(c.PostForm("course_id"))
		course.WeekTime, _ = strconv.Atoi(c.PostForm("week_time"))
		course.StartTime, _ = strconv.Atoi(c.PostForm("start_time"))
		course.EndTime, _ = strconv.Atoi(c.PostForm("end_time"))
		course.StartWeek, _ = strconv.Atoi(c.PostForm("start_week"))
		course.EndWeek, _ = strconv.Atoi(c.PostForm("end_week"))
		course.ClassroomID, _ = strconv.Atoi(c.PostForm("classroom_id"))
		course.TeacherID, _ = strconv.Atoi(c.PostForm("teacher_id"))
		course.FacultyID, _ = strconv.Atoi(c.PostForm("faculty_id"))
		course.SpecialtyID, _ = strconv.Atoi(c.PostForm("specialty_id"))
		course.DirectionID, _ = strconv.Atoi(c.PostForm("direction_id"))
		timeNode := c.PostForm("time")
		flag := c.PostForm("flag")
		if timeNode == "after" || timeNode == "applyAfter" {
			if timeNode == "applyAfter" {
				db.Debug().Model(&Model.Course{}).Select("faculty_id,direction_id,specialty_id,teacher_id,record_id").Joins("left join classrooms c on c.id = courses.classroom_id").Where(" courses.id = ? AND start_time = ? and end_time = ? AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ?", course.ID, course.StartTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime).Group("c.name").Scan(&tempCourse)
				course.FacultyID = tempCourse.FacultyID
				course.DirectionID = tempCourse.DirectionID
				course.SpecialtyID = tempCourse.SpecialtyID
				course.TeacherID = tempCourse.TeacherID
				course.RecordID = tempCourse.RecordID
				if flag != "not" {
					if course.DirectionID == 0 && course.SpecialtyID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = 0))) AND record_id <> ?) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.RecordID).Scan(&tempClassrooms1)
					} else if course.DirectionID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = 0))) AND record_id <> ? ) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.RecordID).Scan(&tempClassrooms1)
					} else if course.SpecialtyID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = ? OR direction_id = 0))) AND record_id <> ? ) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.DirectionID, course.RecordID).Scan(&tempClassrooms1)
					} else {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = ? OR direction_id = 0))) AND record_id <> ? ) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.DirectionID, course.RecordID).Scan(&tempClassrooms1)
					}
				} else {
					if course.DirectionID == 0 && course.SpecialtyID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID).Scan(&tempClassrooms1)
					} else if course.DirectionID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID).Scan(&tempClassrooms1)
					} else if course.SpecialtyID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.DirectionID).Scan(&tempClassrooms1)
					} else {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.DirectionID).Scan(&tempClassrooms1)
					}
				}
			} else if course.DirectionID == 0 && course.SpecialtyID == 0 {
				db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID).Scan(&tempClassrooms1)
			} else if course.DirectionID == 0 {
				db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID).Scan(&tempClassrooms1)
			} else if course.SpecialtyID == 0 {
				db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.DirectionID).Scan(&tempClassrooms1)
			} else {
				db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.DirectionID).Scan(&tempClassrooms1)
			}
			db.Debug().Model(&Model.Course{}).Select("courses.classroom_id as id,c.name").Joins("left join classrooms c on c.id = courses.classroom_id").Where("(IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ?", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime).Group("c.name").Scan(&tempClassrooms2)
			db.Debug().Model(&Model.Classroom{}).Not(tempClassrooms1, tempClassrooms2).Scan(&classrooms)
		} else {
			db.Debug().Model(&Model.Course{}).Select("courses.classroom_id as id,c.name").Joins("left join classrooms c on c.id = courses.classroom_id").Where(" courses.id = ? AND start_time = ? and end_time = ? AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (record_id = ? OR copy_flag = ?)", course.ID, course.StartTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.RecordID, course.RecordID).Group("c.name").Scan(&classrooms)
		}
		c.JSON(http.StatusOK, gin.H{
			"classrooms": classrooms,
		})
	}
}

func GetClassSheet(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var classSheet [7][14]string
		courses := Model.ClassSheets{}
		teacherId := c.Query("teacher_id")
		studentId := c.Query("student_id")
		week := c.Query("week")
		if teacherId == "" {
			db.Debug().Model(&Model.Student2Course{}).Select("courses.*,t.name as teacher_name,c.name as classroom_name").Joins("left join courses on courses.record_id = student2_Courses.record_id").Joins("left join teachers t on t.id = courses.teacher_id").Joins("left join classrooms c on c.id = courses.classroom_id ").Where("student_id = ? AND start_week <= ? AND end_week >= ?", studentId, week, week).Scan(&courses)
		} else {
			db.Debug().Model(&Model.Course{}).Select("courses.*,t.name as teacher_name,c.name as classroom_name").Joins("left join teachers t on t.id = courses.teacher_id").Joins("left join classrooms c on c.id = courses.classroom_id ").Where("teacher_id = ? AND start_week <= ? AND end_week >= ?", teacherId, week, week).Scan(&courses)
		}
		for _, v := range courses {
			for i := v.StartTime; i <= v.EndTime; i++ {
				weekTime := v.WeekTime
				courseName := v.CourseName
				courseClassroom := v.ClassroomName
				teacherName := v.TeacherName
				startWeek := utils.ToString(v.StartWeek)
				endWeek := utils.ToString(v.EndWeek)
				classSheet[weekTime-1][i-1] = courseName + "\n教师：" + teacherName + "\n教室：" + courseClassroom + "\n第" + startWeek + "周-第" + endWeek + "周"
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"classSheet": classSheet,
		})
	}
}

func GetAvailableCourses(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		studentId := c.Query("student_id")
		facultyId := c.Query("faculty_id")
		directionId := c.Query("direction_id")
		specialtyId := c.Query("specialty_id")
		courses := &Model.CourseForChooses{}
		tempCourses := Model.Student2CourseForChooses{}
		db.Debug().Model(&Model.Student2Course{}).Select("course_id").Where("student_id = ?", studentId).Scan(&tempCourses)
		tempIDs := make(Model.SelectedCourseIDs, len(tempCourses))
		for i, v := range tempCourses {
			tempIDs[i].ID = v.CourseID
		}
		fmt.Println(tempCourses)
		if directionId == "0" {
			db.Debug().Table("(?) as F", db.Model(&Model.Course{}).Select("courses.*,count(s2c.record_id) as selected").Joins("left join student2_courses s2c on courses.record_id = s2c.record_id").Where("courses.deleted_at is null AND (courses.faculty_id = ? OR courses.faculty_id = 0) AND (courses.direction_id = 0) AND (courses.specialty_id = ? OR courses.specialty_id = 0) AND copy_flag = 0", facultyId, specialtyId).Group("s2c.record_id")).Joins("left join teachers t on f.teacher_id = t.id").Where("F.selected_num <max_choose_num").Not(tempIDs).Scan(&courses)
			//db.Debug().Raw("select * from (select student2_courses.record_id,count(student2_courses.record_id) from student2_courses group by record_id) F left join courses c on c.record_id = F.record_id where deleted_at is null AND (faculty_id = ? OR faculty_id = 0) AND (direction_id = 0) AND (specialty_id = ? OR specialty_id = 0) AND copy_flag = 0",facultyId,specialtyId).Joins("left join teachers on teachers.teacher_id = courses.teacher_id").Joins("left join teachers on teachers.teacher_id = courses.teacher_id").Not(tempCourses).Scan(&courses)
		} else {
			db.Debug().Table("(?) as F", db.Model(&Model.Course{}).Select("courses.*,count(s2c.record_id) as selected").Joins("left join student2_courses s2c on courses.record_id = s2c.record_id").Where("courses.deleted_at is null AND (courses.faculty_id = ? OR courses.faculty_id = 0) AND (courses.direction_id = 0 OR direction_id = ? ) AND (courses.specialty_id = ? OR courses.specialty_id = 0) AND copy_flag = 0", facultyId, directionId, specialtyId).Group("s2c.record_id")).Joins("left join teachers t on f.teacher_id = t.id").Where("AND F.selected_num <max_choose_num").Not(tempIDs).Scan(&courses)
			//db.Debug().Table("(?) as F",db.Model(&Model.Student2Course{}).Select("student2_courses.record_id,count(student2_courses.record_id) as selected_num").Group("record_id")).Select("F.*,c.*,t.*,F.selected_num selected_num").Joins("left join courses c on c.record_id = F.record_id").Joins("left join teachers t on t.id = c.teacher_id").Where("c.deleted_at is null AND (c.faculty_id = ? OR c.faculty_id = 0) AND (c.direction_id = 0 OR direction_id = ? ) AND (c.specialty_id = ? OR c.specialty_id = 0) AND copy_flag = 0 AND F.selected_num <c.max_choose_num",facultyId,directionId,specialtyId).Not(tempCourses).Scan(&courses)
		}
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
		})
	}
}

func ChooseCourse(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var count int64
		record := &Model.Student2Course{}
		course := &Model.CourseForSelect{}
		record.RecordID, _ = strconv.Atoi(c.PostForm("record_id"))
		record.StudentID, _ = strconv.Atoi(c.PostForm("student_id"))
		record.CourseID, _ = strconv.Atoi(c.PostForm("course_id"))
		db.Debug().Model(&Model.Course{}).Select("max_choose_num").Where("record_id = ?", record.RecordID).Scan(&course)
		db.Debug().Create(&record)
		db.Debug().Model(&Model.Student2Course{}).Select("count(1)").Where("record_id = ?", record.RecordID).Count(&count)
		fmt.Println(count, int64(course.MaxChooseNum))
		if count <= int64(course.MaxChooseNum) {
			c.JSON(http.StatusOK, gin.H{
				"snackbar": "setSuccess",
				"msg":      "添加课程成功",
			})
		} else {
			db.Unscoped().Model(&Model.Student2Course{}).Delete(record)
			c.JSON(http.StatusOK, gin.H{
				"snackbar": "setError",
				"msg":      "添加课程失败!人数已满！",
			})
		}

	}
}

func GetChosenCourse(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		studentId := c.Query("student_id")
		courses := &Model.CourseForChooses{}
		db.Debug().Model(&Model.Student2Course{}).Select("*").Joins("left join courses on courses.record_id = student2_courses.record_id left join teachers on teachers.id = courses.teacher_id").Where("courses.deleted_at IS NULL AND student_id = ?", studentId).Scan(&courses)
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
		})
	}
}

func QuitCourse(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		StudentID := c.Query("student_id")
		CourseID := c.Query("course_id")
		RecordID := c.Query("record_id")
		db.Unscoped().Where("student_id = ? AND course_id = ? AND record_id = ?",
			StudentID, CourseID, RecordID).Delete(&Model.Student2Course{})
		fmt.Println(db.RowsAffected)
		c.JSON(http.StatusOK, gin.H{
			"snackbar": "setSuccess",
			"msg":      "退课成功!",
		})
	}
}

func GetStudentScore (db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var count int64
		scores := &Model.StudentHomeworkScores{}
		recordID := c.Query("record_id")
		db.Model(&Model.Homework{}).Select("count(distinct id)").Where("record_id = ?",recordID).Count(&count)
		fmt.Println(count)
		db.Debug().Table(" (?) as T",db.Model(&Model.Homework{}).Joins("left join homework_upload_records on homework_id = homeworks.id AND homeworks.question_id = homework_upload_records.question_id").Select("homework_upload_records.student_id,homework_upload_records.homework_id,sum(score) total_score ").Group("id")).Select("name,T.student_id,T.homework_id,round(sum(total_score)/?,0) homework_score",count).Joins("left join students on student_id = students.id").Group("student_id").Scan(&scores)
		c.JSON(http.StatusOK,gin.H{
			"scores":scores,
		})
	}
}

func ReturnOK(c *gin.Context) {
	c.JSON(http.StatusOK, 0)
}
