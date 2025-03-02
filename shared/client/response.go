// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
)

func JsonResponseOrError(response *http.Response) ([]byte, error) {
	contentType := response.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		if response.StatusCode > 299 {
			var errorResponse errorpkg.ErrorResponse
			if err := json.Unmarshal(b, &errorResponse); err != nil {
				return nil, err
			}
			return nil, &errorResponse
		} else {
			return b, nil
		}
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("unexpected response Content-Type: %s", contentType))
	}
}

func TextResponseOrError(response *http.Response) ([]byte, error) {
	contentType := response.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/plain") {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		if response.StatusCode > 299 {
			var errorResponse errorpkg.ErrorResponse
			if err := json.Unmarshal(b, &errorResponse); err != nil {
				return nil, err
			}
			return nil, &errorResponse
		} else {
			return b, nil
		}
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("unexpected response Content-Type: %s", contentType))
	}
}

func ByteResponseOrError(response *http.Response) ([]byte, error) {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 299 {
		var errorResponse errorpkg.ErrorResponse
		if err := json.Unmarshal(b, &errorResponse); err != nil {
			return nil, err
		}
		return nil, &errorResponse
	} else {
		return b, nil
	}
}

func ByteResponseWithWriterOrError(response *http.Response, w io.Writer) error {
	if response.StatusCode > 299 {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var errorResponse errorpkg.ErrorResponse
		if err := json.Unmarshal(b, &errorResponse); err != nil {
			return err
		}
		return &errorResponse
	} else {
		if _, err := io.Copy(w, response.Body); err != nil {
			return err
		}
		return nil
	}
}

func SuccessfulResponseOrError(response *http.Response) error {
	if response.StatusCode > 299 {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var errorResponse errorpkg.ErrorResponse
		if err := json.Unmarshal(b, &errorResponse); err != nil {
			return err
		}
		return &errorResponse
	} else {
		return nil
	}
}
