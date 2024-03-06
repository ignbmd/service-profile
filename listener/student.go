package listener

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func ListenStudentBinding(msg *amqp.Delivery) bool {
	switch msg.RoutingKey {
	case "user.created", "user.import", "user.updated":
		if StudentUpsert(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "user.deleted":
		if StudentDelete(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "user.create-elastic":
		// if StudentElastic(msg.Body, msg.RoutingKey) {
		lib.MessageAck(msg)
		return true
		// }
	case "user.upsert-profile-elastic":
		if UpsertProfileElastic(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "user.upsert-compmap-elastic":
		if UpsertCompMapProfileElastic(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "user.binsus.sync":
		if SyncBinsusProfile(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	case "user.binsus.final-sync":
		if SyncBinsusFinalProfile(msg.Body, msg.RoutingKey) {
			lib.MessageAck(msg)
			return true
		}
	}
	return false
}

func StudentUpsert(body []byte, key string) bool {
	msg, err := request.UnmarshalMessageStudentBody(body)

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

		return strings.Contains(err.Error(), "user_tryout_id")
	}

	if msg.Version < 1 {
		text := fmt.Sprintf("[RabbitMQ][%s] Invalid Version", key)
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}

	if msg.Data.AccountType == "" {
		msg.Data.AccountType = "smartbtw"
	}

	_, err = lib.UpsertStudent(&msg.Data)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when creating data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if key == "user.import" {
		text := fmt.Sprintf("User with smartbtw_id %d has been imported", msg.Data.ID)
		log.Println(text)
	} else if key == "user.created" {
		payPar := request.CreateParentData{
			SmartBtwID:   msg.Data.ID,
			ParentName:   msg.Data.ParentName,
			ParentNumber: msg.Data.ParentNumber,
		}
		err = lib.CreateParentData(&payPar)
		if err != nil {
			text := fmt.Sprintf("Error insert parent data for smartbtw_id %d", msg.Data.ID)
			golog.Slack.ErrorWithData(text, body, err)
		}
		text := fmt.Sprintf("User with smartbtw_id %d has been created", msg.Data.ID)
		log.Println(text)
	} else if key == "user.updated" {
		text := fmt.Sprintf("User with smartbtw_id %d has been updated", msg.Data.ID)
		log.Println(text)
	}
	return true
}

func StudentDelete(body []byte, key string) bool {
	msg, err := request.UnmarshalDeleteStudentBodyMessage(body)

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

	_, err = lib.DeleteStudent(&msg)

	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when deleting data", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func StudentElastic(body []byte, key string) bool {
	var (
		sPtkID   uint
		sPtkName string
		sPtkC    time.Time
		mPtkID   uint
		mPtkName string

		sPtnID   uint
		sPtnName string
		sPtnC    time.Time
		mPtnID   uint
		mPtnName string

		gender string
	)
	msg, err := request.UnmarshalCreateStudentProfileElastic(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(err)
		return false
	}

	std, err := lib.GetStudentOnlyBySmartBTWID(msg.Data.SmartbtwID)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error get student's profile", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	if len(std) == 0 {
		text := fmt.Sprintf("[RabbitMQ][%s] Error get student's profile", key)
		golog.Slack.WarningWithData(text, body, errors.New("student raw data not found, please trigger user.created"))
		return true
	}

	//getTargetPTK
	tPtk, err := lib.GetStudentTargetByCustom(msg.Data.SmartbtwID, "PTK")
	if err != nil {
		sPtkID = 0
		sPtkName = ""
		sPtkC = time.Now()
		mPtkID = 0
		mPtkName = ""
	} else {
		sPtkID = uint(tPtk.SchoolID)
		sPtkName = tPtk.SchoolName
		sPtkC = tPtk.CreatedAt
		mPtkID = uint(tPtk.MajorID)
		mPtkName = tPtk.MajorName
	}

	//getTargetPTN
	tPtn, err := lib.GetStudentTargetByCustom(msg.Data.SmartbtwID, "PTN")
	if err != nil {
		sPtnID = 0
		sPtnName = ""
		sPtnC = time.Now()
		mPtnID = 0
		mPtnName = ""
	} else {
		sPtnID = uint(tPtn.SchoolID)
		sPtnName = tPtn.SchoolName
		sPtnC = tPtn.CreatedAt
		mPtnID = uint(tPtn.MajorID)
		mPtnName = tPtn.MajorName
	}

	if std[0].Gender == 1 {
		gender = "L"
	} else {
		gender = "P"
	}

	phoneNumber := "0"
	if std[0].Phone != nil {
		phoneNumber = *std[0].Phone
	}

	re := request.StudentProfileElastic{
		SmartbtwID:     msg.Data.SmartbtwID,
		Name:           std[0].Name,
		Email:          std[0].Email,
		Photo:          std[0].Photo,
		Phone:          phoneNumber,
		Gender:         gender,
		BranchCode:     std[0].BranchCode,
		BirthDate:      msg.Data.BirthDate,
		Province:       msg.Data.Province,
		ProvinceID:     msg.Data.ProvinceID,
		Region:         msg.Data.Region,
		RegionID:       msg.Data.RegionID,
		LastEdID:       msg.Data.LastEdID,
		LastEdName:     msg.Data.LastEdName,
		LastEdType:     msg.Data.LastEdType,
		LastEdMajor:    msg.Data.LastEdMajor,
		LastEdMajorID:  msg.Data.LastEdMajorID,
		LastEdRegion:   msg.Data.LastEdRegion,
		LastEdRegionID: msg.Data.LastEdRegionID,
		EyeColorBlind:  msg.Data.EyeColorBlind,
		Height:         msg.Data.Height,
		Weight:         msg.Data.Weight,
		AccountType:    msg.Data.AccountType,
		SchoolPTKID:    sPtkID,
		SchoolNamePTK:  sPtkName,
		MajorNamePTK:   mPtkName,
		MajorPTKID:     mPtkID,
		CreatedAtPTK:   sPtkC,
		SchoolPTNID:    sPtnID,
		SchoolNamePTN:  sPtnName,
		MajorPTNID:     mPtnID,
		MajorNamePTN:   mPtnName,
		CreatedAtPTN:   sPtnC,
		CreatedAt:      std[0].CreatedAt,
	}

	err = lib.InsertStudentProfileElastic(&re)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func SyncBinsusProfile(body []byte, key string) bool {

	msg, err := request.UnmarshalSyncBinsusStudentProfileElastic(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(err)
		return false
	}

	err = lib.UpsertBinsusStudentProfile(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func SyncBinsusFinalProfile(body []byte, key string) bool {

	msg, err := request.UnmarshalSyncBinsusFinalStudentProfileElastic(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(err)
		return false
	}

	err = lib.UpsertBinsusFinalStudentProfile(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func UpsertElastic(body []byte, key string) bool {

	msg, err := request.UnmarshalCacheStudentProfileElastic(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(err)
		return false
	}

	err = lib.CacheStudentProfileToElastic(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when inserting student's profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}

	return true
}

func UpsertProfileElastic(body []byte, key string) bool {

	msg, err := request.UnmarshalCacheStudentProfileElastic(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(err)
		return false
	}

	err = lib.CacheStudentProfileToElastic(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when updating student's profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	fmt.Printf("Student %d cached to elastic\n", msg.Data.SmartBTWID)

	return true
}

func UpsertCompMapProfileElastic(body []byte, key string) bool {

	msg, err := request.UnmarshalCacheStudentProfileElastic(body)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Failed to decode from json", key)

		golog.Slack.ErrorWithData(text, body, err)

		log.Println(err)
		return false
	}

	err = lib.CacheStudentCompMapProfileToElastic(&msg.Data)
	if err != nil {
		text := fmt.Sprintf("[RabbitMQ][%s] Error when updating compmap student's profile to elasticsearch", key)
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
		return false
	}
	fmt.Printf("Student %d comp map cached to elastic\n", msg.Data.SmartBTWID)

	return true
}
