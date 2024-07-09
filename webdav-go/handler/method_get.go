package handler

import (
	"net/http"
)

func (h *Handler) methodGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
