package listener

import (
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/mockstruct"
)

type ProgressResultRaportBody struct {
	SmartbtwID int    `json:"smartbtw_id"`
	Program    string `json:"program"`
	UKAType    string `json:"uka_type"`
	StageType  string `json:"stage_type"`
	Link       string `json:"raport_link"`
}
type CreateProgressResultRaportMessageBody struct {
	Version int                      `json:"version"`
	Data    ProgressResultRaportBody `json:"data" valid:"required"`
}

type BuildProgressRaportMessageBody struct {
	Version int                                      `json:"version"`
	Data    mockstruct.GenerateProgressReportMessage `json:"data" valid:"required"`
}

func UnmarshalCreateProgressResultRaportBody(data []byte) (CreateProgressResultRaportMessageBody, error) {
	var decoded CreateProgressResultRaportMessageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalBuildProgressResultRaportBody(data []byte) (BuildProgressRaportMessageBody, error) {
	var decoded BuildProgressRaportMessageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func ListenProgressResultRaportBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "progress-result-raport.generated":
		if CreateProgressResultRaport(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "progress-result-raport.build.queue":
		if BuildProgressRaport(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}

	return false
}

func CreateProgressResultRaport(body []byte, key string) bool {
	msg, err := UnmarshalCreateProgressResultRaportBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Version >= 1 {
		// return true
		// bd := models.ProgressResultRaport{
		// 	SmartbtwID: msg.Data.SmartbtwID,
		// 	Program:    msg.Data.Program,
		// 	Link:       msg.Data.Link,
		// 	UKAType:    msg.Data.UKAType,
		// 	StageType:  msg.Data.StageType,
		// }
		err = lib.CacheProgressReportToRedis(uint(msg.Data.SmartbtwID), msg.Data.UKAType, msg.Data.StageType, msg.Data.Program, msg.Data.Link)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

	}
	return true
}

func BuildProgressRaport(body []byte, key string) bool {
	msg, err := UnmarshalBuildProgressResultRaportBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Version >= 1 {
		err = lib.StoreBuildProcess(msg.Data)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when build progress raport data", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

	}
	return true
}
