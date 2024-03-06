package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/request"
)

func TestCreateParentDataSuccess(t *testing.T) {
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

	payload := request.CreateParentData{
		SmartBtwID: 69,
	}
	payload.SetParentName("GANTENG")
	payload.SetParentNumber("087862216989")

	err := lib.CreateParentData(&payload)

	assert.Equal(t, nil, err)
}

func TestCreateParentDataUserNotFoundError(t *testing.T) {
	Init()

	payload := request.CreateParentData{
		SmartBtwID: 696969696,
	}
	payload.SetParentName("GANTENG")
	payload.SetParentNumber("087862216989")

	err := lib.CreateParentData(&payload)
	assert.NotNil(t, err)
}

func TestCreateParentDataIgnoreExisting(t *testing.T) {
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

	paylod := request.CreateParentData{
		SmartBtwID: 69,
	}
	paylod.SetParentName("GANTENG")
	paylod.SetParentNumber("087862216989")

	lib.CreateParentData(&paylod)

	payload := request.CreateParentData{
		SmartBtwID: 69,
	}
	payload.SetParentName("GANTENG")
	payload.SetParentNumber("087862216989")

	err := lib.CreateParentData(&payload)

	assert.Nil(t, err)
}
