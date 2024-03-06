package config

import "os"

const (
	EnvTest = "test"
	EnvDev  = "development"
	EnvProd = "production"
)

func IsTest() bool {
	return os.Getenv("ENV") == EnvTest
}

func IsDev() bool {
	return os.Getenv("ENV") == EnvDev
}

func IsProd() bool {
	return os.Getenv("ENV") == EnvProd
}
