<template>
  <el-dialog
    :model-value="courseStore.showCourseDialog"
    :title="courseStore.isEditMode ? '编辑课程' : '创建新课程'"
    width="700px"
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

      <el-divider>课程时间安排</el-divider>

      <div class="schedules-container">
        <div
          v-for="(schedule, index) in courseStore.courseForm.schedules"
          :key="index"
          class="schedule-item"
        >
          <div class="schedule-header">
            <span class="schedule-title">时间段 {{ index + 1 }}</span>
            <el-button
              type="danger"
              size="small"
              text
              @click="removeSchedule(index)"
            >
              删除
            </el-button>
          </div>

          <el-row :gutter="12">
            <el-col :span="8">
              <el-form-item label="星期" :prop="`schedules.${index}.day_of_week`">
                <el-select
                  v-model="schedule.day_of_week"
                  placeholder="选择星期"
                  size="large"
                  class="full-width"
                >
                  <el-option label="周一" :value="1" />
                  <el-option label="周二" :value="2" />
                  <el-option label="周三" :value="3" />
                  <el-option label="周四" :value="4" />
                  <el-option label="周五" :value="5" />
                </el-select>
              </el-form-item>
            </el-col>

            <el-col :span="8">
              <el-form-item label="节次" :prop="`schedules.${index}.time_slot`">
                <el-select
                  v-model="schedule.time_slot"
                  placeholder="选择节次"
                  size="large"
                  class="full-width"
                >
                  <el-option label="上午第一节" :value="1" />
                  <el-option label="上午第二节" :value="2" />
                  <el-option label="下午第一节" :value="3" />
                  <el-option label="下午第二节" :value="4" />
                </el-select>
              </el-form-item>
            </el-col>

            <el-col :span="8">
              <el-form-item label="教室" :prop="`schedules.${index}.classroom`">
                <el-input
                  v-model="schedule.classroom"
                  placeholder="如: A101"
                  size="large"
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="12">
            <el-col :span="12">
              <el-form-item label="开始周次" :prop="`schedules.${index}.start_week`">
                <el-input-number
                  v-model="schedule.start_week"
                  :min="1"
                  :max="20"
                  size="large"
                  controls-position="right"
                  class="full-width"
                />
              </el-form-item>
            </el-col>

            <el-col :span="12">
              <el-form-item label="结束周次" :prop="`schedules.${index}.end_week`">
                <el-input-number
                  v-model="schedule.end_week"
                  :min="1"
                  :max="20"
                  size="large"
                  controls-position="right"
                  class="full-width"
                />
              </el-form-item>
            </el-col>
          </el-row>
        </div>

        <el-button
          type="primary"
          plain
          @click="addSchedule"
          class="add-schedule-btn"
        >
          + 添加上课时间
        </el-button>
      </div>
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
import { ElMessage } from 'element-plus'

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

const addSchedule = () => {
  courseStore.courseForm.schedules.push({
    day_of_week: null,
    time_slot: null,
    start_week: 1,
    end_week: 16,
    classroom: ''
  })
}

const removeSchedule = (index) => {
  courseStore.courseForm.schedules.splice(index, 1)
}

const validateSchedules = () => {
  const schedules = courseStore.courseForm.schedules

  for (let i = 0; i < schedules.length; i++) {
    const schedule = schedules[i]

    if (!schedule.day_of_week) {
      ElMessage.error(`时间段 ${i + 1}: 请选择星期`)
      return false
    }

    if (!schedule.time_slot) {
      ElMessage.error(`时间段 ${i + 1}: 请选择节次`)
      return false
    }

    if (!schedule.start_week || !schedule.end_week) {
      ElMessage.error(`时间段 ${i + 1}: 请设置周次范围`)
      return false
    }

    if (schedule.end_week < schedule.start_week) {
      ElMessage.error(`时间段 ${i + 1}: 结束周次不能小于开始周次`)
      return false
    }
  }

  return true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      // 验证课程时间
      if (!validateSchedules()) {
        return
      }

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
.capacity-input,
.full-width {
  width: 100%;
}

.schedules-container {
  margin-top: 16px;
}

.schedule-item {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 12px;
}

.schedule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.schedule-title {
  font-weight: 600;
  color: #303133;
}

.add-schedule-btn {
  width: 100%;
  margin-top: 8px;
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
