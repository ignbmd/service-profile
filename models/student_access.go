package models

import (
	"time"

	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentAccess struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID       int                `json:"smartbtw_id" bson:"smartbtw_id"`
	DisallowedAccess string             `json:"disallowed_access" bson:"disallowed_access"`
	AppType          string             `json:"app_type" bson:"app_type"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt        *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type FlattenStudentAccess struct {
	SmartbtwID       int       `json:"smartbtw_id" bson:"smartbtw_id"`
	DisallowedAccess string    `json:"disallowed_access" bson:"disallowed_access"`
	AppType          string    `json:"app_type" bson:"app_type"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

func UnmarshalStudentAccess(data []byte) (StudentAccess, error) {
	var r StudentAccess
	err := sonic.Unmarshal(data, &r)
	return r, err
}

func (r *StudentAccess) Marshal() ([]byte, error) {
	return sonic.Marshal(r)
}
