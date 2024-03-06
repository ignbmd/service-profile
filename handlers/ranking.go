package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/wghttp"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/request"
)

func GetStudentSchoolRanking(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	query := new(request.GetStudentSchoolRanking)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse query",
			"error":   err,
		})
	}

	taskId, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "id is not valid",
		})
	}

	if query.Page < 1 || query.Limit < 1 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "limit and page should more than 0",
		})
	}

	//TODO: add PTN and CPNS handlers
	prog := c.Params("program")
	if prog == "ptk" {
		res, err := lib.FetchPTKRankingSchoolPurposes(uint(taskId), query.SchoolID, query.Limit, query.Page, query.SearchKeyword)
		if err != nil {
			golog.Slack.ErrorWithData("error get ptk student school ranking", nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error get student school ranking",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	} else if prog == "cpns" {
		res, err := lib.FetchCPNSRankingSchoolPurposes(uint(taskId), query.SchoolID, query.Limit, query.Page, query.SearchKeyword)
		if err != nil {
			golog.Slack.ErrorWithData("error get ptk student school ranking", nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error get student school ranking",
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	} else {
		res, err := lib.FetchPTNRankingSchoolPurposes(uint(taskId), query.SchoolID, query.Limit, query.Page, query.SearchKeyword)
		if err != nil {
			golog.Slack.ErrorWithData("error get ptn student school ranking", nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error get student school ranking",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    res,
			"message": "Success",
		})
	}
}

func GetStudentSchoolRankingWithInformation(c *fiber.Ctx) error {
	wghttp.HttpWG.Add(1)
	defer wghttp.HttpWG.Done()
	query := new(request.GetStudentSchoolRanking)
	payload := new(mockstruct.FetchRankingUkaCodeBody)
	payloadPTN := new(mockstruct.FetchRankingPTNUkaCodeBody)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "failed to parse query",
			"error":   err,
		})
	}

	taskId, err := strconv.Atoi(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "id is not valid",
		})
	}

	if query.Page < 1 || query.Limit < 1 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "limit and page should more than 0",
		})
	}

	//TODO: add PTN and CPNS handlers
	prog := c.Params("program")
	if prog != "ptn" {
		res, err := lib.FetchPTKRankingSchoolPurposes(uint(taskId), query.SchoolID, query.Limit, query.Page, query.SearchKeyword)
		if err != nil {
			golog.Slack.ErrorWithData("error get ptk student school ranking", nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error get student school ranking",
				"error":   err.Error(),
			})
		}

		totalPTK, totalPTN, totalPTKandPTN, total, err := lib.CountTargetCompetition(query.SchoolID)

		if err != nil {
			golog.Slack.ErrorWithData("error count ptk student school ranking", nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error count student school ranking",
				"error":   err.Error(),
			})
		} else {
			payload = &mockstruct.FetchRankingUkaCodeBody{
				FetchRankingBase: res.FetchRankingBase,
				RankingData:      res.RankingData,
				RankingDataInformation: mockstruct.InformationRankingData{
					TargetPtn:       totalPTN,
					TargetPtk:       totalPTK,
					TargetPTNandPTK: totalPTKandPTN,
					Total:           total,
				},
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    payload,
			"message": "Success",
		})
	} else {
		res, err := lib.FetchPTNRankingSchoolPurposes(uint(taskId), query.SchoolID, query.Limit, query.Page, query.SearchKeyword)
		if err != nil {
			golog.Slack.ErrorWithData("error get ptn student school ranking", nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error get student school ranking",
				"error":   err.Error(),
			})
		}

		totalPTK, totalPTN, totalPTKandPTN, total, err := lib.CountTargetCompetition(query.SchoolID)

		if err != nil {
			golog.Slack.ErrorWithData("error count ptk student school ranking", nil, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "error count student school ranking",
				"error":   err.Error(),
			})
		} else {
			payloadPTN = &mockstruct.FetchRankingPTNUkaCodeBody{
				FetchRankingBase: res.FetchRankingBase,
				RankingData:      res.RankingData,
				RankingDataInformation: mockstruct.InformationRankingData{
					TargetPtn:       totalPTN,
					TargetPtk:       totalPTK,
					TargetPTNandPTK: totalPTKandPTN,
					Total:           total,
				},
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    payloadPTN,
			"message": "Success",
		})
	}
}
