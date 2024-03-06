package handlers_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func TestGetHistoryScoreSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  79111,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	if err1 != nil {
		assert.NotNil(t, err1)
	} else {
		assert.Nil(t, err1)
	}

	payload := request.CreateHistoryPtn{
		SmartBtwID:              79111,
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

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-scores/%s?smartbtw_id=%d&task_id=%d", "ptn", 79111, 2),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")
	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)
	assert.Nil(t, err)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var jBody interface{}
	err = sonic.Unmarshal(body, &jBody)
	assert.Equal(t, nil, err)
	assert.NotNil(t, jBody)

	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetHistoryScoreInvalidTargetType(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7259,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	if err1 != nil {
		assert.NotNil(t, err1)
	} else {
		assert.Nil(t, err1)
	}

	payload := request.CreateHistoryPtk{
		SmartBtwID: 7259,
		TaskID:     1,
		ModuleCode: "M-001",
		ModuleType: string(models.UkaFree),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "test",
		Grade:      string(models.Basic),
	}

	_, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-scores/%s?smartbtw_id=%d&task_id=%d", "nact", 79111, 2),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")
	assert.Equal(t, nil, e)
	response, err := app.Test(request, -1)
	assert.Equal(t, nil, err)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var jBody interface{}
	err = sonic.Unmarshal(body, &jBody)
	assert.Equal(t, nil, err)
	assert.NotNil(t, jBody)

	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func TestGetLatestHistorySuccess(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7259,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	if err1 != nil {
		assert.NotNil(t, err1)
	} else {
		assert.Nil(t, err1)
	}

	payload := request.CreateHistoryPtk{
		SmartBtwID: 7259,
		TaskID:     1,
		ModuleCode: "M-001",
		ModuleType: string(models.UkaFree),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "test",
		Grade:      string(models.Basic),
	}

	_, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-scores/%s?smartbtw_id=%d&task_id=%d&only_latest_score=%v", "ptk", 7259, 1, true),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")
	assert.Equal(t, nil, e)
	response, err := app.Test(request, -1)
	assert.Equal(t, nil, err)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var jBody interface{}
	err = sonic.Unmarshal(body, &jBody)
	assert.Equal(t, nil, err)
	assert.NotNil(t, jBody)

	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetLatestHistoryScoreIncompleteParams(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7259,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	if err1 != nil {
		assert.NotNil(t, err1)
	} else {
		assert.Nil(t, err1)
	}

	payload := request.CreateHistoryPtk{
		SmartBtwID: 7259,
		TaskID:     1,
		ModuleCode: "M-001",
		ModuleType: string(models.UkaFree),
		Twk:        100,
		Tiu:        100,
		Tkp:        100,
		Total:      300,
		Repeat:     1,
		ExamName:   "test",
		Grade:      string(models.Basic),
	}

	_, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-scores/%s?only_latest_score=%v", "ptk", true),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")
	assert.Equal(t, nil, e)
	response, err := app.Test(request, -1)
	assert.Equal(t, nil, err)

	body, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	var jBody interface{}
	err = sonic.Unmarshal(body, &jBody)
	assert.Equal(t, nil, err)
	assert.NotNil(t, jBody)

	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}
