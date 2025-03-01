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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kouprlabs/voltaserve/shared/errorpkg"
)

func JsonResponseOrError(resp *http.Response) ([]byte, error) {
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode > 299 {
			var errorResponse ErrorResponse
			if err := json.Unmarshal(body, &errorResponse); err != nil {
				return nil, err
			}
			return nil, &errorResponse
		} else {
			return body, nil
		}
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("unexpected response Content-Type: %s", contentType))
	}
}

func TextResponseOrError(resp *http.Response) ([]byte, error) {
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/plain") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode > 299 {
			var errorResponse ErrorResponse
			if err := json.Unmarshal(body, &errorResponse); err != nil {
				return nil, err
			}
			return nil, &errorResponse
		} else {
			return body, nil
		}
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("unexpected response Content-Type: %s", contentType))
	}
}

func OctetStreamResponseOrError(resp *http.Response) ([]byte, error) {
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/octet-stream") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode > 299 {
			var errorResponse ErrorResponse
			if err := json.Unmarshal(body, &errorResponse); err != nil {
				return nil, err
			}
			return nil, &errorResponse
		} else {
			buf := &bytes.Buffer{}
			_, err := io.Copy(buf, resp.Body)
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
	} else {
		return nil, errorpkg.NewInternalServerError(fmt.Errorf("unexpected response Content-Type: %s", contentType))
	}
}

func OctetStreamResponseWithWriterOrThrow(resp *http.Response, w io.Writer) error {
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/octet-stream") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode > 299 {
			var errorResponse ErrorResponse
			if err := json.Unmarshal(body, &errorResponse); err != nil {
				return err
			}
			return &errorResponse
		} else {
			if _, err := io.Copy(w, resp.Body); err != nil {
				return err
			}
			return nil
		}
	} else {
		return errorpkg.NewInternalServerError(fmt.Errorf("unexpected response Content-Type: %s", contentType))
	}
}

func SuccessfulResponseOrThrow(resp *http.Response) error {
	if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return err
		}
		return &errorResponse
	} else {
		return nil
	}
}
