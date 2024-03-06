package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Avatar struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID int                `json:"smartbtw_id" bson:"smartbtw_id"`
	AvaType    string             `json:"ava_type" bson:"ava_type"`
	Style      int                `json:"style" bson:"style"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

