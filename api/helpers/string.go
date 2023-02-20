package helpers

import (
	"regexp"
	"strings"
)

func RemoveNonAlphanumeric(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(reg.ReplaceAllString(s, ""))
}

func RemoveNonNumeric(s string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(reg.ReplaceAllString(s, ""))
}
