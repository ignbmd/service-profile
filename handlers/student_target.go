package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateStudentTarget(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.CreateStudentTarget)

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

	res, err := lib.CreateStudentTargetRest(req)
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error insert student target", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error insert student target",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Student target created",
	})
}

func UpdateStudentTargetByID(c *fiber.Ctx) error {
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

	req := new(request.UpdateStudentTarget)
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

	err = lib.UpdateStudentTargetByID(req, id)
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error update student target", c.Body(), err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "error update student target",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Student target updated",
	})
}

func UpdateStudentTargetBySmartbtwID(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.UpdateStudentTarget)
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

	err1 := lib.UpdateStudentTargetBySmartbtwID(req, req.SmartbtwID)
	if err1 != nil {
		errSign := err1.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err1.Error(),
			})
		}
		golog.Slack.ErrorWithData("error update student target", c.Body(), err1)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "error update student target",
			"error":   err1.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Student target updated",
	})
}

func UpdateStudentPolbitTargetBySmartbtwID(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.UpdatePolbitType)
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

	err1 := lib.UpdateStudentPolbit(req, context.Background())
	if err1 != nil {
		errSign := err1.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err1.Error(),
			})
		}
		golog.Slack.ErrorWithData("error update student target polbit", c.Body(), err1)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "error update student target polbit",
			"error":   err1.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Student target polbit updated",
	})
}

func UpdateBulkStudentPolbitTargetBySmartbtwID(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.UpdatBulkPolbit)
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

	err1 := lib.UpdateBulkStudentPolbit(req, context.Background())
	if err1 != nil {
		errSign := err1.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err1.Error(),
			})
		}
		golog.Slack.ErrorWithData("error update bulk student target polbit", c.Body(), err1)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "error update bulk student target polbit",
			"error":   err1.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Student target polbit bulk update requested",
	})
}

func DeleteStudentTarget(c *fiber.Ctx) error {
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

	err = lib.DeleteStudentTarget(id)
	if err != nil {
		golog.Slack.ErrorWithData("delete student target", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "delete student target",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func GetStudentTargetByID(c *fiber.Ctx) error {
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

	res, err := lib.GetStudentTargetByID(id)
	if err != nil {
		golog.Slack.ErrorWithData("error get history ptk", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get history ptk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentTargetElastic(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.GetStudentTargetElasticBody)
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

	res, err := lib.GetStudentTargetElastic(req.SmartbtwID, req.TargetType, req.ProgramKey)
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error get student target from elastic", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student target from elastic",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentTarget(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.GetStudentsCompletedModules)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse body",
			"error":   err,
		})
	}

	if req.SmartBTWID == 0 || req.TargetType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "please specify the smartbtw_id & target_type",
			"error":   "request not correct",
		})

	}
	if req.TargetType == string(models.CPNS) {
		res, err := lib.GetStudentTargetCPNS(req.SmartBTWID)
		if err != nil {
			errSign := err.Error()[0:1]
			if errSign == "^" {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"success": false,
					"message": err.Error(),
				})
			}
			golog.Slack.ErrorWithData("error get student target", c.Body(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error get student target",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	}
	res, err := lib.GetStudentTargetByCustom(req.SmartBTWID, req.TargetType)
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error get student target", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student target",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetAllStudentTarget(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.GetStudentsCompletedModules)

	if err := c.QueryParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse body",
			"error":   err,
		})
	}

	if req.SmartBTWID == 0 || req.TargetType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "please specify the smartbtw_id & target_type",
			"error":   "request not correct",
		})

	}
	if strings.ToLower(req.TargetType) == "cpns" {
		res, err := lib.GetAllStudentTargetCPNS(req.SmartBTWID)
		if err != nil {
			errSign := err.Error()[0:1]
			if errSign == "^" {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
					"success": false,
					"message": err.Error(),
				})
			}
			golog.Slack.ErrorWithData("error get student target", c.Body(), err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error get student target from elastic",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	}
	res, err := lib.GetAllStudentTargetByCustom(req.SmartBTWID, req.TargetType)
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		golog.Slack.ErrorWithData("error get student target from elastic", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get student target from elastic",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func UpdateBulkStudentTargetBySmartbtwID(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	smId, err := strconv.Atoi(c.Params("smartbtw_id"))

	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to ID", c.Params("smartbtw_id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	req := new(request.UpdateStudentTargetRequest)
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

	err1 := lib.UpdateBulkStudentTargetBySmartbtwID(req, smId)
	if err1 != nil {
		errSign := err1.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err1.Error(),
			})
		}
		golog.Slack.ErrorWithData("error update student target", c.Body(), err1)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "error update student target",
			"error":   err1.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Student target updated",
	})
}

func UpdateStudentTarget(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.UpdateStudentTargetOne)
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

	err1 := lib.UpdateStudentTargetOne(req)
	if err1 != nil {
		errSign := err1.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err1.Error(),
			})
		}
		golog.Slack.ErrorWithData("error update student target", c.Body(), err1)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "error update student target",
			"error":   err1.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Student target updated",
	})
}

func GetSchoolCompetition(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	reqQuery := new(request.GetCompetitonList)
	if err := c.QueryParser(reqQuery); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err.Error(),
		})
	}

	res, err := lib.GetSchoolCompetitionList(c.Params("school_origin_id"), *reqQuery)
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		// golog.Slack.ErrorWithData("error get school competition", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get school competition",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func CountStudentWithUKACodeBySchoolOriginID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.CountUKACodeBySchoolOriginID(c.Params("school_origin_id"))
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error count student with uka code",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}
