<template>
  <v-app>
    <v-main class="flex">
      <div>
        <v-snackbar v-model="$store.state.successFlag" :timeout="-1" color="green" top>
          <v-icon>
            mdi-check-circle
          </v-icon>
          <span>
        {{ $store.state.successMsg }}
      </span>
        </v-snackbar>
        <v-snackbar v-model="$store.state.errorFlag" :timeout="-1" color="red" top>
          <v-icon>
            mdi-minus-circle
          </v-icon>
          <span>
        {{ $store.state.errorMsg }}
      </span>
        </v-snackbar>

      </div>
    </v-main>

    <v-app id="blur">
      <v-card
        class="mx-auto"
        width="256"
        tile
      >
        <v-navigation-drawer permanent app expand-on-hover :mini-variant.sync="mini">
          <v-system-bar></v-system-bar>
          <v-toolbar-title class="text-center">{{ $store.state.recordId }}</v-toolbar-title>
          <v-list>
            <v-list-item>
              <v-list-item-avatar>
                <v-avatar color="red" size="36">
                  <span class="white--text headline">{{ avatar }}</span>
                </v-avatar>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title class="title">
                  {{ $store.state.username }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <v-list
            nav
          >
            <v-list-item-group
              v-model="selectedItem"
              color="primary"
            >
              <v-list-item to="/UserInfo">
                <v-icon>
                  mdi-account-box
                </v-icon>
                <v-list-item-content >
                  <v-list-item-title>个人信息</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/SelectCourse">
                <v-icon>
                  mdi-playlist-edit
                </v-icon>
                <v-list-item-content>
                  <v-list-item-title>课程管理</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/SetClass">
                <v-icon>
                  mdi-table-cog
                </v-icon>
                <v-list-item-content>
                  <v-list-item-title>安排课程</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/Apply">
                <v-icon>
                  mdi-alert-circle-check
                </v-icon>
                <v-list-item-content>
                  <v-list-item-title>课程申请审核</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/ClassSheet">
                <v-icon>
                  mdi-table-clock
                </v-icon>
                <v-list-item-content>
                  <v-list-item-title>查看课表</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/ChooseCourse">
                <v-icon>
                  mdi-playlist-check
                </v-icon>
                <v-list-item-content >
                  <v-list-item-title>选课</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/ChosenCourse">
                <v-icon>
                  mdi-playlist-remove
                </v-icon>
                <v-list-item-content >
                  <v-list-item-title>查看已选课程/退课</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/QueryResults">
                <v-icon>
                  mdi-file-search
                </v-icon>
                <v-list-item-content >
                  <v-list-item-title>查询成绩</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/Manage">
                <v-icon>
                  mdi-cog
                </v-icon>
                <v-list-item-content >
                  <v-list-item-title>管理</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
          <template v-slot:append>
            <v-list-item-group>
            <v-list-item @click.native="exit">
              <v-icon>
                mdi-exit-to-app
              </v-icon>
              <v-list-item-content >
                <v-list-item-title>退出登录</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            </v-list-item-group>
          </template>
        </v-navigation-drawer>

      </v-card>
      <v-main>
        <router-view v-on:func="jwtInvalid"></router-view>
      </v-main>
    </v-app>
    <div>
      <v-dialog
        persistent
        overlay-opacity="0.92"
        v-model="dialog"
        width="500"
        transition="slide-y-transition"
      >

        <v-card class="text-center">
          <v-card-text class="font-weight-bold title pt-4">欢迎使用本教务系统</v-card-text>
          <v-card-text class="font-weight-bold title">登录</v-card-text>
            <v-avatar color="indigo" size="128">
              <v-icon dark size="80">
                mdi-account-circle
              </v-icon>
            </v-avatar>
          <v-card-text class="py-5">
            <v-text-field v-model="account" label="教务系统账号"></v-text-field>
            <v-text-field v-model="password" label="教务系统密码"></v-text-field>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions>
            <v-btn
              block
              class="mx-auto"
              color="primary"
              text
              @click="accept()"
            >
              登录
            </v-btn>
          </v-card-actions>
        </v-card>

      </v-dialog>
    </div>

  </v-app>
</template>

<script>
import HelloWorld from './components/HelloWorld';
import Canvas from "./components/Canvas";
import ClassSheet from "./components/SetClass";
import SelectCourse from "./components/SelectCourse";
import SelectHomework from "./components/SelectHomework";
import ChooseCourse from "./components/ChooseCourse";
import ChosenCourse from "./components/ChosenCourse"
import UserInfo from "./components/UserInfo";
import Manage from "./components/Manage";

export default {
  name: 'App',
  mounted() {
    if (this.$store.state.Jwt !== null){
      this.$axios({
        url:"isExpire",
        method:"post",
        headers:{
          "Content-Type": "multipart/form-data"
        }
      }).then((response)=>{
        if (response.data.msg === "Token无效" ){
          this.jwtInvalid()
        }else{
          this.dialog = false
          this.$store.commit('jwtValid')
          this.$nextTick(()=> {
            document.getElementById("blur").style.filter = "blur(0px)"
          });
          this.avatar = this.$store.state.username.substr(0,1)
        }
      })
    }
  },
  components: {
    HelloWorld,
    Canvas,
    ClassSheet,
    SelectCourse,
    SelectHomework,
    ChooseCourse,
    ChosenCourse,
    UserInfo,
    Manage
  },

  data: () => ({
    avatar: null,
    account:null,
    password:null,
    dialog:true,
    selectedItem: 0,
    mini: true,
    isActive: true,
  }),
  methods:{
    jwtInvalid() {
      document.getElementById("blur").style.filter = "blur(0px)"
      this.$store.commit('setError','登录失效，请重新登录！')
      this.$router.replace('/HelloWorld')
      this.$store.commit('jwtInvalid')
      setTimeout(()=>{
        this.exit()
      },1000)
    },
    exit(){
      this.dialog = true
      this.$store.commit('clear')
      document.getElementById("blur").style.filter = "blur(10px)"
    },
    accept(){
      const formData = new FormData()
      formData.append("id",this.account)
      formData.append("password",this.password)
      this.$axios({
        url:"login",
        method:"post",
        data:formData,
        headers:{
          "Content-Type": "multipart/form-data"
        }
      }).then((response)=>{
        if (response.data.msg === "success"){
          this.avatar = response.data.username.substr(0,1)
          console.log(response)
          this.$store.commit("setSuccess","登录成功")
          this.$store.commit('jwtValid')
          setTimeout(()=>{
            this.$store.commit("closeSuccess")
          },3000)
          this.$store.commit("setLevel",response.data.level)
          this.$store.commit("setUsername",response.data.username)
          this.$store.commit("setJwt",response.data.data.token)
          if (response.data.level === 1){
            this.$store.commit("setStudentID",response.data.ID)
            this.$store.commit("setFacultyId",response.data.faculty_id)
            this.$store.commit("setSpecialtyId",response.data.specialty_id)
            this.$store.commit("setDirectionId",response.data.direction_id)
          }else if (response.data.level === 2){
            this.$store.commit("setFacultyId",response.data.faculty_id)
            this.$store.commit("setTeacherID",response.data.ID)
          }
          this.dialog = false
          document.getElementById("blur").style.filter = "blur(0px)"
        }else if (response.data.msg === "failed"){
          this.$store.commit("setError","用户名或密码错误！")
          setTimeout(()=>{
            this.$store.commit("closeError")
          },3000)
        }
      })
    }
  }
};
</script>

<style scoped>
#blur {
  -webkit-filter: blur(10px);
}
</style>
