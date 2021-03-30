<template>
  <v-card>
    <v-card-title>个人信息</v-card-title>
    <v-row class="mt-4 mb-4">
      <v-col cols="1">
      </v-col>
      <v-col cols="3">
        <v-avatar color="indigo" size="160" class="float-right">
          <v-icon dark size="120">
            mdi-account-circle
          </v-icon>
        </v-avatar>
      </v-col>
      <v-col cols="2">
      </v-col>
      <v-col cols="5">
        <v-row>
          <v-card-title class="font-weight-bold">姓名：{{ name }}</v-card-title>
        </v-row>
        <v-row>
          <v-card-title class="font-weight-bold">学院：{{ faculty }}</v-card-title>
        </v-row>
        <v-row v-if="specialty_flag">
          <v-card-title class="font-weight-bold">专业：{{ specialty }}</v-card-title>
        </v-row>
      </v-col>
    </v-row>
    <v-card class="d-flex justify-center" v-if="specialty_flag">
      <v-col cols="4">
        <v-card>
          <v-card-title>选择方向</v-card-title>
          <v-select outlined :items="directions" item-text="direction_name" item-value="direction_id"
                    v-model="direction_id"></v-select>
          <v-btn block x-large class="ma-0" @click="postDirection">提交</v-btn>
          <v-overlay :absolute="true" :value="direction_invalid">
            <v-card-title class="font-weight-bold">在第四学期6月开放</v-card-title>
          </v-overlay>
        </v-card>
      </v-col>
      <v-col cols="4">
        <v-card>
          <v-card-title>选择导师</v-card-title>
          <v-select outlined :items="teachers" item-text="name" item-value="value" v-model="teacher_id"></v-select>
          <v-btn block x-large class="ma-0" @click="postTeacher">提交</v-btn>
          <v-overlay :absolute="true" :value="teacher_invalid">
            <v-card-title class="font-weight-bold">在第七学期9月开放</v-card-title>
          </v-overlay>

        </v-card>
      </v-col>
      <v-col cols="4">
        <v-card>
          <v-card-title>选择实习方式</v-card-title>
          <v-select outlined :items="practices" v-model="practice" item-value="value" item-text="name"></v-select>
          <v-btn block x-large class="ma-0" :disabled="practice === 1" @click="postPractice">提交</v-btn>
          <v-overlay :absolute="true" :value="practice_invalid">
            <v-card-title class="font-weight-bold">在第七学期9~12月开放</v-card-title>

          </v-overlay>
        </v-card>
      </v-col>
    </v-card>

    <v-card v-if="practice === 1 " class="d-flex flex-wrap">
      <v-col cols="12">
        <v-row>
          <v-card-title>自主实习申请表</v-card-title>
          <v-spacer></v-spacer>
          <v-card-title>
            <v-btn color="primary" @click="postPractice">提交</v-btn>
          </v-card-title>
        </v-row>
      </v-col>
      <v-col cols="2">
        <v-text-field label="在外联系方式" outlined v-model="phone"></v-text-field>
      </v-col>
      <v-col cols="5">
        <v-menu
          ref="menu"
          v-model="menu"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
          min-width="auto">
          <template v-slot:activator="{ on, attrs }">
            <v-text-field
              outlined
              v-model="independent_practice.start_time"
              label="开始实习日期"
              v-bind="attrs"
              v-on="on"
              :disabled="$store.state.modifyHomeworkFlag"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="independent_practice.start_time"
            @input="menu1 = false"
            color="green lighten-1"
          ></v-date-picker>
        </v-menu>
      </v-col>
      <v-col cols="5">
        <v-menu
          ref="menu"
          v-model="independent_practice.menu"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
          min-width="auto">
          <template v-slot:activator="{ on, attrs }">
            <v-text-field
              outlined
              v-model="independent_practice.end_time"
              label="结束实习日期"
              v-bind="attrs"
              v-on="on"
              :disabled="$store.state.modifyHomeworkFlag"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="independent_practice.end_time"
            @input="menu1 = false"
            color="green lighten-1"
          ></v-date-picker>
        </v-menu>
      </v-col>
      <v-col cols="12">
        <v-textarea label="申请理由" outlined v-model="independent_practice.reason"></v-textarea>
      </v-col>
      <v-col cols="6">
        <v-text-field label="实习单位名称" outlined v-model="independent_practice.company_name"></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field label="实习单位地址" outlined v-model="independent_practice.company_address"></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field label="实习单位联系人" outlined v-model="independent_practice.company_person"></v-text-field>
      </v-col>
      <v-col cols="6">
        <v-text-field label="实习单位电话" outlined v-model="independent_practice.company_phone"></v-text-field>
      </v-col>
      <v-col cols="12">
        <v-text-field label="在外住宿地址" outlined v-model="independent_practice.address"></v-text-field>
      </v-col>

    </v-card>
  </v-card>
</template>

<script>
export default {
  name: "UserInfo",
  data() {
    return {
      direction_id: null,
      teacher_id: null,
      directions: [],
      teachers: [],
      practice: 0,
      practices: [{value: 0, name: '统一实习'}, {value: 1, name: '自主实习'}],
      direction_invalid: false,
      teacher_invalid: false,
      practice_invalid: false,
      specialty_flag: true,
      student_year: null,
      teacher: {},
      student: {},
      name: null,
      faculty: null,
      specialty: null,
      independent_practice:{}
    }
  },
  created() {
    this.$axios({
      method: 'get',
      url: "http://127.0.0.1:9000/getUserInfo",
      params: {
        student_id: this.$store.state.studentId,
        teacher_id: this.$store.state.teacherId,
      }
    }).then((response) => {
      console.log(response)
      if (response.data.teacher.name === "") {
        //学生查询
        this.student = response.data.student
        this.name = response.data.student.name
        this.faculty = response.data.student.faculty_name
        this.specialty = response.data.student.specialty_name
        this.student_year = parseInt(response.data.student.created.substr(0, 4))
        let d = new Date()
        // d.setMonth(9)
        // d.setFullYear(2020)
        if (d.getFullYear() === this.student_year + 2 && d.getMonth() + 1 === 6) {
          this.getDirectionList();
          this.direction_invalid = false
        }
        if (d.getFullYear() === this.student_year + 3 && d.getMonth() + 1 === 9) {
          this.getTeacherList();
          this.teacher_invalid = false
        }
        if ((d.getFullYear() === this.student_year + 3 && d.getMonth() + 1 >= 9 && d.getMonth() + 1 <= 12) || (d.getFullYear() === this.student_year + 4 && d.getMonth() + 1 === 1)) this.practice_invalid = false
        if (response.data.student.direction_id !== 0) {
          this.getDirectionList();
          this.direction_id = response.data.student.direction_id
        }
        if (response.data.student.teacher_id !== 0) {
          this.getTeacherList();
          this.teacher_id = response.data.student.teacher_id
        }
        if (response.data.student.practice !== 0) {
          this.practice = parseInt(response.data.student.practice)
          console.log(this.practice)
        }
      } else {
        //老师查询
        this.specialty_flag = false
        this.teacher = response.data.teacher
        this.name = response.data.teacher.name
        this.faculty = response.data.teacher.faculty_name
      }
    })
  },
  methods: {
    getDirectionList() {
      this.$axios({
        url: "http://127.0.0.1:9000/getDirectionList",
        method: "GET",
        params: {
          faculty_id: this.student.faculty_id,
          specialty_id: this.student.specialty_id,
        }
      }).then((response) => {
        this.directions = response.data.directions
      })
    },
    getTeacherList() {
      this.$axios({
        url: "http://127.0.0.1:9000/getTeacherList",
        method: "GET",
        params: {
          faculty_id: this.student.faculty_id,
        }
      }).then((response) => {
        this.teachers = response.data.teachers
      })
    },
    postDirection() {
      const formData = new FormData()
      formData.append("direction_id", this.direction_id)
      formData.append("student_id", this.$store.state.studentId)
      this.$axios({
        url: "http://127.0.0.1:9000/postDirection",
        method: "post",
        data: formData,
        headers:{
          "Content-Type": "multipart/form-data"
        }
      }).then((response)=>{
        this.$store.commit(response.data.snackbar,response.data.msg)
      })
    },
    postTeacher(){
      const formData = new FormData()
      formData.append("teacher_id", this.teacher_id)
      formData.append("student_id", this.$store.state.studentId)
      this.$axios({
        url: "http://127.0.0.1:9000/postTeacher",
        method: "post",
        data: formData,
        headers:{
          "Content-Type": "multipart/form-data"
        }
      }).then((response)=>{
        this.$store.commit(response.data.snackbar,response.data.msg)
      })
    },
    postPractice(){
      const formData = new FormData()
      formData.append("practice", this.practice)
      formData.append("student_id", this.$store.state.studentId)
      this.$axios({
        url:"http://127.0.0.1:9000/postPractice",
        method:"post",
        data:formData,
        headers:{
          "Content-Type": "multipart/form-data"
        }
      }).then((response)=>{
        this.$store.commit(response.data.snackbar,response.data.msg)
      })
      if(this.practice === 1){
        const formData1 = new FormData()
        formData1.append("independent_practice", this.independent_practice)
        formData1.append("student_id", this.$store.state.studentId)

        this.$axios({
          url:"http://127.0.0.1:9000/postIndependentPractice",
          method:"post",
          data:formData1,
          headers:{
            "Content-Type": "multipart/form-data"
          }
        })
      }
    }
  },
}
</script>

<style scoped>

</style>
