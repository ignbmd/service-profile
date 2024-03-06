package lib

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

var (
	studentAccessCollection string = "student_access"
)

type StudentAllowedAccessStruct struct {
	DisallowedAccess string `json:"disallowed_access" bson:"disallowed_access"`
}

func CreateStudentAccess(c *request.CreateStudentAccess) (*mongo.InsertOneResult, error) {
	walCol := db.Mongodb.Collection(studentAccessCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "disallowed_access": c.DisallowedAccess, "deleted_at": nil}
	walModels := models.StudentAccess{}
	err := walCol.FindOne(ctx, filter).Decode(&walModels)
	if err == nil {
		return nil, fmt.Errorf("record already exist")
	}

	payload := models.StudentAccess{
		SmartbtwID:       c.SmartBtwID,
		DisallowedAccess: c.DisallowedAccess,
		AppType:          c.AppType,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		DeletedAt:        nil,
	}

	res, err1 := walCol.InsertOne(ctx, payload)

	if err1 != nil {
		return nil, err1
	}

	mdl := models.FlattenStudentAccess{
		SmartbtwID:       payload.SmartbtwID,
		DisallowedAccess: payload.DisallowedAccess,
		AppType:          payload.AppType,
		CreatedAt:        payload.CreatedAt,
		UpdatedAt:        payload.UpdatedAt,
	}
	err2 := InsertStudentDisallowedAccessElastic(&mdl)

	if err2 != nil {
		js, _ := sonic.Marshal(mdl)
		golog.Slack.ErrorWithData("failed to cache disallowed access to elastic", js, err2)
	}
	return res, nil
}

func CreateStudentAccessBulk(c *request.CreateStudentAccessBulk) ([]*mongo.InsertOneResult, error) {
	walCol := db.Mongodb.Collection(studentAccessCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res := []*mongo.InsertOneResult{}
	filter := bson.M{"smartbtw_id": c.SmartBtwID, "disallowed_access": c.DisallowedAccess, "deleted_at": nil}
	walModels := models.StudentAccess{}
	err := walCol.FindOne(ctx, filter).Decode(&walModels)
	if err == nil {
		return nil, fmt.Errorf("record already exist")
	}

	for _, k := range c.DisallowedAccess {

		payload := models.StudentAccess{
			SmartbtwID:       c.SmartBtwID,
			DisallowedAccess: k,
			AppType:          c.AppType,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			DeletedAt:        nil,
		}

		resIns, err1 := walCol.InsertOne(ctx, payload)

		if err1 != nil {
			return nil, err1
		}
		mdl := models.FlattenStudentAccess{
			SmartbtwID:       payload.SmartbtwID,
			DisallowedAccess: payload.DisallowedAccess,
			AppType:          payload.AppType,
			CreatedAt:        payload.CreatedAt,
			UpdatedAt:        payload.UpdatedAt,
		}
		err2 := InsertStudentDisallowedAccessElastic(&mdl)

		if err2 != nil {
			js, _ := sonic.Marshal(mdl)
			golog.Slack.ErrorWithData("failed to cache disallowed access to elastic", js, err2)
		}
		res = append(res, resIns)

	}

	return res, nil
}

func GetStudentAllowedAccess(SmartBTWID int, appType string) ([]models.StudentAccess, error) {
	var allowedAccess []models.StudentAccess
	collection := db.Mongodb.Collection(studentAccessCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentAccess(SmartBTWID, appType)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []models.StudentAccess{}, err
	}

	err = cursor.All(ctx, &allowedAccess)
	if err != nil {
		return []models.StudentAccess{}, err
	}
	return allowedAccess, nil
}

func GetStudentAllowedAccessByCode(SmartBTWID int, accessCode string, appType string) (*StudentAllowedAccessStruct, error) {
	collection := db.Mongodb.Collection(studentAccessCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": SmartBTWID, "disallowed_access": accessCode, "app_type": appType, "deleted_at": nil}
	stdModels := StudentAllowedAccessStruct{}
	err := collection.FindOne(ctx, filter).Decode(&stdModels)

	if err != nil {
		return nil, err
	}
	return &stdModels, nil
}

func DeleteStudentAccess(smId int, accessCode string, appType string) error {
	htnCol := db.Mongodb.Collection(studentAccessCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"smartbtw_id": smId, "disallowed_access": accessCode, "app_type": appType, "deleted_at": nil}
	fmt.Println(fil)
	htnModel := models.StudentAccess{}
	err := htnCol.FindOne(ctx, fil).Decode(&htnModel)
	if err != nil {
		return err
	}

	var (
		tn  time.Time  = time.Now()
		tmn *time.Time = &tn
	)

	payload := models.StudentAccess{
		SmartbtwID:       htnModel.SmartbtwID,
		DisallowedAccess: htnModel.DisallowedAccess,
		AppType:          htnModel.AppType,
		CreatedAt:        htnModel.CreatedAt,
		UpdatedAt:        htnModel.UpdatedAt,
		DeletedAt:        tmn,
	}

	update := bson.M{"$set": payload}
	_, err = htnCol.UpdateByID(ctx, htnModel.ID, update)

	if err != nil {
		return err
	}
	err2 := DeleteStudentDisallowedAccessElastic(payload.SmartbtwID, payload.DisallowedAccess, payload.AppType)

	if err2 != nil {
		golog.Slack.Error("failed to delete disallowed access cache in elastic", err2)
	}
	return nil
}

func DeleteStudentAccessBulk(smId int, accessCode []string, appType string) error {
	htnCol := db.Mongodb.Collection(studentAccessCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for _, accCode := range accessCode {
		fil := bson.M{"smartbtw_id": smId, "disallowed_access": accCode, "app_type": appType, "deleted_at": nil}
		htnModel := models.StudentAccess{}
		err := htnCol.FindOne(ctx, fil).Decode(&htnModel)
		if err != nil {
			continue
		}

		var (
			tn  time.Time  = time.Now()
			tmn *time.Time = &tn
		)

		payload := models.StudentAccess{
			SmartbtwID:       htnModel.SmartbtwID,
			DisallowedAccess: htnModel.DisallowedAccess,
			AppType:          htnModel.AppType,
			CreatedAt:        htnModel.CreatedAt,
			UpdatedAt:        htnModel.UpdatedAt,
			DeletedAt:        tmn,
		}

		update := bson.M{"$set": payload}
		_, err = htnCol.UpdateByID(ctx, htnModel.ID, update)

		if err != nil {
			return err
		}

		err2 := DeleteStudentDisallowedAccessElastic(payload.SmartbtwID, payload.DisallowedAccess, payload.AppType)

		if err2 != nil {
			golog.Slack.Error("failed to delete disallowed access cache in elastic", err2)
		}
	}
	return nil
}

func InsertStudentDisallowedAccessElastic(data *models.FlattenStudentAccess) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentDisallowedAccessIndexName()).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func DeleteStudentDisallowedAccessElastic(smartBtwId int, access string, appType string) error {
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("disallowed_access", access), elastic.NewMatchQuery("smartbtw_id", smartBtwId), elastic.NewMatchQuery("app_type", appType))
	_, err := elastic.NewDeleteByQueryService(db.ElasticClient).
		Index(db.GetStudentDisallowedAccessIndexName()).
		Query(bq).
		Do(context.Background())
	if err != nil {
		return err
	}

	_, errs := db.ElasticClient.Flush().Index(db.GetStudentDisallowedAccessIndexName()).Do(context.Background())
	if errs != nil {
		return err
	}
	return nil
}

func GetStudentAccessElastic(smID int, appType string) ([]models.FlattenStudentAccess, error) {
	ctx := context.Background()

	var t models.FlattenStudentAccess
	var gres []models.FlattenStudentAccess

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))
	if appType != "" {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("app_type", appType))
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentDisallowedAccessIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []models.FlattenStudentAccess{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []models.FlattenStudentAccess{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(models.FlattenStudentAccess))
	}

	return gres, nil
}

func GetStudentAccessByCodeElastic(req *request.GetStudentAccessElastic) ([]models.FlattenStudentAccess, error) {
	ctx := context.Background()

	var t models.FlattenStudentAccess
	var gres []models.FlattenStudentAccess

	var elasticQuery []elastic.Query

	if req.Code != "" {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("disallowed_access", req.Code))
	}
	if req.AppType != "" {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("app_type", req.AppType))
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentDisallowedAccessIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []models.FlattenStudentAccess{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []models.FlattenStudentAccess{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(models.FlattenStudentAccess))
	}

	return gres, nil
}

func GetSingleStudentAccessElastic(smID int, accessCode string, appType string) (*models.FlattenStudentAccess, error) {
	var err error
	sa := models.FlattenStudentAccess{}
	ctx := context.Background()

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))
	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("disallowed_access.keyword", accessCode))
	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("app_type.keyword", appType))

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentDisallowedAccessIndexName()).
		// SortBy(sort).
		Query(query).
		Size(1).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var t models.FlattenStudentAccess
	for _, item := range res.Each(reflect.TypeOf(t)) {
		sa = item.(models.FlattenStudentAccess)
		return &sa, nil
	}
	return nil, nil
}
