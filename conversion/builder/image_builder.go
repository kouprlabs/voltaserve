package builder

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/processor"
)

type imageBuilder struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	imageProc          *processor.ImageProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
}

func NewImageBuilder() core.Builder {
	return &imageBuilder{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		imageProc:          processor.NewImageProcessor(),
		s3:                 infra.NewS3Manager(),
		apiClient:          client.NewAPIClient(),
	}
}

func (p *imageBuilder) Build(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	thumbnail, err := p.imageProc.Base64Thumbnail(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.UpdateSnapshot(&core.SnapshotUpdateOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
