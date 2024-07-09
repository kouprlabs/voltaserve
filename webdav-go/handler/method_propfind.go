package handler

import (
	"fmt"
	"net/http"
	"voltaserve/client"
	"voltaserve/helper"
	"voltaserve/infra"
)

/*
This method retrieves properties and metadata of a resource.

Example implementation:

- Extract the file path from the URL.
- Use fs.stat() to retrieve the file metadata.
- Format the response body in the desired XML format with the properties and metadata.
- Set the response status code to 207 if successful or an appropriate error code if the file is not found or encountered an error.
- Set the Content-Type header to indicate the XML format.
- Return the response.
*/
func (h *Handler) methodPropfind(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*infra.Token)
	if !ok {
		infra.HandleError(fmt.Errorf("missing token"), w)
		return
	}
	apiClient := client.NewAPIClient(token)
	file, err := apiClient.GetFileByPath(helper.DecodeURIComponent(r.URL.Path))
	if err != nil {
		infra.HandleError(err, w)
		return
	}
	if file.Type == client.FileTypeFile {
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
			func() int {
				if file.Type == client.FileTypeFile && file.Snapshot != nil && file.Snapshot.Original != nil {
					return file.Snapshot.Original.Size
				}
				return 0
			}(),
			helper.ToUTCString(&file.CreateTime),
			helper.ToUTCString(file.UpdateTime),
		)
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		w.WriteHeader(http.StatusMultiStatus)
		if _, err := w.Write([]byte(responseXml)); err != nil {
			infra.HandleError(err, w)
			return
		}
	} else if file.Type == client.FileTypeFolder {
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
		list, err := apiClient.ListFilesByPath(helper.DecodeURIComponent(r.URL.Path))
		if err != nil {
			infra.HandleError(err, w)
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
					if item.Type == client.FileTypeFolder {
						return "<D:collection/>"
					}
					return ""
				}(),
				func() int {
					if item.Type == client.FileTypeFile && item.Snapshot != nil && item.Snapshot.Original != nil {
						return item.Snapshot.Original.Size
					}
					return 0
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
			infra.HandleError(err, w)
			return
		}
	}
}
