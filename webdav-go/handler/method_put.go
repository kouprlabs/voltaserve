package handler

import (
	"net/http"
)

func (h *Handler) methodPut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
