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
	"encoding/json"
	"net/http"

	apicache "github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/client/apiclient"
	"github.com/kouprlabs/voltaserve/api/client/idpclient"
	apiinfra "github.com/kouprlabs/voltaserve/api/infra"
)

type Handler struct {
	s3             apiinfra.S3Manager
	workspaceCache *apicache.WorkspaceCache
}

func NewHandler() *Handler {
	return &Handler{
		s3:             apiinfra.NewS3Manager(),
		workspaceCache: apicache.NewWorkspaceCache(),
	}
}

func (h *Handler) Dispatch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		h.methodOptions(w, r)
	case "GET":
		h.methodGet(w, r)
	case "HEAD":
		h.methodHead(w, r)
	case "PUT":
		h.methodPut(w, r)
	case "DELETE":
		h.methodDelete(w, r)
	case "MKCOL":
		h.methodMkcol(w, r)
	case "COPY":
		h.methodCopy(w, r)
	case "MOVE":
		h.methodMove(w, r)
	case "PROPFIND":
		h.methodPropfind(w, r)
	case "PROPPATCH":
		h.methodProppatch(w, r)
	default:
		http.Error(w, "Method not implemented", http.StatusNotImplemented)
	}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	apiClient := apiclient.NewHealthClient()
	apiHealth, err := apiClient.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	idpClient := idpclient.NewHealthClient()
	idpHealth, err := idpClient.Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if apiHealth == "OK" && idpHealth == "OK" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func (h *Handler) Version(w http.ResponseWriter, _ *http.Request) {
	versionInfo := map[string]string{
		"version": "3.0.0",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(versionInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
