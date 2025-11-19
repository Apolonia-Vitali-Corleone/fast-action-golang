<template>
  <el-form :model="form" label-width="200px" style="max-width: 400px">
    <el-form-item label="手机号">
      <el-input v-model="form.phone"/>
    </el-form-item>
    <el-form-item label="密码">
      <el-input v-model="form.password"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="onSubmit">登录</el-button>
      <el-button>注册</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup>

import {reactive} from 'vue'

// do not use same name with ref
const form = reactive({
  phone: '',
  password: '',
})

const onSubmit = () => {
  console.log('submit!')
  const resp = axios.post('http://localhost:8000/api/student/login', form)
  console.log(resp)
}

import {useRouter} from 'vue-router'
import StudentLoginForm from '@/components/Auth/StudentLoginForm.vue'
import axios from "axios";

const router = useRouter()

const handleLoginSuccess = (user) => {
  // 保存用户信息到store（如果使用）或直接跳转
  // 由于登录表单已经保存了token，这里直接跳转即可
  router.push('/')
}
</script>

<style scoped>

</style>
