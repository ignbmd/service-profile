package request

import "time"

type GetStudentPresenceResponse struct {
	Data map[string]any `json:"data"`
}

type ClassSchedule struct {
	ID                  string    `json:"_id"`
	ClassScheduleID     string    `json:"class_schedule_id"`
	ClassroomTitle      string    `json:"classroom_title"`
	Comment             string    `json:"comment"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	DurationPercentages float64   `json:"duration_percentages"`
	Presence            string    `json:"presence"`
	ScheduleTopic       string    `json:"schedule_topic"`
	SmartbtwID          int       `json:"smartbtw_id"`
	UpdatedAt           time.Time `json:"updated_at"`
}
