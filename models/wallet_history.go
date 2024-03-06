package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string
type TransactionType string

const (
	IN       Status          = "IN"
	OUT      Status          = "OUT"
	TOPUP    TransactionType = "TOPUP"
	PURCHASE TransactionType = "PURCHASE"
	REWARD   TransactionType = "REWARD"
	GENERAL  TransactionType = "GENERAL"
)

type WalletHistory struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	WalletID    primitive.ObjectID `json:"wallet_id" bson:"wallet_id"`
	SmartbtwID  int                `json:"smartbtw_id" bson:"smartbtw_id"`
	AmountPay   float32            `json:"amount_pay" bson:"amount_pay"`
	Point       float32            `json:"point" bson:"point"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status"`
	Type        string             `json:"type" bson:"type"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
