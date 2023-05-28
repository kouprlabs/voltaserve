package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type videoPipeline struct {
	cmd       *infra.Command
	imageProc *infra.ImageProcessor
	videoProc *infra.VideoProcessor
	s3        *infra.S3Manager
}

func NewVideoPipeline() *videoPipeline {
	return &videoPipeline{
		cmd:       infra.NewCommand(),
		imageProc: infra.NewImageProcessor(),
		videoProc: infra.NewVideoProcessor(),
		s3:        infra.NewS3Manager(),
	}
}

func (p *videoPipeline) Run(opts core.PipelineOptions) (core.PipelineResponse, error) {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return core.PipelineResponse{}, err
	}
	thumbnail, err := p.videoProc.ThumbnailBase64(inputPath)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return core.PipelineResponse{}, err
		}
	}
	return core.PipelineResponse{
		Thumbnail: &thumbnail,
	}, nil
}
