<template>
  <el-dialog
    v-model="courseStore.showStudentsDialog"
    title="选课学生名单"
    width="800px"
    :close-on-click-modal="false"
    class="student-dialog"
  >
    <div class="dialog-content">
      <div class="course-info">
        <p><strong>课程名称：</strong>{{ courseStore.currentCourse.name }}</p>
        <p><strong>总人数：</strong>{{ courseStore.courseStudents.total }}</p>
      </div>

      <el-divider />

      <el-empty
        v-if="courseStore.courseStudents.students.length === 0"
        description="暂无学生选课"
        :image-size="120"
      />

      <el-table
        v-else
        :data="courseStore.courseStudents.students"
        stripe
        height="400"
        class="students-table"
        header-cell-class-name="table-header-cell"
      >
        <el-table-column prop="id" label="学生ID" width="100" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="enrolled_at" label="选课时间" width="180" />
      </el-table>
    </div>

    <template #footer>
      <el-button
        type="primary"
        round
        @click="courseStore.closeStudentsDialog"
      >
        关闭
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { useCourseStore } from '@/stores/course'

const courseStore = useCourseStore()
</script>

<style scoped>
.dialog-content {
  padding: 0;
}

.course-info p {
  margin: 8px 0;
  font-size: 15px;
  color: #1c1c1c;
}

.course-info strong {
  color: #00AFF0;
  font-weight: 600;
}

.students-table {
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

:deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, rgba(0, 175, 240, 0.05) 0%, rgba(102, 126, 234, 0.05) 100%);
  padding: 20px;
}

:deep(.el-dialog__title) {
  font-size: 20px;
  font-weight: 700;
  color: #1c1c1c;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 16px 24px;
  background: #f7f7f7;
}
</style>
