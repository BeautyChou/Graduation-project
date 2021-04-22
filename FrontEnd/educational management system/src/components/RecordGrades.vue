<template>
  <v-card>
    <v-card-title>
      <v-btn icon to="/SelectCourse" @click.native="$store.commit('previousPage');$store.commit('setRecordId',null);">
        <v-icon>mdi-arrow-left-bold</v-icon>
        <span>返回</span>
      </v-btn>
      <div class="mx-8">录入成绩  (输入的所有数据仅限正整数!)</div>
      <v-spacer></v-spacer>
      <v-btn color="green" :disabled="uploadFlag()||scores.length===0" @click="postScore">上传成绩</v-btn>
    </v-card-title>
    <v-data-table
      :items="scores"
      :items-per-page="5"
      class="elevation-1"
    >
      <template v-slot:body="{ items }">
        <thead :style="{}">
        <tr>
          <th rowspan="2" colspan="1" :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}">学生姓名</th>
          <th rowspan="1" colspan="2" :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}">
            <span>平时成绩 (占比:{{ 100 - percentage }}%)</span>
          </th>
          <th rowspan="2" colspan="1" :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}">
            <span>考试成绩 (占比:</span>
            <v-col
              cols="1"
              class="d-inline-block ma-0 pa-0"
            >
              <v-text-field dense v-model.number="percentage"></v-text-field>
            </v-col>
            <span>%)</span>
          </th>
          <th rowspan="2" colspan="1" :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}">
            <span>总分</span>
          </th>
        </tr>
        <tr>
          <th :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}">
            <span>作业成绩 (占比:</span>
            <v-col
              cols="1"
              class="d-inline-block ma-0 pa-0"
            >
              <v-text-field dense v-model.number="homework_percentage"></v-text-field>
            </v-col>
            <span>%)</span></th>
          <th :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}">
            <span>平时表现成绩 (占比:{{ 100 - homework_percentage }}%)</span></th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="(item,v) in items" :key="v" class="font-weight-bold text-center">
          <td>{{ item.student_name }}</td>
          <td>{{ item.homework_score }}</td>
          <td>
            <v-text-field v-model.number="item.behavior_score"></v-text-field>
          </td>
          <td>
            <v-text-field v-model.number="item.test_score"></v-text-field>
          </td>
          <td>
            <v-text-field :disabled="true" v-model.number="item.total" :value="totalScore(item,v)"></v-text-field>

          </td>
        </tr>
        </tbody>
      </template>
    </v-data-table>
  </v-card>
</template>

<script>
export default {
  name: "RecordGrades",
  data() {
    return {
      scores: [{homework_score: 30, student_id: 1, student_name: "123"}],
      headers2: [
        {text: '作业成绩', align: "center"},
        {text: '平时表现成绩', align: "center"},
      ],
      headers: [
        {text: '学生姓名', align: 'start', sortable: false, value: 'name'},
        {text: '平时成绩', sortable: false, value: 'score'},
        {text: '考试成绩', sortable: false, value: 'test_score'},
        {text: '总分', sortable: false, value: 'total_score'},
      ],
      homework_percentage: 0,
      percentage: 0,
    }
  },
  methods:{
    postScore(){
      this.scores.forEach((v,i)=>{
        v.record_id = this.$store.state.recordId
        v.course_id = this.$store.state.courseId
        v.percentage = this.percentage
        v.homework_percentage = this.homework_percentage
      })
      let scores = JSON.stringify(this.scores)
      console.log(this.scores)
      this.$axios({
        method:"post",
        url:"Score",
        data:scores,
        headers:{
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response)=>{
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.$store.commit(response.data.snackbar,response.data.msg)
        setTimeout(()=>{
          this.$store.commit(response.data.snackbar2)
        },3000)
      })
    }
  },
  computed: {
    totalScore() {
      return function (obj,v) {
        let homework = obj.homework_score * this.homework_percentage * (100 - this.percentage) / 10000
        let behavior = obj.behavior_score * (100 - this.homework_percentage) * (100 - this.percentage) / 10000
        let test = obj.test_score * this.percentage / 100
        this.$set(this.scores[v],'total',(homework + behavior + test).toFixed(2))
        return (homework + behavior + test).toFixed(2)
      }
    },
    uploadFlag() {
      return function () {
        let flag = false
        this.scores.some((v, i) => {
          console.log(v.total,i)
          if (isNaN(v.total)) {
            flag = true
          }
        })
        console.log(this.scores)
        return flag
      }
    }
  },
  watch: {
    '$store.state.recordId': {
      handler(newValue, oldValue) {
        this.$axios({
          method: "get",
          url: 'StudentScore',
          params: {
            'record_id': newValue,
          },
          headers:{
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response) => {
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.scores = response.data.scores
          console.log(response)
        })
      },
    },
  },
  created() {
    this.$axios({
      method: "get",
      url: 'StudentScore',
      params: {
        'record_id': this.$store.state.recordId,
      },
      headers:{
        'Token': "8a54sh " + this.$store.state.Jwt
      }
    }).then((response) => {
      if (response.data.msg === "Token无效") {
        this.$emit('func')
        return
      }
      this.scores = response.data.scores
      console.log(response)
    })
  },
  beforeRouteLeave(to,from,next) {
    this.$store.commit('setRecordId', null);
    next()
  },
}
</script>

<style scoped>
</style>
