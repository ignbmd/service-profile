package handlers_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/pandeptwidyaop/golog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
	"smartbtw.com/services/profile/server"
)

func Init() {
	InitEnv()
	InitSlack()
	ConnectMongoDB()
	db.NewElastic()
}

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

func Test_GetAllStudents(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads := []request.RequestCreateStudent{
		{
			SmartbtwID:        130009,
			Name:              "Akun Siswa Satu",
			Email:             "akunsiswasatu@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130009,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130010,
			Name:              "Akun Siswa Dua",
			Email:             "akunsiswadua@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130010,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130011,
			Name:              "Akun Siswa Tiga",
			Email:             "akunsiswatiga@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130011,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130012,
			Name:              "Akun Siswa Empat",
			Email:             "akunsiswaempat@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130012,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130013,
			Name:              "Akun Siswa Lima",
			Email:             "akunsiswalima@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130013,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
	}

	for _, payload := range payloads {
		filter := bson.M{"smartbtw_id": payload.SmartbtwID}
		update := bson.M{"$set": payload}
		collection.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	req, _ := http.NewRequest("GET", "/students", nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	_, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}

	assert.Equal(t, nil, err)
	assert.Equalf(t, fiber.StatusOK, res.StatusCode, "Get all students data must return 200 status code")
}

func Test_GetStudentsBySmartBTWID(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads := []request.RequestCreateStudent{
		{
			SmartbtwID:        130009,
			Name:              "Akun Siswa Satu",
			Email:             "akunsiswasatu@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Interest:          "SEKDIN",
			Photo:             "photo.jpg",
			UserTryoutId:      130009,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130010,
			Name:              "Akun Siswa Dua",
			Email:             "akunsiswadua@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Interest:          "PTN",
			Photo:             "photo.jpg",
			UserTryoutId:      130010,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130011,
			Name:              "Akun Siswa Tiga",
			Email:             "akunsiswatiga@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Interest:          "CPNS",
			Photo:             "photo.jpg",
			UserTryoutId:      130011,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130012,
			Name:              "Akun Siswa Empat",
			Email:             "akunsiswaempat@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Interest:          "SEKDIN",
			Photo:             "photo.jpg",
			UserTryoutId:      130012,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130013,
			Name:              "Akun Siswa Lima",
			Email:             "akunsiswalima@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Interest:          "CPNS",
			Photo:             "photo.jpg",
			UserTryoutId:      130013,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
	}

	for _, payload := range payloads {
		filter := bson.M{"smartbtw_id": payload.SmartbtwID}
		update := bson.M{"$set": payload}
		collection.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	url := fmt.Sprintf("/students?smartbtw_id=%d&smartbtw_id=%d", 130010, 130011)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	_, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}

	assert.Equal(t, nil, err)
	assert.Equalf(t, fiber.StatusOK, res.StatusCode, "Get students by SmartBTWID must return 200 status code")
}

func Test_GetNonexistentStudentsBySmartBTWID(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads := []request.RequestCreateStudent{
		{
			SmartbtwID:        130009,
			Name:              "Akun Siswa Satu",
			Email:             "akunsiswasatu@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130009,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130010,
			Name:              "Akun Siswa Dua",
			Email:             "akunsiswadua@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130010,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130011,
			Name:              "Akun Siswa Tiga",
			Email:             "akunsiswatiga@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130011,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130012,
			Name:              "Akun Siswa Empat",
			Email:             "akunsiswaempat@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130012,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130013,
			Name:              "Akun Siswa Lima",
			Email:             "akunsiswalima@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130013,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
	}

	for _, payload := range payloads {
		filter := bson.M{"smartbtw_id": payload.SmartbtwID}
		update := bson.M{"$set": payload}
		collection.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	url := fmt.Sprintf("/students?smartbtw_id=%d", 999999)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	_, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}

	assert.Equal(t, nil, err)
	assert.Equalf(t, fiber.StatusNotFound, res.StatusCode, "Get students by SmartBTWID with nonexistent ID must return 400 status code")
}

func Test_GetStudentsBySmartBTWIDWithNonIntegerIDParameter(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads := []request.RequestCreateStudent{
		{
			SmartbtwID:        130009,
			Name:              "Akun Siswa Satu",
			Email:             "akunsiswasatu@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130009,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130010,
			Name:              "Akun Siswa Dua",
			Email:             "akunsiswadua@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130010,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130011,
			Name:              "Akun Siswa Tiga",
			Email:             "akunsiswatiga@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130011,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130012,
			Name:              "Akun Siswa Empat",
			Email:             "akunsiswaempat@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130012,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130013,
			Name:              "Akun Siswa Lima",
			Email:             "akunsiswalima@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130013,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
	}

	for _, payload := range payloads {
		filter := bson.M{"smartbtw_id": payload.SmartbtwID}
		update := bson.M{"$set": payload}
		collection.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	url := fmt.Sprintf("/students?smartbtw_id=%s", "asdf")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	_, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}

	assert.Equal(t, nil, err)
	assert.Equalf(t, fiber.StatusBadRequest, res.StatusCode, "Get students by SmartBTWID with non integer parameter value must return 400 status code")
}

func Test_GetSingleStudent(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads := []request.RequestCreateStudent{
		{
			SmartbtwID:        130009,
			Name:              "Akun Siswa Satu",
			Email:             "akunsiswasatu@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130009,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130010,
			Name:              "Akun Siswa Dua",
			Email:             "akunsiswadua@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130010,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130011,
			Name:              "Akun Siswa Tiga",
			Email:             "akunsiswatiga@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130011,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130012,
			Name:              "Akun Siswa Empat",
			Email:             "akunsiswaempat@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130012,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130013,
			Name:              "Akun Siswa Lima",
			Email:             "akunsiswalima@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130013,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
	}

	for _, payload := range payloads {
		filter := bson.M{"smartbtw_id": payload.SmartbtwID}
		update := bson.M{"$set": payload}
		collection.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	req, _ := http.NewRequest("GET", "/students/130009", nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, -1)
	_, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(res.Body)
	assert.Equal(t, nil, err)
	assert.Equalf(t, fiber.StatusOK, res.StatusCode, "Get single student data must return 200 status code")
}

func Test_GetSingleNonexistentStudent(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads := []request.RequestCreateStudent{
		{
			SmartbtwID:        130009,
			Name:              "Akun Siswa Satu",
			Email:             "akunsiswasatu@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130009,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130010,
			Name:              "Akun Siswa Dua",
			Email:             "akunsiswadua@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130010,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130011,
			Name:              "Akun Siswa Tiga",
			Email:             "akunsiswatiga@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130011,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130012,
			Name:              "Akun Siswa Empat",
			Email:             "akunsiswaempat@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130012,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130013,
			Name:              "Akun Siswa Lima",
			Email:             "akunsiswalima@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130013,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
	}

	for _, payload := range payloads {
		filter := bson.M{"smartbtw_id": payload.SmartbtwID}
		update := bson.M{"$set": payload}
		collection.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	req, _ := http.NewRequest("GET", "/students/1300450", nil)
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req, -1)
	_, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}
	assert.Equalf(t, fiber.StatusNotFound, res.StatusCode, "Get single student data with nonexistent ID must return 404 status code")
}

func Test_GetSingleStudentWithNonIntegerIDParameter(t *testing.T) {
	Init()

	opts := options.Update().SetUpsert(true)

	collection := db.Mongodb.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payloads := []request.RequestCreateStudent{
		{
			SmartbtwID:        130009,
			Name:              "Akun Siswa Satu",
			Email:             "akunsiswasatu@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130009,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130010,
			Name:              "Akun Siswa Dua",
			Email:             "akunsiswadua@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130010,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130011,
			Name:              "Akun Siswa Tiga",
			Email:             "akunsiswatiga@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130011,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130012,
			Name:              "Akun Siswa Empat",
			Email:             "akunsiswaempat@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130012,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
		{
			SmartbtwID:        130013,
			Name:              "Akun Siswa Lima",
			Email:             "akunsiswalima@email.com",
			Gender:            1,
			BirthDateLocation: "TTL",
			Phone:             "081444555666",
			SchoolOrigin:      "Asal Sekolah",
			Intention:         "CPNS_TEST_PREPARATION",
			LastEd:            "Pendidikan Terakhir",
			Major:             "Jurusan",
			Profession:        "Profesi",
			Address:           "Alamat",
			ProvinceId:        2,
			RegionId:          28,
			ParentName:        "Nama Orang Tua",
			ParentNumber:      "6285737573551",
			Photo:             "photo.jpg",
			UserTryoutId:      130013,
			Status:            true,
			IsPhoneVerified:   true,
			IsEmailVerified:   true,
			IsDataComplete:    true,
			BranchCode:        "PT0000",
			AffiliateCode:     "AFFILIATE_CODE",
			AdditionalInfo:    "Know about btw from social media",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			DeletedAt:         nil,
		},
	}

	for _, payload := range payloads {
		filter := bson.M{"smartbtw_id": payload.SmartbtwID}
		update := bson.M{"$set": payload}
		collection.UpdateOne(ctx, filter, update, opts)
	}

	app := server.SetupFiber()
	req, _ := http.NewRequest("GET", "/students/asdf", nil)
	req.Header.Set("Content-Type", "application/json")
	res, _ := app.Test(req, -1)
	_, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatal(e)
	}
	assert.Equalf(t, fiber.StatusBadRequest, res.StatusCode, "Get one data with non integer parameter value must return 400 status code")
}

func Test_GetStudentsWithBranchCodeNotFound(t *testing.T) {
	Init()

	app := server.SetupFiber()

	request, e := http.NewRequest(
		"GET",
		"/students-with-branch?branch_code=GANTENG&limit=3&skip=1&page=1",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetStudentsWithBranchCodeSuccess(t *testing.T) {
	Init()

	app := server.SetupFiber()

	request, e := http.NewRequest(
		"GET",
		"/students-with-branch?branch_code=PT0000&limit=3&skip=1&page=1&search=satu",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func Test_GetStudentsWithManyBranchCodeSuccess(t *testing.T) {
	Init()

	app := server.SetupFiber()

	request, e := http.NewRequest(
		"GET",
		"/students-with-many-branch?branch_code=PT0000,KB0001&limit=3&skip=1&page=1&search=satu",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}

func TestGetStudentCompletedModulesSuccess(t *testing.T) {
	Init()

	smId := int(time.Now().Unix())
	payload1 := request.CreateStudentTarget{
		SmartbtwID:  smId,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "UI",
		MajorName:   "KEDOKTERAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload1)
	if err1 != nil {
		assert.NotNil(t, err1)
	} else {
		assert.Nil(t, err1)
	}

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

	app := server.SetupFiber()
	request, e := http.NewRequest(
		"GET",
		fmt.Sprintf("/student-completed-modules?smartbtw_id=%d&target_type=%s", smId, "ptn"),
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

func TestGetStudentBranch(t *testing.T) {
	Init()

	app := server.SetupFiber()

	request, e := http.NewRequest(
		"GET",
		"/students-branch?emails=akunsiswasatu@email.com",
		nil,
	)
	request.Header.Add("Content-Type", "application/json")

	assert.Equal(t, nil, e)

	response, err := app.Test(request, -1)

	// We expected no error when request to server, the error expected to nil
	assert.Equal(t, nil, err)

	// We expected when post some data to api got 500 created
	assert.Equal(t, fiber.StatusOK, response.StatusCode)
}
