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
This method creates a new collection (directory) at the specified URL.

Example implementation:

- Extract the directory path from the URL.
- Use fs.mkdir() to create the directory.
- Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
- Return the response.
*/
func (h *Handler) methodMkcol(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	apiClient := client.NewAPIClient(token)
	directoryPath := helper.DecodeURIComponent(helper.Dirname(r.URL.Path))
	directory, err := apiClient.GetFileByPath(directoryPath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	if _, err = apiClient.CreateFolder(client.FileCreateFolderOptions{
		Type:        client.FileTypeFolder,
		WorkspaceID: directory.WorkspaceID,
		ParentID:    directory.ID,
		Name:        helper.DecodeURIComponent(path.Base(r.URL.Path)),
	}); err != nil {
		infra.HandleError(err, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
