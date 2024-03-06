package request

import (
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateWalletHistoryPremium struct {
	WalletID    primitive.ObjectID `json:"wallet_id"`
	SmartbtwID  int                `json:"smartbtw_id"`
	Point       float32            `json:"point"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	Price       float64            `json:"price"`
}

type CreateWalletHistory struct {
	WalletID    primitive.ObjectID `json:"wallet_id"`
	SmartbtwID  int                `json:"smartbtw_id"`
	Point       float32            `json:"point"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
}

type GetWalletHistory struct {
	Type *string `query:"type"`
}

type MessageWalletHistoryBody struct {
	Version int                 `json:"version"`
	Data    CreateWalletHistory `json:"data" valid:"required"`
}
type MessageWalletHistoryPremiumPackageBody struct {
	Version int                        `json:"version"`
	Data    CreateWalletHistoryPremium `json:"data" valid:"required"`
}

func UnmarshalMessageWalletHistoryBody(data []byte) (MessageWalletHistoryBody, error) {
	var decoded MessageWalletHistoryBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalMessageWalletHistoryPremiumPackageBody(data []byte) (MessageWalletHistoryPremiumPackageBody, error) {
	var decoded MessageWalletHistoryPremiumPackageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
