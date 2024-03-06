package listener_test

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/listener"
	"smartbtw.com/services/profile/request"
)

func TestCreateClassMemberSuccess(t *testing.T) {
	Init()
	pyl1 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res, err := lib.CreateClassroom(&pyl1)
	assert.Nil(t, err)

	payload := []byte(`
	{
		"version": 1,
		"data": {
			"smartbtw_id": 1,
			"classroom_id": "` + res.InsertedID.(primitive.ObjectID).Hex() + `"
		}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey:   "class-member.created",
		Body:         payload,
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenClassMemberBinding(&msg))
}

func TestUpdateClassMemberSuccess(t *testing.T) {
	Init()

	pyl1 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res, err := lib.CreateClassroom(&pyl1)
	assert.Nil(t, err)

	pyl2 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res2, err := lib.CreateClassroom(&pyl2)
	assert.Nil(t, err)

	payload := []byte(`
	{
		"version": 1,
		"data": {
			"smartbtw_id": 1,
			"classroom_id_before": "` + res.InsertedID.(primitive.ObjectID).Hex() + `",
			"classroom_id_after": "` + res2.InsertedID.(primitive.ObjectID).Hex() + `"
		}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey:   "class-member.updated",
		Body:         payload,
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenClassMemberBinding(&msg))
}

func TestSwitchClassMemberSuccess(t *testing.T) {
	Init()

	pyl1 := request.CreateClassroom{
		ID:          primitive.NewObjectID(),
		BranchCode:  "B001",
		Quota:       10,
		QuotaFilled: 0,
		Description: "Classroom for testing",
		Tags:        []string{"test", "testing", "testcase"},
		Year:        2021,
		Status:      "active",
		Title:       "Test Classroom",
		ProductID:   "prod_1234567890",
	}
	res, err := lib.CreateClassroom(&pyl1)
	assert.Nil(t, err)

	pyl := request.CreateClassMember{
		SmartbtwID:  87654,
		ClassroomID: res.InsertedID.(primitive.ObjectID),
	}
	err = lib.CreateClassMember(&pyl)
	assert.Nil(t, err)

	payload := []byte(`
	{
		"version": 1,
		"data": {
			"classroom_id": "` + res.InsertedID.(primitive.ObjectID).Hex() + `",
			"class_members": [
				{
					"smartbtw_id": 87654,
					"btwedutech_id": 12345
				}
			]
		}
	}`)
	msg := amqp.Delivery{
		RoutingKey:   "class-member.switch",
		Body:         payload,
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenClassMemberBinding(&msg))

}
