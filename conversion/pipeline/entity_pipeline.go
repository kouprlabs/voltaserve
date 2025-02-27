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
	"encoding/json"
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

type entityPipeline struct {
	imageProc      *processor.ImageProcessor
	pdfProc        *processor.PDFProcessor
	ocrProc        *processor.OCRProcessor
	fileIdent      *infra.FileIdentifier
	s3             infra.S3Manager
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
	languageClient *client.LanguageClient
}

func NewEntityPipeline() Pipeline {
	return &entityPipeline{
		imageProc:      processor.NewImageProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		ocrProc:        processor.NewOCRProcessor(),
		fileIdent:      infra.NewFileIdentifier(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		languageClient: client.NewLanguageClient(config.GetConfig().LanguageURL),
	}
}

func (p *entityPipeline) Run(opts dto.PipelineRunOptions) error {
	if opts.Payload == nil || opts.Payload["language"] == "" {
		return errors.New("language is undefined")
	}
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

func (p *entityPipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	if err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Extracting text."),
	}); err != nil {
		return err
	}
	text, err := p.extractText(inputPath, opts)
	if err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Collecting entities."),
	}); err != nil {
		return err
	}
	if err := p.collectEntities(*text, opts); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName, model.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(model.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}

func (p *entityPipeline) extractText(inputPath string, opts dto.PipelineRunOptions) (*string, error) {
	/* Generate PDF/A */
	var pdfPath string
	if p.fileIdent.IsImage(opts.Key) {
		/* Get DPI */
		dpi, err := p.imageProc.DPIFromImage(inputPath)
		if err != nil {
			dpi = helper.ToPtr(72)
		}
		/* Remove alpha channel */
		noAlphaImagePath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
		if err := p.imageProc.RemoveAlphaChannel(inputPath, noAlphaImagePath); err != nil {
			return nil, err
		}
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				logger.GetLogger().Error(err)
			}
		}(noAlphaImagePath)
		/* Convert to PDF/A */
		pdfPath = filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".pdf")
		if err := p.ocrProc.SearchablePDFFromFile(noAlphaImagePath, opts.Payload["language"], *dpi, pdfPath); err != nil {
			return nil, err
		}
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				logger.GetLogger().Error(err)
			}
		}(pdfPath)
		/* Set OCR S3 object */
		stat, err := os.Stat(pdfPath)
		if err != nil {
			return nil, err
		}
		s3Object := model.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.SnapshotID + "/ocr.pdf",
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(s3Object.Key, pdfPath, helper.DetectMimeFromFile(pdfPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
			return nil, err
		}
		if err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{model.SnapshotFieldOCR},
			OCR:     &s3Object,
		}); err != nil {
			return nil, err
		}
	} else if p.fileIdent.IsPDF(opts.Key) || p.fileIdent.IsOffice(opts.Key) || p.fileIdent.IsPlainText(opts.Key) {
		pdfPath = inputPath
	} else {
		return nil, errors.New("unsupported file type")
	}
	/* Extract text */
	text, err := p.pdfProc.TextFromPDF(pdfPath)
	if text == nil || err != nil {
		return nil, err
	}
	/* Set text S3 object */
	s3Object := model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/text.txt",
		Size:   int64(len(*text)),
	}
	if err := p.s3.PutText(s3Object.Key, *text, "text/plain", s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{model.SnapshotFieldText},
		Text:    &s3Object,
	}); err != nil {
		return nil, err
	}
	return text, nil
}

func (p *entityPipeline) collectEntities(text string, opts dto.PipelineRunOptions) error {
	if len(text) == 0 {
		return errors.New("text is empty")
	}
	if len(text) > 1000000 {
		return errors.New("text exceeds supported limit of 1000000 characters")
	}
	res, err := p.languageClient.GetEntities(client.GetEntitiesOptions{
		Text:     text,
		Language: opts.Payload["language"],
	})
	if err != nil {
		return err
	}
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	content := string(b)
	s3Object := model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/entities.json",
		Size:   int64(len(content)),
	}
	if err := p.s3.PutText(s3Object.Key, content, "application/json", s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options:  opts,
		Fields:   []string{model.SnapshotFieldEntities},
		Entities: &s3Object,
	}); err != nil {
		return err
	}
	return nil
}
