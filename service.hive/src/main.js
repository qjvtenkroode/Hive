import Vue from 'vue'
import VueRouter from 'vue-router'
import App from './App.vue'
import './main.css'
import HiveRooms from './components/HiveRooms'
import HiveSwitches from './components/HiveSwitches'
import HiveRegistry from './components/HiveRegistry'

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

Vue.use(VueRouter);

Vue.config.productionTip = false

window.axios = require('axios');
window.host = location.hostname;

const router = new VueRouter({
    routes: [
        { path: '/', component: HiveSwitches },
        { path: '/rooms', component: HiveRooms },
        { path: '/registry', component: HiveRegistry }
    ],
    mode: 'history'
});

new Vue({
    router,
    render: h => h(App),
}).$mount('#app')
