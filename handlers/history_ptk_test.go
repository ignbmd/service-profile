package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func TestCreateHistoryPtkHandlerSuccess(t *testing.T) {
	Init()

	payload := request.CreateStudentTarget{
		SmartbtwID:  666999,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}
	_, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 666999,
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
		}`)

	request, e := http.NewRequest(
		"POST",
		"/history-ptk/",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 200 created
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func TestUpdateHistoryPtkHandlerSuccess(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7749,
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
		SmartBtwID: 7749,
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

	res, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 7749,
		"task_id": 2,
		"module_code": "M-002",
		"module_type": "uka_free",
		"twk": 100,
		"tiu": 100,
		"tkp": 100,
		"total": 300,
		"repeat": 2,
		"exam_name": "test",
		"grade": "platinum"
	}`)

	url := fmt.Sprintf("/history-ptk/%s", res.InsertedID.(primitive.ObjectID).Hex())
	request, e := http.NewRequest(
		"PUT",
		url,
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 200 created
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func TestDeleteHistoryPtk(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7859,
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
		SmartBtwID: 7859,
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
		TargetID:   "63311dbf5f63b4d4885a6e4f",
	}

	res, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/history-ptk/%s", res.InsertedID.(primitive.ObjectID).Hex()),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetHistoryPtkByID(t *testing.T) {
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
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

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

	res, err := lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/detail/%s", res.InsertedID.(primitive.ObjectID).Hex()),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetHistoryPtkBySmartBTWID(t *testing.T) {
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
	_, err := lib.CreateStudentTarget(&payload1)
	if err != nil {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
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

	_, err = lib.CreateHistoryPtk(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student/%d", payload.SmartBtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetStudentAveragePtk(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7759,
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
		SmartBtwID: 7759,
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
		fmt.Sprintf("/history-ptk/average/%d", payload.SmartBtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetStudentLastScore(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7959,
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
		SmartBtwID: 7959,
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
		fmt.Sprintf("/history-ptk/last-score/%d", payload.SmartBtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetStudentPtkLast10ScoreSuccess(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7999,
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
		SmartBtwID: 7999,
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
		fmt.Sprintf("/history-ptk/last-ten-score/%d", payload.SmartBtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetStudentFreePTKHistory(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  66666661,
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
		SmartBtwID: 66666661,
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

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-free/%d", payload.SmartBtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetStudentFreeNotFound(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-free/%d", 66666666),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func Test_GetStudentFreeError(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-free/%s", "GANTENG"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_GetStudentPremiumUKAPTKHistory(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  66666662,
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
		SmartBtwID: 66666662,
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

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-premium/%d", payload.SmartBtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetStudentPremiumUKANotFound(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-premium/%d", 66666666),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func Test_GetStudentPremiumUKAError(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-premium/%s", "GANTENG"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_GetStudentPackageUKAPTKHistory(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  66666665,
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
		SmartBtwID: 66666665,
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

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-package/%d", payload.SmartBtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetStudentPackageUKANotFound(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-package/%d", 66666666),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func Test_GetStudentPackageUKAError(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptk/student-package/%s", "GANTENG"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_AllStudentPtkScore(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  88668865,
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
		SmartBtwID: 88668865,
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
	app := server.SetupFiber()

	body := []byte(
		`{
			"smartbtw_id": 88668865,
			"limit": 5,
			"page": 1
		}`)

	request, e := http.NewRequest(
		"POST",
		"/history-ptk/all-score",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}
