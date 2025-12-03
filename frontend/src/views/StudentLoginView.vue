<template>
  <el-form :model="form" label-width="200px" style="max-width: 400px">
    <el-form-item>
      <el-text class="mx-1" type="primary">学生登录</el-text>
    </el-form-item>
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
import axios from "axios";
import {ElMessage} from 'element-plus'
import {useRouter} from 'vue-router'
import StudentLoginForm from '@/components/Auth/StudentLoginForm.vue'

const router = useRouter()

const form = reactive({
  phone: '',
  password: '',
})

const onSubmit = async () => {
  try {
    const resp = await axios.post('http://localhost:8000/api/student/login', form)

    console.log(resp)

    if (resp.status === 200) {
      ElMessage.success("登录成功！")

      localStorage.setItem('role', resp.data.user.role)
      localStorage.setItem('token', resp.data.token)

      handleLoginSuccess()
    } else {
      ElMessage.error(resp.data.msg || "登录失败，请检查账号密码！")
    }
  } catch (error) {
    ElMessage.error("网络异常或服务器错误，请稍后再试！")
    console.log("登录请求失败" + error)
  }
}

const handleLoginSuccess = () => {
  router.push('/')
}
</script>

<style scoped>

</style>
