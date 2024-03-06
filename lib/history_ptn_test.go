package lib_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestCreateHistoryPtnSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  112,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              112,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)
}

func TestUpdateHistoryPtnSuccess(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  912349,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              912349,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	res, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	payload1 := request.UpdateHistoryPtn{
		SmartBtwID:              912349,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		ExamName:                "test exam update",
		Grade:                   string(models.Basic),
	}

	err = lib.UpdateHistoryPtn(&payload1, res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}

func TestDeleteHistoryPtnSuccess(t *testing.T) {
	Init()
	payload2 := request.CreateStudentTarget{
		SmartbtwID:  943219,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              943219,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	res, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	err = lib.DeleteHistoryPtn(res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}

func TestGetHistoryPtnByID(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  943439,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              943439,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	res, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	result, err1 := lib.GetHistoryPtnByID(res.InsertedID.(primitive.ObjectID))
	fmt.Println(result)
	assert.Nil(t, err1)
}

func TestGetHistoryPtnBySmartBTWID(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  943439,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	if err1 != nil {
		assert.NotNil(t, err1)
	} else {
		assert.Nil(t, err1)
	}

	payload := request.CreateHistoryPtn{
		SmartBtwID:              943439,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	params := new(request.HistoryPTNQueryParams)
	result, err1 := lib.GetHistoryPtnBySmartBTWID(payload.SmartBtwID, params)
	fmt.Println(result)
	assert.Nil(t, err1)
}

func TestGetLastScorePtn(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  943229,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              943229,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	result, err1 := lib.GetStudentPtnLastScore(payload.SmartBtwID, "")
	fmt.Println(result)
	assert.Nil(t, err1)
}

func TestGetLast10ScorePtn(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  943228,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              943228,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	result, err1 := lib.GetLast10StudentScorePtn(payload.SmartBtwID, "")
	fmt.Println(result)
	assert.Nil(t, err1)
}

func TestGetAveragePtn(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  943118,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              943118,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	result, err1 := lib.GetStudentPtnAverage(payload.SmartBtwID, "")
	fmt.Println(result)
	assert.Nil(t, err1)
}

func TestGetAllRecordStudentPTN(t *testing.T) {
	Init()

	res, err1 := lib.GetALLStudentScorePtn(500)
	fmt.Printf("%v", res)
	assert.Nil(t, err1)
}

func Test_CreateStudentPtnProfileElastic(t *testing.T) {
	Init()

	body := request.StudentProfilePtnElastic{
		SmartbtwID:              500,
		Name:                    "Testing Student",
		Photo:                   "http://www.twitter.com",
		SchoolID:                10,
		SchoolName:              "Testing School",
		MajorID:                 11,
		MajorName:               "Testing Major",
		TargetType:              "PTN",
		TargetScore:             400,
		PkAvgScore:              100,
		PmAvgScore:              100,
		LbindAvgScore:           100,
		LbingAvgScore:           100,
		PkAvgPercentScore:       10,
		PmAvgPercentScore:       10,
		LbindAvgPercentScore:    10,
		LbingAvgPercentScore:    10,
		TotalAvgPercentScore:    0,
		TotalAvgScore:           0,
		LatestTotalScore:        100,
		LatestTotalPercentScore: 25,
		ModuleDone:              2,
	}
	indexID := fmt.Sprintf("%d_PTN", body.SmartbtwID)
	err := lib.InsertStudentPtnProfileElastic(&body, indexID)
	assert.Nil(t, err)
}

func Test_CreateStudentPtnHistoryElastic(t *testing.T) {
	Init()

	body := request.CreateHistoryPtn{
		SmartBtwID:              500,
		TaskID:                  1,
		ModuleCode:              "MD-102",
		ModuleType:              "TRIAL",
		PotensiKognitif:         100,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 100,
		LiterasiBahasaInggris:   100,
		Total:                   400,
		Repeat:                  1,
		ExamName:                "Coba Modul PTN",
		Grade:                   "NONE",
		TargetID:                "633ba06dd04fc827027ed1ad",
	}
	objID := primitive.NewObjectID().Hex()
	err := lib.InsertStudentHistoryPtnElastic(&body, objID)
	assert.Nil(t, err)
}

func Test_GetHistoryFreePTNSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  944490000,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              944490000,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	res, err := lib.GetHistoryFreeSingleStudentPTN(944490000)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_GetHistoryUKAFreePTNEmpty(t *testing.T) {
	Init()

	res, err := lib.GetHistoryFreeSingleStudentPTN(6666666666)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("mongo: no documents in result"), err)
	fmt.Println(err)
	fmt.Println(res)
}

func Test_GetHistoryPremiumUKAPTNSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  944490010,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              944490010,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaPremium),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	res, err := lib.GetHistoryPremiumUKASingleStudentPTN(944490010)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_GetHistoryPremiumUKAPTNEmpty(t *testing.T) {
	Init()

	res, err := lib.GetHistoryPremiumUKASingleStudentPTN(6666666666)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("mongo: no documents in result"), err)
	fmt.Println(err)
	fmt.Println(res)
}

func Test_GetHistoryPackageUKAPTNSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  944490011,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              944490011,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.Package),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	res, err := lib.GetHistoryPackageUKASingleStudentPTN(944490011)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_GetHistoryPackageUKAPTNEmpty(t *testing.T) {
	Init()

	res, err := lib.GetHistoryPackageUKASingleStudentPTN(6666666666)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Errorf("mongo: no documents in result"), err)
	fmt.Println(err)
	fmt.Println(res)
}

func Test_UpdateDurationPTN(t *testing.T) {
	Init()
	payload2 := request.CreateStudentTarget{
		SmartbtwID:  9123491,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              9123491,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	payload1 := request.BodyUpdateStudentDuration{
		SmartbtwID: 9123491,
		TaskID:     2,
		Repeat:     1,
		Program:    "tps",
		Start:      time.Now(),
		End:        time.Now(),
	}
	err = lib.UpdateDurationtPTN(&payload1)
	assert.Nil(t, err)
}
