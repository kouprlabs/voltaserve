// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package pipeline

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/client/mosaic_client"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type mosaicPipeline struct {
	videoProc      *processor.VideoProcessor
	fileIdent      *identifier.FileIdentifier
	s3             *infra.S3Manager
	taskClient     *api_client.TaskClient
	snapshotClient *api_client.SnapshotClient
	mosaicClient   *mosaic_client.MosaicClient
}

func NewMosaicPipeline() model.Pipeline {
	return &mosaicPipeline{
		videoProc:      processor.NewVideoProcessor(),
		fileIdent:      identifier.NewFileIdentifier(),
		s3:             infra.NewS3Manager(),
		taskClient:     api_client.NewTaskClient(),
		snapshotClient: api_client.NewSnapshotClient(),
		mosaicClient:   mosaic_client.NewMosaicClient(),
	}
}

func (p *mosaicPipeline) Run(opts api_client.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket, minio.GetObjectOptions{}); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(inputPath)
	return p.RunFromLocalPath(inputPath, opts)
}

func (p *mosaicPipeline) RunFromLocalPath(inputPath string, opts api_client.PipelineRunOptions) error {
	if !p.fileIdent.IsImage(opts.Key) {
		return errors.New("unsupported file type")
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Creating mosaic."),
	}); err != nil {
		return err
	}
	metadata, err := p.mosaicClient.Create(mosaic_client.MosaicCreateOptions{
		Path:     inputPath,
		S3Key:    filepath.FromSlash(opts.SnapshotID),
		S3Bucket: opts.Bucket,
	})
	if err != nil {
		return err
	}
	var zoomLevels []api_client.ZoomLevel
	for _, level := range metadata.ZoomLevels {
		zoomLevels = append(zoomLevels, api_client.ZoomLevel{
			Index:               level.Index,
			Width:               level.Width,
			Height:              level.Height,
			Rows:                level.Rows,
			Cols:                level.Cols,
			ScaleDownPercentage: level.ScaleDownPercentage,
			Tile: api_client.Tile{
				Width:         level.Tile.Width,
				Height:        level.Tile.Height,
				LastColWidth:  level.Tile.LastColWidth,
				LastRowHeight: level.Tile.LastRowHeight,
			},
		})
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldMosaic},
		Mosaic: &api_client.S3Object{
			Key:    filepath.FromSlash(opts.SnapshotID + "/mosaic"),
			Bucket: opts.Bucket,
			Image: &api_client.ImageProps{
				Width:      metadata.Width,
				Height:     metadata.Height,
				ZoomLevels: zoomLevels,
			},
		},
	}); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName, api_client.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(api_client.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}
