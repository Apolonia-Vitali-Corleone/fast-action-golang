<template>
  <el-card class="table-card" shadow="hover">
    <template #header>
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
    </template>

    <el-empty
      v-if="courseStore.teacherCourses.length === 0"
      description="您还没有创建课程"
      :image-size="160"
    />

    <el-table
      v-else
      :data="courseStore.teacherCourses"
      stripe
      class="courses-table"
      header-cell-class-name="table-header-cell"
    >
      <el-table-column prop="name" label="课程名称" min-width="150" />
      <el-table-column prop="description" label="描述" min-width="200">
        <template #default="{ row }">
          {{ row.description || '暂无描述' }}
        </template>
      </el-table-column>
      <el-table-column prop="capacity" label="容量" width="100" align="center" />
      <el-table-column label="已选人数" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="row.enrolled >= row.capacity ? 'danger' : 'success'">
            {{ row.enrolled }}/{{ row.capacity }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            round
            @click="handleViewStudents(row.id)"
          >
            查看学生
          </el-button>
          <el-button
            type="danger"
            size="small"
            round
            @click="handleDelete(row.id)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { useCourseStore } from '@/stores/course'

const courseStore = useCourseStore()
const loading = ref(false)

const handleRefresh = async () => {
  loading.value = true
  try {
    await courseStore.fetchTeacherCourses()
  } finally {
    loading.value = false
  }
}

const handleViewStudents = async (courseId) => {
  await courseStore.viewStudents(courseId)
}

const handleDelete = async (courseId) => {
  await courseStore.deleteCourse(courseId)
}

onMounted(() => {
  courseStore.fetchTeacherCourses()
})
</script>

<style scoped>
.table-card {
  border-radius: 16px;
}

.table-card :deep(.el-card__header) {
  background: linear-gradient(135deg, rgba(0, 175, 240, 0.05) 0%, rgba(102, 126, 234, 0.05) 100%);
}

.table-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.table-header h2 {
  margin: 0;
  font-size: 20px;
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
</style>
