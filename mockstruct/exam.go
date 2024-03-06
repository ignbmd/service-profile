package mockstruct

import (
	"time"

	"github.com/lib/pq"
)

type GetUKACodeByschoolID struct {
	PackageID  uint   `json:"packages_id"`
	TryoutCode string `json:"tryout_code"`
	Program    string `json:"program"`
}

type Packages struct {
	ID              uint           `json:"id"`
	Title           string         `json:"title"`
	ModulesID       uint           `json:"modules_id"`
	LegacyTaskID    uint           `json:"legacy_task_id"`
	MaxRepeat       *uint          `json:"max_repeat"`
	Duration        uint           `json:"duration"`
	PackageType     string         `json:"package_type"`
	ProductCode     *string        `json:"product_code"`
	Status          *bool          `json:"status"`
	PrivacyType     string         `json:"privacy_type"`
	BranchCode      string         `json:"branch_code"`
	InstructionsID  uint           `json:"instructions_id"`
	StartDate       *time.Time     `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	Program         string         `json:"program"`
	Description     *string        `json:"description"`
	Requirement     *string        `json:"requirement"`
	AllowDiscussion *bool          `json:"allow_discussion"`
	Tags            pq.StringArray `json:"tags"`
}

type GetUKACodeBySchoolIDResponse struct {
	Data []GetUKACodeByschoolID `json:"data"`
}

type GetPackageByIDResponse struct {
	Data Packages `json:"data"`
}

type SendHistoryStageBody struct {
	Version int                    `json:"version"`
	Data    map[string]interface{} `json:"data"`
}
