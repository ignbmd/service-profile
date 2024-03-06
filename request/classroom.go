package request

import (
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateClassroom struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	BranchCode  string             `json:"branch_code"`
	Year        int32              `json:"year"`
	Quota       int32              `json:"quota"`
	QuotaFilled int32              `json:"quota_filled"`
	Tags        []string           `json:"tags"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	IsOnline    bool               `json:"is_online"`
	ProductID   string             `json:"product_id"`
}

type UpdateClassroom struct {
	ID          primitive.ObjectID `json:"id"`
	Title       string             `json:"title"`
	BranchCode  string             `json:"branch_code"`
	Year        int32              `json:"year"`
	Quota       int32              `json:"quota"`
	QuotaFilled int32              `json:"quota_filled"`
	Tags        []string           `json:"tags"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	IsOnline    bool               `json:"is_online"`
	ProductID   string             `json:"product_id"`
}

type MessageBodyCreateClassroom struct {
	Version int             `json:"version" valid:"type(int), required"`
	Data    CreateClassroom `json:"data" valid:"required"`
}

type MessageBodyUpdateClassroom struct {
	Version int             `json:"version" valid:"type(int), required"`
	Data    UpdateClassroom `json:"data" valid:"required"`
}

func UnmarshalBodyCreateClassroom(data []byte) (MessageBodyCreateClassroom, error) {
	var decoded MessageBodyCreateClassroom
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalBodyUpdateClassroom(data []byte) (MessageBodyUpdateClassroom, error) {
	var decoded MessageBodyUpdateClassroom
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
