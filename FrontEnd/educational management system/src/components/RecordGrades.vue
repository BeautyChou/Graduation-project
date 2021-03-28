<template>
  <v-card>
    <v-card-title>
      <v-btn icon to="/SelectCourse" @click.native="$store.commit('previousPage');$store.commit('setRecordId',null);">
        <v-icon>mdi-arrow-left-bold</v-icon>
        <span>返回</span>
      </v-btn>
      <div class="mx-8">录入成绩</div>
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
          </th >
          <th rowspan="2" colspan="1" :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}">
            <span>总分</span>
          </th>
        </tr>
        <tr>
          <th :style="{textAlign:'center',border:'thin solid rgba(0,0,0,.12)'}" >
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
          <td >{{ item.student_name }}</td>
          <td >{{ item.homework_score }}</td>
          <td >
            <v-text-field v-model.number="item.behavior_score"></v-text-field>
          </td>
          <td >
            <v-text-field v-model.number="item.test_score"></v-text-field>
          </td>
          <td >{{ totalScore(item) }}</td>
        </tr>
        </tbody>
      </template>
      <template v-slot:item.operation="{ item }">
        <v-tooltip v-if="$store.state.level===2||true" bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="red" v-bind="attrs" v-on="on" x-large
                   @click="quitCourse(item)">
              <v-icon>
                mdi-exit-to-app
              </v-icon>
            </v-btn>
          </template>
          <span>退出课程</span>
        </v-tooltip>
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
  computed: {
    totalScore() {
      return function (obj) {
        let homework = obj.homework_score * this.homework_percentage * (100 - this.percentage) / 10000
        let behavior = obj.behavior_score * (100 - this.homework_percentage) * (100 - this.percentage) / 10000
        let test = obj.test_score * this.percentage / 100
        return (homework + behavior + test).toFixed(2)
      }
    },
  },
  watch: {
    '$store.state.recordId': {
      handler(newValue, oldValue) {
        this.$axios({
          method: "get",
          url: 'http://127.0.0.1:9000/getStudentScore',
          params: {
            'record_id': newValue,
          }
        }).then((response) => {
          this.scores = response.data.scores
          console.log(response)
        })
      },
    },
  },
}
</script>

<style scoped>
</style>
