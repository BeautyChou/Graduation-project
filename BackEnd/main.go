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

func main(){

	dsn := "root:123456@tcp(127.0.0.1:3306)/tttt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open("mysql", "root:123456@/tttt?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	//defer db.Close()
	Model.CreateDatabase(db)
	r:=gin.Default()
	r.Use(Middleware.Cors())
	r.POST("/test",Controller.Test)
	r.POST("/login",Controller.Auth)
	r.POST("/home",Middleware.JWTAuthMiddleWare(),Controller.HoeHandler)
	r.POST("/image",Controller.CheckedImage(db))
	r.GET("/image",Controller.ImageSend)
	r.GET("/imageTest",Controller.ImageSendTest)
	r.GET("/homeworkJudgeList",Controller.HomeworkJudgeList(db))
	r.GET("/getCourseList",Controller.GetCourseList(db))
	r.GET("/getHomeworkList",Controller.GetHomeworkList(db))
	r.DELETE("/deleteHomework",Controller.DeleteHomework(db))
	r.OPTIONS("/deleteHomework",Controller.ReturnOK)
	r.POST("/postQuestion",Controller.PostQuestion(db))
	r.GET("/getContentList",Controller.GetContentList(db))
	r.POST("/createHomework",Controller.CreateHomework(db))
	r.GET("/getNewHomeworkID",Controller.GetNewHomeworkID(db))
	r.GET("/getQuestionList",Controller.GetQuestionList(db))
	r.POST("/postHomework",Controller.PostHomework(db))
	r.GET("/getReviewList",Controller.GetReviewList(db))
	r.GET("/getClassesList",Controller.GetClassesList(db))
	r.POST("/postNewClass",Controller.PostNewClass(db))
	r.DELETE("/deleteClass",Controller.DeleteClass(db))
	r.OPTIONS("/deleteClass",Controller.ReturnOK)
	r.POST("/uploadClass",Controller.UploadClass(db))
	r.POST("/postApply",Controller.PostApply(db))
	r.GET("/getApply",Controller.GetApply(db))
	r.POST("/operateApply",Controller.OperateApply(db))
	r.POST("/validClassrooms",Controller.ValidClassrooms(db))
	r.Run(":9000")
}
