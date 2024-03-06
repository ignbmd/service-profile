package lib_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func TestCreateStudentTargetCpns(t *testing.T) {
	Init()

	py := request.CreateStudentTargetCpns{
		SmartbtwID:        140015,
		InstanceID:        1,
		PositionID:        3,
		InstanceName:      "Kementrian Agama",
		PositionName:      "testing",
		FormationLocation: "PROVINCE",
		FormationType:     "GENERAL",
		FormationCode:     "BALADU",
		CompetitionID:     1,
		TargetScore:       430,
	}
	_, err := lib.CreateStudentTargetCpns(&py)
	assert.Nil(t, err)
}

func TestGetStudentTargetCPNS(t *testing.T) {
	Init()
	res, err := lib.GetStudentTargetCPNS(140015)
	assert.Nil(t, err)
	fmt.Println(res)
}
