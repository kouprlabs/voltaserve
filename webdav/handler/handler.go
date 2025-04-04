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

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/infra"

	"github.com/kouprlabs/voltaserve/webdav/config"
)

type Handler struct {
	s3 infra.S3Manager
}

func NewHandler() *Handler {
	return &Handler{
		s3: infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
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
	apiHealth, err := client.NewHealthClient(config.GetConfig().APIURL).Get()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	idpHealth, err := client.NewHealthClient(config.GetConfig().IdPURL).Get()
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
