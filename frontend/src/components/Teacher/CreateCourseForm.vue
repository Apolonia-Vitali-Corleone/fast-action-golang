<template>
  <el-dialog
    :model-value="courseStore.showCourseDialog"
    :title="courseStore.isEditMode ? '编辑课程' : '创建新课程'"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
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
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="loading"
        @click="handleSubmit"
      >
        {{ courseStore.isEditMode ? '保存' : '创建' }}
      </el-button>
    </template>
  </el-dialog>
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
        if (courseStore.isEditMode) {
          await courseStore.updateCourse()
        } else {
          await courseStore.createCourse()
        }
        formRef.value.resetFields()
      } finally {
        loading.value = false
      }
    }
  })
}

const handleClose = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  courseStore.closeCourseDialog()
}
</script>

<style scoped>
.capacity-input {
  width: 100%;
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
