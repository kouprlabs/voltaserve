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

	"go.uber.org/zap"
)

type videoPipeline struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	videoProc          *processor.VideoProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
	logger             *zap.SugaredLogger
}

func NewVideoPipeline() core.Pipeline {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &videoPipeline{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		videoProc:          processor.NewVideoProcessor(),
		s3:                 infra.NewS3Manager(),
		apiClient:          client.NewAPIClient(),
		logger:             logger,
	}
}

func (p *videoPipeline) Run(opts core.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	defer func() {
		_, err := os.Stat(inputPath)
		if os.IsExist(err) {
			if err := os.Remove(inputPath); err != nil {
				p.logger.Error(err)
			}
		}
	}()
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
	return nil
}
