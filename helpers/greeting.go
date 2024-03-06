package helpers

import (
	"time"
)

func GetGreeting() string {
	// loc, err := time.LoadLocation("Asia/Jakarta")
	// if err != nil {
	// 	fmt.Println("Error loading location:", err)
	// 	return ""
	// }

	currentTime := time.Now().Add(8 * time.Hour)
	hour := currentTime.Hour()

	var greeting string

	if hour >= 0 && hour < 12 {
		greeting = "pagi"
	} else if hour >= 12 && hour < 15 {
		greeting = "siang"
	} else if hour >= 15 && hour < 18 {
		greeting = "sore"
	} else {
		greeting = "malam"
	}

	return greeting
}
