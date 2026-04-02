import { createRouter, createWebHistory } from 'vue-router'
import LogIn from '../components/LogIn.vue'
import SignIn from '../components/SignIn.vue' 

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: LogIn
  },
  {
    path: '/signup',
    name: 'Signup',
    component: SignIn
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router