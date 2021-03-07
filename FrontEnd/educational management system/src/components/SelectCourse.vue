<template>
<v-card>
  <v-card-title>选择课程</v-card-title>
  <v-data-table
    :headers="headers"
    :items="courses"
    :items-per-page="5"
    class="elevation-1"
  >
    <template v-slot:item.course_name="{ item }">
        <router-link to="/SelectHomework" @click.native="$store.commit('setCourseId',item.course_id)">{{ item.course_name }}</router-link>
    </template>
  </v-data-table>
</v-card>
</template>

<script>
export default {
  name: "SelectCourse",
  data(){
    return {
      course:null,
      courses:[],
      headers:[
        {text:'课程名',align:'start',sortable:false,value:'course_name'},
        {text:'教师',sortable:false,value:'name'},
        //学年由上传年份决定
        {text:'开始周',sortable:false,value:'start_time'},
        {text:'结束周',sortable:false,value:'end_time'},
      ],
    }
  },
  created() {
    this.$axios({
      method:"get",
      url:'http://127.0.0.1:9000/getCourseList',
      params:{
        'teacher_id':this.$store.state.teacherId
      }
    }).then((response)=>{
      this.courses = response.data.courses
      console.log(this.courses)
    })
  },
  watch:{
    '$store.state.refreshFlag':{
      handler(newValue, oldValue) {
        console.log(newValue,oldValue)
        if (newValue !== 1 ) return
        this.$axios({
          method:"get",
          url:'http://127.0.0.1:9000/getCourseList',
          params:{
            'teacher_id':this.$store.state.teacherId
          }
        }).then((response)=>{
          this.courses = response.data.courses
          console.log(this.courses)
        })
      },
    },
  }
}
</script>

<style scoped>

</style>
