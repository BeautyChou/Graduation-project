<template>
  <v-card>
    <v-card-title>
      <v-row>
        <v-select
          class="col-3"
          outlined
          dense
          v-model="faculty_id"
          :items="faculties"
          item-text="name"
          item-value="id"
          label="学院"
          @change="getDirectionSpecialtyFacultyList"
        ></v-select>
        <v-select
          class="col-3"
          outlined
          dense
          v-model="specialty_id"
          :items="specialties[faculty_id]"
          item-text="specialty_name"
          item-value="specialty_id"
          label="专业"
          :disabled="faculty_id ==  null"
          @change="getDirectionSpecialtyFacultyList"></v-select>
        <v-btn
          icon
          @click.native="faculty_id = null;specialty_id = null;getDirectionSpecialtyFacultyList()">
          <v-icon>
            mdi-trash-can-outline
          </v-icon>
        </v-btn>
        <v-spacer></v-spacer>

        <v-menu offset-y>
          <template v-slot:activator="{ on,attrs }">
            <v-btn color="primary" v-bind="attrs" v-on="on">
              <v-icon>mdi-plus-thick</v-icon>
              <span>
          添加...
        </span>
            </v-btn>
          </template>
          <v-list>
            <v-list-item @click="add_faculty = true">
              <v-list-item-title>学院</v-list-item-title>
            </v-list-item>
            <v-list-item @click="add_specialty = true">
              <v-list-item-title>专业</v-list-item-title>
            </v-list-item>
            <v-list-item @click="add_direction = true">
              <v-list-item-title>方向</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>

        <v-menu offset-y>
          <template v-slot:activator="{ on,attrs }">
            <v-btn color="red" v-bind="attrs" v-on="on">
              <v-icon>mdi-delete</v-icon>
              <span>
          删除...
        </span>
            </v-btn>
          </template>
          <v-list>
            <v-list-item @click="delete_faculty = true">
              <v-list-item-title>学院</v-list-item-title>
            </v-list-item>
            <v-list-item @click="delete_specialty = true">
              <v-list-item-title>专业</v-list-item-title>
            </v-list-item>
            <v-list-item @click="delete_direction = true">
              <v-list-item-title>方向</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>

      </v-row>
    </v-card-title>
    <v-data-table
      :headers="faculty_headers"
      :items="faculty_specialty_direction"
      :options.sync="options"
      :server-items-length="total"
    >
      <template v-slot:item.direction_name="{ item }">
        {{ item.direction_id === 0 ? "无" : item.direction_name }}
      </template>
      <template v-slot:item.specialty_name="{ item }">
        {{ item.specialty_id === 0 ? "无" : item.specialty_name }}
      </template>
    </v-data-table>
    <v-dialog v-model="add_faculty" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加学院
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-text-field
              dense
              outlined
              v-model="new_faculty"
              label="学院名"
            >
            </v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :disabled="new_faculty == null"
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addFacultySpecialtyDirection(1)"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_faculty = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_specialty" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加专业
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              label="学院"
            ></v-select>
            <v-text-field
              dense
              outlined
              v-model="new_specialty"
              label="专业名"
            >
            </v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :disabled="new_specialty == null"
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addFacultySpecialtyDirection(2)"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_specialty = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="add_direction" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          添加方向
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              label="学院"
            ></v-select>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_specialty"
              :items="specialties[new_faculty]"
              item-text="specialty_name"
              item-value="specialty_id"
              label="专业"
              :disabled="new_faculty ==  null">
            </v-select>
            <v-text-field
              outlined
              v-model="new_direction"
              :disabled="new_specialty == null"
              label="方向名"
            >

            </v-text-field>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :disabled="new_direction == null"
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="addFacultySpecialtyDirection(3)"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="add_direction = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="delete_faculty" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          删除学院
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              label="学院"
            ></v-select>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :disabled="new_faculty == null"
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteFacultySpecialtyDirection(1)"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_faculty = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="delete_specialty" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          删除专业
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              label="学院"
            ></v-select>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_specialty"
              :items="specialties[new_faculty]"
              item-text="specialty_name"
              item-value="specialty_id"
              label="专业"
              :disabled="new_faculty ==  null">
            </v-select>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :disabled="new_specialty == null"
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteFacultySpecialtyDirection(2)"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_specialty = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="delete_direction" width="500" persistent>
      <v-card>
        <v-card-title class="headline font-weight-bold">
          删除方向
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_faculty"
              :items="faculties"
              item-text="name"
              item-value="id"
              label="学院"
            ></v-select>
            <v-select
              class="col-6"
              outlined
              dense
              v-model="new_specialty"
              :items="specialties[new_faculty]"
              item-text="specialty_name"
              item-value="specialty_id"
              label="专业"
              :disabled="new_faculty ==  null">
            </v-select>
            <v-select
              class="col-12"
              outlined
              dense
              v-model="new_direction"
              :items="directions[new_specialty]"
              item-text="direction_name"
              item-value="direction_id"
              label="专业"
              :disabled="new_specialty ==  null">
            </v-select>
          </v-row>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            :disabled="new_direction == null"
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="deleteFacultySpecialtyDirection(3)"
          >
            是
          </v-btn>
          <v-btn
            class="font-weight-bold"
            color="green darken-1"
            text
            @click="delete_direction = false"
          >
            否
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>

</template>

<script>
export default {
  name: "ManageFaculty",
  data() {
    return {
      options: {},
      total: null,
      new_faculty: null,
      new_specialty: null,
      new_direction: null,
      delete_direction: false,
      delete_faculty: false,
      delete_specialty: false,
      add_direction: false,
      add_faculty: false,
      add_specialty: false,
      faculty_id: null,
      specialty_id: null,
      faculty_specialty_direction: [],
      directions: [],
      specialties: [],
      faculties: [],
      faculty_headers: [
        {text: '学院名', align: 'start', sortable: false, value: 'faculty_name'},
        {text: '专业名', align: 'start', sortable: false, value: 'specialty_name'},
        {text: '方向名', align: 'start', sortable: false, value: 'direction_name'},
      ]
    }
  },
  created() {
    this.getDirectionSpecialtyFacultyList()
  },
  methods: {
    getDirectionSpecialtyFacultyList() {
      this.$axios({
        url: "DirectionSpecialtyFaculty",
        params: {
          faculty_id: this.faculty_id,
          specialty_id: this.specialty_id,
          "page": this.options.page,
          "items": this.options.itemsPerPage,
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        },
        method: "get"
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.specialties = response.data.specialties
        this.faculties = response.data.faculties
        this.directions = response.data.directions
        this.faculty_specialty_direction = response.data.faculty_specialty_direction
        this.total = response.data.total
      })
    },
    addFacultySpecialtyDirection(flag) {
      const formData = new FormData()
      if (flag === 1) {
        formData.append("faculty_name", this.new_faculty)
      } else if (flag === 2) {
        formData.append("faculty_id", this.new_faculty)
        formData.append("specialty_name", this.new_specialty)
      } else if (flag === 3) {
        this.specialties[this.new_faculty].some((v, i) => {
          if (v.specialty_id === this.new_specialty) {
            formData.append("specialty_name", v.specialty_name)
          }
        })
        formData.append("faculty_id", this.new_faculty)
        formData.append("specialty_id", this.new_specialty)
        formData.append("direction_name", this.new_direction)
      }
      this.$axios({
        url: "DirectionSpecialtyFaculty",
        method: "post",
        data: formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.add_faculty = this.add_specialty = this.add_direction = false
        this.new_faculty = this.new_specialty = this.new_direction = null
        this.getDirectionSpecialtyFacultyList()
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    deleteFacultySpecialtyDirection() {
      this.$axios({
        url: "DirectionSpecialtyFaculty",
        method: "delete",
        params: {
          faculty_id: this.new_faculty,
          specialty_id: this.new_specialty,
          direction_id: this.new_direction,
        },
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        if (response.data.msg === "Token无效") {
          this.expire()
          return
        }
        this.delete_faculty = this.delete_specialty = this.delete_direction = false
        this.new_faculty = this.new_specialty = this.new_direction = null
        this.getDirectionSpecialtyFacultyList()
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    }
  },
  watch: {
    options: {
      handler() {
        this.getDirectionSpecialtyFacultyList()
      },
      deep: true,
    }
  },
  inject: ['expire']
}
</script>

<style scoped>

</style>
