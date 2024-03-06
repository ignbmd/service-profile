package listener

import (
	"fmt"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func ListenClassroomBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "classroom.created":
		if CreateClassroom(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "classroom.updated":
		if UpdateClassroom(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}
	return false
}

func CreateClassroom(body []byte, key string) bool {
	msg, err := request.UnmarshalBodyCreateClassroom(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)

		return false
	}

	result, err := govalidator.ValidateStruct(&msg)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msg.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	_, err = lib.CreateClassroom(&msg.Data)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)

		return false
	}

	return true
}

func UpdateClassroom(body []byte, key string) bool {
	msg, err := request.UnmarshalBodyUpdateClassroom(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)

		return false
	}

	result, err := govalidator.ValidateStruct(&msg)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msg.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	err = lib.UpdateClassroom(&msg.Data)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)

		return false
	}

	return true
}
