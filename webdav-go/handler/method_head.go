package handler

import (
	"net/http"
)

func (h *Handler) methodHead(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
