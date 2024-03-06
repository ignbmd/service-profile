package aggregates

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetStudents(smartbtwIds []int, fields string) []bson.M {
	var smartBtwIdFilter bson.A
	var projectFilter = bson.M{}
	var projectFields = bson.M{}

	trimmedFields := strings.ReplaceAll(fields, " ", "")
	fieldSlice := strings.Split(trimmedFields, ",")

	if len(fieldSlice) > 0 && fieldSlice[0] != "" {
		for _, v := range fieldSlice {
			projectFields[v] = 1
		}
		projectFilter = bson.M{
			"$project": projectFields,
		}
	}

	if len(smartbtwIds) > 0 {
		smartBtwIdFilter = bson.A{
			bson.M{"$in": bson.A{"$smartbtw_id", smartbtwIds}},
			bson.M{"$eq": bson.A{"$deleted_at", nil}},
		}
	} else {
		smartBtwIdFilter = bson.A{
			bson.M{"$eq": bson.A{"$deleted_at", nil}},
		}

	}

	pipelines := []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": smartBtwIdFilter,
				},
			},
		},
		{"$sort": bson.M{"smartbtw_id": 1}},
		{
			"$lookup": bson.M{
				"from":         "parent_datas",
				"localField":   "_id",
				"foreignField": "student_id",
				"as":           "parent_datas",
			},
		},
		{"$unwind": bson.M{"path": "$parent_datas", "preserveNullAndEmptyArrays": true}},
		{
			"$lookup": bson.M{
				"from":         "branchs",
				"localField":   "branch_code",
				"foreignField": "branch_code",
				"as":           "branchs",
			},
		},
		{"$unwind": bson.M{"path": "$branchs", "preserveNullAndEmptyArrays": true}},
		projectFilter,
	}

	if len(fieldSlice) > 0 && fieldSlice[0] == "" {
		pipelines = pipelines[:len(pipelines)-1]
	}
	return pipelines
}

func GetStudentWithParents(smartBtwId int) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", smartBtwId}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "parent_datas",
				"localField":   "_id",
				"foreignField": "student_id",
				"as":           "parent_datas",
			},
		},
		{"$unwind": bson.M{"path": "$parent_datas", "preserveNullAndEmptyArrays": true}},
		{
			"$lookup": bson.M{
				"from":         "branchs",
				"localField":   "branch_code",
				"foreignField": "branch_code",
				"as":           "branchs",
			},
		},
		{"$unwind": bson.M{"path": "$branchs", "preserveNullAndEmptyArrays": true}},
	}
}

func GetStudentOnly(smartBtwId int) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$smartbtw_id", smartBtwId}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
	}
}

func GetStudentByBranchCodeNoLimit(bc string, skip int, limit int, page int, sc *string) []bson.M {
	search := ""
	if sc != nil {
		search = *sc
	}

	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$cond": bson.A{
						bc == "PT0000",
						bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$deleted_at", nil}},
							},
						},
						bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$branch_code", bc}},
								bson.M{"$eq": bson.A{"$deleted_at", nil}},
							},
						},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$cond": bson.A{
						search != "",
						bson.M{
							"$or": bson.A{
								bson.M{
									"$regexMatch": bson.M{
										"input": "$name",
										"regex": primitive.Regex{
											Pattern: search,
											Options: "i",
										},
									},
								},
								bson.M{
									"$regexMatch": bson.M{
										"input": "$email",
										"regex": primitive.Regex{
											Pattern: search,
											Options: "i",
										},
									},
								},
							},
						},
						"",
					},
				},
			},
		},
		// {
		// 	"$limit": 500,
		// },
		{
			"$sort": bson.M{"created_at": -1},
		},
		{
			"$lookup": bson.M{
				"from":         "parent_datas",
				"localField":   "_id",
				"foreignField": "student_id",
				"as":           "parent_datas",
			},
		},
		{"$unwind": bson.M{"path": "$parent_datas", "preserveNullAndEmptyArrays": true}},
		{
			"$lookup": bson.M{
				"from":         "branchs",
				"localField":   "branch_code",
				"foreignField": "branch_code",
				"as":           "branchs",
			},
		},
		{"$unwind": bson.M{"path": "$branchs", "preserveNullAndEmptyArrays": true}},
		{
			"$facet": bson.M{
				"all": bson.A{},
				"filter": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$cond": bson.A{
									search != "",
									bson.M{
										"$or": bson.A{
											bson.M{
												"$regexMatch": bson.M{
													"input": "$name",
													"regex": primitive.Regex{
														Pattern: search,
														Options: "i",
													},
												},
											},
											bson.M{
												"$regexMatch": bson.M{
													"input": "$email",
													"regex": primitive.Regex{
														Pattern: search,
														Options: "i",
													},
												},
											},
										},
									},
									"",
								},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				// "data":     bson.M{"$slice": bson.A{"$filter", skip, limit}},
				"total":    bson.M{"$size": "$all"},
				"filtered": bson.M{"$size": "$filter"},
				"page":     page,
			},
		},
		{
			"$project": bson.M{
				// "all":    false,
				"filter": false,
			},
		},
	}
}

func GetStudentByBranchCodeAndPagination(bc string, skip int, limit int, page int, sc *string) []bson.M {
	search := ""
	if sc != nil {
		search = *sc
	}

	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$cond": bson.A{
						bc == "PT0000",
						bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$deleted_at", nil}},
							},
						},
						bson.M{
							"$and": bson.A{
								bson.M{"$eq": bson.A{"$branch_code", bc}},
								bson.M{"$eq": bson.A{"$deleted_at", nil}},
							},
						},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$cond": bson.A{
						search != "",
						bson.M{
							"$or": bson.A{
								bson.M{
									"$regexMatch": bson.M{
										"input": "$name",
										"regex": primitive.Regex{
											Pattern: search,
											Options: "i",
										},
									},
								},
								bson.M{
									"$regexMatch": bson.M{
										"input": "$email",
										"regex": primitive.Regex{
											Pattern: search,
											Options: "i",
										},
									},
								},
							},
						},
						"",
					},
				},
			},
		},
		{
			"$limit": 50,
		},
		{
			"$sort": bson.M{"created_at": -1},
		},
		{
			"$lookup": bson.M{
				"from":         "parent_datas",
				"localField":   "_id",
				"foreignField": "student_id",
				"as":           "parent_datas",
			},
		},
		{"$unwind": bson.M{"path": "$parent_datas", "preserveNullAndEmptyArrays": true}},
		{
			"$lookup": bson.M{
				"from":         "branchs",
				"localField":   "branch_code",
				"foreignField": "branch_code",
				"as":           "branchs",
			},
		},
		{"$unwind": bson.M{"path": "$branchs", "preserveNullAndEmptyArrays": true}},
		{
			"$facet": bson.M{
				"all": bson.A{},
				"filter": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$cond": bson.A{
									search != "",
									bson.M{
										"$or": bson.A{
											bson.M{
												"$regexMatch": bson.M{
													"input": "$name",
													"regex": primitive.Regex{
														Pattern: search,
														Options: "i",
													},
												},
											},
											bson.M{
												"$regexMatch": bson.M{
													"input": "$email",
													"regex": primitive.Regex{
														Pattern: search,
														Options: "i",
													},
												},
											},
										},
									},
									"",
								},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"data":     bson.M{"$slice": bson.A{"$filter", skip, limit}},
				"total":    bson.M{"$size": "$all"},
				"filtered": bson.M{"$size": "$filter"},
				"page":     page,
			},
		},
		{
			"$project": bson.M{
				"all":    false,
				"filter": false,
			},
		},
	}
}

func GetStudentByBranchCodeArrayAndPagination(bc []string, skip int, limit int, page int, sc *string) []bson.M {
	search := ""
	if sc != nil {
		search = *sc
	}

	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$in": bson.A{"$branch_code", bc}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$cond": bson.A{
						search != "",
						bson.M{
							"$or": bson.A{
								bson.M{
									"$regexMatch": bson.M{
										"input": "$name",
										"regex": primitive.Regex{
											Pattern: search,
											Options: "i",
										},
									},
								},
								bson.M{
									"$regexMatch": bson.M{
										"input": "$email",
										"regex": primitive.Regex{
											Pattern: search,
											Options: "i",
										},
									},
								},
							},
						},
						"",
					},
				},
			},
		},
		{
			"$limit": 50,
		},
		{
			"$sort": bson.M{"created_at": -1},
		},
		{
			"$lookup": bson.M{
				"from":         "parent_datas",
				"localField":   "_id",
				"foreignField": "student_id",
				"as":           "parent_datas",
			},
		},
		{"$unwind": bson.M{"path": "$parent_datas", "preserveNullAndEmptyArrays": true}},
		{
			"$lookup": bson.M{
				"from":         "branchs",
				"localField":   "branch_code",
				"foreignField": "branch_code",
				"as":           "branchs",
			},
		},
		{"$unwind": bson.M{"path": "$branchs", "preserveNullAndEmptyArrays": true}},
		{
			"$facet": bson.M{
				"all": bson.A{},
				"filter": bson.A{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$cond": bson.A{
									search != "",
									bson.M{
										"$or": bson.A{
											bson.M{
												"$regexMatch": bson.M{
													"input": "$name",
													"regex": primitive.Regex{
														Pattern: search,
														Options: "i",
													},
												},
											},
											bson.M{
												"$regexMatch": bson.M{
													"input": "$email",
													"regex": primitive.Regex{
														Pattern: search,
														Options: "i",
													},
												},
											},
										},
									},
									"",
								},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"data":     bson.M{"$slice": bson.A{"$filter", skip, limit}},
				"total":    bson.M{"$size": "$all"},
				"filtered": bson.M{"$size": "$filter"},
				"page":     page,
			},
		},
		{
			"$project": bson.M{
				"all":    false,
				"filter": false,
			},
		},
	}
}

func GetStudentCompletedModules(smartBTWID int) []bson.M {
	return []bson.M{
		{"$sort": bson.M{"created_at": -1}},
		{
			"$match": bson.M{
				"smartbtw_id": smartBTWID,
				"deleted_at":  nil,
			},
		},
		{
			"$group": bson.M{
				"_id":         bson.A{"$_id"},
				"module_code": bson.M{"$first": "$module_code"},
			},
		},
		{
			"$project": bson.M{
				"_id":         0,
				"module_code": 1,
			},
		},
	}
}

func GetStudentBranchByEmail(emails []string) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$in": bson.A{"$email", emails}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "branchs",
				"localField":   "branch_code",
				"foreignField": "branch_code",
				"as":           "branch",
			},
		},
		{"$sort": bson.M{"updated_at": -1}},
		{"$unwind": bson.M{"path": "$branchs", "preserveNullAndEmptyArrays": true}},
		{
			"$project": bson.M{
				"_id":        false,
				"email":      true,
				"updated_at": true,
				"created_at": true,
				"branch":     true,
			},
		},
		{"$unwind": bson.M{"path": "$branch", "preserveNullAndEmptyArrays": false}},
		{"$addFields": bson.D{
			{"branch_code", "$branch.branch_code"},
			{"branch_name", "$branch.branch_name"},
			{"branch", "$$REMOVE"},
		},
		},
	}
}
