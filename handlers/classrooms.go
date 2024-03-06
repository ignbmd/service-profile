package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
)

func GetClassroomsByBranchCodes(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.GetClassroomsByBranchCodes(c.Params("code"))
	if err != nil {
		if err.Error() != "record not found" {
			golog.Slack.ErrorWithData("error get classrooms", c.Body(), err)
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "error get classrooms",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}
