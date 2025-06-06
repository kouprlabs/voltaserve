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
	"fmt"
	"net/http"
	"path"

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"

	"github.com/kouprlabs/voltaserve/webdav/config"
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
	token, ok := r.Context().Value("token").(*dto.Token)
	if !ok {
		handleError(fmt.Errorf("missing token"), w)
		return
	}
	cl := client.NewFileClient(token, config.GetConfig().APIURL, config.GetConfig().Security.APIKey)
	sourcePath := helper.DecodeURIComponent(r.URL.Path)
	targetPath := helper.DecodeURIComponent(helper.GetTargetPath(r))
	sourceFile, err := cl.GetByPath(sourcePath)
	if err != nil {
		handleError(err, w)
		return
	}
	targetDir := helper.DecodeURIComponent(helper.Dirname(helper.GetTargetPath(r)))
	targetFile, err := cl.GetByPath(targetDir)
	if err != nil {
		handleError(err, w)
		return
	}
	if sourceFile.Workspace.ID != targetFile.Workspace.ID {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("Source and target files are in different workspaces")); err != nil {
			return
		}
	} else {
		clone, err := cl.CopyOne(sourceFile.ID, targetFile.ID)
		if err != nil {
			handleError(err, w)
			return
		}
		if _, err = cl.PatchName(clone.ID, dto.FilePatchNameOptions{Name: path.Base(targetPath)}); err != nil {
			handleError(err, w)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
