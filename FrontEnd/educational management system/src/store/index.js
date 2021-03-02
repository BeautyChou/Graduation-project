import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    homeworkId:null,
    courseId:null,
    Jwt:null,
    teacherId:1,
    studentId:null,
    homeworkFlag:false,
    courseFlag:false,
    level:1,
  },
  mutations: {
    setJwt(state,str){
      state.Jwt = str
    },
    setHomeworkId(state,id){
      state.homeworkId = id
      state.homeworkFlag = true
    },
    setCourseId(state,id){
      state.courseId = id
      state.courseFlag = true
      return 0
    },
    setTeacherId(state,id){
      state.teacherId = id
    },
    jumpToHomeworkSelect(state){
      state.courseFlag = false
    },
    jumpToCanvas(state){
      state.homeworkFlag = false
    }
  },
  actions: {
    async setCourseId({ commit },id){
      commit('setCourseId',id)
    }
  },
  modules: {
  }
})
