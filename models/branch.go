package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Branch struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BranchCode string             `json:"branch_code" bson:"branch_code"`
	BranchName string             `json:"branch_name" bson:"branch_name"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
}
