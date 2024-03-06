package request

type UpsertSamaptaScore struct {
	SmartBtwID   int     `json:"smartbtw_id" bson:"smartbtw_id"`
	Gender       bool    `json:"gender" bson:"gender"`
	RunScore     float32 `json:"run_score" bson:"run_score"`
	PullUpScore  float32 `json:"pull_up_score" bson:"pull_up_score"`
	PushUpScore  float32 `json:"push_up_score" bson:"push_up_score"`
	SitUpScore   float32 `json:"sit_up_score" bson:"sit_up_score"`
	ShuttleScore float32 `json:"shuttle_score" bson:"shuttle_score"`
	Total        float32 `json:"total" bson:"total"`
	Year         uint16  `json:"year" bson:"year"`
}

type GetSamaptaScoreByArrEmail struct {
	Email []string `json:"email" bson:"email"`
	Year  uint16   `json:"year" bson:"year"`
}
