<template>
  <v-card>


    <v-card-title>
      <v-btn icon to="/SelectHomework"
             @click.native="$store.commit('previousPage');$store.commit('setHomeworkId',null)">
        <v-icon>mdi-arrow-left-bold</v-icon>
        <span>返回</span>
      </v-btn>

      <div class="mx-8" v-if="!$store.state.modifyHomeworkFlag">添加作业</div>
      <div class="mx-8" v-if="$store.state.modifyHomeworkFlag">查看作业</div>

      <v-menu
        ref="menu"
        v-model="menu"
        :close-on-content-click="false"
        transition="scale-transition"
        offset-y
        min-width="auto">
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
            v-model="deadlineDate"
            label="截止日期"
            prepend-icon="mdi-calendar"
            v-bind="attrs"
            v-on="on"
            :disabled="$store.state.modifyHomeworkFlag"
          ></v-text-field>
        </template>
        <v-date-picker
          v-model="deadlineDate"
          @input="menu1 = false"
          color="green lighten-1"
          :allowed-dates="allowedDates"
        ></v-date-picker>
      </v-menu>

      <v-menu
        ref="menu1"
        v-model="menu1"
        :close-on-content-click="false"
        :nudge-right="40"
        :return-value.sync="deadlineTime"
        transition="scale-transition"
        offset-y
        max-width="290px"
        min-width="290px"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-text-field
            v-model="deadlineTime"
            label="截止时间"
            prepend-icon="mdi-clock-time-four-outline"
            v-bind="attrs"
            v-on="on"
            :disabled="$store.state.modifyHomeworkFlag"
          ></v-text-field>
        </template>
        <v-time-picker
          v-if="menu1"
          v-model="deadlineTime"
          full-width
          scrollable
          format="24hr"
          @click:minute="$refs.menu1.save(deadlineTime)"
          color="green lighten-1"
        ></v-time-picker>
      </v-menu>

      <v-text-field
        v-model="homeworkTitle"
        prepend-icon="mdi-format-align-justify"
        label="作业标题"
        :disabled="$store.state.modifyHomeworkFlag"></v-text-field>


      <v-spacer></v-spacer>

      <v-btn
        :disabled="total!==100||deadlineDate===null||deadlineTime===null||homeworkTitle===null||$store.state.modifyHomeworkFlag"
        color="green"
        class="mr-2"
        @click.native="createHomework();"
      >

        <v-icon>mdi-upload</v-icon>
        <span>
          添加作业
        </span>
      </v-btn>
      <v-menu offset-y>
        <template v-slot:activator="{ on,attrs }">
          <v-btn color="primary" v-bind="attrs" v-on="on" :disabled="$store.state.modifyHomeworkFlag">
            <v-icon>mdi-plus-thick</v-icon>
            <span>
          添加题目
        </span>
          </v-btn>
        </template>
        <v-list>
          <v-list-item @click.stop="alreadyQuestion = true;getContentList();">
            <v-list-item-title>使用已有题目</v-list-item-title>
          </v-list-item>
          <v-list-item @click.stop="newQuestion = true">
            <v-list-item-title>创建新题目</v-list-item-title>
          </v-list-item>
        </v-list>
        <v-dialog v-model="alreadyQuestion" width="500" persistent>
          <v-card>
            <v-card-title class="headline font-weight-bold">
              选择已有题目
            </v-card-title>
            <v-card-text>
              <v-autocomplete
                v-model="selectQuestion"
                :items="contentList"
                item-text="content"
                item-value="id"
                return-object
              ></v-autocomplete>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                class="font-weight-bold"
                color="green darken-1"
                text
                @click="addQuestionToHomework()"
              >
                是
              </v-btn>
              <v-btn
                class="font-weight-bold"
                color="green darken-1"
                text
                @click="alreadyQuestion = false"
              >
                否
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <v-dialog v-model="newQuestion" width="500" persistent>
          <v-card>
            <v-card-title class="headline font-weight-bold">
              添加新题目
            </v-card-title>
            <v-card-text class="font-weight-bold">
              <v-textarea outlined label="题目" v-model="content"></v-textarea>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                class="font-weight-bold"
                color="green darken-1"
                text
                @click="submitQuestion()"
              >
                提交
              </v-btn>
              <v-btn
                class="font-weight-bold"
                color="green darken-1"
                text
                @click="newQuestion = false"
              >
                取消
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-menu>


    </v-card-title>

    <v-data-table
      :headers="headers"
      :items="questions"
      :options.sync="options"
      :server-items-length="totals"
      class="elevation-1"
      calculate-widths
    >
      <template v-slot:item.content="{ item }">
        <td>{{ item.content }}</td>
      </template>
      <template v-slot:item.question_max_score="{ item }">
        <td>
          <v-text-field
            hide-details
            v-model.number="item.question_max_score"
            :disabled="$store.state.modifyHomeworkFlag"></v-text-field>
        </td>
      </template>


      <template v-slot:item.operation="{ item }">
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="error" v-bind="attrs" v-on="on" x-large @click="deleteQuestion(item.question_id);"
                   :disabled="$store.state.modifyHomeworkFlag">
              <v-icon>
                mdi-delete
              </v-icon>
            </v-btn>
          </template>
          <span>删除题目</span>
        </v-tooltip>
      </template>

    </v-data-table>
  </v-card>
</template>

<script>
export default {
  name: "AddHomework",
  created() {
    this.getQuestionList()
    this.getContentList()
  },
  data() {
    return {
      options: {},
      totals: 0,
      successText: null,
      errorText: null,
      homeworkTitle: null,
      errorStatus: false,
      menu: false,
      menu1: false,
      deadlineTime: null,
      deadlineDate: null,
      selectQuestion: null,
      content: null,
      alreadyQuestion: false,
      newQuestion: false,
      newQuestionStatus: false,
      headers: [
        {text: '题目', align: 'start', sortable: false, value: 'content', width: '60%'},
        {text: '分值(仅限正整数)(%)', sortable: false, value: 'question_max_score', width: '30%'},
        {text: '操作', sortable: false, value: 'operation', width: '10%'},
      ],
      questions: [],
      contentList: null,
      question: null,
    }
  },
  methods: {
    getQuestionList() {
      this.$axios({
        method: "get",
        url: 'SelectedQuestion',
        params: {
          'homework_id': this.$store.state.homeworkId,
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
        if (response.data.questions !== null) {
          this.questions = response.data.questions
          this.deadlineTime = response.data.questions[0].DeadLine.substr(11, 5)
          this.deadlineDate = response.data.questions[0].DeadLine.substr(0, 10)
          this.homeworkTitle = response.data.questions[0].homework_title
          this.totals = response.data.total
        }
      })
    },
    allowedDates(val) {
      var date = new Date();
      var y = date.getFullYear();
      var m = date.getMonth() + 1;
      var d = date.getDate()
      if (m < 10) {
        m = '0' + m
      }
      if (d < 10) {
        d = '0' + d
      }
      var dateStr = y + '-' + m + '-' + d
      return val >= dateStr
    },
    createHomework() {
      this.$axios({
        method: "get",
        url: "http://127.0.0.1:9000/getNewHomeworkID",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        for (let i = 0; i < this.questions.length; i++) {
          this.questions[i].id = response.data.id
          this.questions[i].deadline = this.deadlineDate + 'T' + this.deadlineTime + ':00+08:00'
          this.questions[i].homework_title = this.homeworkTitle
          this.questions[i].course_id = this.$store.state.courseId
          this.questions[i].record_id = this.$store.state.recordId
        }
        var questions = JSON.stringify(this.questions)
        this.$axios({
          method: "post",
          url: "Homework",
          data: questions,
          headers: {
            "Content-Type": "multipart/form-data",
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response) => {
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.$store.commit('previousPage')
          this.$store.commit("setSuccess", "作业布置成功！")
          setTimeout(() => {
            this.$store.commit("closeSuccess")
          }, 3000)
          this.$router.push({path: '/selectHomework'})
        })
          .catch((res) => {
          })
      })


    },
    addQuestionToHomework() {
      let flag = false;
      this.questions.some((item, i) => {
        if (item.question_id === this.selectQuestion.id) {
          flag = true
        }
      })
      if (flag) {
        this.$store.commit("setError", '你的作业中已经有这道题了！')
        setTimeout(() => {
          this.$store.commit("closeError")
        }, 3000)
        return
      }
      this.questions.push({
        question_id: this.selectQuestion.id,
        content: this.selectQuestion.content,
        question_max_score: 0
      })
      this.totals++
      this.selectQuestion = null
      this.alreadyQuestion = false
      this.$store.commit("setSuccess", "添加至作业成功！")
      setTimeout(() => {
        this.$store.commit("closeSuccess")
      }, 3000)
    },
    getContentList() {
      this.$axios({
        method: "get",
        url: 'Question',
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.contentList = response.data.questions
      })
    },
    deleteQuestion(id) {
      this.questions.some((item, i) => {
        if (item.question_id === id) {
          this.questions.splice(i, 1)
        }
      })
    },
    submitQuestion() {
      const formData = new FormData();
      formData.append("content", this.content)
      this.$axios({
        method: 'post',
        url: 'Question',
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
        this.newQuestion = false
        this.$store.commit("setSuccess", "添加题目成功！")
        setTimeout(() => {
          this.$store.commit("closeSuccess")
        }, 3000)
        this.questions.push({question_id: response.data.id, content: this.content, question_max_score: 0})
        this.totals++
        this.content = null
      })
    },
  },
  computed: {
    'total': function () {
      var sum = 0
      for (let i = 0; i < this.questions.length; i++) {
        if (this.questions[i].question_max_score <= 0) return (sum = -1)
        else sum += this.questions[i].question_max_score
      }
      return sum
    },
  },
  watch: {
    '$store.state.homeworkId': {
      handler(newValue, oldValue) {
        this.$axios({
          method: "get",
          url: 'SelectedQuestion',
          params: {
            'homework_id': newValue,
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
          if (response.data.questions !== null) {
            this.questions = response.data.questions
            this.deadlineTime = response.data.questions[0].DeadLine.substr(11, 5)
            this.deadlineDate = response.data.questions[0].DeadLine.substr(0, 10)
            this.homeworkTitle = response.data.questions[0].homework_title
            this.totals = response.data.total
          }
        })
        this.getContentList()
      }
    },
    '$store.state.refreshFlag': {
      handler(newValue, oldValue) {
        if (this.$store.state.modifyHomeworkFlag !== false) return
        if (newValue === 1) this.getQuestionList()
      },
      immediate: true
    },
    options: {
      handler() {
        this.getQuestionList()
      },
      deep: true,
    }
  }
}
</script>

<style scoped>

</style>
