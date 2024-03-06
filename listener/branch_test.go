package listener_test

import (
	"fmt"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/listener"
)

type MockAcknowledger struct {
}

func (m *MockAcknowledger) Ack(tag uint64, multiple bool) error {
	return nil
}

func (m *MockAcknowledger) Nack(tag uint64, multiple bool, requeue bool) error {
	return nil
}

func (m *MockAcknowledger) Reject(tag uint64, requeue bool) error {
	return nil
}

func Test_CreateBranchSuccess(t *testing.T) {
	Init()

	payload := []byte(`
	{
		"version": 1,
		"data": {
			"branch_code": "KB0003",
			"branch_name": "Cabang 3",
			"created_at": "2022-01-10T07:53:16Z",
			"updated_at": "2022-01-10T07:53:16Z"
		}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey: "branch.created",
		Body:       payload,
	}

	assert.True(t, listener.ListenBranchBinding(&msg))
}

func Test_UpdateBranchSuccess(t *testing.T) {
	Init()

	payload := []byte(`
	{
		"version": 1,
		"data": {
			"branch_code": "KB0001",
			"branch_name": "Cabang 111",
			"created_at": "2022-01-10T07:53:16Z",
			"updated_at": "2022-01-10T07:53:16Z"
		}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey: "branch.updated",
		Body:       payload,
	}
	fmt.Printf("%s", payload)
	assert.True(t, listener.ListenBranchBinding(&msg))
}

func Test_CreateBranchErrorBody(t *testing.T) {
	Init()

	payload := []byte(`
	{
		"version": 1,
		"data": {
			"branch_code": "",
			"branch_name": "",
			"created_at": "2022-01-10T07:53:16Z",
			"updated_at": "2022-01-10T07:53:16Z"
		}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey: "branch.created",
		Body:       payload,
	}

	assert.False(t, listener.ListenBranchBinding(&msg))
}

func Test_CreateBranchErrorBodyDataNull(t *testing.T) {
	Init()

	payload := []byte(`
	{
		"version": 1,
		"data": {}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey: "branch.created",
		Body:       payload,
	}

	assert.False(t, listener.ListenBranchBinding(&msg))
}

func Test_CreateBranchErrorVersion0(t *testing.T) {
	Init()

	payload := []byte(`
	{
		"version": 0,
		"data": {}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey: "branch.created",
		Body:       payload,
	}

	assert.False(t, listener.ListenBranchBinding(&msg))
}
