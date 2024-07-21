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
	"fmt"
	"github.com/kouprlabs/voltaserve/webdav/client/api_client"
	"net/http"

	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
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
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	cl := api_client.NewFileClient(token)
	inputPath := helper.DecodeURIComponent(r.URL.Path)
	file, err := cl.GetByPath(inputPath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	if file.Type == api_client.FileTypeFile {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Snapshot.Original.Size))
	}
	w.WriteHeader(http.StatusOK)
}
