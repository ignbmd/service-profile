package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BKNScore struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartBtwID     int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Twk            float64            `json:"twk" bson:"twk"`
	Tiu            float64            `json:"tiu" bson:"tiu"`
	Tkp            float64            `json:"tkp" bson:"tkp"`
	Total          float64            `json:"total" bson:"total"`
	Year           uint16             `json:"year" bson:"year"`
	SurveyStatus   bool               `json:"survey_status" bson:"survey_status"`
	Reason         ReasonFailed       `json:"reason" bson:"reason"`
	Suggestion     string             `json:"suggestion" bson:"suggestion"`
	ReturnedResult string             `json:"returned_result" bson:"returned_result"`
	IsContinue     bool               `json:"is_continue" bson:"is_continue"`
	BKNRank        uint32             `json:"bkn_rank" bson:"bkn_rank"`
	PtkSchoolId    uint32             `json:"ptk_school_id" bson:"ptk_school_id"`
	PtkSchool      string             `json:"ptk_school" bson:"ptk_school"`
	PtkMajorId     uint32             `json:"ptk_major_id" bson:"ptk_major_id"`
	PtkMajor       string             `json:"ptk_major" bson:"ptk_major"`
	BknTestNumber  string             `json:"bkn_test_number" bson:"bkn_test_number"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt      *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type ReasonFailed struct {
	TestType     string `json:"test_type" bson:"test_type"`
	FailedReason string `json:"failed_reason" bson:"failed_reason"`
}

type BKNScoreEdutech struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartBtwID   int                `json:"smartbtw_id" bson:"smartbtw_id"`
	BTWEdutechID int                `json:"btwedutech_id" bson:"btwedutech_id"`
	Twk          float64            `json:"twk" bson:"twk"`
	Tiu          float64            `json:"tiu" bson:"tiu"`
	Tkp          float64            `json:"tkp" bson:"tkp"`
	Total        float64            `json:"total" bson:"total"`
	Year         uint16             `json:"year" bson:"year"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt    *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type BKNScoreEmailEdutech struct {
	SmartBtwID   int `json:"smartbtw_id" bson:"smartbtw_id"`
	BTWEdutechID int `json:"btwedutech_id" bson:"btwedutech_id"`

	SchoolName       string `json:"school_name" bson:"school_name"`
	SchoolID         uint   `json:"school_id" bson:"school_id"`
	MajorName        string `json:"major_name" bson:"major_name"`
	MajorID          uint   `json:"major_id" bson:"major_id"`
	OriginSchoolID   string `json:"origin_school_id" bson:"origin_school_id"`
	OriginSchoolName string `json:"origin_school_name" bson:"origin_school_name"`

	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email" bson:"email"`
	Phone       *string   `json:"phone" bson:"phone"`
	BranchCode  *string   `json:"branch_code" bson:"branch_code"`
	AccountType string    `json:"account_type" bson:"account_type"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	BKNScore    *BKNScore `json:"bkn_score" bson:"bkn_score"`
}

type BKNScoreEmailEdutechGDS struct {
	SmartBtwID   int `json:"smartbtw_id" bson:"smartbtw_id"`
	BTWEdutechID int `json:"btwedutech_id" bson:"btwedutech_id"`

	SchoolName            string `json:"school_name" bson:"school_name"`
	SchoolID              uint   `json:"school_id" bson:"school_id"`
	MajorName             string `json:"major_name" bson:"major_name"`
	MajorID               uint   `json:"major_id" bson:"major_id"`
	OriginSchoolID        string `json:"origin_school_id" bson:"origin_school_id"`
	OriginSchoolName      string `json:"origin_school_name" bson:"origin_school_name"`
	FormationType         string `json:"formation_type" bson:"formation_type"`
	FormationDesc         string `json:"formation_desc" bson:"formation_desc"`
	PolbitCompetitionType string `json:"polbit_competition_type" bson:"polbit_competition_type"`
	PolbitCompetitionID   uint   `json:"polbit_competition_id" bson:"polbit_competition_id"`

	Name        string    `json:"name" bson:"name"`
	Email       string    `json:"email" bson:"email"`
	Phone       *string   `json:"phone" bson:"phone"`
	BranchCode  *string   `json:"branch_code" bson:"branch_code"`
	AccountType string    `json:"account_type" bson:"account_type"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
	BKNScore    *BKNScore `json:"bkn_score" bson:"bkn_score"`
}
