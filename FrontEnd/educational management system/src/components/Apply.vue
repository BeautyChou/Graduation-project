<template>
  <v-card>
    <v-card-title>申请列表</v-card-title>
    <v-expansion-panels>
      <v-expansion-panel
        v-for="item in applies"
        :key="item.id"
        >
        <v-expansion-panel-header :color="item.result===0?'white':item.result===1?'green':'red'">
          教师：{{ item.name }} 课程：{{ item.course_name }}
          <v-spacer></v-spacer>
          <v-tooltip v-if="$store.state.level===2||true" bottom>
            <template v-slot:activator="{ on,attrs }">
              <v-btn
                icon
                color="green"
                v-bind="attrs"
                v-on="on"
                @click.native="pass(item)"
                small
                class="col-1"
                :disabled="item.result!==0">
                <v-icon>
                  mdi-check-bold
                </v-icon>
              </v-btn>
            </template>
            <span>通过</span>
          </v-tooltip>
          <v-tooltip v-if="$store.state.level===2||true" bottom>
            <template v-slot:activator="{ on,attrs }">
              <v-btn
                icon
                color="red"
                v-bind="attrs"
                v-on="on"
                @click.native="denied(item)"
                small
                class="col-1"
                :disabled="item.result!==0">
                <v-icon>
                  mdi-close-thick
                </v-icon>
              </v-btn>
            </template>
            <span>不通过</span>
          </v-tooltip>
        </v-expansion-panel-header>
        <v-expansion-panel-content class="mt-6">
          <v-row class="font-weight-bold">
            <v-select
              disabled
              v-model="item.before_week_time"
              :items="week_time"
              item-text="name"
              item-value="value"
              label="上课时间"
              outlined
              class="col-4"></v-select>
            <v-select
              disabled
              outlined
              v-model="item.before_start_time"
              :items="start_time"
              item-text="name"
              item-value="value"
              label="起始节次"
              class="col-4"></v-select>
            <v-select
              disabled
              outlined
              v-model="item.before_end_time"
              :items="end_time"
              item-text="name"
              item-value="value"
              label="结束节次"
              class="col-4"></v-select>
            <v-select
              disabled
              v-model="item.before_classroom_id"
              :items="classrooms"
              item-text="name"
              item-value="value"
              outlined
              label="上课教室"
              class="col-4"></v-select>
            <v-select
              disabled
              outlined
              v-model="item.before_start_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="开始周"
              class="col-4"></v-select>
            <v-select
              disabled
              outlined
              v-model="item.before_end_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="结束周"
              class="col-4"></v-select>
          </v-row>
          <v-row class="my-0 font-weight-bold"><span class="mx-auto my-0"><v-icon x-large>mdi-arrow-down-bold</v-icon>修改为</span></v-row>
          <v-row class="font-weight-bold">
            <v-select
              disabled
              v-model="item.after_week_time"
              :items="week_time"
              item-text="name"
              item-value="value"
              label="上课时间"
              outlined
              class="col-4"></v-select>
            <v-select
              disabled
              class="col-4"
              outlined
              v-model="item.after_start_time"
              :items="start_time"
              item-text="name"
              item-value="value"
              label="起始节次"></v-select>
            <v-select
              disabled
              class="col-4"
              outlined
              v-model="item.after_end_time"
              :items="end_time"
              item-text="name"
              item-value="value"
              label="结束节次"></v-select>
            <v-select
              disabled
              class="col-4"
              v-model="item.after_classroom_id"
              :items="classrooms"
              item-text="name"
              item-value="value"
              outlined
              label="上课教室"
            ></v-select>
            <v-select
              disabled
              class="col-4"
              outlined
              v-model="item.after_start_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="开始周"></v-select>
            <v-select
              disabled
              class="col-4"
              outlined
              v-model="item.after_end_week"
              :items="weeks"
              item-text="name"
              item-value="value"
              label="结束周"></v-select>
          </v-row>
          <v-row class="font-weight-bold">
            <v-textarea outlined disabled v-model="item.reason"></v-textarea>
          </v-row>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-card>
</template>

<script>
export default {
  name: "Apply",
  data() {
    return {
      applies: [],
      classrooms: [],
      start_time: [
        {name: "第一节", value: 1},
        {name: "第二节", value: 2},
        {name: "第三节", value: 3},
        {name: "第四节", value: 4},
        {name: "第五节", value: 5},
        {name: "第六节", value: 6},
        {name: "第七节", value: 7},
        {name: "第八节", value: 8},
        {name: "第九节", value: 9},
        {name: "第十节", value: 10},
        {name: "第十一节", value: 11},
        {name: "第十二节", value: 12},
        {name: "第十三节", value: 13},
        {name: "第十四节", value: 14},
      ],
      end_time: [
        {name: "第一节", value: 1},
        {name: "第二节", value: 2},
        {name: "第三节", value: 3},
        {name: "第四节", value: 4},
        {name: "第五节", value: 5},
        {name: "第六节", value: 6},
        {name: "第七节", value: 7},
        {name: "第八节", value: 8},
        {name: "第九节", value: 9},
        {name: "第十节", value: 10},
        {name: "第十一节", value: 11},
        {name: "第十二节", value: 12},
        {name: "第十三节", value: 13},
        {name: "第十四节", value: 14},
      ],
      weeks: [
        {name: "第一周", value: 1},
        {name: "第二周", value: 2},
        {name: "第三周", value: 3},
        {name: "第四周", value: 4},
        {name: "第五周", value: 5},
        {name: "第六周", value: 6},
        {name: "第七周", value: 7},
        {name: "第八周", value: 8},
        {name: "第九周", value: 9},
        {name: "第十周", value: 10},
        {name: "第十一周", value: 11},
        {name: "第十二周", value: 12},
        {name: "第十三周", value: 13},
        {name: "第十四周", value: 14},
        {name: "第十五周", value: 15},
        {name: "第十六周", value: 16},
        {name: "第十七周", value: 17},
        {name: "第十八周", value: 18},
        {name: "第十九周", value: 19},
        {name: "第二十周", value: 20},
      ],
      week_time: [
        {name: "星期一", value: 1},
        {name: "星期二", value: 2},
        {name: "星期三", value: 3},
        {name: "星期四", value: 4},
        {name: "星期五", value: 5},
        {name: "星期六", value: 6},
        {name: "星期日", value: 7},
      ],
    }
  },
  methods: {
    pass(applyObj) {
      const formdata = new FormData()
      formdata.append('result', 1)
      formdata.append('id', applyObj.id)
      this.$axios({
        method: "PUT",
        url: "Apply",
        data: formdata,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.applies.some((apply, index) => {
          if (apply.id === applyObj.id) {
            apply.result = 1
          }
        })
      })
    },
    denied(applyObj) {
      const formdata = new FormData()
      formdata.append('result', 2)
      formdata.append('id', applyObj.id)
      this.$axios({
        method: "PUT",
        url: "Apply",
        data: formdata,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        this.applies.some((apply, index) => {
          if (apply.id === applyObj.id) {
            apply.result = 2
          }
        })
      })
    }
  },
  created() {
    this.$axios({
      method: "get",
      url: "Apply?level=" + this.$store.state.level + '&teacher_id=' + this.$store.state.teacherId,
      headers:{
        'Token': "8a54sh " + this.$store.state.Jwt
      }
    }).then((response) => {
      if (response.data.msg === "Token无效") {
        this.$emit('func')
        return
      }
      this.classrooms = response.data.classrooms
      this.applies = response.data.applies
      console.log(response)
    })
  }
}
</script>

<style scoped>

</style>
