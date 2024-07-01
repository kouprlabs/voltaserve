package helper

import (
	"regexp"
	"strings"
)

func RemoveNonNumeric(s string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(reg.ReplaceAllString(s, ""))
}
