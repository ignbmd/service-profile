package listener

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	"golang.org/x/sync/errgroup"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func InsertScoreAssessment(body []byte, key string) bool {
	lib.LogEvent(
		"InsertScoreAssessment",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data",
		"INFO",
		fmt.Sprintf("profile-%s", key))
	// TODO: Nyesuain Multiple Target
	msg, err := request.UnmarshalMessageHistoryAssessmentsBody(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return true
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
		text := fmt.Sprintf("[RabbitMQ][%s] Version data < 1", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	studentData, err := lib.GetStudentProfileElastic(msg.Data.SmartBtwID)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when get profile elastic data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		// return false
	}

	msg.Data.StudentEmail = studentData.Email
	msg.Data.StudentName = studentData.Name

	err = lib.SaveResultAssessments(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data to firebase", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	err = lib.UpsertHistoryAssessments(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	text := fmt.Sprintf("[RabbitMQ][%s] insert history Assessments data for user_id %d and package_id %d success", key, msg.Data.SmartBtwID, msg.Data.PackageID)
	log.Println(text)
	lib.LogEvent(
		"InsertScoreAssessments",
		body,
		fmt.Sprintf("CONSUME:%s", key),
		"consumed data successfully",
		"INFO",
		fmt.Sprintf("profile-%s", key))

	srp := map[string]any{
		"version": 1,
		"data": map[string]any{
			"assessment_code": msg.Data.AssessmentCode,
			"package_id":      msg.Data.PackageID,
			"student_id":      msg.Data.SmartBtwID,
		},
	}

	srpJson, err := sonic.Marshal(srp)
	if err != nil {
		return false
	}

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		time.Sleep(10 * time.Second)

		if err = db.Broker.Publish(
			"exam.central-assessment.generate-pdf",
			"application/json",
			[]byte(srpJson), // message to publish
		); err != nil {
			return err
		}
		return nil
	})

	go func() {
		err := g.Wait()
		if err != nil {
			fmt.Println("At least one calculation encountered an error: ", err.Error())
			return
		}
	}()

	if err := g.Wait(); err != nil {
		return false
	}

	return true
}
