package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/request"
)

func GetStudents(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	params := new(request.GetStudentsParams)
	if err := c.QueryParser(params); err != nil {
		return err
	}

	smartbtw_ids := []int{}
	for _, v := range params.SmartBTWID {
		id, err := strconv.Atoi(v)
		if err != nil {
			message := fmt.Sprintf("Parameter smartbtw_id of value: %s cannot be converted to int", v)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": message,
				"error":   err,
			})
		}
		smartbtw_ids = append(smartbtw_ids, id)
	}

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []bson.M
	pipel := aggregates.GetStudents(smartbtw_ids, params.Fields)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    results,
	})
}

func GetSingleStudent(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	collection := db.Mongodb.Collection("students")
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

	pipel := aggregates.GetStudentWithParents(studentId)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    results[0],
		})
	}

}

func GetStudentCaching(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	studentData, err := lib.GetStudentProfileElastic(studentId)
	if err != nil {
		if err.Error() == "student data not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    studentData,
	})

}

func GetSingleStudentFromElastic(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	studentData, err := lib.GetStudentProfileElastic(studentId)
	if err != nil {
		if err.Error() == "student data not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    studentData,
	})

}

func GetStudentJoinedClassType(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	studentId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		message := fmt.Sprintf("Paramer smartbtw_id of value: %s cannot be converted to int", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	studentData, err := lib.GetStudentJoinedClassType(studentId)
	if err != nil {
		if err.Error() == "student data not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    studentData,
	})

}

func GetStudentsByBranchCode(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	var results []bson.M

	req := new(request.GetStudentByBranchCode)

	if err := c.QueryParser(req); err != nil {
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

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentByBranchCodeAndPagination(req.BranchCode, req.Skip, req.Limit, req.Page, req.Search)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    results[0],
	})
}

func GetStudentsByBranchCodeNoLimit(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	var results []bson.M

	req := new(request.GetStudentByBranchCode)

	if err := c.QueryParser(req); err != nil {
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

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentByBranchCodeNoLimit(req.BranchCode, req.Skip, req.Limit, req.Page, req.Search)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    results[0],
	})
}

func GetStudentsByArrayBranchCode(c *fiber.Ctx) error {

	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	var results []bson.M

	req := new(request.GetStudentByBranchCodeArray)

	if err := c.QueryParser(req); err != nil {
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
	pars := strings.Split(req.BranchCode, ",")

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipel := aggregates.GetStudentByBranchCodeArrayAndPagination(pars, req.Skip, req.Limit, req.Page, req.Search)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    results[0],
	})
}

func GetStudentCompletedModules(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	query := new(request.GetStudentsCompletedModules)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse body",
			"error":   err,
		})
	}

	if query.SmartBTWID == 0 || query.TargetType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "please specify the smartbtw_id & target_type",
			"error":   "request not correct",
		})

	}

	res, err := lib.GetStudentCompletedModulesBySmartBTWID(query.SmartBTWID, query.TargetType)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentBranch(c *fiber.Ctx) error {
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

	res, err := lib.GetStudentBrachByEmails(query.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})

}

func GetStudentResultForPerformaSiswa(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.GetPerformaSiswa)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	res, avg, err := lib.GetPerformaSiswa(req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"average": avg,
		"message": "Success",
	})
}

func GetStudentResultForPerformaSiswaUKA(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.GetPerformaSiswaUKA)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}
	if req.Program == "ptk" {
		res, avg, err := lib.GetPerformaSiswaPTK(req)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"average": avg,
			"message": "Success",
		})
	} else if req.Program == "cpns" {
		res, avg, err := lib.GetPerformaSiswaCPNS(req)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"average": avg,
			"message": "Success",
		})
	} else if req.Program == "ptn" {
		res, avg, err := lib.GetPerformaSiswaPTN(req)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"average": avg,
			"message": "Success",
		})
	}
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"success": false,
		"message": "program not valid",
	})
}

func FetchRaport(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	id, err := strconv.Atoi(c.Params("smId"))
	if err != nil {
		message := fmt.Sprintf("Parameter smartbtw_id of value: %s cannot be converted to int", c.Params("smId"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	fil := c.Query("filter")

	switch strings.ToLower(c.Params("program")) {
	case "ptn":

		res, err := lib.FetchPTNRaport(uint(id), fil)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	case "cpns":

		res, err := lib.FetchCPNSRaport(uint(id), fil)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	}

	res, err := lib.FetchPTKRaport(uint(id), fil, c.Query("raport_type"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func FetchRaportUKA(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	id, err := strconv.Atoi(c.Params("smId"))
	if err != nil {
		message := fmt.Sprintf("Parameter smartbtw_id of value: %s cannot be converted to int", c.Params("smId"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	fil := c.Query("filter")
	typStg := c.Query("type_stage")
	switch strings.ToLower(c.Params("program")) {
	case "ptn":

		res, err := lib.FetchPTNUKARaport(uint(id), typStg, fil)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	case "cpns":

		res, err := lib.FetchCPNSUKARaport(uint(id), typStg, fil)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	}

	res, err := lib.FetchPTKUKARaport(uint(id), typStg, fil)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func FetchRaportPrePostTest(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	id, err := strconv.Atoi(c.Params("smId"))
	if err != nil {
		message := fmt.Sprintf("Parameter smartbtw_id of value: %s cannot be converted to int", c.Params("smId"))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": message,
			"error":   err,
		})
	}

	switch strings.ToLower(c.Params("program")) {
	case "ptn":

		res, err := lib.FetchPTNPrePostRaport(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	case "cpns":

		res, err := lib.FetchCPNSPrePostRaport(uint(id))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Student data not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	}

	res, err := lib.FetchPTKPrePostRaport(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentSchoolCount(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()

	res, err := lib.GetStudentSchoolCount(c.Params("school_id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}

func GetStudentBySmartBtwID(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	req := new(request.SmartBtwIDArray)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	res, err := lib.GetStudentProfileByArrayOfSmartbtwIDMongo(req.SmartbtwID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Student data not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    res,
		"message": "Success",
	})
}
