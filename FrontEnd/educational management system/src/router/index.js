import Vue from 'vue'
import Router from 'vue-router'
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
import Manage from "../components/Manage";
import store from "../store/index"

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/canvas',
      name: 'canvas',
      component: canvas,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 2){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path: '/SetClass',
      name:'SetClass',
      component:SetClass,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 3){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/SelectCourse',
      name:'SelectCourse',
      component:SelectCourse,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1||store.state.level === 2){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/SelectHomework',
      name:'SelectHomework',
      component:SelectHomework,
      meta:{
        keepAlive: true
      },
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1||store.state.level === 2){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/AddHomework',
      name:'AddHomework',
      component:AddHomework,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 2){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/SubmitHomework',
      name:'SubmitHomework',
      component:SubmitHomework,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/CheckReview',
      name:'CheckReview',
      component:CheckReview,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/Apply',
      name:'Apply',
      component:Apply,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 3){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/ClassSheet',
      name:'ClassSheet',
      component:ClassSheet,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1||store.state.level === 2){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/ChooseCourse',
      name:'ChooseCourse',
      component:ChooseCourse,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/ChosenCourse',
      name:'ChosenCourse',
      component:ChosenCourse,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/RecordGrades',
      name:'RecordGrades',
      component:RecordGrades,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 2){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/QueryResults',
      name:'QueryResults',
      component:QueryResults,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 1){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'/UserInfo',
      name:'UserInfo',
      component:UserInfo,
      alias:'/',
    },
    {
      path:'/Manage',
      name:'Manage',
      component:Manage,
      beforeEnter:(to, from, next)=>{
        console.log(store)
        if (store.state.level === 3){
          next()
        }else {
          store.commit('setError',"无权进入该页面！")
          next("/UserInfo")
        }
      },
    },
    {
      path:'*',
      redirect:'/UserInfo',
    }
  ]
})
