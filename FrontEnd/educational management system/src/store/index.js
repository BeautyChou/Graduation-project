import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from "vuex-persistedstate"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    homeworkId: null,
    courseId: null,
    JwtFlag:true,
    Jwt: "",
    teacherId: 1,
    studentId: 1,
    homeworkFlag: false,
    courseFlag: false,
    level: 3,
    successFlag: false,
    successMsg: null,
    errorFlag: false,
    errorMsg: null,
    modifyHomeworkFlag: false,
    refreshFlag: 0,
    recordId: null,
    facultyId: 1,
    specialtyId: 1,
    directionId: 0,
    username: "123",
  },
  mutations: {
    jwtInvalid(state){
      state.JwtFlag = true
    },
    jwtValid(state){
      state.JwtFlag = false
    },
    setStudentID(state, studentID) {
      state.studentId = studentID
    },
    setTeacherID(state, teacherID) {
      state.teacherId = teacherID
    },
    setUsername(state, username) {
      state.username = username
    },
    setFacultyId(state, facultyID) {
      state.facultyId = facultyID
    },
    setSpecialtyId(state, specialtyID) {
      state.specialtyId = specialtyID
    },
    setDirectionId(state, directionID) {
      state.directionId = directionID
    },
    setRecordId(state, recordId) {
      state.recordId = recordId
    },
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
    closeSuccess(state) {
      state.successFlag = false
      state.successMsg = null
    },
    closeError(state) {
      state.errorFlag = false
      state.errorMsg = null
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
      state.refreshFlag = 1
    },
    setLevel(state, level) {
      if (state.level === level) ++state.refreshFlag
      state.level = level
    },
    clear(state) {
      state.homeworkId = null
      state.courseId = null
      state.Jwt = null
      state.teacherId = null
      state.studentId = null
      state.homeworkFlag = false
      state.courseFlag = false
      state.level = null
      state.successFlag = false
      state.successMsg = null
      state.errorFlag = false
      state.errorMsg = null
      state.modifyHomeworkFlag = false
      state.refreshFlag = 0
      state.recordId = null
      state.facultyId = null
      state.specialtyId = null
      state.directionId = null
      state.username = " "
    }
  },
  actions: {},
  modules: {},
  plugins: [createPersistedState()]
})
