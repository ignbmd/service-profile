package listener_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pandeptwidyaop/golog"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/listener"
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

func InitSlack() {
	golog.New()
}

func Init() {
	InitEnv()
	InitSlack()
	ConnectMongoDB()
	db.NewElastic()
	db.NewFirebase()
}

func Test_UpsertStudentWithCompleteBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": 130018,
			"name": "Akun Siswa Satu",
			"email": "akunsiswasatu@email.com",
			"gender": 1,
			"birth_date_location": "TTL",
			"phone": "081444555666",
			"school_origin": "Asal Sekolah",
			"school_origin_id" : "63688ee59f1f93681a559e5b",
			"intention": "CPNS_TEST_PREPARATION",
			"last_ed": "Pendidikan Terakhir",
			"major": "Jurusan",
			"profession": null,
			"address": "Alamat",
			"province_id": 2,
			"region_id": 28,
			"parent_name": "Nama Orang Tua",
			"parent_number": "6285737573551",
			"interest": "SEKDIN",
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
	msg := amqp.Delivery{
		RoutingKey:   "user.created",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenStudentBinding(&msg))
}

func Test_UpsertStudentWithInvalidBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": "130009",
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
			"interest": "PTS",
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
	msg := amqp.Delivery{
		RoutingKey:   "user.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.False(t, listener.ListenStudentBinding(&msg))
}

func Test_UpsertStudentWithIncompleteBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
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
	msg := amqp.Delivery{
		RoutingKey:   "user.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.False(t, listener.ListenStudentBinding(&msg))
}

func Test_UpsertStudentWithEmptyBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"user_tryout_id":9
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "user.updated",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.False(t, listener.ListenStudentBinding(&msg))
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
	msg := amqp.Delivery{
		RoutingKey:   "user.deleted",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenStudentBinding(&msg))
}

func Test_DeleteStudentWithIncompleteBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"deleted_at": "2022-01-13T07:55:59Z"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "user.deleted",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.False(t, listener.ListenStudentBinding(&msg))
}

func Test_DeleteStudentWithEmptyBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "user.deleted",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.False(t, listener.ListenStudentBinding(&msg))
}

func Test_DeleteStudentWithInvalidBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 2,
		"data": {
			"id": "1",
			"deleted_at": "2022-01-10T07:55:59Z"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "user.deleted",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.False(t, listener.ListenStudentBinding(&msg))
}

func Test_CreateStudentElasticBodyData(t *testing.T) {
	Init()
	json := `
	{
		"version" : 1,
		"data": {
			"smartbtw_id": 500,
			"birth_date": "2012-12-21T18:24:12.311346+08:00",
			"province_id": 2,
			"province": "surakarta",
			"region_id": 28,
			"region": "pekalongan",
			"last_ed_id": "8h12vhg421",
			"last_ed_name": "SMA Tawuran",
			"last_ed_type": "SMA",
			"last_ed_major": "IPA",
			"last_ed_major_id": 1,
			"last_ed_region": "ujung kulon",
			"last_ed_region_id": 1,
			"eye_color_blind": false,
			"height": 169,
			"weight": 70,
			"account_type": "btwedutech"
		}
	}
	`
	msg := amqp.Delivery{
		RoutingKey:   "user.create-elastic",
		Body:         []byte(json),
		Acknowledger: &MockAcknowledger{},
	}

	assert.True(t, listener.ListenStudentBinding(&msg))
}
