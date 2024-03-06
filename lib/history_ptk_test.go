package lib_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestCreateHisoryPtk(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  5579,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	json := `{
	"version": 2,
	"data" : {
		"smartbtw_id": 5579,
		"task_id": 1,
		"module_code": "M-001",
		"module_type": "uka_free",
		"twk": 500,
		"tiu": 400,
		"tkp": 300,
		"total": 800,
		"repeat": 1,
		"exam_name": "test",
		"grade": "basic",
		"created_at": "2022-01-10T07:53:16Z",
		"updated_at": "2022-01-10T07:53:16Z"
	}
}`
	body := []byte(json)
	payload, _ := request.UnmarshalMessageHistoryPtkBody(body)
	_, err := lib.CreateHistoryPtk(&payload.Data)
	assert.Nil(t, err)
}

func TestUpsertHisoryPtk(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  55822,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	json := `{
	"version": 2,
	"data" : {
		"smartbtw_id": 55822,
		"task_id": 1,
		"module_code": "M-001",
		"module_type": "uka_free",
		"twk": 500,
		"tiu": 400,
		"tkp": 300,
		"total": 800,
		"repeat": 1,
		"exam_name": "test",
		"grade": "basic",
		"created_at": "2022-01-10T07:53:16Z",
		"updated_at": "2022-01-10T07:53:16Z",
		"is_live": true
	}
}`
	body := []byte(json)
	payload, _ := request.UnmarshalMessageHistoryPtkBody(body)
	_, err := lib.UpsertHistoryPtk(&payload.Data)
	assert.Nil(t, err)
}

func TestUpdateHistoryPtkSuccess(t *testing.T) {
	Init()
	payload2 := request.CreateStudentTarget{
		SmartbtwID:  55009,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 55009,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
	}
	res, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	payload1 := request.UpdateHistoryPtk{
		SmartBtwID: 5579,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     2,
		ExamName:   "testing",
		Grade:      string(models.Gold),
	}
	err = lib.UpdateHistoryPtk(&payload1, res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}

func TestDeleteHistoryPtkSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  51119,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 51119,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
		TargetID:   "633115f63b4d4886e4f",
	}
	res, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)
	err1 := lib.DeleteHistoryPtk(res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err1)
}

func TestGetHistoryPtkByIDSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  53339,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 53339,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
		TargetID:   "63311dbf5f63b4d4885a6e4f",
	}
	res, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	_, err = lib.GetHistoryPtkByID(res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}

func TestGetHistoryPtkBySmartBTWIDSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  53339,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	if err != nil {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}

	payload := request.CreateHistoryPtk{
		SmartBtwID: 53339,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
		TargetID:   "63311dbf5f63b4d4885a6e4f",
	}
	_, err = lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	params := new(request.HistoryPTKQueryParams)
	_, err = lib.GetHistoryPtkBySmartBTWID(payload.SmartBtwID, params)
	assert.Nil(t, err)
}

func TestGetStudentAverage(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  93339,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 93339,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      380,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
	}
	_, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	_, err1 = lib.GetStudentAveragePtk(93339)
	assert.Nil(t, err1)
}

func TestGetStudentLastScore(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  92229,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 92229,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      777,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
	}
	_, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	res, err1 := lib.GetStudentLastScore(92229)
	fmt.Println(res)
	assert.Nil(t, err1)
}

func TestGetLast10RecordStudentPtk(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  94449,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 94449,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      7977,
		Repeat:     1,
		ExamName:   "new exam2",
		Grade:      string(models.Basic),
	}
	_, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	res, err1 := lib.GetLast10StudentScorePtk(94449)
	fmt.Println(res)
	assert.Nil(t, err1)
}

func TestGetAllRecordStudentPTK(t *testing.T) {
	Init()

	res, err1 := lib.GetALLStudentScorePtk(94449)
	fmt.Printf("%v", res)
	assert.Nil(t, err1)
}

func Test_CreateStudentPtkProfileElastic(t *testing.T) {
	Init()

	body := request.StudentProfilePtkElastic{
		Name:                    "Siswa Testing",
		Photo:                   "http://google.com",
		TargetType:              "PTK",
		SchoolID:                1,
		SchoolName:              "Testing School",
		MajorID:                 2,
		MajorName:               "Testing Major",
		TargetScore:             400,
		SmartbtwID:              500,
		TwkAvgScore:             100,
		TiuAvgScore:             100,
		TkpAvgScore:             100,
		TwkAvgPercentScore:      66.67,
		TiuAvgPercentScore:      57.14,
		TkpAvgPercentScore:      44.44,
		TotalAvgPercentScore:    44.44,
		TotalAvgScore:           100,
		LatestTotalScore:        100,
		LatestTotalPercentScore: 25,
		ModuleDone:              2,
	}
	indexID := fmt.Sprintf("%d_PTK", body.SmartbtwID)
	err := lib.InsertStudentPtkProfileElastic(&body, indexID)
	assert.Nil(t, err)
}

func Test_CreateStudentPtkHistoryElastic(t *testing.T) {
	Init()

	body := request.CreateHistoryPtk{
		SmartBtwID: 500,
		TaskID:     1,
		ModuleCode: "MD-100",
		ModuleType: "TRIAL",
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		ExamName:   "Coba Modul PTK",
		Grade:      "PLATINUM",
		TargetID:   "633ba0784b72a120947073bd",
	}
	objID := primitive.NewObjectID().Hex()
	err := lib.InsertStudentHistoryPtkElastic(&body, objID)
	assert.Nil(t, err)
}

func Test_GetHistoryFreeSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  6666666,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 6666666,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaFree),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
		TargetID:   "63311dbf5f63b4d4885a6e4f",
	}
	_, err = lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	res, err := lib.GetHistoryFreeSingleStudentPTK(6666666)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_GetHistoryUKAEmpty(t *testing.T) {
	Init()

	res, err := lib.GetHistoryFreeSingleStudentPTK(6666666666)
	assert.NotNil(t, err)
	assert.Equal(t, mongo.ErrNoDocuments, err)
	fmt.Println(err)
	fmt.Println(res)
}

func Test_GetHistoryPremiumUKASuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  6666666999,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 6666666999,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.UkaPremium),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
		TargetID:   "63311dbf5f63b4d4885a6e4f",
	}
	_, err = lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	res, err := lib.GetHistoryPremiumUKASingleStudentPTK(6666666999)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_GetHistoryPremiumUKAEmpty(t *testing.T) {
	Init()

	res, err := lib.GetHistoryPremiumUKASingleStudentPTK(6666666666)
	assert.NotNil(t, err)
	assert.Equal(t, mongo.ErrNoDocuments, err)
	fmt.Println(err)
	fmt.Println(res)
}

func Test_GetHistoryPackageUKASuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  6666667,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

	payload := request.CreateHistoryPtk{
		SmartBtwID: 6666667,
		TaskID:     2,
		ModuleCode: "M-002",
		ModuleType: string(models.Package),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "testing",
		Grade:      string(models.Basic),
		TargetID:   "63311dbf5f63b4d4885a6e4f",
	}
	_, err = lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	res, err := lib.GetHistoryPackageUKASingleStudentPTK(6666667)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_GetHistoryPackageUKAEmpty(t *testing.T) {
	Init()

	res, err := lib.GetHistoryPackageUKASingleStudentPTK(6666666699)
	assert.NotNil(t, err)
	assert.Equal(t, mongo.ErrNoDocuments, err)
	fmt.Println(err)
	fmt.Println(res)
}

func TestGetStudentHistoryPTKOnlyStage(t *testing.T) {
	Init()

	// payload1 := request.CreateStudentTarget{
	// 	SmartbtwID:  6666666,
	// 	SchoolID:    1,
	// 	MajorID:     1,
	// 	SchoolName:  "UI",
	// 	MajorName:   "KEDOKTERAN",
	// 	TargetScore: 700,
	// 	TargetType:  string(models.PTK),
	// }
	// _, err := lib.CreateStudentTarget(&payload1)
	// assert.Nil(t, err)

	// payload := request.CreateHistoryPtk{
	// 	SmartBtwID: 6666666,
	// 	TaskID:     2,
	// 	ModuleCode: "M-002",
	// 	ModuleType: string(models.UkaFree),
	// 	Twk:        100,
	// 	Tiu:        100,
	// 	Tkp:        100,
	// 	Total:      300,
	// 	Repeat:     1,
	// 	ExamName:   "testing",
	// 	Grade:      string(models.Basic),
	// 	TargetID:   "63311dbf5f63b4d4885a6e4f",
	// }
	// _, err = lib.CreateHistoryPtk(&payload)
	// assert.Nil(t, err)

	res, err := lib.GetStudentHistoryPTKOnlyStage()
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_UpdateDuration(t *testing.T) {
	Init()
	payload2 := request.CreateStudentTarget{
		SmartbtwID:  550092,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "POLTEKIM",
		MajorName:   "KELAUTAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtk{
		SmartBtwID:    550092,
		TaskID:        2,
		ModuleCode:    "M-002",
		ModuleType:    string(models.UkaPremium),
		Twk:           100,
		Tiu:           100,
		Tkp:           100,
		Total:         300,
		Repeat:        1,
		ExamName:      "testing",
		Grade:         string(models.Basic),
		PackageType:   "testing",
		TwkPass:       200,
		TiuPass:       200,
		TkpPass:       200,
		TwkPassStatus: true,
		TiuPassStatus: true,
		TkpPassStatus: true,
		AllPassStatus: true,
	}
	_, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	payload1 := request.BodyUpdateStudentDuration{
		SmartbtwID: 550092,
		TaskID:     2,
		Repeat:     1,
		Program:    "skd",
		Start:      time.Now(),
		End:        time.Now(),
	}
	err = lib.UpdateDurationtPTK(&payload1)
	assert.Nil(t, err)
}
