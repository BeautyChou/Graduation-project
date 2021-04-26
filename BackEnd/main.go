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
	"io"
	"os"
)

type Config struct {
	Port          int    `mapstructure:"port"`
	InitDatabase  int    `mapstructure:"init_database"`
	DatabaseName  string `mapstructure:"database_name"`
	MysqlUsername string `mapstructure:"mysql_username"`
	MysqlPassword string `mapstructure:"mysql_password"`
	MysqlPort     string `mapstructure:"mysql_port"`
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
	dsn := Conf.MysqlUsername + ":" + Conf.MysqlPassword + "@tcp(127.0.0.1:" + Conf.MysqlPort + ")/" + Conf.DatabaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	if Conf.InitDatabase == 1 {
		Model.CreateDatabase(db)
		viper.Set("init_database", 0)
		err = viper.WriteConfigAs("config.yaml")
	}
	//定期清理上学期课程
	go Controller.StartPeriodicityTask(db)

	//启动服务
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.Use(Middleware.Cors())
	r.POST("/login", Controller.Auth(db))
	TestConnection := r.Group("/", Middleware.Cors(), Controller.ReturnOK)
	{
		TestConnection.OPTIONS("/Homework")
		TestConnection.OPTIONS("/Class")
		TestConnection.OPTIONS("/Apply")
		TestConnection.OPTIONS("/Course")
		TestConnection.OPTIONS("/Practice")
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
		//批改作业学生列表
		AuthGroup.GET("/homeworkJudgeList", Controller.HomeworkJudgeList(db))
		//教师获取课程列表
		AuthGroup.GET("/getCourseList", Controller.GetCourseList(db))
		//学生提交作业
		AuthGroup.POST("/postHomework", Controller.PostHomework(db))
		//学生获取老师批改作业
		AuthGroup.GET("/getReviewList", Controller.GetReviewList(db))
		//查询可安排作业的课程
		AuthGroup.POST("/validClassrooms", Controller.ValidClassrooms(db))
		//查询课程表
		AuthGroup.GET("/getClassSheet", Controller.GetClassSheet(db))
		//检验用户JWT是否过期
		AuthGroup.POST("/isExpire")
		//题目及其相关
		Question := AuthGroup.Group("/")
		{
			Question.GET("/SelectedQuestion", Controller.GetQuestionList(db))
			Question.POST("/Question", Controller.PostQuestion(db))
			Question.GET("/Question", Controller.GetContentList(db))
		}
		//作业及其相关
		Homework := AuthGroup.Group("/")
		{
			Homework.GET("/Homework", Controller.GetHomeworkList(db))
			Homework.DELETE("/Homework", Controller.DeleteHomework(db))
			Homework.POST("/Homework", Controller.CreateHomework(db))
			Homework.GET("/getNewHomeworkID", Controller.GetNewHomeworkID(db))
		}
		//设置课程相关
		Class := AuthGroup.Group("/")
		{
			Class.GET("/Class", Controller.GetClassesList(db))
			Class.POST("/Class", Controller.PostNewClass(db))
			Class.DELETE("/Class", Controller.DeleteClass(db))
			Class.PUT("/Class", Controller.UploadClass(db))
		}
		//课程修改申请相关
		Apply := AuthGroup.Group("/")
		{
			Apply.POST("/Apply", Controller.PostApply(db))
			Apply.GET("/Apply", Controller.GetApply(db))
			Apply.PUT("/Apply", Controller.OperateApply(db))
		}
		//选、退课相关
		Course := AuthGroup.Group("/")
		{
			Course.GET("/Course", Controller.GetAvailableCourses(db))
			Course.POST("/Course", Controller.ChooseCourse(db))
			Course.GET("/ChosenCourse", Controller.GetChosenCourse(db))
			Course.DELETE("/Course", Controller.QuitCourse(db))
		}
		//提交、查询课程成绩相关
		Score := AuthGroup.Group("/")
		{
			Score.GET("/StudentScore", Controller.GetStudentScore(db))
			Score.POST("/Score", Controller.PostScore(db))
			Score.GET("/Score", Controller.GetScore(db))
		}
		//用户信息相关
		UserInfo := AuthGroup.Group("/UserInfo")
		{
			UserInfo.GET("/UserInfo", Controller.GetUserInfo(db))
			UserInfo.GET("/Direction", Controller.GetDirectionList(db))
			UserInfo.POST("/Direction", Controller.PostDirection(db))
			UserInfo.GET("/Teacher", Controller.GetTeacherList(db))
			UserInfo.POST("/Teacher", Controller.PostTeacher(db))
			UserInfo.POST("/ApplyTeacher", Controller.ApplyTeacher(db))
		}
		//学生实习相关
		Practice := AuthGroup.Group("/")
		{
			Practice.GET("/Practice", Controller.GetPracticeInfo(db))
			Practice.PUT("/Practice", Controller.PutPractice(db))
			Practice.POST("/Practice", Controller.PostPractice(db))
		}
		//学生管理相关
		Student := AuthGroup.Group("/")
		{
			Student.GET("/Student", Controller.GetStudentList(db))
			Student.DELETE("/Student", Controller.DeleteStudent(db))
			Student.PUT("/Student", Controller.PutStudent(db))
			Student.POST("/Student", Controller.AddStudent(db))
		}
		//处分管理相关
		Punishment := AuthGroup.Group("/")
		{
			Punishment.POST("/Punishment", Controller.PunishStudent(db))
			Punishment.GET("/Punishment", Controller.GetPunishment(db))
			Punishment.DELETE("/Punishment", Controller.DeletePunishment(db))
		}
		//教师管理相关
		Teacher := AuthGroup.Group("/")
		{
			Teacher.GET("/Teacher", Controller.GetTeacherListForManage(db))
			Teacher.DELETE("/Teacher", Controller.DeleteTeacher(db))
			Teacher.PUT("/Teacher", Controller.PutTeacher(db))
			Teacher.POST("/Teacher", Controller.AddTeacher(db))
		}
		//管理员管理相关
		Admin := AuthGroup.Group("/")
		{
			Admin.GET("/Admin", Controller.GetAdminList(db))
			Admin.DELETE("/Admin", Controller.DeleteAdmin(db))
			Admin.PUT("/Admin", Controller.PutAdmin(db))
			Admin.POST("/Admin", Controller.AddAdmin(db))
		}
		//职称管理相关
		Title := AuthGroup.Group("/")
		{
			Title.GET("/Title", Controller.GetTitleList(db))
			Title.DELETE("/Title", Controller.DeleteTitle(db))
			Title.PUT("/Title", Controller.PutTitle(db))
			Title.POST("/Title", Controller.AddTitle(db))
		}
		//处分等级相关
		PunishLevel := AuthGroup.Group("/")
		{
			PunishLevel.GET("/PunishLevel", Controller.GetPunishLevelList(db))
			PunishLevel.DELETE("/PunishLevel", Controller.DeletePunishLevel(db))
			PunishLevel.PUT("/PunishLevel", Controller.PutPunishLevel(db))
			PunishLevel.POST("/PunishLevel", Controller.AddPunishLevel(db))
		}
		//教室管理相关
		Classroom := AuthGroup.Group("/")
		{
			Classroom.GET("/Classroom", Controller.GetClassroomList(db))
			Classroom.DELETE("/Classroom", Controller.DeleteClassroom(db))
			Classroom.PUT("/Classroom", Controller.PutClassroom(db))
			Classroom.POST("/Classroom", Controller.AddClassroom(db))
		}
		//学院、专业、方向相关
		DirectionSpecialtyFaculty := AuthGroup.Group("/")
		{
			DirectionSpecialtyFaculty.GET("/DirectionSpecialtyFaculty", Controller.GetDirectionSpecialtyFacultyList(db))
			DirectionSpecialtyFaculty.POST("/DirectionSpecialtyFaculty", Controller.AddFacultySpecialtyDirection(db))
			DirectionSpecialtyFaculty.DELETE("/DirectionSpecialtyFaculty", Controller.DeleteFacultySpecialtyDirection(db))
		}
		//系统通知相关
		Notification := AuthGroup.Group("/")
		{
			Notification.GET("/Notification", Controller.GetNotification(db))
			Notification.PUT("/Notification", Controller.PutNotification(db))
		}
	}
	r.Run(":" + utils.ToString(Conf.Port))

}
