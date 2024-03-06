package aggregates

import (
	"go.mongodb.org/mongo-driver/bson"
	"smartbtw.com/services/profile/request"
)

func GetHistoryScoreByTargetType(targetType string, params *request.HistoryScoreQueryParams) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$ne": bson.A{params.SmartBTWID, nil}},
								"then": bson.M{"$eq": bson.A{"$smartbtw_id", params.SmartBTWID}},
								"else": "",
							},
						},
						bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$ne": bson.A{params.TaskID, nil}},
								"then": bson.M{"$eq": bson.A{"$task_id", params.TaskID}},
								"else": "",
							},
						},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "student_targets",
				"localField":   "target_id",
				"foreignField": "_id",
				"as":           "target",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$target",
			},
		},
		{
			"$sort": bson.M{
				"task_id":    -1,
				"created_at": -1,
			},
		},
		{
			"$addFields": bson.M{
				"target_id":          "$$REMOVE",
				"target.smartbtw_id": "$$REMOVE",
			},
		},
	}
}

func GetHistoryScoreCPNS(params *request.HistoryScoreQueryParams) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$ne": bson.A{params.SmartBTWID, nil}},
								"then": bson.M{"$eq": bson.A{"$smartbtw_id", params.SmartBTWID}},
								"else": "",
							},
						},
						bson.M{
							"$cond": bson.M{
								"if":   bson.M{"$ne": bson.A{params.TaskID, nil}},
								"then": bson.M{"$eq": bson.A{"$task_id", params.TaskID}},
								"else": "",
							},
						},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "student_target_cpns",
				"localField":   "target_id",
				"foreignField": "_id",
				"as":           "target",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$target",
			},
		},
		{
			"$sort": bson.M{
				"task_id":    -1,
				"created_at": -1,
			},
		},
		{
			"$addFields": bson.M{
				"target_id":          "$$REMOVE",
				"target.smartbtw_id": "$$REMOVE",
			},
		},
	}
}
