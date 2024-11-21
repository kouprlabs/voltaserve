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
	"os"
	"path/filepath"

	"github.com/kouprlabs/voltaserve/mosaic/infra"
)

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

func (mb *MosaicBuilder) Build() (*Metadata, error) {
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

	var zoomLevels []ZoomLevel
	for _, index := range zoomLevelsIndexes {
		mb.CreateZoomLevelDirectory(index)
		scaled, err := mb.Scale(index)
		if err != nil {
			return nil, err
		}
		zoomLevel := mb.Decompose(scaled, index, Region{})
		zoomLevels = append(zoomLevels, zoomLevel)
	}

	metadata := &Metadata{
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

func (mb *MosaicBuilder) Decompose(image *Image, zoomLevel int, region Region) ZoomLevel {
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
				return ZoomLevel{}
			}
			if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, tileMetadata.Row, tileMetadata.Col)); err != nil {
				return ZoomLevel{}
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
				return ZoomLevel{}
			}
			if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, totalRows-1, c)); err != nil {
				return ZoomLevel{}
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
				return ZoomLevel{}
			}
			if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, r, totalCols-1)); err != nil {
				return ZoomLevel{}
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
			return ZoomLevel{}
		}
		if err := cropped.Save(mb.GetTileOutputPath(zoomLevel, totalRows-1, totalCols-1)); err != nil {
			return ZoomLevel{}
		}
	}

	return ZoomLevel{
		Index:               zoomLevel,
		Width:               image.Width(),
		Height:              image.Height(),
		Rows:                totalRows,
		Cols:                totalCols,
		ScaleDownPercentage: float32(mb.GetScaleDownPercentage(zoomLevel)),
		Tile: Tile{
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
