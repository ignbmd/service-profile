package lib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	requests "smartbtw.com/services/profile/request"
)

func CreateStudentTargetCpns(c *request.CreateStudentTargetCpns) (*mongo.InsertOneResult, error) {
	stdCol := db.Mongodb.Collection("student_target_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "deleted_at": nil}
	stdModels := models.StudentTargetCpns{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
	}

	if err == nil {
		return nil, nil
	}
	payload := models.StudentTargetCpns{
		SmartbtwID:        c.SmartbtwID,
		InstanceID:        c.InstanceID,
		InstanceName:      c.InstanceName,
		PositionID:        c.PositionID,
		PositionName:      c.PositionName,
		TargetScore:       c.TargetScore,
		FormationType:     c.FormationType,
		FormationLocation: c.FormationLocation,
		FormationCode:     c.FormationCode,
		CompetitionID:     c.CompetitionID,
		TargetType:        "CPNS",
		CanUpdate:         true,
		IsActive:          true,
		Position:          0,
		Type:              "PRIMARY",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		DeletedAt:         nil,
	}

	res, err1 := stdCol.InsertOne(ctx, payload)

	if err1 != nil {
		return nil, err1
	}

	err2 := InsertStudentTargetCpnsElastic(&request.StudentTargetCpnsElastic{
		SmartbtwID:        c.SmartbtwID,
		TargetScore:       c.TargetScore,
		TargetType:        c.TargetType,
		InstanceID:        c.InstanceID,
		InstanceName:      c.InstanceName,
		PositionID:        c.PositionID,
		PositionName:      c.PositionName,
		FormationType:     c.FormationType,
		FormationLocation: c.FormationLocation,
		FormationCode:     c.FormationCode,
		CompetitionID:     c.CompetitionID,
	})
	if err2 != nil {
		text := "Error when creating PTK data to elastic"
		js, _ := sonic.Marshal(c)
		golog.Slack.ErrorWithData(text, js, err)
	}

	ctxEls := context.Background()
	_, errEls := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", c.SmartbtwID)).
		Doc(map[string]interface{}{
			"instance_cpns_id":        c.InstanceID,
			"instance_cpns_name":      c.InstanceName,
			"position_cpns_id":        c.PositionID,
			"position_cpns_name":      c.PositionName,
			"formation_cpns_type":     c.FormationType,
			"formation_cpns_location": c.FormationLocation,
			"formation_cpns_code":     c.FormationCode,
			"competition_cpns_id":     c.CompetitionID,
			"created_at_cpns":         time.Now(),
		}).
		DocAsUpsert(true).
		Do(ctxEls)

	if errEls != nil {
		js, _ := sonic.Marshal(c)
		golog.Slack.ErrorWithData("error update elastic data", js, err)
	}

	return res, nil
}

func GetStudentTargetCPNS(smID int) (models.StudentTargetCpns, error) {
	var results models.StudentTargetCpns

	collection := db.Mongodb.Collection("student_target_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": smID,
		"is_active":   true,
		"deleted_at":  nil,
	}

	err := collection.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		return models.StudentTargetCpns{}, fmt.Errorf("^data not found")
	}

	return results, nil
}

func InsertStudentTargetCpnsElastic(data *request.StudentTargetCpnsElastic) error {
	var idst string
	ctx := context.Background()

	idst = fmt.Sprintf("%d_CPNS", data.SmartbtwID)

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentTargetCpnsIndexName()).
		Id(idst).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func GetAllStudentTargetCPNS(smID int) ([]models.StudentTargetCpns, error) {
	var results []models.StudentTargetCpns

	collection := db.Mongodb.Collection("student_target_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": smID,
		"is_active":   true,
		"deleted_at":  nil,
	}

	// err := collection.Find(ctx, filter).Decode(&results)
	// if err != nil {
	// 	return models.StudentTarget{}, fmt.Errorf("^data not found")
	// }
	sort := bson.M{"position": 1}
	opts := options.Find()
	opts.SetSort(sort)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		stdTarget := models.StudentTargetCpns{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		results = append(results, stdTarget)
	}

	return results, nil
}

func UpdateStudentTargetDataCPNSElastic(data *request.StudentTargetCpnsElastic) error {
	ctxEls := context.Background()
	_, err1 := db.ElasticClient.Update().
		Index(db.GetStudentTargetCpnsIndexName()).
		Id(fmt.Sprintf("%d_CPNS", data.SmartbtwID)).
		Doc(map[string]interface{}{
			"instance_id":        data.InstanceID,
			"position_id":        data.PositionID,
			"instance_name":      data.InstanceName,
			"position_name":      data.PositionName,
			"target_score":       data.TargetScore,
			"formation_type":     data.FormationType,
			"formation_location": data.FormationLocation,
			"formation_code":     data.FormationCode,
			"competition_id":     data.CompetitionID,
		}).
		DocAsUpsert(true).
		Do(ctxEls)
	if err1 != nil {
		return err1
	}

	_, errCreateStudentProfile := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", data.SmartbtwID)).
		Doc(map[string]interface{}{
			"instance_cpns_id":        data.InstanceID,
			"instance_cpns_name":      data.InstanceName,
			"position_cpns_id":        data.PositionID,
			"position_cpns_name":      data.PositionName,
			"formation_cpns_type":     data.FormationType,
			"formation_cpns_location": data.FormationLocation,
			"formation_cpns_code":     data.FormationCode,
			"competition_cpns_id":     data.CompetitionID,
			"created_at_cpns":         time.Now(),
		}).
		DocAsUpsert(true).
		Do(ctxEls)

	if errCreateStudentProfile != nil {
		return errCreateStudentProfile
	}

	return nil
}

func UpdateStudentTargetDataCPNSFormationYearElastic(data *request.StudentTargetCpnsElastic) error {
	ctxEls := context.Background()
	_, err1 := db.ElasticClient.Update().
		Index(db.GetStudentTargetCpnsIndexName()).
		Id(fmt.Sprintf("%d_CPNS", data.SmartbtwID)).
		Doc(map[string]interface{}{
			"instance_id":        data.InstanceID,
			"position_id":        data.PositionID,
			"instance_name":      data.InstanceName,
			"position_name":      data.PositionName,
			"target_score":       data.TargetScore,
			"formation_type":     data.FormationType,
			"formation_location": data.FormationLocation,
			"formation_code":     data.FormationCode,
			"competition_id":     data.CompetitionID,
			"formation_year":     data.FormationYear,
		}).
		DocAsUpsert(true).
		Do(ctxEls)
	if err1 != nil {
		return err1
	}

	_, errCreateStudentProfile := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", data.SmartbtwID)).
		Doc(map[string]interface{}{
			"instance_cpns_id":        data.InstanceID,
			"instance_cpns_name":      data.InstanceName,
			"position_cpns_id":        data.PositionID,
			"position_cpns_name":      data.PositionName,
			"formation_cpns_type":     data.FormationType,
			"formation_cpns_location": data.FormationLocation,
			"formation_cpns_code":     data.FormationCode,
			"competition_cpns_id":     data.CompetitionID,
			"created_at_cpns":         time.Now(),
		}).
		DocAsUpsert(true).
		Do(ctxEls)

	if errCreateStudentProfile != nil {
		return errCreateStudentProfile
	}

	return nil
}

func GetCompetitionChances(formationCode string, positionID uint, formationType string) ([]request.CompMapCompetitionCPNSChances, error) {
	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	bd := map[string]any{
		"formation_code": formationCode,
		"formation_type": formationType,
		"position_id":    positionID,
	}
	ns, _ := sonic.Marshal(bd)
	url := fmt.Sprintf("%s/competition-cpns-chance", conn)

	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(ns))

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map " + err.Error())
	}

	st := requests.ResponseContentCPNS{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map " + errs.Error())
	}

	return st.Data, nil
}

func DeleteStudentTargetCPNS(smartbtwID uint, id primitive.ObjectID) error {
	stdCol := db.Mongodb.Collection("student_target_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smartbtwID, "_id": id}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
			"is_active":  false},
	}

	_, err := stdCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func GetAllStudentTargetCPNSForUpdateTarget(smID int) ([]models.StudentTargetCpns, error) {
	var results []models.StudentTargetCpns

	collection := db.Mongodb.Collection("student_target_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": smID,
		"deleted_at":  nil,
	}

	// err := collection.Find(ctx, filter).Decode(&results)
	// if err != nil {
	// 	return models.StudentTarget{}, fmt.Errorf("^data not found")
	// }
	sort := bson.M{"position": 1}
	opts := options.Find()
	opts.SetSort(sort)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		stdTarget := models.StudentTargetCpns{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		results = append(results, stdTarget)
	}

	return results, nil
}
