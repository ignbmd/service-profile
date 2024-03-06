package lib

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/request"
)

func GetHistoryScoreByTargetType(targetType string, params *request.HistoryScoreQueryParams) ([]bson.M, error) {
	var result []bson.M
	var historyCollection string

	if targetType != "ptk" && targetType != "ptn" && targetType != "cpns" {
		return result, errors.New("target type is not valid, must be either ptk, ptn or cpns")
	}

	if targetType == "ptk" {
		historyCollection = "history_ptk"
	} else if targetType == "ptn" {
		historyCollection = "history_ptn"
	} else {
		historyCollection = "history_cpns"
	}

	collection := db.Mongodb.Collection(historyCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var pipel []primitive.M
	if targetType == "cpns" {
		pipel = aggregates.GetHistoryScoreCPNS(params)
	} else {
		pipel = aggregates.GetHistoryScoreByTargetType(targetType, params)
	}
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []bson.M{}, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return []bson.M{}, err
	}

	return result, nil
}

func SendAssessmentCompleted(req *request.BodySendAssessmentCompleted) error {
	msgBody := request.SendAssessmentCompleted{
		Version: 1,
		Data:    *req,
	}

	msgJson, err := sonic.Marshal(msgBody)
	if err != nil {
		return errors.New("error on marshaling json student recommendation body " + err.Error())
	}
	if db.Broker == nil {
		return errors.New("rabbit mq not available " + err.Error())
	}

	LogEvent(
		"InsertScorePTNandPTN-SendAssesmentComplated",
		msgJson,
		"PUBLISH:stages.student-assessment.upsert",
		"publishing student assesment completed data",
		"INFO",
		fmt.Sprintf("profile-publish-event-%s", "stages.student-assessment.upsert"))
	// Attempt to publish a message to the queue.
	if err = db.Broker.Publish(
		"stages.student-assessment.upsert",
		"application/json",
		[]byte(msgJson), // message to publish
	); err != nil {
		return errors.New("error on publishing mq for student recommendation " + err.Error())
	}
	return nil
}

func SaveResultFirebase(c *request.FirebaseHistoryScores, program string, targetType string) error {

	if db.FbDB == nil {
		fmt.Println("firebase not initialized")
		return nil
	}

	refName := fmt.Sprintf("/student-results-%s", strings.ToLower(targetType))

	ref := db.FbDB.NewRef(refName)

	strct := map[string]interface{}{
		"exam_name":        c.ExamName,
		"grade":            c.Grade,
		"target_type":      c.TargetType,
		"target_score":     c.TargetScore,
		"target":           c.Target,
		"proficiency":      c.Proficiency,
		"exam_proficiency": c.ExamProficiency,
		"stages_done":      c.StagesDone,
	}

	if strings.ToLower(targetType) == "ptn" {
		strct["summary"] = c.SummaryPTN
		strct["result"] = c.ResultPTN
		strct["total"] = c.Total
	} else {
		strct["summary"] = c.Summary
		strct["result"] = c.Result
	}

	if strings.ToLower(targetType) == "cpns" {
		strct["instance_id"] = c.InstanceID
		strct["instance_name"] = c.InstanceName
		strct["position_id"] = c.PositionID
		strct["position_name"] = c.PositionName
		strct["target"] = c.TargetCPNS
		strct["total_stages"] = c.TotalStages
		strct["uka_passed"] = c.UKAPassed
	} else {
		strct["school_id"] = c.SchoolID
		strct["major_id"] = c.MajorID
		strct["school_name"] = c.SchoolName
		strct["major_name"] = c.MajorName
		strct["target"] = c.Target
	}

	err := ref.Child(fmt.Sprintf("%d_%d_%s", c.SmartbtwID, c.TaskID, strings.ToUpper(program))).Update(db.Ctx, strct)
	if err != nil {
		return err
	}

	ptnScores := []map[string]any{}

	err = ref.Child(fmt.Sprintf("%d_%d_%s/scores", c.SmartbtwID, c.TaskID, strings.ToUpper(program))).Get(db.Ctx, &ptnScores)
	if err != nil {
		return err
	}

	for _, val := range ptnScores {
		switch fmt.Sprintf("%v", val["category"]) {
		case "penalaran_umum":
			val["score"] = c.NewScorePTN.PenalaranUmum
		case "pengetahuan_umum":
			val["score"] = c.NewScorePTN.PengetahuanUmum
		case "pemahaman_bacaan":
			val["score"] = c.NewScorePTN.PemahamanBacaan
		case "pengetahuan_kuantitatif":
			val["score"] = c.NewScorePTN.PengetahuanKuantitatif
		case "literasi_bahasa_indonesia":
			val["score"] = c.NewScorePTN.LiterasiBahasaIndonesia
		case "literasi_bahasa_inggris":
			val["score"] = c.NewScorePTN.LiterasiBahasaInggris
		case "penalaran_matematika":
			val["score"] = c.NewScorePTN.PenalaranMatematika
		}
	}

	err = ref.Child(fmt.Sprintf("%d_%d_%s", c.SmartbtwID, c.TaskID, strings.ToUpper(program))).Update(db.Ctx, map[string]any{
		"scores": ptnScores,
	})

	if err != nil {
		return err
	}

	return nil
}

func UpdateCPNSTimeConsumed(smartBtwId int, taskId int, program string, data *request.UpdateHistoryCpnsTime) error {

	if db.FbDB == nil {
		fmt.Println("firebase not initialized")
		return nil
	}

	refName := "/student-results-cpns"

	ref := db.FbDB.NewRef(refName)

	err := ref.Child(fmt.Sprintf("%d_%d_%s/summary/score_values/TIU", smartBtwId, taskId, strings.ToUpper(program))).Update(db.Ctx, map[string]any{
		"category_attempt_time": data.TiuTimeConsumed,
	})

	if err != nil {
		return err
	}

	err = ref.Child(fmt.Sprintf("%d_%d_%s/summary/score_values/TWK", smartBtwId, taskId, strings.ToUpper(program))).Update(db.Ctx, map[string]any{
		"category_attempt_time": data.TwkTimeConsumed,
	})

	if err != nil {
		return err
	}

	err = ref.Child(fmt.Sprintf("%d_%d_%s/summary/score_values/TKP", smartBtwId, taskId, strings.ToUpper(program))).Update(db.Ctx, map[string]any{
		"category_attempt_time": data.TkpTimeConsumed,
	})

	if err != nil {
		return err
	}

	return nil
}

func GetUkaCodeScoresByEmail(email string) (map[string]any, error) {
	student, err := GetStudentProfileElasticByEmail(email)
	if err != nil {
		return nil, err
	}

	hisPTK, err := GetStudentHistoryPTKElasticSpecific(student.SmartbtwID, "with_code")
	if err != nil {
		return nil, err
	}

	hisPTN, err := GetStudentHistoryPTNElasticFilter(student.SmartbtwID, "with_code", "utbk")
	if err != nil {
		return nil, err
	}

	for i, ptn := range hisPTN {
		if ptn.StudentName == "" {
			hisPTN[i].StudentName = student.Name
		}
	}
	for i, ptk := range hisPTK {
		if ptk.StudentName == "" {
			hisPTK[i].StudentName = student.Name
		}
	}
	results := map[string]any{
		"student_email": student.Email,
		"ptn_histories": hisPTN,
		"ptk_histories": hisPTK,
	}

	return results, nil

}
