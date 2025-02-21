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

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type imagePipeline struct {
	mosaicPipeline model.Pipeline
	imageProc      *processor.ImageProcessor
	s3             *infra.S3Manager
	taskClient     *api_client.TaskClient
	snapshotClient *api_client.SnapshotClient
	fileIdent      *identifier.FileIdentifier
	imageIdent     *identifier.ImageIdentifier
	config         *config.Config
}

func NewImagePipeline() model.Pipeline {
	return &imagePipeline{
		mosaicPipeline: NewMosaicPipeline(),
		imageProc:      processor.NewImageProcessor(),
		s3:             infra.NewS3Manager(),
		taskClient:     api_client.NewTaskClient(),
		snapshotClient: api_client.NewSnapshotClient(),
		fileIdent:      identifier.NewFileIdentifier(),
		imageIdent:     identifier.NewImageIdentifier(),
		config:         config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts api_client.PipelineRunOptions) error {
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

func (p *imagePipeline) RunFromLocalPath(inputPath string, opts api_client.PipelineRunOptions) error {
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Measuring image dimensions."),
	}); err != nil {
		return err
	}
	imageProps, err := p.measureImageDimensions(inputPath, opts)
	if err != nil {
		return err
	}
	var imagePath string
	if p.imageIdent.IsTIFF(inputPath) {
		if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
			Fields: []string{api_client.TaskFieldName},
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
			infra.GetLogger().Error(err)
		}
	}(imagePath)
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
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
		if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
			Fields: []string{api_client.TaskFieldName, api_client.TaskFieldStatus},
			Name:   helper.ToPtr("Done."),
			Status: helper.ToPtr(api_client.TaskStatusSuccess),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p *imagePipeline) measureImageDimensions(inputPath string, opts api_client.PipelineRunOptions) (*api_client.ImageProps, error) {
	imageProps, err := p.imageProc.MeasureImage(inputPath)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldOriginal},
		Original: &api_client.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   helper.ToPtr(stat.Size()),
			Image:  imageProps,
		},
	}); err != nil {
		return nil, err
	}
	return imageProps, nil
}

func (p *imagePipeline) createThumbnail(inputPath string, opts api_client.PipelineRunOptions) error {
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
				infra.GetLogger().Error(err)
			}
		}(outputPath)
	} else {
		outputPath = inputPath
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	s3Object := &api_client.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(outputPath),
		Image: &api_client.ImageProps{
			Width:  res.Width,
			Height: res.Height,
		},
		Size: helper.ToPtr(stat.Size()),
	}
	if err := p.s3.PutFile(s3Object.Key, outputPath, helper.DetectMimeFromFile(outputPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options:   opts,
		Fields:    []string{api_client.SnapshotFieldThumbnail},
		Thumbnail: s3Object,
	}); err != nil {
		return err
	}
	return nil
}

func (p *imagePipeline) convertTIFFToJPEG(inputPath string, imageProps api_client.ImageProps, opts api_client.PipelineRunOptions) (*string, error) {
	jpegPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
	if err := p.imageProc.ConvertImage(inputPath, jpegPath); err != nil {
		return nil, err
	}
	stat, err := os.Stat(jpegPath)
	if err != nil {
		return nil, err
	}
	s3Object := &api_client.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/preview.jpg",
		Size:   helper.ToPtr(stat.Size()),
		Image:  &imageProps,
	}
	if err := p.s3.PutFile(s3Object.Key, jpegPath, helper.DetectMimeFromFile(jpegPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldPreview},
		Preview: s3Object,
	}); err != nil {
		return nil, err
	}
	return &jpegPath, nil
}

func (p *imagePipeline) saveOriginalAsPreview(inputPath string, imageProps api_client.ImageProps, opts api_client.PipelineRunOptions) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldPreview},
		Preview: &api_client.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   helper.ToPtr(stat.Size()),
			Image:  &imageProps,
		},
	}); err != nil {
		return err
	}
	return nil
}
