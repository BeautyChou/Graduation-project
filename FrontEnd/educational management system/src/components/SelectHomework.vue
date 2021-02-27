<template>
  <v-card>
    <v-alert
      :value="alert"
      type="success"
      dark
      icon="mdi-check-circle"
      transition="scroll-y-transition"
    >
      删除作业成功!
    </v-alert>
    <v-card-title>选择作业</v-card-title>
    <v-data-table
      :headers="headers"
      :items="homeworks"
      :items-per-page="5"
      class="elevation-1"
    >
      <template v-slot:item.DeadLine="{ item }">
        {{ item.DeadLine | dateFormat }}
      </template>
      <template v-slot:item.operation="{ item }">
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="green" v-bind="attrs" v-on="on" to="/canvas" @click.native="$store.commit('setHomeworkId',item.id)" x-large>
              <v-icon>
                mdi-check-bold
              </v-icon>
            </v-btn>
          </template>
          <span>批改作业</span>
        </v-tooltip>
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large>
              <v-icon>
                mdi-lead-pencil
                <router-link to="/SelectHomework" @click.native="$store.commit('setCourseId',item.course_id)">修改作业内容
                </router-link>
              </v-icon>
            </v-btn>
          </template>
          <span>修改作业内容</span>
        </v-tooltip>
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="error" v-bind="attrs" v-on="on" x-large @click="dialog=true;selectId = item.id">
              <v-icon>
                mdi-delete
              </v-icon>
            </v-btn>
          </template>
          <span>删除作业</span>
        </v-tooltip>
        <v-dialog v-if="$store.state.level===2||true" v-model="dialog" width="500" persistent>
          <v-card>
            <v-card-title class="headline font-weight-bold">
              警告!
            </v-card-title>
            <v-card-text class="font-weight-bold">
              真的要删除本次作业吗?这个操作是不可逆转的!
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                class="font-weight-bold"
                color="green darken-1"
                text
                @click="deleteHomework(selectId)"
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
        <v-tooltip v-if="$store.state.level===1||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" :disabled=" time > item.DeadLine " x-large>
              <v-icon>
                mdi-square-edit-outline
                <router-link to="/SelectHomework" @click.native="$store.commit('setCourseId',item.course_id)">修改作业内容
                </router-link>
              </v-icon>
            </v-btn>
          </template>
          <span>写作业</span>
        </v-tooltip>
        <v-tooltip v-if="$store.state.level===1||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="green" v-bind="attrs" v-on="on" x-large>
              <v-icon>
                mdi-file-check
              </v-icon>
            </v-btn>
          </template>
          <span>查看老师批注</span>
        </v-tooltip>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
export default {
  name: "SelectHomework",
  data() {
    return {
      selectId: null,
      alert: false,
      dialog: false,
      time: null,
      homeworks: [],
      headers: [
        {text: '作业名', align: 'start', sortable: false, value: 'homework_title'},
        {text: '截止时间', sortable: false, value: 'DeadLine'},
        //学年由上传年份决定
        {text: '操作', sortable: false, value: 'operation'},
      ],
    }
  },
  methods: {
    deleteHomework(id) {
      this.dialog = false
      this.$axios({
        method: "DELETE",
        url: "http://127.0.0.1:9000/deleteHomework?id=" + id,
      }).then((response) => {
        this.homeworks.some((item, i) => {
          if (item.id === id) {
            this.homeworks.splice(i, 1)
          }
        })
        this.alert = true
        setTimeout(() => {
          this.alert = false
        }, 3000)
      })
    }
  },
  watch: {
    '$store.state.courseId': {
      handler(newValue, oldValue) {
        this.$axios({
          method: "get",
          url: 'http://127.0.0.1:9000/getHomeworkList',
          params: {
            'course_id': newValue,
          }
        }).then((response) => {
          this.homeworks = response.data.homeworks
          this.time = response.data.time
          console.log(response)
        })
      },
      immediate: true
    }
  },
  filters: {
    dateFormat(str) {
      var date = str.substring(0, 10)
      var time = str.substring(11, 19)
      return `${date}` + '　' + `${time}`
    }
  },
}
</script>

<style scoped>

</style>
