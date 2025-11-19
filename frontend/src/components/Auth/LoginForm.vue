<template>
  <div class="login-container">
    <!-- 身份选择 -->
    <el-form-item label="身份">
      <el-select
        v-model="loginRole"
        placeholder="请选择您的身份"
        size="large"
        @change="handleRoleChange"
      >
        <el-option label="学生" value="student" />
        <el-option label="老师" value="teacher" />
      </el-select>
    </el-form-item>

    <!-- 学生登录表单（短信验证码） -->
    <StudentLoginForm
      v-if="loginRole === 'student'"
      @login-success="handleLoginSuccess"
    />

    <!-- 教师登录表单（用户名密码） -->
    <el-form
      v-else-if="loginRole === 'teacher'"
      ref="formRef"
      :model="teacherForm"
      :rules="teacherRules"
      label-position="top"
      class="teacher-form"
    >
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="teacherForm.username"
          placeholder="请输入用户名"
          size="large"
          clearable
        />
      </el-form-item>

      <el-form-item label="密码" prop="password">
        <el-input
          v-model="teacherForm.password"
          type="password"
          placeholder="请输入密码"
          size="large"
          show-password
        />
      </el-form-item>

      <el-button
        type="primary"
        size="large"
        class="submit-btn"
        :loading="loading"
        @click="handleTeacherLogin"
      >
        登录
      </el-button>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import StudentLoginForm from './StudentLoginForm.vue'

const API_BASE = 'http://localhost:8000/api'

const emit = defineEmits(['login-success'])
const loginRole = ref('')
const formRef = ref(null)
const loading = ref(false)

const teacherForm = reactive({
  username: '',
  password: ''
})

const teacherRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不少于 6 个字符', trigger: 'blur' }
  ]
}

const handleRoleChange = () => {
  // 切换身份时清空表单
  if (formRef.value) {
    formRef.value.resetFields()
  }
  teacherForm.username = ''
  teacherForm.password = ''
}

const handleLoginSuccess = (user) => {
  emit('login-success', user)
}

const handleTeacherLogin = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await axios.post(`${API_BASE}/teacher/login/`, {
          username: teacherForm.username,
          password: teacherForm.password
        })

        // 保存token
        if (res.data.token) {
          localStorage.setItem('token', res.data.token)
        }

        ElMessage.success('登录成功')
        emit('login-success', res.data.user)
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>

</style>
