<template>
  <div>
    <v-card-title>
      <v-btn icon to="/SelectHomework" @click.native="$store.commit('previousPage')">
        <v-icon>mdi-arrow-left-bold</v-icon>
        <span>返回</span>
      </v-btn>
      <div class="mx-8">
        选择作业
      </div>
      <div class="mx-4">你的总分:<span :class="colorSelect(sum,100)+' lighten-1'">{{ sum }}/100</span></div>
      <div v-if="unReviewedFlag" class="warning mx-4">老师还没有批改完你的作业，该分数为目前老师已批改题目的分数！！！！</div>
    </v-card-title>
    <div v-for="(review,index) in reviews" :key="index">
      <v-card class="ma-4" color="white" elevation="5">
        <v-card-title>第{{ index + 1 }}题</v-card-title>
        <v-card-subtitle class="title">
          <span :class="review.is_scored?colorSelect(review.score,review.question_max_score):' '+'lighten-1 d-inline'">分数:{{ review.score }}/{{
              review.question_max_score
            }}分
          </span>
        </v-card-subtitle>
        <v-card-text>{{ review.content }}</v-card-text>
        <v-card-actions>
          <v-expansion-panels :disabled="!review.is_scored">
            <v-expansion-panel>
              <v-expansion-panel-header v-text="review.is_scored?'查看老师批改':'老师还未批改这道题，请稍后再来查看'">

              </v-expansion-panel-header>
              <v-expansion-panel-content>
                <v-img
                  :src="'http://127.0.0.1:9000/imageChecked?homework='+$store.state.homeworkId+'&question='+review.question_id+'&student='+$store.state.studentId"
                  alt="老师还没有批改这道题"></v-img>
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-card-actions>
      </v-card>
      <v-divider></v-divider>
    </div>
  </div>
</template>

<script>
export default {
  name: "CheckReview",
  data() {
    return {
      reviews: [],
      homeworkId: null,
      sum: 0,
      unReviewedFlag: false,
    }
  },
  methods: {
    colorSelect(score, maxScore) {
      var r = score * 100 / maxScore
      if (r >= 80) return 'green'
      else if (r >= 60) return 'yellow'
      else return 'red'
    }
  },
  created() {
    this.$store.commit('nextPage')
    this.$axios({
      method: "get",
      url: 'http://127.0.0.1:9000/getReviewList',
      params: {
        'homework_id': this.$store.state.homeworkId,
        'student_id': this.$store.state.studentId,
      },
      headers: {
        'Token': "8a54sh " + this.$store.state.Jwt
      }
    }).then((response) => {
      if (response.data.msg === "Token无效") {
        this.$emit('func')
        return
      }
      this.reviews = response.data.reviews
      this.sum = 0
      this.reviews.forEach((item, i) => {
        this.sum += item.score
        if (item.is_scored === false) this.unReviewedFlag = true
      })
      this.homeworkTitle = response.data.questions[0].homework_title
      this.homeworkId = this.$store.state.homeworkId
    })
  },
  watch: {
    '$store.state.homeworkId': {
      handler(newValue, oldValue) {
        this.$axios({
          method: "get",
          url: 'http://127.0.0.1:9000/getReviewList',
          params: {
            'homework_id': newValue,
            'student_id': this.$store.state.studentId,
          },
          headers: {
            'Token': "8a54sh " + this.$store.state.Jwt
          }
        }).then((response) => {
          if (response.data.msg === "Token无效") {
            this.$emit('func')
            return
          }
          this.reviews = response.data.reviews
          this.sum = 0
          response.data.reviews.forEach((item, i) => {
            this.sum += item.score
            if (item.is_scored === false) this.unReviewedFlag = true
          })
          this.homeworkTitle = response.data.questions[0].homework_title
          this.homeworkId = newValue
        })
      },
    },
    '$store.state.refreshFlag': {
      handler(newValue, oldValue) {
        if (newValue === 1) {
          this.$store.commit('nextPage')
          this.$axios({
            method: "get",
            url: 'http://127.0.0.1:9000/getReviewList',
            params: {
              'homework_id': this.$store.state.homeworkId,
              'student_id': this.$store.state.studentId,
            },
            headers: {
              'Token': "8a54sh " + this.$store.state.Jwt
            }
          }).then((response) => {
            if (response.data.msg === "Token无效") {
              this.$emit('func')
              return
            }
            this.reviews = response.data.reviews
            this.sum = 0
            response.data.reviews.forEach((item, i) => {
              this.sum += item.score
              if (item.is_scored === false) this.unReviewedFlag = true
            })
            this.homeworkTitle = response.data.questions[0].homework_title
            this.homeworkId = this.$store.state.homeworkId
          })
        } else {
          return
        }
      },
      immediate: true
    },
  },
}
</script>

<style scoped>

</style>
