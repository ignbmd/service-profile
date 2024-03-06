package mockstruct

import "time"

type FetchRankingBase struct {
	RankingInformation struct {
		DataTotal         int64 `json:"data_total"`
		CurrentCountTotal int   `json:"current_count_total"`
		Page              int   `json:"page"`
		PageTotal         int   `json:"page_total"`
	} `json:"ranking_information"`
}

type FetchRankingStudentBase struct {
	Name          string  `json:"name" bson:"name"`
	Email         string  `json:"email" bson:"email"`
	LastEdID      string  `json:"last_ed_id" bson:"last_ed_id"`
	LastEdName    string  `json:"last_ed_name" bson:"last_ed_name"`
	BranchCode    string  `json:"branch_code" bson:"branch_code"`
	BranchName    string  `json:"branch_name" bson:"branch_name"`
	Rank          int     `json:"rank"`
	SchoolID      int     `json:"school_id"`
	MajorID       int     `json:"major_id"`
	SchoolName    string  `json:"school_name"`
	MajorName     string  `json:"major_name"`
	PassingChance float64 `json:"passing_chance"`
	SmartBtwID    int     `json:"smartbtw_id"`
	TaskID        int     `json:"task_id"`
	PackageID     int     `json:"package_id"`
	IsSameSchool  bool    `json:"is_same_school"`
}

type FetchRankingPTK struct {
	FetchRankingStudentBase
	ModuleCode    string    `json:"module_code"`
	ModuleType    string    `json:"module_type"`
	PackageType   string    `json:"package_type"`
	Twk           float64   `json:"twk"`
	Tiu           float64   `json:"tiu"`
	Tkp           float64   `json:"tkp"`
	TwkPassStatus bool      `json:"twk_pass_status"`
	TiuPassStatus bool      `json:"tiu_pass_status"`
	TkpPassStatus bool      `json:"tkp_pass_status"`
	AllPassStatus bool      `json:"is_all_passed"`
	Title         string    `json:"title"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	Total         float64   `json:"total"`
}

type FetchRankingCPNS struct {
	FetchRankingStudentBase
	ModuleCode    string    `json:"module_code"`
	ModuleType    string    `json:"module_type"`
	PackageType   string    `json:"package_type"`
	Twk           float64   `json:"twk"`
	Tiu           float64   `json:"tiu"`
	Tkp           float64   `json:"tkp"`
	TwkPassStatus bool      `json:"twk_pass_status"`
	TiuPassStatus bool      `json:"tiu_pass_status"`
	TkpPassStatus bool      `json:"tkp_pass_status"`
	AllPassStatus bool      `json:"is_all_passed"`
	Title         string    `json:"title"`
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	Total         float64   `json:"total"`
}

type FetchRankingPTN struct {
	FetchRankingStudentBase
	ModuleCode              string    `json:"module_code"`
	ModuleType              string    `json:"module_type"`
	PackageType             string    `json:"package_type"`
	PotensiKognitif         float64   `json:"potensi_kognitif"`
	PenalaranMatematika     float64   `json:"penalaran_matematika"`
	LiterasiBahasaIndonesia float64   `json:"literasi_bahasa_indonesia"`
	LiterasiBahasaInggris   float64   `json:"literasi_bahasa_inggris"`
	PenalaranUmum           float64   `json:"penalaran_umum"`
	PengetahuanUmum         float64   `json:"pengetahuan_umum"`
	PemahamanBacaan         float64   `json:"pemahaman_bacaan"`
	PengetahuanKuantitatif  float64   `json:"pengetahuan_kuantitatif"`
	ProgramKey              string    `json:"program_key"`
	Title                   string    `json:"title"`
	Total                   float64   `json:"total"`
	Start                   time.Time `json:"start"`
	End                     time.Time `json:"end"`
}

type InformationRankingData struct {
	TargetPtn       int `json:"target_ptn"`
	TargetPtk       int `json:"target_ptk"`
	TargetPTNandPTK int `json:"target_ptnptk"`
	Total           int `json:"total"`
}

type FetchRankingPTNBody struct {
	FetchRankingBase
	RankingData []FetchRankingPTN `json:"ranking_data"`
}

type FetchRankingPTKBody struct {
	FetchRankingBase
	RankingData []FetchRankingPTK `json:"ranking_data"`
}

type FetchRankingCPNSBody struct {
	FetchRankingBase
	RankingData []FetchRankingCPNS `json:"ranking_data"`
}

type FetchRankingUkaCodeBody struct {
	FetchRankingBase
	RankingData            []FetchRankingPTK      `json:"ranking_data"`
	RankingDataInformation InformationRankingData `json:"ranking_data_information"`
}

type FetchRankingPTNUkaCodeBody struct {
	FetchRankingBase
	RankingData            []FetchRankingPTN      `json:"ranking_data"`
	RankingDataInformation InformationRankingData `json:"ranking_data_information"`
}
