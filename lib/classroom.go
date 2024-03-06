package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateClassroom(req *request.CreateClassroom) (*mongo.InsertOneResult, error) {
	stdCol := db.Mongodb.Collection("classrooms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pyl := models.Classroom{
		ID:          req.ID,
		BranchCode:  req.BranchCode,
		Quota:       req.Quota,
		QuotaFilled: req.QuotaFilled,
		Description: req.Description,
		Tags:        req.Tags,
		Year:        req.Year,
		Status:      req.Status,
		Title:       req.Title,
		ClassCode:   "",
		ProductID:   req.ProductID,
		IsOnline:    req.IsOnline,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := stdCol.InsertOne(ctx, pyl)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateClassroom(req *request.UpdateClassroom) error {
	stdCol := db.Mongodb.Collection("classrooms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": req.ID}
	update := bson.M{"$set": bson.M{
		"branch_code":  req.BranchCode,
		"quota":        req.Quota,
		"quota_filled": req.QuotaFilled,
		"description":  req.Description,
		"tags":         req.Tags,
		"year":         req.Year,
		"status":       req.Status,
		"title":        req.Title,
		"product_id":   req.ProductID,
		"is_online":    req.IsOnline,
		"updated_at":   time.Now(),
	}}

	_, err := stdCol.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	//bulk update to elastic
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("classroom_id", req.ID))
	scpr := elastic.NewScript(`
	ctx._source.branch_code = params.branch_code; ctx._source.quota = params.quota;
	ctx._source.quota_filled = params.quota_filled; ctx._source.description = params.description; 
	ctx._source.tags = params.tags; ctx._source.year = params.year; ctx._source.status = params.status; 
	ctx._source.title = params.title; ctx._source.product_id = params.product_id; ctx._source.is_online = params.is_online`).
		Params(map[string]interface{}{
			"branch_code":  req.BranchCode,
			"quota":        req.Quota,
			"quota_filled": req.QuotaFilled,
			"description":  req.Description,
			"tags":         req.Tags,
			"year":         req.Year,
			"status":       req.Status,
			"title":        req.Title,
			"product_id":   req.ProductID,
			"is_online":    req.IsOnline,
		})

	_, err = elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetClassMemberIndexName()).
		Query(bq).
		Script(scpr).
		DoAsync(ctx)

	if err != nil {
		return err
	}

	return nil

}

func GetClassroomsByBranchCodes(bc string) ([]bson.M, error) {
	var results []bson.M

	collection := db.Mongodb.Collection("classrooms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"branch_code": bc, "status": "ONGOING"})
	if err != nil {
		return []bson.M{}, fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err = cursor.All(ctx, &results)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}
