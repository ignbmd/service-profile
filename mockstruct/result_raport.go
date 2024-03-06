package mockstruct

type Subtest struct {
	CorrectAnswer     int      `json:"correct_answer"`
	EmptyAnswer       int      `json:"empty_answer"`
	PassingIndex      float64  `json:"passing_index"`
	PassingPercentage float64  `json:"passing_percentage"`
	Position          int      `json:"position"`
	Scores            []Scores `json:"scores"`
	SubAlias          string   `json:"sub_alias"`
	SubID             int      `json:"sub_id"`
	SubName           string   `json:"sub_name"`
	WrongAnswer       int      `json:"wrong_answer"`
}

type Scores struct {
	ScoreType string `json:"score_type"`
	Value     int    `json:"value"`
}

type Score struct {
	CategoryAlias     string    `json:"category_alias"`
	CategoryID        int       `json:"category_id"`
	CategoryName      string    `json:"category_name"`
	CorrectAnswer     int       `json:"correct_answer"`
	EmptyAnswer       int       `json:"empty_answer"`
	IsPass            bool      `json:"is_pass"`
	PassingGrade      int       `json:"passing_grade"`
	PassingIndex      float64   `json:"passing_index"`
	PassingPercentage float64   `json:"passing_percentage"`
	Position          int       `json:"position"`
	Scores            []Scores  `json:"scores"`
	Subtests          []Subtest `json:"subtests"`
	WrongAnswer       int       `json:"wrong_answer"`
}

type ScreeningTarget struct {
	DomicileProvince    string `json:"domicile_province"`
	DomicileProvinceID  int    `json:"domicile_province_id"`
	DomicileRegion      string `json:"domicile_region"`
	DomicileRegionID    int    `json:"domicile_region_id"`
	MajorID             int    `json:"major_id"`
	MajorName           string `json:"major_name"`
	PolbitCompetitionID int    `json:"polbit_competition_id"`
	PolbitLocationID    int    `json:"polbit_location_id"`
	PolbitType          string `json:"polbit_type"`
	SchoolID            int    `json:"school_id"`
	SchoolName          string `json:"school_name"`
	TargetScore         int    `json:"target_score"`
}

type ScreeningTargetCPNS struct {
	DomicileProvince   string `json:"domicile_province"`
	DomicileProvinceID int    `json:"domicile_province_id"`
	DomicileRegion     string `json:"domicile_region"`
	DomicileRegionID   int    `json:"domicile_region_id"`
	PositionName       string `json:"position_name"`
	FormationType      string `json:"formation_type"`
	InstanceName       string `json:"instance_name"`
	FormationLocation  string `json:"formation_location"`
	TargetScore        int    `json:"target_score"`
}

type Bio struct {
	BirthDate string `json:"birth_date"`
	Date      string `json:"date"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	Name      string `json:"name"`
	Origin    string `json:"origin"`
	Phone     string `json:"phone"`
}

type Screening struct {
	AssessmentCode  string      `json:"assessment_code"`
	Bio             Bio         `json:"bio"`
	CreatedAt       string      `json:"created_at"`
	PackageID       int         `json:"package_id"`
	ProgramType     string      `json:"program_type"`
	ScreeningTarget interface{} `json:"screening_target"`
	SmartbtwID      int         `json:"smartbtw_id"`
	UpdatedAt       string      `json:"updated_at"`
}

type Report struct {
	File   string `json:"file"`
	Status string `json:"status"`
}

type Assessments struct {
	AssessmentCode        string  `json:"assessment_code"`
	AvgIndex              int     `json:"avg_index"`
	CreatedAt             string  `json:"created_at"`
	Date                  string  `json:"date"`
	DateCertificateSigned string  `json:"date_certificate_signed"`
	End                   string  `json:"end"`
	ExamName              string  `json:"exam_name"`
	IsLive                bool    `json:"is_live"`
	IsPass                bool    `json:"is_pass"`
	ModuleCode            string  `json:"module_code"`
	ModuleType            string  `json:"module_type"`
	PackageID             int     `json:"package_id"`
	TaskID                int     `json:"task_id"`
	PackageType           string  `json:"package_type"`
	Program               string  `json:"program"`
	ProgramType           string  `json:"program_type"`
	ProgramVariant        string  `json:"program_variant"`
	ProgramVersion        int     `json:"program_version"`
	ScoreType             string  `json:"score_type"`
	Scores                []Score `json:"scores"`
	SmartbtwID            int     `json:"smartbtw_id"`
	Start                 string  `json:"start"`
	StudentEmail          string  `json:"student_email"`
	StudentName           string  `json:"student_name"`
	Total                 float64 `json:"total"`
	UpdatedAt             string  `json:"updated_at"`
}
type StudentAnswerCompetition struct {
	CategoryAlias            string  `json:"category_alias"`
	CategoryDisplay          string  `json:"category_display"`
	SubCategoryAlias         string  `json:"sub_category_alias"`
	CategoryName             string  `json:"category_name"`
	SubCategoryName          string  `json:"sub_category_name"`
	CategoryOrder            uint    `json:"category_order"`
	Order                    uint    `json:"order"`
	ChoosenAnswer            string  `json:"choosen_answer"`
	CorrectAnswer            string  `json:"correct_answer"`
	IsTrue                   bool    `json:"is_true"`
	Point                    float64 `json:"point"`
	CorrectStudentPercentage float64 `json:"correct_student_percentage"`
	CorrectStudent           int     `json:"correct_student"`
	TotalStudent             int     `json:"total_student"`
	AnswerType               string  `json:"answer_type"`
	ChoosenMultiAnswerChoice []bool  `json:"choosen_multi_answer_choice"`
	CorrectMultiAnswerChoice []bool  `json:"correct_multi_answer_choice"`
	AnswerHeaderTrue         string  `json:"answer_header_true"`
	AnswerHeaderFalse        string  `json:"answer_header_false"`
	Essay                    string  `json:"essay"`
	AnsweredEssay            string  `json:"answered_essay"`
}

type StudentAnswerCategory struct {
	Category        string                     `json:"category"`
	CategoryAlias   string                     `json:"category_alias"`
	CategoryDisplay string                     `json:"category_display"`
	Order           uint                       `json:"order"`
	Answers         []StudentAnswerCompetition `json:"answers"`
}

type ResultRaportBody struct {
	Assessments                    Assessments                `json:"assessments"`
	RecordAnswer                   []StudentAnswerCompetition `json:"record_answer"`
	RecordAnswerFormatted          []StudentAnswerCategory    `json:"record_answer_formatted"`
	RecordAnswerMapped             []StudentAnswerCategory    `json:"record_answer_mapped"`
	Repeat                         int                        `json:"repeat"`
	Report                         Report                     `json:"report"`
	Screening                      Screening                  `json:"screening"`
	StudentRecommendationFormatted []Recommendation           `json:"student_recommendation_formatted"`
	TransactionID                  int                        `json:"transaction_id"`
	TaskID                         uint                       `json:"task_id"`
}

type Recommendation struct {
	SubCategoryName           string   `json:"sub_category_name"`
	SubCategoryAlias          string   `json:"sub_category_alias"`
	CategoryOrder             int      `json:"category_order"`
	CategoryAlias             string   `json:"category_alias"`
	CategoryName              string   `json:"category_name"`
	PassingIndex              float64  `json:"passing_index"`
	MaterialDescription       string   `json:"material_description"`
	RecommendedTopic          []string `json:"recommended_topic"`
	RecommendedTopicFormatted string   `json:"recommended_topic_formatted"`
}

type ListRaort struct {
	SmartbtwID  uint    `json:"smartbtw_id"`
	StudentName string  `json:"student_name"`
	ExamName    string  `json:"exam_name"`
	ModuleName  string  `json:"module_name"`
	TWK         float64 `json:"twk"`
	TIU         float64 `json:"tiu"`
	TKP         float64 `json:"tkp"`
	Total       float64 `json:"total"`
	IsPass      bool    `json:"is_pass"`
}

type BodyRequestBuildRaport struct {
	SmartbtwID uint   `json:"smartbtw_id"`
	Program    string `json:"program"`
}

type GenerateProgressReportMessage struct {
	SmartbtwID uint   `json:"smartbtw_id"`
	Program    string `json:"program"`
	UKAType    string `json:"uka_type"`
	StageType  string `json:"stage_type"`
}

type ProgressReport struct {
	SmartbtwID int    `json:"smartbtw_id" bson:"smartbtw_id"`
	Program    string `json:"program" bson:"program"`
	Link       string `json:"link" bson:"link"`
	UKAType    string `json:"uka_type" bson:"uka_type"`
	StageType  string `json:"stage_type" bson:"stage_type"`
}
