<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-position="top"
    class="login-form"
  >
    <el-form-item label="手机号" prop="phone">
      <el-input
        v-model="form.phone"
        placeholder="请输入手机号"
        size="large"
        maxlength="11"
        clearable
      />
    </el-form-item>

    <el-form-item label="图形验证码" prop="captchaCode">
      <CaptchaInput
        v-model:captcha-code="form.captchaCode"
        v-model:captcha-id="form.captchaId"
      />
    </el-form-item>

    <el-form-item label="短信验证码" prop="smsCode">
      <div class="sms-code-input">
        <el-input
          v-model="form.smsCode"
          placeholder="请输入短信验证码"
          size="large"
          maxlength="6"
          clearable
        />
        <el-button
          type="primary"
          size="large"
          :disabled="countdown > 0 || !form.phone || !form.captchaCode"
          :loading="sending"
          @click="sendSMSCode"
          class="send-btn"
        >
          {{ countdown > 0 ? `${countdown}秒后重试` : '发送验证码' }}
        </el-button>
      </div>
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
import { ref, reactive } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import CaptchaInput from './CaptchaInput.vue'

const API_BASE = 'http://localhost:8000/api'

const emit = defineEmits(['login-success'])
const formRef = ref(null)
const loading = ref(false)
const sending = ref(false)
const countdown = ref(0)

const form = reactive({
  phone: '',
  captchaId: '',
  captchaCode: '',
  smsCode: ''
})

const rules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  captchaCode: [
    { required: true, message: '请输入图形验证码', trigger: 'blur' },
    { len: 4, message: '验证码为4位', trigger: 'blur' }
  ],
  smsCode: [
    { required: true, message: '请输入短信验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

// 发送短信验证码
const sendSMSCode = async () => {
  // 先验证手机号和图形验证码
  try {
    await formRef.value.validateField(['phone', 'captchaCode'])
  } catch {
    return
  }

  sending.value = true
  try {
    await axios.post(`${API_BASE}/sms/send/`, {
      phone: form.phone,
      purpose: 'login',
      captcha_id: form.captchaId,
      captcha_code: form.captchaCode
    })

    ElMessage.success('验证码已发送，请查收短信')

    // 开始倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '发送验证码失败')
  } finally {
    sending.value = false
  }
}

// 提交登录
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await axios.post(`${API_BASE}/student/login/`, {
          phone: form.phone,
          sms_code: form.smsCode
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
.login-form {
  width: 100%;
}

.sms-code-input {
  display: flex;
  gap: 8px;
}

.send-btn {
  flex-shrink: 0;
  min-width: 120px;
}

.submit-btn {
  width: 100%;
  margin-top: 10px;
  border-radius: 12px;
  font-weight: 600;
}

:deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 175, 240, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 4px 16px rgba(0, 175, 240, 0.15);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 4px 16px rgba(0, 175, 240, 0.25);
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #1c1c1c;
}
</style>
