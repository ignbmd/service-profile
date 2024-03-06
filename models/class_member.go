package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassMember struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID  int32              `json:"smartbtw_id" bson:"smartbtw_id"`
	ClassroomID primitive.ObjectID `json:"classroom_id" bson:"classroom_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type ClassMemberElastic struct {
	ID          string    `json:"id"`
	SmartbtwID  int32     `json:"smartbtw_id"`
	ClassroomID string    `json:"classroom_id"`
	BranchCode  string    `json:"branch_code"`
	Quota       int32     `json:"quota"`
	QuotaFilled int32     `json:"quota_filled"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Year        int32     `json:"year"`
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	ClassCode   string    `json:"class_code"`
	ProductID   string    `json:"product_id"`
	IsOnline    bool      `json:"is_online"`
	CreatedAt   time.Time `json:"created_at"`
}
