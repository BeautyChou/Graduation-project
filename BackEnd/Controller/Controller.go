package Controller

import (
	"BackEnd/Middleware"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
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

func ImageTest(c *gin.Context){
	file,err := c.FormFile("image")
	if err!= nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":err.Error(),
		})
		return
	}
	file.Filename = "blob2"
	log.Println(file.Filename)
	dst := fmt.Sprintf("D:/FinProject/BackEnd/images/%s",file.Filename)
	c.SaveUploadedFile(file,dst)
	c.JSON(http.StatusOK,gin.H{
		"message":fmt.Sprintf("'%s' uploaded!",file.Filename),
	})
}
func ParseBlob()  {
	file,err := os.Open("./images/blob")
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer file.Close()
	var content []byte
	var tmp = make([]byte,128)
	for{
		n,err := file.Read(tmp)
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		content = append(content,tmp[:n]...)
	}
	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := os.Create("./images/img.jpeg")
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 1

	err = jpeg.Encode(out, img, &opts)
	//jpeg.Encode(out, img, nil)
	if err != nil {
		log.Println(err)
	}
}
func ImageSendTest(c *gin.Context){
	homeworkid := c.Query("homework")
	studentid := c.Query("student")
	c.File("./images/"+homeworkid+"/"+studentid)
	//c.JSON(http.StatusOK,gin.H{
	//	"images":content,
	//})
}