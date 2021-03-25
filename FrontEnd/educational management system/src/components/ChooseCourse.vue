<template>
  <v-card>
    <v-card-title>选择课程</v-card-title>
    <v-data-table
      :headers="headers"
      :items="courses"
      :items-per-page="5"
      class="elevation-1"
    >
      <template v-slot:item.operation="{ item }">
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="changeClass = true ;selectId = item.course_id;getInfo()">
              <v-icon>
                mdi-lead-pencil
              </v-icon>
            </v-btn>
          </template>
          <span>申请修改课程</span>
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
      headers: [
        {text: '课程名', align: 'start', sortable: false, value: 'course_name'},
        {text: '教师', sortable: false, value: 'name'},
        {text: '学分', sortable: false, value: 'credit'},
        {text: '已选人数', sortable: false, value: 'selected_num'},
        //学年由上传年份决定
        {text: '开始周', sortable: false, value: 'start_time'},
        {text: '结束周', sortable: false, value: 'end_time'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      courses:[{}],
    }
  },
  created() {
    this.$axios({
      url: "http://127.0.0.1:9000/getAvailableCourses",
      method: "get",
      params: {
        'student_id': this.$store.state.studentId,
        'faculty_id': this.$store.state.facultyId,
        'direction_id': this.$store.state.directionId,
        'specialty_id': this.$store.state.specialtyId,
      }
    }).then((response)=>{
      this.courses = response.data.courses
      console.log(response)
    })
  }
}
</script>

<style scoped>

</style>
