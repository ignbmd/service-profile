package lib_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func Test_CreateAvatar(t *testing.T) {
	Init()

	payload := request.CreateAvatar{
		SmartbtwID: 115,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	res, err := lib.CreateAvatar(&payload)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_UpdateAvatarBySmartbtwID(t *testing.T) {
	Init()
	payload1 := request.CreateAvatar{
		SmartbtwID: 333,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err := lib.CreateAvatar(&payload1)
	assert.Nil(t, err)

	payload := request.UpdateAvatar{
		SmartbtwID: 333,
		AvaType:    "PTK",
		Style:      2,
	}
	err1 := lib.UpdateAvatarSmartbtwID(&payload)
	assert.Nil(t, err1)
}

func Test_GetAvatarBySmartbtwIDAndType(t *testing.T) {
	Init()
	payload := request.CreateAvatar{
		SmartbtwID: 118,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err := lib.CreateAvatar(&payload)
	assert.Nil(t, err)

	req := request.BodyRequestAvatar{
		SmartbtwID: payload.SmartbtwID,
		AvaType:    payload.AvaType,
	}
	res, err := lib.GetAvatarBySmartBtwIDAndType(&req)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_DeleteAvatarBySmartbtwID(t *testing.T) {
	Init()
	payload := request.CreateAvatar{
		SmartbtwID: 117,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err1 := lib.CreateAvatar(&payload)
	assert.Nil(t, err1)

	req := request.BodyRequestAvatar{
		SmartbtwID: payload.SmartbtwID,
		AvaType:    payload.AvaType,
	}

	err := lib.DeleteAvatarBySmartbtwID(&req)
	assert.Nil(t, err)
}

func Test_GetAvatar(t *testing.T) {
	Init()
	payload := request.CreateAvatar{
		SmartbtwID: 888,
		AvaType:    string(models.PTK),
		Style:      1,
	}
	_, err1 := lib.CreateAvatar(&payload)
	assert.Nil(t, err1)

	res, err := lib.GetAvatarBySmartbtwID(888)
	assert.Nil(t, err)
	fmt.Println(res)
}
