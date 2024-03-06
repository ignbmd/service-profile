package lib

import (
	"context"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/klauspost/lctime"
	"github.com/pandeptwidyaop/gorabbit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateRaportResult(req models.ResultRaport) error {
	collection := db.Mongodb.Collection("result_raports")
	filter := bson.M{
		"smartbtw_id":  req.SmartbtwID,
		"program":      req.Program,
		"task_id":      req.TaskID,
		"stage_type":   req.StageType,
		"module_type":  req.ModuleType,
		"package_type": req.PackageType,
	}

	update := bson.M{
		"$set": bson.M{
			"student_name": req.StudentName,
			"exam_name":    req.ExamName,
			"module_code":  req.ModuleCode,
			"score":        req.Score,
			"score_ptn":    req.ScorePTN,
			"link":         req.Link,
			"updated_at":   time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func CreateFinalRaport(req models.ResultRaport) error {
	switch req.Program {
	case "PTK":
		his, err := GetHistoryPTKElasticByTaskID(uint(req.SmartbtwID), uint(req.TaskID))
		if err != nil {
			return err
		}

		stg := "UMUM"
		if strings.ToLower(his.PackageType) == "multi-stages-uka" {
			stg = "KELAS"
			err := SendRaport(uint(req.SmartbtwID), req.Link)
			if err != nil {
				return err
			}
		}
		req.ExamName = his.ExamName
		req.ModuleCode = his.ModuleCode
		req.StudentName = his.StudentName
		req.Score.TWK = his.Twk
		req.Score.TIU = his.Tiu
		req.Score.TKP = his.Tkp
		req.Score.Total = his.Total
		req.Score.IsPass = his.AllPassStatus
		req.StudentName = his.StudentName
		req.StageType = stg
		req.ModuleType = strings.ToLower(his.ModuleType)
		req.PackageType = strings.ToLower(his.PackageType)

		err = CreateRaportResult(req)
		if err != nil {
			return err
		}
	case "PTN":
		his, err := GetHistoryPTNElasticByTaskID(uint(req.SmartbtwID), uint(req.TaskID))
		if err != nil {
			return err
		}
		stg := "UMUM"
		if strings.ToLower(his.PackageType) == "multi-stages-uka" {
			stg = "KELAS"
			err := SendRaport(uint(req.SmartbtwID), req.Link)
			if err != nil {
				return err
			}
		}
		req.ExamName = his.ExamName
		req.ModuleCode = his.ModuleCode
		req.StudentName = his.StudentName
		req.ScorePTN.PenalaranUmum = his.PenalaranUmum
		req.ScorePTN.PotensiKognitif = his.PotensiKognitif
		req.ScorePTN.LiterasiBahasaIndonesia = his.LiterasiBahasaIndonesia
		req.ScorePTN.LiterasiBahasaInggris = his.LiterasiBahasaInggris
		req.ScorePTN.PemahamanBacaan = his.PemahamanBacaan
		req.ScorePTN.PenalaranMatematika = his.PenalaranMatematika
		req.ScorePTN.PengetahuanKuantitatif = his.PengetahuanKuantitatif
		req.ScorePTN.Total = his.Total
		isPass := his.Total >= his.TargetScore
		req.ScorePTN.IsPass = isPass
		req.StudentName = his.StudentName
		req.StageType = stg
		req.ModuleType = strings.ToLower(his.ModuleType)
		req.PackageType = strings.ToLower(his.PackageType)

		err = CreateRaportResult(req)
		if err != nil {
			return err
		}
	case "CPNS":
		his, err := GetHistoryCPNSElastic(uint(req.SmartbtwID), req.TaskID)
		if err != nil {
			return err
		}
		if strings.ToLower(his.ModuleType) == "with_code" {
			err := SendRaportUKACodeCPNS(uint(req.SmartbtwID), req.Link, his.ExamName)
			if err != nil {
				return err
			}
		}
		stg := "UMUM"
		if strings.ToLower(his.PackageType) == "multi-stages-uka" {
			stg = "KELAS"
			err := SendRaport(uint(req.SmartbtwID), req.Link)
			if err != nil {
				return err
			}
		}
		req.ExamName = his.ExamName
		req.ModuleCode = his.ModuleCode
		req.StudentName = his.StudentName
		req.Score.TWK = his.Twk
		req.Score.TIU = his.Tiu
		req.Score.TKP = his.Tkp
		req.Score.Total = his.Total
		req.Score.IsPass = his.AllPassStatus
		req.StudentName = his.StudentName
		req.StageType = stg
		req.ModuleType = strings.ToLower(his.ModuleType)
		req.PackageType = strings.ToLower(his.PackageType)

		err = CreateRaportResult(req)
		if err != nil {
			return err
		}

	}
	return nil
}

func BuildPDFRaportPTK(bd request.CreateHistoryPtk) error {
	if strings.ToUpper(bd.ModuleType) == "PRE_TEST" || strings.ToUpper(bd.ModuleType) == "POST_TEST" || strings.ToUpper(bd.ModuleType) == "TESTING" {
		return nil
	}

	student, err := GetStudentProfileElastic(bd.SmartBtwID)
	if err != nil {
		return err
	}
	dateTest := time.Now().Format("02-01-2006")
	formattedBirtdate := student.BirthDate.Format("02-01-2006")

	scores, ansres, recommendation, err := CalculateScorePTK(uint(bd.TaskID), uint(bd.SmartBtwID), bd)
	if err != nil {
		return err
	}
	totalIndex := int(0)
	totalItemIndex := int(0)
	totalScore := float64(0)
	for _, sc := range scores {
		for _, scr := range sc.Scores {
			totalScore += float64(scr.Value)
		}
		totalItemIndex++
		totalIndex += int(sc.PassingIndex)
	}
	avgIndex := float64(totalIndex) / float64(totalItemIndex)
	if math.IsNaN(avgIndex) || math.IsInf(avgIndex, 0) {
		avgIndex = 0
	}

	lctime.SetLocale("id_ID")
	formattedAns := splitAnswersIntoCategories(ansres)
	payloadRaport := mockstruct.ResultRaportBody{
		Assessments: mockstruct.Assessments{
			AssessmentCode:        bd.ExamName,
			AvgIndex:              int(math.Round(avgIndex)),
			CreatedAt:             bd.CreatedAt.String(),
			Date:                  dateTest,
			DateCertificateSigned: lctime.Strftime("%A, %d %B %Y", time.Now()),
			ExamName:              bd.ExamName,
			ModuleCode:            bd.ModuleCode,
			ModuleType:            bd.ModuleType,
			PackageID:             bd.PackageID,
			TaskID:                bd.TaskID,
			PackageType:           bd.PackageType,
			Program:               "skd",
			ProgramType:           "PTK",
			ProgramVariant:        "",
			ProgramVersion:        1,
			ScoreType:             "CLASSICAL",
			Scores:                scores,
			SmartbtwID:            bd.SmartBtwID,
			StudentEmail:          student.Email,
			StudentName:           student.Name,
			Total:                 bd.Total,
		},
		RecordAnswerFormatted: formattedAns,
		RecordAnswerMapped:    ansres,
		Screening: mockstruct.Screening{
			AssessmentCode: bd.ExamName,
			Bio: mockstruct.Bio{
				BirthDate: formattedBirtdate,
				Date:      dateTest,
				Email:     student.Email,
				Gender:    student.Gender,
				Name:      student.Name,
				Origin:    student.LastEdName,
				Phone:     student.Phone,
			},
			PackageID:   bd.PackageID,
			ProgramType: "PTK",
			ScreeningTarget: mockstruct.ScreeningTarget{
				DomicileProvince:    student.DomicileProvince,
				DomicileProvinceID:  int(student.DomicileProvinceID),
				DomicileRegion:      student.DomicileRegion,
				DomicileRegionID:    int(student.DomicileRegionID),
				MajorID:             bd.MajorID,
				MajorName:           bd.MajorName,
				PolbitCompetitionID: bd.PolbitCompetitionID,
				PolbitLocationID:    bd.PolbitLocationID,
				PolbitType:          bd.PolbitType,
				SchoolID:            bd.SchoolID,
				SchoolName:          bd.SchoolName,
				TargetScore:         int(bd.TargetScore),
			},
		},
		StudentRecommendationFormatted: recommendation,
	}

	if len(payloadRaport.RecordAnswerFormatted) == 0 && len(payloadRaport.Assessments.Scores) == 0 {
		return nil
	}

	msgBodys := map[string]any{
		"version": 1,
		"data":    payloadRaport,
	}
	msgJsons, errs := sonic.Marshal(msgBodys)
	if errs != nil {
		return errs
	}
	if errs == nil && db.Broker != nil {
		err = db.Broker.Publish(
			"result-raport.build.delivery",
			"application/json",
			[]byte(msgJsons), // message to publish
		)
	}

	// err = SendToGenerateRaport("PTK", uint(bd.SmartBtwID), bd.PackageType)
	// if err != nil {
	// 	return err
	// }

	// srpJsonAns, err := sonic.Marshal(payloadRaport)
	// if err != nil {
	// 	return errors.New("marshalling " + err.Error())
	// }
	// os.WriteFile(fmt.Sprintf("test_json_%s.json", strings.ToLower("ptk-results")), srpJsonAns, 0644)
	return nil
}

func BuildPDFRaportCPNS(bd request.CreateHistoryCpns) error {
	if strings.ToUpper(bd.ModuleType) == "PRE_TEST" || strings.ToUpper(bd.ModuleType) == "POST_TEST" || strings.ToUpper(bd.ModuleType) == "TESTING" {
		return nil
	}
	student, err := GetStudentProfileElastic(bd.SmartBtwID)
	if err != nil {
		return err
	}
	dateTest := time.Now().Format("02-01-2006")
	formattedBirtdate := student.BirthDate.Format("02-01-2006")

	scores, ansres, recommendation, err := CalculateScoreCPNS(uint(bd.TaskID), uint(bd.SmartBtwID), bd)
	if err != nil {
		return err
	}
	totalIndex := int(0)
	totalItemIndex := int(0)
	totalScore := float64(0)
	for _, sc := range scores {
		for _, scr := range sc.Scores {
			totalScore += float64(scr.Value)
		}
		totalItemIndex++
		totalIndex += int(sc.PassingIndex)
	}
	avgIndex := float64(totalIndex) / float64(totalItemIndex)
	if math.IsNaN(avgIndex) || math.IsInf(avgIndex, 0) {
		avgIndex = 0
	}
	origin := ""
	if student.OriginUniversity != nil {
		origin = *student.OriginUniversity
	}
	lctime.SetLocale("id_ID")
	formattedAns := splitAnswersIntoCategories(ansres)
	payloadRaport := mockstruct.ResultRaportBody{
		Assessments: mockstruct.Assessments{
			AssessmentCode:        bd.ExamName,
			AvgIndex:              int(math.Round(avgIndex)),
			CreatedAt:             bd.CreatedAt.String(),
			Date:                  dateTest,
			DateCertificateSigned: lctime.Strftime("%A, %d %B %Y", time.Now()),
			ExamName:              bd.ExamName,
			ModuleCode:            bd.ModuleCode,
			ModuleType:            bd.ModuleType,
			PackageID:             bd.PackageID,
			TaskID:                bd.TaskID,
			PackageType:           bd.PackageType,
			Program:               "skd",
			ProgramType:           "CPNS",
			ProgramVariant:        "",
			ProgramVersion:        1,
			ScoreType:             "CLASSICAL",
			Scores:                scores,
			SmartbtwID:            bd.SmartBtwID,
			StudentEmail:          student.Email,
			StudentName:           student.Name,
			Total:                 bd.Total,
		},
		RecordAnswerFormatted: formattedAns,
		RecordAnswerMapped:    ansres,
		Screening: mockstruct.Screening{
			AssessmentCode: bd.ExamName,
			Bio: mockstruct.Bio{
				BirthDate: formattedBirtdate,
				Date:      dateTest,
				Email:     student.Email,
				Gender:    student.Gender,
				Name:      student.Name,
				Origin:    origin,
				Phone:     student.Phone,
			},
			PackageID:   bd.PackageID,
			ProgramType: "PTK",
			ScreeningTarget: mockstruct.ScreeningTargetCPNS{
				DomicileProvince:   student.DomicileProvince,
				DomicileProvinceID: int(student.DomicileProvinceID),
				DomicileRegion:     student.DomicileRegion,
				DomicileRegionID:   int(student.DomicileRegionID),
				InstanceName:       student.InstanceCPNSName,
				PositionName:       student.PositionCPNSName,
				FormationLocation:  student.FormationCPNSLocation,
				FormationType:      student.FormationCPNSType,
				TargetScore:        int(bd.TargetScore),
			},
		},
		StudentRecommendationFormatted: recommendation,
	}

	if len(payloadRaport.RecordAnswerFormatted) == 0 && len(payloadRaport.Assessments.Scores) == 0 {
		return nil
	}

	msgBodys := map[string]any{
		"version": 1,
		"data":    payloadRaport,
	}
	msgJsons, errs := sonic.Marshal(msgBodys)
	if errs != nil {
		return errs
	}
	if errs == nil && db.Broker != nil {
		err = db.Broker.Publish(
			"result-raport.build.delivery",
			"application/json",
			[]byte(msgJsons), // message to publish
		)
	}
	// err = SendToGenerateRaport("CPNS", uint(bd.SmartBtwID), bd.PackageType)
	// if err != nil {
	// 	return err
	// }

	// srpJsonAns, err := sonic.Marshal(payloadRaport)
	// if err != nil {
	// 	return errors.New("marshalling " + err.Error())
	// }
	// os.WriteFile(fmt.Sprintf("test_json_%s.json", strings.ToLower("ptk-results")), srpJsonAns, 0644)
	return nil
}

func BuildPDFRaportPTN(bd request.CreateHistoryPtn) error {
	if strings.ToUpper(bd.ModuleType) == "PRE_TEST" || strings.ToUpper(bd.ModuleType) == "POST_TEST" || strings.ToUpper(bd.ModuleType) == "TESTING" {
		return nil
	}
	student, err := GetStudentProfileElastic(bd.SmartBtwID)
	if err != nil {
		return err
	}

	dateTest := time.Now().Format("02-01-2006")
	formattedBirtdate := student.BirthDate.Format("02-01-2006")

	scores, ansres, recomendation, err := CalculateScorePTN(uint(bd.TaskID), uint(bd.SmartBtwID), bd)
	if err != nil {
		return err
	}

	formattedAns := splitAnswersIntoCategories(ansres)

	totalIndex := int(0)
	totalItemIndex := int(0)
	totalScore := float64(0)
	for _, sc := range scores {
		for _, scr := range sc.Scores {
			totalScore += float64(scr.Value)
		}
		totalItemIndex++
		totalIndex += int(sc.PassingIndex)
	}
	avgIndex := float64(totalIndex) / float64(totalItemIndex)
	if math.IsNaN(avgIndex) || math.IsInf(avgIndex, 0) {
		avgIndex = 0
	}
	lctime.SetLocale("id_ID")
	payloadRaport := mockstruct.ResultRaportBody{
		Assessments: mockstruct.Assessments{
			AssessmentCode:        bd.ExamName,
			AvgIndex:              int(math.Round(avgIndex)),
			CreatedAt:             bd.CreatedAt.String(),
			Date:                  dateTest,
			DateCertificateSigned: lctime.Strftime("%A, %d %B %Y", time.Now()),
			ExamName:              bd.ExamName,
			ModuleCode:            bd.ModuleCode,
			ModuleType:            bd.ModuleType,
			PackageID:             bd.PackageID,
			TaskID:                bd.TaskID,
			PackageType:           bd.PackageType,
			Program:               "utbk",
			ProgramType:           "PTN",
			ProgramVariant:        "",
			ProgramVersion:        1,
			Scores:                scores,
			SmartbtwID:            bd.SmartBtwID,
			StudentEmail:          student.Email,
			StudentName:           student.Name,
			Total:                 bd.Total,
		},
		RecordAnswerFormatted: formattedAns,
		RecordAnswerMapped:    ansres,
		Screening: mockstruct.Screening{
			AssessmentCode: bd.ExamName,
			Bio: mockstruct.Bio{
				BirthDate: formattedBirtdate,
				Date:      dateTest,
				Email:     student.Email,
				Gender:    student.Gender,
				Name:      student.Name,
				Origin:    student.LastEdName,
				Phone:     student.Phone,
			},
			PackageID:   bd.PackageID,
			ProgramType: "PTN",
			ScreeningTarget: mockstruct.ScreeningTarget{
				DomicileProvince:   student.DomicileProvince,
				DomicileProvinceID: int(student.DomicileProvinceID),
				DomicileRegion:     student.DomicileRegion,
				DomicileRegionID:   int(student.DomicileRegionID),
				MajorID:            bd.MajorID,
				MajorName:          bd.MajorName,
				SchoolID:           bd.SchoolID,
				SchoolName:         bd.SchoolName,
				TargetScore:        int(bd.TargetScore),
			},
		},
		StudentRecommendationFormatted: recomendation,
	}

	if len(payloadRaport.RecordAnswerFormatted) == 0 && len(payloadRaport.Assessments.Scores) == 0 {
		return nil
	}

	msgBodys := map[string]any{
		"version": 1,
		"data":    payloadRaport,
	}
	msgJsons, errs := sonic.Marshal(msgBodys)
	if errs != nil {
		return errs
	}
	if errs == nil && db.Broker != nil {
		err = db.Broker.Publish(
			"result-raport.build.delivery",
			"application/json",
			[]byte(msgJsons), // message to publish
		)
	}

	// err = SendToGenerateRaport("PTN", uint(bd.SmartBtwID), bd.PackageType)
	// if err != nil {
	// 	return err
	// }

	// srpJsonAns, err := sonic.Marshal(payloadRaport)
	// if err != nil {
	// 	return errors.New("marshalling " + err.Error())
	// }
	// os.WriteFile(fmt.Sprintf("test_json_%s.json", strings.ToLower("ptn-results")), srpJsonAns, 0644)
	return nil
}

func CalculateScoreCPNS(taskID uint, smID uint, bd request.CreateHistoryCpns) ([]mockstruct.Score, []mockstruct.StudentAnswerCategory, []mockstruct.Recommendation, error) {
	exp, err := GetStudentExplanation(smID, taskID, "CPNS")
	if err != nil {
		return nil, nil, nil, err
	}
	if len(exp) < 1 {
		err := StoreBuildProcessRaport(models.GetResultRaportBody{
			SmartbtwID:  bd.SmartBtwID,
			Program:     "CPNS",
			TaskID:      bd.TaskID,
			Link:        "",
			StudentName: bd.StudentName,
			ExamName:    bd.ExamName,
			ModuleCode:  bd.ModuleCode,
			StageType:   "",
			ModuleType:  bd.ModuleType,
			PackageType: bd.PackageType,
		})
		if err != nil {
			return nil, nil, nil, err
		}
		return nil, nil, nil, nil
	}

	categoryMap := make(map[string]map[string][]mockstruct.Explanation)

	for _, ex := range exp {
		category := ex.Category
		subcategory := *ex.SubCategory

		if _, ok := categoryMap[category]; !ok {
			categoryMap[category] = make(map[string][]mockstruct.Explanation)
		}

		if _, ok := categoryMap[category][subcategory]; !ok {
			categoryMap[category][subcategory] = []mockstruct.Explanation{}
		}

		categoryMap[category][subcategory] = append(categoryMap[category][subcategory], ex)
	}
	//calculate score
	scores, resultAnswer, recommendation, err := calculateScoreForExplanation(categoryMap, "CPNS")
	if err != nil {
		return nil, nil, nil, err
	}

	for i := range resultAnswer {
		switch resultAnswer[i].CategoryAlias {
		case "TWK":
			resultAnswer[i].Category = "Tes Wawasan Kebangsaan"
			resultAnswer[i].CategoryDisplay = "Tes Wawasan Kebangsaan (TWK)"

		case "TIU":
			resultAnswer[i].Category = "Tes Intelegensi Umum"
			resultAnswer[i].CategoryDisplay = "Tes Intelegensi Umum (TIU)"

		case "TKP":
			resultAnswer[i].Category = "Tes Karakteristik Pribadi"
			resultAnswer[i].CategoryDisplay = "Tes Karakteristik Pribadi (TKP)"
		}
		for j := range resultAnswer[i].Answers {
			switch resultAnswer[i].Answers[j].CategoryName {
			case "TWK":
				resultAnswer[i].Answers[j].CategoryAlias = "Tes Wawasan Kebangsaan"
				if j == 0 {
					resultAnswer[i].Answers[j].CategoryDisplay = "Tes Wawasan Kebangsaan (TWK)"
				} else {
					resultAnswer[i].Answers[j].CategoryDisplay = ""
				}

			case "TIU":
				resultAnswer[i].Answers[j].CategoryAlias = "Tes Intelegensi Umum"
				if j == 0 {
					resultAnswer[i].Answers[j].CategoryDisplay = "Tes Intelegensi Umum (TIU)"
				} else {
					resultAnswer[i].Answers[j].CategoryDisplay = ""
				}

			case "TKP":
				resultAnswer[i].Answers[j].CategoryAlias = "Tes Karakteristik Pribadi"
				if j == 0 {
					resultAnswer[i].Answers[j].CategoryDisplay = "Tes Karakteristik Pribadi (TKP)"
				} else {
					resultAnswer[i].Answers[j].CategoryDisplay = ""
				}
			}
		}
	}

	for i := range recommendation {
		switch recommendation[i].CategoryAlias {
		case "TWK":
			recommendation[i].CategoryOrder = 1
		case "TIU":
			recommendation[i].CategoryOrder = 2
		case "TKP":
			recommendation[i].CategoryOrder = 3
		}
	}

	recommendation = sortAndAdjustRecommendations(recommendation)

	for i := range scores {
		switch scores[i].CategoryAlias {
		case "TWK":
			scores[i].CategoryName = "Tes Wawasan Kebangsaan"
			scores[i].PassingGrade = int(bd.TwkPass)
			scores[i].IsPass = scores[i].Scores[0].Value >= int(bd.TwkPass)
			scores[i].Position = 1
		case "TIU":
			scores[i].CategoryName = "Tes Intelegensi Umum"
			scores[i].PassingGrade = int(bd.TiuPass)
			scores[i].IsPass = scores[i].Scores[0].Value >= int(bd.TiuPass)
			scores[i].Position = 2
			for j := range scores[i].Subtests {
				if scores[i].Subtests[j].SubName == "Verbal Silogis" {
					scores[i].Subtests[j].SubName = "Verbal Silogisme"
					scores[i].Subtests[j].SubAlias = "Verbal Silogisme"
				}
			}
		case "TKP":
			scores[i].CategoryName = "Tes Karakteristik Pribadi"
			scores[i].PassingGrade = int(bd.TkpPass)
			scores[i].IsPass = scores[i].Scores[0].Value >= int(bd.Tkp)
			scores[i].Position = 3
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Position < scores[j].Position
	})

	return scores, resultAnswer, recommendation, nil
}
func CalculateScorePTK(taskID uint, smID uint, bd request.CreateHistoryPtk) ([]mockstruct.Score, []mockstruct.StudentAnswerCategory, []mockstruct.Recommendation, error) {
	exp, err := GetStudentExplanation(smID, taskID, "PTK")
	if err != nil {
		return nil, nil, nil, err
	}
	if len(exp) < 1 {
		err := StoreBuildProcessRaport(models.GetResultRaportBody{
			SmartbtwID:  bd.SmartBtwID,
			Program:     "PTK",
			TaskID:      bd.TaskID,
			Link:        "",
			StudentName: bd.StudentName,
			ExamName:    bd.ExamName,
			ModuleCode:  bd.ModuleCode,
			StageType:   "",
			ModuleType:  bd.ModuleType,
			PackageType: bd.PackageType,
		})
		if err != nil {
			return nil, nil, nil, err
		}
		return nil, nil, nil, nil
	}

	categoryMap := make(map[string]map[string][]mockstruct.Explanation)

	for _, ex := range exp {
		category := ex.Category
		subcategory := *ex.SubCategory

		if _, ok := categoryMap[category]; !ok {
			categoryMap[category] = make(map[string][]mockstruct.Explanation)
		}

		if _, ok := categoryMap[category][subcategory]; !ok {
			categoryMap[category][subcategory] = []mockstruct.Explanation{}
		}

		categoryMap[category][subcategory] = append(categoryMap[category][subcategory], ex)
	}
	//calculate score
	scores, resultAnswer, recommendation, err := calculateScoreForExplanation(categoryMap, "PTK")
	if err != nil {
		return nil, nil, nil, err
	}

	for i := range resultAnswer {
		switch resultAnswer[i].CategoryAlias {
		case "TWK":
			resultAnswer[i].Category = "Tes Wawasan Kebangsaan"
			resultAnswer[i].CategoryDisplay = "Tes Wawasan Kebangsaan (TWK)"

		case "TIU":
			resultAnswer[i].Category = "Tes Intelegensi Umum"
			resultAnswer[i].CategoryDisplay = "Tes Intelegensi Umum (TIU)"

		case "TKP":
			resultAnswer[i].Category = "Tes Karakteristik Pribadi"
			resultAnswer[i].CategoryDisplay = "Tes Karakteristik Pribadi (TKP)"
		}
		for j := range resultAnswer[i].Answers {
			switch resultAnswer[i].Answers[j].CategoryName {
			case "TWK":
				resultAnswer[i].Answers[j].CategoryAlias = "Tes Wawasan Kebangsaan"
				if j == 0 {
					resultAnswer[i].Answers[j].CategoryDisplay = "Tes Wawasan Kebangsaan (TWK)"
				} else {
					resultAnswer[i].Answers[j].CategoryDisplay = ""
				}

			case "TIU":
				resultAnswer[i].Answers[j].CategoryAlias = "Tes Intelegensi Umum"
				if j == 0 {
					resultAnswer[i].Answers[j].CategoryDisplay = "Tes Intelegensi Umum (TIU)"
				} else {
					resultAnswer[i].Answers[j].CategoryDisplay = ""
				}

			case "TKP":
				resultAnswer[i].Answers[j].CategoryAlias = "Tes Karakteristik Pribadi"
				if j == 0 {
					resultAnswer[i].Answers[j].CategoryDisplay = "Tes Karakteristik Pribadi (TKP)"
				} else {
					resultAnswer[i].Answers[j].CategoryDisplay = ""
				}
			}
		}
	}

	for i := range recommendation {
		switch recommendation[i].CategoryAlias {
		case "TWK":
			recommendation[i].CategoryOrder = 1
		case "TIU":
			recommendation[i].CategoryOrder = 2
		case "TKP":
			recommendation[i].CategoryOrder = 3
		}
	}

	recommendation = sortAndAdjustRecommendations(recommendation)

	for i := range scores {
		switch scores[i].CategoryAlias {
		case "TWK":
			scores[i].CategoryName = "Tes Wawasan Kebangsaan"
			scores[i].PassingGrade = int(bd.TwkPass)
			scores[i].IsPass = scores[i].Scores[0].Value >= int(bd.TwkPass)
			scores[i].Position = 1
		case "TIU":
			scores[i].CategoryName = "Tes Intelegensi Umum"
			scores[i].PassingGrade = int(bd.TiuPass)
			scores[i].IsPass = scores[i].Scores[0].Value >= int(bd.TiuPass)
			scores[i].Position = 2
			for j := range scores[i].Subtests {
				if scores[i].Subtests[j].SubName == "Verbal Silogis" {
					scores[i].Subtests[j].SubName = "Verbal Silogisme"
					scores[i].Subtests[j].SubAlias = "Verbal Silogisme"
				}
			}
		case "TKP":
			scores[i].CategoryName = "Tes Karakteristik Pribadi"
			scores[i].PassingGrade = int(bd.TkpPass)
			scores[i].IsPass = scores[i].Scores[0].Value >= int(bd.Tkp)
			scores[i].Position = 3
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Position < scores[j].Position
	})

	return scores, resultAnswer, recommendation, nil
}

func CalculateScorePTN(taskID uint, smID uint, bd request.CreateHistoryPtn) ([]mockstruct.Score, []mockstruct.StudentAnswerCategory, []mockstruct.Recommendation, error) {
	exp, err := GetStudentExplanation(smID, taskID, "PTN")
	if err != nil {
		return nil, nil, nil, err
	}
	if len(exp) < 1 {
		err := StoreBuildProcessRaport(models.GetResultRaportBody{
			SmartbtwID:  bd.SmartBtwID,
			Program:     "PTN",
			TaskID:      bd.TaskID,
			Link:        "",
			StudentName: bd.StudentName,
			ExamName:    bd.ExamName,
			ModuleCode:  bd.ModuleCode,
			StageType:   "",
			ModuleType:  bd.ModuleType,
			PackageType: bd.PackageType,
		})
		if err != nil {
			return nil, nil, nil, err
		}
		return nil, nil, nil, nil
	}

	categoryMap := make(map[string][]mockstruct.Explanation)

	for _, ex := range exp {
		category := ex.Category

		categoryMap[category] = append(categoryMap[category], ex)
	}
	scores, resultAnswer, recommendation, err := calculateScoreForExplanationPTN(categoryMap, "PTN")
	if err != nil {
		return nil, nil, nil, err
	}
	recommendation = sortAndAdjustRecommendations(recommendation)

	for i := range scores {
		switch scores[i].CategoryAlias {
		case "pengetahuan_kuantitatif":
			scores[i].CategoryName = "Pengetahuan Kuantitatif"
			scores[i].Position = 4
			for x := range scores[i].Scores {
				scores[i].Scores[x].Value = int(bd.PengetahuanKuantitatif)
			}
		case "penalaran_umum":
			scores[i].CategoryName = "Penalaran Umum"
			scores[i].Position = 1
			for x := range scores[i].Scores {
				scores[i].Scores[x].Value = int(bd.PenalaranUmum)
			}
		case "literasi_bahasa_indonesia":
			scores[i].CategoryName = "Literasi Bahasa Indonesia"
			scores[i].Position = 5
			for x := range scores[i].Scores {
				scores[i].Scores[x].Value = int(bd.LiterasiBahasaIndonesia)
			}
		case "literasi_bahasa_inggris":
			scores[i].CategoryName = "Literasi Bahasa Inggris"
			scores[i].Position = 6
			for x := range scores[i].Scores {
				scores[i].Scores[x].Value = int(bd.LiterasiBahasaInggris)
			}
		case "penalaran_matematika":
			scores[i].CategoryName = "Penalaran Matematika"
			scores[i].Position = 7
			for x := range scores[i].Scores {
				scores[i].Scores[x].Value = int(bd.PenalaranMatematika)
			}
		case "pengetahuan_umum":
			scores[i].CategoryName = "Pengetahuan Umum"
			scores[i].Position = 2
			for x := range scores[i].Scores {
				scores[i].Scores[x].Value = int(bd.PengetahuanUmum)
			}
		case "pemahaman_bacaan":
			scores[i].CategoryName = "Pemahaman Bacaan"
			scores[i].Position = 3
			for x := range scores[i].Scores {
				scores[i].Scores[x].Value = int(bd.PemahamanBacaan)
			}
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Position < scores[j].Position
	})

	for j := range resultAnswer {
		for k := range resultAnswer[j].Answers {
			{
				switch resultAnswer[j].Answers[k].CategoryName {
				case "pengetahuan_kuantitatif":
					if k == 0 {
						resultAnswer[j].Answers[k].CategoryDisplay = "Pengetahuan Kuantitatif"
					} else {
						resultAnswer[j].Answers[k].CategoryDisplay = ""
					}
				case "penalaran_umum":
					if k == 0 {
						resultAnswer[j].Answers[k].CategoryDisplay = "Penalaran Umum"
					} else {
						resultAnswer[j].Answers[k].CategoryDisplay = ""
					}
				case "literasi_bahasa_indonesia":
					if k == 0 {
						resultAnswer[j].Answers[k].CategoryDisplay = "Literasi Bahasa Indonesia"
					} else {
						resultAnswer[j].Answers[k].CategoryDisplay = ""
					}
				case "literasi_bahasa_inggris":
					if k == 0 {
						resultAnswer[j].Answers[k].CategoryDisplay = "Literasi Bahasa Inggris"
					} else {
						resultAnswer[j].Answers[k].CategoryDisplay = ""
					}
				case "penalaran_matematika":
					if k == 0 {
						resultAnswer[j].Answers[k].CategoryDisplay = "Penalaran Matematika"
					} else {
						resultAnswer[j].Answers[k].CategoryDisplay = ""
					}
				case "pengetahuan_umum":
					if k == 0 {
						resultAnswer[j].Answers[k].CategoryDisplay = "Pengetahuan Umum"
					} else {
						resultAnswer[j].Answers[k].CategoryDisplay = ""
					}
				case "pemahaman_bacaan":
					if k == 0 {
						resultAnswer[j].Answers[k].CategoryDisplay = "Pemahaman Bacaan"
					} else {
						resultAnswer[j].Answers[k].CategoryDisplay = ""
					}
				}
			}
		}

	}

	for l := range recommendation {
		switch recommendation[l].CategoryName {
		case "pengetahuan_kuantitatif":
			recommendation[l].CategoryName = "Pengetahuan Kuantitatif"
		case "penalaran_umum":
			recommendation[l].CategoryName = "Penalaran Umum"
		case "literasi_bahasa_indonesia":
			recommendation[l].CategoryName = "Literasi Bahasa Indonesia"
		case "literasi_bahasa_inggris":
			recommendation[l].CategoryName = "Literasi Bahasa Inggris"
		case "penalaran_matematika":
			recommendation[l].CategoryName = "Penalaran Matematika"
		case "pengetahuan_umum":
			recommendation[l].CategoryName = "Pengetahuan Umum"
		case "pemahaman_bacaan":
			recommendation[l].CategoryName = "Pemahaman Bacaan"

		}
	}

	return scores, resultAnswer, recommendation, nil
}
func calculateScoreForExplanation(ex map[string]map[string][]mockstruct.Explanation, program string) ([]mockstruct.Score, []mockstruct.StudentAnswerCategory, []mockstruct.Recommendation, error) {

	var resultScores []mockstruct.Score
	var resultAnswer []mockstruct.StudentAnswerCategory
	var resultRecommendation []mockstruct.Recommendation

	for category, subcategories := range ex {
		ansOps := []mockstruct.StudentAnswerCompetition{}

		ans := mockstruct.StudentAnswerCategory{
			Category:        category,
			CategoryAlias:   category,
			CategoryDisplay: category,
		}

		catVal := int(0)
		var categoryScore mockstruct.Score

		var correctAnswer, wrongAnswer, emptyAnswer int

		subscoreVal := []mockstruct.Scores{}

		for subc, subcategory := range subcategories {
			var subCorrectAnswer, subWrongAnswer, subEmptyAnswer int
			subcatVal := int(0)
			recommendations := mockstruct.Recommendation{
				SubCategoryName:  subc,
				SubCategoryAlias: subc,
				CategoryAlias:    category,
				CategoryName:     category,
			}

			recTopic := []string{}

			for _, explanation := range subcategory {

				if explanation.Answered != explanation.AnswerKey {
					masterQuestion, err := GetSingleQuestionFromMaster(explanation.Code, program)
					if err != nil {
						return nil, nil, nil, err
					}
					recTopic = append(recTopic, masterQuestion.QuestionKeyword...)
				}

				totalTrueQuestion := int(0)
				questions, err := GetStudentExplanationPerQuestion(explanation.Code, uint(explanation.LegacyTaskID), program)
				if err != nil {
					return nil, nil, nil, err
				}
				for _, qs := range questions {
					if qs.AnswerKey == qs.Answered {
						totalTrueQuestion++
					}
				}

				correctStudentPercentage := math.Round(float64(totalTrueQuestion) / float64(len(questions)) * 100)
				if math.IsNaN(correctStudentPercentage) || math.IsInf(correctStudentPercentage, 0) {
					correctStudentPercentage = 0
				}

				choosenOps := "-"
				optionIndex := indexOf(explanation.OptionID, explanation.Answered)
				if optionIndex != -1 {
					choosenOps = explanation.OptionTypes[optionIndex]
				}

				corectOps := ""
				corectOpsIndex := indexOf(explanation.OptionID, explanation.AnswerKey)
				if corectOpsIndex != -1 {
					corectOps = explanation.OptionTypes[corectOpsIndex]
				}

				point := 0
				if explanation.Answered == explanation.AnswerKey {
					if optionIndex != -1 {
						subcatVal += explanation.OptionValues[optionIndex]
						catVal += explanation.OptionValues[optionIndex]
						point = explanation.OptionValues[optionIndex]
					}
					subCorrectAnswer++
				} else if explanation.Answered == 0 {
					subEmptyAnswer++
				} else {
					point = explanation.OptionValues[optionIndex]
					catVal += explanation.OptionValues[optionIndex]
					subcatVal += explanation.OptionValues[optionIndex]
					subWrongAnswer++
				}
				subscoreVal = []mockstruct.Scores{
					{
						ScoreType: "",
						Value:     subcatVal,
					},
				}

				an := mockstruct.StudentAnswerCompetition{
					CategoryAlias:            category,
					CategoryDisplay:          category,
					SubCategoryAlias:         *explanation.SubCategory,
					CategoryName:             category,
					SubCategoryName:          *explanation.SubCategory,
					AnswerType:               explanation.AnswerType,
					Order:                    uint(explanation.Order),
					ChoosenAnswer:            choosenOps,
					CorrectAnswer:            corectOps,
					CorrectStudent:           totalTrueQuestion,
					CorrectStudentPercentage: correctStudentPercentage,
					TotalStudent:             len(questions),
					Point:                    float64(point),
					IsTrue:                   explanation.Answered == explanation.AnswerKey,
				}
				sort.Slice(ansOps, func(i, j int) bool {
					return ansOps[i].Order < ansOps[j].Order
				})
				ansOps = append(ansOps, an)

			}
			topic := helpers.RemoveDuplicates(recTopic)
			recommendations.RecommendedTopic = topic
			recommendations.RecommendedTopicFormatted = strings.Join(topic, ", ")
			matDesc := fmt.Sprintf("%v", mockstruct.SubCategoryData[helpers.ToLowerAndUnderscore(subc)])
			resMatDesc, l := mockstruct.SubCategoryData[fmt.Sprintf("%s_%s", helpers.ToLowerAndUnderscore(subc), strings.ToLower(program))]
			if l {
				matDesc = fmt.Sprintf("%v", resMatDesc)
			}

			percentage := float64(0)
			indexCategory := float64(0)

			percentage = float64(subCorrectAnswer) / float64(subWrongAnswer+subEmptyAnswer+subCorrectAnswer) * 100
			indexCategory = float64(subCorrectAnswer) / float64(subWrongAnswer+subEmptyAnswer+subCorrectAnswer) * 9

			if math.IsNaN(percentage) {
				percentage = 0
			}
			if math.IsNaN(indexCategory) {
				indexCategory = 0
			}
			if math.IsInf(percentage, 0) {
				percentage = 0
			}
			if math.IsInf(indexCategory, 0) {
				indexCategory = 0
			}

			recommendations.MaterialDescription = matDesc
			recommendations.PassingIndex = math.Round(indexCategory)
			resultRecommendation = append(resultRecommendation, recommendations)
			if len(topic) < 1 {
				resultRecommendation = nil
			}

			categoryScore.Subtests = append(categoryScore.Subtests, mockstruct.Subtest{
				CorrectAnswer:     subCorrectAnswer,
				WrongAnswer:       subWrongAnswer,
				EmptyAnswer:       subEmptyAnswer,
				PassingIndex:      math.Round(indexCategory),
				PassingPercentage: math.Round(percentage),
				Scores:            subscoreVal,
				SubAlias:          *subcategory[0].SubCategory,
				SubName:           *subcategory[0].SubCategory,
			})

			correctAnswer += subCorrectAnswer
			wrongAnswer += subWrongAnswer
			emptyAnswer += subEmptyAnswer

		}

		sort.Slice(ansOps, func(i, j int) bool {
			return ansOps[i].Order < ansOps[j].Order
		})

		ans.Answers = ansOps

		scoreValCat := []mockstruct.Scores{
			{
				ScoreType: "",
				Value:     catVal,
			},
		}

		percentage := float64(0)
		indexCategory := float64(0)

		percentage = float64(correctAnswer) / float64(wrongAnswer+emptyAnswer+correctAnswer) * 100
		indexCategory = float64(correctAnswer) / float64(wrongAnswer+emptyAnswer+correctAnswer) * 9

		if math.IsNaN(percentage) {
			percentage = 0
		}
		if math.IsNaN(indexCategory) {
			indexCategory = 0
		}
		if math.IsInf(percentage, 0) {
			percentage = 0
		}
		if math.IsInf(indexCategory, 0) {
			indexCategory = 0
		}

		resultScores = append(resultScores, mockstruct.Score{
			CategoryAlias:     category,
			CategoryID:        0,
			CategoryName:      category,
			CorrectAnswer:     correctAnswer,
			WrongAnswer:       wrongAnswer,
			EmptyAnswer:       emptyAnswer,
			PassingGrade:      0,
			PassingIndex:      math.Round(indexCategory),
			PassingPercentage: math.Round(percentage),
			Position:          0,
			Scores:            scoreValCat,
			Subtests:          categoryScore.Subtests,
		})
		resultAnswer = append(resultAnswer, ans)

		resultRecommendation = append(resultRecommendation)

	}

	return resultScores, resultAnswer, resultRecommendation, nil
}

func calculateScoreForExplanationPTN(ex map[string][]mockstruct.Explanation, program string) ([]mockstruct.Score, []mockstruct.StudentAnswerCategory, []mockstruct.Recommendation, error) {

	var resultScores []mockstruct.Score
	var resultAnswer []mockstruct.StudentAnswerCategory
	var resultRecommendation []mockstruct.Recommendation

	for category, expl := range ex {

		recTopic := []string{}
		catVal := int(0)
		var correctAnswer, wrongAnswer, emptyAnswer int
		ansOps := []mockstruct.StudentAnswerCompetition{}

		for _, answ := range expl {
			totalTrueQuestion := int(0)
			questions, err := GetStudentExplanationPerQuestion(answ.Code, uint(answ.LegacyTaskID), program)
			if err != nil {
				return nil, nil, nil, err
			}

			for _, qs := range questions {
				if qs.AnswerKey == qs.Answered {
					totalTrueQuestion++
				}
			}

			correctStudentPercentage := math.Round(float64(totalTrueQuestion) / float64(len(questions)) * 100)
			if math.IsNaN(correctStudentPercentage) || math.IsInf(correctStudentPercentage, 0) {
				correctStudentPercentage = 0
			}

			ns := mockstruct.StudentAnswerCompetition{}
			ns.AnswerType = answ.AnswerType
			ns.CategoryAlias = category
			ns.CategoryDisplay = category
			ns.CategoryName = category
			ns.TotalStudent = len(questions)
			ns.CorrectStudentPercentage = correctStudentPercentage
			ns.CorrectStudent = totalTrueQuestion

			if answ.AnswerType != "MULTIPLE_CHOICES" && answ.AnswerType != "NUMBER" {

				userAnsw := []bool{}
				correctAnsw := []bool{}

				orderedList := []uint{}

				answKey := map[uint]bool{}

				for _, k := range answ.OptionKeys {
					orderedList = append(orderedList, uint(answ.OptionID[k]))
				}

				for id, k := range answ.OptionValues {
					answKey[uint(answ.OptionID[id])] = k == 1
				}
				for _, k := range orderedList {

					answAdded := false

					correctAnsw = append(correctAnsw, answKey[k])

					if len(answ.AnsweredFalseItem) > 0 || len(answ.AnsweredTrueItem) > 0 {
						for _, l := range answ.AnsweredTrueItem {
							if l == k {
								if answKey[l] {
									userAnsw = append(userAnsw, true)
									answAdded = true
								}
							}
						}
						if !answAdded {
							userAnsw = append(userAnsw, false)
						}
					}

				}

				if len(answ.AnsweredFalseItem) < 1 && len(answ.AnsweredTrueItem) < 1 {
					ns.Point = 0
					ns.IsTrue = false
					ns.ChoosenMultiAnswerChoice = []bool{}
					for range correctAnsw {
						ns.ChoosenMultiAnswerChoice = append(ns.ChoosenMultiAnswerChoice, false)
					}
				} else {

					ns.Point = 1
					ns.IsTrue = true

					for idx, ks := range correctAnsw {
						if ks != userAnsw[idx] {
							ns.Point = 0
							ns.IsTrue = false
						}
					}
					ns.ChoosenMultiAnswerChoice = userAnsw
				}
				ns.CorrectMultiAnswerChoice = correctAnsw

				ns.AnswerHeaderFalse = "SALAH"
				ns.AnswerHeaderTrue = "BENAR"

				if answ.AnswerHeaderFalse != nil {
					ns.AnswerHeaderFalse = *answ.AnswerHeaderFalse
				}

				if answ.AnswerHeaderTrue != nil {
					ns.AnswerHeaderTrue = *answ.AnswerHeaderTrue
				}

				ns.Order = uint(answ.Order)
				if len(answ.AnsweredFalseItem) < 1 && len(answ.AnsweredTrueItem) < 1 {
					emptyAnswer++
				}

				if len(answ.AnsweredFalseItem) > 0 && len(answ.AnsweredTrueItem) > 0 && !ns.IsTrue {
					wrongAnswer++
				}
				if ns.IsTrue {
					correctAnswer++
				}

			} else if answ.AnswerType == "NUMBER" {
				choosenEssay := "-"
				isTrue := false
				point := 0
				corect := ""
				if answ.AnsweredEssay == nil {
					emptyAnswer++
				} else if *answ.AnsweredEssay == "" {
					emptyAnswer++
				} else {
					stdAns, err := helpers.ParseStringToInt(*answ.AnsweredEssay)
					if err != nil {
						return nil, nil, nil, err
					}
					ansKey, err := helpers.ParseStringToInt(*answ.Essay)
					if err != nil {
						return nil, nil, nil, err
					}
					corect = fmt.Sprintf("%d", ansKey)

					if stdAns != ansKey {
						choosenEssay = fmt.Sprintf("%d", stdAns)
						wrongAnswer++
					}
					if stdAns == ansKey {
						choosenEssay = fmt.Sprintf("%d", stdAns)
						isTrue = true
						point = 1
						correctAnswer++
					}
				}

				ns = mockstruct.StudentAnswerCompetition{
					CategoryAlias:            category,
					CategoryDisplay:          category,
					SubCategoryAlias:         *answ.SubCategory,
					CategoryName:             category,
					SubCategoryName:          *answ.SubCategory,
					CorrectAnswer:            corect,
					ChoosenAnswer:            choosenEssay,
					Point:                    float64(point),
					Order:                    uint(answ.Order),
					IsTrue:                   isTrue,
					AnswerType:               answ.AnswerType,
					TotalStudent:             len(questions),
					CorrectStudentPercentage: correctStudentPercentage,
					CorrectStudent:           totalTrueQuestion,
				}

			} else {
				if answ.Answered == 0 {
					emptyAnswer++
				}
				if answ.Answered != 0 && answ.Answered != answ.AnswerKey {
					wrongAnswer++
				}
				if answ.Answered == answ.AnswerKey {
					correctAnswer++
				}
				trueIdx := -1

				for idx, k := range answ.OptionID {
					if k == answ.AnswerKey {
						trueIdx = idx
						break
					}
				}

				trueDataIdx := trueIdx
				for randIdx, grpRealAnsIdx := range answ.OptionKeys {
					if grpRealAnsIdx == trueIdx {
						trueDataIdx = randIdx
						break
					}
				}

				answIdx := -1
				isTrue := answ.Answered == answ.AnswerKey
				if answ.Answered > 0 {
					for idx, k := range answ.OptionID {
						if k == answ.Answered {
							answIdx = idx
							break
						}
					}
				}
				if answIdx != -1 {
					ansDataIdx := answIdx
					for randIdx, grpRealAnsIdx := range answ.OptionKeys {
						if grpRealAnsIdx == answIdx {
							ansDataIdx = randIdx
							break
						}
					}
					ns = mockstruct.StudentAnswerCompetition{
						CategoryAlias:            category,
						CategoryDisplay:          category,
						SubCategoryAlias:         *answ.SubCategory,
						CategoryName:             category,
						SubCategoryName:          *answ.SubCategory,
						CorrectAnswer:            strings.ToUpper(answ.OptionTypes[trueDataIdx]),
						ChoosenAnswer:            strings.ToUpper(answ.OptionTypes[ansDataIdx]),
						Point:                    float64(answ.OptionValues[answIdx]),
						Order:                    uint(answ.Order),
						IsTrue:                   isTrue,
						AnswerType:               answ.AnswerType,
						TotalStudent:             len(questions),
						CorrectStudentPercentage: correctStudentPercentage,
						CorrectStudent:           totalTrueQuestion,
					}
				} else {
					ns = mockstruct.StudentAnswerCompetition{
						CategoryAlias:            category,
						CategoryDisplay:          category,
						SubCategoryAlias:         *answ.SubCategory,
						CategoryName:             category,
						SubCategoryName:          *answ.SubCategory,
						CorrectAnswer:            strings.ToUpper(answ.OptionTypes[trueDataIdx]),
						ChoosenAnswer:            "-",
						Point:                    0,
						Order:                    uint(answ.Order),
						IsTrue:                   false,
						AnswerType:               answ.AnswerType,
						TotalStudent:             len(questions),
						CorrectStudentPercentage: correctStudentPercentage,
						CorrectStudent:           totalTrueQuestion,
					}
				}

			}
			if !ns.IsTrue {
				masterQuestion, err := GetSingleQuestionFromMaster(answ.Code, program)
				if err != nil {
					return nil, nil, nil, err
				}
				recTopic = append(recTopic, masterQuestion.QuestionKeyword...)
			}

			catVal += int(ns.Point)
			ansOps = append(ansOps, ns)
		}

		ans := mockstruct.StudentAnswerCategory{
			Category:        category,
			CategoryAlias:   category,
			CategoryDisplay: category,
		}

		sort.Slice(ansOps, func(i, j int) bool {
			return ansOps[i].Order < ansOps[j].Order
		})

		ans.Answers = ansOps

		scoreValCat := []mockstruct.Scores{
			{
				ScoreType: "",
				Value:     catVal,
			},
		}

		percentage := float64(0)
		indexCategory := float64(0)

		percentage = float64(correctAnswer) / float64(wrongAnswer+emptyAnswer+correctAnswer) * 100
		indexCategory = float64(correctAnswer) / float64(wrongAnswer+emptyAnswer+correctAnswer) * 9

		if math.IsNaN(percentage) {
			percentage = 0
		}
		if math.IsNaN(indexCategory) {
			indexCategory = 0
		}
		if math.IsInf(percentage, 0) {
			percentage = 0
		}
		if math.IsInf(indexCategory, 0) {
			indexCategory = 0
		}

		resultScores = append(resultScores, mockstruct.Score{
			CategoryAlias:     category,
			CategoryID:        0,
			CategoryName:      category,
			CorrectAnswer:     correctAnswer,
			WrongAnswer:       wrongAnswer,
			EmptyAnswer:       emptyAnswer,
			PassingGrade:      0,
			PassingIndex:      math.Round(indexCategory),
			PassingPercentage: math.Round(percentage),
			Position:          0,
			Scores:            scoreValCat,
		})
		resultAnswer = append(resultAnswer, ans)

		recommendations := mockstruct.Recommendation{
			CategoryAlias: category,
			CategoryName:  category,
		}

		topic := helpers.RemoveDuplicates(recTopic)
		if len(topic) < 1 {
			resultRecommendation = nil
		} else {
			recommendations.RecommendedTopic = topic
			recommendations.RecommendedTopicFormatted = strings.Join(topic, ", ")
			matDesc := fmt.Sprintf("%v", mockstruct.CategoryData[helpers.ToLowerAndUnderscore(category)])
			resMatDesc, l := mockstruct.CategoryData[fmt.Sprintf("%s_%s", helpers.ToLowerAndUnderscore(category), strings.ToLower(program))]
			if l {
				matDesc = fmt.Sprintf("%v", resMatDesc)
			}
			recommendations.MaterialDescription = matDesc
			recommendations.PassingIndex = math.Round(indexCategory)
			resultRecommendation = append(resultRecommendation, recommendations)
		}

	}

	return resultScores, resultAnswer, resultRecommendation, nil
}

func indexOf(slice []int, value int) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func splitAnswersIntoCategories(studentAnswerCategories []mockstruct.StudentAnswerCategory) []mockstruct.StudentAnswerCategory {
	const maxAnswersPerCategory = 13
	var allAnswers []mockstruct.StudentAnswerCompetition

	// Collect all answers from all categories
	for _, category := range studentAnswerCategories {
		allAnswers = append(allAnswers, category.Answers...)
	}

	// Sort all answers by the Order field
	sort.Slice(allAnswers, func(i, j int) bool {
		return allAnswers[i].Order < allAnswers[j].Order
	})

	var updatedCategories []mockstruct.StudentAnswerCategory

	// Split the sorted answers into categories with a maximum of 15 answers per category
	for i := 0; i < len(allAnswers); i += maxAnswersPerCategory {
		end := i + maxAnswersPerCategory
		if end > len(allAnswers) {
			end = len(allAnswers)
		}

		updatedCategory := mockstruct.StudentAnswerCategory{
			Category:        allAnswers[0].CategoryName,
			CategoryAlias:   allAnswers[0].CategoryAlias,
			CategoryDisplay: allAnswers[0].CategoryDisplay,
			Answers:         allAnswers[i:end],
		}

		updatedCategories = append(updatedCategories, updatedCategory)
	}
	return updatedCategories
}

func sortAndAdjustRecommendations(data []mockstruct.Recommendation) []mockstruct.Recommendation {
	groupedData := make(map[string][]mockstruct.Recommendation)
	for _, item := range data {
		groupedData[item.CategoryAlias] = append(groupedData[item.CategoryAlias], item)
	}

	var result []mockstruct.Recommendation
	for _, group := range groupedData {

		for i := range group {
			if i == 0 {
				group[i].CategoryName = group[i].CategoryAlias
			} else {
				group[i].CategoryName = ""
			}

			result = append(result, group[i])
		}
	}

	return result

}

func TriggerBuildRaport(smartbtwID uint, program string) error {
	switch program {
	case "PTK":
		histories, err := GetHistoryPTK(smartbtwID)
		if err != nil {
			return err
		}
		if len(histories) == 0 {
			return nil
		} else {
			stageKelas := 0
			for _, his := range histories {
				if strings.ToLower(his.PackageType) == "multi-stages-uka" {
					stageKelas++
				}
				err := BuildPDFRaportPTK(his)
				if err != nil {
					continue
				}
				time.Sleep(3 * time.Second)
			}
			if stageKelas > 0 {
				err = SendToGenerateRaport("PTK", smartbtwID, "multi-stages-uka")
				if err != nil {
					if err != nil {
						return err
					}
				}
			} else {
				err = SendToGenerateRaport("PTK", smartbtwID, "UMUM")
				if err != nil {
					if err != nil {
						return err
					}
				}
			}
		}
	case "PTN":
		histories, err := GetHistoryPTN(smartbtwID)
		if err != nil {
			return err
		}
		if len(histories) == 0 {
			return nil
		} else {
			stageKelas := 0
			for _, his := range histories {
				if strings.ToLower(his.PackageType) == "multi-stages-uka" {
					stageKelas++
				}
				err := BuildPDFRaportPTN(his)
				if err != nil {
					continue
				}
				time.Sleep(3 * time.Second)
			}
			if stageKelas > 0 {
				err = SendToGenerateRaport("PTN", smartbtwID, "multi-stages-uka")
				if err != nil {
					if err != nil {
						return err
					}
				}
			} else {
				err = SendToGenerateRaport("PTN", smartbtwID, "UMUM")
				if err != nil {
					if err != nil {
						return err
					}
				}
			}
		}
	case "CPNS":
		histories, err := GetHistoryCPNS(smartbtwID)
		if err != nil {
			return err
		}
		if len(histories) == 0 {
			return nil
		} else {
			stageKelas := 0
			for _, his := range histories {
				if strings.ToLower(his.PackageType) == "multi-stages-uka" {
					stageKelas++
				}
				err := BuildPDFRaportCPNS(his)
				if err != nil {
					continue
				}
				time.Sleep(3 * time.Second)
			}
			if stageKelas > 0 {
				err = SendToGenerateRaport("CPNS", smartbtwID, "multi-stages-uka")
				if err != nil {
					if err != nil {
						return err
					}
				}
			} else {
				err = SendToGenerateRaport("CPNS", smartbtwID, "UMUM")
				if err != nil {
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func BuildRaportByTaskID(smartbtwID uint, program string, taskID int) error {
	switch program {
	case "PTK":
		his, err := GetHistoryPTKElasticByTaskID(smartbtwID, uint(taskID))
		if err != nil {
			return err
		}
		err = BuildPDFRaportPTK(his)
		if err != nil {
			return err
		}
	case "PTN":
		his, err := GetHistoryPTNElasticByTaskID(smartbtwID, uint(taskID))
		if err != nil {
			return err
		}
		err = BuildPDFRaportPTN(his)
		if err != nil {
			return err
		}
	case "CPNS":
		his, err := GetHistoryCPNSElastic(smartbtwID, taskID)
		if err != nil {
			return err
		}
		err = BuildPDFRaportCPNS(his)
		if err != nil {
			return err
		}

	}
	return nil
}

func GetListingRaport(smartbtwID uint, program string, ukaType string, stageType string, search *string) ([]models.GetResultRaportBody, error) {
	collection := db.Mongodb.Collection("result_raports")
	filter := bson.M{"smartbtw_id": smartbtwID, "program": program}

	var mdltype string
	switch strings.ToUpper(ukaType) {
	case "PRE_UKA":
		mdltype = "pre-uka"
	case "ALL_MODULE":
		mdltype = "all-module"
	case "UKA_STAGE":
		mdltype = "challenge-uka"
	}

	if stageType == "UMUM" && (mdltype == "pre-uka" || mdltype == "all-module") {
		filter["$or"] = []bson.M{
			{"package_type": "pre-uka", "module_type": "platinum"},
			{"package_type": "challenge-uka", "module_type": "premium_tryout"},
		}
	} else if stageType == "KELAS" && (mdltype == "pre-uka" || mdltype == "all-module") {
		filter["$and"] = []bson.M{
			{"module_type": "platinum"},
			{"package_type": "multi-stages-uka"},
		}
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var raports []models.GetResultRaportBody
	if err := cursor.All(context.Background(), &raports); err != nil {
		return nil, err
	}

	return raports, nil
}

func SendRequestBuildRaportBulk(smIDs []uint, program string) error {
	for _, id := range smIDs {
		bd := mockstruct.BodyRequestBuildRaport{
			SmartbtwID: id,
			Program:    program,
		}
		msgBodys := map[string]any{
			"version": 1,
			"data":    bd,
		}
		msgJsons, errs := sonic.Marshal(msgBodys)
		if errs != nil {
			return errs
		}
		if errs == nil && db.Broker != nil {
			err := db.Broker.Publish(
				"result-raport.build-bulk.request",
				"application/json",
				[]byte(msgJsons), // message to publish
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CacheReportPTKToRedis(bd request.CreateHistoryPtk) error {
	key := fmt.Sprintf("profile:report:%s:%d", "PTK", bd.SmartBtwID)
	k := fmt.Sprintf("task_id:%d", bd.TaskID)
	body := map[string]any{
		k: bd,
	}
	errs := db.NewRedisCluster().HSet(context.Background(), key, body).Err()
	if errs != nil {
		return errs
	}
	return nil
}

func SendRaport(smartbtwID uint, link string) error {
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
				CustomMessage: fmt.Sprintf("dokumen rapor UKA Stage siswa atas nama %s", parent[0].Name),
				FileName:      fmt.Sprintf("raport_uka_%d.pdf", time.Now().UnixMilli()),
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

func SendRaportUKACodeCPNS(smartbtwID uint, link string, examName string) error {
	type SendRaportBody struct {
		To            string `json:"to"`
		Name          string `json:"name"`
		Greeting      string `json:"greeting"`
		CustomMessage string `json:"custom_message"`
		FileName      string `json:"file_name"`
		FileUrl       string `json:"file_url"`
	}

	student, err := GetStudentProfileElastic(int(smartbtwID))
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("hasil_assessment_%d.pdf", time.Now().UnixMilli())

	if student.Phone != "" {
		payload := SendRaportBody{
			To:            student.Phone,
			Name:          student.Name,
			CustomMessage: fmt.Sprintf("dokumen hasil Tes Asesmen Anda atas nama %s", student.Name),
			FileName:      fileName,
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

	msgBody := map[string]any{
		"version": 1,
		"data": map[string]any{
			"student_name": student.Name,
			"email":        student.Email,
			"phone":        student.Phone,
			// "assessment_code": assGraph.AssessmentCode,
			"assessment_name": examName,
			"greeting":        helpers.GetGreeting(),
			"file_name":       fileName,
			"report_url":      link,
		},
	}

	j, errs := sonic.Marshal(msgBody)
	if errs == nil {
		_ = gorabbit.Context.Publish("message-gateway.email.central-assessment.deliver-raport", "application/json", j)

	}

	return nil
}

func GetResultRaport(smID uint, program string, stgType string) ([]models.GetResultRaportBody, error) {
	collection := db.Mongodb.Collection("result_raports")
	filter := bson.M{"smartbtw_id": smID, "program": program, "stage_type": stgType}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var raports []models.GetResultRaportBody
	if err := cursor.All(context.Background(), &raports); err != nil {
		return nil, err
	}

	return raports, nil
}

func StoreBuildProcessRaport(payload models.GetResultRaportBody) error {
	if os.Getenv("ENV") == "test" {
		return nil
	}
	ns, _ := sonic.Marshal(payload)
	return db.NewRedisCluster().RPush(context.Background(), db.RAPORT_RESULT_REDIS_QUEUE_BUILDY_KEY, string(ns)).Err()
}

func StoreFailedRaport(payload models.GetResultRaportBody) error {
	if os.Getenv("ENV") == "test" {
		return nil
	}
	ns, _ := sonic.Marshal(payload)
	return db.NewRedisCluster().RPush(context.Background(), db.RAPORT_RESULT_REDIS_QUEUE_FAILED_BUILD_KEY, string(ns)).Err()
}

func AddToQueueRegenerateRaport(smIDs []uint, program string, stgType string) error {
	for _, sm := range smIDs {
		res, err := GetResultRaport(sm, program, stgType)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			continue
		}
		for _, t := range res {
			err = StoreBuildProcessRaport(t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
