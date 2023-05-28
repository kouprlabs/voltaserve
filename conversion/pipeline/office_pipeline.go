package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/infra"
)

type officePipeline struct {
	pdfPipeline core.Pipeline
	officeProc  *infra.OfficeProcessor
	s3          *infra.S3Manager
	config      config.Config
}

func NewOfficePipeline() core.Pipeline {
	return &officePipeline{
		pdfPipeline: NewPDFPipeline(),
		officeProc:  infra.NewOfficeProcessor(),
		s3:          infra.NewS3Manager(),
		config:      config.GetConfig(),
	}
}

func (p *officePipeline) Run(opts core.PipelineOptions) (core.PipelineResponse, error) {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket); err != nil {
		return core.PipelineResponse{}, err
	}
	outputPath, err := p.officeProc.PDF(inputPath)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return core.PipelineResponse{}, err
	}
	s3Object := core.S3Object{
		Bucket: opts.Bucket,
		Key:    filepath.FromSlash(opts.FileID + "/" + opts.SnapshotID + "/preview.pdf"),
		Size:   stat.Size(),
	}
	if err := p.s3.PutFile(s3Object.Key, outputPath, infra.DetectMimeFromFile(outputPath), s3Object.Bucket); err != nil {
		return core.PipelineResponse{}, err
	}
	res, err := p.pdfPipeline.Run(core.PipelineOptions{
		Bucket:     s3Object.Bucket,
		Key:        s3Object.Key,
		FileID:     opts.FileID,
		SnapshotID: opts.SnapshotID,
	})
	if err != nil {
		return core.PipelineResponse{}, err
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return core.PipelineResponse{}, err
		}
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return core.PipelineResponse{}, err
		}
	}
	return core.PipelineResponse{
		Preview:   &s3Object,
		Thumbnail: res.Thumbnail,
		OCR:       res.OCR,
		Text:      res.Text,
	}, nil
}
