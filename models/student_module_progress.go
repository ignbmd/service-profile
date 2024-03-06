package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentModuleProgress struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	SmartbtwID  int                `json:"smartbtw_id" bson:"smartbtw_id"`
	TaskID      int                `json:"task_id" bson:"task_id"`
	ModuleNo    int                `json:"module_no" bson:"module_no"`
	Repeat      int                `json:"repeat" bson:"repeat"`
	ModuleTotal int                `json:"module_total" bson:"module_total"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt   *time.Time         `json:"deleted_at" bson:"deleted_at"`
}

type UpsertLiveRankData struct {
	SmartBTWID             uint    `json:"smart_btw_id" valid:"type(uint),required"`
	PackageID              uint    `json:"package_id" valid:"type(uint)"`
	ClusterID              *uint   `json:"cluster_id" valid:"type(*uint)"`
	Slug                   string  `json:"slug" valid:"type(string)"`
	Program                string  `json:"program" valid:"type(string),required"`
	PassingScorePercentage float64 `json:"passing_score_percentage" valid:"optional"`
}

type MessageUpdateLiveRanking struct {
	LRData UpsertLiveRankData `json:"live_ranking" bson:"live_ranking"`
}

type UpsertLiveRankDataStruct struct {
	Version uint                     `json:"version"`
	Data    MessageUpdateLiveRanking `json:"data"`
}
