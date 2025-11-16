<template>
  <div class="auth-container">
    <div class="auth-card">
      <h1 class="title">选课系统</h1>

      <el-tabs v-model="activeTab" class="auth-tabs" stretch>
        <el-tab-pane label="登录" name="login">
          <LoginForm />
        </el-tab-pane>
        <el-tab-pane label="注册" name="register">
          <RegisterForm @register-success="handleRegisterSuccess" />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import LoginForm from '@/components/Auth/LoginForm.vue'
import RegisterForm from '@/components/Auth/RegisterForm.vue'

const activeTab = ref('login')

const handleRegisterSuccess = () => {
  activeTab.value = 'login'
}
</script>

<style scoped>
.auth-container {
  height: 100%;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow-y: auto;
  overflow-x: hidden;
}

/* 背景装饰 */
.auth-container::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(0, 175, 240, 0.1) 0%, transparent 70%);
  animation: float 20s infinite ease-in-out;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-50px, 50px); }
}

.auth-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 48px 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 450px;
  animation: slideIn 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.title {
  text-align: center;
  font-size: 32px;
  font-weight: 700;
  color: #1c1c1c;
  margin-bottom: 36px;
  background: linear-gradient(135deg, #00AFF0 0%, #667eea 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.auth-tabs {
  margin-top: 20px;
}

:deep(.el-tabs__header) {
  margin-bottom: 32px;
}

:deep(.el-tabs__nav-wrap::after) {
  display: none;
}

:deep(.el-tabs__active-bar) {
  height: 3px;
  border-radius: 3px;
  background: linear-gradient(90deg, #00AFF0 0%, #667eea 100%);
}

:deep(.el-tabs__item) {
  font-size: 16px;
  font-weight: 600;
  color: #909399;
  transition: all 0.3s;
}

:deep(.el-tabs__item:hover) {
  color: #00AFF0;
}

:deep(.el-tabs__item.is-active) {
  color: #00AFF0;
}

/* 响应式 */
@media (max-width: 600px) {
  .auth-card {
    padding: 32px 24px;
  }

  .title {
    font-size: 28px;
    margin-bottom: 24px;
  }
}
</style>
