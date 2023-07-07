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
	s3          *infra.S3Manager
	config      config.Config
	apiClient   *client.APIClient
}

func NewOfficePipeline() core.Pipeline {
	return &officePipeline{
		pdfPipeline: NewPDFPipeline(),
		officeProc:  processor.NewOfficeProcessor(),
		s3:          infra.NewS3Manager(),
		config:      config.GetConfig(),
		apiClient:   client.NewAPIClient(),
	}
}

func (p *officePipeline) Run(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	outputPath, err := p.officeProc.PDF(inputPath)
	if err != nil {
		return err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	preview := core.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.FileID + "/" + opts.SnapshotID + "/preview.pdf",
		Size:   stat.Size(),
	}
	if err := p.s3.PutFile(preview.Key, outputPath, helper.DetectMimeFromFile(outputPath), preview.Bucket); err != nil {
		return err
	}
	res := core.PipelineResponse{
		Options: opts,
		Preview: &preview,
	}
	if err := p.apiClient.UpdateSnapshot(&res); err != nil {
		return err
	}
	if err := p.pdfPipeline.Run(core.PipelineOptions{
		Bucket:     preview.Bucket,
		Key:        preview.Key,
		FileID:     opts.FileID,
		SnapshotID: opts.SnapshotID,
	}); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}
