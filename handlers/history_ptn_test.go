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

func TestCreateHistoryPtnSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  7999,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err1)

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 7999,
		"task_id": 1,
		"module_code": "M-001",
		"module_type": "uka_free",
		"potensi_kognitif": 500,
		"penalaran_matematika": 400,
		"literasi_bahasa_indonesia": 300,
		"literasi_bahasa_inggris": 300,
		"total": 800,
		"exam_name": "test",
		"grade": "basic",
		"repeat": 1,
		"created_at": "2022-01-10T07:53:16Z",
		"updated_at": "2022-01-10T07:53:16Z"
		}`)

	request, e := http.NewRequest(
		"POST",
		"/history-ptn/",
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

func TestUpdateHistoryPtnSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  79129,
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
		SmartBtwID:              79129,
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

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 79129,
		"task_id": 1,
		"module_code": "M-001",
		"module_type": "uka_free",
		"potensi_kognitif": 500,
		"penalaran_matematika": 400,
		"literasi_bahasa_indonesia": 300,
		"literasi_bahasa_inggris": 300,
		"total": 800,
		"repeat": 1,
		"exam_name": "test new",
		"grade": "basic"
		}`)

	url := fmt.Sprintf("/history-ptn/%s", res.InsertedID.(primitive.ObjectID).Hex())
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

func TestDeleteHistoryPtnSuccess(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  791444,
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
		SmartBtwID:              791444,
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

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/history-ptn/%s", res.InsertedID.(primitive.ObjectID).Hex()),
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

func TestGetHistoryPtnByIDSuccess(t *testing.T) {
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
	assert.Nil(t, err1)

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

	res, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/detail/%s", res.InsertedID.(primitive.ObjectID).Hex()),
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

func TestGetHistoryPtnBySmartBTWIDSuccess(t *testing.T) {
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
		fmt.Sprintf("/history-ptn/student/%d", payload.SmartBtwID),
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

func TestGetStudentPtnAverageSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  79333,
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
		SmartBtwID:              79333,
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
		fmt.Sprintf("/history-ptn/average/%d", payload.SmartBtwID),
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

func TestGetStudentPtnLastScoreSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  79078,
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
		SmartBtwID:              79078,
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
		fmt.Sprintf("/history-ptn/last-score/%d", payload.SmartBtwID),
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

func TestGetStudentPtnLast10ScoreSuccess(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  79008,
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
		SmartBtwID:              79008,
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
		fmt.Sprintf("/history-ptn/last-ten-score/%d", payload.SmartBtwID),
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

func Test_GetStudentFreePTNHistory(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  9444900001,
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
		SmartBtwID:              9444900001,
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

	_, err = lib.GetHistoryFreeSingleStudentPTN(9444900001)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-free/%d", payload.SmartBtwID),
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

func Test_GetStudentFreeNotFoundPTN(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-free/%d", 66666666),
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

func Test_GetStudentFreeErrorPTN(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-free/%s", "GANTENG"),
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

func Test_GetStudentPremiumUKAPTNHistory(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  9444900101,
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
		SmartBtwID:              9444900101,
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

	_, err = lib.GetHistoryPremiumUKASingleStudentPTN(9444900101)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-premium/%d", payload.SmartBtwID),
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

func Test_GetStudentPremiumUKANotFoundPTN(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-premium/%d", 7778889999),
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

func Test_GetStudentPremiumUKAErrorPTN(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-premium/%s", "GANTENG"),
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

func Test_GetStudentPackageUKAPTNHistory(t *testing.T) {
	Init()

	payload1 := request.CreateStudentTarget{
		SmartbtwID:  9444900102,
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
		SmartBtwID:              9444900102,
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

	_, err = lib.GetHistoryPackageUKASingleStudentPTN(9444900102)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-package/%d", payload.SmartBtwID),
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

func Test_GetStudentPackageUKANotFoundPTN(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-package/%d", 7778889999),
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

func Test_GetStudentPackageUKAErrorPTN(t *testing.T) {
	Init()
	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/history-ptn/student-package/%s", "GANTENG"),
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

func Test_AllStudentPtnScore(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  9999101,
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
		SmartBtwID:              9999101,
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
	body := []byte(
		`{
			"smartbtw_id": 9999101,
			"limit": 5,
			"page": 2
		}`)

	request, e := http.NewRequest(
		"POST",
		"/history-ptn/all-score",
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
