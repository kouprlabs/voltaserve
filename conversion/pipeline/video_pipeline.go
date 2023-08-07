package pipeline

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

type videoPipeline struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	videoProc          *processor.VideoProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
}

func NewVideoPipeline() core.Pipeline {
	return &videoPipeline{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		videoProc:          processor.NewVideoProcessor(),
		s3:                 infra.NewS3Manager(),
		apiClient:          client.NewAPIClient(),
	}
}

func (p *videoPipeline) Run(opts core.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	thumbnail, err := p.videoProc.Base64Thumbnail(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	if err := os.Remove(inputPath); err != nil {
		return err
	}
	return nil
}
