package aggregates

import "go.mongodb.org/mongo-driver/bson"

func GetStudentAccess(smartBTWID int, appType string) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"smartbtw_id": smartBTWID,
				"app_type":    appType,
				"deleted_at":  nil,
			},
		},
	}
}
