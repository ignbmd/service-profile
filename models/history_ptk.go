package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModuleType string
type Grade string

const (
	UkaFree    ModuleType = "TESTING"
	UkaPremium ModuleType = "PREMIUM_TRYOUT"
	Package    ModuleType = "PLATINUM"
	UkaCode    ModuleType = "WITH_CODE"
	Basic      Grade      = "basic"
	Bronze     Grade      = "bronze"
	Gold       Grade      = "gold"
	Platinum   Grade      = "platinum"
)

type HistoryPtk struct {
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartBtwID          int                `json:"smartbtw_id" bson:"smartbtw_id"`
	TaskID              int                `json:"task_id" bson:"task_id"`
	PackageID           int                `json:"package_id" bson:"package_id"`
	PackageType         string             `json:"package_type" bson:"package_type"`
	ModuleCode          string             `json:"module_code" bson:"module_code"`
	ModuleType          string             `json:"module_type" bson:"module_type"`
	Twk                 float64            `json:"twk" bson:"twk"`
	Tiu                 float64            `json:"tiu" bson:"tiu"`
	Tkp                 float64            `json:"tkp" bson:"tkp"`
	TwkPass             float64            `json:"twk_pass" bson:"twk_pass"`
	TiuPass             float64            `json:"tiu_pass" bson:"tiu_pass"`
	TkpPass             float64            `json:"tkp_pass" bson:"tkp_pass"`
	Total               float64            `json:"total" bson:"total"`
	Repeat              int                `json:"repeat" bson:"repeat"`
	ExamName            string             `json:"exam_name" bson:"exam_name"`
	Grade               string             `json:"grade" bson:"grade"`
	TargetID            primitive.ObjectID `json:"target_id" bson:"target_id"`
	Start               *time.Time         `json:"start" bson:"start"`
	End                 *time.Time         `json:"end" bson:"end"`
	IsLive              bool               `json:"is_live" bson:"is_live"`
	StudentName         string             `json:"student_name"`
	SchoolOriginID      string             `json:"school_origin_id"`
	SchoolOrigin        string             `json:"school_origin"`
	TargetScore         float64            `json:"target_score" bson:"target_score"`
	SchoolID            int                `json:"school_id"`
	MajorID             int                `json:"major_id"`
	SchoolName          string             `json:"school_name"`
	MajorName           string             `json:"major_name"`
	PolbitType          string             `json:"polbit_type"`
	PolbitCompetitionID int                `json:"polbit_competition_id"`
	PolbitLocationID    int                `json:"polbit_location_id"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt           *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
