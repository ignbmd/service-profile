package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TargetType string

const (
	PTN       TargetType = "PTN"
	PTK       TargetType = "PTK"
	CPNS      TargetType = "CPNS"
	PRIMARY   Type       = "PRIMARY"
	SECONDARY Type       = "SECONDARY"
)

type StudentTarget struct {
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID          int                `json:"smartbtw_id" bson:"smartbtw_id"`
	SchoolID            int                `json:"school_id" bson:"school_id"`
	MajorID             int                `json:"major_id" bson:"major_id"`
	SchoolName          string             `json:"school_name" bson:"school_name"`
	MajorName           string             `json:"major_name" bson:"major_name"`
	TargetScore         float64            `json:"target_score" bson:"target_score"`
	TargetType          string             `json:"target_type" bson:"target_type"`
	PolbitType          string             `json:"polbit_type" bson:"polbit_type"`
	PolbitCompetitionID *int               `json:"polbit_competition_id" bson:"polbit_competition_id"`
	PolbitLocationID    *int               `json:"polbit_location_id" bson:"polbit_location_id"`
	Position            uint               `json:"position" bson:"position"`
	CanUpdate           bool               `json:"can_update" bson:"can_update"`
	IsActive            bool               `json:"is_active" bson:"is_active"`
	Type                string             `json:"type" bson:"type"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt           *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
