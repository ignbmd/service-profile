package lib_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func TestCreateStudentAccessSuccess(t *testing.T) {
	Init()

	payload := request.CreateStudentAccess{
		SmartBtwID:       515,
		DisallowedAccess: "ACCESS_BOARD_OF_SUCCESS",
		AppType:          "btwedutech",
	}

	_, err := lib.CreateStudentAccess(&payload)
	assert.Nil(t, err)
}

func TestGetStudentAccessSuccess(t *testing.T) {
	Init()

	tmUnx := time.Now().Unix()

	payload := request.CreateStudentAccess{
		SmartBtwID:       int(tmUnx),
		DisallowedAccess: "ACCESS_PROFILE",
		AppType:          "btwedutech",
	}

	_, err := lib.CreateStudentAccess(&payload)
	if err != nil {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}

	res, err := lib.GetStudentAllowedAccess(int(tmUnx), "btwedutech")
	assert.Nil(t, err)
	fmt.Println(res)

}
