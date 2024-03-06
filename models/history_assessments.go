package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistoryAssessmentsScores struct {
	Position          uint                          `json:"position" bson:"position"`
	CategoryID        uint                          `json:"category_id" bson:"category_id"`
	CategoryName      string                        `json:"category_name" bson:"category_name"`
	CategoryAlias     string                        `json:"category_alias" bson:"category_alias"`
	Subtests          []HistoryAssessmentsSubScores `json:"subtests" bson:"subtests"`
	WrongAnswer       uint                          `json:"wrong_answer" bson:"wrong_answer"`
	CorrectAnswer     uint                          `json:"correct_answer" bson:"correct_answer"`
	EmptyAnswer       uint                          `json:"empty_answer" bson:"empty_answer"`
	PassingGrade      float64                       `json:"passing_grade" bson:"passing_grade"`
	PassingPercentage float64                       `json:"passing_percentage" bson:"passing_percentage"`
	PassingIndex      float64                       `json:"passing_index" bson:"passing_index"`
	IsPass            bool                          `json:"is_pass" bson:"is_pass"`
	Scores            []HistoryScores               `json:"scores" bson:"scores"`
}
type HistoryAssessmentsSubScores struct {
	SubID             uint            `json:"sub_id" bson:"sub_id"`
	Position          uint            `json:"position" bson:"position"`
	SubName           string          `json:"sub_name" bson:"sub_name"`
	SubAlias          string          `json:"sub_alias" bson:"sub_alias"`
	WrongAnswer       uint            `json:"wrong_answer" bson:"wrong_answer"`
	CorrectAnswer     uint            `json:"correct_answer" bson:"correct_answer"`
	EmptyAnswer       uint            `json:"empty_answer" bson:"empty_answer"`
	PassingPercentage float64         `json:"passing_percentage" bson:"passing_percentage"`
	PassingIndex      float64         `json:"passing_index" bson:"passing_index"`
	Scores            []HistoryScores `json:"scores" bson:"scores"`
}

type HistoryScores struct {
	ScoreType string  `json:"score_type" bson:"score_type"`
	Value     float64 `json:"value" bson:"value"`
}

type HistoryAssessments struct {
	ID             primitive.ObjectID         `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartBtwID     int                        `json:"smartbtw_id" bson:"smartbtw_id"`
	PackageID      int                        `json:"package_id" bson:"package_id"`
	AssessmentCode string                     `json:"assessment_code" bson:"assessment_code"`
	ModuleCode     string                     `json:"module_code" bson:"module_code"`
	ModuleType     string                     `json:"module_type" bson:"module_type"`
	PackageType    string                     `json:"package_type" bson:"package_type"`
	ScoreType      string                     `json:"score_type" bson:"score_type"`
	Total          float64                    `json:"total" bson:"total"`
	ExamName       string                     `json:"exam_name" bson:"exam_name"`
	Start          time.Time                  `json:"start" bson:"start"`
	End            time.Time                  `json:"end" bson:"end"`
	IsLive         bool                       `json:"is_live" bson:"is_live"`
	StudentName    string                     `json:"student_name" bson:"student_name"`
	StudentEmail   string                     `json:"student_email" bson:"student_email"`
	Program        string                     `json:"program" bson:"program"`
	ProgramType    string                     `json:"program_type" bson:"program_type"`
	ProgramVersion int                        `json:"program_version" bson:"program_version"`
	Scores         []HistoryAssessmentsScores `json:"scores" bson:"scores"`
	CreatedAt      time.Time                  `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at" bson:"updated_at"`
	DeletedAt      *time.Time                 `json:"deleted_at" bson:"deleted_at"`
}
