// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package pipeline

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/logger"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type mosaicPipeline struct {
	videoProc      *processor.VideoProcessor
	imageProc      *processor.ImageProcessor
	fileIdent      *infra.FileIdentifier
	s3             infra.S3Manager
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
	mosaicClient   *client.MosaicClient
}

func NewMosaicPipeline() Pipeline {
	return &mosaicPipeline{
		videoProc:      processor.NewVideoProcessor(),
		imageProc:      processor.NewImageProcessor(),
		fileIdent:      infra.NewFileIdentifier(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		mosaicClient:   client.NewMosaicClient(config.GetConfig().MosaicURL),
	}
}

func (p *mosaicPipeline) Run(opts dto.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket, minio.GetObjectOptions{}); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(inputPath)
	return p.RunFromLocalPath(inputPath, opts)
}

func (p *mosaicPipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	if !p.fileIdent.IsImage(opts.Key) {
		return errors.New("unsupported file type")
	}
	if err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Creating mosaic."),
	}); err != nil {
		return err
	}
	if !p.imageProc.IsSupportedByBild(inputPath) {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
		if err := p.imageProc.ConvertImage(inputPath, outputPath); err != nil {
			return err
		}
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				logger.GetLogger().Error(err)
			}
		}(outputPath)
		inputPath = outputPath
	}
	if _, err := p.mosaicClient.Create(client.MosaicCreateOptions{
		Path:     inputPath,
		S3Key:    filepath.FromSlash(opts.SnapshotID),
		S3Bucket: opts.Bucket,
	}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{model.SnapshotFieldMosaic},
		Mosaic: &model.S3Object{
			Key:    filepath.FromSlash(opts.SnapshotID + "/mosaic"),
			Bucket: opts.Bucket,
		},
	}); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName, model.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(model.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}
