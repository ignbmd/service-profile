package lib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
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

func CreateHistoryCPNS(c *request.CreateHistoryCpns) (*mongo.InsertOneResult, error) {
	htpCol := db.Mongodb.Collection("history_cpns")
	stdCol := db.Mongodb.Collection("student_target_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "is_active": true, "deleted_at": nil}
	stdModels := models.StudentTargetCpns{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	payload := models.HistoryCpns{
		SmartBtwID:        c.SmartBtwID,
		TaskID:            c.TaskID,
		PackageID:         c.PackageID,
		PackageType:       c.PackageType,
		ModuleCode:        c.ModuleCode,
		ModuleType:        c.ModuleType,
		Twk:               c.Twk,
		Tiu:               c.Tiu,
		Tkp:               c.Tkp,
		TwkPass:           c.TwkPass,
		TiuPass:           c.TiuPass,
		TkpPass:           c.TkpPass,
		Total:             c.Total,
		Repeat:            c.Repeat,
		ExamName:          c.ExamName,
		Grade:             c.Grade,
		TargetID:          stdModels.ID,
		SchoolOriginID:    c.SchoolOriginID,
		SchoolOrigin:      c.SchoolOrigin,
		InstanceID:        c.InstanceID,
		InstanceName:      c.InstanceName,
		PositionID:        c.PositionID,
		PositionName:      c.PositionName,
		CompetitionCpnsID: c.CompetitionCpnsID,
		FormationType:     c.FormationType,
		FormationCode:     c.FormationCode,
		FormationLocation: c.FormationLocation,
		StudentName:       c.StudentName,
		TargetScore:       c.TargetScore,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	res, err := htpCol.InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateHistoryCPNS(c *request.UpdateHistoryCPNS, id primitive.ObjectID) error {
	opts := options.Update().SetUpsert(true)
	htkCol := db.Mongodb.Collection("history_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "deleted_at": nil}
	stdTarget := models.HistoryCpns{}
	err := htkCol.FindOne(ctx, filter).Decode(&stdTarget)
	if err != nil {
		return fmt.Errorf("data not found")
	}

	payload := models.HistoryCpns{
		SmartBtwID:        stdTarget.SmartBtwID,
		TaskID:            c.TaskID,
		PackageID:         c.PackageID,
		PackageType:       c.PackageType,
		ModuleCode:        c.ModuleCode,
		ModuleType:        c.ModuleType,
		Twk:               c.Twk,
		Tiu:               c.Tiu,
		Tkp:               c.Tkp,
		Total:             c.Total,
		Repeat:            c.Repeat,
		ExamName:          c.ExamName,
		Grade:             c.Grade,
		TargetID:          stdTarget.TargetID,
		CreatedAt:         stdTarget.CreatedAt,
		SchoolOriginID:    stdTarget.SchoolOriginID,
		SchoolOrigin:      stdTarget.SchoolOrigin,
		InstanceID:        stdTarget.InstanceID,
		InstanceName:      stdTarget.InstanceName,
		PositionID:        stdTarget.PositionID,
		PositionName:      stdTarget.PositionName,
		CompetitionCpnsID: stdTarget.CompetitionCpnsID,
		FormationType:     stdTarget.FormationType,
		FormationCode:     stdTarget.FormationCode,
		FormationLocation: stdTarget.FormationLocation,
		StudentName:       stdTarget.StudentName,
		TargetScore:       stdTarget.TargetScore,
		UpdatedAt:         time.Now(),
		DeletedAt:         nil,
	}

	update := bson.M{"$set": payload}
	_, err1 := htkCol.UpdateByID(ctx, stdTarget.ID, update, opts)
	if err1 != nil {
		return err1
	}
	return nil
}

func DeleteHistoryCPNS(id primitive.ObjectID) error {
	htkCol := db.Mongodb.Collection("history_cpns")
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

func GetHistoryCPNSByID(id primitive.ObjectID) ([]bson.M, error) {
	var results []bson.M

	collection := db.Mongodb.Collection("history_cpns")
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

func GetHistoryCPNSBySmartBTWID(SmartBTWID int, params *request.HistoryCPNSQueryParams) ([]bson.M, error) {
	var results []bson.M
	scCol := db.Mongodb.Collection("history_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if params.Limit != nil {
		if *params.Limit <= 0 {
			return nil, fmt.Errorf("limit must be a positive number")
		}
	}

	pipel := aggregates.GetStudentCPNSHistoryScores(SmartBTWID, params)
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

func GetHistoryFreeSingleStudentCPNS(smID int) (models.HistoryCpns, error) {
	scCol := db.Mongodb.Collection("history_cpns")
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
	stdModel := models.HistoryCpns{}
	err := scCol.FindOne(ctx, filter).Decode(&stdModel)
	if err != nil {
		return models.HistoryCpns{}, err
	}

	return stdModel, nil
}

func GetHistoryPremiumUKASingleStudentCPNS(smID int) (models.HistoryCpns, error) {
	scCol := db.Mongodb.Collection("history_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smID, "module_type": models.UkaPremium, "deleted_at": nil}
	stdModel := models.HistoryCpns{}
	err := scCol.FindOne(ctx, filter).Decode(&stdModel)
	if err != nil {
		return models.HistoryCpns{}, err
	}

	return stdModel, nil
}

func GetHistoryPackageUKASingleStudentCPNS(smID int) (models.HistoryCpns, error) {
	scCol := db.Mongodb.Collection("history_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": smID, "module_type": models.Package, "deleted_at": nil}
	stdModel := models.HistoryCpns{}
	err := scCol.FindOne(ctx, filter).Decode(&stdModel)
	if err != nil {
		return models.HistoryCpns{}, err
	}

	return stdModel, nil
}

func GetALLStudentScoreCPNSPagination(smID int, limit *int64, page *int64) ([]models.HistoryCpns, int64, error) {
	scCol := db.Mongodb.Collection("history_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"smartbtw_id": smID, "deleted_at": nil, "module_type": bson.M{"$nin": []string{"PRE_TEST", "POST_TEST"}}}
	var scrModel = make([]models.HistoryCpns, 0)

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
			return []models.HistoryCpns{}, totalData, err
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var model models.HistoryCpns
			e := cur.Decode(&model)
			if e != nil {
				log.Fatal(e)
			}
			scrModel = append(scrModel, model)
		}
	} else {
		cur, err := scCol.Find(ctx, fil, opts)
		if err != nil {
			return []models.HistoryCpns{}, totalData, err
		}

		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var model models.HistoryCpns
			e := cur.Decode(&model)
			if e != nil {
				log.Fatal(e)
			}
			scrModel = append(scrModel, model)
		}
	}

	return scrModel, totalData, nil

}

func GetStudentHistoryCPNSElastic(smID int, isStagesHistory bool) ([]request.CreateHistoryCpns, error) {
	ctx := context.Background()

	var t request.CreateHistoryCpns
	var gres []request.CreateHistoryCpns

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
		Index(db.GetStudentHistoryCpnsIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryCpns{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryCpns{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryCpns))
	}

	return gres, nil
}

func GetStudentHistoryCPNS(smID int, isStagesHistory bool) ([]request.HistoryCpnsElastic, error) {
	ctx := context.Background()

	var t request.HistoryCpnsElastic
	var gres []request.HistoryCpnsElastic

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if isStagesHistory {
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type", "pre-uka"))
		elasticQuery = append(elasticQuery, elastic.NewMatchQuery("package_type", "challenge-uka"))
	}

	query := elastic.NewBoolQuery().Must(elasticQuery...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryCpnsIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.HistoryCpnsElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.HistoryCpnsElastic{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.HistoryCpnsElastic))
	}

	return gres, nil
}

func GetStudentHistoryCPNSElasticFilter(smID int, filter string) ([]request.HistoryCpnsElastic, error) {
	ctx := context.Background()

	var t request.HistoryCpnsElastic
	var gres []request.HistoryCpnsElastic

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
		Index(db.GetStudentHistoryCpnsIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.HistoryCpnsElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.HistoryCpnsElastic{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.HistoryCpnsElastic))
	}

	return gres, nil
}

func InsertStudentCPNSProfileElastic(data *request.StudentProfileCPNSElastic, indexID string) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentTargetCpnsIndexName()).
		Id(indexID).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func UpsertHistoryCPNS(c *request.CreateHistoryCpns) (*string, error) {
	var upid string
	opts := options.Update().SetUpsert(true)
	htpCol := db.Mongodb.Collection("history_cpns")
	stdCol := db.Mongodb.Collection("student_target_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartBtwID, "is_active": true, "deleted_at": nil}
	stdModels := models.StudentTargetCpns{}
	htsModels := models.HistoryCpns{}
	filll := bson.M{"smartbtw_id": c.SmartBtwID, "task_id": c.TaskID}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return nil, fmt.Errorf("data not found")
	}

	update := bson.M{"$set": bson.M{
		"smartbtw_id":  c.SmartBtwID,
		"task_id":      c.TaskID,
		"package_id":   c.PackageID,
		"package_type": c.PackageType,
		"module_code":  c.ModuleCode,
		"module_type":  c.ModuleType,
		"twk":          c.Twk,
		"tiu":          c.Tiu,
		"tkp":          c.Tkp,
		"twk_pass":     c.TwkPass,
		"tiu_pass":     c.TiuPass,
		"tkp_pass":     c.TkpPass,
		"total":        c.Total,
		"repeat":       c.Repeat,
		"exam_name":    c.ExamName,
		"grade":        c.Grade,
		"start":        &c.Start,
		"end":          &c.End,
		"is_live":      c.IsLive,
		"target_id":    stdModels.ID,
		"student_name": c.StudentName,

		"target_score":       c.TargetScore,
		"school_origin_id":   c.SchoolOriginID,
		"school_origin":      c.SchoolOrigin,
		"instance_id":        c.InstanceID,
		"instance_name":      c.InstanceName,
		"position_id":        c.PositionID,
		"position_name":      c.PositionName,
		"formation_type":     c.FormationType,
		"formation_location": c.FormationLocation,
		"formation_code":     c.FormationCode,
		"competition_id":     c.CompetitionCpnsID,

		"created_at": time.Now(),
		"updated_at": time.Now(),
		"deleted_at": nil,
	}}

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
		switch reflect.TypeOf(res.UpsertedID).Kind() {
		case reflect.String:
			upid = fmt.Sprintf("%v", res.UpsertedID)
		default:
			upid = res.UpsertedID.(primitive.ObjectID).Hex()
		}
	}

	return &upid, nil
}

func InsertStudentHistoryCPNSElastic(data *request.CreateHistoryCpns, historyCPNSID string) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Update().
		Index(db.GetStudentHistoryCpnsIndexName()).
		Id(historyCPNSID).
		Doc(data).
		DocAsUpsert(true).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func UpsertHistoryCPNSCategoryTimeConsumed(data *request.UpdateHistoryCpnsTime) (*string, error) {
	var upid string
	opts := options.Update().SetUpsert(true)
	htpCol := db.Mongodb.Collection("history_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	htsModels := models.HistoryCpns{}
	filll := bson.M{"smartbtw_id": data.SmartBtwID, "task_id": data.TaskID}

	update := bson.M{"$set": bson.M{
		"twk_time_consumed": data.TwkTimeConsumed,
		"tiu_time_consumed": data.TiuTimeConsumed,
		"tkp_time_consumed": data.TkpTimeConsumed,
	}}

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
		switch reflect.TypeOf(res.UpsertedID).Kind() {
		case reflect.String:
			upid = fmt.Sprintf("%v", res.UpsertedID)
		default:
			upid = res.UpsertedID.(primitive.ObjectID).Hex()
		}
	}

	_ = UpdateCPNSTimeConsumed(data.SmartBtwID, data.TaskID, "SKD", data)

	return &upid, nil
}

func UpsyncStudentHistoryCPNSTimeConsumedElastic(data *request.UpdateHistoryCpnsTime, historyCPNSID string, retryCount int) error {
	ctx := context.Background()

	_, err := db.ElasticClient.Update().
		Index(db.GetStudentHistoryCpnsIndexName()).
		Id(historyCPNSID).
		Doc(data).
		DocAsUpsert(true).
		Do(ctx)

	if err != nil {
		if retryCount >= 5 {
			return err
		}
		if elastic.IsConflict(err) {
			time.Sleep(500 * time.Millisecond)
			return UpsyncStudentHistoryCPNSTimeConsumedElastic(data, historyCPNSID, retryCount+1)
		} else {
			return err
		}
	}

	return nil
}

func FetchStudentLearningRecordCPNS(smID uint) (*request.StudentLearningRecordHistory, error) {
	conn := os.Getenv("SERVICE_EXAM_CPNS_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		requestx *http.Request
		err      error
	)
	requestx, err = http.NewRequest("POST", conn+fmt.Sprintf("/explanation/learning-record/%d", smID), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to product " + err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(requestx)
	if err != nil {
		return nil, errors.New("doing request to product " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of product " + err.Error())
	}

	type responseBody struct {
		Message any `json:"message"`
		Data    *request.StudentLearningRecordHistory
	}

	st := responseBody{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of product " + errs.Error())
	}

	return st.Data, nil
}

func GetHistoryCpnsByTaskID(taskID uint) ([]request.GetHistoryCpnsResultElastic, error) {
	resData := []request.GetHistoryCpnsResultElastic{}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("task_id", taskID),
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetProfileCPNSExamResult()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var t request.GetHistoryCpnsResultElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		resData = append(resData, item.(request.GetHistoryCpnsResultElastic))
	}

	return resData, nil

}

func GetStudentProfileBySmartBtwID(smartbtwid int) ([]request.StudentProfileCPNSElastic, error) {
	resData := []request.StudentProfileCPNSElastic{}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smartbtwid),
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentTargetCpnsIndexName()).
		Query(query).
		Size(1000).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	var t request.StudentProfileCPNSElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		resData = append(resData, item.(request.StudentProfileCPNSElastic))
	}

	return resData, nil
}

func DateFormat(timePtr *time.Time) string {
	if timePtr == nil {
		return "Invalid Date"
	}
	day := timePtr.Day()
	month := getIndonesianMonthName(timePtr.Month())
	year := timePtr.Year()

	formattedDate := fmt.Sprintf("%d %s %d", day, month, year)

	return formattedDate
}

func getIndonesianMonthName(month time.Month) string {
	monthNames := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	// Mengurangi 1 dari nilai bulan karena index dimulai dari 0 pada package time
	indonesianMonth := monthNames[month-1]

	return indonesianMonth
}

func GetTotalWorkTime(start time.Time, end time.Time) string {
	duration := end.Sub(start)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	var date string

	if minutes == 0 {
		date = fmt.Sprintf("%d Detik", seconds)
	} else if hours == 0 {
		date = fmt.Sprintf("%d Menit, %d Detik", minutes, seconds)
	} else {
		date = fmt.Sprintf("%d Jam, %d Menit", hours, minutes)
	}

	return date
}

func GetRankCPNSByTaskID(task_id uint) ([]map[string]interface{}, error) {
	resCol, err := GetHistoryCpnsByTaskID(task_id)
	if err != nil {
		return nil, err
	}

	sort.Slice(resCol, func(i, j int) bool {
		return resCol[i].Total > resCol[j].Total
	})

	var pay []map[string]interface{}

	for i, v := range resCol {
		resEl, err := GetStudentProfileBySmartBtwID(v.SmartBtwID)
		if err != nil {
			return nil, err
		}

		date := DateFormat(v.Start)

		rank := i + 1

		for _, e := range resEl {

			duration := GetTotalWorkTime(*v.Start, *v.End)

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

			if v.Twk >= 65 && v.Tiu >= 80 && v.Tkp >= 156 {
				status = true
			} else {
				status = false
			}

			pay = append(pay, map[string]interface{}{
				"rankcpns": map[string]interface{}{
					"name":          v.Name,
					"exam_name":     v.Title,
					"task_id":       v.TaskID,
					"instance_name": e.InstanceName,
					"position_name": e.PositionName,
					"start":         fmt.Sprintf("%s WIB", v.Start.Format("15:04:05")),
					"end":           fmt.Sprintf("%s WIB", v.End.Format("15:04:05")),
					"duration":      duration,
					"twk":           v.Twk,
					"tiu":           v.Tiu,
					"tkp":           v.Tkp,
					"twk_status":    statusTwk,
					"tiu_status":    statusTiu,
					"tkp_status":    statusTkp,
					"date":          date,
					"rank":          rank,
					"total":         v.Total,
					"status":        status,
				},
			})

		}

	}

	return pay, nil
}

func GetHistoryCPNSElasticPeforma(smID uint, typStg, mdltype string) ([]request.CreateHistoryCpns, error) {
	ctx := context.Background()

	var t request.CreateHistoryCpns
	var gres []request.CreateHistoryCpns

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))

	if typStg == "UMUM" && mdltype == "pre-uka" {
		elasticQuery = append(elasticQuery, elastic.NewBoolQuery().
			Must(elastic.NewMatchQuery("module_type.keyword", "PLATINUM")).
			Must(elastic.NewMatchQuery("package_type.keyword", "pre-uka")))
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
		Index(db.GetStudentHistoryCpnsIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.CreateHistoryCpns{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.CreateHistoryCpns{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.CreateHistoryCpns))
	}

	return gres, nil
}

func GetStudentHistoryCPNSElasticFilterOld(smID int, filter string) ([]request.HistoryCpnsElastic, error) {
	ctx := context.Background()

	var t request.HistoryCpnsElastic
	var gres []request.HistoryCpnsElastic

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
		Index(db.GetStudentHistoryCpnsIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.HistoryCpnsElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.HistoryCpnsElastic{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.HistoryCpnsElastic))
	}

	return gres, nil
}

func GetHistoryCPNSElasticFetchStudent(smID uint, typStg, mdltype string) ([]request.HistoryCpnsElastic, error) {
	ctx := context.Background()

	var t request.HistoryCpnsElastic
	var gres []request.HistoryCpnsElastic

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
		Index(db.GetStudentHistoryCpnsIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return []request.HistoryCpnsElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []request.HistoryCpnsElastic{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(request.HistoryCpnsElastic))
	}

	return gres, nil
}

func FetchCPNSRankingSchoolPurposes(taskId uint, schoolId string, limit int, page int, keyword string) (mockstruct.FetchRankingCPNSBody, error) {
	ptkRankBody := mockstruct.FetchRankingCPNSBody{}

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
		Index("student_history_cpns").
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
		Index("student_history_cpns").
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

				smDataProf, err := GetStudentProfileCPNSElastic(t.SmartBtwID)
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

		ptkRankBody.RankingData = append(ptkRankBody.RankingData, mockstruct.FetchRankingCPNS{
			FetchRankingStudentBase: mockstruct.FetchRankingStudentBase{
				SmartBtwID:    smData.SmartbtwID,
				Email:         smData.Email,
				TaskID:        t.TaskID,
				PackageID:     t.PackageID,
				Name:          smData.Name,
				MajorID:       int(smData.PositionCPNSID),
				MajorName:     smData.PositionCPNSName,
				SchoolID:      int(smData.InstanceCPNSID),
				SchoolName:    smData.InstanceCPNSName,
				LastEdID:      smData.LastEdID,
				LastEdName:    smData.LastEdName,
				PassingChance: percATT,
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

func GetHistoryCPNSElastic(smID uint, taskID int) (request.CreateHistoryCpns, error) {
	ctx := context.Background()

	var t request.CreateHistoryCpns

	var elasticQuery []elastic.Query

	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("smartbtw_id", smID))
	elasticQuery = append(elasticQuery, elastic.NewMatchQuery("task_id", taskID))

	query := elastic.NewBoolQuery().Must(elasticQuery...)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentHistoryCpnsIndexName()).
		Query(query).
		Size(1). // Set the size to 1 to retrieve only one record
		Do(ctx)

	if err != nil {
		return request.CreateHistoryCpns{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.CreateHistoryCpns{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		return item.(request.CreateHistoryCpns), nil
	}

	// This line should not be reached, as we are returning within the loop
	return request.CreateHistoryCpns{}, nil
}

func GetHistoryCPNS(smartbtwID uint) ([]request.CreateHistoryCpns, error) {
	htkCol := db.Mongodb.Collection("history_cpns")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": smartbtwID,
		"module_type": bson.M{"$nin": []string{"TESTING", "WITH_CODE"}},
		"deleted_at":  nil,
	}

	cursor, err := htkCol.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var historyCPNSList []request.CreateHistoryCpns
	for cursor.Next(ctx) {
		var historyCPNS models.HistoryCpns
		if err := cursor.Decode(&historyCPNS); err != nil {
			return nil, err
		}

		bd := request.CreateHistoryCpns{
			SmartBtwID:     historyCPNS.SmartBtwID,
			TaskID:         historyCPNS.TaskID,
			PackageID:      historyCPNS.PackageID,
			ModuleCode:     historyCPNS.ModuleCode,
			ModuleType:     historyCPNS.ModuleType,
			PackageType:    historyCPNS.PackageType,
			Twk:            historyCPNS.Twk,
			Tiu:            historyCPNS.Tiu,
			Tkp:            historyCPNS.Tkp,
			TwkPass:        historyCPNS.TwkPass,
			TiuPass:        historyCPNS.TiuPass,
			TkpPass:        historyCPNS.TkpPass,
			Total:          historyCPNS.Total,
			Repeat:         historyCPNS.Repeat,
			ExamName:       historyCPNS.ExamName,
			Grade:          historyCPNS.Grade,
			StudentName:    historyCPNS.StudentName,
			SchoolOrigin:   historyCPNS.SchoolOrigin,
			SchoolOriginID: historyCPNS.SchoolOriginID,
			InstanceName:   historyCPNS.InstanceName,
			InstanceID:     historyCPNS.InstanceID,
			PositionName:   historyCPNS.PositionName,
			PositionID:     historyCPNS.PositionID,
			TargetScore:    historyCPNS.TargetScore,
		}
		historyCPNSList = append(historyCPNSList, bd)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return historyCPNSList, nil
}
