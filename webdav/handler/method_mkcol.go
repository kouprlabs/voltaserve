// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package handler

import (
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/kouprlabs/voltaserve/webdav/client/api_client"
	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
)

/*
This method creates a new collection (directory) at the specified URL.

Example implementation:

- Extract the directory path from the URL.
- Create the directory.
- Set the response status code to 201 if created or an appropriate error code if the directory already exists or encountered an error.
- Return the response.
*/
func (h *Handler) methodMkcol(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	cl := api_client.NewFileClient(token)
	wantedPath := helper.DecodeURIComponent(helper.Dirname(r.URL.Path))
	directory, err := cl.GetByPath(wantedPath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	if directory.Name != "/" && directory.WorkspaceID != "" {
		if _, err = cl.CreateFolder(api_client.FileCreateFolderOptions{
			Type:        api_client.FileTypeFolder,
			WorkspaceID: directory.WorkspaceID,
			ParentID:    directory.ID,
			Name:        helper.DecodeURIComponent(path.Base(r.URL.Path)),
		}); err != nil {
			var apiError *infra.APIError
			if errors.As(err, &apiError) {
				if apiError.Value.Code == "file_with_similar_name_exists" && apiError.Value.Status == http.StatusForbidden {
					// No-op
					return
				} else {
					infra.HandleError(err, w)
				}
			} else {
				infra.HandleError(err, w)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
