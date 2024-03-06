package handlers_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

// success
func Test_CreateStudentModuleProgress_Success(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 1121,
		"task_id": 2,
		"module_no": 35,
		"repeat": 2,
		"module_total": 35
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-module-progress/",
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

func Test_CreateStudentModuleProgress_AlreadyExist(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  59,
		TaskID:      1,
		ModuleNo:    28,
		Repeat:      1,
		ModuleTotal: 28,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 59,
		"task_id": 1,
		"module_no": 28,
		"repeat": 1,
		"module_total": 28
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-module-progress/",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 422 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_CreateStudentModuleProgress_ModuleNoGreaterThanModuleTotal(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 511,
		"task_id": 2,
		"module_no": 35,
		"repeat": 2,
		"module_total": 5
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-module-progress/",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 200 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_CreateStudentModuleProgress_ErrorBodyNull(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 500,
		"task_id": null,
		"module_no": null,
		"repeat": null,
		"module_total": null
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-module-progress/",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 422 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_CreateStudentModuleProgress_ErrorBodyString(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 500,
		"task_id": 1,
		"module_no": "20",
		"repeat": "1",
		"module_total": "25"
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-module-progress/",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 422 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_CreateStudentModuleProgress_ErrorNegativeValue(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 600,
		"task_id": -1,
		"module_no": -20,
		"repeat": -1,
		"module_total": -25
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-module-progress/",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 422 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_CreateStudentModuleProgress_ErrorZeroValue(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(
		`{
		"smartbtw_id": 0,
		"task_id": 0,
		"module_no": 0,
		"repeat": 0,
		"module_total": 0
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-module-progress/",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 422 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_UpdateStudentModuleProgress_Success(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  509,
		TaskID:      1,
		ModuleNo:    30,
		Repeat:      1,
		ModuleTotal: 38,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		"task_id": 2,
		"module_no": 35,
		"repeat": 1,
		"module_total": 38
	}`)

	url := fmt.Sprintf("/student-module-progress/%s", res.InsertedID.(primitive.ObjectID).Hex())
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

func Test_UpdateStudentModuleProgress_NotFound(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  511,
		TaskID:      1,
		ModuleNo:    30,
		Repeat:      1,
		ModuleTotal: 38,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		"task_id": 2,
		"module_no": 35,
		"repeat": 1,
		"module_total": 39
	}`)

	url := fmt.Sprintf("/student-module-progress/%s", "62b442ab0f7f1f46df2baab9")
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

	// We expected when post some data to api got 404 created
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func Test_UpdateStudentModuleProgress_BodyNull(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  4900,
		TaskID:      1,
		ModuleNo:    30,
		Repeat:      1,
		ModuleTotal: 38,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		
	}`)

	url := fmt.Sprintf("/student-module-progress/%s", res.InsertedID.(primitive.ObjectID).Hex())
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

	// We expected when post some data to api got 422 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_UpdateStudentModuleProgress_FailedIDConvertID(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  681,
		TaskID:      1,
		ModuleNo:    28,
		Repeat:      1,
		ModuleTotal: 28,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		"task_id": 2,
		"module_no": 35,
		"repeat": 1,
		"module_total": 40
	}`)

	request, e := http.NewRequest(
		"PUT",
		"/student-module-progress/ajdkajshfahaof",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func Test_UpdateStudentModuleProgress_BodyNegativeValue(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  31,
		TaskID:      1,
		ModuleNo:    30,
		Repeat:      1,
		ModuleTotal: 38,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		"task_id": -2,
		"module_no": -35,
		"repeat": -1,
		"module_total": -4
	}`)

	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/student-module-progress/%s", res.InsertedID.(primitive.ObjectID).Hex()),
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_UpdateStudentModuleProgress_ModuleNoGreaterThanModuleTotal(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  49,
		TaskID:      1,
		ModuleNo:    30,
		Repeat:      1,
		ModuleTotal: 38,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	body := []byte(`{
		"task_id": 2,
		"module_no": 35,
		"repeat": 1,
		"module_total": 4
	}`)

	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/student-module-progress/%s", res.InsertedID.(primitive.ObjectID).Hex()),
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

// success
func Test_GetStudentModuleProgressBySmartBtwIdSuccess(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := []request.StudentModuleProgress{
		{
			SmartBtwID:  130009,
			TaskID:      1,
			ModuleNo:    1,
			Repeat:      2,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130010,
			TaskID:      1,
			ModuleNo:    8,
			Repeat:      2,
			ModuleTotal: 30,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130011,
			TaskID:      3,
			ModuleNo:    10,
			Repeat:      3,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, payload := range payload {
		filter := bson.M{"smartbtw_id": payload.SmartBtwID}
		update := bson.M{"$set": payload}
		stdCol.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-module-progress/detail/%s", "130011"),
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

//Data Not found
func Test_GetStudentModuleProgressBySmartBtwIdNotFound(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := []request.StudentModuleProgress{
		{
			SmartBtwID:  130009,
			TaskID:      1,
			ModuleNo:    1,
			Repeat:      2,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130010,
			TaskID:      1,
			ModuleNo:    8,
			Repeat:      2,
			ModuleTotal: 30,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130011,
			TaskID:      3,
			ModuleNo:    10,
			Repeat:      3,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, payload := range payload {
		filter := bson.M{"smartbtw_id": payload.SmartBtwID}
		update := bson.M{"$set": payload}
		stdCol.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-module-progress/detail/%s", "1"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	assert.Equalf(t, fiber.StatusNotFound, response.StatusCode, "Get student module progress data with nonexistent ID must return 404 status code")
}

//Failed Convert to int
func Test_GetStudentModuleProgressBySmartBtwIdError(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := []request.StudentModuleProgress{
		{
			SmartBtwID:  130009,
			TaskID:      1,
			ModuleNo:    1,
			Repeat:      2,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130010,
			TaskID:      1,
			ModuleNo:    8,
			Repeat:      2,
			ModuleTotal: 30,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130011,
			TaskID:      3,
			ModuleNo:    10,
			Repeat:      3,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, payload := range payload {
		filter := bson.M{"smartbtw_id": payload.SmartBtwID}
		update := bson.M{"$set": payload}
		stdCol.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-module-progress/detail/%s", "abcdef123"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	assert.Equalf(t, fiber.StatusBadRequest, response.StatusCode, "Get student module progress data with non integer parameter value must return 400 status code")
}

// success
func Test_GetStudentModuleProgressByTaskId_Success(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := []request.StudentModuleProgress{
		{
			SmartBtwID:  130009,
			TaskID:      1,
			ModuleNo:    1,
			Repeat:      2,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130010,
			TaskID:      2,
			ModuleNo:    8,
			Repeat:      2,
			ModuleTotal: 30,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130011,
			TaskID:      10,
			ModuleNo:    10,
			Repeat:      3,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, payload := range payload {
		filter := bson.M{"smartbtw_id": payload.SmartBtwID}
		update := bson.M{"$set": payload}
		stdCol.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-module-progress/task/%s", "10"),
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

// Error converted task_id to int
func Test_GetStudentModuleProgressByTaskIdError(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := []request.StudentModuleProgress{
		{
			SmartBtwID:  130009,
			TaskID:      1,
			ModuleNo:    1,
			Repeat:      2,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130010,
			TaskID:      1,
			ModuleNo:    8,
			Repeat:      2,
			ModuleTotal: 30,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130011,
			TaskID:      3,
			ModuleNo:    10,
			Repeat:      3,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, payload := range payload {
		filter := bson.M{"smartbtw_id": payload.SmartBtwID}
		update := bson.M{"$set": payload}
		stdCol.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-module-progress/task/%s", "abc"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	assert.Equalf(t, fiber.StatusBadRequest, response.StatusCode, "Get student module progress data with non integer parameter value must return 400 status code")
}

// Error Not found
func Test_GetStudentModuleProgressByTaskIdErrorNotFound(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	stdCol := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := []request.StudentModuleProgress{
		{
			SmartBtwID:  130009,
			TaskID:      1,
			ModuleNo:    1,
			Repeat:      2,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130010,
			TaskID:      1,
			ModuleNo:    8,
			Repeat:      2,
			ModuleTotal: 30,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			SmartBtwID:  130011,
			TaskID:      3,
			ModuleNo:    10,
			Repeat:      3,
			ModuleTotal: 48,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, payload := range payload {
		filter := bson.M{"smartbtw_id": payload.SmartBtwID}
		update := bson.M{"$set": payload}
		stdCol.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-module-progress/task/%s", "1000"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}
