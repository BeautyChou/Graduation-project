<template>
  <v-card>
    <v-overlay :value="overlay&&false">
      <v-card-title class="title font-weight-bold">当前时间不在可选课时间内！</v-card-title>
    </v-overlay>
    <v-card-title>课程列表</v-card-title>
    <v-data-table
      :headers="headers"
      :items="courses"
      :options.sync="options"
      :server-items-length="total"
      class="elevation-1"
    >
      <template v-slot:item.operation="{ item }">
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="chooseCourse(item)">
              <v-icon>
                mdi-plus-thick
              </v-icon>
            </v-btn>
          </template>
          <span>申请课程</span>
        </v-tooltip>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
export default {
  name: "ChooseCourse",
  data() {
    return {
      options:{},
      total:null,
      headers: [
        {text: '课程名', align: 'start', sortable: false, value: 'course_name'},
        {text: '教师', sortable: false, value: 'name'},
        {text: '学分', sortable: false, value: 'credit'},
        {text: '最大人数', sortable: false, value: 'max_choose_num'},
        {text: '已选人数', sortable: false, value: 'selected'},
        //学年由上传年份决定
        {text: '开始周', sortable: false, value: 'start_time'},
        {text: '结束周', sortable: false, value: 'end_time'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      courses: [{}],
      overlay:true,
    }
  },
  created() {
    this.getCourse()
  },
  methods: {
    chooseCourse(courseOBJ) {
      const formData = new FormData()
      formData.append("student_id", this.$store.state.studentId)
      formData.append("record_id", courseOBJ.record_id)
      formData.append("course_id", courseOBJ.course_id)
      this.$axios({
        method: "post",
        url: "Course",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        },
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(()=>{
          this.$store.commit(response.data.snackbar2)
        },3000)
        this.$axios({
          url: "Course",
          method: "get",
          params: {
            'student_id': this.$store.state.studentId,
            'faculty_id': this.$store.state.facultyId,
            'direction_id': this.$store.state.directionId,
            'specialty_id': this.$store.state.specialtyId,
          },
          headers:{
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response) => {
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.courses = response.data.courses
          console.log(response)
        })
      })
    },
    getCourse(){
      this.$axios({
        url: "Course",
        method: "get",
        params: {
          'student_id': this.$store.state.studentId,
          'faculty_id': this.$store.state.facultyId,
          'direction_id': this.$store.state.directionId,
          'specialty_id': this.$store.state.specialtyId,
          "page":this.options.page,
          "items":this.options.itemsPerPage,
        },
        headers:{
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.courses = response.data.courses
        this.total = response.data.total
        var time = response.data.time
        var month = time.substr(5,2)
        var day = time.substr(8,2)
        console.log(month,day)
        if ( (month==='07'&&(day>='05'||day<='15'))||(month==='02'&&(day>='05'||day<='15')) ) this.overlay = false
        console.log(response)
      })
    }
  },
  watch: {
    "$route.path": {
      handler(newVal, oldVal) {
        this.getCourse()
      },
    },
    options:{
      handler(){
        this.getCourse()
      },
      deep:true,
    }
  }
}
</script>

<style scoped>

</style>
