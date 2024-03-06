package request

import (
	"time"

	"github.com/bytedance/sonic"
)

type CreateHistoryCpns struct {
	SmartBtwID        int       `json:"smartbtw_id"`
	TaskID            int       `json:"task_id"`
	PackageID         int       `json:"package_id"`
	ModuleCode        string    `json:"module_code"`
	ModuleType        string    `json:"module_type"`
	PackageType       string    `json:"package_type"`
	Twk               float64   `json:"twk"`
	Tiu               float64   `json:"tiu"`
	Tkp               float64   `json:"tkp"`
	TwkPass           float64   `json:"twk_passing_grade"`
	TiuPass           float64   `json:"tiu_passing_grade"`
	TkpPass           float64   `json:"tkp_passing_grade"`
	TwkPassStatus     bool      `json:"twk_pass_status"`
	TiuPassStatus     bool      `json:"tiu_pass_status"`
	TkpPassStatus     bool      `json:"tkp_pass_status"`
	AllPassStatus     bool      `json:"is_all_passed"`
	Total             float64   `json:"total"`
	Repeat            int       `json:"repeat"`
	ExamName          string    `json:"exam_name"`
	Grade             string    `json:"grade"`
	TargetID          string    `json:"target_id"`
	Start             time.Time `json:"start"`
	End               time.Time `json:"end"`
	IsLive            bool      `json:"is_live"`
	StudentName       string    `json:"student_name"`
	SchoolOriginID    string    `json:"school_origin_id"`
	SchoolOrigin      string    `json:"school_origin"`
	InstanceName      string    `json:"instance_name"`
	InstanceID        int       `json:"instance_id"`
	PositionName      string    `json:"position_name"`
	PositionID        int       `json:"position_id"`
	CompetitionCpnsID uint      `json:"competition_id"`
	FormationType     string    `json:"formation_type"`
	FormationCode     string    `json:"formation_code"`
	FormationLocation string    `json:"formation_location"`
	TargetScore       float64   `json:"target_score" bson:"target_score"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GetHistoryCpnsResultElastic struct {
	SmartBtwID      int        `json:"smartbtw_id"`
	TaskID          int        `json:"task_id"`
	Name            string     `json:"name"`
	PackageID       int        `json:"package_id"`
	ModuleCode      string     `json:"module_code"`
	ModuleType      string     `json:"module_type"`
	PackageType     string     `json:"package_type"`
	Twk             float64    `json:"twk_score"`
	Tiu             float64    `json:"tiu_score"`
	Tkp             float64    `json:"tkp_score"`
	TwkPass         float64    `json:"twk_passing_grade"`
	TiuPass         float64    `json:"tiu_passing_grade"`
	TkpPass         float64    `json:"tkp_passing_grade"`
	TwkPassStatus   bool       `json:"twk_pass_status"`
	TiuPassStatus   bool       `json:"tiu_pass_status"`
	TkpPassStatus   bool       `json:"tkp_pass_status"`
	AllPassStatus   bool       `json:"is_all_passed"`
	Total           float64    `json:"total_score"`
	Repeat          int        `json:"repeat"`
	Title           string     `json:"title"`
	Grade           string     `json:"grade"`
	TargetID        string     `json:"target_id"`
	TargetScore     float64    `json:"target_score"`
	Start           *time.Time `json:"start"`
	End             *time.Time `json:"end"`
	IsLive          bool       `json:"is_live"`
	TwkTimeConsumed int        `json:"twk_time_consumed"`
	TiuTimeConsumed int        `json:"tiu_time_consumed"`
	TkpTimeConsumed int        `json:"tkp_time_consumed"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type GetHistoryUKAResultElastic struct {
	SmartBtwID      int        `json:"smartbtw_id"`
	TaskID          int        `json:"task_id"`
	ExamName        string     `json:"exam_name"`
	PackageID       int        `json:"package_id"`
	MajorName       string     `json:"major_name"`
	ModuleCode      string     `json:"module_code"`
	ModuleType      string     `json:"module_type"`
	PackageType     string     `json:"package_type"`
	Twk             float64    `json:"twk"`
	Tiu             float64    `json:"tiu"`
	Tkp             float64    `json:"tkp"`
	TwkPass         float64    `json:"twk_passing_grade"`
	TiuPass         float64    `json:"tiu_passing_grade"`
	TkpPass         float64    `json:"tkp_passing_grade"`
	TwkPassStatus   bool       `json:"twk_pass_status"`
	TiuPassStatus   bool       `json:"tiu_pass_status"`
	TkpPassStatus   bool       `json:"tkp_pass_status"`
	AllPassStatus   bool       `json:"is_all_passed"`
	Total           float64    `json:"total"`
	Repeat          int        `json:"repeat"`
	Title           string     `json:"title"`
	Grade           string     `json:"grade"`
	TargetID        string     `json:"target_id"`
	TargetScore     float64    `json:"target_score"`
	Start           *time.Time `json:"start"`
	End             *time.Time `json:"end"`
	IsLive          bool       `json:"is_live"`
	TwkTimeConsumed int        `json:"twk_time_consumed"`
	TiuTimeConsumed int        `json:"tiu_time_consumed"`
	TkpTimeConsumed int        `json:"tkp_time_consumed"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type HistoryCpnsElastic struct {
	SmartBtwID      int       `json:"smartbtw_id"`
	TaskID          int       `json:"task_id"`
	PackageID       int       `json:"package_id"`
	ModuleCode      string    `json:"module_code"`
	ModuleType      string    `json:"module_type"`
	PackageType     string    `json:"package_type"`
	Twk             float64   `json:"twk"`
	Tiu             float64   `json:"tiu"`
	Tkp             float64   `json:"tkp"`
	TwkPass         float64   `json:"twk_passing_grade"`
	TiuPass         float64   `json:"tiu_passing_grade"`
	TkpPass         float64   `json:"tkp_passing_grade"`
	TwkPassStatus   bool      `json:"twk_pass_status"`
	TiuPassStatus   bool      `json:"tiu_pass_status"`
	TkpPassStatus   bool      `json:"tkp_pass_status"`
	AllPassStatus   bool      `json:"is_all_passed"`
	Total           float64   `json:"total"`
	Repeat          int       `json:"repeat"`
	ExamName        string    `json:"exam_name"`
	Grade           string    `json:"grade"`
	TargetID        string    `json:"target_id"`
	Start           time.Time `json:"start"`
	End             time.Time `json:"end"`
	IsLive          bool      `json:"is_live"`
	TwkTimeConsumed int       `json:"twk_time_consumed"`
	TiuTimeConsumed int       `json:"tiu_time_consumed"`
	TkpTimeConsumed int       `json:"tkp_time_consumed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UpdateHistoryCpnsTime struct {
	SmartBtwID      int `json:"smartbtw_id"`
	TaskID          int `json:"task_id"`
	TwkTimeConsumed int `json:"twk_time_consumed"`
	TiuTimeConsumed int `json:"tiu_time_consumed"`
	TkpTimeConsumed int `json:"tkp_time_consumed"`
}

type UpdateHistoryCPNS struct {
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

type GetRankHistoryCpns struct {
	SmartBtwID    int       `json:"smartbtw_id"`
	TaskID        int       `json:"task_id"`
	InstanceName  string    `json:"instance_name"`
	JurusanName   string    `json:"jurusan_name"`
	PackageID     int       `json:"package_id"`
	ModuleCode    string    `json:"module_code"`
	ModuleType    string    `json:"module_type"`
	PackageType   string    `json:"package_type"`
	Twk           float64   `json:"twk"`
	Tiu           float64   `json:"tiu"`
	Tkp           float64   `json:"tkp"`
	TwkPass       float64   `json:"twk_passing_grade"`
	TiuPass       float64   `json:"tiu_passing_grade"`
	TkpPass       float64   `json:"tkp_passing_grade"`
	TwkPassStatus bool      `json:"twk_pass_status"`
	TiuPassStatus bool      `json:"tiu_pass_status"`
	TkpPassStatus bool      `json:"tkp_pass_status"`
	AllPassStatus bool      `json:"is_all_passed"`
	Total         float64   `json:"total"`
	Repeat        int       `json:"repeat"`
	ExamName      string    `json:"exam_name"`
	Grade         string    `json:"grade"`
	TargetID      string    `json:"target_id"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	IsLive        bool      `json:"is_live"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RequestRankCpns struct {
	ExamName     string  `json:"exam_name"`
	TaskID       int     `json:"task_id"`
	InstanceName string  `json:"instance_name"`
	PositionName string  `json:"position_name"`
	Start        string  `json:"start"`
	End          string  `json:"end"`
	Duration     string  `json:"duration"`
	Twk          float64 `json:"twk"`
	Tiu          float64 `json:"tiu"`
	Tkp          float64 `json:"tkp"`
	TwkStatus    bool    `json:"twk_status"`
	TiuStatus    bool    `json:"tiu_status"`
	TkpStatus    bool    `json:"tkp_status"`
	Date         string  `json:"date"`
	Rank         int     `json:"rank"`
	Total        float64 `json:"total"`
	Status       bool    `json:"status"`
}

type HistoryCPNSQueryParams struct {
	Limit *int `json:"limit"`
}

type MessageHistoryCpnsBody struct {
	Version int               `json:"version"`
	Data    CreateHistoryCpns `json:"data" valid:"required"`
}
type MessageTimeConsumedHistoryCpnsBody struct {
	Version int                   `json:"version"`
	Data    UpdateHistoryCpnsTime `json:"data" valid:"required"`
}

func UnmarshalMessageHistoryCpnsBody(data []byte) (MessageHistoryCpnsBody, error) {
	var decoded MessageHistoryCpnsBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalMessageHistoryCpnsTimeConsumedBody(data []byte) (MessageTimeConsumedHistoryCpnsBody, error) {
	var decoded MessageTimeConsumedHistoryCpnsBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
