<template>
  <v-card>
    <v-card-title>已选课程</v-card-title>
    <v-data-table
      :headers="headers"
      :items="courses"
      :items-per-page="5"
      class="elevation-1"
    >
      <template v-slot:item.operation="{ item }">
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red" v-bind="attrs" v-on="on" x-large
                   @click="quitCourse(item)">
              <v-icon>
                mdi-exit-to-app
              </v-icon>
            </v-btn>
          </template>
          <span>退出课程</span>
        </v-tooltip>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
export default {
name: "ChosenCourse",
  data() {
    return {
      headers: [
        {text: '课程名', align: 'start', sortable: false, value: 'course_name'},
        {text: '教师', sortable: false, value: 'name'},
        {text: '学分', sortable: false, value: 'credit'},
        {text: '最大人数', sortable: false, value: 'max_choose_num'},
        //学年由上传年份决定
        {text: '开始周', sortable: false, value: 'start_time'},
        {text: '结束周', sortable: false, value: 'end_time'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      courses:[],
    }
  },
  created() {
    this.$axios({
      url: "http://127.0.0.1:9000/getChosenCourse",
      method: "get",
      params: {
        'student_id': this.$store.state.studentId,
      },
      headers:{
        'Token': "8a54sh " + this.$store.state.Jwt
      }
    }).then((response)=>{
      if (response.data.msg === "Token无效") {
        this.$emit('func')
        return
      }
      this.courses = response.data.courses
      console.log(response)
    })
  },
  methods:{
    quitCourse(courseOBJ){
      this.$axios({
        method:"delete",
        url:"http://127.0.0.1:9000/quitCourse",
        params:{
          "student_id":this.$store.state.studentId,
          "record_id":courseOBJ.record_id,
          "course_id":courseOBJ.course_id,
        },
        headers:{
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response)=>{
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.$store.commit(response.data.snackbar,response.data.msg)
        setTimeout(()=>{
          this.$store.commit(response.data.snackbar2)
        },3000)
        this.$axios({
          url: "http://127.0.0.1:9000/getChosenCourse",
          method: "get",
          params: {
            'student_id': this.$store.state.studentId,
          },
          headers:{
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response)=>{
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.courses = response.data.courses
          console.log(response)
        })
      })
    },
  },
  watch: {
    "$route.path": {
      handler(newVal, oldVal) {
        this.$axios({
          url: "http://127.0.0.1:9000/getChosenCourse",
          method: "get",
          params: {
            'student_id': this.$store.state.studentId,
          },
          headers:{
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response)=>{
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.courses = response.data.courses
          console.log(response)
        })
      },
    }
  }
}
</script>

<style scoped>

</style>
