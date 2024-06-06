package pipeline

import (
	"errors"
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/processor"
)

type moasicPipeline struct {
	videoProc    *processor.VideoProcessor
	fileIdent    *identifier.FileIdentifier
	s3           *infra.S3Manager
	apiClient    *client.APIClient
	mosaicClient *client.MosaicClient
}

func NewMosaicPipeline() model.Pipeline {
	return &moasicPipeline{
		videoProc:    processor.NewVideoProcessor(),
		fileIdent:    identifier.NewFileIdentifier(),
		s3:           infra.NewS3Manager(),
		apiClient:    client.NewAPIClient(),
		mosaicClient: client.NewMosaicClient(),
	}
}

func (p *moasicPipeline) Run(opts client.PipelineRunOptions) error {
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
		Name: helper.ToPtr("Creating mosaic."),
	}); err != nil {
		return err
	}
	if err := p.create(inputPath, opts); err != nil {
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

func (p *moasicPipeline) create(inputPath string, opts client.PipelineRunOptions) error {
	if p.fileIdent.IsImage(opts.Key) {
		if _, err := p.mosaicClient.Create(client.MosaicCreateOptions{
			Path:     inputPath,
			S3Key:    filepath.FromSlash(opts.SnapshotID),
			S3Bucket: opts.Bucket,
		}); err != nil {
			return err
		}
		if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Mosaic: &client.S3Object{
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
