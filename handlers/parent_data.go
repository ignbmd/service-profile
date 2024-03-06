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
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateParentData(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.CreateParentData)

	if err := c.BodyParser(req); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	err := lib.CreateParentData(req)

	if err != nil {
		golog.Slack.ErrorWithData("error insert parent data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error insert parent data",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    nil,
		"message": "Parent data created",
	})
}

func UpdateParentData(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.CreateParentData)
	if err := c.BodyParser(req); err != nil {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	valRes, er := govalidator.ValidateStruct(req)
	if !valRes {
		return c.Status(422).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   er,
		})
	}

	collection := db.Mongodb.Collection("students")
	coll := db.Mongodb.Collection("parent_datas")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var student models.Student
	err := collection.FindOne(ctx, bson.M{"smartbtw_id": req.SmartBtwID, "deleted_at": nil}).Decode(&student)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	payload := models.ParentData{
		StudentID:    student.ID,
		ParentName:   req.ParentName,
		ParentNumber: req.ParentNumber,
		CreatedAt:    student.CreatedAt,
		UpdatedAt:    time.Now(),
	}
	filter := bson.M{"student_id": student.ID}
	update := bson.M{"$set": payload}

	result, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		golog.Slack.ErrorWithData("error update parent data", c.Body(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "error update parent data",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "Parent data updated",
	})
}

func GetParentDataByStudentID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	collection := db.Mongodb.Collection("students")
	coll := db.Mongodb.Collection("parent_datas")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var student models.Student
	var parent models.ParentData
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	err = collection.FindOne(ctx, bson.M{"smartbtw_id": studentId, "deleted_at": nil}).Decode(&student)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	err = coll.FindOne(ctx, bson.M{"student_id": student.ID}).Decode(&parent)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Parent data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    parent,
	})
}
