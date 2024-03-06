package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func GetStudentWalletBalance(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	var balance interface{}
	SmartBTWID, err := strconv.Atoi(c.Params("smartbtw_id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	wallet, err := lib.GetStudentWalletTotalBalance(SmartBTWID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student wallet balance not found",
			"error":   err,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student wallet balance not found",
			"error":   err,
		})
	}

	if len(wallet) > 0 {
		balance = wallet[0].Balance
	} else {
		balance = nil
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    balance,
		"message": "Get wallet balance",
	})
}

func GetStudentWalletDetailBalance(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	SmartBTWID, err := strconv.Atoi(c.Params("smartbtw_id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	wallet, err := lib.GetStudentWalletBalance(SmartBTWID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student wallet balance not found",
			"error":   err,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student wallet balance not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    wallet,
		"message": "Get wallet balance",
	})
}

func GetStudentWalletHistory(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	SmartBTWID, err := strconv.Atoi(c.Params("smartbtw_id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	params := new(request.GetWalletHistory)
	if err := c.QueryParser(params); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse query param",
			"error":   err,
		})
	}

	walletHistoryData, err := lib.GetStudentWalletHistory(SmartBTWID, params.Type)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student wallet history not found",
			"error":   err,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Student wallet history not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    walletHistoryData,
		"message": "Get wallet history",
	})
}

func ChargeStudentWallet(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.ChargeWallet)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	err := lib.ChargeStudentWallet(req)
	if err != nil {
		errSign := err.Error()[0:1]
		if errSign == "^" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  false,
				"message": strings.Replace(err.Error(), "^", "", 1),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed charge student wallet",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Wallet Balance successfully deducted",
	})
}

func CheckCoin(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	SmartBTWID, err := strconv.Atoi(c.Params("smartbtw_id"))
	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	err = lib.CheckCoinDrillo(SmartBTWID)
	if err != nil {
		errSign := err.Error()
		if errSign == "BTW coin anda tidak cukup, silahkan top up" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
				"message": "silahkan top up coin",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to process data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
