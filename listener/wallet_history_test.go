package listener_test

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/listener"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func TestCreateWalletHistoryPremiumPackage(t *testing.T) {
	Init()
	payload1 := request.CreateWallet{
		SmartbtwID: 999,
		Point:      100,
		Type:       string(models.BONUS),
	}

	_, err := lib.CreateWallet(&payload1)
	assert.Nil(t, err)

	payload := []byte(`{
		"version": 1,
		"data": {
			"smartbtw_id": 999,
			"price": 2310,
			"description": "beli paket premium"
		}
		}`)

	msg := amqp.Delivery{
		RoutingKey:   "wallet-history-premium.created",
		Body:         payload,
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenWalletHistoryBinding(&msg))
}

// func TestCreateWalletHistoryPremiumPackageWithEmptyBodyData(t *testing.T) {
// 	Init()

// 	json := `{
// 		"version": 2,
// 		"data": {}
// 	}`

// 	msg := amqp.Delivery{
// 		RoutingKey: "wallet-history-premium.created",
// 		Body:       []byte(json),
// 	}

// 	assert.False(t, listener.ListenWalletHistoryBinding(&msg))
// }

func TestCreateWalletHistoryInvitePeople(t *testing.T) {
	Init()
	payload1 := request.CreateWallet{
		SmartbtwID: 999111,
		Point:      100,
		Type:       string(models.BONUS),
	}

	_, err := lib.CreateWallet(&payload1)
	assert.Nil(t, err)

	payload := []byte(`{
		"version": 1,
		"data": {
			"smartbtw_id": 999111,
			"description": "mengundang teman"
		}
		}`)

	msg := amqp.Delivery{
		RoutingKey:   "wallet-history-invite.created",
		Body:         payload,
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenWalletHistoryBinding(&msg))
}
