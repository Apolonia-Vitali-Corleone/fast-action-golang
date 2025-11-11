<template>
  <div class="course-list-container">
    <div class="list-header">
      <h2>所有课程</h2>
      <el-button
        type="primary"
        :icon="Refresh"
        circle
        @click="handleRefresh"
        :loading="loading"
      />
    </div>

    <el-empty
      v-if="courseStore.courses.length === 0"
      description="暂无课程"
      :image-size="160"
    />

    <div v-else class="course-grid">
      <CourseCard
        v-for="course in courseStore.courses"
        :key="course.id"
        :course="course"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { useCourseStore } from '@/stores/course'
import CourseCard from './CourseCard.vue'

const courseStore = useCourseStore()
const loading = ref(false)
let refreshTimer = null

const handleRefresh = async () => {
  loading.value = true
  try {
    await courseStore.fetchAvailableCourses()
  } finally {
    loading.value = false
  }
}

// 自动刷新功能
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    courseStore.fetchAvailableCourses()
  }, 30000) // 每30秒刷新一次
}

onMounted(() => {
  courseStore.fetchAvailableCourses()
  startAutoRefresh()
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.course-list-container {
  width: 100%;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.list-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #1c1c1c;
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

@media (max-width: 768px) {
  .course-grid {
    grid-template-columns: 1fr;
  }
}
</style>
