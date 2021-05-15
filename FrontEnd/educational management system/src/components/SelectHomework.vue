<template>
  <v-card>
    <v-card-title>
      <v-btn
        icon
        to="/SelectCourse"
        @click.native="$store.commit('previousPage');">
        <v-icon>mdi-arrow-left-bold</v-icon>
        <span>返回</span>
      </v-btn>
      <div class="ml-8 mr-4">
        选择作业
      </div>
      <span v-if="unchangeable" class="red lighten-3">本课程已经提交成绩，不能再添加、批改、删除作业!</span>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        to="/AddHomework"
        @click.native="$store.commit('addHomework')"
        v-if="$store.state.level===2"
        :disabled="unchangeable">
        <v-icon>mdi-plus-thick</v-icon>
        <span>
          添加作业
        </span>
      </v-btn>
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="homeworks"
      :options.sync="options"
      :server-items-length="total"
      class="elevation-1"
    >
      <template v-slot:item.DeadLine="{ item }">
        <div v-if="item.DeadLine">
          {{ dateFormat(item.DeadLine) }}
        </div>
      </template>
      <template v-slot:item.operation="{ item }">
        <v-tooltip v-if="($store.state.level===2)" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn
              icon
              color="green"
              v-bind="attrs"
              v-on="on"
              to="/canvas"
              @click.native="$store.commit('setHomeworkId',item.id)"
              :disabled=" time < item.DeadLine || unchangeable"
              x-large>
              <v-icon>
                mdi-check-bold
              </v-icon>
            </v-btn>
          </template>
          <span>批改作业</span>
        </v-tooltip>
        <v-tooltip v-if="($store.state.level===2)" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn
              icon
              color="primary"
              v-bind="attrs"
              v-on="on"
              x-large
              to="/AddHomework"
              @click.native="$store.commit('setHomeworkId',item.id);$store.commit('modifyHomework')"
            >
              <v-icon>
                mdi-text-box-search-outline
              </v-icon>
            </v-btn>
          </template>
          <span>查看作业内容</span>
        </v-tooltip>

        <v-tooltip v-if="($store.state.level===2)" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn
              icon
              color="error"
              v-bind="attrs"
              v-on="on"
              x-large
              @click="dialog=true;selectId = item.id"
              :disabled="unchangeable">
              <v-icon>
                mdi-delete
              </v-icon>
            </v-btn>
          </template>
          <span>删除作业</span>
        </v-tooltip>
        <v-dialog
          v-if="($store.state.level===2)"
          v-model="dialog"
          width="500"
          persistent>
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
        <v-tooltip
          v-if="($store.state.level===1)"
          bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn
              icon
              color="primary"
              v-bind="attrs"
              v-on="on"
              :disabled=" time > item.DeadLine "
              x-large
              to="/SubmitHomework"
              @click.native="$store.commit('setHomeworkId',item.id)">
              <v-icon>
                mdi-square-edit-outline
              </v-icon>
            </v-btn>
          </template>
          <span>写作业</span>
        </v-tooltip>
        <v-tooltip
          v-if="($store.state.level===1)"
          bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn
              icon
              color="green"
              to="/CheckReview"
              v-bind="attrs"
              v-on="on"
              x-large
              @click.native="$store.commit('setHomeworkId',item.id)">
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
      unchangeable: false,
      options: {},
      total: null,
      selectId: null,
      alert: false,
      dialog: false,
      time: null,
      homeworkBackup: [{'homework_title': '暂时还没有布置作业'}],
      homeworks: [],
      headers: [
        {text: '作业名', align: 'start', sortable: false, value: 'homework_title'},
        {text: '截止时间', sortable: false, value: 'DeadLine'},
        //学年由上传年份决定
        {text: '操作', sortable: false, value: 'operation'},
      ],
      loadingFlag: true,
      noHomeworkFlag: true,
    }
  },
  methods: {
    deleteHomework(id) {
      this.dialog = false
      this.$axios({
        method: "DELETE",
        url: "Homework?id=" + id,
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.homeworks.some((item, i) => {
          if (item.id === id) {
            this.homeworks.splice(i, 1)
          }
          this.$store.commit('setSuccess', "删除作业成功!")
          this.setTimeout(() => {
            this.$store.commit('closeSuccess')
          }, 3000)
        })
      })
    },
    getHomework() {
      this.$axios({
        method: "get",
        url: 'Homework',
        params: {
          'course_id': this.$store.state.courseId,
          'record_id': this.$store.state.recordId,
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
        if (response.data.homeworks === null) {
          this.homeworks = this.homeworkBackup
          this.noHomeworkFlag = false
        } else {
          this.homeworks = response.data.homeworks
          this.noHomeworkFlag = true
        }
        this.unchangeable = response.data.unchangeable
        this.total = response.data.total
        this.time = response.data.time
        this.loadingFlag = false
      })
    }
  },
  created() {
    this.getHomework()
  },
  watch: {
    '$store.state.courseId': {
      handler(newValue, oldValue) {
        this.$axios({
          method: "get",
          url: 'Homework',
          params: {
            'course_id': newValue,
            'record_id': this.$store.state.recordId,
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
          if (response.data.homeworks === null) {
            this.homeworks = this.homeworkBackup
            this.noHomeworkFlag = false
          } else {
            this.homeworks = response.data.homeworks
            this.noHomeworkFlag = true
          }
          this.unchangeable = response.data.unchangeable
          this.total = response.data.total
          this.time = response.data.time
          this.loadingFlag = false
        })
      },
    },
    '$store.state.refreshFlag': {
      handler(newValue, oldValue) {
        if (newValue < 1) return
        this.getHomework()
      },
    },
    options: {
      handler() {
        this.getHomework()
      },
      deep: true,
    },
  },
  computed: {
    dateFormat() {
      return function (str) {
        return str.substring(0, 10) + '　' + str.substring(11, 19)
      }
    },
  }
}
</script>

<style scoped>

</style>
