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

func (p *officePipeline) Run(opts core.PipelineOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewId() + filepath.Ext(opts.Key))
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
	if err := p.s3.PutFile(preview.Key, outputPath, infra.DetectMimeFromFile(outputPath), preview.Bucket); err != nil {
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
