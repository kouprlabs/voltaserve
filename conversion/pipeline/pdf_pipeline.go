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

	"github.com/kouprlabs/voltaserve/api/client/apiclient"
	apiinfra "github.com/kouprlabs/voltaserve/api/infra"
	apimodel "github.com/kouprlabs/voltaserve/api/model"
	apiservice "github.com/kouprlabs/voltaserve/api/service"

	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type pdfPipeline struct {
	pdfProc        *processor.PDFProcessor
	imageProc      *processor.ImageProcessor
	s3             apiinfra.S3Manager
	taskClient     *apiclient.TaskClient
	snapshotClient *apiclient.SnapshotClient
	fileIdent      *apiinfra.FileIdentifier
	config         *config.Config
}

func NewPDFPipeline() model.Pipeline {
	return &pdfPipeline{
		pdfProc:        processor.NewPDFProcessor(),
		imageProc:      processor.NewImageProcessor(),
		s3:             apiinfra.NewS3Manager(),
		taskClient:     apiclient.NewTaskClient(),
		snapshotClient: apiclient.NewSnapshotClient(),
		fileIdent:      apiinfra.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts model.PipelineRunOptions) error {
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

func (p *pdfPipeline) RunFromLocalPath(inputPath string, opts model.PipelineRunOptions) error {
	count, err := p.pdfProc.CountPages(inputPath)
	if err != nil {
		return err
	}
	document := apimodel.DocumentProps{
		Pages: &apimodel.PagesProps{
			Count:     *count,
			Extension: ".pdf",
		},
	}
	if err := p.patchSnapshotPreviewField(inputPath, &document, opts); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
		Name:   helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	// We don't consider failing the creation of the thumbnail an error
	_ = p.createThumbnail(inputPath, opts)
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
		Name:   helper.ToPtr("Saving preview."),
	}); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
		Name:   helper.ToPtr("Extracting text."),
	}); err != nil {
		return err
	}
	if err := p.extractText(inputPath, opts); err != nil {
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

func (p *pdfPipeline) createThumbnail(inputPath string, opts model.PipelineRunOptions) error {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	if err := p.pdfProc.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
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
	s3Object := &apimodel.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(outputPath),
		Image:  props,
		Size:   stat.Size(),
	}
	if err := p.s3.PutFile(s3Object.Key, outputPath, helper.DetectMimeFromFile(outputPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
		Options:   opts,
		Fields:    []string{apimodel.SnapshotFieldThumbnail},
		Thumbnail: s3Object,
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) extractText(inputPath string, opts model.PipelineRunOptions) error {
	text, err := p.pdfProc.TextFromPDF(inputPath)
	if err != nil {
		infra.GetLogger().Named(infra.StrPipeline).Errorw(err.Error())
	}
	key := opts.SnapshotID + "/text.txt"
	if text == nil || err != nil {
		return err
	}
	if err := p.s3.PutText(key, *text, "text/plain", opts.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{apimodel.SnapshotFieldText},
		Text: &apimodel.S3Object{
			Bucket: opts.Bucket,
			Key:    key,
			Size:   int64(len(*text)),
		},
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) patchSnapshotPreviewField(inputPath string, document *apimodel.DocumentProps, opts model.PipelineRunOptions) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if filepath.Ext(inputPath) == filepath.Ext(opts.Key) {
		/* The original is a PDF file */
		if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{apimodel.SnapshotFieldPreview},
			Preview: &apimodel.S3Object{
				Bucket:   opts.Bucket,
				Key:      opts.Key,
				Size:     stat.Size(),
				Document: document,
			},
		}); err != nil {
			return err
		}
	} else {
		/* The original is an office file */
		if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{apimodel.SnapshotFieldPreview},
			Preview: &apimodel.S3Object{
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
