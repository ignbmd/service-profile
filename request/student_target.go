package request

import (
	"time"

	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateStudentTarget struct {
	SmartbtwID          int     `json:"smartbtw_id"`
	SchoolID            int     `json:"school_id" `
	MajorID             int     `json:"major_id" `
	SchoolName          string  `json:"school_name" `
	MajorName           string  `json:"major_name" `
	PolbitType          string  `json:"polbit_type"`
	PolbitCompetitionID *int    `json:"polbit_competition_id"`
	PolbitLocationID    *int    `json:"polbit_location_id"`
	TargetScore         float64 `json:"target_score"`
	TargetType          string  `json:"target_type"`
	Gender              int     `json:"gender"`
}

type NewUpdateStudentTarget struct {
	SmartbtwID          int     `json:"smartbtw_id"`
	SchoolID            int     `json:"school_id" `
	MajorID             int     `json:"major_id" `
	SchoolName          string  `json:"school_name" `
	MajorName           string  `json:"major_name" `
	TargetScore         float64 `json:"target_score" `
	TargetType          string  `json:"target_type"`
	PolbitType          string  `json:"polbit_type"`
	PolbitCompetitionID *int    `json:"polbit_competition_id"`
	PolbitLocationID    *int    `json:"polbit_location_id"`
	Position            uint    `json:"position"`
	Type                string  `json:"type"`
}

type UpdateStudentTargetRequest struct {
	TargetType        string                    `json:"target_type" valid:"required"`
	StudentTargetList []UpdateStudentTargetBody `json:"student_target_list" valid:"required"`
}

type UpdateStudentTarget struct {
	ID                  primitive.ObjectID `json:"id" bson:"id"`
	SmartbtwID          int                `json:"smartbtw_id"`
	SchoolID            int                `json:"school_id" `
	MajorID             int                `json:"major_id" `
	SchoolName          string             `json:"school_name" `
	MajorName           string             `json:"major_name" `
	PolbitType          string             `json:"polbit_type"`
	PolbitCompetitionID *int               `json:"polbit_competition_id"`
	PolbitLocationID    *int               `json:"polbit_location_id"`
	Position            uint               `json:"position"`
	Type                string             `json:"type"`
	TargetScore         float64            `json:"target_score"`
	TargetType          string             `json:"target_type"`
	FormationYear       int                `json:"formation_year"`
	Gender              int                `json:"gender"`
}

type UpdateStudentTargetBody struct {
	SmartbtwID          int     `json:"smartbtw_id"`
	SchoolID            int     `json:"school_id" `
	MajorID             int     `json:"major_id" `
	SchoolName          string  `json:"school_name" `
	MajorName           string  `json:"major_name" `
	PolbitType          string  `json:"polbit_type"`
	PolbitCompetitionID *int    `json:"polbit_competition_id"`
	PolbitLocationID    *int    `json:"polbit_location_id"`
	Position            uint    `json:"position"`
	Type                string  `json:"type"`
	TargetScore         float64 `json:"target_score"`
	TargetType          string  `json:"target_type"`
	InstanceID          int     `json:"instance_id"`
	InstanceName        string  `json:"instance_name"`
	PositionID          int     `json:"position_id"`
	PositionName        string  `json:"position_name"`
	FormationType       string  `json:"formation_type"`
	FormationLocation   string  `json:"formation_location" bson:"formation_location"`
	FormationYear       int     `json:"formation_year"`
	CompetitionID       int     `json:"competition_id" bson:"competition_id"`
	FormationCode       string  `json:"formation_code" bson:"formation_code"`
	NewestFormation     bool    `json:"newest_formation" bson:"newest_formation"`
}

type UpdateStudentTargetOne struct {
	SmartbtwID          int `json:"smartbtw_id"`
	PolbitCompetitionID int `json:"polbit_competition_id"`
}

type StudentTargetPtkElastic struct {
	SmartbtwID              int     `json:"smartbtw_id"`
	Name                    string  `json:"name"`
	Photo                   string  `json:"photo"`
	SchoolName              string  `json:"school_name"`
	SchoolID                int     `json:"school_id"`
	MajorName               string  `json:"major_name"`
	MajorID                 int     `json:"major_id"`
	PolbitType              string  `json:"polbit_type"`
	PolbitCompetitionID     *int    `json:"polbit_competition_id"`
	PolbitLocationID        *int    `json:"polbit_location_id"`
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
	FormationYear           int     `json:"formation_year"`
}

type StudentTargetPtnElastic struct {
	SmartbtwID              int     `json:"smartbtw_id"`
	Name                    string  `json:"name"`
	Photo                   string  `json:"photo"`
	SchoolID                int     `json:"school_id"`
	SchoolName              string  `json:"school_name"`
	MajorID                 int     `json:"major_id"`
	MajorName               string  `json:"major_name"`
	PolbitType              string  `json:"polbit_type"`
	PolbitCompetitionID     *int    `json:"polbit_competition_id"`
	PolbitLocationID        *int    `json:"polbit_location_id"`
	TargetType              string  `json:"target_type"`
	ProgramKey              string  `json:"program_key"`
	TargetScore             float64 `json:"target_score"`
	ModuleDone              int     `json:"module_done"`
	PkAvgScore              float64 `json:"pk_avg_score"`
	PmAvgScore              float64 `json:"pm_avg_score"`
	LbindAvgScore           float64 `json:"lbind_avg_score"`
	LbingAvgScore           float64 `json:"lbing_avg_score"`
	TotalAvgScore           float64 `json:"total_avg_score"`
	PkAvgPercentScore       float64 `json:"pk_avg_percent_score"`
	PmAvgPercentScore       float64 `json:"pm_avg_percent_score"`
	LbindAvgPercentScore    float64 `json:"lbind_avg_percent_score"`
	LbingAvgPercentScore    float64 `json:"lbing_avg_percent_score"`
	TotalAvgPercentScore    float64 `json:"total_avg_percent_score"`
	LatestTotalScore        float64 `json:"latest_total_score"`
	LatestTotalPercentScore float64 `json:"latest_total_percent_score"`
	Proficiency             string  `json:"proficiency"`
	FormationYear           int     `json:"formation_year"`
}

type GetStudentTargetElasticBody struct {
	SmartbtwID int    `json:"smartbtw_id"`
	TargetType string `json:"target_type"`
	ProgramKey string `json:"program_key"`
}

type UpdateUserData struct {
	SmartbtwID int    `json:"smartbtw_id"`
	Name       string `json:"name"`
	Photo      string `json:"photo"`
}

type UpdateSchool struct {
	SchoolID   int    `json:"school_id"`
	SchoolName string `json:"school_name"`
	TargetType string `json:"target_type"`
}
type UpdateStudyProgram struct {
	MajorID    int    `json:"major_id"`
	MajorName  string `json:"major_name"`
	TargetType string `json:"target_type"`
}

type UpdatBulkPolbit struct {
	StudentTargets []UpdatePolbitType `json:"student_targets"`
}

type UpdatePolbitType struct {
	MajorID             int     `json:"major_id"`
	SmartBTWID          int     `json:"smartbtw_id"`
	TargetType          string  `json:"target_type"`
	PolbitType          string  `json:"polbit_type"`
	PolbitCompetitionID int     `json:"polbit_competition_id"`
	PolbitLocationID    int     `json:"polbit_location_id"`
	TargetScore         float64 `json:"target_score"`
}

type UpdateTargetScore struct {
	MajorID     int     `json:"major_id"`
	TargetScore float64 `json:"target_score"`
	TargetType  string  `json:"target_type"`
}

type UpdateSpecificStudyProgram struct {
	MajorID    []int  `json:"major_id"`
	TargetType string `json:"target_type"`
	MajorName  string `json:"major_name"`
	NewMajorID int    `json:"new_major_id"`
}

type StudentSchoolData struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Logo     *string `json:"logo"`
	Location *string `json:"location"`
}
type ResponseContent struct {
	Data []Data `json:"data"`
}
type ResponseContentPTN struct {
	Data DataPTN `json:"data"`
}
type Data struct {
	ID              int    `json:"id"`
	Quota           int    `json:"quota"`
	Registered      int    `json:"registered"`
	Year            int    `json:"year"`
	CompetitionType string `json:"competition_type"`
}
type DataPTN struct {
	ID              int    `json:"id"`
	Quota           int    `json:"sbmptn_capacity"`
	Registered      int    `json:"registered"`
	Year            int    `json:"year"`
	CompetitionType string `json:"competition_type"`
}

type MajorCompData struct {
	MajorQuota       int     `json:"major_quota"`
	MajorRegistered  int     `json:"major_registered"`
	MajorYear        int     `json:"major_year"`
	MajorQuotaYear   int     `json:"major_quota_year"`
	MajorQuotaChance int     `json:"major_quota_chance"`
	CompetitionType  *string `json:"competition_type"`
}

type GetCompetitonList struct {
	Page    *int    `json:"page" query:"page"`
	PerPage *int    `json:"per_page" query:"per_page"`
	Search  *string `json:"search" query:"search"`
}
type SchoolCompetition struct {
	UKALevel   int     `json:"uka_level"`
	SchoolID   uint    `json:"school_id"`
	SchoolName string  `json:"school_name"`
	MajorID    uint    `json:"major_id"`
	MajorName  string  `json:"major_name"`
	Score      float64 `json:"score"`
	TotalUKA   int     `json:"total_uka"`
}

type Competitions struct {
	SmartbtwID     uint               `json:"smartbtw_id"`
	Name           string             `json:"name"`
	PTKCompetition *SchoolCompetition `json:"ptk"`
	PTNCompetition *SchoolCompetition `json:"ptn"`
}
type CompetitionList struct {
	TotalTargetPTN       int            `json:"total_target_ptn"`
	TotalTargetPTK       int            `json:"total_target_ptk"`
	TotalTargetPTKAndPTN int            `json:"total_target_ptk_and_ptn"`
	Total                int            `json:"total"`
	Competitions         []Competitions `json:"competitions"`
}

type StudentSchoolDataResponse struct {
	Data StudentSchoolData `json:"data"`
}

type StagesStudentLevel struct {
	Level      int    `json:"level"`
	Stage      int    `json:"stage"`
	Type       string `json:"type"`
	SmartbtwID int    `json:"smartbtw_id"`
}

type StudentStage struct {
	Data StagesStudentLevel `json:"data"`
}

type MessageStudentTargetBody struct {
	Version int                 `json:"version"`
	Data    CreateStudentTarget `json:"data" valid:"required"`
}

type MessageUpdateStudentTargetBody struct {
	Version int                 `json:"version"`
	Data    UpdateStudentTarget `json:"data" valid:"required"`
}

type MessageStudentTargetPtkElasticBody struct {
	Version int                     `json:"version"`
	Data    StudentTargetPtkElastic `json:"data" valid:"required"`
}

type MessageStudentTargetPtnElasticBody struct {
	Version int                     `json:"version"`
	Data    StudentTargetPtnElastic `json:"data" valid:"required"`
}

type MessageUpdateUserDataElasticBody struct {
	Version int            `json:"version"`
	Data    UpdateUserData `json:"data" valid:"required"`
}
type MessageUpdateSchoolBody struct {
	Version int          `json:"version"`
	Data    UpdateSchool `json:"data" valid:"required"`
}
type MessageUpdateStudyProgramBody struct {
	Version int                `json:"version"`
	Data    UpdateStudyProgram `json:"data" valid:"required"`
}

type MessageUpdateTargetScoreElasticBody struct {
	Version int               `json:"version"`
	Data    UpdateTargetScore `json:"data" valid:"required"`
}
type DeleteStudentTarget struct {
	ID        int        `json:"id" bson:"id" valid:"type(int), required"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}

type DeleteStudentTargetBodyMessage struct {
	Version int                 `json:"version"`
	Data    DeleteStudentTarget `json:"data" valid:"required"`
}

type MessageUpdateStudentTargetElasticBody struct {
	Version int                            `json:"version"`
	Data    MessageUpdateBulkStudentTarget `json:"data" valid:"required"`
}

type MessageUpdateBulkStudentTarget struct {
	StudentData []*NewUpdateStudentTarget `json:"student_data"`
}

func UnmarshalMessageStudentTargetBody(data []byte) (MessageStudentTargetBody, error) {
	var decoded MessageStudentTargetBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalMessageUpdateStudentTargetBody(data []byte) (MessageUpdateStudentTargetBody, error) {
	var decoded MessageUpdateStudentTargetBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalStudentTargetPtkElastic(data []byte) (MessageStudentTargetPtkElasticBody, error) {
	var decoded MessageStudentTargetPtkElasticBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalStudentTargetPtnElastic(data []byte) (MessageStudentTargetPtnElasticBody, error) {
	var decoded MessageStudentTargetPtnElasticBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalUpdateUserDataElastic(data []byte) (MessageUpdateUserDataElasticBody, error) {
	var decoded MessageUpdateUserDataElasticBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalUpdateSchoolBodyMessage(data []byte) (MessageUpdateSchoolBody, error) {
	var decoded MessageUpdateSchoolBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
func UnmarshalUpdateStudyProgramBodyMessage(data []byte) (MessageUpdateStudyProgramBody, error) {
	var decoded MessageUpdateStudyProgramBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalUpdateTargetScoreElastic(data []byte) (MessageUpdateTargetScoreElasticBody, error) {
	var decoded MessageUpdateTargetScoreElasticBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalDeleteStudentTargetBodyMessage(data []byte) (DeleteStudentTargetBodyMessage, error) {
	var decoded DeleteStudentTargetBodyMessage
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalUpdateStudentTargetElastic(data []byte) (MessageUpdateStudentTargetElasticBody, error) {
	var decoded MessageUpdateStudentTargetElasticBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
