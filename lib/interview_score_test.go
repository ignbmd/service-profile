package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func Test_CreateInterviewScore_Success(t *testing.T) {
	Init()

	// Insert Interview Session
	interviewSessionPayload := request.InterviewSessionRequest{
		Name:        "Testing Sesi Interview",
		Description: "Testing Sesi Interview",
		Number:      1,
	}
	interviewSession, err := lib.CreateInterviewSession(&interviewSessionPayload)
	assert.Nil(t, err)

	insertedSessionID := interviewSession.InsertedID
	insertedSessionObjectID := insertedSessionID.(primitive.ObjectID)

	// Insert Interview Score
	note := "Catatan Tambahan"
	payload := request.UpsertInterviewScore{
		SmartBtwID:                        500,
		Penampilan:                        4,
		CaraDudukDanBerjabat:              4,
		PraktekBarisBerbaris:              4,
		PenampilanSopanSantun:             4,
		KepercayaanDiriDanStabilitasEmosi: 4,
		Komunikasi:                        4,
		PengembanganDiri:                  4,
		Integritas:                        4,
		Kerjasama:                         4,
		MengelolaPerubahan:                4,
		PerekatBangsa:                     4,
		PelayananPublik:                   4,
		PengambilanKeputusan:              4,
		OrientasiHasil:                    4,
		PrestasiAkademik:                  4,
		PrestasiNonAkademik:               4,
		BahasaAsing:                       4,
		FinalScore:                        100,
		SessionID:                         insertedSessionObjectID,
		Year:                              2023,
		CreatedBy: request.InterviewScoreCreatedUpdatedBy{
			ID:   "aa03bb90-047a-11eb-9b21-53f214d6ba9c",
			Name: "Bina Taruna Wiratama",
		},
		Note: &note,
	}

	res, err := lib.CreateInterviewScore(&payload)
	assert.Nil(t, err)

	insertedScoreID := res.InsertedID
	insertedScoreObjectID := insertedScoreID.(primitive.ObjectID)

	// Hard Delete Interview Session & Score
	var expectedDeletedCountMoreThan int64 = 0

	deleteRes, err := lib.HardDeleteInterviewSession(insertedSessionObjectID)
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
	assert.Nil(t, err)

	deleteRes, err = lib.HardDeleteInterviewScore(insertedScoreObjectID.Hex())
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
	assert.Nil(t, err)
}

func Test_UpdateInterviewScore_Success(t *testing.T) {
	Init()

	// Insert Interview Session
	interviewSessionPayload := request.InterviewSessionRequest{
		Name:        "Testing Sesi Interview",
		Description: "Testing Sesi Interview",
		Number:      1,
	}
	interviewSession, err := lib.CreateInterviewSession(&interviewSessionPayload)
	assert.Nil(t, err)

	insertedSessionID := interviewSession.InsertedID
	insertedSessionObjectID := insertedSessionID.(primitive.ObjectID)

	// Insert Interview Score
	initialNote := "Catatan Tambahan"
	updatedNote := "Catatan Tambahan Updated"

	createPayload := request.UpsertInterviewScore{
		SmartBtwID:                        500,
		Penampilan:                        4,
		CaraDudukDanBerjabat:              4,
		PraktekBarisBerbaris:              4,
		PenampilanSopanSantun:             4,
		KepercayaanDiriDanStabilitasEmosi: 4,
		Komunikasi:                        4,
		PengembanganDiri:                  4,
		Integritas:                        4,
		Kerjasama:                         4,
		MengelolaPerubahan:                4,
		PerekatBangsa:                     4,
		PelayananPublik:                   4,
		PengambilanKeputusan:              4,
		OrientasiHasil:                    4,
		PrestasiAkademik:                  4,
		PrestasiNonAkademik:               4,
		BahasaAsing:                       4,
		FinalScore:                        100,
		SessionID:                         insertedSessionObjectID,
		Year:                              2023,
		CreatedBy: request.InterviewScoreCreatedUpdatedBy{
			ID:   "aa03bb90-047a-11eb-9b21-53f214d6ba9c",
			Name: "Bina Taruna Wiratama",
		},
		Note: &initialNote,
	}

	interviewScore, err := lib.CreateInterviewScore(&createPayload)
	assert.Nil(t, err)

	insertedScoreID := interviewScore.InsertedID
	insertedScoreObjectID := insertedScoreID.(primitive.ObjectID)

	createdInterviewScore, err := lib.GetSingleInterviewScoreByID(insertedScoreObjectID)
	assert.Nil(t, err)

	// Update Interview Score
	updatePayload := request.UpsertInterviewScore{
		SmartBtwID:                        500,
		Penampilan:                        4,
		CaraDudukDanBerjabat:              4,
		PraktekBarisBerbaris:              4,
		PenampilanSopanSantun:             4,
		KepercayaanDiriDanStabilitasEmosi: 4,
		Komunikasi:                        4,
		PengembanganDiri:                  4,
		Integritas:                        4,
		Kerjasama:                         4,
		MengelolaPerubahan:                4,
		PerekatBangsa:                     4,
		PelayananPublik:                   4,
		PengambilanKeputusan:              4,
		OrientasiHasil:                    4,
		PrestasiAkademik:                  4,
		PrestasiNonAkademik:               4,
		BahasaAsing:                       4,
		FinalScore:                        100,
		SessionID:                         insertedSessionObjectID,
		Year:                              2023,
		CreatedBy: request.InterviewScoreCreatedUpdatedBy{
			ID:   "aa03bb90-047a-11eb-9b21-53f214d6ba9c",
			Name: "Bina Taruna Wiratama",
		},
		Note: &updatedNote,
	}
	err = lib.UpdateInterviewScore(insertedScoreObjectID.Hex(), &updatePayload)
	assert.Nil(t, err)

	updatedInterviewScore, err := lib.GetSingleInterviewScoreByID(insertedScoreObjectID)
	assert.Nil(t, err)

	var updatePayloadCreatedByID *string = &updatePayload.CreatedBy.ID
	var updatedInterviewScoreCreatedByID *string = &updatedInterviewScore.CreatedBy.ID

	var updatePayloadCreatedByName *string = &updatePayload.CreatedBy.Name
	var updatedInterviewScoreCreatedByName *string = &updatedInterviewScore.CreatedBy.Name

	assert.Nil(t, createdInterviewScore.UpdatedBy)
	assert.Equal(t, updatedInterviewScoreCreatedByID, updatePayloadCreatedByID)
	assert.Equal(t, updatePayloadCreatedByName, updatedInterviewScoreCreatedByName)
	assert.NotNil(t, updatedInterviewScore.UpdatedBy)

	// Hard Delete Interview Session & Score
	var expectedDeletedCountMoreThan int64 = 0
	deleteRes, err := lib.HardDeleteInterviewSession(insertedSessionObjectID)
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
	assert.Nil(t, err)

	deleteRes, err = lib.HardDeleteInterviewScore(insertedScoreObjectID.Hex())
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
	assert.Nil(t, err)
}

func Test_GetSingleInterviewScoreByStudentIDAndYear_Success(t *testing.T) {
	Init()

	// Insert Interview Session
	interviewSessionPayload := request.InterviewSessionRequest{
		Name:        "Testing Sesi Interview",
		Description: "Testing Sesi Interview",
		Number:      1,
	}
	interviewSession, err := lib.CreateInterviewSession(&interviewSessionPayload)
	assert.Nil(t, err)

	insertedSessionID := interviewSession.InsertedID
	insertedSessionObjectID := insertedSessionID.(primitive.ObjectID)

	// Insert Interview Score
	initialNote := "Catatan Tambahan"
	createPayload := request.UpsertInterviewScore{
		SmartBtwID:                        500,
		Penampilan:                        4,
		CaraDudukDanBerjabat:              4,
		PraktekBarisBerbaris:              4,
		PenampilanSopanSantun:             4,
		KepercayaanDiriDanStabilitasEmosi: 4,
		Komunikasi:                        4,
		PengembanganDiri:                  4,
		Integritas:                        4,
		Kerjasama:                         4,
		MengelolaPerubahan:                4,
		PerekatBangsa:                     4,
		PelayananPublik:                   4,
		PengambilanKeputusan:              4,
		OrientasiHasil:                    4,
		PrestasiAkademik:                  4,
		PrestasiNonAkademik:               4,
		BahasaAsing:                       4,
		FinalScore:                        100,
		SessionID:                         insertedSessionObjectID,
		Year:                              2023,
		CreatedBy: request.InterviewScoreCreatedUpdatedBy{
			ID:   "aa03bb90-047a-11eb-9b21-53f214d6ba9c",
			Name: "Bina Taruna Wiratama",
		},
		Note: &initialNote,
	}

	interviewScore, err := lib.CreateInterviewScore(&createPayload)
	assert.Nil(t, err)

	insertedScoreID := interviewScore.InsertedID
	insertedScoreObjectID := insertedScoreID.(primitive.ObjectID)

	res, err := lib.GetSingleInterviewScoreByStudentIDAndYear(500, 2023)
	assert.Nil(t, err)
	assert.NotEmpty(t, res)

	// Hard Delete Interview Session & Score
	var expectedDeletedCountMoreThan int64 = 0
	deleteRes, err := lib.HardDeleteInterviewSession(insertedSessionObjectID)
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
	assert.Nil(t, err)

	deleteRes, err = lib.HardDeleteInterviewScore(insertedScoreObjectID.Hex())
	assert.Greater(t, deleteRes.DeletedCount, expectedDeletedCountMoreThan)
	assert.Nil(t, err)
}

func Test_GetSingleInterviewScoreByStudentIDAndYear_DataNotFound(t *testing.T) {
	Init()

	_, err := lib.GetSingleInterviewScoreByStudentIDAndYear(998877, 2023)
	assert.NotNil(t, err)
}
