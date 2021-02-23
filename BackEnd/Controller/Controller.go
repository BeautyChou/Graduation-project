package Controller

import (
	"BackEnd/Middleware"
	"BackEnd/Model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type HomeWork struct {
	HomeworkID int `form:"homework_id" json:"homework_id" binding:"required"`
}
var HomeworkJudgeLists Model.HomeworkUploadRecordsForSelects

func Test(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"msg":"ok",
	})
}

func Auth(c *gin.Context){
	var user UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"msg":"无效的参数",
		})
	}
	if user.Username == "123" && user.Password=="456"{
		TokenString,_ := Middleware.GenToken(user.Username)
		c.JSON(http.StatusOK,gin.H{
			"msg":"success",
			"data":gin.H{"token":TokenString},
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"failed",
	})
	return
}

func HoeHandler(c *gin.Context){
	username := c.MustGet("username").(string)
	c.JSON(http.StatusOK,gin.H{
		"msg":"success",
		"data":gin.H{
			"username":username,
		},
	})

}


func CheckedImage (db *gorm.DB) func(ctx *gin.Context){
	return  func(c *gin.Context){
		file,err := c.FormFile("image")
		studentID := c.PostForm("student_id")
		homeworkID := c.PostForm("homework_id")
		questionID := c.PostForm("question_id")
		fmt.Println(studentID,homeworkID,questionID)
		score := c.PostForm("score")
		if err!= nil {
			c.JSON(http.StatusInternalServerError,gin.H{
				"message":err.Error(),
			})
			return
		}
		file.Filename = studentID+"-checked"
		log.Println(file.Filename)
		dst := fmt.Sprintf("D:/FinProject/BackEnd/images/%s/%s/%s",homeworkID,questionID,file.Filename)
		c.SaveUploadedFile(file,dst)
		db.Raw("UPDATE homework_upload_records SET score = ?,is_scored = 1 where homework_id = ? AND question_id = ? AND student_id = ?",score,homeworkID,questionID,studentID)
		c.JSON(http.StatusOK,gin.H{
			"message":fmt.Sprintf("'%s' uploaded!",file.Filename),
		})
	}
}

func ImageSendTest(c *gin.Context){
	homeworkid := c.Query("homework")
	questionid := c.Query("question")
	studentid := c.Query("student")
	c.File("./images/"+homeworkid+"/"+questionid+"/"+studentid)
}

func HomeworkJudgeList (db *gorm.DB) func(*gin.Context){
	return func(c *gin.Context) {
		var checked int64
		homework := c.Query("homework_id")
		fmt.Println(homework)
		db.Raw("SELECT name,homework_upload_records.question_id,student_id,score,is_scored,content,question_max_score from homework_upload_records RIGHT JOIN questions ON homework_upload_records.question_id = questions.id LEFT JOIN students On students.id=homework_upload_records.student_id RIGHT JOIN homeworks on homeworks.id = homework_upload_records.homework_id AND homeworks.question_id = homework_upload_records.question_id where homework_id = ?",homework).Scan(&HomeworkJudgeLists)
		db.Model(&Model.HomeworkUploadRecord{}).Where("homework_id = ? AND is_scored = 1",homework).Count(&checked)
		c.JSON(http.StatusOK,gin.H{
			"lists":HomeworkJudgeLists,
			"checked":checked,
		})
	}
}