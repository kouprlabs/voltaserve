package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/processor"
)

type videoPipeline struct {
	videoProc *processor.VideoProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
}

func NewVideoPipeline() core.Pipeline {
	return &videoPipeline{
		videoProc: processor.NewVideoProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
	}
}

func (p *videoPipeline) Run(opts core.PipelineRunOptions) error {
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

func (p *videoPipeline) create(inputPath string, opts core.PipelineRunOptions) error {
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
