package request

type HistoryScoreQueryParams struct {
	SmartBTWID      *int  `query:"smartbtw_id"`
	TaskID          *int  `query:"task_id"`
	OnlyLatestScore *bool `query:"only_latest_score"`
}

type BodySendAssessmentCompleted struct {
	PackageID  uint    `json:"package_id"`
	SmartbtwID uint    `json:"smartbtw_id"`
	Grade      string  `json:"grade"`
	Score      float64 `json:"score"`
	TaskID     uint    `json:"task_id"`
	Program    string  `json:"program"`
}

type SendAssessmentCompleted struct {
	Version uint                        `json:"version"`
	Data    BodySendAssessmentCompleted `json:"data"`
}

type ProfileStudentResultHistory struct {
	Date  string `json:"date"`
	Title string `json:"title"`
	Total int    `json:"total"`
}
type ProfileStudentResultHistoryPTN struct {
	Date  string  `json:"date"`
	Title string  `json:"title"`
	Total float64 `json:"total"`
}

type FirebaseHistoryScores struct {
	TaskID          int                              `json:"task_id"`
	SmartbtwID      uint                             `json:"smartbtw_id"`
	SchoolID        int                              `json:"school_id"`
	MajorID         int                              `json:"major_id"`
	SchoolName      string                           `json:"school_name"`
	MajorName       string                           `json:"major_name"`
	InstanceID      int                              `json:"instance_id"`
	InstanceName    string                           `json:"instance_name"`
	PositionID      int                              `json:"position_id"`
	PositionName    string                           `json:"position_name"`
	ExamName        string                           `json:"exam_name"`
	Grade           string                           `json:"grade"`
	TargetType      string                           `json:"target_type"`
	TargetScore     float64                          `json:"target_score"`
	Summary         FirebaseHistoryScoresSummary     `json:"summary"`
	SummaryPTN      FirebaseHistoryScoresSummaryPTN  `json:"summary_ptn"`
	Result          []ProfileStudentResultHistory    `json:"result"`
	ResultPTN       []ProfileStudentResultHistoryPTN `json:"result_ptn"`
	NewScorePTN     CreateHistoryPtn                 `json:"history_ptn"`
	Target          FirebaseTarget                   `json:"target"`
	TargetCPNS      FirebaseTargetCPNS               `json:"target_cpns"`
	Proficiency     string                           `json:"proficiency"`
	ExamProficiency string                           `json:"exam_proficiency"`
	StagesDone      int                              `json:"stages_done"`
	TotalStages     int                              `json:"total_stages"`
	UKAPassed       int                              `json:"uka_passed"`
	Total           float64                          `json:"total"`
}

type FirebaseHistoryScoresSummary struct {
	ScoreKeys                 []string                         `json:"score_keys"`
	ScoreValues               FirebaseHistoryScoresScoreValues `json:"score_values"`
	LatestTotalScore          float64                          `json:"latest_total_score"`
	LatestTotalPercentScore   float64                          `json:"latest_total_percent_score"`
	CurrentTargetTotalScore   float64                          `json:"current_target_total_score"`
	CurrentTargetPercentScore float64                          `json:"current_target_percent_score"`
}
type FirebaseHistoryScoresSummaryPTN struct {
	ScoreKeys                 []string                            `json:"score_keys"`
	ScoreValues               FirebaseHistoryScoresScoreValuesPTN `json:"score_values"`
	LatestTotalScore          float64                             `json:"latest_total_score"`
	LatestTotalPercentScore   float64                             `json:"latest_total_percent_score"`
	CurrentTargetTotalScore   float64                             `json:"current_target_total_score"`
	CurrentTargetPercentScore float64                             `json:"current_target_percent_score"`
}

type FirebaseHistoryScoresScoreValues struct {
	Twk   FirebaseHistoryScoresScoreValue `json:"TWK"`
	Tiu   FirebaseHistoryScoresScoreValue `json:"TIU"`
	Tkp   FirebaseHistoryScoresScoreValue `json:"TKP"`
	Total FirebaseHistoryScoresScoreValue `json:"Total"`
}

type FirebaseHistoryScoresScoreValuesPTN struct {
	LBING FirebaseHistoryScoresScoreValue `json:"LBING"`
	LBIN  FirebaseHistoryScoresScoreValue `json:"LBIND"`
	PBM   FirebaseHistoryScoresScoreValue `json:"PBM"`
	PK    FirebaseHistoryScoresScoreValue `json:"PK"`
	PM    FirebaseHistoryScoresScoreValue `json:"PM"`
	PU    FirebaseHistoryScoresScoreValue `json:"PU"`
	PPU   FirebaseHistoryScoresScoreValue `json:"PPU"`
	Total FirebaseHistoryScoresScoreValue `json:"Total"`
}

type FirebaseHistoryScoresScoreValue struct {
	AvgScore                        float64 `json:"avg_score"`
	AvgPercentScore                 float64 `json:"avg_percent_score"`
	ExplanationRecordTime           uint    `json:"explanation_record_time"`
	ExplanationApproached           uint    `json:"explanation_approached"`
	ExplanationQuestionsTotal       uint    `json:"explanation_questions_total"`
	ExplanationApproachedPercentage float64 `json:"explanation_approached_percentage"`
	CategoryAttemptTime             uint    `json:"category_attempt_time"`
}

type FirebaseTarget struct {
	Location            *string `json:"location"`
	MajorID             uint    `json:"major_id"`
	MajorName           string  `json:"major_name"`
	MaximumScore        int     `json:"maximum_score"`
	PolbitCompetitionID *int    `json:"polbit_competition_id"`
	PolbitLocationID    *int    `json:"polbit_location_id"`
	PolbitType          *string `json:"polbit_type"`
	SchoolID            uint    `json:"school_id"`
	SchoolName          string  `json:"school_name"`
	SchoolLogo          *string `json:"school_logo"`
	TargetScore         int     `json:"target_score"`
	Type                string  `json:"type"`
	MajorChances        *string `json:"major_chances"`
	MajorCompYear       *int    `json:"major_comp_year"`
	MajorQuota          *int    `json:"major_quota"`
	MajorQuotaYear      *int    `json:"major_quota_year"`
	MajorReqistrant     *int    `json:"major_registrant"`
}

type FirebaseTargetCPNS struct {
	CompetitionID          uint    `json:"competition_id"`
	FormationCode          string  `json:"formation_code"`
	FormationLocation      string  `json:"formation_location"`
	FormationType          string  `json:"formation_type"`
	InstanceID             uint    `json:"instance_id"`
	InstanceLogo           string  `json:"instance_logo"`
	InstanceName           string  `json:"instance_name"`
	MaximumScore           int     `json:"maximum_score"`
	PositionID             uint    `json:"position_id"`
	PositionName           string  `json:"position_name"`
	TargetChancePercentage float32 `json:"target_chance_percentage"`
	TargetScore            int     `json:"target_score"`
	TargetType             string  `json:"target_type"`
	Type                   string  `json:"type"`
	MajorChances           *string `json:"major_chances"`
	MajorCompYear          *int    `json:"major_comp_year"`
	MajorQuota             *int    `json:"major_quota"`
	MajorQuotaYear         *int    `json:"major_quota_year"`
	MajorReqistrant        *int    `json:"major_registrant"`
}

type StudentLearningRecordHistory struct {
	TWKLearningDuration   uint `json:"twk_learning_duration"`
	TWKLearningApproached uint `json:"twk_learning_approached"`
	TIULearningDuration   uint `json:"tiu_learning_duration"`
	TIULearningApproached uint `json:"tiu_learning_approached"`
	TKPLearningDuration   uint `json:"tkp_learning_duration"`
	TKPLearningApproached uint `json:"tkp_learning_approached"`
}
