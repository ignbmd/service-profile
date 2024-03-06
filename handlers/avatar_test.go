package handlers_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func Test_CreateAvatar(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
	"smartbtw_id": 133158,
	"ava_type": "PTK",
	"style": 1
	}`)
	request, e := http.NewRequest(
		"POST",
		"/avatar",
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

func Test_CreateAvatarFailedBody(t *testing.T) {
	Init()

	payload := request.CreateAvatar{
		SmartbtwID: 119,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err1 := lib.CreateAvatar(&payload)
	assert.Nil(t, err1)

	app := server.SetupFiber()
	body := []byte(`{
	"smartbtw_id": 119,
	"ava_type": "PTN",
	"style": "1"
	}`)
	request, e := http.NewRequest(
		"POST",
		"/avatar",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 422
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_UpdateAvatar(t *testing.T) {
	Init()

	payload := request.CreateAvatar{
		SmartbtwID: 120,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err1 := lib.CreateAvatar(&payload)
	assert.Nil(t, err1)

	app := server.SetupFiber()
	body := []byte(`{
	"smartbtw_id": 120,
	"ava_type": "PTK",
	"style": 2
	}`)
	request, e := http.NewRequest(
		"PUT",
		"/avatar",
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

func Test_GetAvatar(t *testing.T) {
	Init()

	payload := request.CreateAvatar{
		SmartbtwID: 121,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err1 := lib.CreateAvatar(&payload)
	assert.Nil(t, err1)

	app := server.SetupFiber()
	body := []byte(`{
	"smartbtw_id": 121,
	"ava_type": "PTK"
	}`)
	request, e := http.NewRequest(
		"POST",
		"/student-avatar",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 200
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_DeleteAvatar(t *testing.T) {
	Init()
	payload := request.CreateAvatar{
		SmartbtwID: 123,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err1 := lib.CreateAvatar(&payload)
	assert.Nil(t, err1)

	app := server.SetupFiber()
	body := []byte(`{
	"smartbtw_id": 123,
	"ava_type": "PTK"
	}`)
	request, e := http.NewRequest(
		"DELETE",
		"/avatar",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)

}

func Test_GetAvatarBySmartbtwID(t *testing.T) {
	Init()

	payload := request.CreateAvatar{
		SmartbtwID: 555,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err1 := lib.CreateAvatar(&payload)
	assert.Nil(t, err1)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		"/avatar/555",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 200
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}
