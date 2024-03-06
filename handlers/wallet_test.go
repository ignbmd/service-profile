package handlers_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func TestGetStudentWalletBalance(t *testing.T) {
	Init()

	payload := request.CreateWallet{
		SmartbtwID: 133182,
		Point:      0,
		Type:       string(models.BONUS),
	}

	_, err2 := lib.CreateWallet(&payload)
	if err2 != nil {
		assert.NotNil(t, err2)
	} else {
		assert.Nil(t, err2)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/wallet/%d/balance", payload.SmartbtwID),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")
	assert.Equal(t, nil, e)
	response, err := app.Test(request, -1)
	assert.Equal(t, nil, err)
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetStudentWalletBalanceInvalidDataType(t *testing.T) {
	Init()

	payload := request.CreateWallet{
		SmartbtwID: 133182,
		Point:      0,
		Type:       string(models.BONUS),
	}

	_, err2 := lib.CreateWallet(&payload)
	if err2 != nil {
		assert.NotNil(t, err2)
	} else {
		assert.Nil(t, err2)
	}

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/wallet/%s/balance", "133sss"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")
	assert.Equal(t, nil, e)
	response, err := app.Test(request, -1)
	assert.Equal(t, nil, err)
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}
