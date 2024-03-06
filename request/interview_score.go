package request

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpsertInterviewScore struct {
	SmartBtwID                        int                            `json:"smartbtw_id" bson:"smartbtw_id" valid:"required~ID Siswa harus diisi"`
	Penampilan                        float64                        `json:"penampilan" bson:"penampilan" valid:"required~Nilai Penampilan harus diisi"`
	CaraDudukDanBerjabat              float64                        `json:"cara_duduk_dan_berjabat" bson:"cara_duduk_dan_berjabat" valid:"required~Nilai Cara Duduk Dan Berjabat harus diisi"`
	PraktekBarisBerbaris              float64                        `json:"praktek_baris_berbaris" bson:"praktek_baris_berbaris" valid:"required~Nilai Praktek Baris Berbaris harus diisi"`
	PenampilanSopanSantun             float64                        `json:"penampilan_sopan_santun" bson:"penampilan_sopan_santun" valid:"required~Nilai Penampilan Sopan Santun"`
	KepercayaanDiriDanStabilitasEmosi float64                        `json:"kepercayaan_diri_dan_stabilitas_emosi" bson:"kepercayaan_diri_dan_stabilitas_emosi" valid:"required~Nilai Kepercayaan Diri Dan Stabilitas Emosi harus diisi"`
	Komunikasi                        float64                        `json:"komunikasi" bson:"komunikasi" valid:"required~Nilai Komunikasi harus diisi"`
	PengembanganDiri                  float64                        `json:"pengembangan_diri" bson:"pengembangan_diri" valid:"required~Nilai Pengembangan Diri harus diisi"`
	Integritas                        float64                        `json:"integritas" bson:"integritas" valid:"required~Nilai Integritas harus diisi"`
	Kerjasama                         float64                        `json:"kerjasama" bson:"kerjasama" valid:"required~Nilai Kerjasama harus diisi"`
	MengelolaPerubahan                float64                        `json:"mengelola_perubahan" bson:"mengelola_perubahan" valid:"required~Nilai Mengelola Perubahan harus diisi"`
	PerekatBangsa                     float64                        `json:"perekat_bangsa" bson:"perekat_bangsa" valid:"required~Nilai Perekat Bangsa harus diisi"`
	PelayananPublik                   float64                        `json:"pelayanan_publik" bson:"pelayanan_publik" valid:"required~Nilai Pelayanan Publik harus diisi"`
	PengambilanKeputusan              float64                        `json:"pengambilan_keputusan" bson:"pengambilan_keputusan" valid:"required~Nilai Pengambilan Keputusan harus diisi"`
	OrientasiHasil                    float64                        `json:"orientasi_hasil" bson:"orientasi_hasil" valid:"required~Nilai Orientasi Hasil harus diisi"`
	PrestasiAkademik                  float64                        `json:"prestasi_akademik" bson:"prestasi_akademik" valid:"required~Nilai Prestasi Akademik harus diisi"`
	PrestasiNonAkademik               float64                        `json:"prestasi_non_akademik" bson:"prestasi_non_akademik" valid:"required~Nilai Prestasi Non Akademik harus diisi"`
	BahasaAsing                       float64                        `json:"bahasa_asing" bson:"bahasa_asing" valid:"required~Nilai Bahasa Asing harus diisi"`
	BersediaPindahJurusan             bool                           `json:"bersedia_pindah_jurusan" bson:"bersedia_pindah_jurusan"`
	SessionID                         primitive.ObjectID             `json:"session_id" bson:"session_id" valid:"required~Sesi harus diisi"`
	FinalScore                        float64                        `json:"final_score" bson:"final_score" valid:"required~Nilai Akhir harus diisi"`
	Year                              uint16                         `json:"year" bson:"year" valid:"required~Tahun harus diisi"`
	Note                              *string                        `json:"note" bson:"note"`
	ClosingStatement                  bool                           `json:"closing_statement" bson:"closing_statement"`
	CreatedBy                         InterviewScoreCreatedUpdatedBy `json:"created_by" bson:"created_by" valid:"required~Data pengisi data nilai wawancara harus diisi"`
}

type GetInterviewScoreByArrEmail struct {
	Email []string `json:"email" bson:"email" valid:"required~Email harus diisi"`
	Year  uint16   `json:"year" bson:"year" valid:"required~Tahun harus diisi"`
}

type InterviewScoreCreatedUpdatedBy struct {
	ID   string `json:"id" bson:"id" valid:"required~ID User Pengisi Nilai Wawancara harus diisi"`
	Name string `json:"name" bson:"name" valid:"required~Nama User Pengisi Nilai Wawancara harus diisi"`
}
