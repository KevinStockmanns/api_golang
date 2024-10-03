package utils

import "regexp"

func IsInteger(num string) bool {
	matched, _ := regexp.MatchString("^[0-9]+$", num)
	return matched
}
