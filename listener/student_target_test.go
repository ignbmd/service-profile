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

func Test_CreateStudentTarget(t *testing.T) {
	Init()

	json := `
	{
		"version": 2,
		"data": {
			"smartbtw_id": 99000,
     		"name": "John Doe",
			"school_id": 1,
			"major_id":1,
			"school_name": "UI",
			"major_name": "Kedokteran",
			"target_score": 700,
			"target_type": "PTN"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "student.target.created",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenStudentTargetBinding(&msg))
}

func TestCreateStudentTargetInvalidBody(t *testing.T) {
	Init()

	json := `
	{
		"version": 2,
		"data": {
			"smartbtw_id": 69953,
     		"name": "asdasd",
			"school_id": 1,
			"major_id":1,
			"school_name": "UI",
			"major_name": "Kedokteran",
			"target_score": "700",
			"target_type": "PTK"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "student.target.created",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.False(t, listener.ListenStudentTargetBinding(&msg))
}

func TestCreateStudentTargetEmptyBody(t *testing.T) {
	Init()

	json := `
	{
		"version": 2,
		"data": {}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "student.target.created",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.False(t, listener.ListenStudentTargetBinding(&msg))
}

func TestUpdateUserData(t *testing.T) {
	Init()
	json1 := `
	{
		"version": 2,
		"data": {
			"smartbtw_id": 239118197,
	 		"name": "iiii",
			"school_id": 1,
			"major_id":1,
			"school_name": "UI",
			"major_name": "Kedokteran",
			"target_score": 700,
			"target_type": "PTK"
		}
	}
	`
	msg1 := amqp.Delivery{
		RoutingKey:   "student.target.created",
		Body:         []byte(json1),
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenStudentTargetBinding(&msg1))

	json := `
	{
		"version": 2,
		"data": {
			"smartbtw_id": 239118197,
     		"name": "asdasd",
			"photo": "https://test.png"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "user.data.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenStudentTargetBinding(&msg))
}

func TestUpdateSchool(t *testing.T) {
	Init()

	json := `
	{
		"version": 2,
		"data": {
			"school_id": 1,
     		"school_name": "STIN",
			"target_type": "PTK"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "school.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenStudentTargetBinding(&msg))
}

func TestUpdateStudyProgram(t *testing.T) {
	Init()

	json := `
	{
		"version": 2,
		"data": {
			"major_id": 1,
			"major_name": "Nautika",
			"target_type": "PTK"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "study.program.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenStudentTargetBinding(&msg))
}

func TestUpdateTargetScore(t *testing.T) {
	Init()

	json := `
	{
		"version": 2,
		"data": {
			"major_id":1,
			"target_score": 900,
			"target_type": "PTK"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "target.score.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenStudentTargetBinding(&msg))
}

func TestUpdateStudentTarget(t *testing.T) {
	Init()
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  338899,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err := lib.CreateStudentTarget(&payload1)
	assert.Nil(t, err)

	json := `
	{
		"version": 2,
		"data": {
			"student_data" :
			[{
				"smartbtw_id": 338899,
				"school_id": 1,
				"major_id":1,
				"school_name": "UNES",
				"major_name": "TEKNIK MESIN",
				"target_score": 700,
				"target_type": "PTN",
				"position": 1,
				"type": "PRIMARY"
			},
			{
				"smartbtw_id": 338899,
				"school_id": 1,
				"major_id":1,
				"school_name": "UI",
				"major_name": "KEDOKTERAN",
				"target_score": 700,
				"target_type": "PTN",
				"position": 2,
				"type": "SECONDARY"
			}]
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "student.target.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenStudentTargetBinding(&msg))
}
