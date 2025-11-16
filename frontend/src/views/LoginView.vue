<template>
  <div class="auth-container">
    <div class="auth-card">
      <h1 class="title">选课系统</h1>
      <p class="subtitle">欢迎回来</p>

      <LoginForm @login-success="handleLoginSuccess" />

      <div class="auth-footer">
        <span class="footer-text">还没有账号？</span>
        <router-link to="/register" class="footer-link">注册一个</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import LoginForm from '@/components/Auth/LoginForm.vue'

const router = useRouter()

const handleLoginSuccess = () => {
  router.push('/')
}
</script>

<style scoped>
.auth-container {
  height: 100vh;
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
  position: relative;
  z-index: 1;
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
  margin-bottom: 8px;
  background: linear-gradient(135deg, #00AFF0 0%, #667eea 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  text-align: center;
  font-size: 16px;
  color: #606266;
  margin-bottom: 36px;
}

.auth-footer {
  margin-top: 24px;
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

.footer-text {
  color: #606266;
  font-size: 14px;
  margin-right: 8px;
}

.footer-link {
  color: #00AFF0;
  font-weight: 600;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.3s;
}

.footer-link:hover {
  color: #667eea;
  text-decoration: underline;
}

/* 响应式 */
@media (max-width: 600px) {
  .auth-card {
    padding: 32px 24px;
  }

  .title {
    font-size: 28px;
  }

  .subtitle {
    font-size: 14px;
    margin-bottom: 24px;
  }
}
</style>
