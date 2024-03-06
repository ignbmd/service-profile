package scripts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
)

func SyncStudentProfileBranchCodes(offset int) error {

	fmt.Println("Fetching all branches")
	brn, err := lib.GetBranches()

	if err != nil {
		return err
	}

	if brn == nil {

		return errors.New("branches should not empty")
	}
	branch := map[string]string{}

	for _, k := range *brn {
		branch[k.BranchCode] = k.BranchName

	}
	fmt.Println("Fetched total ", len(branch), " of data")
	fmt.Println("Fetching all students")
	scCol := db.Mongodb.Collection("students")
	ctx := context.Background()

	fil := bson.M{"deleted_at": nil, "smartbtw_id": bson.M{"$gt": offset}}
	var scrModel = make([]models.Student, 0)

	sort := bson.M{"smartbtw_id": 1}
	opts := options.Find()
	opts.SetSort(sort)

	cur, err := scCol.Find(ctx, fil, opts)

	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.Student
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}
	st := time.Now()
	fmt.Println("Fetched total ", len(scrModel), " of data")

	for _, value := range scrModel {
		bc := "PT0000"
		bn := branch[bc]
		if value.BranchCode != nil {
			bc = *value.BranchCode
			bn = branch[bc]
		}
		_, err2 := db.ElasticClient.Update().
			Index(db.GetStudentProfileIndexName()).
			Id(fmt.Sprintf("%d", value.SmartbtwID)).
			Doc(map[string]interface{}{
				"branch_code": bc,
				"branch_name": bn,
			}).
			DocAsUpsert(true).
			Do(context.Background())
			// 		bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", value.SmartbtwID))
			// 		script1 := elastic.NewScript(`
			// ctx._source.branch_code = params.bc;
			// ctx._source.branch_name = params.bn;
			// `).Params(map[string]interface{}{
			// 			"bc": bc,
			// 			"bn": bn,
			// 		})
			// 		_, err2 := elastic.NewUpdateByQueryService(db.ElasticClient).
			// 			Index(db.GetStudentProfileIndexName()).
			// 			Query(bq).
			// 			Script(script1).
			// 			DoAsync(ctx)
		if err2 != nil {
			fmt.Printf("Student ID %d branch data failed to update student profile in elastic with error %s\n", value.SmartbtwID, err2.Error())
			continue
		}

		fmt.Printf("Student ID %d branch data updated\n", value.SmartbtwID)

	}
	elapsed := time.Since(st)
	et := int64(elapsed / time.Second)
	fmt.Println(len(scrModel), " of data updated, took ", et, " seconds")
	return nil
}
