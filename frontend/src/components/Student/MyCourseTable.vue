<template>
  <div class="my-course-container">
    <div class="table-header">
      <h2>我的课程</h2>
      <el-button
        type="primary"
        :icon="Refresh"
        circle
        @click="handleRefresh"
        :loading="loading"
      />
    </div>

    <el-empty
      v-if="courseStore.myCourses.length === 0"
      description="您还没有选课"
      :image-size="160"
    />

    <el-table
      v-else
      :data="courseStore.myCourses"
      stripe
      class="courses-table"
      header-cell-class-name="table-header-cell"
    >
      <el-table-column prop="course_name" label="课程名称" min-width="150" />
      <el-table-column prop="description" label="描述" min-width="200">
        <template #default="{ row }">
          {{ row.description || '暂无描述' }}
        </template>
      </el-table-column>
      <el-table-column prop="teacher" label="教师" width="120" />
      <el-table-column prop="enrolled_at" label="选课时间" width="180" />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button
            type="danger"
            size="small"
            round
            @click="handleDrop(row.course_id)"
          >
            退课
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { useCourseStore } from '@/stores/course'

const courseStore = useCourseStore()
const loading = ref(false)
let refreshTimer = null

const handleRefresh = async () => {
  loading.value = true
  try {
    await courseStore.fetchMyCourses()
  } finally {
    loading.value = false
  }
}

const handleDrop = (courseId) => {
  courseStore.dropCourse(courseId)
}

// 自动刷新功能
const startAutoRefresh = () => {
  refreshTimer = setInterval(() => {
    courseStore.fetchMyCourses()
  }, 30000) // 每30秒刷新一次
}

onMounted(() => {
  courseStore.fetchMyCourses()
  startAutoRefresh()
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.my-course-container {
  width: 100%;
}

.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.table-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #1c1c1c;
}

.courses-table {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.table-header-cell) {
  background: linear-gradient(135deg, #00AFF0 0%, #667eea 100%);
  color: white;
  font-weight: 600;
}

:deep(.el-table__row:hover) {
  background: rgba(0, 175, 240, 0.05);
}

:deep(.el-table) {
  border-radius: 12px;
}
</style>
