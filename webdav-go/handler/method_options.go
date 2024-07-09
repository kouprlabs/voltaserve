package handler

import (
	"net/http"
)

/*
This method should respond with the allowed methods and capabilities of the server.

Example implementation:

- Set the response status code to 200.
- Set the Allow header to specify the supported methods, such as OPTIONS, GET, PUT, DELETE, etc.
- Return the response.
*/
func (h *Handler) methodOptions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Allow", "OPTIONS, GET, HEAD, PUT, DELETE, MKCOL, COPY, MOVE, PROPFIND, PROPPATCH")
	w.WriteHeader(http.StatusOK)
}
