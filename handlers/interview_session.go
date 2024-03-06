package handlers

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func GetAllInterviewSessions(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.GetAllInterviewSessions()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "interview sessions not found",
			"error":   err,
		})
	}

	if len(res) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "interview sessions not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}

func GetSingleInterviewSessionByID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	interviewSessionID, interviewSessionIDErr := primitive.ObjectIDFromHex(c.Params("id"))
	if interviewSessionIDErr != nil {
		message := fmt.Sprintf("Parameter id of value: %s is not a valid object ID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   interviewSessionIDErr.Error(),
		})
	}

	interviewSession, interviewSessionErr := lib.GetSingleInterviewSessionByID(interviewSessionID)
	if interviewSessionErr != nil {
		if interviewSessionErr.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "fail to soft delete interview session",
				"error":   "data not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "interview session not found",
			"error":   interviewSessionErr,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    interviewSession,
	})
}

func CreateInterviewSession(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.InterviewSessionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
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

	res, err := lib.CreateInterviewSession(req)
	if err != nil {
		golog.Slack.ErrorWithData("error create interview session", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error create interview session",
			"error":   err,
		})
	}

	insertedID := res.InsertedID
	insertedObjectID := insertedID.(primitive.ObjectID).Hex()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"id": insertedObjectID,
		},
		"message": "interview session created successfully",
	})
}

func UpdateInterviewSession(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	interviewSessionID, interviewSessionIDErr := primitive.ObjectIDFromHex(c.Params("id"))
	if interviewSessionIDErr != nil {
		message := fmt.Sprintf("Parameter ID of value: %s is not a valid object ID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   interviewSessionIDErr.Error(),
		})
	}

	req := new(request.InterviewSessionRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
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

	upd, err := lib.UpdateInterviewSession(interviewSessionID, req)
	if err != nil {
		golog.Slack.ErrorWithData("error update interview session", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error update interview session",
			"error":   err,
		})
	}

	if upd.ModifiedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "error update interview session",
			"error":   "data not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "interview session has been updated",
	})
}

func SoftDeleteInterviewSession(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	interviewSessionID, interviewSessionIDErr := primitive.ObjectIDFromHex(c.Params("id"))
	if interviewSessionIDErr != nil {
		message := fmt.Sprintf("Parameter ID of value: %s is not a valid object ID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   interviewSessionIDErr.Error(),
		})
	}

	_, interviewSessionErr := lib.GetSingleInterviewSessionByID(interviewSessionID)
	if interviewSessionErr != nil {
		if interviewSessionErr.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "fail to soft delete interview session",
				"error":   "data not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "fail to soft delete interview session",
			"error":   interviewSessionErr,
		})
	}

	_, err := lib.SoftDeleteInterviewSession(interviewSessionID)
	if err != nil {
		golog.Slack.ErrorWithData("error soft deleting interview session", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error soft deleting interview session",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "interview session has been soft deleted",
	})
}
