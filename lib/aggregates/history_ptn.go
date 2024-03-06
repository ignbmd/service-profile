package aggregates

import (
	"go.mongodb.org/mongo-driver/bson"
	"smartbtw.com/services/profile/request"
)

func GetStudentPtnLastScore(SmartBtwID int, programKey string) []bson.M {
	return []bson.M{
		{"$sort": bson.M{"created_at": -1}},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", SmartBtwID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
						bson.M{
							"$cond": bson.M{
								"if":   programKey != "",
								"then": bson.M{"$eq": bson.A{"$program_key", programKey}},
								"else": "",
							},
						},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":  bson.A{"$smartbtw_id", SmartBtwID},
				"docs": bson.M{"$push": "$$ROOT"},
			},
		},
		{
			"$project": bson.M{
				"_id":       0,
				"last_exam": bson.M{"$slice": bson.A{"$docs", -1}},
			},
		},
		{"$unwind": bson.M{"path": "$history_ptn", "preserveNullAndEmptyArrays": true}},
	}
}

func GetStudentPtnAverage(SmartBtwID int, programKey string) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", SmartBtwID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
						bson.M{
							"$cond": bson.M{
								"if":   programKey != "",
								"then": bson.M{"$eq": bson.A{"$program_key", programKey}},
								"else": "",
							},
						},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":                       bson.A{"$smartbtw_id", SmartBtwID},
				"potensi_kognitif":          bson.M{"$avg": "$potensi_kognitif"},
				"penalaran_matematika":      bson.M{"$avg": "$penalaran_matematika"},
				"literasi_bahasa_indonesia": bson.M{"$avg": "$literasi_bahasa_indonesia"},
				"literasi_bahasa_inggris":   bson.M{"$avg": "$literasi_bahasa_inggris"},
				"pengetahuan_kuantitatif":   bson.M{"$avg": "$pengetahuan_kuantitatif"},
				"pemahaman_bacaan":          bson.M{"$avg": "$pemahaman_bacaan"},
				"pengetahuan_umum":          bson.M{"$avg": "$pengetahuan_umum"},
				"penalaran_umum":            bson.M{"$avg": "$penalaran_umum"},
				"total":                     bson.M{"$avg": "$total"},
			},
		},
		{
			"$project": bson.M{
				"_id":                       0,
				"potensi_kognitif":          1,
				"penalaran_matematika":      1,
				"literasi_bahasa_indonesia": 1,
				"literasi_bahasa_inggris":   1,
				"pengetahuan_kuantitatif":   1,
				"pemahaman_bacaan":          1,
				"pengetahuan_umum":          1,
				"penalaran_umum":            1,
				"total":                     1,
			},
		},
	}
}

func GetLast10StudentPtnScore(SmartBtwID int, programKey string) []bson.M {
	return []bson.M{
		{"$sort": bson.M{"created_at": 1}},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", SmartBtwID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
						bson.M{"$not": bson.M{"$in": bson.A{"$module_type", bson.A{"PRE_TEST", "POST_TEST"}}}},
						bson.M{
							"$cond": bson.M{
								"if":   programKey != "",
								"then": bson.M{"$eq": bson.A{"$program_key", programKey}},
								"else": "",
							},
						},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":  bson.A{"$smartbtw_id", SmartBtwID},
				"docs": bson.M{"$push": "$$ROOT"},
			},
		},
		{
			"$project": bson.M{
				"_id":       0,
				"last_exam": bson.M{"$slice": bson.A{"$docs", -10}},
			},
		},
		{"$unwind": bson.M{"path": "$history_ptn", "preserveNullAndEmptyArrays": true}},
	}
}

func GetStudentPTNHistoryScores(SmartBTWID int, params *request.HistoryPTNQueryParams) []bson.M {
	pipelines := []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", SmartBTWID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
						bson.M{
							"$cond": bson.M{
								"if": params.ProgramKey != nil,
								"then": bson.M{"$or": bson.A{
									bson.M{"$eq": bson.A{"$program_key", params.ProgramKey}},
									bson.M{"$eq": bson.A{"$module_type", "TESTING"}},
								}},
								"else": "",
							},
						},
					},
				},
			},
		},
	}

	if params.Limit != nil {
		if *params.Limit > 0 {
			pipelines = append(pipelines, bson.M{"$limit": params.Limit})
		}
	}

	return pipelines
}

func GetRecordOnlyStagePTN() []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$ne": bson.A{"$module_type", "TESTING"}},
						bson.M{"$ne": bson.A{"$module_type", "WITH_CODE"}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
						bson.M{"$eq": bson.A{"$program_key", "utbk"}},
					},
				},
			},
		},
	}
}
