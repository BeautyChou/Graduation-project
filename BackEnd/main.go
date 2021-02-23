package main

import (
	"BackEnd/Controller"
	"BackEnd/Middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main(){
	db, err := gorm.Open("mysql", "root:123456@/tttt?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	r:=gin.Default()
	r.Use(Middleware.Cors())
	r.POST("/test",Controller.Test)
	r.POST("/login",Controller.Auth)
	r.POST("/home",Middleware.JWTAuthMiddleWare(),Controller.HoeHandler)
	r.POST("/image",Controller.CheckedImage(db))
	r.GET("/image",Controller.ImageSendTest)
	r.GET("/homeworkJudgeList",Controller.HomeworkJudgeList(db))
	r.Run(":9000")
}
