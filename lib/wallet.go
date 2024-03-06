package lib

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

type WalletTotalBalance struct {
	Balance int `json:"balance" bson:"balance"`
}
type WalletBalance struct {
	Balance int    `json:"balance" bson:"balance"`
	Type    string `json:"type" bson:"type"`
}

func CreateWallet(c *request.CreateWallet) (*mongo.InsertOneResult, error) {
	walCol := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "type": c.Type, "deleted_at": nil}
	walModels := models.StudentTarget{}
	err := walCol.FindOne(ctx, filter).Decode(&walModels)
	if err == nil {
		return nil, nil
	}

	payload := models.Wallet{
		SmartbtwID: c.SmartbtwID,
		Point:      c.Point,
		Type:       c.Type,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}
	res, err1 := walCol.InsertOne(ctx, payload)

	if err1 != nil {
		return nil, err1
	}

	return res, nil
}

func GetStudentWalletTotalBalance(SmartBTWID int) ([]WalletTotalBalance, error) {
	var wallet []WalletTotalBalance
	collection := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentWalletTotalBalance(SmartBTWID)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func GetStudentWalletBalance(SmartBTWID int) ([]WalletBalance, error) {
	var wallet []WalletBalance
	collection := db.Mongodb.Collection("wallets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentWalletBalance(SmartBTWID, true)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
