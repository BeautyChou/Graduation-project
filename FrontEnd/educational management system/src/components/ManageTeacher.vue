<template>
  <v-card>
    <v-card-title>
      <v-row>
        <v-select
          class="col-3"
          outlined
          dense
          v-model="faculty_id"
          :items="faculties"
          item-text="name"
          item-value="id"
          label="教师所属学院"
          @change="getTeacherList"
        ></v-select>
        <v-select
          class="col-3"
          outlined
          dense
          v-model="title_id"
          :items="titles"
          item-text="name"
          item-value="id"
          label="教师职称"
          @change="getTeacherList"
        ></v-select>
        <v-btn
          icon
          @click.native="faculty_id = null;title_id = null;getTeacherList()">
          <v-icon>
            mdi-trash-can-outline
          </v-icon>
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          dark
          class="mb-2"
          @click="add_teacher = true"
        >
          <v-icon>
            mdi-plus-thick
          </v-icon>
          添加教师
        </v-btn>
      </v-row>
    </v-card-title>
    <v-data-table
      :headers="teachers_headers"
      :items="teachers"
      :options.sync="options"
      :server-items-length="total"
    >
      <template v-slot:item.operation="{ item }">

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="modify_teacher = true ;selectOBJ = item">
              <v-icon>
                mdi-lead-pencil
              </v-icon>
            </v-btn>
          </template>
          <span>修改教师信息</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="delete_teacher = true;selectOBJ = item">
              <v-icon>
                mdi-account-remove
              </v-icon>
            </v-btn>
          </template>
          <span>删除教师</span>
        </v-tooltip>

      </template>
    </v-data-table>
    <v-dialog v-if="($store.state.level===2||true)" v-model="delete_teacher" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          警告!
        </v-card-title>
        <v-card-text class="font-weight-bold">
          真的要删除这个教师吗?这个操作是不可逆转的!
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteTeacher"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_teacher = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>

    </v-dialog>
    <v-dialog v-model="modify_teacher" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          修改教师信息:{{ selectOBJ.name }} (仅填写需要修改的地方)
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="密码"
              v-model="new_password"
              class="col-12"></v-text-field>
            <v-select
              class="col-6"
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              outlined
              label="学院"
            ></v-select>
            <v-select
              class="col-6"
              v-model="new_title"
              :items="titles"
              item-text="name"
              item-value="id"
              outlined
              label="职称"
            ></v-select>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modifyTeacher"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modify_teacher = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_teacher" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加教师信息
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
              class="col-6"
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              outlined
              label="学院"
            ></v-select>
            <v-select
              class="col-6"
              v-model="new_title"
              :items="titles"
              item-text="name"
              item-value="id"
              outlined
              label="职称"
            ></v-select>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addTeacher"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_teacher = false"
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
  name: "ManageTeacher",
  created() {
    this.getTeacherList()
  },
  data() {
    return {
      options: {},
      total: null,
      new_name: null,
      new_title: null,
      new_faculty: null,
      new_password: null,
      selectOBJ: {},
      delete_teacher: false,
      modify_teacher: false,
      add_teacher: false,
      faculty_id: null,
      title_id: null,
      faculties: [],
      teachers: [],
      titles: [],
      teachers_headers: [
        {text: '姓名', align: 'start', sortable: false, value: 'name'},
        {text: '学院', sortable: false, value: 'faculty_name'},
        {text: '职称', sortable: false, value: 'title_name'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
    }
  },
  methods: {
    getTeacherList() {
      this.$axios({
        url: "Teacher",
        params: {
          faculty_id: this.faculty_id,
          title_id: this.title_id,
          "page": this.options.page,
          "items": this.options.itemsPerPage,
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get"
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.faculties = response.data.faculties
        this.teachers = response.data.teachers
        this.titles = response.data.titles
        this.total = response.data.total
      })
    },
    modifyTeacher() {
      const formData = new FormData()
      formData.append("teacher_id", this.selectOBJ.id)
      formData.append("password", this.new_password)
      formData.append("faculty_id", this.new_faculty)
      formData.append("title_id", this.new_title)
      this.$axios({
        method: "put",
        url: "Teacher",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.modify_teacher = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.new_password = null
        this.new_faculty = null
        this.new_title = null
        this.getTeacherList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    deleteTeacher() {
      this.$axios({
        url: "Teacher",
        method: "delete",
        params: {
          teacher_id: this.selectOBJ.id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.delete_teacher = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getTeacherList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    addTeacher() {
      const formData = new FormData()
      formData.append("name", this.new_name)
      formData.append("password", this.new_password)
      formData.append("faculty_id", this.new_faculty)
      formData.append("title_id", this.new_title)
      this.$axios({
        method: "post",
        url: "Teacher",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.add_teacher = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getTeacherList()
        this.new_name = null
        this.new_password = null
        this.new_faculty = null
        this.new_title = null
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
  },
  watch: {
    options: {
      handler() {
        this.getTeacherList()
      },
      deep: true,
    }
  },
  inject: ['expire']
}
</script>

<style scoped>

</style>
