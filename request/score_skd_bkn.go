package request

type ScoreSkdBkn struct {
	SmartBtwID int     `json:"smartbtw_id" bson:"smartbtw_id" valid:"type(int),required"`
	Year       int     `json:"year" bson:"year" valid:"type(int),required"`
	ScoreTWK   float32 `json:"score_twk" bson:"score_twk" valid:"type(float32),required"`
	ScoreTIU   float32 `json:"score_tiu" bson:"score_tiu" valid:"type(float32),required"`
	ScoreTKP   float32 `json:"score_tkp" bson:"score_tkp" valid:"type(float32),required"`
	ScoreSKD   float32 `json:"score_skd" bson:"score_skd" valid:"type(float32),required"`
}

type UpdateScoreSKDBKN struct {
	SmartBtwID int     `json:"smartbtw_id" bson:"smartbtw_id" valid:"type(int),required"`
	ScoreTWK   float32 `json:"score_twk" bson:"score_twk" valid:"type(float32),required"`
	ScoreTIU   float32 `json:"score_tiu" bson:"score_tiu" valid:"type(float32),required"`
	ScoreTKP   float32 `json:"score_tkp" bson:"score_tkp" valid:"type(float32),required"`
	ScoreSKD   float32 `json:"score_skd" bson:"score_skd" valid:"type(float32),required"`
}

type GetManyRecordScore struct {
	SmartBtwID []int `json:"smartbtw_id" bson:"smartbtw_id" valid:"type([]int),required"`
	Year       int   `json:"year" bson:"year" valid:"type(int),required"`
}
