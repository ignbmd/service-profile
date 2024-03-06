package lib

import (
	"time"

	"smartbtw.com/services/profile/request"
)

func CreateAssessmentScreening(c *request.CreateAssessmentScreening) error {
	req := &request.AssessmentScreening{
		SmartBtwID:      c.SmartBtwID,
		PackageID:       c.PackageID,
		AssessmentCode:  c.AssessmentCode,
		Bio:             c.Bio,
		ScreeningTarget: c.ScreeningTarget,
		ProgramType:     c.ProgramType,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := StoreAssessmentScreening(req)
	if err != nil {
		return err
	}
	return nil
}
