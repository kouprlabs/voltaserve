package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/config"
	"voltaserve/core"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/processor"
)

type pdfPipeline struct {
	pdfProc   *processor.PDFProcessor
	imageProc *processor.ImageProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
	fileIdent *identifier.FileIdentifier
	config    config.Config
}

func NewPDFPipeline() core.Pipeline {
	return &pdfPipeline{
		pdfProc:   processor.NewPDFProcessor(),
		imageProc: processor.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
		fileIdent: identifier.NewFileIdentifier(),
		config:    config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts client.PipelineRunOptions) error {
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

func (p *pdfPipeline) create(inputPath string, opts client.PipelineRunOptions) error {
	thumbnail, err := p.pdfProc.Base64Thumbnail(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options:   opts,
		Thumbnail: &thumbnail,
	}); err != nil {
		return err
	}
	text, err := p.pdfProc.TextFromPDF(inputPath)
	if err != nil {
		infra.GetLogger().Named(infra.StrPipeline).Errorw(err.Error())
	}
	textKey := opts.SnapshotID + "/text.txt"
	if text != "" && err == nil {
		if err := p.s3.PutText(textKey, text, "text/plain", opts.Bucket); err != nil {
			return err
		}
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Preview: &client.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   helper.ToPtr(stat.Size()),
		},
		Text: &client.S3Object{
			Bucket: opts.Bucket,
			Key:    textKey,
			Size:   helper.ToPtr(int64(len(text))),
		},
	}); err != nil {
		return err
	}
	return nil
}
