<template>
  <div class="home-view">
    <!-- 根据角色显示不同的界面 -->
    <StudentView v-if="userStore.currentUser?.role === 'student'" />
    <TeacherView v-else-if="userStore.currentUser?.role === 'teacher'" />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useCourseStore } from '@/stores/course'
import StudentView from './StudentView.vue'
import TeacherView from './TeacherView.vue'

const userStore = useUserStore()
const courseStore = useCourseStore()

onMounted(async () => {
  // 加载用户数据
  if (userStore.currentUser) {
    if (userStore.currentUser.role === 'student') {
      await courseStore.fetchAvailableCourses()
    } else if (userStore.currentUser.role === 'teacher') {
      await courseStore.fetchTeacherCourses()
    }
  }
})
</script>

<style scoped>
.home-view {
  height: 100%;
  width: 100%;
  overflow: hidden;
}
</style>
