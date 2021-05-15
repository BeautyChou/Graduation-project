<template>
  <div>
    <v-card>
      <v-card-title>
        安排课程
        <v-spacer></v-spacer>
        <v-btn color="green" :disabled="uploadFlag" @click="uploadClasses">
          <v-icon>mdi-upload</v-icon>
          <span>
            上传课程
          </span>
        </v-btn>
        <v-btn color="primary" class="mx-8" @click="addClass = true;copyFlag = false">
          <v-icon>mdi-plus-thick</v-icon>
          <span>
            添加课程
          </span>
        </v-btn>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="classes"
        :options.sync="options"
        :server-items-length="total"
        class="elevation-1"
      >

        <template v-slot:item.week_time="{ item }">
          <v-select
            v-model="item.week_time"
            :items="week_time"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.start_time="{ item }">
          <v-select
            v-model="item.start_time"
            :items="start_time"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.end_time="{ item }">
          <v-select
            v-model="item.end_time"
            :items="end_time"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.max_choose_num="{ item }">
          <v-text-field
            :value="item.max_choose_num"
            @change="pushing(item)"></v-text-field>
        </template>
        <template v-slot:item.credit="{ item }">
          <v-text-field
            v-model="item.credit"
            :disabled="item.copy_flag > 0"
            :value="item.credit"
            @change="pushing(item)"></v-text-field>
        </template>
        <template v-slot:item.teacher="{ item }">
          <v-select
            v-model="item.teacher_id"
            :items="teachers"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.classes="{ item }">
          <v-select
            v-model="item.classes"
            :items="classTable"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.classroom="{ item }">
          <v-select
            v-model="item.classroom_id"
            :items="classrooms"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.faculties="{ item }">
          <v-select
            v-model="item.faculty_id"
            :items="faculties"
            item-text="name"
            item-value="id"
            @change="pushing(item,'')"
          ></v-select>
        </template>
        <template v-slot:item.specialties="{ item }">
          <v-select
            v-model="item.specialty_id"
            :items="specialties[item.faculty_id]"
            item-text="specialty_name"
            item-value="specialty_id"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.directions="{ item }">
          <v-select
            v-model="item.direction_id"
            :items="directions[item.specialty_id]"
            item-text="direction_name"
            item-value="direction_id"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.startWeek="{ item }">
          <v-select
            v-model="item.start_week"
            :items="weeks"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.endWeek="{ item }">
          <v-select
            v-model="item.end_week"
            :items="weeks"
            item-text="name"
            item-value="value"
            @change="pushing(item)"
          ></v-select>
        </template>
        <template v-slot:item.operation="{ item }">
          <v-tooltip v-if="$store.state.level===2||true" bottom>
            <template v-slot:activator="{ on,attrs }">
              <v-btn
                icon
                color="error"
                v-bind="attrs"
                v-on="on"
                x-large
                @click="dialog = true ;selectRecordID = item.record_id;"
                :disabled="item.copy_flag > 0">
                <v-icon>
                  mdi-delete
                </v-icon>
              </v-btn>
            </template>
            <span>删除课程</span>
          </v-tooltip>
          <v-tooltip v-if="$store.state.level===2||true" bottom>
            <template v-slot:activator="{ on,attrs }">
              <v-btn
                icon
                color="primary"
                v-bind="attrs"
                v-on="on"
                x-large
                @click="copyClass(item)"
                :disabled="item.copy_flag > 0">
                <v-icon>
                  mdi-plus-thick
                </v-icon>
              </v-btn>
            </template>
            <span>添加当周课程</span>
          </v-tooltip>
          <v-tooltip v-if="$store.state.level===2||true" bottom>
            <template v-slot:activator="{ on,attrs }">
              <v-btn
                icon
                color="primary"
                v-bind="attrs"
                v-on="on"
                x-large
                @click="copyCourse(item)"
                :disabled="item.copy_flag > 0">
                <v-icon>
                  mdi-plus-thick
                </v-icon>
              </v-btn>
            </template>
            <span>添加班级</span>
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
              @click="deleteClass(selectRecordID)"
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
      <v-dialog v-model="addClass" width="900" persistent>
        <v-card>
          <v-card-title class="headline font-weight-bold">
            添加新科目
          </v-card-title>
          <v-card-text class="font-weight-bold">
            <v-row class="mt-2">
              <v-select
                :disabled="newClass.week_time == null || newClass.start_time == null || newClass.end_time == null || newClass.start_week == null || newClass.end_week == null||newClass.faculty_id == null|| newClass.specialty_id == null||newClass.direction_id == null"
                @click="validClassrooms()"
                class="col-4"
                v-model="newClass.classroom_id"
                :items="valClassrooms"
                item-text="name"
                item-value="value"
                outlined
                label="上课教室"
              ></v-select>
              <v-select
                class="col-4"
                v-model="newClass.selectable"
                :disabled="copyFlag||copyCourseFlag"
                :items="selectable"
                item-text="name"
                item-value="value"
                outlined
                label="是否可被学生选择"
              ></v-select>
              <v-select
                class="col-4"
                v-model="newClass.grades"
                :disabled="copyFlag||copyCourseFlag"
                :items="grades"
                item-text="name"
                item-value="id"
                outlined
                label="学生年级"
              ></v-select>
              <v-text-field
                outlined
                label="科目名称"
                v-model="newClass.course_name"
                :disabled="copyFlag||copyCourseFlag"
                class="col-4"></v-text-field>
              <v-text-field
                class="col-4"
                outlined
                label="学分"
                v-model="newClass.credit"
                :disabled="copyFlag||copyCourseFlag"></v-text-field>
              <v-text-field
                class="col-4"
                outlined
                label="最大人数"
                v-model="newClass.max_choose_num"
                :disabled="copyFlag||copyCourseFlag"></v-text-field>
              <v-select
                @change="valClassrooms"
                class="col-4"
                v-model="newClass.week_time"
                :items="week_time"
                item-text="name"
                item-value="value"
                label="上课时间"
                outlined></v-select>
              <v-select
                @change="valClassrooms"
                class="col-4"
                outlined
                v-model="newClass.start_time"
                :items="start_time"
                item-text="name"
                item-value="value"
                label="起始节次"></v-select>
              <v-select
                @change="valClassrooms"
                class="col-4"
                outlined
                v-model="newClass.end_time"
                :items="end_time"
                item-text="name"
                item-value="value"
                label="结束节次"></v-select>
              <v-select
                class="col-4"
                outlined
                v-model="newClass.teacher_id"
                :items="teachers"
                item-text="name"
                item-value="value"
                label="教师"
                :disabled="copyFlag"></v-select>
              <v-select
                @change="valClassrooms"
                class="col-4"
                outlined
                v-model="newClass.start_week"
                :items="weeks"
                item-text="name"
                item-value="value"
                label="开始周"></v-select>
              <v-select
                @change="valClassrooms"
                class="col-4"
                outlined
                v-model="newClass.end_week"
                :items="weeks"
                item-text="name"
                item-value="value"
                label="结束周"></v-select>
              <v-select
                class="col-4"
                outlined
                v-model="newClass.faculty_id"
                :items="faculties"
                item-text="name"
                item-value="id"
                label="学生所属学院"
                :disabled="copyFlag||copyCourseFlag"></v-select>
              <v-select
                class="col-4"
                outlined
                v-model="newClass.specialty_id"
                :items="specialties[newClass.faculty_id]"
                item-text="specialty_name"
                item-value="specialty_id"
                label="学生所属专业"
                :disabled="copyFlag||copyCourseFlag||newClass.faculty_id ==  null"></v-select>
              <v-select
                class="col-4"
                outlined
                v-model="newClass.direction_id"
                :items="directions[newClass.specialty_id]"
                item-text="direction_name"
                item-value="direction_id"
                label="学生所属方向"
                :disabled="copyFlag||copyCourseFlag||newClass.faculty_id == null|| newClass.specialty_id == null"></v-select>

            </v-row>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              class="font-weight-bold"
              color="green darken-1"
              text
              @click="submitNewClass()"
              :disabled="newClass.course_name==null||newClass.credit==null||newClass.max_choose_num==null||newClass.week_time==null||newClass.start_time==null||newClass.end_time==null||newClass.teacher_id==null||newClass.classroom_id==null||newClass.start_week==null||newClass.end_week==null||newClass.faculty_id==null||newClass.direction_id==null||newClass.selectable==null"
            >
              提交
            </v-btn>
            <v-btn
              class="font-weight-bold"
              color="green darken-1"
              text
              @click="addClass = false;copyCourseFlag = false;copyFlag = false;newClass = {}"
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
  name: "SetClass",
  data() {
    return {
      options: {},
      total: null,
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
        {text: '学生所属专业', sortable: false, value: 'specialties'},
        {text: '学生所属方向', sortable: false, value: 'directions'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      classes: [],
      drag: false,
      selectable: [
        {name: "是", value: true},
        {name: "否", value: false},
      ],
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
        {name: "第十二节", value: 12},
        {name: "第十三节", value: 13},
        {name: "第十四节", value: 14},
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
        {name: "第十二节", value: 12},
        {name: "第十三节", value: 13},
        {name: "第十四节", value: 14},
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
      valClassrooms: [],
      faculties: [],
      specialties: [],
      directions: [],
      addClass: false,
      newClass: {},
      dialog: false,
      selectId: null,
      changedClasses: [],
      uploadFlag: true,
      copyFlag: false,
      str: '',
      copyCourseFlag: false,
      selectRecordID: null,
      grades: [
        {id: 0, name: "不限"},
        {id: 1, name: "大一"},
        {id: 2, name: "大二"},
        {id: 3, name: "大三"},
        {id: 4, name: "大四"}
      ],
    }
  },
  created() {
    this.getClass()
  },
  components: {},
  methods: {
    copyCourse(classOBJ) { //复制班级(允许更多人上这门课)
      this.copyCourseFlag = true
      this.addClass = true
      this.selectId = classOBJ.course_id
      this.newClass.course_name = classOBJ.course_name
      this.newClass.credit = classOBJ.credit
      this.newClass.max_choose_num = classOBJ.max_choose_num
      this.newClass.faculty_id = classOBJ.faculty_id
      this.newClass.specialty_id = classOBJ.specialty_id
      this.newClass.direction_id = classOBJ.direction_id
      this.newClass.grades = classOBJ.grades
      this.newClass.selectable = classOBJ.selectable
      this.str = 'course'
    },
    copyClass(classOBJ) {
      this.copyFlag = true
      this.addClass = true
      this.selectId = classOBJ.course_id
      this.selectRecordID = classOBJ.record_id
      this.newClass.course_name = classOBJ.course_name
      this.newClass.credit = classOBJ.credit
      this.newClass.max_choose_num = classOBJ.max_choose_num
      this.newClass.teacher_id = classOBJ.teacher_id
      this.newClass.faculty_id = classOBJ.faculty_id
      this.newClass.specialty_id = classOBJ.specialty_id
      this.newClass.direction_id = classOBJ.direction_id
      this.newClass.selectable = classOBJ.selectable
      this.newClass.grades = classOBJ.grades
      this.str = 'copy'
    },
    submitNewClass() {
      const formData = new FormData()
      if (this.str === 'copy') {
        this.$set(this.newClass, "copy_flag", this.selectRecordID);
      } else this.$set(this.newClass, "copy_flag", 0);
      formData.append("name", this.newClass.course_name)
      formData.append("credit", this.newClass.credit)
      formData.append("max_choose_num", this.newClass.max_choose_num)
      formData.append("week_time", this.newClass.week_time)
      formData.append("start_time", this.newClass.start_time)
      formData.append("end_time", this.newClass.end_time)
      formData.append("teacher_id", this.newClass.teacher_id)
      formData.append("start_week", this.newClass.start_week)
      formData.append("end_week", this.newClass.end_week)
      formData.append("faculty_id", this.newClass.faculty_id)
      formData.append("classroom_id", this.newClass.classroom_id)
      formData.append("specialty_id", this.newClass.specialty_id)
      formData.append("direction_id", this.newClass.direction_id)
      formData.append("copy_flag", this.newClass.copy_flag)
      formData.append("flag", this.str)
      formData.append("selectable", this.newClass.selectable)
      formData.append("grades", this.newClass.grades)
      this.str = ''
      this.$axios({
        method: "post",
        url: "Class",
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
        this.addClass = false
        this.copyFlag = false
        this.copyCourseFlag = false
        if (response.data.snackbar === 'setSuccess') {
          this.classes.push(this.newClass);
          this.newClass = {};
          this.getClass()
        }
      })
    },
    deleteClass(RecordID) {
      this.dialog = false
      this.$axios({
        method: "DELETE",
        url: "Class?record_id=" + RecordID,
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        for (let i = 0; i < this.changedClasses.length; i++) {
          if (this.changedClasses[i].record_id === RecordID || this.changedClasses[i].copy_flag === RecordID) {
            this.changedClasses.splice(i, 1);
            i--
          }
        }
        if (this.changedClasses.length === 0) this.uploadFlag = true
        for (let i = 0; i < this.classes.length; i++) {
          if (this.classes[i].record_id === RecordID || this.classes[i].copy_flag === RecordID) {
            this.classes.splice(i, 1);
            i--
          }
        }
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
        this.dialog = false
      })
    },
    uploadClasses() {
      this.$axios({
        method: "put",
        url: "Class",
        data: this.changedClasses,
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
        if (response.data.snackbar === "setSuccess") {
          this.changedClasses = []
          this.uploadFlag = true
        }
      })
    },
    pushing(classOBJ, str) {
      var flag = true
      this.changedClasses.some((item, i) => {
        if (item.record_id === classOBJ.record_id) {
          flag = false
          this.changedClasses.splice(i, 1, item)
        }
      })
      if (str !== 'recursion') {
        this.classes.forEach((item, i) => {
          if (item.course_id === classOBJ.course_id) {
            if (item.credit !== classOBJ.credit) {
              item.credit = classOBJ.credit
              this.pushing(item, 'recursion')
            }
          }
        })
      }
      if (flag) this.changedClasses.push(classOBJ)
      this.uploadFlag = false
    },
    validClassrooms() {
      if (this.newClass.week_time == null || this.newClass.start_time == null || this.newClass.end_time == null || this.newClass.start_week == null || this.newClass.end_week == null) return
      const formdata = new FormData
      formdata.append('course_id', this.selectId)
      formdata.append('selectable', this.newClass.selectable)
      formdata.append("max_choose_num", this.newClass.max_choose_num)
      formdata.append('week_time', this.newClass.week_time)
      formdata.append('start_time', this.newClass.start_time)
      formdata.append('end_time', this.newClass.end_time)
      formdata.append('start_week', this.newClass.start_week)
      formdata.append('end_week', this.newClass.end_week)
      formdata.append('classroom_id', this.newClass.classroom_id)
      formdata.append('teacher_id', this.newClass.teacher_id)
      formdata.append('faculty_id', this.newClass.faculty_id)
      formdata.append('direction_id', this.newClass.direction_id)
      formdata.append('specialty_id', this.newClass.specialty_id)
      formdata.append('time', "after")
      this.$axios({
        method: "post",
        url: "http://127.0.0.1:9000/validClassrooms",
        data: formdata,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.valClassrooms = response.data.classrooms
        if (this.valClassrooms.length !== 0) {
          this.valClassrooms.forEach((v, i) => {
            v.name = "教室号:" + v.name + "  人数：" + v.max_num
          })
        } else this.$set(this.newClass, 'classroom_id', null)
      })
    },
    getClass() {
      this.$axios({
        method: "get",
        url: 'Class',
        params: {
          'teacher_id': this.$store.state.teacherId,
          'level': this.$store.state.level,
          "page": this.options.page,
          "items": this.options.itemsPerPage,
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.classes = response.data.classes
        this.teachers = response.data.teachers
        this.classrooms = response.data.classrooms
        this.specialties = response.data.specialties
        this.directions = response.data.directions
        this.faculties = response.data.faculties
        this.total = response.data.total
      })
    }
  },
  watch: {
    options: {
      handler() {
        this.getClass()
      },
      deep: true,
    }
  }
}
</script>

<style scoped>

</style>
