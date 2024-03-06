package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentTargetCpns struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID        int                `json:"smartbtw_id" bson:"smartbtw_id"`
	InstanceID        int                `json:"instance_id" bson:"instance_id"`
	InstanceName      string             `json:"instance_name" bson:"instance_name"`
	PositionID        int                `json:"position_id" bson:"position_id"`
	PositionName      string             `json:"position_name" bson:"position_name"`
	TargetScore       float64            `json:"target_score" bson:"target_score"`
	FormationType     string             `json:"formation_type" bson:"formation_type"`
	FormationLocation string             `json:"formation_location" bson:"formation_location"`
	FormationCode     string             `json:"formation_code" bson:"formation_code"`
	CompetitionID     int                `json:"competition_id" bson:"competition_id"`
	Position          uint               `json:"position" bson:"position"`
	CanUpdate         bool               `json:"can_update" bson:"can_update"`
	IsActive          bool               `json:"is_active" bson:"is_active"`
	Type              string             `json:"type" bson:"type"`
	TargetType        string             `json:"target_type" bson:"target_type"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt         *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
