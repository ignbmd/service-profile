package handlers

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func UpsertBKNScore(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.UpsertBKNScore)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	err := lib.UpsertBKNScore(req)
	if err != nil {
		golog.Slack.ErrorWithData("error upsert bkn score", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error upsert bkn score",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "bkn score created or updated",
	})
}

func UpdateBKNScoreSurvey(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.UpdateBKNScoreForSurvey)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	err := lib.UpdateForSurvey(req)
	if err != nil {
		golog.Slack.ErrorWithData("error update bkn score", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error upsert bkn score",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "bkn score updated",
	})
}

func UpdateBKNScoreProdi(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.UpdateBKNScoreForProdi)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	err := lib.UpdateForProdi(req)
	if err != nil {
		golog.Slack.ErrorWithData("error update bkn score prodi", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error upsert bkn score",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "bkn score prodi updated",
	})
}

func GetBKNScoreByArrayOfSmartbtwID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.GetBKNScoreByArrSmID)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	res, err := lib.GetBKNScoreByArrayOfIDStudent(req.SmartBtwID, req.Year)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "bkn score not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}

func GetBKNScoreByArrayOfEmail(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.GetBKNScoreByArrEmail)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	res, err := lib.GetBKNScoreByArrayOfStudentEmail(req.Email, req.Year)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "bkn score not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}

func GetBKNScoreByArrayOfEmailForGDS(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	req := new(request.GetBKNScoreByArrEmail)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	res, err := lib.GetBKNScoreByArrayOfStudentEmailGDS(req.Email, req.Year)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "bkn score not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}

func GetSingleBKNScoreBySMIDAndYear(c *fiber.Ctx) error {
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

	res, err := lib.GetSingleBKNScoreByYearAndStudent(studentId, uint16(year))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}
