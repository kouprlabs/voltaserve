package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/processor"
)

type officePipeline struct {
	pdfPipeline core.Pipeline
	officeProc  *processor.OfficeProcessor
	pdfProc     *processor.PDFProcessor
	s3          *infra.S3Manager
	config      config.Config
	apiClient   *client.APIClient
}

func NewOfficePipeline() core.Pipeline {
	return &officePipeline{
		pdfPipeline: NewPDFPipeline(),
		officeProc:  processor.NewOfficeProcessor(),
		pdfProc:     processor.NewPDFProcessor(),
		s3:          infra.NewS3Manager(),
		config:      config.GetConfig(),
		apiClient:   client.NewAPIClient(),
	}
}

func (p *officePipeline) Run(opts client.PipelineRunOptions) error {
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

func (p *officePipeline) create(inputPath string, opts client.PipelineRunOptions) error {
	outputPath, err := p.officeProc.PDF(inputPath)
	if err != nil {
		return err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(outputPath)
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	thumbnail, err := p.pdfProc.Base64Thumbnail(outputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	previewKey := opts.SnapshotID + "/preview.pdf"
	if err := p.s3.PutFile(previewKey, outputPath, helper.DetectMimeFromFile(outputPath), opts.Bucket); err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Preview: &client.S3Object{
			Bucket: opts.Bucket,
			Key:    previewKey,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return err
	}
	if err := p.pdfPipeline.Run(client.PipelineRunOptions{
		Bucket:     opts.Bucket,
		Key:        previewKey,
		SnapshotID: opts.SnapshotID,
	}); err != nil {
		return err
	}
	return nil
}
