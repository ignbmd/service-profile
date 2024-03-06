package lib_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestInsertHistoryAssessments(t *testing.T) {
	Init()

	py := request.CreateHistoryAssessment{
		SmartBtwID:     123,
		PackageID:      1,
		ModuleCode:     "MD-test",
		ModuleType:     "SKD",
		PackageType:    "tes",
		ScoreType:      "CLASSICAL",
		Total:          300,
		ExamName:       "Assessment Test",
		Start:          time.Now(),
		End:            time.Now().Add(1 * time.Hour),
		IsLive:         true,
		StudentName:    "Janu",
		StudentEmail:   "janu@btwe.com",
		Program:        "skd",
		ProgramType:    "PTK",
		ProgramVersion: 1,
		Scores: []models.HistoryAssessmentsScores{
			{
				Position:      1,
				CategoryName:  "Tes Wawasan",
				CategoryAlias: "TWK",
				Subtests: []models.HistoryAssessmentsSubScores{
					{
						Position:      1,
						SubName:       "Nasionalisme",
						SubAlias:      "nasionalisme",
						WrongAnswer:   0,
						CorrectAnswer: 1,
						EmptyAnswer:   1,
						Scores: []models.HistoryScores{
							{
								ScoreType: "CLASSICAL",
								Value:     20,
							},
						},
					},
				},
				Scores: []models.HistoryScores{
					{
						ScoreType: "CLASSICAL",
						Value:     20,
					},
				},
				WrongAnswer:   0,
				CorrectAnswer: 2,
				EmptyAnswer:   2,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := lib.UpsertHistoryAssessments(&py)
	assert.Nil(t, err)
}
