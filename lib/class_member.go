package lib

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func CreateClassMember(bd *request.CreateClassMember) error {
	stdCol := db.Mongodb.Collection("class_members")
	stdCol1 := db.Mongodb.Collection("classrooms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// get classroom data
	filter := bson.M{"_id": bd.ClassroomID}
	stdModels := models.Classroom{}
	err := stdCol1.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return err
	}

	// get class member data
	cmfilter := bson.M{"classroom_id": bd.ClassroomID, "smartbtw_id": bd.SmartbtwID, "deleted_at": nil}
	cmModels := models.ClassMember{}
	err = stdCol.FindOne(ctx, cmfilter).Decode(&cmModels)
	if err == nil {
		return errors.New("siswa ini terdaftar pada kelas")
	}

	// upsert class member
	pyl := models.ClassMember{
		SmartbtwID:  int32(bd.SmartbtwID),
		ClassroomID: bd.ClassroomID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	cmoption := options.Update().SetUpsert(true)
	cmupdate := bson.M{"$set": pyl}
	_, err = stdCol.UpdateOne(ctx, cmfilter, cmupdate, cmoption)
	if err != nil {
		return err
	}

	isOl := false
	if stdModels.IsOnline {
		isOl = true
	}

	//insert to elastic search
	bdElastic := models.ClassMemberElastic{
		ID:          fmt.Sprintf("%d_%s", bd.SmartbtwID, bd.ClassroomID.Hex()),
		SmartbtwID:  int32(bd.SmartbtwID),
		ClassroomID: bd.ClassroomID.Hex(),
		BranchCode:  stdModels.BranchCode,
		Quota:       stdModels.Quota,
		QuotaFilled: stdModels.QuotaFilled,
		Description: stdModels.Description,
		Tags:        stdModels.Tags,
		Year:        stdModels.Year,
		Status:      stdModels.Status,
		Title:       stdModels.Title,
		ClassCode:   stdModels.ClassCode,
		ProductID:   stdModels.ProductID,
		IsOnline:    isOl,
		CreatedAt:   time.Now(),
	}
	ctx1 := context.Background()
	_, err = db.ElasticClient.Index().
		Index(db.GetClassMemberIndexName()).
		Id(fmt.Sprintf("%d_%s", bd.SmartbtwID, bd.ClassroomID.Hex())).
		BodyJson(bdElastic).
		Do(ctx1)

	if err != nil {
		return err
	}

	return nil
}

func SyncClassMemberToElastic(smID uint) error {
	stdCol := db.Mongodb.Collection("class_members")
	stdCol1 := db.Mongodb.Collection("classrooms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//get class member
	filter := bson.M{"smartbtw_id": smID, "deleted_at": nil}
	stdModels := models.ClassMember{}
	err := stdCol.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return err
	}

	//get classroom
	filter1 := bson.M{"_id": stdModels.ClassroomID}
	stdModels1 := models.Classroom{}
	err = stdCol1.FindOne(ctx, filter1).Decode(&stdModels1)
	if err != nil {
		return err
	}

	isOl := false
	if stdModels1.IsOnline {
		isOl = true
	}

	//insert to elastic search
	bdElastic := models.ClassMemberElastic{
		ID:          fmt.Sprintf("%d_%s", smID, stdModels.ClassroomID.Hex()),
		SmartbtwID:  int32(smID),
		ClassroomID: stdModels.ClassroomID.Hex(),
		BranchCode:  stdModels1.BranchCode,
		Quota:       stdModels1.Quota,
		QuotaFilled: stdModels1.QuotaFilled,
		Description: stdModels1.Description,
		Tags:        stdModels1.Tags,
		Year:        stdModels1.Year,
		Status:      stdModels1.Status,
		Title:       stdModels1.Title,
		ClassCode:   stdModels1.ClassCode,
		ProductID:   stdModels1.ProductID,
		IsOnline:    isOl,
		CreatedAt:   stdModels.CreatedAt,
	}
	ctx1 := context.Background()
	_, err = db.ElasticClient.Index().
		Index(db.GetClassMemberIndexName()).
		Id(fmt.Sprintf("%d_%s", smID, stdModels.ClassroomID.Hex())).
		BodyJson(bdElastic).
		Do(ctx1)

	if err != nil {
		return err
	}

	return nil
}

func GetAllClassMember() ([]models.ClassMember, error) {
	stdCol := db.Mongodb.Collection("class_members")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var stdModels []models.ClassMember
	cursor, err := stdCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &stdModels)
	if err != nil {
		return nil, err
	}

	return stdModels, nil
}

func UpdateClassMember(bd *request.UpdateClassMember) error {
	stdCol := db.Mongodb.Collection("class_members")
	stdCol1 := db.Mongodb.Collection("classrooms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": bd.ClassroomIDAfter}
	stdModels := models.Classroom{}
	err := stdCol1.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return err
	}

	//soft delete old class member
	filter1 := bson.M{"smartbtw_id": bd.SmartbtwID, "classroom_id": bd.ClassroomIDBefore}
	update := bson.M{"$set": bson.M{
		"deleted_at": time.Now(),
		"updated_at": time.Now(),
	}}

	_, err1 := stdCol.UpdateOne(ctx, filter1, update)
	if err1 != nil {
		return err1
	}

	//delete old class member elastic
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("classroom_id", bd.ClassroomIDBefore.Hex()), elastic.NewMatchQuery("smartbtw_id", bd.SmartbtwID))
	_, err = elastic.NewDeleteByQueryService(db.ElasticClient).
		Index(db.GetClassMemberIndexName()).
		Query(bq).
		Do(context.Background())
	if err != nil {
		return err
	}

	_, errs := db.ElasticClient.Flush().Index(db.GetClassMemberIndexName()).Do(context.Background())
	if errs != nil {
		return err
	}

	//insert new class member
	pyl := models.ClassMember{
		SmartbtwID:  int32(bd.SmartbtwID),
		ClassroomID: bd.ClassroomIDAfter,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	_, err = stdCol.InsertOne(ctx, pyl)

	if err != nil {
		return err
	}

	isOl := false
	if stdModels.IsOnline {
		isOl = true
	}

	//insert to elastic search classmember after
	bdElastic := models.ClassMemberElastic{
		ID:          fmt.Sprintf("%d_%s", bd.SmartbtwID, bd.ClassroomIDAfter.Hex()),
		SmartbtwID:  int32(bd.SmartbtwID),
		ClassroomID: bd.ClassroomIDAfter.Hex(),
		BranchCode:  stdModels.BranchCode,
		Quota:       stdModels.Quota,
		QuotaFilled: stdModels.QuotaFilled,
		Description: stdModels.Description,
		Tags:        stdModels.Tags,
		Year:        stdModels.Year,
		Status:      stdModels.Status,
		Title:       stdModels.Title,
		ClassCode:   stdModels.ClassCode,
		ProductID:   stdModels.ProductID,
		IsOnline:    isOl,
		CreatedAt:   time.Now(),
	}
	ctx1 := context.Background()
	_, err = db.ElasticClient.Index().
		Index(db.GetClassMemberIndexName()).
		Id(fmt.Sprintf("%d_%s", bd.SmartbtwID, bd.ClassroomIDAfter.Hex())).
		BodyJson(bdElastic).
		Do(ctx1)

	if err != nil {
		return err
	}

	return nil
}

func GetSingleClassMemberFromElastic(smID int32, isOnline bool) (models.ClassMemberElastic, error) {
	classMember := models.ClassMemberElastic{}
	ctx := context.Background()

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smID),
		elastic.NewMatchQuery("is_online", isOnline),
	)
	sort := elastic.NewFieldSort("created_at").Desc()

	// Current query still use 200 as data limit to load all question
	// TODO: Try to find better solution to load all question
	res, err := db.ElasticClient.Search().
		Index(db.GetClassMemberIndexName()).
		SortBy(sort).
		Query(query).
		Size(1).
		Do(ctx)
	if err != nil {
		return models.ClassMemberElastic{}, err
	}

	var t models.ClassMemberElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		classMember = item.(models.ClassMemberElastic)
	}
	if classMember.ID == "" {
		err = errors.New("record not found")
		return models.ClassMemberElastic{}, err
	}
	return classMember, nil
}

func GetClassMemberFromElastic(clsId string) (students []request.StudentProfileElastic, err error) {
	students = []request.StudentProfileElastic{}
	ctx := context.Background()

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("classroom_id", clsId),
	)
	sort := elastic.NewFieldSort("created_at").Desc()

	// Current query still use 200 as data limit to load all question
	// TODO: Try to find better solution to load all question
	res, err := db.ElasticClient.Search().
		Index(db.GetClassMemberIndexName()).
		SortBy(sort).
		Query(query).
		Size(10000).
		Do(ctx)
	if err != nil {
		return
	}

	var t models.ClassMemberElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		classMember := item.(models.ClassMemberElastic)
		std, err := GetStudentProfileElastic(int(classMember.SmartbtwID))
		if err == nil {
			std.CreatedAt = classMember.CreatedAt
			students = append(students, std)
		}
	}
	return
}

func GetStudentJoinedClassType(stId int) (students []string, err error) {
	students = []string{}
	ctx := context.Background()

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", stId),
	)
	sort := elastic.NewFieldSort("created_at").Desc()

	// Current query still use 200 as data limit to load all question
	// TODO: Try to find better solution to load all question
	res, err := db.ElasticClient.Search().
		Index(db.GetClassMemberIndexName()).
		SortBy(sort).
		Query(query).
		Size(10000).
		Do(ctx)
	if err != nil {
		return
	}

	var t models.ClassMemberElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		classMember := item.(models.ClassMemberElastic)
		students = append(students, classMember.Tags...)
	}
	students = removeDuplicates(students)
	return
}

func removeDuplicates(arr []string) []string {
	uniqueMap := make(map[string]bool)
	for _, str := range arr {
		uniqueMap[str] = true
	}

	uniqueArr := make([]string, 0, len(uniqueMap))
	for str := range uniqueMap {
		uniqueArr = append(uniqueArr, str)
	}

	return uniqueArr
}

func SoftDeleteClassMemberBySmartbtwIDAndClassroomID(bd *request.CreateClassMember) error {
	stdCol := db.Mongodb.Collection("class_members")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"smartbtw_id": bd.SmartbtwID, "classroom_id": bd.ClassroomID}
	update := bson.M{"$set": bson.M{
		"deleted_at": time.Now(),
		"updated_at": time.Now(),
	}}

	_, err := stdCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	//delete class member elastic
	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("classroom_id", bd.ClassroomID.Hex()), elastic.NewMatchQuery("smartbtw_id", bd.SmartbtwID))
	_, err = elastic.NewDeleteByQueryService(db.ElasticClient).
		Index(db.GetClassMemberIndexName()).
		Query(bq).
		Do(context.Background())
	if err != nil {
		return err
	}

	_, errs := db.ElasticClient.Flush().Index(db.GetClassMemberIndexName()).Do(context.Background())
	if errs != nil {
		return err
	}

	return nil
}

func SwitchClassMember(bd *request.SwitchClassMember) error {
	stdCol := db.Mongodb.Collection("class_members")
	stdCol1 := db.Mongodb.Collection("classrooms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"_id": bd.ClassroomID}
	stdModels := models.Classroom{}
	err := stdCol1.FindOne(ctx, filter).Decode(&stdModels)
	if err != nil {
		return err
	}

	isOl := false
	if stdModels.IsOnline {
		isOl = true
	}

	for _, v := range bd.ClassMembers {
		//soft delete old class member
		filter1 := bson.M{"smartbtw_id": v.SmartbtwID, "classroom_id": bd.ClassroomID}
		update := bson.M{"$set": bson.M{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}}
		_, err1 := stdCol.UpdateOne(ctx, filter1, update)
		if err1 != nil {
			return err1
		}

		//delete old class member elastic
		bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("classroom_id", bd.ClassroomID.Hex()), elastic.NewMatchQuery("smartbtw_id", v.SmartbtwID))
		_, err := elastic.NewDeleteByQueryService(db.ElasticClient).
			Index(db.GetClassMemberIndexName()).
			Query(bq).
			Do(context.Background())
		if err != nil {
			return err
		}
		_, errs := db.ElasticClient.Flush().Index(db.GetClassMemberIndexName()).Do(context.Background())
		if errs != nil {
			return err
		}

		//insert new class member
		pyl := models.ClassMember{
			SmartbtwID:  int32(v.BtwedutechID),
			ClassroomID: bd.ClassroomID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		}

		_, err = stdCol.InsertOne(ctx, pyl)

		if err != nil {
			return err
		}

		//insert to elastic search classmember after
		bdElastic := models.ClassMemberElastic{
			ID:          fmt.Sprintf("%d_%s", v.BtwedutechID, bd.ClassroomID.Hex()),
			SmartbtwID:  int32(v.BtwedutechID),
			ClassroomID: bd.ClassroomID.Hex(),
			BranchCode:  stdModels.BranchCode,
			Quota:       stdModels.Quota,
			QuotaFilled: stdModels.QuotaFilled,
			Description: stdModels.Description,
			Tags:        stdModels.Tags,
			Year:        stdModels.Year,
			Status:      stdModels.Status,
			Title:       stdModels.Title,
			ClassCode:   stdModels.ClassCode,
			ProductID:   stdModels.ProductID,
			IsOnline:    isOl,
			CreatedAt:   time.Now(),
		}
		ctx1 := context.Background()
		_, err = db.ElasticClient.Index().
			Index(db.GetClassMemberIndexName()).
			Id(fmt.Sprintf("%d_%s", v.BtwedutechID, bd.ClassroomID.Hex())).
			BodyJson(bdElastic).
			Do(ctx1)

		if err != nil {
			return err
		}
	}

	return nil

}

func GetStudentJoinedClassList(smID int32, year *int, isOnline bool) ([]models.ClassMemberElastic, error) {
	classMember := []models.ClassMemberElastic{}
	ctx := context.Background()

	queries := []elastic.Query{
		elastic.NewMatchQuery("smartbtw_id", smID),
	}

	if year != nil {
		queries = append(queries,
			elastic.NewMatchQuery("year", *year))
	}
	query := elastic.NewBoolQuery().Must(
		queries...,
	)

	sort := elastic.NewFieldSort("created_at").Desc()

	// Current query still use 200 as data limit to load all question
	// TODO: Try to find better solution to load all question
	res, err := db.ElasticClient.Search().
		Index(db.GetClassMemberIndexName()).
		SortBy(sort).
		Query(query).
		Size(200).
		Do(ctx)
	if err != nil {
		return []models.ClassMemberElastic{}, err
	}

	var t models.ClassMemberElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		if item.(models.ClassMemberElastic).ID == "" {
			continue
		}
		classMember = append(classMember, item.(models.ClassMemberElastic))
	}
	// 	err = errors.New("record not found")
	// 	return []models.ClassMemberElastic{}, err
	// }
	return classMember, nil
}
