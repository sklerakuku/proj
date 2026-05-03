import { createRouter, createWebHistory } from 'vue-router'
import LogIn from '../components/LogIn.vue'
import SignIn from '../components/SignIn.vue'
import Dashboard from '../components/TaskList.vue'
import ProcessView from '../components/ProcessView.vue'
import AdminPanel from '../components/AdminPanel.vue'  // Добавить импорт

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', name: 'Login', component: LogIn },
  { path: '/signup', name: 'Signup', component: SignIn },
  { 
    path: '/dashboard', 
    name: 'Dashboard', 
    component: Dashboard, 
    meta: { requiresAuth: true } 
  },
  { 
    path: '/process/:id', 
    name: 'Process', 
    component: ProcessView, 
    meta: { requiresAuth: true } 
  },
  { 
    path: '/admin', 
    name: 'Admin', 
    component: AdminPanel, 
    meta: { requiresAuth: true, requiresAdmin: true } 
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('auth_token')
  const user = JSON.parse(localStorage.getItem('user') || '{}')
  
  // Проверка авторизации
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
    return
  }
  
  // Проверка прав админа
  if (to.meta.requiresAdmin && user.role !== 'admin') {
    console.warn('Access denied: admin role required')
    next('/dashboard')
    return
  }
  
  // Редирект с логина если уже авторизован
  if ((to.path === '/login' || to.path === '/signup') && isAuthenticated) {
    if (user.role === 'admin') {
      next('/admin')
    } else {
      next('/dashboard')
    }
    return
  }
  
  next()
})

export default router