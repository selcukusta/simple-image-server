package helper

import "strconv"

//ValidateRange is using to value is valid range
func validateRange(value int, minValue int, maxValue int) bool {
	return value >= minValue && value <= maxValue
}

//ValidateRange is using to value is valid range
func ValidateRange(value string, minValue int, maxValue int) bool {
	numeric, _ := strconv.Atoi(value)
	return validateRange(numeric, minValue, maxValue)
}
