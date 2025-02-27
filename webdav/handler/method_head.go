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

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/webdav/config"
)

/*
This method is similar to GET but only retrieves the metadata of a resource, without returning the actual content.

Example implementation:

- Extract the file path from the URL.
- Retrieve the file metadata.
- Set the response status code to 200 if successful or an appropriate error code if the file is not found.
- Set the Content-Length header with the file size.
- Return the response.
*/
func (h *Handler) methodHead(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*dto.Token)
	if !ok {
		handleError(fmt.Errorf("missing token"), w)
		return
	}
	cl := client.NewFileClient(token, config.GetConfig().APIURL, config.GetConfig().Security.APIKey)
	inputPath := helper.DecodeURIComponent(r.URL.Path)
	file, err := cl.GetByPath(inputPath)
	if err != nil {
		handleError(err, w)
		return
	}
	if file.Type == model.FileTypeFile {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Snapshot.Original.Size))
	}
	w.WriteHeader(http.StatusOK)
}
