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
//   - 在同一天（DayOfWeek相同）
//   - 在同一节次（TimeSlot相同）
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
			// 检查是否在同一天且同一节次
			if newSchedule.DayOfWeek == existingSchedule.DayOfWeek &&
				newSchedule.TimeSlot == existingSchedule.TimeSlot {
				// 发现冲突，查询课程信息以返回详细提示
				var conflictCourse models.Course
				config.DB.First(&conflictCourse, existingSchedule.CourseID)

				timeSlotName := GetTimeSlotName(existingSchedule.TimeSlot)
				conflictMsg := fmt.Sprintf(
					"时间冲突：与已选课程《%s》冲突（周%s %s）",
					conflictCourse.Name,
					GetDayOfWeekName(existingSchedule.DayOfWeek),
					timeSlotName,
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
//   - day: 星期几（1-5）
// 返回:
//   - string: 中文星期名称
func GetDayOfWeekName(day int) string {
	days := []string{"", "一", "二", "三", "四", "五"}
	if day < 1 || day > 5 {
		return "未知"
	}
	return days[day]
}

// GetTimeSlotName 将节次数字转换为中文名称
// 参数:
//   - slot: 节次（1-4）
// 返回:
//   - string: 节次名称
func GetTimeSlotName(slot int) string {
	slots := []string{"", "上午第一节", "上午第二节", "下午第一节", "下午第二节"}
	if slot < 1 || slot > 4 {
		return "未知"
	}
	return slots[slot]
}

// ScheduleCell 课表单元格
type ScheduleCell struct {
	CourseID    int    `json:"course_id"`
	CourseName  string `json:"course_name"`
	TeacherName string `json:"teacher_name"`
	Classroom   string `json:"classroom"`
	StartWeek   int    `json:"start_week"`
	EndWeek     int    `json:"end_week"`
}

// GetStudentScheduleTable 获取学生的课表（二维数组）
// 返回一个 4x5 的二维数组，表示 4个节次 x 5天（周一到周五）
// 参数:
//   - studentID: 学生ID
//   - currentWeek: 当前周次（可选，用于过滤不在当前周次的课程）
// 返回:
//   - [][]*ScheduleCell: 课表二维数组 [timeSlot-1][dayOfWeek-1]
//   - error: 查询错误
func GetStudentScheduleTable(studentID int, currentWeek int) ([][]*ScheduleCell, error) {
	// 初始化 4x5 的二维数组（4个节次 x 5天）
	schedule := make([][]*ScheduleCell, 4)
	for i := range schedule {
		schedule[i] = make([]*ScheduleCell, 5)
	}

	// 查询学生已选课程
	var enrollments []models.Enrollment
	if err := config.DB.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("查询已选课程失败: %v", err)
	}

	if len(enrollments) == 0 {
		return schedule, nil // 返回空课表
	}

	// 提取课程ID
	var courseIDs []int
	for _, enrollment := range enrollments {
		courseIDs = append(courseIDs, enrollment.CourseID)
	}

	// 查询课程信息（预加载教师信息）
	var courses []models.Course
	if err := config.DB.Where("id IN ?", courseIDs).Find(&courses).Error; err != nil {
		return nil, fmt.Errorf("查询课程信息失败: %v", err)
	}

	// 建立课程ID到课程信息的映射
	courseMap := make(map[int]models.Course)
	for _, course := range courses {
		courseMap[course.ID] = course
	}

	// 查询教师信息
	var teacherIDs []int
	for _, course := range courses {
		teacherIDs = append(teacherIDs, course.TeacherID)
	}
	var teachers []models.Teacher
	if err := config.DB.Where("id IN ?", teacherIDs).Find(&teachers).Error; err != nil {
		return nil, fmt.Errorf("查询教师信息失败: %v", err)
	}

	// 建立教师ID到教师名称的映射
	teacherMap := make(map[int]string)
	for _, teacher := range teachers {
		teacherMap[teacher.ID] = teacher.Username
	}

	// 查询课程时间表
	var schedules []models.CourseSchedule
	if err := config.DB.Where("course_id IN ?", courseIDs).Find(&schedules).Error; err != nil {
		return nil, fmt.Errorf("查询课程时间表失败: %v", err)
	}

	// 填充课表
	for _, sch := range schedules {
		// 如果指定了当前周次，过滤不在当前周次的课程
		if currentWeek > 0 && (currentWeek < sch.StartWeek || currentWeek > sch.EndWeek) {
			continue
		}

		course, exists := courseMap[sch.CourseID]
		if !exists {
			continue
		}

		teacherName := teacherMap[course.TeacherID]

		cell := &ScheduleCell{
			CourseID:    course.ID,
			CourseName:  course.Name,
			TeacherName: teacherName,
			Classroom:   sch.Classroom,
			StartWeek:   sch.StartWeek,
			EndWeek:     sch.EndWeek,
		}

		// 放入对应的位置 [节次-1][星期-1]
		if sch.TimeSlot >= 1 && sch.TimeSlot <= 4 && sch.DayOfWeek >= 1 && sch.DayOfWeek <= 5 {
			schedule[sch.TimeSlot-1][sch.DayOfWeek-1] = cell
		}
	}

	return schedule, nil
}
