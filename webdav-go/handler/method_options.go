package handler

import (
	"net/http"
)

func (h *Handler) methodOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
