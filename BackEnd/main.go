package main

import "github.com/gin-gonic/gin"
import "BackEnd/Controller"

func main(){
	r:=gin.Default()
	r.Use(Controller.Cors())
	r.POST("/test",Controller.Test)
	r.Run(":9000")
}
