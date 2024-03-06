package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func GetBranchList(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	query := new(request.Emails)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse body",
			"error":   err,
		})
	}

	res, err := lib.GetBranches()
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Branches data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})

}
