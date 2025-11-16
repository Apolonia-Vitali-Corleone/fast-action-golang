/**
 * Vue Router 配置
 * 定义应用的路由规则和导航守卫
 */
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/RegisterView.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局前置守卫 - 处理认证逻辑
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const token = localStorage.getItem('token')

  // 如果去往需要认证的页面
  if (to.meta.requiresAuth) {
    if (!token) {
      // 无 token，静默跳转到登录页
      next('/login')
      return
    }

    // 有 token，但还没有用户信息，尝试恢复会话
    if (!userStore.currentUser) {
      const restored = await userStore.restoreSession()
      if (!restored) {
        // 恢复失败（token 无效），静默跳转到登录页
        next('/login')
        return
      }
    }

    // 认证通过，允许访问
    next()
  } else {
    // 去往不需要认证的页面（登录/注册）
    // 如果已经登录了，直接跳转到首页
    if (token && userStore.currentUser) {
      next('/')
    } else {
      next()
    }
  }
})

export default router
