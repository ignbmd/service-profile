package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InterviewScore struct {
	ID                                primitive.ObjectID              `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartBtwID                        int                             `json:"smartbtw_id" bson:"smartbtw_id"`
	Penampilan                        float64                         `json:"penampilan" bson:"penampilan"`
	CaraDudukDanBerjabat              float64                         `json:"cara_duduk_dan_berjabat" bson:"cara_duduk_dan_berjabat"`
	PraktekBarisBerbaris              float64                         `json:"praktek_baris_berbaris" bson:"praktek_baris_berbaris"`
	PenampilanSopanSantun             float64                         `json:"penampilan_sopan_santun" bson:"penampilan_sopan_santun"`
	KepercayaanDiriDanStabilitasEmosi float64                         `json:"kepercayaan_diri_dan_stabilitas_emosi" bson:"kepercayaan_diri_dan_stabilitas_emosi"`
	Komunikasi                        float64                         `json:"komunikasi" bson:"komunikasi"`
	PengembanganDiri                  float64                         `json:"pengembangan_diri" bson:"pengembangan_diri"`
	Integritas                        float64                         `json:"integritas" bson:"integritas"`
	Kerjasama                         float64                         `json:"kerjasama" bson:"kerjasama"`
	MengelolaPerubahan                float64                         `json:"mengelola_perubahan" bson:"mengelola_perubahan"`
	PerekatBangsa                     float64                         `json:"perekat_bangsa" bson:"perekat_bangsa"`
	PelayananPublik                   float64                         `json:"pelayanan_publik" bson:"pelayanan_publik"`
	PengambilanKeputusan              float64                         `json:"pengambilan_keputusan" bson:"pengambilan_keputusan"`
	OrientasiHasil                    float64                         `json:"orientasi_hasil" bson:"orientasi_hasil"`
	PrestasiAkademik                  float64                         `json:"prestasi_akademik" bson:"prestasi_akademik"`
	PrestasiNonAkademik               float64                         `json:"prestasi_non_akademik" bson:"prestasi_non_akademik"`
	BahasaAsing                       float64                         `json:"bahasa_asing" bson:"bahasa_asing"`
	BersediaPindahJurusan             bool                            `json:"bersedia_pindah_jurusan" bson:"bersedia_pindah_jurusan"`
	FinalScore                        float64                         `json:"final_score" bson:"final_score"`
	Year                              uint16                          `json:"year" bson:"year"`
	SessionID                         primitive.ObjectID              `json:"session_id" bson:"session_id"`
	SessionName                       string                          `json:"session_name" bson:"session_name"`
	SessionDescription                string                          `json:"session_description" bson:"session_description"`
	SessionNumber                     int                             `json:"session_number" bson:"session_number"`
	Note                              *string                         `json:"note" bson:"note"`
	ClosingStatement                  bool                            `json:"closing_statement" bson:"closing_statement"`
	CreatedBy                         InterviewScoreCreatedUpdatedBy  `json:"created_by" bson:"created_by"`
	UpdatedBy                         *InterviewScoreCreatedUpdatedBy `json:"updated_by" bson:"updated_by"`
	CreatedAt                         time.Time                       `json:"created_at" bson:"created_at"`
	UpdatedAt                         time.Time                       `json:"updated_at" bson:"updated_at"`
	DeletedAt                         *time.Time                      `json:"deleted_at" bson:"deleted_at"`
}

type InterviewScoreEmailEdutech struct {
	SmartBtwID     int                    `json:"smartbtw_id" bson:"smartbtw_id"`
	BTWEdutechID   int                    `json:"btwedutech_id" bson:"btwedutech_id"`
	Name           string                 `json:"name" bson:"name"`
	AccountType    string                 `json:"account_type" bson:"account_type"`
	InterviewScore *InterviewAverageScore `json:"interview_score" bson:"interview_score"`
}

type InterviewScoreCreatedUpdatedBy struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type InterviewAverageScore struct {
	SmartBtwID                        int     `json:"smartbtw_id" bson:"smartbtw_id"`
	Name                              string  `json:"name" bson:"name"`
	Penampilan                        float64 `json:"penampilan" bson:"penampilan"`
	CaraDudukDanBerjabat              float64 `json:"cara_duduk_dan_berjabat" bson:"cara_duduk_dan_berjabat"`
	PraktekBarisBerbaris              float64 `json:"praktek_baris_berbaris" bson:"praktek_baris_berbaris"`
	PenampilanSopanSantun             float64 `json:"penampilan_sopan_santun" bson:"penampilan_sopan_santun"`
	KepercayaanDiriDanStabilitasEmosi float64 `json:"kepercayaan_diri_dan_stabilitas_emosi" bson:"kepercayaan_diri_dan_stabilitas_emosi"`
	Komunikasi                        float64 `json:"komunikasi" bson:"komunikasi"`
	PengembanganDiri                  float64 `json:"pengembangan_diri" bson:"pengembangan_diri"`
	Integritas                        float64 `json:"integritas" bson:"integritas"`
	Kerjasama                         float64 `json:"kerjasama" bson:"kerjasama"`
	MengelolaPerubahan                float64 `json:"mengelola_perubahan" bson:"mengelola_perubahan"`
	PerekatBangsa                     float64 `json:"perekat_bangsa" bson:"perekat_bangsa"`
	PelayananPublik                   float64 `json:"pelayanan_publik" bson:"pelayanan_publik"`
	PengambilanKeputusan              float64 `json:"pengambilan_keputusan" bson:"pengambilan_keputusan"`
	OrientasiHasil                    float64 `json:"orientasi_hasil" bson:"orientasi_hasil"`
	PrestasiAkademik                  float64 `json:"prestasi_akademik" bson:"prestasi_akademik"`
	PrestasiNonAkademik               float64 `json:"prestasi_non_akademik" bson:"prestasi_non_akademik"`
	BahasaAsing                       float64 `json:"bahasa_asing" bson:"bahasa_asing"`
	FinalScore                        float64 `json:"final_score" bson:"final_score"`
}
