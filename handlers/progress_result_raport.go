package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
)

func GetProgressResultRaport(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.GetProgressResultRaport(c.Params("program"), c.Params("type"))
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
			"message": "error get progress result raport",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func TriggerBuildRaport(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID uint   `json:"smartbtw_id"`
		Program    string `json:"program"`
	}

	req := new(BodyRequest)
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

	err := lib.TriggerBuildRaport(req.SmartbtwID, req.Program)
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
			"message": "error trigger build raport",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func TriggerBuildRaportByTaskID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID uint   `json:"smartbtw_id"`
		Program    string `json:"program"`
		TaskID     int    `json:"task_id"`
	}

	req := new(BodyRequest)
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

	err := lib.BuildRaportByTaskID(req.SmartbtwID, req.Program, req.TaskID)
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
			"message": "error trigger build raport",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func ListingRaport(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID uint   `json:"smartbtw_id"`
		Program    string `json:"program"`
		StageType  string `json:"stage_type"`
		UKAType    string `json:"uka_type"`
	}

	type QueryRequest struct {
		Search string `json:"search" query:"search"`
	}

	req := new(BodyRequest)
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

	query := new(QueryRequest)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse body",
			"error":   err,
		})
	}

	res, err := lib.GetListingRaport(req.SmartbtwID, req.Program, req.UKAType, req.StageType, &query.Search)
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
			"message": "error get raport",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetProgressRaport(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID []uint `json:"smartbtw_id"`
		Program    string `json:"program"`
		UKAType    string `json:"uka_type"`
		StageType  string `json:"stage_type"`
	}

	req := new(BodyRequest)
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

	res, err := lib.GetProgressRaportList(req.SmartbtwID, req.UKAType, req.StageType, req.Program)
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
			"message": "error get raport",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func TriggerBuildProgressResult(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID uint   `json:"smartbtw_id"`
		Program    string `json:"program"`
		StageType  string `json:"stage_type"`
	}

	req := new(BodyRequest)
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

	err := lib.SendToGenerateRaport(req.Program, req.SmartbtwID, req.StageType)
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
			"message": "error trigger build progress raport",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func BuildRaportBulk(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID []uint `json:"smartbtw_id"`
	}

	req := new(BodyRequest)
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

	err := lib.SendRequestBuildRaportBulk(req.SmartbtwID, c.Params("program"))
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
			"message": "error trigger build raport bulk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func RequestBuildProgressRaportBulk(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID []uint `json:"smartbtw_id"`
		Program    string `json:"program"`
		StageType  string `json:"stage_type"`
	}

	req := new(BodyRequest)
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

	err := lib.RequestGenerateProgressRaport(req.SmartbtwID, req.Program, req.StageType)
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
			"message": "error trigger build progress raport bulk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func ReGenerateRaportBulk(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	type BodyRequest struct {
		SmartbtwID []uint `json:"smartbtw_id"`
		Program    string `json:"program"`
		StageType  string `json:"stage_type"`
	}

	req := new(BodyRequest)
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

	err := lib.AddToQueueRegenerateRaport(req.SmartbtwID, req.Program, req.StageType)
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
			"message": "error trigger build progress raport bulk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}
