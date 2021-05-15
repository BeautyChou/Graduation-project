<template>
  <v-card>
    <v-card-title>
      <v-row>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          dark
          class="mb-2"
          @click="add_title = true"
        >
          <v-icon>
            mdi-plus-thick
          </v-icon>
          添加职称
        </v-btn>
      </v-row>
    </v-card-title>
    <v-data-table
      :headers="titles_headers"
      :items="titles"
      hide-default-footer
    >
      <template v-slot:item.operation="{ item }">

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="modify_title = true ;selectOBJ = item">
              <v-icon>
                mdi-lead-pencil
              </v-icon>
            </v-btn>
          </template>
          <span>修改职称名</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="delete_title = true;selectOBJ = item">
              <v-icon>
                mdi-account-remove
              </v-icon>
            </v-btn>
          </template>
          <span>删除职称</span>
        </v-tooltip>

      </template>
    </v-data-table>
    <v-dialog v-if="($store.state.level===2||true)" v-model="delete_title" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          警告!
        </v-card-title>
        <v-card-text class="font-weight-bold">
          真的要删除这个职称吗?这个操作是不可逆转的!
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteTitle"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_title = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>

    </v-dialog>
    <v-dialog v-model="modify_title" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          修改职称信息
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="职称名"
              v-model="new_name"
              class="col-12"></v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modifyTitle"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modify_title = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_title" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加职称信息
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="职称名"
              v-model="new_name"
              class="col-12"></v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addTitle"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_title = false"
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
  name: "ManageTitle",
  created() {
    this.getTitleList()
  },
  data() {
    return {
      new_name: null,
      selectOBJ: {},
      delete_title: false,
      modify_title: false,
      add_title: false,
      title_id: null,
      titles: [],
      titles_headers: [
        {text: '职称名', align: 'start', sortable: false, value: 'name'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
    }
  },
  methods: {
    getTitleList() {
      this.$axios({
        url: "Title",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get"
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.titles = response.data.titles
      })
    },
    modifyTitle() {
      const formData = new FormData()
      formData.append("title_id", this.selectOBJ.id)
      formData.append("name", this.new_name)
      this.$axios({
        method: "put",
        url: "Title",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.modify_title = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.new_name = null
        this.getTitleList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    deleteTitle() {
      this.$axios({
        url: "Title",
        method: "delete",
        params: {
          title_id: this.selectOBJ.id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.delete_title = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getTitleList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    addTitle() {
      const formData = new FormData()
      formData.append("name", this.new_name)
      this.$axios({
        method: "post",
        url: "Title",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.add_title = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getTitleList()
        this.new_name = null
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
