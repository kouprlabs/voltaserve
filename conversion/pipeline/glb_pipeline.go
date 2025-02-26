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

	"github.com/kouprlabs/voltaserve/api/client/apiclient"
	apiinfra "github.com/kouprlabs/voltaserve/api/infra"
	apimodel "github.com/kouprlabs/voltaserve/api/model"
	apiservice "github.com/kouprlabs/voltaserve/api/service"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type glbPipeline struct {
	glbProc        *processor.GLBProcessor
	imageProc      *processor.ImageProcessor
	s3             apiinfra.S3Manager
	taskClient     *apiclient.TaskClient
	snapshotClient *apiclient.SnapshotClient
	config         *config.Config
}

func NewGLBPipeline() model.Pipeline {
	return &glbPipeline{
		glbProc:        processor.NewGLBProcessor(),
		imageProc:      processor.NewImageProcessor(),
		s3:             apiinfra.NewS3Manager(),
		taskClient:     apiclient.NewTaskClient(),
		snapshotClient: apiclient.NewSnapshotClient(),
		config:         config.GetConfig(),
	}
}

func (p *glbPipeline) Run(opts model.PipelineRunOptions) error {
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

func (p *glbPipeline) RunFromLocalPath(inputPath string, opts model.PipelineRunOptions) error {
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
		Name:   helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	// We don't consider failing to create the thumbnail an error
	_ = p.createThumbnail(inputPath, opts)
	if err := p.patchSnapshotPreviewField(inputPath, opts); err != nil {
		return err
	}
	return nil
}

func (p *glbPipeline) patchSnapshotPreviewField(inputPath string, opts model.PipelineRunOptions) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if filepath.Ext(inputPath) == filepath.Ext(opts.Key) {
		/* The original is a .glb */
		if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{apimodel.SnapshotFieldPreview},
			Preview: &apimodel.S3Object{
				Bucket: opts.Bucket,
				Key:    opts.Key,
				Size:   stat.Size(),
			},
		}); err != nil {
			return err
		}
	} else {
		/* The original is likely an .zip glTF file */
		if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{apimodel.SnapshotFieldPreview},
			Preview: &apimodel.S3Object{
				Bucket: opts.Bucket,
				Key:    filepath.FromSlash(opts.SnapshotID + "/preview" + filepath.Ext(inputPath)),
				Size:   stat.Size(),
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p *glbPipeline) createThumbnail(inputPath string, opts model.PipelineRunOptions) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(outputPath)
	if err := p.glbProc.Thumbnail(inputPath, p.config.Limits.ImagePreviewMaxWidth, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return err
	}
	stat, err := os.Stat(outputPath)
	if err == nil {
		props, err := p.imageProc.MeasureImage(outputPath)
		if err != nil {
			return err
		}
		s3Object := &apimodel.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(outputPath),
			Image:  props,
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(s3Object.Key, outputPath, helper.DetectMimeFromFile(outputPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
			return err
		}
		if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options:   opts,
			Fields:    []string{apimodel.SnapshotFieldThumbnail},
			Thumbnail: s3Object,
		}); err != nil {
			return err
		}
	}
	return nil
}
