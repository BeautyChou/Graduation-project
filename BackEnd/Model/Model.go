package Model

import (
	"gorm.io/gorm"
	"time"
)

type MyModel struct {
	ID        int `json:"id" gorm:"primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

//管理员表
type Admin struct {
	gorm.Model
	Name     string `json:"name" gorm:"name;size:50;uniqueIndex:name;not null"` // 管理员名字
	Password string `json:"password" gorm:"password;size:50;not null"`          // 密码
}

//学院表
type Faculty struct {
	MyModel
	Name string `json:"name" gorm:"size:30;uniqueIndex:name;not null"` // 学院名
}

//职称表
type Title struct {
	MyModel
	Name string `json:"name" gorm:"size:30;not null;uniqueIndex:name"` // 职称名
}

//教室表
type Classroom struct {
	MyModel
	Name string `json:"name" gorm:"size:30;not null;uniqueIndex:name;"` // 教室名
}

//教师表
type Teacher struct {
	MyModel
	Avatar    string  `json:"avatar" gorm:"size:50"`         // 头像
	Name      string  `json:"name" gorm:"size:50;not null"`  // 教师姓名
	Email     string  `json:"email" gorm:"size:50;"`         // 邮箱
	Phone     string  `json:"phone" gorm:"size:50;"`         // 电话
	Password  string  `json:"password" gorm:"size:50"`       // 密码
	Faculty   Faculty `gorm:"foreignKey:FacultyID;"`         // 学院外键
	FacultyID int     `json:"faculty_id" gorm:"index:f_id;"` // 学院ID
	Title     Title   `gorm:"foreignKey:TitleID;"`           // 职称外键
	TitleID   int     `json:"title_id" gorm:"index:t_id"`    // 职称ID
}

//所有课程
type Course struct {
	MyModel
	CourseName   string    `json:"course_name" gorm:"type:varchar(60);not null;"`             // 课程名
	Credit       string    `json:"credit" gorm:"type:varchar(50);"`                           // 学分
	Teacher      Teacher   `gorm:"ForeignKey:TeacherID;"`                                     // 教师外键
	TeacherID    int       `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"` // 教师ID
	Classroom    Classroom `gorm:"ForeignKey:ClassroomID;"`                                   // 教室外键
	ClassroomID  int       `json:"classroom_id" gorm:"type:int(11);index:classroom_id;"`      // 教室号
	MaxChooseNum int       `json:"max_choose_num" gorm:"type:int(11);not null;"`              // 最大可选课人数
	SelectedNum  int       `json:"selected_num" gorm:"type:int(11);default:0;"`               // 已选人数
	StartTime    time.Time `json:"start_time" gorm:"type:datetime;"`                          // 开始时间
	EndTime      time.Time `json:"end_time" gorm:"type:datetime;"`                            // 结束时间

}

//学生表
type Student struct {
	MyModel
	Name         string  `json:"name" gorm:"type:varchar(50);not null;"` // 学生姓名
	Sex          string  `json:"sex" gorm:"type:varchar(10);"`           // 性别
	Password     string  `json:"password" gorm:"type:varchar(20);"`      // 密码
	Faculty      Faculty `gorm:"ForeignKey:FacultyID;"`                  // 院系外键
	FacultyID    int     `json:"faculty_id" gorm:"index:faculty_id;"`    // 所在院系
	NativePlace  string  `json:"native_place" gorm:"type:varchar(60);"`  // 籍贯
	Credit       float32 `json:"mark" gorm:"type:float(5,2);"`           // 累计学分
	Email        string  `json:"email" gorm:"type:varchar(50);"`         // 电子邮件
	Avatar       string  `json:"avatar" gorm:"type:varchar(50);"`        // 头像
	Phone        string  `json:"phone" gorm:"type:varchar(50);"`         // 手机号码
	MaxChooseNum int     `json:"max_choose_num" gorm:"type:int(11);"`    // 最大可选课数
}

//学生成绩
type Elective struct {
	ID            int `json:"id" gorm:"primary_key;" sql:"type:INT(11) NOT NULL"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	Student       Student    `gorm:"ForeignKey:StudentNumber;"`                                                         // 学生外键
	StudentID     int        `json:"student_id" gorm:"primary_Key;uniqueIndex:student_id;" sql:"type:INT(11) NOT NULL"` // 学号
	Course        Course     `gorm:"ForeignKey:CourseID;"`                                                              // 课程外键
	CourseID      int        `json:"course_id" gorm:"primary_Key;uniqueIndex:course_id;" sql:"type:INT(11) NOT NULL"`   // 课程号
	TestScore     int        `json:"test_score" gorm:"size:3;"`                                                         // 考试成绩
	BehaviorScore int        `json:"behavior_score" gorm:"size:3;"`                                                     // 平时成绩
	Percentage    string     `json:"Percentage" gorm:"size:3"`                                                          // 考试成绩所占比例
}

type HomeworkUploadRecord struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	Homework   Homework   `gorm:"ForeignKey:HomeworkID;-"`
	HomeworkID int        `json:"homework_id" gorm:"primary_Key;uniqueIndex:homework_id;" sql:"type:INT(11) NOT NULL"` // 课程号
	QuestionID int        `json:"question_id" gorm:"primary_Key;uniqueIndex:homework_id;" sql:"type:INT(11) NOT NULL"` // 课程号
	Question   Question   `gorm:"ForeignKey:QuestionID;-"`
	StudentID  int        `json:"student_id" gorm:"primary_Key;uniqueIndex:student_id;" sql:"type:INT(11) NOT NULL"`
	Student    Student    `gorm:"ForeignKey:StudentNumber;-"`
	Score      int        `json:"score" gorm:"uniqueIndex:score;" sql:"type:INT(11) NOT NULL"` // 课程号
	IsScored   bool       `json:"is_scored" gorm:"uniqueIndex:is_scored;"`                     // 课程号
	IsUpload   bool       `json:"is_upload" gorm:"uniqueIndex:is_upload;"`
}

type Homework struct {
	ID               int `json:"id" gorm:"primary_key;" sql:"type:INT(11) NOT NULL"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	Question         Question   `gorm:"ForeignKey:QuestionID;-"`
	QuestionID       int        `json:"question_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	QuestionMaxScore int        `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time  `json:"deadline"`
	CourseID         int        `json:"course_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	Course           Course     `gorm:"ForeignKey:CourseID;-"`
	HomeworkTitle    string     `json:"homework_title" gorm:"uniqueIndex:homework_title"`
}

type Question struct {
	ID        int `json:"id" gorm:"primary_key;autoIncrement" sql:"type:INT(11) NOT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	Content   string     `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
}

type Student2Course struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	StudentID int        `json:"student_id" gorm:"primary_Key;uniqueIndex:student_id;" sql:"type:INT(11) NOT NULL"`
	Student   Student    `gorm:"ForeignKey:StudentID;-"`
	Course    Course     `gorm:"ForeignKey:CourseID;-"`
	CourseID  int        `json:"course_id" gorm:"primary_Key;uniqueIndex:course_id;" sql:"type:INT(11) NOT NULL"`
}

type HomeworkUploadRecordsForSelect struct {
	QuestionID       int    `json:"question_id" gorm:"primary_Key;uniqueIndex:homework_id;" sql:"type:INT(11) NOT NULL"` // 课程号
	Name             string `json:"name" gorm:"type:varchar(50);not null;"`                                              // 学生姓名
	StudentID        int    `json:"student_id" gorm:"primary_Key;uniqueIndex:student_id;" sql:"type:INT(11) NOT NULL"`
	Score            int    `json:"score" gorm:"uniqueIndex:score;" sql:"type:INT(11) NOT NULL"` // 课程号
	IsScored         bool   `json:"is_scored" gorm:"uniqueIndex:is_scored;"`                     // 课程号
	IsUpload         bool   `json:"is_upload" gorm:"uniqueIndex:is_upload;"`
	Content          string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
	QuestionMaxScore int    `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
}

type HomeworkForSelect struct {
	ID               int `json:"id" gorm:"primary_key;" sql:"type:INT(11) NOT NULL"`
	QuestionID       int `json:"question_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	QuestionMaxScore int `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time
	CourseID         int    `json:"course_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
}

type CourseForSelect struct {
	ID           int       `json:"course_id" gorm:"primary_key;"`
	CourseName   string    `json:"course_name" gorm:"type:varchar(60);not null;"`             // 课程名
	Credit       string    `json:"credit" gorm:"type:varchar(50);"`                           // 学分
	TeacherID    int       `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"` // 教师ID
	ClassroomID  int       `json:"classroom_id" gorm:"type:int(11);index:classroom_id;"`      // 教室号
	MaxChooseNum int       `json:"max_choose_num" gorm:"type:int(11);not null;"`              // 最大可选课人数
	SelectedNum  int       `json:"selected_num" gorm:"type:int(11);default:0;"`               // 已选人数
	StartTime    time.Time `json:"start_time" gorm:"type:datetime;"`                          // 开始时间
	EndTime      time.Time `json:"end_time" gorm:"type:datetime;"`                            // 结束时间
	Name         string    `json:"name" gorm:"size:50;not null"`                              // 教师姓名
}

type QuestionForSelect struct {
	ID      int    `json:"id" gorm:"primary_key;autoIncrement" sql:"type:INT(11) NOT NULL"`
	Content string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
}

type HomeworkForCreate struct {
	ID               int       `json:"id" gorm:"primary_key;" sql:"type:INT(11) NOT NULL"`
	QuestionID       int       `json:"question_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	QuestionMaxScore int       `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time `json:"deadline"`
	CourseID         int       `json:"course_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	HomeworkTitle    string    `json:"homework_title" gorm:"uniqueIndex:homework_title"`
}

type QuestionForCheck struct {
	QuestionMaxScore int `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
	Content          string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
	QuestionId       int    `json:"question_id"`
}

type QuestionForStudent struct {
	QuestionMaxScore int `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
	Content          string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
	QuestionId       int    `json:"question_id"`
	Uploaded         bool   `json:"uploaded"`
}

type RecordForReview struct {
	QuestionID       int    `json:"question_id" gorm:"primary_Key;uniqueIndex:homework_id;" sql:"type:INT(11) NOT NULL"` // 课程号
	Score            int    `json:"score" gorm:"uniqueIndex:score;" sql:"type:INT(11) NOT NULL"` // 课程号
	IsScored         bool   `json:"is_scored" gorm:"uniqueIndex:is_scored;"`                     // 课程号
	Content          string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
	QuestionMaxScore int    `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
}
type HomeworkUploadRecordsForSelects []HomeworkUploadRecordsForSelect
type HomeworkUploadRecords []HomeworkUploadRecord
type HomeworkForSelects []HomeworkForSelect
type CourseForSelects []CourseForSelect
type QuestionForSelects []QuestionForSelect
type Homeworks []Homework
type QuestionForChecks []QuestionForCheck
type QuestionForStudents []QuestionForStudent
type RecordForReviews []RecordForReview

func CreateDatabase(db *gorm.DB) {
	db.AutoMigrate(&Question{})
	//db.AutoMigrate(&Title{}, &Faculty{}, &Teacher{}, &Course{}, &Elective{}, &Admin{}, &Classroom{}, &Student{}, &HomeworkUploadRecord{}, &Homework{}, &Question{},&Student2Course{})
	//db.Model(&Teacher{}).AddForeignKey("title_id", "titles(id)", "no action", "no action")
	//db.Model(&Teacher{}).AddForeignKey("faculty_id", "faculties(id)", "cascade", "no action")
	//db.Model(&Course{}).AddForeignKey("teacher_id", "teachers(id)", "cascade", "no action")
	//db.Model(&Course{}).AddForeignKey("classroom_id", "classrooms(id)", "no action", "no action")
	//db.Model(&Student{}).AddForeignKey("faculty_id", "faculties(id)", "cascade", "no action")
	//db.Model(&Elective{}).AddForeignKey("student_id", "students(id)", "cascade", "no action")
	//db.Model(&Elective{}).AddForeignKey("course_id", "courses(id)", "cascade", "no action")
	//db.Model(&Elective{}).AddForeignKey("course_id", "courses(id)", "cascade", "no action")
	//db.Model(&HomeworkUploadRecord{}).AddForeignKey("homework_id", "homeworks(id)", "no action", "no action")
	//db.Model(&HomeworkUploadRecord{}).AddForeignKey("question_id", "questions(id)", "no action", "no action")
	//db.Model(&HomeworkUploadRecord{}).AddForeignKey("student_id", "students(id)", "no action", "no action")
	//db.Model(&Homework{}).AddForeignKey("question_id", "questions(id)", "no action", "no action")
	//db.Model(&Homework{}).AddForeignKey("course_id", "courses(id)", "no action", "no action")
	//db.Model(&Student2Course{}).AddForeignKey("course_id", "courses(id)", "no action", "no action")
	//db.Model(&Student2Course{}).AddForeignKey("student_id", "students(id)", "no action", "no action")
}
