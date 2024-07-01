package errorpkg

import "fmt"

type ErrorResponse struct {
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
	MoreInfo    string `json:"moreInfo"`
	Err         error  `json:"-"`
}

func NewErrorResponse(code string, status int, message string, userMessage string, err error) *ErrorResponse {
	return &ErrorResponse{
		Code:        code,
		Status:      status,
		Message:     message,
		UserMessage: userMessage,
		MoreInfo:    fmt.Sprintf("https://voltaserve.com/docs/api/errors/%s", code),
		Err:         err,
	}
}

func (err ErrorResponse) Error() string {
	return fmt.Sprintf("%s %s", err.Code, err.Message)
}

func (err ErrorResponse) Unwrap() error {
	return err.Err
}
