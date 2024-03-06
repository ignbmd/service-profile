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
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func TestCreateStudentAccessSuccess(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 5999,
		"disallowed_access": "CODE_PROFILE",
		"app_type":"btwedutech"
	}`)

	request, e := http.NewRequest(
		"POST",
		"/student-access",
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

func TestDeleteStudentAccess(t *testing.T) {
	Init()

	payload := request.CreateStudentAccess{
		SmartBtwID:       912912,
		DisallowedAccess: "CODE_PROFILE",
		AppType:          "btwedutech",
	}

	res, err := lib.CreateStudentAccess(&payload)
	assert.Nil(t, err)
	fmt.Println(res.InsertedID.(primitive.ObjectID))
	body := []byte(`{
		"disallowed_access": "CODE_PROFILE",
		"app_type":"btwedutech"
	}`)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/student-access/student-id/%d", 912912),
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

func TestGetStudentAccessByID(t *testing.T) {
	Init()

	payload := request.CreateStudentAccess{
		SmartBtwID:       91291223,
		DisallowedAccess: "CODE_PROFILE",
		AppType:          "btwedutech",
	}

	_, err := lib.CreateStudentAccess(&payload)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-access/by-student-id/%d", 91291223),
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
