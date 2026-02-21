import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'

import Login from './views/Login.vue'
import Dashboard from './views/Dashboard.vue'
import NewOrder from './views/NewOrder.vue'
import ActiveOrders from './views/ActiveOrders.vue'
import History from './views/History.vue'
import ManagePackages from './views/ManagePackages.vue'
import ManageMenu from './views/ManageMenu.vue'
import ManageDiscounts from './views/ManageDiscounts.vue'
import ManageUsers from './views/ManageUsers.vue'
import ManageShifts from './views/ManageShifts.vue'
import './assets/style.css'
const routes = [
  { path: '/login', name: 'Login', component: Login, meta: { noAuth: true } },
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/new-order', name: 'NewOrder', component: NewOrder },
  { path: '/active', name: 'ActiveOrders', component: ActiveOrders },
  { path: '/history', name: 'History', component: History },
  { path: '/manage/packages', name: 'ManagePackages', component: ManagePackages, meta: { admin: true } },
  { path: '/manage/menu', name: 'ManageMenu', component: ManageMenu, meta: { admin: true } },
  { path: '/manage/discounts', name: 'ManageDiscounts', component: ManageDiscounts, meta: { admin: true } },
  { path: '/manage/users', name: 'ManageUsers', component: ManageUsers, meta: { admin: true } },
  { path: '/manage/shifts', name: 'ManageShifts', component: ManageShifts, meta: { admin: true } },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

// Auth guard
router.beforeEach((to, from, next) => {
  const user = JSON.parse(localStorage.getItem('currentUser') || 'null')
  if (!to.meta.noAuth && !user) {
    next('/login')
  } else if (to.meta.admin && user?.role !== 'superadmin') {
    next('/')
  } else {
    next()
  }
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
