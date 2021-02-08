package Controller

import (
	"BackEnd/Middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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
