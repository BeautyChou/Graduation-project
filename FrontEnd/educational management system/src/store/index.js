import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    homeworkId: null,
    courseId: null,
    Jwt: null,
    teacherId: 1,
    studentId: 171020101,
    homeworkFlag: false,
    courseFlag: false,
    level: 1,
    successFlag: false,
    successMsg: null,
    errorFlag: false,
    errorMsg: null,
    modifyHomeworkFlag: false,
    refreshFlag: 0,
  },
  mutations: {
    setJwt(state, str) {
      state.Jwt = str
    },
    setHomeworkId(state, id) {
      if (state.homeworkId === id) ++state.refreshFlag;
      state.homeworkId = id
      state.homeworkFlag = true
    },
    setCourseId(state, id) {
      if (state.courseId === id) ++state.refreshFlag;
      state.courseId = id
      state.courseFlag = true
      return 0
    },
    setTeacherId(state, id) {
      state.teacherId = id
    },
    jumpToHomeworkSelect(state) {
      state.courseFlag = false
    },
    jumpToCanvas(state) {
      state.homeworkFlag = false
    },
    setSuccess(state, str) {
      state.successFlag = true
      state.successMsg = str
    },
    setError(state, str) {
      state.errorFlag = true
      state.errorMsg = str
    },
    modifyHomework(state) {
      state.modifyHomeworkFlag = true
    },
    addHomework(state) {
      state.modifyHomeworkFlag = false
    },
    previousPage(state) {
      state.refreshFlag++
    },
    nextPage(state) {
      state.refreshFlag = 0
    }
  }
  ,
  actions: {},
  modules: {}
})
