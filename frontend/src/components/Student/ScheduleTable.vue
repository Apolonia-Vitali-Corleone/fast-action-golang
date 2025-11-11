<template>
  <div class="schedule-container">
    <div class="schedule-header">
      <h2 class="title">我的课表</h2>
      <el-button type="primary" @click="refreshSchedule">
        <el-icon><Refresh /></el-icon>
        刷新
      </el-button>
    </div>

    <div class="schedule-table-wrapper">
      <table class="schedule-table">
        <thead>
          <tr>
            <th class="corner-cell">节次/星期</th>
            <th v-for="day in weekDays" :key="day.value" class="day-header">
              {{ day.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(slot, slotIndex) in timeSlots" :key="slot.value">
            <td class="time-slot-header">
              <div class="slot-name">{{ slot.label }}</div>
              <div class="slot-time">{{ slot.time }}</div>
            </td>
            <td
              v-for="(day, dayIndex) in weekDays"
              :key="day.value"
              class="course-cell"
              :class="{ 'has-course': schedule[slotIndex]?.[dayIndex] }"
            >
              <div v-if="schedule[slotIndex]?.[dayIndex]" class="course-info">
                <div class="course-name">{{ schedule[slotIndex][dayIndex].course_name }}</div>
                <div class="course-detail">
                  <span class="teacher">
                    <el-icon><User /></el-icon>
                    {{ schedule[slotIndex][dayIndex].teacher_name }}
                  </span>
                  <span class="classroom">
                    <el-icon><Location /></el-icon>
                    {{ schedule[slotIndex][dayIndex].classroom || '未安排' }}
                  </span>
                </div>
                <div class="course-weeks">
                  第{{ schedule[slotIndex][dayIndex].start_week }}-{{ schedule[slotIndex][dayIndex].end_week }}周
                </div>
              </div>
              <div v-else class="empty-cell">
                <span class="empty-text">—</span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Refresh, User, Location } from '@element-plus/icons-vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const API_BASE = 'http://localhost:8000/api'

const schedule = ref([])

const weekDays = [
  { value: 1, label: '周一' },
  { value: 2, label: '周二' },
  { value: 3, label: '周三' },
  { value: 4, label: '周四' },
  { value: 5, label: '周五' }
]

const timeSlots = [
  { value: 1, label: '上午第一节', time: '08:00-09:40' },
  { value: 2, label: '上午第二节', time: '10:00-11:40' },
  { value: 3, label: '下午第一节', time: '14:00-15:40' },
  { value: 4, label: '下午第二节', time: '16:00-17:40' }
]

const fetchSchedule = async () => {
  try {
    const res = await axios.get(`${API_BASE}/student/schedule/`)
    schedule.value = res.data.schedule || []
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '获取课表失败')
  }
}

const refreshSchedule = () => {
  fetchSchedule()
  ElMessage.success('课表已刷新')
}

onMounted(() => {
  fetchSchedule()
})
</script>

<style scoped>
.schedule-container {
  padding: 24px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.schedule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  color: #1c1c1c;
  margin: 0;
}

.schedule-table-wrapper {
  overflow-x: auto;
}

.schedule-table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #e4e7ed;
  padding: 12px;
  text-align: center;
}

.corner-cell {
  background: linear-gradient(135deg, #00aff0 0%, #0085c7 100%);
  color: #fff;
  font-weight: 600;
  font-size: 14px;
}

.day-header {
  background: linear-gradient(135deg, #00aff0 0%, #0085c7 100%);
  color: #fff;
  font-weight: 600;
  font-size: 16px;
  padding: 16px 12px;
  min-width: 140px;
}

.time-slot-header {
  background: #f5f7fa;
  font-weight: 600;
  color: #303133;
  min-width: 120px;
  vertical-align: middle;
}

.slot-name {
  font-size: 14px;
  margin-bottom: 4px;
}

.slot-time {
  font-size: 12px;
  color: #909399;
  font-weight: normal;
}

.course-cell {
  min-width: 140px;
  min-height: 100px;
  vertical-align: middle;
  background: #fff;
  transition: all 0.3s;
}

.course-cell:hover {
  background: #f5f7fa;
}

.course-cell.has-course {
  background: linear-gradient(135deg, #e6f7ff 0%, #f0f9ff 100%);
}

.course-cell.has-course:hover {
  background: linear-gradient(135deg, #d4f0ff 0%, #e1f4ff 100%);
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 175, 240, 0.15);
}

.course-info {
  padding: 8px;
}

.course-name {
  font-size: 15px;
  font-weight: 600;
  color: #00aff0;
  margin-bottom: 8px;
  line-height: 1.4;
}

.course-detail {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 6px;
}

.course-detail span {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 13px;
  color: #606266;
}

.course-detail .el-icon {
  font-size: 14px;
}

.course-weeks {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.empty-cell {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-text {
  color: #dcdfe6;
  font-size: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .schedule-table-wrapper {
    overflow-x: scroll;
  }

  .day-header {
    min-width: 100px;
    font-size: 14px;
  }

  .course-cell {
    min-width: 100px;
  }
}
</style>
