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
		db.Raw("Select courses.id,course_name,credit,teacher_id,classroom_id,max_choose_num,selected_num,start_time,end_time,name from courses right join teachers on courses.teacher_id = teachers.id where teacher_id = ? AND courses.deleted_at is NULL AND teachers.deleted_at is null", teacherId).Scan(&courses)
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
		})
	}
}

func GetHomeworkList(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var courseId = c.Query("course_id")
		var homeworks Model.HomeworkForSelects
		db.Raw("Select * from homeworks where course_id = ? AND deleted_at is NULL GROUP BY id", courseId).Scan(&homeworks)
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
		id := utils.ToString(questions[0].ID)
		for _, item := range questions {
			err := os.MkdirAll("./images/"+id+"/"+utils.ToString(item.QuestionID), 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		db.Where("id = ?", questions[0].ID).Delete(&Model.Homework{})
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
			for _,value := range records{
				for i,v := range questions{
					fmt.Println(v.QuestionId == value.QuestionID)
					if v.QuestionId == value.QuestionID {
						questions[i].Uploaded = true
						break
					}
				}
			}
		}
		c.JSON(http.StatusOK,gin.H{
			"questions":questions,
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
		c.SaveUploadedFile(file, dst)
		db.Debug().Create(&record)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	}
}

func GetReviewList(db *gorm.DB) func(c *gin.Context){
	return func(c *gin.Context) {
		var reviwes Model.RecordForReviews
		homeworkID := c.Query("homework_id")
		studentId := c.Query("student_id")
		db.Raw("select * from homework_upload_records right join homeworks h on h.id = homework_upload_records.homework_id AND h.question_id = homework_upload_records.question_id right join questions q on q.id = h.question_id where homework_id =? AND student_id = ?",homeworkID,studentId).Scan(&reviwes)
		c.JSON(http.StatusOK,gin.H{
			"reviews":reviwes,
		})
	}
}

func ReturnOK(c *gin.Context) {
	c.JSON(http.StatusOK, 0)
}
