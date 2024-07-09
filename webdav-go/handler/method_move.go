package handler

import (
	"net/http"
)

func (h *Handler) methodMove(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
