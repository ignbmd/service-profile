package lib_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func Test_UpsertBranchSuccess(t *testing.T) {
	Init()

	payload := request.UpsertBranchData{
		BranchCode: "PT0000",
		BranchName: "Cabang Pusat",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	fmt.Println(payload)
	err := lib.UpsertBranchData(&payload)
	assert.Equal(t, nil, err)
}

func Test_UpsertBranchNullDataError(t *testing.T) {
	Init()
	payload := request.UpsertBranchData{
		BranchCode: "",
		BranchName: "",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	err := lib.UpsertBranchData(&payload)
	assert.NotNil(t, err)
}

func Test_UpsertBranchOnlyNameNullDataError(t *testing.T) {
	Init()
	payload := request.UpsertBranchData{
		BranchCode: "PT0000",
		BranchName: "",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := lib.UpsertBranchData(&payload)
	assert.NotNil(t, err)
}
