package request

import "github.com/bytedance/sonic"

type CreateWallet struct {
	SmartbtwID int     `json:"smartbtw_id"`
	Point      float32 `json:"point"`
	Type       string  `json:"type"`
}

type ReceiveWallet struct {
	SmartbtwID      int     `json:"smartbtw_id"`
	AmountPay       float32 `json:"amount_pay"`
	Point           float32 `json:"point"`
	Description     string  `json:"description"`
	Type            string  `json:"type"`
	TransactionType string  `json:"transaction_type"`
}

type ChargeWallet struct {
	SmartbtwID      int     `json:"smartbtw_id"`
	Point           float32 `json:"point"`
	Description     string  `json:"description"`
	TransactionType string  `json:"transaction_type"`
}

type CuttingWallet struct {
	SmartbtwID int `json:"smartbtw_id"`
}

type MessageWalletBody struct {
	Version int          `json:"version"`
	Data    CreateWallet `json:"data" valid:"required"`
}

type MessageReceiveWalletBody struct {
	Version int           `json:"version"`
	Data    ReceiveWallet `json:"data" valid:"required"`
}
type MessageCuttingWalletBody struct {
	Version int           `json:"version"`
	Data    CuttingWallet `json:"data" valid:"required"`
}

func UnmarshalMessageWalletBody(data []byte) (MessageWalletBody, error) {
	var decoded MessageWalletBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
func UnmarshalMessageReceiveWalletBody(data []byte) (MessageReceiveWalletBody, error) {
	var decoded MessageReceiveWalletBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalMessageCuttingWalletBody(data []byte) (MessageCuttingWalletBody, error) {
	var decoded MessageCuttingWalletBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
