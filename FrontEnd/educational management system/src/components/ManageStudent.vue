<template>
  <v-card>
    <v-card-title>
      <v-btn
        color="primary"
        dark
        class="mb-2"
        large
        @click="add_student = true"
      >
        <v-icon>
          mdi-plus-thick
        </v-icon>
        添加学生
      </v-btn>
      <v-select
        class="col-4"
        outlined
        dense
        v-model="faculty_id"
        :items="faculties"
        item-text="name"
        item-value="id"
        label="学生所属学院"
        @change="getStudentList"
      ></v-select>
      <v-select
        class="col-4"
        outlined
        dense
        v-model="specialty_id"
        :items="specialties[faculty_id]"
        item-text="specialty_name"
        item-value="specialty_id"
        label="学生所属专业"
        :disabled="faculty_id ==  null"
        @change="getStudentList"></v-select>
    </v-card-title>
    <v-data-table
      :headers="student_headers"
      :items="students"
    >
      <template v-slot:item.direction="{ item }">
        {{ item.direction_id === 0 ? "无方向" : item.direction_name }}
      </template>
      <template v-slot:item.practice="{ item }">
        {{ item.practice === 0 ? " 统一实习" : "自主实习" }}
      </template>
      <template v-slot:item.operation="{ item }">

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="modify_student = true ;selectOBJ = item">
              <v-icon>
                mdi-lead-pencil
              </v-icon>
            </v-btn>
          </template>
          <span>修改学生信息</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="delete_student = true;selectOBJ = item">
              <v-icon>
                mdi-account-remove
              </v-icon>
            </v-btn>
          </template>
          <span>删除学生</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="punish_student = true;selectOBJ = item">
              <v-icon>
                mdi-alert
              </v-icon>
            </v-btn>
          </template>
          <span>给予处分</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="selectOBJ = item;getPunishments();cancel_punishment = true;">
              <v-icon>
                mdi-alert-minus
              </v-icon>
            </v-btn>
          </template>
          <span>查看/取消处分</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="selectOBJ = item;getPracticeInfo();practice = true" :disabled="item.practice === 0">
              <v-icon>
                mdi-factory
              </v-icon>
            </v-btn>
          </template>
          <span>查看实习单位</span>
        </v-tooltip>

      </template>
    </v-data-table>
    <v-dialog v-if="($store.state.level===2||true)" v-model="delete_student" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          警告!
        </v-card-title>
        <v-card-text class="font-weight-bold">
          真的要删除这个学生吗?这个操作是不可逆转的!
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteStudent"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_student = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>

    </v-dialog>
    <v-dialog v-model="modify_student" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          修改学生信息:{{ selectOBJ.name }} (仅填写需要修改的地方)
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="密码"
              v-model="new_password"
              class="col-12"></v-text-field>
            <v-select
              class="col-4"
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              outlined
              label="学院"
            ></v-select>
            <v-select
              class="col-4"
              v-model="new_specialty"
              :disabled="new_faculty === null"
              :items="specialties[new_faculty]"
              item-text="specialty_name"
              item-value="specialty_id"
              outlined
              label="专业"
            ></v-select>
            <v-select
              class="col-4"
              v-model="new_direction"
              :disabled="new_specialty === null"
              :items="directions[new_specialty]"
              item-text="direction_name"
              item-value="direction_id"
              outlined
              label="方向"
            ></v-select>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modifyStudent"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modify_student = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="punish_student" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          给予学生 {{ selectOBJ.name }} 处分
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-select
              class="col-12"
              v-model="new_punishment_level"
              :items="punishment_level"
              item-text="level"
              item-value="id"
              outlined
              label="处分级别"
            ></v-select>
            <v-textarea
              outlined
              v-model="punish_reason"
              label="给予处分理由"
            ></v-textarea>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            :disabled="punish_reason == null || new_punishment_level == null"
            @click="punishStudent"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="punish_student = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="cancel_punishment" width="1200">
      <v-card>
        <v-card-title class="headline font-weight-bold">
          查看/取消{{ selectOBJ.name }}处分
        </v-card-title>
        <v-data-table
          :headers="punishment_headers"
          :items="punishments">
          <template v-slot:item.operation="{ item }">
            <v-tooltip v-if="$store.state.level===2||true" bottom>
              <template v-slot:activator="{ on,attrs }">
                <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                       @click="cancelPunishment();selectOBJ = item" :disabled="item.is_cancelled === true">
                  <v-icon>
                    mdi-alert-remove
                  </v-icon>
                </v-btn>
              </template>
              <span>取消处分</span>
            </v-tooltip>
          </template>
        </v-data-table>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_student" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加学生信息
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="姓名"
              v-model="new_name"
              class="col-6"></v-text-field>
            <v-text-field
              outlined
              label="密码"
              v-model="new_password"
              class="col-6"></v-text-field>
            <v-select
              class="col-3"
              v-model="new_grade"
              :items="grades"
              item-text="name"
              item-value="id"
              outlined
              label="年级"
            ></v-select>
            <v-select
              class="col-3"
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              outlined
              label="学院"
            ></v-select>
            <v-select
              class="col-3"
              v-model="new_specialty"
              :disabled="new_faculty === null"
              :items="specialties[new_faculty]"
              item-text="specialty_name"
              item-value="specialty_id"
              outlined
              label="专业"
            ></v-select>
            <v-select
              class="col-3"
              v-model="new_direction"
              :disabled="new_specialty === null"
              :items="directions[new_specialty]"
              item-text="direction_name"
              item-value="direction_id"
              outlined
              label="方向"
            ></v-select>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addStudent"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_student = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="practice" width="1200">
      <v-card>
        <v-card-title class="headline font-weight-bold">
          自主实习申请表
        </v-card-title>
        <v-col cols="2">
          <v-text-field
            label="在外联系方式"
            outlined
            v-model.number="independent_practice.phone"
          readonly></v-text-field>
        </v-col>
        <v-col cols="5">
          <v-text-field
            outlined
            v-model="independent_practice.start_time"
            label="开始实习日期"
            :disabled="$store.state.modifyHomeworkFlag"
            readonly
          ></v-text-field>
        </v-col>
        <v-col cols="5">
          <v-text-field
            outlined
            v-model="independent_practice.end_time"
            label="结束实习日期"
            :disabled="$store.state.modifyHomeworkFlag"
            readonly
          ></v-text-field>
        </v-col>
        <v-col cols="12">
          <v-textarea
            label="申请理由"
            outlined
            v-model="independent_practice.reason"
            readonly
          ></v-textarea>
        </v-col>
        <v-col cols="6">
          <v-text-field
            label="实习单位名称"
            outlined
            v-model="independent_practice.company_name"
            readonly
          ></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-text-field
            label="实习单位地址"
            outlined v-model="independent_practice.company_address"
            readonly
          ></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-text-field
            label="实习单位联系人"
            outlined v-model="independent_practice.company_person"
            readonly
          ></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-text-field
            label="实习单位电话"
            outlined
            v-model="independent_practice.company_phone"
            readonly
          ></v-text-field>
        </v-col>
        <v-col cols="12">
          <v-text-field
            label="在外住宿地址"
            outlined
            v-model="independent_practice.address"
            readonly
          ></v-text-field>
        </v-col>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
export default {
  name: "ManageStudent",
  data() {
    return {
      independent_practice: {},
      practice: false,
      new_grade: null,
      new_name: null,
      add_student: null,
      punishments: [],
      cancel_punishment: false,
      punish_reason: null,
      new_punishment_level: null,
      punish_student: false,
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
      grades: [
        {id: 1, name: "大一"},
        {id: 2, name: "大二"},
        {id: 3, name: "大三"},
        {id: 4, name: "大四"}
      ],
      student_headers: [
        {text: '姓名', align: 'start', sortable: false, value: 'name'},
        {text: '学院', sortable: false, value: 'faculty_name'},
        {text: '专业', sortable: false, value: 'specialty_name'},
        {text: '方向', sortable: false, value: 'direction'},
        {text: '实习方式', sortable: false, value: 'practice'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
      punishment_headers: [
        {text: '处分等级', align: 'start', sortable: false, value: 'punishment_name'},
        {text: '处分原因', sortable: false, value: 'reason'},
        {text: '操作', sortable: false, value: 'operation'},
      ],
    }
  },
  created() {
    this.getStudentList()
  },
  methods: {
    addStudent() {
      const formData = new FormData()
      formData.append("name", this.new_name)
      formData.append("grade", this.new_grade)
      formData.append("password", this.new_password)
      formData.append("faculty_id", this.new_faculty)
      formData.append("specialty_id", this.new_specialty)
      formData.append("direction_id", this.new_direction)
      this.$axios({
        method: "post",
        url: "addStudent",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        this.add_student = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getStudentList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
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
    modifyStudent() {
      const formData = new FormData()
      formData.append("student_id", this.selectOBJ.student_id)
      formData.append("password", this.new_password)
      formData.append("faculty_id", this.new_faculty)
      formData.append("specialty_id", this.new_specialty)
      formData.append("direction_id", this.new_direction)
      this.$axios({
        method: "put",
        url: "putStudent",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        this.modify_student = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getStudentList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    punishStudent() {
      var formData = new FormData()
      formData.append("student_id", this.selectOBJ.student_id)
      formData.append("punishment_level", this.new_punishment_level)
      formData.append("punishment_content", this.punish_reason)
      this.$axios({
        method: "post",
        url: "punishStudent",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        this.punish_student = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    getPunishments() {
      this.$axios({
        url: "getPunishment",
        params: {
          student_id: this.selectOBJ.student_id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get"
      }).then((response) => {
        console.log(response)
        this.punishments = response.data.punishments
      })
    },
    cancelPunishment() {
      this.$axios({
        url: "deletePunishment",
        method: "delete",
        params: {
          punishment_id: this.selectOBJ.id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        this.cancel_punishment = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    getPracticeInfo() {
      this.$axios({
        url: "getPracticeInfo",
        params: {
          student_id: this.selectOBJ.student_id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get"
      }).then((response) => {
        console.log(response)
        this.independent_practice = response.data.independent_practice
        this.independent_practice.start_time = this.independent_practice.start_time.substr(0,10)
        this.independent_practice.end_time = this.independent_practice.end_time.substr(0,10)
      })
    },
  },
}
</script>

<style scoped>

</style>
