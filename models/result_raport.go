package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResultRaport struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID  int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Program     string             `json:"program" bson:"program"`
	TaskID      int                `json:"task_id" bson:"task_id"`
	Link        string             `json:"link" bson:"link"`
	StudentName string             `json:"student_name" bson:"student_name"`
	ExamName    string             `json:"exam_name" bson:"exam_name"`
	ModuleCode  string             `json:"module_code" bson:"module_code"`
	Score       Score              `json:"score" bson:"score"`
	ScorePTN    ScorePTN           `json:"score_ptn" bson:"score_ptn"`
	StageType   string             `json:"stage_type" bson:"stage_type"`
	ModuleType  string             `json:"module_type" bson:"module_type"`
	PackageType string             `json:"package_type" bson:"package_type"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type GetResultRaportBody struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID  int                `json:"smartbtw_id" bson:"smartbtw_id"`
	Program     string             `json:"program" bson:"program"`
	TaskID      int                `json:"task_id" bson:"task_id"`
	Link        string             `json:"link" bson:"link"`
	StudentName string             `json:"student_name" bson:"student_name"`
	ExamName    string             `json:"exam_name" bson:"exam_name"`
	ModuleCode  string             `json:"module_code" bson:"module_code"`
	StageType   string             `json:"stage_type" bson:"stage_type"`
	ModuleType  string             `json:"module_type" bson:"module_type"`
	PackageType string             `json:"package_type" bson:"package_type"`
}

type Score struct {
	TWK    float64 `json:"twk" bson:"twk"`
	TIU    float64 `json:"tiu" bson:"tiu"`
	TKP    float64 `json:"tkp" bson:"tkp"`
	Total  float64 `json:"total" bson:"total"`
	IsPass bool    `json:"is_pass" bson:"is_pass"`
}

type ScorePTN struct {
	PotensiKognitif         float64 `json:"potensi_kognitif" bson:"potensi_kognitif"`
	PenalaranMatematika     float64 `json:"penalaran_matematika" bson:"penalaran_matematika"`
	LiterasiBahasaIndonesia float64 `json:"literasi_bahasa_indonesia" bson:"literasi_bahasa_indonesia"`
	LiterasiBahasaInggris   float64 `json:"literasi_bahasa_inggris" bson:"literasi_bahasa_inggris"`
	PenalaranUmum           float64 `json:"penalaran_umum" bson:"penalaran_umum"`
	PengetahuanUmum         float64 `json:"pengetahuan_umum" bson:"pengetahuan_umum"`
	PemahamanBacaan         float64 `json:"pemahaman_bacaan" bson:"pemahaman_bacaan"`
	PengetahuanKuantitatif  float64 `json:"pengetahuan_kuantitatif" bson:"pengetahuan_kuantitatif"`
	Total                   float64 `json:"total" bson:"total"`
	IsPass                  bool    `json:"is_pass" bson:"is_pass"`
}
