package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func Test_CreateRecordBKNSuccess(t *testing.T) {
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
	payload, _ := request.UnmarshalMessageStudentBody(bdy)
	lib.UpsertStudent(&payload.Data)

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 500,
		"year": 2018,
		"score_twk": 30,
		"score_tiu": 25,
		"score_tkp": 100,
		"score_skd": 200
	}`)

	request, e := http.NewRequest(
		"POST",
		"/score-skd",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func Test_CreateRecordBKNFaliedBody(t *testing.T) {
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
	payload, _ := request.UnmarshalMessageStudentBody(bdy)
	lib.UpsertStudent(&payload.Data)

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 500,
		"year": 2017,
		"score_twk": null,
		"score_tiu": 25,
		"score_tkp": 100,
		"score_skd": null
	}`)

	request, e := http.NewRequest(
		"POST",
		"/score-skd",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func Test_GetScoreByStudentSuccess(t *testing.T) {
	Init()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		"/score-skd/student/500",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

}

func Test_GetScoreByStudentError(t *testing.T) {
	Init()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		"/score-skd/student/asdas",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

}

func Test_GetScoreByStudentErrorNotFound(t *testing.T) {
	Init()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		"/score-skd/student/6666",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)

}

func Test_GetScoreDetailByIdSuccess(t *testing.T) {
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
	payload, _ := request.UnmarshalMessageStudentBody(bdy)
	lib.UpsertStudent(&payload.Data)

	payl := request.ScoreSkdBkn{
		SmartBtwID: 500,
		Year:       2031,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	res1, err := lib.CreateScoreData(&payl)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/score-skd/detail/%s", res1.InsertedID.(primitive.ObjectID).Hex()),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

}

func Test_GetScoreDetailByIdError(t *testing.T) {
	Init()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/score-skd/detail/%s", "awnjgnwaj217893y4"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func Test_UpdateScoreSuccess(t *testing.T) {
	Init()

	jso := `
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
	bdy := []byte(jso)
	payload, _ := request.UnmarshalMessageStudentBody(bdy)
	lib.UpsertStudent(&payload.Data)

	payl := request.ScoreSkdBkn{
		SmartBtwID: 500,
		Year:       2010,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	res1, err := lib.CreateScoreData(&payl)
	assert.Nil(t, err)

	bdy1 := request.UpdateScoreSKDBKN{
		SmartBtwID: 500,
		ScoreTWK:   40,
		ScoreTIU:   79,
		ScoreTKP:   123,
		ScoreSKD:   111,
	}
	enc, err := sonic.Marshal(bdy1)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"PUT",
		fmt.Sprintf("/score-skd/%s", res1.InsertedID.(primitive.ObjectID).Hex()),
		bytes.NewBuffer(enc),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusCreated, response.StatusCode)
}

func Test_UpdateDataScoreFail(t *testing.T) {
	Init()
	bdy1 := request.UpdateScoreSKDBKN{
		SmartBtwID: 500,
		ScoreTWK:   40,
		ScoreTIU:   79,
		ScoreTKP:   123,
		ScoreSKD:   111,
	}
	enc, err := sonic.Marshal(bdy1)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"PUT",
		"/score-skd/62833e3f7807a1d34cf497fd",
		bytes.NewBuffer(enc),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func Test_DeleteScoreByIdSuccess(t *testing.T) {
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
	payload, _ := request.UnmarshalMessageStudentBody(bdy)
	lib.UpsertStudent(&payload.Data)

	payl := request.ScoreSkdBkn{
		SmartBtwID: 500,
		Year:       2032,
		ScoreTWK:   35,
		ScoreTIU:   75,
		ScoreTKP:   120,
		ScoreSKD:   230,
	}
	res1, err := lib.CreateScoreData(&payl)
	assert.Nil(t, err)

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/score-skd/%s", res1.InsertedID.(primitive.ObjectID).Hex()),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

}

func Test_DeleteScoreByIdError(t *testing.T) {
	Init()

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"DELETE",
		fmt.Sprintf("/score-skd/%s", "awnjgnwaj217893y4"),
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
}

func Test_GetManyRecordScore(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": [500, 501, 502],
		"year": 2020
	}`)

	request, e := http.NewRequest(
		"POST",
		"/score-skd/year",
		bytes.NewBuffer(body),
	)

	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 201 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)

}
