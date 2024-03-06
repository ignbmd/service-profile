package listener

import (
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

// import (
// 	"fmt"
// 	"log"

// 	"github.com/bytedance/sonic"
// 	"github.com/pandeptwidyaop/golog"
// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"smartbtw.com/services/profile/lib"
// 	"smartbtw.com/services/profile/mockstruct"
// 	"smartbtw.com/services/profile/models"
// 	"smartbtw.com/services/profile/request"
// )

type ResultRaportBody struct {
	SmartbtwID int    `json:"smartbtw_id"`
	Program    string `json:"program"`
	TaskID     int    `json:"task_id"`
	Link       string `json:"link"`
}
type CreateResultRaportMessageBody struct {
	Version int              `json:"version"`
	Data    ResultRaportBody `json:"data" valid:"required"`
}

type BuildRaportPTKMessageBody struct {
	Version int                      `json:"version"`
	Data    request.CreateHistoryPtk `json:"data" valid:"required"`
}

type BuildRaportPTNMessageBody struct {
	Version int                      `json:"version"`
	Data    request.CreateHistoryPtn `json:"data" valid:"required"`
}

type BuildRaportCPNSMessageBody struct {
	Version int                       `json:"version"`
	Data    request.CreateHistoryCpns `json:"data" valid:"required"`
}

type RequestBuildRaportBulkMessageBody struct {
	Version int                               `json:"version"`
	Data    mockstruct.BodyRequestBuildRaport `json:"data" valid:"required"`
}

func UnmarshalCreateResultRaportBody(data []byte) (CreateResultRaportMessageBody, error) {
	var decoded CreateResultRaportMessageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalBuildRaportPTK(data []byte) (BuildRaportPTKMessageBody, error) {
	var decoded BuildRaportPTKMessageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalBuildRaportPTN(data []byte) (BuildRaportPTNMessageBody, error) {
	var decoded BuildRaportPTNMessageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
func UnmarshalBuildRaportCPNS(data []byte) (BuildRaportCPNSMessageBody, error) {
	var decoded BuildRaportCPNSMessageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalRequestBuildRaport(data []byte) (RequestBuildRaportBulkMessageBody, error) {
	var decoded RequestBuildRaportBulkMessageBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func ListenResultRaportBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "result-raport.generated":
		if CreateResultRaport(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "raport-ptk.build":
		if BuildRaportPTK(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "raport-ptn.build":
		if BuildRaportPTN(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "raport-cpns.build":
		if BuildRaportCPNS(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "result-raport.build-bulk.request":
		if RequestBuildRaport(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}

	return false
}

func CreateResultRaport(body []byte, key string) bool {
	msg, err := UnmarshalCreateResultRaportBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Version >= 1 {
		// return true
		bd := models.ResultRaport{
			SmartbtwID: msg.Data.SmartbtwID,
			Program:    msg.Data.Program,
			TaskID:     msg.Data.TaskID,
			Link:       msg.Data.Link,
		}
		err := lib.CreateFinalRaport(bd)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

	}
	return true
}

func BuildRaportPTK(body []byte, key string) bool {

	msg, err := UnmarshalBuildRaportPTK(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Version >= 1 {
		if msg.Data.ModuleType == "WITH_CODE" {
			return true
		}
		err := lib.StoreBuildProcessRaport(models.GetResultRaportBody{
			SmartbtwID:  msg.Data.SmartBtwID,
			Program:     "PTK",
			TaskID:      msg.Data.TaskID,
			Link:        "",
			StudentName: msg.Data.StudentName,
			ExamName:    msg.Data.ExamName,
			ModuleCode:  msg.Data.ModuleCode,
			StageType:   "",
			ModuleType:  msg.Data.ModuleType,
			PackageType: msg.Data.PackageType,
		})
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error storing data build raport to queue", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

	}
	return true
}

func BuildRaportPTN(body []byte, key string) bool {

	msg, err := UnmarshalBuildRaportPTN(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Version >= 1 {
		if msg.Data.ModuleType == "WITH_CODE" {
			return true
		}
		err := lib.StoreBuildProcessRaport(models.GetResultRaportBody{
			SmartbtwID:  msg.Data.SmartBtwID,
			Program:     "PTN",
			TaskID:      msg.Data.TaskID,
			Link:        "",
			StudentName: msg.Data.StudentName,
			ExamName:    msg.Data.ExamName,
			ModuleCode:  msg.Data.ModuleCode,
			StageType:   "",
			ModuleType:  msg.Data.ModuleType,
			PackageType: msg.Data.PackageType,
		})
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error storing data build raport to queue", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

	}
	return true
}

func BuildRaportCPNS(body []byte, key string) bool {

	msg, err := UnmarshalBuildRaportCPNS(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Version >= 1 {
		if msg.Data.ModuleType == "WITH_CODE" {
			return true
		}
		err := lib.StoreBuildProcessRaport(models.GetResultRaportBody{
			SmartbtwID:  msg.Data.SmartBtwID,
			Program:     "CPNS",
			TaskID:      msg.Data.TaskID,
			Link:        "",
			StudentName: msg.Data.StudentName,
			ExamName:    msg.Data.ExamName,
			ModuleCode:  msg.Data.ModuleCode,
			StageType:   "",
			ModuleType:  msg.Data.ModuleType,
			PackageType: msg.Data.PackageType,
		})
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error storing data build raport to queue", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

	}
	return true
}

func RequestBuildRaport(body []byte, key string) bool {
	msg, err := UnmarshalRequestBuildRaport(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msg.Version >= 1 {
		// return true
		err := lib.TriggerBuildRaport(msg.Data.SmartbtwID, msg.Data.Program)
		if err != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}

	}
	return true
}
