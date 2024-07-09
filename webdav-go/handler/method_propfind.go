package handler

import (
	"net/http"
)

func (h *Handler) methodPropfind(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
