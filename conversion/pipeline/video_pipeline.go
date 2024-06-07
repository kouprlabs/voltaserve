package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/processor"
)

type videoPipeline struct {
	videoProc *processor.VideoProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
}

func NewVideoPipeline() model.Pipeline {
	return &videoPipeline{
		videoProc: processor.NewVideoProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
	}
}

func (p *videoPipeline) Run(opts client.PipelineRunOptions) error {
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
		Fields: []string{client.TaskFieldName},
		Name:   helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	if err := p.createThumbnail(inputPath, opts); err != nil {
		return err
	}
	if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Fields: []string{client.TaskFieldName, client.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(client.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}

func (p *videoPipeline) createThumbnail(inputPath string, opts client.PipelineRunOptions) error {
	thumbnail, err := p.videoProc.Base64Thumbnail(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options:   opts,
		Fields:    []string{client.SnapshotFieldThumbnail},
		Thumbnail: thumbnail,
	}); err != nil {
		return err
	}
	return nil
}
