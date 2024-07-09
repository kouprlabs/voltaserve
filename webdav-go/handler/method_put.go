package handler

import (
	"fmt"
	"net/http"
	"path"
	"voltaserve/client"
	"voltaserve/helper"
	"voltaserve/infra"
)

/*
This method creates or updates a resource with the provided content.

Example implementation:

- Extract the file path from the URL.
- Create a write stream to the file.
- Listen for the data event to write the incoming data to the file.
- Listen for the end event to indicate the completion of the write stream.
- Set the response status code to 201 if created or 204 if updated.
- Return the response.
*/
func (h *Handler) methodPut(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	name := helper.DecodeURIComponent(path.Base(r.URL.Path))
	if helper.IsMicrosoftOfficeLockFile(name) || helper.IsOpenOfficeOfficeLockFile(name) {
		w.WriteHeader(http.StatusOK)
		return
	}
	apiClient := client.NewAPIClient(token)
	directory, err := apiClient.GetFileByPath(helper.DecodeURIComponent(helper.Dirname(r.URL.Path)))
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	existingFile, err := apiClient.GetFileByPath(r.URL.Path)
	if err == nil {
		if _, err = apiClient.PatchFile(client.FilePatchOptions{
			ID:     existingFile.ID,
			Reader: r.Body,
			Name:   name,
		}); err != nil {
			infra.HandleError(err, w)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		if _, err = apiClient.CreateFile(client.FileCreateOptions{
			Type:        client.FileTypeFile,
			WorkspaceID: directory.WorkspaceID,
			ParentID:    directory.ID,
			Reader:      r.Body,
			Name:        name,
		}); err != nil {
			infra.HandleError(err, w)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}
