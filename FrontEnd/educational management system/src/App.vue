<template>
  <v-app>
    <v-card
      class="mx-auto"
      width="256"
      tile
    >

      <v-navigation-drawer permanent app expand-on-hover :mini-variant.sync="mini">
        <v-system-bar></v-system-bar>
        <v-toolbar-title class="text-center">周美丽专属教务系统{{ $store.state.homeworkId }}</v-toolbar-title>
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
          dense
        >
          <v-list-group>
            <template v-slot:activator>
              <v-list-item-content >
                <v-list-item-title>报名申请</v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>学籍异动申请</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>

          <v-list-item>
            <v-list-item-content>
              <router-link to="/SelectCourse">
                <v-list-item-title>作业管理</v-list-item-title>
              </router-link>
            </v-list-item-content>
          </v-list-item>

          <v-list-item>
            <v-list-item-content>
              <router-link to="/SetClass">
                <v-list-item-title>安排课程</v-list-item-title>
              </router-link>
            </v-list-item-content>
          </v-list-item>

          <v-list-item>
            <v-list-item-content>
              <router-link to="/Apply">
                <v-list-item-title>修改课程申请</v-list-item-title>
              </router-link>
            </v-list-item-content>
          </v-list-item>

          <v-list-item>
            <v-list-item-content>
              <router-link to="/ClassSheet">
                <v-list-item-title>查看课表</v-list-item-title>
              </router-link>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <router-link to="/ChooseCourse">
                <v-list-item-title>选课</v-list-item-title>
              </router-link>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-content>
              <router-link to="/ChosenCourse">
                <v-list-item-title>查看已选课程/退课</v-list-item-title>
              </router-link>
            </v-list-item-content>
          </v-list-item>

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
    selectedItem: 0,
    mini:true,
    isActive:true
  }),
};
</script>
