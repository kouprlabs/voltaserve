package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/processor"
)

type officePipeline struct {
	pdfPipeline model.Pipeline
	officeProc  *processor.OfficeProcessor
	pdfProc     *processor.PDFProcessor
	s3          *infra.S3Manager
	config      config.Config
	apiClient   *client.APIClient
}

func NewOfficePipeline() model.Pipeline {
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
	if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Name: helper.ToPtr("Converting to PDF."),
	}); err != nil {
		return err
	}
	pdfKey, err := p.convertToPDF(inputPath, opts)
	if err != nil {
		return err
	}
	if err := p.pdfPipeline.Run(client.PipelineRunOptions{
		Bucket:     opts.Bucket,
		Key:        *pdfKey,
		SnapshotID: opts.SnapshotID,
	}); err != nil {
		return err
	}
	return nil
}

func (p *officePipeline) convertToPDF(inputPath string, opts client.PipelineRunOptions) (*string, error) {
	outputPath, err := p.officeProc.PDF(inputPath)
	if err != nil {
		return nil, err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(*outputPath)
	stat, err := os.Stat(*outputPath)
	if err != nil {
		return nil, err
	}
	pdfKey := opts.SnapshotID + "/preview.pdf"
	if err := p.s3.PutFile(pdfKey, *outputPath, helper.DetectMimeFromFile(*outputPath), opts.Bucket); err != nil {
		return nil, err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Preview: &client.S3Object{
			Bucket: opts.Bucket,
			Key:    pdfKey,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return nil, err
	}
	return &pdfKey, nil
}
