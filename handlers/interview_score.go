package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateInterviewScore(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.UpsertInterviewScore)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		errs := er.(govalidator.Errors).Errors()
		for _, e := range errs {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": e.Error(),
			})
		}
	}

	data, _ := lib.GetSingleInterviewScoreBySessionIDSSOIDAndStudentID(req.SessionID, req.CreatedBy.ID, req.SmartBtwID)
	if (data != models.InterviewScore{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "data nilai wawancara siswa sudah ada di sesi ini",
		})
	}

	profile, _ := lib.GetStudentProfileElastic(req.SmartBtwID)
	if (profile == request.StudentProfileElastic{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("data siswa dengan id %d tidak ditemukan", req.SmartBtwID),
		})
	}

	_, err := lib.CreateInterviewScore(req)
	if err != nil {
		golog.Slack.ErrorWithData("error create interview score", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error create interview score",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "interview score has been created",
	})
}

func UpdateInterviewScore(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.UpsertInterviewScore)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		errs := er.(govalidator.Errors).Errors()
		for _, e := range errs {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": e.Error(),
			})
		}
	}

	profile, _ := lib.GetStudentProfileElastic(req.SmartBtwID)
	if (profile == request.StudentProfileElastic{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("data siswa dengan id %d tidak ditemukan", req.SmartBtwID),
		})
	}

	err := lib.UpdateInterviewScore(c.Params("id"), req)
	if err != nil {
		if err.Error() == "interview score data not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error update interview score", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error update interview score",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "interview score has been updated",
	})
}

func GetInterviewScoreByArrayOfEmail(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.GetInterviewScoreByArrEmail)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		errs := er.(govalidator.Errors).Errors()
		for _, e := range errs {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": e.Error(),
			})
		}
	}

	for _, email := range req.Email {
		if !govalidator.IsEmail(email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "email tidak valid",
			})
		}
	}

	res, err := lib.GetInterviewScoreByArrayOfStudentEmail(req.Email, req.Year)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "interview score not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}

func GetSingleInterviewScoreBySMIDAndYear(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	studentId, err := strconv.Atoi(c.Params("smartbtw_id"))
	if err != nil {
		message := fmt.Sprintf("Parameter smartbtw_id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	year, err := strconv.Atoi(c.Params("year"))
	if err != nil {
		message := fmt.Sprintf("Parameter year of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	collection := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var interviewScore []bson.M
	pipelines := aggregates.GetSingleInterviewScoreBySMIDAndYear(studentId, year)
	aggregateOptions := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipelines, aggregateOptions...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Interview score not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &interviewScore)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Interview score not found",
			"error":   err,
		})
	}

	if len(interviewScore) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Interview score not found",
			"data":    interviewScore,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    interviewScore[0],
		})
	}
}

func GetInterviewScoresByInterviewSessionIDAndSSOID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	collection := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	interviewSessionID, err := primitive.ObjectIDFromHex(c.Params("session_id"))
	if err != nil {
		message := fmt.Sprintf("Could not convert session_id of value: %s to objectID", c.Params("session_id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	var results []bson.M

	pipel := aggregates.GetInterviewScoresByInterviewSessionIDAndSSOID(interviewSessionID, c.Params("sso_id"))
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "could not find interview score data",
			"error":   err.Error(),
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "could not find interview score data",
			"error":   err.Error(),
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "interview score data was not found",
			"error":   nil,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    results,
		})
	}
}

func GetSingleInterviewScoreByID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	collection := db.Mongodb.Collection("interview_score")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	interviewScoreID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		message := fmt.Sprintf("Could not convert interview_score_id of value: %s to objectID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	var results []bson.M

	pipel := aggregates.GetSingleInterviewScoreByID(interviewScoreID)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "could not find interview score data",
			"error":   err.Error(),
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "could not find interview score data",
			"error":   err.Error(),
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "interview score data was not found",
			"error":   nil,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    results[0],
		})
	}
}
