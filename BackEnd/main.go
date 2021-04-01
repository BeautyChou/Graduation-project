package main

import (
	"BackEnd/Controller"
	"BackEnd/Middleware"
	"BackEnd/Model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/tttt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open("mysql", "root:123456@/tttt?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	//defer db.Close()
	Model.CreateDatabase(db)
	r := gin.Default()
	r.Use(Middleware.Cors())
	r.POST("/test", Controller.Test)
	r.POST("/login", Controller.Auth(db))
	r.POST("/home", Middleware.JWTAuthMiddleWare(), Controller.HoeHandler)
	r.OPTIONS("/deleteHomework", Controller.ReturnOK)
	r.OPTIONS("/deleteClass", Controller.ReturnOK)
	r.OPTIONS("/quitCourse", Controller.ReturnOK)
	AuthGroup := r.Group("/",Middleware.Cors(),Middleware.JWTAuthMiddleWare())
	{
		AuthGroup.POST("/image", Controller.CheckedImage(db))
		AuthGroup.GET("/image", Controller.ImageSend)
		AuthGroup.GET("/imageTest", Controller.ImageSendTest)
		AuthGroup.GET("/homeworkJudgeList", Controller.HomeworkJudgeList(db))
		AuthGroup.GET("/getCourseList", Controller.GetCourseList(db))
		AuthGroup.GET("/getHomeworkList", Controller.GetHomeworkList(db))
		AuthGroup.DELETE("/deleteHomework", Controller.DeleteHomework(db))
		AuthGroup.POST("/postQuestion", Controller.PostQuestion(db))
		AuthGroup.GET("/getContentList", Controller.GetContentList(db))
		AuthGroup.POST("/createHomework", Controller.CreateHomework(db))
		AuthGroup.GET("/getNewHomeworkID", Controller.GetNewHomeworkID(db))
		AuthGroup.GET("/getQuestionList", Controller.GetQuestionList(db))
		AuthGroup.POST("/postHomework", Controller.PostHomework(db))
		AuthGroup.GET("/getReviewList", Controller.GetReviewList(db))
		AuthGroup.GET("/getClassesList", Controller.GetClassesList(db))
		AuthGroup.POST("/postNewClass", Controller.PostNewClass(db))
		AuthGroup.DELETE("/deleteClass", Controller.DeleteClass(db))
		AuthGroup.POST("/uploadClass", Controller.UploadClass(db))
		AuthGroup.POST("/postApply", Controller.PostApply(db))
		AuthGroup.GET("/getApply", Controller.GetApply(db))
		AuthGroup.POST("/operateApply", Controller.OperateApply(db))
		AuthGroup.POST("/validClassrooms", Controller.ValidClassrooms(db))
		AuthGroup.GET("/getClassSheet", Controller.GetClassSheet(db))
		AuthGroup.GET("/getAvailableCourses", Controller.GetAvailableCourses(db))
		AuthGroup.POST("/chooseCourse", Controller.ChooseCourse(db))
		AuthGroup.GET("/getChosenCourse", Controller.GetChosenCourse(db))
		AuthGroup.DELETE("/quitCourse", Controller.QuitCourse(db))
		AuthGroup.GET("/getStudentScore", Controller.GetStudentScore(db))
		AuthGroup.POST("/postScore", Controller.PostScore(db))
		AuthGroup.GET("/getScore", Controller.GetScore(db))
		AuthGroup.GET("/getUserInfo", Controller.GetUserInfo(db))
		AuthGroup.GET("/getDirectionList", Controller.GetDirectionList(db))
		AuthGroup.GET("/getTeacherList", Controller.GetTeacherList(db))
		AuthGroup.POST("/postDirection", Controller.PostDirection(db))
		AuthGroup.POST("/postTeacher", Controller.PostTeacher(db))
		AuthGroup.POST("/postPractice", Controller.PostPractice(db))
		AuthGroup.POST("/postIndependentPractice", Controller.PostIndependentPractice(db))
		AuthGroup.POST("/ApplyTeacher", Controller.ApplyTeacher(db))
		AuthGroup.POST("/isExpire")
	}
	r.Run(":9000")
}
