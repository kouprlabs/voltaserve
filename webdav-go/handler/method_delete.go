package handler

import (
	"fmt"
	"net/http"
	"voltaserve/client"
	"voltaserve/helper"
	"voltaserve/infra"
)

/*
This method deletes a resource identified by the URL.

Example implementation:

- Extract the file path from the URL.
- Use fs.unlink() to delete the file.
- Set the response status code to 204 if successful or an appropriate error code if the file is not found.
- Return the response.
*/
func (h *Handler) methodDelete(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	apiClient := client.NewAPIClient(token)
	file, err := apiClient.GetFileByPath(helper.DecodeURIComponent(r.URL.Path))
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	if _, err = apiClient.DeleteFile(file.ID); err != nil {
		infra.HandleError(err, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
