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
This method retrieves properties and metadata of a resource.

Example implementation:

- Extract the file path from the URL.
- Retrieve the file metadata.
- Format the response body in the desired XML format with the properties and metadata.
- Set the response status code to 207 if successful or an appropriate error code if the file is not found or encountered an error.
- Set the Content-Type header to indicate the XML format.
- Return the response.
*/
func (h *Handler) methodPropfind(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*dto.Token)
	if !ok {
		handleError(fmt.Errorf("missing token"), w)
		return
	}
	cl := client.NewFileClient(token, config.GetConfig().APIURL, config.GetConfig().Security.APIKey)
	file, err := cl.GetByPath(helper.DecodeURIComponent(r.URL.Path))
	if err != nil {
		handleError(err, w)
		return
	}
	if file.Type == model.FileTypeFile {
		responseXml := fmt.Sprintf(
			`<D:multistatus xmlns:D="DAV:">
				<D:response>
					<D:href>%s</D:href>
					<D:propstat>
						<D:prop>
							<D:resourcetype></D:resourcetype>
							<D:getcontentlength>%d</D:getcontentlength>
							<D:creationdate>%s</D:creationdate>
							<D:getlastmodified>%s</D:getlastmodified>
						</D:prop>
						<D:status>HTTP/1.1 200 OK</D:status>
					</D:propstat>
				</D:response>
			</D:multistatus>`,
			helper.EncodeURIComponent(file.Name),
			func() int64 {
				if file.Type == model.FileTypeFile && file.Snapshot != nil && file.Snapshot.Original != nil {
					return file.Snapshot.Original.Size
				} else {
					return 0
				}
			}(),
			helper.ToUTCString(&file.CreateTime),
			helper.ToUTCString(file.UpdateTime),
		)
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusMultiStatus)
		if _, err := w.Write([]byte(responseXml)); err != nil {
			handleError(err, w)
			return
		}
	} else if file.Type == model.FileTypeFolder {
		responseXml := fmt.Sprintf(
			`<D:multistatus xmlns:D="DAV:">
				<D:response>
					<D:href>%s</D:href>
					<D:propstat>
						<D:prop>
							<D:resourcetype><D:collection/></D:resourcetype>
							<D:getcontentlength>0</D:getcontentlength>
							<D:getlastmodified>%s</D:getlastmodified>
							<D:creationdate>%s</D:creationdate>
						</D:prop>
						<D:status>HTTP/1.1 200 OK</D:status>
					</D:propstat>
				</D:response>`,
			helper.EncodeURIComponent(r.URL.Path),
			helper.ToUTCString(file.UpdateTime),
			helper.ToUTCString(&file.CreateTime),
		)
		list, err := cl.ListByPath(helper.DecodeURIComponent(r.URL.Path))
		if err != nil {
			handleError(err, w)
			return
		}
		for _, item := range list {
			itemXml := fmt.Sprintf(
				`<D:response>
					<D:href>%s</D:href>
					<D:propstat>
						<D:prop>
							<D:resourcetype>%s</D:resourcetype>
							<D:getcontentlength>%d</D:getcontentlength>
							<D:getlastmodified>%s</D:getlastmodified>
							<D:creationdate>%s</D:creationdate>
						</D:prop>
						<D:status>HTTP/1.1 200 OK</D:status>
					</D:propstat>
				</D:response>`,
				helper.EncodeURIComponent(r.URL.Path+item.Name),
				func() string {
					if item.Type == model.FileTypeFolder {
						return "<D:collection/>"
					}
					return ""
				}(),
				func() int64 {
					if item.Type == model.FileTypeFile && item.Snapshot != nil && item.Snapshot.Original != nil {
						return item.Snapshot.Original.Size
					} else {
						return 0
					}
				}(),
				helper.ToUTCString(item.UpdateTime),
				helper.ToUTCString(&item.CreateTime),
			)
			responseXml += itemXml
		}
		responseXml += `</D:multistatus>`
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusMultiStatus)
		if _, err := w.Write([]byte(responseXml)); err != nil {
			handleError(err, w)
			return
		}
	}
}
