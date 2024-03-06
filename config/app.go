package config

import (
	"fmt"
	"os"
)

var ServiceVersion string = fmt.Sprintf("1.0-%s", os.Getenv("APP_SERVICE_VERSION"))
