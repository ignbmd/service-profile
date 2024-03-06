package lib

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func UpsertStudentModuleProgress(c *request.StudentModuleProgress) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	collection := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if c.TaskID <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.Repeat <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.ModuleNo <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.ModuleTotal <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.ModuleNo > c.ModuleTotal {
		return nil, fmt.Errorf("data module number greater than module total")
	}

	payload := models.StudentModuleProgress{
		SmartbtwID:  c.SmartBtwID,
		TaskID:      c.TaskID,
		ModuleNo:    c.ModuleNo,
		Repeat:      c.Repeat,
		ModuleTotal: c.ModuleTotal,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "task_id": c.TaskID}
	update := bson.M{"$set": payload}

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	return result, err
}

func CreateStudentModuleProgress(c *request.StudentModuleProgress) (*mongo.InsertOneResult, error) {
	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Get Student Module by SmartbtwID Exists
	filter := bson.M{"smartbtw_id": c.SmartBtwID, "deleted_at": nil}
	stdModels := models.StudentModuleProgress{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err == nil {
		return nil, fmt.Errorf("already exist")
	}

	if c.SmartBtwID <= 0 {
		return nil, fmt.Errorf("smartbtw_id not valid")
	}
	if c.TaskID <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.Repeat <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.ModuleNo <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.ModuleTotal <= 0 {
		return nil, fmt.Errorf("must be greater than 0")
	}

	if c.ModuleNo > c.ModuleTotal {
		return nil, fmt.Errorf("data module number greater than module total")
	}

	payload := models.StudentModuleProgress{
		SmartbtwID:  c.SmartBtwID,
		TaskID:      c.TaskID,
		ModuleNo:    c.ModuleNo,
		Repeat:      c.Repeat,
		ModuleTotal: c.ModuleTotal,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	res, err := stdCol.InsertOne(ctx, payload)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateStudentModuleProgress(c *request.UpdateStudentModuleProgress, id primitive.ObjectID) error {
	opts := options.Update().SetUpsert(true)
	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get Students Module Progress by ID
	filter := bson.M{"_id": id, "deleted_at": nil}
	stdModule := models.StudentModuleProgress{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModule)
	if err != nil {
		return fmt.Errorf("data not found")
	}

	payload := models.StudentModuleProgress{
		SmartbtwID:  stdModule.SmartbtwID,
		TaskID:      c.TaskID,
		ModuleNo:    c.ModuleNo,
		Repeat:      c.Repeat,
		ModuleTotal: c.ModuleTotal,
		CreatedAt:   stdModule.CreatedAt,
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	if c.TaskID <= 0 {
		return fmt.Errorf("must be greater than 0")
	}

	if c.Repeat <= 0 {
		return fmt.Errorf("must be greater than 0")
	}

	if c.ModuleNo <= 0 {
		return fmt.Errorf("must be greater than 0")
	}

	if c.ModuleTotal <= 0 {
		return fmt.Errorf("must be greater than 0")
	}

	if c.ModuleNo > c.ModuleTotal {
		return fmt.Errorf("data module number greater than module total")
	}

	update := bson.M{"$set": payload}
	_, err = stdCol.UpdateByID(ctx, stdModule.ID, update, opts)

	if err != nil {
		return err
	}

	return nil
}

func GetStudentModuleProgressByTaskID(TaskID int) ([]bson.M, error) {
	var StudentModuleProgress []bson.M

	collection := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"task_id": TaskID})
	if err != nil {
		return []bson.M{}, fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err = cursor.All(ctx, &StudentModuleProgress)
		if err != nil {
			return []bson.M{}, err
		}
	}
	sonic.Marshal(StudentModuleProgress)
	log.Println(err)

	return StudentModuleProgress, nil
}

func SendExamReward(smId uint, score float64, targetScore float64, packageType string) error {
	codeName := "COIN_"
	isChallengeUKA := false
	isNotSupportedPackage := false

	switch strings.ToLower(packageType) {
	case "pre-uka":
		codeName = "COIN_500"
	case "challenge-uka":
		isChallengeUKA = true
	default:
		isNotSupportedPackage = true

	}

	if isNotSupportedPackage {
		return nil
	}

	if isChallengeUKA {
		passPercentage := (float64(score) / float64(targetScore)) * 100
		if passPercentage >= 100 {
			codeName = "COIN_10000"
		} else if passPercentage >= 51 && passPercentage < 76 {
			codeName = "COIN_7500"
		} else if passPercentage >= 26 && passPercentage < 51 {
			codeName = "COIN_5000"
		} else if passPercentage >= 25 {
			codeName = "COIN_2500"
		} else {
			codeName = "COIN_1500"
		}
	}

	msgBody := models.WalletRewardStruct{
		Version: 1,
		Data: models.WalletRewardBody{
			SmartbtwID: smId,
			CodeName:   codeName,
		},
	}

	msgJson, err := sonic.Marshal(msgBody)
	if err != nil {
		return errors.New("error on marshaling json student recommendation body " + err.Error())
	}
	if db.Broker == nil {
		return errors.New("rabbit mq not available " + err.Error())
	}

	LogEvent(
		"InsertScorePTN-SendExamReward",
		msgJson,
		"PUBLISH:history-reward.created",
		"publishing reward data",
		"INFO",
		fmt.Sprintf("profile-publish-event-%s", "history-reward.created"))
	// Attempt to publish a message to the queue.
	if err = db.Broker.Publish(
		"history-reward.created",
		"application/json",
		[]byte(msgJson), // message to publish
	); err != nil {
		return errors.New("error on publishing mq for student recommendation " + err.Error())
	}
	return nil
}

func SendUKAFreeReward(smId uint, ty string) error {
	var codeName string
	if ty == "PTK" {
		codeName = "CHANGEMAJOR_PTK_1"
	} else if ty == "CPNS" {
		codeName = "CHANGEMAJOR_CPNS_1"
	} else {
		codeName = "CHANGEMAJOR_PTN_1"
	}

	msgBody := models.WalletRewardStruct{
		Version: 1,
		Data: models.WalletRewardBody{
			SmartbtwID: smId,
			CodeName:   codeName,
		},
	}

	msgJson, err := sonic.Marshal(msgBody)
	if err != nil {
		return errors.New("error on marshaling json student reward body " + err.Error())
	}
	if db.Broker == nil {
		return errors.New("rabbit mq not available " + err.Error())
	}
	LogEvent(
		"InsertScore-SendUKAFreeReward",
		msgJson,
		"PUBLISH:history-reward.created",
		"publishing reward data",
		"INFO",
		fmt.Sprintf("profile-publish-event-%s", "history-reward.created"))
	// Attempt to publish a message to the queue.
	if err = db.Broker.Publish(
		"history-reward.created",
		"application/json",
		[]byte(msgJson), // message to publish
	); err != nil {
		return errors.New("error on publishing mq for reward " + err.Error())
	}
	return nil
}

func UpdatePassingPercentage(smId uint, program string, pkgId uint, percentage float64, slug string) error {

	msgBody := models.UpsertLiveRankDataStruct{
		Version: 1,
		Data: models.MessageUpdateLiveRanking{
			LRData: models.UpsertLiveRankData{
				SmartBTWID:             smId,
				Program:                program,
				PackageID:              pkgId,
				PassingScorePercentage: percentage,
				Slug:                   slug,
			},
		},
	}

	msgJson, err := sonic.Marshal(msgBody)
	if err != nil {
		return errors.New("error on marshaling json student passing update body " + err.Error())
	}
	if db.Broker == nil {
		return errors.New("rabbit mq not available " + err.Error())
	}

	LogEvent(
		fmt.Sprintf("InsertScore%s-UpdatePassingPercentage", strings.ToUpper(program)),
		msgJson,
		"PUBLISH:exam.live-ranking.update",
		"publishing live rank data",
		"INFO",
		fmt.Sprintf("profile-publish-event-%s", "exam.live-ranking.update"))
	// Attempt to publish a message to the queue.
	if err = db.Broker.Publish(
		"exam.live-ranking.update",
		"application/json",
		[]byte(msgJson), // message to publish
	); err != nil {
		return errors.New("error on publishing mq for student passing update " + err.Error())
	}
	return nil
}

func UpdatePassingPercentageCPNS(smId uint, program string, pkgId uint, percentage float64, slug string) error {

	msgBody := models.UpsertLiveRankDataStruct{
		Version: 1,
		Data: models.MessageUpdateLiveRanking{
			LRData: models.UpsertLiveRankData{
				SmartBTWID:             smId,
				Program:                program,
				PackageID:              pkgId,
				PassingScorePercentage: percentage,
				Slug:                   slug,
			},
		},
	}

	msgJson, err := sonic.Marshal(msgBody)
	if err != nil {
		return errors.New("error on marshaling json student passing update body " + err.Error())
	}
	if db.Broker == nil {
		return errors.New("rabbit mq not available " + err.Error())
	}

	LogEvent(
		fmt.Sprintf("InsertScore%s-UpdatePassingPercentage", strings.ToUpper(program)),
		msgJson,
		"PUBLISH:exam-cpns.live-ranking.update",
		"publishing live rank data",
		"INFO",
		fmt.Sprintf("profile-publish-event-%s", "exam-cpns.live-ranking.update"))
	// Attempt to publish a message to the queue.
	if err = db.Broker.Publish(
		"exam-cpns.live-ranking.update",
		"application/json",
		[]byte(msgJson), // message to publish
	); err != nil {
		return errors.New("error on publishing mq for student passing update " + err.Error())
	}
	return nil
}
