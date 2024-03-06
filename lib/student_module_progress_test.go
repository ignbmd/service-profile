package lib_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func Test_UpsertStudentModuleProgress(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  130009,
		TaskID:      1,
		ModuleNo:    30,
		Repeat:      2,
		ModuleTotal: 30,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.UpsertStudentModuleProgress(&payload)
	assert.Nil(t, err)
	log.Println(err)
}

func Test_UpsertStudentModuleProgress_ModuleNoGreaterThanModuleTotal(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  130009,
		TaskID:      1,
		ModuleNo:    30,
		Repeat:      1,
		ModuleTotal: 28,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.UpsertStudentModuleProgress(&payload)
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_CreateStudentModuleProgress_Success(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  1300066,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 20,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)
	log.Println(err)
}

func Test_CreateStudentModuleProgress_AlreadyExist(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  130009,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 20,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_CreateStudentModuleProgress_ErrorModuleNoGreaterThanModuleTotal(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  1300109,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 19,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_CreateStudentModuleProgress_ErrorNegativeValue(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  -1300109,
		TaskID:      -2,
		ModuleNo:    -20,
		Repeat:      -1,
		ModuleTotal: -19,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_CreateStudentModuleProgress_ErrorZeroValue(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  0,
		TaskID:      0,
		ModuleNo:    0,
		Repeat:      0,
		ModuleTotal: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_UpdateStudentModuleProgress_Success(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  700,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 20,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	payloadUpdate := request.UpdateStudentModuleProgress{
		TaskID:      2,
		ModuleNo:    30,
		Repeat:      1,
		ModuleTotal: 50,
	}

	err = lib.UpdateStudentModuleProgress(&payloadUpdate, res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
	log.Println(err)
}

func Test_UpdateStudentModuleProgress_ErrorZeroValue(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  8001,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 20,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	payloadUpdate := request.UpdateStudentModuleProgress{
		TaskID:      0,
		ModuleNo:    0,
		Repeat:      0,
		ModuleTotal: 0,
	}

	err = lib.UpdateStudentModuleProgress(&payloadUpdate, res.InsertedID.(primitive.ObjectID))
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_UpdateStudentModuleProgress_ErrorNotFound(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  130091,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 20,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	payloadUpdate := request.UpdateStudentModuleProgress{
		TaskID:      2,
		ModuleNo:    25,
		Repeat:      1,
		ModuleTotal: 30,
	}

	err = lib.UpdateStudentModuleProgress(&payloadUpdate, primitive.NewObjectID())
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_UpdateStudentModuleProgress_ErrorModuleNoGreaterThanModuleTotal(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  130098,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 21,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	payloadUpdate := request.UpdateStudentModuleProgress{
		TaskID:      2,
		ModuleNo:    25,
		Repeat:      1,
		ModuleTotal: 3,
	}

	err = lib.UpdateStudentModuleProgress(&payloadUpdate, res.InsertedID.(primitive.ObjectID))
	assert.NotNil(t, err)
	log.Println(err)
}

func Test_UpdateStudentModuleProgress_ErrorNegativeValue(t *testing.T) {
	Init()

	payload := request.StudentModuleProgress{
		SmartBtwID:  1300988,
		TaskID:      2,
		ModuleNo:    20,
		Repeat:      1,
		ModuleTotal: 21,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	res, err := lib.CreateStudentModuleProgress(&payload)
	assert.Nil(t, err)

	payloadUpdate := request.UpdateStudentModuleProgress{
		TaskID:      -2,
		ModuleNo:    -5,
		Repeat:      -1,
		ModuleTotal: -3,
	}

	err = lib.UpdateStudentModuleProgress(&payloadUpdate, res.InsertedID.(primitive.ObjectID))
	assert.NotNil(t, err)
	log.Println(err)
}