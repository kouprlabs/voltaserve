package helper

import "github.com/gabriel-vasile/mimetype"

func DetectMimeFromFile(path string) string {
	mime, err := mimetype.DetectFile(path)
	if err != nil {
		return "application/octet-stream"
	}
	return mime.String()
}
