package builder

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type imageBuilder struct {
	pipelineIdentifier *infra.PipelineIdentifier
	imageProc          *infra.ImageProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
}

func NewImageBuilder() core.Builder {
	return &imageBuilder{
		pipelineIdentifier: infra.NewPipelineIdentifier(),
		imageProc:          infra.NewImageProcessor(),
		s3:                 infra.NewS3Manager(),
		apiClient:          client.NewAPIClient(),
	}
}

func (p *imageBuilder) Build(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	thumbnail, err := p.imageProc.ThumbnailBase64(inputPath)
	if err != nil {
		return err
	}
	res := core.PipelineResponse{
		Options:   opts,
		Thumbnail: &thumbnail,
	}
	if err := p.apiClient.UpdateSnapshot(&res); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
