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
This method copies a resource from a source URL to a destination URL.

Example implementation:

- Extract the source and destination paths from the headers or request body.
- Use fs.copyFile() to copy the file from the source to the destination.
- Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
- Return the response.
*/
func (h *Handler) methodCopy(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	apiClient := client.NewAPIClient(token)
	sourcePath := helper.DecodeURIComponent(r.URL.Path)
	targetPath := helper.DecodeURIComponent(helper.GetTargetPath(r))
	sourceFile, err := apiClient.GetFileByPath(sourcePath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	targetDir := helper.DecodeURIComponent(helper.Dirname(helper.GetTargetPath(r)))
	targetFile, err := apiClient.GetFileByPath(targetDir)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	if sourceFile.WorkspaceID != targetFile.WorkspaceID {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("Source and target files are in different workspaces")); err != nil {
			return
		}
	} else {
		clones, err := apiClient.CopyFile(targetFile.ID, client.FileCopyOptions{
			IDs: []string{sourceFile.ID},
		})
		if err != nil {
			infra.HandleError(err, w)
			return
		}
		_, err = apiClient.PatchFileName(clones[0].ID, client.FileRenameOptions{
			Name: path.Base(targetPath),
		})
		if err != nil {
			infra.HandleError(err, w)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
