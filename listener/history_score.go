package listener

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func ListenScoreHistoryBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "history-ptk.created":
		if InsertScorePTK(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "history-ptn.created":
		if InsertScorePTN(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "history-cpns.created":
		if InsertScoreCPNS(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "history-assessment.created":
		if InsertScoreAssessment(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "history-cpns.time-consume.update":
		if UpdateCPNSTimeConsumed(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "profile.syncResult":
		if UpdateStudentDuration(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}

	return false
}

func InsertScorePTK(body []byte, key string) bool {
	lib.LogEvent(
		"InsertScorePTK",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data",
		"INFO",
		fmt.Sprintf("profile-%s", key))
	// TODO: Nyesuain Multiple Target
	msg, err := request.UnmarshalMessageHistoryPtkBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return true
	}

	result, err := govalidator.ValidateStruct(&msg)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msg.Version < 2 {
		text := fmt.Sprintf("[RabbitMQ][%s] Version data < 2", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	var isLive bool = false
	if msg.Data.IsLive {
		isLive = true
	}

	type schoolOriginData struct {
		SchoolOriginID string
		SchoolOrigin   string
	}

	sco := schoolOriginData{}
	// get student profile
	student, err := lib.GetStudentProfilePTKElastic(msg.Data.SmartBtwID)
	if err != nil {

		// get student profile
		studentRawData, err := lib.GetStudentBySmartBTWID(msg.Data.SmartBtwID)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Student data with id: %d not found on service profile", key, msg.Data.SmartBtwID)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
		if len(studentRawData) > 0 {
			if studentRawData[0].SchoolOriginID != nil && studentRawData[0].SchoolOrigin != nil {
				sco.SchoolOrigin = *studentRawData[0].SchoolOrigin
				sco.SchoolOriginID = *studentRawData[0].SchoolOriginID
			}
		}
		getDetailTargetData, err := lib.GetStudentTargetByCustom(msg.Data.SmartBtwID, "PTK")
		if err != nil {
			return true
		}

		var sPhoto string
		if studentRawData[0].Photo == nil {
			tempPhoto := ""
			sPhoto = tempPhoto
		} else {
			sPhoto = *studentRawData[0].Photo
		}

		student = request.StudentProfilePtkElastic{
			SmartbtwID:  studentRawData[0].SmartbtwID,
			Name:        studentRawData[0].Name,
			Photo:       sPhoto,
			SchoolID:    getDetailTargetData.SchoolID,
			SchoolName:  getDetailTargetData.SchoolName,
			MajorID:     getDetailTargetData.MajorID,
			MajorName:   getDetailTargetData.MajorName,
			TargetType:  getDetailTargetData.TargetType,
			TargetScore: getDetailTargetData.TargetScore,
			Proficiency: "BEGINNER",
		}
	} else {
		sts, err := lib.GetStudentProfileElastic(msg.Data.SmartBtwID)
		if err == nil {
			sco.SchoolOrigin = sts.LastEdName
			sco.SchoolOriginID = sts.LastEdID
		}
	}

	isBinsus := false
	isSpecificStage := false

	joinedClass, _ := lib.GetStudentJoinedClassType(msg.Data.SmartBtwID)

	for _, k := range joinedClass {
		if strings.Contains(strings.ToLower(k), "binsus") {
			isBinsus = true
			break
		}
	}

	studentProduct, err := lib.GetStudentAOP(uint(msg.Data.SmartBtwID))
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] error get activated product student with smartbtw_id : %d", key, msg.Data.SmartBtwID)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	totalModule := 0
	for _, prod := range studentProduct.Data {
		for _, tag := range prod.ProductTags {
			if strings.HasPrefix(strings.ToLower(tag), "ptk_stage_") {
				isSpecificStage = true
				totalModule += 1
			}
		}
	}
	totalModule = totalModule * 32 //TODO : pastikan jumlah module per stage

	var (
		percentScore    float64 = 0
		grade           string
		proficiency     string
		examProficiency string
	)

	if !msg.Data.AllPassStatus {
		percentScore = 0
	} else {
		percentScore = (msg.Data.Total / float64(student.TargetScore)) * 100
	}

	if percentScore >= 100 {
		percentScore = 99
	}

	if percentScore > 90 {
		grade = "DIAMOND"
	} else if percentScore > 75 && percentScore <= 90 {
		grade = "PLATINUM"
	} else if percentScore > 50 && percentScore <= 75 {
		grade = "GOLD"
	} else if percentScore > 25 && percentScore <= 50 {
		grade = "BRONZE"
	} else if percentScore <= 25 {
		grade = "NONE"
	} else {
		text := fmt.Sprintf("[RabbitMQ][%s] error searching percent with percent: %f", key, percentScore)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	insHis := request.CreateHistoryPtk{
		SmartBtwID:     msg.Data.SmartBtwID,
		TaskID:         msg.Data.TaskID,
		PackageID:      msg.Data.PackageID,
		PackageType:    msg.Data.PackageType,
		ModuleCode:     msg.Data.ModuleCode,
		ModuleType:     msg.Data.ModuleType,
		Twk:            msg.Data.Twk,
		Tiu:            msg.Data.Tiu,
		Tkp:            msg.Data.Tkp,
		TwkPass:        msg.Data.TwkPass,
		TiuPass:        msg.Data.TiuPass,
		TkpPass:        msg.Data.TkpPass,
		Total:          msg.Data.Total,
		Repeat:         msg.Data.Repeat,
		ExamName:       msg.Data.ExamName,
		Grade:          grade,
		Start:          msg.Data.Start,
		End:            msg.Data.End,
		IsLive:         isLive,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		SchoolOriginID: sco.SchoolOriginID,
		SchoolOrigin:   sco.SchoolOrigin,
		SchoolID:       student.SchoolID,
		SchoolName:     student.SchoolName,
		MajorID:        student.MajorID,
		MajorName:      student.MajorName,
		PolbitType:     student.PolbitType,
		StudentName:    student.Name,
		TiuPassStatus:  msg.Data.TiuPassStatus,
		TkpPassStatus:  msg.Data.TkpPassStatus,
		TwkPassStatus:  msg.Data.TwkPassStatus,
		AllPassStatus:  msg.Data.AllPassStatus,
		TargetScore:    student.TargetScore,
		// TargetID:   getDetailTarget.ID.Hex(),
	}

	if student.PolbitCompetitionID != nil {
		insHis.PolbitCompetitionID = *student.PolbitCompetitionID
	}

	if student.PolbitLocationID != nil {
		insHis.PolbitLocationID = *student.PolbitLocationID
	}
	//! Preparing data for update to elastic start from here
	// Get how much user doing module
	averages, err := lib.GetStudentHistoryPTKElastic(msg.Data.SmartBtwID, false)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when get average data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	countModule := len(averages) + 1

	isTaskExist := false
	for ns, k := range averages {
		if k.TaskID == insHis.TaskID {
			averages[ns].Twk = insHis.Twk
			averages[ns].TwkPass = insHis.TwkPass
			averages[ns].TwkPassStatus = insHis.TwkPassStatus
			averages[ns].Tiu = insHis.Tiu
			averages[ns].TiuPass = insHis.TiuPass
			averages[ns].TiuPassStatus = insHis.TiuPassStatus
			averages[ns].Tkp = insHis.Tkp
			averages[ns].TkpPass = insHis.TkpPass
			averages[ns].TkpPassStatus = insHis.TkpPassStatus
			averages[ns].AllPassStatus = insHis.AllPassStatus
			averages[ns].Total = insHis.Total
			averages[ns].Grade = insHis.Grade
			averages[ns].Repeat = insHis.Repeat
			isTaskExist = true
		}
	}

	if !isTaskExist {
		averages = append(averages, insHis)
	}

	stagesRecord := []request.CreateHistoryPtk{}
	challengeRecord := []request.CreateHistoryPtk{}
	stagesPreUKARecord := []request.CreateHistoryPtk{}
	specificStageRecord := []request.CreateHistoryPtk{}
	// stagesChallengeUKARecord := []request.CreateHistoryPtk{}

	for _, k := range averages {
		if strings.ToLower(k.PackageType) == "pre-uka" {
			stagesPreUKARecord = append(stagesPreUKARecord, k)
		}
		// if strings.ToLower(k.PackageType) == "challenge-uka" {
		// 	stagesChallengeUKARecord = append(stagesChallengeUKARecord, k)
		// }
		if strings.ToLower(k.PackageType) == "pre-uka" || strings.ToLower(k.PackageType) == "challenge-uka" {
			stagesRecord = append(stagesRecord, k)
		}
		if strings.ToLower(k.PackageType) == "challenge-uka" || strings.ToUpper(k.ModuleType) == "WITH_CODE" {
			challengeRecord = append(challengeRecord, k)
		}

		if strings.ToLower(k.PackageType) == "multi-stages-uka" {
			specificStageRecord = append(specificStageRecord, k)
		}

	}

	countStagesCompleted := len(stagesRecord)
	countPreStagesCompleted := len(stagesPreUKARecord)

	// Get student proficiency
	if countStagesCompleted >= 43 {
		proficiency = "EXPERT"
	} else if countStagesCompleted >= 15 && countStagesCompleted < 43 {
		proficiency = "ADVANCE"
	} else if countStagesCompleted >= 4 && countStagesCompleted < 15 {
		proficiency = "INTERMEDIATE"
	} else if countStagesCompleted < 4 {
		proficiency = "BEGINNER"
	}

	if countPreStagesCompleted >= 30 {
		examProficiency = "EXPERT"
	} else if countPreStagesCompleted >= 13 && countPreStagesCompleted < 31 {
		examProficiency = "ADVANCE"
	} else if countPreStagesCompleted >= 4 && countPreStagesCompleted < 13 {
		examProficiency = "INTERMEDIATE"
	} else if countPreStagesCompleted < 4 {
		examProficiency = "BEGINNER"
	}

	twkScore := float64(0)
	tiuScore := float64(0)
	tkpScore := float64(0)
	totalScore := float64(0)

	if isBinsus {
		for _, k := range challengeRecord {
			twkScore += k.Twk
			tiuScore += k.Tiu
			tkpScore += k.Tkp
			totalScore += k.Total
		}
	} else {
		for _, k := range averages {
			twkScore += k.Twk
			tiuScore += k.Tiu
			tkpScore += k.Tkp
			totalScore += k.Total
		}
	}

	passingTotalScore := float64(0)
	passingTotalItem := 0

	if !isBinsus {
		for _, k := range averages {
			if strings.ToLower(k.PackageType) != "pre-uka" {
				if (k.Tiu >= k.TiuPass) && (k.Twk >= k.TwkPass) && (k.Tkp >= k.TkpPass) {
					passingTotalItem += 1
					passingTotalScore += k.Total
				}
			}
		}
	} else {
		if len(challengeRecord) < 11 {
			passingTotalItem = 10
		} else {
			passingTotalItem = len(challengeRecord)
		}
	}

	atwk := float64(0)
	atiu := float64(0)
	atkp := float64(0)
	att := float64(0)
	pAtt := float64(0)

	if isBinsus {
		atwk = math.Round(helpers.RoundFloat((twkScore / float64(passingTotalItem)), 2))
		atiu = math.Round(helpers.RoundFloat((tiuScore / float64(passingTotalItem)), 2))
		atkp = math.Round(helpers.RoundFloat((tkpScore / float64(passingTotalItem)), 2))
		att = math.Round(helpers.RoundFloat((totalScore / float64(passingTotalItem)), 2))
		// att = atwk + atiu + atkp
		pAtt = att
	} else if isSpecificStage {
		atwk = helpers.RoundFloat((twkScore / float64(totalModule)), 2)
		atiu = helpers.RoundFloat((tiuScore / float64(totalModule)), 2)
		atkp = helpers.RoundFloat((tkpScore / float64(totalModule)), 2)
		att = helpers.RoundFloat((totalScore / float64(totalModule)), 2)
		pAtt = helpers.RoundFloat((passingTotalScore / float64(totalModule)), 2)
	} else {
		atwk = helpers.RoundFloat((twkScore / float64(len(averages))), 2)
		atiu = helpers.RoundFloat((tiuScore / float64(len(averages))), 2)
		atkp = helpers.RoundFloat((tkpScore / float64(len(averages))), 2)
		att = helpers.RoundFloat((totalScore / float64(len(averages))), 2)
		pAtt = helpers.RoundFloat((passingTotalScore / float64(passingTotalItem)), 2)
	}

	if math.IsNaN(pAtt) {
		pAtt = 0
	}

	percTWK := helpers.RoundFloat((atwk/150)*100, 2)
	percTIU := helpers.RoundFloat((atiu/175)*100, 2)
	percTKP := helpers.RoundFloat((atkp/225)*100, 2)
	percTT := helpers.RoundFloat((att/225)*100, 2)
	percATT := helpers.RoundFloat((pAtt/student.TargetScore)*100, 2)

	if percATT > 99 {
		percATT = 99
	}

	complatedStage := len(stagesRecord)
	if isSpecificStage {
		complatedStage += len(specificStageRecord)
	}

	stUpdate := request.StudentProfilePtkElastic{
		SmartbtwID:                           msg.Data.SmartBtwID,
		Name:                                 student.Name,
		Photo:                                student.Photo,
		ModuleDone:                           countModule,
		TwkAvgScore:                          atwk,
		TiuAvgScore:                          atiu,
		TkpAvgScore:                          atkp,
		TwkAvgPercentScore:                   percTWK,
		TiuAvgPercentScore:                   percTIU,
		TkpAvgPercentScore:                   percTKP,
		TotalAvgPercentScore:                 percTT,
		TotalAvgScore:                        att,
		LatestTotalScore:                     msg.Data.Total,
		LatestTotalPercentScore:              percentScore,
		SchoolID:                             student.SchoolID,
		SchoolName:                           student.SchoolName,
		MajorID:                              student.MajorID,
		MajorName:                            student.MajorName,
		TargetType:                           student.TargetType,
		PolbitType:                           student.PolbitType,
		PolbitCompetitionID:                  student.PolbitCompetitionID,
		PolbitLocationID:                     student.PolbitLocationID,
		TargetScore:                          float64(student.TargetScore),
		Proficiency:                          proficiency,
		ExamProficiency:                      examProficiency,
		StagesCompleted:                      complatedStage,
		PassingRecommendationAvgScore:        pAtt,
		PassingRecommendationAvgPercentScore: percATT,
	}
	// fmt.Println(stUpdate)

	// ! NOTE: insert lib update to elastic here
	err = lib.InsertStudentPtkProfileElastic(&stUpdate, fmt.Sprintf("%d_PTK", msg.Data.SmartBtwID))
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's PTK profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	school, err := lib.GetStudentSchoolData(uint(student.SchoolID), "PTK")
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when get school data from service comp map data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	summary := request.FirebaseHistoryScoresSummary{
		ScoreKeys: []string{"TWK", "TIU", "TKP"},
		ScoreValues: request.FirebaseHistoryScoresScoreValues{
			Twk: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        atwk,
				AvgPercentScore: percTWK,
			},
			Tiu: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        atiu,
				AvgPercentScore: percTIU,
			},
			Tkp: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        atkp,
				AvgPercentScore: percTKP,
			},
			Total: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        att,
				AvgPercentScore: percTT,
			},
		},
		LatestTotalScore:          msg.Data.Total,
		LatestTotalPercentScore:   percentScore,
		CurrentTargetTotalScore:   pAtt,
		CurrentTargetPercentScore: percATT,
	}

	historyResults := []request.ProfileStudentResultHistory{}
	for _, h := range averages {
		formattedTime := h.CreatedAt.Format("2006-01-02")
		res := request.ProfileStudentResultHistory{
			Date:  formattedTime,
			Title: h.ExamName,
			Total: int(h.Total),
		}
		historyResults = append(historyResults, res)
	}
	locID := 0
	if student.PolbitLocationID != nil {
		locID = *student.PolbitLocationID
	}

	studentData, err := lib.GetStudentProfileElastic(msg.Data.SmartBtwID)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when get profile elastic data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		// return false
	}
	majorData, err := lib.FetchMajorCompetitionData("ptk", uint(student.MajorID), uint(locID), studentData.Gender, student.PolbitType)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when fetch major chances comp map data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		// return false
	}

	chances := lib.GetMajorChances(majorData)

	target := request.FirebaseTarget{
		Location:            school.Location,
		MajorID:             uint(student.MajorID),
		MajorName:           student.MajorName,
		MaximumScore:        550,
		PolbitCompetitionID: student.PolbitCompetitionID,
		PolbitLocationID:    student.PolbitLocationID,
		PolbitType:          &student.PolbitType,
		SchoolID:            uint(student.SchoolID),
		SchoolName:          student.SchoolName,
		SchoolLogo:          school.Logo,
		TargetScore:         int(student.TargetScore),
		Type:                "PTK",
	}

	if majorData != nil {
		target.MajorChances = &chances
		target.MajorCompYear = &majorData.MajorYear
		target.MajorQuota = &majorData.MajorQuota
		target.MajorQuotaYear = &majorData.MajorQuotaYear
		target.MajorReqistrant = &majorData.MajorRegistered
	}

	_ = lib.SaveResultFirebase(&request.FirebaseHistoryScores{
		TaskID:          insHis.TaskID,
		SmartbtwID:      uint(insHis.SmartBtwID),
		ExamName:        insHis.ExamName,
		Grade:           insHis.Grade,
		TargetType:      student.TargetType,
		TargetScore:     student.TargetScore,
		SchoolID:        student.SchoolID,
		MajorID:         student.MajorID,
		SchoolName:      student.SchoolName,
		MajorName:       student.MajorName,
		Summary:         summary,
		Result:          historyResults,
		Target:          target,
		Proficiency:     proficiency,
		ExamProficiency: examProficiency,
		StagesDone:      complatedStage,
	}, "skd", "ptk")

	historyPTKID, err := lib.UpsertHistoryPtk(&insHis)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	err = lib.InsertStudentHistoryPtkElastic(&insHis, *historyPTKID)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's PTK score history to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Data.ModuleType == string(models.UkaPremium) || msg.Data.ModuleType == string(models.Package) {
		//insert poin to wallet history
		// w := request.CreateWalletHistory{
		// 	SmartbtwID:  msg.Data.SmartBtwID,
		// 	Description: fmt.Sprintf("Mengerjakan UKA %s (%s)", msg.Data.ExamName, msg.Data.ModuleCode),
		// }
		// _, er := lib.CreateWalletHistoryUKA(&w)
		// if er != nil {
		// 	text := fmt.Sprintf("[RabbitMQ][%s] Error when creating wallet history", key)
		// 	golog.Slack.ErrorWithData(text, body, er)
		// 	log.Println(er)
		// 	return false
		// }

		erReward := lib.SendExamReward(uint(msg.Data.SmartBtwID), msg.Data.Total, student.TargetScore, msg.Data.PackageType)
		if erReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d reward of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, erReward)
			return false
		}

		if msg.Data.IsLive {
			erReward := lib.UpdatePassingPercentage(uint(msg.Data.SmartBtwID), "ptk", uint(msg.Data.PackageID), percentScore, "stages-uka")
			if erReward != nil {
				text := fmt.Sprintf("[RabbitMQ][%s] Error when updating student %d passing percentage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
				golog.Slack.ErrorWithData(text, body, erReward)
			}
		}
		msgBody := request.BodySendAssessmentCompleted{
			PackageID:  uint(msg.Data.PackageID),
			SmartbtwID: uint(msg.Data.SmartBtwID),
			Grade:      grade,
			Score:      msg.Data.Total,
			TaskID:     uint(msg.Data.TaskID),
			Program:    "PTK",
		}
		errPub := lib.SendAssessmentCompleted(&msgBody)
		if errPub != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when publish student %d PTK assessment complete of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, errPub)
		}
	}

	if msg.Data.ModuleType == string(models.UkaCode) {
		erReward := lib.UpdatePassingPercentage(uint(msg.Data.SmartBtwID), "ptk", uint(msg.Data.PackageID), percentScore, "uka-code")
		if erReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when updating student %d passing percentage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, erReward)
		}
	}

	if msg.Data.ModuleType == string(models.UkaFree) {
		errReward := lib.SendUKAFreeReward(uint(msg.Data.SmartBtwID), "PTK")
		if errReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d reward changemajor of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, errReward)
			return false
		}
	}

	err = lib.SendHistoryStage(uint(msg.Data.SmartBtwID), uint(msg.Data.PackageID), "PTK")
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d history stage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] insert history PTK data for user_id %d and task_id %d success", key, msg.Data.SmartBtwID, msg.Data.TaskID)
	log.Println(text)
	lib.LogEvent(
		"InsertScorePTK",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data successfully",
		"INFO",
		fmt.Sprintf("profile-%s", key))

	if strings.ToUpper(msg.Data.ModuleType) != "PRE_TEST" || strings.ToUpper(msg.Data.ModuleType) != "POST_TEST" {
		msgBodys := map[string]any{
			"version": 1,
			"data":    insHis,
		}
		msgJsons, errs := sonic.Marshal(msgBodys)
		if errs != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when publish to generate pdf raport student %d task_id %d", key, msg.Data.SmartBtwID, msg.Data.TaskID)
			golog.Slack.ErrorWithData(text, body, err)
			return true
		}
		if errs == nil && db.Broker != nil {
			err = db.Broker.Publish(
				"raport-ptk.build",
				"application/json",
				[]byte(msgJsons), // message to publish
			)
		}
	}

	err = lib.SendToGenerateRaport("PTK", uint(msg.Data.SmartBtwID), msg.Data.PackageType)
	if err != nil {
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when generating pdf raport progress student %d task_id %d", key, msg.Data.SmartBtwID, msg.Data.TaskID)
			golog.Slack.ErrorWithData(text, body, err)
			return true
		}
	}
	return true
}

func InsertScorePTN(body []byte, key string) bool {
	// TODO: Nyesuain Multiple Target
	lib.LogEvent(
		"InsertScorePTN",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data",
		"INFO",
		fmt.Sprintf("profile-%s", key))
	msg, err := request.UnmarshalMessageHistoryPtnBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return true
	}

	result, err := govalidator.ValidateStruct(&msg)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	var isLive bool = false
	if msg.Data.IsLive {
		isLive = true
	}

	type schoolOriginData struct {
		SchoolOriginID string
		SchoolOrigin   string
	}

	sco := schoolOriginData{}
	// get student profile
	student, err := lib.GetStudentProfilePTNElastic(msg.Data.SmartBtwID)
	if err != nil {
		// get student profile
		studentRawData, err := lib.GetStudentBySmartBTWID(msg.Data.SmartBtwID)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Student data with id: %d not found on service profile", key, msg.Data.SmartBtwID)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
		if len(studentRawData) > 0 {
			if studentRawData[0].SchoolOriginID != nil && studentRawData[0].SchoolOrigin != nil {
				sco.SchoolOrigin = *studentRawData[0].SchoolOrigin
				sco.SchoolOriginID = *studentRawData[0].SchoolOriginID
			}
		}
		getDetailTargetData, err := lib.GetStudentTargetByCustom(msg.Data.SmartBtwID, "PTN")
		if err != nil {
			return true
		}

		var sPhoto string
		if studentRawData[0].Photo == nil {
			tempPhoto := ""
			sPhoto = tempPhoto
		} else {
			sPhoto = *studentRawData[0].Photo
		}

		student = request.StudentProfilePtnElastic{
			SmartbtwID:  studentRawData[0].SmartbtwID,
			Name:        studentRawData[0].Name,
			Photo:       sPhoto,
			SchoolID:    getDetailTargetData.SchoolID,
			SchoolName:  getDetailTargetData.SchoolName,
			MajorID:     getDetailTargetData.MajorID,
			MajorName:   getDetailTargetData.MajorName,
			TargetType:  getDetailTargetData.TargetType,
			TargetScore: getDetailTargetData.TargetScore,
			Proficiency: "BEGINNER",
		}
	} else {
		sts, err := lib.GetStudentProfileElastic(msg.Data.SmartBtwID)
		if err == nil {
			sco.SchoolOrigin = sts.LastEdName
			sco.SchoolOriginID = sts.LastEdID
		}
	}

	var (
		percentScore    float64 = 0
		grade           string
		proficiency     string
		examProficiency string
	)
	percentScore = (msg.Data.Total / student.TargetScore) * 100
	if percentScore >= 100 {
		percentScore = 99
	}

	if percentScore > 90 {
		grade = "DIAMOND"
	} else if percentScore > 75 && percentScore <= 90 {
		grade = "PLATINUM"
	} else if percentScore > 50 && percentScore <= 75 {
		grade = "GOLD"
	} else if percentScore > 25 && percentScore <= 50 {
		grade = "BRONZE"
	} else if percentScore <= 25 {
		grade = "NONE"
	} else {
		text := fmt.Sprintf("[RabbitMQ][%s] error searching percent with percent: %f", key, percentScore)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	progKey := "tps"
	if strings.ToLower(msg.Data.ProgramKey) != "" {
		progKey = strings.ToLower(msg.Data.ProgramKey)
	}

	insHis := request.CreateHistoryPtn{
		SmartBtwID:     msg.Data.SmartBtwID,
		TaskID:         msg.Data.TaskID,
		PackageID:      msg.Data.PackageID,
		PackageType:    msg.Data.PackageType,
		ModuleCode:     msg.Data.ModuleCode,
		ModuleType:     msg.Data.ModuleType,
		Total:          msg.Data.Total,
		Repeat:         msg.Data.Repeat,
		ProgramKey:     progKey,
		ExamName:       msg.Data.ExamName,
		Grade:          grade,
		Start:          msg.Data.Start,
		End:            msg.Data.End,
		IsLive:         isLive,
		SchoolOriginID: sco.SchoolOriginID,
		SchoolOrigin:   sco.SchoolOrigin,
		SchoolID:       student.SchoolID,
		SchoolName:     student.SchoolName,
		MajorID:        student.MajorID,
		MajorName:      student.MajorName,
		StudentName:    student.Name,
		TargetScore:    student.TargetScore,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if strings.ToLower(msg.Data.ProgramKey) == "utbk" {
		insHis.PengetahuanKuantitatif = msg.Data.PengetahuanKuantitatif
		insHis.PenalaranUmum = msg.Data.PenalaranUmum
		insHis.PengetahuanUmum = msg.Data.PengetahuanUmum
		insHis.PemahamanBacaan = msg.Data.PemahamanBacaan
	} else {
		insHis.PotensiKognitif = msg.Data.PotensiKognitif
	}

	insHis.PenalaranMatematika = msg.Data.PenalaranMatematika
	insHis.LiterasiBahasaIndonesia = msg.Data.LiterasiBahasaIndonesia
	insHis.LiterasiBahasaInggris = msg.Data.LiterasiBahasaInggris
	//! Preparing data for update to elastic start from here

	// Get Current Percent Score for drill progress

	averages, err := lib.GetStudentHistoryPTNElastic(msg.Data.SmartBtwID, false, msg.Data.ProgramKey)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when get average data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	countModule := len(averages) + 1

	isTaskExist := false
	for ns, k := range averages {
		if k.TaskID == insHis.TaskID {
			if strings.ToLower(msg.Data.ProgramKey) == "utbk" {
				averages[ns].PengetahuanKuantitatif = insHis.PengetahuanKuantitatif
				averages[ns].PenalaranUmum = insHis.PenalaranUmum
				averages[ns].PengetahuanUmum = insHis.PengetahuanUmum
				averages[ns].PemahamanBacaan = insHis.PemahamanBacaan
			} else {
				averages[ns].PotensiKognitif = insHis.PotensiKognitif
			}
			averages[ns].PenalaranMatematika = insHis.PenalaranMatematika
			averages[ns].LiterasiBahasaIndonesia = insHis.LiterasiBahasaIndonesia
			averages[ns].LiterasiBahasaInggris = insHis.LiterasiBahasaInggris
			averages[ns].Total = insHis.Total
			averages[ns].Grade = insHis.Grade
			averages[ns].Repeat = insHis.Repeat
			isTaskExist = true
		}
	}

	if !isTaskExist {
		averages = append(averages, insHis)
	}

	isSpecificStage := false
	totalModule := 0

	activatedProduct, err := lib.GetStudentAOP(uint(msg.Data.SmartBtwID))
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] error get activated product student with smartbtw_id : %d", key, msg.Data.SmartBtwID)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	for _, prod := range activatedProduct.Data {
		for _, tag := range prod.ProductTags {
			if strings.HasPrefix(strings.ToLower(tag), "ptn_stage_") {
				isSpecificStage = true
				totalModule += 1
			}
		}
	}
	totalModule = totalModule * 32 //TODO : pastikan jumlah module per stage

	stagesRecord := []request.CreateHistoryPtn{}
	stagesPreUKARecord := []request.CreateHistoryPtn{}
	stagesSpecificRecord := []request.CreateHistoryPtn{}
	// stagesChallengeUKARecord := []request.CreateHistoryPtn{}

	for _, k := range averages {
		if strings.ToLower(k.PackageType) == "pre-uka" {
			stagesPreUKARecord = append(stagesPreUKARecord, k)
		}
		// if strings.ToLower(k.PackageType) == "challenge-uka" {
		// 	stagesChallengeUKARecord = append(stagesChallengeUKARecord, k)
		// }
		if strings.ToLower(k.PackageType) == "pre-uka" || strings.ToLower(k.PackageType) == "challenge-uka" {
			stagesRecord = append(stagesRecord, k)
		}

		if strings.Contains("multi-stages-uka", k.PackageType) {
			stagesSpecificRecord = append(stagesSpecificRecord, k)
		}
	}

	countStagesCompleted := len(stagesRecord)

	if isSpecificStage {
		countStagesCompleted += len(stagesSpecificRecord)
	}

	countPreStagesCompleted := len(stagesPreUKARecord)

	// Get student proficiency
	if countStagesCompleted >= 43 {
		proficiency = "EXPERT"
	} else if countStagesCompleted >= 15 && countStagesCompleted < 43 {
		proficiency = "ADVANCE"
	} else if countStagesCompleted >= 4 && countStagesCompleted < 15 {
		proficiency = "INTERMEDIATE"
	} else if countStagesCompleted < 4 {
		proficiency = "BEGINNER"
	}

	if countPreStagesCompleted >= 30 {
		examProficiency = "EXPERT"
	} else if countPreStagesCompleted >= 13 && countPreStagesCompleted < 31 {
		examProficiency = "ADVANCE"
	} else if countPreStagesCompleted >= 4 && countPreStagesCompleted < 13 {
		examProficiency = "INTERMEDIATE"
	} else if countPreStagesCompleted < 4 {
		examProficiency = "BEGINNER"
	}

	pkScore := float64(0)
	puScore := float64(0)
	ppuScore := float64(0)
	pbmScore := float64(0)
	pmScore := float64(0)
	lbIndScore := float64(0)
	lbIngScore := float64(0)
	totalScore := float64(0)

	for _, k := range averages {
		if strings.ToLower(msg.Data.ProgramKey) == "utbk" {
			pkScore += k.PengetahuanKuantitatif
			puScore += k.PenalaranUmum
			ppuScore += k.PengetahuanUmum
			pbmScore += k.PemahamanBacaan
		} else {
			pkScore += k.PotensiKognitif
		}
		pmScore += k.PenalaranMatematika
		lbIndScore += k.LiterasiBahasaIndonesia
		lbIngScore += k.LiterasiBahasaInggris
		totalScore += k.Total
	}

	passingTotalScore := float64(0)

	for _, k := range averages {
		if strings.ToLower(k.PackageType) != "pre-uka" {
			passingTotalScore += k.Total
		}
	}
	apu := float64(0)
	appu := float64(0)
	apbm := float64(0)
	apk := float64(0)
	apm := float64(0)
	abi := float64(0)
	abing := float64(0)
	att := float64(0)
	pAtt := float64(0)

	if isSpecificStage {
		if strings.ToLower(msg.Data.ProgramKey) == "utbk" {
			apu = helpers.RoundFloat((puScore / float64(totalModule)), 3)
			appu = helpers.RoundFloat((ppuScore / float64(totalModule)), 3)
			apbm = helpers.RoundFloat((pbmScore / float64(totalModule)), 3)
		}
		apk = helpers.RoundFloat((pkScore / float64(totalModule)), 3)
		apm = helpers.RoundFloat((pmScore / float64(totalModule)), 3)
		abi = helpers.RoundFloat((lbIndScore / float64(totalModule)), 3)
		abing = helpers.RoundFloat((lbIngScore / float64(totalModule)), 3)
		att = helpers.RoundFloat((totalScore / float64(totalModule)), 3)
		pAtt = helpers.RoundFloat((passingTotalScore / float64(totalModule)), 3)
	} else {
		if strings.ToLower(msg.Data.ProgramKey) == "utbk" {
			apu = helpers.RoundFloat((puScore / float64(len(averages))), 3)
			appu = helpers.RoundFloat((ppuScore / float64(len(averages))), 3)
			apbm = helpers.RoundFloat((pbmScore / float64(len(averages))), 3)
		}
		apk = helpers.RoundFloat((pkScore / float64(len(averages))), 3)
		apm = helpers.RoundFloat((pmScore / float64(len(averages))), 3)
		abi = helpers.RoundFloat((lbIndScore / float64(len(averages))), 3)
		abing = helpers.RoundFloat((lbIngScore / float64(len(averages))), 3)
		att = helpers.RoundFloat((totalScore / float64(len(averages))), 3)
		pAtt = helpers.RoundFloat((passingTotalScore / float64(len(averages)-len(stagesPreUKARecord))), 3)
	}

	if math.IsNaN(pAtt) {
		pAtt = 0
	}

	percAPU := float64(0)
	percAPPU := float64(0)
	percAPBM := float64(0)
	if strings.ToLower(msg.Data.ProgramKey) == "utbk" {
		percAPU = helpers.RoundFloat((apu/1000)*100, 3)
		percAPPU = helpers.RoundFloat((appu/1000)*100, 3)
		percAPBM = helpers.RoundFloat((apbm/1000)*100, 3)
	}
	percAPK := helpers.RoundFloat((apk/1000)*100, 3)
	percAPM := helpers.RoundFloat((apm/1000)*100, 3)
	percABI := helpers.RoundFloat((abi/1000)*100, 3)
	percABING := helpers.RoundFloat((abing/1000)*100, 3)
	percTT := helpers.RoundFloat((att/1000)*100, 3)
	percATT := helpers.RoundFloat((pAtt/student.TargetScore)*100, 3)

	if percATT > 99 {
		percATT = 99
	}
	bdy := request.StudentProfilePtnElastic{
		SmartbtwID:                           msg.Data.SmartBtwID,
		Name:                                 student.Name,
		Photo:                                student.Photo,
		SchoolID:                             student.SchoolID,
		SchoolName:                           student.SchoolName,
		MajorID:                              student.MajorID,
		MajorName:                            student.MajorName,
		TargetType:                           student.TargetType,
		TargetScore:                          float64(student.TargetScore),
		ModuleDone:                           countModule,
		PkAvgScore:                           float64(apk),
		PmAvgScore:                           float64(apm),
		LbindAvgScore:                        float64(abi),
		LbingAvgScore:                        float64(abing),
		TotalAvgScore:                        float64(att),
		PkAvgPercentScore:                    float64(percAPK),
		PmAvgPercentScore:                    float64(percAPM),
		LbindAvgPercentScore:                 float64(percABI),
		LbingAvgPercentScore:                 float64(percABING),
		TotalAvgPercentScore:                 float64(percTT),
		LatestTotalScore:                     msg.Data.Total,
		LatestTotalPercentScore:              helpers.RoundFloat(percentScore, 3),
		Proficiency:                          proficiency,
		ExamProficiency:                      examProficiency,
		ProgramKey:                           msg.Data.ProgramKey,
		StagesCompleted:                      countStagesCompleted,
		PassingRecommendationAvgScore:        pAtt,
		PassingRecommendationAvgPercentScore: percATT,
	}

	if strings.ToLower(msg.Data.ProgramKey) == "utbk" {
		bdy.PuAvgScore = float64(apu)
		bdy.PpuAvgScore = float64(appu)
		bdy.PbmAvgScore = float64(apbm)
		bdy.PuAvgPercentScore = float64(percAPU)
		bdy.PpuAvgPercentScore = float64(percAPPU)
		bdy.PbmAvgPercentScore = float64(percAPBM)
	}
	// fmt.Println(bdy)

	//! NOTE: insert lib update to elastic here
	err = lib.InsertStudentPtnProfileElastic(&bdy, fmt.Sprintf("%d_PTN", msg.Data.SmartBtwID))
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's PTN profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	summary := request.FirebaseHistoryScoresSummaryPTN{
		ScoreKeys: []string{"PU", "PPU", "PBM", "PK", "LBIND", "LBING", "PM"},
		ScoreValues: request.FirebaseHistoryScoresScoreValuesPTN{
			PU: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        apu,
				AvgPercentScore: percAPU,
			},
			PPU: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        appu,
				AvgPercentScore: percAPPU,
			},
			PBM: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        apbm,
				AvgPercentScore: percAPBM,
			},
			PK: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        apk,
				AvgPercentScore: percAPK,
			},
			LBING: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        abing,
				AvgPercentScore: percABING,
			},
			LBIN: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        abi,
				AvgPercentScore: percABI,
			},
			PM: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        apm,
				AvgPercentScore: percAPM,
			},
		},
		LatestTotalScore:          msg.Data.Total,
		LatestTotalPercentScore:   percentScore,
		CurrentTargetTotalScore:   pAtt,
		CurrentTargetPercentScore: percATT,
	}

	school, err := lib.GetStudentSchoolData(uint(student.SchoolID), "PTN")
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when get school data from service comp map data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	historyResults := []request.ProfileStudentResultHistoryPTN{}
	for _, h := range averages {
		formattedTime := h.CreatedAt.Format("2006-01-02")
		res := request.ProfileStudentResultHistoryPTN{
			Date:  formattedTime,
			Title: h.ExamName,
			Total: h.Total,
		}
		historyResults = append(historyResults, res)
	}

	majorData, err := lib.FetchMajorCompetitionData("ptn", uint(student.MajorID), 0, "", "")
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when fetch major chances comp map data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
	}

	chances := lib.GetMajorChances(majorData)

	target := request.FirebaseTarget{
		Location:            nil,
		MajorID:             uint(student.MajorID),
		MajorName:           student.MajorName,
		MaximumScore:        1000,
		PolbitCompetitionID: nil,
		PolbitLocationID:    nil,
		PolbitType:          nil,
		SchoolID:            uint(student.SchoolID),
		SchoolName:          student.SchoolName,
		SchoolLogo:          school.Logo,
		TargetScore:         int(student.TargetScore),
		Type:                "PTN",
	}

	if majorData != nil {
		target.MajorChances = &chances
		target.MajorCompYear = &majorData.MajorYear
		target.MajorQuota = &majorData.MajorQuota
		target.MajorQuotaYear = &majorData.MajorQuotaYear
		target.MajorReqistrant = &majorData.MajorRegistered
	}

	_ = lib.SaveResultFirebase(&request.FirebaseHistoryScores{
		TaskID:          insHis.TaskID,
		SmartbtwID:      uint(insHis.SmartBtwID),
		ExamName:        insHis.ExamName,
		Grade:           insHis.Grade,
		TargetType:      student.TargetType,
		TargetScore:     student.TargetScore,
		SchoolID:        student.SchoolID,
		MajorID:         student.MajorID,
		SchoolName:      student.SchoolName,
		MajorName:       student.MajorName,
		SummaryPTN:      summary,
		ResultPTN:       historyResults,
		Target:          target,
		Proficiency:     proficiency,
		ExamProficiency: examProficiency,
		StagesDone:      countStagesCompleted,
		Total:           msg.Data.Total,
		NewScorePTN:     msg.Data,
	}, "utbk", "ptn")

	historyPTNID, err := lib.UpsertHistoryPtn(&insHis)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	err = lib.InsertStudentHistoryPtnElastic(&insHis, *historyPTNID)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's PTN history score to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Data.ModuleType == string(models.UkaPremium) || msg.Data.ModuleType == string(models.Package) {
		//insert poin to wallet history
		// w := request.CreateWalletHistory{
		// 	SmartbtwID:  msg.Data.SmartBtwID,
		// 	Description: fmt.Sprintf("Mengerjakan UKA %s (%s)", msg.Data.ExamName, msg.Data.ModuleCode),
		// }
		// _, er := lib.CreateWalletHistoryUKA(&w)
		// if er != nil {
		// 	text := fmt.Sprintf("[RabbitMQ][%s] Error when creating wallet history", key)
		// 	golog.Slack.ErrorWithData(text, body, er)
		// 	log.Println(er)
		// 	return false
		// }
		if !isLive {
			erReward := lib.SendExamReward(uint(msg.Data.SmartBtwID), msg.Data.Total, student.TargetScore, msg.Data.PackageType)
			if erReward != nil {
				text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d reward of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
				golog.Slack.ErrorWithData(text, body, erReward)
				return false
			}
		}
		if msg.Data.IsLive && strings.ToLower(msg.Data.ProgramKey) == "utbk" {
			erReward := lib.UpdatePassingPercentage(uint(msg.Data.SmartBtwID), "utbk", uint(msg.Data.PackageID), percentScore, "stages-uka")
			if erReward != nil {
				text := fmt.Sprintf("[RabbitMQ][%s] Error when updating student %d passing percentage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
				golog.Slack.ErrorWithData(text, body, erReward)
			}

		}
		msgBody := request.BodySendAssessmentCompleted{
			PackageID:  uint(msg.Data.PackageID),
			SmartbtwID: uint(msg.Data.SmartBtwID),
			Grade:      grade,
			Score:      msg.Data.Total,
			TaskID:     uint(msg.Data.TaskID),
			Program:    "PTN",
		}
		errPub := lib.SendAssessmentCompleted(&msgBody)
		if errPub != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when publish student %d PTN assessment complete of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, errPub)
		}

	}

	if msg.Data.ModuleType == string(models.UkaCode) {
		erReward := lib.UpdatePassingPercentage(uint(msg.Data.SmartBtwID), "utbk", uint(msg.Data.PackageID), percentScore, "uka-code")
		if erReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when updating student %d passing percentage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, erReward)
		}
	}

	if msg.Data.ModuleType == string(models.UkaFree) {
		errReward := lib.SendUKAFreeReward(uint(msg.Data.SmartBtwID), "PTN")
		if errReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d reward changemajor of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, errReward)
			return false
		}
	}

	err = lib.SendHistoryStage(uint(msg.Data.SmartBtwID), uint(msg.Data.PackageID), "PTN")
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d history stage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] insert history PTN data for user_id %d and task_id %d success", key, msg.Data.SmartBtwID, msg.Data.TaskID)
	log.Println(text)

	lib.LogEvent(
		"InsertScorePTN",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data successfully",
		"INFO",
		fmt.Sprintf("profile-%s", key))

	if strings.ToUpper(msg.Data.ModuleType) != "PRE_TEST" || strings.ToUpper(msg.Data.ModuleType) != "POST_TEST" {
		msgBodys := map[string]any{
			"version": 1,
			"data":    insHis,
		}
		msgJsons, errs := sonic.Marshal(msgBodys)
		if errs != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when publish to generate pdf raport student %d task_id %d", key, msg.Data.SmartBtwID, msg.Data.TaskID)
			golog.Slack.ErrorWithData(text, body, err)
			return true
		}
		if errs == nil && db.Broker != nil {
			err = db.Broker.Publish(
				"raport-ptn.build",
				"application/json",
				[]byte(msgJsons), // message to publish
			)
		}
	}

	err = lib.SendToGenerateRaport("PTN", uint(msg.Data.SmartBtwID), msg.Data.PackageType)
	if err != nil {
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when generating pdf raport progress student %d task_id %d", key, msg.Data.SmartBtwID, msg.Data.TaskID)
			golog.Slack.ErrorWithData(text, body, err)
			return true
		}
	}
	return true
}

func UpdateStudentDuration(body []byte, key string) bool {
	msg, err := request.UnmarshalUpdateDurationBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	result, err := govalidator.ValidateStruct(&msg)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msg.Version < 1 {
		return false
	}
	if msg.Data.Program == "skd" {
		err := lib.UpdateDurationtPTK(&msg.Data)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when updating data to db", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
		er := lib.UpdateHistoryPtkElastic(&msg.Data)
		if er != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when updating data to elastic", key)
			golog.Slack.ErrorWithData(text, body, er)
			log.Println(er)
			return false
		}
	} else {
		err := lib.UpdateDurationtPTN(&msg.Data)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when updating data to db", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
		er := lib.UpdateHistoryPtnElastic(&msg.Data)
		if er != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when updating data to elastic", key)
			golog.Slack.ErrorWithData(text, body, er)
			log.Println(er)
			return false
		}
	}
	return true
}

func UpdateCPNSTimeConsumed(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageHistoryCpnsTimeConsumedBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return true
	}

	if msg.Version < 1 {
		return true
	}
	fmt.Printf("[RabbitMQ][%s] Start updating Time Consumed for %d and task ID %d\n", key, msg.Data.SmartBtwID, msg.Data.TaskID)
	res, err := lib.UpsertHistoryCPNSCategoryTimeConsumed(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when updating data to db", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if res == nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when updating data to db, res ID not found for %d (%d)", key, msg.Data.TaskID, msg.Data.SmartBtwID)
		golog.Slack.Error(text, nil)
		return false
	}

	errR := lib.UpsyncStudentHistoryCPNSTimeConsumedElastic(&msg.Data, *res, 0)
	if errR != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when updating cache data", key)
		golog.Slack.ErrorWithData(text, body, errR)
		return false
	}

	fmt.Printf("[RabbitMQ][%s] Success updating Time Consumed for %d and task ID %d\n", key, msg.Data.SmartBtwID, msg.Data.TaskID)
	return true
}

func InsertScoreCPNS(body []byte, key string) bool {
	lib.LogEvent(
		"InsertScoreCPNS",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data",
		"INFO",
		fmt.Sprintf("profile-%s", key))
	// TODO: Nyesuain Multiple Target
	msg, err := request.UnmarshalMessageHistoryCpnsBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return true
	}

	result, err := govalidator.ValidateStruct(&msg)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msg.Version < 2 {
		text := fmt.Sprintf("[RabbitMQ][%s] Version data < 2", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	var isLive bool = false
	if msg.Data.IsLive {
		isLive = true
	}

	type schoolOriginData struct {
		SchoolOriginID string
		SchoolOrigin   string
	}

	sco := schoolOriginData{}
	student, err := lib.GetStudentProfileCPNSElastic(msg.Data.SmartBtwID)
	if err != nil {

		studentRawData, err := lib.GetStudentBySmartBTWID(msg.Data.SmartBtwID)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Student data with id: %d not found on service profile", key, msg.Data.SmartBtwID)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

		if len(studentRawData) > 0 {
			if studentRawData[0].SchoolOriginID != nil && studentRawData[0].SchoolOrigin != nil {
				sco.SchoolOrigin = *studentRawData[0].SchoolOrigin
				sco.SchoolOriginID = *studentRawData[0].SchoolOriginID
			}
		}
		getDetailTargetData, err := lib.GetStudentTargetCPNS(msg.Data.SmartBtwID)
		if err != nil {
			log.Println(err)
			if msg.Data.SmartBtwID < 172000 {
				log.Println("Skipping user ID ", msg.Data.SmartBtwID)
				return true
			}
			// text := fmt.Sprintf("[RabbitMQ][%s] Student data with id: %d target not found", key, msg.Data.SmartBtwID)
			// golog.Slack.ErrorWithData(text, body, err)
			return false
		}

		var sPhoto string
		if studentRawData[0].Photo == nil {
			tempPhoto := ""
			sPhoto = tempPhoto
		} else {
			sPhoto = *studentRawData[0].Photo
		}

		student = request.StudentProfileCPNSElastic{
			SmartbtwID:        studentRawData[0].SmartbtwID,
			Name:              studentRawData[0].Name,
			Photo:             sPhoto,
			InstanceID:        getDetailTargetData.InstanceID,
			InstanceName:      getDetailTargetData.InstanceName,
			PositionID:        getDetailTargetData.PositionID,
			PositionName:      getDetailTargetData.PositionName,
			TargetType:        getDetailTargetData.TargetType,
			TargetScore:       getDetailTargetData.TargetScore,
			FormationType:     getDetailTargetData.FormationType,
			CompetitionCpnsID: uint(getDetailTargetData.CompetitionID),
			FormationCode:     getDetailTargetData.FormationCode,
			FormationLocation: getDetailTargetData.FormationLocation,
			Proficiency:       "BEGINNER",
		}
	} else {
		sts, err := lib.GetStudentProfileElastic(msg.Data.SmartBtwID)
		if err == nil {
			sco.SchoolOrigin = sts.LastEdName
			sco.SchoolOriginID = sts.LastEdID
		}
	}

	isBinsus := false
	isSpecificStage := false

	joinedClass, _ := lib.GetStudentJoinedClassType(msg.Data.SmartBtwID)

	for _, k := range joinedClass {
		if strings.Contains(strings.ToLower(k), "binsus") {
			isBinsus = true
			break
		}
	}

	studentProduct, err := lib.GetStudentAOP(uint(msg.Data.SmartBtwID))
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] error get activated product student with smartbtw_id : %d", key, msg.Data.SmartBtwID)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	totalModule := 0
	moduleMultiStage := 32
	for _, prod := range studentProduct.Data {
		for _, tag := range prod.ProductTags {
			if strings.Contains(strings.ToUpper(tag), "PACKAGE_MULTISTAGE_CPNS") {
				isSpecificStage = true
				totalModule += 1
			}
			if strings.Contains(tag, "MAX_STAGE_") {
				mxStg := strings.Split(tag, "_")
				num, err := strconv.Atoi(mxStg[2])
				if err != nil {
					text := fmt.Sprintf("[RabbitMQ][%s] parse stage to int student with smartbtw_id : %d", key, msg.Data.SmartBtwID)
					golog.Slack.ErrorWithData(text, body, err)
					return false
				}
				moduleMultiStage = num
			}
		}
	}
	totalModule = totalModule * moduleMultiStage

	var (
		percentScore    float64 = 0
		grade           string
		proficiency     string
		examProficiency string
	)

	if !msg.Data.AllPassStatus {
		percentScore = 0
	} else {
		if student.TargetScore < 1 {
			percentScore = 99
		} else {
			percentScore = (msg.Data.Total / float64(student.TargetScore)) * 100
		}
	}

	if percentScore >= 100 {
		percentScore = 99
	}

	if percentScore > 90 {
		grade = "DIAMOND"
	} else if percentScore > 75 && percentScore <= 90 {
		grade = "PLATINUM"
	} else if percentScore > 50 && percentScore <= 75 {
		grade = "GOLD"
	} else if percentScore > 25 && percentScore <= 50 {
		grade = "BRONZE"
	} else if percentScore <= 25 {
		grade = "NONE"
	} else {
		text := fmt.Sprintf("[RabbitMQ][%s] error searching percent with percent: %f", key, percentScore)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	insHis := request.CreateHistoryCpns{
		SmartBtwID:        msg.Data.SmartBtwID,
		TaskID:            msg.Data.TaskID,
		PackageID:         msg.Data.PackageID,
		PackageType:       msg.Data.PackageType,
		ModuleCode:        msg.Data.ModuleCode,
		ModuleType:        msg.Data.ModuleType,
		Twk:               msg.Data.Twk,
		Tiu:               msg.Data.Tiu,
		Tkp:               msg.Data.Tkp,
		TwkPass:           msg.Data.TwkPass,
		TiuPass:           msg.Data.TiuPass,
		TkpPass:           msg.Data.TkpPass,
		Total:             msg.Data.Total,
		Repeat:            msg.Data.Repeat,
		ExamName:          msg.Data.ExamName,
		Grade:             grade,
		Start:             msg.Data.Start,
		End:               msg.Data.End,
		IsLive:            isLive,
		SchoolOriginID:    sco.SchoolOriginID,
		SchoolOrigin:      sco.SchoolOrigin,
		InstanceID:        student.InstanceID,
		InstanceName:      student.InstanceName,
		PositionID:        student.PositionID,
		PositionName:      student.PositionName,
		CompetitionCpnsID: student.CompetitionCpnsID,
		FormationType:     student.FormationType,
		FormationCode:     student.FormationCode,
		FormationLocation: student.FormationLocation,
		StudentName:       student.Name,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		TiuPassStatus:     msg.Data.TiuPassStatus,
		TkpPassStatus:     msg.Data.TkpPassStatus,
		TwkPassStatus:     msg.Data.TwkPassStatus,
		TargetScore:       student.TargetScore,
		AllPassStatus:     msg.Data.AllPassStatus,
		// TargetID:   getDetailTarget.ID.Hex(),
	}

	//! Preparing data for update to elastic start from here
	// Get how much user doing module
	averages, err := lib.GetStudentHistoryCPNSElastic(msg.Data.SmartBtwID, false)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when get average data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	countModule := len(averages) + 1

	isTaskExist := false
	for ns, k := range averages {
		if k.TaskID == insHis.TaskID {
			averages[ns].Twk = insHis.Twk
			averages[ns].TwkPass = insHis.TwkPass
			averages[ns].TwkPassStatus = insHis.TwkPassStatus
			averages[ns].Tiu = insHis.Tiu
			averages[ns].TiuPass = insHis.TiuPass
			averages[ns].TiuPassStatus = insHis.TiuPassStatus
			averages[ns].Tkp = insHis.Tkp
			averages[ns].TkpPass = insHis.TkpPass
			averages[ns].TkpPassStatus = insHis.TkpPassStatus
			averages[ns].AllPassStatus = insHis.AllPassStatus
			averages[ns].Total = insHis.Total
			averages[ns].Grade = insHis.Grade
			averages[ns].Repeat = insHis.Repeat
			isTaskExist = true
		}
	}

	if !isTaskExist {
		averages = append(averages, insHis)
	}

	stagesRecord := []request.CreateHistoryCpns{}
	stagesPreUKARecord := []request.CreateHistoryCpns{}
	challengeRecord := []request.CreateHistoryCpns{}
	specificStageRecord := []request.CreateHistoryCpns{}
	// stagesChallengeUKARecord := []request.CreateHistoryPtk{}

	for _, k := range averages {
		if strings.ToLower(k.PackageType) == "pre-uka" {
			stagesPreUKARecord = append(stagesPreUKARecord, k)
		}
		// if strings.ToLower(k.PackageType) == "challenge-uka" {
		// 	stagesChallengeUKARecord = append(stagesChallengeUKARecord, k)
		// }
		if strings.ToLower(k.PackageType) == "pre-uka" || strings.ToLower(k.PackageType) == "challenge-uka" {
			stagesRecord = append(stagesRecord, k)
		}
		if strings.ToLower(k.PackageType) == "multi-stages-uka" {
			specificStageRecord = append(specificStageRecord, k)
		}
	}

	countStagesCompleted := len(stagesRecord)
	countPreStagesCompleted := len(stagesPreUKARecord)

	// Get student proficiency
	if countStagesCompleted >= 43 {
		proficiency = "EXPERT"
	} else if countStagesCompleted >= 15 && countStagesCompleted < 43 {
		proficiency = "ADVANCE"
	} else if countStagesCompleted >= 4 && countStagesCompleted < 15 {
		proficiency = "INTERMEDIATE"
	} else if countStagesCompleted < 4 {
		proficiency = "BEGINNER"
	}

	if countPreStagesCompleted >= 30 {
		examProficiency = "EXPERT"
	} else if countPreStagesCompleted >= 13 && countPreStagesCompleted < 31 {
		examProficiency = "ADVANCE"
	} else if countPreStagesCompleted >= 4 && countPreStagesCompleted < 13 {
		examProficiency = "INTERMEDIATE"
	} else if countPreStagesCompleted < 4 {
		examProficiency = "BEGINNER"
	}

	twkScore := float64(0)
	tiuScore := float64(0)
	tkpScore := float64(0)
	totalScore := float64(0)

	if isBinsus {
		for _, k := range challengeRecord {
			twkScore += k.Twk
			tiuScore += k.Tiu
			tkpScore += k.Tkp
			totalScore += k.Total
		}
	} else {
		for _, k := range averages {
			twkScore += k.Twk
			tiuScore += k.Tiu
			tkpScore += k.Tkp
			totalScore += k.Total
		}
	}

	passingTotalScore := float64(0)
	passingTotalItem := 0

	if !isBinsus {
		for _, k := range averages {
			if strings.ToLower(k.PackageType) != "pre-uka" {
				if (k.Tiu >= k.TiuPass) && (k.Twk >= k.TwkPass) && (k.Tkp >= k.TkpPass) {
					passingTotalItem += 1
					passingTotalScore += k.Total
				}
			}
		}
	} else {
		if len(challengeRecord) < 11 {
			passingTotalItem = 10
		} else {
			passingTotalItem = len(challengeRecord)
		}
	}

	atwk := float64(0)
	atiu := float64(0)
	atkp := float64(0)
	att := float64(0)
	pAtt := float64(0)

	if isBinsus {
		atwk = math.Round(helpers.RoundFloat((twkScore / float64(passingTotalItem)), 2))
		atiu = math.Round(helpers.RoundFloat((tiuScore / float64(passingTotalItem)), 2))
		atkp = math.Round(helpers.RoundFloat((tkpScore / float64(passingTotalItem)), 2))
		att = math.Round(helpers.RoundFloat((totalScore / float64(passingTotalItem)), 2))
		// att = atwk + atiu + atkp
		pAtt = att
	} else if isSpecificStage {
		atwk = helpers.RoundFloat((twkScore / float64(totalModule)), 2)
		atiu = helpers.RoundFloat((tiuScore / float64(totalModule)), 2)
		atkp = helpers.RoundFloat((tkpScore / float64(totalModule)), 2)
		att = helpers.RoundFloat((totalScore / float64(totalModule)), 2)
		pAtt = helpers.RoundFloat((passingTotalScore / float64(totalModule)), 2)
	} else {
		atwk = helpers.RoundFloat((twkScore / float64(len(averages))), 2)
		atiu = helpers.RoundFloat((tiuScore / float64(len(averages))), 2)
		atkp = helpers.RoundFloat((tkpScore / float64(len(averages))), 2)
		att = helpers.RoundFloat((totalScore / float64(len(averages))), 2)
		pAtt = helpers.RoundFloat((passingTotalScore / float64(passingTotalItem)), 2)
	}

	if math.IsNaN(pAtt) {
		pAtt = 0
	}

	percTWK := helpers.RoundFloat((atwk/150)*100, 2)
	percTIU := helpers.RoundFloat((atiu/175)*100, 2)
	percTKP := helpers.RoundFloat((atkp/225)*100, 2)
	percTT := helpers.RoundFloat((att/225)*100, 2)
	percATT := float64(0)
	if student.TargetScore < 1 {
		percATT = 99
	} else {
		percATT = helpers.RoundFloat((pAtt/student.TargetScore)*100, 2)
	}

	if percATT > 99 {
		percATT = 99
	}

	complatedStage := len(stagesRecord)
	if isSpecificStage {
		complatedStage += len(specificStageRecord)
	}

	stUpdate := request.StudentProfileCPNSElastic{
		SmartbtwID:                           msg.Data.SmartBtwID,
		Name:                                 student.Name,
		Photo:                                student.Photo,
		ModuleDone:                           countModule,
		TwkAvgScore:                          atwk,
		TiuAvgScore:                          atiu,
		TkpAvgScore:                          atkp,
		TwkAvgPercentScore:                   percTWK,
		TiuAvgPercentScore:                   percTIU,
		TkpAvgPercentScore:                   percTKP,
		TotalAvgPercentScore:                 percTT,
		TotalAvgScore:                        att,
		LatestTotalScore:                     msg.Data.Total,
		LatestTotalPercentScore:              percentScore,
		InstanceID:                           student.InstanceID,
		InstanceName:                         student.InstanceName,
		FormationCode:                        student.FormationCode,
		PositionID:                           student.PositionID,
		PositionName:                         student.PositionName,
		TargetType:                           "CPNS",
		TargetScore:                          float64(student.TargetScore),
		FormationType:                        student.FormationType,
		CompetitionCpnsID:                    student.CompetitionCpnsID,
		FormationLocation:                    student.FormationLocation,
		Proficiency:                          proficiency,
		ExamProficiency:                      examProficiency,
		StagesCompleted:                      len(stagesRecord),
		PassingRecommendationAvgScore:        pAtt,
		PassingRecommendationAvgPercentScore: percATT,
	}

	// ! NOTE: insert lib update to elastic here
	err = lib.InsertStudentCPNSProfileElastic(&stUpdate, fmt.Sprintf("%d_CPNS", msg.Data.SmartBtwID))
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's CPNS profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	summary := request.FirebaseHistoryScoresSummary{
		ScoreKeys: []string{"TWK", "TIU", "TKP"},
		ScoreValues: request.FirebaseHistoryScoresScoreValues{
			Twk: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        atwk,
				AvgPercentScore: percTWK,
			},
			Tiu: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        atiu,
				AvgPercentScore: percTIU,
			},
			Tkp: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        atkp,
				AvgPercentScore: percTKP,
			},
			Total: request.FirebaseHistoryScoresScoreValue{
				AvgScore:        att,
				AvgPercentScore: percTT,
			},
		},
		LatestTotalScore:          msg.Data.Total,
		LatestTotalPercentScore:   percentScore,
		CurrentTargetTotalScore:   pAtt,
		CurrentTargetPercentScore: percATT,
	}
	totalStages := 32
	ukaPassed := 0
	historyResults := []request.ProfileStudentResultHistory{}
	for _, h := range averages {
		if h.PackageType != "challenge-uka" && h.PackageType != "pre-uka" {
			totalStages += 1
		}

		if h.Tiu >= h.TiuPass && h.Tkp >= h.TkpPass && h.Twk >= h.TwkPass {
			ukaPassed += 1
		}

		formattedTime := h.CreatedAt.Format("2006-01-02")
		res := request.ProfileStudentResultHistory{
			Date:  formattedTime,
			Title: h.ExamName,
			Total: int(h.Total),
		}
		historyResults = append(historyResults, res)
	}

	target := request.FirebaseTargetCPNS{
		CompetitionID:          student.CompetitionCpnsID,
		FormationCode:          student.FormationCode,
		FormationLocation:      student.FormationLocation,
		FormationType:          student.FormationType,
		InstanceID:             uint(student.InstanceID),
		InstanceName:           student.InstanceName,
		InstanceLogo:           "-",
		PositionID:             uint(student.PositionID),
		PositionName:           student.PositionName,
		TargetChancePercentage: float32(percATT),
		TargetScore:            int(student.TargetScore),
		TargetType:             "CPNS",
		Type:                   "CPNS",
		MaximumScore:           550,
	}

	chance, err := lib.GetCompetitionChances(student.FormationCode, uint(student.PositionID), student.FormationType)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when fetch cpns chances comp map data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
	}

	valMajorQuota := 0
	valMajorRegistrant := 0
	valMajorQuotaYear := 0
	valMajorChances := "1:1"
	valMajorCompYear := 0
	if len(chance) > 0 {
		valMajorQuota = int(chance[0].Quota)
		valMajorRegistrant = int(chance[0].Registered)
		valMajorQuotaYear = int(chance[0].Year)
		valMajorCompYear = int(chance[0].Year)

		for _, k := range chance {
			if valMajorRegistrant > 0 {
				break
			}
			if k.Registered > 0 {
				valMajorRegistrant = int(k.Registered)
				valMajorCompYear = int(k.Year)
			}
		}
		thg := math.Round(float64(valMajorRegistrant) / float64(valMajorQuota))
		valMajorChances = fmt.Sprintf("1:%.0f", thg)
	}
	target.MajorQuota = &valMajorQuota
	target.MajorReqistrant = &valMajorRegistrant
	target.MajorQuotaYear = &valMajorQuotaYear
	target.MajorCompYear = &valMajorCompYear
	target.MajorChances = &valMajorChances

	_ = lib.SaveResultFirebase(&request.FirebaseHistoryScores{
		TaskID:          insHis.TaskID,
		SmartbtwID:      uint(insHis.SmartBtwID),
		ExamName:        insHis.ExamName,
		Grade:           insHis.Grade,
		TargetType:      student.TargetType,
		TargetScore:     student.TargetScore,
		InstanceID:      student.InstanceID,
		InstanceName:    student.InstanceName,
		PositionID:      student.PositionID,
		PositionName:    student.PositionName,
		Summary:         summary,
		Result:          historyResults,
		TargetCPNS:      target,
		Proficiency:     proficiency,
		ExamProficiency: examProficiency,
		TotalStages:     totalStages,
		UKAPassed:       ukaPassed,
		StagesDone:      len(stagesRecord),
	}, "skd", "cpns")

	historyCPNSID, err := lib.UpsertHistoryCPNS(&insHis)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	err = lib.InsertStudentHistoryCPNSElastic(&insHis, *historyCPNSID)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's CPNS score history to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Data.ModuleType == string(models.UkaPremium) || msg.Data.ModuleType == string(models.Package) {
		//insert poin to wallet history
		// w := request.CreateWalletHistory{
		// 	SmartbtwID:  msg.Data.SmartBtwID,
		// 	Description: fmt.Sprintf("Mengerjakan UKA %s (%s)", msg.Data.ExamName, msg.Data.ModuleCode),
		// }
		// _, er := lib.CreateWalletHistoryUKA(&w)
		// if er != nil {
		// 	text := fmt.Sprintf("[RabbitMQ][%s] Error when creating wallet history", key)
		// 	golog.Slack.ErrorWithData(text, body, er)
		// 	log.Println(er)
		// 	return false
		// }

		erReward := lib.SendExamReward(uint(msg.Data.SmartBtwID), msg.Data.Total, student.TargetScore, msg.Data.PackageType)
		if erReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d reward of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, erReward)
			return false
		}

		if msg.Data.IsLive {
			erReward := lib.UpdatePassingPercentageCPNS(uint(msg.Data.SmartBtwID), "cpns", uint(msg.Data.PackageID), percentScore, "stages-uka")
			if erReward != nil {
				text := fmt.Sprintf("[RabbitMQ][%s] Error when updating student %d passing percentage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
				golog.Slack.ErrorWithData(text, body, erReward)
			}
		}
		msgBody := request.BodySendAssessmentCompleted{
			PackageID:  uint(msg.Data.PackageID),
			SmartbtwID: uint(msg.Data.SmartBtwID),
			Grade:      grade,
			Score:      msg.Data.Total,
			TaskID:     uint(msg.Data.TaskID),
			Program:    "CPNS",
		}
		errPub := lib.SendAssessmentCompleted(&msgBody)
		if errPub != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when publish student %d CPNS assessment complete of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, errPub)
		}
	}

	if msg.Data.ModuleType == string(models.UkaCode) {
		erReward := lib.UpdatePassingPercentageCPNS(uint(msg.Data.SmartBtwID), "cpns", uint(msg.Data.PackageID), percentScore, "uka-code")
		if erReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when updating student %d passing percentage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, erReward)
		}
	}

	if msg.Data.ModuleType == string(models.UkaFree) {
		errReward := lib.SendUKAFreeReward(uint(msg.Data.SmartBtwID), "CPNS")
		if errReward != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d reward changemajor of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
			golog.Slack.ErrorWithData(text, body, errReward)
			return false
		}
	}

	err = lib.SendHistoryStage(uint(msg.Data.SmartBtwID), uint(msg.Data.PackageID), "CPNS")
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when sending student %d history stage of package id %d", key, msg.Data.SmartBtwID, msg.Data.PackageID)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] insert history CPNS data for user_id %d and task_id %d success", key, msg.Data.SmartBtwID, msg.Data.TaskID)
	log.Println(text)
	lib.LogEvent(
		"InsertScoreCPNS",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data successfully",
		"INFO",
		fmt.Sprintf("profile-%s", key))

	if strings.ToUpper(msg.Data.ModuleType) != "PRE_TEST" || strings.ToUpper(msg.Data.ModuleType) != "POST_TEST" {
		msgBodys := map[string]any{
			"version": 1,
			"data":    insHis,
		}
		msgJsons, errs := sonic.Marshal(msgBodys)
		if errs != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when publish to generate pdf raport student %d task_id %d", key, msg.Data.SmartBtwID, msg.Data.TaskID)
			golog.Slack.ErrorWithData(text, body, err)
			return true
		}
		if errs == nil && db.Broker != nil {
			err = db.Broker.Publish(
				"raport-cpns.build",
				"application/json",
				[]byte(msgJsons), // message to publish
			)
		}
	}

	err = lib.SendToGenerateRaport("CPNS", uint(msg.Data.SmartBtwID), msg.Data.PackageType)
	if err != nil {
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when generating pdf raport progress student %d task_id %d", key, msg.Data.SmartBtwID, msg.Data.TaskID)
			golog.Slack.ErrorWithData(text, body, err)
			return true
		}
	}
	return true
}
