package scripts

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func UpdateIPDNHistoryScore() error {

	fmt.Println("Fetching all students with IPDN target scores")
	scCol := db.Mongodb.Collection("student_targets")
	ctx := context.Background()

	fil := bson.M{"school_name": "IPDN", "deleted_at": nil}
	var scrModel = make([]models.StudentTarget, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)

	cur, err := scCol.Find(ctx, fil, opts)

	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.StudentTarget
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}
	st := time.Now()
	fmt.Println("Fetched total ", len(scrModel), " of data")

	for _, value := range scrModel {

		filter := bson.M{"_id": value.ID}
		update := bson.M{"$set": bson.M{"target_score": 446, "updated_at": time.Now()}}
		_, err1 := scCol.UpdateOne(ctx, filter, update)
		if err1 != nil {
			fmt.Printf("Student ID %d IPDN data failed to update with error %s\n", value.SmartbtwID, err1.Error())
			continue
		}
		isBinsus := false

		joinedClass, _ := lib.GetStudentJoinedClassType(value.SmartbtwID)

		for _, k := range joinedClass {
			if strings.Contains(strings.ToLower(k), "binsus") {
				isBinsus = true
				break
			}
		}
		//! Preparing data for update to elastic start from here
		// Get how much user doing module
		averages, err := lib.GetStudentHistoryPTKElastic(value.SmartbtwID, false)
		if err != nil {
			fmt.Printf("Student ID %d IPDN data failed to fetch history ptk from elastic with error %s\n", value.SmartbtwID, err.Error())
			continue
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

		percATT := helpers.RoundFloat((pAtt/446)*100, 2)

		if percATT > 99 {
			percATT = 99
		}
		// fmt.Println(stUpdate)

		bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", value.SmartbtwID), elastic.NewMatchQuery("school_id", value.SchoolID))
		script1 := elastic.NewScript(`
ctx._source.passing_recommendation_avg_score = params.ps;
ctx._source.passing_recommendation_avg_percent_score = params.aps;
ctx._source.target_score = 446;
`).Params(map[string]interface{}{
			"ps":  pAtt,
			"aps": percATT,
		})
		_, err2 := elastic.NewUpdateByQueryService(db.ElasticClient).
			Index(db.GetStudentTargetPtkIndexName()).
			Query(bq).
			Script(script1).
			DoAsync(ctx)
		if err2 != nil {
			fmt.Printf("Student ID %d IPDN data failed to update student ptk profile in elastic with error %s\n", value.SmartbtwID, err2.Error())
			continue
		}

		fmt.Printf("Student ID %d IPDN data updated\n", value.SmartbtwID)

	}
	elapsed := time.Since(st)
	et := int64(elapsed / time.Second)
	fmt.Println(len(scrModel), " of data updated, took ", et, " seconds")
	return nil
}

func UpdateSTDHistoryScore() error {

	fmt.Println("Fetching all students with STD target scores")
	scCol := db.Mongodb.Collection("student_targets")
	ctx := context.Background()

	fil := bson.M{"polbit_location_id": 17, "major_id": 214}
	var scrModel = make([]models.StudentTarget, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)

	cur, err := scCol.Find(ctx, fil, opts)

	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.StudentTarget
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}
	st := time.Now()
	fmt.Println("Fetched total ", len(scrModel), " of data")

	for _, value := range scrModel {

		filter := bson.M{"_id": value.ID}
		update := bson.M{"$set": bson.M{"target_score": 442, "updated_at": time.Now()}}
		_, err1 := scCol.UpdateOne(ctx, filter, update)
		if err1 != nil {
			fmt.Printf("Student ID %d STD data failed to update with error %s\n", value.SmartbtwID, err1.Error())
			continue
		}
		isBinsus := false

		joinedClass, _ := lib.GetStudentJoinedClassType(value.SmartbtwID)

		for _, k := range joinedClass {
			if strings.Contains(strings.ToLower(k), "binsus") {
				isBinsus = true
				break
			}
		}
		//! Preparing data for update to elastic start from here
		// Get how much user doing module
		averages, err := lib.GetStudentHistoryPTKElastic(value.SmartbtwID, false)
		if err != nil {
			fmt.Printf("Student ID %d STD data failed to fetch history ptk from elastic with error %s\n", value.SmartbtwID, err.Error())
			continue
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

		percATT := helpers.RoundFloat((pAtt/442)*100, 2)

		if percATT > 99 {
			percATT = 99
		}
		// fmt.Println(stUpdate)

		bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", value.SmartbtwID), elastic.NewMatchQuery("school_id", value.SchoolID), elastic.NewMatchQuery("major_id", 214), elastic.NewMatchQuery("polbit_location_id", 17))
		script1 := elastic.NewScript(`
ctx._source.passing_recommendation_avg_score = params.ps;
ctx._source.passing_recommendation_avg_percent_score = params.aps;
ctx._source.target_score = 442;
`).Params(map[string]interface{}{
			"ps":  pAtt,
			"aps": percATT,
		})
		_, err2 := elastic.NewUpdateByQueryService(db.ElasticClient).
			Index(db.GetStudentTargetPtkIndexName()).
			Query(bq).
			Script(script1).
			DoAsync(ctx)
		if err2 != nil {
			fmt.Printf("Student ID %d STD data failed to update student ptk profile in elastic with error %s\n", value.SmartbtwID, err2.Error())
			continue
		}

		fmt.Printf("Student ID %d STD data updated\n", value.SmartbtwID)

	}
	elapsed := time.Since(st)
	et := int64(elapsed / time.Second)
	fmt.Println(len(scrModel), " of data updated, took ", et, " seconds")
	return nil
}

func ResyncPolbitHistoryScore(polbitType string) error {

	fmt.Println("Fetching all students with polbit target scores")
	scCol := db.Mongodb.Collection("student_targets")
	ctx := context.Background()

	fil := bson.M{"polbit_type": fmt.Sprintf("DAERAH_%s", strings.ToUpper(polbitType))}
	var scrModel = make([]models.StudentTarget, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)

	cur, err := scCol.Find(ctx, fil, opts)

	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model models.StudentTarget
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}
	st := time.Now()
	fmt.Println("Fetched total ", len(scrModel), " of data")

	for _, value := range scrModel {
		skdRankData, err := lib.GetSKDRankFromCompMap((value.MajorID), (value.SchoolID), value.PolbitLocationID)
		if err != nil {
			fmt.Printf("Student ID %d  (majorId %d schoolId %d) data failed to fetch skd RANK with error %s\n", value.SmartbtwID, value.MajorID, value.SchoolID, err.Error())
			continue
		}
		filter := bson.M{"_id": value.ID}
		update := bson.M{"$set": bson.M{"target_score": skdRankData.StudyProgramPassingGrade, "updated_at": time.Now()}}
		_, err1 := scCol.UpdateOne(ctx, filter, update)
		if err1 != nil {
			fmt.Printf("Student ID %d STD data failed to update with error %s\n", value.SmartbtwID, err1.Error())
			continue
		}

		if value.Position == 0 && value.Type == "PRIMARY" && value.IsActive {
			isBinsus := false

			joinedClass, _ := lib.GetStudentJoinedClassType(value.SmartbtwID)

			for _, k := range joinedClass {
				if strings.Contains(strings.ToLower(k), "binsus") {
					isBinsus = true
					break
				}
			}
			//! Preparing data for update to elastic start from here
			// Get how much user doing module
			averages, err := lib.GetStudentHistoryPTKElastic(value.SmartbtwID, false)
			if err != nil {
				fmt.Printf("Student ID %d STD data failed to fetch history ptk from elastic with error %s\n", value.SmartbtwID, err.Error())
				continue
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

			percATT := helpers.RoundFloat((pAtt/float64(skdRankData.StudyProgramPassingGrade))*100, 2)

			if percATT > 99 {
				percATT = 99
			}
			// fmt.Println(stUpdate)

			//	bq := elastic.NewBoolQuery().Must(elastic.NewMatchQuery("smartbtw_id", value.SmartbtwID), elastic.NewMatchQuery("school_id", value.SchoolID), elastic.NewMatchQuery("major_id", value.MajorID), elastic.NewMatchQuery("polbit_location_id", *value.PolbitLocationID))
			//	script1 := elastic.NewScript(`
			//
			// ctx._source.passing_recommendation_avg_score = params.ps;
			// ctx._source.passing_recommendation_avg_percent_score = params.aps;
			// ctx._source.target_score = params.tss;
			//
			//	`).Params(map[string]interface{}{
			//					"ps":  pAtt,
			//					"aps": percATT,
			//					"tss": skdRankData.StudyProgramPassingGrade,
			//				})
			//				_, err2 := elastic.NewUpdateByQueryService(db.ElasticClient).
			//					Index(db.GetStudentTargetPtkIndexName()).
			//					Query(bq).
			//					Script(script1).
			//					DoAsync(ctx)
			_, err1 := db.ElasticClient.Update().
				Index(db.GetStudentTargetPtkIndexName()).
				Id(fmt.Sprintf("%d_PTK", value.SmartbtwID)).
				Doc(map[string]interface{}{
					"school_id":                        value.SchoolID,
					"major_id":                         value.MajorID,
					"school_name":                      value.SchoolName,
					"major_name":                       value.MajorName,
					"polbit_type":                      value.PolbitType,
					"polbit_competition_id":            value.PolbitCompetitionID,
					"polbit_location_id":               value.PolbitLocationID,
					"passing_recommendation_avg_score": pAtt,
					"passing_recommendation_avg_percent_score": percATT,
					"target_score": skdRankData.StudyProgramPassingGrade,
				}).
				DocAsUpsert(true).
				Do(context.Background())
			if err1 != nil {
				return err1
			}

			_, err2 := db.ElasticClient.Update().
				Index(db.GetStudentProfileIndexName()).
				Id(fmt.Sprintf("%d", value.SmartbtwID)).
				Doc(map[string]interface{}{
					"school_ptk_id":             value.SchoolID,
					"major_ptk_id":              value.MajorID,
					"school_name_ptk":           value.SchoolName,
					"major_name_ptk":            value.MajorName,
					"target_score_ptk":          skdRankData.StudyProgramPassingGrade,
					"polbit_type_ptk":           value.PolbitType,
					"polbit_competition_ptk_id": value.PolbitCompetitionID,
					"polbit_location_ptk_id":    value.PolbitLocationID,
					"created_at_ptk":            time.Now(),
				}).
				DocAsUpsert(true).
				Do(context.Background())
			if err2 != nil {
				fmt.Printf("Student ID %d STD data failed to update student ptk profile in elastic with error %s\n", value.SmartbtwID, err2.Error())
				continue
			}
			fmt.Printf("Student ID %d STD elastic data updated\n", value.SmartbtwID)

		}
		fmt.Printf("Student ID %d STD data updated\n", value.SmartbtwID)

	}
	elapsed := time.Since(st)
	et := int64(elapsed / time.Second)
	fmt.Println(len(scrModel), " of data updated, took ", et, " seconds")
	return nil
}

func ResyncCentralTargetScore() error {

	fmt.Println("Fetching all students with central target scores")
	scCol := db.Mongodb.Collection("student_targets")
	ctx := context.Background()
	smId := []uint{
		198384,
		191429,
		198387,
		198352,
		198281,
		198383,
		198467,
		194031,
		198579,
		198446,
		198498,
		190977,
		198967,
		198367,
		198420,
		198332,
		198276,
		178404,
		198371,
		198458,
		198322,
		198867,
		198268,
		198575,
		198335,
		198406,
		198528,
		198306,
		198417,
		194422,
		198318,
		198340,
		173666,
		178360,
		198871,
		175636,
		175854,
		175755,
		175552,
		198427,
		175741,
		198451,
		198418,
		198424,
		198954,
		198516,
		198393,
		198714,
		198447,
		198496,
		198328,
		198180,
		198412,
		198324,
		198182,
		198533,
		197638,
		198859,
		198160,
		198266,
		198993,
		198369,
		198171,
		179364,
		198597,
		198145,
		198126,
		197659,
		195653,
		197629,
		198124,
		198150,
		198958,
		197647,
		198343,
		197653,
		198128,
		197627,
		197651,
		198118,
		198462,
		198290,
		197625,
		197665,
		197690,
		184172,
		175828,
		198518,
		175526,
		198339,
		175584,
		175517,
		199068,
		198556,
		175578,
		175619,
		198585,
		175889,
		175857,
		175628,
		175782,
		198203,
		175613,
		175608,
		198943,
		175885,
		175519,
		199092,
		198192,
		175561,
		175810,
		198260,
		175781,
		175703,
		175604,
		175769,
		198883,
		198503,
		198296,
		175582,
		198284,
		198287,
		198500,
		198262,
		198521,
		198445,
		198316,
		198444,
		198510,
		198559,
		198440,
		198278,
		198285,
		199210,
		198710,
		198244,
		180754,
		180712,
		198156,
		180377,
		190527,
		198161,
		198142,
		193047,
		198143,
		198141,
		198139,
		190495,
		175869,
		175816,
		198966,
		198164,
		175617,
		198992,
		198975,
		178068,
		198955,
		198950,
		175838,
		175364,
		199030,
		198228,
		198903,
		198934,
		194845,
		198929,
		198301,
		198480,
		198298,
		198501,
		198310,
		198270,
		198312,
		198321,
		198345,
		198376,
		198542,
		198282,
		198319,
		198292,
		198633,
		198479,
		198302,
		198377,
		198277,
		198364,
		199892,
		198489,
		198506,
		198381,
		198372,
		198342,
		175401,
		175653,
		198486,
		198481,
		198474,
		198478,
		175739,
		198465,
		198951,
		198436,
		198487,
		179407,
		198898,
		181849,
		179397,
		178298,
		175880,
		175844,
		198971,
		198880,
		198241,
		198905,
		198250,
		198906,
		198816,
		198223,
		198215,
		198368,
		198240,
		198457,
		198197,
		198295,
		198263,
		198415,
		198454,
		189728,
		198235,
		198237,
		198207,
		198488,
		198434,
		175507,
		198361,
		198378,
		198294,
		175947,
		198362,
		198502,
		175836,
		198274,
		198429,
		198190,
		175500,
		175365,
		197704,
		197713,
		197698,
		197733,
		198910,
		175474,
		199527,
		174199,
		174200,
		175397,
		174221,
		198267,
		197721,
		174526,
		197739,
		198902,
		198495,
		179284,
		198112,
		197735,
		197750,
		197661,
		197738,
		197640,
		197664,
		197631,
		198574,
		198131,
		179383,
		198116,
		197632,
		197730,
		197682,
		198115,
		197646,
		197726,
		198106,
		197734,
		197688,
		197776,
		197786,
		195556,
		197660,
		198108,
		198114,
		197703,
		197636,
		172898,
		175775,
		175771,
		172919,
		175774,
		198933,
		198931,
		198561,
		176286,
		198870,
		198546,
		198578,
		198863,
		198610,
		175717,
		179750,
		198862,
		198549,
		175738,
		198543,
		198661,
		198403,
		175693,
		198565,
		198395,
		198408,
		198402,
		198398,
		198571,
		198400,
		198624,
		198895,
		198558,
		198630,
		198536,
		198961,
		198915,
		198523,
		199190,
		199849,
		201077,
		201415,
		199115,
		198728,
		202344,
		198470,
		201241,
		203992,
		202465,
		204934,
		206003,
		206597,
		207028,
		199032,
		182165,
		182063,
		182148,
		182011,
		182149,
		205906,
		182086,
		182022,
		182100,
		181989,
		182114,
		182095,
		182146,
		206048,
		206167,
		205902,
		182176,
		182079,
		182002,
		182142,
		205833,
		193875,
		205249,
		205831,
		207569,
		178018,
		184257,
		186251,
		184006,
		205093,
		205260,
		205000,
		210813,
		203752,
		205313,
		205329,
		205901,
		205635,
		205491,
		205224,
		205090,
		206071,
		206087,
		205327,
		203877,
		205373,
		206063,
		185593,
		206070,
		206100,
		206080,
		206177,
		203871,
		206187,
		206314,
		206338,
		206340,
		205472,
		203860,
		204206,
		206993,
		206795,
		207695,
		207696,
		207394,
		184021,
		207707,
		207685,
		207840,
		207206,
		207382,
		207598,
		207723,
		207596,
		207642,
		207603,
		207766,
		207719,
		207869,
		208058,
		207851,
		203692,
		201462,
		204699,
		205109,
		198469,
		199053,
		207999,
		208523,
		191367,
		190579,
		205184,
		204994,
		204996,
		204995,
		206010,
		205870,
		204959,
		204909,
		204960,
		198923,
		207998,
		199850,
		204977,
		204345,
		205038,
		205204,
		207007,
		198989,
		210260,
		210354,
		210459,
		210474,
		210568,
		211113,
		211098,
		210787,
		206561,
		208576,
		211119,
		189913,
		207640,
		207641,
		207638,
		207882,
		207844,
		207724,
		208453,
		199603,
		199569,
		199556,
		206929,
		199574,
		206930,
		207385,
		205112,
		205115,
		205113,
		205111,
		205108,
		205117,
		205128,
		205110,
		205043,
		206497,
		206148,
		206648,
		207720,
		207665,
		207721,
		207393,
		208564,
		211307,
		201088,
		199137,
		206319,
		199875,
		198995,
		172808,
		176531,
		176516,
		202345,
		172812,
		172813,
		211049,
		178075,
	}
	for _, ls := range smId {
		fil := bson.M{"smartbtw_id": ls, "target_type": "PTK", "is_active": true, "position": 0, "type": "PRIMARY"}
		// fil := bson.M{"smartbtw_id": bson.M{"$gt": 172000}, "target_type": "PTK", "is_active": true, "position": 0, "type": "PRIMARY"}

		var scrModel = make([]models.StudentTarget, 0)

		sort := bson.M{"smartbtw_id": 1}
		opts := options.Find()
		opts.SetSort(sort)

		cur, err := scCol.Find(ctx, fil, opts)

		if err != nil {
			return err
		}

		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var model models.StudentTarget
			e := cur.Decode(&model)
			if e != nil {
				log.Fatal(e)
			}
			scrModel = append(scrModel, model)
		}
		st := time.Now()
		fmt.Println("Fetched total ", len(scrModel), " of data")
		count := 0
		for idx, value := range scrModel {

			// if value.SchoolID != 23 {
			// 	continue
			// }

			// gt, err := lib.GetStudentProfileElastic(value.SmartbtwID)

			// if err != nil {
			// 	continue
			// }

			// mjrId := value.MajorID
			// if value.SchoolID == 4 {
			// 	mjrId = 347
			// }
			// if value.SchoolID == 6 {
			// 	mjrId = 349
			// }
			// if value.SchoolID == 5 {
			// 	mjrId = 348
			// }
			// if value.SchoolID == 8 {
			// 	mjrId = 350
			// }

			// compMapData, err := lib.GetCompetitionMapData((mjrId), (value.SchoolID), value.PolbitLocationID, gt.Gender)
			// // skdRankData, err := lib.GetSKDRankFromCompMap((value.MajorID), (value.SchoolID), value.PolbitLocationID)
			// if err != nil {
			// 	// fmt.Printf("Student ID %d  (majorId %d schoolId %d) data failed to fetch skd RANK with error %s\n", value.SmartbtwID, value.MajorID, value.SchoolID, err.Error())
			// 	continue
			// }

			// if compMapData.LowestScore == 0 {
			// 	// fmt.Printf("Student ID %d  (majorId %d schoolId %d) data failed to fetch skd RANK no data available\n", value.SmartbtwID, value.MajorID, value.SchoolID)
			// 	continue
			// }
			// filter := bson.M{"_id": value.ID}
			// update := bson.M{"$set": bson.M{"target_score": compMapData.LowestScore, "updated_at": time.Now()}}
			// _, err1 := scCol.UpdateOne(ctx, filter, update)
			// if err1 != nil {
			// 	fmt.Printf("Student ID %d STD data failed to update with error %s\n", value.SmartbtwID, err1.Error())
			// 	continue
			// }

			// if value.Position == 0 && value.Type == "PRIMARY" && value.IsActive {
			isBinsus := false

			joinedClass, _ := lib.GetStudentJoinedClassType(value.SmartbtwID)

			for _, k := range joinedClass {
				if strings.Contains(strings.ToLower(k), "binsus") {
					isBinsus = true
					break
				}
			}

			if !isBinsus {
				fmt.Printf("Student ID %d is not binsus\n", value.SmartbtwID)
				continue
			}
			//! Preparing data for update to elastic start from here
			// Get how much user doing module
			averages, err := lib.GetStudentHistoryPTKElastic(value.SmartbtwID, false)
			if err != nil {
				fmt.Printf("Student ID %d STD data failed to fetch history ptk from elastic with error %s\n", value.SmartbtwID, err.Error())
				continue
			}

			passingTotalScore := float64(0)
			passingTotalItem := 0

			twkScore := float64(0)
			tiuScore := float64(0)
			tkpScore := float64(0)
			atwk := float64(0)
			atiu := float64(0)
			atkp := float64(0)
			pAtt := float64(0)
			percTT := float64(0)
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
				atwk = math.Round(helpers.RoundFloat((twkScore / float64(passingTotalItem)), 2))
				atiu = math.Round(helpers.RoundFloat((tiuScore / float64(passingTotalItem)), 2))
				atkp = math.Round(helpers.RoundFloat((tkpScore / float64(passingTotalItem)), 2))
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

			percTWK := helpers.RoundFloat((atwk/150)*100, 2)
			percTIU := helpers.RoundFloat((atiu/175)*100, 2)
			percTKP := helpers.RoundFloat((atkp/225)*100, 2)
			percTT = helpers.RoundFloat((pAtt/225)*100, 2)
			percATT := helpers.RoundFloat((pAtt/float64(value.TargetScore))*100, 2)

			if percATT > 99 {
				percATT = 99
			}
			_, err1 := db.ElasticClient.Update().
				Index(db.GetStudentTargetPtkIndexName()).
				Id(fmt.Sprintf("%d_PTK", value.SmartbtwID)).
				Doc(map[string]interface{}{
					"passing_recommendation_avg_score":         pAtt,
					"passing_recommendation_avg_percent_score": percATT,
					"total_avg_score":                          pAtt,
					"total_avg_percent_score":                  percTT,
					"tkp_avg_score":                            atkp,
					"tkp_avg_percent_score":                    percTKP,
					"tiu_avg_score":                            atiu,
					"tiu_avg_percent_score":                    percTIU,
					"twk_avg_score":                            atwk,
					"twk_avg_percent_score":                    percTWK,

					// "target_score": compMapData.LowestScore,
				}).
				DocAsUpsert(true).
				Do(context.Background())
			if err1 != nil {
				return err1
			}

			// _, err2 := db.ElasticClient.Update().
			// 	Index(db.GetStudentProfileIndexName()).
			// 	Id(fmt.Sprintf("%d", value.SmartbtwID)).
			// 	Doc(map[string]interface{}{
			// 		"target_score_ptk": compMapData.LowestScore,
			// 		"created_at_ptk":   time.Now(),
			// 	}).
			// 	DocAsUpsert(true).
			// 	Do(context.Background())
			// if err2 != nil {
			// 	fmt.Printf("Student ID %d STD data failed to update student ptk profile in elastic with error %s\n", value.SmartbtwID, err2.Error())
			// 	continue
			// }
			// fmt.Printf("Student ID %d %s elastic data updated\n", value.SmartbtwID, value.PolbitType)

			// }
			fmt.Printf("Student ID %d (%d - %d) %s data updated from score %d to %f(%d of %d)\n", value.SmartbtwID, value.SchoolID, value.MajorID, value.PolbitType, int(value.TargetScore), pAtt, idx+1, len(scrModel))
			count++
		}

		elapsed := time.Since(st)
		et := int64(elapsed / time.Second)
		fmt.Println(len(scrModel), " of data updated, took ", et, " seconds (Processed : ", count, " of data)")
	}
	return nil
}

func ResyncStudentHistoryAdditionalData(program string, isforce bool) error {
	type HistoryData struct {
		ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		SmartBtwID     int                `json:"smartbtw_id" bson:"smartbtw_id"`
		SchoolOriginID string             `json:"school_origin_id" bson:"school_origin_id"`
		InstanceID     int                `json:"instance_id" bson:"instance_id"`
		SchoolID       int                `json:"school_id" bson:"school_id"`
	}

	scCol := db.Mongodb.Collection(fmt.Sprintf("history_%s", program))
	stTarget := db.Mongodb.Collection("student_targets")
	stCol := db.Mongodb.Collection("students")
	if program == "cpns" {
		stTarget = db.Mongodb.Collection("student_target_cpns")
	}
	ctx := context.Background()

	fil := bson.M{"deleted_at": nil}
	var scrModel = make([]HistoryData, 0)

	sort := bson.M{"created_at": -1}
	opts := options.Find()
	opts.SetSort(sort)

	cur, err := scCol.Find(ctx, fil, opts)

	if err != nil {
		return err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var model HistoryData
		e := cur.Decode(&model)
		if e != nil {
			log.Fatal(e)
		}
		scrModel = append(scrModel, model)
	}
	fmt.Println("Fetched total ", len(scrModel), " of data")

	for idx, value := range scrModel {
		if !isforce {
			if program == "cpns" {
				if value.InstanceID > 0 {
					fmt.Println("Skipping student history ", value.SmartBtwID, " id ", value.ID.Hex(), " as data already added")
					continue
				}
			} else {
				if value.SchoolID > 0 {
					fmt.Println("Skipping student history ", value.SmartBtwID, " id ", value.ID.Hex(), " as data already added")
					continue
				}
			}
		}

		fmt.Printf("(%d/%d) Updating %s of ID: %d with history ID: %s\n", idx+1, len(scrModel), program, value.SmartBtwID, value.ID.Hex())
		var stData models.Student

		filterStData := bson.M{
			"smartbtw_id": value.SmartBtwID,
			"deleted_at":  nil,
		}

		err = stCol.FindOne(ctx, filterStData).Decode(&stData)
		if err != nil {
			fmt.Println("Skipping student history ", value.SmartBtwID, " id ", value.ID.Hex(), " as student not found")
			continue
		}

		type StudentTarget struct {
			ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
			SchoolID            int                `json:"school_id" bson:"school_id"`
			MajorID             int                `json:"major_id" bson:"major_id"`
			SchoolName          string             `json:"school_name" bson:"school_name"`
			MajorName           string             `json:"major_name" bson:"major_name"`
			PolbitType          string             `json:"polbit_type" bson:"polbit_type"`
			PolbitCompetitionID *int               `json:"polbit_competition_id" bson:"polbit_competition_id"`
			PolbitLocationID    *int               `json:"polbit_location_id" bson:"polbit_location_id"`
			Position            uint               `json:"position" bson:"position"`
			InstanceID          int                `json:"instance_id" bson:"instance_id"`
			InstanceName        string             `json:"instance_name" bson:"instance_name"`
			TargetScore         float64            `json:"target_score" bson:"target_score"`
			PositionID          int                `json:"position_id" bson:"position_id"`
			PositionName        string             `json:"position_name" bson:"position_name"`
			FormationType       string             `json:"formation_type" bson:"formation_type"`
			FormationLocation   string             `json:"formation_location" bson:"formation_location"`
			FormationCode       string             `json:"formation_code" bson:"formation_code"`
			CompetitionID       int                `json:"competition_id" bson:"competition_id"`
		}

		var results StudentTarget

		filter := bson.M{
			"smartbtw_id": value.SmartBtwID,
			"target_type": strings.ToUpper(program),
			"is_active":   true,
			"position":    0,
			"type":        "PRIMARY",
			"deleted_at":  nil,
		}

		err := stTarget.FindOne(ctx, filter).Decode(&results)
		if err != nil {
			fmt.Println("Skipping student history ", value.SmartBtwID, " id ", value.ID.Hex(), " as target data not found")
			continue
		}

		payl := bson.M{}

		scId := ""
		scNm := ""
		if stData.SchoolOriginID != nil {
			scId = *stData.SchoolOriginID
		}
		if stData.SchoolOrigin != nil {
			scNm = *stData.SchoolOrigin
		}
		if program == "ptk" {

			pcId := 0
			plId := 0

			if results.PolbitCompetitionID != nil {
				pcId = *results.PolbitCompetitionID
			}
			if results.PolbitLocationID != nil {
				plId = *results.PolbitLocationID
			}

			payl = bson.M{
				"student_name":          stData.Name,
				"school_origin_id":      scId,
				"school_origin":         scNm,
				"school_id":             results.SchoolID,
				"major_id":              results.MajorID,
				"school_name":           results.SchoolName,
				"major_name":            results.MajorName,
				"polbit_type":           results.PolbitType,
				"polbit_competition_id": pcId,
				"polbit_location_id":    plId,
				"target_score":          results.TargetScore,
				"updated_at":            time.Now(),
			}
		}
		if program == "ptn" {
			payl = bson.M{
				"student_name":     stData.Name,
				"school_origin_id": scId,
				"school_origin":    scNm,
				"school_id":        results.SchoolID,
				"major_id":         results.MajorID,
				"school_name":      results.SchoolName,
				"major_name":       results.MajorName,
				"target_score":     results.TargetScore,
				"updated_at":       time.Now(),
			}
		}

		if program == "cpns" {
			payl = bson.M{
				"student_name":       stData.Name,
				"school_origin_id":   scId,
				"school_origin":      scNm,
				"instance_id":        results.InstanceID,
				"instance_name":      results.InstanceName,
				"position_id":        results.PositionID,
				"position_name":      results.PositionName,
				"formation_type":     results.FormationType,
				"formation_location": results.FormationLocation,
				"formation_code":     results.FormationCode,
				"competition_id":     results.CompetitionID,
				"target_score":       results.TargetScore,
				"updated_at":         time.Now(),
			}
		}

		update := bson.M{"$set": payl}

		_, err1 := scCol.UpdateByID(ctx, value.ID, update)

		if err1 != nil {
			fmt.Println("Skipping student history ", value.SmartBtwID, " id ", value.ID.Hex(), " as error occured :", err1.Error())
			continue
		}

		delete(payl, "updated_at")

		_, err2 := db.ElasticClient.Update().
			Index(fmt.Sprintf("student_history_%s", program)).
			Id(value.ID.Hex()).
			Doc(payl).
			Do(context.Background())
		if err2 != nil {
			fmt.Println("Skipping student history ", value.SmartBtwID, " id ", value.ID.Hex(), " on elastic as error occured :", err2.Error())
			continue
		}
	}
	return nil
}
