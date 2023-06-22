package builder

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type officeBuilder struct {
	pipelineIdentifier *infra.PipelineIdentifier
	pdfProc            *infra.PDFProcessor
	officeProc         *infra.OfficeProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
}

func NewOfficeBuilder() core.Builder {
	return &officeBuilder{
		pipelineIdentifier: infra.NewPipelineIdentifier(),
		pdfProc:            infra.NewPDFProcessor(),
		officeProc:         infra.NewOfficeProcessor(),
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
	thumbnail, err := p.pdfProc.ThumbnailBase64(outputPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	res := core.PipelineResponse{
		Options:   opts,
		Thumbnail: &thumbnail,
	}
	if err := p.apiClient.UpdateSnapshot(&res); err != nil {
		return err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}
