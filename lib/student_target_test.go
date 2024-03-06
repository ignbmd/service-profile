package lib_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestCreateStudentTargetSuccess(t *testing.T) {
	Init()

	json := `{
		"version" : 2,
		"data": {
			"smartbtw_id": 12432,
			"school_id": 1,
			"major_id":1,
			"school_name": "UI",
			"major_name": "Kedokteran",
			"target_score": 700,
			"target_type": "PTK"
		}
	}`

	body := []byte(json)
	payload, _ := request.UnmarshalMessageStudentTargetBody(body)
	_, err := lib.CreateStudentTarget(&payload.Data)
	assert.Nil(t, err)
}

func TestDeleteStudentTargetSuccess(t *testing.T) {
	Init()
	payload := request.CreateStudentTarget{
		SmartbtwID:  5559,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	res, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	err = lib.DeleteStudentTarget(res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}

func TestUpdateStudentTargetByIDSuccess(t *testing.T) {
	Init()

	payload := request.CreateStudentTarget{
		SmartbtwID:  4739,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	res, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)
	fmt.Println(res.InsertedID.(primitive.ObjectID))

	payload1 := request.UpdateStudentTarget{
		SmartbtwID:  4739,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "Teknik Sipil",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	err = lib.UpdateStudentTargetByID(&payload1, res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}

func TestUpdateStudentTargetBySmartbtwIDSuccess(t *testing.T) {
	Init()

	payload := request.CreateStudentTarget{
		SmartbtwID:  99999999,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	res, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)
	fmt.Println(res.InsertedID.(primitive.ObjectID))

	payload1 := request.UpdateStudentTarget{
		SmartbtwID:  99999999,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "Teknik Sipil",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	err = lib.UpdateStudentTargetBySmartbtwID(&payload1, payload.SmartbtwID)
	assert.Nil(t, err)
}

func TestGetStudentTargetByID(t *testing.T) {
	Init()

	payload := request.CreateStudentTarget{
		SmartbtwID:  5189,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	res, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	_, err1 := lib.GetStudentTargetByID(res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err1)
}

func Test_CreateStudentTargetPtkElastic(t *testing.T) {
	Init()

	body := request.StudentTargetPtkElastic{
		SmartbtwID:  987646,
		Name:        "Bambang",
		Photo:       "https://test.png",
		ModuleDone:  12,
		SchoolName:  "STAN",
		SchoolID:    1,
		MajorName:   "test",
		MajorID:     1,
		TargetScore: 500,
		TargetType:  "PTK",
	}
	err := lib.InsertStudentTargetPtkElastic(&body)
	assert.Nil(t, err)
}

func TestGetStudentBySmartbtwIDandTargetType(t *testing.T) {
	Init()
	body := request.StudentTargetPtkElastic{
		SmartbtwID:  22336655,
		Name:        "Bambang",
		Photo:       "https://test.png",
		ModuleDone:  12,
		SchoolName:  "STAN",
		SchoolID:    1,
		MajorName:   "test",
		MajorID:     1,
		TargetScore: 500,
		TargetType:  "PTK",
	}
	err := lib.InsertStudentTargetPtkElastic(&body)
	assert.Nil(t, err)

	res, err := lib.GetStudentTargetElastic(body.SmartbtwID, "PTK", "")
	assert.Nil(t, err)
	fmt.Println(res)
}

func TestUpdateUserData(t *testing.T) {
	Init()

	req := request.UpdateUserData{
		Name:  "Siti",
		Photo: "test",
	}

	err := lib.UpdateUserData(&req, 65113, context.Background())
	assert.Nil(t, err)
}

func TestUpdateSchool(t *testing.T) {
	Init()
	payload := request.CreateStudentTarget{
		SmartbtwID:  8654445,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	req := request.UpdateSchool{
		SchoolID:   payload.SchoolID,
		SchoolName: "IPDN",
	}

	err1 := lib.UpdateSchool(&req, payload.SchoolID, string(models.PTN), context.Background())
	assert.Nil(t, err1)
}

func TestUpdateStudyProgram(t *testing.T) {
	Init()
	payload := request.CreateStudentTarget{
		SmartbtwID:  86544675,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "STAN",
		MajorName:   "STATISTIK",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	req := request.UpdateStudyProgram{
		MajorID:   payload.MajorID,
		MajorName: "NAUTIKA",
	}

	err1 := lib.UpdateStudyProgram(&req, payload.MajorID, string(models.PTK), context.Background())
	assert.Nil(t, err1)
}

func TestUpdatePolbit(t *testing.T) {
	Init()
	payload := request.CreateStudentTarget{
		SmartbtwID:  865446725,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "STAN",
		MajorName:   "STATISTIK",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	req := request.UpdatePolbitType{
		MajorID:             payload.MajorID,
		SmartBTWID:          865446725,
		PolbitType:          "DAERAH_PROVINCE",
		PolbitCompetitionID: 1,
		PolbitLocationID:    1,
		TargetScore:         650,
		TargetType:          string(models.PTK),
	}

	err1 := lib.UpdateStudentPolbit(&req, context.Background())
	assert.Nil(t, err1)
}

func TestUpdateTargetScore(t *testing.T) {
	Init()
	payload := request.CreateStudentTarget{
		SmartbtwID:  86500675,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "STAN",
		MajorName:   "STATISTIK",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	req := request.UpdateTargetScore{
		MajorID:     payload.MajorID,
		TargetScore: 600,
	}

	err1 := lib.UpdateTargetScore(&req, payload.MajorID, string(models.PTK), context.Background())
	assert.Nil(t, err1)
}

func TestGetStudentTargetByCustomSuccess(t *testing.T) {
	Init()

	json := `{
		"version" : 2,
		"data": {
			"smartbtw_id": 123123123,
			"school_id": 1,
			"major_id":1,
			"school_name": "UI",
			"major_name": "Kedokteran",
			"target_score": 700,
			"target_type": "PTK"
		}
	}`

	body := []byte(json)
	payload, _ := request.UnmarshalMessageStudentTargetBody(body)
	_, err := lib.CreateStudentTarget(&payload.Data)
	assert.Nil(t, err)

	get, err := lib.GetStudentTargetByCustom(123123123, "PTK")
	assert.Nil(t, err)
	fmt.Println(get)
}

func TestGetStudentTargetByCustomFailedNotFound(t *testing.T) {
	Init()

	_, err := lib.GetStudentTargetByCustom(12312312333, "PTK")
	assert.NotNil(t, err)
}

func TestNewUpdateStudentTarget(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  9812333,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

	payload := []*request.NewUpdateStudentTarget{
		{
			SmartbtwID:  9812333,
			SchoolID:    4,
			MajorID:     4,
			SchoolName:  "UGM",
			MajorName:   "AKUNTANSI",
			TargetScore: 700,
			TargetType:  string(models.PTN),
			Position:    1,
			Type:        string(models.SECONDARY),
		},
		{
			SmartbtwID:  9812333,
			SchoolID:    8,
			MajorID:     2,
			SchoolName:  "UI",
			MajorName:   "KEDOKTERAN",
			TargetScore: 700,
			TargetType:  string(models.PTN),
			Position:    2,
			Type:        string(models.PRIMARY),
		},
	}

	err1 := lib.UpdateStudentTarget(payload, context.Background())
	assert.Nil(t, err1)
}

func TestSycnStudyProgramSSN(t *testing.T) {
	Init()
	md := []int{311, 312, 313}

	py := request.UpdateSpecificStudyProgram{
		MajorID:    md,
		TargetType: "PTK",
		MajorName:  "Keamanan Siber 1, 2, 3",
		NewMajorID: 311,
	}
	err := lib.UpdateSpecificStudyProgram(context.Background(), &py)
	assert.Nil(t, err)
}
