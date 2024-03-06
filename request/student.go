package request

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/mockstruct"
)

type GetStudentsParams struct {
	SmartBTWID []string `query:"smartbtw_id"`
	Fields     string   `query:"fields"`
}

type BodyMessageCreateUpdateStudent struct {
	ID                 int        `json:"id" bson:"id" valid:"type(int), required"`
	Name               string     `json:"name" bson:"name" valid:"type(string), required"`
	Email              string     `json:"email" bson:"email" valid:"type(string), required"`
	Gender             int        `json:"gender" bson:"gender" valid:"optional"`
	BirthDateLocation  *string    `json:"birth_date_location" bson:"birth_date_location" valid:"optional"`
	Phone              *string    `json:"phone" bson:"phone" valid:"optional"`
	SchoolOrigin       *string    `json:"school_origin" bson:"school_origin" valid:"optional"`
	SchoolOriginID     *string    `json:"school_origin_id" bson:"school_origin_id" valid:"optional"`
	Intention          *string    `json:"intention" bson:"intention" valid:"optional"`
	LastEd             *string    `json:"last_ed" bson:"last_ed" valid:"optional"`
	Major              *string    `json:"major" bson:"major" valid:"optional"`
	Profession         *string    `json:"profession" bson:"profession" valid:"optional"`
	Address            *string    `json:"address" bson:"address" valid:"optional"`
	ProvinceId         *int       `json:"province_id" bson:"province_id" valid:"optional"`
	RegionId           *int       `json:"region_id" bson:"region_id" valid:"optional"`
	DomicileProvinceId *int       `json:"domicile_province_id" bson:"domicile_province_id" valid:"optional"`
	DomicileRegionId   *int       `json:"domicile_region_id" bson:"domicile_region_id" valid:"optional"`
	ParentName         *string    `json:"parent_name" bson:"parent_name" valid:"optional"`
	ParentNumber       *string    `json:"parent_number" bson:"parent_number" valid:"optional"`
	Interest           *string    `json:"interest" bson:"interest" valid:"optional"`
	Photo              *string    `json:"photo" bson:"photo" valid:"optional"`
	OriginUniversity   *string    `json:"origin_university" bson:"origin_university" valid:"optional"`
	UserTryoutId       int        `json:"user_tryout_id" bson:"user_tryout_id" valid:"type(int), required"`
	Status             bool       `json:"status" bson:"status" valid:"type(bool)"`
	IsPhoneVerified    bool       `json:"is_phone_verified" bson:"is_phone_verified" valid:"type(bool)"`
	IsEmailVerified    bool       `json:"is_email_verified" bson:"is_email_verified" valid:"type(bool)"`
	IsDataComplete     bool       `json:"is_data_complete" bson:"is_data_complete" valid:"type(bool)"`
	BranchCode         *string    `json:"branch_code" bson:"branch_code" valid:"optional"`
	AffiliateCode      *string    `json:"affiliate_code" bson:"affiliate_code" valid:"optional"`
	AdditionalInfo     *string    `json:"additional_info" bson:"additional_info" valid:"optional"`
	AccountType        string     `json:"account_type" bson:"account_type" valid:"optional"`
	BirthMotherName    *string    `json:"birth_mother_name" bson:"birth_mother_name" valid:"optional"`
	BirthPlace         *string    `json:"birth_place" bson:"birth_place" valid:"optional"`
	NIK                *string    `json:"nik" bson:"nik" valid:"optional"`
	CreatedAt          time.Time  `json:"created_at" bson:"created_at" valid:"type(time.Time)"`
	UpdatedAt          time.Time  `json:"updated_at" bson:"updated_at" valid:"type(time.Time)"`
	DeletedAt          *time.Time `json:"deleted_at" bson:"deleted_at" valid:"optional"`
}

type RequestCreateStudent struct {
	SmartbtwID         int         `json:"smartbtw_id" bson:"smartbtw_id" valid:"type(int), required"`
	Name               string      `json:"name" bson:"name" valid:"type(string), required"`
	Email              string      `json:"email" bson:"email" valid:"type(string), required"`
	Gender             int         `json:"gender" bson:"gender" valid:"type(int), required"`
	BirthDateLocation  string      `json:"birth_date_location" bson:"birth_date_location" valid:"optional"`
	Phone              string      `json:"phone" bson:"phone" valid:"optional"`
	SchoolOrigin       string      `json:"school_origin" bson:"school_origin" valid:"optional"`
	Intention          string      `json:"intention" bson:"intention" valid:"optional"`
	LastEd             string      `json:"last_ed" bson:"last_ed" valid:"optional"`
	Major              string      `json:"major" bson:"major" valid:"optional"`
	Profession         string      `json:"profession" bson:"profession" valid:"optional"`
	Address            string      `json:"address" bson:"address" valid:"optional"`
	ProvinceId         int         `json:"province_id" bson:"province_id" valid:"optional"`
	RegionId           int         `json:"region_id" bson:"region_id" valid:"optional"`
	DomicileProvinceId int         `json:"domicile_province_id" bson:"domicile_province_id" valid:"optional"`
	DomicileRegionId   int         `json:"domicile_region_id" bson:"domicile_region_id" valid:"optional"`
	OriginUniversity   string      `json:"origin_university" bson:"origin_university" valid:"optional"`
	ParentName         string      `json:"parent_name" bson:"parent_name" valid:"optional"`
	ParentNumber       string      `json:"parent_number" bson:"parent_number" valid:"optional"`
	Interest           string      `json:"interest" bson:"interest" valid:"optional"`
	Photo              string      `json:"photo" bson:"photo" valid:"optional"`
	UserTryoutId       int         `json:"user_tryout_id" bson:"user_tryout_id" valid:"type(int), required"`
	Status             bool        `json:"status" bson:"status" valid:"type(bool)"`
	IsPhoneVerified    bool        `json:"is_phone_verified" bson:"is_phone_verified" valid:"type(bool)"`
	IsEmailVerified    bool        `json:"is_email_verified" bson:"is_email_verified" valid:"type(bool)"`
	IsDataComplete     bool        `json:"is_data_complete" bson:"is_data_complete" valid:"type(bool)"`
	BranchCode         string      `json:"branch_code" bson:"branch_code" valid:"optional"`
	AffiliateCode      string      `json:"affiliate_code" bson:"affiliate_code" valid:"optional"`
	AdditionalInfo     string      `json:"additional_info" bson:"additional_info" valid:"optional"`
	BirthMotherName    string      `json:"birth_mother_name" bson:"birth_mother_name" valid:"optional"`
	BirthPlace         string      `json:"birth_place" bson:"birth_place" valid:"optional"`
	NIK                string      `json:"nik" bson:"nik" valid:"optional"`
	CreatedAt          time.Time   `json:"created_at" bson:"created_at" valid:"type(time.Time)"`
	UpdatedAt          time.Time   `json:"updated_at" bson:"updated_at" valid:"type(time.Time)"`
	DeletedAt          interface{} `json:"deleted_at" bson:"deleted_at" valid:"optional"`
}

type StudentProfileBody struct {
	SmartbtwID int `json:"smartbtw_id"`
	// Name           string    `json:"name"`
	// Email          string    `json:"email"`
	// Phone          string    `json:"phone"`
	// Gender         string    `json:"gender"`
	BranchCode         *string   `json:"branch_code"`
	BirthDate          time.Time `json:"birth_date"`
	Province           string    `json:"province"`
	ProvinceID         uint      `json:"province_id"`
	Region             string    `json:"region"`
	RegionID           uint      `json:"region_id"`
	DomicileProvinceId int       `json:"domicile_province_id"`
	DomicileRegionId   int       `json:"domicile_region_id"`
	LastEdID           string    `json:"last_ed_id"`
	LastEdName         string    `json:"last_ed_name"`
	LastEdType         string    `json:"last_ed_type"`
	LastEdMajor        string    `json:"last_ed_major"`
	LastEdMajorID      uint      `json:"last_ed_major_id"`
	LastEdRegion       string    `json:"last_ed_region"`
	LastEdRegionID     uint      `json:"last_ed_region_id"`
	EyeColorBlind      bool      `json:"eye_color_blind"`
	Height             float64   `json:"height"`
	Weight             float64   `json:"weight"`
	// SchoolPTKID    uint      `json:"school_ptk_id"`
	// SchoolNamePTK  string    `json:"school_name_ptk"`
	// MajorNamePTK   string    `json:"major_name_ptk"`
	// MajorPTKID     uint      `json:"major_ptk_id"`
	// CreatedAtPTK   time.Time `json:"created_at_ptk"`
	// SchoolPTNID    uint      `json:"school_ptn_id"`
	// SchoolNamePTN  string    `json:"school_name_ptn"`
	// MajorPTNID     uint      `json:"major_ptn_id"`
	// MajorNamePTN   string    `json:"major_name_ptn"`
	// CreatedAtPTN   time.Time `json:"created_at_ptn"`
	AccountType string `json:"account_type"`
	// CreatedAt      time.Time `json:"created_at"`
}

type GetStudentsCompletedModules struct {
	SmartBTWID int    `query:"smartbtw_id"`
	TargetType string `query:"target_type"`
}

type Emails struct {
	Email string `json:"emails" query:"emails"`
}

type ActivatedOnlineProductFlatten struct {
	SmartBTWID          uint               `json:"smartbtw_id" bson:"smartbtw_id"`
	Name                string             `json:"name" bson:"name"`
	Email               string             `json:"email" bson:"email"`
	Phone               string             `json:"phone" bson:"phone"`
	Status              bool               `json:"status" bson:"status"`
	DateActivated       time.Time          `json:"date_activated" bson:"date_activated"`
	DateExpired         time.Time          `json:"date_expired" bson:"date_expired"`
	ProductID           primitive.ObjectID `json:"product_id" bson:"product_id"`
	ProductTitle        string             `json:"title" bson:"title"`
	ProductTags         pq.StringArray     `json:"tags" bson:"tags"`
	ProductLegacyID     *uint              `json:"legacy_id" bson:"legacy_id"`
	ProductIAPProductID *string            `json:"iap_product_id" bson:"iap_product_id"`
	ProductProgram      string             `json:"program" bson:"program"`
	ProductParentID     *string            `json:"parent_product_id" bson:"parent_product_id"`
	ProductCode         *string            `json:"product_code" bson:"product_code"`
	CreatedAt           time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" bson:"updated_at"`
}

type StudentAOPResults struct {
	Data []ActivatedOnlineProductFlatten `json:"data"`
}

type StagesRegularClass struct {
	ID            primitive.ObjectID `json:"_id,omitempty"`
	Stage         uint               `json:"stage"`
	Level         uint               `json:"level"`
	RequiredStage uint               `json:"required_stage"`
	ModuleType    string             `json:"module_type"`
	PackageID     uint               `json:"package_id"`
	Type          string             `json:"type"`
	IsLocked      bool               `json:"is_locked"`
	StageType     string             `json:"stage_type"`
	StartDate     *time.Time         `json:"start_date"`
	EndDate       *time.Time         `json:"end_date"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	DeletedAt     *time.Time         `json:"deleted_at"`
}

type StudentStagesClassResult struct {
	Data []StagesRegularClass `json:"data"`
}

type GetPerformaSiswa struct {
	SmartBtwID []uint `json:"smartbtw_id"`
	Type       string `json:"type"`
	ClassTags  string `json:"class_tags"`
	ClassYear  *int   `json:"class_year"`
}

type GetPerformaSiswaUKA struct {
	SmartBtwID []uint `json:"smartbtw_id"`
	Program    string `json:"program"`
	TypeStages string `json:"type_stages"`
	TypeModule string `json:"type_module"`
}

type Values struct {
	Total         int     `json:"total"`
	Passed        int     `json:"passed"`
	Failed        int     `json:"failed"`
	IsPass        bool    `json:"is_pass"`
	TotalScore    float32 `json:"total_score"`
	AverageScore  float32 `json:"average_score"`
	PassedPercent float32 `json:"passed_percent"`
	FailedPercent float32 `json:"failed_percent"`
}
type ScoreValues struct {
	IsPass bool   `json:"is_pass"`
	TWK    Values `json:"TWK"`
	TIU    Values `json:"TIU"`
	TKP    Values `json:"TKP"`
}

type ScoreValuesPTN struct {
	PU    Values `json:"PU"`
	PPU   Values `json:"PPU"`
	PBM   Values `json:"PBM"`
	PK    Values `json:"PK"`
	LBIND Values `json:"LBIND"`
	LBING Values `json:"LBING"`
	PM    Values `json:"PM"`
}

type Summary struct {
	AverageScore         float64        `json:"average_score"`
	Passed               int            `json:"passed"`
	Failed               int            `json:"failed"`
	Total                int            `json:"total"`
	TotalScore           float32        `json:"total_score"`
	PassedPercent        float32        `json:"passed_percent"`
	ScoreKeys            []string       `json:"score_keys"`
	DonePercent          float32        `json:"done_percent"`
	Owned                int            `json:"owned"`
	Done                 float32        `json:"done"`
	AverageDone          float32        `json:"average_done"`
	ScoreValues          ScoreValues    `json:"score_values"`
	Deviation            map[string]any `json:"deviation"`
	FinalPassingPercent  float64        `json:"final_passed_percent"`
	TargetPassingPercent float64        `json:"target_passed_percent"`
}

type SummaryPTN struct {
	AverageScore         float64        `json:"average_score"`
	Passed               int            `json:"passed"`
	Failed               int            `json:"failed"`
	Total                int            `json:"total"`
	TotalScore           float32        `json:"total_score"`
	PassedPercent        float32        `json:"passed_percent"`
	ScoreKeys            []string       `json:"score_keys"`
	DonePercent          float32        `json:"done_percent"`
	Owned                int            `json:"owned"`
	Done                 float32        `json:"done"`
	AverageDone          float32        `json:"average_done"`
	ScoreValues          ScoreValuesPTN `json:"score_values"`
	Deviation            map[string]any `json:"deviation"`
	FinalPassingPercent  float64        `json:"final_passed_percent"`
	TargetPassingPercent float64        `json:"target_passed_percent"`
}

type StudentTargetDataPerforma struct {
	SchoolID                  int     `json:"school_id"`
	MajorID                   int     `json:"major_id"`
	SchoolName                string  `json:"school_name"`
	MajorName                 string  `json:"major_name"`
	PolbitType                string  `json:"polbit_type"`
	PolbitCompetitionID       *int    `json:"polbit_competition_id"`
	PolbitLocationID          *int    `json:"polbit_location_id"`
	TargetScore               int     `json:"target_score"`
	CurrentTargetPercentScore float64 `json:"current_target_percent_score"`
}

type StudentTargetDataPerformaCPNS struct {
	SchoolID                  int     `json:"school_id"`
	MajorID                   int     `json:"major_id"`
	SchoolName                string  `json:"school_name"`
	MajorName                 string  `json:"major_name"`
	PolbitType                string  `json:"polbit_type"`
	PolbitCompetitionID       *int    `json:"polbit_competition_id"`
	PolbitLocationID          *int    `json:"polbit_location_id"`
	TargetScore               int     `json:"target_score"`
	CurrentTargetPercentScore float64 `json:"current_target_percent_score"`
}

type StudentClassInformation struct {
	JoinedClass     bool      `json:"joined_class"`
	ClassTitle      string    `json:"class_title"`
	ClassType       string    `json:"class_type"`
	ClassYear       int       `json:"class_year"`
	ClassJoined     time.Time `json:"class_joined"`
	ClassStatus     string    `json:"class_status"`
	ClassBranchCode string    `json:"class_branch_code"`
}

type ResultsPerformaSiswa struct {
	BranchCode       string                    `json:"branch_code"`
	BranchName       string                    `json:"branch_name"`
	Name             string                    `json:"name"`
	Email            string                    `json:"email"`
	SmartBtwID       int                       `json:"smartbtw_id"`
	PTKTarget        StudentTargetDataPerforma `json:"student_target_ptk"`
	Summary          Summary                   `json:"summary"`
	BKNScore         map[string]any            `json:"bkn_score"`
	ClassInformation StudentClassInformation   `json:"class_information"`
	HistoryRecord    []CreateHistoryPtk        `json:"history_record"`
}

type ResultsPerformaSiswaCPNS struct {
	BranchCode       string                        `json:"branch_code"`
	BranchName       string                        `json:"branch_name"`
	Name             string                        `json:"name"`
	Email            string                        `json:"email"`
	SmartBtwID       int                           `json:"smartbtw_id"`
	CPNSTarget       StudentTargetDataPerformaCPNS `json:"student_target_cpns"`
	Summary          Summary                       `json:"summary"`
	BKNScore         map[string]any                `json:"bkn_score"`
	ClassInformation StudentClassInformation       `json:"class_information"`
	HistoryRecord    []CreateHistoryCpns           `json:"history_record"`
}

type ResultsPerformaSiswaPTN struct {
	BranchCode       string                    `json:"branch_code"`
	BranchName       string                    `json:"branch_name"`
	Name             string                    `json:"name"`
	Email            string                    `json:"email"`
	SmartBtwID       int                       `json:"smartbtw_id"`
	PTNTarget        StudentTargetDataPerforma `json:"student_target_ptn"`
	Summary          SummaryPTN                `json:"summary"`
	ClassInformation StudentClassInformation   `json:"class_information"`
	HistoryRecord    []CreateHistoryPtn        `json:"history_record"`
}

type SmartBtwIDArray struct {
	SmartbtwID []uint `json:"smartbtw_id"`
}

type MessageStudentBody struct {
	Version int                            `json:"version"`
	Data    BodyMessageCreateUpdateStudent `json:"data" valid:"required"`
}

type MessageStudentElasticBody struct {
	Version int                `json:"version"`
	Data    StudentProfileBody `json:"data" valid:"required"`
}

type MessageStudentCacheElasticBody struct {
	Version int                                   `json:"version"`
	Data    mockstruct.StudentElasticCacheRequest `json:"data" valid:"required"`
}
type MessageSyncBinsusStudentProfileBody struct {
	Version int                          `json:"version"`
	Data    mockstruct.SyncBinsusProfile `json:"data" valid:"required"`
}
type MessageSyncBinsusFinalStudentProfileBody struct {
	Version int                               `json:"version"`
	Data    mockstruct.SyncBinsusFinalProfile `json:"data" valid:"required"`
}

type DeleteStudentData struct {
	ID        int        `json:"id" bson:"id" valid:"type(int), required"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}

type DeleteStudentBodyMessage struct {
	Version int               `json:"version"`
	Data    DeleteStudentData `json:"data" valid:"required"`
}

type GetStudentByBranchCode struct {
	BranchCode string  `json:"branch_code" query:"branch_code" valid:"type(string), required"`
	Page       int     `json:"page" query:"page" valid:"type(int), optional"`
	Limit      int     `json:"limit" query:"limit" valid:"type(int), optional"`
	Skip       int     `json:"skip" query:"skip" valid:"type(int), optional"`
	Search     *string `json:"search" query:"search" valid:"type(*string), optional"`
}
type GetStudentByBranchCodeArray struct {
	BranchCode string  `json:"branch_code" query:"branch_code" valid:"type(string), required"`
	Page       int     `json:"page" query:"page" valid:"type(int), optional"`
	Limit      int     `json:"limit" query:"limit" valid:"type(int), optional"`
	Skip       int     `json:"skip" query:"skip" valid:"type(int), optional"`
	Search     *string `json:"search" query:"search" valid:"type(*string), optional"`
}

func UnmarshalMessageStudentBody(data []byte) (MessageStudentBody, error) {
	var decoded MessageStudentBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalDeleteStudentBodyMessage(data []byte) (DeleteStudentBodyMessage, error) {
	var decoded DeleteStudentBodyMessage
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalCreateStudentProfileElastic(data []byte) (MessageStudentElasticBody, error) {
	var decoded MessageStudentElasticBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalCacheStudentProfileElastic(data []byte) (MessageStudentCacheElasticBody, error) {
	var decoded MessageStudentCacheElasticBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalSyncBinsusStudentProfileElastic(data []byte) (MessageSyncBinsusStudentProfileBody, error) {
	var decoded MessageSyncBinsusStudentProfileBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
func UnmarshalSyncBinsusFinalStudentProfileElastic(data []byte) (MessageSyncBinsusFinalStudentProfileBody, error) {
	var decoded MessageSyncBinsusFinalStudentProfileBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
