// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import vuetify from './plugins/vuetify';
import '@babel/polyfill'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import '@mdi/font/css/materialdesignicons.css'
import axios from 'axios'
import './assets/index.css'
import store from './store'
import config from 'raw-loader!../../../BackEnd/config.yaml'


Vue.prototype.$axios = axios
Vue.config.productionTip = false
/* eslint-disable no-new */
//读取配置文件中的baseURL
const configArray = config.split("\n")
let baseUrl
configArray.some((val, index) => {
  if (val.indexOf("baseurl") !== -1) {
    baseUrl = val
  }
  return val.indexOf("baseurl") !== -1
})
console.log(baseUrl.split(": ")[1])
axios.defaults.baseURL = baseUrl.split(": ")[1]

new Vue({
  el: '#app',
  router,
  components: {App},
  vuetify,
  store,
  template: '<App/>'
})
