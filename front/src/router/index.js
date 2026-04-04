import { createRouter, createWebHistory } from 'vue-router'
import LogIn from '../components/LogIn.vue'
import SignIn from '../components/SignIn.vue'
import Dashboard from '../components/TaskList.vue'

// DEVELOPMENT ONLY - set to true to bypass authentication
const DEV_BYPASS_AUTH = true  // 👈 Change to false for production!

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/login', name: 'Login', component: LogIn },
  { path: '/signup', name: 'Signup', component: SignIn },
  { path: '/dashboard', name: 'Dashboard', component: Dashboard, meta: { requiresAuth: true } },
  { path: '/process/:id', name: 'Process', component: () => import('../components/ProcessView.vue'), meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  // DEVELOPMENT BYPASS
  if (DEV_BYPASS_AUTH) {
    // Set mock user if not exists
    if (!localStorage.getItem('auth_token')) {
      localStorage.setItem('auth_token', 'dev-bypass-token')
      localStorage.setItem('user', JSON.stringify({ 
        id: 1, 
        username: 'DevUser',
        email: 'dev@example.com'
      }))
    }
    
    // Allow access to all routes
    if (to.path === '/login' || to.path === '/signup') {
      next('/dashboard')
    } else {
      next()
    }
    return
  }

  // Normal authentication for production
  const isAuthenticated = localStorage.getItem('auth_token')
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if ((to.path === '/login' || to.path === '/signup') && isAuthenticated) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router