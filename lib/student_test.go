package lib_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

func InitEnv() {
	e := godotenv.Load()
	if e != nil {
		log.Println(".env file not found, using global variable")
	}
}

func ConnectMongoDB() {
	connection := os.Getenv("MONGODB_CONNECTION")
	database := os.Getenv("MONGODB_DATABASE")
	db.Connect(connection, database)
}

func Init() {
	InitEnv()
	ConnectMongoDB()
	db.NewElastic()
}

func Test_UpsertStudentWithCompleteBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": 130009,
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
			"interest": "CPNS",
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
	body := []byte(json)
	payload, _ := request.UnmarshalMessageStudentBody(body)
	_, err := lib.UpsertStudent(&payload.Data)
	assert.Nil(t, err)
}

func Test_DeleteStudentWithCompleteBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": 130013,
			"deleted_at": "2022-01-13T07:55:59Z"
		}
	}
	`
	body := []byte(json)
	payload, _ := request.UnmarshalDeleteStudentBodyMessage(body)
	_, err := lib.DeleteStudent(&payload)
	assert.Nil(t, err)
}

func Test_GetStudentBySmartBTWIDSuccess(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": 130009,
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
			"interest": "CPNS",
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
	body := []byte(json)
	payload, _ := request.UnmarshalMessageStudentBody(body)
	_, err := lib.UpsertStudent(&payload.Data)
	assert.Nil(t, err)
	st, err := lib.GetStudentBySmartBTWID(payload.Data.ID)
	assert.Nil(t, err)
	assert.NotNil(t, st)
}

func Test_GetStudentCompletedModules(t *testing.T) {
	Init()
	smId := int(time.Now().Unix() + 90123801)

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  smId,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)
	payload := request.CreateHistoryPtn{
		SmartBtwID:              smId,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}

	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	st, err := lib.GetStudentCompletedModulesBySmartBTWID(smId, "ptn")
	assert.Nil(t, err)
	assert.NotNil(t, st)
	assert.Equal(t, len(st), 1)
}
