package handler

import (
	"net/http"
)

func (h *Handler) methodDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
