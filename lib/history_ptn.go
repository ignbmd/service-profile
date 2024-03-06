package lib

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateHistoryPtn(c *request.CreateHistoryPtn) (*mongo.InsertOneResult, error) {
	htpCol := db.Mongodb.Collection("history_ptn")
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "target_type": models.PTN, "is_active": true, "deleted_at": nil}
	stdModels := models.StudentTarget{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	payload := models.HistoryPtn{
		SmartBtwID:              c.SmartBtwID,
		TaskID:                  c.TaskID,
		PackageID:               c.PackageID,
		PackageType:             c.PackageType,
		ModuleCode:              c.ModuleCode,
		ModuleType:              c.ModuleType,
		PotensiKognitif:         c.PotensiKognitif,
		PenalaranMatematika:     c.PenalaranMatematika,
		LiterasiBahasaIndonesia: c.LiterasiBahasaIndonesia,
		LiterasiBahasaInggris:   c.LiterasiBahasaInggris,
		PengetahuanKuantitatif:  c.PengetahuanKuantitatif,
		PenalaranUmum:           c.PenalaranUmum,
		PengetahuanUmum:         c.PengetahuanUmum,
		PemahamanBacaan:         c.PemahamanBacaan,
		ProgramKey:              c.ProgramKey,
		Total:                   c.Total,
		Repeat:                  c.Repeat,
		ExamName:                c.ExamName,
		Grade:                   c.Grade,
		TargetID:                stdModels.ID,

		TargetScore:    c.TargetScore,
		SchoolOriginID: c.SchoolOriginID,
		SchoolOrigin:   c.SchoolOrigin,
		SchoolID:       c.SchoolID,
		SchoolName:     c.SchoolName,
		MajorID:        c.MajorID,
		MajorName:      c.MajorName,
		StudentName:    c.StudentName,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	res, err := htpCol.InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpsertHistoryPtn(c *request.CreateHistoryPtn) (*string, error) {
	var upid string
	opts := options.Update().SetUpsert(true)
	htpCol := db.Mongodb.Collection("history_ptn")
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "target_type": models.PTN, "is_active": true, "deleted_at": nil}
	stdModels := models.StudentTarget{}
	htsModels := models.HistoryPtn{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	filll := bson.M{"smartbtw_id": c.SmartBtwID, "task_id": c.TaskID}
	payload := models.HistoryPtn{
		SmartBtwID:              c.SmartBtwID,
		TaskID:                  c.TaskID,
		PackageID:               c.PackageID,
		PackageType:             c.PackageType,
		ModuleCode:              c.ModuleCode,
		ModuleType:              c.ModuleType,
		PotensiKognitif:         c.PotensiKognitif,
		PenalaranMatematika:     c.PenalaranMatematika,
		LiterasiBahasaIndonesia: c.LiterasiBahasaIndonesia,
		LiterasiBahasaInggris:   c.LiterasiBahasaInggris,
		PengetahuanKuantitatif:  c.PengetahuanKuantitatif,
		PenalaranUmum:           c.PenalaranUmum,
		PengetahuanUmum:         c.PengetahuanUmum,
		PemahamanBacaan:         c.PemahamanBacaan,
		Total:                   c.Total,
		Repeat:                  c.Repeat,
		ProgramKey:              c.ProgramKey,
		ExamName:                c.ExamName,
		Grade:                   c.Grade,
		TargetID:                stdModels.ID,
		Start:                   &c.Start,
		End:                     &c.End,
		IsLive:                  c.IsLive,
		SchoolOriginID:          c.SchoolOriginID,
		SchoolOrigin:            c.SchoolOrigin,
		SchoolID:                c.SchoolID,
		SchoolName:              c.SchoolName,
		MajorID:                 c.MajorID,
		MajorName:               c.MajorName,
		StudentName:             c.StudentName,
		TargetScore:             c.TargetScore,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}

	update := bson.M{"$set": payload}

	res, err := htpCol.UpdateOne(ctx, filll, update, opts)
	if err != nil {
		return nil, err
	}

	if res.UpsertedID == nil {
		err = htpCol.FindOne(ctx, filll).Decode(&htsModels)
		if err != nil {
			return nil, fmt.Errorf("data not found")
		}
		upid = htsModels.ID.Hex()
	} else {
		upid = res.UpsertedID.(primitive.ObjectID).Hex()
	}

	return &upid, nil
}

func UpdateHistoryPtn(c *request.UpdateHistoryPtn, id primitive.ObjectID) error {
	opts := options.Update().SetUpsert(true)
	htkCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "deleted_at": nil}
	stdTarget := models.HistoryPtn{}
	err := htkCol.FindOne(ctx, filter).Decode(&stdTarget)
	if err != nil {
		return fmt.Errorf("data not found")
	}

	payload := models.HistoryPtn{
		SmartBtwID:              stdTarget.SmartBtwID,
		TaskID:                  c.TaskID,
		PackageID:               c.PackageID,
		PackageType:             c.PackageType,
		ModuleCode:              c.ModuleCode,
		ModuleType:              c.ModuleType,
		PotensiKognitif:         c.PotensiKognitif,
		PenalaranMatematika:     c.PenalaranMatematika,
		LiterasiBahasaIndonesia: c.LiterasiBahasaIndonesia,
		LiterasiBahasaInggris:   c.LiterasiBahasaInggris,
		PengetahuanKuantitatif:  c.PengetahuanKuantitatif,
		PenalaranUmum:           c.PenalaranUmum,
		PengetahuanUmum:         c.PengetahuanUmum,
		PemahamanBacaan:         c.PemahamanBacaan,
		Total:                   c.Total,
		ProgramKey:              c.ProgramKey,
		Repeat:                  c.Repeat,
		ExamName:                c.ExamName,
		Grade:                   c.Grade,
		TargetID:                stdTarget.TargetID,
		CreatedAt:               stdTarget.CreatedAt,
		SchoolOriginID:          stdTarget.SchoolOriginID,
		SchoolOrigin:            stdTarget.SchoolOrigin,
		SchoolID:                stdTarget.SchoolID,
		SchoolName:              stdTarget.SchoolName,
		MajorID:                 stdTarget.MajorID,
		MajorName:               stdTarget.MajorName,
		StudentName:             stdTarget.StudentName,
		TargetScore:             stdTarget.TargetScore,
		UpdatedAt:               time.Now(),
		DeletedAt:               nil,
	}

	update := bson.M{"$set": payload}
	_, err1 := htkCol.UpdateByID(ctx, stdTarget.ID, update, opts)
	if err1 != nil {
		return err1
	}
	return nil
}

func DeleteHistoryPtn(id primitive.ObjectID) error {
	htnCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"deleted_at": time.Now(),
	}}
	_, err := htnCol.UpdateByID(ctx, id, update)

	if err != nil {
		return err
	}

	return nil
}

func GetHistoryPtnByID(id primitive.ObjectID) ([]bson.M, error) {
	var result []bson.M

	collection := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return []bson.M{}, fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err = cursor.All(ctx, &result)
	}

	sonic.Marshal(result)
	log.Println(err)

	return result, nil
}

func GetHistoryPtnBySmartBTWID(SmartBTWID int, params *request.HistoryPTNQueryParams) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if params.Limit != nil {
		if *params.Limit <= 0 {
			return nil, fmt.Errorf("limit must be a positive number")
		}
	}

	pipel := aggregates.GetStudentPTNHistoryScores(SmartBTWID, params)
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

func GetStudentPtnLastScore(SmartBtwID int, programKey string) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentPtnLastScore(SmartBtwID, programKey)
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

func GetStudentPtnAverage(SmartBtwID int, programKey string) ([]bson.M, error) {
	var result []bson.M
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentPtnAverage(SmartBtwID, programKey)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []bson.M{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []bson.M{}, err
	}

	return result, nil
}

func GetLast10StudentScorePtn(SmartBtwID int, programKey string) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetLast10StudentPtnScore(SmartBtwID, programKey)
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

func GetALLStudentScorePtn(SmartBtwID int) ([]models.HistoryPtn, error) {
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"smartbtw_id": SmartBtwID, "deleted_at": nil}
	var scrModel = make([]models.HistoryPtn, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)

	cur, err := scCol.Find(ctx, fil, opts)

	if err != nil {
		return []models.HistoryPtn{}, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.HistoryPtn
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}

	return scrModel, nil
}

func GetStudentAveragePtn(SmartBtwID int, programKey string) ([]bson.M, error) {
	var result []bson.M
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentPtnAverage(SmartBtwID, programKey)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []bson.M{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []bson.M{}, err
	}

	return result, nil
}

func InsertStudentPtnProfileElastic(data *request.StudentProfilePtnElastic, indexID string) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentTargetPtnIndexName()).
		Id(indexID).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func InsertStudentHistoryPtnElastic(data *request.CreateHistoryPtn, historyPTNID string) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentHistoryPtnIndexName()).
		Id(historyPTNID).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

// TODO: Add support for specific task_id for future
func GetStudentHistoryPTNElastic(smID int, isStagesHistory bool, programKey string) ([]request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtn
	var gres []request.CreateHistoryPtn

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if isStagesHistory {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type", "pre-uka"))
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type", "challenge-uka"))
	}

	if programKey != "" {
		if strings.ToLower(programKey) == "utbk" {
			elasticQuery = append(elasticQuery, elastic.NewBoolQuery().Should(elastic.NewMatchQuery("program_key", "utbk"), elastic.NewMatchQuery("module_type", "TESTING")))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("program_key", programKey))
		}
	}

	notInModuleType := elastic.NewBoolQuery().
		MustNot(elastic.NewTermsQuery("module_type", "PRE_TEST", "POST_TEST"))

	elasticQuery = append(elasticQuery, notInModuleType)

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtn{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtn{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtn))
	}

	return gres, nil
}
func GetStudentHistoryPTNElasticFilter(smID int, filter string, programKey string) ([]request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtn
	var gres []request.CreateHistoryPtn

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if filter != "" {
		if filter == "with_code" {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("module_type.keyword", "WITH_CODE"))
		} else if filter == "all-module" {
			boolQuery := elastic.NewBoolQuery()
			boolQuery.Should(
				elastic.NewTermsQuery("package_type.keyword", "challenge-uka", "pre-uka"),
			)

			elasticQuery = append(elasticQuery, boolQuery)
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type.keyword", filter))
		}
	}

	if programKey != "" {
		if strings.ToLower(programKey) == "utbk" {
			elasticQuery = append(elasticQuery, elastic.NewBoolQuery().Should(elastic.NewMatchQuery("program_key.keyword", "utbk"), elastic.NewMatchQuery("module_type.keyword", "TESTING")))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("program_key.keyword", programKey))
		}
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtn{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtn{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtn))
	}

	return gres, nil
}

func GetStudentHistoryPTNElasticFilterOffice(smID int, filter string, programKey string) ([]request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtn
	var gres []request.CreateHistoryPtn

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if filter != "" {

		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type.keyword", filter))
	}

	if programKey != "" {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("program_key.keyword", programKey))
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtn{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtn{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtn))
	}

	return gres, nil
}

func GetHistoryFreeSingleStudentPTN(smID int) (models.HistoryPtk, error) {
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "$or",
			Value: bson.A{
				bson.D{{Key: "module_type", Value: models.UkaFree}},
				bson.D{{Key: "module_type", Value: models.UkaCode}},
			},
		},
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "smartbtw_id", Value: smID}},
				bson.D{{Key: "deleted_at", Value: nil}},
			},
		},
	}
	// filter := bson.M{"smartbtw_id": smID, "module_type": models.UkaFree, "deleted_at": nil}
	stdModel := models.HistoryPtk{}
	err := scCol.FindOne(ctx, filter).Decode(&stdModel)
	if err != nil {
		return models.HistoryPtk{}, err
	}

	return stdModel, nil
}

func GetHistoryPremiumUKASingleStudentPTN(smID int) (models.HistoryPtk, error) {
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smID, "module_type": models.UkaPremium, "deleted_at": nil}
	stdModel := models.HistoryPtk{}
	err := scCol.FindOne(ctx, filter).Decode(&stdModel)
	if err != nil {
		return models.HistoryPtk{}, err
	}

	return stdModel, nil
}

func GetHistoryPackageUKASingleStudentPTN(smID int) (models.HistoryPtk, error) {
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smID, "module_type": models.Package, "deleted_at": nil}
	stdModel := models.HistoryPtk{}
	err := scCol.FindOne(ctx, filter).Decode(&stdModel)
	if err != nil {
		return models.HistoryPtk{}, err
	}

	return stdModel, nil
}

func GetALLStudentScorePtnPagination(smID int, limit *int64, page *int64, progKey string) ([]models.HistoryPtn, int64, error) {
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"smartbtw_id": smID, "deleted_at": nil, "module_type": bson.M{"$nin": []string{"PRE_TEST", "POST_TEST"}}}
	if progKey != "" {
		fil["program_key"] = progKey
	}
	var scrModel = make([]models.HistoryPtn, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)

	totalData, err := scCol.CountDocuments(ctx, fil)
	if err != nil {
		panic(err)
	}

	if limit != nil && page != nil {
		var itemLimit int64
		var itemPage int64
		if limit != nil {
			itemLimit = *limit
		}
		if page != nil {
			itemPage = *page
		}
		skip := ((itemPage * itemLimit) - itemLimit)

		fOpt := options.FindOptions{Limit: &itemLimit, Skip: &skip}

		cur, err := scCol.Find(ctx, fil, opts, &fOpt)
		if err != nil {
			return []models.HistoryPtn{}, totalData, err
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var model models.HistoryPtn
			e := cur.Decode(&model)
			if e != nil {
				log.Fatal(e)
			}
			scrModel = append(scrModel, model)
		}
	} else {
		cur, err := scCol.Find(ctx, fil, opts)
		if err != nil {
			return []models.HistoryPtn{}, totalData, err
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var model models.HistoryPtn
			e := cur.Decode(&model)
			if e != nil {
				log.Fatal(e)
			}
			scrModel = append(scrModel, model)
		}
	}

	return scrModel, totalData, nil

}

func RecommendStudentPTNSchool(smId uint, score float64, currentPTNStudyProgramID uint) error {
	msgBody := models.PTNStudentRecommendationStruct{
		Version: 2,
		Data: models.PTNStudentRecommendation{
			SmartbtwID:        smId,
			Score:             int(score),
			PtnStudyProgramID: currentPTNStudyProgramID,
		},
	}
	msgJson, err := sonic.Marshal(msgBody)
	if err != nil {
		return errors.New("error on marshaling json student recommendation body " + err.Error())
	}
	if db.Broker == nil {
		return errors.New("rabbit mq not available " + err.Error())
	}
	// Attempt to publish a message to the queue.
	if err = db.Broker.Publish(
		"exam.redis.generate-modules",
		"application/json",
		[]byte(msgJson), // message to publish
	); err != nil {
		return errors.New("error on publishing mq for student recommendation " + err.Error())
	}
	return nil
}

func UpdateDurationtPTN(req *request.BodyUpdateStudentDuration) error {
	htkCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": req.SmartbtwID, "task_id": req.TaskID, "repeat": req.Repeat, "deleted_at": nil}
	stdTarget := models.HistoryPtk{}
	err := htkCol.FindOne(ctx, filter).Decode(&stdTarget)
	if err != nil {
		return fmt.Errorf("data not found")
	}

	update := bson.M{"$set": bson.M{"start": req.Start, "end": req.End, "updated_at": time.Now()}}
	_, err1 := htkCol.UpdateOne(ctx, filter, update)
	if err1 != nil {
		return err1
	}
	return nil
}

func GetStudentHistoryPTNOnlyStage() ([]models.HistoryPtn, error) {
	var scrModel = make([]models.HistoryPtn, 0)
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetRecordOnlyStagePTN()
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []models.HistoryPtn{}, err
	}

	err = cursor.All(ctx, &scrModel)
	if err != nil {
		return []models.HistoryPtn{}, err
	}

	return scrModel, nil
}

func UpdateHistoryPtnElastic(req *request.BodyUpdateStudentDuration) error {
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", req.SmartbtwID), elastic.NewMatchQuery("task_id", req.TaskID), elastic.NewMatchQuery("repeat", req.Repeat))
	script1 := elastic.NewScript(`
ctx._source.start = params.start;
ctx._source.end = params.end;
`).Params(map[string]interface{}{
		"start": req.Start,
		"end":   req.End,
	})
	ctx := context.Background()

	_, err := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index("student_history_ptn").
		Query(bq).
		Script(script1).
		DoAsync(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetAllStudentHistoryPTN() ([]models.HistoryPtn, error) {
	// var result []bson.M
	var scrModel = make([]models.HistoryPtn, 0)
	scCol := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	pipel := aggregates.GetNotDeletedStuff()
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []models.HistoryPtn{}, err
	}

	err = cursor.All(ctx, &scrModel)
	if err != nil {
		return []models.HistoryPtn{}, err
	}

	return scrModel, nil
}

func UpdateTimestampHistoryPtnElastic(req *request.BodyUpdateStudentDuration) error {
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", req.SmartbtwID), elastic.NewMatchQuery("task_id", req.TaskID), elastic.NewMatchQuery("repeat", req.Repeat))
	script1 := elastic.NewScript(`
ctx._source.created_at = params.created_at;
ctx._source.updated_at = params.updated_at;
`).Params(map[string]interface{}{
		"created_at": req.Start,
		"updated_at": req.End,
	})
	ctx := context.Background()

	_, err := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(bq).
		Script(script1).
		DoAsync(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetStudentHistoryPTNElasticSpecific(smID int, filter string, programKey string) ([]request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtn
	var gres []request.CreateHistoryPtn

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if filter != "" {
		if filter == "with_code" {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("module_type.keyword", "WITH_CODE"))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("module_type.keyword", filter))
		}
	}

	if programKey != "" {
		if strings.ToLower(programKey) == "utbk" {
			elasticQuery = append(elasticQuery, elastic.NewBoolQuery().Should(elastic.NewMatchQuery("program_key.keyword", "utbk"), elastic.NewMatchQuery("module_type.keyword", "TESTING")))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("program_key.keyword", programKey))
		}
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtn{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtn{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtn))
	}

	return gres, nil
}
func FetchPTNRankingSchoolPurposes(taskId uint, schoolId string, limit int, page int, keyword string) (mockstruct.FetchRankingPTNBody, error) {
	ptnRankBody := mockstruct.FetchRankingPTNBody{}

	if limit > 1000 {
		return ptnRankBody, errors.New("limit cannot be more than 100 currently")
	}

	ctx := context.Background()

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("task_id", taskId),
	)

	if len(keyword) > 3 {
		query = query.Should(elastic.NewQueryStringQuery(fmt.Sprintf("*%s*", keyword)).Field("student_name.keyword")).
			Should(elastic.NewRegexpQuery("student_name.keyword", fmt.Sprintf("(?i).*%s.*", keyword)))
	}

	totalData, err := db.ElasticClient.Count().
		Index("student_history_ptn").
		Query(query).
		Do(ctx)

	if err != nil {
		return ptnRankBody, err
	}

	ptnRankBody.FetchRankingBase.RankingInformation.DataTotal = totalData
	ptnRankBody.FetchRankingBase.RankingInformation.Page = page

	from := (page - 1) * limit

	searchSource := elastic.NewSearchSource().
		Query(query).
		Size(int(limit)).
		From(int(from))

	searchResult, err := db.ElasticClient.Search().
		Index("student_history_ptn").
		SearchSource(searchSource).
		Sort("total", false).
		Do(context.Background()) // Sort by total in descending order Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	stPost := (page - 1) * limit
	for idx, hit := range searchResult.Hits.Hits {
		t := request.CreateHistoryPtn{}
		sonic.Unmarshal(hit.Source, &t)
		smData, err := GetStudentProfileElastic(t.SmartBtwID)
		if err != nil {
			fmt.Println("Error: ", t.SmartBtwID, " : ", err.Error())
			continue
		}

		bCode := "PT0000"
		bName := "Bimbel BTW (Kantor Pusat)"
		if smData.BranchCode != nil {
			bCode = *smData.BranchCode
			bName = *smData.BranchName
		}

		pAtt := helpers.RoundFloat((t.Total / float64(1)), 2)

		if math.IsNaN(pAtt) {
			pAtt = 0
		}
		if t.TargetScore == 0 {

			smDataProf, err := GetStudentProfilePTNElastic(t.SmartBtwID)
			if err != nil {
				fmt.Println("Error: ", t.SmartBtwID, " : ", err.Error())
				continue
			}
			t.TargetScore = smDataProf.TargetScore
		}
		percATT := helpers.RoundFloat((pAtt/t.TargetScore)*100, 2)
		if math.IsNaN(percATT) || math.IsInf(percATT, 0) {
			percATT = 0
		}

		stRank := stPost + (idx + 1)
		ptnRankBody.RankingData = append(ptnRankBody.RankingData, mockstruct.FetchRankingPTN{
			FetchRankingStudentBase: mockstruct.FetchRankingStudentBase{
				SmartBtwID:    smData.SmartbtwID,
				Email:         smData.Email,
				TaskID:        t.TaskID,
				PackageID:     t.PackageID,
				Name:          smData.Name,
				MajorID:       int(smData.MajorPTNID),
				MajorName:     smData.MajorNamePTN,
				SchoolID:      int(smData.SchoolPTNID),
				SchoolName:    smData.SchoolNamePTN,
				LastEdID:      smData.LastEdID,
				LastEdName:    smData.LastEdName,
				PassingChance: percATT,
				IsSameSchool:  smData.LastEdID == schoolId,
				BranchCode:    bCode,
				BranchName:    bName,
				Rank:          stRank,
			},
			ModuleCode:              t.ModuleCode,
			ModuleType:              t.ModuleType,
			PackageType:             t.PackageType,
			PotensiKognitif:         t.PotensiKognitif,
			PenalaranMatematika:     t.PenalaranMatematika,
			LiterasiBahasaIndonesia: t.LiterasiBahasaIndonesia,
			LiterasiBahasaInggris:   t.LiterasiBahasaInggris,
			PenalaranUmum:           t.PenalaranUmum,
			PengetahuanUmum:         t.PengetahuanUmum,
			PemahamanBacaan:         t.PemahamanBacaan,
			PengetahuanKuantitatif:  t.PengetahuanKuantitatif,
			ProgramKey:              t.ProgramKey,
			Title:                   t.Title,
			Start:                   t.Start,
			End:                     t.End,
			Total:                   t.Total,
		})

	}

	totalPages := math.Ceil(float64(totalData) / float64(limit))

	if math.IsNaN(float64(totalPages)) || math.IsInf(float64(totalPages), 0) {
		totalPages = 1
	}

	ptnRankBody.FetchRankingBase.RankingInformation.CurrentCountTotal = len(ptnRankBody.RankingData)
	ptnRankBody.FetchRankingBase.RankingInformation.PageTotal = int(totalPages)
	return ptnRankBody, nil
}

func GetRankPTNByTaskID(task_id uint) ([]map[string]interface{}, error) {
	resCol, err := GetHistoryPTNByTaskID(task_id)
	if err != nil {
		return nil, err
	} // fmt.Println("ini value dari get history \n", resCol)

	sort.Slice(resCol, func(i, j int) bool {
		return resCol[i].Total > resCol[j].Total
	})

	var pay []map[string]interface{}

	for i, v := range resCol {
		resEl, err := GetStudentProfileUKABySmartBtwID(v.SmartBtwID)
		if err != nil {
			return nil, err
		}

		// fmt.Println("ini value dari smartbtw: \n", resEl)
		var date string
		var start string
		var end string
		if v.Start != nil {
			date = DateFormat(v.Start)
			start = fmt.Sprintf("%s WIB", v.Start.Format("15:04:05"))
			end = fmt.Sprintf("%s WIB", v.End.Format("15:04:05"))
		} else {
			start = "-"
			end = "-"
			date = "-"
		}

		rank := i + 1

		for _, e := range resEl {
			var duration string
			if v.Start != nil && v.End != nil {
				duration = GetTotalWorkTime(*v.Start, *v.End)
			} else {
				duration = "-"
			}

			// var status bool
			// var statusTiu bool
			// var statusTkp bool
			// var statusTwk bool

			// if v.Twk >= 65 {
			// 	statusTwk = true
			// } else {
			// 	statusTwk = false
			// }

			// if v.Tkp >= 156 {
			// 	statusTkp = true
			// } else {
			// 	statusTkp = false
			// }

			// if v.Tiu >= 80 {
			// 	statusTiu = true
			// } else {
			// 	statusTiu = false
			// }

			// percATT := float64(0)
			// if v.Twk >= 65 && v.Tiu >= 80 && v.Tkp >= 156 {
			// 	status = true
			// 	pAtt := helpers.RoundFloat((v.Total / float64(1)), 2)

			// 	if math.IsNaN(pAtt) {
			// 		pAtt = 0
			// 	}

			// 	if v.TargetScore == 0 {
			// 		smDataProf, err := GetStudentProfilePTKElastic(v.SmartBtwID)
			// 		if err != nil {
			// 			fmt.Println("Error: ", v.SmartBtwID, " : ", err.Error())
			// 			continue
			// 		}
			// 		v.TargetScore = smDataProf.TargetScore
			// 	}
			// 	percATT = helpers.RoundFloat((pAtt/v.TargetScore)*100, 2)
			// 	if math.IsNaN(percATT) || math.IsInf(percATT, 0) {
			// 		percATT = 0
			// 	}
			// } else {
			// 	status = false
			// }

			pay = append(pay, map[string]interface{}{
				"rankptn": map[string]interface{}{
					"name":                      e.Name,
					"exam_name":                 v.ExamName,
					"task_id":                   v.TaskID,
					"instance_name":             e.SchoolNamePTN,
					"major_name":                e.MajorNamePTN,
					"start":                     start,
					"end":                       end,
					"duration":                  duration,
					"literasi_bahasa_indonesia": v.LiterasiBahasaIndonesia,
					"literasi_bahasa_inggris":   v.LiterasiBahasaInggris,
					"pemahaman_bacaan":          v.PemahamanBacaan,
					"penalaran_matematika":      v.PenalaranMatematika,
					"penalaran_umum":            v.PenalaranUmum,
					"pengetahuan_umum":          v.PengetahuanUmum,
					"pengetahuan_kuantitatif":   v.PengetahuanKuantitatif,
					"date":                      date,
					"rank":                      rank,
					"total":                     v.Total,
				},
			})

		}

	}

	return pay, nil
}

func GetStudentHistoryPTNElasticFilterOld(smID int, filter string, programKey string) ([]request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtn
	var gres []request.CreateHistoryPtn

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if filter != "" {
		if filter == "with_code" {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("module_type.keyword", "WITH_CODE"))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type.keyword", filter))
		}
	}

	if programKey != "" {
		if strings.ToLower(programKey) == "utbk" {
			elasticQuery = append(elasticQuery, elastic.NewBoolQuery().Should(elastic.NewMatchQuery("program_key.keyword", "utbk"), elastic.NewMatchQuery("module_type.keyword", "TESTING")))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("program_key.keyword", programKey))
		}
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtn{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtn{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtn))
	}

	return gres, nil
}

func GetHistoryPTNElasticPeforma(smID uint, typStg, mdltype string) ([]request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtn
	var gres []request.CreateHistoryPtn

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if typStg == "UMUM" && mdltype == "pre-uka" {
		boolQuery := elastic.NewBoolQuery()

		boolQuery.Should(
			elastic.NewBoolQuery().
				Must(
					elastic.NewMatchQuery("package_type.keyword", "pre-uka"),
					elastic.NewMatchQuery("module_type.keyword", "PLATINUM"),
				),
			elastic.NewBoolQuery().
				Must(
					elastic.NewMatchQuery("package_type.keyword", "general"),
					elastic.NewMatchQuery("module_type.keyword", "PLATINUM"),
				),
		)
		elasticQuery = append(elasticQuery, boolQuery)
	} else if typStg == "KELAS" && mdltype == "pre-uka" {
		elasticQuery = append(elasticQuery, elastic.NewBoolQuery().
			Must(elastic.NewMatchQuery("module_type.keyword", "PLATINUM")).
			Must(elastic.NewMatchQuery("package_type.keyword", "multi-stages-uka")))
	} else if typStg == "UMUM" && mdltype == "challenge-uka" {
		elasticQuery = append(elasticQuery, elastic.NewBoolQuery().
			Must(elastic.NewMatchQuery("module_type.keyword", "PREMIUM_TRYOUT")).
			Must(elastic.NewMatchQuery("package_type.keyword", "challenge-uka")))
	} else if typStg == "KELAS" && mdltype == "challenge-uka" {
		elasticQuery = append(elasticQuery, elastic.NewBoolQuery().
			Must(elastic.NewMatchQuery("module_type.keyword", "PREMIUM_TRYOUT")).
			Must(elastic.NewMatchQuery("package_type.keyword", "multi-stages-uka")))
	} else if typStg == "UMUM" && mdltype == "all-module" {
		boolQuery := elastic.NewBoolQuery()

		boolQuery.Should(
			elastic.NewBoolQuery().
				Must(
					elastic.NewMatchQuery("package_type.keyword", "pre-uka"),
					elastic.NewMatchQuery("module_type.keyword", "PLATINUM"),
				),
			elastic.NewBoolQuery().
				Must(
					elastic.NewMatchQuery("package_type.keyword", "general"),
					elastic.NewMatchQuery("module_type.keyword", "PLATINUM"),
				),
			elastic.NewBoolQuery().
				Must(
					elastic.NewMatchQuery("package_type.keyword", "challenge-uka"),
					elastic.NewMatchQuery("module_type.keyword", "PREMIUM_TRYOUT"),
				),
		)

		elasticQuery = append(elasticQuery, boolQuery)
	} else if typStg == "KELAS" && mdltype == "all-module" {
		boolQuery := elastic.NewBoolQuery()

		boolQuery.Should(
			elastic.NewBoolQuery().
				Must(
					elastic.NewMatchQuery("package_type.keyword", "multi-stages-uka"),
					elastic.NewMatchQuery("module_type.keyword", "PLATINUM"),
				),
			elastic.NewBoolQuery().
				Must(
					elastic.NewMatchQuery("package_type.keyword", "multi-stages-uka"),
					elastic.NewMatchQuery("module_type.keyword", "PREMIUM_TRYOUT"),
				),
		)

		elasticQuery = append(elasticQuery, boolQuery)
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtn{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtn{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtn))
	}

	return gres, nil
}

func GetHistoryPTNElasticByTaskID(smID uint, tskID uint) (request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtn

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))
	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("task_id", tskID))

	query := elastic.NewBoolQuery().Must(elasticQuery...)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1). // Set the size to 1 to retrieve only one record
		Do(ctx)

	if err != nil {
		return request.CreateHistoryPtn{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.CreateHistoryPtn{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		return item.(request.CreateHistoryPtn), nil
	}

	// This line should not be reached, as we are returning within the loop
	return request.CreateHistoryPtn{}, nil
}

func GetHistoryPTN(smartbtwID uint) ([]request.CreateHistoryPtn, error) {
	ctx := context.Background()

	var result []request.CreateHistoryPtn

	elasticQuery := elastic.NewBoolQuery().
		Must(elastic.NewMatchQuery("smartbtw_id", smartbtwID)).
		MustNot(elastic.NewTermsQuery("module_type", "TESTING", "WITH_CODE"))

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(elasticQuery).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return nil, nil
	}

	for _, item := range res.Each(reflect.TypeOf(request.CreateHistoryPtn{})) {
		result = append(result, item.(request.CreateHistoryPtn))
	}

	return result, nil
}

func GetHistoryPTNByPackageID(pckID uint) ([]request.CreateHistoryPtn, error) {
	resData := []request.CreateHistoryPtn{}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("package_id", pckID),
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var t request.CreateHistoryPtn
	for _, item := range res.Each(reflect.TypeOf(t)) {
		resData = append(resData, item.(request.CreateHistoryPtn))
	}

	return resData, nil

}
