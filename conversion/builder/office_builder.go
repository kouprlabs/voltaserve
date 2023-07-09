package builder

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

type officeBuilder struct {
	pipelineIdentifier *identifier.PipelineIdentifier
	pdfProc            *processor.PDFProcessor
	officeProc         *processor.OfficeProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
}

func NewOfficeBuilder() core.Builder {
	return &officeBuilder{
		pipelineIdentifier: identifier.NewPipelineIdentifier(),
		pdfProc:            processor.NewPDFProcessor(),
		officeProc:         processor.NewOfficeProcessor(),
		s3:                 infra.NewS3Manager(),
		apiClient:          client.NewAPIClient(),
	}
}

func (p *officeBuilder) Build(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	outputPath, err := p.officeProc.PDF(inputPath)
	if err != nil {
		return err
	}
	thumbnail, err := p.pdfProc.Base64Thumbnail(outputPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	if err := p.apiClient.UpdateSnapshot(&core.SnapshotUpdateOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
