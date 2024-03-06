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

func CreateStudentAccess(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.CreateStudentAccess)
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

	res, err := lib.CreateStudentAccess(req)
	// fmt.Println(err)
	if err != nil {
		if err.Error() == "record already exist" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": "user already has this record",
				"error":   err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error insert create student access", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error insert student data data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Student access created",
	})
}

func CreateStudentAccessBulk(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.CreateStudentAccessBulk)
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

	res, err := lib.CreateStudentAccessBulk(req)
	// fmt.Println(err)
	if err != nil {
		if err.Error() == "record already exist" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": "user already has one of the record",
				"error":   err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error insert create student access", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error insert student data data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Student access created",
	})
}

func DeleteStudentAccess(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id := c.Params("id")

	smId, err := strconv.Atoi(id)

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	req := new(request.DeleteStudentAccess)
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

	err = lib.DeleteStudentAccess(smId, req.DisallowedAccess, req.AppType)
	if err != nil {
		golog.Slack.ErrorWithData("delete student access", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "delete student access",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func GetStudentAccessListBySmartBTWID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id := c.Params("id")
	appType := c.Query("app_type")

	smId, err := strconv.Atoi(id)

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	res, err := lib.GetStudentAllowedAccess(smId, appType)
	if err != nil {
		golog.Slack.ErrorWithData("error get student allowed access list", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student allowed access list",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentAccessListFromElastic(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id := c.Params("id")
	appType := c.Query("app_type")

	smId, err := strconv.Atoi(id)

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	res, err := lib.GetStudentAccessElastic(smId, appType)
	if err != nil {
		golog.Slack.ErrorWithData("error get student allowed access list", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student allowed access list",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentAccessListByCodeFromElastic(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	query := new(request.GetStudentAccessElastic)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse body",
			"error":   err,
		})
	}
	fmt.Println(query)

	res, err := lib.GetStudentAccessByCodeElastic(query)
	if err != nil {
		golog.Slack.ErrorWithData("error get student allowed access list", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student allowed access list",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentAccessListElastic(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id := c.Params("id")
	appType := c.Query("app_type")

	smId, err := strconv.Atoi(id)

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	res, err := lib.GetStudentAccessElastic(smId, appType)
	if err != nil {
		golog.Slack.ErrorWithData("error get student allowed access list", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student allowed access list",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentAccessCode(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id := c.Params("id")
	accessCode := c.Query("access_code")
	appType := c.Query("app_type")

	smId, err := strconv.Atoi(id)

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	res, err := lib.GetSingleStudentAccessElastic(smId, accessCode, appType)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"success": true,
				"data":    nil,
				"message": "Success",
			})
		}
		golog.Slack.ErrorWithData("error get student allowed access list", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student allowed access list",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func DeleteStudentAccessBulk(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	id := c.Params("id")

	smId, err := strconv.Atoi(id)

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	req := new(request.DeleteStudentAccess)
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

	err = lib.DeleteStudentAccessBulk(smId, req.DisallowedAccessBulk, req.AppType)
	if err != nil {
		golog.Slack.ErrorWithData("delete student access bulk", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "delete student access bulk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}
