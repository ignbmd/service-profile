package request

import (
	"time"

	"github.com/bytedance/sonic"
)

type UpsertBranchData struct {
	BranchCode string    `json:"branch_code" valid:"type(string),required"`
	BranchName string    `json:"branch_name" valid:"type(string),required"`
	CreatedAt  time.Time `json:"created_at" valid:"type(time.Time), required"`
	UpdatedAt  time.Time `json:"updated_at" valid:"type(time.Time), required"`
}

type MessageBranchBody struct {
	Version int              `json:"version" valid:"type(int), required"`
	Data    UpsertBranchData `json:"data" valid:"required"`
}

func UnmarshalMessageBranchBody(data []byte) (MessageBranchBody, error) {
	var decoded MessageBranchBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
