package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InterviewSession struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Number      int                `json:"number" bson:"number"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
