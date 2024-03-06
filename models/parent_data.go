package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParentData struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StudentID    primitive.ObjectID `json:"student_id" bson:"student_id"`
	ParentName   *string            `json:"parent_name" bson:"parent_name"`
	ParentNumber *string            `json:"parent_number" bson:"parent_number"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}
