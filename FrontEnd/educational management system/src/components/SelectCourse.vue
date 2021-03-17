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
        <router-link to="/SelectHomework" @click.native="$store.commit('setCourseId',item.course_id)">
          {{ item.course_name }}
        </router-link>
      </template>
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
    <v-dialog v-model="changeClass" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          修改科目
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-4">
            <v-select
              v-model="apply.before_week_time"
              :items="week_time"
              item-text="name"
              item-value="value"
              label="上课时间"
              outlined
              class="col-4"></v-select>
            <v-select
              outlined
              v-model="apply.before_start_time"
              :items="start_time"
              item-text="name"
              item-value="value"
              label="起始节次"
              class="col-4"></v-select>
            <v-select
              outlined
              v-model="apply.before_end_time"
              :items="end_time"
              item-text="name"
              item-value="value"
              label="结束节次"
              class="col-4"></v-select>
            <v-select
              v-model="apply.before_classroom_id"
              :items="classrooms"
              item-text="name"
              item-value="value"
              outlined
              label="上课教室"
              class="col-4"></v-select>
            <v-select
              outlined
              v-model="apply.before_start_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="开始周"
              class="col-4"></v-select>
            <v-select
              outlined
              v-model="apply.before_end_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="结束周"
              class="col-4"></v-select>
          </v-row>
          <v-row class="my-0"><span class="mx-auto my-0"><v-icon x-large>mdi-arrow-down-bold</v-icon>修改为</span></v-row>
          <v-row>
            <v-select
              v-model="apply.after_week_time"
              :items="week_time"
              item-text="name"
              item-value="value"
              label="上课时间"
              outlined
              class="col-4"></v-select>
            <v-select
              class="col-4"
              outlined
              v-model="apply.after_start_time"
              :items="start_time"
              item-text="name"
              item-value="value"
              label="起始节次"></v-select>
            <v-select
              class="col-4"
              outlined
              v-model="apply.after_end_time"
              :items="end_time"
              item-text="name"
              item-value="value"
              label="结束节次"></v-select>
            <v-select
              class="col-4"
              v-model="apply.after_classroom_id"
              :items="classrooms"
              item-text="name"
              item-value="value"
              outlined
              label="上课教室"
            ></v-select>
            <v-select
              class="col-4"
              outlined
              v-model="apply.after_start_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="开始周"></v-select>
            <v-select
              class="col-4"
              outlined
              v-model="apply.after_end_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="结束周"></v-select>
          </v-row>
          <v-row>
            <v-textarea outlined label="请说明更换理由" v-model="apply.reason"></v-textarea>
          </v-row>

        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="submitApplyForClassChange()"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="changeClass = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
export default {
  name: "SelectCourse",
  data() {
    return {
      course: null,
      courses: [],
      headers: [
        {text: '课程名', align: 'start', sortable: false, value: 'course_name'},
        {text: '教师', sortable: false, value: 'name'},
        //学年由上传年份决定
        {text: '开始周', sortable: false, value: 'start_time'},
        {text: '结束周', sortable: false, value: 'end_time'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      changeClass: false,
      start_time: [
        {name: "第一节", value: 1},
        {name: "第二节", value: 2},
        {name: "第三节", value: 3},
        {name: "第四节", value: 4},
        {name: "第五节", value: 5},
        {name: "第六节", value: 6},
        {name: "第七节", value: 7},
        {name: "第八节", value: 8},
        {name: "第九节", value: 9},
        {name: "第十节", value: 10},
        {name: "第十一节", value: 11},
      ],
      end_time: [
        {name: "第一节", value: 1},
        {name: "第二节", value: 2},
        {name: "第三节", value: 3},
        {name: "第四节", value: 4},
        {name: "第五节", value: 5},
        {name: "第六节", value: 6},
        {name: "第七节", value: 7},
        {name: "第八节", value: 8},
        {name: "第九节", value: 9},
        {name: "第十节", value: 10},
        {name: "第十一节", value: 11},
      ],
      weeks: [
        {name: "第一周", value: 1},
        {name: "第二周", value: 2},
        {name: "第三周", value: 3},
        {name: "第四周", value: 4},
        {name: "第五周", value: 5},
        {name: "第六周", value: 6},
        {name: "第七周", value: 7},
        {name: "第八周", value: 8},
        {name: "第九周", value: 9},
        {name: "第十周", value: 10},
        {name: "第十一周", value: 11},
        {name: "第十二周", value: 12},
        {name: "第十三周", value: 13},
        {name: "第十四周", value: 14},
        {name: "第十五周", value: 15},
        {name: "第十六周", value: 16},
        {name: "第十七周", value: 17},
        {name: "第十八周", value: 18},
        {name: "第十九周", value: 19},
        {name: "第二十周", value: 20},
      ],
      week_time: [
        {name: "星期一", value: 1},
        {name: "星期二", value: 2},
        {name: "星期三", value: 3},
        {name: "星期四", value: 4},
        {name: "星期五", value: 5},
        {name: "星期六", value: 6},
        {name: "星期日", value: 7},
      ],
      teachers: [],
      classrooms: [],
      faculties: [],
      directions: [],
      apply: {},
      beforeClass: {},
      afterClass: {},
      selectId: null,
    }
  },
  created() {
    this.$axios({
      method: "get",
      url: 'http://127.0.0.1:9000/getCourseList',
      params: {
        'teacher_id': this.$store.state.teacherId
      }
    }).then((response) => {
      this.courses = response.data.courses
      console.log(this.courses)
    })
  },
  watch: {
    '$store.state.refreshFlag': {
      handler(newValue, oldValue) {
        console.log(newValue, oldValue)
        if (newValue !== 1) return
        this.$axios({
          method: "get",
          url: 'http://127.0.0.1:9000/getCourseList',
          params: {
            'teacher_id': this.$store.state.teacherId
          }
        }).then((response) => {
          this.courses = response.data.courses
          console.log(this.courses)
        })
      },
    },
  },
  methods: {
    submitApplyForClassChange() {
      this.apply.teacher_id = this.$store.state.teacherId
      this.apply.course_id = this.selectId
      this.$axios({
        method: "post",
        url: "http://127.0.0.1:9000/postApply",
        data: this.apply,
        headers: {
          "Content-Type": "multipart/form-data"
        }
      }).then((response)=>{
        this.$store.commit(response.data.snackbar,response.data.msg)
        this.changeClass = false
        console.log(response)
      })
    },
    getInfo() {
      this.$axios({
        method: "get",
        url: 'http://127.0.0.1:9000/getClassesList',
        params: {
          'level': this.$store.state.level
        }
      }).then((response) => {
        this.classrooms = response.data.classrooms
      })
    }
  }
}
</script>

<style scoped>

</style>
