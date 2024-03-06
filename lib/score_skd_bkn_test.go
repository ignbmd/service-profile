package lib_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func Test_CreateScoreSuccess(t *testing.T) {
	Init()

	json := `
	{
		"version" : 2,
		"data": {
			"id": 69,
			"name": "Akun Siswa Satu",
			"email": "akunsiswasatu@email.com",
			"gender": 1,
			"birth_date_location": "TTL",
			"phone": "081444555666",
			"school_origin": "Asal Sekolah",
			"intention": "CPNS_TEST_PREPARATION",
			"last_ed": "Pendidikan Terakhir",
			"major": "Jurusan",
			"profession": null,
			"address": "Alamat",
			"province_id": 2,
			"region_id": 28,
			"parent_name": "Nama Orang Tua",
			"parent_number": "6285737573551",
			"photo": null,
			"user_tryout_id": 130347,
			"status": false,
			"is_phone_verified": false,
			"is_email_verified": false,
			"is_data_complete": false,
			"branch_code": "PT0000",
			"affiliate_code": null,
			"additional_info": "Dapat info soal btw dari media sosial",
			"created_at": "2022-01-10T07:53:16Z",
			"updated_at": "2022-01-10T07:53:16Z"
		}
	}
	`
	bdy := []byte(json)
	paylo, _ := request.UnmarshalMessageStudentBody(bdy)
	lib.UpsertStudent(&paylo.Data)

	payload := request.ScoreSkdBkn{
		SmartBtwID: 69,
		Year:       2020,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	_, err := lib.CreateScoreData(&payload)
	assert.Nil(t, err)

}

func Test_CreateScoreErrorExist(t *testing.T) {
	Init()
	payload := request.ScoreSkdBkn{
		SmartBtwID: 69,
		Year:       2019,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	_, err := lib.CreateScoreData(&payload)
	assert.Nil(t, err)

	_, err = lib.CreateScoreData(&payload)
	assert.NotNil(t, err)
}

func Test_UpdateScoreSuccess(t *testing.T) {
	Init()
	payload := request.ScoreSkdBkn{
		SmartBtwID: 69,
		Year:       2016,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	res, err := lib.CreateScoreData(&payload)
	assert.Nil(t, err)

	payload1 := request.UpdateScoreSKDBKN{
		SmartBtwID: 69,
		ScoreTWK:   66,
		ScoreTIU:   66,
		ScoreTKP:   66,
		ScoreSKD:   111,
	}

	err = lib.UpdateScoreData(&payload1, res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}

// func TestHoaHoe(t *testing.T) {
// 	Init()
// 	res, err := lib.GetAllStudentStageClass()
// 	fmt.Println(res)
// 	assert.Nil(t, err)
// }

func Test_UpdateScoreErrorNotFound(t *testing.T) {
	Init()

	payload1 := request.UpdateScoreSKDBKN{
		SmartBtwID: 70,
		ScoreTWK:   66,
		ScoreTIU:   66,
		ScoreTKP:   66,
		ScoreSKD:   111,
	}

	err := lib.UpdateScoreData(&payload1, primitive.NewObjectID())
	assert.NotNil(t, err)
}

func Test_GetScoreStudentSuccess(t *testing.T) {
	Init()
	res, err := lib.GetScoreDataByStudent(69)
	assert.Nil(t, err)
	fmt.Println(res)
}

func Test_GetSingleScoreSuccess(t *testing.T) {
	Init()

	payload := request.ScoreSkdBkn{
		SmartBtwID: 69,
		Year:       2029,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	res, err := lib.CreateScoreData(&payload)
	assert.Nil(t, err)

	res1, err := lib.GetSingleScoreData(res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)

	fmt.Println(res1)
}

func Test_GetScoreByManyStudents(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": 500,
			"name": "Akun Siswa Satu",
			"email": "akunsiswasatu@email.com",
			"gender": 1,
			"birth_date_location": "TTL",
			"phone": "081444555666",
			"school_origin": "Asal Sekolah",
			"intention": "CPNS_TEST_PREPARATION",
			"last_ed": "Pendidikan Terakhir",
			"major": "Jurusan",
			"profession": null,
			"address": "Alamat",
			"province_id": 2,
			"region_id": 28,
			"parent_name": "Nama Orang Tua",
			"parent_number": "6285737573551",
			"photo": null,
			"user_tryout_id": 130347,
			"status": false,
			"is_phone_verified": false,
			"is_email_verified": false,
			"is_data_complete": false,
			"branch_code": "PT0000",
			"affiliate_code": null,
			"additional_info": "Dapat info soal btw dari media sosial",
			"created_at": "2022-01-10T07:53:16Z",
			"updated_at": "2022-01-10T07:53:16Z"
		}
	}
	`
	bdy := []byte(json)
	paylo, _ := request.UnmarshalMessageStudentBody(bdy)
	lib.UpsertStudent(&paylo.Data)

	payload := request.ScoreSkdBkn{
		SmartBtwID: 500,
		Year:       2020,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	_, err := lib.CreateScoreData(&payload)
	assert.Nil(t, err)

	res, err := lib.GetManyStudentScoreByYear([]int{69, 500}, 2020)
	assert.Nil(t, err)

	fmt.Println(res)

}

func Test_DeleteScoreByID(t *testing.T) {
	Init()

	payload := request.ScoreSkdBkn{
		SmartBtwID: 69,
		Year:       2015,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	res, err := lib.CreateScoreData(&payload)
	assert.Nil(t, err)

	err = lib.DeleteScoreSingleRecord(res.InsertedID.(primitive.ObjectID))
	assert.Nil(t, err)
}
