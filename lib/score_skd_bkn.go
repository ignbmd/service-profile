package lib

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateScoreData(c *request.ScoreSkdBkn) (*mongo.InsertOneResult, error) {
	stdCol := db.Mongodb.Collection("score_skd_bkn")
	stdColStu := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Get Students by SmartbtwID
	filter := bson.M{"smartbtw_id": c.SmartBtwID, "deleted_at": nil}
	stdModels := models.Student{}
	err := stdColStu.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return nil, err
	}

	fil := bson.M{"student_id": stdModels.ID, "year": c.Year, "deleted_at": nil}
	stdSc := models.ScoreSkdBkn{}
	err = stdCol.FindOne(ctx, fil).Decode(&stdSc)
	if err == nil {
		return nil, fmt.Errorf("data tahun sudah ada sebelumnya")
	}

	payload := models.ScoreSkdBkn{
		StudentID: stdModels.ID,
		Year:      c.Year,
		ScoreTWK:  c.ScoreTWK,
		ScoreTIU:  c.ScoreTIU,
		ScoreTKP:  c.ScoreTKP,
		ScoreSKD:  c.ScoreSKD,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	res, err := stdCol.InsertOne(ctx, payload)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateScoreData(c *request.UpdateScoreSKDBKN, id primitive.ObjectID) error {
	opts := options.Update().SetUpsert(true)
	stdCol := db.Mongodb.Collection("score_skd_bkn")
	stdColStu := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Get Students by SmartbtwID
	filter := bson.M{"smartbtw_id": c.SmartBtwID, "deleted_at": nil}
	stdModels := models.Student{}
	err := stdColStu.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return err
	}

	//Get Data Score by ID
	fil := bson.M{"_id": id, "deleted_at": nil}
	scrModel := models.ScoreSkdBkn{}
	err = stdCol.FindOne(ctx, fil).Decode(&scrModel)
	if err != nil {
		return err
	}

	if scrModel.StudentID != stdModels.ID {
		return fmt.Errorf("data score with student not found")
	}

	payload := models.ScoreSkdBkn{
		StudentID: stdModels.ID,
		Year:      scrModel.Year,
		ScoreTWK:  c.ScoreTWK,
		ScoreTIU:  c.ScoreTIU,
		ScoreTKP:  c.ScoreTKP,
		ScoreSKD:  c.ScoreSKD,
		CreatedAt: scrModel.CreatedAt,
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
	// filUpd := bson.M{"_id": c.ID}
	update := bson.M{"$set": payload}
	_, err = stdCol.UpdateByID(ctx, scrModel.ID, update, opts)

	if err != nil {
		return err
	}

	return nil
}

func GetScoreDataByStudent(smID int) ([]models.ScoreSkdBkn, error) {
	scCol := db.Mongodb.Collection("score_skd_bkn")
	stdColStu := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Get Students by SmartbtwID
	filter := bson.M{"smartbtw_id": smID, "deleted_at": nil}
	stdModels := models.Student{}
	err := stdColStu.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return []models.ScoreSkdBkn{}, err
	}

	//Get Data Score by ID
	fil := bson.M{"student_id": stdModels.ID, "deleted_at": nil}
	scrModel := []models.ScoreSkdBkn{}
	cur, err := scCol.Find(ctx, fil)
	if err != nil {
		return []models.ScoreSkdBkn{}, err
	}

	if err = cur.All(ctx, &scrModel); err != nil {
		return []models.ScoreSkdBkn{}, err
	}

	return scrModel, nil
}

func GetSingleScoreData(ID primitive.ObjectID) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("score_skd_bkn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Get Data Score by ID
	pipel := aggregates.GetSingleScore(ID)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []bson.M{}, err
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return []bson.M{}, err
	}

	return results, nil
}

func GetManyStudentScoreByYear(ids []int, year int) ([]bson.M, error) {
	var (
		results    []bson.M
		studentIds []primitive.ObjectID
	)
	scCol := db.Mongodb.Collection("score_skd_bkn")
	stdColStu := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, val := range ids {
		filter := bson.M{"smartbtw_id": val, "deleted_at": nil}
		stdModels := models.Student{}
		err := stdColStu.FindOne(ctx, filter).Decode(&stdModels)

		if err != nil {
			// return []bson.M{}, err
			continue
		}
		studentIds = append(studentIds, stdModels.ID)
	}

	//Get Data
	pipel := aggregates.GetByManyStudent(studentIds, year)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []bson.M{}, err
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return []bson.M{}, err
	}
	return results, nil
}

func DeleteScoreSingleRecord(id primitive.ObjectID) error {
	scCol := db.Mongodb.Collection("score_skd_bkn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Get Data Score by ID
	fil := bson.M{"_id": id, "deleted_at": nil}
	scrModel := models.ScoreSkdBkn{}
	err := scCol.FindOne(ctx, fil).Decode(&scrModel)
	if err != nil {
		return err
	}
	var (
		tn  time.Time  = time.Now()
		tmn *time.Time = &tn
	)

	payload := models.ScoreSkdBkn{
		StudentID: scrModel.ID,
		Year:      scrModel.Year,
		ScoreTWK:  scrModel.ScoreTWK,
		ScoreTIU:  scrModel.ScoreTIU,
		ScoreTKP:  scrModel.ScoreTKP,
		ScoreSKD:  scrModel.ScoreSKD,
		CreatedAt: scrModel.CreatedAt,
		UpdatedAt: scrModel.UpdatedAt,
		DeletedAt: tmn,
	}

	update := bson.M{"$set": payload}
	_, err = scCol.UpdateByID(ctx, id, update)

	if err != nil {
		return err
	}

	return nil
}
