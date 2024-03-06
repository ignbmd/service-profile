package lib_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func Test_CreateBKNScore(t *testing.T) {
	Init()

	payload := request.UpsertBKNScore{
		SmartBtwID: 500,
		Twk:        45,
		Tiu:        100,
		Tkp:        180,
		Total:      666,
		Year:       2023,
	}

	err := lib.UpsertBKNScore(&payload)
	assert.Nil(t, err)
}

func Test_GetBKNScore(t *testing.T) {
	Init()
	payload := request.UpsertBKNScore{
		SmartBtwID: 501,
		Twk:        45,
		Tiu:        100,
		Tkp:        180,
		Total:      666,
		Year:       2023,
	}

	err := lib.UpsertBKNScore(&payload)
	assert.Nil(t, err)

	res, err := lib.GetBKNScoreByArrayOfIDStudent([]int{500, 501, 502}, 2023)
	assert.Nil(t, err)
	fmt.Println(res)
}
