<template>
  <div id="app">
    <!-- 未登录：显示认证页面 -->
    <AuthView v-if="!userStore.currentUser" />

    <!-- 已登录：根据角色显示不同的界面 -->
    <StudentView v-else-if="userStore.currentUser.role === 'student'" />
    <TeacherView v-else-if="userStore.currentUser.role === 'teacher'" />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from './stores/user'
import { useCourseStore } from './stores/course'
import AuthView from './views/AuthView.vue'
import StudentView from './views/StudentView.vue'
import TeacherView from './views/TeacherView.vue'

const userStore = useUserStore()
const courseStore = useCourseStore()

onMounted(async () => {
  // 尝试从localStorage恢复登录状态
  const restored = await userStore.restoreSession()

  if (restored && userStore.currentUser) {
    // 根据角色加载相应的数据
    if (userStore.currentUser.role === 'student') {
      await courseStore.fetchAvailableCourses()
    } else if (userStore.currentUser.role === 'teacher') {
      await courseStore.fetchTeacherCourses()
    }
  }
})
</script>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  height: 100%;
}

/* Element Plus全局样式覆盖 - OnlyFans风格 */
:root {
  --el-color-primary: #00AFF0;
  --el-color-primary-light-3: #4dc3f5;
  --el-color-primary-light-5: #80d4f9;
  --el-color-primary-light-7: #b3e5fc;
  --el-color-primary-light-8: #d1f0fd;
  --el-color-primary-light-9: #e8f7fe;
  --el-color-primary-dark-2: #008ac0;

  --el-border-radius-base: 12px;
  --el-border-radius-small: 8px;
  --el-border-radius-round: 20px;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #00AFF0 0%, #667eea 100%);
  border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #008ac0 0%, #5569d3 100%);
}

/* 全局按钮样式增强 */
.el-button {
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 175, 240, 0.25);
}

.el-button:active {
  transform: translateY(0);
}

/* 全局卡片样式 */
.el-card {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 全局表格行悬停效果 */
.el-table__row {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 全局message样式 */
.el-message {
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

/* 全局dialog样式 */
.el-dialog {
  border-radius: 16px;
}

/* 全局空状态样式 */
.el-empty__description {
  color: #909399;
  font-size: 15px;
}
</style>
