package request

import "smartbtw.com/services/profile/models"

type UpsertBKNScore struct {
	SmartBtwID    int     `json:"smartbtw_id" bson:"smartbtw_id"`
	Twk           float64 `json:"twk" bson:"twk"`
	Tiu           float64 `json:"tiu" bson:"tiu"`
	Tkp           float64 `json:"tkp" bson:"tkp"`
	Total         float64 `json:"total" bson:"total"`
	Year          uint16  `json:"year" bson:"year"`
	IsContinue    bool    `json:"is_continue" bson:"is_continue"`
	BKNRank       uint32  `json:"bkn_rank" bson:"bkn_rank"`
	PtkSchoolId   uint32  `json:"ptk_school_id" bson:"ptk_school_id"`
	PtkSchool     string  `json:"ptk_school" bson:"ptk_school"`
	PtkMajorId    uint32  `json:"ptk_major_id" bson:"ptk_major_id"`
	PtkMajor      string  `json:"ptk_major" bson:"ptk_major"`
	BknTestNumber string  `json:"bkn_test_number" bson:"bkn_test_number"`
}

type GetBKNScoreByArrSmID struct {
	SmartBtwID []int  `json:"smartbtw_id" bson:"smartbtw_id"`
	Year       uint16 `json:"year" bson:"year"`
}

type GetBKNScoreByArrEmail struct {
	Email []string `json:"email" bson:"email"`
	Year  uint16   `json:"year" bson:"year"`
}

type UpdateBKNScoreForSurvey struct {
	SmartBtwID     int                 `json:"smartbtw_id" bson:"smartbtw_id"`
	Year           uint16              `json:"year" bson:"year"`
	SurveyStatus   bool                `json:"survey_status" bson:"survey_status"`
	Suggestion     string              `json:"suggestion" bson:"suggestion"`
	Reason         models.ReasonFailed `json:"reason" bson:"reason"`
	ReturnedResult string              `json:"returned_result" bson:"returned_result"`
}
type UpdateBKNScoreForProdi struct {
	SmartBtwID       int    `json:"smartbtw_id" bson:"smartbtw_id"`
	Year             uint16 `json:"year" bson:"year"`
	PtkSchoolId      uint32 `json:"ptk_school_id" bson:"ptk_school_id"`
	PtkSchool        string `json:"ptk_school" bson:"ptk_school"`
	PtkMajorId       uint32 `json:"ptk_major_id" bson:"ptk_major_id"`
	PtkMajor         string `json:"ptk_major" bson:"ptk_major"`
	PtkCompetitionID uint   `json:"ptk_competition_id" bson:"ptk_competition_id"`
	BknTestNumber    string `json:"bkn_test_number" bson:"bkn_test_number"`
}
