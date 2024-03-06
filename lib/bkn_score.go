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

func UpsertBKNScore(c *request.UpsertBKNScore) error {
	var ca time.Time = time.Now()
	col := db.Mongodb.Collection("bkn_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := GetSingleBKNScoreByYearAndStudent(c.SmartBtwID, c.Year)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}

	if (res != models.BKNScore{}) {
		ca = res.CreatedAt
	}

	opts := options.Update().SetUpsert(true)
	pyl := bson.M{
		"smartbtw_id":     c.SmartBtwID,
		"twk":             c.Twk,
		"tiu":             c.Tiu,
		"tkp":             c.Tkp,
		"total":           c.Total,
		"year":            c.Year,
		"created_at":      ca,
		"updated_at":      time.Now(),
		"is_continue":     c.IsContinue,
		"bkn_rank":        c.BKNRank,
		"ptk_school_id":   c.PtkSchoolId,
		"ptk_school":      c.PtkSchool,
		"ptk_major_id":    c.PtkMajorId,
		"ptk_major":       c.PtkMajor,
		"bkn_test_number": c.BknTestNumber,
	}

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "year": c.Year, "deleted_at": nil}
	update := bson.M{"$set": pyl}

	_, err = col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil

}

func UpdateForSurvey(c *request.UpdateBKNScoreForSurvey) error {
	var ca time.Time = time.Now()
	col := db.Mongodb.Collection("bkn_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := GetSingleBKNScoreByYearAndStudent(c.SmartBtwID, c.Year)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "year": c.Year, "deleted_at": nil}
	payload := bson.M{
		"survey_status":   c.SurveyStatus,
		"reason":          c.Reason,
		"updated_at":      ca,
		"suggestion":      c.Suggestion,
		"returned_result": c.ReturnedResult,
	}
	update := bson.M{"$set": payload}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil

}

func UpdateForProdi(c *request.UpdateBKNScoreForProdi) error {
	var ca time.Time = time.Now()
	col := db.Mongodb.Collection("bkn_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := GetSingleBKNScoreByYearAndStudent(c.SmartBtwID, c.Year)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "year": c.Year, "deleted_at": nil}
	payload := bson.M{
		"ptk_school_id":      c.PtkSchoolId,
		"ptk_school":         c.PtkSchool,
		"ptk_major_id":       c.PtkMajorId,
		"ptk_major":          c.PtkMajor,
		"ptk_competition_id": c.PtkCompetitionID,
		"bkn_test_number":    c.BknTestNumber,
		"updated_at":         ca,
	}
	update := bson.M{"$set": payload}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil

}

func GetSingleBKNScoreByYearAndStudent(smID int, yr uint16) (models.BKNScore, error) {
	var result models.BKNScore
	col := db.Mongodb.Collection("bkn_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smID, "year": yr, "deleted_at": nil}
	err := col.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.BKNScore{}, err
	}

	return result, nil
}

func GetSingleBKNScoreByYearAndStudentUKA(smID int) (models.BKNScore, error) {
	var result models.BKNScore
	col := db.Mongodb.Collection("bkn_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smID, "deleted_at": nil}
	err := col.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.BKNScore{}, err
	}

	return result, nil
}

func GetBKNScoreByArrayOfIDStudent(smid []int, year uint16) (map[int]models.BKNScore, error) {
	var (
		scrModel = make([]models.BKNScore, 0)
		temp     = make(map[int]models.BKNScore)
	)
	col := db.Mongodb.Collection("bkn_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": bson.M{
			"$in": smid,
		},
		"year":       year,
		"deleted_at": nil,
	}

	cur, err := col.Find(ctx, filter)
	if err != nil {
		return map[int]models.BKNScore{}, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.BKNScore
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}

	for _, dt := range scrModel {
		temp[dt.SmartBtwID] = dt
	}

	return temp, nil
}

func GetBKNScoreByArrayOfStudentEmail(email []string, year uint16) ([]models.BKNScoreEmailEdutech, error) {
	var (
		scrModel   = make([]models.BKNScore, 0)
		scrStModel = make([]models.StudentSimpleData, 0)
		stData     = make(map[string][]models.StudentSimpleData)
		studentIds = []int{}
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	colStudents := db.Mongodb.Collection("students")
	col := db.Mongodb.Collection("bkn_score")

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
		return []models.BKNScoreEmailEdutech{}, err
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
			dt.SchoolName = res.SchoolNamePTK
			dt.SchoolID = res.SchoolPTKID
			dt.MajorName = res.MajorNamePTK
			dt.MajorID = res.MajorPTKID
			dt.OriginSchoolName = res.LastEdName
			dt.OriginSchoolID = res.LastEdID
			dt.AccountType = res.AccountType
		} else {
			dt.AccountType = "smartbtw"
		}
		stData[strings.ToLower(dt.Email)] = append(stData[strings.ToLower(dt.Email)], dt)
	}
	if year == 0 {
		year = uint16(time.Now().Year())
	}

	studentProfileData := []models.BKNScoreEmailEdutech{}

	filter := bson.M{
		"smartbtw_id": bson.M{
			"$in": studentIds,
		},
		"year":       year,
		"deleted_at": nil,
	}

	curBkn, err := col.Find(ctx, filter)
	if err != nil {
		return []models.BKNScoreEmailEdutech{}, err
	}

	defer curBkn.Close(ctx)

	for curBkn.Next(ctx) {
		var model models.BKNScore
		e := curBkn.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}

	for _, ems := range stData {
		stData := models.BKNScoreEmailEdutech{}

		for _, valStd := range ems {
			if valStd.AccountType == "btwedutech" {
				stData.Name = valStd.Name
				stData.Email = valStd.Email
				stData.Phone = valStd.Phone
				stData.BranchCode = valStd.BranchCode
				stData.AccountType = "btwedutech"
				stData.SchoolName = valStd.SchoolName
				stData.SchoolID = valStd.SchoolID
				stData.MajorName = valStd.MajorName
				stData.MajorID = valStd.MajorID
				stData.OriginSchoolName = valStd.OriginSchoolName
				stData.OriginSchoolID = valStd.OriginSchoolID
				stData.BTWEdutechID = valStd.SmartbtwID
			} else {
				if stData.Name == "" {
					stData.Name = valStd.Name
					stData.Email = valStd.Email
					stData.Phone = valStd.Phone
					stData.BranchCode = valStd.BranchCode
				}
				stData.AccountType = "btwedutech"
				stData.SmartBtwID = valStd.SmartbtwID
			}

			for _, dt := range scrModel {
				ns := dt
				if valStd.SmartbtwID == dt.SmartBtwID {
					stData.BKNScore = &ns
				}
			}
		}

		studentProfileData = append(studentProfileData, stData)
	}

	return studentProfileData, nil
}

func GetBKNScoreByArrayOfStudentEmailGDS(email []string, year uint16) ([]models.BKNScoreEmailEdutechGDS, error) {
	var (
		scrModel   = make([]models.BKNScore, 0)
		scrStModel = make([]models.StudentSimpleData, 0)
		stData     = make(map[string][]models.StudentSimpleData)
		studentIds = []int{}
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	colStudents := db.Mongodb.Collection("students")
	col := db.Mongodb.Collection("bkn_score")

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
		return []models.BKNScoreEmailEdutechGDS{}, err
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

			polbitType := "PUSAT"
			polbitCompetitionType := ""

			if res.PolbitCompetitionPTKID != 0 {
				locId := 0
				if res.PolbitLocationPTKID != 0 {
					locId = int(res.PolbitLocationPTKID)
				}
				chRes, err := GetCompetitionDataPTK(uint(res.MajorPTKID), uint(locId), res.Gender, res.PolbitTypePTK)
				if err == nil {

					if res.PolbitLocationPTKID != 0 {
						if res.PolbitTypePTK == "DAERAH_REGION" {
							polbitType = res.DomicileRegion
						} else if res.PolbitTypePTK == "DAERAH_PROVINCE" {
							polbitType = fmt.Sprintf("Provinsi %s", res.DomicileProvince)
						} else if strings.Contains(res.PolbitTypePTK, "AFIRMASI") {
							if res.PolbitTypePTK == "DAERAH_AFIRMASI_PROVINCE" {
								polbitType = fmt.Sprintf("%s (Afirmasi)", res.DomicileProvince)
							} else if strings.Contains(res.PolbitTypePTK, "PUSAT_AFIRMASI_PROVINCE") {
								polbitType = fmt.Sprintf("%s (Afirmasi)", res.DomicileProvince)
							} else {
								polbitType = fmt.Sprintf("%s (Afirmasi)", res.DomicileProvince)
							}
						}
					}

					if chRes.CompetitionType != nil {
						polbitCompetitionType = *chRes.CompetitionType
					}
				}
			}

			dt.SchoolName = res.SchoolNamePTK
			dt.SchoolID = res.SchoolPTKID
			dt.MajorName = res.MajorNamePTK
			dt.MajorID = res.MajorPTKID
			dt.OriginSchoolName = res.LastEdName
			dt.OriginSchoolID = res.LastEdID
			dt.AccountType = res.AccountType
			dt.FormationDesc = polbitType
			dt.FormationType = res.PolbitTypePTK
			dt.PolbitCompetitionType = polbitCompetitionType
			dt.PolbitCompetitionID = res.PolbitCompetitionPTKID
		} else {
			dt.AccountType = "smartbtw"
		}
		stData[strings.ToLower(dt.Email)] = append(stData[strings.ToLower(dt.Email)], dt)
	}
	if year == 0 {
		year = uint16(time.Now().Year())
	}

	studentProfileData := []models.BKNScoreEmailEdutechGDS{}

	filter := bson.M{
		"smartbtw_id": bson.M{
			"$in": studentIds,
		},
		"year":       year,
		"deleted_at": nil,
	}

	curBkn, err := col.Find(ctx, filter)
	if err != nil {
		return []models.BKNScoreEmailEdutechGDS{}, err
	}

	defer curBkn.Close(ctx)

	for curBkn.Next(ctx) {
		var model models.BKNScore
		e := curBkn.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}

	for _, ems := range stData {
		stData := models.BKNScoreEmailEdutechGDS{}

		for _, valStd := range ems {
			if valStd.AccountType == "btwedutech" {
				stData.Name = valStd.Name
				stData.Email = valStd.Email
				stData.Phone = valStd.Phone
				stData.BranchCode = valStd.BranchCode
				stData.AccountType = "btwedutech"
				stData.SchoolName = valStd.SchoolName
				stData.SchoolID = valStd.SchoolID
				stData.MajorName = valStd.MajorName
				stData.MajorID = valStd.MajorID
				stData.OriginSchoolName = valStd.OriginSchoolName
				stData.OriginSchoolID = valStd.OriginSchoolID
				stData.BTWEdutechID = valStd.SmartbtwID
				stData.FormationDesc = valStd.FormationDesc
				stData.FormationType = valStd.FormationType
				stData.PolbitCompetitionType = valStd.PolbitCompetitionType
				stData.PolbitCompetitionID = valStd.PolbitCompetitionID
			} else {
				if stData.Name == "" {
					stData.Name = valStd.Name
					stData.Email = valStd.Email
					stData.Phone = valStd.Phone
					stData.BranchCode = valStd.BranchCode
				}
				stData.AccountType = "btwedutech"
				stData.SmartBtwID = valStd.SmartbtwID
			}

			for _, dt := range scrModel {
				ns := dt
				if valStd.SmartbtwID == dt.SmartBtwID {
					stData.BKNScore = &ns
				}
			}
		}

		studentProfileData = append(studentProfileData, stData)
	}

	return studentProfileData, nil
}
