package helper

import (
	"net/url"
	"strings"
)

func DecodeURIComponent(value string) string {
	res, err := url.PathUnescape(value)
	if err != nil {
		return ""
	}
	return res
}

func EncodeURIComponent(value string) string {
	encoded := url.QueryEscape(value)
	encoded = strings.ReplaceAll(encoded, "%2F", "/")
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	return encoded
}
