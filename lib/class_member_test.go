package lib_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func TestCreateNewStudentClassMemberSuccess(t *testing.T) {
	Init()
	pyl1 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res, err := lib.CreateClassroom(&pyl1)
	assert.Nil(t, err)

	pyl := request.CreateClassMember{
		SmartbtwID:  1,
		ClassroomID: res.InsertedID.(primitive.ObjectID),
	}
	err = lib.CreateClassMember(&pyl)
	assert.Nil(t, err)
}

func TestGetAllClassMember(t *testing.T) {
	Init()
	res, err := lib.GetAllClassMember()
	assert.Nil(t, err)
	fmt.Println(res)
}

func TestUpdateClassMemberSuccess(t *testing.T) {
	Init()

	pyl1 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res, err := lib.CreateClassroom(&pyl1)
	assert.Nil(t, err)

	pyl2 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res2, err := lib.CreateClassroom(&pyl2)
	assert.Nil(t, err)

	pyl := request.UpdateClassMember{
		SmartbtwID:        1,
		ClassroomIDBefore: res.InsertedID.(primitive.ObjectID),
		ClassroomIDAfter:  res2.InsertedID.(primitive.ObjectID),
	}

	err = lib.UpdateClassMember(&pyl)
	assert.Nil(t, err)
}

func TestGetSingleClassMemberFromElasticSuccess(t *testing.T) {
	Init()

	pyl1 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
		IsOnline:    true,
	}
	res, err := lib.CreateClassroom(&pyl1)
	assert.Nil(t, err)

	pyl := request.CreateClassMember{
		SmartbtwID:  1,
		ClassroomID: res.InsertedID.(primitive.ObjectID),
	}
	err = lib.CreateClassMember(&pyl)
	assert.Nil(t, err)

	res1, err := lib.GetSingleClassMemberFromElastic(1, true)
	assert.Nil(t, err)
	fmt.Println(res1)
}

func TestSwitchClassMemberSuccess(t *testing.T) {
	Init()

	pyl1 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res, err := lib.CreateClassroom(&pyl1)
	assert.Nil(t, err)

	pyl := request.CreateClassMember{
		SmartbtwID:  87654,
		ClassroomID: res.InsertedID.(primitive.ObjectID),
	}
	err = lib.CreateClassMember(&pyl)
	assert.Nil(t, err)

	rPyl := request.SwitchClassMember{
		ClassroomID: res.InsertedID.(primitive.ObjectID),
		ClassMembers: []request.ArrayStudentSwitchClass{
			{
				SmartbtwID:   87654,
				BtwedutechID: 888888,
			},
		}}
	err = lib.SwitchClassMember(&rPyl)
	assert.Nil(t, err)
}
