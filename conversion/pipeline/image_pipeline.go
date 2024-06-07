package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/processor"
)

type imagePipeline struct {
	imageProc *processor.ImageProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	config    config.Config
}

func NewImagePipeline() model.Pipeline {
	return &imagePipeline{
		imageProc: processor.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		config:    config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts client.PipelineRunOptions) error {
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
	if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Name: helper.ToPtr("Measuring image dimensions."),
	}); err != nil {
		return err
	}
	imageProps, err := p.measureImageDimensions(inputPath, opts)
	if err != nil {
		return err
	}
	var imagePath string
	if filepath.Ext(inputPath) == ".tiff" {
		if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
			Name: helper.ToPtr("Converting TIFF image to JPEG format."),
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
	}
	if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Name: helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	if err := p.createThumbnail(imagePath, opts); err != nil {
		return err
	}
	if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(client.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}

func (p *imagePipeline) measureImageDimensions(inputPath string, opts client.PipelineRunOptions) (*client.ImageProps, error) {
	imageProps, err := p.imageProc.MeasureImage(inputPath)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Original: &client.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   helper.ToPtr(stat.Size()),
			Image:  &imageProps,
		},
	}); err != nil {
		return nil, err
	}
	return &imageProps, nil
}

func (p *imagePipeline) createThumbnail(inputPath string, opts client.PipelineRunOptions) error {
	thumbnail, err := p.imageProc.Base64Thumbnail(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	return nil
}

func (p *imagePipeline) convertTIFFToJPEG(inputPath string, imageProps client.ImageProps, opts client.PipelineRunOptions) (*string, error) {
	jpegPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
	if err := p.imageProc.ConvertImage(inputPath, jpegPath); err != nil {
		return nil, err
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
		return nil, err
	}
	s3Object := &client.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/preview.jpg",
		Size:   helper.ToPtr(stat.Size()),
		Image:  &imageProps,
	}
	if err := p.s3.PutFile(s3Object.Key, jpegPath, helper.DetectMimeFromFile(jpegPath), s3Object.Bucket); err != nil {
		return nil, err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Preview: s3Object,
	}); err != nil {
		return nil, err
	}
	return &jpegPath, nil
}
