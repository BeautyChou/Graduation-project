package Controller

import (
	"BackEnd/Middleware"
	"BackEnd/Model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Test(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func Auth(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := &Model.User{}
		temp := &Model.User{}
		var err error
		user.ID, err = strconv.Atoi(c.PostForm("id"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "无效的参数",
			})
		}
		user.Password = c.PostForm("password")

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		fmt.Println(string(hash))
		fmt.Println(user)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"msg": "加密失败！",
			})
		}
		db.Model(&Model.User{}).Where("id = ?", user.ID).Find(&temp)
		if err = bcrypt.CompareHashAndPassword([]byte(temp.Password), []byte(user.Password)); err == nil {
			TokenString, _ := Middleware.GenToken(strconv.Itoa(user.ID))
			if user.ID > 0 && user.ID < 100000 {
				admin := &Model.Admin{}
				db.Model(&Model.Admin{}).Where("id = ?", temp.ID).Scan(&admin)
				c.JSON(http.StatusOK, gin.H{
					"ID":       admin.ID,
					"level":    3,
					"username": admin.Name,
					"msg":      "success",
					"data":     gin.H{"token": TokenString},
				})
			} else if user.ID >= 100000 && user.ID < 1000000 {
				teacher := &Model.TeacherForSelect{}
				db.Model(&Model.Teacher{}).Where("id = ?", temp.ID).Scan(&teacher)
				c.JSON(http.StatusOK, gin.H{
					"ID":         teacher.ID,
					"level":      2,
					"username":   teacher.Name,
					"faculty_id": teacher.FacultyID,
					"msg":        "success",
					"data":       gin.H{"token": TokenString},
				})
			} else {
				student := &Model.StudentForSelect{}
				db.Model(&Model.Student{}).Where("id = ?", temp.ID).Scan(&student)
				c.JSON(http.StatusOK, gin.H{
					"ID":           student.ID,
					"level":        1,
					"username":     student.Name,
					"faculty_id":   student.FacultyID,
					"specialty_id": student.SpecialtyID,
					"direction_id": student.DirectionID,
					"msg":          "success",
					"data":         gin.H{"token": TokenString},
				})
			}
			return
		} else {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "failed",
		})
		return
	}
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
		dst := fmt.Sprintf("./images/%s/%s/%s", homeworkID, questionID, file.Filename)
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

func SendCheckedImage(c *gin.Context) {
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
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		var teacherId = c.Query("teacher_id")
		var studentId = c.Query("student_id")
		var courses Model.CourseForSelects
		if studentId == "" {
			db.Debug().Model(&Model.Course{}).Select("courses.id,course_name,credit,teacher_id,classroom_id,max_choose_num,selected_num,start_time,end_time,name,courses.record_id,start_week,end_week").Joins("right join teachers on courses.teacher_id = teachers.id").Where("teacher_id = ? AND copy_flag = 0", teacherId).Group("record_id").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&courses)
		} else {
			db.Debug().Model(&Model.Student2Course{}).Select("courses.id,course_name,credit,teacher_id,classroom_id,max_choose_num,selected_num,start_time,end_time,name,courses.record_id,start_week,end_week").Joins("left join courses on courses.record_id = student2_courses.record_id right join teachers on courses.teacher_id = teachers.id").Where("student_id = ?", studentId).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&courses)
		}
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
			"total":   count,
		})
	}
}

func GetHomeworkList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		var courseId = c.Query("course_id")
		var recordId = c.Query("record_id")
		var homeworks Model.HomeworkForSelects
		db.Model(&[]Model.Homework{}).Where("course_id = ? AND record_id = ?", courseId, recordId).Group("id").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&homeworks)
		//db.Raw("Select * from homeworks where course_id = ? AND deleted_at is NULL GROUP BY id", courseId).Scan(&homeworks)
		c.JSON(http.StatusOK, gin.H{
			"homeworks": homeworks,
			"time":      time.Now(),
			"total":     count,
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
		fmt.Println(homework.ID)
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
		// 创建文件夹并修改权限
		err1 := os.MkdirAll("./images/"+id, 0777)
		err2 := os.Chmod("./images/"+id, 0777)
		if err1 != nil || err2 != nil {
			fmt.Println(err1, err2)
			return
		}
		for _, item := range questions {
			err1 := os.MkdirAll("./images/"+id+"/"+utils.ToString(item.QuestionID), 0666)
			err2 := os.Chmod("./images/"+id+"/"+utils.ToString(item.QuestionID), 0777)
			if err1 != nil || err2 != nil {
				fmt.Println(err1, err2)
				return
			}
		}

		db.Create(&questions)
	}
}

func GetQuestionList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		var questions Model.QuestionForStudents
		var records Model.HomeworkUploadRecordsForSelects
		homeworkID := c.Query("homework_id")
		studentId := c.Query("student_id")
		db.Model(&Model.Homework{}).Select("dead_line,homework_title,question_max_score,content,question_id,record_id").Joins("right join questions on homeworks.question_id = questions.id ").Where("homeworks.id = ?", homeworkID).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&questions)
		//db.Raw("select dead_line,homework_title,question_max_score,content,question_id,record_id from homeworks right join questions on homeworks.question_id = questions.id where homeworks.id = ?", homeworkID).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&questions)
		if studentId != "" {
			db.Model(&Model.HomeworkUploadRecord{}).Where("homework_id = ? AND student_id = ?", homeworkID, studentId).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&records)
			//db.Raw("select * from homework_upload_records where homework_id = ? AND student_id = ?", homeworkID, studentId).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&records)
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
			"total":     count,
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
		record.RecordID, err = strconv.Atoi(recordID)
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
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			fmt.Println(err)
			return
		}
		db.Create(&record)
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
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
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
			db.Find("teacher_id = ?", tId, &Model.Course{}).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&classes)
			c.JSON(http.StatusOK, gin.H{
				"classes": classes,
			})
		} else {
			db.Find(&[]Model.Teacher{}).Where("title_id <> 0 AND teachers.id <> 100000").Scan(&teachers)
			db.Find(&[]Model.Course{}).Order("id asc").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&classes)
			db.Find(&[]Model.Faculty{}).Where("id <> 1").Scan(&faculties)
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
				"total":       count,
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
		course.Selectable = c.PostForm("selectable") == "true"
		flag := c.PostForm("flag")
		if flag == "copy" || flag == "course" {
			db.Select("id").Where("course_name = ?", course.CourseName).Last(&courseTemp)
		} else {
			db.Unscoped().Select("id").Order("id desc").Limit(1).Find(&courseTemp)
			courseTemp.ID++
		}
		db.Unscoped().Select("record_id").Order("record_id desc").Limit(1).Find(&course)
		course.RecordID++
		course.ID = courseTemp.ID
		fmt.Println(course)
		db.Create(&course)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "添加课程成功",
			"snackbar2": "closeSuccess",
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
			"snackbar":  "setSuccess",
			"msg":       "删除课程成功！",
			"id":        id,
			"snackbar2": "closeSuccess",
		})
	}
}

func UploadClass(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		classes := Model.CourseForSelects{}
		err := c.ShouldBindJSON(&classes)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"snackbar":  "setError",
				"msg":       "修改课程失败!",
				"snackbar2": "closeSuccess",
			})
			return
		}
		for _, v := range classes {
			fmt.Println(v)
			db.Where("record_id = ?", v.RecordID).Updates(&Model.Course{
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
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改课程成功",
			"snackbar2": "closeSuccess",
		})
	}
}

func PostApply(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		apply := &Model.ApplyForCourseChange{}
		c.ShouldBindJSON(&apply)
		db.Create(&apply)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "提交请求成功！",
			"snackbar2": "closeSuccess",
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
		db.Model(&Model.ApplyForCourseChange{}).Where("id = ?", id).Update("result", result)
		if result == 1 {
			apply := &Model.ApplyForCourseChange{}
			course := &Model.Course{}
			tempCourse := &Model.Course{}
			db.Where("id = ?", id).Find(&apply)
			db.Where("id = ? AND week_time = ? AND ( IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?) )", apply.CourseID, apply.BeforeWeekTime, apply.BeforeStartTime, apply.BeforeStartTime, apply.BeforeEndTime, apply.BeforeEndTime).Find(&course)

			if course.StartWeek == apply.BeforeStartWeek && course.EndWeek == apply.BeforeEndWeek {
				db.Debug().Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("start_time", apply.AfterStartTime).Update("end_time", apply.AfterEndTime).Update("week_time", apply.AfterWeekTime)
				c.JSON(http.StatusOK, 0)
				return
			} else if course.StartWeek == apply.BeforeStartWeek || course.EndWeek == apply.BeforeEndWeek {
				if course.StartWeek == apply.BeforeStartWeek {
					db.Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("start_week", apply.BeforeEndWeek+1)
				} else {
					db.Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("end_week", apply.BeforeStartWeek-1)
				}
			} else {
				fmt.Println(course)
				courseOBJ1 := course
				courseOBJ1.StartWeek = apply.BeforeEndWeek + 1
				courseOBJ1.EndWeek = course.EndWeek
				courseOBJ1.CopyFlag = course.RecordID
				db.Unscoped().Select("record_id").Order("record_id desc").Limit(1).Find(&tempCourse)
				db.Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Update("end_week", apply.BeforeStartWeek-1)
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
			db.Unscoped().Select("record_id").Order("record_id desc").Limit(1).Find(&course)
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
		tempClassrooms3 := &Model.ClassroomForSelects{}
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
		course.MaxChooseNum, _ = strconv.Atoi(c.PostForm("max_choose_num"))
		fmt.Println("course_num:", course.MaxChooseNum)
		if course.MaxChooseNum == 0 {
			// 申请换课
			class := &Model.Course{}
			db.Model(&Model.Course{}).Where("record_id = ?", course.RecordID).Scan(&class)
			db.Model(&Model.Classroom{}).Where("max_num < ?", class.MaxChooseNum).Scan(&tempClassrooms3)
		} else {
			// 设置课程
			db.Model(&Model.Classroom{}).Where("max_num < ?", course.MaxChooseNum).Scan(&tempClassrooms3)
		}
		timeNode := c.PostForm("time")
		flag := c.PostForm("flag")
		if timeNode == "after" || timeNode == "applyAfter" {
			if timeNode == "applyAfter" {
				db.Debug().Model(&Model.Course{}).Select("faculty_id,direction_id,specialty_id,teacher_id,record_id").Where(" courses.id = ?", course.ID).Scan(&tempCourse)
				course.FacultyID = tempCourse.FacultyID
				course.DirectionID = tempCourse.DirectionID
				course.SpecialtyID = tempCourse.SpecialtyID
				course.TeacherID = tempCourse.TeacherID
				course.RecordID = tempCourse.RecordID
				if flag != "not" {
					//判断教师、上课同学时间是否有冲突
					if course.DirectionID == 0 && course.SpecialtyID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = 0))) AND record_id <> ?) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.RecordID).Scan(&tempClassrooms1)
					} else if course.DirectionID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = 0))) AND record_id <> ? ) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.RecordID).Scan(&tempClassrooms1)
					} else if course.SpecialtyID == 0 {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = ? OR direction_id = 0))) AND record_id <> ? ) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.DirectionID, course.RecordID).Scan(&tempClassrooms1)
					} else {
						db.Debug().Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = ? OR direction_id = 0))) AND record_id <> ? ) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.DirectionID, course.RecordID).Scan(&tempClassrooms1)
					}
				} else {
					//不可能有空教室的情况
					if course.DirectionID == 0 && course.SpecialtyID == 0 {
						db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID).Scan(&tempClassrooms1)
					} else if course.DirectionID == 0 {
						db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID).Scan(&tempClassrooms1)
					} else if course.SpecialtyID == 0 {
						db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.DirectionID).Scan(&tempClassrooms1)
					} else {
						db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.DirectionID).Scan(&tempClassrooms1)
					}
				}
			} else if course.DirectionID == 0 && course.SpecialtyID == 0 {
				db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID).Scan(&tempClassrooms1)
			} else if course.DirectionID == 0 {
				db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID).Scan(&tempClassrooms1)
			} else if course.SpecialtyID == 0 {
				db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.DirectionID).Scan(&tempClassrooms1)
			} else {
				db.Model(&Model.Classroom{}).Where("(select count(1) from courses where (IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ? AND courses.deleted_at IS NULL AND (teacher_id = ? OR(faculty_id = ? AND (specialty_id = ? OR specialty_id = 0) AND (direction_id = ? OR direction_id = 0)))) > 0", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime, course.TeacherID, course.FacultyID, course.SpecialtyID, course.DirectionID).Scan(&tempClassrooms1)
			}
			db.Model(&Model.Course{}).Select("courses.classroom_id as id,c.name").Joins("left join classrooms c on c.id = courses.classroom_id").Where("(IF(start_time > ? ,start_time,?)<=IF(end_time<?,end_time,?)) AND (IF(start_week > ? ,start_week,?)<=IF(end_week<?,end_week,?)) AND week_time = ?", course.StartTime, course.StartTime, course.EndTime, course.EndTime, course.StartWeek, course.StartWeek, course.EndWeek, course.EndWeek, course.WeekTime).Group("c.name").Scan(&tempClassrooms2)
			db.Model(&Model.Classroom{}).Not(tempClassrooms1, tempClassrooms2, tempClassrooms3).Scan(&classrooms)
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
			db.Debug().Model(&Model.Student2Course{}).Select("courses.*,t.name as teacher_name,c.name as classroom_name").Joins("left join courses on (courses.record_id = student2_courses.record_id OR courses.copy_flag = student2_courses.record_id)").Joins("left join teachers t on t.id = courses.teacher_id").Joins("left join classrooms c on c.id = courses.classroom_id ").Where("student_id = ? AND start_week <= ? AND end_week >= ?", studentId, week, week).Scan(&courses)
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
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		studentId := c.Query("student_id")
		facultyId := c.Query("faculty_id")
		directionId := c.Query("direction_id")
		specialtyId := c.Query("specialty_id")
		courses := &Model.CourseForChooses{}
		tempCourses := Model.Student2CourseForChooses{}
		db.Model(&Model.Student2Course{}).Select("course_id").Where("student_id = ?", studentId).Scan(&tempCourses)
		tempIDs := make(Model.SelectedCourseIDs, len(tempCourses))
		for i, v := range tempCourses {
			tempIDs[i].ID = v.CourseID
		}
		fmt.Println(tempCourses)
		if directionId == "0" {
			db.Debug().Table("(?) as F", db.Model(&Model.Course{}).Select("courses.*,count(s2c.record_id) as selected,courses.id as course_id").Joins("left join student2_courses s2c on courses.record_id = s2c.record_id").Where("courses.deleted_at is null AND (courses.faculty_id = ? OR courses.faculty_id = 0) AND (courses.direction_id = 0) AND (courses.specialty_id = ? OR courses.specialty_id = 0) AND copy_flag = 0 AND selectable = 1", facultyId, specialtyId).Group("courses.record_id")).Joins("left join teachers t on F.teacher_id = t.id").Where("F.selected_num <max_choose_num").Not(tempIDs).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&courses)
			//db.Raw("select * from (select student2_courses.record_id,count(student2_courses.record_id) from student2_courses group by record_id) F left join courses c on c.record_id = F.record_id where deleted_at is null AND (faculty_id = ? OR faculty_id = 0) AND (direction_id = 0) AND (specialty_id = ? OR specialty_id = 0) AND copy_flag = 0",facultyId,specialtyId).Joins("left join teachers on teachers.teacher_id = courses.teacher_id").Joins("left join teachers on teachers.teacher_id = courses.teacher_id").Not(tempCourses).Scan(&courses)
		} else {
			db.Debug().Table("(?) as F", db.Model(&Model.Course{}).Select("courses.*,count(s2c.record_id) as selected,courses.id as course_id").Joins("left join student2_courses s2c on courses.record_id = s2c.record_id").Where("courses.deleted_at is null AND (courses.faculty_id = ? OR courses.faculty_id = 0) AND (courses.direction_id = 0 OR direction_id = ? ) AND (courses.specialty_id = ? OR courses.specialty_id = 0) AND copy_flag = 0 AND selectable = 1", facultyId, directionId, specialtyId).Group("courses.record_id")).Joins("left join teachers t on F.teacher_id = t.id").Where("AND F.selected_num <max_choose_num").Not(tempIDs).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&courses)
			//db.Table("(?) as F",db.Model(&Model.Student2Course{}).Select("student2_courses.record_id,count(student2_courses.record_id) as selected_num").Group("record_id")).Select("F.*,c.*,t.*,F.selected_num selected_num").Joins("left join courses c on c.record_id = F.record_id").Joins("left join teachers t on t.id = c.teacher_id").Where("c.deleted_at is null AND (c.faculty_id = ? OR c.faculty_id = 0) AND (c.direction_id = 0 OR direction_id = ? ) AND (c.specialty_id = ? OR c.specialty_id = 0) AND copy_flag = 0 AND F.selected_num <c.max_choose_num",facultyId,directionId,specialtyId).Not(tempCourses).Scan(&courses)
		}
		c.JSON(http.StatusOK, gin.H{
			"total":   count,
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
		db.Model(&Model.Course{}).Select("max_choose_num").Where("record_id = ?", record.RecordID).Scan(&course)
		db.Create(&record)
		db.Model(&Model.Student2Course{}).Select("count(1)").Where("record_id = ?", record.RecordID).Count(&count)
		fmt.Println(count, int64(course.MaxChooseNum))
		if count <= int64(course.MaxChooseNum) {
			c.JSON(http.StatusOK, gin.H{
				"snackbar":  "setSuccess",
				"msg":       "添加课程成功",
				"snackbar2": "closeSuccess",
			})
		} else {
			db.Unscoped().Model(&Model.Student2Course{}).Delete(record)
			c.JSON(http.StatusOK, gin.H{
				"snackbar":  "setError",
				"msg":       "添加课程失败!人数已满！",
				"snackbar2": "closeError",
			})
		}
	}
}

func GetChosenCourse(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		now := time.Now()
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		studentId := c.Query("student_id")
		courses := &Model.CourseForChooses{}
		db.Model(&Model.Student2Course{}).Select("courses.*,student2_courses.*,teachers.*").Joins("left join courses on courses.record_id = student2_courses.record_id left join teachers on teachers.id = courses.teacher_id").Where("courses.deleted_at IS NULL AND student_id = ?", studentId).Count(&count).Offset((page - 1) * items).Limit(items).Scan(&courses)
		c.JSON(http.StatusOK, gin.H{
			"total":   count,
			"courses": courses,
			"time":    now,
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
			"snackbar":  "setSuccess",
			"msg":       "退课成功!",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetStudentScore(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var count int64
		scores := &Model.StudentHomeworkScores{}
		recordID := c.Query("record_id")
		db.Model(&Model.Homework{}).Select("count(distinct id)").Where("record_id = ?", recordID).Count(&count)
		fmt.Println(count)
		db.Table(" (?) as T", db.Model(&Model.Homework{}).Joins("left join homework_upload_records on homework_id = homeworks.id AND homeworks.question_id = homework_upload_records.question_id").Select("homework_upload_records.student_id,homework_upload_records.homework_id,sum(score) total_score,homeworks.record_id").Group("id")).Select("name,T.student_id,T.homework_id,round(sum(total_score)/?,0) homework_score", count).Joins("left join students on student_id = students.id").Where("record_id = ?", recordID).Group("student_id").Scan(&scores)
		c.JSON(http.StatusOK, gin.H{
			"scores": scores,
		})
	}
}

func PostScore(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var scores Model.Electives
		temp := &Model.Elective{}
		if err := c.ShouldBindJSON(&scores); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"snackbar":  "setError",
				"msg":       "上传分数失败，请稍后重试！",
				"snackbar2": "closeError",
			})
			return
		}
		if err := db.Where("course_id = ? AND student_id = ?", scores[0].CourseID, scores[0].StudentID).First(&temp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&scores)
			}
		} else {
			for _, v := range scores {
				db.Model(&Model.Elective{}).Where("course_id = ? AND student_id = ?", v.CourseID, v.StudentID).Updates(map[string]interface{}{
					"HomeworkScore":      v.HomeworkScore,
					"TestScore":          v.TestScore,
					"BehaviorScore":      v.BehaviorScore,
					"Percentage":         v.Percentage,
					"HomeworkPercentage": v.HomeworkPercentage,
				})
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "上传分数成功!",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetScore(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		studentId := c.Query("student_id")
		electives := &Model.ElectiveForSelects{}
		db.Debug().Model(&Model.Elective{}).Select("student_id,course_id,course_name,homework_score,test_score,behavior_score,percentage,homework_percentage").Joins("left join courses on courses.id = electives.course_id").Where("student_id = ?", studentId).Group("course_id").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&electives)
		c.JSON(http.StatusOK, gin.H{
			"total":     count,
			"electives": electives,
		})
	}
}

func GetUserInfo(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		studentID := c.Query("student_id")
		teacherID := c.Query("teacher_id")
		student := &Model.StudentForUserInfo{}
		teacher := &Model.TeacherForUserInfo{}
		students := &Model.StudentForSelects{}
		currentTime := time.Now()
		if studentID == "" {
			bTime := time.Date(currentTime.Year()-4, 6, 1, 0, 0, 0, 00, time.Local)
			tTime := time.Date(currentTime.Year()-4, 8, 31, 23, 59, 59, 99, time.Local)
			db.Model(&Model.Teacher{}).Select("teachers.*,faculties.name faculty_name").Where("teachers.id = ?", teacherID).Joins("left join faculties on faculties.id = teachers.faculty_id").Find(&teacher)
			db.Model(&Model.Student{}).Where("teacher_id = ? AND created_at >= ? AND created_at <= ?", teacherID, bTime, tTime).Scan(&students)
		} else {
			db.Model(&Model.Student{}).Select("students.created_at created,faculties.name faculty_name,direction_to_specialties.*,students.*").Where("students.id = ?", studentID).Joins("left join faculties on faculties.id = students.faculty_id left join direction_to_specialties on direction_to_specialties.faculty_id = students.faculty_id AND direction_to_specialties.specialty_id = students.specialty_id AND direction_to_specialties.direction_id = students.direction_id").Find(&student)
		}
		c.JSON(http.StatusOK, gin.H{
			"student":      student,
			"teacher":      teacher,
			"current_time": currentTime,
			"students":     students,
		})
	}
}

func GetDirectionList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		facultyID := c.Query("faculty_id")
		specialtyID := c.Query("specialty_id")
		directions := &Model.DirectionToSpecialtyForSelects{}
		db.Model(&Model.DirectionToSpecialty{}).Where("faculty_id = ? AND specialty_id = ? AND direction_id <> 0", facultyID, specialtyID).Scan(&directions)
		c.JSON(http.StatusOK, gin.H{
			"directions": directions,
		})
	}
}

func GetTeacherList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		facultyID := c.Query("faculty_id")
		teachers := &Model.TeacherForSelects{}
		db.Model(&Model.Teacher{}).Where("faculty_id = ? AND teachers.id <> 100000", facultyID).Scan(&teachers)
		c.JSON(http.StatusOK, gin.H{
			"teachers": teachers,
		})
	}
}

func PostDirection(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		directionID, _ := strconv.Atoi(c.PostForm("direction_id"))
		studentID := c.PostForm("student_id")
		db.Model(&Model.Student{}).Where("id = ?", studentID).Updates(&Model.Student{
			DirectionID: directionID,
		})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改方向成功!",
			"snackbar2": "closeSuccess",
		})
	}
}

func PostTeacher(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		teacherID, _ := strconv.Atoi(c.PostForm("teacher_id"))
		studentID := c.PostForm("student_id")
		db.Model(&Model.Student{}).Where("id = ?", studentID).Updates(&Model.Student{
			TeacherID: teacherID,
		})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "导师修改成功!",
			"snackbar2": "closeSuccess",
		})
	}
}

func PutPractice(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		practice, _ := strconv.Atoi(c.PostForm("practice"))
		studentID := c.PostForm("student_id")
		db.Model(&Model.Student{}).Where("id = ?", studentID).Updates(&Model.Student{
			Practice: practice,
		})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "实习方式修改成功!",
			"snackbar2": "closeSuccess",
		})
	}
}

func PostPractice(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		practice := &Model.IndependentPractice{}
		temp := &Model.IndependentPractice{}
		if err := c.ShouldBindJSON(&practice); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"snackbar":  "setError",
				"msg":       "申请表格式有错，请重新填写！",
				"snackbar2": "closeError",
			})
			return
		}
		fmt.Println(practice.Phone)
		if err := db.Where("student_id = ?", practice.StudentID).First(&temp).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Model(&Model.IndependentPractice{}).Create(&practice)
			}
		} else {
			db.Model(&Model.IndependentPractice{}).Where("student_id = ?", practice.StudentID).Updates(&practice)
		}
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "申请表上传成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func ApplyTeacher(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		studentID := c.PostForm("student_id")
		teacherID := c.PostForm("teacher_id")
		state := c.PostForm("operation")
		if state == "pass" {
			course := &Model.CourseForSelect{}
			record := &Model.Student2Course{}
			var count int64
			db.Model(&Model.Course{}).Where("teacher_id = ? AND selectable = 0", teacherID).Order("record_id desc").Limit(1).Scan(&course)
			record.StudentID, _ = strconv.Atoi(studentID)
			record.CourseID = course.ID
			record.RecordID = course.RecordID
			db.Model(&Model.Course{}).Select("max_choose_num").Where("record_id = ?", course.RecordID).Scan(&course)
			db.Create(&record)
			db.Model(&Model.Student2Course{}).Select("count(1)").Where("record_id = ?", record.RecordID).Count(&count)
			fmt.Println(count, int64(course.MaxChooseNum))
			if count <= int64(course.MaxChooseNum) {
				c.JSON(http.StatusOK, gin.H{
					"snackbar":  "setSuccess",
					"msg":       "添加课程成功",
					"snackbar2": "closeSuccess",
				})
				db.Model(&Model.Student{}).Where("id = ?", studentID).Update("teacher_flag", 1)
			} else {
				db.Unscoped().Model(&Model.Student2Course{}).Delete(record)
				c.JSON(http.StatusOK, gin.H{
					"snackbar":  "setError",
					"msg":       "添加课程失败!人数已满！",
					"snackbar2": "closeError",
				})
				db.Model(&Model.Student{}).Where("id = ?", studentID).Update("teacher_id", 0)
			}
		} else {
			db.Model(&Model.Student{}).Where("id = ?", studentID).Update("teacher_id", 0)
		}
	}
}

func GetStudentList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		specialtyId := c.Query("specialty_id")
		facultyId := c.Query("faculty_id")
		students := &Model.StudentForUserInfos{}
		faculties := Model.Faculties{}
		specialties := Model.DirectionToSpecialtyForSelects{}
		directions := Model.DirectionToSpecialtyForSelects{}
		facultyToSpecialty := make(map[int]Model.DirectionToSpecialtyForSelects)
		SpecialtyToDirection := make(map[int]Model.DirectionToSpecialtyForSelects)
		PunishmentLevel := &Model.PunishmentLevels{}
		db.Find(&[]Model.DirectionToSpecialty{}).Where("direction_id <> 0").Scan(&directions)
		db.Model(&Model.PunishmentLevel{}).Scan(&PunishmentLevel)
		db.Find(&[]Model.DirectionToSpecialty{}).Where("specialty_id <> 0").Scan(&specialties)
		db.Find(&[]Model.Faculty{}).Where("id <> 1").Scan(&faculties)
		if facultyId != "" {
			if specialtyId != "" {
				db.Model(&Model.Student{}).Select("students.ID student_id,students.created_at created,f.name faculty_name,dts.*,students.*").Where(" students.faculty_id = ? AND students.specialty_id = ?", facultyId, specialtyId).Joins("left join direction_to_specialties dts on( dts.direction_id = students.direction_id and students.specialty_id = dts.specialty_id and students.direction_id = dts.direction_id and dts.faculty_id = students.faculty_id) left join faculties f on f.id = students.faculty_id ").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&students)
			} else {
				db.Model(&Model.Student{}).Select("students.ID student_id,students.created_at created,f.name faculty_name,dts.*,students.*").Where(" students.faculty_id = ?", facultyId).Joins("left join direction_to_specialties dts on( dts.direction_id = students.direction_id and students.specialty_id = dts.specialty_id and students.direction_id = dts.direction_id and dts.faculty_id = students.faculty_id) left join faculties f on f.id = students.faculty_id ").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&students)
			}
		} else {
			db.Model(&Model.Student{}).Select("students.ID student_id,students.created_at created,f.name faculty_name,dts.*,students.*").Joins("left join direction_to_specialties dts on( dts.direction_id = students.direction_id and students.specialty_id = dts.specialty_id and students.direction_id = dts.direction_id and dts.faculty_id = students.faculty_id) left join faculties f on f.id = students.faculty_id ").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&students)
		}
		for _, v := range specialties {
			facultyToSpecialty[v.FacultyID] = append(facultyToSpecialty[v.FacultyID], v)
		}
		for _, v := range directions {
			SpecialtyToDirection[v.SpecialtyID] = append(SpecialtyToDirection[v.SpecialtyID], v)
		}
		c.JSON(http.StatusOK, gin.H{
			"total":            count,
			"specialties":      facultyToSpecialty,
			"students":         students,
			"faculties":        faculties,
			"punishment_level": PunishmentLevel,
			"directions":       SpecialtyToDirection,
		})
	}
}

func DeleteStudent(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		StudentID := c.Query("student_id")
		db.Unscoped().Where("id = ?", StudentID).Delete(&Model.Student{})
		db.Where("id = ?", StudentID).Delete(&Model.User{})

		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "删除学生成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetPunishment(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		StudentID := c.Query("student_id")
		Punishments := &Model.PunishmentForSelects{}
		db.Model(&Model.Punishment{}).Select("punishments.*,pl.level Level").Joins("left join punishment_levels pl on pl.id = punishments.punishment_level_id").Where("student_id = ? ", StudentID).Scan(&Punishments)
		c.JSON(http.StatusOK, gin.H{
			"punishments": Punishments,
		})
	}
}

func PutStudent(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		studentId := c.PostForm("student_id")
		password := c.PostForm("password")
		facultyID := c.PostForm("faculty_id")
		specialtyID := c.PostForm("specialty_id")
		directionID := c.PostForm("direction_id")
		if password != "" {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			db.Model(&Model.User{}).Where("id = ?", studentId).Update("password", string(hash))
		}
		if facultyID != "null" {
			if specialtyID != "null" {
				if directionID != "null" {
					db.Model(&Model.Student{}).Where("id = ?", studentId).Updates(map[string]interface{}{"faculty_id": facultyID, "specialty_id": specialtyID, "direction_id": directionID})
				} else {
					db.Model(&Model.Student{}).Where("id = ?", studentId).Updates(map[string]interface{}{"faculty_id": facultyID, "specialty_id": specialtyID})
				}
			} else {
				db.Model(&Model.Student{}).Where("id = ?", studentId).Updates(map[string]interface{}{"faculty_id": facultyID})
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改学生信息成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func PunishStudent(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		Punishment := &Model.Punishment{}
		Punishment.StudentID, _ = strconv.Atoi(c.PostForm("student_id"))
		Punishment.PunishmentLevelID, _ = strconv.Atoi(c.PostForm("punishment_level"))
		Punishment.Reason = c.PostForm("punishment_content")
		db.Create(&Punishment)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "添加处分纪录成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func DeletePunishment(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		PunishmentID := c.Query("punishment_id")
		db.Model(&Model.Punishment{}).Where("id = ?", PunishmentID).Update("is_cancelled", true)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "取消处分成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func AddStudent(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		currentTime := time.Now()
		student := &Model.Student{}
		user := &Model.User{}
		Grade, _ := strconv.Atoi(c.PostForm("grade"))
		date := time.Date(currentTime.Year()-5+Grade, 6, 11, 0, 0, 0, 00, time.Local)
		student.Name = c.PostForm("name")
		Password, _ := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
		student.CreatedAt = date
		student.FacultyID, _ = strconv.Atoi(c.PostForm("faculty_id"))
		student.SpecialtyID, _ = strconv.Atoi(c.PostForm("specialty_id"))
		student.DirectionID, _ = strconv.Atoi(c.PostForm("direction_id"))
		student.TeacherID = 100000
		db.Create(&student)
		user.ID = student.ID
		user.Password = string(Password)
		db.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "创建学生成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetPracticeInfo(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		StudentID := c.Query("student_id")
		independentPractice := &Model.IndependentPracticeForSelect{}
		db.Model(&Model.IndependentPractice{}).Where("student_id = ?", StudentID).Scan(&independentPractice)
		c.JSON(http.StatusOK, gin.H{
			"independent_practice": independentPractice,
		})
	}
}

func GetTeacherListForManage(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		facultyId := c.Query("faculty_id")
		titleID := c.Query("title_id")
		teachers := &Model.TeacherForUserInfos{}
		titles := &Model.Titles{}
		faculties := Model.Faculties{}
		PunishmentLevel := &Model.PunishmentLevels{}
		db.Model(&Model.PunishmentLevel{}).Scan(&PunishmentLevel)
		db.Find(&[]Model.Faculty{}).Where("id <> 1").Scan(&faculties)
		db.Model(&Model.Title{}).Where("id <> 0").Scan(&titles)
		if facultyId != "" {
			if titleID != "" {
				db.Model(&Model.Teacher{}).Select("teachers.*,f.*,t.*,teachers.ID ID,teachers.name name,f.name faculty_name,t.name title_name").Where(" teachers.faculty_id = ? AND teachers.title_id = ? AND teachers.id <> 100000", facultyId, titleID).Joins("left join faculties f on f.id = teachers.faculty_id left join titles t on t.id = teachers.title_id ").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&teachers)
			} else {
				db.Model(&Model.Teacher{}).Select("teachers.*,f.*,t.*,teachers.ID ID,teachers.name name,f.name faculty_name,t.name title_name").Where(" teachers.faculty_id = ? AND teachers.id <> 100000", facultyId).Joins("left join faculties f on f.id = teachers.faculty_id left join titles t on t.id = teachers.title_id ").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&teachers)
			}
		} else {
			if titleID != "" {
				db.Model(&Model.Teacher{}).Select("teachers.*,f.*,t.*,teachers.ID ID,teachers.name name,f.name faculty_name,t.name title_name").Where("teachers.title_id = ? AND teachers.id <> 100000", titleID).Joins("left join faculties f on f.id = teachers.faculty_id left join titles t on t.id = teachers.title_id ").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&teachers)
			} else {
				db.Model(&Model.Teacher{}).Select("teachers.*,f.*,t.*,teachers.ID ID,teachers.name name,f.name faculty_name,t.name title_name").Where("teachers.id <> 100000").Joins("left join faculties f on f.id = teachers.faculty_id left join titles t on t.id = teachers.title_id ").Count(&count).Offset((page - 1) * items).Limit(items).Scan(&teachers)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"total":     count,
			"teachers":  teachers,
			"faculties": faculties,
			"titles":    titles,
		})
	}
}

func DeleteTeacher(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		teacherID := c.Query("teacher_id")
		db.Unscoped().Where("id = ?", teacherID).Delete(&Model.Teacher{})
		db.Where("id = ?", teacherID).Delete(&Model.User{})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "删除教师成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func PutTeacher(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		teacherID := c.PostForm("teacher_id")
		password := c.PostForm("password")
		facultyID := c.PostForm("faculty_id")
		titleID := c.PostForm("title_id")
		if password != "" {
			hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			db.Model(&Model.User{}).Where("id = ?", teacherID).Update("password", string(hash))
		}
		if facultyID != "null" {
			if titleID != "null" {
				db.Model(&Model.Teacher{}).Where("id = ?", teacherID).Updates(map[string]interface{}{"faculty_id": facultyID, "title_id": titleID})
			} else {
				db.Model(&Model.Teacher{}).Where("id = ?", teacherID).Updates(map[string]interface{}{"faculty_id": facultyID})
			}
		} else {
			if titleID != "null" {
				db.Model(&Model.Teacher{}).Where("id = ?", teacherID).Updates(map[string]interface{}{"title_id": titleID})
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改教师信息成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func AddTeacher(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		teacher := &Model.Teacher{}
		user := &Model.User{}
		teacher.Name = c.PostForm("name")
		Password, _ := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
		teacher.FacultyID, _ = strconv.Atoi(c.PostForm("faculty_id"))
		teacher.TitleID, _ = strconv.Atoi(c.PostForm("title_id"))
		db.Create(&teacher)
		user.ID = teacher.ID
		user.Password = string(Password)
		db.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "创建教师成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetAdminList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		admins := &Model.Admins{}
		db.Model(&Model.Admin{}).Count(&count).Offset((page - 1) * items).Limit(items).Find(&admins)
		c.JSON(http.StatusOK, gin.H{
			"admins": admins,
			"total":  count,
		})
	}
}

func DeleteAdmin(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println(c.Copy())
		adminID := c.Query("admin_id")
		fmt.Println(adminID)
		db.Where("id = ?", adminID).Delete(&Model.Admin{})
		db.Where("id = ?", adminID).Delete(&Model.User{})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "删除管理员成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func PutAdmin(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		adminID := c.PostForm("admin_id")
		password := c.PostForm("password")
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		db.Model(&Model.User{}).Where("id = ?", adminID).Update("password", string(hash))
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改管理员信息成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func AddAdmin(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println(c)
		admin := &Model.Admin{}
		user := &Model.User{}
		admin.Name = c.PostForm("name")
		Password, _ := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
		db.Create(&admin)
		user.ID = int(admin.ID)
		user.Password = string(Password)
		db.Create(&user)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "创建管理员成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetTitleList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		titles := &Model.Titles{}
		db.Model(&Model.Title{}).Where("id <> 0").Find(&titles)
		c.JSON(http.StatusOK, gin.H{
			"titles": titles,
		})
	}
}

func DeleteTitle(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		titleID := c.Query("title_id")
		db.Unscoped().Where("id = ?", titleID).Delete(&Model.Title{})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "删除职称成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func PutTitle(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		titleID := c.PostForm("title_id")
		titleName := c.PostForm("name")
		db.Model(&Model.Title{}).Where("id = ?", titleID).Update("name", titleName)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改职称成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func AddTitle(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		title := &Model.Title{}
		title.Name = c.PostForm("name")
		db.Create(&title)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "创建职称成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetPunishLevelList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		punishLevels := &Model.PunishmentLevels{}
		db.Model(&Model.PunishmentLevel{}).Find(&punishLevels)
		c.JSON(http.StatusOK, gin.H{
			"punish_levels": punishLevels,
		})
	}
}

func DeletePunishLevel(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		levelID := c.Query("level_id")
		db.Where("id = ?", levelID).Delete(&Model.PunishmentLevel{})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "删除处分等级成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func PutPunishLevel(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		levelID := c.PostForm("level_id")
		levelName := c.PostForm("name")
		db.Model(&Model.PunishmentLevel{}).Where("id = ?", levelID).Update("level", levelName)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改处分等级成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func AddPunishLevel(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		punishLevel := &Model.PunishmentLevel{}
		punishLevel.Level = c.PostForm("name")
		db.Create(&punishLevel)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "创建处分等级成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetClassroomList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		classrooms := &Model.Classrooms{}
		db.Model(&Model.Classroom{}).Count(&count).Offset((page - 1) * items).Limit(items).Find(&classrooms)
		c.JSON(http.StatusOK, gin.H{
			"classrooms": classrooms,
			"total":      count,
		})
	}
}

func DeleteClassroom(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		classroomID := c.Query("classroom_id")
		db.Unscoped().Where("id = ?", classroomID).Delete(&Model.Classroom{})
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "删除教室成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func PutClassroom(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		classroomID := c.PostForm("classroom_id")
		classroomName := c.PostForm("name")
		classroomMaxNum, _ := strconv.Atoi(c.PostForm("max_num"))
		if classroomName == "null" {
			db.Model(&Model.Classroom{}).Where("id = ?", classroomID).Updates(map[string]interface{}{"max_num": classroomMaxNum})
		} else {
			db.Model(&Model.Classroom{}).Where("id = ?", classroomID).Updates(map[string]interface{}{"name": classroomName, "max_num": classroomMaxNum})
		}
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改教室成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func AddClassroom(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		classroom := &Model.Classroom{}
		classroom.Name = c.PostForm("name")
		classroom.MaxNum, _ = strconv.Atoi(c.PostForm("max_num"))
		db.Create(&classroom)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "创建教室成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetDirectionSpecialtyFacultyList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.Query("page"))
		items, _ := strconv.Atoi(c.Query("items"))
		var count int64
		specialtyId := c.Query("specialty_id")
		facultyId := c.Query("faculty_id")
		facultySpecialtyDirection := &Model.DirectionToSpecialtyForSelects{}
		faculties := Model.Faculties{}
		specialties := Model.DirectionToSpecialtyForSelects{}
		directions := Model.DirectionToSpecialtyForSelects{}
		facultyToSpecialty := make(map[int]Model.DirectionToSpecialtyForSelects)
		SpecialtyToDirection := make(map[int]Model.DirectionToSpecialtyForSelects)
		db.Find(&[]Model.DirectionToSpecialty{}).Where("direction_id <> 0").Scan(&directions)
		db.Find(&[]Model.DirectionToSpecialty{}).Where("specialty_id <> 0").Scan(&specialties)
		db.Find(&[]Model.Faculty{}).Where("id <> 1").Scan(&faculties)
		db.Model(&Model.DirectionToSpecialty{}).Count(&count)
		if facultyId != "" {
			if specialtyId != "" {
				db.Find(&[]Model.DirectionToSpecialty{}).Select("f.name faculty_name,direction_to_specialties.*").Where("faculty_id = ? AND specialty_id = ?", facultyId, specialtyId).Joins("left join faculties f on f.id = direction_to_specialties.faculty_id").Offset((page - 1) * items).Limit(items).Scan(&facultySpecialtyDirection)
			} else {
				db.Find(&[]Model.DirectionToSpecialty{}).Select("f.name faculty_name,direction_to_specialties.*").Where("faculty_id = ?", facultyId).Joins("left join faculties f on f.id = direction_to_specialties.faculty_id").Offset((page - 1) * items).Limit(items).Scan(&facultySpecialtyDirection)
			}
		} else {
			db.Find(&[]Model.DirectionToSpecialty{}).Select("f.name faculty_name,direction_to_specialties.*").Joins("left join faculties f on f.id = direction_to_specialties.faculty_id").Offset((page - 1) * items).Limit(items).Scan(&facultySpecialtyDirection)
		}
		for _, v := range specialties {
			facultyToSpecialty[v.FacultyID] = append(facultyToSpecialty[v.FacultyID], v)
		}
		for _, v := range directions {
			SpecialtyToDirection[v.SpecialtyID] = append(SpecialtyToDirection[v.SpecialtyID], v)
		}
		c.JSON(http.StatusOK, gin.H{
			"total":                       count,
			"specialties":                 facultyToSpecialty,
			"faculties":                   faculties,
			"directions":                  SpecialtyToDirection,
			"faculty_specialty_direction": facultySpecialtyDirection,
		})
	}
}

func AddFacultySpecialtyDirection(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		facultySpecialtyDirection := &Model.DirectionToSpecialty{}
		facultySpecialtyDirectionTemp := &Model.DirectionToSpecialtyForSelect{}
		_ = c.ShouldBind(&facultySpecialtyDirection)
		if facultySpecialtyDirection.FacultyID != 0 {
			if facultySpecialtyDirection.SpecialtyID != 0 {
				//新建方向
				db.Model(&Model.DirectionToSpecialty{}).Where("faculty_id = ? AND specialty_id = ?", facultySpecialtyDirection.FacultyID, facultySpecialtyDirection.SpecialtyID).Order("direction_id desc").Limit(1).Scan(&facultySpecialtyDirectionTemp)
				facultySpecialtyDirection.DirectionID = facultySpecialtyDirectionTemp.DirectionID + 1
			} else {
				db.Model(&Model.DirectionToSpecialty{}).Where("faculty_id = ?", facultySpecialtyDirection.FacultyID).Order("specialty_id desc").Limit(1).Scan(&facultySpecialtyDirectionTemp)
				facultySpecialtyDirection.SpecialtyID = facultySpecialtyDirectionTemp.SpecialtyID + 1
				facultySpecialtyDirection.DirectionName = "全体专业学生"
				//新建专业
			}
		} else {
			facultySpecialtyDirection.SpecialtyName = "全体学院学生"
			facultySpecialtyDirection.DirectionName = "全体专业学生"
		}
		db.Model(&Model.DirectionToSpecialty{}).Create(&facultySpecialtyDirection)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "添加成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func DeleteFacultySpecialtyDirection(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		facultyID := c.Query("faculty_id")
		specialtyID := c.Query("specialty_id")
		directionID := c.Query("direction_id")
		fmt.Println(facultyID, specialtyID, directionID)
		if specialtyID == "" {
			//删除学院
			db.Model(&Model.Faculty{}).Delete("id = ?", facultyID)
			db.Where("faculty_id = ?", facultyID).Delete(&Model.DirectionToSpecialty{})
		} else if directionID == "" {
			//删除专业
			db.Where("faculty_id = ? AND specialty_id = ?", facultyID, specialtyID).Delete(&Model.DirectionToSpecialty{})
		} else {
			//删除方向
			db.Where("faculty_id = ? AND specialty_id = ? AND direction_id = ?", facultyID, specialtyID, directionID).Delete(&Model.DirectionToSpecialty{})
		}
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "删除成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func GetNotification(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		notification := &Model.Notification{}
		db.Model(&Model.Notification{}).First(&notification)
		c.JSON(http.StatusOK, gin.H{
			"notification": notification,
		})
	}
}

func PutNotification(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		content := c.PostForm("notification")
		db.Model(&Model.Notification{}).Where("id = 1").Update("notification", content)
		c.JSON(http.StatusOK, gin.H{
			"snackbar":  "setSuccess",
			"msg":       "修改通知成功！",
			"snackbar2": "closeSuccess",
		})
	}
}

func ReturnOK(c *gin.Context) {
	c.JSON(http.StatusOK, 1)
}
