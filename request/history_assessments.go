package request

import (
	"time"

	"github.com/bytedance/sonic"
	"smartbtw.com/services/profile/models"
)

type CreateAssessmentScreening struct {
	SmartBtwID      int                    `json:"smartbtw_id" valid:"required"`
	PackageID       int                    `json:"package_id" valid:"required"`
	AssessmentCode  string                 `json:"assessment_code" valid:"required"`
	Bio             AssessmentScreeningBio `json:"bio" valid:"required"`
	ScreeningTarget any                    `json:"screening_target"`
	ProgramType     string                 `json:"program_type" valid:"required"`
}

type AssessmentScreening struct {
	SmartBtwID      int                    `json:"smartbtw_id"`
	PackageID       int                    `json:"package_id"`
	AssessmentCode  string                 `json:"assessment_code"`
	Bio             AssessmentScreeningBio `json:"bio"`
	ScreeningTarget any                    `json:"screening_target"`
	ProgramType     string                 `json:"program_type"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type AssessmentScreeningTargetPTK struct {
	SchoolID            int     `json:"school_id"`
	MajorID             int     `json:"major_id"`
	SchoolName          string  `json:"school_name"`
	MajorName           string  `json:"major_name"`
	PolbitType          string  `json:"polbit_type"`
	PolbitCompetitionID int     `json:"polbit_competition_id"`
	PolbitLocationID    int     `json:"polbit_location_id"`
	TargetScore         float64 `json:"target_score" bson:"target_score"`
	DomicileProvince    string  `json:"domicile_province" bson:"domicile_province"`
	DomicileRegion      string  `json:"domicile_region" bson:"domicile_region"`
	DomicileProvinceId  int     `json:"domicile_province_id" bson:"domicile_province_id"`
	DomicileRegionId    int     `json:"domicile_region_id" bson:"domicile_region_id"`
}

type AssessmentScreeningTargetPTN struct {
	SchoolID   int    `json:"school_id"`
	MajorID    int    `json:"major_id"`
	SchoolName string `json:"school_name"`
	MajorName  string `json:"major_name"`
}

type AssessmentScreeningTargetCPNS struct {
	InstanceName      string  `json:"instance_name"`
	InstanceID        int     `json:"instance_id"`
	PositionName      string  `json:"position_name"`
	PositionID        int     `json:"position_id"`
	CompetitionCpnsID uint    `json:"competition_id"`
	TargetScore       float64 `json:"target_score"`
	FormationType     string  `json:"formation_type"`
	FormationCode     string  `json:"formation_code"`
	FormationLocation string  `json:"formation_location"`
}

type AssessmentScreeningBio struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
	Origin    string `json:"origin"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type CreateHistoryAssessment struct {
	SmartBtwID     int                               `json:"smartbtw_id"`
	PackageID      int                               `json:"package_id"`
	AssessmentCode string                            `json:"assessment_code"`
	ModuleCode     string                            `json:"module_code"`
	ModuleType     string                            `json:"module_type"`
	PackageType    string                            `json:"package_type"`
	ScoreType      string                            `json:"score_type"`
	Total          float64                           `json:"total"`
	IsPass         bool                              `json:"is_pass"`
	ExamName       string                            `json:"exam_name"`
	Start          time.Time                         `json:"start"`
	End            time.Time                         `json:"end"`
	IsLive         bool                              `json:"is_live"`
	StudentName    string                            `json:"student_name"`
	StudentEmail   string                            `json:"student_email"`
	Program        string                            `json:"program"`
	ProgramType    string                            `json:"program_type"`
	ProgramVersion int                               `json:"program_version"`
	Scores         []models.HistoryAssessmentsScores `json:"scores"`
	CreatedAt      time.Time                         `json:"created_at"`
	UpdatedAt      time.Time                         `json:"updated_at"`
}

type MessageHistoryAssessmentsBody struct {
	Version int                     `json:"version"`
	Data    CreateHistoryAssessment `json:"data" valid:"required"`
}

func UnmarshalMessageHistoryAssessmentsBody(data []byte) (MessageHistoryAssessmentsBody, error) {
	var decoded MessageHistoryAssessmentsBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
