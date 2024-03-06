package lib

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func GetSingleInterviewScoreByID(id primitive.ObjectID) (models.InterviewScore, error) {
	var result models.InterviewScore
	col := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "deleted_at": nil}
	err := col.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.InterviewScore{}, err
	}

	return result, nil
}

func GetSingleInterviewScoreBySessionIDSSOIDAndStudentID(sessionID primitive.ObjectID, ssoID string, smartbtwID int) (models.InterviewScore, error) {
	var result models.InterviewScore
	col := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"session_id": sessionID, "smartbtw_id": smartbtwID, "created_by.id": ssoID, "deleted_at": nil}
	err := col.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.InterviewScore{}, err
	}

	return result, nil
}

func GetSingleInterviewScoreByStudentIDAndYear(id int, year uint16) (models.InterviewScore, error) {
	var result models.InterviewScore
	col := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": id, "year": year, "deleted_at": nil}
	err := col.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.InterviewScore{}, err
	}

	return result, nil
}

func GetInterviewScoreByArrayOfStudentEmail(email []string, year uint16) ([]models.InterviewScoreEmailEdutech, error) {
	var (
		scrModel   = make([]models.InterviewAverageScore, 0)
		scrStModel = make([]models.StudentSimpleData, 0)
		stData     = make(map[string][]models.StudentSimpleData)
		studentIds = []int{}
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	colStudents := db.Mongodb.Collection("students")
	col := db.Mongodb.Collection("interview_score")

	// Construct case-insensitive regular expression patterns for each email in the array
	regexPatterns := make([]primitive.Regex, len(email))
	for i, e := range email {
		pattern := fmt.Sprintf("(?i)%s", regexp.QuoteMeta(e))
		regexPatterns[i] = primitive.Regex{Pattern: pattern}
	}

	filStudents := bson.M{
		"$or": []bson.M{
			{"email": bson.M{"$in": regexPatterns}},
		},
	}

	cur, err := colStudents.Find(ctx, filStudents)
	if err != nil {
		return []models.InterviewScoreEmailEdutech{}, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.StudentSimpleData
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrStModel = append(scrStModel, model)
	}
	for _, dt := range scrStModel {
		studentIds = append(studentIds, dt.SmartbtwID)
		_, k := stData[strings.ToLower(dt.Email)]
		if !k {
			stData[strings.ToLower(dt.Email)] = []models.StudentSimpleData{}
		}
		res, err := GetStudentProfileElastic(dt.SmartbtwID)
		if err == nil {
			dt.SchoolName = res.SchoolNamePTK
			dt.SchoolID = res.SchoolPTKID
			dt.MajorName = res.MajorNamePTK
			dt.MajorID = res.MajorPTKID
			dt.OriginSchoolName = res.LastEdName
			dt.OriginSchoolID = res.LastEdID
			dt.AccountType = res.AccountType
		} else {
			dt.AccountType = "smartbtw"
		}
		stData[strings.ToLower(dt.Email)] = append(stData[strings.ToLower(dt.Email)], dt)
	}
	if year == 0 {
		year = uint16(time.Now().Year())
	}

	studentProfileData := []models.InterviewScoreEmailEdutech{}

	pipel := aggregates.GetInterviewAverageScoresByArrayOfStudentIDAndYear(studentIds, year)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	curBkn, err := col.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return []models.InterviewScoreEmailEdutech{}, err
	}
	defer curBkn.Close(ctx)

	for curBkn.Next(ctx) {
		var model models.InterviewAverageScore
		e := curBkn.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}

	for _, ems := range stData {
		stData := models.InterviewScoreEmailEdutech{}

		for _, valStd := range ems {
			if valStd.AccountType == "btwedutech" {
				stData.Name = valStd.Name
				stData.AccountType = "btwedutech"
				stData.BTWEdutechID = valStd.SmartbtwID
			} else {
				if stData.Name == "" {
					stData.Name = valStd.Name
				}
				stData.AccountType = "btwedutech"
				stData.SmartBtwID = valStd.SmartbtwID
			}

			for _, dt := range scrModel {
				ns := dt
				if valStd.SmartbtwID == dt.SmartBtwID {
					stData.InterviewScore = &ns
				}
			}
		}

		studentProfileData = append(studentProfileData, stData)
	}

	return studentProfileData, nil
}

func CreateInterviewScore(req *request.UpsertInterviewScore) (*mongo.InsertOneResult, error) {
	var (
		createdAt time.Time = time.Now()
	)
	col := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	interviewSession, err := GetSingleInterviewSessionByID(req.SessionID)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			errMessage := fmt.Sprintf("interview session data with ID: %s was not found", req.SessionID.Hex())
			customErr := errors.New(errMessage)
			return nil, customErr
		}
		return nil, err
	}

	payload := bson.M{
		"smartbtw_id":                           req.SmartBtwID,
		"penampilan":                            req.Penampilan,
		"cara_duduk_dan_berjabat":               req.CaraDudukDanBerjabat,
		"praktek_baris_berbaris":                req.PraktekBarisBerbaris,
		"penampilan_sopan_santun":               req.PenampilanSopanSantun,
		"kepercayaan_diri_dan_stabilitas_emosi": req.KepercayaanDiriDanStabilitasEmosi,
		"komunikasi":                            req.Komunikasi,
		"pengembangan_diri":                     req.PengembanganDiri,
		"integritas":                            req.Integritas,
		"kerjasama":                             req.Kerjasama,
		"mengelola_perubahan":                   req.MengelolaPerubahan,
		"perekat_bangsa":                        req.PerekatBangsa,
		"pelayanan_publik":                      req.PelayananPublik,
		"pengambilan_keputusan":                 req.PengambilanKeputusan,
		"orientasi_hasil":                       req.OrientasiHasil,
		"prestasi_akademik":                     req.PrestasiAkademik,
		"prestasi_non_akademik":                 req.PrestasiNonAkademik,
		"bahasa_asing":                          req.BahasaAsing,
		"final_score":                           req.FinalScore,
		"year":                                  req.Year,
		"session_id":                            req.SessionID,
		"session_name":                          interviewSession.Name,
		"session_description":                   interviewSession.Description,
		"session_number":                        interviewSession.Number,
		"bersedia_pindah_jurusan":               req.BersediaPindahJurusan,
		"closing_statement":                     req.ClosingStatement,
		"created_by": bson.M{
			"id":   req.CreatedBy.ID,
			"name": req.CreatedBy.Name,
		},
		"updated_by": nil,
		"created_at": createdAt,
		"note":       req.Note,
		"updated_at": time.Now(),
		"deleted_at": nil,
	}

	res, err := col.InsertOne(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func UpdateInterviewScore(interviewScoreID string, req *request.UpsertInterviewScore) error {
	var (
		createdAt         time.Time = time.Now()
		updatedByUserId   *string
		updatedByUserName *string
		createdByUserId   string
		createdByUserName string
	)
	col := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	interviewSession, err := GetSingleInterviewSessionByID(req.SessionID)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			errMessage := fmt.Sprintf("interview session data with ID: %s was not found", req.SessionID.Hex())
			customErr := errors.New(errMessage)
			return customErr
		}
		return err
	}

	interviewScoreObjectID, err := primitive.ObjectIDFromHex(interviewScoreID)
	if err != nil {
		return err
	}

	interviewScore, err := GetSingleInterviewScoreByID(interviewScoreObjectID)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}

	if (interviewScore == models.InterviewScore{}) {
		newErr := errors.New("interview score data not found")
		return newErr
	}

	if (interviewScore != models.InterviewScore{}) {
		createdAt = interviewScore.CreatedAt
		createdByUserId = interviewScore.CreatedBy.ID
		createdByUserName = interviewScore.CreatedBy.Name
		updatedByUserId = &req.CreatedBy.ID
		updatedByUserName = &req.CreatedBy.Name
	} else {
		createdByUserId = req.CreatedBy.ID
		createdByUserName = req.CreatedBy.Name
	}

	payload := bson.M{
		"smartbtw_id":                           req.SmartBtwID,
		"penampilan":                            req.Penampilan,
		"cara_duduk_dan_berjabat":               req.CaraDudukDanBerjabat,
		"praktek_baris_berbaris":                req.PraktekBarisBerbaris,
		"penampilan_sopan_santun":               req.PenampilanSopanSantun,
		"kepercayaan_diri_dan_stabilitas_emosi": req.KepercayaanDiriDanStabilitasEmosi,
		"komunikasi":                            req.Komunikasi,
		"pengembangan_diri":                     req.PengembanganDiri,
		"integritas":                            req.Integritas,
		"kerjasama":                             req.Kerjasama,
		"mengelola_perubahan":                   req.MengelolaPerubahan,
		"perekat_bangsa":                        req.PerekatBangsa,
		"pelayanan_publik":                      req.PelayananPublik,
		"pengambilan_keputusan":                 req.PengambilanKeputusan,
		"orientasi_hasil":                       req.OrientasiHasil,
		"prestasi_akademik":                     req.PrestasiAkademik,
		"prestasi_non_akademik":                 req.PrestasiNonAkademik,
		"bahasa_asing":                          req.BahasaAsing,
		"final_score":                           req.FinalScore,
		"year":                                  req.Year,
		"session_id":                            req.SessionID,
		"session_name":                          interviewSession.Name,
		"session_description":                   interviewSession.Description,
		"session_number":                        interviewSession.Number,
		"bersedia_pindah_jurusan":               req.BersediaPindahJurusan,
		"closing_statement":                     req.ClosingStatement,
		"created_by": bson.M{
			"id":   createdByUserId,
			"name": createdByUserName,
		},
		"created_at": createdAt,
		"note":       req.Note,
		"updated_at": time.Now(),
	}

	if updatedByUserId != nil && updatedByUserName != nil {
		payload["updated_by"] = bson.M{
			"id":   updatedByUserId,
			"name": updatedByUserName,
		}
	} else {
		payload["updated_by"] = nil
	}

	filter := bson.M{"_id": interviewScoreObjectID, "deleted_at": nil}
	update := bson.M{"$set": payload}

	_, err = col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func HardDeleteInterviewScore(interviewSessionID string) (*mongo.DeleteResult, error) {
	col := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	interviewScoreObjectID, err := primitive.ObjectIDFromHex(interviewSessionID)
	if err != nil {
		return nil, err
	}

	interviewScore, err := GetSingleInterviewScoreByID(interviewScoreObjectID)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
	}

	filter := bson.M{"_id": interviewScore.ID}
	res, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}
