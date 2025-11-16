<template>
  <el-form
    ref="formRef"
    :model="userStore.authForm"
    :rules="rules"
    label-position="top"
    class="login-form"
  >
    <el-form-item label="身份">
      <el-select
        v-model="userStore.loginRole"
        placeholder="请选择您的身份"
        size="large"
      >
        <el-option label="学生" value="student" />
        <el-option label="老师" value="teacher" />
      </el-select>
    </el-form-item>

    <el-form-item label="用户名" prop="username">
      <el-input
        v-model="userStore.authForm.username"
        placeholder="请输入用户名"
        size="large"
        clearable
      />
    </el-form-item>

    <el-form-item label="密码" prop="password">
      <el-input
        v-model="userStore.authForm.password"
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
      @click="handleSubmit"
    >
      登录
    </el-button>
  </el-form>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const emit = defineEmits(['login-success'])
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 20, message: '用户名长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不少于 6 个字符', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  // 先手动验证身份选择
  if (!userStore.loginRole) {
    ElMessage.warning('请选择登录身份')
    return
  }

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const success = await userStore.login()
        if (success) {
          emit('login-success')
        }
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-form {
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
