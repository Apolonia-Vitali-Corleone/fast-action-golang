<template>
  <div class="captcha-container">
    <el-input
      v-model="captchaCode"
      placeholder="请输入图形验证码"
      size="large"
      :maxlength="4"
      clearable
      @input="handleInput"
    >
      <template #append>
        <div class="captcha-image-wrapper" @click="refreshCaptcha">
          <el-image
            v-if="captchaImage"
            :src="captchaImage"
            fit="contain"
            class="captcha-image"
          />
          <span v-else>加载中...</span>
        </div>
      </template>
    </el-input>
    <div class="captcha-hint">点击图片刷新验证码</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const API_BASE = 'http://localhost:8000/api'

const emit = defineEmits(['update:captchaCode', 'update:captchaId'])
const props = defineProps({
  captchaCode: String,
  captchaId: String
})

const captchaCode = ref('')
const captchaImage = ref('')
const captchaId = ref('')

// 获取图形验证码
const refreshCaptcha = async () => {
  try {
    const res = await axios.get(`${API_BASE}/captcha/`)
    captchaId.value = res.data.captcha_id
    captchaImage.value = res.data.image
    emit('update:captchaId', res.data.captcha_id)
  } catch (error) {
    ElMessage.error('获取验证码失败')
  }
}

const handleInput = (value) => {
  emit('update:captchaCode', value)
}

onMounted(() => {
  refreshCaptcha()
})

defineExpose({ refreshCaptcha })
</script>

<style scoped>

</style>
