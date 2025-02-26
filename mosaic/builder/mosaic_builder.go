// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package builder

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"

	"github.com/kouprlabs/voltaserve/mosaic/infra"
	"github.com/kouprlabs/voltaserve/mosaic/model"
)

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
		return nil, fmt.Errorf("%s", (MinimumScaleSize{}).GetAcceptanceCriteria())
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
		return nil, fmt.Errorf("%s", s.GetAcceptanceCriteria())
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
		return nil, fmt.Errorf("%s", t.GetAcceptanceCriteria())
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

type Image struct {
	img  image.Image
	file string
}

func NewImage(file string) (*Image, error) {
	img, err := imgio.Open(file)
	if err != nil {
		return nil, err
	}
	return &Image{
		img:  img,
		file: file,
	}, nil
}

func NewImageFromSource(source *Image) (*Image, error) {
	if source == nil {
		return nil, fmt.Errorf("source image is nil")
	}
	return &Image{
		img:  source.img,
		file: source.file,
	}, nil
}

func (img *Image) Width() int {
	return img.img.Bounds().Dx()
}

func (img *Image) Height() int {
	return img.img.Bounds().Dy()
}

func (img *Image) Extension() string {
	return filepath.Ext(img.file)
}

func (img *Image) Crop(x, y, width, height int) error {
	img.img = transform.Crop(
		img.img,
		image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + width, Y: y + height},
		})
	return nil
}

func (img *Image) ScaleWithAspectRatio(width, height int) error {
	img.img = transform.Resize(img.img, width, height, transform.Lanczos)
	return nil
}

func (img *Image) Save(file string) error {
	var encoder imgio.Encoder
	if strings.HasSuffix(file, ".png") {
		encoder = imgio.PNGEncoder()
	} else if strings.HasSuffix(file, ".jpg") {
		encoder = imgio.JPEGEncoder(100)
	} else {
		return fmt.Errorf("unsupported image format: %s", file)
	}
	return imgio.Save(file, img.img, encoder)
}

const (
	ActionOnExistingDirectoryDelete = "delete"
	ActionOnExistingDirectorySkip   = "skip"
)

type MosaicBuilder struct {
	image                     *Image
	scaleDownPercentage       *ScaleDownPercentage
	minimumScaleSize          *MinimumScaleSize
	tileSize                  *TileSize
	options                   MosaicBuilderOptions
	actionOnExistingDirectory string
}

type MosaicBuilderOptions struct {
	File            string
	OutputDirectory string
}

func NewMosaicBuilder(opts MosaicBuilderOptions) *MosaicBuilder {
	return &MosaicBuilder{
		options: opts,
	}
}

func (mb *MosaicBuilder) ScaleDownPercentage() *ScaleDownPercentage {
	if mb.scaleDownPercentage == nil {
		mb.scaleDownPercentage, _ = NewScaleDownPercentage(70)
	}
	return mb.scaleDownPercentage
}

func (mb *MosaicBuilder) SetScaleDownPercentage(scaleDownPercentage *ScaleDownPercentage) {
	mb.scaleDownPercentage = scaleDownPercentage
}

func (mb *MosaicBuilder) MinimumScaleSize() *MinimumScaleSize {
	if mb.minimumScaleSize == nil {
		size := Size{Width: 500, Height: 500}
		mb.minimumScaleSize, _ = NewMinimumScaleSize(size)
	}
	return mb.minimumScaleSize
}

func (mb *MosaicBuilder) SetMinimumScaleSize(minimumScaleSize *MinimumScaleSize) {
	mb.minimumScaleSize = minimumScaleSize
}

func (mb *MosaicBuilder) TileSize() *TileSize {
	if mb.tileSize == nil {
		size := Size{Width: 300, Height: 300}
		mb.tileSize, _ = NewTileSize(size)
	}
	return mb.tileSize
}

func (mb *MosaicBuilder) SetTileSize(tileSize *TileSize) {
	mb.tileSize = tileSize
}

func (mb *MosaicBuilder) SetActionOnExistingDirectory(action string) {
	mb.actionOnExistingDirectory = action
}

func (mb *MosaicBuilder) Build() (*model.Metadata, error) {
	cleanupIfFails := false
	if _, err := os.Stat(mb.options.OutputDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(mb.options.OutputDirectory, 0o750); err != nil {
			return nil, err
		}
		cleanupIfFails = true
	}
	defer func() {
		if r := recover(); r != nil {
			if cleanupIfFails {
				if err := os.RemoveAll(mb.options.OutputDirectory); err != nil {
					infra.GetLogger().Error(err)
					return
				}
			}
		}
	}()

	image, err := NewImage(mb.options.File)
	if err != nil {
		return nil, err
	}
	mb.image = image

	zoomLevelsIndexes := mb.RequiredZoomLevelIndexes()
	if len(zoomLevelsIndexes) == 0 {
		return nil, fmt.Errorf("creating zoom levels is not required for this image")
	}

	var zoomLevels []model.ZoomLevel
	for _, index := range zoomLevelsIndexes {
		mb.CreateZoomLevelDirectory(index)
		scaled, err := mb.Scale(index)
		if err != nil {
			return nil, err
		}
		zoomLevel := mb.Decompose(scaled, index, Region{})
		zoomLevels = append(zoomLevels, zoomLevel)
	}

	metadata := &model.Metadata{
		Width:      mb.image.Width(),
		Height:     mb.image.Height(),
		Extension:  filepath.Ext(mb.options.File),
		ZoomLevels: zoomLevels,
	}

	metadataFilePath := mb.GetMetadataFilePath()
	metadataBytes, _ := json.MarshalIndent(metadata, "", "  ")
	if err := os.WriteFile(metadataFilePath, metadataBytes, 0o600); err != nil {
		return nil, err
	}

	return metadata, nil
}

func (mb *MosaicBuilder) Decompose(image *Image, zoomLevel int, region Region) model.ZoomLevel {
	tileWidthExceeded := image.Width() > mb.TileSize().Width()
	tileHeightExceeded := image.Height() > mb.TileSize().Height()

	cols := 1
	if tileWidthExceeded {
		cols = image.Width() / mb.TileSize().Width()
	}
	rows := 1
	if tileHeightExceeded {
		rows = image.Height() / mb.TileSize().Height()
	}
	remainingWidth := 0
	if tileWidthExceeded {
		remainingWidth = image.Width() - (mb.TileSize().Width() * cols)
	}
	remainingHeight := 0
	if tileHeightExceeded {
		remainingHeight = image.Height() - (mb.TileSize().Height() * rows)
	}
	totalCols := cols
	if remainingWidth != 0 {
		totalCols = cols + 1
	}
	totalRows := rows
	if remainingHeight != 0 {
		totalRows = rows + 1
	}

	adaptedTileSize := *mb.TileSize()
	if !tileWidthExceeded {
		adaptedTileSize.SetWidth(image.Width())
	}
	if !tileHeightExceeded {
		adaptedTileSize.SetHeight(image.Height())
	}

	colStart, colEnd, rowStart, rowEnd := 0, cols-1, 0, rows-1
	includesRemainingTiles := true
	if !region.IsNull() {
		colStart = region.ColStart
		colEnd = region.ColEnd
		rowStart = region.RowStart
		rowEnd = region.RowEnd
		includesRemainingTiles = region.IncludesRemainingTiles
	}

	for c := colStart; c <= colEnd; c++ {
		for r := rowStart; r <= rowEnd; r++ {
			tileMetadata := TileMetadata{
				X:      c * mb.tileSize.Width(),
				Y:      r * mb.tileSize.Height(),
				Width:  mb.tileSize.Width(),
				Height: mb.tileSize.Height(),
				Row:    r,
				Col:    c,
			}
			clippingRect := Rectangle{
				X:      tileMetadata.X,
				Y:      tileMetadata.Y,
				Width:  tileMetadata.Width,
				Height: tileMetadata.Height,
			}
			cropped, _ := NewImageFromSource(image)
			if err := cropped.Crop(clippingRect.X, clippingRect.Y, clippingRect.Width, clippingRect.Height); err != nil {
				return model.ZoomLevel{}
			}
			if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, tileMetadata.Row, tileMetadata.Col)); err != nil {
				return model.ZoomLevel{}
			}
		}
	}

	if includesRemainingTiles && remainingHeight > 0 {
		for c := 0; c < cols; c++ {
			clippingRect := Rectangle{
				X:      c * mb.tileSize.Width(),
				Y:      image.Height() - remainingHeight,
				Width:  mb.tileSize.Width(),
				Height: remainingHeight,
			}
			cropped, _ := NewImageFromSource(image)
			if err := cropped.Crop(clippingRect.X, clippingRect.Y, clippingRect.Width, clippingRect.Height); err != nil {
				return model.ZoomLevel{}
			}
			if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, totalRows-1, c)); err != nil {
				return model.ZoomLevel{}
			}
		}
	}

	if includesRemainingTiles && remainingWidth > 0 {
		for r := 0; r < rows; r++ {
			clippingRect := Rectangle{
				X:      image.Width() - remainingWidth,
				Y:      r * mb.tileSize.Height(),
				Width:  remainingWidth,
				Height: mb.tileSize.Height(),
			}
			cropped, _ := NewImageFromSource(image)
			if err := cropped.Crop(clippingRect.X, clippingRect.Y, clippingRect.Width, clippingRect.Height); err != nil {
				return model.ZoomLevel{}
			}
			if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, r, totalCols-1)); err != nil {
				return model.ZoomLevel{}
			}
		}
	}

	if includesRemainingTiles && remainingWidth > 0 && remainingHeight > 0 {
		clippingRect := Rectangle{
			X:      image.Width() - remainingWidth,
			Y:      image.Height() - remainingHeight,
			Width:  remainingWidth,
			Height: remainingHeight,
		}
		cropped, _ := NewImageFromSource(image)
		if err := cropped.Crop(clippingRect.X, clippingRect.Y, clippingRect.Width, clippingRect.Height); err != nil {
			return model.ZoomLevel{}
		}
		if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, totalRows-1, totalCols-1)); err != nil {
			return model.ZoomLevel{}
		}
	}

	return model.ZoomLevel{
		Index:               zoomLevel,
		Width:               image.Width(),
		Height:              image.Height(),
		Rows:                totalRows,
		Cols:                totalCols,
		ScaleDownPercentage: float32(mb.GetScaleDownPercentage(zoomLevel)),
		Tile: model.Tile{
			Width:         adaptedTileSize.Width(),
			Height:        adaptedTileSize.Height(),
			LastColWidth:  remainingWidth,
			LastRowHeight: remainingHeight,
		},
	}
}

func (mb *MosaicBuilder) GetScaleDownPercentage(zoomLevel int) float64 {
	value := 100.0
	for i := 0; i < zoomLevel; i++ {
		value *= 0.70
	}
	return value
}

func (mb *MosaicBuilder) Scale(zoomLevel int) (*Image, error) {
	imageSizeForZoomLevel := mb.GetImageSizeForZoomLevel(zoomLevel)
	scaled, err := NewImageFromSource(mb.image)
	if err != nil {
		return nil, err
	}
	err = scaled.ScaleWithAspectRatio(imageSizeForZoomLevel.Width, imageSizeForZoomLevel.Height)
	if err != nil {
		return nil, err
	}
	return scaled, nil
}

func (mb *MosaicBuilder) GetImageSizeForZoomLevel(zoomLevel int) Size {
	size := Size{Width: mb.image.Width(), Height: mb.image.Height()}
	counter := 0
	for {
		if counter == zoomLevel {
			break
		}
		size.Width = int(float64(size.Width) * mb.ScaleDownPercentage().Factor())
		size.Height = int(float64(size.Height) * mb.ScaleDownPercentage().Factor())
		counter += 1
	}
	return size
}

func (mb *MosaicBuilder) RequiredZoomLevelIndexes() []int {
	var levels []int
	zoomLevelCount := 0
	imageSize := Size{Width: mb.image.Width(), Height: mb.image.Height()}
	for {
		imageSize.Width = int(float64(imageSize.Width) * mb.ScaleDownPercentage().Factor())
		imageSize.Height = int(float64(imageSize.Height) * mb.ScaleDownPercentage().Factor())
		if imageSize.Width < mb.MinimumScaleSize().Width() && imageSize.Height < mb.MinimumScaleSize().Height() {
			break
		}
		levels = append(levels, zoomLevelCount)
		zoomLevelCount += 1
	}
	return levels
}

func (mb *MosaicBuilder) GetMetadataFilePath() string {
	return filepath.Join(mb.options.OutputDirectory, "mosaic.json")
}

func (mb *MosaicBuilder) GetTileOutputPath(zoomLevel, row, col int) string {
	extension := filepath.Ext(mb.options.File)
	if extension == "" {
		extension = mb.image.Extension()
	}
	if extension[0] == '.' {
		extension = extension[1:]
	}
	return filepath.Join(mb.options.OutputDirectory, fmt.Sprintf("%d/%dx%d.%s", zoomLevel, row, col, extension))
}

func (mb *MosaicBuilder) GetZoomLevelDirectoryPath(zoomLevel int) string {
	return filepath.Join(mb.options.OutputDirectory, fmt.Sprintf("%d", zoomLevel))
}

func (mb *MosaicBuilder) CreateZoomLevelDirectory(zoomLevel int) {
	mb.CreateDirectory(mb.GetZoomLevelDirectoryPath(zoomLevel))
}

func (mb *MosaicBuilder) CreateDirectory(directory string) {
	if _, err := os.Stat(directory); err == nil {
		if mb.actionOnExistingDirectory == ActionOnExistingDirectoryDelete {
			mb.DeleteDirectoryWithContent(directory)
		} else if mb.actionOnExistingDirectory == ActionOnExistingDirectorySkip {
			return
		}
	}
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return
	}
}

func (mb *MosaicBuilder) DeleteDirectoryWithContent(directory string) {
	files, _ := os.ReadDir(directory)
	for _, file := range files {
		filePath := filepath.Join(directory, file.Name())
		if file.IsDir() {
			mb.DeleteDirectoryWithContent(filePath)
		} else {
			if err := os.Remove(filePath); err != nil {
				return
			}
		}
	}
	if err := os.Remove(directory); err != nil {
		return
	}
}
