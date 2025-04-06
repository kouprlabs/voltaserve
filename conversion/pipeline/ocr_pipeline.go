// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package pipeline

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/logger"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type ocrPipeline struct {
	imageProc      *processor.ImageProcessor
	ocrProc        *processor.OCRProcessor
	pdfProc        *processor.PDFProcessor
	s3             infra.S3Manager
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
	fileIdent      *infra.FileIdentifier
	config         *config.Config
}

func NewOCRPipeline() Pipeline {
	return &ocrPipeline{
		imageProc:      processor.NewImageProcessor(),
		ocrProc:        processor.NewOCRProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		fileIdent:      infra.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

func (p *ocrPipeline) Run(opts dto.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket, minio.GetObjectOptions{}); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(inputPath)
	return p.RunFromLocalPath(inputPath, opts)
}

func (p *ocrPipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	if opts.TaskID != nil {
		if _, err := p.taskClient.Patch(*opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName},
			Name:   helper.ToPtr("Extracting text."),
		}); err != nil {
			return err
		}
	}
	if opts.Intent != nil && *opts.Intent == model.SnapshotIntentDocument && opts.Language != nil && *opts.Language != "" {
		_ = p.patchText(inputPath, opts)
	}
	if opts.TaskID != nil {
		if _, err := p.taskClient.Patch(*opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName, model.TaskFieldStatus},
			Name:   helper.ToPtr("Done."),
			Status: helper.ToPtr(model.TaskStatusSuccess),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p *ocrPipeline) patchText(inputPath string, opts dto.PipelineRunOptions) error {
	// Generate PDF/A
	var pdfPath string
	// Get DPI
	dpi, err := p.imageProc.DPIFromImage(inputPath)
	if err != nil {
		dpi = helper.ToPtr(72)
	}
	// Remove alpha channel
	noAlphaImagePath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.imageProc.RemoveAlphaChannel(inputPath, noAlphaImagePath); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(noAlphaImagePath)
	// Convert to PDF/A
	pdfPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".pdf")
	if err := p.ocrProc.SearchablePDFFromFile(noAlphaImagePath, *opts.Language, *dpi, pdfPath); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(pdfPath)
	// Set OCR S3 object
	stat, err := os.Stat(pdfPath)
	if err != nil {
		return err
	}
	count, err := p.pdfProc.CountPages(pdfPath)
	if err != nil {
		return err
	}
	s3Object := model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/ocr.pdf",
		Size:   stat.Size(),
		Document: &model.DocumentProps{
			Page: &model.PageProps{
				Count:     *count,
				Extension: ".pdf",
			},
		},
	}
	if err := p.s3.PutFile(s3Object.Key, pdfPath, helper.DetectMIMEFromPath(pdfPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(opts.SnapshotID, dto.SnapshotPatchOptions{
		Fields: []string{model.SnapshotFieldOCR},
		OCR:    &s3Object,
	}); err != nil {
		return err
	}
	// Extract text
	text, err := p.pdfProc.TextFromPDF(pdfPath)
	if err != nil {
		return err
	}
	if text == nil || len(*text) == 0 {
		return nil
	}
	// Set text S3 object
	s3Object = model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/text.txt",
		Size:   int64(len(*text)),
	}
	if err := p.s3.PutText(s3Object.Key, *text, "text/plain", s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(opts.SnapshotID, dto.SnapshotPatchOptions{
		Fields: []string{model.SnapshotFieldText},
		Text:   &s3Object,
	}); err != nil {
		return err
	}
	return nil
}
