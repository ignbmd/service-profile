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

func ListenWalletHistoryBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "wallet-history-invite.created":
		// return WalletHistoryCreateInvitePeople(msg.Body, msg.RoutingKey)
		err := msg.Ack(false)
		if err != nil {
			golog.Slack.ErrorWithData(fmt.Sprintf("[RabbitMQ][%s] Failed to ack message", msg.RoutingKey), msg.Body, err)
			return false
		}
		return true
	case "wallet-history-premium.created":
		err := msg.Ack(false)
		if err != nil {
			golog.Slack.ErrorWithData(fmt.Sprintf("[RabbitMQ][%s] Failed to ack message", msg.RoutingKey), msg.Body, err)
			return false
		}
		return true
		// return WalletHistoryCreatePremiumPackage(msg.Body, msg.RoutingKey)
	}
	return false
}

func WalletHistoryCreateInvitePeople(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageWalletHistoryBody(body)

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

	_, err = lib.CreateWalletHistoryInvitePeople(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func WalletHistoryCreatePremiumPackage(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageWalletHistoryPremiumPackageBody(body)

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

	_, err = lib.CreateWalletHistoryPremiumPackage(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}
