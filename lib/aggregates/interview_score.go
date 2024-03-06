package aggregates

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSingleInterviewScoreBySMIDAndYear(studentId int, year int) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", studentId}},
						bson.M{"$eq": bson.A{"$year", year}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "students",
				"localField":   "smartbtw_id",
				"foreignField": "smartbtw_id",
				"as":           "student",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$student",
			},
		},
		{
			"$addFields": bson.M{
				"name": bson.M{
					"$cond": bson.M{
						"if":   "$student.name",
						"then": "$student.name",
						"else": "$$REMOVE",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"student": 0,
			},
		},
	}
}

func GetInterviewScoresByInterviewSessionIDAndSSOID(interviewSessionID primitive.ObjectID, ssoID string) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$session_id", interviewSessionID}},
						bson.M{"$eq": bson.A{"$created_by.id", ssoID}},
						bson.M{"$eq": bson.A{"$year", time.Now().Year()}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"a_score": bson.M{
					"$sum": bson.A{
						"$penampilan",
						"$cara_duduk_dan_berjabat",
						"$praktek_baris_berbaris",
						"$penampilan_sopan_santun",
						"$kepercayaan_diri_dan_stabilitas_emosi",
						"$komunikasi",
						"$pengembangan_diri",
						"$integritas",
						"$kerjasama",
						"$mengelola_perubahan",
						"$perekat_bangsa",
						"$pelayanan_publik",
						"$pengambilan_keputusan",
						"$orientasi_hasil",
					},
				},
				"b_score": bson.M{
					"$sum": bson.A{
						"$prestasi_akademik",
						"$prestasi_non_akademik",
						"$bahasa_asing",
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"profil_dan_potensi_calon_taruna": bson.M{
					"$round": bson.A{
						bson.M{"$multiply": bson.A{"$a_score", 0.15}},
						2,
					},
				},
				"prestasi_dan_kemampuan_bahasa_asing": bson.M{
					"$round": bson.A{
						bson.M{"$multiply": bson.A{"$b_score", 0.10}},
						2,
					},
				},
				"a_score": "$$REMOVE",
				"b_score": "$$REMOVE",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "students",
				"localField":   "smartbtw_id",
				"foreignField": "smartbtw_id",
				"as":           "student",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$student",
			},
		},
		{
			"$addFields": bson.M{
				"name": bson.M{
					"$cond": bson.M{
						"if":   "$student.name",
						"then": "$student.name",
						"else": "$$REMOVE",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"student": 0,
			},
		},
	}
}

func GetSingleInterviewScoreByID(interviewScoreID primitive.ObjectID) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$_id", interviewScoreID}},
						bson.M{"$eq": bson.A{"$year", time.Now().Year()}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "students",
				"localField":   "smartbtw_id",
				"foreignField": "smartbtw_id",
				"as":           "student",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$student",
			},
		},
		{
			"$addFields": bson.M{
				"name": bson.M{
					"$cond": bson.M{
						"if":   "$student.name",
						"then": "$student.name",
						"else": "$$REMOVE",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"student": 0,
			},
		},
	}
}

func GetInterviewAverageScoresByArrayOfStudentIDAndYear(ids []int, year uint16) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.A{
					bson.M{"$in": bson.A{"$smartbtw_id", bson.A{ids}}},
					bson.M{"$eq": bson.A{"$year", year}},
					bson.M{"$eq": bson.A{"$deleted_at", nil}},
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$smartbtw_id",
				"penampilan": bson.M{
					"$avg": "$penampilan",
				},
				"cara_duduk_dan_berjabat": bson.M{
					"$avg": "$cara_duduk_dan_berjabat",
				},
				"praktek_baris_berbaris": bson.M{
					"$avg": "$praktek_baris_berbaris",
				},
				"penampilan_sopan_santun": bson.M{
					"$avg": "$penampilan_sopan_santun",
				},
				"kepercayaan_diri_dan_stabilitas_emosi": bson.M{
					"$avg": "$kepercayaan_diri_dan_stabilitas_emosi",
				},
				"komunikasi": bson.M{
					"$avg": "$komunikasi",
				},
				"pengembangan_diri": bson.M{
					"$avg": "$pengembangan_diri",
				},
				"integritas": bson.M{
					"$avg": "$integritas",
				},
				"kerjasama": bson.M{
					"$avg": "$kerjasama",
				},
				"mengelola_perubahan": bson.M{
					"$avg": "$mengelola_perubahan",
				},
				"perekat_bangsa": bson.M{
					"$avg": "$perekat_bangsa",
				},
				"pelayanan_publik": bson.M{
					"$avg": "$pelayanan_publik",
				},
				"pengambilan_keputusan": bson.M{
					"$avg": "$pengambilan_keputusan",
				},
				"orientasi_hasil": bson.M{
					"$avg": "$orientasi_hasil",
				},
				"prestasi_akademik": bson.M{
					"$avg": "$prestasi_akademik",
				},
				"prestasi_non_akademik": bson.M{
					"$avg": "$prestasi_non_akademik",
				},
				"bahasa_asing": bson.M{
					"$avg": "$bahasa_asing",
				},
				"final_score": bson.M{
					"$avg": "$final_score",
				},
			},
		},
		{
			"$addFields": bson.M{
				"_id":         "$$REMOVE",
				"smartbtw_id": "$_id",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "students",
				"localField":   "smartbtw_id",
				"foreignField": "smartbtw_id",
				"as":           "student",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$student",
			},
		},
		{
			"$addFields": bson.M{
				"name":    "$student.name",
				"student": "$$REMOVE",
			},
		},
	}
}
