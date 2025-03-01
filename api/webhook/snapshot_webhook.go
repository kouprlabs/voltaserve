// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/logger"
)

type SnapshotWebhook struct {
	config *config.Config
}

func NewSnapshotWebhook() *SnapshotWebhook {
	return &SnapshotWebhook{
		config: config.GetConfig(),
	}
}

func (wh *SnapshotWebhook) Call(snapshot model.Snapshot, eventType string) error {
	body, err := json.Marshal(dto.SnapshotWebhookOptions{
		EventType: eventType,
		Snapshot:  snapshot,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s?api_key=%s", config.GetConfig().Webhook.Snapshot, wh.config.Security.APIKey),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			logger.GetLogger().Error(err)
		}
	}(resp.Body)
	return client.SuccessfulResponseOrThrow(resp)
}
