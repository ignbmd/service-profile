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

func ListenWalletBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "wallet.created":
		if WalletCreate(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "wallet.received":
		if HandleWalletReceived(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "wallet.cutting-masa-ai":
		if HandleWalletCutting(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}

	return false
}

func WalletCreate(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageWalletBody(body)

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
	_, err = lib.CreateWallet(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	return true
}

func HandleWalletReceived(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageReceiveWalletBody(body)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	_, err = govalidator.ValidateStruct(&msg)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	if msg.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	_, err = lib.ReceiveNewWalletPoint(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func HandleWalletCutting(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageCuttingWalletBody(body)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	_, err = govalidator.ValidateStruct(&msg)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	if msg.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	err = lib.CoinCuttingMasaAI(msg.Data.SmartbtwID)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}
