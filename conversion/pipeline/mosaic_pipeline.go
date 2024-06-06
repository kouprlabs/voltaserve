package pipeline

import (
	"errors"
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/processor"
)

type moasicPipeline struct {
	videoProc    *processor.VideoProcessor
	fileIdent    *identifier.FileIdentifier
	s3           *infra.S3Manager
	apiClient    *client.APIClient
	mosaicClient *client.MosaicClient
}

func NewMosaicPipeline() core.Pipeline {
	return &moasicPipeline{
		videoProc:    processor.NewVideoProcessor(),
		fileIdent:    identifier.NewFileIdentifier(),
		s3:           infra.NewS3Manager(),
		apiClient:    client.NewAPIClient(),
		mosaicClient: client.NewMosaicClient(),
	}
}

func (p *moasicPipeline) Run(opts core.PipelineRunOptions) error {
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

func (p *moasicPipeline) create(inputPath string, opts core.PipelineRunOptions) error {
	if p.fileIdent.IsImage(opts.Key) {
		if _, err := p.mosaicClient.Create(client.MosaicCreateOptions{
			Path:     inputPath,
			S3Key:    filepath.FromSlash(opts.SnapshotID),
			S3Bucket: opts.Bucket,
		}); err != nil {
			return err
		}
		if err := p.apiClient.UpdateSnapshot(core.SnapshotUpdateOptions{
			Options: opts,
			Mosaic: &core.S3Object{
				Key:    filepath.FromSlash(opts.SnapshotID + "/mosaic.json"),
				Bucket: opts.Bucket,
			},
		}); err != nil {
			return err
		}
		return nil
	}
	return errors.New("unsupported file type")
}
