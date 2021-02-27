<template>
  <div>
    <v-data-table
      :headers="headers"
      :items="classes"
      class="elevation-1"
      >
      <template v-slot:item.time="{ item }">
        <v-select
          v-model="item.time"
          :items="classTime"
          item-text="name"
          item-value="value"
        ></v-select>
      </template>
      <template v-slot:item.teacher="{ item }">
        <v-select
          v-model="item.teacher"
          :items="teachers"
          item-text="name"
          item-value="value"
        ></v-select>
      </template>
      <template v-slot:item.classes="{ item }">
        <v-select
          v-model="item.classes"
          :items="classTable"
          item-text="name"
          item-value="value"
        ></v-select>
      </template>
      <template v-slot:item.classroom="{ item }">
        <v-select
          v-model="item.classroom"
          :items="classrooms"
          item-text="name"
          item-value="value"
        ></v-select>
      </template>
      <template v-slot:item.classroom="{ item }">
        <v-select
          v-model="item.classroom"
          :items="classrooms"
          item-text="name"
          item-value="value"
        ></v-select>
      </template>
      <template v-slot:item.startWeek="{ item }">
        <v-select
          v-model="item.startWeek"
          :items="weeks"
          item-text="name"
          item-value="value"
        ></v-select>
      </template>
      <template v-slot:item.endWeek="{ item }">
        <v-select
          v-model="item.endWeek"
          :items="weeks"
          item-text="name"
          item-value="value"
        ></v-select>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import draggable from 'vuedraggable'

export default {
  name: "ClassSheet",
  data() {
    return {
      headers:[
        {text:'课程',align:'start',sortable:false,value:'name'},
        {text:'节次',sortable:false,value:'time'},
        {text:'教师',sortable:false,value:'teacher'},
        {text:'教室',sortable:false,value:'classroom'},
        //学年由上传年份决定
        {text:'开始周',sortable:false,value:'startWeek'},
        {text:'结束周',sortable:false,value:'endWeek'},
      ],
      classes: [
        {name:'计算机导论',time:6,teacher:1,classroom:1,years:2,startWeek:3,endWeek:5}
      ],
      drag: false,
      classTime:[
        {name:"第一节",value:1},
        {name:"第二节",value:2},
        {name:"第三节",value:3},
        {name:"第四节",value:4},
        {name:"第五节",value:5},
        {name:"第六节",value:6},
        {name:"第七节",value:7},
        {name:"第八节",value:8},
        {name:"第九节",value:9},
        {name:"第十节",value:10},
        {name:"第十一节",value:11},
      ],
      teachers:[
        {name:"张三",value:1},
        {name:"李四",value:2},
        {name:"王五",value:3},
        {name:"赵六",value:4},
        {name:"周美丽",value:5},
        {name:"王二狗",value:6},
        {name:"李二狗",value:7},
        {name:"赵铁蛋",value:8},
        {name:"王八蛋",value:9},
        {name:"臭狗",value:10},
        {name:"哈哈哈哈哈",value:11},
      ],
      classrooms:[
        {name:"101",value:1},
        {name:"102",value:2},
        {name:"103",value:3},
        {name:"104",value:4},
        {name:"105",value:5},
        {name:"106",value:6},
        {name:"107",value:7},
        {name:"108",value:8},
        {name:"109",value:9},
        {name:"110",value:10},
        {name:"111",value:11},
      ],
      weeks:[
        {name:"第一周",value:1},
        {name:"第二周",value:2},
        {name:"第三周",value:3},
        {name:"第四周",value:4},
        {name:"第五周",value:5},
        {name:"第六周",value:6},
        {name:"第七周",value:7},
        {name:"第八周",value:8},
        {name:"第九周",value:9},
        {name:"第十周",value:10},
        {name:"第十一周",value:11},
        {name:"第十二周",value:12},
        {name:"第十三周",value:13},
        {name:"第十四周",value:14},
        {name:"第十五周",value:15},
        {name:"第十六周",value:16},
        {name:"第十七周",value:17},
        {name:"第十八周",value:18},
        {name:"第十九周",value:19},
        {name:"第二十周",value:20},
      ],
    }
  },
  components: {
  },
  methods: {
  },
}
</script>

<style scoped>

</style>
