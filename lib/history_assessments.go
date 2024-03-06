package lib

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func UpsertHistoryAssessments(c *request.CreateHistoryAssessment) error {
	opts := options.Update().SetUpsert(true)
	htpCol := db.Mongodb.Collection(fmt.Sprintf("history_assessments_%s", strings.ToLower(c.ProgramType)))
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filll := bson.M{"smartbtw_id": c.SmartBtwID, "package_id": c.PackageID}

	payload := models.HistoryAssessments{
		SmartBtwID:     c.SmartBtwID,
		PackageID:      c.PackageID,
		AssessmentCode: c.AssessmentCode,
		PackageType:    c.PackageType,
		ModuleCode:     c.ModuleCode,
		ModuleType:     c.ModuleType,
		Scores:         c.Scores,
		ScoreType:      c.ScoreType,
		Total:          c.Total,
		ExamName:       c.ExamName,
		Start:          c.Start,
		End:            c.End,
		IsLive:         c.IsLive,
		StudentName:    c.StudentName,
		StudentEmail:   c.StudentEmail,
		Program:        c.Program,
		ProgramType:    c.ProgramType,
		ProgramVersion: c.ProgramVersion,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	update := bson.M{"$set": payload}

	_, err := htpCol.UpdateOne(ctx, filll, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func SaveResultAssessments(c *request.CreateHistoryAssessment) error {

	// if db.FbDBMulti["STAGES"] == nil {
	// 	fmt.Println("firebase not initialized")
	// 	return nil
	// }
	if strings.ToLower(c.ProgramType) == "skb-cpns" {
		c.PackageType = "skb"
	}

	refName := fmt.Sprintf("/cassessments-results-%s", strings.ToLower(c.ProgramType))

	ref := db.FbDBMulti["STAGES"].NewRef(refName)

	strct := map[string]interface{}{
		"assessments": c,
	}

	err := ref.Child(fmt.Sprintf("package_%d/assessment_%s/student_%d", c.PackageID, c.AssessmentCode, c.SmartBtwID)).Update(db.Ctx, strct)
	if err != nil {
		return err
	}
	return nil
}

func StoreAssessmentScreening(c *request.AssessmentScreening) error {

	// if db.FbDBMulti["STAGES"] == nil {
	// 	fmt.Println("firebase not initialized")
	// 	return nil
	// }
	if strings.ToLower(c.ProgramType) == "skb-cpns" {
		c.ProgramType = "skb"
	}

	refName := fmt.Sprintf("/cassessments-results-%s", strings.ToLower(c.ProgramType))

	ref := db.FbDBMulti["STAGES"].NewRef(refName)

	strct := map[string]interface{}{
		"screening": c,
	}

	err := ref.Child(fmt.Sprintf("package_%d/assessment_%s/student_%d", c.PackageID, c.AssessmentCode, c.SmartBtwID)).Update(db.Ctx, strct)
	if err != nil {
		return err
	}
	return nil
}
