package request

import (
	"time"

	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateHistoryPtn struct {
	SmartBtwID              int       `json:"smartbtw_id"`
	TaskID                  int       `json:"task_id"`
	PackageID               int       `json:"package_id"`
	PackageType             string    `json:"package_type"`
	ModuleCode              string    `json:"module_code"`
	ModuleType              string    `json:"module_type"`
	PotensiKognitif         float64   `json:"potensi_kognitif"`
	PenalaranMatematika     float64   `json:"penalaran_matematika"`
	LiterasiBahasaIndonesia float64   `json:"literasi_bahasa_indonesia"`
	LiterasiBahasaInggris   float64   `json:"literasi_bahasa_inggris"`
	PenalaranUmum           float64   `json:"penalaran_umum"`
	PengetahuanUmum         float64   `json:"pengetahuan_umum"`
	PemahamanBacaan         float64   `json:"pemahaman_bacaan"`
	PengetahuanKuantitatif  float64   `json:"pengetahuan_kuantitatif"`
	ProgramKey              string    `json:"program_key"`
	Total                   float64   `json:"total"`
	Repeat                  int       `json:"repeat"`
	ExamName                string    `json:"exam_name"`
	Grade                   string    `json:"grade"`
	TargetID                string    `json:"target_id"`
	Start                   time.Time `json:"start"`
	End                     time.Time `json:"end"`
	IsLive                  bool      `json:"is_live"`
	Title                   string    `json:"title"`
	StudentName             string    `json:"student_name"`
	SchoolOriginID          string    `json:"school_origin_id"`
	SchoolOrigin            string    `json:"school_origin"`
	SchoolID                int       `json:"school_id"`
	MajorID                 int       `json:"major_id"`
	SchoolName              string    `json:"school_name"`
	MajorName               string    `json:"major_name"`
	TargetScore             float64   `json:"target_score" bson:"target_score"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type CreateHistoryPtnRanking struct {
	SmartBtwID              int        `json:"smartbtw_id"`
	TaskID                  int        `json:"task_id"`
	PackageID               int        `json:"package_id"`
	PackageType             string     `json:"package_type"`
	ModuleCode              string     `json:"module_code"`
	ModuleType              string     `json:"module_type"`
	PotensiKognitif         float64    `json:"potensi_kognitif"`
	PenalaranMatematika     float64    `json:"penalaran_matematika"`
	LiterasiBahasaIndonesia float64    `json:"literasi_bahasa_indonesia"`
	LiterasiBahasaInggris   float64    `json:"literasi_bahasa_inggris"`
	PenalaranUmum           float64    `json:"penalaran_umum"`
	PengetahuanUmum         float64    `json:"pengetahuan_umum"`
	PemahamanBacaan         float64    `json:"pemahaman_bacaan"`
	PengetahuanKuantitatif  float64    `json:"pengetahuan_kuantitatif"`
	ProgramKey              string     `json:"program_key"`
	Total                   float64    `json:"total"`
	Repeat                  int        `json:"repeat"`
	ExamName                string     `json:"exam_name"`
	Grade                   string     `json:"grade"`
	TargetID                string     `json:"target_id"`
	Start                   *time.Time `json:"start"`
	End                     *time.Time `json:"end"`
	IsLive                  bool       `json:"is_live"`
	StudentName             string     `json:"student_name"`
	SchoolOriginID          string     `json:"school_origin_id"`
	SchoolOrigin            string     `json:"school_origin"`
	SchoolID                int        `json:"school_id"`
	MajorID                 int        `json:"major_id"`
	SchoolName              string     `json:"school_name"`
	MajorName               string     `json:"major_name"`
	TargetScore             float64    `json:"target_score" bson:"target_score"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

type UpdateScoreOnlyPTN struct {
	SmartBtwID                            int     `json:"smartbtw_id"`
	AveragePK                             float64 `json:"pk_avg_score"`
	AveragePenalaranMatematika            float64 `json:"pm_avg_score"`
	AverageLiterasiBahasaIndonesia        float64 `json:"lbi_avg_score"`
	AverageLiterasiBahasaInggris          float64 `json:"lbing_avg_score"`
	AveragePercentPK                      float64 `json:"pk_avg_percent_score"`
	AveragePercentPenalaranMatematika     float64 `json:"pm_avg_percent_score"`
	AveragePercentLiterasiBahasaIndonesia float64 `json:"lbi_avg_percent_score"`
	AveragePercentLiterasiBahasaInggris   float64 `json:"lbing_avg_percent_score"`
	AveragePenalaranUmum                  float64 `json:"pu_avg_score"`
	AveragePengetahuanUmum                float64 `json:"ppu_avg_score"`
	AveragePemahamanBacaan                float64 `json:"pbm_avg_score"`
	AveragePercentPenalaranUmum           float64 `json:"pu_avg_percent_score"`
	AveragePercentPengetahuanUmum         float64 `json:"ppu_avg_percent_score"`
	AveragePercentPemahamanBacaan         float64 `json:"pbm_avg_percent_score"`
	AveragePercentTotal                   float64 `json:"tt_avg_percent_score"`
	AverageTotal                          float64 `json:"total_avg_score"`
	LatestTotal                           float64 `json:"latest_total_score"`
	LatestTotalPercent                    float64 `json:"latest_total_score_percent"`
	ProgramKey                            string  `json:"program_key"`
	ModuleDone                            int     `json:"module_done"`
}

type UpdateHistoryPtn struct {
	ID                      int                `json:"id" bson:"id"`
	SmartBtwID              int                `json:"smartbtw_id"`
	TaskID                  int                `json:"task_id"`
	PackageID               int                `json:"package_id"`
	PackageType             string             `json:"package_type"`
	ModuleCode              string             `json:"module_code"`
	ModuleType              string             `json:"module_type"`
	PotensiKognitif         float64            `json:"potensi_kognitif"`
	PenalaranMatematika     float64            `json:"penalaran_matematika"`
	LiterasiBahasaIndonesia float64            `json:"literasi_bahasa_indonesia"`
	LiterasiBahasaInggris   float64            `json:"literasi_bahasa_inggris"`
	PenalaranUmum           float64            `json:"penalaran_umum"`
	PengetahuanUmum         float64            `json:"pengetahuan_umum"`
	PemahamanBacaan         float64            `json:"pemahaman_bacaan"`
	PengetahuanKuantitatif  float64            `json:"pengetahuan_kuantitatif"`
	ProgramKey              string             `json:"program_key"`
	Total                   float64            `json:"total"`
	Repeat                  int                `json:"repeat"`
	ExamName                string             `json:"exam_name"`
	Grade                   string             `json:"grade"`
	TargetID                primitive.ObjectID `json:"target_id"`
}

type HistoryPTNQueryParams struct {
	Limit      *int    `json:"limit"`
	ProgramKey *string `json:"program_key"`
}

type HistoryPtnGetAll struct {
	SmartbtwID      int    `json:"smartbtw_id"`
	Limit           *int64 `json:"limit"`
	Page            *int64 `json:"page"`
	FromElastic     bool   `json:"use_elastic"`
	IsStagesHistory bool   `json:"is_stages_history"`
	ProgramKey      string `json:"program_key"`
}

type MessageHistoryPtnBody struct {
	Version int              `json:"version"`
	Data    CreateHistoryPtn `json:"data" valid:"required"`
}

func UnmarshalMessageHistoryPtnBody(data []byte) (MessageHistoryPtnBody, error) {
	var decoded MessageHistoryPtnBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
