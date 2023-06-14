package builder

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type videoBuilder struct {
	pipelineIdentifier *infra.PipelineIdentifier
	videoProc          *infra.VideoProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
}

func NewVideoBuilder() core.Builder {
	return &videoBuilder{
		pipelineIdentifier: infra.NewPipelineIdentifier(),
		videoProc:          infra.NewVideoProcessor(),
		s3:                 infra.NewS3Manager(),
		apiClient:          client.NewAPIClient(),
	}
}

func (p *videoBuilder) Build(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	thumbnail, err := p.videoProc.ThumbnailBase64(inputPath)
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
