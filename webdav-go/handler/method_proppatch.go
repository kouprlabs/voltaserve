package handler

import (
	"net/http"
)

func (h *Handler) methodProppatch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
