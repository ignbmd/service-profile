package lib

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func UpsertBranchData(c *request.UpsertBranchData) error {
	valRes, err := govalidator.ValidateStruct(c)
	if !valRes {
		return err
	}

	opts := options.Update().SetUpsert(true)
	stdCol := db.Mongodb.Collection("branchs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := models.Branch{
		BranchCode: c.BranchCode,
		BranchName: c.BranchName,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}

	filter := bson.M{"branch_code": c.BranchCode}
	update := bson.M{"$set": payload}

	_, err = stdCol.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

func GetBranches() (*[]models.Branch, error) {
	scCol := db.Mongodb.Collection("branchs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var brcModel = make([]models.Branch, 0)

	cur, err := scCol.Find(ctx, bson.D{{}})
	if err != nil {
		fmt.Println("error find")
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.Branch
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		brcModel = append(brcModel, model)
	}

	return &brcModel, nil

}

func GetBranchByBranchCode(bc string) (models.Branch, error) {
	var results models.Branch

	collection := db.Mongodb.Collection("branchs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"branch_code": bc,
	}

	err := collection.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		return models.Branch{}, fmt.Errorf("^data not found")
	}

	return results, nil
}
