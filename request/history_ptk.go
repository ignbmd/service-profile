package request

import (
	"time"

	"github.com/bytedance/sonic"
)

type CreateHistoryPtk struct {
	SmartBtwID          int       `json:"smartbtw_id"`
	TaskID              int       `json:"task_id"`
	PackageID           int       `json:"package_id"`
	ModuleCode          string    `json:"module_code"`
	ModuleType          string    `json:"module_type"`
	PackageType         string    `json:"package_type"`
	Twk                 float64   `json:"twk"`
	Tiu                 float64   `json:"tiu"`
	Tkp                 float64   `json:"tkp"`
	TwkPass             float64   `json:"twk_passing_grade"`
	TiuPass             float64   `json:"tiu_passing_grade"`
	TkpPass             float64   `json:"tkp_passing_grade"`
	TwkPassStatus       bool      `json:"twk_pass_status"`
	TiuPassStatus       bool      `json:"tiu_pass_status"`
	TkpPassStatus       bool      `json:"tkp_pass_status"`
	AllPassStatus       bool      `json:"is_all_passed"`
	Total               float64   `json:"total"`
	Repeat              int       `json:"repeat"`
	ExamName            string    `json:"exam_name"`
	Grade               string    `json:"grade"`
	TargetID            string    `json:"target_id"`
	Start               time.Time `json:"start"`
	End                 time.Time `json:"end"`
	IsLive              bool      `json:"is_live"`
	StudentName         string    `json:"student_name"`
	SchoolOriginID      string    `json:"school_origin_id"`
	SchoolOrigin        string    `json:"school_origin"`
	SchoolID            int       `json:"school_id"`
	MajorID             int       `json:"major_id"`
	SchoolName          string    `json:"school_name"`
	MajorName           string    `json:"major_name"`
	PolbitType          string    `json:"polbit_type"`
	PolbitCompetitionID int       `json:"polbit_competition_id"`
	PolbitLocationID    int       `json:"polbit_location_id"`
	TargetScore         float64   `json:"target_score" bson:"target_score"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
type UpdateScoreOnlyPTK struct {
	SmartBtwID          int     `json:"smartbtw_id"`
	AverageTWK          float64 `json:"twk_avg_score"`
	AverageTIU          float64 `json:"tiu_avg_score"`
	AverageTKP          float64 `json:"tkp_avg_score"`
	AveragePercentTWK   float64 `json:"twk_avg_percent_score"`
	AveragePercentTIU   float64 `json:"tiu_avg_percent_score"`
	AveragePercentTKP   float64 `json:"tkp_avg_percent_score"`
	AveragePercentTotal float64 `json:"tt_avg_percent_score"`
	AverageTotal        float64 `json:"total_avg_score"`
	LatestTotal         float64 `json:"latest_total_score"`
	LatestTotalPercent  float64 `json:"latest_total_score_percent"`
	ModuleDone          int     `json:"module_done"`
}

type UpdateHistoryPtk struct {
	ID          int     `json:"id" bson:"id"`
	SmartBtwID  int     `json:"smartbtw_id"`
	TaskID      int     `json:"task_id"`
	PackageID   int     `json:"package_id"`
	PackageType string  `json:"package_type"`
	ModuleCode  string  `json:"module_code"`
	ModuleType  string  `json:"module_type"`
	Twk         float64 `json:"twk"`
	Tiu         float64 `json:"tiu"`
	Tkp         float64 `json:"tkp"`
	Total       float64 `json:"total"`
	Repeat      int     `json:"repeat"`
	ExamName    string  `json:"exam_name"`
	Grade       string  `json:"grade"`
	TargetID    string  `json:"target_id"`
}

type BodyUpdateStudentDuration struct {
	SmartbtwID int       `json:"smartbtw_id"`
	TaskID     int       `json:"task_id"`
	Repeat     int       `json:"repeat"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Program    string    `json:"program"`
}

type HistoryPTKQueryParams struct {
	Limit *int `json:"limit"`
}
type HistoryPtkGetAll struct {
	SmartbtwID int    `json:"smartbtw_id"`
	Limit      *int64 `json:"limit"`
	Page       *int64 `json:"page"`
}

type HistoryPTKSendToExamResult struct {
	SmartbtwID int    `json:"smartbtw_id"`
	TaskID     int    `json:"task_id"`
	Repeat     int    `json:"repeat"`
	Program    string `json:"program"`
}

type MessageHistoryPTKSendToExamResult struct {
	Version int                        `json:"version"`
	Data    HistoryPTKSendToExamResult `json:"data" valid:"required"`
}

type MessageHistoryPtkBody struct {
	Version int              `json:"version"`
	Data    CreateHistoryPtk `json:"data" valid:"required"`
}
type MessageUpdateDurationBody struct {
	Version int                       `json:"version"`
	Data    BodyUpdateStudentDuration `json:"data" valid:"required"`
}

func UnmarshalMessageHistoryPtkBody(data []byte) (MessageHistoryPtkBody, error) {
	var decoded MessageHistoryPtkBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalUpdateDurationBody(data []byte) (MessageUpdateDurationBody, error) {
	var decoded MessageUpdateDurationBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
