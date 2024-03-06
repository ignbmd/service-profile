package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Classroom struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BranchCode  string             `json:"branch_code" bson:"branch_code"`
	Quota       int32              `json:"quota" bson:"quota"`
	QuotaFilled int32              `json:"quota_filled" bson:"quota_filled"`
	Description string             `json:"description" bson:"description"`
	Tags        []string           `json:"tags" bson:"tags"`
	Year        int32              `json:"year" bson:"year"`
	Status      string             `json:"status" bson:"status"`
	Title       string             `json:"title" bson:"title"`
	ClassCode   string             `json:"class_code" bson:"class_code"`
	ProductID   string             `json:"product_id" bson:"product_id"`
	IsOnline    bool               `json:"is_online" bson:"is_online"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
