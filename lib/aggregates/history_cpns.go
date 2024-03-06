package aggregates

import (
	"go.mongodb.org/mongo-driver/bson"
	"smartbtw.com/services/profile/request"
)

func GetStudentAverageCPNS(SmartBtwID int) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", SmartBtwID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":           bson.A{"$smartbtw_id", SmartBtwID},
				"average_twk":   bson.M{"$avg": "$twk"},
				"average_tiu":   bson.M{"$avg": "$tiu"},
				"average_tkp":   bson.M{"$avg": "$tkp"},
				"average_total": bson.M{"$avg": "$total"},
			},
		},
		{
			"$project": bson.M{
				"_id":           0,
				"average_twk":   1,
				"average_tiu":   1,
				"average_tkp":   1,
				"average_total": 1,
			},
		},
	}
}

func GetStudentLastScoreCPNS(SmartBtwID int) []bson.M {
	return []bson.M{
		{"$sort": bson.M{"created_at": -1}},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", SmartBtwID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
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
		{"$unwind": bson.M{"path": "$history_cpns", "preserveNullAndEmptyArrays": true}},
	}
}

func GetLast10StudentScoreCPNS(SmartBtwID int) []bson.M {
	return []bson.M{
		{"$sort": bson.M{"created_at": 1}},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", SmartBtwID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
						bson.M{"$not": bson.M{"$in": bson.A{"$module_type", bson.A{"PRE_TEST", "POST_TEST"}}}},
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
		{"$unwind": bson.M{"path": "$history_cpns", "preserveNullAndEmptyArrays": true}},
	}
}

func GetStudentCPNSHistoryScores(SmartBTWID int, params *request.HistoryCPNSQueryParams) []bson.M {
	pipelines := []bson.M{
		{
			"$match": bson.M{
				"smartbtw_id": SmartBTWID,
				"deleted_at":  nil,
			},
		},
		{
			"$project": bson.M{
				"_id":               1,
				"smartbtw_id":       1,
				"exam_name":         1,
				"twk":               1,
				"tiu":               1,
				"tkp":               1,
				"twk_pass":          1,
				"tiu_pass":          1,
				"tkp_pass":          1,
				"twk_time_consumed": 1,
				"tiu_time_consumed": 1,
				"tkp_time_consumed": 1,
				"package_type":      1,
				"module_type":       1,
				"total":             1,
				"grade":             1,
				"created_at":        1,
				"updated_at":        1,
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
