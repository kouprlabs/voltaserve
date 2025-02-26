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
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/api/client/apiclient"
	apimodel "github.com/kouprlabs/voltaserve/api/model"

	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
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
	cl := apiclient.NewFileClient(token)
	directory, err := cl.GetByPath(helper.DecodeURIComponent(helper.Dirname(r.URL.Path)))
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	outputPath := filepath.Join(os.TempDir(), uuid.New().String())
	//nolint:gosec // Known safe path
	file, err := os.Create(outputPath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	defer func(path string, file *os.File) {
		if err := file.Close(); err != nil {
			infra.HandleError(err, w)
		}
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(outputPath, file)
	if _, err = io.Copy(file, r.Body); err != nil {
		infra.HandleError(err, w)
		return
	}
	workspace, err := h.workspaceCache.Get(helper.ExtractWorkspaceIDFromPath(r.URL.Path))
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	snapshotID := helper.NewID()
	key := snapshotID + "/original" + strings.ToLower(filepath.Ext(name))
	if err = h.s3.PutFile(key, outputPath, infra.DetectMIMEFromPath(outputPath), workspace.GetBucket(), minio.PutObjectOptions{}); err != nil {
		infra.HandleError(err, w)
		return
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	s3Reference := apiclient.S3Reference{
		Bucket:      workspace.GetBucket(),
		Key:         key,
		SnapshotID:  snapshotID,
		Size:        stat.Size(),
		ContentType: infra.DetectMIMEFromPath(outputPath),
	}
	existingFile, err := cl.GetByPath(r.URL.Path)
	if err == nil {
		if _, err = cl.PatchFromS3(apiclient.FilePatchFromS3Options{
			ID:          existingFile.ID,
			Name:        name,
			S3Reference: s3Reference,
		}); err != nil {
			infra.HandleError(err, w)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	} else {
		if _, err = cl.CreateFromS3(apiclient.FileCreateFromS3Options{
			Type:        apimodel.FileTypeFile,
			WorkspaceID: directory.WorkspaceID,
			ParentID:    directory.ID,
			Name:        name,
			S3Reference: s3Reference,
		}); err != nil {
			infra.HandleError(err, w)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}
