package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/lib/aggregates"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	requests "smartbtw.com/services/profile/request"
)

func GetStudentBySmartBTWID(smartbtw_id int) ([]models.Student, error) {
	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []models.Student
	pipel := aggregates.GetStudentWithParents(smartbtw_id)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetStudentProfileByArrayOfSmartbtwIDMongo(smartbtw_id []uint) ([]models.Student, error) {
	var results []models.Student
	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, val := range smartbtw_id {
		filter := bson.M{"smartbtw_id": val, "deleted_at": nil}
		var hlModel models.Student

		err := collection.FindOne(ctx, filter).Decode(&hlModel)
		if err != nil {
			continue
		}
		results = append(results, hlModel)
	}

	return results, nil

}

func GetStudentOnlyBySmartBTWID(smartbtw_id int) ([]models.Student, error) {
	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []models.Student
	pipel := aggregates.GetStudentOnly(smartbtw_id)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func UpsertStudent(c *request.BodyMessageCreateUpdateStudent) (*mongo.UpdateResult, error) {
	opts := options.Update().SetUpsert(true)
	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payload := models.Student{
		SmartbtwID:         c.ID,
		Name:               c.Name,
		Email:              c.Email,
		Gender:             c.Gender,
		BirthDateLocation:  c.BirthDateLocation,
		Phone:              c.Phone,
		SchoolOrigin:       c.SchoolOrigin,
		SchoolOriginID:     c.SchoolOriginID,
		Intention:          c.Intention,
		LastEd:             c.LastEd,
		Major:              c.Major,
		Profession:         c.Profession,
		Address:            c.Address,
		ProvinceId:         c.ProvinceId,
		RegionId:           c.RegionId,
		ParentName:         c.ParentName,
		ParentNumber:       c.ParentNumber,
		Interest:           c.Interest,
		Photo:              c.Photo,
		UserTryoutId:       c.UserTryoutId,
		Status:             c.Status,
		IsPhoneVerified:    c.IsPhoneVerified,
		IsEmailVerified:    c.IsEmailVerified,
		IsDataComplete:     c.IsDataComplete,
		BranchCode:         c.BranchCode,
		AffiliateCode:      c.AffiliateCode,
		AdditionalInfo:     c.AdditionalInfo,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
		DeletedAt:          c.DeletedAt,
		DomicileProvinceId: c.DomicileProvinceId,
		DomicileRegionId:   c.DomicileRegionId,
		OriginUniversity:   c.OriginUniversity,
		BirthMotherName:    c.BirthMotherName,
		BirthPlace:         c.BirthPlace,
		NIK:                c.NIK,
	}

	filter := bson.M{"smartbtw_id": c.ID}
	update := bson.M{"$set": payload}

	result, err := collection.UpdateOne(ctx, filter, update, opts)

	if c.AccountType != "smartbtw" {
		c.AccountType = "btwedutech"
	}

	msgBody := map[string]any{
		"version": 1,
		"data": map[string]any{
			"smartbtw_id":          c.ID,
			"account_type":         c.AccountType,
			"student_profile_data": payload,
		},
	}

	msgJson, errs := sonic.Marshal(msgBody)
	if errs == nil && db.Broker != nil {
		_ = db.Broker.Publish(
			"user.upsert-profile-elastic",
			"application/json",
			[]byte(msgJson), // message to publish
		)
	}
	return result, err
}

func DeleteStudent(c *request.DeleteStudentBodyMessage) (*mongo.UpdateResult, error) {
	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": c.Data.ID}
	update := bson.M{"$set": bson.M{"deleted_at": c.Data.DeletedAt, "photo": nil}}
	result, err := collection.UpdateOne(ctx, filter, update)
	return result, err
}

func GetStudentCompletedModulesBySmartBTWID(smartbtwId int, targetType string) ([]bson.M, error) {
	collection := db.Mongodb.Collection(fmt.Sprintf("history_%s", strings.ToLower(targetType)))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []bson.M
	pipel := aggregates.GetStudentCompletedModules(smartbtwId)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func InsertStudentProfileElastic(data *request.StudentProfileElastic) error {
	ctx := context.Background()
	_, err := db.ElasticClient.Index().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", data.SmartbtwID)).
		BodyJson(data).
		Do(ctx)

	if err != nil {
		return err
	}

	return nil
}

func UpsyncStudentProfileElastic(data *request.StudentProfileUpsertElastic, retryCount int) error {
	ctx := context.Background()
	_, err := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", data.SmartbtwID)).
		Doc(data).
		DocAsUpsert(true).
		Do(ctx)
	if err != nil {
		if retryCount >= 5 {
			return err
		}
		if elastic.IsConflict(err) {
			time.Sleep(500 * time.Millisecond)
			return UpsyncStudentProfileElastic(data, retryCount+1)
		} else {
			return err
		}
	}

	return nil
}

func UpsyncCompMapStudentProfileElastic(data *request.CreateStudentCompMapProfileElastic, retryCount int) error {
	ctx := context.Background()
	_, err := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", data.SmartbtwID)).
		Doc(data).
		DocAsUpsert(true).
		Do(ctx)

	if err != nil {
		if retryCount >= 5 {
			return err
		}
		if elastic.IsConflict(err) {
			time.Sleep(500 * time.Millisecond)
			return UpsyncCompMapStudentProfileElastic(data, retryCount+1)
		} else {
			return err
		}
	}

	return nil
}

func GetStudentProfilePTNElastic(smID int) (request.StudentProfilePtnElastic, error) {
	ctx := context.Background()

	var t request.StudentProfilePtnElastic
	var gres request.StudentProfilePtnElastic

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smID),
	)
	res, err := db.ElasticClient.Search().
		Index(db.GetStudentTargetPtnIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return request.StudentProfilePtnElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.StudentProfilePtnElastic{}, fmt.Errorf("student data not found")
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = item.(request.StudentProfilePtnElastic)
	}

	return gres, nil
}

func GetStudentProfilePTKElastic(smID int) (request.StudentProfilePtkElastic, error) {
	ctx := context.Background()

	var t request.StudentProfilePtkElastic
	var gres request.StudentProfilePtkElastic

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smID),
	)
	res, err := db.ElasticClient.Search().
		Index(db.GetStudentTargetPtkIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return request.StudentProfilePtkElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.StudentProfilePtkElastic{}, fmt.Errorf("student data not found")
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = item.(request.StudentProfilePtkElastic)
	}

	return gres, nil
}

func GetStudentProfileElastic(smID int) (request.StudentProfileElastic, error) {
	ctx := context.Background()

	var t request.StudentProfileElastic
	var gres request.StudentProfileElastic

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smID),
	)
	res, err := db.ElasticClient.Search().
		Index(db.GetStudentProfileIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return request.StudentProfileElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.StudentProfileElastic{}, fmt.Errorf("student data not found")
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = item.(request.StudentProfileElastic)
	}

	return gres, nil
}

func GetStudentProfileUKAElastic(smID int) (request.StudentProfileElastic, error) {
	ctx := context.Background()

	var t request.StudentProfileElastic
	var gres request.StudentProfileElastic

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smID),
	)
	res, err := db.ElasticClient.Search().
		Index(db.GetStudentProfileIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return request.StudentProfileElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.StudentProfileElastic{}, fmt.Errorf("student data not found")
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = item.(request.StudentProfileElastic)
	}

	return gres, nil
}

func GetStudentProfileElasticByEmail(email string) (request.StudentProfileElastic, error) {
	ctx := context.Background()

	var t request.StudentProfileElastic
	var gres request.StudentProfileElastic

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("email.keyword", email),
		elastic.NewMatchQuery("account_type.keyword", "btwedutech"),
	)
	res, err := db.ElasticClient.Search().
		Index(db.GetStudentProfileIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return request.StudentProfileElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.StudentProfileElastic{}, fmt.Errorf("student data not found")
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = item.(request.StudentProfileElastic)
	}

	return gres, nil
}

func GetStudentBrachByEmails(emails string) ([]bson.M, error) {
	email := emails
	allEmail := strings.Split(email, ",")

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []bson.M
	pipel := aggregates.GetStudentBranchByEmail(allEmail)
	opts := []*options.AggregateOptions{options.Aggregate().SetAllowDiskUse(true)}
	cursor, err := collection.Aggregate(ctx, pipel, opts...)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func CacheStudentProfileToElastic(req *mockstruct.StudentElasticCacheRequest) error {
	var (
		sPtkID                uint
		sPtkName              string
		sPtkC                 time.Time
		mPtkID                uint
		mPtkName              string
		polbitPtkCompetencyID uint
		polbitPtkLocationID   uint
		polbitPtkType         string

		sPtnID   uint
		sPtnName string
		sPtnC    time.Time
		mPtnID   uint
		mPtnName string

		iCpnsID       uint
		iCpnsName     string
		pCpnsID       uint
		pCpnsName     string
		fCpnsType     string
		cCpnsID       uint
		fCpnsLocation string
		fCpnsCode     string
		cCpns         time.Time

		gender string
	)

	if req.StudentProfileData == nil {
		std, err := GetStudentOnlyBySmartBTWID(req.SmartBTWID)
		if err != nil {
			return err
		}

		if len(std) == 0 {
			return nil
		}

		req.StudentProfileData = &std[0]
	}

	//getTargetPTK
	tPtk, err := GetStudentTargetByCustom(req.SmartBTWID, "PTK")
	if err != nil {
		sPtkID = 0
		sPtkName = ""
		sPtkC = time.Now()
		mPtkID = 0
		mPtkName = ""
		polbitPtkCompetencyID = 0
		polbitPtkLocationID = 0
		polbitPtkType = ""
	} else {
		sPtkID = uint(tPtk.SchoolID)
		sPtkName = tPtk.SchoolName
		sPtkC = tPtk.CreatedAt
		mPtkID = uint(tPtk.MajorID)
		mPtkName = tPtk.MajorName
		polbitPtkLocationID = 0
		polbitPtkCompetencyID = 0
		if tPtk.PolbitCompetitionID != nil {
			polbitPtkCompetencyID = uint(*tPtk.PolbitCompetitionID)
		}
		if tPtk.PolbitLocationID != nil {
			polbitPtkLocationID = uint(*tPtk.PolbitLocationID)
		}
		polbitPtkType = tPtk.PolbitType
	}

	//getTargetPTN
	tPtn, err := GetStudentTargetByCustom(req.SmartBTWID, "PTN")
	if err != nil {
		sPtnID = 0
		sPtnName = ""
		sPtnC = time.Now()
		mPtnID = 0
		mPtnName = ""
	} else {
		sPtnID = uint(tPtn.SchoolID)
		sPtnName = tPtn.SchoolName
		sPtnC = tPtn.CreatedAt
		mPtnID = uint(tPtn.MajorID)
		mPtnName = tPtn.MajorName
	}

	if req.StudentProfileData.Gender == 1 {
		gender = "L"
	} else {
		gender = "P"
	}

	phoneNumber := "0"
	if req.StudentProfileData.Phone != nil {
		phoneNumber = *req.StudentProfileData.Phone
	}
	bc := "PT0000"
	bn := "Bimbel BTW (Kantor Pusat)"
	if req.StudentProfileData.BranchCode != nil {
		if *req.StudentProfileData.BranchCode != "PT0000" {
			bc = *req.StudentProfileData.BranchCode
			bnn, err := GetBranchByBranchCode(bc)
			if err == nil {
				bn = bnn.BranchName
			}
		}
	}

	//getTargetCPNS
	tCpns, err := GetStudentTargetCPNS(req.SmartBTWID)
	if err != nil {
		iCpnsID = 0
		iCpnsName = ""
		pCpnsID = 0
		pCpnsName = ""
		fCpnsType = ""
		cCpnsID = 0
		fCpnsLocation = ""
		fCpnsCode = ""
		cCpns = time.Now()

	} else {
		iCpnsID = uint(tCpns.InstanceID)
		iCpnsName = tCpns.InstanceName
		pCpnsID = uint(tCpns.PositionID)
		pCpnsName = tCpns.PositionName
		fCpnsType = tCpns.FormationType
		cCpnsID = uint(tCpns.CompetitionID)
		fCpnsLocation = tCpns.FormationLocation
		fCpnsCode = tCpns.FormationCode
		cCpns = tCpns.CreatedAt
	}

	re := request.StudentProfileUpsertElastic{
		SmartbtwID:             req.SmartBTWID,
		Name:                   req.StudentProfileData.Name,
		Email:                  req.StudentProfileData.Email,
		Photo:                  req.StudentProfileData.Photo,
		Interest:               req.StudentProfileData.Interest,
		OriginUniversity:       req.StudentProfileData.OriginUniversity,
		Phone:                  phoneNumber,
		Gender:                 gender,
		SchoolPTKID:            sPtkID,
		SchoolNamePTK:          sPtkName,
		MajorNamePTK:           mPtkName,
		MajorPTKID:             mPtkID,
		CreatedAtPTK:           sPtkC,
		SchoolPTNID:            sPtnID,
		SchoolNamePTN:          sPtnName,
		MajorPTNID:             mPtnID,
		MajorNamePTN:           mPtnName,
		CreatedAtPTN:           sPtnC,
		CreatedAt:              req.StudentProfileData.CreatedAt,
		AccountType:            req.AccountType,
		PolbitTypePTK:          polbitPtkType,
		PolbitLocationPTKID:    polbitPtkLocationID,
		PolbitCompetitionPTKID: polbitPtkCompetencyID,
		InstanceCPNSID:         iCpnsID,
		InstanceCPNSName:       iCpnsName,
		PositionCPNSID:         pCpnsID,
		PositionCPNSName:       pCpnsName,
		FormationCPNSType:      fCpnsType,
		CompetitionCPNSID:      cCpnsID,
		FormationCPNSLocation:  fCpnsLocation,
		FormationCPNSCode:      fCpnsCode,
		CreatedAtCPNS:          cCpns,
		BranchCode:             bc,
		BranchName:             bn,
	}

	if req.StudentProfileData.SchoolOriginID != nil {
		re.LastEdID = *req.StudentProfileData.SchoolOriginID

		highschoolData, _ := GetHighschoolStudent(re.LastEdID)
		if highschoolData != nil {
			re.LastEdName = highschoolData.Name
			re.LastEdType = highschoolData.Type
			re.LastEdRegionID = highschoolData.LocationID
			schoolRegionData, _ := GetLocationByID(highschoolData.LocationID)
			if schoolRegionData != nil {
				re.LastEdRegion = schoolRegionData.Name
			}
		}
	}

	if req.StudentProfileData.ProvinceId != nil {
		re.ProvinceID = uint(*req.StudentProfileData.ProvinceId)
		locationData, _ := GetLocationByID(re.ProvinceID)
		if locationData != nil {
			re.Province = locationData.Name
		}
	}

	if req.StudentProfileData.RegionId != nil {
		re.RegionID = uint(*req.StudentProfileData.RegionId)
		locationData, _ := GetLocationByID(re.RegionID)
		if locationData != nil {
			re.Region = locationData.Name
		}
	}

	if req.StudentProfileData.DomicileProvinceId != nil {
		re.DomicileProvinceID = uint(*req.StudentProfileData.DomicileProvinceId)
		locationData, _ := GetLocationByID(re.DomicileProvinceID)
		if locationData != nil {
			re.DomicileProvince = locationData.Name
		}
	}

	if req.StudentProfileData.DomicileRegionId != nil {
		re.DomicileRegionID = uint(*req.StudentProfileData.DomicileRegionId)
		locationData, _ := GetLocationByID(re.DomicileRegionID)
		if locationData != nil {
			re.DomicileRegion = locationData.Name
		}
	}

	// if req.StudentProfileData.PolbitTypePTK != nil {
	// 	re.PolbitTypePTK = *req.StudentProfileData.PolbitTypePTK
	// } else {
	// 	re.PolbitTypePTK = "PUSAT"
	// }

	err = UpsyncStudentProfileElastic(&re, 0)
	if err != nil {
		return err
	}
	return nil
}

func CacheStudentCompMapProfileToElastic(req *mockstruct.StudentElasticCacheRequest) error {

	re := request.CreateStudentCompMapProfileElastic{
		SmartbtwID: req.SmartBTWID,
	}

	compMapStudentData, _ := GetStudentProfileFromCompMap(uint(req.SmartBTWID))
	if compMapStudentData != nil {
		if compMapStudentData.BirthDate != nil {
			re.BirthDate = *compMapStudentData.BirthDate
		}

		re.Height = compMapStudentData.Height
		re.Weight = compMapStudentData.Weight

		if compMapStudentData.EyeColorBlind != nil {
			re.EyeColorBlind = *compMapStudentData.EyeColorBlind
		}

		re.LastEdMajor = compMapStudentData.LastEd.Name
		re.LastEdMajorID = compMapStudentData.LastEd.ID

		re.LastUniversityEdMajor = compMapStudentData.LastUniversityEd.Name
		re.LastUniversityEdMajorID = compMapStudentData.LastUniversityEd.ID
	}

	err := UpsyncCompMapStudentProfileElastic(&re, 0)
	if err != nil {
		return err
	}
	return nil
}

func GetStudentProfileCPNSElastic(smID int) (request.StudentProfileCPNSElastic, error) {
	ctx := context.Background()

	var t request.StudentProfileCPNSElastic
	var gres request.StudentProfileCPNSElastic

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smID),
	)
	res, err := db.ElasticClient.Search().
		Index(db.GetStudentTargetCpnsIndexName()).
		Query(query).
		Size(1000).Do(ctx)

	if err != nil {
		return request.StudentProfileCPNSElastic{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return request.StudentProfileCPNSElastic{}, fmt.Errorf("student data not found")
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = item.(request.StudentProfileCPNSElastic)
	}

	return gres, nil
}

func UpsertBinsusStudentProfile(c *mockstruct.SyncBinsusProfile) error {
	opts := options.Update().SetUpsert(true)
	collection := db.Mongodb.Collection("students")
	collectionParent := db.Mongodb.Collection("parent_datas")
	ctx := context.Background()
	std, err := GetStudentOnlyBySmartBTWID(c.SmartbtwID)

	if err != nil {
		return err
	}

	if len(std) == 0 {
		return nil
	}

	gender := 0
	if c.Gender == "L" {
		gender = 1
	}

	filter := bson.M{"smartbtw_id": c.SmartbtwID}
	update := bson.M{"$set": bson.M{
		"name":                 c.Name,
		"email":                c.Email,
		"gender":               gender,
		"school_origin":        c.LastEdName,
		"school_origin_id":     c.LastEdID,
		"address":              c.Address,
		"province_id":          c.ProvinceID,
		"region_id":            c.RegionID,
		"parent_name":          c.ParentName,
		"parent_number":        c.ParentNumber,
		"branch_code":          c.BranchCode,
		"updated_at":           time.Now(),
		"domicile_province_id": c.DomicileProvinceID,
		"domicile_region_id":   c.DomicileRegionID,
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	filterPar := bson.M{"student_id": std[0].ID}
	updatePar := bson.M{"$set": bson.M{
		"parent_name":   c.ParentName,
		"parent_number": c.ParentNumber,
		"updatedAt":     time.Now(),
	}}

	_, err = collectionParent.UpdateOne(ctx, filterPar, updatePar, opts)
	if err != nil {
		return err
	}

	_, err2 := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", c.SmartbtwID)).
		Doc(map[string]interface{}{
			"name":                 c.Name,
			"email":                c.Email,
			"gender":               c.Gender,
			"last_ed_name":         c.LastEdName,
			"last_ed_id":           c.LastEdID,
			"last_ed_type":         c.LastEdType,
			"last_ed_region":       c.LastEdRegion,
			"last_ed_region_id":    c.LastEdRegionID,
			"address":              c.Address,
			"province_id":          c.ProvinceID,
			"province":             c.Province,
			"region_id":            c.RegionID,
			"region":               c.Region,
			"parent_name":          c.ParentName,
			"parent_number":        c.ParentNumber,
			"branch_code":          c.BranchCode,
			"branch_name":          c.BranchName,
			"updated_at":           time.Now(),
			"domicile_province_id": c.DomicileProvinceID,
			"domicile_region_id":   c.DomicileRegionID,
			"domicile_province":    c.DomicileProvince,
			"domicile_region":      c.DomicileRegion,
		}).
		DocAsUpsert(true).
		Do(context.Background())

	uri := fmt.Sprintf("%s/internal/students/sync-binsus", os.Getenv("API_GATEWAY_HOST_URL"))
	ns, _ := sonic.Marshal(c)

	client := &http.Client{}

	request, err := http.NewRequest("POST", uri,
		bytes.NewBuffer(ns))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Office-Token", os.Getenv("API_GATEWAY_TOKEN"))
	if err != nil {
		return errors.New("creating request to api gateway " + err.Error())
	}
	resp, err := client.Do(request)
	if err != nil {
		return errors.New("doing request to api gateway " + err.Error())
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("reading response of api gateway " + err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New("failed to sync to api gateway")
	}

	return err2
}

func UpsertBinsusFinalStudentProfile(c *mockstruct.SyncBinsusFinalProfile) error {
	opts := options.Update().SetUpsert(true)
	collection := db.Mongodb.Collection("students")
	collectionParent := db.Mongodb.Collection("parent_datas")
	collectionStudentTargets := db.Mongodb.Collection("student_targets")

	// Profile Update Section
	ctx := context.Background()

	compData, err := GetCompetitionFromCompMap(uint(c.School.CompetitionID))

	if err != nil {
		return err
	}
	std, err := GetStudentOnlyBySmartBTWID(c.SmartbtwID)

	if err != nil {
		return err
	}

	if len(std) == 0 {
		return nil
	}

	gender := 0
	if c.Profile.Gender == "L" {
		gender = 1
	}

	filter := bson.M{"smartbtw_id": c.SmartbtwID}
	update := bson.M{"$set": bson.M{
		"name":                 c.Profile.Name,
		"email":                c.Profile.Email,
		"gender":               gender,
		"school_origin":        c.Profile.LastEdName,
		"school_origin_id":     c.Profile.LastEdID,
		"address":              c.Profile.Address,
		"province_id":          c.Profile.ProvinceID,
		"region_id":            c.Profile.RegionID,
		"parent_name":          c.Profile.ParentName,
		"parent_number":        c.Profile.ParentNumber,
		"branch_code":          c.Profile.BranchCode,
		"interest":             "BINSUS",
		"updated_at":           time.Now(),
		"domicile_province_id": c.Profile.DomicileProvinceID,
		"domicile_region_id":   c.Profile.DomicileRegionID,
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	filterPar := bson.M{"student_id": std[0].ID}
	updatePar := bson.M{"$set": bson.M{
		"parent_name":   c.Profile.ParentName,
		"parent_number": c.Profile.ParentNumber,
		"updatedAt":     time.Now(),
	}}

	_, err = collectionParent.UpdateOne(ctx, filterPar, updatePar, opts)
	if err != nil {
		return err
	}

	_, err2 := db.ElasticClient.Update().
		Index(db.GetStudentProfileIndexName()).
		Id(fmt.Sprintf("%d", c.Profile.SmartbtwID)).
		Doc(map[string]interface{}{
			"name":                 c.Profile.Name,
			"email":                c.Profile.Email,
			"gender":               c.Profile.Gender,
			"last_ed_name":         c.Profile.LastEdName,
			"last_ed_id":           c.Profile.LastEdID,
			"last_ed_type":         c.Profile.LastEdType,
			"last_ed_region":       c.Profile.LastEdRegion,
			"last_ed_region_id":    c.Profile.LastEdRegionID,
			"address":              c.Profile.Address,
			"province_id":          c.Profile.ProvinceID,
			"province":             c.Profile.Province,
			"region_id":            c.Profile.RegionID,
			"region":               c.Profile.Region,
			"parent_name":          c.Profile.ParentName,
			"parent_number":        c.Profile.ParentNumber,
			"branch_code":          c.Profile.BranchCode,
			"branch_name":          c.Profile.BranchName,
			"updated_at":           time.Now(),
			"domicile_province_id": c.Profile.DomicileProvinceID,
			"domicile_region_id":   c.Profile.DomicileRegionID,
			"domicile_province":    c.Profile.DomicileProvince,
			"domicile_region":      c.Profile.DomicileRegion,
			"interest":             "BINSUS",
		}).
		DocAsUpsert(true).
		Do(context.Background())
	if err2 != nil {
		return err2
	}
	// Student Target Update
	filSt := bson.M{"smartbtw_id": c.SmartbtwID, "target_type": "PTK"}
	updateTar := bson.M{"$set": bson.M{
		"is_active":  false,
		"can_update": false,
	}}
	_, err = collectionStudentTargets.UpdateMany(ctx, filSt, updateTar, opts)
	if err != nil {
		return err
	}
	// polbitType := "PUSAT"
	// if strings.Contains(strings.ToLower(c.School.PolbitName), "kabupaten") {
	// 	polbitType = fmt.Sprintf("%s_%s", c.School.PolbitType, "REGION")
	// } else if c.School.PolbitType != "PUSAT" {
	// 	polbitType = fmt.Sprintf("%s_%s", c.School.PolbitType, "PROVINCE")
	// }
	polbitType := "PUSAT"

	if compData.LocationID != nil {
		polbitType = fmt.Sprintf("%s_%s", compData.PolbitType, compData.Location.Type)
	}

	payload1 := models.StudentTarget{
		SmartbtwID:          c.SmartbtwID,
		SchoolID:            c.School.SchoolID,
		MajorID:             c.School.MajorID,
		SchoolName:          c.School.SchoolName,
		MajorName:           c.School.MajorName,
		TargetScore:         c.School.TargetScore,
		TargetType:          "PTK",
		PolbitType:          polbitType,
		PolbitCompetitionID: &c.School.CompetitionID,
		PolbitLocationID:    compData.LocationID,
		CanUpdate:           true,
		IsActive:            true,
		Position:            0,
		Type:                string(models.PRIMARY),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		DeletedAt:           nil,
	}
	_, err3 := collectionStudentTargets.InsertOne(ctx, payload1)
	if err3 != nil {
		return err3
	}

	averages, err := GetStudentHistoryPTKElastic(c.SmartbtwID, false)
	if err != nil {
		return fmt.Errorf("student ID %d data failed to fetch history ptk from elastic with error %s", c.SmartbtwID, err.Error())

	}

	// twkScore := float64(0)
	// tiuScore := float64(0)
	// tkpScore := float64(0)
	totalScore := float64(0)
	passingTotalItem := 0
	challengeRecord := []request.CreateHistoryPtk{}
	for _, k := range averages {
		if strings.ToLower(k.PackageType) == "challenge-uka" || strings.ToUpper(k.ModuleType) == "WITH_CODE" {
			challengeRecord = append(challengeRecord, k)
		}

	}
	for _, k := range challengeRecord {
		// twkScore += k.Twk
		// tiuScore += k.Tiu
		// tkpScore += k.Tkp
		totalScore += k.Total
	}
	if len(challengeRecord) < 11 {
		passingTotalItem = 10
	} else {
		passingTotalItem = len(challengeRecord)
	}
	// atwk := math.Round(helpers.RoundFloat((twkScore / float64(passingTotalItem)), 2))
	// atiu := math.Round(helpers.RoundFloat((tiuScore / float64(passingTotalItem)), 2))
	// atkp := math.Round(helpers.RoundFloat((tkpScore / float64(passingTotalItem)), 2))
	pAtt := math.Round(helpers.RoundFloat((totalScore / float64(passingTotalItem)), 2))
	// pAtt := atwk + atiu + atkp

	if _isNaNorInf(pAtt) {
		pAtt = 0
	}

	percATT := helpers.RoundFloat((pAtt/float64(c.School.TargetScore))*100, 2)

	if percATT > 99 {
		percATT = 99
	}

	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", c.SmartbtwID))

	// _, err1 := db.ElasticClient.Update().
	// 	Index(db.GetStudentTargetPtkIndexName()).
	// 	Id(fmt.Sprintf("%d_PTK", c.SmartbtwID)).
	// 	Doc(map[string]interface{}{
	// 		"school_id":                        c.School.SchoolID,
	// 		"major_id":                         c.School.MajorID,
	// 		"school_name":                      c.School.SchoolName,
	// 		"major_name":                       c.School.MajorName,
	// 		"polbit_type":                      polbitType,
	// 		"polbit_competition_id":            c.School.CompetitionID,
	// 		"polbit_location_id":               compData.LocationID,
	// 		"passing_recommendation_avg_score": pAtt,
	// 		"passing_recommendation_avg_percent_score": percATT,
	// 		"target_score": c.School.TargetScore,
	// 	}).
	// 	DocAsUpsert(true).
	// 	Do(context.Background())

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
		"school_id":                        c.School.SchoolID,
		"major_id":                         c.School.MajorID,
		"school_name":                      c.School.SchoolName,
		"major_name":                       c.School.MajorName,
		"polbit_type":                      polbitType,
		"polbit_competition_id":            c.School.CompetitionID,
		"polbit_location_id":               compData.LocationID,
		"passing_recommendation_avg_score": pAtt,
		"passing_recommendation_avg_percent_score": percATT,
		"target_score": c.School.TargetScore,
	})

	_, err1 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentTargetPtkIndexName()).
		Query(bq).
		Script(script).
		DoAsync(context.Background())
	if err1 != nil {
		return err1
	}
	// _, err4 := db.ElasticClient.Update().
	// 	Index(db.GetStudentProfileIndexName()).
	// 	Id(fmt.Sprintf("%d", c.SmartbtwID)).
	// 	Doc(map[string]interface{}{
	// 		"school_ptk_id":             c.School.SchoolID,
	// 		"major_ptk_id":              c.School.MajorID,
	// 		"school_name_ptk":           c.School.SchoolName,
	// 		"major_name_ptk":            c.School.MajorName,
	// 		"target_score_ptk":          c.School.TargetScore,
	// 		"polbit_type_ptk":           polbitType,
	// 		"polbit_competition_ptk_id": c.School.CompetitionID,
	// 		"polbit_location_ptk_id":    compData.LocationID,
	// 		"created_at_ptk":            time.Now(),
	// 	}).
	// 	DocAsUpsert(true).
	// 	Do(context.Background())
	script2 := elastic.NewScript(`
	ctx._source.school_ptk_id = params.school_ptk_id;
	ctx._source.school_name_ptk = params.school_name_ptk;
	ctx._source.major_name_ptk = params.major_name_ptk;
	ctx._source.major_ptk_id = params.major_ptk_id;
	ctx._source.polbit_competition_ptk_id = params.polbit_competition_ptk_id;
	ctx._source.polbit_location_ptk_id = params.polbit_location_ptk_id;
	ctx._source.polbit_type_ptk = params.polbit_type_ptk;
	ctx._source.target_score_ptk = params.target_score_ptk;
	ctx._source.created_at_ptk = params.created_at_ptk;
	`).Params(map[string]interface{}{
		"school_ptk_id":             c.School.SchoolID,
		"major_ptk_id":              c.School.MajorID,
		"school_name_ptk":           c.School.SchoolName,
		"major_name_ptk":            c.School.MajorName,
		"target_score_ptk":          c.School.TargetScore,
		"polbit_type_ptk":           polbitType,
		"polbit_competition_ptk_id": c.School.CompetitionID,
		"polbit_location_ptk_id":    compData.Location.ID,
		"created_at_ptk":            time.Now(),
	})

	_, err4 := elastic.NewUpdateByQueryService(db.ElasticClient).
		Index(db.GetStudentProfileIndexName()).
		Query(bq).
		Script(script2).
		DoAsync(context.Background())

	if err4 != nil {
		return err4
	}
	c.Profile.Interest = "BINSUS"

	uri := fmt.Sprintf("%s/internal/students/sync-binsus", os.Getenv("API_GATEWAY_HOST_URL"))
	ns, _ := sonic.Marshal(c.Profile)

	client := &http.Client{}

	request, err := http.NewRequest("POST", uri,
		bytes.NewBuffer(ns))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Office-Token", os.Getenv("API_GATEWAY_TOKEN"))
	if err != nil {
		return errors.New("creating request to api gateway " + err.Error())
	}
	resp, err := client.Do(request)
	if err != nil {
		return errors.New("doing request to api gateway " + err.Error())
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("reading response of api gateway " + err.Error())
	}
	if resp.StatusCode != 200 {
		return errors.New("failed to sync to api gateway")
	}

	return nil
}

func GetStudentProfileByArrayOfSmartbtwID(ids []uint) ([]*request.StudentProfileElastic, error) {
	ctx := context.Background()

	var results []*request.StudentProfileElastic

	var idsInterface []interface{}
	for _, id := range ids {
		idsInterface = append(idsInterface, id)
	}

	query := elastic.NewTermsQuery("smartbtw_id", idsInterface...)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentProfileIndexName()).
		Query(query).
		Size(1000).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if res.TotalHits() == 0 {
		return nil, fmt.Errorf("no student data found for the given smartbtw_ids")
	}

	for _, hit := range res.Hits.Hits {
		var studentProfile request.StudentProfileElastic
		err := json.Unmarshal(hit.Source, &studentProfile)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal student profile: %v", err)
		}
		results = append(results, &studentProfile)
	}

	return results, nil
}

func GetStudentProfileByArrayOfEmail(emails []string) ([]*request.StudentProfileElastic, error) {
	ctx := context.Background()

	var results []*request.StudentProfileElastic

	var idsInterface []interface{}
	for _, id := range emails {
		idsInterface = append(idsInterface, id)
	}

	query := elastic.NewTermsQuery("email.keyword", idsInterface...)

	res, err := db.ElasticClient.Search().
		Index(db.GetStudentProfileIndexName()).
		Query(query).
		Size(1000).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	if res.TotalHits() == 0 {
		return nil, fmt.Errorf("no student data found for the given emails")
	}

	for _, hit := range res.Hits.Hits {
		var studentProfile request.StudentProfileElastic
		err := json.Unmarshal(hit.Source, &studentProfile)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal student profile: %v", err)
		}
		results = append(results, &studentProfile)
	}

	return results, nil
}

func calculatePassProbability(averageScore, passingScore float64) float64 {
	passProbability := averageScore / passingScore

	return passProbability
}

func GetPerformaSiswa(req *request.GetPerformaSiswa) ([]*request.ResultsPerformaSiswa, map[string]any, error) {

	stTotal := 0
	avgAllTwk := float64(0)
	avgAllTiu := float64(0)
	avgAllTkp := float64(0)
	avgAllTotal := float64(0)

	stdRes, err := GetStudentProfileByArrayOfSmartbtwID(req.SmartBtwID)
	if err != nil {
		return nil, nil, err
	}

	stTotal = len(stdRes)

	results := []*request.ResultsPerformaSiswa{}
	for _, res := range stdRes {

		var clsData models.ClassMemberElastic
		stdClsList, errCls := GetStudentJoinedClassList(int32(res.SmartbtwID), req.ClassYear, false)
		if errCls != nil {
			continue
		}

		if req.ClassTags == "BINSUS" {
			isChoosen := false
			for _, k := range stdClsList {
				for _, tagsClass := range k.Tags {
					if tagsClass == "BINSUS" {
						isChoosen = true
						clsData = k
						break
					}
				}
				if isChoosen {
					break
				}
			}
			if !isChoosen {
				continue
			}
		} else {
			joinedClassList := ""
			isChoosen := false
			for _, k := range stdClsList {
				for _, tagsClass := range k.Tags {
					if tagsClass == "BINSUS" || tagsClass == "REGULER" || tagsClass == "INTENSIF" {
						joinedClassList = tagsClass
						// 	isChoosen = true
						// 	clsData = k
						break
					}
				}
				if joinedClassList != "" {
					if req.ClassTags != "" {
						if req.ClassTags != joinedClassList {
							break
						} else {
							clsData = k
							isChoosen = true
						}
					} else {

						clsData = k
						isChoosen = true
					}

				}
				if isChoosen {
					break
				}
			}
			if !isChoosen {
				continue
			}
		}

		trialCount := 0
		hisRes, er := GetHistoryPTKElastic(uint(res.SmartbtwID), req.Type)
		if er != nil {
			return nil, nil, er
		}

		studentTarget, errStudentTarget := GetStudentProfilePTKElastic(res.SmartbtwID)

		var scoreTotal float64
		var scoreTWK float64
		var scoreTIU float64
		var scoreTKP float64
		var avgScore float64
		totalPassed := int(0)
		totalFailed := int(0)
		twkPass := int(0)
		tiuPass := int(0)
		tkpPass := int(0)
		twkFailed := int(0)
		tiuFailed := int(0)
		tkpFailed := int(0)
		withCode := int(0)
		stgTotal := int(0)

		hRes := []request.CreateHistoryPtk{}

		for _, t := range hisRes {
			if t.ExamName == "" {
				continue
			}
			if t.ModuleType == "TESTING" || t.ModuleType == "TRIAL" {
				trialCount += 1
				continue
			}
			scoreTotal += t.Total
			scoreTWK += t.Twk
			scoreTIU += t.Tiu
			scoreTKP += t.Tkp
			if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
				totalPassed += 1
			}
			if t.Twk < t.TwkPass || t.Tiu < t.TiuPass || t.Tkp < t.TkpPass {
				totalFailed += 1
			}
			if t.Twk >= t.TwkPass {
				twkPass += 1
			}
			if t.Tiu >= t.TiuPass {
				tiuPass += 1
			}
			if t.Tkp >= t.TkpPass {
				tkpPass += 1
			}
			if t.Twk < t.TwkPass {
				twkFailed += 1
			}
			if t.Tiu < t.TiuPass {
				tiuFailed += 1
			}
			if t.Tkp < t.TkpPass {
				tkpFailed += 1
			}
			if t.PackageType == "WITH_CODE" {
				withCode += 1
			} else {
				stgTotal += 1
			}
			hRes = append(hRes, t)
		}

		sort.SliceStable(hRes, func(i, j int) bool {
			return hRes[i].Start.After(hRes[j].End)
		})

		totalHis := len(hisRes) - trialCount
		done := len(hisRes) - trialCount
		avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
		avgScoreTWK := helpers.RoundFloat(scoreTWK/float64(totalHis), 1)
		avgScoreTIU := helpers.RoundFloat(scoreTIU/float64(totalHis), 1)
		avgScoreTKP := helpers.RoundFloat(scoreTKP/float64(totalHis), 1)

		passPercent := helpers.RoundFloat(float64(totalPassed)/float64(totalHis)*100, 1)
		passPercentTWK := helpers.RoundFloat(float64(twkPass)/float64(totalHis)*100, 1)
		passPercentTIU := helpers.RoundFloat(float64(tiuPass)/float64(totalHis)*100, 1)
		passPercentTKP := helpers.RoundFloat(float64(tkpPass)/float64(totalHis)*100, 1)

		if _isNaNorInf(avgScore) {
			avgScore = 0
		}
		if _isNaNorInf(avgScoreTWK) {
			avgScoreTWK = 0
		}
		if _isNaNorInf(avgScoreTIU) {
			avgScoreTIU = 0
		}
		if _isNaNorInf(avgScoreTKP) {
			avgScoreTKP = 0
		}

		if _isNaNorInf(passPercent) {
			passPercent = 0
		}
		if _isNaNorInf(passPercentTWK) {
			passPercentTWK = 0
		}
		if _isNaNorInf(passPercentTIU) {
			passPercentTIU = 0
		}
		if _isNaNorInf(passPercentTKP) {
			passPercentTKP = 0
		}
		failePercentTWK := 100 - passPercentTWK
		failePercentTIU := 100 - passPercentTIU
		failePercentTKP := 100 - passPercentTKP

		scoreKeys := []string{"TWK", "TIU", "TKP"}

		scVal := request.ScoreValues{
			TWK: request.Values{
				Total:         totalHis,
				Passed:        twkPass,
				Failed:        twkFailed,
				TotalScore:    float32(scoreTWK),
				AverageScore:  float32(avgScoreTWK),
				PassedPercent: float32(passPercentTWK),
				FailedPercent: float32(failePercentTWK),
			},
			TIU: request.Values{
				Total:         totalHis,
				Passed:        tiuPass,
				Failed:        tiuFailed,
				TotalScore:    float32(scoreTIU),
				AverageScore:  float32(avgScoreTIU),
				PassedPercent: float32(passPercentTIU),
				FailedPercent: float32(failePercentTIU),
			},
			TKP: request.Values{
				Total:         totalHis,
				Passed:        tkpPass,
				Failed:        tkpFailed,
				TotalScore:    float32(scoreTKP),
				AverageScore:  float32(avgScoreTKP),
				PassedPercent: float32(passPercentTKP),
				FailedPercent: float32(failePercentTKP),
			},
		}

		prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
		if err != nil {
			return nil, nil, err
		}

		countStgLv := 0
		for _, obj := range prodRe.Data {
			isSkipped := false
			if strings.ToLower(obj.ProductProgram) != "skd" {
				continue
			}
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
					isSkipped = true
					continue
				}
				if strings.Contains(tag, "MATERIAL") {
					isSkipped = true
					continue
				}
				if strings.Contains(tag, "CPNS") {
					isSkipped = true
					continue
				}
			}
			if isSkipped {
				continue
			}

			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "STAGE_LEVEL_") && obj.ProductProgram == "skd" {
					countStgLv += 1
				}
			}
		}

		ownMod := withCode

		if (countStgLv - stgTotal) < 0 {
			ownMod += stgTotal
		} else {
			ownMod += countStgLv
		}

		donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)
		if _isNaNorInf(donePerct) {
			donePerct = 0
		}

		avgDonePercent := helpers.RoundFloat(float64(done)/float64(ownMod), 2)

		if _isNaNorInf(avgDonePercent) {
			avgDonePercent = 0
		}

		smr := request.Summary{
			AverageScore:  avgScore,
			Passed:        totalPassed,
			Failed:        totalFailed,
			Total:         totalHis,
			TotalScore:    float32(scoreTotal),
			PassedPercent: float32(passPercent),
			ScoreKeys:     scoreKeys,
			ScoreValues:   scVal,
			DonePercent:   float32(donePerct),
			Owned:         ownMod,
			Done:          float32(totalHis),
			AverageDone:   float32(avgDonePercent),
		}

		bc := "PT0000"
		bn := "Bimbel BTW (Kantor Pusat)"

		if res.BranchCode != nil {
			bc = *res.BranchCode
			bn = *res.BranchName

		}

		resScore := request.ResultsPerformaSiswa{
			BranchCode: bc,
			BranchName: bn,
			Name:       res.Name,
			SmartBtwID: res.SmartbtwID,
			Email:      res.Email,
			Summary:    smr,
			ClassInformation: requests.StudentClassInformation{
				JoinedClass:     true,
				ClassTitle:      clsData.Title,
				ClassType:       req.ClassTags,
				ClassYear:       int(clsData.Year),
				ClassJoined:     clsData.CreatedAt,
				ClassStatus:     clsData.Status,
				ClassBranchCode: clsData.BranchCode,
			},
			PTKTarget:     requests.StudentTargetDataPerforma{},
			HistoryRecord: hRes,
		}

		ctPercent := float64(0)

		if !_isNaNorInf(studentTarget.PassingRecommendationAvgPercentScore) {
			ctPercent = studentTarget.PassingRecommendationAvgPercentScore
		}

		if errStudentTarget == nil {
			resScore.PTKTarget = request.StudentTargetDataPerforma{
				SchoolID:                  studentTarget.SchoolID,
				MajorID:                   studentTarget.MajorID,
				SchoolName:                studentTarget.SchoolName,
				MajorName:                 studentTarget.MajorName,
				PolbitType:                studentTarget.PolbitType,
				PolbitCompetitionID:       studentTarget.PolbitCompetitionID,
				PolbitLocationID:          studentTarget.PolbitLocationID,
				TargetScore:               int(studentTarget.TargetScore),
				CurrentTargetPercentScore: ctPercent,
			}
		}

		bknScore := map[string]any{
			"twk":           float64(0),
			"twk_pass":      false,
			"tiu":           float64(0),
			"tiu_pass":      false,
			"tkp":           float64(0),
			"tkp_pass":      false,
			"total":         float64(0),
			"year":          2023,
			"bkn_attempted": false,
			"is_pass":       false,
		}

		finalPassingPercentage := float64(0)

		deviationExist := false
		twkDeviation := float64(0)
		twkPercentage := float64(0)
		tiuDeviation := float64(0)
		tiuPercentage := float64(0)
		tkpDeviation := float64(0)
		tkpPercentage := float64(0)
		totalDeviation := float64(0)
		totalPercentage := float64(0)
		bknTargetDeviation := float64(0)
		bknTargetPercentage := float64(0)
		ukaTargetDeviation := float64(0)
		ukaTargetPercentage := float64(0)

		ukaTargetDeviation = math.Round((avgScore-studentTarget.TargetScore)*10) / 10
		ukaTargetPercentage = math.Round((float64(ukaTargetDeviation) / float64(studentTarget.TargetScore)) * 100)

		if _isNaNorInf(ukaTargetDeviation) {
			ukaTargetDeviation = 0
		}
		if _isNaNorInf(ukaTargetPercentage) {
			ukaTargetPercentage = 0
		}

		yr := time.Now().Year()
		if req.ClassYear != nil {
			yr = *req.ClassYear
		}

		bknSc, err := GetSingleBKNScoreByYearAndStudent(res.SmartbtwID, uint16(yr))
		if err == nil {
			isBknPass := false
			if bknSc.Twk >= 65 && bknSc.Tiu >= 80 && bknSc.Tkp >= 156 {
				isBknPass = true
			}

			bknScore = map[string]any{
				"twk":           bknSc.Twk,
				"twk_pass":      bknSc.Twk >= 65,
				"tiu":           bknSc.Tiu,
				"tiu_pass":      bknSc.Tiu >= 80,
				"tkp":           bknSc.Tkp,
				"tkp_pass":      bknSc.Tkp >= 156,
				"total":         bknSc.Total,
				"year":          bknSc.Year,
				"bkn_attempted": true,
				"is_pass":       isBknPass,
			}

			finalPassingPercentage = helpers.RoundFloat(float64(bknSc.Total)/float64(studentTarget.TargetScore)*100, 1)
			if _isNaNorInf(finalPassingPercentage) {
				finalPassingPercentage = 0
			}

			twkDeviation = math.Round((bknSc.Twk-avgScoreTWK)*10) / 10
			twkPercentage = math.Round((float64(twkDeviation) / float64(avgScoreTWK)) * 100)
			tiuDeviation = math.Round((bknSc.Tiu-avgScoreTIU)*10) / 10
			tiuPercentage = math.Round((float64(tiuDeviation) / float64(avgScoreTIU)) * 100)
			tkpDeviation = math.Round((bknSc.Tkp-avgScoreTKP)*10) / 10
			tkpPercentage = math.Round((float64(tkpDeviation) / float64(avgScoreTKP)) * 100)
			totalDeviation = math.Round((bknSc.Total-avgScore)*10) / 10
			totalPercentage = math.Round((float64(totalDeviation) / float64(avgScore)) * 100)
			bknTargetDeviation = math.Round((bknSc.Total-studentTarget.TargetScore)*10) / 10
			bknTargetPercentage = math.Round((float64(bknTargetDeviation) / float64(studentTarget.TargetScore)) * 100)
			if _isNaNorInf(twkDeviation) {
				twkDeviation = 0
			}
			if _isNaNorInf(twkPercentage) {
				twkPercentage = 0
			}
			if _isNaNorInf(tiuDeviation) {
				tiuDeviation = 0
			}
			if _isNaNorInf(tiuPercentage) {
				tiuPercentage = 0
			}
			if _isNaNorInf(tkpDeviation) {
				tkpDeviation = 0
			}
			if _isNaNorInf(tkpPercentage) {
				tkpPercentage = 0
			}
			if _isNaNorInf(totalDeviation) {
				totalDeviation = 0
			}
			if _isNaNorInf(totalPercentage) {
				totalPercentage = 0
			}
			if _isNaNorInf(bknTargetDeviation) {
				bknTargetDeviation = 0
			}
			if _isNaNorInf(bknTargetPercentage) {
				bknTargetPercentage = 0
			}
			deviationExist = true
		}

		targetPassingPercentage := helpers.RoundFloat(float64(avgScore)/float64(studentTarget.TargetScore)*100, 1)
		if _isNaNorInf(targetPassingPercentage) {
			targetPassingPercentage = 0
		}
		resScore.BKNScore = bknScore
		resScore.Summary.Deviation = map[string]any{
			"twk": map[string]any{
				"percentage":  twkPercentage,
				"differences": twkDeviation,
			},
			"tiu": map[string]any{
				"percentage":  tiuPercentage,
				"differences": tiuDeviation,
			},
			"tkp": map[string]any{
				"percentage":  tkpPercentage,
				"differences": tkpDeviation,
			},
			"total": map[string]any{
				"percentage":  totalPercentage,
				"differences": totalDeviation,
			},
			"bkn": map[string]any{
				"percentage":  bknTargetPercentage,
				"differences": bknTargetDeviation,
			},
			"uka": map[string]any{
				"percentage":  ukaTargetPercentage,
				"differences": ukaTargetDeviation,
			},
			"available": deviationExist,
		}

		resScore.Summary.FinalPassingPercent = finalPassingPercentage
		resScore.Summary.TargetPassingPercent = targetPassingPercentage

		avgAllTwk += float64(scVal.TWK.AverageScore)
		avgAllTiu += float64(scVal.TIU.AverageScore)
		avgAllTkp += float64(scVal.TKP.AverageScore)
		avgAllTotal += float64(resScore.Summary.AverageScore)

		results = append(results, &resScore)

	}

	twkAvg := helpers.RoundFloat(float64(avgAllTwk)/float64(stTotal), 1)
	tiuAvg := helpers.RoundFloat(float64(avgAllTiu)/float64(stTotal), 1)
	tkpAvg := helpers.RoundFloat(float64(avgAllTkp)/float64(stTotal), 1)
	totalAvg := helpers.RoundFloat(float64(avgAllTotal)/float64(stTotal), 1)
	if _isNaNorInf(twkAvg) {
		twkAvg = 0
	}
	if _isNaNorInf(tiuAvg) {
		tiuAvg = 0
	}
	if _isNaNorInf(tkpAvg) {
		tkpAvg = 0
	}
	if _isNaNorInf(totalAvg) {
		totalAvg = 0
	}
	avgData := map[string]any{
		"twk":           twkAvg,
		"tiu":           tiuAvg,
		"tkp":           tkpAvg,
		"total":         totalAvg,
		"student_total": stTotal,
	}

	return results, avgData, nil
}

func GetPerformaSiswaPTK(req *request.GetPerformaSiswaUKA) ([]*request.ResultsPerformaSiswa, map[string]any, error) {

	stTotal := 0
	avgAllTwk := float64(0)
	avgAllTiu := float64(0)
	avgAllTkp := float64(0)
	avgAllTotal := float64(0)
	skip := ""

	stdRes, err := GetStudentProfileByArrayOfSmartbtwID(req.SmartBtwID)
	if err != nil {
		return nil, nil, err
	}

	stTotal = len(stdRes)

	results := []*request.ResultsPerformaSiswa{}
	for _, res := range stdRes {

		var clsData models.ClassMemberElastic
		// stdClsList, errCls := GetStudentJoinedClassList(int32(res.SmartbtwID), nil, false)
		// if errCls != nil {
		// 	continue
		// }

		var filter string
		switch strings.ToUpper(req.TypeModule) {
		case "PRE_UKA":
			filter = "pre-uka"
		case "ALL_MODULE":
			filter = "all-module"
		case "UKA_STAGE":
			filter = "challenge-uka"
		}

		trialCount := 0
		hisRes, er := GetHistoryPTKElasticFetchStudentReport(uint(res.SmartbtwID), req.TypeStages, filter)
		if er != nil {
			return nil, nil, er
		}

		// resPkg, err := GetHistoryPTKElasticTypeStage(uint(res.SmartbtwID), req.TypeStages)
		// if err != nil {
		// 	return nil, nil, err
		// }

		// fmt.Println(resPkg)

		// hisRes = append(hisRes, resPkg...)
		studentTarget, errStudentTarget := GetStudentProfilePTKElastic(res.SmartbtwID)

		var scoreTotal float64
		var scoreTWK float64
		var scoreTIU float64
		var scoreTKP float64
		var avgScore float64
		totalPassed := int(0)
		totalFailed := int(0)
		twkPass := int(0)
		tiuPass := int(0)
		tkpPass := int(0)
		twkFailed := int(0)
		tiuFailed := int(0)
		tkpFailed := int(0)
		withCode := int(0)
		stgTotal := int(0)
		// isPassTiu := false
		// isPassTwk := false
		// isPassTkp := false

		hRes := []request.CreateHistoryPtk{}

		for _, t := range hisRes {
			if t.ExamName == "" {
				continue
			}
			// if strings.Contains(t.ExamName, "Post-Test") || strings.Contains(t.ExamName, "Pre-Test") {
			// 	continue
			// }
			if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
				continue
			}
			if t.ModuleType == "TESTING" || t.ModuleType == "TRIAL" {
				trialCount += 1
				continue
			}
			scoreTotal += t.Total
			scoreTWK += t.Twk
			scoreTIU += t.Tiu
			scoreTKP += t.Tkp
			if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
				totalPassed += 1
				// isPass = true
			}
			if t.Twk < t.TwkPass || t.Tiu < t.TiuPass || t.Tkp < t.TkpPass {
				totalFailed += 1
			}
			if t.Twk >= t.TwkPass {
				twkPass += 1
				// isPassTwk = true
			}
			if t.Tiu >= t.TiuPass {
				tiuPass += 1
				// isPassTiu = true
			}
			if t.Tkp >= t.TkpPass {
				tkpPass += 1
				// isPassTkp = true
			}
			if t.Twk < t.TwkPass {
				twkFailed += 1
			}
			if t.Tiu < t.TiuPass {
				tiuFailed += 1
			}
			if t.Tkp < t.TkpPass {
				tkpFailed += 1
			}
			if t.PackageType == "WITH_CODE" {
				withCode += 1
			} else {
				stgTotal += 1
			}
			hRes = append(hRes, t)
		}

		sort.SliceStable(hRes, func(i, j int) bool {
			return hRes[i].Start.After(hRes[j].End)
		})

		totalHis := len(hisRes) - trialCount
		done := len(hisRes) - trialCount
		avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
		avgScoreTWK := helpers.RoundFloat(scoreTWK/float64(totalHis), 1)
		avgScoreTIU := helpers.RoundFloat(scoreTIU/float64(totalHis), 1)
		avgScoreTKP := helpers.RoundFloat(scoreTKP/float64(totalHis), 1)

		passPercent := helpers.RoundFloat(float64(totalPassed)/float64(totalHis)*100, 1)
		passPercentTWK := helpers.RoundFloat(float64(twkPass)/float64(totalHis)*100, 1)
		passPercentTIU := helpers.RoundFloat(float64(tiuPass)/float64(totalHis)*100, 1)
		passPercentTKP := helpers.RoundFloat(float64(tkpPass)/float64(totalHis)*100, 1)

		if _isNaNorInf(avgScore) {
			avgScore = 0
		}
		if _isNaNorInf(avgScoreTWK) {
			avgScoreTWK = 0
		}
		if _isNaNorInf(avgScoreTIU) {
			avgScoreTIU = 0
		}
		if _isNaNorInf(avgScoreTKP) {
			avgScoreTKP = 0
		}

		if _isNaNorInf(passPercent) {
			passPercent = 0
		}
		if _isNaNorInf(passPercentTWK) {
			passPercentTWK = 0
		}
		if _isNaNorInf(passPercentTIU) {
			passPercentTIU = 0
		}
		if _isNaNorInf(passPercentTKP) {
			passPercentTKP = 0
		}
		failePercentTWK := 100 - passPercentTWK
		failePercentTIU := 100 - passPercentTIU
		failePercentTKP := 100 - passPercentTKP

		scoreKeys := []string{"TWK", "TIU", "TKP"}

		scVal := request.ScoreValues{
			IsPass: avgScoreTWK >= 65 && avgScoreTIU >= 80 && avgScoreTKP >= 156,
			TWK: request.Values{
				Total:         totalHis,
				Passed:        twkPass,
				Failed:        twkFailed,
				IsPass:        avgScoreTWK >= 65,
				TotalScore:    float32(scoreTWK),
				AverageScore:  float32(avgScoreTWK),
				PassedPercent: float32(passPercentTWK),
				FailedPercent: float32(failePercentTWK),
			},
			TIU: request.Values{
				Total:         totalHis,
				Passed:        tiuPass,
				Failed:        tiuFailed,
				IsPass:        avgScoreTKP >= 80,
				TotalScore:    float32(scoreTIU),
				AverageScore:  float32(avgScoreTIU),
				PassedPercent: float32(passPercentTIU),
				FailedPercent: float32(failePercentTIU),
			},
			TKP: request.Values{
				Total:         totalHis,
				Passed:        tkpPass,
				Failed:        tkpFailed,
				IsPass:        avgScoreTKP >= 156,
				TotalScore:    float32(scoreTKP),
				AverageScore:  float32(avgScoreTKP),
				PassedPercent: float32(passPercentTKP),
				FailedPercent: float32(failePercentTKP),
			},
		}

		prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
		if err != nil {
			return nil, nil, err
		}

		// resStages, err := GetAllStudentStageClass("PTK")
		// if err != nil {
		// 	return nil, nil, err
		// }
		countStgLv := 0

		for _, obj := range prodRe.Data {
			isSkipped := false
			if strings.ToLower(obj.ProductProgram) != "skd" {
				continue
			}
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
					isSkipped = true
					continue
				}
				if strings.Contains(tag, "MATERIAL") {
					isSkipped = true
					continue
				}
				if strings.Contains(tag, "CPNS") {
					isSkipped = true
					continue
				}
			}
			if isSkipped {
				continue
			}
			if filter != "" {
				if filter == "pre-uka" {
					for _, tag := range obj.ProductTags {
						if req.TypeStages == "UMUM" {
							if strings.Contains(tag, "STAGE_PRE_UKA") {
								countStgLv += 1
							}
						} else {
							re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
							matches := re.FindStringSubmatch(tag)

							// Check if there are captured substrings
							if len(matches) > 1 {

								// Iterate over captured substrings starting from index 1
								for _, v := range matches[1:] {

									// Convert both strings to lowercase for case-insensitive comparison
									length := (len(skip))
									var name string
									if len(v) == length && name == res.Name {
										continue // Skip processing if already processed
									} else {
										name = res.Name
										if v == "BINSUS" {
											skip = v
											countStgLv += 14
										} else if v == "REGULER" {
											skip = v
											countStgLv += 14
										} else {
											skip = v
											countStgLv += 14
										}
										// Mark the value as processed
										// resStages, err := GetAllStudentStageClass("PTK", v)
										// if err != nil {
										// 	return nil, nil, err
										// }

										// for _, mdl := range resStages.Data {
										// 	if mdl.ModuleType == "PLATINUM" {
										// 		countStgLv += 1
										// 	}
										// }
									}
								}
							}

						}

					}
				} else if filter == "challenge-uka" {
					for _, tag := range obj.ProductTags {
						if req.TypeStages == "UMUM" {
							if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
								countStgLv += 1
							}
						} else {
							re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
							matches := re.FindStringSubmatch(tag)

							// Check if there are captured substrings
							if len(matches) > 1 {

								// Iterate over captured substrings starting from index 1
								for _, v := range matches[1:] {

									// Convert both strings to lowercase for case-insensitive comparison
									length := (len(skip))

									var name string
									if len(v) == length && name == res.Name {
										continue // Skip processing if already processed
									} else {
										name = res.Name
										if v == "BINSUS" {
											skip = v
											countStgLv += 14
										} else if v == "REGULER" {
											skip = v
											countStgLv += 14
										} else {
											skip = v
											countStgLv += 14
										}

										// Mark the value as processed
										// resStages, err := GetAllStudentStageClass("PTK", v)
										// if err != nil {
										// 	return nil, nil, err
										// }

										// for _, mdl := range resStages.Data {
										// 	if mdl.ModuleType == "PREMIUM_TRYOUT" {
										// 		countStgLv += 1
										// 	}
										// }
									}
								}
							}

						}

					}
				} else if filter == "all-module" {
					if req.TypeStages == "UMUM" {
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
								countStgLv += 1
							}
						}
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_PRE_UKA") {
								countStgLv += 1
							}
						}
					} else {
						for _, tag := range obj.ProductTags {
							re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
							matches := re.FindStringSubmatch(tag)

							// Check if there are captured substrings
							if len(matches) > 1 {

								// Iterate over captured substrings starting from index 1
								for _, v := range matches[1:] {

									// Convert both strings to lowercase for case-insensitive comparison
									length := (len(skip))

									var name string
									if len(v) == length && name == res.Name {
										continue // Skip processing if already processed
									} else {
										name = res.Name
										if v == "BINSUS" {
											skip = v
											countStgLv += 28
										} else if v == "REGULER" {
											skip = v
											countStgLv += 28
										} else {
											skip = v
											countStgLv += 28
										}
										// Mark the value as processed
										// resStages, err := GetAllStudentStageClass("PTK", v)
										// if err != nil {
										// 	return nil, nil, err
										// }

										// for _, mdl := range resStages.Data {
										// 	if mdl.ModuleType == "PREMIUM_TRYOUT" || mdl.ModuleType == "PLATINUM" {
										// 		countStgLv += 1
										// 	}
										// }
									}
								}
							}
						}
					}

				}
			} else {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_LEVEL_") {
						countStgLv += 1
					}
				}
			}
		}
		// } else {
		// 	if filter == "pre-uka" {
		// 		for _, mdl := range resStages.Data {
		// 			if mdl.ModuleType == "PLATINUM" {
		// 				countStgLv += 1
		// 			}
		// 		}
		// 	} else if filter == "challenge-uka" {
		// 		for _, mdl := range resStages.Data {
		// 			if mdl.ModuleType == "PREMIUM_TRYOUT" {
		// 				countStgLv += 1
		// 			}
		// 		}

		// 	} else {
		// 		for _, mdl := range resStages.Data {
		// 			if mdl.ModuleType == "PREMIUM_TRYOUT" || mdl.ModuleType == "PLATINUM" {
		// 				countStgLv += 1
		// 			}
		// 		}
		// 	}

		// }

		ownMod := withCode

		if (countStgLv - stgTotal) < 0 {
			ownMod += stgTotal
		} else {
			ownMod += countStgLv
		}

		donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)
		if _isNaNorInf(donePerct) {
			donePerct = 0
		}

		avgDonePercent := helpers.RoundFloat(float64(done)/float64(ownMod), 2)

		if _isNaNorInf(avgDonePercent) {
			avgDonePercent = 0
		}

		smr := request.Summary{
			AverageScore:  avgScore,
			Passed:        totalPassed,
			Failed:        totalFailed,
			Total:         totalHis,
			TotalScore:    float32(scoreTotal),
			PassedPercent: float32(passPercent),
			ScoreKeys:     scoreKeys,
			ScoreValues:   scVal,
			DonePercent:   float32(donePerct),
			Owned:         ownMod,
			Done:          float32(totalHis),
			AverageDone:   float32(avgDonePercent),
		}

		bc := "PT0000"
		bn := "Bimbel BTW (Kantor Pusat)"

		if res.BranchCode != nil {
			bc = *res.BranchCode
			bn = *res.BranchName

		}

		resScore := request.ResultsPerformaSiswa{
			BranchCode: bc,
			BranchName: bn,
			Name:       res.Name,
			SmartBtwID: res.SmartbtwID,
			Email:      res.Email,
			Summary:    smr,
			ClassInformation: requests.StudentClassInformation{
				JoinedClass:     true,
				ClassTitle:      clsData.Title,
				ClassYear:       int(clsData.Year),
				ClassJoined:     clsData.CreatedAt,
				ClassStatus:     clsData.Status,
				ClassBranchCode: clsData.BranchCode,
			},
			PTKTarget:     requests.StudentTargetDataPerforma{},
			HistoryRecord: hRes,
		}

		ctPercent := float64(0)

		if !_isNaNorInf(studentTarget.PassingRecommendationAvgPercentScore) {
			passProbability := calculatePassProbability(avgScore, studentTarget.TargetScore)
			pasPercentage := passProbability * 100
			ctPercent = math.Round(pasPercentage*10) / 10
		}

		if ctPercent > 100 {

			ctPercent = 100
		}

		if errStudentTarget == nil {
			resScore.PTKTarget = request.StudentTargetDataPerforma{
				SchoolID:                  studentTarget.SchoolID,
				MajorID:                   studentTarget.MajorID,
				SchoolName:                studentTarget.SchoolName,
				MajorName:                 studentTarget.MajorName,
				PolbitType:                studentTarget.PolbitType,
				PolbitCompetitionID:       studentTarget.PolbitCompetitionID,
				PolbitLocationID:          studentTarget.PolbitLocationID,
				TargetScore:               int(studentTarget.TargetScore),
				CurrentTargetPercentScore: ctPercent,
			}
		}

		bknScore := map[string]any{
			"twk":           float64(0),
			"twk_pass":      false,
			"tiu":           float64(0),
			"tiu_pass":      false,
			"tkp":           float64(0),
			"tkp_pass":      false,
			"total":         float64(0),
			"year":          2023,
			"bkn_attempted": false,
			"is_pass":       false,
		}

		finalPassingPercentage := float64(0)

		deviationExist := false
		twkDeviation := float64(0)
		twkPercentage := float64(0)
		tiuDeviation := float64(0)
		tiuPercentage := float64(0)
		tkpDeviation := float64(0)
		tkpPercentage := float64(0)
		totalDeviation := float64(0)
		totalPercentage := float64(0)
		bknTargetDeviation := float64(0)
		bknTargetPercentage := float64(0)
		ukaTargetDeviation := float64(0)
		ukaTargetPercentage := float64(0)

		ukaTargetDeviation = math.Round((avgScore-studentTarget.TargetScore)*10) / 10
		ukaTargetPercentage = math.Round((float64(ukaTargetDeviation) / float64(studentTarget.TargetScore)) * 100)

		if _isNaNorInf(ukaTargetDeviation) {
			ukaTargetDeviation = 0
		}
		if _isNaNorInf(ukaTargetPercentage) {
			ukaTargetPercentage = 0
		}

		bknSc, err := GetSingleBKNScoreByYearAndStudentUKA(res.SmartbtwID)
		if err == nil {
			isBknPass := false
			if bknSc.Twk >= 65 && bknSc.Tiu >= 80 && bknSc.Tkp >= 156 {
				isBknPass = true
			}

			bknScore = map[string]any{
				"twk":           bknSc.Twk,
				"twk_pass":      bknSc.Twk >= 65,
				"tiu":           bknSc.Tiu,
				"tiu_pass":      bknSc.Tiu >= 80,
				"tkp":           bknSc.Tkp,
				"tkp_pass":      bknSc.Tkp >= 156,
				"total":         bknSc.Total,
				"year":          bknSc.Year,
				"bkn_attempted": true,
				"is_pass":       isBknPass,
			}

			finalPassingPercentage = helpers.RoundFloat(float64(bknSc.Total)/float64(studentTarget.TargetScore)*100, 1)
			if _isNaNorInf(finalPassingPercentage) {
				finalPassingPercentage = 0
			}

			twkDeviation = math.Round((bknSc.Twk-avgScoreTWK)*10) / 10
			twkPercentage = math.Round((float64(twkDeviation) / float64(avgScoreTWK)) * 100)
			tiuDeviation = math.Round((bknSc.Tiu-avgScoreTIU)*10) / 10
			tiuPercentage = math.Round((float64(tiuDeviation) / float64(avgScoreTIU)) * 100)
			tkpDeviation = math.Round((bknSc.Tkp-avgScoreTKP)*10) / 10
			tkpPercentage = math.Round((float64(tkpDeviation) / float64(avgScoreTKP)) * 100)
			totalDeviation = math.Round((bknSc.Total-avgScore)*10) / 10
			totalPercentage = math.Round((float64(totalDeviation) / float64(avgScore)) * 100)
			bknTargetDeviation = math.Round((bknSc.Total-studentTarget.TargetScore)*10) / 10
			bknTargetPercentage = math.Round((float64(bknTargetDeviation) / float64(studentTarget.TargetScore)) * 100)
			if _isNaNorInf(twkDeviation) {
				twkDeviation = 0
			}
			if _isNaNorInf(twkPercentage) {
				twkPercentage = 0
			}
			if _isNaNorInf(tiuDeviation) {
				tiuDeviation = 0
			}
			if _isNaNorInf(tiuPercentage) {
				tiuPercentage = 0
			}
			if _isNaNorInf(tkpDeviation) {
				tkpDeviation = 0
			}
			if _isNaNorInf(tkpPercentage) {
				tkpPercentage = 0
			}
			if _isNaNorInf(totalDeviation) {
				totalDeviation = 0
			}
			if _isNaNorInf(totalPercentage) {
				totalPercentage = 0
			}
			if _isNaNorInf(bknTargetDeviation) {
				bknTargetDeviation = 0
			}
			if _isNaNorInf(bknTargetPercentage) {
				bknTargetPercentage = 0
			}
			deviationExist = true
		}

		targetPassingPercentage := helpers.RoundFloat(float64(avgScore)/float64(studentTarget.TargetScore)*100, 1)
		if _isNaNorInf(targetPassingPercentage) {
			targetPassingPercentage = 0
		}
		resScore.BKNScore = bknScore
		resScore.Summary.Deviation = map[string]any{
			"twk": map[string]any{
				"percentage":  twkPercentage,
				"differences": twkDeviation,
			},
			"tiu": map[string]any{
				"percentage":  tiuPercentage,
				"differences": tiuDeviation,
			},
			"tkp": map[string]any{
				"percentage":  tkpPercentage,
				"differences": tkpDeviation,
			},
			"total": map[string]any{
				"percentage":  totalPercentage,
				"differences": totalDeviation,
			},
			"bkn": map[string]any{
				"percentage":  bknTargetPercentage,
				"differences": bknTargetDeviation,
			},
			"uka": map[string]any{
				"percentage":  ukaTargetPercentage,
				"differences": ukaTargetDeviation,
			},
			"available": deviationExist,
		}

		resScore.Summary.FinalPassingPercent = finalPassingPercentage
		resScore.Summary.TargetPassingPercent = targetPassingPercentage

		avgAllTwk += float64(scVal.TWK.AverageScore)
		avgAllTiu += float64(scVal.TIU.AverageScore)
		avgAllTkp += float64(scVal.TKP.AverageScore)
		avgAllTotal += float64(resScore.Summary.AverageScore)

		results = append(results, &resScore)

	}

	twkAvg := helpers.RoundFloat(float64(avgAllTwk)/float64(stTotal), 1)
	tiuAvg := helpers.RoundFloat(float64(avgAllTiu)/float64(stTotal), 1)
	tkpAvg := helpers.RoundFloat(float64(avgAllTkp)/float64(stTotal), 1)
	totalAvg := helpers.RoundFloat(float64(avgAllTotal)/float64(stTotal), 1)
	if _isNaNorInf(twkAvg) {
		twkAvg = 0
	}
	if _isNaNorInf(tiuAvg) {
		tiuAvg = 0
	}
	if _isNaNorInf(tkpAvg) {
		tkpAvg = 0
	}
	if _isNaNorInf(totalAvg) {
		totalAvg = 0
	}
	avgData := map[string]any{
		"twk":           twkAvg,
		"tiu":           tiuAvg,
		"tkp":           tkpAvg,
		"total":         totalAvg,
		"student_total": stTotal,
	}

	return results, avgData, nil
}

func GetPerformaSiswaCPNS(req *request.GetPerformaSiswaUKA) ([]*request.ResultsPerformaSiswaCPNS, map[string]any, error) {

	stTotal := 0
	avgAllTwk := float64(0)
	avgAllTiu := float64(0)
	avgAllTkp := float64(0)
	avgAllTotal := float64(0)

	stdRes, err := GetStudentProfileByArrayOfSmartbtwID(req.SmartBtwID)
	if err != nil {
		return nil, nil, err
	}

	stTotal = len(stdRes)

	results := []*request.ResultsPerformaSiswaCPNS{}
	for _, res := range stdRes {

		var clsData models.ClassMemberElastic
		// stdClsList, errCls := GetStudentJoinedClassList(int32(res.SmartbtwID), nil, false)
		// if errCls != nil {
		// 	continue
		// }

		var filter string
		switch strings.ToUpper(req.TypeModule) {
		case "PRE_UKA":
			filter = "pre-uka"
		case "ALL_MODULE":
			filter = "all-module"
		case "UKA_STAGE":
			filter = "challenge-uka"
		}

		trialCount := 0
		hisRes, er := GetHistoryCPNSElasticPeforma(uint(res.SmartbtwID), req.TypeStages, filter)
		if er != nil {
			return nil, nil, er
		}

		// resPkg, err := GetHistoryPTKElasticTypeStage(uint(res.SmartbtwID), req.TypeStages)
		// if err != nil {
		// 	return nil, nil, err
		// }

		// fmt.Println(resPkg)

		// hisRes = append(hisRes, resPkg...)
		studentTarget, errStudentTarget := GetStudentProfileCPNSElastic(res.SmartbtwID)

		var scoreTotal float64
		var scoreTWK float64
		var scoreTIU float64
		var scoreTKP float64
		var avgScore float64
		totalPassed := int(0)
		totalFailed := int(0)
		twkPass := int(0)
		tiuPass := int(0)
		tkpPass := int(0)
		twkFailed := int(0)
		tiuFailed := int(0)
		tkpFailed := int(0)
		withCode := int(0)
		stgTotal := int(0)
		// isPassTwk := false
		// isPassTiu := false
		// isPassTkp := false
		// isPass := false

		hRes := []request.CreateHistoryCpns{}

		for _, t := range hisRes {
			if t.ExamName == "" {
				continue
			}
			// if strings.Contains(t.ExamName, "Post-Test") || strings.Contains(t.ExamName, "Pre-Test") {
			// 	continue
			// }
			if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
				continue
			}
			if t.ModuleType == "TESTING" || t.ModuleType == "TRIAL" {
				trialCount += 1
				continue
			}
			scoreTotal += t.Total
			scoreTWK += t.Twk
			scoreTIU += t.Tiu
			scoreTKP += t.Tkp
			if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
				totalPassed += 1
				// isPass = true
			}
			if t.Twk < t.TwkPass || t.Tiu < t.TiuPass || t.Tkp < t.TkpPass {
				totalFailed += 1
			}
			if t.Twk >= t.TwkPass {
				twkPass += 1
				// isPassTwk = true
			}
			if t.Tiu >= t.TiuPass {
				tiuPass += 1
				// isPassTiu = true
			}
			if t.Tkp >= t.TkpPass {
				tkpPass += 1
				// isPassTkp = true
			}
			if t.Twk < t.TwkPass {
				twkFailed += 1
			}
			if t.Tiu < t.TiuPass {
				tiuFailed += 1
			}
			if t.Tkp < t.TkpPass {
				tkpFailed += 1
			}
			if t.PackageType == "WITH_CODE" {
				withCode += 1
			} else {
				stgTotal += 1
			}
			hRes = append(hRes, t)
		}

		sort.SliceStable(hRes, func(i, j int) bool {
			return hRes[i].Start.After(hRes[j].End)
		})

		totalHis := len(hisRes) - trialCount
		done := len(hisRes) - trialCount
		avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
		avgScoreTWK := helpers.RoundFloat(scoreTWK/float64(totalHis), 1)
		avgScoreTIU := helpers.RoundFloat(scoreTIU/float64(totalHis), 1)
		avgScoreTKP := helpers.RoundFloat(scoreTKP/float64(totalHis), 1)

		passPercent := helpers.RoundFloat(float64(totalPassed)/float64(totalHis)*100, 1)
		passPercentTWK := helpers.RoundFloat(float64(twkPass)/float64(totalHis)*100, 1)
		passPercentTIU := helpers.RoundFloat(float64(tiuPass)/float64(totalHis)*100, 1)
		passPercentTKP := helpers.RoundFloat(float64(tkpPass)/float64(totalHis)*100, 1)

		if _isNaNorInf(avgScore) {
			avgScore = 0
		}
		if _isNaNorInf(avgScoreTWK) {
			avgScoreTWK = 0
		}
		if _isNaNorInf(avgScoreTIU) {
			avgScoreTIU = 0
		}
		if _isNaNorInf(avgScoreTKP) {
			avgScoreTKP = 0
		}

		if _isNaNorInf(passPercent) {
			passPercent = 0
		}
		if _isNaNorInf(passPercentTWK) {
			passPercentTWK = 0
		}
		if _isNaNorInf(passPercentTIU) {
			passPercentTIU = 0
		}
		if _isNaNorInf(passPercentTKP) {
			passPercentTKP = 0
		}
		failePercentTWK := 100 - passPercentTWK
		failePercentTIU := 100 - passPercentTIU
		failePercentTKP := 100 - passPercentTKP

		scoreKeys := []string{"TWK", "TIU", "TKP"}

		scVal := request.ScoreValues{
			IsPass: avgScoreTWK >= 65 && avgScoreTIU >= 80 && avgScoreTKP >= 156,
			TWK: request.Values{
				Total:         totalHis,
				Passed:        twkPass,
				Failed:        twkFailed,
				IsPass:        avgScoreTWK >= 65,
				TotalScore:    float32(scoreTWK),
				AverageScore:  float32(avgScoreTWK),
				PassedPercent: float32(passPercentTWK),
				FailedPercent: float32(failePercentTWK),
			},
			TIU: request.Values{
				Total:         totalHis,
				Passed:        tiuPass,
				Failed:        tiuFailed,
				IsPass:        avgScoreTIU >= 80,
				TotalScore:    float32(scoreTIU),
				AverageScore:  float32(avgScoreTIU),
				PassedPercent: float32(passPercentTIU),
				FailedPercent: float32(failePercentTIU),
			},
			TKP: request.Values{
				Total:         totalHis,
				Passed:        tkpPass,
				Failed:        tkpFailed,
				IsPass:        avgScoreTKP >= 156,
				TotalScore:    float32(scoreTKP),
				AverageScore:  float32(avgScoreTKP),
				PassedPercent: float32(passPercentTKP),
				FailedPercent: float32(failePercentTKP),
			},
		}

		prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
		if err != nil {
			return nil, nil, err
		}

		countStgLv := 0
		if req.TypeStages == "UMUM" {
			for _, obj := range prodRe.Data {
				isSkipped := false
				if strings.ToLower(obj.ProductProgram) != "skd" {
					continue
				}
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
						isSkipped = true
						continue
					}
					if strings.Contains(tag, "MATERIAL") {
						isSkipped = true
						continue
					}
					if strings.Contains(tag, "CPNS") {
						isSkipped = true
						continue
					}
				}
				if isSkipped {
					continue
				}

				if filter != "" {
					if filter == "pre-uka" {
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_PRE_UKA") {
								countStgLv += 1
							}
						}
					} else if filter == "challenge-uka" {
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
								countStgLv += 1
							}
						}
					} else if filter == "all-module" {
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
								countStgLv += 1
							}
						}
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_PRE_UKA") {
								countStgLv += 1
							}
						}
					} else {
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_LEVEL_") && obj.ProductProgram == "skd" {
								countStgLv += 1
							}
						}
					}
				}
			}
		} else {
			countStgLv = 0
		}

		ownMod := withCode

		if (countStgLv - stgTotal) < 0 {
			ownMod += stgTotal
		} else {
			ownMod += countStgLv
		}

		donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)
		if _isNaNorInf(donePerct) {
			donePerct = 0
		}

		avgDonePercent := helpers.RoundFloat(float64(done)/float64(ownMod), 2)

		if _isNaNorInf(avgDonePercent) {
			avgDonePercent = 0
		}

		smr := request.Summary{
			AverageScore:  avgScore,
			Passed:        totalPassed,
			Failed:        totalFailed,
			Total:         totalHis,
			TotalScore:    float32(scoreTotal),
			PassedPercent: float32(passPercent),
			ScoreKeys:     scoreKeys,
			ScoreValues:   scVal,
			DonePercent:   float32(donePerct),
			Owned:         ownMod,
			Done:          float32(totalHis),
			AverageDone:   float32(avgDonePercent),
		}

		bc := "PT0000"
		bn := "Bimbel BTW (Kantor Pusat)"

		if res.BranchCode != nil {
			bc = *res.BranchCode
			bn = *res.BranchName

		}

		resScore := request.ResultsPerformaSiswaCPNS{
			BranchCode: bc,
			BranchName: bn,
			Name:       res.Name,
			SmartBtwID: res.SmartbtwID,
			Email:      res.Email,
			Summary:    smr,
			ClassInformation: requests.StudentClassInformation{
				JoinedClass:     true,
				ClassTitle:      clsData.Title,
				ClassYear:       int(clsData.Year),
				ClassJoined:     clsData.CreatedAt,
				ClassStatus:     clsData.Status,
				ClassBranchCode: clsData.BranchCode,
			},
			CPNSTarget:    requests.StudentTargetDataPerformaCPNS{},
			HistoryRecord: hRes,
		}

		ctPercent := float64(0)

		if !_isNaNorInf(studentTarget.PassingRecommendationAvgPercentScore) {
			passProbability := calculatePassProbability(avgScore, studentTarget.TargetScore)
			pasPercentage := passProbability * 100
			ctPercent = math.Round(pasPercentage*10) / 10
		}

		if ctPercent > 100 {
			ctPercent = 100
		}

		if errStudentTarget == nil {
			resScore.CPNSTarget = request.StudentTargetDataPerformaCPNS{
				SchoolID:                  studentTarget.InstanceID,
				MajorID:                   studentTarget.PositionID,
				SchoolName:                studentTarget.InstanceName,
				MajorName:                 studentTarget.PositionName,
				TargetScore:               int(studentTarget.TargetScore),
				CurrentTargetPercentScore: ctPercent,
			}
		}

		bknScore := map[string]any{
			"twk":           float64(0),
			"twk_pass":      false,
			"tiu":           float64(0),
			"tiu_pass":      false,
			"tkp":           float64(0),
			"tkp_pass":      false,
			"total":         float64(0),
			"year":          2023,
			"bkn_attempted": false,
			"is_pass":       false,
		}

		finalPassingPercentage := float64(0)

		deviationExist := false
		twkDeviation := float64(0)
		twkPercentage := float64(0)
		tiuDeviation := float64(0)
		tiuPercentage := float64(0)
		tkpDeviation := float64(0)
		tkpPercentage := float64(0)
		totalDeviation := float64(0)
		totalPercentage := float64(0)
		bknTargetDeviation := float64(0)
		bknTargetPercentage := float64(0)
		ukaTargetDeviation := float64(0)
		ukaTargetPercentage := float64(0)

		ukaTargetDeviation = math.Round((avgScore-studentTarget.TargetScore)*10) / 10
		ukaTargetPercentage = math.Round((float64(ukaTargetDeviation) / float64(studentTarget.TargetScore)) * 100)

		if _isNaNorInf(ukaTargetDeviation) {
			ukaTargetDeviation = 0
		}
		if _isNaNorInf(ukaTargetPercentage) {
			ukaTargetPercentage = 0
		}

		bknSc, err := GetSingleBKNScoreByYearAndStudentUKA(res.SmartbtwID)
		if err == nil {
			isBknPass := false
			if bknSc.Twk >= 65 && bknSc.Tiu >= 80 && bknSc.Tkp >= 156 {
				isBknPass = true
			}

			bknScore = map[string]any{
				"twk":           bknSc.Twk,
				"twk_pass":      bknSc.Twk >= 65,
				"tiu":           bknSc.Tiu,
				"tiu_pass":      bknSc.Tiu >= 80,
				"tkp":           bknSc.Tkp,
				"tkp_pass":      bknSc.Tkp >= 156,
				"total":         bknSc.Total,
				"year":          bknSc.Year,
				"bkn_attempted": true,
				"is_pass":       isBknPass,
			}

			finalPassingPercentage = helpers.RoundFloat(float64(bknSc.Total)/float64(studentTarget.TargetScore)*100, 1)
			if _isNaNorInf(finalPassingPercentage) {
				finalPassingPercentage = 0
			}

			twkDeviation = math.Round((bknSc.Twk-avgScoreTWK)*10) / 10
			twkPercentage = math.Round((float64(twkDeviation) / float64(avgScoreTWK)) * 100)
			tiuDeviation = math.Round((bknSc.Tiu-avgScoreTIU)*10) / 10
			tiuPercentage = math.Round((float64(tiuDeviation) / float64(avgScoreTIU)) * 100)
			tkpDeviation = math.Round((bknSc.Tkp-avgScoreTKP)*10) / 10
			tkpPercentage = math.Round((float64(tkpDeviation) / float64(avgScoreTKP)) * 100)
			totalDeviation = math.Round((bknSc.Total-avgScore)*10) / 10
			totalPercentage = math.Round((float64(totalDeviation) / float64(avgScore)) * 100)
			bknTargetDeviation = math.Round((bknSc.Total-studentTarget.TargetScore)*10) / 10
			bknTargetPercentage = math.Round((float64(bknTargetDeviation) / float64(studentTarget.TargetScore)) * 100)
			if _isNaNorInf(twkDeviation) {
				twkDeviation = 0
			}
			if _isNaNorInf(twkPercentage) {
				twkPercentage = 0
			}
			if _isNaNorInf(tiuDeviation) {
				tiuDeviation = 0
			}
			if _isNaNorInf(tiuPercentage) {
				tiuPercentage = 0
			}
			if _isNaNorInf(tkpDeviation) {
				tkpDeviation = 0
			}
			if _isNaNorInf(tkpPercentage) {
				tkpPercentage = 0
			}
			if _isNaNorInf(totalDeviation) {
				totalDeviation = 0
			}
			if _isNaNorInf(totalPercentage) {
				totalPercentage = 0
			}
			if _isNaNorInf(bknTargetDeviation) {
				bknTargetDeviation = 0
			}
			if _isNaNorInf(bknTargetPercentage) {
				bknTargetPercentage = 0
			}
			deviationExist = true
		}

		targetPassingPercentage := helpers.RoundFloat(float64(avgScore)/float64(studentTarget.TargetScore)*100, 1)
		if _isNaNorInf(targetPassingPercentage) {
			targetPassingPercentage = 0
		}
		resScore.BKNScore = bknScore
		resScore.Summary.Deviation = map[string]any{
			"twk": map[string]any{
				"percentage":  twkPercentage,
				"differences": twkDeviation,
			},
			"tiu": map[string]any{
				"percentage":  tiuPercentage,
				"differences": tiuDeviation,
			},
			"tkp": map[string]any{
				"percentage":  tkpPercentage,
				"differences": tkpDeviation,
			},
			"total": map[string]any{
				"percentage":  totalPercentage,
				"differences": totalDeviation,
			},
			"bkn": map[string]any{
				"percentage":  bknTargetPercentage,
				"differences": bknTargetDeviation,
			},
			"uka": map[string]any{
				"percentage":  ukaTargetPercentage,
				"differences": ukaTargetDeviation,
			},
			"available": deviationExist,
		}

		resScore.Summary.FinalPassingPercent = finalPassingPercentage
		resScore.Summary.TargetPassingPercent = targetPassingPercentage

		avgAllTwk += float64(scVal.TWK.AverageScore)
		avgAllTiu += float64(scVal.TIU.AverageScore)
		avgAllTkp += float64(scVal.TKP.AverageScore)
		avgAllTotal += float64(resScore.Summary.AverageScore)

		results = append(results, &resScore)

	}

	twkAvg := helpers.RoundFloat(float64(avgAllTwk)/float64(stTotal), 1)
	tiuAvg := helpers.RoundFloat(float64(avgAllTiu)/float64(stTotal), 1)
	tkpAvg := helpers.RoundFloat(float64(avgAllTkp)/float64(stTotal), 1)
	totalAvg := helpers.RoundFloat(float64(avgAllTotal)/float64(stTotal), 1)
	if _isNaNorInf(twkAvg) {
		twkAvg = 0
	}
	if _isNaNorInf(tiuAvg) {
		tiuAvg = 0
	}
	if _isNaNorInf(tkpAvg) {
		tkpAvg = 0
	}
	if _isNaNorInf(totalAvg) {
		totalAvg = 0
	}
	avgData := map[string]any{
		"twk":           twkAvg,
		"tiu":           tiuAvg,
		"tkp":           tkpAvg,
		"total":         totalAvg,
		"student_total": stTotal,
	}

	return results, avgData, nil
}

func GetPerformaSiswaPTN(req *request.GetPerformaSiswaUKA) ([]*request.ResultsPerformaSiswaPTN, map[string]any, error) {

	stTotal := 0
	avgAllPU := float64(0)
	avgAllPPU := float64(0)
	avgAllPBM := float64(0)
	avgAllPK := float64(0)
	avgAllLBINDO := float64(0)
	avgAllLBING := float64(0)
	avgAllPM := float64(0)
	avgAllTotal := float64(0)
	skip := ""

	stdRes, err := GetStudentProfileByArrayOfSmartbtwID(req.SmartBtwID)
	if err != nil {
		return nil, nil, err
	}

	stTotal = len(stdRes)

	results := []*request.ResultsPerformaSiswaPTN{}
	for _, res := range stdRes {

		var clsData models.ClassMemberElastic
		// stdClsList, errCls := GetStudentJoinedClassList(int32(res.SmartbtwID), nil, false)
		// if errCls != nil {
		// 	continue
		// }

		var filter string
		switch strings.ToUpper(req.TypeModule) {
		case "PRE_UKA":
			filter = "pre-uka"
		case "ALL_MODULE":
			filter = "all-module"
		case "UKA_STAGE":
			filter = "challenge-uka"
		}

		trialCount := 0
		hisRes, er := GetHistoryPTNElasticPeforma(uint(res.SmartbtwID), req.TypeStages, filter)
		if er != nil {
			return nil, nil, er
		}

		// resPkg, err := GetHistoryPTKElasticTypeStage(uint(res.SmartbtwID), req.TypeStages)
		// if err != nil {
		// 	return nil, nil, err
		// }

		// fmt.Println(resPkg)

		// hisRes = append(hisRes, resPkg...)
		studentTarget, errStudentTarget := GetStudentProfilePTNElastic(res.SmartbtwID)

		var scoreTotal float64
		var scorePU float64
		var scorePPU float64
		var scorePBM float64
		var scorePK float64
		var scoreLBIND float64
		var scoreLBING float64
		var scorePM float64
		var avgScore float64
		totalPassed := int(0)
		totalFailed := int(0)
		withCode := int(0)
		stgTotal := int(0)

		hRes := []request.CreateHistoryPtn{}

		for _, t := range hisRes {
			// if strings.Contains(t.ExamName, "Post-Test") || strings.Contains(t.ExamName, "Pre-Test") {
			// 	continue
			// }
			if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
				continue
			}
			fmt.Println(t.PenalaranUmum)
			scorePU += t.PenalaranUmum
			scorePPU += t.PengetahuanUmum
			scorePBM += t.PemahamanBacaan
			scorePK += t.PengetahuanKuantitatif
			scoreLBIND += t.LiterasiBahasaIndonesia
			scoreLBING += t.LiterasiBahasaInggris
			scorePM += t.PenalaranMatematika
			if t.ExamName == "" {
				continue
			}
			// if t.ModuleType == "TESTING" || t.ModuleType == "TRIAL" {
			// 	trialCount += 1
			// 	continue
			// }
			scoreTotal += t.Total
			if t.PackageType == "WITH_CODE" {
				withCode += 1
			} else {
				stgTotal += 1
			}
			hRes = append(hRes, t)
		}

		sort.SliceStable(hRes, func(i, j int) bool {
			return hRes[i].Start.After(hRes[j].End)
		})

		totalHis := len(hisRes) - trialCount
		fmt.Println("ini total his", float64(totalHis))
		done := len(hisRes) - trialCount
		avgScore = helpers.RoundFloats(scoreTotal/float64(totalHis), 1)
		avgScorePU := helpers.RoundFloats(scorePU/float64(totalHis), 1)
		avgScorePPU := helpers.RoundFloats(scorePPU/float64(totalHis), 1)
		avgScorePBM := helpers.RoundFloats(scorePBM/float64(totalHis), 1)
		avgScorePK := helpers.RoundFloats(scorePK/float64(totalHis), 1)
		avgScoreLBIND := helpers.RoundFloats(scoreLBIND/float64(totalHis), 1)
		avgScoreLBING := helpers.RoundFloats(scoreLBING/float64(totalHis), 1)
		avgScorePM := helpers.RoundFloats(scorePM/float64(totalHis), 1)
		fmt.Println("skor PU sebelum dirata ratakan: ", avgScore)

		// passPercent := helpers.RoundFloats(float64(totalPassed)/float64(totalHis)*100, 1)
		// passPercentTWK := helpers.RoundFloats(float64(twkPass)/float64(totalHis)*100, 1)
		// passPercentTIU := helpers.RoundFloats(float64(tiuPass)/float64(totalHis)*100, 1)
		// passPercentTKP := helpers.RoundFloats(float64(tkpPass)/float64(totalHis)*100, 1)

		if _isNaNorInf(avgScore) {
			avgScore = 0
		}
		if _isNaNorInf(avgScorePU) {
			avgScorePU = 0
		}
		if _isNaNorInf(avgScorePPU) {
			avgScorePPU = 0
		}
		if _isNaNorInf(avgScorePBM) {
			avgScorePBM = 0
		}
		if _isNaNorInf(avgScorePK) {
			avgScorePK = 0
		}
		if _isNaNorInf(avgScoreLBIND) {
			avgScoreLBIND = 0
		}
		if _isNaNorInf(avgScoreLBING) {
			avgScoreLBING = 0
		}
		if _isNaNorInf(avgScorePM) {
			avgScorePM = 0
		}

		// failePercentTWK := 100 - passPercentTWK
		// failePercentTIU := 100 - passPercentTIU
		// failePercentTKP := 100 - passPercentTKP

		scoreKeys := []string{"PU", "PPU", "PBM", "PK", "BIND", "BING", "PM"}

		scVal := request.ScoreValuesPTN{
			PU: request.Values{
				Total:        totalHis,
				TotalScore:   float32(scorePU),
				AverageScore: float32(avgScorePU),
			},
			PPU: request.Values{
				Total:        totalHis,
				TotalScore:   float32(scorePPU),
				AverageScore: float32(avgScorePPU),
			},
			PBM: request.Values{
				Total:        totalHis,
				TotalScore:   float32(scorePBM),
				AverageScore: float32(avgScorePBM),
			},
			PK: request.Values{
				Total:        totalHis,
				TotalScore:   float32(scorePK),
				AverageScore: float32(avgScorePK),
			},
			LBIND: request.Values{
				Total:        totalHis,
				TotalScore:   float32(scoreLBIND),
				AverageScore: float32(avgScoreLBIND),
			},
			LBING: request.Values{
				Total:        totalHis,
				TotalScore:   float32(scoreLBING),
				AverageScore: float32(avgScoreLBING),
			},
			PM: request.Values{
				Total:        totalHis,
				TotalScore:   float32(scorePM),
				AverageScore: float32(avgScorePM),
			},
		}

		prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
		if err != nil {
			return nil, nil, err
		}

		// resStages, err := GetAllStudentStageClass("PTN")
		// if err != nil {
		// 	return nil, nil, err
		// }

		countStgLv := 0
		for _, obj := range prodRe.Data {
			isSkipped := false
			if strings.ToLower(obj.ProductProgram) != "tps" {
				continue
			}
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
					isSkipped = true
					continue
				}
				if strings.Contains(tag, "MATERIAL") {
					isSkipped = true
					continue
				}
				if strings.Contains(tag, "CPNS") {
					isSkipped = true
					continue
				}
			}
			if isSkipped {
				continue
			}
			if filter != "" {
				if filter == "pre-uka" {
					for _, tag := range obj.ProductTags {
						if req.TypeStages == "UMUM" {
							if strings.Contains(tag, "STAGE_PRE_UKA") {
								countStgLv += 1
							}
						} else {
							re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
							matches := re.FindStringSubmatch(tag)

							// Check if there are captured substrings
							if len(matches) > 1 {

								// Iterate over captured substrings starting from index 1
								for _, v := range matches[1:] {

									// Convert both strings to lowercase for case-insensitive comparison
									length := (len(skip))
									skip = v
									var name string
									if len(v) == length && name == res.Name {
										continue // Skip processing if already processed
									} else {
										name = res.Name
										if v == "BINSUS" {
											countStgLv += 14
										} else if v == "REGULER" {
											countStgLv += 14
										} else {
											countStgLv += 14
										}

										// Mark the value as processed
										// resStages, err := GetAllStudentStageClass("PTN", v)
										// if err != nil {
										// 	return nil, nil, err
										// }

										// for _, mdl := range resStages.Data {
										// 	if mdl.ModuleType == "PLATINUM" {
										// 		countStgLv += 1
										// 	}
										// }
									}
								}
							}

						}

					}
				} else if filter == "challenge-uka" {
					for _, tag := range obj.ProductTags {
						if req.TypeStages == "UMUM" {

							if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
								countStgLv += 1
							}
						} else {
							re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
							matches := re.FindStringSubmatch(tag)

							// Check if there are captured substrings
							if len(matches) > 1 {

								// Iterate over captured substrings starting from index 1
								for _, v := range matches[1:] {

									// Convert both strings to lowercase for case-insensitive comparison
									length := (len(skip))
									skip = v
									var name string
									if len(v) == length && name == res.Name {
										continue // Skip processing if already processed
									} else {
										name = res.Name
										if v == "BINSUS" {
											countStgLv += 14
										} else if v == "REGULER" {
											countStgLv += 14
										} else {
											countStgLv += 14
										}

										// Mark the value as processed
										// resStages, err := GetAllStudentStageClass("PTN", v)
										// if err != nil {
										// 	return nil, nil, err
										// }

										// for _, mdl := range resStages.Data {
										// 	if mdl.ModuleType == "PREMIUM_TRYOUT" {
										// 		countStgLv += 1
										// 	}
										// }
									}
								}
							}

						}

					}
				} else if filter == "all-module" {
					if req.TypeStages == "UMUM" {
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
								countStgLv += 1
							}
						}
						for _, tag := range obj.ProductTags {
							if strings.Contains(tag, "STAGE_PRE_UKA") {
								countStgLv += 1
							}
						}
					} else {
						for _, tag := range obj.ProductTags {
							re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
							matches := re.FindStringSubmatch(tag)

							// Check if there are captured substrings
							if len(matches) > 1 {

								// Iterate over captured substrings starting from index 1
								for _, v := range matches[1:] {

									// Convert both strings to lowercase for case-insensitive comparison
									length := (len(skip))
									skip = v
									var name string
									if len(v) == length && name == res.Name {
										continue // Skip processing if already processed
									} else {
										name = res.Name
										if v == "BINSUS" {
											countStgLv += 28
										} else if v == "REGULER" {
											countStgLv += 28
										} else {
											countStgLv += 28
										}

										// Mark the value as processed
										// resStages, err := GetAllStudentStageClass("PTN", v)
										// if err != nil {
										// 	return nil, nil, err
										// }

										// for _, mdl := range resStages.Data {
										// 	if mdl.ModuleType == "PLATINUM" || mdl.ModuleType == "PREMIUM_TRYOUT" {
										// 		countStgLv += 1
										// 	}
										// }
									}
								}
							}
						}

					}

				}
			} else {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_LEVEL_") {
						countStgLv += 1
					}
				}
			}
		}
		// } else {
		// 	if filter == "pre-uka" {
		// 		for _, mdl := range resStages.Data {
		// 			if mdl.ModuleType == "PLATINUM" {
		// 				countStgLv += 1
		// 			}
		// 		}
		// 	} else if filter == "challenge-uka" {
		// 		for _, mdl := range resStages.Data {
		// 			if mdl.ModuleType == "PREMIUM_TRYOUT" {
		// 				countStgLv += 1
		// 			}
		// 		}

		// 	} else {
		// 		for _, mdl := range resStages.Data {
		// 			if mdl.ModuleType == "PREMIUM_TRYOUT" || mdl.ModuleType == "PLATINUM" {
		// 				countStgLv += 1
		// 			}
		// 		}
		// 	}

		// }

		ownMod := withCode
		fmt.Println(ownMod)
		if (countStgLv - stgTotal) < 0 {
			ownMod += stgTotal
		} else {
			ownMod += countStgLv
		}

		donePerct := helpers.RoundFloats(float64(totalHis)/float64(ownMod)*100, 1)
		if _isNaNorInf(donePerct) {
			donePerct = 0
		}

		avgDonePercent := helpers.RoundFloats(float64(done)/float64(ownMod), 2)

		if _isNaNorInf(avgDonePercent) {
			avgDonePercent = 0
		}

		smr := request.SummaryPTN{
			AverageScore: avgScore,
			Passed:       totalPassed,
			Failed:       totalFailed,
			Total:        totalHis,
			TotalScore:   float32(scoreTotal),
			ScoreKeys:    scoreKeys,
			DonePercent:  float32(donePerct),
			ScoreValues:  scVal,
			Owned:        ownMod,
			Done:         float32(totalHis),
			AverageDone:  float32(avgDonePercent),
		}

		bc := "PT0000"
		bn := "Bimbel BTW (Kantor Pusat)"

		if res.BranchCode != nil {
			bc = *res.BranchCode
			bn = *res.BranchName

		}

		resScore := request.ResultsPerformaSiswaPTN{
			BranchCode: bc,
			BranchName: bn,
			Name:       res.Name,
			SmartBtwID: res.SmartbtwID,
			Email:      res.Email,
			Summary:    smr,
			ClassInformation: requests.StudentClassInformation{
				JoinedClass:     true,
				ClassTitle:      clsData.Title,
				ClassYear:       int(clsData.Year),
				ClassJoined:     clsData.CreatedAt,
				ClassStatus:     clsData.Status,
				ClassBranchCode: clsData.BranchCode,
			},
			PTNTarget:     requests.StudentTargetDataPerforma{},
			HistoryRecord: hRes,
		}

		ctPercent := float64(0)

		if !_isNaNorInf(studentTarget.PassingRecommendationAvgPercentScore) {
			passProbability := calculatePassProbability(avgScore, studentTarget.TargetScore)
			pasPercentage := passProbability * 100
			ctPercent = math.Round(pasPercentage*10) / 10
		}

		if ctPercent > 100 {
			ctPercent = 100
		}

		if errStudentTarget == nil {
			resScore.PTNTarget = request.StudentTargetDataPerforma{
				SchoolID:                  studentTarget.SchoolID,
				MajorID:                   studentTarget.MajorID,
				SchoolName:                studentTarget.SchoolName,
				MajorName:                 studentTarget.MajorName,
				TargetScore:               int(studentTarget.TargetScore),
				CurrentTargetPercentScore: ctPercent,
			}
		}

		// bknScore := map[string]any{
		// 	"pu":    float64(0),
		// 	"ppu":   float64(0),
		// 	"pbm":   float64(0),
		// 	"pk":    float64(0),
		// 	"bindo": float64(0),
		// 	"bing":  float64(0),
		// 	"pm":    float64(0),
		// 	"total": float64(0),
		// 	"year":  2023,
		// }

		finalPassingPercentage := float64(0)

		// deviationExist := false
		// puDeviation := float64(0)
		// puPercentage := float64(0)
		// ppuDeviation := float64(0)
		// ppuPercentage := float64(0)
		// pbmDeviation := float64(0)
		// pbmPercentage := float64(0)
		// pkDeviation := float64(0)
		// pkPercentage := float64(0)
		// lbindDeviation := float64(0)
		// lbindPercentage := float64(0)
		// lbingDeviation := float64(0)
		// lbingPercentage := float64(0)
		// pmDeviation := float64(0)
		// pmPercentage := float64(0)
		// totalDeviation := float64(0)
		// totalPercentage := float64(0)
		// bknTargetDeviation := float64(0)
		// bknTargetPercentage := float64(0)
		ukaTargetDeviation := float64(0)
		ukaTargetPercentage := float64(0)

		ukaTargetDeviation = math.Round((avgScore-studentTarget.TargetScore)*10) / 10
		ukaTargetPercentage = math.Round((float64(ukaTargetDeviation) / float64(studentTarget.TargetScore)) * 100)

		if _isNaNorInf(ukaTargetDeviation) {
			ukaTargetDeviation = 0
		}
		if _isNaNorInf(ukaTargetPercentage) {
			ukaTargetPercentage = 0
		}

		// bknSc, err := GetSingleBKNScoreByYearAndStudentUKA(res.SmartbtwID)
		// if err == nil {
		// 	isBknPass := false
		// 	if bknSc.Twk >= 65 && bknSc.Tiu >= 80 && bknSc.Tkp >= 156 {
		// 		isBknPass = true
		// 	}

		// 	bknScore = map[string]any{
		// 		"twk":           bknSc.Twk,
		// 		"tiu":           bknSc.Tiu,
		// 		"tkp":           bknSc.Tkp,
		// 		"total":         bknSc.Total,
		// 		"year":          bknSc.Year,
		// 		"bkn_attempted": true,
		// 		"is_pass":       isBknPass,
		// 	}

		// 	finalPassingPercentage = helpers.RoundFloats(float64(bknSc.Total)/float64(studentTarget.TargetScore)*100, 1)
		// 	if _isNaNorInf(finalPassingPercentage) {
		// 		finalPassingPercentage = 0
		// 	}

		// 	twkDeviation = math.Round((bknSc.Twk-avgScorePU)*10) / 10
		// 	twkPercentage = math.Round((float64(twskDeviation) / float64(avgScorePU)) * 100)
		// 	tiuDeviation = math.Round((bknSc.Tiu-avgScoreTIU)*10) / 10
		// 	tiuPercentage = math.Round((float64(tiuDeviation) / float64(avgScoreTIU)) * 100)
		// 	tkpDeviation = math.Round((bknSc.Tkp-avgScoreTKP)*10) / 10
		// 	tkpPercentage = math.Round((float64(tkpDeviation) / float64(avgScoreTKP)) * 100)
		// 	totalDeviation = math.Round((bknSc.Total-avgScore)*10) / 10
		// 	totalPercentage = math.Round((float64(totalDeviation) / float64(avgScore)) * 100)
		// 	bknTargetDeviation = math.Round((bknSc.Total-studentTarget.TargetScore)*10) / 10
		// 	bknTargetPercentage = math.Round((float64(bknTargetDeviation) / float64(studentTarget.TargetScore)) * 100)
		// 	if _isNaNorInf(twkDeviation) {
		// 		twkDeviation = 0
		// 	}
		// 	if _isNaNorInf(twkPercentage) {
		// 		twkPercentage = 0
		// 	}
		// 	if _isNaNorInf(tiuDeviation) {
		// 		tiuDeviation = 0
		// 	}
		// 	if _isNaNorInf(tiuPercentage) {
		// 		tiuPercentage = 0
		// 	}
		// 	if _isNaNorInf(tkpDeviation) {
		// 		tkpDeviation = 0
		// 	}
		// 	if _isNaNorInf(tkpPercentage) {
		// 		tkpPercentage = 0
		// 	}
		// 	if _isNaNorInf(totalDeviation) {
		// 		totalDeviation = 0
		// 	}
		// 	if _isNaNorInf(totalPercentage) {
		// 		totalPercentage = 0
		// 	}
		// 	if _isNaNorInf(bknTargetDeviation) {
		// 		bknTargetDeviation = 0
		// 	}
		// 	if _isNaNorInf(bknTargetPercentage) {
		// 		bknTargetPercentage = 0
		// 	}
		// 	deviationExist = true
		// }

		targetPassingPercentage := helpers.RoundFloats(float64(avgScore)/float64(studentTarget.TargetScore)*100, 1)
		if _isNaNorInf(targetPassingPercentage) {
			targetPassingPercentage = 0
		}
		// resScore.BKNScore = bknScore
		// resScore.Summary.Deviation = map[string]any{
		// 	"twk": map[string]any{
		// 		"percentage":  twkPercentage,
		// 		"differences": twkDeviation,
		// 	},
		// 	"tiu": map[string]any{
		// 		"percentage":  tiuPercentage,
		// 		"differences": tiuDeviation,
		// 	},
		// 	"tkp": map[string]any{
		// 		"percentage":  tkpPercentage,
		// 		"differences": tkpDeviation,
		// 	},
		// 	"total": map[string]any{
		// 		"percentage":  totalPercentage,
		// 		"differences": totalDeviation,
		// 	},
		// 	"bkn": map[string]any{
		// 		"percentage":  bknTargetPercentage,
		// 		"differences": bknTargetDeviation,
		// 	},
		// 	"uka": map[string]any{
		// 		"percentage":  ukaTargetPercentage,
		// 		"differences": ukaTargetDeviation,
		// 	},
		// 	"available": deviationExist,
		// }

		resScore.Summary.FinalPassingPercent = finalPassingPercentage
		resScore.Summary.TargetPassingPercent = targetPassingPercentage

		avgAllPU += float64(scVal.PU.AverageScore)
		avgAllPPU += float64(scVal.PPU.AverageScore)
		avgAllPBM += float64(scVal.PBM.AverageScore)
		avgAllPK += float64(scVal.PK.AverageScore)
		avgAllLBINDO += float64(scVal.LBIND.AverageScore)
		avgAllLBING += float64(scVal.LBING.AverageScore)
		avgAllPM += float64(scVal.PM.AverageScore)
		avgAllTotal += float64(resScore.Summary.AverageScore)

		results = append(results, &resScore)

	}

	avgPU := helpers.RoundFloats(float64(avgAllPU)/float64(stTotal), 1)
	avgPPU := helpers.RoundFloats(float64(avgAllPPU)/float64(stTotal), 1)
	avgPBM := helpers.RoundFloats(float64(avgAllPBM)/float64(stTotal), 1)
	avgPK := helpers.RoundFloats(float64(avgAllPK)/float64(stTotal), 1)
	avgLBINDO := helpers.RoundFloats(float64(avgAllLBINDO)/float64(stTotal), 1)
	avgLBING := helpers.RoundFloats(float64(avgAllLBING)/float64(stTotal), 1)
	avgPM := helpers.RoundFloats(float64(avgAllPM)/float64(stTotal), 1)
	totalAvg := helpers.RoundFloats(float64(avgAllTotal)/float64(stTotal), 1)
	if _isNaNorInf(avgPU) {
		avgPU = 0
	}
	if _isNaNorInf(avgPPU) {
		avgPPU = 0
	}
	if _isNaNorInf(avgPBM) {
		avgPBM = 0
	}
	if _isNaNorInf(avgPK) {
		avgPK = 0
	}
	if _isNaNorInf(avgLBINDO) {
		avgLBINDO = 0
	}
	if _isNaNorInf(avgLBING) {
		avgLBING = 0
	}
	if _isNaNorInf(avgPM) {
		avgPM = 0
	}
	if _isNaNorInf(totalAvg) {
		totalAvg = 0
	}
	avgData := map[string]any{
		"PU":            avgPU,
		"PPU":           avgPPU,
		"PBM":           avgPBM,
		"PK":            avgPK,
		"LBINDO":        avgLBINDO,
		"LBING":         avgLBING,
		"PM":            avgPM,
		"total":         totalAvg,
		"student_total": stTotal,
	}

	return results, avgData, nil
}

func GetStudentAOP(smID uint) (*requests.StudentAOPResults, error) {
	conn := os.Getenv("SERVICE_NEW_PRODUCT_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/aop/elastic/get-per-student?smartbtw_id=%d&status=1", smID), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to product " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to product " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of product " + err.Error())
	}

	st := requests.StudentAOPResults{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of product " + errs.Error())
	}

	return &st, nil
}

func GetStudentAOPNew(smID uint) (*requests.StudentAOPResults, error) {
	conn := os.Getenv("SERVICE_NEW_PRODUCT_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/aop/elastic/get-per-student/new?smartbtw_id=%d&status=1", smID), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to product " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to product " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of product " + err.Error())
	}

	st := requests.StudentAOPResults{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of product " + errs.Error())
	}

	return &st, nil
}

func FetchPTKUKARaport(smId uint, typStg, fil string) (map[string]any, error) {

	res, err := GetStudentProfileElastic(int(smId))
	if err != nil {
		return nil, err
	}
	var studentData = map[string]any{}
	var reportData = map[string]any{}
	// var studentTarget requests.StudentProfilePtkElastic
	var resBinsusReport *mockstruct.BinsusScreeningSummary
	// var binsusResult map[string]any

	// yrNow := time.Now().Year()
	studentTarget, err := GetStudentProfilePTKElastic(res.SmartbtwID)
	if err != nil {
		return nil, err
	}
	locId := 0
	if studentTarget.PolbitLocationID != nil {
		locId = *studentTarget.PolbitLocationID
	}
	chRes, err := GetCompetitionDataPTK(uint(studentTarget.MajorID), uint(locId), res.Gender, studentTarget.PolbitType)
	if err != nil {
		return nil, err
	}
	gender := "Perempuan"
	if res.Gender == "L" {
		gender = "Laki-laki"
	}

	targetFormat := "02/01/2006"
	formattedTime := res.BirthDate.Format(targetFormat)

	formation := "PUSAT"
	if res.PolbitTypePTK == "DAERAH" {
		formation = res.DomicileProvince
	}

	studentData = map[string]any{
		"name":                  res.Name,
		"email":                 res.Email,
		"program":               "PTK",
		"school_name":           studentTarget.SchoolName,
		"major_name":            studentTarget.MajorName,
		"target_score":          studentTarget.TargetScore,
		"quota":                 chRes.MajorQuota,
		"raport_date":           time.Now().Format("2006-01-02"),
		"raport_formatted_date": time.Now().Format("Monday, 2 January 2006"),
		"gender":                gender,
		"birth_date":            formattedTime,
		"branch_name":           res.BranchName,
		"selection":             "SKD-PTK",
		"formation":             formation,
		"phone":                 res.Phone,
	}

	filter := ""
	trialCount := 0

	switch strings.ToUpper(fil) {
	case "PRE_UKA":
		filter = "pre-uka"
	case "ALL_MODULE":
		filter = "all-module"
	case "UKA_STAGE":
		filter = "challenge-uka"
	}

	hisRes, er := GetHistoryPTKElasticFetchStudentReport(uint(smId), typStg, filter)
	if er != nil {
		return nil, er
	}

	sort.SliceStable(hisRes, func(i, j int) bool {
		return hisRes[i].Start.After(hisRes[j].End)
	})

	packageIdList := map[int]bool{}

	var scoreTotal float64
	var scoreTWK float64
	var scoreTIU float64
	var scoreTKP float64
	var avgScore float64
	totalPassed := int(0)
	totalFailed := int(0)
	twkPass := int(0)
	tiuPass := int(0)
	tkpPass := int(0)
	twkFailed := int(0)
	tiuFailed := int(0)
	tkpFailed := int(0)
	withCode := int(0)
	stgTotal := int(0)

	studentHistory := []map[string]any{}
	dn := 0
	for _, t := range hisRes {
		if t.ExamName == "" {
			continue
		}
		// if strings.Contains(t.ExamName, "Pre") {
		// 	continue
		// }
		// if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
		// 	continue
		// }

		if filter == "challenge-uka" && resBinsusReport != nil {
			isUKADataExist := false
			for _, k := range resBinsusReport.Record.Score.ScoresData {
				if t.TaskID == k.TaskID {
					isUKADataExist = true
				}
			}
			if !isUKADataExist {
				continue
			}
		}

		dn += 1
		isPass := false
		scoreTotal += t.Total
		scoreTWK += t.Twk
		scoreTIU += t.Tiu
		scoreTKP += t.Tkp
		if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
			totalPassed += 1
			isPass = true
		} else {
			totalFailed += 1
		}

		if t.Twk >= t.TwkPass {
			twkPass += 1
		}
		if t.Tiu >= t.TiuPass {
			tiuPass += 1
		}
		if t.Tkp >= t.TkpPass {
			tkpPass += 1
		}
		if t.Twk < t.TwkPass {
			twkFailed += 1
		}
		if t.Tiu < t.TiuPass {
			tiuFailed += 1
		}
		if t.Tkp < t.TkpPass {
			tkpFailed += 1
		}
		if t.PackageType == "WITH_CODE" {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				withCode += 1
			}
		} else {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				stgTotal += 1
			}
		}

		studentHistory = append(studentHistory, map[string]any{
			"exam_name":          t.ExamName,
			"start_date":         t.Start,
			"task_id":            t.TaskID,
			"duration":           t.End.Sub(t.Start).Milliseconds(),
			"duration_formatted": fmt.Sprintf("%02d:%02d:%02d", int(t.End.Sub(t.Start).Hours()), int(t.End.Sub(t.Start).Minutes())%60, int(t.End.Sub(t.Start).Seconds())%60),
			"date_formatted":     t.Start.Format("02/01/2006"),
			"time_formatted":     t.Start.Format("15:04:05"),
			"total":              t.Total,
			"status":             isPass,
			"category": map[string]any{
				"twk": map[string]any{
					"score":         t.Twk,
					"passing_grade": t.TwkPass,
					"is_pass":       t.Twk >= t.TwkPass,
					"duration":      nil,
				},
				"tiu": map[string]any{
					"score":         t.Tiu,
					"passing_grade": t.TiuPass,
					"is_pass":       t.Tiu >= t.TiuPass,
					"duration":      nil,
				},
				"tkp": map[string]any{
					"score":         t.Tkp,
					"passing_grade": t.TkpPass,
					"is_pass":       t.Tkp >= t.TkpPass,
					"duration":      nil,
				},
			},
		})

	}

	totalHis := dn - trialCount
	done := dn - trialCount
	avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
	avgScoreTWK := helpers.RoundFloat(scoreTWK/float64(totalHis), 1)
	avgScoreTIU := helpers.RoundFloat(scoreTIU/float64(totalHis), 1)
	avgScoreTKP := helpers.RoundFloat(scoreTKP/float64(totalHis), 1)

	passPercent := helpers.RoundFloat(float64(avgScore)/float64(studentTarget.TargetScore)*100, 1)
	passPercentTWK := helpers.RoundFloat(float64(twkPass)/float64(totalHis)*100, 1)
	passPercentTIU := helpers.RoundFloat(float64(tiuPass)/float64(totalHis)*100, 1)
	passPercentTKP := helpers.RoundFloat(float64(tkpPass)/float64(totalHis)*100, 1)

	if _isNaNorInf(avgScore) {
		avgScore = 0
	}
	if _isNaNorInf(avgScoreTWK) {
		avgScoreTWK = 0
	}
	if _isNaNorInf(avgScoreTIU) {
		avgScoreTIU = 0
	}
	if _isNaNorInf(avgScoreTKP) {
		avgScoreTKP = 0
	}

	if _isNaNorInf(passPercent) {
		passPercent = 0
	}
	if _isNaNorInf(passPercentTWK) {
		passPercentTWK = 0
	}
	if _isNaNorInf(passPercentTIU) {
		passPercentTIU = 0
	}
	if _isNaNorInf(passPercentTKP) {
		passPercentTKP = 0
	}
	skip := ""

	prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
	if err != nil {
		return nil, err
	}
	// resStages, err := GetAllStudentStageClass("PTK")
	// if err != nil {
	// 	return nil, err
	// }

	countStgLv := 0

	for _, obj := range prodRe.Data {
		isSkipped := false
		if strings.ToLower(obj.ProductProgram) != "skd" {
			continue
		}
		for _, tag := range obj.ProductTags {
			if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
				isSkipped = true
				continue
			}

			if strings.Contains(tag, "MATERIAL") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "CPNS") {
				isSkipped = true
				continue
			}
		}
		if isSkipped {
			continue
		}
		if fil != "" {
			if filter == "pre-uka" {
				for _, tag := range obj.ProductTags {
					if typStg == "UMUM" {
						if strings.Contains(tag, "STAGE_PRE_UKA") {
							countStgLv += 1
						}
					} else {
						re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
						matches := re.FindStringSubmatch(tag)

						// Check if there are captured substrings
						if len(matches) > 1 {

							// Iterate over captured substrings starting from index 1
							for _, v := range matches[1:] {

								// Convert both strings to lowercase for case-insensitive comparison
								length := (len(skip))
								skip = v
								if len(v) == length {
									continue // Skip processing if already processed
								} else {
									if v == "BINSUS" {
										countStgLv += 14
									} else if v == "REGULER" {
										countStgLv += 14

									} else {
										countStgLv += 14
									}

									// Mark the value as processed
									// resStages, err := GetAllStudentStageClass("PTK", v)
									// if err != nil {
									// 	return nil, err
									// }

									// for _, mdl := range resStages.Data {
									// 	if mdl.ModuleType == "PLATINUM" {
									// 		countStgLv += 1
									// 	}
									// }
								}
							}
						}

					}

				}
			} else if filter == "challenge-uka" {
				for _, tag := range obj.ProductTags {
					if typStg == "UMUM" {
						if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
							countStgLv += 1
						}
					} else {
						re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
						matches := re.FindStringSubmatch(tag)

						// Check if there are captured substrings
						if len(matches) > 1 {

							// Iterate over captured substrings starting from index 1
							for _, v := range matches[1:] {

								// Convert both strings to lowercase for case-insensitive comparison
								length := (len(skip))
								skip = v
								if len(v) == length {
									continue // Skip processing if already processed
								} else {
									if v == "BINSUS" {
										countStgLv += 14
									} else if v == "REGULER" {
										countStgLv += 14

									} else {
										countStgLv += 14
									}

									// Mark the value as processed
									// resStages, err := GetAllStudentStageClass("PTK", v)
									// if err != nil {
									// 	return nil, err
									// }

									// for _, mdl := range resStages.Data {
									// 	if mdl.ModuleType == "PREMIUM_TRYOUT" {
									// 		countStgLv += 1
									// 	}
									// }
								}
							}
						}

					}

				}
			} else if filter == "all-module" {
				if typStg == "UMUM" {
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
							countStgLv += 1
						}
					}
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_PRE_UKA") {
							countStgLv += 1
						}
					}
				} else {
					for _, tag := range obj.ProductTags {
						re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
						matches := re.FindStringSubmatch(tag)

						// Check if there are captured substrings
						if len(matches) > 1 {

							// Iterate over captured substrings starting from index 1
							for _, v := range matches[1:] {

								// Convert both strings to lowercase for case-insensitive comparison
								length := (len(skip))
								skip = v
								if len(v) == length {
									continue // Skip processing if already processed
								} else {
									if v == "BINSUS" {
										countStgLv += 28
									} else if v == "REGULER" {
										countStgLv += 28

									} else {
										countStgLv += 28
									}

									// Mark the value as processed
									// resStages, err := GetAllStudentStageClass("PTK", v)
									// if err != nil {
									// 	return nil, err
									// }

									// for _, mdl := range resStages.Data {
									// 	if mdl.ModuleType == "PREMIUM_TRYOUT" || mdl.ModuleType == "PLATINUM" {
									// 		countStgLv += 1
									// 	}
									// }
								}
							}
						}
					}

				}

			}
		} else {
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "STAGE_LEVEL_") {
					countStgLv += 1
				}
			}
		}
	}
	// } else {
	// 	if filter == "pre-uka" {
	// 		for _, mdl := range resStages.Data {
	// 			if mdl.ModuleType == "PLATINUM" {
	// 				countStgLv += 1
	// 			}
	// 		}
	// 	} else if filter == "challenge-uka" {
	// 		for _, mdl := range resStages.Data {
	// 			if mdl.ModuleType == "PREMIUM_TRYOUT" {
	// 				countStgLv += 1
	// 			}
	// 		}

	// 	} else {
	// 		for _, mdl := range resStages.Data {
	// 			if mdl.ModuleType == "PREMIUM_TRYOUT" || mdl.ModuleType == "PLATINUM" {
	// 				countStgLv += 1
	// 			}
	// 		}
	// 	}
	// }

	ownMod := withCode

	if (countStgLv - stgTotal) < 0 {
		ownMod += stgTotal
	} else {
		ownMod += countStgLv
	}

	donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)

	avgDoneScore := helpers.RoundFloat(scoreTotal/float64(dn), 1)
	if _isNaNorInf(donePerct) {
		donePerct = 0
	}
	if _isNaNorInf(avgDoneScore) {
		avgDoneScore = 0
	}

	// bknScore := map[string]any{
	// 	"twk":           float64(0),
	// 	"twk_pass":      false,
	// 	"tiu":           float64(0),
	// 	"tiu_pass":      false,
	// 	"tkp":           float64(0),
	// 	"tkp_pass":      false,
	// 	"total":         float64(0),
	// 	"year":          2023,
	// 	"bkn_attempted": false,
	// 	"is_continue":   false,
	// 	"is_pass":       false,
	// }

	if donePerct > 100 {
		donePerct = 100
	}

	if passPercent > 100 {
		passPercent = 100
	}

	setMembers, errs := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptk:user_%d:pre-test", smId)).Result()
	if errs != nil {
		return nil, err
	}

	pre, err := helpers.FormatingCategory(setMembers)
	if err != nil {
		return nil, err
	}

	newPre := map[string]helpers.DataStruct{}
	for pe := range pre {
		peN := helpers.ToLowerAndUnderscore(pe)
		newPre[peN] = pre[pe]
	}

	setMembersPost, errs := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptk:user_%d:post-test", smId)).Result()
	if errs != nil {
		return nil, errs
	}
	post, err := helpers.FormatingCategory(setMembersPost)
	if err != nil {
		return nil, err
	}

	newPost := map[string]helpers.DataStruct{}
	for po := range post {
		poN := helpers.ToLowerAndUnderscore(po)
		newPost[poN] = post[po]
	}

	reportData = map[string]any{
		"student": studentData,
		"summary": map[string]any{
			"category": map[string]any{
				"twk": map[string]any{
					"average":            avgScoreTWK,
					"is_pass":            avgScoreTWK >= 65,
					"passing_percentage": passPercentTWK,
				},
				"tiu": map[string]any{
					"average":            avgScoreTIU,
					"is_pass":            avgScoreTIU >= 80,
					"passing_percentage": passPercentTIU,
				},
				"tkp": map[string]any{
					"average":            avgScoreTKP,
					"is_pass":            avgScoreTKP >= 156,
					"passing_percentage": passPercentTKP,
				},
			},
			"received":             ownMod,
			"completed":            done,
			"passed":               totalPassed,
			"completed_percentage": donePerct,
			"passing_percentage":   passPercent,
			"average_score":        avgScore,
			"average_total":        avgDoneScore,
			"average_is_pass":      avgScoreTKP >= 156 && avgScoreTIU >= 80 && avgScoreTWK >= 65,
		},
		"histories": studentHistory,
		// "assessment": overallResult,
		"pre_test":  newPre,
		"post_test": newPost,
	}

	return reportData, nil
}

func FetchPTKRaport(smId uint, fil string, raportType string) (map[string]any, error) {

	res, err := GetStudentProfileElastic(int(smId))
	if err != nil {
		return nil, err
	}
	var studentData = map[string]any{}
	var reportData = map[string]any{}
	var studentTarget requests.StudentProfilePtkElastic
	var resBinsusReport *mockstruct.BinsusScreeningSummary
	var binsusResult map[string]any

	yrNow := time.Now().Year()

	if raportType != "raport" {

		binsusResult = map[string]any{
			"is_available":          false,
			"recommended_with_note": false,
			"is_recommended":        false,
			"year":                  yrNow,
		}

		var errBinsusReport error
		resBinsusReport, errBinsusReport = FetchStudentBinsusSummary(res.SmartbtwID, yrNow)
		if errBinsusReport == nil {
			binsusResult = map[string]any{
				"is_available":          true,
				"recommended_with_note": resBinsusReport.ChoosenSummary.FinalState == "recommended_with_note",
				"is_recommended":        resBinsusReport.ChoosenSummary.FinalPass,
				"year":                  yrNow,
			}
		}
		var errStudentTarget error
		studentTarget, errStudentTarget = GetStudentProfilePTKElastic(res.SmartbtwID)
		if errStudentTarget != nil {
			return nil, errStudentTarget
		}
		joinedClassList := ""
		joinedClassName := ""
		isClassJoined := false
		stdClsList, errCls := GetStudentJoinedClassList(int32(res.SmartbtwID), &yrNow, false)
		if errCls == nil {
			for _, k := range stdClsList {
				for _, tagsClass := range k.Tags {
					if tagsClass == "BINSUS" || tagsClass == "REGULER" || tagsClass == "INTENSIF" {
						joinedClassList = tagsClass
						break
					}
				}
				if joinedClassList != "" {
					joinedClassName = k.Title
					isClassJoined = true
				}
				if isClassJoined {
					break
				}
			}
		}

		polbitType := "PUSAT"
		polbitCompetitionType := ""
		formationQuota := -1

		if studentTarget.PolbitCompetitionID != nil {
			locId := 0
			if studentTarget.PolbitLocationID != nil {
				locId = *studentTarget.PolbitLocationID
			}
			chRes, err := GetCompetitionDataPTK(uint(studentTarget.MajorID), uint(locId), res.Gender, studentTarget.PolbitType)
			if err == nil {

				if studentTarget.PolbitLocationID != nil {
					if studentTarget.PolbitType == "DAERAH_REGION" {
						polbitType = res.DomicileRegion
					} else if studentTarget.PolbitType == "DAERAH_PROVINCE" {
						polbitType = fmt.Sprintf("Provinsi %s", res.DomicileProvince)
					} else if strings.Contains(studentTarget.PolbitType, "AFIRMASI") {
						if studentTarget.PolbitType == "DAERAH_AFIRMASI_PROVINCE" {
							polbitType = fmt.Sprintf("%s (Afirmasi)", res.DomicileProvince)
						} else if strings.Contains(studentTarget.PolbitType, "PUSAT_AFIRMASI_PROVINCE") {
							polbitType = fmt.Sprintf("%s (Afirmasi)", res.DomicileProvince)
						} else {
							polbitType = fmt.Sprintf("%s (Afirmasi)", res.DomicileProvince)
						}
					}
				}
				formationQuota = int(chRes.MajorQuota)
				if chRes.CompetitionType != nil {
					polbitCompetitionType = *chRes.CompetitionType
				}
			}
		}

		branchName := "Pusat"
		if res.BranchName != nil && res.BranchCode != nil {
			if (*res.BranchCode) != "PT0000" {
				branchName = strings.Replace(*res.BranchName, "Bimbel BTW ", "", 1)
			}
		}

		studentData = map[string]any{
			"name":                  res.Name,
			"email":                 res.Email,
			"program":               "PTK",
			"raport_date":           time.Now().Format("2006-01-02"),
			"raport_formatted_date": time.Now().Format("Monday, 2 January 2006"),
			"school_name":           studentTarget.SchoolName,
			"major_name":            studentTarget.MajorName,
			"target_score":          studentTarget.TargetScore,
			"formation":             polbitType,
			"competition_type":      polbitCompetitionType,
			"quota":                 formationQuota,
			"quota_available":       formationQuota != -1,
			"class_joined":          isClassJoined,
			"class_program_type":    joinedClassList,
			"class_name":            joinedClassName,
			"branch_name":           branchName,
		}
	} else {

		studentData = map[string]any{
			"name":                  res.Name,
			"email":                 res.Email,
			"program":               "PTK",
			"raport_date":           time.Now().Format("2006-01-02"),
			"raport_formatted_date": time.Now().Format("Monday, 2 January 2006"),
		}
	}
	filter := ""
	isNoTrial := false
	trialCount := 0

	switch strings.ToUpper(fil) {
	case "PRE_UKA":
		filter = "pre-uka"
	case "UKA_CODE":
		filter = "with_code"
	case "CHALLENGE_UKA":
		filter = "challenge-uka"
	case "NO_TRIAL":
		isNoTrial = true

	}

	hisRes, er := GetStudentHistoryPTKElasticFilterOld(int(smId), filter)
	if er != nil {
		return nil, er
	}

	sort.SliceStable(hisRes, func(i, j int) bool {
		return hisRes[i].Start.After(hisRes[j].End)
	})

	packageIdList := map[int]bool{}

	var scoreTotal float64
	var scoreTWK float64
	var scoreTIU float64
	var scoreTKP float64
	var avgScore float64
	totalPassed := int(0)
	totalFailed := int(0)
	twkPass := int(0)
	tiuPass := int(0)
	tkpPass := int(0)
	twkFailed := int(0)
	tiuFailed := int(0)
	tkpFailed := int(0)
	withCode := int(0)
	stgTotal := int(0)

	studentHistory := []map[string]any{}
	dn := 0
	for _, t := range hisRes {
		if t.ExamName == "" {
			continue
		}

		if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
			continue
		}

		if strings.Contains(t.ExamName, "Pre") || strings.Contains(t.ExamName, "Post") {
			continue
		}

		if isNoTrial {
			if t.ModuleType == "TESTING" || t.ModuleType == "TRIAL" {
				trialCount += 1
				continue
			}
		}

		if filter == "challenge-uka" && resBinsusReport != nil {
			isUKADataExist := false
			for _, k := range resBinsusReport.Record.Score.ScoresData {
				if t.TaskID == k.TaskID {
					isUKADataExist = true
				}
			}
			if !isUKADataExist {
				continue
			}
		}

		dn += 1
		isPass := false
		scoreTotal += t.Total
		scoreTWK += t.Twk
		scoreTIU += t.Tiu
		scoreTKP += t.Tkp
		if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
			totalPassed += 1
			isPass = true
		} else {
			totalFailed += 1
		}

		if t.Twk >= t.TwkPass {
			twkPass += 1
		}
		if t.Tiu >= t.TiuPass {
			tiuPass += 1
		}
		if t.Tkp >= t.TkpPass {
			tkpPass += 1
		}
		if t.Twk < t.TwkPass {
			twkFailed += 1
		}
		if t.Tiu < t.TiuPass {
			tiuFailed += 1
		}
		if t.Tkp < t.TkpPass {
			tkpFailed += 1
		}
		if t.PackageType == "WITH_CODE" {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				withCode += 1
			}
		} else {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				stgTotal += 1
			}
		}

		studentHistory = append(studentHistory, map[string]any{
			"exam_name":          t.ExamName,
			"start_date":         t.Start,
			"duration":           t.End.Sub(t.Start).Milliseconds(),
			"duration_formatted": fmt.Sprintf("%02d:%02d:%02d", int(t.End.Sub(t.Start).Hours()), int(t.End.Sub(t.Start).Minutes())%60, int(t.End.Sub(t.Start).Seconds())%60),
			"date_formatted":     t.Start.Format("02/01/2006"),
			"time_formatted":     t.Start.Format("15:04:05"),
			"total":              t.Total,
			"status":             isPass,
			"category": map[string]any{
				"twk": map[string]any{
					"score":         t.Twk,
					"passing_grade": t.TwkPass,
					"is_pass":       t.Twk >= t.TwkPass,
					"duration":      nil,
				},
				"tiu": map[string]any{
					"score":         t.Tiu,
					"passing_grade": t.TiuPass,
					"is_pass":       t.Tiu >= t.TiuPass,
					"duration":      nil,
				},
				"tkp": map[string]any{
					"score":         t.Tkp,
					"passing_grade": t.TkpPass,
					"is_pass":       t.Tkp >= t.TkpPass,
					"duration":      nil,
				},
			},
		})

	}

	totalHis := dn - trialCount
	done := dn - trialCount
	avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
	avgScoreTWK := helpers.RoundFloat(scoreTWK/float64(totalHis), 1)
	avgScoreTIU := helpers.RoundFloat(scoreTIU/float64(totalHis), 1)
	avgScoreTKP := helpers.RoundFloat(scoreTKP/float64(totalHis), 1)

	passPercent := helpers.RoundFloat(float64(totalPassed)/float64(totalHis)*100, 1)
	passPercentTWK := helpers.RoundFloat(float64(twkPass)/float64(totalHis)*100, 1)
	passPercentTIU := helpers.RoundFloat(float64(tiuPass)/float64(totalHis)*100, 1)
	passPercentTKP := helpers.RoundFloat(float64(tkpPass)/float64(totalHis)*100, 1)

	if _isNaNorInf(avgScore) {
		avgScore = 0
	}
	if _isNaNorInf(avgScoreTWK) {
		avgScoreTWK = 0
	}
	if _isNaNorInf(avgScoreTIU) {
		avgScoreTIU = 0
	}
	if _isNaNorInf(avgScoreTKP) {
		avgScoreTKP = 0
	}

	if _isNaNorInf(passPercent) {
		passPercent = 0
	}
	if _isNaNorInf(passPercentTWK) {
		passPercentTWK = 0
	}
	if _isNaNorInf(passPercentTIU) {
		passPercentTIU = 0
	}
	if _isNaNorInf(passPercentTKP) {
		passPercentTKP = 0
	}

	prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
	if err != nil {
		return nil, err
	}

	countStgLv := 0
	for _, obj := range prodRe.Data {
		isSkipped := false
		if strings.ToLower(obj.ProductProgram) != "skd" {
			continue
		}
		for _, tag := range obj.ProductTags {
			if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "MATERIAL") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "CPNS") {
				isSkipped = true
				continue
			}
		}
		if isSkipped {
			continue
		}
		if fil != "" {
			if filter == "pre-uka" {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_PRE_UKA") {
						countStgLv += 1
					}
				}
			} else if filter == "challenge-uka" {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
						countStgLv += 1
					}
				}
			}
		} else {
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "STAGE_LEVEL_") {
					countStgLv += 1
				}
			}
		}
	}

	ownMod := withCode

	if (countStgLv - stgTotal) < 0 {
		ownMod += stgTotal
	} else {
		ownMod += countStgLv
	}

	donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)

	avgDoneScore := helpers.RoundFloat(scoreTotal/float64(dn), 1)
	if _isNaNorInf(donePerct) {
		donePerct = 0
	}
	if _isNaNorInf(avgDoneScore) {
		avgDoneScore = 0
	}

	bknScore := map[string]any{
		"twk":           float64(0),
		"twk_pass":      false,
		"tiu":           float64(0),
		"tiu_pass":      false,
		"tkp":           float64(0),
		"tkp_pass":      false,
		"total":         float64(0),
		"year":          2023,
		"bkn_attempted": false,
		"is_continue":   false,
		"is_pass":       false,
	}

	finalPassingPercentage := float64(0)

	deviationExist := false
	twkDeviation := float64(0)
	twkPercentage := float64(0)
	tiuDeviation := float64(0)
	tiuPercentage := float64(0)
	tkpDeviation := float64(0)
	tkpPercentage := float64(0)
	totalDeviation := float64(0)
	totalPercentage := float64(0)
	bknTargetDeviation := float64(0)
	bknTargetPercentage := float64(0)
	ukaTargetDeviation := float64(0)
	ukaTargetPercentage := float64(0)

	if raportType != "raport" {
		ukaTargetDeviation = math.Round((avgScore-studentTarget.TargetScore)*10) / 10
		ukaTargetPercentage = math.Round((float64(ukaTargetDeviation) / float64(studentTarget.TargetScore)) * 100)

		if _isNaNorInf(ukaTargetDeviation) {
			ukaTargetDeviation = 0
		}
		if _isNaNorInf(ukaTargetPercentage) {
			ukaTargetPercentage = 0
		}

		bknSc, err := GetSingleBKNScoreByYearAndStudent(res.SmartbtwID, 2023)
		if err == nil {
			isBknPass := false
			if bknSc.Twk >= 65 && bknSc.Tiu >= 80 && bknSc.Tkp >= 156 {
				isBknPass = true
			}

			bknScore = map[string]any{
				"twk":           bknSc.Twk,
				"twk_pass":      bknSc.Twk >= 65,
				"tiu":           bknSc.Tiu,
				"tiu_pass":      bknSc.Tiu >= 80,
				"tkp":           bknSc.Tkp,
				"tkp_pass":      bknSc.Tkp >= 156,
				"total":         bknSc.Total,
				"year":          bknSc.Year,
				"bkn_attempted": true,
				"is_continue":   bknSc.IsContinue,
				"is_pass":       isBknPass,
			}

			finalPassingPercentage = helpers.RoundFloat(float64(bknSc.Total)/float64(studentTarget.TargetScore)*100, 1)
			if _isNaNorInf(finalPassingPercentage) {
				finalPassingPercentage = 0
			}

			twkDeviation = math.Round((bknSc.Twk-avgScoreTWK)*10) / 10
			twkPercentage = math.Round((float64(twkDeviation) / float64(avgScoreTWK)) * 100)
			tiuDeviation = math.Round((bknSc.Tiu-avgScoreTIU)*10) / 10
			tiuPercentage = math.Round((float64(tiuDeviation) / float64(avgScoreTIU)) * 100)
			tkpDeviation = math.Round((bknSc.Tkp-avgScoreTKP)*10) / 10
			tkpPercentage = math.Round((float64(tkpDeviation) / float64(avgScoreTKP)) * 100)
			totalDeviation = math.Round((bknSc.Total-avgScore)*10) / 10
			totalPercentage = math.Round((float64(totalDeviation) / float64(avgScore)) * 100)
			bknTargetDeviation = math.Round((bknSc.Total-studentTarget.TargetScore)*10) / 10
			bknTargetPercentage = math.Round((float64(bknTargetDeviation) / float64(studentTarget.TargetScore)) * 100)
			if _isNaNorInf(twkDeviation) {
				twkDeviation = 0
			}
			if _isNaNorInf(twkPercentage) {
				twkPercentage = 0
			}
			if _isNaNorInf(tiuDeviation) {
				tiuDeviation = 0
			}
			if _isNaNorInf(tiuPercentage) {
				tiuPercentage = 0
			}
			if _isNaNorInf(tkpDeviation) {
				tkpDeviation = 0
			}
			if _isNaNorInf(tkpPercentage) {
				tkpPercentage = 0
			}
			if _isNaNorInf(totalDeviation) {
				totalDeviation = 0
			}
			if _isNaNorInf(totalPercentage) {
				totalPercentage = 0
			}
			if _isNaNorInf(bknTargetDeviation) {
				bknTargetDeviation = 0
			}
			if _isNaNorInf(bknTargetPercentage) {
				bknTargetPercentage = 0
			}
			deviationExist = true
		}

		targetPassingPercentage := helpers.RoundFloat(float64(avgScore)/float64(studentTarget.TargetScore)*100, 1)
		if _isNaNorInf(targetPassingPercentage) {
			targetPassingPercentage = 0
		}

		reportData = map[string]any{
			"student":          studentData,
			"bkn_score":        bknScore,
			"binsus_screening": binsusResult,
			"summary": map[string]any{
				"deviation": map[string]any{
					"twk": map[string]any{
						"percentage":  twkPercentage,
						"differences": twkDeviation,
					},
					"tiu": map[string]any{
						"percentage":  tiuPercentage,
						"differences": tiuDeviation,
					},
					"tkp": map[string]any{
						"percentage":  tkpPercentage,
						"differences": tkpDeviation,
					},
					"total": map[string]any{
						"percentage":  totalPercentage,
						"differences": totalDeviation,
					},
					"bkn": map[string]any{
						"percentage":  bknTargetPercentage,
						"differences": bknTargetDeviation,
					},
					"uka": map[string]any{
						"percentage":  ukaTargetPercentage,
						"differences": ukaTargetDeviation,
					},
					"available": deviationExist,
				},
				"category": map[string]any{
					"twk": map[string]any{
						"average":                avgScoreTWK,
						"is_pass":                avgScoreTWK >= 65,
						"passing_percentage":     passPercentTWK,
						"approached_explanation": nil,
						"total_explanation":      nil,
						"percentage_explanation": nil,
						"calculated_answer_time": nil,
					},
					"tiu": map[string]any{
						"average":                avgScoreTIU,
						"is_pass":                avgScoreTIU >= 80,
						"passing_percentage":     passPercentTIU,
						"approached_explanation": nil,
						"total_explanation":      nil,
						"percentage_explanation": nil,
						"calculated_answer_time": nil,
					},
					"tkp": map[string]any{
						"average":                avgScoreTKP,
						"is_pass":                avgScoreTKP >= 156,
						"passing_percentage":     passPercentTKP,
						"approached_explanation": nil,
						"total_explanation":      nil,
						"percentage_explanation": nil,
						"calculated_answer_time": nil,
					},
				},
				"received":                  ownMod,
				"completed":                 done,
				"passed":                    totalPassed,
				"completed_percentage":      donePerct,
				"passing_percentage":        passPercent,
				"average_score":             avgScore,
				"average_total":             avgDoneScore,
				"average_is_pass":           avgScoreTKP >= 156 && avgScoreTIU >= 80 && avgScoreTWK >= 65,
				"final_passing_percentage":  finalPassingPercentage,
				"target_passing_percentage": targetPassingPercentage,
			},
			"histories": studentHistory,
		}
	} else {

		reportData = map[string]any{
			"student": studentData,
			"summary": map[string]any{
				"category": map[string]any{
					"twk": map[string]any{
						"average":            avgScoreTWK,
						"is_pass":            avgScoreTWK >= 65,
						"passing_percentage": passPercentTWK,
					},
					"tiu": map[string]any{
						"average":            avgScoreTIU,
						"is_pass":            avgScoreTIU >= 80,
						"passing_percentage": passPercentTIU,
					},
					"tkp": map[string]any{
						"average":            avgScoreTKP,
						"is_pass":            avgScoreTKP >= 156,
						"passing_percentage": passPercentTKP,
					},
				},
				"received":             ownMod,
				"completed":            done,
				"passed":               totalPassed,
				"completed_percentage": donePerct,
				"passing_percentage":   passPercent,
				"average_score":        avgScore,
				"average_total":        avgDoneScore,
				"average_is_pass":      avgScoreTKP >= 156 && avgScoreTIU >= 80 && avgScoreTWK >= 65,
			},
			"histories": studentHistory,
		}
	}
	return reportData, nil
}

func FetchPTNUKARaport(smId uint, typStg, fil string) (map[string]any, error) {

	res, err := GetStudentProfileElastic(int(smId))
	if err != nil {
		return nil, err
	}
	filter := ""

	switch strings.ToUpper(fil) {
	case "PRE_UKA":
		filter = "pre-uka"
	case "ALL_MODULE":
		filter = "all-module"
	case "UKA_STAGE":
		filter = "challenge-uka"
	}

	studentTarget, err := GetStudentProfilePTNElastic(res.SmartbtwID)
	if err != nil {
		return nil, err
	}

	comp, err := GetCompetitonPTN(uint(studentTarget.MajorID))
	if err != nil {
		return nil, err
	}
	gender := "Perempuan"
	if res.Gender == "L" {
		gender = "Laki-laki"
	}
	targetFormat := "02/01/2006"
	formattedTime := res.BirthDate.Format(targetFormat)

	studentData := map[string]any{
		"name":                  res.Name,
		"email":                 res.Email,
		"program":               "PTN",
		"school_name":           studentTarget.SchoolName,
		"major_name":            studentTarget.MajorName,
		"target_score":          studentTarget.TargetScore,
		"quota":                 comp.SbmptnCapacity,
		"raport_date":           time.Now().Format("2006-01-02"),
		"raport_formatted_date": time.Now().Format("Monday, 2 January 2006"),
		"gender":                gender,
		"birth_date":            formattedTime,
		"branch_name":           res.BranchName,
		"selection":             "SKD-PTN",
		"phone":                 res.Phone,
	}

	hisRes, er := GetHistoryPTNElasticPeforma(uint(res.SmartbtwID), typStg, filter)
	if er != nil {
		return nil, er
	}

	sort.SliceStable(hisRes, func(i, j int) bool {
		return hisRes[i].Start.After(hisRes[j].End)
	})

	packageIdList := map[int]bool{}

	var scoreTotal float64
	var scorePU float64
	var scorePPU float64
	var scorePBM float64
	var scorePM float64
	var scorePK float64
	var scoreLBIND float64
	var scoreLBING float64
	var avgScore float64
	withCode := int(0)
	stgTotal := int(0)

	studentHistory := []map[string]any{}
	dn := 0
	for _, t := range hisRes {
		if t.ExamName == "" {
			continue
		}
		// if strings.Contains(t.ExamName, "Pre") {
		// 	continue
		// }
		// if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
		// 	continue
		// }
		dn += 1
		isPass := false
		scoreTotal += t.Total
		scorePU += t.PenalaranUmum
		scorePBM += t.PemahamanBacaan
		scorePPU += t.PengetahuanUmum
		scorePM += t.PenalaranMatematika
		scorePK += t.PengetahuanKuantitatif
		scoreLBIND += t.LiterasiBahasaIndonesia
		scoreLBING += t.LiterasiBahasaInggris

		if t.PackageType == "WITH_CODE" {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				withCode += 1
			}
		} else {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				stgTotal += 1
			}
		}

		studentHistory = append(studentHistory, map[string]any{
			"exam_name":          t.ExamName,
			"start_date":         t.Start,
			"task_id":            t.TaskID,
			"duration":           t.End.Sub(t.Start).Milliseconds(),
			"duration_formatted": fmt.Sprintf("%02d:%02d:%02d", int(t.End.Sub(t.Start).Hours()), int(t.End.Sub(t.Start).Minutes())%60, int(t.End.Sub(t.Start).Seconds())%60),
			"date_formatted":     t.Start.Format("02/01/2006"),
			"time_formatted":     t.Start.Format("15:04:05"),
			"total":              t.Total,
			"status":             isPass,
			"category": map[string]any{
				"pu": map[string]any{
					"score": t.PenalaranUmum,
				},
				"ppu": map[string]any{
					"score": t.PengetahuanUmum,
				},
				"pbm": map[string]any{
					"score": t.PemahamanBacaan,
				},
				"pk": map[string]any{
					"score": t.PengetahuanKuantitatif,
				},
				"pm": map[string]any{
					"score": t.PenalaranMatematika,
				},
				"lbind": map[string]any{
					"score": t.LiterasiBahasaIndonesia,
				},
				"lbing": map[string]any{
					"score": t.LiterasiBahasaInggris,
				},
			},
		})

	}

	totalHis := dn
	avgScore = helpers.RoundFloats(scoreTotal/float64(totalHis), 1)
	avgScorePU := helpers.RoundFloats(scorePU/float64(totalHis), 1)
	avgScorePPU := helpers.RoundFloats(scorePPU/float64(totalHis), 1)
	avgScorePBM := helpers.RoundFloats(scorePBM/float64(totalHis), 1)
	avgScorePK := helpers.RoundFloats(scorePK/float64(totalHis), 1)
	avgScorePM := helpers.RoundFloats(scorePM/float64(totalHis), 1)
	avgScoreLBIND := helpers.RoundFloats(scoreLBIND/float64(totalHis), 1)
	avgScoreLBING := helpers.RoundFloats(scoreLBING/float64(totalHis), 1)

	prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
	if err != nil {
		return nil, err
	}
	if _isNaNorInf(avgScore) {
		avgScore = 0
	}
	if _isNaNorInf(avgScorePU) {
		avgScorePU = 0
	}
	if _isNaNorInf(avgScorePPU) {
		avgScorePPU = 0
	}
	if _isNaNorInf(avgScorePK) {
		avgScorePK = 0
	}
	if _isNaNorInf(avgScorePBM) {
		avgScorePBM = 0
	}
	if _isNaNorInf(avgScorePM) {
		avgScorePM = 0
	}
	if _isNaNorInf(avgScoreLBIND) {
		avgScoreLBIND = 0
	}
	if _isNaNorInf(avgScoreLBING) {
		avgScoreLBING = 0
	}

	skip := ""
	// resStages, err := GetAllStudentStageClass("PTN")
	// if err != nil {
	// 	return nil, err
	// }

	countStgLv := 0

	for _, obj := range prodRe.Data {
		isSkipped := false
		if strings.ToLower(obj.ProductProgram) != "tps" {
			continue
		}
		for _, tag := range obj.ProductTags {
			if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "MATERIAL") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "CPNS") {
				isSkipped = true
				continue
			}
		}
		if isSkipped {
			continue
		}
		if fil != "" {
			if filter == "pre-uka" {
				for _, tag := range obj.ProductTags {
					if typStg == "UMUM" {
						if strings.Contains(tag, "STAGE_PRE_UKA") {
							countStgLv += 1
						}
					} else {
						re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
						matches := re.FindStringSubmatch(tag)

						// Check if there are captured substrings
						if len(matches) > 1 {

							// Iterate over captured substrings starting from index 1
							for _, v := range matches[1:] {

								// Convert both strings to lowercase for case-insensitive comparison
								length := (len(skip))
								skip = v
								if len(v) == length {
									continue // Skip processing if already processed
								} else {
									if v == "BINSUS" {
										countStgLv += 14
									} else if v == "REGULER" {
										countStgLv += 14

									} else {
										countStgLv += 14
									}

									// Mark the value as processed
									// resStages, err := GetAllStudentStageClass("PTN", v)
									// if err != nil {
									// 	return nil, err
									// }

									// for _, mdl := range resStages.Data {
									// 	if mdl.ModuleType == "PLATINUM" {
									// 		countStgLv += 1
									// 	}
									// }
								}
							}
						}

					}

				}
			} else if filter == "challenge-uka" {
				for _, tag := range obj.ProductTags {
					if typStg == "UMUM" {

						if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
							countStgLv += 1
						}
					} else {
						re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
						matches := re.FindStringSubmatch(tag)

						// Check if there are captured substrings
						if len(matches) > 1 {

							// Iterate over captured substrings starting from index 1
							for _, v := range matches[1:] {

								// Convert both strings to lowercase for case-insensitive comparison
								length := (len(skip))
								skip = v
								if len(v) == length {
									continue // Skip processing if already processed
								} else {
									if v == "BINSUS" {
										countStgLv += 14
									} else if v == "REGULER" {
										countStgLv += 14

									} else {
										countStgLv += 14
									}

									// Mark the value as processed
									// resStages, err := GetAllStudentStageClass("PTN", v)
									// if err != nil {
									// 	return nil, err
									// }

									// for _, mdl := range resStages.Data {
									// 	if mdl.ModuleType == "PREMIUM_TRYOUT" {
									// 		countStgLv += 1
									// 	}
									// }
								}
							}
						}

					}

				}
			} else if filter == "all-module" {
				if typStg == "UMUM" {
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
							countStgLv += 1
						}
					}
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_PRE_UKA") {
							countStgLv += 1
						}
					}
				} else {
					for _, tag := range obj.ProductTags {
						re := regexp.MustCompile(`MULTI_STAGES_(\w+)`)
						matches := re.FindStringSubmatch(tag)

						// Check if there are captured substrings
						if len(matches) > 1 {

							// Iterate over captured substrings starting from index 1
							for _, v := range matches[1:] {

								// Convert both strings to lowercase for case-insensitive comparison
								length := (len(skip))
								skip = v
								if len(v) == length {
									continue // Skip processing if already processed
								} else {
									if v == "BINSUS" {
										countStgLv += 28
									} else if v == "REGULER" {
										countStgLv += 28

									} else {
										countStgLv += 28
									}

									// Mark the value as processed
									// resStages, err := GetAllStudentStageClass("PTN", v)
									// if err != nil {
									// 	return nil, err
									// }

									// for _, mdl := range resStages.Data {
									// 	if mdl.ModuleType == "PREMIUM_TRYOUT" || mdl.ModuleType == "PLATINUM" {
									// 		countStgLv += 1
									// 	}
									// }
								}
							}
						}
					}

				}

			}
		} else {
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "STAGE_LEVEL_") {
					countStgLv += 1
				}
			}
		}
	}
	// } else {
	// 	if filter == "pre-uka" {
	// 		for _, mdl := range resStages.Data {
	// 			if mdl.ModuleType == "PLATINUM" {
	// 				countStgLv += 1
	// 			}
	// 		}
	// 	} else if filter == "challenge-uka" {
	// 		for _, mdl := range resStages.Data {
	// 			if mdl.ModuleType == "PREMIUM_TRYOUT" {
	// 				countStgLv += 1
	// 			}
	// 		}

	// 	} else {
	// 		for _, mdl := range resStages.Data {
	// 			if mdl.ModuleType == "PREMIUM_TRYOUT" || mdl.ModuleType == "PLATINUM" {
	// 				countStgLv += 1
	// 			}
	// 		}
	// 	}
	// }
	ownMod := withCode

	if (countStgLv - stgTotal) < 0 {
		ownMod += stgTotal
	} else {
		ownMod += countStgLv
	}

	donePerct := helpers.RoundFloats(float64(totalHis)/float64(ownMod)*100, 1)
	passPercent := helpers.RoundFloats(float64(avgScore)/float64(studentTarget.TargetScore)*100, 1)

	if donePerct > 100 {
		donePerct = 100
	}

	if passPercent > 100 {
		passPercent = 100
	}

	avgDoneScore := helpers.RoundFloats(scoreTotal/float64(dn), 1)
	if _isNaNorInf(donePerct) {
		donePerct = 0
	}
	if _isNaNorInf(passPercent) {
		passPercent = 0
	}
	if _isNaNorInf(avgDoneScore) {
		avgDoneScore = 0
	}

	setMembers, errs := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptn:user_%d:pre-test", smId)).Result()
	if errs != nil {
		return nil, err
	}

	pre, err := helpers.FormatingCategoryPTN(setMembers)
	if err != nil {
		return nil, err
	}

	setMembersPost, errs := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptn:user_%d:post-test", smId)).Result()
	if errs != nil {
		return nil, errs
	}
	post, err := helpers.FormatingCategoryPTN(setMembersPost)
	if err != nil {
		return nil, err
	}

	reportData := map[string]any{
		"student": studentData,
		"summary": map[string]any{
			"category": map[string]any{
				"pu": map[string]any{
					"average": avgScorePU,
				},
				"ppu": map[string]any{
					"average": avgScorePPU,
				},
				"pbm": map[string]any{
					"average": avgScorePBM,
				},
				"pk": map[string]any{
					"average": avgScorePK,
				},
				"pm": map[string]any{
					"average": avgScorePM,
				},
				"lbind": map[string]any{
					"average": avgScoreLBIND,
				},
				"lbing": map[string]any{
					"average": avgScoreLBING,
				},
			},
			"received":             ownMod,
			"completed":            totalHis,
			"passed":               totalHis,
			"completed_percentage": donePerct,
			"passing_percentage":   passPercent,
			"average_score":        avgScore,
			"average_total":        avgDoneScore,
		},
		"histories": studentHistory,
		// "assessment": resultMap,
		"pre_test":  pre,
		"post_test": post,
	}
	return reportData, nil
}

func FetchPTNPrePostRaport(smId uint) (map[string]any, error) {
	setMembers, err := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptn:user_%d:pre-test", smId)).Result()
	if err != nil {
		return nil, err
	}
	setMembersPost, err := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptn:user_%d:post-test", smId)).Result()
	if err != nil {
		return nil, err
	}

	resPre, err := helpers.FormatingCategoryPTN(setMembers)
	if err != nil {
		return nil, err
	}

	resPost, err := helpers.FormatingCategoryPTN(setMembersPost)
	if err != nil {
		return nil, err
	}

	type ResultDetails struct {
		MaxScore float64 `json:"max_score"`
		Score    float64 `json:"score"`
		Date     string  `json:"date"`
	}

	finalResult := map[string]any{}

	for materi, pr := range resPre {
		respRe := ResultDetails{
			Score:    float64(pr.Score),
			MaxScore: float64(pr.ScoreMax),
			Date:     pr.DateStart,
		}
		res := map[string]any{
			"category": helpers.ConvertToTitleCase(materi),
			"materi":   helpers.ConvertToTitleCase(pr.Materi),
			"pre_test": respRe,
		}
		finalResult[materi] = res
	}

	for cat, fn := range finalResult {
		for co, po := range resPost {
			if cat == co {
				respRe := ResultDetails{
					Score:    float64(po.Score),
					MaxScore: float64(po.ScoreMax),
					Date:     po.DateStart,
				}
				fn.(map[string]any)["post_test"] = respRe
			}
		}
	}
	mapped := []any{}
	for _, fin := range finalResult {
		mapped = append(mapped, fin)
	}

	total := len(resPost) + len(resPre)
	donePerc := helpers.RoundFloat(float64((total*100))/float64(14), 1)
	if donePerc > 100 {
		donePerc = 100
	}
	if _isNaNorInf(donePerc) {
		donePerc = 0
	}

	reportData := map[string]any{
		"data":            mapped,
		"total_pre_test":  7,
		"total_post_test": 7,
		"count_pre_test":  len(resPre),
		"count_post_test": len(resPost),
		"percentage":      donePerc,
	}

	return reportData, nil
}

func FetchPTKPrePostRaport(smId uint) (map[string]any, error) {
	setMembers, err := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptk:user_%d:pre-test", smId)).Result()
	if err != nil {
		return nil, err
	}
	setMembersPost, err := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_ptk:user_%d:post-test", smId)).Result()
	if err != nil {
		return nil, err
	}
	resPre, err := helpers.FormatingCategory(setMembers)
	if err != nil {
		return nil, err
	}

	resPost, err := helpers.FormatingCategory(setMembersPost)
	if err != nil {
		return nil, err
	}

	type ResultDetails struct {
		MaxScore float64 `json:"max_score"`
		Score    float64 `json:"score"`
		Date     string  `json:"date"`
	}

	finalResult := map[string]any{}

	for materi, pr := range resPre {
		respRe := ResultDetails{
			Score:    float64(pr.SubCategories.Score),
			MaxScore: float64(pr.SubCategories.ScoreMax),
			Date:     pr.SubCategories.DateStart,
		}
		res := map[string]any{
			"category": pr.Category,
			"materi":   helpers.ConvertToTitleCase(materi),
			"pre_test": respRe,
		}
		finalResult[materi] = res
	}

	for cat, fn := range finalResult {
		for co, po := range resPost {
			if cat == co {
				respRe := ResultDetails{
					Score:    float64(po.SubCategories.Score),
					MaxScore: float64(po.SubCategories.ScoreMax),
					Date:     po.SubCategories.DateStart,
				}
				fn.(map[string]any)["post_test"] = respRe
			}
		}
	}
	mapped := []any{}
	for _, fin := range finalResult {
		mapped = append(mapped, fin)
	}
	total := len(resPost) + len(resPre)
	donePerc := helpers.RoundFloat(float64((total*100))/float64(42), 1)
	if donePerc > 100 {
		donePerc = 100
	}
	if _isNaNorInf(donePerc) {
		donePerc = 0
	}

	reportData := map[string]any{
		"data":            mapped,
		"total_pre_test":  21,
		"total_post_test": 21,
		"count_pre_test":  len(resPre),
		"count_post_test": len(resPost),
		"percentage":      donePerc,
	}

	return reportData, nil
}

func FetchCPNSPrePostRaport(smId uint) (map[string]any, error) {
	setMembers, err := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_cpns:user_%d:pre-test", smId)).Result()
	if err != nil {
		return nil, err
	}

	setMembersPost, err := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_cpns:user_%d:post-test", smId)).Result()
	if err != nil {
		return nil, err
	}

	resPre, err := helpers.FormatingCategory(setMembers)
	if err != nil {
		return nil, err
	}

	resPost, err := helpers.FormatingCategory(setMembersPost)
	if err != nil {
		return nil, err
	}

	type ResultDetails struct {
		MaxScore float64 `json:"max_score"`
		Score    float64 `json:"score"`
		Date     string  `json:"date"`
	}

	finalResult := map[string]any{}

	for materi, pr := range resPre {
		respRe := ResultDetails{
			Score:    float64(pr.SubCategories.Score),
			MaxScore: float64(pr.SubCategories.ScoreMax),
			Date:     pr.SubCategories.DateStart,
		}
		res := map[string]any{
			"category": pr.Category,
			"materi":   helpers.ConvertToTitleCase(materi),
			"pre_test": respRe,
		}
		finalResult[materi] = res
	}

	for cat, fn := range finalResult {
		for co, po := range resPost {
			if cat == co {
				respRe := ResultDetails{
					Score:    float64(po.SubCategories.Score),
					MaxScore: float64(po.SubCategories.ScoreMax),
					Date:     po.SubCategories.DateStart,
				}
				fn.(map[string]any)["post_test"] = respRe
			}
		}
	}
	mapped := []any{}
	for _, fin := range finalResult {
		mapped = append(mapped, fin)
	}

	total := len(resPost) + len(resPre)
	donePerc := helpers.RoundFloat(float64((total*100))/float64(42), 1)
	if donePerc > 100 {
		donePerc = 100
	}
	if _isNaNorInf(donePerc) {
		donePerc = 0
	}

	reportData := map[string]any{
		"data":            mapped,
		"total_pre_test":  21,
		"total_post_test": 21,
		"count_pre_test":  len(resPre),
		"count_post_test": len(resPost),
		"percentage":      donePerc,
	}

	return reportData, nil
}

type ExamResult struct {
	DateStart string `json:"date_start"`
	ExamType  string `json:"exam_type"`
	Materi    string `json:"materi"`
	Program   string `json:"program"`
	Score     int    `json:"score"`
	ScoreMax  int    `json:"score_max"`
}

func mergeResults(resultPre map[string]interface{}, resultPost ExamResult) map[string]interface{} {
	return map[string]interface{}{
		"pre_test": map[string]interface{}{
			"date_start": resultPre["date_start"],
			"exam_type":  "pre-test",
			"score":      resultPre["score"],
			"score_max":  resultPre["score_max"],
		},
		"post_test": map[string]interface{}{
			"date_start": resultPost.DateStart,
			"exam_type":  "post-test",
			"score":      resultPost.Score,
			"score_max":  resultPost.ScoreMax,
		},
	}
}

func FetchCPNSUKARaport(smId uint, typStg, fil string) (map[string]any, error) {

	res, err := GetStudentProfileElastic(int(smId))
	if err != nil {
		return nil, err
	}
	filter := ""

	switch strings.ToUpper(fil) {
	case "PRE_UKA":
		filter = "pre-uka"
	case "ALL_MODULE":
		filter = "all-module"
	case "UKA_STAGE":
		filter = "challenge-uka"
	}

	target, err := GetStudentTargetCPNS(int(smId))
	if err != nil {
		return nil, err
	}
	form, err := GetCompetitionFormationCPNS(mockstruct.GetCompetitionCPNS{
		FormationType: target.FormationType,
		PositionID:    uint(target.PositionID),
		FormationCode: target.FormationCode,
	})
	if err != nil {
		return nil, err
	}

	gender := "Perempuan"
	if res.Gender == "L" {
		gender = "Laki-laki"
	}

	targetFormat := "02/01/2006"
	formattedTime := res.BirthDate.Format(targetFormat)

	studentData := map[string]any{
		"name":                  res.Name,
		"email":                 res.Email,
		"program":               "CPNS",
		"instance":              target.InstanceName,
		"position":              target.PositionName,
		"quota":                 form.Quota,
		"target_score":          target.TargetScore,
		"raport_date":           time.Now().Format("2006-01-02"),
		"raport_formatted_date": time.Now().Format("Monday, 2 January 2006"),
		"gender":                gender,
		"branch_name":           res.BranchName,
		"birth_date":            formattedTime,
		"phone":                 res.Phone,
	}

	hisRes, er := GetHistoryCPNSElasticFetchStudent(uint(res.SmartbtwID), typStg, filter)
	if er != nil {
		return nil, er
	}

	lrRecord, er := FetchStudentLearningRecordCPNS(smId)
	if er != nil {
		return nil, er
	}

	sort.SliceStable(hisRes, func(i, j int) bool {
		return hisRes[i].Start.After(hisRes[j].End)
	})

	packageIdList := map[int]bool{}

	var scoreTotal float64
	var scoreTWK float64
	var scoreTIU float64
	var scoreTKP float64
	var avgScore float64
	totalPassed := int(0)
	totalFailed := int(0)
	twkPass := int(0)
	tiuPass := int(0)
	tkpPass := int(0)
	twkFailed := int(0)
	tiuFailed := int(0)
	tkpFailed := int(0)
	twkTimeConsumed := int(0)
	tiuTimeConsumed := int(0)
	tkpTimeConsumed := int(0)
	withCode := int(0)
	stgTotal := int(0)
	totalStages := 32

	studentHistory := []map[string]any{}
	dn := 0
	for _, t := range hisRes {
		if t.ExamName == "" {
			continue
		}
		// if strings.Contains(t.ExamName, "Pre") {
		// 	continue
		// }
		// if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
		// 	continue
		// }
		dn += 1
		isPass := false
		scoreTotal += t.Total
		scoreTWK += t.Twk
		scoreTIU += t.Tiu
		scoreTKP += t.Tkp
		if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
			totalPassed += 1
			isPass = true
		} else {
			totalFailed += 1
		}

		if t.Twk >= t.TwkPass {
			twkPass += 1
		}
		if t.Tiu >= t.TiuPass {
			tiuPass += 1
		}
		if t.Tkp >= t.TkpPass {
			tkpPass += 1
		}
		if t.Twk < t.TwkPass {
			twkFailed += 1
		}
		if t.Tiu < t.TiuPass {
			tiuFailed += 1
		}
		if t.Tkp < t.TkpPass {
			tkpFailed += 1
		}
		if t.PackageType == "WITH_CODE" {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				withCode += 1
			}
		} else {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				stgTotal += 1
			}
		}
		if t.PackageType != "challenge-uka" && t.PackageType != "pre-uka" {
			totalStages += 1
		}
		twkTimeConsumed += t.TwkTimeConsumed
		tiuTimeConsumed += t.TiuTimeConsumed
		tkpTimeConsumed += t.TkpTimeConsumed

		studentHistory = append(studentHistory, map[string]any{
			"exam_name":          t.ExamName,
			"start_date":         t.Start,
			"task_id":            t.TaskID,
			"duration":           t.End.Sub(t.Start).Milliseconds(),
			"duration_formatted": fmt.Sprintf("%02d:%02d:%02d", int(t.End.Sub(t.Start).Hours()), int(t.End.Sub(t.Start).Minutes())%60, int(t.End.Sub(t.Start).Seconds())%60),
			"date_formatted":     t.Start.Format("02/01/2006"),
			"time_formatted":     t.Start.Format("15:04:05"),
			"total":              t.Total,
			"status":             isPass,
			"category": map[string]any{
				"twk": map[string]any{
					"score":              t.Twk,
					"passing_grade":      t.TwkPass,
					"is_pass":            t.Twk >= t.TwkPass,
					"duration":           t.TwkTimeConsumed,
					"formatted_duration": fmt.Sprintf("%02d:%02d", int((time.Duration(t.TwkTimeConsumed)*time.Millisecond).Minutes())%60, int((time.Duration(t.TwkTimeConsumed)*time.Millisecond).Seconds())%60),
				},
				"tiu": map[string]any{
					"score":              t.Tiu,
					"passing_grade":      t.TiuPass,
					"is_pass":            t.Tiu >= t.TiuPass,
					"duration":           t.TiuTimeConsumed,
					"formatted_duration": fmt.Sprintf("%02d:%02d", int((time.Duration(t.TiuTimeConsumed)*time.Millisecond).Minutes())%60, int((time.Duration(t.TiuTimeConsumed)*time.Millisecond).Seconds())%60),
				},
				"tkp": map[string]any{
					"score":              t.Tkp,
					"passing_grade":      t.TkpPass,
					"is_pass":            t.Tkp >= t.TkpPass,
					"duration":           t.TkpTimeConsumed,
					"formatted_duration": fmt.Sprintf("%02d:%02d", int((time.Duration(t.TkpTimeConsumed)*time.Millisecond).Minutes())%60, int((time.Duration(t.TkpTimeConsumed)*time.Millisecond).Seconds())%60),
				},
			},
		})

	}

	totalHis := dn
	avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
	avgScoreTWK := helpers.RoundFloat(scoreTWK/float64(totalHis), 1)
	avgScoreTIU := helpers.RoundFloat(scoreTIU/float64(totalHis), 1)
	avgScoreTKP := helpers.RoundFloat(scoreTKP/float64(totalHis), 1)

	passPercent := helpers.RoundFloat(float64(avgScore)/float64(target.TargetScore)*100, 1)
	passPercentTWK := helpers.RoundFloat(float64(twkPass)/float64(totalHis)*100, 1)
	passPercentTIU := helpers.RoundFloat(float64(tiuPass)/float64(totalHis)*100, 1)
	passPercentTKP := helpers.RoundFloat(float64(tkpPass)/float64(totalHis)*100, 1)

	if _isNaNorInf(avgScore) {
		avgScore = 0
	}
	if _isNaNorInf(avgScoreTWK) {
		avgScoreTWK = 0
	}
	if _isNaNorInf(avgScoreTIU) {
		avgScoreTIU = 0
	}
	if _isNaNorInf(avgScoreTKP) {
		avgScoreTKP = 0
	}

	if _isNaNorInf(passPercent) {
		passPercent = 0
	}
	if _isNaNorInf(passPercentTWK) {
		passPercentTWK = 0
	}
	if _isNaNorInf(passPercentTIU) {
		passPercentTIU = 0
	}
	if _isNaNorInf(passPercentTKP) {
		passPercentTKP = 0
	}
	prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
	if err != nil {
		return nil, err
	}

	countStgLv := 0
	if typStg == "UMUM" {

		for _, obj := range prodRe.Data {
			isSkipped := false
			if strings.ToLower(obj.ProductProgram) != "skd" {
				continue
			}
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
					isSkipped = true
					continue
				}
				if strings.Contains(tag, "MATERIAL") {
					isSkipped = true
					continue
				}
			}

			if isSkipped {
				continue
			}
			isSkipped = true

			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "CPNS") {
					isSkipped = false
					continue
				}
			}
			if isSkipped {
				continue
			}

			if fil != "" {
				if filter == "pre-uka" {
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_PRE_UKA") {
							countStgLv += 1
						}
					}
				} else if filter == "challenge-uka" {
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
							countStgLv += 1
						}
					}
				} else if filter == "all-module" {
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
							countStgLv += 1
						}
					}
					for _, tag := range obj.ProductTags {
						if strings.Contains(tag, "STAGE_PRE_UKA") {
							countStgLv += 1
						}
					}
				}
			} else {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_LEVEL_") {
						countStgLv += 1
					}
				}
			}
			// for _, tag := range obj.ProductTags {
			// 	if strings.Contains(tag, "STAGE_LEVEL_") {
			// 		countStgLv += 1
			// 	}
			// }
		}
	} else {
		countStgLv = 0
	}

	ownMod := withCode

	if (countStgLv - stgTotal) < 0 {
		ownMod += stgTotal
	} else {
		ownMod += countStgLv
	}

	twkQstTotal := 30 * 32
	tiuQstTotal := 35 * 32
	tkpQstTotal := 45 * 32

	twkPercentageApproached := float64(0)
	tiuPercentageApproached := float64(0)
	tkpPercentageApproached := float64(0)

	twkRecommended := false
	tiuRecommended := false
	tkpRecommended := false

	if twkTimeConsumed >= (2100000-210000) && twkTimeConsumed <= (2100000+210000) {
		twkRecommended = true
	}
	if tiuTimeConsumed >= (2100000-210000) && tiuTimeConsumed <= (2100000+210000) {
		tiuRecommended = true
	}
	if tkpTimeConsumed >= (1500000-150000) && tkpTimeConsumed <= (1500000+150000) {
		tkpRecommended = true
	}

	if lrRecord.TWKLearningApproached > 0 {
		twkPercentageApproached = (float64((float64(lrRecord.TWKLearningApproached) / float64(twkQstTotal)) * 100))
		twkPercentageApproached = (math.Round(twkPercentageApproached*1000) / 1000)
	}
	if lrRecord.TIULearningApproached > 0 {
		tiuPercentageApproached = (float64((float64(lrRecord.TIULearningApproached) / float64(tiuQstTotal)) * 100))
		tiuPercentageApproached = (math.Round(tiuPercentageApproached*1000) / 1000)
	}
	if lrRecord.TKPLearningApproached > 0 {
		tkpPercentageApproached = (float64((float64(lrRecord.TKPLearningApproached) / float64(tkpQstTotal)) * 100))
		tkpPercentageApproached = (math.Round(tkpPercentageApproached*1000) / 1000)

	}
	if _isNaNorInf(twkPercentageApproached) {
		twkPercentageApproached = 0
	}
	if _isNaNorInf(tiuPercentageApproached) {
		tiuPercentageApproached = 0
	}
	if _isNaNorInf(tkpPercentageApproached) {
		tkpPercentageApproached = 0
	}

	donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)
	avgDoneScore := helpers.RoundFloat(scoreTotal/float64(dn), 1)
	if _isNaNorInf(donePerct) {
		donePerct = 0
	}

	if _isNaNorInf(avgDoneScore) {
		avgDoneScore = 0
	}
	if donePerct > 100 {
		donePerct = 100
	}

	if passPercent > 100 {
		passPercent = 100
	}

	setMembers, errs := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_cpns:user_%d:pre-test", smId)).Result()
	if errs != nil {
		return nil, err
	}
	setMembersPost, errs := db.NewRedisCluster().SMembers(context.Background(), fmt.Sprintf("exam:result-test:program_cpns:user_%d:post-test", smId)).Result()
	if errs != nil {
		return nil, err
	}

	pre, err := helpers.FormatingCategory(setMembers)
	if err != nil {
		return nil, err
	}
	newPre := map[string]helpers.DataStruct{}
	for pe := range pre {
		peN := helpers.ToLowerAndUnderscore(pe)
		newPre[peN] = pre[pe]
	}

	post, err := helpers.FormatingCategory(setMembersPost)
	if err != nil {
		return nil, err
	}

	newPost := map[string]helpers.DataStruct{}
	for po := range post {
		poN := helpers.ToLowerAndUnderscore(po)
		newPost[poN] = post[po]
	}

	reportData := map[string]any{
		"student": studentData,
		"summary": map[string]any{
			"category": map[string]any{
				"twk": map[string]any{
					"average":                         avgScoreTWK,
					"passing_percentage":              passPercentTWK,
					"approached_explanation":          lrRecord.TWKLearningApproached,
					"total_explanation":               twkQstTotal,
					"percentage_explanation":          twkPercentageApproached,
					"calculated_answer_time":          twkTimeConsumed / totalStages,
					"calculated_recommended":          twkRecommended,
					"calculated_formated_answer_time": fmt.Sprintf("%02d:%02d", int((time.Duration(twkTimeConsumed/totalStages)*time.Millisecond).Minutes())%60, int((time.Duration(twkTimeConsumed/totalStages)*time.Millisecond).Seconds())%60),
				},
				"tiu": map[string]any{
					"average":                         avgScoreTIU,
					"passing_percentage":              passPercentTIU,
					"approached_explanation":          lrRecord.TIULearningApproached,
					"total_explanation":               tiuQstTotal,
					"percentage_explanation":          tiuPercentageApproached,
					"calculated_answer_time":          tiuTimeConsumed / totalStages,
					"calculated_recommended":          tiuRecommended,
					"calculated_formated_answer_time": fmt.Sprintf("%02d:%02d", int((time.Duration(tiuTimeConsumed/totalStages)*time.Millisecond).Minutes())%60, int((time.Duration(tiuTimeConsumed/totalStages)*time.Millisecond).Seconds())%60),
				},
				"tkp": map[string]any{
					"average":                         avgScoreTKP,
					"passing_percentage":              passPercentTKP,
					"approached_explanation":          lrRecord.TKPLearningApproached,
					"total_explanation":               tkpQstTotal,
					"percentage_explanation":          tkpPercentageApproached,
					"calculated_answer_time":          tkpTimeConsumed / totalStages,
					"calculated_recommended":          tkpRecommended,
					"calculated_formated_answer_time": fmt.Sprintf("%02d:%02d", int((time.Duration(tkpTimeConsumed/totalStages)*time.Millisecond).Minutes())%60, int((time.Duration(tkpTimeConsumed/totalStages)*time.Millisecond).Seconds())%60),
				},
			},
			"received":             ownMod,
			"passed":               totalPassed,
			"completed":            dn,
			"completed_percentage": donePerct,
			"passing_percentage":   passPercent,
			"average_score":        avgScore,
			"average_total":        avgDoneScore,
		},
		"histories": studentHistory,
		// "assessment": overallResult,
		"pre_test":  newPre,
		"post_test": newPost,
	}

	return reportData, nil
}

func FetchPTNRaport(smId uint, fil string) (map[string]any, error) {

	res, err := GetStudentProfileElastic(int(smId))
	if err != nil {
		return nil, err
	}
	filter := ""

	switch strings.ToUpper(fil) {
	case "PRE_UKA":
		filter = "pre-uka"
	case "UKA_CODE":
		filter = "with_code"
	case "CHALLENGE_UKA":
		filter = "challenge-uka"
	}

	studentTarget, err := GetStudentProfilePTNElastic(res.SmartbtwID)
	if err != nil {
		return nil, err
	}

	comp, err := GetCompetitonPTN(uint(studentTarget.MajorID))
	if err != nil {
		return nil, err
	}

	studentData := map[string]any{
		"name":                  res.Name,
		"email":                 res.Email,
		"program":               "PTN",
		"major_name":            studentTarget.MajorName,
		"school_name":           studentTarget.SchoolName,
		"quota":                 comp.SbmptnCapacity,
		"target_score":          studentTarget.TargetScore,
		"raport_date":           time.Now().Format("2006-01-02"),
		"raport_formatted_date": time.Now().Format("Monday, 2 January 2006"),
	}

	hisRes, er := GetStudentHistoryPTNElasticFilterOld(int(res.SmartbtwID), filter, "utbk")
	if er != nil {
		return nil, er
	}

	sort.SliceStable(hisRes, func(i, j int) bool {
		return hisRes[i].Start.After(hisRes[j].End)
	})

	packageIdList := map[int]bool{}

	var scoreTotal float64
	var scorePU float64
	var scorePPU float64
	var scorePBM float64
	var scorePM float64
	var scorePK float64
	var scoreLBIND float64
	var scoreLBING float64
	var avgScore float64
	withCode := int(0)
	stgTotal := int(0)

	studentHistory := []map[string]any{}
	dn := 0
	for _, t := range hisRes {
		if t.ExamName == "" {
			continue
		}
		if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
			continue
		}

		if strings.Contains(t.ExamName, "Pre") || strings.Contains(t.ExamName, "Post") {
			continue
		}
		dn += 1
		isPass := false
		scoreTotal += t.Total
		scorePU += t.PenalaranUmum
		scorePBM += t.PemahamanBacaan
		scorePPU += t.PengetahuanUmum
		scorePM += t.PenalaranMatematika
		scorePK += t.PengetahuanKuantitatif
		scoreLBIND += t.LiterasiBahasaIndonesia
		scoreLBING += t.LiterasiBahasaInggris

		if t.PackageType == "WITH_CODE" {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				withCode += 1
			}
		} else {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				stgTotal += 1
			}
		}

		studentHistory = append(studentHistory, map[string]any{
			"exam_name":          t.ExamName,
			"start_date":         t.Start,
			"duration":           t.End.Sub(t.Start).Milliseconds(),
			"duration_formatted": fmt.Sprintf("%02d:%02d:%02d", int(t.End.Sub(t.Start).Hours()), int(t.End.Sub(t.Start).Minutes())%60, int(t.End.Sub(t.Start).Seconds())%60),
			"date_formatted":     t.Start.Format("02/01/2006"),
			"time_formatted":     t.Start.Format("15:04:05"),
			"total":              t.Total,
			"status":             isPass,
			"category": map[string]any{
				"pu": map[string]any{
					"score": t.PenalaranUmum,
				},
				"ppu": map[string]any{
					"score": t.PengetahuanUmum,
				},
				"pbm": map[string]any{
					"score": t.PemahamanBacaan,
				},
				"pk": map[string]any{
					"score": t.PengetahuanKuantitatif,
				},
				"pm": map[string]any{
					"score": t.PenalaranMatematika,
				},
				"lbind": map[string]any{
					"score": t.LiterasiBahasaIndonesia,
				},
				"lbing": map[string]any{
					"score": t.LiterasiBahasaInggris,
				},
			},
		})

	}

	totalHis := dn
	avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
	avgScorePU := helpers.RoundFloat(scorePU/float64(totalHis), 1)
	avgScorePPU := helpers.RoundFloat(scorePPU/float64(totalHis), 1)
	avgScorePBM := helpers.RoundFloat(scorePBM/float64(totalHis), 1)
	avgScorePK := helpers.RoundFloat(scorePK/float64(totalHis), 1)
	avgScorePM := helpers.RoundFloat(scorePM/float64(totalHis), 1)
	avgScoreLBIND := helpers.RoundFloat(scoreLBIND/float64(totalHis), 1)
	avgScoreLBING := helpers.RoundFloat(scoreLBING/float64(totalHis), 1)

	prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
	if err != nil {
		return nil, err
	}
	if _isNaNorInf(avgScore) {
		avgScore = 0
	}
	if _isNaNorInf(avgScorePU) {
		avgScorePU = 0
	}
	if _isNaNorInf(avgScorePPU) {
		avgScorePPU = 0
	}
	if _isNaNorInf(avgScorePK) {
		avgScorePK = 0
	}
	if _isNaNorInf(avgScorePBM) {
		avgScorePBM = 0
	}
	if _isNaNorInf(avgScorePM) {
		avgScorePM = 0
	}
	if _isNaNorInf(avgScoreLBIND) {
		avgScoreLBIND = 0
	}
	if _isNaNorInf(avgScoreLBING) {
		avgScoreLBING = 0
	}

	countStgLv := 0
	for _, obj := range prodRe.Data {
		isSkipped := false
		if strings.ToLower(obj.ProductProgram) != "tps" || strings.ToLower(obj.ProductProgram) != "ptn" || strings.ToLower(obj.ProductProgram) != "utbk" {
			continue
		}
		for _, tag := range obj.ProductTags {
			if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "MATERIAL") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "CPNS") {
				isSkipped = true
				continue
			}
		}
		if isSkipped {
			continue
		}
		if fil != "" {
			if filter == "pre-uka" {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_PRE_UKA") {
						countStgLv += 1
					}
				}
			} else if filter == "challenge-uka" {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
						countStgLv += 1
					}
				}
			}
		} else {
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "STAGE_LEVEL_") {
					countStgLv += 1
				}
			}
		}
	}
	ownMod := withCode

	if (countStgLv - stgTotal) < 0 {
		ownMod += stgTotal
	} else {
		ownMod += countStgLv
	}

	donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)
	passPercent := helpers.RoundFloat(float64(totalHis)/float64(totalHis)*100, 1)

	avgDoneScore := helpers.RoundFloat(scoreTotal/float64(dn), 1)
	if _isNaNorInf(donePerct) {
		donePerct = 0
	}
	if _isNaNorInf(passPercent) {
		passPercent = 0
	}
	if _isNaNorInf(avgDoneScore) {
		avgDoneScore = 0
	}
	reportData := map[string]any{
		"student": studentData,
		"summary": map[string]any{
			"category": map[string]any{
				"pu": map[string]any{
					"average": avgScorePU,
				},
				"ppu": map[string]any{
					"average": avgScorePPU,
				},
				"pbm": map[string]any{
					"average": avgScorePBM,
				},
				"pk": map[string]any{
					"average": avgScorePK,
				},
				"pm": map[string]any{
					"average": avgScorePM,
				},
				"lbind": map[string]any{
					"average": avgScoreLBIND,
				},
				"lbing": map[string]any{
					"average": avgScoreLBING,
				},
			},
			"received":             ownMod,
			"completed":            totalHis,
			"passed":               totalHis,
			"completed_percentage": donePerct,
			"passing_percentage":   passPercent,
			"average_score":        avgScore,
			"average_total":        avgDoneScore,
		},
		"histories": studentHistory,
	}
	return reportData, nil
}

func FetchCPNSRaport(smId uint, fil string) (map[string]any, error) {

	res, err := GetStudentProfileElastic(int(smId))
	if err != nil {
		return nil, err
	}
	filter := ""

	switch strings.ToUpper(fil) {
	case "PRE_UKA":
		filter = "pre-uka"
	case "UKA_CODE":
		filter = "with_code"
	case "CHALLENGE_UKA":
		filter = "challenge-uka"
	}

	target, err := GetStudentTargetCPNS(int(smId))
	if err != nil {
		return nil, err
	}
	form, err := GetCompetitionFormationCPNS(mockstruct.GetCompetitionCPNS{
		FormationType: target.FormationType,
		PositionID:    uint(target.PositionID),
		FormationCode: target.FormationCode,
	})
	if err != nil {
		return nil, err
	}

	studentData := map[string]any{
		"name":                  res.Name,
		"email":                 res.Email,
		"program":               "CPNS",
		"instance":              target.InstanceName,
		"position":              target.Position,
		"quota":                 form.Quota,
		"target_score":          target.TargetScore,
		"raport_date":           time.Now().Format("2006-01-02"),
		"raport_formatted_date": time.Now().Format("Monday, 2 January 2006"),
	}

	hisRes, er := GetStudentHistoryCPNSElasticFilterOld(int(res.SmartbtwID), filter)
	if er != nil {
		return nil, er
	}

	lrRecord, er := FetchStudentLearningRecordCPNS(smId)
	if er != nil {
		return nil, er
	}

	sort.SliceStable(hisRes, func(i, j int) bool {
		return hisRes[i].Start.After(hisRes[j].End)
	})

	packageIdList := map[int]bool{}

	var scoreTotal float64
	var scoreTWK float64
	var scoreTIU float64
	var scoreTKP float64
	var avgScore float64
	totalPassed := int(0)
	totalFailed := int(0)
	twkPass := int(0)
	tiuPass := int(0)
	tkpPass := int(0)
	twkFailed := int(0)
	tiuFailed := int(0)
	tkpFailed := int(0)
	twkTimeConsumed := int(0)
	tiuTimeConsumed := int(0)
	tkpTimeConsumed := int(0)
	withCode := int(0)
	stgTotal := int(0)
	totalStages := 32

	studentHistory := []map[string]any{}
	dn := 0
	for _, t := range hisRes {
		if t.ExamName == "" {
			continue
		}
		if t.ModuleType == "PRE_TEST" || t.ModuleType == "POST_TEST" {
			continue
		}

		if strings.Contains(t.ExamName, "Pre") || strings.Contains(t.ExamName, "Post") {
			continue
		}
		dn += 1
		isPass := false
		scoreTotal += t.Total
		scoreTWK += t.Twk
		scoreTIU += t.Tiu
		scoreTKP += t.Tkp
		if t.Twk >= t.TwkPass && t.Tiu >= t.TiuPass && t.Tkp >= t.TkpPass {
			totalPassed += 1
			isPass = true
		} else {
			totalFailed += 1
		}

		if t.Twk >= t.TwkPass {
			twkPass += 1
		}
		if t.Tiu >= t.TiuPass {
			tiuPass += 1
		}
		if t.Tkp >= t.TkpPass {
			tkpPass += 1
		}
		if t.Twk < t.TwkPass {
			twkFailed += 1
		}
		if t.Tiu < t.TiuPass {
			tiuFailed += 1
		}
		if t.Tkp < t.TkpPass {
			tkpFailed += 1
		}
		if t.PackageType == "WITH_CODE" {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				withCode += 1
			}
		} else {
			_, k := packageIdList[t.PackageID]
			if !k {
				packageIdList[t.PackageID] = true
				stgTotal += 1
			}
		}
		if t.PackageType != "challenge-uka" && t.PackageType != "pre-uka" {
			totalStages += 1
		}
		twkTimeConsumed += t.TwkTimeConsumed
		tiuTimeConsumed += t.TiuTimeConsumed
		tkpTimeConsumed += t.TkpTimeConsumed

		studentHistory = append(studentHistory, map[string]any{
			"exam_name":          t.ExamName,
			"start_date":         t.Start,
			"duration":           t.End.Sub(t.Start).Milliseconds(),
			"duration_formatted": fmt.Sprintf("%02d:%02d:%02d", int(t.End.Sub(t.Start).Hours()), int(t.End.Sub(t.Start).Minutes())%60, int(t.End.Sub(t.Start).Seconds())%60),
			"date_formatted":     t.Start.Format("02/01/2006"),
			"time_formatted":     t.Start.Format("15:04:05"),
			"total":              t.Total,
			"status":             isPass,
			"category": map[string]any{
				"twk": map[string]any{
					"score":              t.Twk,
					"passing_grade":      t.TwkPass,
					"is_pass":            t.Twk >= t.TwkPass,
					"duration":           t.TwkTimeConsumed,
					"formatted_duration": fmt.Sprintf("%02d:%02d", int((time.Duration(t.TwkTimeConsumed)*time.Millisecond).Minutes())%60, int((time.Duration(t.TwkTimeConsumed)*time.Millisecond).Seconds())%60),
				},
				"tiu": map[string]any{
					"score":              t.Tiu,
					"passing_grade":      t.TiuPass,
					"is_pass":            t.Tiu >= t.TiuPass,
					"duration":           t.TiuTimeConsumed,
					"formatted_duration": fmt.Sprintf("%02d:%02d", int((time.Duration(t.TiuTimeConsumed)*time.Millisecond).Minutes())%60, int((time.Duration(t.TiuTimeConsumed)*time.Millisecond).Seconds())%60),
				},
				"tkp": map[string]any{
					"score":              t.Tkp,
					"passing_grade":      t.TkpPass,
					"is_pass":            t.Tkp >= t.TkpPass,
					"duration":           t.TkpTimeConsumed,
					"formatted_duration": fmt.Sprintf("%02d:%02d", int((time.Duration(t.TkpTimeConsumed)*time.Millisecond).Minutes())%60, int((time.Duration(t.TkpTimeConsumed)*time.Millisecond).Seconds())%60),
				},
			},
		})

	}

	totalHis := dn
	avgScore = helpers.RoundFloat(scoreTotal/float64(totalHis), 1)
	avgScoreTWK := helpers.RoundFloat(scoreTWK/float64(totalHis), 1)
	avgScoreTIU := helpers.RoundFloat(scoreTIU/float64(totalHis), 1)
	avgScoreTKP := helpers.RoundFloat(scoreTKP/float64(totalHis), 1)

	passPercent := helpers.RoundFloat(float64(totalPassed)/float64(totalHis)*100, 1)
	passPercentTWK := helpers.RoundFloat(float64(twkPass)/float64(totalHis)*100, 1)
	passPercentTIU := helpers.RoundFloat(float64(tiuPass)/float64(totalHis)*100, 1)
	passPercentTKP := helpers.RoundFloat(float64(tkpPass)/float64(totalHis)*100, 1)

	if _isNaNorInf(avgScore) {
		avgScore = 0
	}
	if _isNaNorInf(avgScoreTWK) {
		avgScoreTWK = 0
	}
	if _isNaNorInf(avgScoreTIU) {
		avgScoreTIU = 0
	}
	if _isNaNorInf(avgScoreTKP) {
		avgScoreTKP = 0
	}

	if _isNaNorInf(passPercent) {
		passPercent = 0
	}
	if _isNaNorInf(passPercentTWK) {
		passPercentTWK = 0
	}
	if _isNaNorInf(passPercentTIU) {
		passPercentTIU = 0
	}
	if _isNaNorInf(passPercentTKP) {
		passPercentTKP = 0
	}
	prodRe, err := GetStudentAOP(uint(res.SmartbtwID))
	if err != nil {
		return nil, err
	}

	countStgLv := 0
	for _, obj := range prodRe.Data {
		isSkipped := false
		if strings.ToLower(obj.ProductProgram) != "skd" {
			continue
		}
		for _, tag := range obj.ProductTags {
			if strings.Contains(tag, "GOLD") || strings.Contains(tag, "DIAMOND") {
				isSkipped = true
				continue
			}
			if strings.Contains(tag, "MATERIAL") {
				isSkipped = true
				continue
			}
		}

		if isSkipped {
			continue
		}
		isSkipped = true

		for _, tag := range obj.ProductTags {
			if strings.Contains(tag, "CPNS") {
				isSkipped = false
				continue
			}
		}
		if isSkipped {
			continue
		}

		if fil != "" {
			if filter == "pre-uka" {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_PRE_UKA") {
						countStgLv += 1
					}
				}
			} else if filter == "challenge-uka" {
				for _, tag := range obj.ProductTags {
					if strings.Contains(tag, "STAGE_CHALLENGE_UKA") || strings.Contains(tag, "STAGE_UKA") {
						countStgLv += 1
					}
				}
			}
		} else {
			for _, tag := range obj.ProductTags {
				if strings.Contains(tag, "STAGE_LEVEL_") {
					countStgLv += 1
				}
			}
		}
		// for _, tag := range obj.ProductTags {
		// 	if strings.Contains(tag, "STAGE_LEVEL_") {
		// 		countStgLv += 1
		// 	}
		// }
	}

	ownMod := withCode

	if (countStgLv - stgTotal) < 0 {
		ownMod += stgTotal
	} else {
		ownMod += countStgLv
	}

	twkQstTotal := 30 * 32
	tiuQstTotal := 35 * 32
	tkpQstTotal := 45 * 32

	twkPercentageApproached := float64(0)
	tiuPercentageApproached := float64(0)
	tkpPercentageApproached := float64(0)

	twkRecommended := false
	tiuRecommended := false
	tkpRecommended := false

	if twkTimeConsumed >= (2100000-210000) && twkTimeConsumed <= (2100000+210000) {
		twkRecommended = true
	}
	if tiuTimeConsumed >= (2100000-210000) && tiuTimeConsumed <= (2100000+210000) {
		tiuRecommended = true
	}
	if tkpTimeConsumed >= (1500000-150000) && tkpTimeConsumed <= (1500000+150000) {
		tkpRecommended = true
	}

	if lrRecord.TWKLearningApproached > 0 {
		twkPercentageApproached = (float64((float64(lrRecord.TWKLearningApproached) / float64(twkQstTotal)) * 100))
		twkPercentageApproached = (math.Round(twkPercentageApproached*1000) / 1000)
	}
	if lrRecord.TIULearningApproached > 0 {
		tiuPercentageApproached = (float64((float64(lrRecord.TIULearningApproached) / float64(tiuQstTotal)) * 100))
		tiuPercentageApproached = (math.Round(tiuPercentageApproached*1000) / 1000)
	}
	if lrRecord.TKPLearningApproached > 0 {
		tkpPercentageApproached = (float64((float64(lrRecord.TKPLearningApproached) / float64(tkpQstTotal)) * 100))
		tkpPercentageApproached = (math.Round(tkpPercentageApproached*1000) / 1000)

	}
	if _isNaNorInf(twkPercentageApproached) {
		twkPercentageApproached = 0
	}
	if _isNaNorInf(tiuPercentageApproached) {
		tiuPercentageApproached = 0
	}
	if _isNaNorInf(tkpPercentageApproached) {
		tkpPercentageApproached = 0
	}

	donePerct := helpers.RoundFloat(float64(totalHis)/float64(ownMod)*100, 1)
	avgDoneScore := helpers.RoundFloat(scoreTotal/float64(dn), 1)
	if _isNaNorInf(donePerct) {
		donePerct = 0
	}

	if _isNaNorInf(avgDoneScore) {
		avgDoneScore = 0
	}
	reportData := map[string]any{
		"student": studentData,
		"summary": map[string]any{
			"category": map[string]any{
				"twk": map[string]any{
					"average":                         avgScoreTWK,
					"passing_percentage":              passPercentTWK,
					"approached_explanation":          lrRecord.TWKLearningApproached,
					"total_explanation":               twkQstTotal,
					"percentage_explanation":          twkPercentageApproached,
					"calculated_answer_time":          twkTimeConsumed / totalStages,
					"calculated_recommended":          twkRecommended,
					"calculated_formated_answer_time": fmt.Sprintf("%02d:%02d", int((time.Duration(twkTimeConsumed/totalStages)*time.Millisecond).Minutes())%60, int((time.Duration(twkTimeConsumed/totalStages)*time.Millisecond).Seconds())%60),
				},
				"tiu": map[string]any{
					"average":                         avgScoreTIU,
					"passing_percentage":              passPercentTIU,
					"approached_explanation":          lrRecord.TIULearningApproached,
					"total_explanation":               tiuQstTotal,
					"percentage_explanation":          tiuPercentageApproached,
					"calculated_answer_time":          tiuTimeConsumed / totalStages,
					"calculated_recommended":          tiuRecommended,
					"calculated_formated_answer_time": fmt.Sprintf("%02d:%02d", int((time.Duration(tiuTimeConsumed/totalStages)*time.Millisecond).Minutes())%60, int((time.Duration(tiuTimeConsumed/totalStages)*time.Millisecond).Seconds())%60),
				},
				"tkp": map[string]any{
					"average":                         avgScoreTKP,
					"passing_percentage":              passPercentTKP,
					"approached_explanation":          lrRecord.TKPLearningApproached,
					"total_explanation":               tkpQstTotal,
					"percentage_explanation":          tkpPercentageApproached,
					"calculated_answer_time":          tkpTimeConsumed / totalStages,
					"calculated_recommended":          tkpRecommended,
					"calculated_formated_answer_time": fmt.Sprintf("%02d:%02d", int((time.Duration(tkpTimeConsumed/totalStages)*time.Millisecond).Minutes())%60, int((time.Duration(tkpTimeConsumed/totalStages)*time.Millisecond).Seconds())%60),
				},
			},
			"received":             ownMod,
			"passed":               totalPassed,
			"completed":            dn,
			"completed_percentage": donePerct,
			"passing_percentage":   passPercent,
			"average_score":        avgScore,
			"average_total":        avgDoneScore,
		},
		"histories": studentHistory,
	}

	return reportData, nil
}

func _isNaNorInf(val float64) bool {
	return math.IsNaN(val) || math.IsInf(val, 0)
}

func GetStudentsBySchoolOriginIDElastic(schoolOriginID string, searchQuery *string, page *int, pageSize *int) ([]request.StudentProfileElastic, int64, error) {
	ctx := context.Background()

	var results []request.StudentProfileElastic

	var from int
	var size int
	if page != nil && pageSize != nil {
		from = (*page - 1) * (*pageSize)
		size = *pageSize
	}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("last_ed_id.keyword", schoolOriginID),
	)

	if searchQuery != nil {
		query = query.Must(elastic.NewMatchQuery("name.keyword", *searchQuery))
	}

	search := db.ElasticClient.Search().
		Index(db.GetStudentProfileIndexName()).
		Query(query)

	if page != nil && pageSize != nil {
		search = search.From(from).Size(size)
	} else {
		search = search.Size(1000)
	}

	searchResult, err := search.Do(ctx)

	if err != nil {
		return nil, 0, err
	}

	totalHits := searchResult.TotalHits()

	if totalHits == 0 {
		return nil, 0, fmt.Errorf("student data not found")
	}

	for _, hit := range searchResult.Hits.Hits {
		var student request.StudentProfileElastic
		err := json.Unmarshal(hit.Source, &student)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, student)
	}

	return results, totalHits, nil
}

func GetStudentsCountBySchoolOriginIDElastic(schoolOriginID string) (int64, error) {
	ctx := context.Background()

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("last_ed_id.keyword", schoolOriginID),
	)
	return db.ElasticClient.Count().
		Index(db.GetStudentProfileIndexName()).
		Query(query).
		Do(ctx)
}

func GetStudentSchoolCount(schoolId string) (mockstruct.SchoolStudentCount, error) {
	res, _, err := GetStudentsBySchoolOriginIDElastic(schoolId, nil, nil, nil)
	if err != nil {
		return mockstruct.SchoolStudentCount{}, err
	}
	dms := mockstruct.SchoolStudentCount{
		StudentTotal:       len(res),
		StudentJoinedClass: 0,
	}
	for _, k := range res {
		stJoined, err := GetStudentJoinedClassType(k.SmartbtwID)
		if err == nil {
			if len(stJoined) > 0 {
				dms.StudentJoinedClass += 1
			}
		}
	}
	return dms, nil
}
