<template>
  <div class="auth-container">
    <div class="auth-card">
      <h1 class="title">选课系统</h1>

      <el-tabs v-model="activeTab" class="auth-tabs" stretch>
        <el-tab-pane label="登录" name="login">
          <LoginForm/>
        </el-tab-pane>
        <el-tab-pane label="注册" name="register">
          <RegisterForm @register-success="handleRegisterSuccess"/>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import {ref} from 'vue'
import LoginForm from '@/components/Auth/LoginForm.vue'
import RegisterForm from '@/components/Auth/RegisterForm.vue'

const activeTab = ref('login')

const handleRegisterSuccess = () => {
  activeTab.value = 'login'
}
</script>

<style scoped>
</style>
