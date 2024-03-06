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

func ListenClassMemberBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "class-member.created":
		if CreateClassMember(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "class-member.updated":
		if UpdateClassMember(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "class-member.deleted":
		if DeleteClassMember(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "class-member.switch":
		if SwitchClassMember(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}

	return false
}

func CreateClassMember(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageCreateClassMemberBody(body)
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

	err = lib.CreateClassMember(&msg.Data)

	if err != nil {
		switch err.Error() {
		case "siswa ini terdaftar pada kelas":
			return true
		case "mongo: no documents in result":
			text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return true
		default:
			text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
	}

	text := fmt.Sprintf("[RabbitMQ][%s] insert data classmember for student id %d and class id %s success", key, msg.Data.SmartbtwID, msg.Data.ClassroomID.Hex())
	log.Println(text)

	return true
}

func UpdateClassMember(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageUpdateClassMemberBody(body)
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

	err = lib.UpdateClassMember(&msg.Data)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)

		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] update data classmember for student id %d and class id %s success", key, msg.Data.SmartbtwID, msg.Data.ClassroomIDAfter.Hex())
	log.Println(text)

	return true
}

func DeleteClassMember(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageCreateClassMemberBody(body)
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

	err = lib.SoftDeleteClassMemberBySmartbtwIDAndClassroomID(&msg.Data)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)

		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] soft delete data classmember for student id %d and class id %s success", key, msg.Data.SmartbtwID, msg.Data.ClassroomID.Hex())
	log.Println(text)

	return true
}

func SwitchClassMember(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageSwitchClassMemberBody(body)
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

	err = lib.SwitchClassMember(&msg.Data)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)

		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] switch data classmember for class id %s success", key, msg.Data.ClassroomID.Hex())
	log.Println(text)

	return true
}
