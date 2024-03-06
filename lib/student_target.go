package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/olivere/elastic/v7"
	"github.com/pandeptwidyaop/golog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	requests "smartbtw.com/services/profile/request"
)

func CreateStudentTargetRest(c *request.CreateStudentTarget) (*mongo.InsertOneResult, error) {
	res, err := CreateStudentTarget(c)
	if err != nil {
		return nil, err
	}
	if c.TargetType == string(models.PTK) {
		polbitType := "PUSAT"

		if c.PolbitType != "" {
			polbitType = c.PolbitType
		}
		err1 := InsertStudentTargetPtkElastic(&request.StudentTargetPtkElastic{
			SmartbtwID:          c.SmartbtwID,
			SchoolID:            c.SchoolID,
			MajorID:             c.MajorID,
			SchoolName:          c.SchoolName,
			MajorName:           c.MajorName,
			TargetScore:         c.TargetScore,
			TargetType:          c.TargetType,
			PolbitType:          polbitType,
			PolbitCompetitionID: c.PolbitCompetitionID,
			PolbitLocationID:    c.PolbitLocationID,
		})
		if err1 != nil {
			text := "Error when creating PTK data to elastic"
			js, _ := sonic.Marshal(c)
			golog.Slack.ErrorWithData(text, js, err)
		}
	} else if c.TargetType == string(models.CPNS) {
		polbitType := "PUSAT"

		if c.PolbitType != "" {
			polbitType = c.PolbitType
		}
		err1 := InsertStudentTargetPtkElastic(&request.StudentTargetPtkElastic{
			SmartbtwID:          c.SmartbtwID,
			SchoolID:            c.SchoolID,
			MajorID:             c.MajorID,
			SchoolName:          c.SchoolName,
			MajorName:           c.MajorName,
			TargetScore:         c.TargetScore,
			TargetType:          c.TargetType,
			PolbitType:          polbitType,
			PolbitCompetitionID: c.PolbitCompetitionID,
			PolbitLocationID:    c.PolbitLocationID,
		})
		if err1 != nil {
			text := "Error when creating CPNS data to elastic"
			js, _ := sonic.Marshal(c)
			golog.Slack.ErrorWithData(text, js, err)
		}
	} else {
		err1 := InsertStudentTargetPtnElastic(&request.StudentTargetPtnElastic{
			SmartbtwID:  c.SmartbtwID,
			SchoolID:    c.SchoolID,
			MajorID:     c.MajorID,
			SchoolName:  c.SchoolName,
			MajorName:   c.MajorName,
			TargetScore: c.TargetScore,
			ProgramKey:  "utbk",
			TargetType:  c.TargetType,
		})
		if err1 != nil {
			text := "Error when creating PTN data to elastic"
			js, _ := sonic.Marshal(c)
			golog.Slack.ErrorWithData(text, js, err)
		}
	}

	return res, nil
}

func CreateStudentTarget(c *request.CreateStudentTarget) (*mongo.InsertOneResult, error) {
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "target_type": c.TargetType, "deleted_at": nil}
	stdModels := models.StudentTarget{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
	}

	if err == nil {
		return nil, nil
	}
	polbitType := "PUSAT"

	if c.PolbitType != "" {
		polbitType = c.PolbitType
	}
	payload := models.StudentTarget{
		SmartbtwID:          c.SmartbtwID,
		SchoolID:            c.SchoolID,
		MajorID:             c.MajorID,
		SchoolName:          c.SchoolName,
		MajorName:           c.MajorName,
		TargetScore:         c.TargetScore,
		TargetType:          c.TargetType,
		PolbitType:          polbitType,
		PolbitCompetitionID: c.PolbitCompetitionID,
		PolbitLocationID:    c.PolbitLocationID,
		CanUpdate:           true,
		IsActive:            true,
		Position:            0,
		Type:                "PRIMARY",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}

	res, err1 := stdCol.InsertOne(ctx, payload)

	if err1 != nil {
		return nil, err1
	}

	ctxEls := context.Background()
	_, errEls := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", c.SmartbtwID)).
		Doc(map[string]interface{}{
			fmt.Sprintf("school_%s_id", strings.ToLower(c.TargetType)):             c.SchoolID,
			fmt.Sprintf("school_name_%s", strings.ToLower(c.TargetType)):           c.SchoolName,
			fmt.Sprintf("major_name_%s", strings.ToLower(c.TargetType)):            c.MajorName,
			fmt.Sprintf("major_%s_id", strings.ToLower(c.TargetType)):              c.MajorID,
			fmt.Sprintf("polbit_type_%s", strings.ToLower(c.TargetType)):           polbitType,
			fmt.Sprintf("polbit_competition_%s_id", strings.ToLower(c.TargetType)): c.PolbitCompetitionID,
			fmt.Sprintf("polbit_location_%s_id", strings.ToLower(c.TargetType)):    c.PolbitLocationID,
			fmt.Sprintf("created_at_%s", strings.ToLower(c.TargetType)):            time.Now(),
		}).
		DocAsUpsert(true).
		Do(ctxEls)

	if errEls != nil {
		js, _ := sonic.Marshal(c)
		golog.Slack.ErrorWithData("error update elastic data", js, err)
	}

	return res, nil
}

func UpdateStudentTargetByID(c *request.UpdateStudentTarget, id primitive.ObjectID) error {
	opts := options.Update().SetUpsert(true)
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "smartbtw_id": c.SmartbtwID, "target_type": c.TargetType, "can_update": true, "deleted_at": nil}
	stdTarget := models.StudentTarget{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdTarget)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}
	if err != nil {
		return fmt.Errorf("^record already updated")
	}

	payload := models.StudentTarget{
		SmartbtwID:          stdTarget.SmartbtwID,
		SchoolID:            stdTarget.SchoolID,
		MajorID:             stdTarget.MajorID,
		SchoolName:          stdTarget.SchoolName,
		MajorName:           stdTarget.MajorName,
		TargetScore:         stdTarget.TargetScore,
		TargetType:          stdTarget.TargetType,
		PolbitType:          stdTarget.PolbitType,
		PolbitCompetitionID: stdTarget.PolbitCompetitionID,
		PolbitLocationID:    stdTarget.PolbitLocationID,
		CanUpdate:           false,
		IsActive:            false,
		CreatedAt:           stdTarget.CreatedAt,
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}

	update := bson.M{"$set": payload}
	_, err1 := stdCol.UpdateByID(ctx, stdTarget.ID, update, opts)

	if err1 != nil {
		return err1
	}
	polbitType := "PUSAT"

	if c.PolbitType != "" {
		polbitType = c.PolbitType
	}
	payload1 := models.StudentTarget{
		SmartbtwID:          stdTarget.SmartbtwID,
		SchoolID:            c.SchoolID,
		MajorID:             c.MajorID,
		SchoolName:          c.SchoolName,
		MajorName:           c.MajorName,
		TargetScore:         c.TargetScore,
		TargetType:          c.TargetType,
		PolbitType:          polbitType,
		PolbitCompetitionID: c.PolbitCompetitionID,
		PolbitLocationID:    c.PolbitLocationID,
		CanUpdate:           false,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}
	_, err2 := stdCol.InsertOne(ctx, payload1)
	if err2 != nil {
		return err2
	}

	return nil
}

func UpdateStudentTargetBySmartbtwID(c *request.UpdateStudentTarget, SmartbtwID int) error {
	opts := options.Update().SetUpsert(true)
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": SmartbtwID, "target_type": c.TargetType, "can_update": true, "deleted_at": nil}
	stdTarget := models.StudentTarget{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdTarget)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}
	if err != nil {
		return fmt.Errorf("^record already updated")
	}

	payload := models.StudentTarget{
		SmartbtwID:          stdTarget.SmartbtwID,
		SchoolID:            stdTarget.SchoolID,
		MajorID:             stdTarget.MajorID,
		SchoolName:          stdTarget.SchoolName,
		MajorName:           stdTarget.MajorName,
		TargetScore:         stdTarget.TargetScore,
		TargetType:          stdTarget.TargetType,
		PolbitType:          stdTarget.PolbitType,
		PolbitCompetitionID: stdTarget.PolbitCompetitionID,
		PolbitLocationID:    stdTarget.PolbitLocationID,
		CanUpdate:           false,
		IsActive:            false,
		CreatedAt:           stdTarget.CreatedAt,
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}

	update := bson.M{"$set": payload}
	_, err1 := stdCol.UpdateByID(ctx, stdTarget.ID, update, opts)

	if err1 != nil {
		return err1
	}
	polbitType := "PUSAT"

	if c.PolbitType != "" {
		polbitType = c.PolbitType
	}
	payload1 := models.StudentTarget{
		SmartbtwID:          stdTarget.SmartbtwID,
		SchoolID:            c.SchoolID,
		MajorID:             c.MajorID,
		SchoolName:          c.SchoolName,
		MajorName:           c.MajorName,
		TargetScore:         c.TargetScore,
		TargetType:          c.TargetType,
		PolbitType:          polbitType,
		PolbitCompetitionID: c.PolbitCompetitionID,
		PolbitLocationID:    c.PolbitLocationID,
		CanUpdate:           false,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}
	_, err2 := stdCol.InsertOne(ctx, payload1)
	if err2 != nil {
		return err2
	}
	if c.TargetType == string(models.PTK) {
		err1 := UpdateStudentTargetDataPTKElastic(&request.StudentTargetPtkElastic{
			SmartbtwID:          c.SmartbtwID,
			SchoolID:            c.SchoolID,
			MajorID:             c.MajorID,
			SchoolName:          c.SchoolName,
			MajorName:           c.MajorName,
			TargetScore:         c.TargetScore,
			TargetType:          c.TargetType,
			PolbitType:          polbitType,
			PolbitCompetitionID: c.PolbitCompetitionID,
			PolbitLocationID:    c.PolbitLocationID,
		})
		if err1 != nil {
			text := "Error when updating PTK data to elastic"
			js, _ := sonic.Marshal(c)
			golog.Slack.ErrorWithData(text, js, err)
		}
	} else {
		err1 := UpdateStudentTargetDataPTNElastic(&request.StudentTargetPtnElastic{
			SmartbtwID:  c.SmartbtwID,
			SchoolID:    c.SchoolID,
			MajorID:     c.MajorID,
			SchoolName:  c.SchoolName,
			MajorName:   c.MajorName,
			TargetScore: c.TargetScore,
			ProgramKey:  "utbk",
			TargetType:  c.TargetType,
		})
		if err1 != nil {
			text := "Error when updating PTN data to elastic"
			js, _ := sonic.Marshal(c)
			golog.Slack.ErrorWithData(text, js, err)
		}
	}

	return nil
}

func UpdateBulkStudentTargetBySmartbtwID(c *request.UpdateStudentTargetRequest, smartbtwID int) error {
	opts := options.Update().SetUpsert(true)
	stdCol := db.Mongodb.Collection("student_targets")
	if strings.ToUpper(c.TargetType) == string(models.CPNS) {
		stdCol = db.Mongodb.Collection("student_target_cpns")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ownedSchool := map[uint]uint{}
	updatePositionSchool := map[uint]request.UpdateStudentTargetBody{}
	for _, k := range c.StudentTargetList {
		_, isPositionExists := updatePositionSchool[k.Position]
		if !isPositionExists {
			updatePositionSchool[k.Position] = request.UpdateStudentTargetBody{}
		} else {
			return errors.New("one or more of the body contains same position")
		}
		_, isSchoolExists := updatePositionSchool[k.Position]
		if !isSchoolExists {
			updatePositionSchool[k.Position] = request.UpdateStudentTargetBody{}
		}

		if strings.ToUpper(c.TargetType) != string(models.CPNS) {
			_, isSchoolCounterExists := ownedSchool[uint(k.SchoolID)]
			if !isSchoolCounterExists {
				ownedSchool[uint(k.SchoolID)] = 0
			}
			ownedSchool[uint(k.SchoolID)] += 1
		} else {
			_, isSchoolCounterExists := ownedSchool[uint(k.InstanceID)]
			if !isSchoolCounterExists {
				ownedSchool[uint(k.InstanceID)] = 0
			}
			ownedSchool[uint(k.InstanceID)] += 1
		}
		updatePositionSchool[k.Position] = k
	}

	for _, k := range ownedSchool {
		if k > 2 {
			return errors.New("having more than 2 same school is not allowed")
		}
	}

	newPosition := make([]uint, len(c.StudentTargetList))
	for i, p := range c.StudentTargetList {
		newPosition[i] = p.Position
	}

	if strings.ToUpper(c.TargetType) != "CPNS" {
		tr, err := GetAllStudentTarget(smartbtwID, strings.ToUpper(c.TargetType))
		if err != nil {
			return err
		}

		for _, t := range tr {
			err := DeleteStudentTarget(t.ID)
			if err != nil {
				return err
			}
		}
	} else {
		tr, err := GetAllStudentTargetCPNSForUpdateTarget(smartbtwID)
		if err != nil {
			return err
		}
		for _, t := range tr {
			err := DeleteStudentTargetCPNS(uint(smartbtwID), t.ID)
			if err != nil {
				return err
			}
		}
	}

	if strings.ToUpper(c.TargetType) != string(models.CPNS) {
		updatePrimary := false
		var polbitLocId *int
		var polbitCompId *int
		var polbitType string = "PUSAT"
		for _, k := range updatePositionSchool {
			filter := bson.M{"smartbtw_id": smartbtwID, "position": k.Position, "target_type": strings.ToUpper(c.TargetType), "is_active": true, "deleted_at": nil}
			stdTarget := models.StudentTarget{}
			err := stdCol.FindOne(ctx, filter).Decode(&stdTarget)
			if err != nil {
				payload1 := models.StudentTarget{
					SmartbtwID:          smartbtwID,
					SchoolID:            k.SchoolID,
					MajorID:             k.MajorID,
					SchoolName:          k.SchoolName,
					MajorName:           k.MajorName,
					PolbitCompetitionID: k.PolbitCompetitionID,
					PolbitLocationID:    k.PolbitLocationID,
					PolbitType:          k.PolbitType,
					TargetScore:         k.TargetScore,
					TargetType:          strings.ToUpper(c.TargetType),
					CanUpdate:           false,
					IsActive:            true,
					Position:            k.Position,
					Type:                k.Type,
					CreatedAt:           time.Now(),
					UpdatedAt:           time.Now(),
					DeletedAt:           nil,
				}
				_, err2 := stdCol.InsertOne(ctx, payload1)
				if err2 != nil {
					return err2
				}
				continue
			}

			if k.Position == stdTarget.Position {
				if k.SchoolID == stdTarget.SchoolID && k.MajorID == stdTarget.MajorID {
					if k.PolbitCompetitionID != nil {
						if stdTarget.PolbitCompetitionID != nil && *stdTarget.PolbitCompetitionID == *k.PolbitCompetitionID {
							continue
						}
					} else {
						continue
					}
				}
			}

			payload := models.StudentTarget{
				SmartbtwID:          stdTarget.SmartbtwID,
				SchoolID:            stdTarget.SchoolID,
				MajorID:             stdTarget.MajorID,
				SchoolName:          stdTarget.SchoolName,
				MajorName:           stdTarget.MajorName,
				TargetScore:         stdTarget.TargetScore,
				TargetType:          stdTarget.TargetType,
				PolbitType:          stdTarget.PolbitType,
				PolbitCompetitionID: stdTarget.PolbitCompetitionID,
				PolbitLocationID:    stdTarget.PolbitLocationID,
				CanUpdate:           false,
				IsActive:            false,
				Position:            stdTarget.Position,
				Type:                stdTarget.Type,
				CreatedAt:           stdTarget.CreatedAt,
				UpdatedAt:           time.Now(),
				DeletedAt:           nil,
			}

			update := bson.M{"$set": payload}
			_, err1 := stdCol.UpdateByID(ctx, stdTarget.ID, update, opts)

			if err1 != nil {
				return err1
			}

			payload1 := models.StudentTarget{
				SmartbtwID:          stdTarget.SmartbtwID,
				SchoolID:            k.SchoolID,
				MajorID:             k.MajorID,
				SchoolName:          k.SchoolName,
				MajorName:           k.MajorName,
				TargetScore:         k.TargetScore,
				PolbitCompetitionID: k.PolbitCompetitionID,
				PolbitLocationID:    k.PolbitLocationID,
				PolbitType:          k.PolbitType,
				TargetType:          strings.ToUpper(c.TargetType),
				CanUpdate:           false,
				IsActive:            true,
				Position:            k.Position,
				Type:                k.Type,
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
				DeletedAt:           nil,
			}
			_, err2 := stdCol.InsertOne(ctx, payload1)
			if err2 != nil {
				return err2
			}

			if k.Type == string(models.PRIMARY) {
				updatePrimary = true
				polbitLocId = k.PolbitLocationID
				polbitCompId = k.PolbitCompetitionID
				polbitType = k.PolbitType
			}
		}

		if updatePrimary {
			if strings.ToUpper(c.TargetType) == string(models.PTK) {

				err1 := UpdateStudentTargetDataPTKElastic(&request.StudentTargetPtkElastic{
					SmartbtwID:          smartbtwID,
					SchoolID:            updatePositionSchool[0].SchoolID,
					MajorID:             updatePositionSchool[0].MajorID,
					SchoolName:          updatePositionSchool[0].SchoolName,
					MajorName:           updatePositionSchool[0].MajorName,
					TargetScore:         updatePositionSchool[0].TargetScore,
					TargetType:          "PRIMARY",
					PolbitType:          polbitType,
					PolbitCompetitionID: polbitCompId,
					PolbitLocationID:    polbitLocId,
					FormationYear:       updatePositionSchool[0].FormationYear,
				})
				if err1 != nil {
					text := "Error when updating PTK data to elastic"
					js, _ := sonic.Marshal(c)
					golog.Slack.ErrorWithData(text, js, err1)
					return err1
				}
			} else {
				err1 := UpdateStudentTargetDataPTNElastic(&request.StudentTargetPtnElastic{
					SmartbtwID:    smartbtwID,
					ProgramKey:    "utbk",
					SchoolID:      updatePositionSchool[0].SchoolID,
					MajorID:       updatePositionSchool[0].MajorID,
					SchoolName:    updatePositionSchool[0].SchoolName,
					MajorName:     updatePositionSchool[0].MajorName,
					TargetScore:   updatePositionSchool[0].TargetScore,
					TargetType:    "PRIMARY",
					FormationYear: updatePositionSchool[0].FormationYear,
				})
				if err1 != nil {
					text := "Error when updating PTN data to elastic"
					js, _ := sonic.Marshal(c)
					golog.Slack.ErrorWithData(text, js, err1)
					return err1
				}
			}
		}
	} else {
		updatePrimary := false
		for _, k := range updatePositionSchool {
			filter := bson.M{"smartbtw_id": smartbtwID, "position": k.Position, "is_active": true, "deleted_at": nil}
			stdTarget := models.StudentTargetCpns{}
			err := stdCol.FindOne(ctx, filter).Decode(&stdTarget)
			if err != nil {
				payload1 := models.StudentTargetCpns{
					SmartbtwID:        smartbtwID,
					InstanceID:        k.InstanceID,
					PositionID:        k.PositionID,
					InstanceName:      k.InstanceName,
					PositionName:      k.PositionName,
					FormationType:     k.FormationType,
					FormationLocation: k.FormationLocation,
					FormationCode:     k.FormationCode,
					CompetitionID:     k.CompetitionID,
					TargetScore:       k.TargetScore,
					TargetType:        strings.ToUpper(c.TargetType),
					CanUpdate:         false,
					IsActive:          true,
					Position:          k.Position,
					Type:              k.Type,
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
					DeletedAt:         nil,
				}
				_, err2 := stdCol.InsertOne(ctx, payload1)
				if err2 != nil {
					return err2
				}
				continue
			}

			if k.Position == stdTarget.Position {
				if k.InstanceID == stdTarget.InstanceID && k.PositionID == stdTarget.PositionID {
					if stdTarget.CompetitionID == k.CompetitionID {
						continue
					}
				}
			}

			payload := models.StudentTargetCpns{
				SmartbtwID:        stdTarget.SmartbtwID,
				InstanceID:        stdTarget.InstanceID,
				PositionID:        stdTarget.PositionID,
				InstanceName:      stdTarget.InstanceName,
				PositionName:      stdTarget.PositionName,
				TargetScore:       stdTarget.TargetScore,
				TargetType:        stdTarget.TargetType,
				FormationType:     stdTarget.FormationType,
				FormationLocation: stdTarget.FormationLocation,
				FormationCode:     stdTarget.FormationCode,
				CompetitionID:     stdTarget.CompetitionID,
				CanUpdate:         false,
				IsActive:          false,
				Position:          stdTarget.Position,
				Type:              stdTarget.Type,
				CreatedAt:         stdTarget.CreatedAt,
				UpdatedAt:         time.Now(),
				DeletedAt:         nil,
			}

			update := bson.M{"$set": payload}
			_, err1 := stdCol.UpdateByID(ctx, stdTarget.ID, update, opts)

			if err1 != nil {
				return err1
			}

			payload1 := models.StudentTargetCpns{
				SmartbtwID:        stdTarget.SmartbtwID,
				InstanceID:        k.InstanceID,
				PositionID:        k.PositionID,
				InstanceName:      k.InstanceName,
				PositionName:      k.PositionName,
				TargetScore:       k.TargetScore,
				CompetitionID:     k.CompetitionID,
				FormationType:     k.FormationType,
				FormationLocation: k.FormationLocation,
				FormationCode:     k.FormationCode,
				TargetType:        strings.ToUpper(c.TargetType),
				CanUpdate:         false,
				IsActive:          true,
				Position:          k.Position,
				Type:              k.Type,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				DeletedAt:         nil,
			}
			_, err2 := stdCol.InsertOne(ctx, payload1)
			if err2 != nil {
				return err2
			}

			if k.Type == string(models.PRIMARY) {
				updatePrimary = true
			}
		}

		if updatePrimary {
			if updatePositionSchool[0].NewestFormation {
				err1 := UpdateStudentTargetDataCPNSFormationYearElastic(&request.StudentTargetCpnsElastic{
					SmartbtwID:        smartbtwID,
					InstanceID:        updatePositionSchool[0].InstanceID,
					PositionID:        updatePositionSchool[0].PositionID,
					InstanceName:      updatePositionSchool[0].InstanceName,
					PositionName:      updatePositionSchool[0].PositionName,
					TargetScore:       updatePositionSchool[0].TargetScore,
					TargetType:        "PRIMARY",
					FormationType:     updatePositionSchool[0].FormationType,
					FormationLocation: updatePositionSchool[0].FormationLocation,
					FormationCode:     updatePositionSchool[0].FormationCode,
					FormationYear:     updatePositionSchool[0].FormationYear,
					CompetitionID:     updatePositionSchool[0].CompetitionID,
				})
				if err1 != nil {
					text := "Error when updating CPNS data to elastic"
					js, _ := sonic.Marshal(c)
					golog.Slack.ErrorWithData(text, js, err1)
					return err1
				}
			} else {
				err1 := UpdateStudentTargetDataCPNSElastic(&request.StudentTargetCpnsElastic{
					SmartbtwID:        smartbtwID,
					InstanceID:        updatePositionSchool[0].InstanceID,
					PositionID:        updatePositionSchool[0].PositionID,
					InstanceName:      updatePositionSchool[0].InstanceName,
					PositionName:      updatePositionSchool[0].PositionName,
					TargetScore:       updatePositionSchool[0].TargetScore,
					TargetType:        "PRIMARY",
					FormationType:     updatePositionSchool[0].FormationType,
					FormationLocation: updatePositionSchool[0].FormationLocation,
					FormationCode:     updatePositionSchool[0].FormationCode,
					CompetitionID:     updatePositionSchool[0].CompetitionID,
				})
				if err1 != nil {
					text := "Error when updating CPNS data to elastic"
					js, _ := sonic.Marshal(c)
					golog.Slack.ErrorWithData(text, js, err1)
					return err1
				}
			}
		}
	}

	return nil
}

func containsPosition(newPosition []uint, pos uint) bool {
	for _, p := range newPosition {
		if p == pos {
			return true
		}
	}
	return false
}

func DeleteStudentTarget(id primitive.ObjectID) error {
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fil := bson.M{"_id": id, "deleted_at": nil}
	stdModel := models.StudentTarget{}
	err := stdCol.FindOne(ctx, fil).Decode(&stdModel)
	if err != nil {
		return err
	}

	var (
		tn  time.Time  = time.Now()
		tmn *time.Time = &tn
	)

	payload := models.StudentTarget{
		SmartbtwID:          stdModel.SmartbtwID,
		SchoolID:            stdModel.SchoolID,
		MajorID:             stdModel.MajorID,
		SchoolName:          stdModel.SchoolName,
		MajorName:           stdModel.MajorName,
		TargetScore:         stdModel.TargetScore,
		PolbitType:          stdModel.PolbitType,
		PolbitCompetitionID: stdModel.PolbitCompetitionID,
		PolbitLocationID:    stdModel.PolbitLocationID,
		TargetType:          stdModel.TargetType,
		CanUpdate:           stdModel.CanUpdate,
		IsActive:            false,
		CreatedAt:           stdModel.CreatedAt,
		UpdatedAt:           stdModel.UpdatedAt,
		DeletedAt:           tmn,
	}

	update := bson.M{"$set": payload}
	_, err = stdCol.UpdateByID(ctx, id, update)

	if err != nil {
		return err
	}

	return nil
}

func GetStudentTargetByID(id primitive.ObjectID) (models.StudentTarget, error) {
	var results models.StudentTarget

	collection := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// err := collection.FindOne(ctx, bson.M{"code": code}).Decode(&model)
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&results)
	if err != nil {
		return models.StudentTarget{}, fmt.Errorf("data not found")
	}

	return results, nil
}

func GetStudentTargetByCustom(smID int, tType string) (models.StudentTarget, error) {
	var results models.StudentTarget

	collection := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": smID,
		"target_type": tType,
		"is_active":   true,
		"deleted_at":  nil,
	}

	err := collection.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		return models.StudentTarget{}, fmt.Errorf("^data not found")
	}

	return results, nil
}

func GetAllStudentTargetByCustom(smID int, tType string) ([]models.StudentTarget, error) {
	var results []models.StudentTarget

	collection := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": smID,
		"target_type": tType,
		"is_active":   true,
		"deleted_at":  nil,
	}

	// err := collection.Find(ctx, filter).Decode(&results)
	// if err != nil {
	// 	return models.StudentTarget{}, fmt.Errorf("^data not found")
	// }
	sort := bson.M{"position": 1}
	opts := options.Find()
	opts.SetSort(sort)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		stdTarget := models.StudentTarget{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		results = append(results, stdTarget)
	}

	return results, nil
}

func InsertStudentTargetPtkElastic(data *request.StudentTargetPtkElastic) error {
	var idst string
	ctx := context.Background()

	idst = fmt.Sprintf("%d_PTK", data.SmartbtwID)

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentTargetPtkIndexName()).
		Id(idst).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func InsertStudentTargetPtnElastic(data *request.StudentTargetPtnElastic) error {
	var idst string
	ctx := context.Background()

	idst = fmt.Sprintf("%d_PTN", data.SmartbtwID)

	_, err := db.ElasticClient.Index().
		Index(db.GetStudentTargetPtnIndexName()).
		Id(idst).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func UpdateUserData(req *request.UpdateUserData, smId int, c context.Context) error {

	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", smId))
	script := elastic.NewScript(`
			ctx._source.name = params.name;
			ctx._source.photo = params.photo;
			`).
		Params(map[string]interface{}{
			"name":  req.Name,
			"photo": req.Photo,
		})
	_, err := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName(), db.GetStudentTargetPtnIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())
	if err != nil {
		return fmt.Errorf("data not found")

	}

	return nil
}

func GetStudentTargetElastic(smID int, tgType string, progKey string) (map[string]interface{}, error) {
	ctx := context.Background()

	var t map[string]interface{}
	var gres map[string]interface{}

	qrs := []elastic.Query{
		elastic.NewMatchQuery("smartbtw_id", smID),
		elastic.NewMatchQuery("target_type", strings.ToUpper(tgType)),
	}

	if progKey != "" {

		qrs = append(qrs, elastic.NewMatchQuery("program_key", progKey))
	}

	query := elastic.NewBoolQuery().Must(
		qrs...,
	)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentTargetPtkIndexName(), db.GetStudentTargetPtnIndexName(), db.GetStudentTargetCpnsIndexName()).
		Query(query).
		Size(fiber.StatusOK).Do(ctx)

	if err != nil {
		return map[string]interface{}{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 || recCount > 1 {
		return map[string]interface{}{}, errors.New("^data not found or got data more than one")
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = item.(map[string]interface{})
	}

	return gres, nil
}

func UpdateSchool(req *request.UpdateSchool, sId int, tgType string, c context.Context) error {

	//update db
	opts := options.Update().SetUpsert(true)
	sCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"school_id": sId, "target_type": tgType}
	cursor, err := sCol.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		stdTarget := models.StudentTarget{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		payload := models.StudentTarget{
			SmartbtwID:          stdTarget.SmartbtwID,
			SchoolID:            stdTarget.SchoolID,
			MajorID:             stdTarget.MajorID,
			SchoolName:          req.SchoolName,
			MajorName:           stdTarget.MajorName,
			TargetScore:         stdTarget.TargetScore,
			PolbitType:          stdTarget.PolbitType,
			PolbitCompetitionID: stdTarget.PolbitCompetitionID,
			PolbitLocationID:    stdTarget.PolbitLocationID,
			TargetType:          stdTarget.TargetType,
			CanUpdate:           stdTarget.CanUpdate,
			IsActive:            stdTarget.IsActive,
			CreatedAt:           stdTarget.CreatedAt,
			UpdatedAt:           time.Now(),
			DeletedAt:           stdTarget.DeletedAt,
		}
		update := bson.M{"$set": payload}
		_, err1 := sCol.UpdateByID(ctx, stdTarget.ID, update, opts)
		if err1 != nil {
			return err1
		}
	}

	//update elastic
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("school_id", sId)).Must(elastic.NewMatchQuery("target_type", tgType))
	script := elastic.NewScript(`
				ctx._source.school_name = params.school_name;
				`).
		Params(map[string]interface{}{
			"school_name": req.SchoolName,
		})
	_, err2 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName(), db.GetStudentTargetPtnIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())
	if err2 != nil {
		return err2

	}

	return nil
}

func UpdateStudyProgram(req *request.UpdateStudyProgram, mId int, tgType string, c context.Context) error {

	//update db
	opts := options.Update().SetUpsert(true)
	sCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"major_id": mId, "target_type": tgType}
	cursor, err := sCol.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		stdTarget := models.StudentTarget{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		payload := models.StudentTarget{
			SmartbtwID:          stdTarget.SmartbtwID,
			SchoolID:            stdTarget.SchoolID,
			MajorID:             stdTarget.MajorID,
			SchoolName:          stdTarget.SchoolName,
			MajorName:           req.MajorName,
			TargetScore:         stdTarget.TargetScore,
			TargetType:          stdTarget.TargetType,
			PolbitType:          stdTarget.PolbitType,
			PolbitCompetitionID: stdTarget.PolbitCompetitionID,
			PolbitLocationID:    stdTarget.PolbitLocationID,
			CanUpdate:           stdTarget.CanUpdate,
			IsActive:            stdTarget.IsActive,
			CreatedAt:           stdTarget.CreatedAt,
			UpdatedAt:           time.Now(),
			DeletedAt:           stdTarget.DeletedAt,
		}
		update := bson.M{"$set": payload}
		_, err1 := sCol.UpdateByID(ctx, stdTarget.ID, update, opts)
		if err1 != nil {
			return err1
		}
	}

	//update elastic
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("major_id", mId)).Must(elastic.NewMatchQuery("target_type", tgType))
	script := elastic.NewScript(`
					ctx._source.major_name = params.major_name;
					`).
		Params(map[string]interface{}{
			"major_name": req.MajorName,
		})
	_, err2 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName(), db.GetStudentTargetPtnIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())
	if err2 != nil {
		return err2

	}

	return nil
}

func UpdateBulkStudentPolbit(req *request.UpdatBulkPolbit, c context.Context) error {
	for _, val := range req.StudentTargets {
		err := UpdateStudentPolbit(&val, c)
		if err != nil {
			switch err.Error() {
			case "invalid target score, must be greater than ????":
				continue
			case "formation already updated":
				continue
			default:
				js, _ := sonic.Marshal(val)
				golog.Slack.ErrorWithData("delete student target", js, err)
				continue
			}
		}
	}
	return nil
}

func UpdateStudentPolbit(req *request.UpdatePolbitType, c context.Context) error {
	//update db
	if req.TargetScore < 340 {
		return errors.New("invalid target score, must be greater than ????")
	}
	ts := uint(req.TargetScore)
	isPrimary := false

	sCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"major_id": req.MajorID, "smartbtw_id": req.SmartBTWID, "is_active": true, "target_type": req.TargetType}
	cursor, err := sCol.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)
	if cursor.RemainingBatchLength() < 1 {
		return fmt.Errorf("datas not found")
	}
	for cursor.Next(ctx) {
		stdTarget := models.StudentTarget{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		isPrimary = stdTarget.Position == 0 && stdTarget.Type == "PRIMARY" && stdTarget.IsActive

		if stdTarget.PolbitCompetitionID != nil {
			return fmt.Errorf("formation already updated")
		}

		payload := bson.D{{Key: "polbit_competition_id", Value: req.PolbitCompetitionID},
			{Key: "polbit_location_id", Value: req.PolbitLocationID},
			{Key: "polbit_type", Value: req.PolbitType},
			{Key: "target_score", Value: ts},
			{Key: "updated_at", Value: time.Now()}}
		update := bson.M{"$set": payload}
		_, err1 := sCol.UpdateByID(ctx, stdTarget.ID, update)
		if err1 != nil {
			return err1
		}
	}

	if !isPrimary {
		return nil
	}
	isBinsus := false

	joinedClass, _ := GetStudentJoinedClassType(req.SmartBTWID)

	for _, k := range joinedClass {
		if strings.Contains(strings.ToLower(k), "binsus") {
			isBinsus = true
			break
		}
	}
	averages, err := GetStudentHistoryPTKElastic(req.SmartBTWID, false)
	if err != nil {
		return err
	}

	passingTotalScore := float64(0)
	passingTotalItem := 0

	twkScore := float64(0)
	tiuScore := float64(0)
	tkpScore := float64(0)
	pAtt := float64(0)
	if isBinsus {
		challengeRecord := []request.CreateHistoryPtk{}
		for _, k := range averages {
			if strings.ToLower(k.PackageType) == "challenge-uka" || strings.ToUpper(k.ModuleType) == "WITH_CODE" {
				challengeRecord = append(challengeRecord, k)
			}

		}
		for _, k := range challengeRecord {
			twkScore += k.Twk
			tiuScore += k.Tiu
			tkpScore += k.Tkp
			passingTotalScore += k.Total
		}
		if len(challengeRecord) < 11 {
			passingTotalItem = 10
		} else {
			passingTotalItem = len(challengeRecord)
		}
		// atwk := math.Round(helpers.RoundFloat((twkScore / float64(passingTotalItem)), 2))
		// atiu := math.Round(helpers.RoundFloat((tiuScore / float64(passingTotalItem)), 2))
		// atkp := math.Round(helpers.RoundFloat((tkpScore / float64(passingTotalItem)), 2))
		pAtt = math.Round(helpers.RoundFloat((passingTotalScore / float64(passingTotalItem)), 2))
		// pAtt = atwk + atiu + atkp
	} else {
		for _, k := range averages {
			if strings.ToLower(k.PackageType) != "pre-uka" {
				if (k.Tiu >= k.TiuPass) && (k.Twk >= k.TwkPass) && (k.Tkp >= k.TkpPass) {
					passingTotalItem += 1
					passingTotalScore += k.Total
				}
			}
		}

		pAtt = helpers.RoundFloat((passingTotalScore / float64(passingTotalItem)), 2)
	}

	if math.IsNaN(pAtt) {
		pAtt = 0
	}

	percATT := helpers.RoundFloat((pAtt/req.TargetScore)*100, 2)

	if percATT > 99 {
		percATT = 99
	}
	//update elastic
	_, err1 := db.ElasticClient.Update().
		Index(db.GetStudentTargetPtkIndexName()).
		Id(fmt.Sprintf("%d_PTK", req.SmartBTWID)).
		Doc(map[string]interface{}{
			"polbit_type":                              req.PolbitType,
			"polbit_competition_id":                    req.PolbitCompetitionID,
			"polbit_location_id":                       req.PolbitLocationID,
			"passing_recommendation_avg_score":         pAtt,
			"passing_recommendation_avg_percent_score": percATT,
			"target_score":                             ts,
		}).
		DocAsUpsert(true).
		Do(context.Background())
	if err1 != nil {
		return err1
	}

	return nil
}

func UpdateTargetScore(req *request.UpdateTargetScore, mId int, tgType string, c context.Context) error {

	//update db
	opts := options.Update().SetUpsert(true)
	sCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"major_id": mId, "target_type": tgType}
	cursor, err := sCol.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		stdTarget := models.StudentTarget{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		payload := models.StudentTarget{
			SmartbtwID:          stdTarget.SmartbtwID,
			SchoolID:            stdTarget.SchoolID,
			MajorID:             stdTarget.MajorID,
			SchoolName:          stdTarget.SchoolName,
			MajorName:           stdTarget.MajorName,
			TargetScore:         float64(req.TargetScore),
			TargetType:          stdTarget.TargetType,
			PolbitType:          stdTarget.PolbitType,
			PolbitCompetitionID: stdTarget.PolbitCompetitionID,
			PolbitLocationID:    stdTarget.PolbitLocationID,
			CanUpdate:           stdTarget.CanUpdate,
			IsActive:            stdTarget.IsActive,
			CreatedAt:           stdTarget.CreatedAt,
			UpdatedAt:           time.Now(),
			DeletedAt:           stdTarget.DeletedAt,
		}
		update := bson.M{"$set": payload}
		_, err1 := sCol.UpdateByID(ctx, stdTarget.ID, update, opts)
		if err1 != nil {
			return err1
		}
	}

	//update elastic
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("major_id", mId)).Must(elastic.NewMatchQuery("target_type", tgType))
	script := elastic.NewScript(`
				ctx._source.target_score = params.target_score;
				`).
		Params(map[string]interface{}{
			"target_score": req.TargetScore,
		})
	_, err1 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName(), db.GetStudentTargetPtnIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())
	if err1 != nil {
		return err1

	}
	return nil
}

func UpdateStudentTargetDataPTKElastic(data *request.StudentTargetPtkElastic) error {
	ctxEls := context.Background()
	_, err1 := db.ElasticClient.Update().
		Index(db.GetStudentTargetPtkIndexName()).
		Id(fmt.Sprintf("%d_PTK", data.SmartbtwID)).
		Doc(map[string]interface{}{
			"school_id":             data.SchoolID,
			"major_id":              data.MajorID,
			"school_name":           data.SchoolName,
			"major_name":            data.MajorName,
			"target_score":          data.TargetScore,
			"polbit_type":           data.PolbitType,
			"polbit_competition_id": data.PolbitCompetitionID,
			"polbit_location_id":    data.PolbitLocationID,
		}).
		DocAsUpsert(true).
		Do(ctxEls)
	if err1 != nil {
		return err1
	}

	_, errCreateStudentProfile := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", data.SmartbtwID)).
		Doc(map[string]interface{}{
			"school_ptk_id":             data.SchoolID,
			"major_ptk_id":              data.MajorID,
			"school_name_ptk":           data.SchoolName,
			"major_name_ptk":            data.MajorName,
			"target_score":              data.TargetScore,
			"polbit_type_ptk":           data.PolbitType,
			"polbit_competition_ptk_id": data.PolbitCompetitionID,
			"polbit_location_ptk_id":    data.PolbitLocationID,
			"created_at_ptk":            time.Now(),
		}).
		DocAsUpsert(true).
		Do(ctxEls)

	if errCreateStudentProfile != nil {
		return errCreateStudentProfile
	}

	return nil
}

func UpdateStudentTargetDataPTNElastic(data *request.StudentTargetPtnElastic) error {
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", data.SmartbtwID))
	script := elastic.NewScript(`
				ctx._source.target_score = params.target_score;
				ctx._source.school_id = params.school_id;
				ctx._source.major_id = params.major_id;
				ctx._source.school_name = params.school_name;
				ctx._source.major_name = params.major_name;
				ctx._source.program_key = params.program_key;
				`).
		Params(map[string]interface{}{
			"school_id":    data.SchoolID,
			"major_id":     data.MajorID,
			"school_name":  data.SchoolName,
			"major_name":   data.MajorName,
			"target_score": data.TargetScore,
			"program_key":  "utbk",
		})
	_, err1 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtnIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())
	if err1 != nil {
		return err1
	}

	scriptCreateStudentProfile := elastic.NewScript(`
			ctx._source.school_ptn_id = params.school_ptn_id;
			ctx._source.school_name_ptn = params.school_name_ptn;
			ctx._source.major_name_ptn = params.major_name_ptn;
			ctx._source.major_ptn_id = params.major_ptn_id;
			`).
		Params(map[string]interface{}{
			"school_ptn_id":   data.SchoolID,
			"major_ptn_id":    data.MajorID,
			"school_name_ptn": data.SchoolName,
			"major_name_ptn":  data.MajorName,
			"target_score":    data.TargetScore,
		})
	_, errCreateStudentProfile := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentProfileIndexName()).
		Query(bq).
		Script(scriptCreateStudentProfile).
		DoAsync(context.Background())
	if errCreateStudentProfile != nil {
		return errCreateStudentProfile
	}

	return nil
}

func UpdateStudentTarget(req []*request.NewUpdateStudentTarget, c context.Context) error {
	opts := options.Update().SetUpsert(true)
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": req[0].SmartbtwID, "target_type": req[0].TargetType, "can_update": true, "deleted_at": nil}
	cursor, err := stdCol.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("data not found")
	}
	defer cursor.Close(ctx)

	std := []models.StudentTarget{}
	for cursor.Next(ctx) {
		stdTarget := models.StudentTarget{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		std = append(std, stdTarget)
	}
	// isRemove := false
	// if len(std) < len(req) {
	// 	isRemove = true
	// }

	//if SmartbtwID, SchoolID, MajorID, TargetType, Position, Type != request do this
	for _, st := range std {
		payload := models.StudentTarget{
			SmartbtwID:          st.SmartbtwID,
			SchoolID:            st.SchoolID,
			MajorID:             st.MajorID,
			SchoolName:          st.SchoolName,
			MajorName:           st.MajorName,
			TargetScore:         st.TargetScore,
			TargetType:          st.TargetType,
			PolbitType:          st.PolbitType,
			PolbitCompetitionID: st.PolbitCompetitionID,
			PolbitLocationID:    st.PolbitLocationID,
			CanUpdate:           true,
			IsActive:            false,
			Position:            st.Position,
			Type:                st.Type,
			CreatedAt:           st.CreatedAt,
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
		}

		update := bson.M{"$set": payload}
		_, err1 := stdCol.UpdateByID(ctx, st.ID, update, opts)

		if err1 != nil {
			return err1
		}

	}

	for _, s := range req {
		payload1 := models.StudentTarget{
			SmartbtwID:          req[0].SmartbtwID,
			SchoolID:            s.SchoolID,
			MajorID:             s.MajorID,
			SchoolName:          s.SchoolName,
			MajorName:           s.MajorName,
			TargetScore:         s.TargetScore,
			TargetType:          s.TargetType,
			PolbitType:          s.PolbitType,
			PolbitCompetitionID: s.PolbitCompetitionID,
			PolbitLocationID:    s.PolbitLocationID,
			CanUpdate:           true,
			IsActive:            true,
			Position:            s.Position,
			Type:                s.Type,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           nil,
		}
		_, err2 := stdCol.InsertOne(ctx, payload1)
		if err2 != nil {
			return err2
		}

	}

	//update elastic
	filterPrimary := bson.M{"smartbtw_id": req[0].SmartbtwID, "type": "PRIMARY", "is_active": true}
	mPrimary := models.StudentTarget{}
	err3 := stdCol.FindOne(ctx, filterPrimary).Decode(&mPrimary)
	if err3 != nil {
		return fmt.Errorf("can not get type primary")
	}

	if mPrimary.TargetType == string(models.PTK) {
		errEls := UpdateStudentTargetDataPTKElastic(&request.StudentTargetPtkElastic{
			SmartbtwID:          req[0].SmartbtwID,
			SchoolName:          mPrimary.SchoolName,
			SchoolID:            mPrimary.SchoolID,
			MajorName:           mPrimary.MajorName,
			MajorID:             mPrimary.MajorID,
			TargetScore:         mPrimary.TargetScore,
			PolbitType:          mPrimary.PolbitType,
			PolbitCompetitionID: mPrimary.PolbitCompetitionID,
			PolbitLocationID:    mPrimary.PolbitLocationID,
		})
		if errEls != nil {
			return errEls
		}
	} else {
		errEls := UpdateStudentTargetDataPTNElastic(&request.StudentTargetPtnElastic{
			SmartbtwID:  req[0].SmartbtwID,
			SchoolName:  mPrimary.SchoolName,
			SchoolID:    mPrimary.SchoolID,
			ProgramKey:  "utbk",
			MajorName:   mPrimary.MajorName,
			MajorID:     mPrimary.MajorID,
			TargetScore: mPrimary.TargetScore,
		})
		if errEls != nil {
			return errEls
		}
	}

	return nil
}

func SyncSchoolPTKUnmatch() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("school_id", 1))
	script1 := elastic.NewScript(`
			ctx._source.school_name = params.school_name;
			`).Params(map[string]interface{}{
		"school_name": "PKN-STAN",
	})

	_, errE := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName()).
		Query(bq).
		Script(script1).
		DoAsync(ctx)
	if errE != nil {
		return errE
	}

	return nil
}

func UpdateSpecificStudyProgram(c context.Context, req *request.UpdateSpecificStudyProgram) error {
	stdCol := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filter := bson.M{"major_id": bson.M{"$in": req.MajorID}, "target_type": req.TargetType}

	update := bson.M{
		"$set": bson.M{
			"major_id":   req.NewMajorID,
			"major_name": req.MajorName,
		},
	}

	_, err := stdCol.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	ids := make([]interface{}, len(req.MajorID))
	for i, id := range req.MajorID {
		ids[i] = id
	}

	bq := elastic.NewBoolQuery().Must(elastic.NewTermsQuery("major_id", ids...))
	script := elastic.NewScript(`
				ctx._source.major_id = params.major_id;
				ctx._source.major_name = params.major_name;
				`).
		Params(map[string]interface{}{
			"major_id":   req.NewMajorID,
			"major_name": req.MajorName,
		})
	_, err1 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())
	if err1 != nil {
		return err1
	}

	bq1 := elastic.NewBoolQuery().Must(elastic.NewTermsQuery("major_ptk_id", ids...))
	script1 := elastic.NewScript(`
				ctx._source.major_ptk_id = params.major_ptk_id;
				ctx._source.major_name_ptk = params.major_name_ptk;
				`).
		Params(map[string]interface{}{
			"major_ptk_id":   req.NewMajorID,
			"major_name_ptk": req.MajorName,
		})
	_, err2 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentProfileIndexName()).
		Query(bq1).
		Script(script1).
		DoAsync(context.Background())
	if err2 != nil {
		return err2
	}
	return nil
}

func UpdateStudentTargetOne(c *request.UpdateStudentTargetOne) error {

	compData, errs := GetCompetitionFromCompMap(uint(c.PolbitCompetitionID))

	if errs != nil {
		return errs
	}

	opts := options.Update().SetUpsert(true)
	collectionStudentTargets := db.Mongodb.Collection("student_targets")
	ctx := context.Background()
	// Student Target Update

	filter := bson.M{"smartbtw_id": c.SmartbtwID, "target_type": "PTK", "is_active": true, "position": 0, "deleted_at": nil}
	stdTarget := models.StudentTarget{}
	err := collectionStudentTargets.FindOne(ctx, filter).Decode(&stdTarget)
	if err != nil {
		return err
	}

	filSt := bson.M{"smartbtw_id": c.SmartbtwID, "target_type": "PTK"}
	updateTar := bson.M{"$set": bson.M{
		"is_active":  false,
		"can_update": false,
	}}

	_, err = collectionStudentTargets.UpdateMany(ctx, filSt, updateTar, opts)
	if err != nil {
		return err
	}

	polbitType := "PUSAT"

	if compData.LocationID != nil {
		polbitType = fmt.Sprintf("%s_%s", compData.PolbitType, compData.Location.Type)
	}

	payload1 := models.StudentTarget{
		SmartbtwID:          stdTarget.SmartbtwID,
		SchoolID:            int(compData.StudyProgram.School.ID),
		MajorID:             int(compData.StudyProgramID),
		SchoolName:          compData.StudyProgram.School.Name,
		MajorName:           compData.StudyProgram.Name,
		TargetScore:         float64(compData.LowestScore),
		TargetType:          "PTK",
		PolbitType:          polbitType,
		PolbitCompetitionID: &c.PolbitCompetitionID,
		PolbitLocationID:    compData.LocationID,
		CanUpdate:           false,
		IsActive:            true,
		Position:            0,
		Type:                "PRIMARY",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}
	_, err2 := collectionStudentTargets.InsertOne(ctx, payload1)
	if err2 != nil {
		return err2
	}
	isBinsus := false

	joinedClass, _ := GetStudentJoinedClassType(stdTarget.SmartbtwID)

	for _, k := range joinedClass {
		if strings.Contains(strings.ToLower(k), "binsus") {
			isBinsus = true
			break
		}
	}
	averages, err := GetStudentHistoryPTKElastic(c.SmartbtwID, false)
	if err != nil {
		return fmt.Errorf("student ID %d data failed to fetch history ptk from elastic with error %s", c.SmartbtwID, err.Error())

	}

	passingTotalScore := float64(0)
	passingTotalItem := 0

	// for _, k := range averages {
	// 	if strings.ToLower(k.PackageType) != "pre-uka" {
	// 		if (k.Tiu >= k.TiuPass) && (k.Twk >= k.TwkPass) && (k.Tkp >= k.TkpPass) {
	// 			passingTotalItem += 1
	// 			passingTotalScore += k.Total
	// 		}
	// 	}
	// }

	// pAtt := helpers.RoundFloat((passingTotalScore / float64(passingTotalItem)), 2)

	twkScore := float64(0)
	tiuScore := float64(0)
	tkpScore := float64(0)
	pAtt := float64(0)
	if isBinsus {
		challengeRecord := []request.CreateHistoryPtk{}
		for _, k := range averages {
			if strings.ToLower(k.PackageType) == "challenge-uka" || strings.ToUpper(k.ModuleType) == "WITH_CODE" {
				challengeRecord = append(challengeRecord, k)
			}

		}
		for _, k := range challengeRecord {
			twkScore += k.Twk
			tiuScore += k.Tiu
			tkpScore += k.Tkp
			passingTotalScore += k.Total
		}
		if len(challengeRecord) < 11 {
			passingTotalItem = 10
		} else {
			passingTotalItem = len(challengeRecord)
		}
		// atwk := math.Round(helpers.RoundFloat((twkScore / float64(passingTotalItem)), 2))
		// atiu := math.Round(helpers.RoundFloat((tiuScore / float64(passingTotalItem)), 2))
		// atkp := math.Round(helpers.RoundFloat((tkpScore / float64(passingTotalItem)), 2))
		pAtt = math.Round(helpers.RoundFloat((passingTotalScore / float64(passingTotalItem)), 2))
		// pAtt = atwk + atiu + atkp
	} else {
		for _, k := range averages {
			if strings.ToLower(k.PackageType) != "pre-uka" {
				if (k.Tiu >= k.TiuPass) && (k.Twk >= k.TwkPass) && (k.Tkp >= k.TkpPass) {
					passingTotalItem += 1
					passingTotalScore += k.Total
				}
			}
		}

		pAtt = helpers.RoundFloat((passingTotalScore / float64(passingTotalItem)), 2)
	}
	if math.IsNaN(pAtt) {
		pAtt = 0
	}

	percATT := helpers.RoundFloat((pAtt/float64(compData.LowestScore))*100, 2)

	if percATT > 99 {
		percATT = 99
	}

	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", c.SmartbtwID))

	script := elastic.NewScript(`
	ctx._source.school_id = params.school_id;
	ctx._source.school_name = params.school_name;
	ctx._source.major_name = params.major_name;
	ctx._source.major_id = params.major_id;
	ctx._source.polbit_competition_id = params.polbit_competition_id;
	ctx._source.polbit_location_id = params.polbit_location_id;
	ctx._source.polbit_type = params.polbit_type;
	ctx._source.target_score = params.target_score;
	`).Params(map[string]interface{}{
		"target_score":                     compData.LowestScore,
		"school_id":                        compData.StudyProgram.School.ID,
		"school_name":                      compData.StudyProgram.School.Name,
		"major_name":                       compData.StudyProgram.Name,
		"major_id":                         compData.StudyProgramID,
		"polbit_competition_id":            c.PolbitCompetitionID,
		"polbit_location_id":               compData.LocationID,
		"polbit_type":                      polbitType,
		"passing_recommendation_avg_score": pAtt,
		"passing_recommendation_avg_percent_score": percATT,
	})

	_, err1 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())

	if err1 != nil {
		return err1
	}

	// script2 := elastic.NewScript(`
	// ctx._source.school_ptk_id = params.school_ptk_id;
	// ctx._source.school_name_ptk = params.school_name_ptk;
	// ctx._source.major_name_ptk = params.major_name_ptk;
	// ctx._source.major_ptk_id = params.major_ptk_id;
	// ctx._source.polbit_type_ptk = params.polbit_type_ptk;
	// ctx._source.target_score_ptk = params.target_score_ptk;
	// `).Params(map[string]interface{}{
	// 	"target_score_ptk": compData.LowestScore,
	// 	"school_ptk_id":    compData.StudyProgram.School.ID,
	// 	"school_name_ptk":  compData.StudyProgram.School.Name,
	// 	"major_name_ptk":   compData.StudyProgram.Name,
	// 	"major_ptk_id":     compData.StudyProgramID,
	// 	"polbit_type_ptk":  polbitType,
	// })

	// _, err4 := elastic.NewUpdateByQueryService(db.ElasticClient).
	// 	Index(db.GetStudentProfileIndexName()).
	// 	Query(bq).
	// 	Script(script2).
	// 	DoAsync(context.Background())

	// if err4 != nil {
	// 	return err4
	// }

	// script3 := elastic.NewScript(`
	// ctx._source.polbit_competition_ptk_id = params.polbit_competition_ptk_id;
	// ctx._source.polbit_location_ptk_id = params.polbit_location_ptk_id;
	// `).Params(map[string]interface{}{
	// 	"polbit_competition_ptk_id": c.PolbitCompetitionID,
	// 	"polbit_location_ptk_id":    compData.Location.ID,
	// })

	// _, err3 := elastic.NewUpdateByQueryService(db.ElasticClient).
	// 	Index(db.GetStudentProfileIndexName()).
	// 	Query(bq).
	// 	Script(script3).
	// 	DoAsync(context.Background())

	// if err3 != nil {
	// 	return err3
	// }
	msgBodys := map[string]any{
		"version": 1,
		"data": map[string]any{
			"smartbtw_id":  c.SmartbtwID,
			"account_type": "btwedutech",
		},
	}

	msgJsons, errs := sonic.Marshal(msgBodys)
	if errs == nil && db.Broker != nil {
		_ = db.Broker.Publish(
			"user.upsert-profile-elastic",
			"application/json",
			[]byte(msgJsons), // message to publish
		)
	}

	// Give reward to make them happy
	codeName := "CHANGEMAJOR_PTK"

	msgBody := models.WalletRewardStruct{
		Version: 1,
		Data: models.WalletRewardBody{
			SmartbtwID: uint(c.SmartbtwID),
			CodeName:   codeName,
		},
	}

	msgJson, err := sonic.Marshal(msgBody)
	if err != nil {
		return errors.New("error on marshaling json student reward body " + err.Error())
	}
	if db.Broker == nil {
		return nil
	}
	if err = db.Broker.Publish(
		"history-reward.level.created",
		"application/json",
		[]byte(msgJson), // message to publish
	); err != nil {
		return errors.New("error on publishing mq for reward " + err.Error())
	}

	return nil
}

func GetStudentSchoolData(scID uint, ty string) (*requests.StudentSchoolData, error) {
	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	url := fmt.Sprintf("/school/%d", scID)
	if ty != "PTK" {
		url = fmt.Sprintf("/ptn-school/%d", scID)
	}
	request, err = http.NewRequest("GET", conn+url, nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map " + err.Error())
	}

	st := requests.StudentSchoolDataResponse{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map " + errs.Error())
	}

	return &st.Data, nil
}

func FetchMajorCompetitionData(program string, majorID uint, locationID uint, gender string, polbit string) (*requests.MajorCompData, error) {
	var tagGender string
	if majorID == 348 || majorID == 349 {
		tagGender = gender
	} else {
		tagGender = ""
	}

	locID := 0
	if locationID != 0 {
		locID = int(locationID)
	}

	var isAfirm bool
	if polbit != "" && strings.Contains(polbit, "_AFIRMASI") {
		isAfirm = true
	} else {
		isAfirm = false
	}
	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}

	var url string
	if program == "ptk" {
		if locationID != 0 && strings.Contains(polbit, "DAERAH") {
			url = conn + fmt.Sprintf("/competition/by-study-program/%d/location/%d?tags=%s&is_afirm=%t",
				majorID,
				locationID,
				tagGender,
				isAfirm)
		} else {
			url = conn + fmt.Sprintf("/competition/by-study-program/%d/%d?tags=%s",
				majorID,
				locID,
				tagGender)
		}
	} else {
		url = conn + fmt.Sprintf("/passing-grade/by-study-program/%d", majorID)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if strings.ToLower(program) == "ptk" {

		var responseContent *requests.ResponseContent
		err = json.NewDecoder(resp.Body).Decode(&responseContent)
		if err != nil {
			return nil, err
		}

		yearIndex := 0
		quotaYearIndex := 0

		for i, data := range responseContent.Data {
			if data.Quota > 0 {
				quotaYearIndex = i
				break
			}
		}

		for i, data := range responseContent.Data {
			if data.Registered > 0 {
				yearIndex = i
				break
			}
		}

		var competitionType *string
		if majorID == 348 || majorID == 349 {
			if quotaYearIndex < len(responseContent.Data) && responseContent.Data[quotaYearIndex].ID > 4266 {
				competitionType = &responseContent.Data[quotaYearIndex].CompetitionType
			}
		}
		if len(responseContent.Data) == 0 {
			return nil, nil
		}
		majorCompData := &requests.MajorCompData{
			MajorQuota:       responseContent.Data[quotaYearIndex].Quota,
			MajorRegistered:  responseContent.Data[yearIndex].Registered,
			MajorYear:        responseContent.Data[yearIndex].Year,
			MajorQuotaYear:   responseContent.Data[quotaYearIndex].Year,
			MajorQuotaChance: responseContent.Data[yearIndex].Quota,
			CompetitionType:  competitionType,
		}
		return majorCompData, nil
	} else {
		var responseContent *requests.ResponseContentPTN
		err = json.NewDecoder(resp.Body).Decode(&responseContent)
		if err != nil {
			return nil, err
		}
		if reflect.DeepEqual(responseContent.Data, requests.DataPTN{}) {
			return nil, nil
		}
		majorCompData := &requests.MajorCompData{
			MajorQuota:       responseContent.Data.Quota,
			MajorRegistered:  responseContent.Data.Registered,
			MajorYear:        responseContent.Data.Year,
			MajorQuotaYear:   responseContent.Data.Year,
			MajorQuotaChance: responseContent.Data.Quota,
			CompetitionType:  nil,
		}
		return majorCompData, nil
	}
}

func GetMajorChances(majorCompData *requests.MajorCompData) string {
	var chance string
	defer func() {
		if r := recover(); r != nil {
			chance = ""
		}
	}()
	scoreChances := int(math.Max(float64(majorCompData.MajorRegistered/majorCompData.MajorQuotaChance), 1))
	if majorCompData.MajorQuotaChance > 0 {
		chance = "1:" + strconv.Itoa(scoreChances)
	}
	return chance
}

func FetchStudentStagesProfileData(smartbtw_id uint, typ string) (*request.StagesStudentLevel, error) {
	conn := os.Getenv("SERVICE_STAGES_HOST")
	bd := map[string]any{
		"smartbtw_id": smartbtw_id,
		"type":        typ,
	}
	ns, _ := sonic.Marshal(bd)
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/student-level", conn),
		bytes.NewBuffer(ns))
	request.Header.Add("Content-Type", "application/json")

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map " + err.Error())
	}

	st := requests.StudentStage{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map " + errs.Error())
	}

	return &st.Data, nil

}

func GetSchoolCompetitionList(schOrgnID string, req request.GetCompetitonList) (*request.CompetitionList, error) {
	compLis := []request.Competitions{}

	res, totalRes, err := GetStudentsBySchoolOriginIDElastic(schOrgnID, req.Search, req.Page, req.PerPage)
	if err != nil {
		if err.Error() == "student data not found" {
			return &requests.CompetitionList{}, nil
		}
		return nil, err
	}

	totalPTK, totalPTN, totalPTKandPTN, _, err := CountTargetCompetition(schOrgnID)
	if err != nil {
		return nil, err
	}

	for _, r := range res {

		var ptk *request.SchoolCompetition
		var ptn *request.SchoolCompetition
		totalScorePTN := float64(0)
		totalScorePTK := float64(0)
		ptkLvl := 0
		ptnLvl := 0
		avgPTK := float64(0)
		avgPTN := float64(0)
		totalUKAPTN := 0
		totalUKAPTK := 0
		if r.SchoolPTKID != 0 {
			ptkLvlRes, err := FetchStudentStagesProfileData(uint(r.SmartbtwID), "PTK")
			if err != nil {
				return nil, err
			}
			ptkLvl = ptkLvlRes.Level
			ptkHis, err := GetStudentHistoryPTKElasticSpecific(r.SmartbtwID, "PREMIUM_TRYOUT")
			if err != nil {
				return nil, err
			}
			if ptkHis != nil {
				for _, a := range ptkHis {
					totalScorePTK += a.Total
				}
				totalUKAPTK = len(ptkHis)
				if len(ptkHis) != 0 {
					avgPTK = helpers.RoundFloat((totalScorePTK / float64(len(ptkHis))), 1)
				} else {
					avgPTK = 0.0
				}
			}

			ptk = &request.SchoolCompetition{
				SchoolID:   uint(r.SchoolPTKID),
				SchoolName: r.SchoolNamePTK,
				MajorID:    r.MajorPTKID,
				MajorName:  r.MajorNamePTK,
				UKALevel:   ptkLvl,
				Score:      avgPTK,
				TotalUKA:   totalUKAPTK,
			}
		}

		if r.SchoolPTNID != 0 {
			ptnLvlRes, err := FetchStudentStagesProfileData(uint(r.SmartbtwID), "PTN")
			if err != nil {
				return nil, err
			}
			ptnLvl = ptnLvlRes.Level
			ptnHis, err := GetStudentHistoryPTNElasticSpecific(r.SmartbtwID, "PREMIUM_TRYOUT", "utbk")
			if err != nil {
				return nil, err
			}
			if ptnHis != nil {
				for _, a := range ptnHis {
					totalScorePTN += a.Total
				}
				totalUKAPTN = len(ptnHis)
				if len(ptnHis) != 0 {
					avgPTN = helpers.RoundFloat((totalScorePTN / float64(len(ptnHis))), 1)
				} else {
					avgPTN = 0.0
				}
			}

			ptn = &request.SchoolCompetition{
				SchoolID:   uint(r.SchoolPTNID),
				SchoolName: r.SchoolNamePTN,
				MajorID:    r.MajorPTNID,
				MajorName:  r.MajorNamePTN,
				UKALevel:   ptnLvl,
				Score:      avgPTN,
				TotalUKA:   totalUKAPTN,
			}
		}

		comp := request.Competitions{
			SmartbtwID:     uint(r.SmartbtwID),
			Name:           r.Name,
			PTKCompetition: ptk,
			PTNCompetition: ptn,
		}

		compLis = append(compLis, comp)
	}

	result := request.CompetitionList{
		TotalTargetPTN:       int(totalPTN),
		TotalTargetPTK:       int(totalPTK),
		TotalTargetPTKAndPTN: int(totalPTKandPTN),
		Total:                int(totalRes),
		Competitions:         compLis,
	}

	return &result, nil
}

func CountTargetCompetition(schOrgnID string) (int, int, int, int, error) {
	students, total, err := GetStudentsBySchoolOriginIDElastic(schOrgnID, nil, nil, nil)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	ptk := 0
	ptn := 0
	ptkAndPtn := 0

	for _, student := range students {
		if student.SchoolPTKID != 0 {
			ptk++
		}
		if student.SchoolPTNID != 0 {
			ptn++
		}
		if student.SchoolPTKID != 0 && student.SchoolPTNID != 0 {
			ptkAndPtn++
		}
	}
	return ptk, ptn, ptkAndPtn, int(total), nil
}

func CountUKACodeBySchoolOriginID(schOrgID string) (int, error) {
	codes, err := GetTryoutCodeBySchoolID(schOrgID)
	if err != nil {
		return 0, err
	}

	uniqueStudentIDs := make(map[uint]struct{})

	for _, s := range codes {
		if strings.ToLower(s.Program) == "skd" {
			ptkHis, err := GetHistoryPTKByPackageID(s.PackageID)
			if err != nil {
				return 0, err
			}
			for _, h := range ptkHis {
				uniqueStudentIDs[uint(h.SmartBtwID)] = struct{}{}
			}
		}
		if strings.ToLower(s.Program) == "utbk" {
			ptnHis, err := GetHistoryPTNByPackageID(s.PackageID)
			if err != nil {
				return 0, err
			}
			for _, h := range ptnHis {
				uniqueStudentIDs[uint(h.SmartBtwID)] = struct{}{}
			}
		}
	}
	return len(uniqueStudentIDs), nil
}

func GetAllStudentTarget(smID int, tType string) ([]models.StudentTarget, error) {
	var results []models.StudentTarget

	collection := db.Mongodb.Collection("student_targets")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"smartbtw_id": smID,
		"target_type": tType,
		"deleted_at":  nil,
	}

	// err := collection.Find(ctx, filter).Decode(&results)
	// if err != nil {
	// 	return models.StudentTarget{}, fmt.Errorf("^data not found")
	// }
	sort := bson.M{"position": 1}
	opts := options.Find()
	opts.SetSort(sort)
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		stdTarget := models.StudentTarget{}
		if err = cursor.Decode(&stdTarget); err != nil {
			log.Fatal(err)
		}
		results = append(results, stdTarget)
	}

	return results, nil
}
