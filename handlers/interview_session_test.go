package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func Test_CreateInterviewSession_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi 1 Testing",
		Description: "Sesi 1 Testing",
		Number:      1,
	}

	marshalPayload, err := json.Marshal(payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"POST",
		"/interview-session",
		bytes.NewBuffer(marshalPayload),
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

	insertedID := jBody.(map[string]interface{})["data"].(map[string]interface{})["id"].(string)
	insertedObjectID, err := primitive.ObjectIDFromHex(insertedID)
	assert.Nil(t, err)

	_, err = lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func Test_CreateInterviewSession_ValidationError(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{}

	marshalPayload, err := json.Marshal(payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"POST",
		"/interview-session",
		bytes.NewBuffer(marshalPayload),
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
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_CreateInterviewSession_StructDataTypeError(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"name": "Sesi 1",
		"description": "Sesi 1 Testing",
		"number":"1"
	}`)

	request, e := http.NewRequest(
		"POST",
		"/interview-session",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	assert.Equal(t, nil, err)

	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func Test_GetAllInterviewSessions_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi 1 Testing",
		Description: "Sesi 1 Testing",
		Number:      1,
	}

	interviewSession, err := lib.CreateInterviewSession(&payload)
	assert.Nil(t, err)
	assert.NotNil(t, interviewSession.InsertedID)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		"/interview-session",
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

	data := jBody.(map[string]interface{})["data"].([]interface{})
	assert.NotNil(t, data)
	assert.Greater(t, len(data), 0)

	insertedID := interviewSession.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)

	_, err = lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetSingleInterviewSessionByID_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi 1 Testing",
		Description: "Sesi 1 Testing",
		Number:      1,
	}

	interviewSession, err := lib.CreateInterviewSession(&payload)
	insertedID := interviewSession.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)
	insertedObjectIDHexString := insertedObjectID.Hex()

	assert.Nil(t, err)
	assert.NotNil(t, interviewSession.InsertedID)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/interview-session/%s", insertedObjectIDHexString),
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

	data := jBody.(map[string]interface{})["data"].(map[string]interface{})
	assert.NotNil(t, data)

	_, err = lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetSingleInterviewSessionByID_InvalidObjectID(t *testing.T) {
	Init()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/interview-session/%s", "InvalidObjectID"),
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

	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func Test_GetSingleInterviewSessionByID_DataNotFound(t *testing.T) {
	Init()

	objectID := primitive.NewObjectID()
	objectIDHexString := objectID.Hex()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/interview-session/%s", objectIDHexString),
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

	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func Test_UpdateInterviewSession_Success(t *testing.T) {
	Init()

	createPayload := request.InterviewSessionRequest{
		Name:        "Sesi 1 Testing",
		Description: "Sesi 1 Testing",
		Number:      1,
	}
	interviewSession, err := lib.CreateInterviewSession(&createPayload)
	assert.Nil(t, err)

	insertedID := interviewSession.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)
	insertedHexString := insertedID.(primitive.ObjectID).Hex()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi 2 Testing Updated",
		Description: "Sesi 2 Testing Updated",
		Number:      2,
	}

	marshalPayload, err := json.Marshal(payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/interview-session/%s", insertedHexString),
		bytes.NewBuffer(marshalPayload),
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

	_, err = lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, err)
}

func Test_UpdateInterviewSession_StructDataTypeError(t *testing.T) {
	Init()

	createPayload := request.InterviewSessionRequest{
		Name:        "Sesi 1 Testing",
		Description: "Sesi 1 Testing",
		Number:      1,
	}
	interviewSession, err := lib.CreateInterviewSession(&createPayload)
	assert.Nil(t, err)

	insertedID := interviewSession.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)
	insertedHexString := insertedID.(primitive.ObjectID).Hex()

	pyl := []byte(`{
		"name": "Sesi 1",
		"description": "Sesi 1 Testing",
		"number":"1"
	}`)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/interview-session/%s", insertedHexString),
		bytes.NewBuffer(pyl),
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

	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	_, err = lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, err)
}

func Test_UpdateInterviewSession_ValidationError(t *testing.T) {
	Init()

	createPayload := request.InterviewSessionRequest{
		Name:        "Sesi 1 Testing",
		Description: "Sesi 1 Testing",
		Number:      1,
	}
	interviewSession, err := lib.CreateInterviewSession(&createPayload)
	assert.Nil(t, err)

	insertedID := interviewSession.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)
	insertedHexString := insertedID.(primitive.ObjectID).Hex()

	payload := request.InterviewSessionRequest{}

	marshalPayload, err := json.Marshal(payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/interview-session/%s", insertedHexString),
		bytes.NewBuffer(marshalPayload),
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

	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)

	_, err = lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, err)
}

func Test_UpdateInterviewSession_InvalidObjectID(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi 2 Testing Updated",
		Description: "Sesi 2 Testing Updated",
		Number:      2,
	}

	marshalPayload, err := json.Marshal(payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/interview-session/%s", "InvalidObjectID"),
		bytes.NewBuffer(marshalPayload),
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

	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func Test_UpdateInterviewSession_DataNotFound(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi 2 Testing Updated",
		Description: "Sesi 2 Testing Updated",
		Number:      2,
	}

	randomObjectID := primitive.NewObjectID()
	randomObjectIDHexString := randomObjectID.Hex()

	marshalPayload, err := json.Marshal(payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/interview-session/%s", randomObjectIDHexString),
		bytes.NewBuffer(marshalPayload),
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

	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func Test_SoftDeleteInterviewSession_Success(t *testing.T) {
	Init()

	createPayload := request.InterviewSessionRequest{
		Name:        "Sesi 1 Testing",
		Description: "Sesi 1 Testing",
		Number:      1,
	}
	interviewSession, err := lib.CreateInterviewSession(&createPayload)
	assert.Nil(t, err)

	insertedID := interviewSession.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)
	insertedHexString := insertedID.(primitive.ObjectID).Hex()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/interview-session/%s", insertedHexString),
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

	_, err = lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, err)
}

func Test_SoftDeleteInterviewSession_DataNotFound(t *testing.T) {
	Init()

	objectID := primitive.NewObjectID()
	objectIDHexString := objectID.Hex()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/interview-session/%s", objectIDHexString),
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

	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}

func Test_SoftDeleteInterviewSession_InvalidObjectID(t *testing.T) {
	Init()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/interview-session/%s", "InvalidObjectID"),
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

	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}
