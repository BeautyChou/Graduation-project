<template>
  <div>
    <v-card>
      <v-card-title>
        安排课程
        <v-spacer></v-spacer>
        <v-btn color="green">
          <v-icon>mdi-upload</v-icon>
          <span>
            上传课程
          </span>
        </v-btn>
        <v-btn color="primary" class="mx-8" @click="addClass = true">
          <v-icon>mdi-plus-thick</v-icon>
          <span>
            添加课程
          </span>
        </v-btn>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="classes"
        class="elevation-1"
      >

        <template v-slot:item.week_time="{ item }">
          <v-select
            v-model="item.week_time"
            :items="week_time"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.start_time="{ item }">
          <v-select
            v-model="item.start_time"
            :items="start_time"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.end_time="{ item }">
          <v-select
            v-model="item.end_time"
            :items="end_time"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.max_choose_num="{ item }">
          <v-text-field :value="item.max_choose_num"></v-text-field>
        </template>
        <template v-slot:item.credit="{ item }">
          <!--          <span>{{item}}</span>-->
          <v-text-field :value="item.credit"></v-text-field>
        </template>
        <template v-slot:item.teacher="{ item }">
          <v-select
            v-model="item.teacher_id"
            :items="teachers"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.classes="{ item }">
          <v-select
            v-model="item.classes"
            :items="classTable"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.classroom="{ item }">
          <v-select
            v-model="item.classroom_id"
            :items="classrooms"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.faculties="{ item }">
          <v-select
            v-model="item.faculty_id"
            :items="faculties"
            item-text="name"
            item-value="id"
          ></v-select>
        </template>
        <template v-slot:item.directions="{ item }">
          <v-select
            v-model="item.direction_id"
            :items="directions[item.faculty_id]"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.startWeek="{ item }">
          <v-select
            v-model="item.start_week"
            :items="weeks"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.endWeek="{ item }">
          <v-select
            v-model="item.end_week"
            :items="weeks"
            item-text="name"
            item-value="value"
          ></v-select>
        </template>
        <template v-slot:item.operation="{ item }">
          <v-tooltip v-if="$store.state.level===2||true" bottom>
            <template v-slot:activator="{ on,attrs }">
              <v-btn icon color="error" v-bind="attrs" v-on="on" x-large @click="dialog = true ;selectId = item.course_id;">
                <v-icon>
                  mdi-delete
                </v-icon>
              </v-btn>
            </template>
            <span>删除作业</span>
          </v-tooltip>
        </template>
      </v-data-table>
      <v-dialog v-if="$store.state.level===2||true" v-model="dialog" width="500" persistent>
        <v-card>
          <v-card-title class="headline font-weight-bold">
            警告!
          </v-card-title>
          <v-card-text class="font-weight-bold">
            真的要删除这个课程吗?这个操作是不可逆转的!
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              class="font-weight-bold"
              color="green darken-1"
              text
              @click="deleteClass(selectId)"
            >
              是
            </v-btn>
            <v-btn
              class="font-weight-bold"
              color="green darken-1"
              text
              @click="dialog = false"
            >
              否
            </v-btn>
          </v-card-actions>
        </v-card>

      </v-dialog>
      <v-dialog v-model="addClass" width="700" persistent>
        <v-card>
          <v-card-title class="headline font-weight-bold">
            添加新科目 {{ newClass }}
          </v-card-title>
          <v-card-text class="font-weight-bold">
            <v-text-field
              outlined
              label="科目名称"
              v-model="newClass.course_name"></v-text-field>
            <v-text-field
              outlined
              label="学分"
              v-model="newClass.credit"></v-text-field>
            <v-text-field
              outlined
              label="最大人数"
              v-model="newClass.max_choose_num"></v-text-field>
            <v-select
              v-model="newClass.week_time"
              :items="week_time"
              item-text="name"
              item-value="value"
              label="上课时间"
              outlined></v-select>
            <v-select
              outlined
              v-model="newClass.start_time"
              :items="start_time"
              item-text="name"
              item-value="value"
              label="起始节次"></v-select>
            <v-select
              outlined
              v-model="newClass.end_time"
              :items="end_time"
              item-text="name"
              item-value="value"
              label="结束节次"></v-select>
            <v-select
              outlined
              v-model="newClass.teacher_id"
              :items="teachers"
              item-text="name"
              item-value="value"
              label="教师"></v-select>
            <v-select
              v-model="newClass.classroom_id"
              :items="classrooms"
              item-text="name"
              item-value="value"
              outlined
              label="上课教室"
            ></v-select>
            <v-select
              outlined
              v-model="newClass.start_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="开始周"></v-select>
            <v-select
              outlined
              v-model="newClass.end_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="结束周"></v-select>
            <v-select
              outlined
              v-model="newClass.faculty_id"
              :items="faculties"
              item-text="name"
              item-value="id"
              label="学生所属学院"></v-select>
            <v-select
              outlined
              v-model="newClass.direction_id"
              :items="directions[newClass.faculty_id]"
              item-text="name"
              item-value="value"
              label="学生所属方向"></v-select>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              class="font-weight-bold"
              color="green darken-1"
              text
              @click="submitNewClass()"
            >
              提交
            </v-btn>
            <v-btn
              class="font-weight-bold"
              color="green darken-1"
              text
              @click="addClass = false"
            >
              取消
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-card>
  </div>
</template>

<script>

export default {
  name: "ClassSheet",
  data() {
    return {
      headers: [
        {text: '课程', align: 'start', sortable: false, value: 'course_name'},
        {text: '学分', sortable: false, value: 'credit'},
        {text: '最大人数', sortable: false, value: 'max_choose_num'},
        {text: '上课时间', sortable: false, value: 'week_time'},
        {text: '起始节次', sortable: false, value: 'start_time'},
        {text: '结束节次', sortable: false, value: 'end_time'},
        {text: '教师', sortable: false, value: 'teacher'},
        {text: '教室', sortable: false, value: 'classroom'},
        {text: '开始周', sortable: false, value: 'startWeek'},
        {text: '结束周', sortable: false, value: 'endWeek'},
        {text: '学生所属学院', sortable: false, value: 'faculties'},
        {text: '学生所属方向', sortable: false, value: 'directions'},
        {text: '操作', sortable: false,value: 'operation'}
      ],
      classes: [],
      drag: false,
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
      teachers: [],
      classrooms: [],
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
      faculties: [],
      directions: [],
      addClass: false,
      newClass: {},
      dialog:false,
      selectId:null,
    }
  },
  created() {
    this.$axios({
      method: "get",
      url: 'http://127.0.0.1:9000/getClassesList',
      params: {
        'teacher_id': this.$store.state.teacherId,
        'level': this.$store.state.level
      }
    }).then((response) => {
      this.classes = response.data.classes
      this.teachers = response.data.teachers
      this.classrooms = response.data.classrooms
      this.directions = response.data.directions
      this.faculties = response.data.faculties
      console.log(response)
    })
  },
  components: {},
  methods: {
    submitNewClass(){
      const formData = new FormData()
      formData.append("name",this.newClass.course_name)
      formData.append("credit",this.newClass.credit)
      formData.append("max_choose_num",this.newClass.max_choose_num)
      formData.append("week_time",this.newClass.week_time)
      formData.append("start_time",this.newClass.start_time)
      formData.append("end_time",this.newClass.end_time)
      formData.append("teacher_id",this.newClass.teacher_id)
      formData.append("start_week",this.newClass.start_week)
      formData.append("end_week",this.newClass.end_week)
      formData.append("faculty_id",this.newClass.faculty_id)
      formData.append("classroom_id",this.newClass.classroom_id)
      formData.append("direction_id",this.newClass.direction_id)
      this.$axios({
        method:"post",
        url:"http://127.0.0.1:9000/postNewClass",
        data:formData,
        headers: {
          "Content-Type": "multipart/form-data"
        }
      }).then((response)=>{
        console.log(response)
        this.$store.commit(response.data.snackbar,response.data.msg)
        this.addClass = false
        if (response.data.snackbar === 'setSuccess') {this.classes.push(this.newClass);this.newClass = {}}
      })
    },
    deleteClass(id){
      this.dialog = false
      this.$axios({
        method: "DELETE",
        url: "http://127.0.0.1:9000/deleteClass?course_id=" + id,
      }).then((response) => {
        this.classes.some((item, i) => {
          if (item.course_id === id) {
            this.classes.splice(i, 1)
          }
          this.$store.commit(response.data.snackbar,response.data.msg)
          this.dialog = false
        })
      })
    },
  },
  watch: {}
}
</script>

<style scoped>

</style>
