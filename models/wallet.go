package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Type string

const (
	DEFAULT Type = "DEFAULT"
	BONUS   Type = "BONUS"
)

type Wallet struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Point      float32            `json:"point" bson:"point"`
	Type       string             `json:"type" bson:"type"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type WalletRewardStruct struct {
	Version uint             `json:"version"`
	Data    WalletRewardBody `json:"data"`
}
type WalletRewardBody struct {
	SmartbtwID uint   `json:"smartbtw_id"`
	CodeName   string `json:"code_name"`
}

type UpdateFilter struct {
	ID   primitive.ObjectID `bson:"_id"`
	Type string             `bson:"type"`
}
