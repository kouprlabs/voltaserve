package handler

import (
	"net/http"
)

func (h *Handler) methodMkcol(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
