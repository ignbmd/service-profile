package lib_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestCreateWalletHistoryUKA(t *testing.T) {
	Init()

	payload1 := request.CreateWallet{
		SmartbtwID: 515,
		Point:      0,
		Type:       string(models.BONUS),
	}

	_, err1 := lib.CreateWallet(&payload1)
	assert.Nil(t, err1)

	payload := request.CreateWalletHistory{
		SmartbtwID:  515,
		Description: "mengerjakan uka",
	}

	_, err := lib.CreateWalletHistoryUKA(&payload)
	assert.Nil(t, err)
}

func TestCreateWalletHistoryPremiumPackage(t *testing.T) {
	Init()

	payload1 := request.CreateWallet{
		SmartbtwID: 12666,
		Point:      0,
		Type:       string(models.BONUS),
	}

	_, err1 := lib.CreateWallet(&payload1)
	assert.Nil(t, err1)
	payload := request.CreateWalletHistoryPremium{
		SmartbtwID:  12666,
		Description: "mpembelian paket premium",
		Price:       5000000,
	}

	_, err := lib.CreateWalletHistoryPremiumPackage(&payload)
	assert.Nil(t, err)
}

func TestCreateWalletHistoryInvitePeople(t *testing.T) {
	Init()

	payload1 := request.CreateWallet{
		SmartbtwID: 12676,
		Point:      0,
		Type:       string(models.BONUS),
	}

	_, err1 := lib.CreateWallet(&payload1)
	assert.Nil(t, err1)
	payload := request.CreateWalletHistory{
		SmartbtwID:  12676,
		Description: "mengundang teman",
	}

	_, err := lib.CreateWalletHistoryInvitePeople(&payload)
	assert.Nil(t, err)
}

func TestGetStudentWalletHistorySuccess(t *testing.T) {
	Init()

	payload1 := request.CreateWallet{
		SmartbtwID: 1267436,
		Point:      0,
		Type:       string(models.BONUS),
	}

	_, err1 := lib.CreateWallet(&payload1)
	assert.Nil(t, err1)
	payload := request.CreateWalletHistory{
		SmartbtwID:  1267436,
		Description: "mengundang teman",
	}
	_, err2 := lib.CreateWalletHistoryInvitePeople(&payload)
	assert.Nil(t, err2)

	_, err := lib.GetStudentWalletHistory(payload.SmartbtwID, nil)
	assert.Nil(t, err)
}

func TestChargeWallet_Success(t *testing.T) {
	Init()

	smbtwId := int(time.Now().Unix() - 921)

	payloadBonus := request.CreateWallet{
		SmartbtwID: smbtwId,
		Point:      101,
		Type:       string(models.DEFAULT),
	}

	_, err1 := lib.CreateWallet(&payloadBonus)

	assert.Nil(t, err1)
	payloadDefault := request.CreateWallet{
		SmartbtwID: smbtwId,
		Point:      50,
		Type:       string(models.BONUS),
	}

	_, err2 := lib.CreateWallet(&payloadDefault)

	assert.Nil(t, err2)

	payload := request.ChargeWallet{
		SmartbtwID:  smbtwId,
		Description: "mpembelian paket wulu wulu",
		Point:       150,
	}

	err := lib.ChargeStudentWallet(&payload)
	assert.Nil(t, err)
}

func TestChargeWallet_FailedInsufficient(t *testing.T) {
	Init()

	smbtwId := int(time.Now().Unix() - 999)

	payloadBonus := request.CreateWallet{
		SmartbtwID: smbtwId,
		Point:      10,
		Type:       string(models.DEFAULT),
	}

	_, err1 := lib.CreateWallet(&payloadBonus)

	assert.Nil(t, err1)
	payloadDefault := request.CreateWallet{
		SmartbtwID: smbtwId,
		Point:      50,
		Type:       string(models.BONUS),
	}

	_, err2 := lib.CreateWallet(&payloadDefault)

	assert.Nil(t, err2)

	payload := request.ChargeWallet{
		SmartbtwID:  smbtwId,
		Description: "mpembelian paket wulu wulu",
		Point:       150,
	}

	err := lib.ChargeStudentWallet(&payload)
	assert.NotNil(t, err)
}

func TestChargeWallet_SuccessBonusOnly(t *testing.T) {
	Init()

	smbtwId := int(time.Now().Unix() - 91)

	payloadBonus := request.CreateWallet{
		SmartbtwID: smbtwId,
		Point:      10,
		Type:       string(models.DEFAULT),
	}

	_, err1 := lib.CreateWallet(&payloadBonus)

	assert.Nil(t, err1)
	payloadDefault := request.CreateWallet{
		SmartbtwID: smbtwId,
		Point:      161,
		Type:       string(models.BONUS),
	}

	_, err2 := lib.CreateWallet(&payloadDefault)

	assert.Nil(t, err2)

	payload := request.ChargeWallet{
		SmartbtwID:  smbtwId,
		Description: "mpembelian paket wulu wulu",
		Point:       150,
	}

	err := lib.ChargeStudentWallet(&payload)
	assert.Nil(t, err)
}

func TestChargeWallet_FailedWalletNotFound(t *testing.T) {
	Init()

	smbtwId := int(time.Now().Unix() - 88888)

	payload := request.ChargeWallet{
		SmartbtwID:  smbtwId,
		Description: "mpembelian paket wulu wulu",
		Point:       150,
	}

	err := lib.ChargeStudentWallet(&payload)
	assert.NotNil(t, err)
}
