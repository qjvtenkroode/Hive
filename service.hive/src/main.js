import Vue from 'vue'
import App from './App.vue'
import './main.css'
// fontawesome
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
// fontawesome icons
import { faRandom } from '@fortawesome/free-solid-svg-icons'
import { faToggleOn } from '@fortawesome/free-solid-svg-icons'
import { faToggleOff } from '@fortawesome/free-solid-svg-icons'
import { faExclamationTriangle } from '@fortawesome/free-solid-svg-icons'
import { faCheckCircle } from '@fortawesome/free-regular-svg-icons'

library.add(faRandom)
library.add(faToggleOn)
library.add(faToggleOff)
library.add(faExclamationTriangle)
library.add(faCheckCircle)

Vue.component('font-awesome-icon', FontAwesomeIcon)

Vue.config.productionTip = false

window.axios = require('axios');
window.host = location.hostname;

new Vue({
  render: h => h(App),
}).$mount('#app')
