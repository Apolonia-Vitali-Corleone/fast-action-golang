<template>
  <el-card class="form-card" shadow="hover">
    <template #header>
      <h2>创建新课程</h2>
    </template>

    <el-form
      ref="formRef"
      :model="courseStore.courseForm"
      :rules="rules"
      label-position="top"
    >
      <el-form-item label="课程名称" prop="name">
        <el-input
          v-model="courseStore.courseForm.name"
          placeholder="请输入课程名称"
          size="large"
          clearable
        />
      </el-form-item>

      <el-form-item label="课程描述" prop="description">
        <el-input
          v-model="courseStore.courseForm.description"
          type="textarea"
          :rows="4"
          placeholder="请输入课程描述"
          size="large"
        />
      </el-form-item>

      <el-form-item label="课程容量" prop="capacity">
        <el-input-number
          v-model="courseStore.courseForm.capacity"
          :min="1"
          :max="500"
          size="large"
          controls-position="right"
          class="capacity-input"
        />
      </el-form-item>

      <el-button
        type="primary"
        size="large"
        class="submit-btn"
        :loading="loading"
        @click="handleSubmit"
      >
        创建课程
      </el-button>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { useCourseStore } from '@/stores/course'

const courseStore = useCourseStore()
const formRef = ref(null)
const loading = ref(false)

const rules = {
  name: [
    { required: true, message: '请输入课程名称', trigger: 'blur' },
    { min: 2, max: 50, message: '课程名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  capacity: [
    { required: true, message: '请输入课程容量', trigger: 'blur' },
    { type: 'number', min: 1, message: '容量至少为 1', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const success = await courseStore.createCourse()
        if (success) {
          formRef.value.resetFields()
        }
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.form-card {
  border-radius: 16px;
  margin-bottom: 32px;
}

.form-card :deep(.el-card__header) {
  background: linear-gradient(135deg, rgba(0, 175, 240, 0.05) 0%, rgba(102, 126, 234, 0.05) 100%);
}

.form-card h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: #1c1c1c;
}

.capacity-input {
  width: 100%;
}

.submit-btn {
  width: 100%;
  margin-top: 10px;
  border-radius: 12px;
  font-weight: 600;
}

:deep(.el-input__wrapper),
:deep(.el-textarea__inner) {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 175, 240, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.el-input__wrapper:hover),
:deep(.el-textarea__inner:hover) {
  box-shadow: 0 4px 16px rgba(0, 175, 240, 0.15);
}

:deep(.el-input__wrapper.is-focus),
:deep(.el-textarea__inner:focus) {
  box-shadow: 0 4px 16px rgba(0, 175, 240, 0.25);
}

:deep(.el-form-item__label) {
  font-weight: 600;
  color: #1c1c1c;
}
</style>
