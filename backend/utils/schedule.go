package utils

import (
	"course-system/config"
	"course-system/models"
	"fmt"
)

// CheckScheduleConflict 检测选课时间冲突
// 参数:
//   - studentID: 学生ID
//   - newCourseID: 要选的新课程ID
// 返回:
//   - bool: true表示有冲突，false表示无冲突
//   - string: 冲突的详细信息（如果有冲突）
//   - error: 数据库查询错误
//
// 工作原理:
//  1. 查询新课程的所有上课时间
//  2. 查询学生已选课程的所有上课时间
//  3. 检测是否有时间重叠
//
// 时间冲突判断：
//   两个时间段 [A_start, A_end] 和 [B_start, B_end] 冲突的条件是：
//   - 在同一天（DayOfWeek相同）
//   - 时间有重叠：A_start < B_end AND B_start < A_end
func CheckScheduleConflict(studentID int, newCourseID int) (bool, string, error) {
	// ========== 步骤1: 查询新课程的上课时间 ==========
	var newCourseSchedules []models.CourseSchedule
	if err := config.DB.Where("course_id = ?", newCourseID).Find(&newCourseSchedules).Error; err != nil {
		return false, "", fmt.Errorf("查询课程时间失败: %v", err)
	}

	// 如果新课程没有设置时间表，则不检测冲突（允许选课）
	if len(newCourseSchedules) == 0 {
		return false, "", nil
	}

	// ========== 步骤2: 查询学生已选课程的ID列表 ==========
	var enrollments []models.Enrollment
	if err := config.DB.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		return false, "", fmt.Errorf("查询已选课程失败: %v", err)
	}

	// 提取已选课程ID
	var enrolledCourseIDs []int
	for _, enrollment := range enrollments {
		enrolledCourseIDs = append(enrolledCourseIDs, enrollment.CourseID)
	}

	// 如果学生还没选任何课程，则无冲突
	if len(enrolledCourseIDs) == 0 {
		return false, "", nil
	}

	// ========== 步骤3: 查询已选课程的所有上课时间 ==========
	var enrolledSchedules []models.CourseSchedule
	if err := config.DB.Where("course_id IN ?", enrolledCourseIDs).Find(&enrolledSchedules).Error; err != nil {
		return false, "", fmt.Errorf("查询已选课程时间失败: %v", err)
	}

	// ========== 步骤4: 检测时间冲突 ==========
	for _, newSchedule := range newCourseSchedules {
		for _, existingSchedule := range enrolledSchedules {
			// 检查是否在同一天
			if newSchedule.DayOfWeek != existingSchedule.DayOfWeek {
				continue // 不在同一天，不会冲突
			}

			// 检查时间是否重叠
			// 时间重叠的条件：new_start < existing_end AND existing_start < new_end
			if newSchedule.StartTime < existingSchedule.EndTime &&
				existingSchedule.StartTime < newSchedule.EndTime {
				// 发现冲突，查询课程信息以返回详细提示
				var conflictCourse models.Course
				config.DB.First(&conflictCourse, existingSchedule.CourseID)

				conflictMsg := fmt.Sprintf(
					"时间冲突：与已选课程《%s》冲突（周%d %s-%s）",
					conflictCourse.Name,
					existingSchedule.DayOfWeek,
					existingSchedule.StartTime,
					existingSchedule.EndTime,
				)
				return true, conflictMsg, nil
			}
		}
	}

	// 没有发现冲突
	return false, "", nil
}

// GetDayOfWeekName 将星期数字转换为中文名称
// 参数:
//   - day: 星期几（1-7）
// 返回:
//   - string: 中文星期名称
func GetDayOfWeekName(day int) string {
	days := []string{"", "一", "二", "三", "四", "五", "六", "日"}
	if day < 1 || day > 7 {
		return "未知"
	}
	return days[day]
}
