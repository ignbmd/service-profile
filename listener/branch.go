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

func ListenBranchBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "branch.created", "branch.updated":
		if BranchUpsert(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}

	return false
}

func BranchUpsert(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageBranchBody(body)

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

	err = lib.UpsertBranchData(&msg.Data)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] Upsert branch data for branch %s success", key, msg.Data.BranchCode)
	log.Println(text)

	return true
}
