package aggregates

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSingleScore(ID primitive.ObjectID) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$_id", ID}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "students",
				"localField":   "student_id",
				"foreignField": "_id",
				"as":           "student_datas",
			},
		},
		{"$unwind": bson.M{"path": "$student_datas", "preserveNullAndEmptyArrays": true}},
	}
}

func GetByManyStudent(ids []primitive.ObjectID, year int) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$in": bson.A{"$student_id", ids}},
						bson.M{"$eq": bson.A{"$year", year}},
						bson.M{"$eq": bson.A{"$deleted_at", nil}},
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "students",
				"localField":   "student_id",
				"foreignField": "_id",
				"as":           "student_datas",
			},
		},
		{"$unwind": bson.M{"path": "$student_datas", "preserveNullAndEmptyArrays": true}},
	}
}
