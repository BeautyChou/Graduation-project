<template>
  <v-tabs grow>
    <v-tab @click="$refs.manageStudent.getStudentList()">学生管理</v-tab>
    <v-tab @click="$refs.manageTeacher.getTeacherList()">教师管理</v-tab>
    <v-tab @click="$refs.manageAdmin.getAdminList()">管理员管理</v-tab>
    <v-tab @click="$refs.manageTitle.getTitleList()">职称管理</v-tab>
    <v-tab @click="$refs.managePunishLevel.getPunishLevelList()">处分等级管理</v-tab>
    <v-tab @click="$refs.manageFaculty.getDirectionSpecialtyFacultyList()">学院管理</v-tab>
    <v-tab @click="$refs.manageClassroom.getClassroomList()">教室管理</v-tab>

    <v-tab-item>
      <v-manage-student ref="manageStudent"></v-manage-student>
    </v-tab-item>
    <v-tab-item>
      <v-manage-teacher ref="manageTeacher"></v-manage-teacher>
    </v-tab-item>
    <v-tab-item>
      <v-manage-admin ref="manageAdmin"></v-manage-admin>
    </v-tab-item>
    <v-tab-item>
      <v-manage-title ref="manageTitle"></v-manage-title>
    </v-tab-item>
    <v-tab-item>
      <v-manage-punish-level ref="managePunishLevel"></v-manage-punish-level>
    </v-tab-item>
    <v-tab-item>
      <v-manage-faculty ref="manageFaculty"></v-manage-faculty>
    </v-tab-item>
    <v-tab-item>
      <v-manage-classroom ref="manageClassroom"></v-manage-classroom>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import ManageStudent from "./ManageStudent";
import ManageTeacher from "./ManageTeacher";
import ManageAdmin from "./ManageAdmin";
import ManageTitle from "./ManageTitle";
import ManagePunishLevel from "./ManagePunishLevel";
import ManageClassroom from "./ManageClassroom";
import ManageFaculty from "./ManageFaculty";
export default {
  name: "Manage",
  components: {
    "v-manage-student":ManageStudent,
    "v-manage-teacher":ManageTeacher,
    "v-manage-admin":ManageAdmin,
    "v-manage-title":ManageTitle,
    "v-manage-punish-level":ManagePunishLevel,
    "v-manage-classroom":ManageClassroom,
    "v-manage-faculty":ManageFaculty,
  },
  data() {
    return {
      admin_header: [
        {text: '姓名', align: 'start', sortable: false, value: 'student_name'},
        {text: '操作', sortable: false, value: 'operation'}
      ]
    }
  },
  methods: {
  },
  watch: {
    "$route.path": {
      handler(newVal, oldVal) {
        this.$refs.manageStudent.getStudentList()
      }
    }
  }
}
</script>

<style scoped>

</style>
