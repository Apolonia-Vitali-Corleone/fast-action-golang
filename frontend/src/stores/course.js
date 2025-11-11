/**
 * 课程状态管理 - 处理学生和教师的课程数据
 * 保持原有业务逻辑不变
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const API_BASE = 'http://localhost:8000/api'

export const useCourseStore = defineStore('course', () => {
  // ========== 学生相关状态 ==========
  const courses = ref([]) // 可选课程
  const myCourses = ref([]) // 我的课程

  // ========== 教师相关状态 ==========
  const teacherCourses = ref([]) // 教师的课程
  const courseForm = ref({ name: '', description: '', capacity: 50 })
  const currentCourse = ref({}) // 当前查看的课程
  const courseStudents = ref({ students: [], total: 0 }) // 课程学生
  const showStudentsDialog = ref(false) // 学生列表弹窗

  // ========== 学生方法 ==========

  /**
   * 获取可选课程列表
   */
  const fetchAvailableCourses = async () => {
    try {
      const res = await axios.get(`${API_BASE}/student/courses/`)
      courses.value = res.data.courses || []
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '获取课程失败')
    }
  }

  /**
   * 获取我的课程
   */
  const fetchMyCourses = async () => {
    try {
      const res = await axios.get(`${API_BASE}/student/my-courses/`)
      myCourses.value = res.data.courses || []
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '获取我的课程失败')
    }
  }

  /**
   * 选课
   */
  const enrollCourse = async (courseId) => {
    try {
      await axios.post(`${API_BASE}/student/enroll/`, { course_id: courseId })
      ElMessage.success('选课成功')
      await fetchAvailableCourses()
      return true
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '选课失败')
      return false
    }
  }

  /**
   * 退课
   */
  const dropCourse = async (courseId) => {
    try {
      await ElMessageBox.confirm('确定要退课吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await axios.post(`${API_BASE}/student/drop/`, { course_id: courseId })
      ElMessage.success('退课成功')
      await fetchMyCourses()
      return true
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error.response?.data?.error || '退课失败')
      }
      return false
    }
  }

  // ========== 教师方法 ==========

  /**
   * 获取教师的课程列表
   */
  const fetchTeacherCourses = async () => {
    try {
      const res = await axios.get(`${API_BASE}/teacher/courses/`)
      teacherCourses.value = res.data.courses || []
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '获取课程失败')
    }
  }

  /**
   * 创建课程
   */
  const createCourse = async () => {
    try {
      await axios.post(`${API_BASE}/teacher/courses/create/`, courseForm.value)
      ElMessage.success('创建成功')
      courseForm.value = { name: '', description: '', capacity: 50 }
      await fetchTeacherCourses()
      return true
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '创建失败')
      return false
    }
  }

  /**
   * 删除课程
   */
  const deleteCourse = async (courseId) => {
    try {
      await ElMessageBox.confirm('确定要删除这门课程吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await axios.delete(`${API_BASE}/teacher/courses/${courseId}/delete/`)
      ElMessage.success('删除成功')
      await fetchTeacherCourses()
      return true
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error(error.response?.data?.error || '删除失败')
      }
      return false
    }
  }

  /**
   * 查看选课学生
   */
  const viewStudents = async (courseId) => {
    try {
      const res = await axios.get(`${API_BASE}/teacher/courses/${courseId}/students/`)
      courseStudents.value = {
        students: res.data.students || [],
        total: res.data.total || 0
      }
      currentCourse.value = res.data.course || {}
      showStudentsDialog.value = true
      return true
    } catch (error) {
      ElMessage.error(error.response?.data?.error || '获取学生列表失败')
      return false
    }
  }

  /**
   * 关闭学生列表弹窗
   */
  const closeStudentsDialog = () => {
    showStudentsDialog.value = false
  }

  return {
    // 学生状态
    courses,
    myCourses,
    // 教师状态
    teacherCourses,
    courseForm,
    currentCourse,
    courseStudents,
    showStudentsDialog,
    // 学生方法
    fetchAvailableCourses,
    fetchMyCourses,
    enrollCourse,
    dropCourse,
    // 教师方法
    fetchTeacherCourses,
    createCourse,
    deleteCourse,
    viewStudents,
    closeStudentsDialog
  }
})
