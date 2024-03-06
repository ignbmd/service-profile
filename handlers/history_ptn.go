package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func CreateHistoryPtn(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.CreateHistoryPtn)
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

	_, err := lib.CreateHistoryPtn(req)
	if err != nil {
		golog.Slack.ErrorWithData("error insert history ptn", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error insert history ptn",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "history ptn created",
	})
}

func UpdateHistoryPtn(c *fiber.Ctx) error {
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

	req := new(request.UpdateHistoryPtn)
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

	err = lib.UpdateHistoryPtn(req, id)
	if err != nil {
		golog.Slack.ErrorWithData("error update history ptn", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error update history ptn",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "history ptn updated",
	})
}

func DeleteHistoryPtn(c *fiber.Ctx) error {
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

	err = lib.DeleteHistoryPtn(id)
	if err != nil {
		golog.Slack.ErrorWithData("delete history ptn", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "delete history ptn",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Success",
	})
}

func GetHistoryPtnByID(c *fiber.Ctx) error {
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

	res, err := lib.GetHistoryPtnByID(id)
	if err != nil {
		golog.Slack.ErrorWithData("error get history ptn", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get history ptn",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res[0],
		"message": "Success",
	})
}

func GetHistoryPtnBySmartBTWID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	studentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		message := fmt.Sprintf("Paramer id of value: %s cannot be converted to integer", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err.Error(),
		})
	}

	params := new(request.HistoryPTNQueryParams)
	if err := c.QueryParser(params); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse query param",
			"error":   err,
		})
	}

	if params.Limit != nil {
		if *params.Limit <= 0 {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"message": "Limit must be a positive number",
				"error":   "Invalid limit value",
			})
		}
	}

	res, err := lib.GetHistoryPtnBySmartBTWID(studentId, params)
	if err != nil {
		golog.Slack.ErrorWithData("error get history ptn", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get history ptn",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentAveragePtn(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	collection := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []bson.M
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	pipel := aggregates.GetStudentPtnAverage(studentId, c.Query("program_key"))
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    results[0],
		})
	}
}

func GetStudentPtnLastScore(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	collection := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []bson.M
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	pipel := aggregates.GetStudentPtnLastScore(studentId, c.Query("program_key"))
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    results,
		})
	}
}

func GetLast10StudentScorePtn(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	collection := db.Mongodb.Collection("history_ptn")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []bson.M
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	pipel := aggregates.GetLast10StudentPtnScore(studentId, c.Query("program_key"))
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "History ptn not found",
			"error":   err,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    results,
		})
	}
}

func GetStudentFreePTN(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	smartbtwId, err := strconv.Atoi(c.Params("smartbtw_id"))
	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("smartbtw_id"))
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	res, err := lib.GetHistoryFreeSingleStudentPTN(smartbtwId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": true,
				"message": "free uka data not found for this smartbtw_id",
				"error":   nil,
				"data":    nil,
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "internal server error",
				"error":   err.Error(),
				"data":    nil,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "ok",
		"error":   nil,
		"data":    res,
	})
}

func GetStudentPremiumUKAPTN(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	smartbtwId, err := strconv.Atoi(c.Params("smartbtw_id"))
	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("smartbtw_id"))
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	res, err := lib.GetHistoryPremiumUKASingleStudentPTN(smartbtwId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": true,
				"message": "premium uka data not found for this smartbtw_id",
				"error":   nil,
				"data":    nil,
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "internal server error",
				"error":   err.Error(),
				"data":    nil,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "ok",
		"error":   nil,
		"data":    res,
	})
}

func GetStudentPackageUKAPTN(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	smartbtwId, err := strconv.Atoi(c.Params("smartbtw_id"))
	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("smartbtw_id"))
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	res, err := lib.GetHistoryPackageUKASingleStudentPTN(smartbtwId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": true,
				"message": "package uka data not found for this smartbtw_id",
				"error":   nil,
				"data":    nil,
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "internal server error",
				"error":   err.Error(),
				"data":    nil,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "ok",
		"error":   nil,
		"data":    res,
	})
}

func GetALLStudentScorePtn(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.HistoryPtnGetAll)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	if !req.FromElastic {
		results, totalData, err := lib.GetALLStudentScorePtnPagination(req.SmartbtwID, req.Limit, req.Page, req.ProgramKey)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": true,
					"message": "data not found for this smartbtw_id",
					"error":   nil,
					"data":    nil,
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "internal server error",
					"error":   err.Error(),
					"data":    nil,
				})
			}
		}
		historyRes := map[string]interface{}{
			"total_histories": totalData,
			"histories":       results,
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "ok",
			"error":   nil,
			"data":    historyRes,
		})
	} else {
		historyRes, err := lib.GetStudentHistoryPTNElastic(req.SmartbtwID, req.IsStagesHistory, req.ProgramKey)
		if err != nil {

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "internal server error",
				"error":   err.Error(),
				"data":    nil,
			})

		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "ok",
			"error":   nil,
			"data":    historyRes,
		})
	}
}

func GetHistoryPTNByTaskID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	taskID, err := strconv.Atoi(c.Params("task_id"))
	if err != nil {
		message := fmt.Sprintf("Paramer task_id of value: %s cannot be converted to int", c.Params("task_id"))
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	res, err := lib.GetRankPTNByTaskID(uint(taskID))
	if err != nil {
		golog.Slack.ErrorWithData("error get rank", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error get rank",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "successfuly get rank",
	})
}
