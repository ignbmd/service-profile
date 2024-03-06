package lib

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func GetAllInterviewSessions() ([]models.InterviewSession, error) {
	col := db.Mongodb.Collection("interview_session")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"deleted_at": nil}
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.InterviewSession
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetSingleInterviewSessionByID(interviewSessionID primitive.ObjectID) (models.InterviewSession, error) {
	var res models.InterviewSession
	col := db.Mongodb.Collection("interview_session")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": interviewSessionID, "deleted_at": nil}
	err := col.FindOne(ctx, filter).Decode(&res)
	return res, err
}

func CreateInterviewSession(req *request.InterviewSessionRequest) (*mongo.InsertOneResult, error) {
	col := db.Mongodb.Collection("interview_session")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := bson.M{
		"name":        req.Name,
		"description": req.Description,
		"number":      req.Number,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
		"deleted_at":  nil,
	}

	res, err := col.InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateInterviewSession(interviewSessionID primitive.ObjectID, req *request.InterviewSessionRequest) (*mongo.UpdateResult, error) {
	col := db.Mongodb.Collection("interview_session")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": interviewSessionID, "deleted_at": nil}
	payload := bson.M{"$set": bson.M{
		"name":        req.Name,
		"description": req.Description,
		"number":      req.Number,
	}}
	res, err := col.UpdateOne(ctx, filter, payload)
	return res, err
}

func SoftDeleteInterviewSession(interviewSessionID primitive.ObjectID) (*mongo.UpdateResult, error) {
	col := db.Mongodb.Collection("interview_session")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": interviewSessionID, "deleted_at": nil}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	res, err := col.UpdateOne(ctx, filter, update)
	return res, err
}

func HardDeleteInterviewSession(interviewSessionID primitive.ObjectID) (*mongo.DeleteResult, error) {
	col := db.Mongodb.Collection("interview_session")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": interviewSessionID}
	res, err := col.DeleteOne(ctx, filter)
	return res, err
}
