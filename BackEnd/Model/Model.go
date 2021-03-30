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
	Name     string `json:"name" gorm:"name;size:50;not null"`         // 管理员名字
	Password string `json:"password" gorm:"password;size:50;not null"` // 密码
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

type DirectionToSpecialty struct {
	MyModel
	Faculty       Faculty `gorm:"ForeignKey:FacultyID;"`                                           // 专业所属学院
	FacultyID     int     `json:"faculty_id" gorm:"primary_key;type:int(11)"`                      // 学院ID
	SpecialtyName string  `json:"specialty_name" gorm:"type:varchar(60)"`                          // 专业名称
	DirectionID   int     `json:"direction_id" gorm:"primary_key;index:direction_id;type:int(11)"` // 学生所属方向
	SpecialtyID   int     `json:"specialty_id" gorm:"primary_key;index:special_id;type:int(11)"`   // 学生所属专业
	DirectionName string  `json:"direction_name" gorm:"type:varchar(60)"`                          // 方向名称
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
	RecordID     int `json:"id" gorm:"primary_key"` //记录id
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt       `sql:"index"`
	ID           int                  `json:"course_id" gorm:"primary_key;type:int(11);not null;index:course_id;"` //课程号
	CourseName   string               `json:"course_name" gorm:"type:varchar(60);not null;"`                       // 课程名
	Credit       string               `json:"credit" gorm:"type:varchar(50);"`                                     // 学分
	Teacher      Teacher              `gorm:"ForeignKey:TeacherID;"`                                               // 教师外键
	TeacherID    int                  `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"`           // 教师ID
	Classroom    Classroom            `gorm:"ForeignKey:ClassroomID;"`                                             // 教室外键
	ClassroomID  int                  `json:"classroom_id" gorm:"type:int(11);index:classroom_id;"`                // 教室号
	MaxChooseNum int                  `json:"max_choose_num" gorm:"type:int(11);not null;"`                        // 最大可选课人数
	SelectedNum  int                  `json:"selected_num" gorm:"type:int(11);default:0;"`                         // 已选人数
	StartTime    int                  `json:"start_time" gorm:"type:int(11);"`                                     // 开始节数
	EndTime      int                  `json:"end_time" gorm:"type:int(11);"`                                       // 结束节数
	WeekTime     int                  `json:"week_time" gorm:"type:int(2)"`                                        // 上课星期数
	Faculty      Faculty              `gorm:"ForeignKey:FacultyID;"`                                               // 上课学生所属学院
	FacultyID    int                  `json:"faculty_id" gorm:"type:int(11)"`                                      // 学院ID
	Direction    DirectionToSpecialty `gorm:"references:DirectionID;"`                                             // 方向外键
	DirectionID  int                  `json:"direction_id" gorm:"index:direction_id;type:int(11)"`                 // 上课学生所属方向
	Specialty    DirectionToSpecialty `gorm:"references:SpecialtyID;"`                                             // 专业外键
	SpecialtyID  int                  `json:"specialty_id" gorm:"index:specialty_id;type:int(11)"`                 // 上课学生所属专业
	StartWeek    int                  `json:"start_week" gorm:"type:int(11);"`                                     // 起始周数
	EndWeek      int                  `json:"end_week" gorm:"type:int(11);"`                                       // 结束周数
	CopyFlag     int                  `json:"copy_flag"`                                                           // 记录复制的Record_ID
	Selectable   bool                 `json:"selectable"`                                                          //区分毕业设计与普通课程的标记
}

type ApplyForCourseChange struct {
	MyModel
	Course            Course  `gorm:"ForeignKey:CourseID;"`                                               // 课程外键
	CourseID          int     `json:"course_id" gorm:"type:int(11)" sql:"type:INT(11) NOT NULL"`          // 课程号
	Teacher           Teacher `gorm:"ForeignKey:TeacherID;"`                                              // 教师外键
	TeacherID         int     `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"`          // 教师ID
	Reason            string  `json:"reason" gorm:"varchar(60)"`                                          // 申请更改的理由
	BeforeClassroomID int     `json:"before_classroom_id" gorm:"type:int(11);index:before_classroom_id;"` // 教室号
	AfterClassroomID  int     `json:"after_classroom_id" gorm:"type:int(11);index:after_classroom_id;"`   // 教室号
	BeforeStartTime   int     `json:"before_start_time" gorm:"type:int(11);"`                             // 开始节数
	AfterStartTime    int     `json:"after_start_time" gorm:"type:int(11);"`                              // 开始节数
	BeforeEndTime     int     `json:"before_end_time" gorm:"type:int(11);"`                               // 结束节数
	AfterEndTime      int     `json:"after_end_time" gorm:"type:int(11);"`                                // 结束节数
	BeforeWeekTime    int     `json:"before_week_time" gorm:"type:int(2)"`                                // 上课星期数
	AfterWeekTime     int     `json:"after_week_time" gorm:"type:int(2)"`                                 // 上课星期数
	BeforeStartWeek   int     `json:"before_start_week" gorm:"type:int(11);"`                             // 起始周数
	AfterStartWeek    int     `json:"after_start_week" gorm:"type:int(11);"`                              // 起始周数
	BeforeEndWeek     int     `json:"before_end_week" gorm:"type:int(11);"`                               // 结束周数
	AfterEndWeek      int     `json:"after_end_week" gorm:"type:int(11);"`                                // 结束周数
	Result            int     `json:"result" gorm:"type:int(2);default:0"`                                // 审批结果
}

//学生表
type Student struct {
	MyModel
	Name         string               `json:"name" gorm:"type:varchar(50);not null;"`                    // 学生姓名
	Sex          string               `json:"sex" gorm:"type:varchar(10);"`                              // 性别
	Password     string               `json:"password" gorm:"type:varchar(20);"`                         // 密码
	Faculty      Faculty              `gorm:"ForeignKey:FacultyID;"`                                     // 院系外键
	FacultyID    int                  `json:"faculty_id" gorm:"index:faculty_id;"`                       // 所在院系
	NativePlace  string               `json:"native_place" gorm:"type:varchar(60);"`                     // 籍贯
	Credit       float32              `json:"mark" gorm:"type:float(5,2);"`                              // 累计学分
	Email        string               `json:"email" gorm:"type:varchar(50);"`                            // 电子邮件
	Avatar       string               `json:"avatar" gorm:"type:varchar(50);"`                           // 头像
	Phone        string               `json:"phone" gorm:"type:varchar(50);"`                            // 手机号码
	MaxChooseNum int                  `json:"max_choose_num" gorm:"type:int(11);"`                       // 最大可选课数
	Direction    DirectionToSpecialty `gorm:"references:DirectionID;"`                                   // 方向外键
	DirectionID  int                  `json:"direction_id" gorm:"index:direction_id;type:int(11)"`       // 学生所属方向
	SpecialtyID  int                  `json:"specialty_id" gorm:"index:special_id;type:int(11)"`         // 学生所属专业
	Specialty    DirectionToSpecialty `gorm:"references:SpecialtyID;"`                                   // 专业外键
	Teacher      Teacher              `gorm:"references:ID;"`                                            //导师外键
	TeacherID    int                  `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"` // 导师ID
	Practice     int               `json:"practice"`                                                  //实习方式
}

//学生成绩
type Elective struct {
	ID                 int `json:"id" gorm:"primary_key;" sql:"type:INT(11) NOT NULL"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `sql:"index"`
	Student            Student        `gorm:"ForeignKey:StudentID;"`                                                 // 学生外键
	StudentID          int            `json:"student_id" gorm:"primary_Key:student_id;" sql:"type:INT(11) NOT NULL"` // 学号
	Course             Course         `gorm:"ForeignKey:CourseID;"`                                                  // 课程外键
	CourseID           int            `json:"course_id" gorm:"primary_Key:course_id;" sql:"type:INT(11) NOT NULL"`   // 课程号
	HomeworkScore      int            `json:"homework_score" gorm:"size:3"`                                          // 作业成绩
	TestScore          int            `json:"test_score" gorm:"size:3;"`                                             // 考试成绩
	BehaviorScore      int            `json:"behavior_score" gorm:"size:3;"`                                         // 平时成绩
	Percentage         int            `json:"percentage" gorm:"size:3"`                                              // 考试成绩所占比例
	HomeworkPercentage int            `json:"homework_percentage" gorm:"size:3"`                                     // 作业分数占比
}

type HomeworkUploadRecord struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `sql:"index"`
	Homework   Homework       `gorm:"ForeignKey:HomeworkID;-"`
	HomeworkID int            `json:"homework_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"` // 课程号
	QuestionID int            `json:"question_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"` // 课程号
	RecordID   int            `json:"record_id" gorm:"primary_key"`                                // 记录id
	Question   Question       `gorm:"ForeignKey:QuestionID;-"`
	StudentID  int            `json:"student_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"`
	Student    Student        `gorm:"ForeignKey:StudentNumber;-"`
	Score      int            `json:"score" gorm:"" sql:"type:INT(11) NOT NULL"` // 课程号
	IsScored   bool           `json:"is_scored" gorm:""`                         // 课程号
	IsUpload   bool           `json:"is_upload" gorm:""`
}

type Homework struct {
	ID               int `json:"id" gorm:"primary_key;" sql:"type:INT(11) NOT NULL"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `sql:"index"`
	Question         Question       `gorm:"ForeignKey:QuestionID;-"`
	QuestionID       int            `json:"question_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"`
	QuestionMaxScore int            `json:"question_max_score" gorm:"" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time      `json:"deadline"`
	CourseID         int            `json:"course_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"`
	Course           Course         `gorm:"ForeignKey:CourseID;-"`
	HomeworkTitle    string         `json:"homework_title" gorm:"size:190"`
	RecordID         int            `json:"record_id" gorm:"primary_key"` // 记录id
}

type Question struct {
	ID        int `json:"id" gorm:"primary_key;autoIncrement" sql:"type:INT(11) NOT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	Content   string         `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
}

type Student2Course struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
	StudentID int            `json:"student_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"` // 学生ID
	Student   Student        `gorm:"ForeignKey:StudentID;"`
	Course    Course         `gorm:"ForeignKey:CourseID;"`
	CourseID  int            `json:"course_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"` //课程ID
	RecordID  int            `json:"record_id" sql:"type:INT(11) NOT NULL"`                     // 课程列表号
}

type IndependentPractice struct {
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `sql:"index"`
	StudentID      int            `json:"student_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"` // 学生ID
	Student        Student        `gorm:"ForeignKey:StudentID;"`                                      // 学生外键
	Phone          int            `json:"phone"`                                                      // 在外联系方式
	StartTime      time.Time      `json:"start_time"`                                                 // 实习开始日期
	EndTime        time.Time      `json:"end_time"`                                                   // 实习结束日期
	Reason         string         `json:"reason" gorm:"varchar(60)"`                                  // 申请理由
	CompanyName    string         `json:"company_name" gorm:"varchar(60)"`                            // 实习单位名称
	CompanyAddress string         `json:"company_address" gorm:"varchar(60)"`                         // 实习单位地址
	CompanyPerson  string         `json:"company_person" gorm:"varchar(60)"`                          // 实习单位联系人
	CompanyPhone   int            `json:"company_phone" gorm:"int(13)"`                               // 实习单位电话
	Address        string         `json:"address" gorm:"varchar(60)"`                                 // 在外住宿地址
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
	RecordID         int    `json:"record_id" gorm:"primary_key"` // 记录id
}

type HomeworkForSelect struct {
	ID               int `json:"id" gorm:"primary_key;" sql:"type:INT(11) NOT NULL"`
	QuestionID       int `json:"question_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	QuestionMaxScore int `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time
	CourseID         int    `json:"course_id" gorm:"primary_Key;uniqueIndex:question_id;" sql:"type:INT(11) NOT NULL"`
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
	RecordID         int    `json:"record_id" gorm:"primary_key"` // 记录id
}

type CourseForSelect struct {
	RecordID     int    `json:"record_id" gorm:"primary_key"` // 记录id
	ID           int    `json:"course_id" gorm:"primary_key;"`
	CourseName   string `json:"course_name" gorm:"type:varchar(60);not null;"`             // 课程名
	Credit       string `json:"credit" gorm:"type:varchar(50);"`                           // 学分
	TeacherID    int    `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"` // 教师ID
	ClassroomID  int    `json:"classroom_id" gorm:"type:int(11);index:classroom_id;"`      // 教室号
	MaxChooseNum int    `json:"max_choose_num" gorm:"type:int(11);not null;"`              // 最大可选课人数
	SelectedNum  int    `json:"selected_num" gorm:"type:int(11);default:0;"`               // 已选人数
	StartTime    int    `json:"start_time" gorm:"type:datetime;"`                          // 开始时间
	EndTime      int    `json:"end_time" gorm:"type:datetime;"`                            // 结束时间
	Name         string `json:"name" gorm:"size:50;not null"`                              // 教师姓名
	WeekTime     int    `json:"week_time" gorm:"type:int(2)"`                              // 上课星期数
	FacultyID    int    `json:"faculty_id" gorm:"type:int(11)"`                            // 学院ID
	StartWeek    int    `json:"start_week" gorm:"type:int(11);"`                           // 起始周数
	EndWeek      int    `json:"end_week" gorm:"type:int(11);"`                             // 结束周数
	DirectionID  int    `json:"direction_id" gorm:"index:direction_id;type:int(11)"`       // 学生所属方向
	SpecialtyID  int    `json:"specialty_id" gorm:"index:special_id;type:int(11)"`         // 学生所属专业
	CopyFlag     int    `json:"copy_flag"`
	Selectable   bool   `json:"selectable"` //区分毕业设计与普通课程的标记
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
	RecordID         int       `json:"record_id" gorm:"primary_key"` // 记录id
}

type QuestionForCheck struct {
	QuestionMaxScore int `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
	Content          string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
	QuestionId       int    `json:"question_id"`
}

type QuestionForStudent struct {
	RecordID         int `json:"record_id" gorm:"primary_key"` // 记录id
	QuestionMaxScore int `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	DeadLine         time.Time
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
	Content          string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
	QuestionId       int    `json:"question_id"`
	Uploaded         bool   `json:"uploaded"`
}

type RecordForReview struct {
	QuestionID       int    `json:"question_id" gorm:"primary_Key;uniqueIndex:homework_id;" sql:"type:INT(11) NOT NULL"` // 课程号
	Score            int    `json:"score" gorm:"uniqueIndex:score;" sql:"type:INT(11) NOT NULL"`                         // 课程号
	IsScored         bool   `json:"is_scored" gorm:"uniqueIndex:is_scored;"`                                             // 课程号
	Content          string `json:"content" gorm:"primary_Key;uniqueIndex:content;"`
	QuestionMaxScore int    `json:"question_max_score" gorm:"uniqueIndex:question_max_score;" sql:"type:INT(11) NOT NULL"`
	HomeworkTitle    string `json:"homework_title" gorm:"uniqueIndex:homework_title"`
}

type TeacherForSelect struct {
	ID        int    `json:"value"`
	Name      string `json:"name" gorm:"size:50;not null"`  // 教师姓名
	FacultyID int    `json:"faculty_id" gorm:"index:f_id;"` // 学院ID
	TitleID   int    `json:"title_id" gorm:"index:t_id"`    // 职称ID
}

type DirectionToSpecialtyForSelect struct {
	ID            int    `json:"value" gorm:"primary_key;type:int(11)"`
	SpecialtyName string `json:"specialty_name" gorm:"type:varchar(60)"`                          // 专业名称
	DirectionID   int    `json:"direction_id" gorm:"primary_key;index:direction_id;type:int(11)"` // 学生所属方向
	SpecialtyID   int    `json:"specialty_id" gorm:"primary_key;index:special_id;type:int(11)"`   // 学生所属专业
	DirectionName string `json:"direction_name" gorm:"type:varchar(60)"`                          // 方向名称
	FacultyID     int    `json:"faculty_id" gorm:"primary_key;type:int(11)"`                      // 学院ID
}

type ClassroomForSelect struct {
	ID   int    `json:"value" gorm:"primary_key;"`
	Name string `json:"name" gorm:"size:30;not null;uniqueIndex:name;"` // 教室名
}

type Apply struct {
	ID                int    `json:"id"`
	CourseName        string `json:"course_name" gorm:"type:varchar(60);not null;"`                      // 课程名
	Name              string `json:"name" gorm:"size:50;not null"`                                       // 教师姓名
	Reason            string `json:"reason"`                                                             // 申请更改的理由
	BeforeClassroomID int    `json:"before_classroom_id" gorm:"type:int(11);index:before_classroom_id;"` // 教室号
	AfterClassroomID  int    `json:"after_classroom_id" gorm:"type:int(11);index:after_classroom_id;"`   // 教室号
	BeforeStartTime   int    `json:"before_start_time" gorm:"type:int(11);"`                             // 开始节数
	AfterStartTime    int    `json:"after_start_time" gorm:"type:int(11);"`                              // 开始节数
	BeforeEndTime     int    `json:"before_end_time" gorm:"type:int(11);"`                               // 结束节数
	AfterEndTime      int    `json:"after_end_time" gorm:"type:int(11);"`                                // 结束节数
	BeforeWeekTime    int    `json:"before_week_time" gorm:"type:int(2)"`                                // 上课星期数
	AfterWeekTime     int    `json:"after_week_time" gorm:"type:int(2)"`                                 // 上课星期数
	BeforeStartWeek   int    `json:"before_start_week" gorm:"type:int(11);"`                             // 起始周数
	AfterStartWeek    int    `json:"after_start_week" gorm:"type:int(11);"`                              // 起始周数
	BeforeEndWeek     int    `json:"before_end_week" gorm:"type:int(11);"`                               // 结束周数
	AfterEndWeek      int    `json:"after_end_week" gorm:"type:int(11);"`                                // 结束周数
	Result            int    `json:"result" gorm:"type:int(2);default:0"`                                // 审批结果
}

type ClassSheet struct {
	CourseName    string `json:"course_name" gorm:"type:varchar(60);not null;"`            // 课程名
	TeacherName   string `json:"teacher_name" gorm:"size:30;not null;uniqueIndex:name;"`   // 教室名
	ClassroomName string `json:"classroom_name" gorm:"size:30;not null;uniqueIndex:name;"` // 教室名
	StartTime     int    `json:"start_time" gorm:"type:datetime;"`                         // 开始时间
	EndTime       int    `json:"end_time" gorm:"type:datetime;"`                           // 结束时间
	Name          string `json:"name" gorm:"size:50;not null"`                             // 教师姓名
	WeekTime      int    `json:"week_time" gorm:"type:int(2)"`                             // 上课星期数
	StartWeek     int    `json:"start_week" gorm:"type:int(11);"`                          // 起始周数
	EndWeek       int    `json:"end_week" gorm:"type:int(11);"`                            // 结束周数
}

type CourseForChoose struct {
	RecordID     int    `json:"record_id" gorm:"primary_key"` // 记录id
	ID           int    `json:"course_id" gorm:"primary_key;"`
	CourseName   string `json:"course_name" gorm:"type:varchar(60);not null;"`             // 课程名
	Credit       string `json:"credit" gorm:"type:varchar(50);"`                           // 学分
	TeacherID    int    `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"` // 教师ID
	TeacherName  string `json:"teacher_name"`
	ClassroomID  int    `json:"classroom_id" gorm:"type:int(11);index:classroom_id;"` // 教室号
	MaxChooseNum int    `json:"max_choose_num" gorm:"type:int(11);not null;"`         // 最大可选课人数
	Selected     int    `json:"selected" gorm:"column:selected"`                      // 已选人数
	StartTime    int    `json:"start_time" gorm:"type:datetime;"`                     // 开始时间
	EndTime      int    `json:"end_time" gorm:"type:datetime;"`                       // 结束时间
	Name         string `json:"name" gorm:"size:50;not null"`                         // 教师姓名
	WeekTime     int    `json:"week_time" gorm:"type:int(2)"`                         // 上课星期数
	FacultyID    int    `json:"faculty_id" gorm:"type:int(11)"`                       // 学院ID
	StartWeek    int    `json:"start_week" gorm:"type:int(11);"`                      // 起始周数
	EndWeek      int    `json:"end_week" gorm:"type:int(11);"`                        // 结束周数
	DirectionID  int    `json:"direction_id" gorm:"index:direction_id;type:int(11)"`  // 学生所属方向
	SpecialtyID  int    `json:"specialty_id" gorm:"index:special_id;type:int(11)"`    // 学生所属专业
	CopyFlag     int    `json:"copy_flag"`
}

type Student2CourseForChoose struct {
	StudentID int `json:"student_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"` // 学生ID
	CourseID  int `json:"course_id" gorm:"primary_Key;" sql:"type:INT(11) NOT NULL"`  //课程ID
	RecordID  int `json:"record_id" sql:"type:INT(11) NOT NULL"`                      // 课程列表号
}

type SelectedCourseID struct {
	ID int
}

type StudentHomeworkScore struct {
	Name          string `json:"student_name"`
	StudentID     int    `json:"student_id"`
	HomeworkScore int    `json:"homework_score"`
}

type ElectiveForSelect struct {
	StudentID          int    `json:"student_id" gorm:"primary_Key:student_id;" sql:"type:INT(11) NOT NULL"` // 学号
	CourseID           int    `json:"course_id" gorm:"primary_Key:course_id;" sql:"type:INT(11) NOT NULL"`   // 课程号
	CourseName         string `json:"course_name" gorm:"type:varchar(60);not null;"`                         // 课程名
	HomeworkScore      int    `json:"homework_score" gorm:"size:3"`                                          // 作业成绩
	TestScore          int    `json:"test_score" gorm:"size:3;"`                                             // 考试成绩
	BehaviorScore      int    `json:"behavior_score" gorm:"size:3;"`                                         // 平时成绩
	Percentage         int    `json:"percentage" gorm:"size:3"`                                              // 考试成绩所占比例
	HomeworkPercentage int    `json:"homework_percentage" gorm:"size:3"`                                     // 作业分数占比
}

type TeacherForUserInfo struct {
	MyModel
	Avatar      string `json:"avatar" gorm:"size:50"`         // 头像
	Name        string `json:"name" gorm:"size:50;not null"`  // 教师姓名
	Email       string `json:"email" gorm:"size:50;"`         // 邮箱
	Phone       string `json:"phone" gorm:"size:50;"`         // 电话
	Password    string `json:"password" gorm:"size:50"`       // 密码
	FacultyID   int    `json:"faculty_id" gorm:"index:f_id;"` // 学院ID
	FacultyName string `json:"faculty_name"`                  //学院名称
	Title       Title  `gorm:"foreignKey:TitleID;"`           // 职称外键
	TitleID     int    `json:"title_id" gorm:"index:t_id"`    // 职称ID
}

type StudentForUserInfo struct {
	Created       time.Time `json:"created"`
	Name          string    `json:"name" gorm:"type:varchar(50);not null;"`                    // 学生姓名
	FacultyID     int       `json:"faculty_id" gorm:"index:faculty_id;"`                       // 所在院系
	FacultyName   string    `json:"faculty_name"`                                              //学院名称
	Credit        float32   `json:"mark" gorm:"type:float(5,2);"`                              // 累计学分
	DirectionID   int       `json:"direction_id" gorm:"index:direction_id;type:int(11)"`       // 学生所属方向
	SpecialtyID   int       `json:"specialty_id" gorm:"index:special_id;type:int(11)"`         // 学生所属专业
	SpecialtyName string    `json:"specialty_name"`                                            //专业名称
	TeacherID     int       `json:"teacher_id" gorm:"type:int(11);not null;index:teacher_id;"` // 导师ID
	Practice      int       `json:"practice"`                                                  //实习方式
}

type CourseForChooses []CourseForChoose
type HomeworkUploadRecordsForSelects []HomeworkUploadRecordsForSelect
type HomeworkUploadRecords []HomeworkUploadRecord
type HomeworkForSelects []HomeworkForSelect
type CourseForSelects []CourseForSelect
type QuestionForSelects []QuestionForSelect
type Homeworks []Homework
type QuestionForChecks []QuestionForCheck
type QuestionForStudents []QuestionForStudent
type RecordForReviews []RecordForReview
type TeacherForSelects []TeacherForSelect
type Faculties []Faculty
type DirectionToSpecialtyForSelects []DirectionToSpecialtyForSelect
type ClassroomForSelects []ClassroomForSelect
type Applies []Apply
type ClassSheets []ClassSheet
type Student2CourseForChooses []Student2CourseForChoose
type SelectedCourseIDs []SelectedCourseID
type StudentHomeworkScores []StudentHomeworkScore
type Electives []Elective
type ElectiveForSelects []ElectiveForSelect

func CreateDatabase(db *gorm.DB) {
	//db.AutoMigrate(&Elective{})
	db.AutoMigrate(&Title{}, &Faculty{}, &Teacher{}, &Elective{}, &Admin{}, &Classroom{}, &DirectionToSpecialty{}, &Course{}, &Student{}, &HomeworkUploadRecord{}, &Homework{}, &Question{}, &Student2Course{}, &ApplyForCourseChange{}, &IndependentPractice{})
}
