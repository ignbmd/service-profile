package lib_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestGetHistoryScoreSuccess(t *testing.T) {
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

	sID := 943439
	tID := 2

	query := request.HistoryScoreQueryParams{
		SmartBTWID: &sID,
		TaskID:     &tID,
	}
	result, err1 := lib.GetHistoryScoreByTargetType("ptn", &query)
	fmt.Println(result)
	assert.Nil(t, err1)
}

func TestGetHistoryScoreInvalidTargetType(t *testing.T) {
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

	sID := 943439
	tID := 2

	query := request.HistoryScoreQueryParams{
		SmartBTWID: &sID,
		TaskID:     &tID,
	}
	result, err1 := lib.GetHistoryScoreByTargetType("nact", &query)
	fmt.Println(result)
	assert.Len(t, result, 0)
	assert.NotNil(t, err1)
}
