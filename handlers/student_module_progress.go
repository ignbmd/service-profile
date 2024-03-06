package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func UpsertStudentModuleProgress(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.StudentModuleProgress)

	if err := c.BodyParser(req); err != nil {
		log.Println(err)
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	_, er := govalidator.ValidateStruct(req)
	if er != nil {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	res, err := lib.UpsertStudentModuleProgress(req)

	if err != nil {
		switch err.Error() {
		case "already exist":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "id already exist",
				"error":   err,
			})
		case "data module number greater than module total":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "data module number greater than module total",
				"error":   err,
			})
		default:
			golog.Slack.ErrorWithData("error insert student module progress data", c.Body(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error insert student module progress",
				"error":   err,
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Student Module Progress created",
	})
}

func CreateStudentModuleProgress(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.StudentModuleProgress)

	if err := c.BodyParser(req); err != nil {
		log.Println(err)
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	_, er := govalidator.ValidateStruct(req)
	if er != nil {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body validator",
			"error":   er,
		})
	}

	res, err := lib.CreateStudentModuleProgress(req)

	if err != nil {
		switch err.Error() {
		case "already exist":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "id already exist",
				"error":   err,
			})
		case "smartbtw_id not valid":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "smartbtw_id not valid",
				"error":   err,
			})
		case "must be greater than 0":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "data must be greater than 0",
				"error":   err,
			})
		case "data module number greater than module total":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "data module number greater than module total",
				"error":   err,
			})
		default:
			golog.Slack.ErrorWithData("error insert student module progress data", c.Body(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error insert student module progress data",
				"error":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Student Module Progress created",
	})
}

func UpdateStudentModuleProgress(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to Object ID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	req := new(request.UpdateStudentModuleProgress)
	if err := c.BodyParser(req); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	_, er := govalidator.ValidateStruct(req)
	if er != nil {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	err = lib.UpdateStudentModuleProgress(req, id)

	if err != nil {
		switch err.Error() {
		case "already exist":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "id already exist",
				"error":   err,
			})
		case "must be greater than 0":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "data must be greater than 0",
				"error":   err,
			})
		case "data module number greater than module total":
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": "data module number greater than module total",
				"error":   err,
			})
		case "data not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  false,
				"message": "data not found",
				"error":   err,
			})
		default:
			golog.Slack.ErrorWithData("error update student module progress data", c.Body(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error update student module progress data",
				"error":   err.Error(),
			})
		}

	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "student module progress data updated",
	})
}

func GetStudentModuleProgressBySmartBtwID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	collection := db.Mongodb.Collection("student_module_progress")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var StdModuleProgress models.StudentModuleProgress
	smartBtwId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	err = collection.FindOne(ctx, bson.M{"smartbtw_id": smartBtwId}).Decode(&StdModuleProgress)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student Module Progress data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    StdModuleProgress,
	})
}

func GetStudentModuleProgressByTaskID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	TaskId, err := strconv.Atoi(c.Params("task_id"))

	if err != nil {
		message := fmt.Sprintf("Paramer task_id of value: %s cannot be converted to int", c.Params("task_id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	res, err := lib.GetStudentModuleProgressByTaskID(TaskId)
	if res == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "student module progress data not found",
			"error":   err,
		})
	}

	if err != nil {
		golog.Slack.ErrorWithData("error get student module progress data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student module progress data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success get student module progress data",
	})
}
