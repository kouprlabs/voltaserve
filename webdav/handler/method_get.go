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
	"strconv"
	"strings"

	"github.com/kouprlabs/voltaserve/webdav/client/api_client"
	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
)

/*
This method retrieves the content of a resource identified by the URL.

Example implementation:

- Extract the file path from the URL.
- Create a read stream from the file and pipe it to the response stream.
- Set the response status code to 200 if successful or an appropriate error code if the file is not found.
- Return the response.
*/
func (h *Handler) methodGet(w http.ResponseWriter, r *http.Request) {
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
	size := file.Snapshot.Original.Size
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		rangeHeader = strings.Replace(rangeHeader, "bytes=", "", 1)
		parts := strings.Split(rangeHeader, "-")
		rangeStart, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			rangeStart = 0
		}
		rangeEnd := size - 1
		if len(parts) > 1 && parts[1] != "" {
			rangeEnd, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				rangeEnd = size - 1
			}
		}
		chunkSize := rangeEnd - rangeStart + 1
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, size))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", chunkSize))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusPartialContent)
		if err := cl.DownloadOriginal(file, w, &rangeHeader); err != nil {
			if !strings.Contains(err.Error(), "write: broken pipe") && !strings.Contains(err.Error(), "write: connection reset by peer") {
				infra.GetLogger().Error(err)
			}
			return
		}
	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		if err := cl.DownloadOriginal(file, w, &rangeHeader); err != nil {
			if !strings.Contains(err.Error(), "write: broken pipe") && !strings.Contains(err.Error(), "write: connection reset by peer") {
				infra.GetLogger().Error(err)
			}
			return
		}
	}
}
