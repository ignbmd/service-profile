package mockstruct

import (
	"time"

	"smartbtw.com/services/profile/models"
)

type StudentElasticCacheRequest struct {
	SmartBTWID         int             `json:"smartbtw_id"`
	AccountType        string          `json:"account_type"`
	StudentProfileData *models.Student `json:"student_profile_data"`
}

type SyncBinsusProfile struct {
	SmartbtwID         int       `json:"smartbtw_id" bson:"smartbtw_id"`
	Name               string    `json:"name" bson:"name"`
	Email              string    `json:"email" bson:"email"`
	Address            string    `json:"address" bson:"address"`
	Gender             string    `json:"gender" bson:"gender"`
	BirthDate          time.Time `json:"birth_date" bson:"birth_date"`
	Province           string    `json:"province" bson:"province"`
	ProvinceID         int       `json:"province_id" bson:"province_id"`
	Region             string    `json:"region" bson:"region"`
	RegionID           int       `json:"region_id" bson:"region_id"`
	DomicileProvince   string    `json:"domicile_province" bson:"domicile_province"`
	DomicileProvinceID int       `json:"domicile_province_id" bson:"domicile_province_id"`
	DomicileRegion     string    `json:"domicile_region" bson:"domicile_region"`
	DomicileRegionID   int       `json:"domicile_region_id" bson:"domicile_region_id"`
	LastEdID           string    `json:"last_ed_id" bson:"last_ed_id"`
	LastEdName         string    `json:"last_ed_name" bson:"last_ed_name"`
	LastEdType         string    `json:"last_ed_type" bson:"last_ed_type"`
	LastEdMajor        string    `json:"last_ed_major" bson:"last_ed_major"`
	LastEdMajorID      uint      `json:"last_ed_major_id" bson:"last_ed_major_id"`
	LastEdRegion       string    `json:"last_ed_region" bson:"last_ed_region"`
	LastEdRegionID     uint      `json:"last_ed_region_id" bson:"last_ed_region_id"`
	BranchCode         string    `json:"branch_code" bson:"branch_code"`
	BranchName         string    `json:"branch_name" bson:"branch_name"`
	ParentName         string    `json:"parent_name" bson:"parent_name"`
	ParentNumber       string    `json:"parent_number" bson:"parent_number"`
	Interest           string    `json:"interest,omitempty" bson:"interest"`
}
type BinsusSchool struct {
	SmartbtwID    int     `json:"smartbtw_id" bson:"smartbtw_id"`
	SchoolID      int     `json:"school_id" bson:"school_id"`
	MajorID       int     `json:"major_id" bson:"major_id"`
	SchoolName    string  `json:"school_name" bson:"school_name"`
	MajorName     string  `json:"major_name" bson:"major_name"`
	TargetScore   float64 `json:"target_score" bson:"target_score"`
	CompetitionID int     `json:"competition_id" bson:"competition_id"`
	PolbitType    string  `json:"polbit_type" bson:"polbit_type"`
	PolbitName    string  `json:"polbit_name" bson:"polbit_name"`
}
type SyncBinsusFinalProfile struct {
	SmartbtwID int               `json:"smartbtw_id" bson:"smartbtw_id"`
	Profile    SyncBinsusProfile `json:"profile" bson:"profile"`
	School     BinsusSchool      `json:"school" bson:"school"`
}

type BinsusScreeningSummary struct {
	SmartBTWID int `json:"smartbtw_id" bson:"smartbtw_id"`
	Year       int `json:"year" bson:"year"`
	Record     struct {
		Score struct {
			ScoresData []BinsusScoreData `json:"scores_data"`
		} `json:"score"`
	} `json:"record"`
	ChoosenSummary struct {
		FinalPass   bool   `json:"final_pass" bson:"final_pass"`
		FinalReason string `json:"final_reason" bson:"final_reason"`
		FinalState  string `json:"final_state" bson:"final_state"`
	} `json:"choosen_school" bson:"choosen_school"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type BinsusScoreData struct {
	TaskID    int `json:"task_id"`
	PackageID int `json:"package_id"`
}

type CreateHistorySimpleData struct {
	SmartBtwID     int    `json:"smartbtw_id"`
	TaskID         int    `json:"task_id"`
	PackageID      int    `json:"package_id"`
	SchoolOriginID string `json:"school_origin_id"`
	SchoolOrigin   string `json:"school_origin"`
	SchoolID       int    `json:"school_id"`
	MajorID        int    `json:"major_id"`
	SchoolName     string `json:"school_name"`
	MajorName      string `json:"major_name"`
	StudentName    string `json:"student_name"`
	StudentEmail   string `json:"student_email"`
}

type SchoolStudentCount struct {
	StudentTotal       int `json:"student_total"`
	StudentJoinedClass int `json:"student_joined_class"`
}
