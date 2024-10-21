package helper

import (
	"encoding/base64"
	"strings"
)

func Base64ToBytes(value string) ([]byte, error) {
	var withoutPrefix string
	if strings.Contains(value, ",") {
		withoutPrefix = strings.Split(value, ",")[1]
	} else {
		withoutPrefix = value
	}
	b, err := base64.StdEncoding.DecodeString(withoutPrefix)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Base64ToMIME(value string) string {
	if !strings.HasPrefix(value, "data:image/") {
		return ""
	}
	colonIndex := strings.Index(value, ":")
	semicolonIndex := strings.Index(value, ";")
	if colonIndex == -1 || semicolonIndex == -1 {
		return ""
	}
	return value[colonIndex+1 : semicolonIndex]
}

func Base64ToExtension(value string) string {
	switch Base64ToMIME(value) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}
