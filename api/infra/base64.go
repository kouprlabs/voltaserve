package infra

import (
	"encoding/base64"
	"net/http"
	"os"
)

func ImageToBase64(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}
	var res string
	mimeType := http.DetectContentType(bytes)
	switch mimeType {
	case "image/jpeg":
		res += "data:image/jpeg;base64,"
	case "image/png":
		res += "data:image/png;base64,"
	}
	res += base64.StdEncoding.EncodeToString(bytes)
	return res, nil
}
