package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/klauspost/lctime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func UpsertProgressResultRaport(req models.ProgressResultRaport) error {
	collection := db.Mongodb.Collection("progress_result_raports")
	filter := bson.M{
		"smartbtw_id": req.SmartbtwID,
		"program":     req.Program,
		"uka_type":    req.UKAType,
		"stage_type":  req.StageType,
	}

	update := bson.M{
		"$set": bson.M{
			"link":       req.Link,
			"updated_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func SendToGenerateRaport(program string, smID uint, stageType string) error {
	//stage type UMUM dan KELAS
	stg := "UMUM"
	if strings.ToLower(stageType) == "multi-stages-uka" {
		stg = "KELAS"
	}
	switch program {
	case "PTK":
		ukaType := []string{"ALL_MODULE", "PRE_UKA", "UKA_STAGE"}
		for _, ty := range ukaType {

			msgBodys := map[string]any{
				"version": 1,
				"data": map[string]any{
					"smartbtw_id": smID,
					"uka_type":    ty,
					"program":     "PTK",
					"stage_type":  stg,
				},
			}

			msgJsons, errs := sonic.Marshal(msgBodys)
			if errs == nil && db.Broker != nil {
				_ = db.Broker.Publish(
					"progress-result-raport.build.queue",
					"application/json",
					[]byte(msgJsons), // message to publish
				)
			}

		}
	case "PTN":
		ukaType := []string{"ALL_MODULE", "PRE_UKA", "UKA_STAGE"}
		for _, ty := range ukaType {
			msgBodys := map[string]any{
				"version": 1,
				"data": map[string]any{
					"smartbtw_id": smID,
					"uka_type":    ty,
					"program":     "PTN",
					"stage_type":  stg,
				},
			}

			msgJsons, errs := sonic.Marshal(msgBodys)
			if errs == nil && db.Broker != nil {
				_ = db.Broker.Publish(
					"progress-result-raport.build.queue",
					"application/json",
					[]byte(msgJsons), // message to publish
				)
			}

		}
	case "CPNS":
		ukaType := []string{"ALL_MODULE", "PRE_UKA", "UKA_STAGE"}
		for _, ty := range ukaType {
			msgBodys := map[string]any{
				"version": 1,
				"data": map[string]any{
					"smartbtw_id": smID,
					"uka_type":    ty,
					"program":     "CPNS",
					"stage_type":  stg,
				},
			}

			msgJsons, errs := sonic.Marshal(msgBodys)
			if errs == nil && db.Broker != nil {
				_ = db.Broker.Publish(
					"progress-result-raport.build.queue",
					"application/json",
					[]byte(msgJsons), // message to publish
				)
			}
		}
	}
	return nil
}

func GetProgressResultRaport(program string, ukaType string) ([]models.ProgressResultRaport, error) {
	collection := db.Mongodb.Collection("progress_result_raports")

	filter := bson.M{
		"program":  strings.ToUpper(program),
		"uka_type": strings.ToUpper(ukaType),
	}

	var results []models.ProgressResultRaport

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var result models.ProgressResultRaport
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetProgressRaport(smartbtwID []uint, ukaType string, stageType string, program string) ([]models.ProgressResultRaport, error) {
	collection := db.Mongodb.Collection("progress_result_raports")

	filter := bson.M{
		"smartbtw_id": bson.M{"$in": smartbtwID},
		"uka_type":    ukaType,
		"stage_type":  stageType,
		"program":     program,
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var results []models.ProgressResultRaport
	for cur.Next(context.TODO()) {
		var result models.ProgressResultRaport
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func StoreBuildProcess(payload mockstruct.GenerateProgressReportMessage) error {
	if os.Getenv("ENV") == "test" {
		return nil
	}
	ns, _ := sonic.Marshal(payload)
	return db.NewRedisCluster().RPush(context.Background(), db.RAPORT_REDIS_QUEUE_BUILDY_KEY, string(ns)).Err()
}

func StoreFailed(payload mockstruct.GenerateProgressReportMessage) error {
	if os.Getenv("ENV") == "test" {
		return nil
	}
	ns, _ := sonic.Marshal(payload)
	return db.NewRedisCluster().RPush(context.Background(), db.RAPORT_REDIS_QUEUE_FAILED_BUILD_KEY, string(ns)).Err()
}

func GetProgressRaportSingle(smartbtwID uint, ukaType string, stageType string, program string) (*models.ProgressResultRaport, error) {
	collection := db.Mongodb.Collection("progress_result_raports")

	filter := bson.M{
		"smartbtw_id": smartbtwID,
		"uka_type":    ukaType,
		"stage_type":  stageType,
		"program":     program,
	}

	var result models.ProgressResultRaport
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func CacheProgressReportToRedis(smartbtwID uint, ukaType string, stageType string, program string, link string) error {

	key := fmt.Sprintf("profile:progress-report:%s:%d:%s:%s", program, smartbtwID, ukaType, stageType)
	body := map[string]any{
		"program":    program,
		"uka_type":   ukaType,
		"stage_type": stageType,
		"link":       link,
	}

	errs := db.NewRedisCluster().HSet(context.Background(), key, body).Err()
	if errs != nil {
		return errs
	}
	if ukaType == "ALL_MODULE" && stageType == "KELAS" {
		err := SendProgressRaport(smartbtwID, link)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetCachedProgressReportFromRedis(smartbtwID uint, ukaType string, stageType string, program string) (map[string]string, error) {
	key := fmt.Sprintf("profile:progress-report:%s:%d:%s:%s", program, smartbtwID, ukaType, stageType)

	result, err := db.NewRedisCluster().HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

type HistoryGroup struct {
	Group  []map[string]any `json:"group"`
	IsLast bool             `json:"is_last"`
}

func BuildProgressRaport(req mockstruct.GenerateProgressReportMessage) error {

	type FormattedAbsent struct {
		Number int    `json:"number"`
		Date   string `json:"date"`
		Topic  string `json:"topic"`
		PIC    string `json:"pic"`
	}

	type ClassRooms struct {
		BranchCode  string   `json:"branch_code"`
		ClassCode   string   `json:"class_code"`
		Quota       int      `json:"quota"`
		QuotaFilled int      `json:"quota_filled"`
		Tags        []string `json:"tags"`
		Year        int      `json:"year"`
		Status      string   `json:"status"`
		IsOnline    bool     `json:"is_online"`
		Title       string   `json:"title"`
	}

	res := map[string]any{}
	switch strings.ToUpper(req.Program) {
	case "PTK":
		result, err := FetchPTKUKARaport(req.SmartbtwID, req.StageType, req.UKAType)
		if err != nil {
			return err
		}
		res = result
	case "PTN":
		result, err := FetchPTNUKARaport(req.SmartbtwID, req.StageType, req.UKAType)
		if err != nil {
			return err
		}
		res = result
	case "CPNS":
		result, err := FetchCPNSUKARaport(req.SmartbtwID, req.StageType, req.UKAType)
		if err != nil {
			return err
		}
		res = result
	}

	isIPDN := false
	if req.Program == "PTK" {
		studentData, _ := res["student"].(map[string]interface{})
		schoolName, _ := studentData["school_name"].(string)
		if schoolName == "IPDN" {
			isIPDN = true
		}
	}
	histories, ok := res["histories"].([]map[string]interface{})
	if !ok || len(histories) == 0 {
		return nil
	}

	his, _ := res["histories"].([]map[string]interface{})

	var first10History []map[string]interface{}
	var remainingHistory []map[string]interface{}

	if isIPDN {
		if len(his) > 8 {
			res["is_page_break"] = true
			first10History = his[:8]
			remainingHistory = his[8:]
		} else {
			first10History = his
		}
	} else {
		if len(his) > 10 {
			res["is_page_break"] = true
			first10History = his[:10]
			remainingHistory = his[10:]
		} else {
			first10History = his
		}
	}

	res["first_histoy"] = first10History

	transactionGroups := createHistoryGroups(remainingHistory)

	res["history_group"] = transactionGroups
	his = first10History

	presence, err := GetStudentPresence(req.SmartbtwID)
	if err != nil {
		return err
	}

	absent, _ := presence["absent_presences"].([]interface{})
	jsonData, err := json.Marshal(absent)
	if err != nil {
		fmt.Println("Error marshaling map to JSON:", err)
		return err
	}

	var abs []request.ClassSchedule
	err = json.Unmarshal(jsonData, &abs)
	if err != nil {
		fmt.Println("Error unmarshaling JSON to struct:", err)
		return err
	}
	var absents []FormattedAbsent

	if len(abs) > 0 {
		for j, a := range abs {
			pic, err := GetUserSSO(a.CreatedBy)
			if err != nil {
				return err
			}
			date := a.CreatedAt.Add(7 * time.Hour)
			targetFormat := "02/01/06 15:04"
			formattedTime := date.Format(targetFormat)
			absents = append(absents, FormattedAbsent{
				Number: j + 1,
				Date:   formattedTime,
				Topic:  a.ScheduleTopic,
				PIC:    pic.Name,
			})
		}

	}

	class, _ := presence["classrooms"].([]interface{})
	jsonDataClass, err := json.Marshal(class)
	if err != nil {
		fmt.Println("Error marshaling map to JSON:", err)
		return err
	}

	var classRooms []ClassRooms
	err = json.Unmarshal(jsonDataClass, &classRooms)
	if err != nil {
		fmt.Println("Error unmarshaling JSON to struct:", err)
		return err
	}

	programClass := "Reguler"
	for _, cl := range classRooms {
		for _, tag := range cl.Tags {
			if strings.Contains(strings.ToLower(tag), "binsus") {
				programClass = "Binsus"
			}

		}
	}

	lctime.SetLocale("id_ID")

	res["program"] = strings.ToUpper(req.Program)
	res["stage_type"] = req.StageType
	res["uka_type"] = req.UKAType
	res["smartbtw_id"] = req.SmartbtwID
	res["presence"] = presence
	res["absent_formatted"] = absents
	res["raport_date"] = lctime.Strftime("%A, %d %B %Y", time.Now())
	res["program_class"] = programClass

	msgBodys := map[string]any{
		"version": 1,
		"data":    res,
	}

	// srpJsonAns, err := sonic.Marshal(msgBodys)
	// if err != nil {
	// 	return errors.New("marshalling " + err.Error())
	// }
	// os.WriteFile(fmt.Sprintf("test_json_%s.json", strings.ToLower("ptn-progress")), srpJsonAns, 0644)

	msgJsons, errs := sonic.Marshal(msgBodys)
	if errs == nil && db.Broker != nil {
		_ = db.Broker.Publish(
			"progress-result-raport.build.delivery",
			"application/json",
			[]byte(msgJsons), // message to publish
		)
	}

	return nil
}

func createHistoryGroups(data []map[string]interface{}) []HistoryGroup {
	var HistoryGroups []HistoryGroup

	group := HistoryGroup{
		Group:  make([]map[string]interface{}, 0),
		IsLast: false,
	}

	for i, transaction := range data {
		if len(group.Group) == 10 || i == len(data)-1 {
			group.Group = append(group.Group, transaction)
			group.IsLast = i == len(data)-1

			HistoryGroups = append(HistoryGroups, group)

			group = HistoryGroup{
				Group:  make([]map[string]interface{}, 0),
				IsLast: false,
			}
		} else {
			group.Group = append(group.Group, transaction)
		}
	}

	return HistoryGroups
}

func GetProgressRaportList(smartbtwID []uint, ukaType string, stageType string, program string) ([]mockstruct.ProgressReport, error) {
	results := []mockstruct.ProgressReport{}

	for _, sm := range smartbtwID {
		redisRes, err := GetCachedProgressReportFromRedis(sm, ukaType, stageType, program)
		if err != nil {
			return nil, err
		}
		if len(redisRes) == 0 {
			dbRes, err := GetProgressRaportSingle(sm, ukaType, stageType, program)
			if err != nil {
				continue
			} else {
				results = append(results, mockstruct.ProgressReport{
					SmartbtwID: dbRes.SmartbtwID,
					Program:    dbRes.Program,
					UKAType:    dbRes.UKAType,
					StageType:  dbRes.StageType,
					Link:       dbRes.Link,
				})
			}
		} else {
			jsonData, err := json.Marshal(redisRes)
			if err != nil {
				fmt.Println("Error marshaling map to JSON:", err)
				return nil, err
			}

			var perkembangan mockstruct.ProgressReport
			err = json.Unmarshal(jsonData, &perkembangan)
			if err != nil {
				fmt.Println("Error unmarshaling JSON to struct:", err)
				return nil, err
			}
			perkembangan.SmartbtwID = int(sm)
			results = append(results, perkembangan)
		}

	}

	return results, nil
}

func SendProgressRaport(smartbtwID uint, link string) error {
	type SendRaportBody struct {
		To            string `json:"to"`
		Name          string `json:"name"`
		Greeting      string `json:"greeting"`
		CustomMessage string `json:"custom_message"`
		FileName      string `json:"file_name"`
		FileUrl       string `json:"file_url"`
	}

	parent, err := GetStudentBySmartBTWID(int(smartbtwID))
	if err != nil {
		return err
	}

	// greeting := helpers.GetGreeting()
	if len(parent) > 0 {
		if parent[0].ParentNumber != nil {
			payload := SendRaportBody{
				To:            *parent[0].ParentNumber,
				Name:          *parent[0].ParentName,
				CustomMessage: fmt.Sprintf("dokumen rapor perkembangan siswa atas nama %s", parent[0].Name),
				FileName:      fmt.Sprintf("raport_perkembangan_%d.pdf", time.Now().UnixMilli()),
				FileUrl:       link,
				Greeting:      helpers.GetGreeting(),
			}

			msgBodys := map[string]any{
				"version": 1,
				"data":    payload,
			}
			msgJsons, errs := sonic.Marshal(msgBodys)
			if errs == nil && db.Broker != nil {
				_ = db.Broker.Publish(
					"message-gateway.whatsapp.raport-result",
					"application/json",
					[]byte(msgJsons), // message to publish
				)
			}
		}
	}

	return nil
}

func RequestGenerateProgressRaport(smid []uint, program string, stgtype string) error {
	for _, sm := range smid {
		err := SendToGenerateRaport(program, sm, stgtype)
		if err != nil {
			return err
		}
	}
	return nil
}
