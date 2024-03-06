package request

type CreateStudentTargetCpns struct {
	SmartbtwID        int     `json:"smartbtw_id"`
	InstanceID        int     `json:"instance_id"`
	InstanceName      string  `json:"instance_name"`
	PositionID        int     `json:"position_id"`
	PositionName      string  `json:"position_name"`
	FormationType     string  `json:"formation_type"`
	FormationLocation string  `json:"formation_location" bson:"formation_location"`
	CompetitionID     int     `json:"competition_id" bson:"competition_id"`
	FormationCode     string  `json:"formation_code" bson:"formation_code"`
	TargetScore       float64 `json:"target_score"`
	TargetType        string  `json:"target_type"`
	Gender            int     `json:"gender"`
}

type StudentTargetCpnsElastic struct {
	SmartbtwID              int     `json:"smartbtw_id"`
	Name                    string  `json:"name"`
	Photo                   string  `json:"photo"`
	InstanceID              int     `json:"instance_id"`
	InstanceName            string  `json:"instance_name"`
	PositionID              int     `json:"position_id"`
	PositionName            string  `json:"position_name"`
	FormationType           string  `json:"formation_type"`
	FormationLocation       string  `json:"formation_location" bson:"formation_location"`
	FormationYear           int     `json:"formation_year"`
	CompetitionID           int     `json:"competition_id" bson:"competition_id"`
	FormationCode           string  `json:"formation_code" bson:"formation_code"`
	TargetType              string  `json:"target_type"`
	TargetScore             float64 `json:"target_score"`
	ModuleDone              int     `json:"module_done"`
	TwkAvgScore             float64 `json:"twk_avg_score"`
	TiuAvgScore             float64 `json:"tiu_avg_score"`
	TkpAvgScore             float64 `json:"tkp_avg_score"`
	TotalAvgScore           float64 `json:"total_avg_score"`
	TwkAvgPercentScore      float64 `json:"twk_avg_percent_score"`
	TiuAvgPercentScore      float64 `json:"tiu_avg_percent_score"`
	TkpAvgPercentScore      float64 `json:"tkp_avg_percent_score"`
	TotalAvgPercentScore    float64 `json:"total_avg_percent_score"`
	LatestTotalScore        float64 `json:"latest_total_score"`
	LatestTotalPercentScore float64 `json:"latest_total_percent_score"`
	Proficiency             string  `json:"proficiency"`
}

type CompMapCompetitionCPNSChances struct {
	PositionID        uint   `json:"position_id"`
	FormationType     string `json:"formation_type"`
	FormationLocation string `json:"formation_location"`
	FormationCode     string `json:"formation_code"`
	Year              uint   `json:"year"`
	Quota             uint   `json:"quota"`
	Registered        uint   `json:"registered"`
	LowestScore       uint   `json:"lowest_score"`
	LowestPosition    uint   `json:"lowest_position"`
	LowestStatus      string `json:"lowest_status"`
	LastEdType        string `json:"last_ed_type"`
}

type ResponseContentCPNS struct {
	Data []CompMapCompetitionCPNSChances `json:"data"`
}
