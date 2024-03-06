package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
)

func GetClassMemberBySmIDElastic(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	studentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		message := fmt.Sprintf("Parameter smartbtw_id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}
	isOL := false
	isOLParam := c.Query("is_online")
	if isOLParam == "true" {
		isOL = true
	}

	res, err := lib.GetSingleClassMemberFromElastic(int32(studentId), isOL)
	if err != nil {
		if err.Error() != "record not found" {
			golog.Slack.ErrorWithData("error get class member", c.Body(), err)
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "error get class member",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}

func GetClassMemberByClassroomIDElastic(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.GetClassMemberFromElastic(c.Params("id"))
	if err != nil {
		if err.Error() != "record not found" {
			golog.Slack.ErrorWithData("error get class member", c.Body(), err)
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "error get class member",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "success",
	})
}
