<template>
  <v-card>
    <v-card-title>
      <v-row>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          dark
          class="mb-2"
          @click="add_classroom = true"
        >
          <v-icon>
            mdi-plus-thick
          </v-icon>
          添加教室
        </v-btn>
      </v-row>
    </v-card-title>
    <v-data-table
      :headers="classrooms_headers"
      :items="classrooms"
      :options.sync="options"
      :server-items-length="total"
    >
      <template v-slot:item.operation="{ item }">

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="modify_classroom = true ;selectOBJ = item">
              <v-icon>
                mdi-lead-pencil
              </v-icon>
            </v-btn>
          </template>
          <span>修改教室信息</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="delete_classroom = true;selectOBJ = item">
              <v-icon>
                mdi-account-remove
              </v-icon>
            </v-btn>
          </template>
          <span>删除教室</span>
        </v-tooltip>

      </template>
    </v-data-table>
    <v-dialog v-if="($store.state.level===2||true)" v-model="delete_classroom" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          警告!
        </v-card-title>
        <v-card-text class="font-weight-bold">
          真的要删除这个教室吗?这个操作是不可逆转的!
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteClassroom"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_classroom = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>

    </v-dialog>
    <v-dialog v-model="modify_classroom" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          修改教室信息
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="教室名"
              v-model="new_name"
              class="col-6"></v-text-field>
            <v-text-field
              outlined
              label="最大人数"
              v-model="new_max_num"
              class="col-6"></v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modifyClassroom"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modify_classroom = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_classroom" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加教室信息
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="教室名"
              v-model="new_name"
              class="col-6"></v-text-field>
            <v-text-field
              outlined
              label="最大人数"
              v-model="new_max_num"
              class="col-6"></v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addClassroom"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_classroom = false"
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
  name: "ManageClassroom",
  created() {
    this.getClassroomList()
  },
  data() {
    return {
      options: {},
      total: null,
      new_name: null,
      new_max_num: null,
      selectOBJ: {},
      delete_classroom: false,
      modify_classroom: false,
      add_classroom: false,
      title_id: null,
      classrooms: [],
      classrooms_headers: [
        {text: '教室名', align: 'start', sortable: false, value: 'name'},
        {text: '最大容纳人数', align: 'start', sortable: false, value: 'max_number'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
    }
  },
  watch: {
    options: {
      handler() {
        this.getClassroomList()
      },
      deep: true,
    }
  },
  methods: {
    getClassroomList() {
      this.$axios({
        url: "Classroom",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        params: {
          "page": this.options.page,
          "items": this.options.itemsPerPage,
        },
        method: "get"
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.classrooms = response.data.classrooms
        this.total = response.data.total
      })
    },
    modifyClassroom() {
      const formData = new FormData()
      formData.append("classroom_id", this.selectOBJ.id)
      formData.append("name", this.new_name)
      formData.append("max_num", this.new_max_num)
      this.$axios({
        method: "put",
        url: "Classroom",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.modify_classroom = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.new_name = null
        this.new_max_num = null
        this.getClassroomList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    deleteClassroom() {
      this.$axios({
        url: "Classroom",
        method: "delete",
        params: {
          classroom_id: this.selectOBJ.id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.delete_classroom = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getClassroomList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    addClassroom() {
      const formData = new FormData()
      formData.append("name", this.new_name)
      formData.append("max_num", this.new_max_num)
      this.$axios({
        method: "post",
        url: "Classroom",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.add_classroom = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getClassroomList()
        this.new_name = null
        this.new_max_num = null
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
  },
  inject: ['expire']
}
</script>

<style scoped>

</style>
