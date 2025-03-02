// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package dto

type MosaicMetadata struct {
	Width      int               `json:"width"`
	Height     int               `json:"height"`
	Extension  string            `json:"extension"`
	ZoomLevels []MosaicZoomLevel `json:"zoomLevels"`
}

type MosaicZoomLevel struct {
	Index               int        `json:"index"`
	Width               int        `json:"width"`
	Height              int        `json:"height"`
	Rows                int        `json:"rows"`
	Cols                int        `json:"cols"`
	ScaleDownPercentage float32    `json:"scaleDownPercentage"`
	Tile                MosaicTile `json:"tile"`
}

type MosaicTile struct {
	Width         int `json:"width"`
	Height        int `json:"height"`
	LastColWidth  int `json:"lastColWidth"`
	LastRowHeight int `json:"lastRowHeight"`
}
