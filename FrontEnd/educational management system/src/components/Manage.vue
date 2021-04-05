<template>
  <v-tabs grow>
    <v-tab>学生管理</v-tab>
    <v-tab>教师管理</v-tab>
    <v-tab>管理员管理</v-tab>
    <v-tab>职称管理</v-tab>
    <v-tab>处分等级管理</v-tab>
    <v-tab>学院管理</v-tab>
    <v-tab>教师管理</v-tab>

    <v-tab-item>
      <v-manage-student></v-manage-student>
    </v-tab-item>
    <v-tab-item>
      <v-data-table
        :headers="teacher_headers"
        :items="teachers"
      ></v-data-table>
    </v-tab-item>
    <v-tab-item>
      <v-data-table
        :headers="admin_header"
        :items="admins"
      ></v-data-table>
    </v-tab-item>
    <v-tab-item>
      <v-data-table
        :headers="admin_header"
        :items="admins"
      ></v-data-table>
    </v-tab-item>
    <v-tab-item>
      <v-data-table
        :headers="admin_header"
        :items="admins"
      ></v-data-table>
    </v-tab-item>
    <v-tab-item>
      <v-data-table
        :headers="admin_header"
        :items="admins"
      ></v-data-table>
    </v-tab-item>
    <v-tab-item>
      <v-data-table
        :headers="admin_header"
        :items="admins"
      ></v-data-table>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import ManageStudent from "./ManageStudent";
export default {
  name: "Manage",
  components: {
    "v-manage-student":ManageStudent,
  },
  data() {
    return {
      directions: [],
      new_faculty: null,
      new_specialty: null,
      new_direction: null,
      new_password: null,
      punishment_level: [],
      selectOBJ: {},
      modify_student: false,
      delete_student: false,
      specialty_id: null,
      faculty_id: null,
      faculties: [],
      specialties: [],
      students: [],
      teachers: [],
      admins: [],
      student_headers: [
        {text: '姓名', align: 'start', sortable: false, value: 'name'},
        {text: '学院', sortable: false, value: 'faculty_name'},
        {text: '专业', sortable: false, value: 'specialty_name'},
        {text: '方向', sortable: false, value: 'direction'},
        {text: '实习方式', sortable: false, value: 'practice'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      teacher_headers: [
        {text: '姓名', align: 'start', sortable: false, value: 'student_name'},
        {text: '学院', sortable: false, value: 'faculty'},
        {text: '专业', sortable: false, value: 'specialty'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      admin_header: [
        {text: '姓名', align: 'start', sortable: false, value: 'student_name'},
        {text: '操作', sortable: false, value: 'operation'}
      ]
    }
  },
  created() {
    this.getStudentList()
  },
  methods: {
    getStudentList() {
      this.$axios({
        url: "http://127.0.0.1:9000/getStudentList",
        params: {
          faculty_id: this.faculty_id,
          specialty_id: this.specialty_id,
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get"
      }).then((response) => {
        console.log(response)
        this.specialties = response.data.specialties
        this.faculties = response.data.faculties
        this.students = response.data.students
        this.punishment_level = response.data.punishment_level
        this.directions = response.data.directions
      })
    },
    deleteStudent() {
      this.$axios({
        url: "deleteStudent",
        method: "delete",
        params: {
          student_id: this.selectOBJ.student_id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        this.delete_student = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getStudentList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    modifyStudent(){
      var formData = new FormData()
      formData.append("student_id",this.selectOBJ.student_id)
      formData.append("password",this.new_password)
      formData.append("faculty_id",this.new_faculty)
      formData.append("specialty_id",this.new_specialty)
      formData.append("direction_id",this.new_direction)
      this.$axios({
        method:"put",
        url:"putStudent",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data:formData
      }).then((response)=>{
        this.modify_student = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getStudentList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    }
  },
  watch: {
    "$route.path": {
      handler(newVal, oldVal) {
        this.getStudentList()
      }
    }
  }
}
</script>

<style scoped>

</style>
