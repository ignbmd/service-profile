package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
)

func GetStagesCompetitionList(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.GetStagesCompetitionList(c.Params("program"), c.Query("school_id"), c.Query("filter_type"))
	if err != nil {
		golog.Slack.ErrorWithData(fmt.Sprintf("error GetStagesCompetitionList (%s , %s, %s)", c.Params("program"), c.Query("school_id"), c.Query("filter_type")), nil, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "stages competition data had problem",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}
