package models

import (
	"time"

	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID                 primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID         int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Name               string             `json:"name" bson:"name"`
	Email              string             `json:"email" bson:"email"`
	Gender             int                `json:"gender" bson:"gender"`
	BirthDateLocation  *string            `json:"birth_date_location" bson:"birth_date_location"`
	Phone              *string            `json:"phone" bson:"phone"`
	SchoolOrigin       *string            `json:"school_origin" bson:"school_origin"`
	SchoolOriginID     *string            `json:"school_origin_id" bson:"school_origin_id"`
	Intention          *string            `json:"intention" bson:"intention"`
	LastEd             *string            `json:"last_ed" bson:"last_ed"`
	Major              *string            `json:"major" bson:"major"`
	Profession         *string            `json:"profession" bson:"profession"`
	Address            *string            `json:"address" bson:"address"`
	ProvinceId         *int               `json:"province_id" bson:"province_id"`
	RegionId           *int               `json:"region_id" bson:"region_id"`
	DomicileProvinceId *int               `json:"domicile_province_id" bson:"domicile_province_id"`
	DomicileRegionId   *int               `json:"domicile_region_id" bson:"domicile_region_id"`
	ParentName         *string            `json:"parent_name" bson:"parent_name"`
	ParentNumber       *string            `json:"parent_number" bson:"parent_number"`
	Interest           *string            `json:"interest" bson:"interest"`
	Photo              *string            `json:"photo" bson:"photo"`
	UserTryoutId       int                `json:"user_tryout_id" bson:"user_tryout_id"`
	Status             bool               `json:"status" bson:"status"`
	IsPhoneVerified    bool               `json:"is_phone_verified" bson:"is_phone_verified"`
	IsEmailVerified    bool               `json:"is_email_verified" bson:"is_email_verified"`
	IsDataComplete     bool               `json:"is_data_complete" bson:"is_data_complete"`
	BranchCode         *string            `json:"branch_code" bson:"branch_code"`
	AffiliateCode      *string            `json:"affiliate_code" bson:"affiliate_code"`
	AdditionalInfo     *string            `json:"additional_info" bson:"additional_info"`
	OriginUniversity   *string            `json:"origin_university" bson:"origin_university"`
	BirthMotherName    *string            `json:"birth_mother_name" bson:"birth_mother_name"`
	BirthPlace         *string            `json:"birth_place" bson:"birth_place"`
	NIK                *string            `json:"nik" bson:"nik"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt          *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type StudentSimpleData struct {
	ID                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID            int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Name                  string             `json:"name" bson:"name"`
	Email                 string             `json:"email" bson:"email"`
	Phone                 *string            `json:"phone" bson:"phone"`
	BranchCode            *string            `json:"branch_code" bson:"branch_code"`
	AccountType           string             `json:"account_type" bson:"account_type"`
	SchoolName            string             `json:"school_name" bson:"school_name"`
	SchoolID              uint               `json:"school_id" bson:"school_id"`
	MajorName             string             `json:"major_name" bson:"major_name"`
	MajorID               uint               `json:"major_id" bson:"major_id"`
	OriginSchoolID        string             `json:"origin_school_id" bson:"origin_school_id"`
	OriginSchoolName      string             `json:"origin_school_name" bson:"origin_school_name"`
	FormationType         string             `json:"formation_type" bson:"formation_type"`
	FormationDesc         string             `json:"formation_desc" bson:"formation_desc"`
	PolbitCompetitionID   uint               `json:"polbit_competition_id" bson:"polbit_competition_id"`
	PolbitCompetitionType string             `json:"polbit_competition_type" bson:"polbit_competition_type"`
	CreatedAt             time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at" bson:"updated_at"`
}

func UnmarshalStudent(data []byte) (Student, error) {
	var r Student
	err := sonic.Unmarshal(data, &r)
	return r, err
}

func (r *Student) Marshal() ([]byte, error) {
	return sonic.Marshal(r)
}
