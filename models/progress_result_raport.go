package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProgressResultRaport struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Program    string             `json:"program" bson:"program"`
	Link       string             `json:"link" bson:"link"`
	UKAType    string             `json:"uka_type" bson:"uka_type"`
	StageType  string             `json:"stage_type" bson:"stage_type"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
