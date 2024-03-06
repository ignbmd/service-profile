package handlers

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func CreateRecordScore(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.ScoreSkdBkn)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er.Error(),
		})
	}

	res, err := lib.CreateScoreData(req)
	// fmt.Println(err)
	if err != nil {
		golog.Slack.ErrorWithData("error insert score data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error insert score data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Score data created",
	})
}

func UpdateRecordScore(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to Object ID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	req := new(request.UpdateScoreSKDBKN)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er.Error(),
		})
	}

	err = lib.UpdateScoreData(req, id)

	if err != nil {
		golog.Slack.ErrorWithData("error update score data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error update score data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Score data updated",
	})
}

func GetRecordScoreByStudent(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	res, err := lib.GetScoreDataByStudent(studentId)

	if err != nil {
		golog.Slack.ErrorWithData("error get score data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get score data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetRecordScoreById(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to Object ID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	res, err := lib.GetSingleScoreData(id)
	if err != nil {
		golog.Slack.ErrorWithData("error get score data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get score data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res[0],
		"message": "Success",
	})
}

func DeleteRecordScoreById(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to Object ID", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	err = lib.DeleteScoreSingleRecord(id)
	if err != nil {
		golog.Slack.ErrorWithData("delete score data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "delete score data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func GetManyRecordScore(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.GetManyRecordScore)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er.Error(),
		})
	}
	res, err := lib.GetManyStudentScoreByYear(req.SmartBtwID, req.Year)

	if err != nil {
		golog.Slack.ErrorWithData("error get score data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get score data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})

}
