package listener_test

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/listener"
)

func TestCreateWallet(t *testing.T) {
	Init()

	payload := []byte(`{
		"version": 1,
		"data": {
			"smartbtw_id": 1234,
			"point": 0,
			"type": "DEFAULT"
			}
	}`)

	msg := amqp.Delivery{
		RoutingKey:   "wallet.created",
		Body:         payload,
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenWalletBinding(&msg))
}

func TestReceiveWallet(t *testing.T) {
	Init()

	payload := []byte(`{
		"version": 1,
		"data": {
			"smartbtw_id": 1234,
			"point": 5000,
			"type": "DEFAULT",
			"description": "Pembelian in App"
			}
	}`)

	msg := amqp.Delivery{
		RoutingKey:   "wallet.received",
		Body:         payload,
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenWalletBinding(&msg))
}

// func TestCuttingWallet(t *testing.T) {
// 	Init()

// 	payload := []byte(`{
// 		"version": 1,
// 		"data": 142557
// 	}`)

// 	msg := amqp.Delivery{
// 		RoutingKey:   "wallet.cutting-masa-ai",
// 		Body:         payload,
// 		Acknowledger: &MockAcknowledger{},
// 	}

// 	assert.True(t, listener.ListenWalletBinding(&msg))
// }
