package helpers

import (
	"math"
	"strconv"
	"unicode"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func RoundFloats(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	number := val * ratio / ratio
	return math.Round(number*1000) / 1000
}

func ParseStringToInt(input string) (int, error) {
	// Filter out non-numeric characters
	var numericChars []rune
	for _, char := range input {
		if unicode.IsDigit(char) {
			numericChars = append(numericChars, char)
		}
	}

	// Convert the numeric characters to a string and then parse to int
	numericString := string(numericChars)
	result, err := strconv.Atoi(numericString)
	if err != nil {
		return 0, err
	}

	return result, nil
}
