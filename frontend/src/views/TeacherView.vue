<template>
  <div class="teacher-container">
    <!-- 顶部Header -->
    <div class="page-header">
      <h1>教师管理系统</h1>
      <div class="user-actions">
        <span class="username">欢迎，{{ userStore.currentUser.username }}</span>
        <el-button
          type="info"
          round
          @click="handleLogout"
        >
          退出
        </el-button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="content-wrapper">
      <!-- 创建课程表单 -->
      <CreateCourseForm />

      <!-- 课程管理表格 -->
      <CourseManageTable />
    </div>

    <!-- 学生列表弹窗 -->
    <StudentListDialog />
  </div>
</template>

<script setup>
import { useUserStore } from '@/stores/user'
import CreateCourseForm from '@/components/Teacher/CreateCourseForm.vue'
import CourseManageTable from '@/components/Teacher/CourseManageTable.vue'
import StudentListDialog from '@/components/Teacher/StudentListDialog.vue'

const userStore = useUserStore()

const handleLogout = async () => {
  await userStore.logout()
}
</script>

<style scoped>
.teacher-container {
  min-height: 100vh;
  padding: 24px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.page-header {
  background: white;
  padding: 24px 32px;
  border-radius: 16px;
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.page-header h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #00AFF0 0%, #667eea 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.username {
  font-weight: 600;
  color: #1c1c1c;
  font-size: 16px;
}

.content-wrapper {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

@media (max-width: 768px) {
  .teacher-container {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .page-header h1 {
    font-size: 24px;
  }
}
</style>
