package builder

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type thumbnailBuilder struct {
	pipelineIdentifier *infra.PipelineIdentifier
	imageProc          *infra.ImageProcessor
	pdfProc            *infra.PDFProcessor
	videoProc          *infra.VideoProcessor
	s3                 *infra.S3Manager
	apiClient          *client.APIClient
}

func NewThunmbnailBuilder() core.Builder {
	return &thumbnailBuilder{
		pipelineIdentifier: infra.NewPipelineIdentifier(),
		imageProc:          infra.NewImageProcessor(),
		pdfProc:            infra.NewPDFProcessor(),
		videoProc:          infra.NewVideoProcessor(),
		s3:                 infra.NewS3Manager(),
		apiClient:          client.NewAPIClient(),
	}
}

func (p *thumbnailBuilder) Build(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return err
	}
	var thumbnail core.Thumbnail
	var err error
	pipeline := p.pipelineIdentifier.Identify(opts)
	if pipeline == core.PipelinePDF {
		thumbnail, err = p.pdfProc.ThumbnailBase64(inputPath)
		if err != nil {
			return err
		}
	} else if pipeline == core.PipelineImage {
		thumbnail, err = p.imageProc.ThumbnailBase64(inputPath)
		if err != nil {
			return err
		}
	} else if pipeline == core.PipelineVideo {
		thumbnail, err = p.videoProc.ThumbnailBase64(inputPath)
		if err != nil {
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
