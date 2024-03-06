package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScoreSkdBkn struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StudentID primitive.ObjectID `json:"student_id" bson:"student_id"`
	Year      int                `json:"year" bson:"year"`
	ScoreTWK  float32            `json:"score_twk" bson:"score_twk"`
	ScoreTIU  float32            `json:"score_tiu" bson:"score_tiu"`
	ScoreTKP  float32            `json:"score_tkp" bson:"score_tkp"`
	ScoreSKD  float32            `json:"score_skd" bson:"score_skd"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at" bson:"deleted_at"`
}
