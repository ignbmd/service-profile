package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func TestInsertHistoryCPNSElastic(t *testing.T) {
	Init()

	py := request.CreateHistoryCpns{
		SmartBtwID:  123,
		TaskID:      1,
		PackageID:   1,
		ModuleCode:  "MD-test",
		ModuleType:  "SKD",
		PackageType: "tes",
		Twk:         122,
		Tiu:         123,
		Tkp:         111,
		Total:       300,
	}
	err := lib.InsertStudentHistoryCPNSElastic(&py, "640eb64f5c36377c27cc47ed")
	assert.Nil(t, err)
}
