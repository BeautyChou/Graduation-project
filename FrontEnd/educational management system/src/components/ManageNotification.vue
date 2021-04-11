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
      <editor
        api-key="2wktxhha9m0ig3njfyxvd50y8us108lenerm7fv0y6xgp0vl"
        :init="{
      height:'500',
      selector: 'textarea',
      language:'zh_CN',
      plugins: [
        'a11ychecker advcode casechange formatpainter linkchecker autolink lists checklist media mediaembed pageembed permanentpen powerpaste table advtable tinycomments tinymcespellchecker',
        'advlist autolink lists link image charmap print preview anchor',
        'searchreplace visualblocks code fullscreen',
        'insertdatetime media table paste code help wordcount'
        ],
      toolbar: 'undo redo | formatselect | bold italic backcolor | \
           alignleft aligncenter alignright alignjustify | \
           bullist numlist outdent indent | removeformat  | a11ycheck addcomment showcomments casechange checklist code formatpainter pageembed permanentpen table| help',
      toolbar_mode: 'floating',
      tinycomments_mode: 'embedded',
      tinycomments_author: 'Author name',
       }"
        v-model="Alert"
      />
    </v-card>
  </div>
</template>
<script>
import Editor from '@tinymce/tinymce-vue'

export default {
  name: "ManageAlert",
  created() {
    this.getNotification()
  },
  components: {
    'editor': Editor
  },
  data(){
    return{
      Alert:null
    }
  },
  methods:{
    getNotification(){
      this.$axios({
        url:"getNotification",
        method:"get",
        headers:{
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response)=>{
        console.log(response)
        this.Alert = response.data.notification.notification
      })
    },
    putNotification(){
      const formData = new FormData()
      formData.append("notification",this.Alert)
      this.$axios({
        url:"putNotification",
        method:"put",
        data:formData,
        headers: {
          "Content-Type": "multipart/form-data",
          'Token': "8a54sh " + this.$store.state.Jwt
        }
      }).then((response)=>{
        if (response.data.msg === "Token无效") {
          this.$emit('func')
          return
        }
        console.log(response)
        this.$store.commit(response.data.snackbar,response.data.msg)
        setTimeout(()=>{
          this.$store.commit(response.data.snackbar2)
        },3000)
      })
    }
  }
}
</script>

<style scoped>

</style>
