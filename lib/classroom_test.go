package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func TestCreateClassroomSuccess(t *testing.T) {
	Init()
	pyl := request.CreateClassroom{
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
	_, err := lib.CreateClassroom(&pyl)
	assert.Nil(t, err)
}

func TestUpdateClassroomSuccess(t *testing.T) {
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

	pyl := request.UpdateClassroom{
		ID:          res.InsertedID.(primitive.ObjectID),
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

	err = lib.UpdateClassroom(&pyl)
	assert.Nil(t, err)

}
