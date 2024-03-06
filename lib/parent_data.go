package lib

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateParentData(c *request.CreateParentData) error {
	stdCol := db.Mongodb.Collection("students")
	stdColPar := db.Mongodb.Collection("parent_datas")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Get Students by SmartbtwID
	filter := bson.M{"smartbtw_id": c.SmartBtwID, "deleted_at": nil}
	stdModels := models.Student{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return err
	}

	// Determine createdAt timestamp
	var createdAt time.Time
	parentFilter := bson.M{"student_id": stdModels.ID}
	parentModel := models.ParentData{}

	parentErr := stdColPar.FindOne(ctx, parentFilter).Decode(&parentModel)
	if parentErr != nil {
		createdAt = time.Now()
	} else {
		createdAt = parentModel.CreatedAt
	}

	// Upsert the data using the given payload
	payload := models.ParentData{
		StudentID:    stdModels.ID,
		ParentName:   c.ParentName,
		ParentNumber: c.ParentNumber,
		CreatedAt:    createdAt,
		UpdatedAt:    time.Now(),
	}
	update := bson.M{"$set": payload}
	opts := options.Update().SetUpsert(true)
	_, err = stdColPar.UpdateOne(ctx, parentFilter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
