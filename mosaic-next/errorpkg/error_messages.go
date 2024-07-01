package errorpkg

var messages = map[string]string{
	"mosaic not found":                                   "Mosaic not found.",
	"no matching pipeline found":                         "This file type cannot be processed.",
	"language is undefined":                              "Language is undefined.",
	"unsupported file type":                              "Unsupported file type.",
	"text is empty":                                      "Text is empty.",
	"text exceeds supported limit of 1000000 characters": "Text exceeds supported limit of 1000000 characters.",
	"missing query param api_key":                        "Missing query param api_key.",
	"invalid query param api_key":                        "Invalid query param api_key.",
	"invalid content type":                               "Invalid content type.",
}

const FallbackMessage = "An error occurred while processing the file."

func GetUserFriendlyMessage(code string, fallback string) string {
	res, ok := messages[code]
	if !ok {
		return fallback
	}
	return res
}
