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

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/client/language_client"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type insightsPipeline struct {
	imageProc      *processor.ImageProcessor
	pdfProc        *processor.PDFProcessor
	ocrProc        *processor.OCRProcessor
	fileIdent      *identifier.FileIdentifier
	s3             *infra.S3Manager
	taskClient     *api_client.TaskClient
	snapshotClient *api_client.SnapshotClient
	languageClient *language_client.LanguageClient
}

func NewInsightsPipeline() model.Pipeline {
	return &insightsPipeline{
		imageProc:      processor.NewImageProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		ocrProc:        processor.NewOCRProcessor(),
		fileIdent:      identifier.NewFileIdentifier(),
		s3:             infra.NewS3Manager(),
		taskClient:     api_client.NewTaskClient(),
		snapshotClient: api_client.NewSnapshotClient(),
		languageClient: language_client.NewLanguageClient(),
	}
}

func (p *insightsPipeline) Run(opts api_client.PipelineRunOptions) error {
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
			infra.GetLogger().Error(err)
		}
	}(inputPath)
	return p.RunFromLocalPath(inputPath, opts)
}

func (p *insightsPipeline) RunFromLocalPath(inputPath string, opts api_client.PipelineRunOptions) error {
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Extracting text."),
	}); err != nil {
		return err
	}
	text, err := p.createText(inputPath, opts)
	if err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Collecting entities."),
	}); err != nil {
		return err
	}
	if err := p.createEntities(*text, opts); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName, api_client.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(api_client.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}

func (p *insightsPipeline) createText(inputPath string, opts api_client.PipelineRunOptions) (*string, error) {
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
				infra.GetLogger().Error(err)
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
				infra.GetLogger().Error(err)
			}
		}(pdfPath)
		/* Set OCR S3 object */
		stat, err := os.Stat(pdfPath)
		if err != nil {
			return nil, err
		}
		s3Object := api_client.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.SnapshotID + "/ocr.pdf",
			Size:   helper.ToPtr(stat.Size()),
		}
		if err := p.s3.PutFile(s3Object.Key, pdfPath, helper.DetectMimeFromFile(pdfPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
			return nil, err
		}
		if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{api_client.SnapshotFieldOCR},
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
	s3Object := api_client.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/text.txt",
		Size:   helper.ToPtr(int64(len(*text))),
	}
	if err := p.s3.PutText(s3Object.Key, *text, "text/plain", s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldText},
		Text:    &s3Object,
	}); err != nil {
		return nil, err
	}
	return text, nil
}

func (p *insightsPipeline) createEntities(text string, opts api_client.PipelineRunOptions) error {
	if len(text) == 0 {
		return errors.New("text is empty")
	}
	if len(text) > 1000000 {
		return errors.New("text exceeds supported limit of 1000000 characters")
	}
	res, err := p.languageClient.GetEntities(language_client.GetEntitiesOptions{
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
	s3Object := api_client.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/entities.json",
		Size:   helper.ToPtr(int64(len(content))),
	}
	if err := p.s3.PutText(s3Object.Key, content, "application/json", s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options:  opts,
		Fields:   []string{api_client.SnapshotFieldEntities},
		Entities: &s3Object,
	}); err != nil {
		return err
	}
	return nil
}
