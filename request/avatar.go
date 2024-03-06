package request

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateAvatar struct {
	SmartbtwID int    `json:"smartbtw_id"`
	AvaType    string `json:"ava_type"`
	Style      int    `json:"style"`
}

type UpdateAvatar struct {
	ID         primitive.ObjectID `json:"_id"`
	SmartbtwID int                `json:"smartbtw_id"`
	AvaType    string             `json:"ava_type"`
	Style      int                `json:"style"`
}

type BodyRequestAvatar struct {
	SmartbtwID int    `json:"smartbtw_id"`
	AvaType    string `json:"ava_type"`
}
