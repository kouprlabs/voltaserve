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
	"net/http"
	"voltaserve/cache"
	"voltaserve/client"
	"voltaserve/infra"
)

type Handler struct {
	s3             *infra.S3Manager
	workspaceCache *cache.WorkspaceCache
}

func NewHandler() *Handler {
	return &Handler{
		s3:             infra.NewS3Manager(),
		workspaceCache: cache.NewWorkspaceCache(),
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
	apiClient := client.NewHealthAPIClient()
	apiHealth, err := apiClient.GetHealth()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	idpClient := client.NewHealthIdPClient()
	idpHealth, err := idpClient.GetHealth()
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
