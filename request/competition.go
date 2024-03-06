package request

type StagesCompetitionList struct {
	Level                        int    `json:"level"`
	PackageID                    int    `json:"package_id"`
	TaskID                       int    `json:"task_id"`
	Program                      string `json:"program"`
	Type                         string `json:"type"`
	AttemptedStudent             int    `json:"attempted_student"`
	AttemptedSchool              int    `json:"attempted_school"`
	AttemptedStudentOriginSchool int    `json:"attempted_student_origin_school"`
}
