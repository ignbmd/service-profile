package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/asaskevich/govalidator"
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func ListenStudentTargetBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "student.target.created":
		if StudentTargetCreate(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "user.data.updated":
		if UpdateUserData(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "school.updated":
		if UpdateSchoolData(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "study.program.updated":
		if UpdateStudyProgramData(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "target.score.updated":
		if UpdateTargetScore(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "student.target.updated":
		if UpdateStudentTarget(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}
	return false
}

func StudentTargetCreate(body []byte, key string) bool {
	msgDb, err := request.UnmarshalMessageStudentTargetBody(body)
	msgE, _ := request.UnmarshalStudentTargetPtkElastic(body)
	msg, _ := request.UnmarshalStudentTargetPtnElastic(body)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	result, err := govalidator.ValidateStruct(&msgDb)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msgDb.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	_, err = lib.CreateStudentTarget(&msgDb.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data to db", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if msgDb.Data.TargetType == string(models.PTK) {
		err1 := lib.InsertStudentTargetPtkElastic(&msgE.Data)
		if err1 != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data to elastic", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
	} else {
		err1 := lib.InsertStudentTargetPtnElastic(&msg.Data)
		if err1 != nil {
			text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data to elastic", key)
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
	}

	return true
}

func UpdateUserData(body []byte, key string) bool {
	msgE, err := request.UnmarshalUpdateUserDataElastic(body)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	result, err := govalidator.ValidateStruct(&msgE)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msgE.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	err1 := lib.UpdateUserData(&msgE.Data, msgE.Data.SmartbtwID, context.Background())
	if err1 != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when update data to elastic", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func UpdateSchoolData(body []byte, key string) bool {
	msgE, err := request.UnmarshalUpdateSchoolBodyMessage(body)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	result, err := govalidator.ValidateStruct(&msgE)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msgE.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	err1 := lib.UpdateSchool(&msgE.Data, msgE.Data.SchoolID, msgE.Data.TargetType, context.Background())
	if err1 != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when update data to elastic", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func UpdateStudyProgramData(body []byte, key string) bool {
	msgE, err := request.UnmarshalUpdateStudyProgramBodyMessage(body)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	result, err := govalidator.ValidateStruct(&msgE)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msgE.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	err1 := lib.UpdateStudyProgram(&msgE.Data, msgE.Data.MajorID, msgE.Data.TargetType, context.Background())
	if err1 != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when update data to elastic", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func UpdateTargetScore(body []byte, key string) bool {
	msgE, err := request.UnmarshalUpdateTargetScoreElastic(body)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	result, err := govalidator.ValidateStruct(&msgE)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Validation Error", key)
		validationResult := fmt.Sprintf("Has data been validated: %t", result)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(validationResult)
		log.Println(err)

		return false
	}

	if msgE.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	err1 := lib.UpdateTargetScore(&msgE.Data, msgE.Data.MajorID, msgE.Data.TargetType, context.Background())
	if err1 != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when update data to elastic", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func UpdateStudentTarget(body []byte, key string) bool {
	msg, err := request.UnmarshalUpdateStudentTargetElastic(body)
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
	ctx := request.MessageUpdateBulkStudentTarget{
		StudentData: msg.Data.StudentData,
	}

	err1 := lib.UpdateStudentTarget(ctx.StudentData, context.Background())
	if err1 != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when update data student target", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}
