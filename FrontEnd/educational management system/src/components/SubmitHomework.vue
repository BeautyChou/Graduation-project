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
    </v-card-title>
    <div v-for="(question,index) in questions" :key="index">
      <v-card class="ma-4" :color="question.uploaded?'green':'white'" elevation="5">
        <v-card-title>第{{ index + 1 }}题</v-card-title>
        <v-card-subtitle>分值:{{ question.question_max_score }}分</v-card-subtitle>
        <v-card-text>{{ question.content }}</v-card-text>
        <v-card-actions>
          <v-file-input
            v-model="question.image"
            accept="image/*"
            label="上传文件"
            show-size
            :rules="[v => !!v || '文件必选']"
            @change="uploadHomework(index)"></v-file-input>
        </v-card-actions>
      </v-card>
      <v-divider></v-divider>
    </div>
  </div>
</template>

<script>
export default {
  name: "SubmitHomework",
  created() {
    this.$axios({
      method: "get",
      url: 'SelectedQuestion',
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
      this.questions = response.data.questions
      this.uploaded = response.data.uploaded
      this.homeworkTitle = response.data.questions[0].homework_title
      this.homeworkId = this.$store.state.homeworkId
      this.$store.commit('nextPage')
    })
  },
  data() {
    return {
      questions: [],
      homeworkId: null,
      uploaded: [],
    }
  },
  methods: {
    uploadHomework(id) {
      const formData = new FormData();
      formData.append("image", this.questions[id].image);
      formData.append("student_id", this.$store.state.studentId);
      formData.append("homework_id", this.homeworkId);
      formData.append("question_id", this.questions[id].question_id);
      formData.append("record_id", this.questions[id].record_id);
      this.$axios({
        method: "post",
        url: "http://127.0.0.1:9000/postHomework",
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
        this.$store.commit('setSuccess', "作业提交成功！")
        this.questions[id].uploaded = true;
      });
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
          this.questions = response.data.questions
          this.uploaded = response.data.uploaded
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
            url: 'SelectedQuestion',
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
            this.questions = response.data.questions
            this.uploaded = response.data.uploaded
            this.homeworkTitle = response.data.questions[0].homework_title
            this.homeworkId = this.$store.state.homeworkId
            this.$store.commit('nextPage')
          })
        } else {
          return
        }
      },
      immediate: true
    },
  }
}
</script>

<style scoped>

</style>
