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
	"net/http"
)

/*
This method updates the properties of a resource.

Example implementation:

- Parse the request body to extract the properties to be updated.
- Read the existing data from the file.
- Parse the existing properties.
- Merge the updated properties with the existing ones.
- Format the updated properties and store them back in the file.
- Set the response status code to 204 if successful or an appropriate error code if the file is not found or encountered an error.
- Return the response.

In this example implementation, the handleProppatch() method first parses the XML
payload containing the properties to be updated. Then, it reads the existing data from the file,
parses the existing properties (assuming an XML format),
merges the updated properties with the existing ones, and formats
the properties back into the desired format (e.g., XML).

Finally, the updated properties are written back to the file.
You can customize the parseProperties() and formatProperties()
functions to match the specific property format you are using in your WebDAV server.

Note that this implementation assumes a simplified example and may require further
customization based on your specific property format and requirements.
*/
func (h *Handler) methodProppatch(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	if _, err := w.Write([]byte(http.StatusText(http.StatusNotImplemented))); err != nil {
		return
	}
}
