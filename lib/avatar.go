package lib

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateAvatar(c *request.CreateAvatar) (*mongo.InsertOneResult, error) {
	avaCol := db.Mongodb.Collection("avatars")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "ava_type": c.AvaType, "deleted_at": nil}
	avaModels := models.Avatar{}
	err := avaCol.FindOne(ctx, filter).Decode(&avaModels)
	if err == nil {
		return nil, fmt.Errorf("record already exist")
	}

	payload := models.Avatar{
		SmartbtwID: c.SmartbtwID,
		AvaType:    c.AvaType,
		Style:      c.Style,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}

	result, err1 := avaCol.InsertOne(ctx, &payload)
	if err1 != nil {
		return nil, err1
	}

	return result, nil
}

func UpdateAvatarSmartbtwID(c *request.UpdateAvatar) error {
	opt := options.Update().SetUpsert(true)
	avaCol := db.Mongodb.Collection("avatars")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "deleted_at": nil}
	ct := time.Now()
	avaModels := models.Avatar{}
	err := avaCol.FindOne(ctx, filter).Decode(&avaModels)
	if err == nil {
		ct = avaModels.CreatedAt
		// return fmt.Errorf("data not found")
	}

	payload := models.Avatar{
		SmartbtwID: c.SmartbtwID,
		AvaType:    c.AvaType,
		Style:      c.Style,
		CreatedAt:  ct,
		UpdatedAt:  time.Now(),
		DeletedAt:  nil,
	}

	update := bson.M{"$set": payload}
	_, err1 := avaCol.UpdateOne(ctx, filter, update, opt)
	if err1 != nil {
		return err1
	}

	return nil
}

func GetAvatarBySmartBtwIDAndType(req *request.BodyRequestAvatar) (models.Avatar, error) {
	var result models.Avatar

	collection := db.Mongodb.Collection("avatars")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": req.SmartbtwID, "ava_type": req.AvaType, "deleted_at": nil}
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return models.Avatar{}, fmt.Errorf("data not found")
	}

	return result, nil
}

func DeleteAvatarBySmartbtwID(req *request.BodyRequestAvatar) error {
	avaCol := db.Mongodb.Collection("avatars")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": req.SmartbtwID, "ava_type": req.AvaType, "deleted_at": nil}
	avaModels := models.Avatar{}
	err := avaCol.FindOne(ctx, filter).Decode(&avaModels)
	if err != nil {
		return err
	}

	var (
		tn  time.Time  = time.Now()
		tmn *time.Time = &tn
	)

	payload := models.Avatar{
		SmartbtwID: avaModels.SmartbtwID,
		AvaType:    avaModels.AvaType,
		Style:      avaModels.Style,
		CreatedAt:  avaModels.CreatedAt,
		UpdatedAt:  avaModels.UpdatedAt,
		DeletedAt:  tmn,
	}
	update := bson.M{"$set": payload}
	_, err1 := avaCol.UpdateOne(ctx, filter, update)

	if err1 != nil {
		return err1
	}
	return nil
}

func GetAvatarBySmartbtwID(smID int) ([]models.Avatar, error) {
	scCol := db.Mongodb.Collection("avatars")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"smartbtw_id": smID, "deleted_at": nil}
	var avaModel = make([]models.Avatar, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)

	cur, err := scCol.Find(ctx, fil, opts)
	if err != nil {
		return []models.Avatar{}, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.Avatar
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		avaModel = append(avaModel, model)
	}

	return avaModel, nil
}
