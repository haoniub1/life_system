import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue')
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth) {
    // Already logged in (store has user)
    if (userStore.isLoggedIn) {
      return true
    }

    // Not logged in but have token in localStorage — try to restore session
    const token = localStorage.getItem('token')
    if (token) {
      try {
        await userStore.fetchMe()
        return true
      } catch {
        // Token is invalid/expired, clean up
        localStorage.removeItem('token')
      }
    }

    // No valid session, redirect to login
    return '/login'
  }

  // If already logged in, don't show login/register pages
  if (userStore.isLoggedIn && (to.path === '/login' || to.path === '/register')) {
    return '/'
  }

  return true
})

export default router
