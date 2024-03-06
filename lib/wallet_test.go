package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestCreateWalletSuccess(t *testing.T) {
	Init()

	payload := request.CreateWallet{
		SmartbtwID: 515,
		Point:      100,
		Type:       string(models.DEFAULT),
	}

	_, err := lib.CreateWallet(&payload)
	assert.Nil(t, err)
}

func TestGetStudentWalletBalanceSuccess(t *testing.T) {
	Init()

	payload := request.CreateWallet{
		SmartbtwID: 333666999,
		Point:      100,
		Type:       string(models.DEFAULT),
	}

	_, err := lib.CreateWallet(&payload)
	if err != nil {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}

	_, err = lib.GetStudentWalletBalance(payload.SmartbtwID)
	assert.Nil(t, err)

}
