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

	"github.com/kouprlabs/voltaserve/api/client/apiclient"
	"github.com/kouprlabs/voltaserve/api/client/languageclient"
	apiinfra "github.com/kouprlabs/voltaserve/api/infra"
	apimodel "github.com/kouprlabs/voltaserve/api/model"
	apiservice "github.com/kouprlabs/voltaserve/api/service"

	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type entityPipeline struct {
	imageProc      *processor.ImageProcessor
	pdfProc        *processor.PDFProcessor
	ocrProc        *processor.OCRProcessor
	fileIdent      *apiinfra.FileIdentifier
	s3             apiinfra.S3Manager
	taskClient     *apiclient.TaskClient
	snapshotClient *apiclient.SnapshotClient
	languageClient *languageclient.LanguageClient
}

func NewEntityPipeline() model.Pipeline {
	return &entityPipeline{
		imageProc:      processor.NewImageProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		ocrProc:        processor.NewOCRProcessor(),
		fileIdent:      apiinfra.NewFileIdentifier(),
		s3:             apiinfra.NewS3Manager(),
		taskClient:     apiclient.NewTaskClient(),
		snapshotClient: apiclient.NewSnapshotClient(),
		languageClient: languageclient.NewLanguageClient(),
	}
}

func (p *entityPipeline) Run(opts model.PipelineRunOptions) error {
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

func (p *entityPipeline) RunFromLocalPath(inputPath string, opts model.PipelineRunOptions) error {
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
		Name:   helper.ToPtr("Extracting text."),
	}); err != nil {
		return err
	}
	text, err := p.extractText(inputPath, opts)
	if err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
		Name:   helper.ToPtr("Collecting entities."),
	}); err != nil {
		return err
	}
	if err := p.collectEntities(*text, opts); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName, apimodel.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(apimodel.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}

func (p *entityPipeline) extractText(inputPath string, opts model.PipelineRunOptions) (*string, error) {
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
		s3Object := apimodel.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.SnapshotID + "/ocr.pdf",
			Size:   stat.Size(),
		}
		if err := p.s3.PutFile(s3Object.Key, pdfPath, helper.DetectMimeFromFile(pdfPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
			return nil, err
		}
		if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{apimodel.SnapshotFieldOCR},
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
	s3Object := apimodel.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/text.txt",
		Size:   int64(len(*text)),
	}
	if err := p.s3.PutText(s3Object.Key, *text, "text/plain", s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{apimodel.SnapshotFieldText},
		Text:    &s3Object,
	}); err != nil {
		return nil, err
	}
	return text, nil
}

func (p *entityPipeline) collectEntities(text string, opts model.PipelineRunOptions) error {
	if len(text) == 0 {
		return errors.New("text is empty")
	}
	if len(text) > 1000000 {
		return errors.New("text exceeds supported limit of 1000000 characters")
	}
	res, err := p.languageClient.GetEntities(languageclient.GetEntitiesOptions{
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
	s3Object := apimodel.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/entities.json",
		Size:   int64(len(content)),
	}
	if err := p.s3.PutText(s3Object.Key, content, "application/json", s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
		Options:  opts,
		Fields:   []string{apimodel.SnapshotFieldEntities},
		Entities: &s3Object,
	}); err != nil {
		return err
	}
	return nil
}
