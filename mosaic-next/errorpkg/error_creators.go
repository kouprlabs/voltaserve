package errorpkg

import (
	"net/http"
)

func NewInternalServerError(err error) *ErrorResponse {
	return NewErrorResponse(
		"internal_server_error",
		http.StatusInternalServerError,
		"Internal server error.",
		MsgSomethingWentWrong,
		err,
	)
}

func NewResourceNotFoundError(err error) *ErrorResponse {
	return &ErrorResponse{
		Code:        "resource_not_found",
		Status:      http.StatusNotFound,
		Message:     "Resource not found.",
		UserMessage: "The requested resource could not be found.",
		MoreInfo:    err.Error(),
		Err:         err,
	}
}
