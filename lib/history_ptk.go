package lib

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"sort"
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

func CreateHistoryPtk(c *request.CreateHistoryPtk) (*mongo.InsertOneResult, error) {
	htpCol := db.Mongodb.Collection("history_ptk")
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "target_type": models.PTK, "is_active": true, "deleted_at": nil}
	stdModels := models.StudentTarget{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	payload := models.HistoryPtk{
		SmartBtwID:          c.SmartBtwID,
		TaskID:              c.TaskID,
		PackageID:           c.PackageID,
		PackageType:         c.PackageType,
		ModuleCode:          c.ModuleCode,
		ModuleType:          c.ModuleType,
		Twk:                 c.Twk,
		Tiu:                 c.Tiu,
		Tkp:                 c.Tkp,
		TwkPass:             c.TwkPass,
		TiuPass:             c.TiuPass,
		TkpPass:             c.TkpPass,
		Total:               c.Total,
		Repeat:              c.Repeat,
		ExamName:            c.ExamName,
		Grade:               c.Grade,
		SchoolOriginID:      c.SchoolOriginID,
		SchoolOrigin:        c.SchoolOrigin,
		SchoolID:            c.SchoolID,
		SchoolName:          c.SchoolName,
		MajorID:             c.MajorID,
		MajorName:           c.MajorName,
		PolbitType:          c.PolbitType,
		PolbitCompetitionID: c.PolbitCompetitionID,
		PolbitLocationID:    c.PolbitLocationID,
		TargetID:            stdModels.ID,
		StudentName:         c.StudentName,
		TargetScore:         c.TargetScore,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	res, err := htpCol.InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpsertHistoryPtk(c *request.CreateHistoryPtk) (*string, error) {
	var upid string
	opts := options.Update().SetUpsert(true)
	htpCol := db.Mongodb.Collection("history_ptk")
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "target_type": models.PTK, "is_active": true, "deleted_at": nil}
	stdModels := models.StudentTarget{}
	htsModels := models.HistoryPtk{}
	filll := bson.M{"smartbtw_id": c.SmartBtwID, "task_id": c.TaskID}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	payload := models.HistoryPtk{
		SmartBtwID:          c.SmartBtwID,
		TaskID:              c.TaskID,
		PackageID:           c.PackageID,
		PackageType:         c.PackageType,
		ModuleCode:          c.ModuleCode,
		ModuleType:          c.ModuleType,
		Twk:                 c.Twk,
		Tiu:                 c.Tiu,
		Tkp:                 c.Tkp,
		TwkPass:             c.TwkPass,
		TiuPass:             c.TiuPass,
		TkpPass:             c.TkpPass,
		Total:               c.Total,
		Repeat:              c.Repeat,
		ExamName:            c.ExamName,
		Grade:               c.Grade,
		TargetID:            stdModels.ID,
		Start:               &c.Start,
		End:                 &c.End,
		IsLive:              c.IsLive,
		SchoolOriginID:      c.SchoolOriginID,
		SchoolOrigin:        c.SchoolOrigin,
		SchoolID:            c.SchoolID,
		SchoolName:          c.SchoolName,
		MajorID:             c.MajorID,
		MajorName:           c.MajorName,
		PolbitType:          c.PolbitType,
		PolbitCompetitionID: c.PolbitCompetitionID,
		PolbitLocationID:    c.PolbitLocationID,
		StudentName:         c.StudentName,
		TargetScore:         c.TargetScore,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
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

func UpdateHistoryPtk(c *request.UpdateHistoryPtk, id primitive.ObjectID) error {
	opts := options.Update().SetUpsert(true)
	htkCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "deleted_at": nil}
	stdTarget := models.HistoryPtk{}
	err := htkCol.FindOne(ctx, filter).Decode(&stdTarget)
	if err != nil {
		return fmt.Errorf("data not found")
	}

	payload := models.HistoryPtk{
		SmartBtwID:  stdTarget.SmartBtwID,
		TaskID:      c.TaskID,
		PackageID:   c.PackageID,
		PackageType: c.PackageType,
		ModuleCode:  c.ModuleCode,
		ModuleType:  c.ModuleType,
		Twk:         c.Twk,
		Tiu:         c.Tiu,
		Tkp:         c.Tkp,
		Total:       c.Total,
		Repeat:      c.Repeat,
		ExamName:    c.ExamName,
		Grade:       c.Grade,
		TargetID:    stdTarget.TargetID,
		CreatedAt:   stdTarget.CreatedAt,

		TargetScore:         stdTarget.TargetScore,
		SchoolOriginID:      stdTarget.SchoolOriginID,
		SchoolOrigin:        stdTarget.SchoolOrigin,
		SchoolID:            stdTarget.SchoolID,
		SchoolName:          stdTarget.SchoolName,
		MajorID:             stdTarget.MajorID,
		MajorName:           stdTarget.MajorName,
		PolbitType:          stdTarget.PolbitType,
		PolbitCompetitionID: stdTarget.PolbitCompetitionID,
		PolbitLocationID:    stdTarget.PolbitLocationID,
		StudentName:         stdTarget.StudentName,
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}

	update := bson.M{"$set": payload}
	_, err1 := htkCol.UpdateByID(ctx, stdTarget.ID, update, opts)
	if err1 != nil {
		return err1
	}
	return nil
}

func DeleteHistoryPtk(id primitive.ObjectID) error {
	htkCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{
		"deleted_at": time.Now(),
	}}

	_, err1 := htkCol.UpdateByID(ctx, id, update)

	if err1 != nil {
		return err1
	}

	return nil
}

func GetHistoryPtkByID(id primitive.ObjectID) ([]bson.M, error) {
	var results []bson.M

	collection := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return []bson.M{}, fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err = cursor.All(ctx, &results)
	}

	sonic.Marshal(results)
	log.Println(err)

	return results, nil
}

func GetStudentAveragePtk(SmartBtwID int) ([]bson.M, error) {
	var result []bson.M
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentAveragePtk(SmartBtwID)
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

func GetStudentLastScore(SmartBtwID int) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentLastScore(SmartBtwID)
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

func GetLast10StudentScorePtk(SmartBtwID int) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetLast10StudentScore(SmartBtwID)
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

func GetHistoryPtkBySmartBTWID(SmartBTWID int, params *request.HistoryPTKQueryParams) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if params.Limit != nil {
		if *params.Limit <= 0 {
			return nil, fmt.Errorf("limit must be a positive number")
		}
	}

	pipel := aggregates.GetStudentPTKHistoryScores(SmartBTWID, params)
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

func GetALLStudentScorePtk(SmartBtwID int) ([]models.HistoryPtk, error) {
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"smartbtw_id": SmartBtwID, "deleted_at": nil}
	var scrModel = make([]models.HistoryPtk, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)
	cur, err := scCol.Find(ctx, fil, opts)

	if err != nil {
		return []models.HistoryPtk{}, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.HistoryPtk
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}

	return scrModel, nil
}

func InsertStudentPtkProfileElastic(data *request.StudentProfilePtkElastic, indexID string) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentTargetPtkIndexName()).
		Id(indexID).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func InsertStudentHistoryPtkElastic(data *request.CreateHistoryPtk, historyPTNID string) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentHistoryPtkIndexName()).
		Id(historyPTNID).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func GetHistoryFreeSingleStudentPTK(smID int) (models.HistoryPtk, error) {
	scCol := db.Mongodb.Collection("history_ptk")
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

func GetHistoryPremiumUKASingleStudentPTK(smID int) (models.HistoryPtk, error) {
	scCol := db.Mongodb.Collection("history_ptk")
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

func GetHistoryPackageUKASingleStudentPTK(smID int) (models.HistoryPtk, error) {
	scCol := db.Mongodb.Collection("history_ptk")
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

func GetALLStudentScorePtkPagination(smID int, limit *int64, page *int64) ([]models.HistoryPtk, int64, error) {
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"smartbtw_id": smID, "deleted_at": nil, "module_type": bson.M{"$nin": []string{"PRE_TEST", "POST_TEST"}}}
	var scrModel = make([]models.HistoryPtk, 0)

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
			return []models.HistoryPtk{}, totalData, err
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var model models.HistoryPtk
			e := cur.Decode(&model)
			if e != nil {
				log.Fatal(e)
			}
			scrModel = append(scrModel, model)
		}
	} else {
		cur, err := scCol.Find(ctx, fil, opts)
		if err != nil {
			return []models.HistoryPtk{}, totalData, err
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var model models.HistoryPtk
			e := cur.Decode(&model)
			if e != nil {
				log.Fatal(e)
			}
			scrModel = append(scrModel, model)
		}
	}

	return scrModel, totalData, nil

}

// TODO: Add support for specific task_id for future
func GetStudentHistoryPTKElastic(smID int, isStagesHistory bool) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk
	var gres []request.CreateHistoryPtk

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if isStagesHistory {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type", "pre-uka"))
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type", "challenge-uka"))
	}
	notInModuleType := elastic.NewBoolQuery().
		MustNot(elastic.NewTermsQuery("module_type", "PRE_TEST", "POST_TEST"))

	elasticQuery = append(elasticQuery, notInModuleType)

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtk))
	}

	return gres, nil
}

func GetStudentHistoryPTKElasticFilter(smID int, filter string) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk
	var gres []request.CreateHistoryPtk

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

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtk))
	}

	return gres, nil
}

func GetStudentHistoryPTKOnlyStage() ([]models.HistoryPtk, error) {
	// var result []bson.M
	var scrModel = make([]models.HistoryPtk, 0)
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetRecordOnlyStage()
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []models.HistoryPtk{}, err
	}

	err = cursor.All(ctx, &scrModel)
	if err != nil {
		return []models.HistoryPtk{}, err
	}

	return scrModel, nil
}

func GetAllStudentHistoryPTK() ([]models.HistoryPtk, error) {
	// var result []bson.M
	var scrModel = make([]models.HistoryPtk, 0)
	scCol := db.Mongodb.Collection("history_ptk")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	pipel := aggregates.GetNotDeletedStuff()
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := scCol.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []models.HistoryPtk{}, err
	}

	err = cursor.All(ctx, &scrModel)
	if err != nil {
		return []models.HistoryPtk{}, err
	}

	return scrModel, nil
}

func UpdateTimestampHistoryPtkElastic(req *request.BodyUpdateStudentDuration) error {
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
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(bq).
		Script(script1).
		DoAsync(ctx)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDurationtPTK(req *request.BodyUpdateStudentDuration) error {
	htkCol := db.Mongodb.Collection("history_ptk")
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

func UpdateHistoryPtkElastic(req *request.BodyUpdateStudentDuration) error {
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
		Index("student_history_ptk").
		Query(bq).
		Script(script1).
		DoAsync(ctx)
	if err != nil {
		return err
	}
	return nil
}

func GetHistoryPTKElastic(smID uint, pckgType string) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk
	var gres []request.CreateHistoryPtk

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if pckgType == "ALL" {
		elasticQuery = append(elasticQuery, elastic.NewBoolQuery().
			MustNot(elastic.NewMatchQuery("module_type.keyword", "TESTING")))
	} else if pckgType == "all-module" {
		boolQuery := elastic.NewBoolQuery()
		boolQuery.Should(
			elastic.NewTermsQuery("package_type.keyword", "challenge-uka", "pre-uka"),
		)

		elasticQuery = append(elasticQuery, boolQuery)
	} else {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type.keyword", pckgType))
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtk))
	}

	return gres, nil
}

func GetStudentHistoryPTKElasticFilterOld(smID int, filter string) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk
	var gres []request.CreateHistoryPtk

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if filter != "" {
		if filter == "with_code" {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("module_type.keyword", "WITH_CODE"))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type.keyword", filter))
		}
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtk))
	}

	return gres, nil
}

func GetHistoryPTKElasticFetchStudentReport(smID uint, typStg, mdltype string) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk
	var gres []request.CreateHistoryPtk

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
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtk))
	}

	return gres, nil
}

func GetHistoryPTKElasticTypeStage(smID uint, modulType string) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk
	var gres []request.CreateHistoryPtk

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if modulType == "UMUM" {
		challengeQuery := elastic.NewBoolQuery().Must(elastic.NewTermQuery("module_type.keyword", "PREMIUM_TRYOUT"))

		preUkaQuery := elastic.NewBoolQuery().Must(elastic.NewTermQuery("package_type.keyword", "multi-stages-uka"))

		finalQuery := elastic.NewBoolQuery().
			Must(elasticQuery...).
			Should(challengeQuery, preUkaQuery)

		elasticQuery = append(elasticQuery, finalQuery)
	} else if modulType == "KELAS" {
		challengeQuery := elastic.NewBoolQuery().Must(elastic.NewTermQuery("module_type.keyword", "PREMIUM_TRYOUT"))

		preUkaQuery := elastic.NewBoolQuery().Must(elastic.NewTermQuery("package_type.keyword", "challenge-uka"))

		finalQuery := elastic.NewBoolQuery().
			Must(elasticQuery...).
			Should(challengeQuery, preUkaQuery)

		elasticQuery = append(elasticQuery, finalQuery)
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtk))
	}

	return gres, nil
}

func GetStudentHistoryPTKElasticSpecific(smID int, filter string) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk
	var gres []request.CreateHistoryPtk

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if filter != "" {
		if filter == "with_code" {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("module_type.keyword", "WITH_CODE"))
		} else {
			elasticQuery = append(elasticQuery, elastic.NewMatchQuery("module_type.keyword", filter))
		}
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryPtk))
	}

	return gres, nil
}
func FetchPTKRankingSchoolPurposes(taskId uint, schoolId string, limit int, page int, keyword string) (mockstruct.FetchRankingPTKBody, error) {
	ptkRankBody := mockstruct.FetchRankingPTKBody{}

	if limit > 1000 {
		return ptkRankBody, errors.New("limit cannot be more than 100 currently")
	}

	ctx := context.Background()

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("task_id", taskId),
	)

	if len(keyword) > 3 {
		query = query.Must(elastic.NewQueryStringQuery(fmt.Sprintf("*%s*", keyword)).Field("student_name.keyword"))
		// Should(elastic.NewRegexpQuery("student_name.keyword", fmt.Sprintf("(?i).*%s.*", keyword)))
	}

	totalData, err := db.ElasticClient.Count().
		Index("student_history_ptk").
		Query(query).
		Do(ctx)

	if err != nil {
		return ptkRankBody, err
	}

	ptkRankBody.FetchRankingBase.RankingInformation.DataTotal = totalData
	ptkRankBody.FetchRankingBase.RankingInformation.Page = page

	from := (page - 1) * limit

	searchSource := elastic.NewSearchSource().
		Query(query).
		Size(int(limit)).
		From(int(from))

	searchResult, err := db.ElasticClient.Search().
		Index("student_history_ptk").
		SearchSource(searchSource).
		SortBy(elastic.NewFieldSort("total").Desc(), elastic.NewFieldSort("is_all_passed").Desc(), elastic.NewFieldSort("updated_at")).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	stPost := (page - 1) * limit
	for idx, hit := range searchResult.Hits.Hits {
		t := request.CreateHistoryPtk{}
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

		percATT := float64(0)
		if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
			pAtt := helpers.RoundFloat((t.Total / float64(1)), 2)

			if math.IsNaN(pAtt) {
				pAtt = 0
			}

			if t.TargetScore == 0 {

				smDataProf, err := GetStudentProfilePTKElastic(t.SmartBtwID)
				if err != nil {
					fmt.Println("Error: ", t.SmartBtwID, " : ", err.Error())
					continue
				}
				t.TargetScore = smDataProf.TargetScore
			}
			percATT = helpers.RoundFloat((pAtt/t.TargetScore)*100, 2)
			if math.IsNaN(percATT) || math.IsInf(percATT, 0) {
				percATT = 0
			}

		}

		stRank := stPost + (idx + 1)

		ptkRankBody.RankingData = append(ptkRankBody.RankingData, mockstruct.FetchRankingPTK{
			FetchRankingStudentBase: mockstruct.FetchRankingStudentBase{
				SmartBtwID:    smData.SmartbtwID,
				Email:         smData.Email,
				TaskID:        t.TaskID,
				PackageID:     t.PackageID,
				Name:          smData.Name,
				MajorID:       int(smData.MajorPTKID),
				MajorName:     smData.MajorNamePTK,
				SchoolID:      int(smData.SchoolPTKID),
				SchoolName:    smData.SchoolNamePTK,
				LastEdID:      smData.LastEdID,
				LastEdName:    smData.LastEdName,
				PassingChance: percATT,
				IsSameSchool:  smData.LastEdID == schoolId,
				BranchCode:    bCode,
				BranchName:    bName,
				Rank:          stRank,
			},
			ModuleCode:    t.ModuleCode,
			ModuleType:    t.ModuleType,
			PackageType:   t.PackageType,
			Twk:           t.Twk,
			Tiu:           t.Tiu,
			Tkp:           t.Tkp,
			TwkPassStatus: t.Twk >= t.TwkPass,
			TiuPassStatus: t.Tiu >= t.TiuPass,
			TkpPassStatus: t.Tkp >= t.TkpPass,
			AllPassStatus: t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass,
			Title:         t.ExamName,
			Start:         t.Start,
			End:           t.End,
			Total:         t.Total,
		})

	}

	totalPages := math.Ceil(float64(totalData) / float64(limit))

	if math.IsNaN(float64(totalPages)) || math.IsInf(float64(totalPages), 0) {
		totalPages = 1
	}

	ptkRankBody.FetchRankingBase.RankingInformation.CurrentCountTotal = len(ptkRankBody.RankingData)
	ptkRankBody.FetchRankingBase.RankingInformation.PageTotal = int(totalPages)
	return ptkRankBody, nil
}

func GetHistoryPTKByTaskID(taskID uint) ([]request.GetHistoryUKAResultElastic, error) {
	resData := []request.GetHistoryUKAResultElastic{}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("task_id", taskID),
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var t request.GetHistoryUKAResultElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		resData = append(resData, item.(request.GetHistoryUKAResultElastic))
	}

	return resData, nil

}

func GetHistoryPTNByTaskID(taskID uint) ([]request.CreateHistoryPtnRanking, error) {
	resData := []request.CreateHistoryPtnRanking{}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("task_id", taskID),
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtnIndexName()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var t request.CreateHistoryPtnRanking
	for _, item := range res.Each(reflect.TypeOf(t)) {
		resData = append(resData, item.(request.CreateHistoryPtnRanking))
	}

	return resData, nil

}

func GetStudentProfileUKABySmartBtwID(smartbtwid int) ([]request.StudentProfileUKAElastic, error) {
	resData := []request.StudentProfileUKAElastic{}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smartbtwid),
	)
	res, err := db.ElasticClient.Search().
		Index(db.GetStudentProfileIndexName()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var t request.StudentProfileUKAElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		resData = append(resData, item.(request.StudentProfileUKAElastic))
	}

	return resData, nil
}

func GetRankPTKByTaskID(task_id uint) ([]map[string]interface{}, error) {
	resCol, err := GetHistoryPTKByTaskID(task_id)
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

			var status bool
			var statusTiu bool
			var statusTkp bool
			var statusTwk bool

			if v.Twk >= 65 {
				statusTwk = true
			} else {
				statusTwk = false
			}

			if v.Tkp >= 156 {
				statusTkp = true
			} else {
				statusTkp = false
			}

			if v.Tiu >= 80 {
				statusTiu = true
			} else {
				statusTiu = false
			}

			percATT := float64(0)
			if v.Twk >= 65 && v.Tiu >= 80 && v.Tkp >= 156 {
				status = true
				pAtt := helpers.RoundFloat((v.Total / float64(1)), 2)

				if math.IsNaN(pAtt) {
					pAtt = 0
				}

				if v.TargetScore == 0 {
					smDataProf, err := GetStudentProfilePTKElastic(v.SmartBtwID)
					if err != nil {
						fmt.Println("Error: ", v.SmartBtwID, " : ", err.Error())
						continue
					}
					v.TargetScore = smDataProf.TargetScore
				}
				percATT = helpers.RoundFloat((pAtt/v.TargetScore)*100, 2)
				if math.IsNaN(percATT) || math.IsInf(percATT, 0) {
					percATT = 0
				}
			} else {
				status = false
			}

			pay = append(pay, map[string]interface{}{
				"rankuka": map[string]interface{}{
					"name":           e.Name,
					"exam_name":      v.ExamName,
					"task_id":        v.TaskID,
					"instance_name":  e.SchoolNamePTK,
					"major_name":     e.MajorNamePTK,
					"start":          start,
					"end":            end,
					"duration":       duration,
					"twk":            v.Twk,
					"tiu":            v.Tiu,
					"tkp":            v.Tkp,
					"twk_status":     statusTwk,
					"tiu_status":     statusTiu,
					"tkp_status":     statusTkp,
					"date":           date,
					"rank":           rank,
					"total":          v.Total,
					"status":         status,
					"passing_chance": percATT,
				},
			})

		}

	}

	return pay, nil
}

func GetHistoryPTKElasticByTaskID(smID uint, tskID uint) (request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var t request.CreateHistoryPtk

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))
	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("task_id", tskID))

	query := elastic.NewBoolQuery().Must(elasticQuery...)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1). // Set the size to 1 to retrieve only one record
		Do(ctx)

	if err != nil {
		return request.CreateHistoryPtk{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.CreateHistoryPtk{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		return item.(request.CreateHistoryPtk), nil
	}

	// This line should not be reached, as we are returning within the loop
	return request.CreateHistoryPtk{}, nil
}

func GetHistoryPTK(smartbtwID uint) ([]request.CreateHistoryPtk, error) {
	ctx := context.Background()

	var result []request.CreateHistoryPtk

	elasticQuery := elastic.NewBoolQuery().
		Must(elastic.NewMatchQuery("smartbtw_id", smartbtwID)).
		MustNot(elastic.NewTermsQuery("module_type", "TESTING", "WITH_CODE"))

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(elasticQuery).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return nil, nil
	}

	for _, item := range res.Each(reflect.TypeOf(request.CreateHistoryPtk{})) {
		result = append(result, item.(request.CreateHistoryPtk))
	}

	return result, nil
}

func GetHistoryPTKByPackageID(pckID uint) ([]request.CreateHistoryPtk, error) {
	resData := []request.CreateHistoryPtk{}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("package_id", pckID),
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryPtkIndexName()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var t request.CreateHistoryPtk
	for _, item := range res.Each(reflect.TypeOf(t)) {
		resData = append(resData, item.(request.CreateHistoryPtk))
	}

	return resData, nil

}
