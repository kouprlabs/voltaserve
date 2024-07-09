package handler

import (
	"net/http"
)

func (h *Handler) methodCopy(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
