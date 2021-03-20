<template>
  <div>
    <v-card>
      <v-card-title>选择课程</v-card-title>
      <div id="coursesTable" class="text-center font-weight-bold"></div>
      <v-dialog v-model="classInfo" width="500">
        <v-card>
          <v-card-title class="font-weight-bold py">课程详情</v-card-title>
          <v-card-text class="font-weight-bold mx-auto" >
            <div v-html="msg" class="mx-auto text-center"></div>
          </v-card-text>
        </v-card>
      </v-dialog>
    </v-card>
  </div>

</template>

<script>


import Timetables from "timetables";

export default {
  name: "ClassSheet",
  mounted() {
    this.$axios({
      method:"GET",
      url:"http://127.0.0.1:9000/getClassSheet?teacher_id="+2+"&week="+5
    }).then((response)=>{
      this.timetables = response.data.classSheet
      this.timetable = new Timetables({
        el: '#coursesTable',
        timetables: this.timetables,
        week: this.week,
        timetableType: this.timetableType,
        highlightWeek: this.highlightWeek,
        gridOnClick:  (e)=> {
          this.classInfo = true
          var reg = new RegExp("\n", "g")
          e.name = e.name.replace(reg, "<br/>")
          console.log(e.name)
          this.msg = e.name + '<br/>' + e.week + '<br/> 第' + e.index + '节课<br/> 课长' + e.length + '节'
        },
        styles: this.styles
      });
    })
  },
  data() {
    return {
      classInfo:false,
      msg:null,
      timetable: null,
      timetables: [],
      timetableType: [
        [{index: '1', name: '8:50'}, 1],
        [{index: '2', name: '9:45'}, 1],
        [{index: '3', name: '10:45'}, 1],
        [{index: '4', name: '11:40'}, 1],
        [{index: '5', name: '午餐'}, 1],
        [{index: '6', name: '午休'}, 1],
        [{index: '7', name: '13:50'}, 1],
        [{index: '8', name: '14:45'}, 1],
        [{index: '9', name: '15:45'}, 1],
        [{index: '10', name: '16:40'}, 1],
        [{index: '11', name: '18:30'}, 1],
        [{index: '12', name: '19:25'}, 1],
        [{index: '13', name: '20:20'}, 1],
        [{index: '14', name: '21:15'}, 1],
      ],
      week: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      highlightWeek: new Date().getDay(),
      styles: {
        Gheight: 50,
        leftHandWidth: 50,
      },
    }
  },
}
</script>

<style scoped>
</style>
