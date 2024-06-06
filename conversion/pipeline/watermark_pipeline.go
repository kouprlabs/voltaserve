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

type watermarkPipeline struct {
	videoProc       *processor.VideoProcessor
	fileIdent       *identifier.FileIdentifier
	s3              *infra.S3Manager
	apiClient       *client.APIClient
	watermarkClient *client.WatermarkClient
}

func NewWatermarkPipeline() core.Pipeline {
	return &watermarkPipeline{
		videoProc:       processor.NewVideoProcessor(),
		fileIdent:       identifier.NewFileIdentifier(),
		s3:              infra.NewS3Manager(),
		apiClient:       client.NewAPIClient(),
		watermarkClient: client.NewWatermarkClient(),
	}
}

func (p *watermarkPipeline) Run(opts client.PipelineRunOptions) error {
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

func (p *watermarkPipeline) create(inputPath string, opts client.PipelineRunOptions) error {
	var category string
	if p.fileIdent.IsImage(opts.Key) {
		category = "image"
	} else if p.fileIdent.IsPDF(opts.Key) {
		category = "document"
	} else if p.fileIdent.IsOffice(opts.Key) || p.fileIdent.IsPlainText(opts.Key) {
		category = "document"
	} else {
		return errors.New("unsupported file type")
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	key := filepath.FromSlash(opts.SnapshotID + "/watermark" + filepath.Ext(opts.Key))
	if err := p.watermarkClient.Create(client.WatermarkCreateOptions{
		Path:     inputPath,
		S3Key:    key,
		S3Bucket: opts.Bucket,
		Category: category,
		Values:   opts.Values,
	}); err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Watermark: &client.S3Object{
			Key:    key,
			Bucket: opts.Bucket,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return err
	}
	return nil
}
