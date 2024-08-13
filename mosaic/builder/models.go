// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package builder

import (
	"fmt"
	"path/filepath"

	"github.com/h2non/bimg"
)

type Metadata struct {
	Width      int         `json:"width"`
	Height     int         `json:"height"`
	Extension  string      `json:"extension"`
	ZoomLevels []ZoomLevel `json:"zoomLevels"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Rectangle struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type MinimumScaleSize struct {
	Value Size `json:"value"`
}

func NewMinimumScaleSize(value Size) (*MinimumScaleSize, error) {
	if !IsValidSize(value) {
		return nil, fmt.Errorf((MinimumScaleSize{}).GetAcceptanceCriteria())
	}
	return &MinimumScaleSize{Value: value}, nil
}

func (m MinimumScaleSize) Width() int {
	return m.Value.Width
}

func (m MinimumScaleSize) Height() int {
	return m.Value.Height
}

func IsValidSize(value Size) bool {
	return value.Width > 0 && value.Height > 0
}

func (m MinimumScaleSize) GetAcceptanceCriteria() string {
	return "Width and Height of MinimumScaleSize should be bigger than zero."
}

type Region struct {
	ColStart               int  `json:"colStart"`
	ColEnd                 int  `json:"colEnd"`
	RowStart               int  `json:"rowStart"`
	RowEnd                 int  `json:"rowEnd"`
	IncludesRemainingTiles bool `json:"includesRemainingTiles"`
}

func (r *Region) IsNull() bool {
	return r.ColStart == 0 && r.ColEnd == 0 && r.RowStart == 0 && r.RowEnd == 0
}

type ScaleDownPercentage struct {
	Value  uint16 `json:"value"`
	factor *float64
}

func NewScaleDownPercentage(value uint16) (*ScaleDownPercentage, error) {
	s := &ScaleDownPercentage{Value: value}
	if !s.isValid() {
		return nil, fmt.Errorf(s.GetAcceptanceCriteria())
	}
	return s, nil
}

func (s ScaleDownPercentage) Factor() float64 {
	if s.factor == nil {
		factor := float64(s.Value) * 0.01
		s.factor = &factor
	}
	return *s.factor
}

func (s ScaleDownPercentage) isValid() bool {
	return s.Value > 0 && s.Value < 100
}

func (s ScaleDownPercentage) GetAcceptanceCriteria() string {
	return "ScaleDownPercentage should be exclusively more than 0, and exclusively less than 100."
}

type Tile struct {
	Width         int `json:"width"`
	Height        int `json:"height"`
	LastColWidth  int `json:"lastColWidth"`
	LastRowHeight int `json:"lastRowHeight"`
}

type TileMetadata struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
	Row    int `json:"row"`
	Col    int `json:"col"`
}

type TileSize struct {
	Value Size `json:"value"`
}

func NewTileSize(value Size) (*TileSize, error) {
	t := &TileSize{Value: value}
	if !t.IsValid() {
		return nil, fmt.Errorf(t.GetAcceptanceCriteria())
	}
	return t, nil
}

func (t *TileSize) Width() int {
	return t.Value.Width
}

func (t *TileSize) SetWidth(width int) {
	t.Value = Size{Width: width, Height: t.Value.Height}
}

func (t *TileSize) Height() int {
	return t.Value.Height
}

func (t *TileSize) SetHeight(height int) {
	t.Value = Size{Width: t.Value.Width, Height: height}
}

func (t *TileSize) IsValid() bool {
	return t.IsValidWidth(t.Value.Width) && t.IsValidHeight(t.Value.Height)
}

func (t *TileSize) IsValidWidth(width int) bool {
	return width > 0
}

func (t *TileSize) IsValidHeight(height int) bool {
	return height > 0
}

func (t *TileSize) GetAcceptanceCriteria() string {
	return "Width and Height of TileSize should be greater than zero."
}

type ZoomLevel struct {
	Index               int     `json:"index"`
	Width               int     `json:"width"`
	Height              int     `json:"height"`
	Rows                int     `json:"rows"`
	Cols                int     `json:"cols"`
	ScaleDownPercentage float32 `json:"scaleDownPercentage"`
	Tile                Tile    `json:"tile"`
}

type Image struct {
	buffer []byte
	file   string
	width  *int
	height *int
}

func NewImage(file string) (*Image, error) {
	b, err := bimg.Read(file)
	if err != nil {
		return nil, err
	}
	return &Image{
		buffer: b,
		file:   file,
	}, nil
}

func NewImageFromSource(source *Image) (*Image, error) {
	if source == nil {
		return nil, fmt.Errorf("source image is nil")
	}
	return &Image{
		buffer: source.buffer,
		file:   source.file,
	}, nil
}

func (img *Image) Width() int {
	if img.width == nil {
		size, err := bimg.Size(img.buffer)
		if err != nil {
			return -1
		}
		img.width = new(int)
		*img.width = size.Width
	}
	return *img.width
}

func (img *Image) Height() int {
	if img.height == nil {
		size, err := bimg.Size(img.buffer)
		if err != nil {
			return -1
		}
		img.height = new(int)
		*img.height = size.Height
	}
	return *img.height
}

func (img *Image) Extension() string {
	return filepath.Ext(img.file)
}

func (img *Image) Crop(x, y, width, height int) error {
	var err error
	img.buffer, err = bimg.NewImage(img.buffer).Extract(y, x, width, height)
	if err != nil {
		return err
	}
	return nil
}

func (img *Image) ScaleWithAspectRatio(width, height int) error {
	var err error
	img.buffer, err = bimg.Resize(img.buffer, bimg.Options{
		Width:  width,
		Height: height,
	})
	if err != nil {
		return err
	}
	return nil
}

func (img *Image) Save(file string) error {
	if err := bimg.Write(file, img.buffer); err != nil {
		return err
	}
	return nil
}
