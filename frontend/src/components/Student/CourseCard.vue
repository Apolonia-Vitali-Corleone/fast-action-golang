<template>
  <el-card class="course-card" shadow="hover">
    <template #header>
      <div class="card-header">
        <h3>{{ course.name }}</h3>
        <el-tag v-if="course.is_enrolled" type="success" round>已选</el-tag>
        <el-tag v-else-if="course.is_full" type="info" round>已满</el-tag>
      </div>
    </template>

    <p class="description">{{ course.description || '暂无描述' }}</p>

    <div class="course-info">
      <div class="info-item">
        <el-icon><User /></el-icon>
        <span>{{ course.teacher }}</span>
      </div>
      <div class="info-item">
        <el-icon><ChatDotRound /></el-icon>
        <span>{{ course.enrolled }}/{{ course.capacity }}</span>
      </div>
    </div>

    <el-button
      v-if="!course.is_enrolled"
      type="primary"
      :disabled="course.is_full"
      :loading="loading"
      class="enroll-btn"
      round
      @click="handleEnroll"
    >
      {{ course.is_full ? '课程已满' : '选择课程' }}
    </el-button>
  </el-card>
</template>

<script setup>
import { ref } from 'vue'
import { User, ChatDotRound } from '@element-plus/icons-vue'
import { useCourseStore } from '@/stores/course'

const props = defineProps({
  course: {
    type: Object,
    required: true
  }
})

const courseStore = useCourseStore()
const loading = ref(false)

const handleEnroll = async () => {
  loading.value = true
  try {
    await courseStore.enrollCourse(props.course.id)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.course-card {
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid transparent;
}

.course-card:hover {
  transform: translateY(-4px);
  border-color: #00AFF0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1c1c1c;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.description {
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
  margin-bottom: 16px;
  min-height: 44px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.course-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
  padding: 12px;
  background: #f7f7f7;
  border-radius: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #909399;
  font-size: 14px;
}

.info-item .el-icon {
  color: #00AFF0;
}

.enroll-btn {
  width: 100%;
  font-weight: 600;
}

:deep(.el-card__header) {
  padding: 16px 20px;
  background: linear-gradient(135deg, rgba(0, 175, 240, 0.05) 0%, rgba(102, 126, 234, 0.05) 100%);
}

:deep(.el-card__body) {
  padding: 20px;
}
</style>
