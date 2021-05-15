<template>
  <div>
    <v-card>
      <v-card-title>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="putNotification">
          <v-icon>
            mdi-save
          </v-icon>
          保存通知
        </v-btn>
      </v-card-title>
      <textarea id="tinyMCE" v-model="Alert"/>
    </v-card>
  </div>
</template>
<script>

export default {
  name: "ManageAlert",
  created() {
    this.getNotification()
  },
  mounted() {
    this.initTinymce()
  },
  data() {
    return {
      Alert: null,
    }
  },
  methods: {
    getNotification() {
      this.$axios({
        url: "Notification",
        method: "get",
        headers: {
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response) => {
        this.Alert = response.data.notification.notification
      })
    },
    putNotification() {
      const formData = new FormData()
      formData.append("notification", this.Alert)
      this.$axios({
        url: "Notification",
        method: "put",
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
        this.$store.commit(response.data.snackbar, response.data.msg)
        setTimeout(() => {
          this.$store.commit(response.data.snackbar2)
        }, 3000)
      })
    },
    initTinymce() {
      window.tinymce.init({
        branding: false,
        height: '500',
        selector: 'textarea',
        language_url: '../../static/zh_CN.js',
        language: 'zh_CN',
        plugins: [
          'checklist pageembed permanentpen advtable tinymcespellchecker',
          'link preview anchor',
          'table code help wordcount'
        ],
        toolbar: 'undo redo | formatselect | bold italic backcolor | \
           alignleft aligncenter alignright alignjustify | \
           bullist numlist outdent indent | removeformat  | code table| help',
        toolbar_mode: 'floating',
        init_instance_callback: editor => {
          this.hasInit = true
          editor.on('KeyUp', () => {
            this.hasChange = true
            this.Alert = editor.getContent()
          })
        },
      })
    },
  },
  inject: ['expire']
}
</script>

<style scoped>

</style>
