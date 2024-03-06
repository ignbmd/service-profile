package lib

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func UpsertSamaptaScore(c *request.UpsertSamaptaScore) error {
	var ca time.Time = time.Now()
	col := db.Mongodb.Collection("samapta_score")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := GetSingleSamaptaScoreByYearAndStudent(c.SmartBtwID, c.Year)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}

	if (res != models.SamaptaScore{}) {
		ca = res.CreatedAt
	}

	opts := options.Update().SetUpsert(true)
	pyl := bson.M{
		"smartbtw_id":   c.SmartBtwID,
		"gender":        c.Gender,
		"run_score":     c.RunScore,
		"pull_up_score": c.PullUpScore,
		"push_up_score": c.PushUpScore,
		"sit_up_score":  c.SitUpScore,
		"shuttle_score": c.ShuttleScore,
		"total":         c.Total,
		"year":          c.Year,
		"created_at":    ca,
		"updated_at":    time.Now(),
	}

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "year": c.Year, "deleted_at": nil}
	update := bson.M{"$set": pyl}

	_, err = col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func GetSingleSamaptaScoreByYearAndStudent(smID int, yr uint16) (models.SamaptaScore, error) {
	var result models.SamaptaScore
	col := db.Mongodb.Collection("samapta_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smID, "year": yr, "deleted_at": nil}
	err := col.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.SamaptaScore{}, err
	}

	return result, nil
}
func GetSamaptaScoreByArrayOfStudentEmail(email []string, year uint16) ([]models.SamaptaScoreEmailEdutech, error) {
	var (
		scrModel   = make([]models.SamaptaScore, 0)
		scrStModel = make([]models.StudentSimpleData, 0)
		stData     = make(map[string][]models.StudentSimpleData)
		studentIds = []int{}
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	colStudents := db.Mongodb.Collection("students")
	col := db.Mongodb.Collection("samapta_score")

	// Construct case-insensitive regular expression patterns for each email in the array
	regexPatterns := make([]primitive.Regex, len(email))
	for i, e := range email {
		pattern := fmt.Sprintf("(?i)%s", regexp.QuoteMeta(e))
		regexPatterns[i] = primitive.Regex{Pattern: pattern}
	}

	filStudents := bson.M{
		"$or": []bson.M{
			{"email": bson.M{"$in": regexPatterns}},
		},
	}

	cur, err := colStudents.Find(ctx, filStudents)
	if err != nil {
		return []models.SamaptaScoreEmailEdutech{}, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.StudentSimpleData
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrStModel = append(scrStModel, model)
	}
	for _, dt := range scrStModel {
		studentIds = append(studentIds, dt.SmartbtwID)
		_, k := stData[strings.ToLower(dt.Email)]
		if !k {
			stData[strings.ToLower(dt.Email)] = []models.StudentSimpleData{}
		}
		res, err := GetStudentProfileElastic(dt.SmartbtwID)
		if err == nil {
			dt.AccountType = res.AccountType
		} else {
			dt.AccountType = "smartbtw"
		}
		stData[strings.ToLower(dt.Email)] = append(stData[strings.ToLower(dt.Email)], dt)
	}
	if year == 0 {
		year = uint16(time.Now().Year())
	}

	studentProfileData := []models.SamaptaScoreEmailEdutech{}

	filter := bson.M{
		"smartbtw_id": bson.M{
			"$in": studentIds,
		},
		"year":       year,
		"deleted_at": nil,
	}

	curSamapta, err := col.Find(ctx, filter)
	if err != nil {
		return []models.SamaptaScoreEmailEdutech{}, err
	}

	defer curSamapta.Close(ctx)

	for curSamapta.Next(ctx) {
		var model models.SamaptaScore
		e := curSamapta.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}

	for _, ems := range stData {
		stData := models.SamaptaScoreEmailEdutech{}

		for _, valStd := range ems {
			if valStd.AccountType == "btwedutech" {
				stData.Name = valStd.Name
				stData.BTWEdutechID = valStd.SmartbtwID
			} else {
				if stData.Name == "" {
					stData.Name = valStd.Name
				}
				stData.AccountType = "btwedutech"
				stData.SmartBtwID = valStd.SmartbtwID
			}

			for _, dt := range scrModel {
				ns := dt
				if valStd.SmartbtwID == dt.SmartBtwID {
					stData.SamaptaScore = &ns
				}
			}
		}

		studentProfileData = append(studentProfileData, stData)
	}

	return studentProfileData, nil
}
