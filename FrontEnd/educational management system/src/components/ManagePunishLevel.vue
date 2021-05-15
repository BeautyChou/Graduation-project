<template>
  <v-card>
    <v-card-title>
      <v-row>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          dark
          class="mb-2"
          @click="add_punish_level = true"
        >
          <v-icon>
            mdi-plus-thick
          </v-icon>
          添加处分等级
        </v-btn>
      </v-row>
    </v-card-title>
    <v-data-table
      :headers="punish_levels_headers"
      :items="punish_levels"
      hide-default-footer
    >
      <template v-slot:item.operation="{ item }">

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="modify_punish_level = true ;selectOBJ = item">
              <v-icon>
                mdi-lead-pencil
              </v-icon>
            </v-btn>
          </template>
          <span>修改处分等级名</span>
        </v-tooltip>

        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red darken-4" v-bind="attrs" v-on="on" x-large
                   @click="delete_punish_level = true;selectOBJ = item">
              <v-icon>
                mdi-account-remove
              </v-icon>
            </v-btn>
          </template>
          <span>删除处分等级</span>
        </v-tooltip>

      </template>
    </v-data-table>
    <v-dialog v-if="($store.state.level===2||true)" v-model="delete_punish_level" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          警告!
        </v-card-title>
        <v-card-text class="font-weight-bold">
          真的要删除这个处分等级吗?这个操作是不可逆转的!
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deletePunishLevel"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_punish_level = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>

    </v-dialog>
    <v-dialog v-model="modify_punish_level" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          修改处分等级信息
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="处分等级名"
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
            @click="modifyPunishLevel"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="modify_punish_level = false"
          >
            取消
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_punish_level" width="900" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加处分等级信息
        </v-card-title>
        <v-card-text class="font-weight-bold">
          <v-row class="mt-2">
            <v-text-field
              outlined
              label="处分等级名"
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
            @click="addPunishLevel"
          >
            提交
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_punish_level = false"
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
  name: "ManagePunishLevel",
  created() {
    this.getPunishLevelList()
  },
  data() {
    return {
      new_name: null,
      selectOBJ: {},
      delete_punish_level: false,
      modify_punish_level: false,
      add_punish_level: false,
      level_id: null,
      punish_levels: [],
      punish_levels_headers: [
        {text: '处分等级名', align: 'start', sortable: false, value: 'level'},
        {text: '操作', sortable: false, value: 'operation'}
      ],
    }
  },
  methods: {
    getPunishLevelList() {
      this.$axios({
        url: "PunishLevel",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get"
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.punish_levels = response.data.punish_levels
      })
    },
    modifyPunishLevel() {
      const formData = new FormData()
      formData.append("level_id", this.selectOBJ.id)
      formData.append("name", this.new_name)
      this.$axios({
        method: "put",
        url: "PunishLevel",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.modify_punish_level = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.new_name = null
        this.getPunishLevelList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    deletePunishLevel() {
      this.$axios({
        url: "PunishLevel",
        method: "delete",
        params: {
          level_id: this.selectOBJ.id
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.delete_punish_level = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getPunishLevelList()
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    addPunishLevel() {
      const formData = new FormData()
      formData.append("name", this.new_name)
      this.$axios({
        method: "post",
        url: "PunishLevel",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        data: formData
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.add_punish_level = false
        this.$store.commit(response.data.snackbar, response.data.msg)
        this.getPunishLevelList()
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
