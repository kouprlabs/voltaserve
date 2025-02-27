// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/webdav/config"
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
	token, ok := r.Context().Value("token").(*dto.Token)
	if !ok {
		handleError(fmt.Errorf("missing token"), w)
		return
	}
	cl := client.NewFileClient(token, config.GetConfig().APIURL, config.GetConfig().Security.APIKey)
	rootPath := helper.DecodeURIComponent(getRootPath(r.URL.Path))
	rootDir, err := cl.GetByPath(rootPath)
	if err != nil {
		handleError(err, w)
		return
	}
	if rootDir.Name != "/" && rootDir.WorkspaceID != "" {
		if _, err = cl.CreateFolder(client.FileCreateFolderOptions{
			Type:        model.FileTypeFolder,
			WorkspaceID: rootDir.WorkspaceID,
			ParentID:    rootDir.ID,
			Name:        helper.DecodeURIComponent(getSubPath(r.URL.Path)),
		}); err != nil {
			var errorResponse *errorpkg.ErrorResponse
			if errors.As(err, &errorResponse) {
				if errorResponse.Code == "file_with_similar_name_exists" && errorResponse.Status == http.StatusForbidden {
					// No-op
					return
				} else {
					handleError(err, w)
				}
			} else {
				handleError(err, w)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getRootPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return "/" + parts[1]
	}
	return path
}

func getSubPath(path string) string {
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}
