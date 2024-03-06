package handlers_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func TestCreateParentDataSuccess(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": 501,
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
		"smartbtw_id": 501,
		"parent_name": "ganteng",
		"parent_number": "087862216989"
	}
	`)

	request, e := http.NewRequest(
		"POST",
		"/parent-data",
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

func TestUpdateParentDataSuccess(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": 502,
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

	paylo := request.CreateParentData{
		SmartBtwID: 502,
	}
	paylo.SetParentName("GANTENG")
	paylo.SetParentNumber("087862216989")

	lib.CreateParentData(&paylo)

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 502,
		"parent_name": "ganteng",
		"parent_number": "087862216988"
	}
	`)

	request, e := http.NewRequest(
		"PUT",
		"/parent-data",
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

func TestCreateParentDataUserNotFound(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 89217816294,
		"parent_name": "ganteng",
		"parent_number": "087862216989"
	}
	`)

	request, e := http.NewRequest(
		"POST",
		"/parent-data",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
}

func TestCreateParentDataBodyNull(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": null,
		"parent_name": null,
		"parent_number": null
	}
	`)

	request, e := http.NewRequest(
		"POST",
		"/parent-data",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func TestUpdateParentDataBodyNull(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": null,
		"parent_name": null,
		"parent_number": null
	}
	`)

	request, e := http.NewRequest(
		"PUT",
		"/parent-data",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusUnprocessableEntity, response.StatusCode)
}

func TestUpdateParentDataNotFound(t *testing.T) {
	Init()

	app := server.SetupFiber()
	body := []byte(`{
		"smartbtw_id": 98217390124,
		"parent_name": "ganteng",
		"parent_number": "087862216989"
	}
	`)

	request, e := http.NewRequest(
		"PUT",
		"/parent-data",
		bytes.NewBuffer(body),
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 404 created
	assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
}
