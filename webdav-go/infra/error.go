package infra

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type IdPErrorResponse struct {
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
	MoreInfo    string `json:"moreInfo"`
}

type IdPError struct {
	Value IdPErrorResponse
}

func (e *IdPError) Error() string {
	return fmt.Sprintf("IdPError: %v", e.Value)
}

type APIErrorResponse struct {
	Code        string `json:"code"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	UserMessage string `json:"userMessage"`
	MoreInfo    string `json:"moreInfo"`
}

type APIError struct {
	Value APIErrorResponse
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIError: %v", e.Value)
}

func HandleError(err error, w http.ResponseWriter) {
	var apiErr *APIError
	var idpErr *IdPError
	switch {
	case errors.As(err, &apiErr):
		w.WriteHeader(apiErr.Value.Status)
		if _, err := w.Write([]byte(apiErr.Value.UserMessage)); err != nil {
			GetLogger().Error(err)
			return
		}
	case errors.As(err, &idpErr):
		w.WriteHeader(idpErr.Value.Status)
		if _, err := w.Write([]byte(idpErr.Value.UserMessage)); err != nil {
			GetLogger().Error(err)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Internal Server Error")); err != nil {
			return
		}
	}
	log.Println(err)
}
