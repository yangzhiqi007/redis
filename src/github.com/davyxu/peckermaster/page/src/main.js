// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import App from './App'
import routes from './routes'
import VueRouter from 'vue-router'


Vue.config.productionTip = false

Vue.use(ElementUI)
Vue.use(VueRouter)

const router = new VueRouter({
  routes
})

router.beforeEach((to, from, next) => {

  //console.log(from.path + " -> " + to.path)

  if (to.path === '/login') {
    sessionStorage.removeItem('user');
  }
  let user = sessionStorage.getItem('user');
  if (!user && to.path !== '/login') {
    next({ path: '/login' })
  } else {
    next()
  }
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')


