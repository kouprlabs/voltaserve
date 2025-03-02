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

type imagePipeline struct {
	mosaicPipeline Pipeline
	imageProc      *processor.ImageProcessor
	s3             infra.S3Manager
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
	fileIdent      *infra.FileIdentifier
	config         *config.Config
}

func NewImagePipeline() Pipeline {
	return &imagePipeline{
		mosaicPipeline: NewMosaicPipeline(),
		imageProc:      processor.NewImageProcessor(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		fileIdent:      infra.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts dto.PipelineRunOptions) error {
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

func (p *imagePipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Measuring image dimensions."),
	}); err != nil {
		return err
	}
	imageProps, err := p.measureImageDimensions(inputPath, opts)
	if err != nil {
		return err
	}
	var imagePath string
	if p.fileIdent.IsTIFF(inputPath) {
		if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName},
			Name:   helper.ToPtr("Converting TIFF image to JPEG format."),
		}); err != nil {
			return err
		}
		jpegPath, err := p.convertTIFFToJPEG(inputPath, *imageProps, opts)
		if err != nil {
			return err
		}
		imagePath = *jpegPath
	} else {
		imagePath = inputPath
		if err := p.saveOriginalAsPreview(imagePath, *imageProps, opts); err != nil {
			return err
		}
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(imagePath)
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	// We don't consider failing the creation of the thumbnail an error
	_ = p.createThumbnail(imagePath, opts)
	// Automatically trigger mosaic pipeline if the image exceeds the pixels threshold
	if imageProps.Width >= p.config.Limits.ImageMosaicTriggerThresholdPixels ||
		imageProps.Height >= p.config.Limits.ImageMosaicTriggerThresholdPixels {
		if err := p.mosaicPipeline.RunFromLocalPath(imagePath, opts); err != nil {
			return err
		}
	} else {
		if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName, model.TaskFieldStatus},
			Name:   helper.ToPtr("Done."),
			Status: helper.ToPtr(model.TaskStatusSuccess),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p *imagePipeline) measureImageDimensions(inputPath string, opts dto.PipelineRunOptions) (*model.ImageProps, error) {
	imageProps, err := p.imageProc.MeasureImage(inputPath)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}
	if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{model.SnapshotFieldOriginal},
		Original: &model.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   stat.Size(),
			Image:  imageProps,
		},
	}); err != nil {
		return nil, err
	}
	return imageProps, nil
}

func (p *imagePipeline) createThumbnail(inputPath string, opts dto.PipelineRunOptions) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(inputPath))
	res, err := p.imageProc.Thumbnail(inputPath, p.config.Limits.ImagePreviewMaxWidth, p.config.Limits.ImagePreviewMaxHeight, outputPath)
	if err != nil {
		return err
	}
	if res.IsCreated {
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				logger.GetLogger().Error(err)
			}
		}(outputPath)
	} else {
		outputPath = inputPath
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	s3Object := &model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(outputPath),
		Image: &model.ImageProps{
			Width:  res.Width,
			Height: res.Height,
		},
		Size: stat.Size(),
	}
	if err := p.s3.PutFile(s3Object.Key, outputPath, helper.DetectMIMEFromPath(outputPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options:   opts,
		Fields:    []string{model.SnapshotFieldThumbnail},
		Thumbnail: s3Object,
	}); err != nil {
		return err
	}
	return nil
}

func (p *imagePipeline) convertTIFFToJPEG(inputPath string, imageProps model.ImageProps, opts dto.PipelineRunOptions) (*string, error) {
	jpegPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
	if err := p.imageProc.ConvertImage(inputPath, jpegPath); err != nil {
		return nil, err
	}
	stat, err := os.Stat(jpegPath)
	if err != nil {
		return nil, err
	}
	s3Object := &model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/preview.jpg",
		Size:   stat.Size(),
		Image:  &imageProps,
	}
	if err := p.s3.PutFile(s3Object.Key, jpegPath, helper.DetectMIMEFromPath(jpegPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{model.SnapshotFieldPreview},
		Preview: s3Object,
	}); err != nil {
		return nil, err
	}
	return &jpegPath, nil
}

func (p *imagePipeline) saveOriginalAsPreview(inputPath string, imageProps model.ImageProps, opts dto.PipelineRunOptions) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{model.SnapshotFieldPreview},
		Preview: &model.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   stat.Size(),
			Image:  &imageProps,
		},
	}); err != nil {
		return err
	}
	return nil
}
