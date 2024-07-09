package handler

import (
	"net/http"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Dispatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		h.methodOptions(w, r)
	case "GET":
		h.methodGet(w, r)
	case "HEAD":
		h.methodHead(w, r)
	case "PUT":
		h.methodPut(w, r)
	case "DELETE":
		h.methodDelete(w, r)
	case "MKCOL":
		h.methodMkcol(w, r)
	case "COPY":
		h.methodCopy(w, r)
	case "MOVE":
		h.methodMove(w, r)
	case "PROPFIND":
		h.methodPropfind(w, r)
	case "PROPPATCH":
		h.methodProppatch(w, r)
	default:
		http.Error(w, "Method not implemented", http.StatusNotImplemented)
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
