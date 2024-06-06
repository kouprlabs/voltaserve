package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/processor"
)

type imagePipeline struct {
	imageProc *processor.ImageProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	config    config.Config
}

func NewImagePipeline() core.Pipeline {
	return &imagePipeline{
		imageProc: processor.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		config:    config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts core.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(inputPath)
	if err := p.create(inputPath, opts); err != nil {
		return err
	}
	return nil
}

func (p *imagePipeline) create(inputPath string, opts core.PipelineRunOptions) error {
	imageProps, err := p.imageProc.MeasureImage(inputPath)
	if err != nil {
		return err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	updateOpts := core.SnapshotUpdateOptions{
		Options: opts,
		Original: &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   helper.ToPtr(stat.Size()),
			Image:  &imageProps,
		},
	}
	if filepath.Ext(inputPath) == ".tiff" {
		jpegPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
		if err := p.imageProc.ConvertImage(inputPath, jpegPath); err != nil {
			return err
		}
		defer func(path string) {
			_, err := os.Stat(path)
			if os.IsExist(err) {
				if err := os.Remove(path); err != nil {
					infra.GetLogger().Error(err)
				}
			}
		}(jpegPath)
		stat, err := os.Stat(jpegPath)
		if err != nil {
			return err
		}
		thumbnail, err := p.imageProc.Base64Thumbnail(jpegPath)
		if err != nil {
			return err
		}
		updateOpts.Thumbnail = &thumbnail
		updateOpts.Preview = &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.SnapshotID + "/preview.jpg",
			Size:   helper.ToPtr(stat.Size()),
			Image:  &imageProps,
		}
		if err := p.s3.PutFile(updateOpts.Preview.Key, jpegPath, helper.DetectMimeFromFile(jpegPath), updateOpts.Preview.Bucket); err != nil {
			return err
		}
	} else {
		updateOpts.Preview = updateOpts.Original
		thumbnail, err := p.imageProc.Base64Thumbnail(inputPath)
		if err != nil {
			return err
		}
		updateOpts.Thumbnail = &thumbnail
	}
	if err := p.apiClient.UpdateSnapshot(updateOpts); err != nil {
		return err
	}
	return nil
}
