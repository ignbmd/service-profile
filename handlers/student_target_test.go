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

func TestCreateStudentTargetSuccess(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 575,
		"school_id": 6,
		"major_id":1,
		"school_name": "UI",
		"major_name": "Kedokteran",
		"target_score": 700,
		"target_type": "PTN"
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-target",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func TestUpdateStudentTargetByIDSuccess(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  56565656,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "Kedokteran",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}

	res, err := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err)
	fmt.Println(res.InsertedID.(primitive.ObjectID))
	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 56565656,
		"school_id": 6,
		"major_id":1,
		"school_name": "UI",
		"major_name": "Teknik",
		"target_score": 700,
		"target_type": "PTN"
	}`)
	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/student-target/by-id/%s", res.InsertedID.(primitive.ObjectID).Hex()),
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func TestDeleteStudentTarget(t *testing.T) {
	Init()

	payload := request.CreateStudentTarget{
		SmartbtwID:  875,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "Kedokteran",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}

	res, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)
	fmt.Println(res.InsertedID.(primitive.ObjectID))

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/student-target/%s", res.InsertedID.(primitive.ObjectID).Hex()),
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

func TestGetStudentTargetByID(t *testing.T) {
	Init()

	payload := request.CreateStudentTarget{
		SmartbtwID:  999199,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "Kedokteran",
		TargetScore: 700,
		TargetType:  string(models.PTK),
	}

	res, err := lib.CreateStudentTarget(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-target/detail/%s", res.InsertedID.(primitive.ObjectID).Hex()),
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

func TestUpdateStudentTargetBySmartbtwIDSuccess(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  44447443,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "Kedokteran",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}

	res, err := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err)
	fmt.Println(res.InsertedID.(primitive.ObjectID))
	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 44447443,
		"school_id": 6,
		"major_id":1,
		"school_name": "UI",
		"major_name": "Teknik",
		"target_score": 700,
		"target_type": "PTN"
	}`)
	request, e := http.NewRequest(
		"PUT",
		"/student-target/by-student",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func TestGetStudentTargetElastic(t *testing.T) {
	Init()

	body1 := request.StudentTargetPtkElastic{
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
	err := lib.InsertStudentTargetPtkElastic(&body1)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 22336655,
		"target_type": "PTK"
	}`)
	request, e := http.NewRequest(
		"GET",
		"/student-target/elastic",
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
