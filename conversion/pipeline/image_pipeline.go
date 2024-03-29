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

	"go.uber.org/zap"
)

type imagePipeline struct {
	imageProc *processor.ImageProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	logger    *zap.SugaredLogger
	config    config.Config
}

func NewImagePipeline() core.Pipeline {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &imagePipeline{
		imageProc: processor.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		logger:    logger,
		config:    config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts core.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	imageProps, err := p.imageProc.MeasureImage(inputPath)
	if err != nil {
		return err
	}
	updateOpts := core.SnapshotUpdateOptions{
		Options: opts,
		Original: &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Image:  &imageProps,
			Size:   stat.Size(),
		},
	}
	if filepath.Ext(inputPath) == ".tiff" {
		jpegPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
		if err := p.imageProc.ConvertImage(inputPath, jpegPath); err != nil {
			return err
		}
		thumbnail, err := p.imageProc.Base64Thumbnail(jpegPath)
		if err != nil {
			return err
		}
		updateOpts.Thumbnail = &thumbnail
		updateOpts.Preview = &core.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.FileID + "/" + opts.SnapshotID + "/preview.jpg",
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(updateOpts.Preview.Key, jpegPath, helper.DetectMimeFromFile(jpegPath), updateOpts.Preview.Bucket); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = jpegPath
	} else {
		thumbnail, err := p.imageProc.Base64Thumbnail(inputPath)
		if err != nil {
			return err
		}
		updateOpts.Thumbnail = &thumbnail
	}
	if err := p.apiClient.UpdateSnapshot(updateOpts); err != nil {
		return err
	}
	if err := os.Remove(inputPath); err != nil {
		return err
	}
	return nil
}
