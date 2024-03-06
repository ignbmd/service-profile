package aggregates

import (
	"go.mongodb.org/mongo-driver/bson"
)

func GetStudentWalletTotalBalance(smartBTWID int) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"smartbtw_id": smartBTWID,
				"deleted_at":  nil,
			},
		},
		{
			"$group": bson.M{
				"_id":     nil,
				"balance": bson.M{"$sum": "$point"},
			},
		},
		{
			"$project": bson.M{
				"_id": 0,
			},
		},
	}
}

func GetStudentWalletBalance(smartBTWID int, isMultipliedPoint bool) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"smartbtw_id": smartBTWID,
				"deleted_at":  nil,
			},
		},
		{
			"$project": bson.M{
				"_id":         1,
				"smartbtw_id": 1,
				"type":        1,
				"created_at":  1,
				"updated_at":  1,
				"point": bson.M{
					"$cond": bson.M{
						"if": isMultipliedPoint,
						"then": bson.M{
							"$divide": bson.A{
								bson.M{
									"$sum": "$point",
								},
								100,
							}},
						"else": "$point",
					},
				},
				"balance": bson.M{
					"$cond": bson.M{
						"if": isMultipliedPoint,
						"then": bson.M{
							"$divide": bson.A{
								bson.M{
									"$sum": "$point",
								},
								100,
							}},
						"else": "$point",
					},
				},
			},
		},
	}
}

func GetStudentWalletHistoryByWalletType(smartBTWID int, walletType *string) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"smartbtw_id": smartBTWID,
				"deleted_at":  nil,
				"$expr": bson.M{
					"$cond": bson.M{
						"if":   walletType != nil && *walletType != "",
						"then": bson.M{"$eq": bson.A{"$type", walletType}},
						"else": "",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "wallet_histories",
				"localField":   "_id",
				"foreignField": "wallet_id",
				"as":           "wallet_history",
			},
		},
		{
			"$addFields": bson.M{
				"wallet_history": bson.M{
					"$map": bson.M{
						"input": "$wallet_history",
						"as":    "i",
						"in": bson.M{
							"$mergeObjects": bson.A{
								"$$i",
								bson.M{
									"type": "$type",
								},
								bson.M{
									"transaction_type": bson.M{
										"$ifNull": bson.A{
											"$$i.type", "GENERAL",
										},
									},
								},
								bson.M{
									"amount_pay": bson.M{
										"$ifNull": bson.A{
											"$$i.amount_pay", 0,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{"$unwind": bson.M{"path": "$wallet_history"}},
		{"$replaceRoot": bson.M{"newRoot": "$wallet_history"}},
		{
			"$sort": bson.M{"created_at": -1},
		},
		{
			"$addFields": bson.M{
				"point": bson.M{
					"$divide": bson.A{
						"$point",
						100,
					},
				},
			},
		},
		{
			"$project": bson.M{
				"_id":         0,
				"wallet_id":   0,
				"smartbtw_id": 0,
				"deleted_at":  0,
			},
		},
	}
}
