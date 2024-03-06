package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SamaptaScore struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartBtwID   int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Gender       bool               `json:"gender" bson:"gender"`
	RunScore     float32            `json:"run_score" bson:"run_score"`
	PullUpScore  float32            `json:"pull_up_score" bson:"pull_up_score"`
	PushUpScore  float32            `json:"push_up_score" bson:"push_up_score"`
	SitUpScore   float32            `json:"sit_up_score" bson:"sit_up_score"`
	ShuttleScore float32            `json:"shuttle_score" bson:"shuttle_score"`
	Total        float32            `json:"total" bson:"total"`
	Year         uint16             `json:"year" bson:"year"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt    *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type SamaptaScoreEmailEdutech struct {
	SmartBtwID   int           `json:"smartbtw_id" bson:"smartbtw_id"`
	BTWEdutechID int           `json:"btwedutech_id" bson:"btwedutech_id"`
	Name         string        `json:"name" bson:"name"`
	AccountType  string        `json:"account_type" bson:"account_type"`
	SamaptaScore *SamaptaScore `json:"samapta_score" bson:"samapta_score"`
}
