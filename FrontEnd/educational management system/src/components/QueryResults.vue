<template>
  <v-card>
    <v-card-title>查看分数</v-card-title>
    <v-data-table
      :headers="headers"
      :items="electives"
      class="elevation-1"
    >
      <template v-slot:item.total_score="{ item }">
        <span>{{ totalScore(item) }}</span>
      </template>
      <template v-slot:item.grade_point="{ item }">
        <span>{{ gradePoint(item) }}</span>
      </template>
      <template v-slot:item.operation="{ item }">
        <v-tooltip bottom>
          <template v-slot:activator="{ on,attrs }">
            <v-btn icon color="primary" v-bind="attrs" v-on="on" x-large
                   @click="dialog = true;selectOBJ = item;">
              <v-icon>
                mdi-format-align-justify
              </v-icon>
            </v-btn>
          </template>
          <span>查看分数明细</span>
        </v-tooltip>
      </template>
    </v-data-table>
    <v-dialog v-model="dialog" width="500">
      <v-card>
        <v-card-title class="font-weight-bold py">分数详情</v-card-title>
        <v-card-text class="font-weight-bold mx-auto text-center">
          作业分数:{{ selectOBJ.homework_score }} 占比:{{
            selectOBJ.homework_percentage * (100 - selectOBJ.percentage) / 100
          }}%
        </v-card-text>
        <v-card-text class="font-weight-bold mx-auto text-center">
          平时表现分数:{{ selectOBJ.behavior_score }}
          占比:{{ (100 - selectOBJ.homework_percentage) * (100 - selectOBJ.percentage) / 100 }}%
        </v-card-text>
        <v-card-text class="font-weight-bold mx-auto text-center">
          期末分数:{{ selectOBJ.test_score }} 占比:{{ selectOBJ.percentage }}%
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
import Timetables from "timetables";

export default {
  name: "QueryResults",
  created() {
    this.$axios({
      method: "get",
      url: "http://127.0.0.1:9000/getScore",
      params: {
        student_id: this.$store.state.studentId,
      },
      headers:{
        'Token': "8a54sh " + this.$store.state.Jwt
      }
    }).then((response) => {
      if (response.data.msg === "Token无效") {
        this.$emit('func')
        return
      }
      this.electives = response.data.electives
      console.log(response)
    })
  },
  data() {
    return {
      dialog: false,
      selectOBJ: {},
      electives: [],
      headers: [
        {text: '课程名', align: 'start', sortable: false, value: 'course_name'},
        {text: '总分', sortable: false, value: 'total_score'},
        {text: '绩点', sortable: false, value: 'grade_point'},
        {text: '操作', sortable: false, value: 'operation'},
      ],
    }
  },
  computed: {
    totalScore() {
      return function (obj) {
        console.log(obj)
        let homework = (obj.homework_score * obj.homework_percentage * (100 - obj.percentage)) / 10000
        let behavior = (obj.behavior_score * (100 - obj.homework_percentage) * (100 - obj.percentage)) / 10000
        let test = obj.test_score * obj.percentage / 100
        console.log(homework, behavior, test)
        return (homework + behavior + test).toFixed(2)
      }
    },
    gradePoint() {
      return function (obj) {
        let homework = (obj.homework_score * obj.homework_percentage * (100 - obj.percentage)) / 10000
        let behavior = (obj.behavior_score * (100 - obj.homework_percentage) * (100 - obj.percentage)) / 10000
        let test = obj.test_score * obj.percentage / 100
        let total = (homework + behavior + test).toFixed(0)
        let grade_point = 0
        console.log()
        if (total>95) return 4.5
        else if (total>=90&&total<=95) return 4.0;
        else if (total>=85&&total<=89) return 3.5;
        else if (total>=80&&total<=84) return 3.0;
        else if (total>=75&&total<=79) return 2.5;
        else if (total>=70&&total<=74) return 2.0;
        else if (total>=60&&total<=70) return 1.0;
        else return 0
      }
    }
  },
  watch: {
    "$route.path": {
      handler(newVal, oldVal) {
        this.$axios({
          method: "get",
          url: "http://127.0.0.1:9000/getScore",
          params: {
            student_id: this.$store.state.studentId,
          },
          headers:{
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response) => {
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.electives = response.data.electives
          console.log(response)
        })
      },
    }
  }
}
</script>

<style scoped>

</style>
