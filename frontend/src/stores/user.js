/**
 * 用户状态管理 - 简化版
 * 处理用户会话恢复和退出登录
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const API_BASE = 'http://localhost:8000/api'

// 配置Axios
axios.defaults.withCredentials = true

// 配置Axios请求拦截器 - 自动添加JWT token
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 配置Axios响应拦截器 - 处理token刷新和401错误
axios.interceptors.response.use(
  (response) => {
    const newToken = response.headers['x-new-token']
    if (newToken) {
      localStorage.setItem('token', newToken)
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // 静默清除 token 和用户状态，不显示提示
      localStorage.removeItem('token')
      // 路由守卫会自动处理跳转到登录页
    }
    return Promise.reject(error)
  }
)

export const useUserStore = defineStore('user', () => {
  // 状态
  const currentUser = ref(null)

  /**
   * 退出登录
   */
  const logout = async () => {
    try {
      await axios.post(`${API_BASE}/logout/`)
      localStorage.removeItem('token')
      currentUser.value = null
      ElMessage.success('已退出')
    } catch (error) {
      // 即使请求失败也要清除本地token
      localStorage.removeItem('token')
      currentUser.value = null
      console.error(error)
    }
  }

  /**
   * 恢复登录状态（从localStorage）
   */
  const restoreSession = async () => {
    const savedToken = localStorage.getItem('token')
    if (!savedToken) {
      return false
    }

    try {
      const res = await axios.get(`${API_BASE}/current-user/`)
      currentUser.value = res.data.user
      return true
    } catch (error) {
      localStorage.removeItem('token')
      currentUser.value = null
      return false
    }
  }

  return {
    // 状态
    currentUser,
    // 方法
    logout,
    restoreSession
  }
})
