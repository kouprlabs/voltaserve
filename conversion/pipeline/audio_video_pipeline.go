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

type audioVideoPipeline struct {
	videoProc      *processor.VideoProcessor
	imageProc      *processor.ImageProcessor
	s3             infra.S3Manager
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
	config         *config.Config
}

func NewAudioVideoPipeline() Pipeline {
	return &audioVideoPipeline{
		videoProc:      processor.NewVideoProcessor(),
		imageProc:      processor.NewImageProcessor(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		config:         config.GetConfig(),
	}
}

func (p *audioVideoPipeline) Run(opts dto.PipelineRunOptions) error {
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

func (p *audioVideoPipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	if opts.TaskID != nil {
		if _, err := p.taskClient.Patch(*opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName},
			Name:   helper.ToPtr("Creating thumbnail."),
		}); err != nil {
			return err
		}
	}
	// Here we intentionally ignore the error, as the media file may contain just audio
	// Additionally, we don't consider failing to create the thumbnail an error
	_ = p.patchThumbnail(inputPath, opts)
	if opts.TaskID != nil {
		if _, err := p.taskClient.Patch(*opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName},
			Name:   helper.ToPtr("Saving preview."),
		}); err != nil {
			return err
		}
	}
	if err := p.patchPreviewWithOriginal(inputPath, opts); err != nil {
		return err
	}
	if opts.TaskID != nil {
		if _, err := p.taskClient.Patch(*opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName, model.TaskFieldStatus},
			Name:   helper.ToPtr("Done."),
			Status: helper.ToPtr(model.TaskStatusSuccess),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p *audioVideoPipeline) patchThumbnail(inputPath string, opts dto.PipelineRunOptions) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}(outputPath)
	if err := p.videoProc.Thumbnail(inputPath, p.config.Limits.ImagePreviewMaxWidth, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return err
	}
	props, err := p.imageProc.MeasureImage(outputPath)
	if err != nil {
		return err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	s3Object := &model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(outputPath),
		Image:  props,
		Size:   stat.Size(),
	}
	if err := p.s3.PutFile(s3Object.Key, outputPath, helper.DetectMIMEFromPath(outputPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(opts.SnapshotID, dto.SnapshotPatchOptions{
		Fields:    []string{model.SnapshotFieldThumbnail},
		Thumbnail: s3Object,
	}); err != nil {
		return err
	}
	return nil
}

func (p *audioVideoPipeline) patchPreviewWithOriginal(inputPath string, opts dto.PipelineRunOptions) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(opts.SnapshotID, dto.SnapshotPatchOptions{
		Fields: []string{model.SnapshotFieldPreview},
		Preview: &model.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   stat.Size(),
		},
	}); err != nil {
		return err
	}
	return nil
}
