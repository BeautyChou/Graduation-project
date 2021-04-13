package main

import (
	"BackEnd/Controller"
	"BackEnd/Middleware"
	"BackEnd/Model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

type Config struct {
	Port         int `mapstructure:"port"`
	InitDatabase int `mapstructure:"init_database"`
	DatabaseName string `mapstructure:"database_name"`
}

var Conf = new(Config)

func main() {

	//读取配置文件并进行配置
	viper.AddConfigPath("./")
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}

	//连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/"+Conf.DatabaseName+"?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	if Conf.InitDatabase == 1 {
		Model.CreateDatabase(db)
		viper.Set("init_database",0)
		err = viper.WriteConfigAs("config.yaml")
	}
	//启动服务
	r := gin.Default()
	r.Use(Middleware.Cors())
	r.POST("/test", Controller.Test)
	r.POST("/login", Controller.Auth(db))
	TestConnection := r.Group("/", Middleware.Cors(), Controller.ReturnOK)
	{
		TestConnection.OPTIONS("/deleteHomework")
		TestConnection.OPTIONS("/deleteClass")
		TestConnection.OPTIONS("/quitCourse")
		TestConnection.OPTIONS("/Student")
		TestConnection.OPTIONS("/Punishment")
		TestConnection.OPTIONS("/Teacher")
		TestConnection.OPTIONS("/Admin")
		TestConnection.OPTIONS("/Title")
		TestConnection.OPTIONS("/PunishLevel")
		TestConnection.OPTIONS("/Classroom")
		TestConnection.OPTIONS("/DirectionSpecialtyFaculty")
		TestConnection.OPTIONS("/Notification")
	}
	Image := r.Group("/", Middleware.Cors())
	{
		Image.POST("/image", Controller.CheckedImage(db))
		Image.GET("/image", Controller.ImageSend)
		Image.GET("/imageChecked", Controller.SendCheckedImage)
	}
	AuthGroup := r.Group("/", Middleware.Cors(), Middleware.JWTAuthMiddleWare())
	{
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
		AuthGroup.POST("/ApplyTeacher", Controller.ApplyTeacher(db))
		Practice := AuthGroup.Group("/")
		{
			Practice.GET("/Practice", Controller.GetPracticeInfo(db))
			Practice.PUT("/Practice", Controller.PutPractice(db))
			Practice.POST("/Practice", Controller.PostPractice(db))
		}
		Student := AuthGroup.Group("/")
		{
			Student.GET("/Student", Controller.GetStudentList(db))
			Student.DELETE("/Student", Controller.DeleteStudent(db))
			Student.PUT("/Student", Controller.PutStudent(db))
			Student.POST("/Student", Controller.AddStudent(db))
		}
		Punishment := AuthGroup.Group("/")
		{
			Punishment.POST("/Punishment", Controller.PunishStudent(db))
			Punishment.GET("/Punishment", Controller.GetPunishment(db))
			Punishment.DELETE("/Punishment", Controller.DeletePunishment(db))
		}
		Teacher := AuthGroup.Group("/")
		{
			Teacher.GET("/Teacher", Controller.GetTeacherListForManage(db))
			Teacher.DELETE("/Teacher", Controller.DeleteTeacher(db))
			Teacher.PUT("/Teacher", Controller.PutTeacher(db))
			Teacher.POST("/Teacher", Controller.AddTeacher(db))
		}
		Admin := AuthGroup.Group("/")
		{
			Admin.GET("/Admin", Controller.GetAdminList(db))
			Admin.DELETE("/Admin", Controller.DeleteAdmin(db))
			Admin.PUT("/Admin", Controller.PutAdmin(db))
			Admin.POST("/Admin", Controller.AddAdmin(db))
		}
		Title := AuthGroup.Group("/")
		{
			Title.GET("/Title", Controller.GetTitleList(db))
			Title.DELETE("/Title", Controller.DeleteTitle(db))
			Title.PUT("/Title", Controller.PutTitle(db))
			Title.POST("/Title", Controller.AddTitle(db))
		}
		PunishLevel := AuthGroup.Group("/")
		{
			PunishLevel.GET("/PunishLevel", Controller.GetPunishLevelList(db))
			PunishLevel.DELETE("/PunishLevel", Controller.DeletePunishLevel(db))
			PunishLevel.PUT("/PunishLevel", Controller.PutPunishLevel(db))
			PunishLevel.POST("/PunishLevel", Controller.AddPunishLevel(db))
		}
		Classroom := AuthGroup.Group("/", Middleware.JWTAuthMiddleWare())
		{
			Classroom.GET("/Classroom", Controller.GetClassroomList(db))
			Classroom.DELETE("/Classroom", Controller.DeleteClassroom(db))
			Classroom.PUT("/Classroom", Controller.PutClassroom(db))
			Classroom.POST("/Classroom", Controller.AddClassroom(db))
		}
		DirectionSpecialtyFaculty := AuthGroup.Group("/")
		{
			DirectionSpecialtyFaculty.GET("/DirectionSpecialtyFaculty", Controller.GetDirectionSpecialtyFacultyList(db))
			DirectionSpecialtyFaculty.POST("/DirectionSpecialtyFaculty", Controller.AddFacultySpecialtyDirection(db))
			DirectionSpecialtyFaculty.DELETE("/DirectionSpecialtyFaculty", Controller.DeleteFacultySpecialtyDirection(db))
		}
		Notification := AuthGroup.Group("/")
		{
			Notification.GET("/Notification", Controller.GetNotification(db))
			Notification.PUT("/Notification", Controller.PutNotification(db))
		}
		AuthGroup.POST("/isExpire")
	}
	r.Run(":" + utils.ToString(Conf.Port))

}
