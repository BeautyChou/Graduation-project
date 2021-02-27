import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '../components/HelloWorld'
import canvas from '../components/Canvas'
import ClassSheet from "../components/ClassSheet";
import canvastest from "../components/canvastest";
import SelectCourse from "../components/SelectCourse";
import SelectHomework from "../components/SelectHomework";

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'HelloWorld',
      component: HelloWorld
    },
    {
      path: '/canvas',
      name: 'canvas',
      component: canvas,
    },
    {
      path: '/Canvastest',
      name: 'Canvastest',
      component: canvastest,
    },
    {
      path: '/ClassSheet',
      name:'ClassSheet',
      component:ClassSheet,
    },
    {
      path:'/SelectCourse',
      name:'SelectCourse',
      component:SelectCourse,
    },
    {
      path:'/SelectHomework',
      name:'SelectHomework',
      component:SelectHomework
    }
  ]
})
