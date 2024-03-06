package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func GetHistoryScoreByTargetType(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	targetType := strings.ToLower(c.Params("target_type"))
	query := new(request.HistoryScoreQueryParams)
	var finalRes interface{}

	if targetType != "ptk" && targetType != "ptn" && targetType != "cpns" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "target type must be either ptk, ptn or cpns",
			"error":   "invalid target type",
		})
	}

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse body",
			"error":   err,
		})
	}

	if query.OnlyLatestScore != nil && (query.SmartBTWID == nil || query.TaskID == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "please specify the smartbtw_id & task_id",
			"error":   "incomplete required query string",
		})

	}

	res, err := lib.GetHistoryScoreByTargetType(targetType, query)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   res,
			"message": fmt.Sprintf("failed on getting %s history score data", targetType),
		})
	}
	finalRes = res
	if query.OnlyLatestScore != nil {
		if *query.OnlyLatestScore {
			if len(res) > 0 {
				finalRes = res[0]
			} else {
				finalRes = nil
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    finalRes,
		"message": "Success",
	})
}

func GetStudentUKAScores(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.GetUkaCodeScoresByEmail(c.Params("email"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "History score not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}
