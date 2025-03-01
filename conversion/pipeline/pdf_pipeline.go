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

type pdfPipeline struct {
	pdfProc        *processor.PDFProcessor
	imageProc      *processor.ImageProcessor
	s3             infra.S3Manager
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
	fileIdent      *infra.FileIdentifier
	config         *config.Config
}

func NewPDFPipeline() Pipeline {
	return &pdfPipeline{
		pdfProc:        processor.NewPDFProcessor(),
		imageProc:      processor.NewImageProcessor(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		fileIdent:      infra.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts dto.PipelineRunOptions) error {
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

func (p *pdfPipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	count, err := p.pdfProc.CountPages(inputPath)
	if err != nil {
		return err
	}
	document := model.DocumentProps{
		Pages: &model.PagesProps{
			Count:     *count,
			Extension: ".pdf",
		},
	}
	if err := p.patchSnapshotPreviewField(inputPath, &document, opts); err != nil {
		return err
	}
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	// We don't consider failing the creation of the thumbnail an error
	_ = p.createThumbnail(inputPath, opts)
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Saving preview."),
	}); err != nil {
		return err
	}
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Extracting text."),
	}); err != nil {
		return err
	}
	if err := p.extractText(inputPath, opts); err != nil {
		return err
	}
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName, model.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(model.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) createThumbnail(inputPath string, opts dto.PipelineRunOptions) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.pdfProc.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(outputPath)
	props, err := p.imageProc.MeasureImage(outputPath)
	if err != nil {
		return err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	s3Object := &model.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(outputPath),
		Image:  props,
		Size:   stat.Size(),
	}
	if err := p.s3.PutFile(s3Object.Key, outputPath, helper.DetectMimeFromFile(outputPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options:   opts,
		Fields:    []string{model.SnapshotFieldThumbnail},
		Thumbnail: s3Object,
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) extractText(inputPath string, opts dto.PipelineRunOptions) error {
	text, err := p.pdfProc.TextFromPDF(inputPath)
	if err != nil {
		logger.GetLogger().Named(logger.StrPipeline).Errorw(err.Error())
	}
	key := opts.SnapshotID + "/text.txt"
	if text == nil || err != nil {
		return err
	}
	if err := p.s3.PutText(key, *text, "text/plain", opts.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{model.SnapshotFieldText},
		Text: &model.S3Object{
			Bucket: opts.Bucket,
			Key:    key,
			Size:   int64(len(*text)),
		},
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) patchSnapshotPreviewField(inputPath string, document *model.DocumentProps, opts dto.PipelineRunOptions) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if filepath.Ext(inputPath) == filepath.Ext(opts.Key) {
		// The original is a PDF file
		if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{model.SnapshotFieldPreview},
			Preview: &model.S3Object{
				Bucket:   opts.Bucket,
				Key:      opts.Key,
				Size:     stat.Size(),
				Document: document,
			},
		}); err != nil {
			return err
		}
	} else {
		// The original is an office file
		if _, err := p.snapshotClient.Patch(dto.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{model.SnapshotFieldPreview},
			Preview: &model.S3Object{
				Bucket:   opts.Bucket,
				Key:      filepath.FromSlash(opts.SnapshotID + "/preview" + filepath.Ext(inputPath)),
				Size:     stat.Size(),
				Document: document,
			},
		}); err != nil {
			return err
		}
	}
	return nil
}
