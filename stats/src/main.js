import Vue from 'vue'
import BootstrapVue from "bootstrap-vue"
import App from './App.vue'
import './assets/bootstrap.min.css'

Vue.use(BootstrapVue)

new Vue({
  el: '#app',
  render: h => h(App)
})
