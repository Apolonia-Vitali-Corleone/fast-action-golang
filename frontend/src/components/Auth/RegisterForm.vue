<template>
  <div class="register-container">
    <!-- 身份选择 -->
    <el-form-item label="身份" prop="role">
      <el-select
        v-model="registerRole"
        placeholder="请选择您的身份"
        size="large"
        @change="handleRoleChange"
      >
        <el-option label="学生" value="student" />
        <el-option label="老师" value="teacher" />
      </el-select>
    </el-form-item>

    <!-- 学生注册表单（手机号+短信验证码） -->
    <StudentRegisterForm
      v-if="registerRole === 'student'"
      @register-success="handleRegisterSuccess"
    />

    <!-- 教师注册表单（用户名+密码+邮箱） -->
    <el-form
      v-else-if="registerRole === 'teacher'"
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

      <el-form-item label="邮箱" prop="email">
        <el-input
          v-model="teacherForm.email"
          type="email"
          placeholder="请输入邮箱"
          size="large"
          clearable
        />
      </el-form-item>

      <el-form-item label="密码" prop="password">
        <el-input
          v-model="teacherForm.password"
          type="password"
          placeholder="请输入密码（至少6位）"
          size="large"
          show-password
        />
      </el-form-item>

      <el-button
        type="primary"
        size="large"
        class="submit-btn"
        :loading="loading"
        @click="handleTeacherRegister"
      >
        注册
      </el-button>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import StudentRegisterForm from './StudentRegisterForm.vue'

const API_BASE = 'http://localhost:8000/api'

const emit = defineEmits(['register-success'])
const registerRole = ref('')
const formRef = ref(null)
const loading = ref(false)

const teacherForm = reactive({
  username: '',
  email: '',
  password: ''
})

const teacherRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
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
  teacherForm.email = ''
  teacherForm.password = ''
}

const handleRegisterSuccess = (user) => {
  emit('register-success', user)
}

const handleTeacherRegister = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // 教师注册不直接返回token和用户信息，需要注册后再登录
        await axios.post(`${API_BASE}/teacher/register/`, {
          username: teacherForm.username,
          password: teacherForm.password,
          email: teacherForm.email
        })

        ElMessage.success('注册成功，正在为您自动登录...')

        // 自动登录
        const loginRes = await axios.post(`${API_BASE}/teacher/login/`, {
          username: teacherForm.username,
          password: teacherForm.password
        })

        // 保存token
        if (loginRes.data.token) {
          localStorage.setItem('token', loginRes.data.token)
        }

        emit('register-success', loginRes.data.user)
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '注册失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.register-container {
  width: 100%;
}

.teacher-form {
  width: 100%;
}

.submit-btn {
  width: 100%;
  margin-top: 10px;
  border-radius: 12px;
  font-weight: 600;
}

:deep(.el-input__wrapper),
:deep(.el-select__wrapper) {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 175, 240, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.el-input__wrapper:hover),
:deep(.el-select__wrapper:hover) {
  box-shadow: 0 4px 16px rgba(0, 175, 240, 0.15);
}

:deep(.el-input__wrapper.is-focus),
:deep(.el-select__wrapper.is-focused) {
  box-shadow: 0 4px 16px rgba(0, 175, 240, 0.25);
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #1c1c1c;
}
</style>
