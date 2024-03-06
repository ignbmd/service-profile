package mockstruct

import "time"

type StudentProfileCompMapRequest struct {
	Data StudentProfileCompMapBody `json:"data"`
}
type SKDRankCompMapRequest struct {
	Data SKDRankCompMapBody `json:"data"`
}

type SKDRankCompMapBody struct {
	StudyProgramPassingGrade uint `json:"study_program_passing_grade"`
}

type StudentProfileCompMapBody struct {
	SmartbtwID         uint                          `json:"smartbtw_id"`
	Name               string                        `json:"name"`
	LastEdID           *uint                         `json:"last_ed_id"`
	LastUniversityEdID *uint                         `json:"last_university_ed_id"`
	LocationID         *uint                         `json:"location_id"`
	BirthDate          *time.Time                    `json:"birth_date"`
	Gender             string                        `json:"gender"`
	Height             float64                       `json:"height"`
	Weight             float64                       `json:"weight"`
	LegX               float64                       `json:"leg_x"`
	LegY               float64                       `json:"leg_y"`
	EyeMin             float64                       `json:"eye_min"`
	EyePlus            float64                       `json:"eye_plus"`
	EyeCylinder        float64                       `json:"eye_cylinder"`
	EyeColorBlind      *bool                         `json:"eye_color_blind"`
	LastEd             StudentProfileCompMapLastEd   `json:"last_ed"`
	LastUniversityEd   StudentProfileCompMapLastEd   `json:"last_university_ed"`
	Location           StudentProfileCompMapLocation `json:"location"`
}

type StudentProfileCompMapLastEd struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type StudentProfileCompMapLocation struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentId *uint  `json:"parent_id"`
	Type     string `json:"type"`
}
type CompetitionCompMapRequest struct {
	Data CompMapCompetition `json:"data"`
}
type CompMapCompetition struct {
	ID             uint   `json:"id"`
	StudyProgramID uint   `json:"study_program_id"`
	LocationID     *int   `json:"location_id"`
	PolbitType     string `json:"polbit_type"`
	Year           uint16 `json:"year"`
	LowestScore    uint16 `json:"lowest_score"`
	StudyProgram   struct {
		ID     uint   `json:"id"`
		Name   string `json:"name"`
		School struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"school"`
	} `json:"study_program"`
	Location struct {
		ID   uint   `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"location"`
}

type CompetitionNewCompMapRequest struct {
	Data []CompMapNewCompetition `json:"data"`
}
type CompMapNewCompetition struct {
	LowestScore uint16 `gorm:"not null" json:"lowest_score"`
}
type CompetitionPTKPTNData struct {
	MajorQuota       uint32  `json:"major_quota"`
	MajorRegistered  uint32  `json:"major_registered"`
	MajorYear        uint32  `json:"major_year"`
	MajorQuotaYear   uint32  `json:"major_quota_year"`
	MajorQuotaChance string  `json:"major_quota_chance"`
	MajorOldQuota    uint32  `json:"major_old_quota"`
	CompetitionType  *string `json:"competition_type"`
}
type CompetitionPTK struct {
	ID              int       `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	StudyProgramID  uint      `gorm:"not null" json:"study_program_id"`
	LocationID      *uint     `gorm:"null" json:"location_id"`
	PolbitType      string    `gorm:"not null" json:"polbit_type"`
	CompetitionType string    `gorm:"not null;default:'ALL'" json:"competition_type"`
	Flag            string    `gorm:"not null;default:''" json:"flag"`
	Year            uint16    `gorm:"not null" json:"year"`
	Quota           uint      `gorm:"not null" json:"quota"`
	Registered      uint      `gorm:"not null" json:"registered"`
	LowestScore     uint16    `gorm:"not null" json:"lowest_score"`
	LowestPosition  uint16    `gorm:"not null" json:"lowest_position"`
	LowestStatus    string    `gorm:"not null" json:"lowest_status"`
	AvgReportScore  uint8     `gorm:"not null;default:0" json:"avg_report_score"`
	AvgDiplomaScore uint8     `gorm:"not null;default:0" json:"avg_diploma_score"`
	StudyProgram    struct {
		ID       int    `json:"id"`
		SchoolID uint   `gorm:"not null" json:"school_id"`
		Name     string `gorm:"not null" json:"name"`
		School   struct {
			ID       int     `json:"id"`
			Name     string  `gorm:"not null" json:"name" valid:"type(string)"`
			SubName  *string `json:"sub_name"  valid:"-"`
			Ministry string  `gorm:"not null" json:"ministry" valid:"type(string)"`
			Address  string  `gorm:"not null" json:"address" valid:"type(string)"`
			Logo     *string `json:"logo" valid:"-"`
			Location *string `json:"location" valid:"-"`
		} `json:"school"`
	} `json:"study_program"`
	Location struct {
		ID       int    `json:"id"`
		Name     string `gorm:"not null" json:"name"`
		ParentId *uint  `gorm:"null" json:"parent_id"`
		Type     string `gorm:"not null" json:"type"`
	} `json:"location"`
}

type CompetitionPTN struct {
	PtnStudyProgramID uint `gorm:"null" json:"ptn_study_program_id"`
	Year              int  `gorm:"null" json:"year"`
	SnmptnCapacity    int  `gorm:"null" json:"snmptn_capacity"`
	SnmptnRegistered  int  `gorm:"null" json:"snmptn_registered"`
	SbmptnCapacity    int  `gorm:"null" json:"sbmptn_capacity"`
	SbmptnRegistered  int  `gorm:"null" json:"sbmptn_registered"`
	Average           int  `gorm:"null" json:"average"`
}
type CompetitionPTNRequest struct {
	Data CompetitionPTN `json:"data"`
}

type GetCompetitionCPNS struct {
	FormationType string `json:"formation_type"`
	PositionID    uint   `json:"position_id"`
	FormationCode string `json:"formation_code"`
}

type CompetitionFormationCPNS struct {
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

type CompetitionFormationCPNSRequest struct {
	Data []CompetitionFormationCPNS `json:"data"`
}
