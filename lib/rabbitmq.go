package lib

import (
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
)

func MessageAck(msg *amqp.Delivery) {
	if err := msg.Ack(false); err != nil {
		golog.Slack.Error("Failed to ack message ", err)
	}
}
