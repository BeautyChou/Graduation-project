import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '../components/HelloWorld'
import canvas from '../components/Canvas'
import SetClass from "../components/SetClass";
import SelectCourse from "../components/SelectCourse";
import SelectHomework from "../components/SelectHomework";
import AddHomework from "../components/AddHomework";
import SubmitHomework from "../components/SubmitHomework";
import CheckReview from "../components/CheckReview";
import Apply from "../components/Apply";
import ClassSheet from "../components/ClassSheet";
import ChooseCourse from "../components/ChooseCourse";
import ChosenCourse from "../components/ChosenCourse";
import RecordGrades from "../components/RecordGrades";
import QueryResults from "../components/QueryResults";
import UserInfo from "../components/UserInfo";

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
      path: '/SetClass',
      name:'SetClass',
      component:SetClass,

    },
    {
      path:'/SelectCourse',
      name:'SelectCourse',
      component:SelectCourse,

    },
    {
      path:'/SelectHomework',
      name:'SelectHomework',
      component:SelectHomework,
      meta:{
        keepAlive: true
      }

    },
    {
      path:'/AddHomework',
      name:'AddHomework',
      component:AddHomework,

    },
    {
      path:'/SubmitHomework',
      name:'SubmitHomework',
      component:SubmitHomework,

    },
    {
      path:'/CheckReview',
      name:'CheckReview',
      component:CheckReview,
    },
    {
      path:'/Apply',
      name:'Apply',
      component:Apply,
    },
    {
      path:'/ClassSheet',
      name:'ClassSheet',
      component:ClassSheet,
    },
    {
      path:'/ChooseCourse',
      name:'ChooseCourse',
      component:ChooseCourse,
    },
    {
      path:'/ChosenCourse',
      name:'ChosenCourse',
      component:ChosenCourse,
    },
    {
      path:'/RecordGrades',
      name:'RecordGrades',
      component:RecordGrades
    },
    {
      path:'/QueryResults',
      name:'QueryResults',
      component:QueryResults
    },
    {
      path:'/UserInfo',
      name:'UserInfo',
      component:UserInfo
    },
  ]
})
