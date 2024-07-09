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
	"github.com/google/uuid"
	"github.com/kouprlabs/voltaserve/webdav/client"
	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	apiClient := client.NewAPIClient(token)
	filePath := helper.DecodeURIComponent(r.URL.Path)
	file, err := apiClient.GetFileByPath(filePath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	outputPath := filepath.Join(os.TempDir(), uuid.New().String())
	err = apiClient.DownloadOriginal(file, outputPath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		rangeHeader = strings.Replace(rangeHeader, "bytes=", "", 1)
		parts := strings.Split(rangeHeader, "-")
		rangeStart, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			rangeStart = 0
		}
		rangeEnd := stat.Size() - 1
		if len(parts) > 1 && parts[1] != "" {
			rangeEnd, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				rangeEnd = stat.Size() - 1
			}
		}
		chunkSize := rangeEnd - rangeStart + 1
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", rangeStart, rangeEnd, stat.Size()))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", chunkSize))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusPartialContent)
		file, err := os.Open(outputPath)
		if err != nil {
			infra.HandleError(err, w)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				infra.HandleError(err, w)
			}
		}(file)
		if _, err := file.Seek(rangeStart, 0); err != nil {
			infra.HandleError(err, w)
			return
		}
		if _, err := io.CopyN(w, file, chunkSize); err != nil {
			return
		}
		if err := os.Remove(outputPath); err != nil {
			return
		}
	} else {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		file, err := os.Open(outputPath)
		if err != nil {
			infra.HandleError(err, w)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				infra.HandleError(err, w)
			}
		}(file)
		if _, err := io.Copy(w, file); err != nil {
			return
		}
		if err := os.Remove(outputPath); err != nil {
			return
		}
	}
}
