<template>
  <v-app>
    <v-app id="blur">
      <v-card
        class="mx-auto"
        width="256"
        tile
      >

        <v-navigation-drawer permanent app expand-on-hover :mini-variant.sync="mini">
          <v-system-bar></v-system-bar>
          <v-toolbar-title class="text-center">周美丽专属教务系统{{ $store.state.recordId }}</v-toolbar-title>
          <v-list>
            <v-list-item>
              <v-list-item-avatar>
                <v-avatar color="red" size="36">
                  <span class="white--text headline">CJ</span>
                </v-avatar>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title class="title">
                  周美丽
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
              <v-list-item to="/SelectCourse">
                <v-list-item-content>
                  <v-list-item-title>作业管理</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/SetClass">
                <v-list-item-content>
                  <v-list-item-title>安排课程</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/Apply">
                <v-list-item-content>
                  <v-list-item-title>修改课程申请</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/ClassSheet">
                <v-list-item-content>
                  <v-list-item-title>查看课表</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/ChooseCourse">
                <v-list-item-content >
                  <v-list-item-title>选课</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/ChosenCourse">
                <v-list-item-content >
                  <v-list-item-title>查看已选课程/退课</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

              <v-list-item to="/QueryResults">
                <v-list-item-content >
                  <v-list-item-title>查询成绩</v-list-item-title>
                </v-list-item-content>
              </v-list-item>

            </v-list-item-group>
          </v-list>

        </v-navigation-drawer>
      </v-card>
      <v-main>
        <v-snackbar v-model="$store.state.successFlag" :timeout="3000" color="green" top>
          <v-icon>
            mdi-check-circle
          </v-icon>
          <span>
        {{ $store.state.successMsg }}
      </span>
        </v-snackbar>
        <v-snackbar v-model="$store.state.errorFlag" :timeout="3000" color="red" top>
          <v-icon>
            mdi-minus-circle
          </v-icon>
          <span>
        {{ $store.state.errorMsg }}
      </span>
        </v-snackbar>
        <router-view></router-view>
      </v-main>
    </v-app>
    <v-app>
      <v-dialog
        persistent
        overlay-opacity="0.92"
        v-model="dialog"
        width="500"
        transition="dialog-top-transition"
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
    </v-app>
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

export default {
  name: 'App',

  components: {
    HelloWorld,
    Canvas,
    ClassSheet,
    SelectCourse,
    SelectHomework,
    ChooseCourse,
    ChosenCourse
  },

  data: () => ({
    account:null,
    password:null,
    dialog:true,
    selectedItem: 0,
    mini: true,
    isActive: true,

  }),
  methods:{
    accept(){
      this.dialog = false
      document.getElementById("blur").style.filter = "blur(0px)"
    }
  }
};
</script>

<style scoped>
#blur {
  -webkit-filter: blur(10px);
}
</style>
