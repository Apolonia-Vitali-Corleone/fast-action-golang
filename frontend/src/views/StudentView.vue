<template>
  <div class="student-container">
    <!-- 顶部Header -->
    <div class="page-header">
      <h1>学生选课系统</h1>
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

    <!-- Tabs切换 -->
    <el-card class="main-card" shadow="never">
      <el-tabs v-model="activeTab" class="view-tabs" stretch>
        <el-tab-pane label="课程表" name="schedule">
          <ScheduleTable />
        </el-tab-pane>
        <el-tab-pane label="可选课程" name="courses">
          <CourseList />
        </el-tab-pane>
        <el-tab-pane label="我的课程" name="my">
          <MyCourseTable />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'
import CourseList from '@/components/Student/CourseList.vue'
import MyCourseTable from '@/components/Student/MyCourseTable.vue'
import ScheduleTable from '@/components/Student/ScheduleTable.vue'

const userStore = useUserStore()
const activeTab = ref('schedule')

const handleLogout = async () => {
  await userStore.logout()
}
</script>

<style scoped>
.student-container {
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

.main-card {
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.main-card :deep(.el-card__body) {
  padding: 32px;
}

.view-tabs {
  margin-top: 0;
}

.view-tabs :deep(.el-tabs__header) {
  margin-bottom: 32px;
}

.view-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.view-tabs :deep(.el-tabs__active-bar) {
  height: 4px;
  border-radius: 4px;
  background: linear-gradient(90deg, #00AFF0 0%, #667eea 100%);
}

.view-tabs :deep(.el-tabs__item) {
  font-size: 16px;
  font-weight: 600;
  color: #909399;
  transition: all 0.3s;
}

.view-tabs :deep(.el-tabs__item:hover) {
  color: #00AFF0;
}

.view-tabs :deep(.el-tabs__item.is-active) {
  color: #00AFF0;
}

@media (max-width: 768px) {
  .student-container {
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

  .main-card :deep(.el-card__body) {
    padding: 20px;
  }
}
</style>
