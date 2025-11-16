/**
 * 用户状态管理 - 处理认证和用户信息
 * 保持原有业务逻辑不变
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

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
      const userStore = useUserStore()
      userStore.currentUser = null
      // 路由守卫会自动处理跳转到登录页
    }
    return Promise.reject(error)
  }
)

export const useUserStore = defineStore('user', () => {
  // 状态
  const currentUser = ref(null)
  const authForm = ref({ username: '', password: '', email: '', role: '' })
  const loginRole = ref('') // 登录时选择的角色

  // ========== 认证方法 ==========

  /**
   * 登录
   */
  const login = async () => {
    try {
      if (!loginRole.value) {
        ElMessage.warning('请选择登录身份')
        return false
      }

      const endpoint = loginRole.value === 'student' ? '/student/login/' : '/teacher/login/'
      const res = await axios.post(`${API_BASE}${endpoint}`, {
        username: authForm.value.username,
        password: authForm.value.password
      })

      if (res.data.token) {
        localStorage.setItem('token', res.data.token)
      }

      currentUser.value = res.data.user
      ElMessage.success('登录成功')

      // 重置表单
      authForm.value = { username: '', password: '', email: '', role: '' }
      loginRole.value = ''

      return true
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '登录失败')
      return false
    }
  }

  /**
   * 注册 - 注册成功后自动登录
   */
  const register = async () => {
    try {
      if (!authForm.value.role) {
        ElMessage.warning('请选择注册身份')
        return false
      }

      // 保存注册信息用于自动登录
      const username = authForm.value.username
      const password = authForm.value.password
      const role = authForm.value.role

      const endpoint = role === 'student' ? '/student/register/' : '/teacher/register/'
      await axios.post(`${API_BASE}${endpoint}`, {
        username: authForm.value.username,
        password: authForm.value.password,
        email: authForm.value.email
      })

      // 注册成功，自动登录
      ElMessage.success('注册成功')

      // 准备登录数据
      authForm.value = { username, password, email: '', role: '' }
      loginRole.value = role

      // 自动登录
      const loginEndpoint = role === 'student' ? '/student/login/' : '/teacher/login/'
      const loginRes = await axios.post(`${API_BASE}${loginEndpoint}`, {
        username,
        password
      })

      if (loginRes.data.token) {
        localStorage.setItem('token', loginRes.data.token)
      }

      currentUser.value = loginRes.data.user

      // 清空表单
      authForm.value = { username: '', password: '', email: '', role: '' }
      loginRole.value = ''

      return true
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '注册失败')
      return false
    }
  }

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
    authForm,
    loginRole,
    // 方法
    login,
    register,
    logout,
    restoreSession
  }
})
