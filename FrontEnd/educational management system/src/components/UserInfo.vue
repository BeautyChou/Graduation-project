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
          <v-card-title class="font-weight-bold">姓名：{{ $store.state.username }}</v-card-title>
        </v-row>
        <v-row>
          <v-card-title class="font-weight-bold" v-if="this.$store.state.level === 1 || this.$store.state.level === 2">学院：{{ faculty }}</v-card-title>
        </v-row>
        <v-row>
          <v-card-title class="font-weight-bold" v-if="this.$store.state.level === 3">身份：管理员</v-card-title>
        </v-row>
        <v-row v-if="this.$store.state.level === 1">
          <v-card-title class="font-weight-bold">专业：{{ specialty }}</v-card-title>
        </v-row>
      </v-col>
    </v-row>
    <v-card class="d-flex justify-center" v-if="this.$store.state.level === 1">
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
          <v-btn block x-large class="ma-0" @click="postTeacher" :disabled="(student.teacher_id!==0)">提交</v-btn>
          <v-overlay :absolute="true" :value="teacher_invalid||teacher_flag">
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
    <v-card v-if="this.$store.state.level === 2">
      <v-card-title>毕业设计学生申请</v-card-title>
      <v-simple-table>
        <template v-slot:default>
          <thead>
          <tr>
            <th class="text-left">
              学生姓名
            </th>
            <th class="text-left">
              操作
            </th>
          </tr>
          </thead>
          <tbody>
          <tr
            v-for="item in students"
            :key="item.student_id"
          >
            <td>{{ item.name }}</td>
            <td>
              <v-tooltip v-if="$store.state.level===2||true" bottom>
                <template v-slot:activator="{ on,attrs }">
                  <v-btn
                    icon
                    color="green"
                    v-bind="attrs"
                    v-on="on"
                    @click.native="pass(item)"
                    small
                    class="col-1">
                    <v-icon>
                      mdi-check-bold
                    </v-icon>
                  </v-btn>
                </template>
                <span>同意</span>
              </v-tooltip>
              <v-tooltip v-if="$store.state.level===2||true" bottom>
                <template v-slot:activator="{ on,attrs }">
                  <v-btn
                    icon
                    color="red"
                    v-bind="attrs"
                    v-on="on"
                    @click.native="denied(item)"
                    small
                    class="col-1">
                    <v-icon>
                      mdi-close-thick
                    </v-icon>
                  </v-btn>
                </template>
                <span>不同意</span>
              </v-tooltip>
            </td>
          </tr>
          </tbody>
        </template>
      </v-simple-table>

    </v-card>
    <v-card v-if="practice === 1 " class="d-flex flex-wrap">
      <v-overlay :absolute="true" :value="teacher_invalid">
        <v-card-title class="font-weight-bold">在第七学期9月开放</v-card-title>
      </v-overlay>
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
        <v-text-field label="在外联系方式" outlined v-model.number="independent_practice.phone"></v-text-field>
      </v-col>
      <v-col cols="5">
        <v-menu
          ref="menu"
          v-model="menu1"
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
          v-model="menu2"
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
            @input="menu2 = false"
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
        <v-text-field label="实习单位电话" outlined v-model.number="independent_practice.company_phone"></v-text-field>
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
      menu1: false,
      menu2: false,
      direction_id: null,
      teacher_id: null,
      directions: [],
      teachers: [],
      practice: 0,
      practices: [{value: 0, name: '统一实习'}, {value: 1, name: '自主实习'}],
      direction_invalid: true,
      teacher_invalid: true,
      practice_invalid: true,
      student_year: null,
      teacher: {},
      student: {},
      name: null,
      faculty: null,
      specialty: null,
      independent_practice: {},
      teacher_flag: null,
      current_time: null,
      students: [],
    }
  },
  created() {
    this.$axios({
      method: 'get',
      url: "UserInfo/UserInfo",
      params: {
        student_id: this.$store.state.studentId,
        teacher_id: this.$store.state.teacherId,
      },
      headers: {
        'Token': "8a54sh " + this.$store.state.Jwt
      }
    }).then((response) => {
      console.log(response)
      if (response.data.msg === "Token无效") {
        this.$emit('func')
        return
      }
      if (response.data.teacher.name === "") {
        //学生查询
        this.student = response.data.student
        this.name = response.data.student.name
        this.faculty = response.data.student.faculty_name
        this.specialty = response.data.student.specialty_name
        this.student_year = parseInt(response.data.student.created.substr(0, 4))
        this.teacher_flag = response.data.student.teacher_flag
        this.current_time = response.data.current_time
        let year = this.current_time.substr(0, 4)
        let month = this.current_time.substr(5, 2)
        this.getDirectionList();
        this.getTeacherList();
        if (year === (this.student_year + 2).toString() && month === "06") {
          this.direction_invalid = false
        }
        if (year === (this.student_year + 3).toString() && month === "09") {
          this.teacher_invalid = false
        }
        if ((year === (this.student_year + 3).toString() && month >= "09" && month <= "12") || (year === (this.student_year + 4).toString() && month === "01")) {
          this.practice_invalid = false
        }
        this.direction_id = response.data.student.direction_id
        this.teacher_id = response.data.student.teacher_id
        this.practice = parseInt(response.data.student.practice)
      } else {
        //老师查询
        this.teacher = response.data.teacher
        this.name = response.data.teacher.name
        this.faculty = response.data.teacher.faculty_name
        this.students = response.data.students

      }
    })
  },
  methods: {
    getDirectionList() {
      this.$axios({
        url: "UserInfo/Direction",
        method: "GET",
        params: {
          faculty_id: this.student.faculty_id,
          specialty_id: this.student.specialty_id,
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.directions = response.data.directions
      })
    },
    getTeacherList() {
      this.$axios({
        url: "UserInfo/Teacher",
        method: "GET",
        params: {
          faculty_id: this.student.faculty_id,
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.teachers = response.data.teachers
      })
    },
    postDirection() {
      const formData = new FormData()
      formData.append("direction_id", this.direction_id)
      formData.append("student_id", this.$store.state.studentId)
      this.$axios({
        url: "UserInfo/Direction",
        method: "post",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    postTeacher() {
      const formData = new FormData()
      formData.append("teacher_id", this.teacher_id)
      formData.append("student_id", this.$store.state.studentId)
      this.$axios({
        url: "UserInfo/Teacher",
        method: "post",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
        this.student.teacher_id = this.teacher_id
      })
    },
    postPractice() {
      const formData = new FormData()
      formData.append("practice", this.practice)
      formData.append("student_id", this.$store.state.studentId)
      this.$axios({
        url: "Practice",
        method: "put",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
      if (this.practice === 1) {
        this.independent_practice.student_id = this.$store.state.studentId
        this.independent_practice.start_time = this.independent_practice.start_time + 'T00:00:00+08:00'
        this.independent_practice.end_time = this.independent_practice.end_time + 'T00:00:00+08:00'
        var p = JSON.stringify(this.independent_practice)
        this.$axios({
          url: "Practice",
          method: "post",
          data: p,
          headers: {
            "Content-Type": "multipart/form-data",
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response) => {
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.$store.commit(response.data.snackbar, response.data.msg)
          setTimeout(() => {
            this.$store.commit(response.data.snackbar2)
          }, 3000)
          this.independent_practice = {}
        })
      }
    },
    pass(studentOBJ) {
      const formData = new FormData()
      formData.append("student_id", studentOBJ.student_id)
      formData.append("teacher_id", this.$store.state.teacherId)

      formData.append("operation", "pass")
      this.$axios({
        url: "UserInfo/ApplyTeacher",
        method: "post",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.students.some((v, i) => {
          if (v.student_id === studentOBJ.student_id) {
            this.students.splice(i, 1)
          }
        })
      })
    },
    denied(studentOBJ) {
      const formData = new FormData()
      formData.append("student_id", studentOBJ.student_id)
      formData.append("operation", "denied")
      this.$axios({
        url: "UserInfo/ApplyTeacher",
        method: "post",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.students.some((v, i) => {
          if (v.student_id === studentOBJ.student_id) {
            this.students.splice(i, 1)
          }
        })
      })
    },
  },
  watch: {
    "$route.path": {
      handler(newVal, oldVal) {
        this.$axios({
          method: 'get',
          url: "UserInfo/UserInfo",
          params: {
            student_id: this.$store.state.studentId,
            teacher_id: this.$store.state.teacherId,
          },
          headers: {
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response) => {
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          console.log(response)
          if (response.data.teacher.name === "") {
            //学生查询
            this.student = response.data.student
            this.name = response.data.student.name
            this.faculty = response.data.student.faculty_name
            this.specialty = response.data.student.specialty_name
            this.student_year = parseInt(response.data.student.created.substr(0, 4))
            this.teacher_flag = response.data.student.teacher_flag
            this.current_time = response.data.current_time
            let year = this.current_time.substr(0, 4)
            let month = this.current_time.substr(5, 2)
            year = "2020"
            month = "09"
            this.getDirectionList();
            this.getTeacherList();
            if (year === (this.student_year + 2).toString() && month === "06") {
              this.direction_invalid = false
            }
            if (year === (this.student_year + 3).toString() && month === "09") {
              this.teacher_invalid = false
            }
            if ((year === (this.student_year + 3).toString() && month >= "09" && month <= "12") || (year === (this.student_year + 4).toString() && month === "01")) {
              this.practice_invalid = false
            }
            this.direction_id = response.data.student.direction_id
            this.teacher_id = response.data.student.teacher_id
            this.practice = parseInt(response.data.student.practice)
          } else {
            //老师查询
            this.specialty_flag = false
            this.teacher = response.data.teacher
            this.name = response.data.teacher.name
            this.faculty = response.data.teacher.faculty_name
            this.students = response.data.students

          }
        })
      },
    }
  }
}
</script>

<style scoped>

</style>
