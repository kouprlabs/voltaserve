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

	"github.com/kouprlabs/voltaserve/api/client/apiclient"

	"github.com/kouprlabs/voltaserve/webdav/helper"
	"github.com/kouprlabs/voltaserve/webdav/infra"
)

/*
This method deletes a resource identified by the URL.

Example implementation:

- Extract the file path from the URL.
- Delete the file.
- Set the response status code to 204 if successful or an appropriate error code if the file is not found.
- Return the response.
*/
func (h *Handler) methodDelete(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	cl := apiclient.NewFileClient(token)
	file, err := cl.GetByPath(helper.DecodeURIComponent(r.URL.Path))
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	if err = cl.DeleteOne(file.ID); err != nil {
		infra.HandleError(err, w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
