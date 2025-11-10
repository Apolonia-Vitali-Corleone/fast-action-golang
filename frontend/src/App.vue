<template>
  <div class="app">
    <!-- 登录注册页面 -->
    <div v-if="!currentUser" class="auth-container">
      <h1>选课系统</h1>

      <div class="auth-box">
        <!-- 切换登录/注册 -->
        <div class="tabs">
          <button
            :class="['tab', { active: authMode === 'login' }]"
            @click="authMode = 'login'"
          >
            登录
          </button>
          <button
            :class="['tab', { active: authMode === 'register' }]"
            @click="authMode = 'register'"
          >
            注册
          </button>
        </div>

        <!-- 登录表单 -->
        <form v-if="authMode === 'login'" @submit.prevent="handleLogin" class="form">
          <div class="form-group">
            <label>身份</label>
            <select v-model="loginRole" required>
              <option value="">请选择</option>
              <option value="student">学生</option>
              <option value="teacher">老师</option>
            </select>
          </div>
          <div class="form-group">
            <label>用户名</label>
            <input v-model="authForm.username" required>
          </div>
          <div class="form-group">
            <label>密码</label>
            <input v-model="authForm.password" type="password" required>
          </div>
          <button type="submit" class="btn btn-primary">登录</button>
        </form>

        <!-- 注册表单 -->
        <form v-if="authMode === 'register'" @submit.prevent="handleRegister" class="form">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="authForm.username" required>
          </div>
          <div class="form-group">
            <label>邮箱</label>
            <input v-model="authForm.email" type="email" required>
          </div>
          <div class="form-group">
            <label>密码</label>
            <input v-model="authForm.password" type="password" required>
          </div>
          <div class="form-group">
            <label>身份</label>
            <select v-model="authForm.role" required>
              <option value="">请选择</option>
              <option value="student">学生</option>
              <option value="teacher">老师</option>
            </select>
          </div>
          <button type="submit" class="btn btn-primary">注册</button>
        </form>
      </div>
    </div>

    <!-- 学生界面 -->
    <div v-else-if="currentUser.role === 'student'" class="main-container">
      <div class="header">
        <h1>学生选课系统</h1>
        <div class="user-info">
          <span>欢迎，{{ currentUser.username }}</span>
          <button @click="handleLogout" class="btn btn-secondary">退出</button>
        </div>
      </div>

      <!-- 切换视图 -->
      <div class="tabs">
        <button
          :class="['tab', { active: studentView === 'courses' }]"
          @click="studentView = 'courses'; fetchAvailableCourses()"
        >
          可选课程
        </button>
        <button
          :class="['tab', { active: studentView === 'my' }]"
          @click="studentView = 'my'; fetchMyCourses()"
        >
          我的课程
        </button>
      </div>

      <!-- 可选课程列表 -->
      <div v-if="studentView === 'courses'" class="content">
        <h2>所有课程</h2>
        <p v-if="courses.length === 0" class="empty">暂无课程</p>
        <div v-else class="course-grid">
          <div v-for="course in courses" :key="course.id" class="course-card">
            <h3>{{ course.name }}</h3>
            <p class="description">{{ course.description || '暂无描述' }}</p>
            <div class="course-info">
              <span>教师：{{ course.teacher }}</span>
              <span>人数：{{ course.enrolled }}/{{ course.capacity }}</span>
            </div>
            <button
              v-if="!course.is_enrolled"
              @click="enrollCourse(course.id)"
              :disabled="course.is_full"
              class="btn btn-primary"
            >
              {{ course.is_full ? '已满' : '选课' }}
            </button>
            <span v-else class="badge">已选</span>
          </div>
        </div>
      </div>

      <!-- 我的课程列表 -->
      <div v-if="studentView === 'my'" class="content">
        <h2>我的课程</h2>
        <p v-if="myCourses.length === 0" class="empty">您还没有选课</p>
        <table v-else>
          <thead>
            <tr>
              <th>课程名称</th>
              <th>描述</th>
              <th>教师</th>
              <th>选课时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="course in myCourses" :key="course.course_id">
              <td>{{ course.course_name }}</td>
              <td>{{ course.description || '暂无描述' }}</td>
              <td>{{ course.teacher }}</td>
              <td>{{ course.enrolled_at }}</td>
              <td>
                <button @click="dropCourse(course.course_id)" class="btn btn-delete">退课</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 教师界面 -->
    <div v-else-if="currentUser.role === 'teacher'" class="main-container">
      <div class="header">
        <h1>教师管理系统</h1>
        <div class="user-info">
          <span>欢迎，{{ currentUser.username }}</span>
          <button @click="handleLogout" class="btn btn-secondary">退出</button>
        </div>
      </div>

      <div class="content">
        <!-- 创建课程表单 -->
        <div class="form-box">
          <h2>创建新课程</h2>
          <form @submit.prevent="createCourse" class="form">
            <div class="form-group">
              <label>课程名称</label>
              <input v-model="courseForm.name" required>
            </div>
            <div class="form-group">
              <label>课程描述</label>
              <textarea v-model="courseForm.description" rows="3"></textarea>
            </div>
            <div class="form-group">
              <label>容量</label>
              <input v-model.number="courseForm.capacity" type="number" min="1" required>
            </div>
            <button type="submit" class="btn btn-primary">创建课程</button>
          </form>
        </div>

        <!-- 我的课程列表 -->
        <div class="table-box">
          <h2>我的课程</h2>
          <p v-if="teacherCourses.length === 0" class="empty">您还没有创建课程</p>
          <table v-else>
            <thead>
              <tr>
                <th>课程名称</th>
                <th>描述</th>
                <th>容量</th>
                <th>已选人数</th>
                <th>创建时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="course in teacherCourses" :key="course.id">
                <td>{{ course.name }}</td>
                <td>{{ course.description || '暂无描述' }}</td>
                <td>{{ course.capacity }}</td>
                <td>{{ course.enrolled }}</td>
                <td>{{ course.created_at }}</td>
                <td>
                  <button @click="viewStudents(course.id)" class="btn btn-edit">查看学生</button>
                  <button @click="deleteCourse(course.id)" class="btn btn-delete">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 查看选课学生弹窗 -->
        <div v-if="showStudentsModal" class="modal" @click="showStudentsModal = false">
          <div class="modal-content" @click.stop>
            <h2>选课学生名单</h2>
            <p><strong>课程：</strong>{{ currentCourse.name }}</p>
            <p><strong>总人数：</strong>{{ courseStudents.total }}</p>
            <table v-if="courseStudents.students.length > 0">
              <thead>
                <tr>
                  <th>学生ID</th>
                  <th>用户名</th>
                  <th>邮箱</th>
                  <th>选课时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="student in courseStudents.students" :key="student.id">
                  <td>{{ student.id }}</td>
                  <td>{{ student.username }}</td>
                  <td>{{ student.email }}</td>
                  <td>{{ student.enrolled_at }}</td>
                </tr>
              </tbody>
            </table>
            <p v-else class="empty">暂无学生选课</p>
            <button @click="showStudentsModal = false" class="btn btn-secondary">关闭</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

// API配置
const API_BASE = 'http://localhost:8000/api'
axios.defaults.withCredentials = true  // 支持Session

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

// 配置Axios响应拦截器 - 处理token刷新
axios.interceptors.response.use(
  (response) => {
    // 检查响应头中是否有新的token
    const newToken = response.headers['x-new-token']
    if (newToken) {
      localStorage.setItem('token', newToken)
    }
    return response
  },
  (error) => {
    // 处理401错误 - token过期或无效
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      currentUser.value = null
      alert('登录已过期，请重新登录')
    }
    return Promise.reject(error)
  }
)

// 全局状态
const currentUser = ref(null)
const authMode = ref('login')
const studentView = ref('courses')
const loginRole = ref('')  // 登录时选择的角色

// 表单数据
const authForm = ref({ username: '', password: '', email: '', role: '' })
const courseForm = ref({ name: '', description: '', capacity: 50 })

// 数据列表
const courses = ref([])  // 可选课程
const myCourses = ref([])  // 我的课程
const teacherCourses = ref([])  // 教师课程
const courseStudents = ref({ students: [], total: 0 })  // 课程学生
const currentCourse = ref({})  // 当前查看的课程
const showStudentsModal = ref(false)

// ========== 认证相关 ==========

const handleLogin = async () => {
  try {
    // 登录时需要先选择身份
    if (!loginRole.value) {
      alert('请选择登录身份')
      return
    }

    // 根据角色调用不同的登录接口
    const endpoint = loginRole.value === 'student' ? '/student/login/' : '/teacher/login/'
    const res = await axios.post(`${API_BASE}${endpoint}`, {
      username: authForm.value.username,
      password: authForm.value.password
    })

    // 保存JWT token到localStorage
    if (res.data.token) {
      localStorage.setItem('token', res.data.token)
    }

    currentUser.value = res.data.user
    alert('登录成功')
    authForm.value = { username: '', password: '', email: '', role: '' }
    loginRole.value = ''

    // 根据角色加载数据
    if (currentUser.value.role === 'student') {
      fetchAvailableCourses()
    } else {
      fetchTeacherCourses()
    }
  } catch (error) {
    alert(error.response?.data?.error || '登录失败')
  }
}

const handleRegister = async () => {
  try {
    if (!authForm.value.role) {
      alert('请选择注册身份')
      return
    }

    // 根据角色调用不同的注册接口
    const endpoint = authForm.value.role === 'student' ? '/student/register/' : '/teacher/register/'
    await axios.post(`${API_BASE}${endpoint}`, {
      username: authForm.value.username,
      password: authForm.value.password,
      email: authForm.value.email
    })

    alert('注册成功，请登录')
    authMode.value = 'login'
    authForm.value = { username: '', password: '', email: '', role: '' }
  } catch (error) {
    alert(error.response?.data?.error || '注册失败')
  }
}

const handleLogout = async () => {
  try {
    await axios.post(`${API_BASE}/logout/`)
    // 清除localStorage中的token
    localStorage.removeItem('token')
    currentUser.value = null
    alert('已退出')
  } catch (error) {
    // 即使请求失败也要清除本地token
    localStorage.removeItem('token')
    currentUser.value = null
    console.error(error)
  }
}

// ========== 学生相关 ==========

const fetchAvailableCourses = async () => {
  try {
    const res = await axios.get(`${API_BASE}/student/courses/`)
    courses.value = res.data.courses || []
  } catch (error) {
    alert(error.response?.data?.error || '获取课程失败')
  }
}

const fetchMyCourses = async () => {
  try {
    const res = await axios.get(`${API_BASE}/student/my-courses/`)
    myCourses.value = res.data.courses || []
  } catch (error) {
    alert(error.response?.data?.error || '获取我的课程失败')
  }
}

const enrollCourse = async (courseId) => {
  try {
    await axios.post(`${API_BASE}/student/enroll/`, { course_id: courseId })
    alert('选课成功')
    fetchAvailableCourses()
  } catch (error) {
    alert(error.response?.data?.error || '选课失败')
  }
}

const dropCourse = async (courseId) => {
  if (!confirm('确定要退课吗？')) return
  try {
    await axios.post(`${API_BASE}/student/drop/`, { course_id: courseId })
    alert('退课成功')
    fetchMyCourses()
  } catch (error) {
    alert(error.response?.data?.error || '退课失败')
  }
}

// ========== 教师相关 ==========

const fetchTeacherCourses = async () => {
  try {
    const res = await axios.get(`${API_BASE}/teacher/courses/`)
    teacherCourses.value = res.data.courses || []
  } catch (error) {
    alert(error.response?.data?.error || '获取课程失败')
  }
}

const createCourse = async () => {
  try {
    await axios.post(`${API_BASE}/teacher/courses/create/`, courseForm.value)
    alert('创建成功')
    courseForm.value = { name: '', description: '', capacity: 50 }
    fetchTeacherCourses()
  } catch (error) {
    alert(error.response?.data?.error || '创建失败')
  }
}

const deleteCourse = async (courseId) => {
  if (!confirm('确定要删除这门课程吗？')) return
  try {
    await axios.delete(`${API_BASE}/teacher/courses/${courseId}/delete/`)
    alert('删除成功')
    fetchTeacherCourses()
  } catch (error) {
    alert(error.response?.data?.error || '删除失败')
  }
}

const viewStudents = async (courseId) => {
  try {
    const res = await axios.get(`${API_BASE}/teacher/courses/${courseId}/students/`)
    courseStudents.value = {
      students: res.data.students || [],
      total: res.data.total || 0
    }
    currentCourse.value = res.data.course || {}
    showStudentsModal.value = true
  } catch (error) {
    alert(error.response?.data?.error || '获取学生列表失败')
  }
}

// ========== 生命周期 ==========

onMounted(async () => {
  // 检查localStorage中是否有token
  const savedToken = localStorage.getItem('token')
  if (!savedToken) {
    // 没有token，显示登录页面
    return
  }

  // 有token，尝试恢复登录状态
  try {
    const res = await axios.get(`${API_BASE}/current-user/`)
    currentUser.value = res.data.user

    // 根据角色加载数据
    if (currentUser.value.role === 'student') {
      fetchAvailableCourses()
    } else {
      fetchTeacherCourses()
    }
  } catch (error) {
    // token无效或过期，清除它
    localStorage.removeItem('token')
    currentUser.value = null
  }
})
</script>

<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.app {
  min-height: 100vh;
  background: #f5f5f5;
}

/* ========== 登录注册页面 ========== */
.auth-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
}

.auth-container h1 {
  margin-bottom: 30px;
  color: #333;
}

.auth-box {
  background: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  width: 100%;
  max-width: 400px;
}

/* ========== 主容器 ========== */
.main-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  background: white;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}

.header h1 {
  color: #333;
  font-size: 24px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.content {
  background: white;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}

/* ========== 标签页 ========== */
.tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.tab {
  padding: 10px 20px;
  border: none;
  background: #e0e0e0;
  cursor: pointer;
  border-radius: 5px;
  font-size: 14px;
  transition: all 0.3s;
}

.tab.active {
  background: #4CAF50;
  color: white;
}

.tab:hover {
  background: #d0d0d0;
}

.tab.active:hover {
  background: #45a049;
}

/* ========== 表单 ========== */
.form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group label {
  margin-bottom: 5px;
  font-weight: bold;
  color: #555;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #4CAF50;
}

.form-box {
  margin-bottom: 30px;
  padding: 20px;
  background: #f9f9f9;
  border-radius: 5px;
}

.form-box h2 {
  margin-bottom: 15px;
  color: #555;
}

/* ========== 按钮 ========== */
.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s;
}

.btn-primary {
  background: #4CAF50;
  color: white;
}

.btn-primary:hover {
  background: #45a049;
}

.btn-primary:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.btn-secondary {
  background: #999;
  color: white;
}

.btn-secondary:hover {
  background: #888;
}

.btn-edit {
  background: #2196F3;
  color: white;
  margin-right: 5px;
}

.btn-edit:hover {
  background: #0b7dda;
}

.btn-delete {
  background: #f44336;
  color: white;
}

.btn-delete:hover {
  background: #da190b;
}

/* ========== 课程卡片 ========== */
.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.course-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  background: #fafafa;
  transition: box-shadow 0.3s;
}

.course-card:hover {
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.course-card h3 {
  color: #333;
  margin-bottom: 10px;
}

.course-card .description {
  color: #666;
  margin-bottom: 15px;
  font-size: 14px;
  min-height: 40px;
}

.course-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  font-size: 14px;
  color: #777;
}

.badge {
  display: inline-block;
  padding: 8px 16px;
  background: #4CAF50;
  color: white;
  border-radius: 4px;
  font-size: 14px;
}

/* ========== 表格 ========== */
.table-box {
  margin-top: 30px;
}

.table-box h2 {
  margin-bottom: 15px;
  color: #555;
  border-bottom: 2px solid #4CAF50;
  padding-bottom: 10px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 15px;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

th {
  background: #4CAF50;
  color: white;
  font-weight: bold;
}

tr:hover {
  background: #f5f5f5;
}

/* ========== 弹窗 ========== */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 30px;
  border-radius: 8px;
  max-width: 800px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-content h2 {
  margin-bottom: 20px;
  color: #333;
}

.modal-content p {
  margin-bottom: 10px;
}

/* ========== 其他 ========== */
.empty {
  text-align: center;
  color: #999;
  padding: 40px 0;
}

h2 {
  color: #555;
  margin-bottom: 20px;
}
</style>
