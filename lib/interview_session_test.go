package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func Test_GetAllInterviewSessions_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1",
		Description: "Sesi pertama dari jam 9:30 - 12:30 WIB",
		Number:      1,
	}

	createRes, createErr := lib.CreateInterviewSession(&payload)
	assert.Nil(t, createErr)
	assert.NotNil(t, createRes.InsertedID)

	insertedID := createRes.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)

	interviewSessions, interviewSessionsErr := lib.GetAllInterviewSessions()
	assert.Nil(t, interviewSessionsErr)
	assert.NotNil(t, interviewSessions)
	assert.Greater(t, len(interviewSessions), 0)
	assert.Len(t, interviewSessions, 1)

	var expectedDeletedCountMoreThan int64 = 0
	deleteRes, deleteErr := lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, deleteErr)
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
}

func Test_GetSingleInterviewSessionByID_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1",
		Description: "Sesi pertama dari jam 9:30 - 12:30 WIB",
		Number:      1,
	}

	createRes, createErr := lib.CreateInterviewSession(&payload)
	assert.Nil(t, createErr)
	assert.NotNil(t, createRes.InsertedID)

	insertedID := createRes.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)

	interviewSession, interviewSessionErr := lib.GetSingleInterviewSessionByID(insertedObjectID)
	assert.Nil(t, interviewSessionErr)
	assert.NotNil(t, interviewSession)

	var expectedDeletedCountMoreThan int64 = 0
	deleteRes, deleteErr := lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, deleteErr)
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
}

func Test_GetSingleInterviewSessionByID_DataNotFound(t *testing.T) {
	Init()
	objectID := primitive.NewObjectID()
	_, interviewSessionErr := lib.GetSingleInterviewSessionByID(objectID)
	assert.NotNil(t, interviewSessionErr)
}

func Test_CreateInterviewSession_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1",
		Description: "Sesi pertama dari jam 9:30 - 12:30 WIB",
		Number:      1,
	}

	createRes, createErr := lib.CreateInterviewSession(&payload)
	assert.Nil(t, createErr)
	assert.NotNil(t, createRes.InsertedID)

	insertedID := createRes.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)

	var expectedDeletedCountMoreThan int64 = 0
	deleteRes, deleteErr := lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, deleteErr)
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
}

func Test_UpdateInterviewSession_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1",
		Description: "Sesi pertama dari jam 9:30 - 12:30 WIB",
		Number:      1,
	}

	createRes, createErr := lib.CreateInterviewSession(&payload)
	assert.Nil(t, createErr)
	assert.NotNil(t, createRes.InsertedID)

	insertedID := createRes.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)

	oldInterviewSession, oldInterviewSessionErr := lib.GetSingleInterviewSessionByID(insertedObjectID)
	assert.Nil(t, oldInterviewSessionErr)
	assert.NotNil(t, oldInterviewSession)

	var expectedModifiedCountMoreThan int64 = 0
	updatePayload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1 Updated",
		Description: "Sesi kedua dari jam 13:00 - 15:00 WIB",
		Number:      2,
	}

	updateRes, updateErr := lib.UpdateInterviewSession(insertedObjectID, &updatePayload)
	assert.Nil(t, updateErr)
	assert.Greater(t, updateRes.ModifiedCount, expectedModifiedCountMoreThan)

	newInterviewSession, newInterviewSessionErr := lib.GetSingleInterviewSessionByID(insertedObjectID)
	assert.Nil(t, newInterviewSessionErr)
	assert.NotNil(t, newInterviewSession)

	assert.NotEqualValues(t, oldInterviewSession.Name, newInterviewSession.Name)
	assert.NotEqualValues(t, oldInterviewSession.Description, newInterviewSession.Description)
	assert.NotEqualValues(t, oldInterviewSession.Number, newInterviewSession.Number)

	var expectedDeletedCountMoreThan int64 = 0
	hardDeleteRes, hardDeleteErr := lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, hardDeleteErr)
	assert.Greater(t, hardDeleteRes.DeletedCount, expectedDeletedCountMoreThan)
}

func Test_UpdateInterviewSession_DataNotFound(t *testing.T) {
	Init()
	objectID := primitive.NewObjectID()

	var expectedModifiedCountEquals int64 = 0
	updatePayload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1 Updated",
		Description: "Sesi kedua dari jam 13:00 - 15:00 WIB",
		Number:      2,
	}

	updateRes, _ := lib.UpdateInterviewSession(objectID, &updatePayload)
	assert.Equal(t, updateRes.ModifiedCount, expectedModifiedCountEquals)
}

func Test_SoftDeleteInterviewSession_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1",
		Description: "Sesi pertama dari jam 9:30 - 12:30 WIB",
		Number:      1,
	}

	createRes, createErr := lib.CreateInterviewSession(&payload)
	assert.Nil(t, createErr)
	assert.NotNil(t, createRes.InsertedID)

	insertedID := createRes.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)

	var expectedModifiedCountMoreThan int64 = 0
	softDeleteRes, softDeleteErr := lib.SoftDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, softDeleteErr)
	assert.Greater(t, softDeleteRes.ModifiedCount, expectedModifiedCountMoreThan)

	var expectedDeletedCountMoreThan int64 = 0
	hardDeleteRes, hardDeleteErr := lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, hardDeleteErr)
	assert.Greater(t, hardDeleteRes.DeletedCount, expectedDeletedCountMoreThan)
}

func Test_SoftDeleteInterviewSession_DataNotFound(t *testing.T) {
	Init()
	objectID := primitive.NewObjectID()

	var expectedModifiedCountEquals int64 = 0

	updateRes, _ := lib.SoftDeleteInterviewSession(objectID)
	assert.Equal(t, updateRes.ModifiedCount, expectedModifiedCountEquals)
}

func Test_HardDeleteInterviewSession_Success(t *testing.T) {
	Init()

	payload := request.InterviewSessionRequest{
		Name:        "Sesi Wawancara Testing #1",
		Description: "Sesi pertama dari jam 9:30 - 12:30 WIB",
		Number:      1,
	}

	createRes, createErr := lib.CreateInterviewSession(&payload)
	assert.Nil(t, createErr)
	assert.NotNil(t, createRes.InsertedID)

	insertedID := createRes.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID)

	var expectedDeletedCountMoreThan int64 = 0
	hardDeleteRes, hardDeleteErr := lib.HardDeleteInterviewSession(insertedObjectID)
	assert.Nil(t, hardDeleteErr)
	assert.Greater(t, hardDeleteRes.DeletedCount, expectedDeletedCountMoreThan)
}
