package helper

import (
	"net/http"
	"net/url"
	"path"
	"strings"
)

func GetTargetPath(req *http.Request) string {
	destination := req.Header.Get("Destination")
	if destination == "" {
		return ""
	}
	/* Check if the destination header is a full URL */
	if strings.HasPrefix(destination, "http://") || strings.HasPrefix(destination, "https://") {
		parsedURL, err := url.Parse(destination)
		if err != nil {
			return ""
		}
		return parsedURL.Path
	}
	/* Extract the path from the destination header */
	startIndex := strings.Index(destination, req.Host) + len(req.Host)
	if startIndex < len(req.Host) {
		return ""
	}
	return destination[startIndex:]
}

func Dirname(value string) string {
	trimmedValue := strings.TrimSuffix(value, "/")
	return path.Dir(trimmedValue)
}
