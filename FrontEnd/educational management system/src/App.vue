<template>
  <v-app>
    <v-card
      class="mx-auto"
      width="256"
      tile
    >

      <v-navigation-drawer permanent app expand-on-hover :mini-variant.sync="mini">
        <v-system-bar></v-system-bar>
        <v-toolbar-title class="text-center">周美丽专属教务系统</v-toolbar-title>
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
                <v-list-item-title>重补修报名</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>考级项目报名</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>教学项目报名</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>缓考|补考确认</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>创新创业报名</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>学生资格证书申请</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title>学籍异动申请</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-group>
        </v-list>
      </v-navigation-drawer>
    </v-card>
    <v-main>
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script>
import HelloWorld from './components/HelloWorld';
import Canvas from "./components/Canvas";
import ClassSheet from "./components/ClassSheet";

export default {
  name: 'App',

  components: {
    HelloWorld,
    Canvas,
    ClassSheet

  },

  data: () => ({
    selectedItem: 0,
    mini:true,
    isActive:true
  }),
};
</script>
