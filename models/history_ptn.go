package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HistoryPtn struct {
	ID                      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartBtwID              int                `json:"smartbtw_id" bson:"smartbtw_id"`
	TaskID                  int                `json:"task_id" bson:"task_id"`
	PackageID               int                `json:"package_id" bson:"package_id"`
	PackageType             string             `json:"package_type" bson:"package_type"`
	ModuleCode              string             `json:"module_code" bson:"module_code"`
	ModuleType              string             `json:"module_type" bson:"module_type"`
	PotensiKognitif         float64            `json:"potensi_kognitif" bson:"potensi_kognitif"`
	PenalaranMatematika     float64            `json:"penalaran_matematika" bson:"penalaran_matematika"`
	LiterasiBahasaIndonesia float64            `json:"literasi_bahasa_indonesia" bson:"literasi_bahasa_indonesia"`
	LiterasiBahasaInggris   float64            `json:"literasi_bahasa_inggris" bson:"literasi_bahasa_inggris"`
	PenalaranUmum           float64            `json:"penalaran_umum" bson:"penalaran_umum"`
	PengetahuanUmum         float64            `json:"pengetahuan_umum" bson:"pengetahuan_umum"`
	PemahamanBacaan         float64            `json:"pemahaman_bacaan" bson:"pemahaman_bacaan"`
	PengetahuanKuantitatif  float64            `json:"pengetahuan_kuantitatif" bson:"pengetahuan_kuantitatif"`
	ProgramKey              string             `json:"program_key" bson:"program_key"`
	Total                   float64            `json:"total" bson:"total"`
	Repeat                  int                `json:"repeat" bson:"repeat"`
	ExamName                string             `json:"exam_name" bson:"exam_name"`
	Grade                   string             `json:"grade" bson:"grade"`
	IsLive                  bool               `json:"is_live" bson:"is_live"`
	TargetID                primitive.ObjectID `json:"target_id" bson:"target_id"`
	Start                   *time.Time         `json:"start" bson:"start"`
	End                     *time.Time         `json:"end" bson:"end"`
	StudentName             string             `json:"student_name"`
	SchoolOriginID          string             `json:"school_origin_id"`
	SchoolOrigin            string             `json:"school_origin"`
	TargetScore             float64            `json:"target_score" bson:"target_score"`
	SchoolID                int                `json:"school_id"`
	MajorID                 int                `json:"major_id"`
	SchoolName              string             `json:"school_name"`
	MajorName               string             `json:"major_name"`
	CreatedAt               time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt               time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt               *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type PTNStudentRecommendationStruct struct {
	Version uint                     `json:"version"`
	Data    PTNStudentRecommendation `json:"data"`
}
type PTNStudentRecommendation struct {
	SmartbtwID        uint `json:"smartbtw_id"`
	PtnStudyProgramID uint `json:"ptn_study_program_id"`
	Score             int  `json:"score"`
}
