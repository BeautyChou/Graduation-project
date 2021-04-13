<template>
  <v-card>
    <v-card-title>
      <v-row>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          dark
          class="mb-2"
          @click="add_admin = true"
        >
          <v-icon>
            mdi-plus-thick
          </v-icon>
          添加管理员
        </v-btn>
      </v-row>
    </v-card-title>
    <v-data-table
      :headers="admins_headers"
      :items="admins"
      :options.sync="options"
      :server-items-length="total"
    >
      <template v-slot:item.operation="{ item }">

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="modify_admin = true ;selectOBJ = item">
              <v-icon>
                mdi-lead-pencil
              </v-icon>
            </v-btn>
          </template>
          <span>修改管理员信息</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="delete_admin = true;selectOBJ = item">
              <v-icon>
                mdi-account-remove
              </v-icon>
            </v-btn>
          </template>
          <span>删除管理员</span>
        </v-tooltip>

      </template>
    </v-data-table>
    <v-dialog v-if="($store.state.level===2||true)" v-model="delete_admin" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          警告!
        </v-card-title>
        <v-card-text class="font-weight-bold">
          真的要删除这个管理员吗?这个操作是不可逆转的!
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteAdmin"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_admin = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>

    </v-dialog>
    <v-dialog v-model="modify_admin" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          修改管理员信息:{{ selectOBJ.name }}
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="密码"
              v-model="new_password"
              class="col-12"></v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modifyAdmin"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modify_admin = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_admin" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加管理员信息
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
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addAdmin"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_admin = false"
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
  name: "ManageAdmin",
  data(){
    return{
      options:{},
      total:null,
      new_name:null,
      new_password:null,
      selectOBJ:{},
      delete_admin:false,
      modify_admin:false,
      add_admin:false,
      admins:[],
      admins_headers: [
        {text: '姓名', align: 'start', sortable: false, value: 'name'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
    }
  },
  methods:{
    getAdminList() {
      this.$axios({
        url: "Admin",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get",
        params:{
          "page":this.options.page,
          "items":this.options.itemsPerPage,
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        console.log(response)
        this.admins = response.data.admins
        this.total = response.data.total
      })
    },
    modifyAdmin(){
      const formData = new FormData()
      formData.append("admin_id", this.selectOBJ.ID)
      formData.append("password", this.new_password)
      this.$axios({
        method: "put",
        url: "Admin",
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.modify_admin = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.new_password = null
        this.getAdminList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    deleteAdmin(){
      this.$axios({
        url: "Admin",
        method: "delete",
        params: {
          "admin_id": this.selectOBJ.id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        console.log(this.selectOBJ)
        this.delete_admin = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getAdminList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    addAdmin(){
      const formData = new FormData()
      formData.append("name", this.new_name)
      formData.append("password", this.new_password)
      this.$axios({
        method: "post",
        url: "Admin",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.add_admin = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getAdminList()
        this.new_name = null
        this.new_password = null
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
  },
  created() {
    this.getAdminList()
  },
  watch:{
    options:{
      handler(){
        this.getAdminList()
      },
      deep:true,
    }
  },
  inject:['expire']
}
</script>

<style scoped>

</style>
