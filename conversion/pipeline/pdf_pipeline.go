// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package pipeline

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type pdfPipeline struct {
	pdfProc        *processor.PDFProcessor
	imageProc      *processor.ImageProcessor
	s3             *infra.S3Manager
	taskClient     *api_client.TaskClient
	snapshotClient *api_client.SnapshotClient
	fileIdent      *identifier.FileIdentifier
	config         *config.Config
}

func NewPDFPipeline() model.Pipeline {
	return &pdfPipeline{
		pdfProc:        processor.NewPDFProcessor(),
		imageProc:      processor.NewImageProcessor(),
		s3:             infra.NewS3Manager(),
		taskClient:     api_client.NewTaskClient(),
		snapshotClient: api_client.NewSnapshotClient(),
		fileIdent:      identifier.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

func (p *pdfPipeline) Run(opts api_client.PipelineRunOptions) error {
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
	count, err := p.pdfProc.CountPages(inputPath)
	if err != nil {
		return err
	}
	document := api_client.DocumentProps{
		Pages: &api_client.PagesProps{
			Count:     *count,
			Extension: filepath.Ext(opts.Key),
		},
	}
	if err := p.updateSnapshot(inputPath, &document, opts); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	if err := p.createThumbnail(inputPath, opts); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Saving preview."),
	}); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Extracting text."),
	}); err != nil {
		return err
	}
	if err := p.extractText(inputPath, opts); err != nil {
		return err
	}
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Performing segmentation."),
	}); err != nil {
		return err
	}
	if err := p.performSegmentation(inputPath, &document, opts); err != nil {
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

func (p *pdfPipeline) createThumbnail(inputPath string, opts api_client.PipelineRunOptions) error {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".png")
	// We don't consider failing the creation of the thumbnail as an error
	_ = p.pdfProc.Thumbnail(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, tmpPath)
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(tmpPath)
	props, err := p.imageProc.MeasureImage(tmpPath)
	if err != nil {
		return err
	}
	stat, err := os.Stat(tmpPath)
	if err != nil {
		return err
	}
	s3Object := &api_client.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(tmpPath),
		Image:  props,
		Size:   helper.ToPtr(stat.Size()),
	}
	if err := p.s3.PutFile(s3Object.Key, tmpPath, helper.DetectMimeFromFile(tmpPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options:   opts,
		Fields:    []string{api_client.SnapshotFieldThumbnail},
		Thumbnail: s3Object,
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) extractText(inputPath string, opts api_client.PipelineRunOptions) error {
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
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldText},
		Text: &api_client.S3Object{
			Bucket: opts.Bucket,
			Key:    key,
			Size:   helper.ToPtr(int64(len(*text))),
		},
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) updateSnapshot(inputPath string, document *api_client.DocumentProps, opts api_client.PipelineRunOptions) error {
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	s3Object := &api_client.S3Object{
		Bucket:   opts.Bucket,
		Key:      opts.Key,
		Size:     helper.ToPtr(stat.Size()),
		Document: document,
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields: []string{
			api_client.SnapshotFieldOriginal,
			api_client.SnapshotFieldPreview,
		},
		Original: s3Object,
		Preview:  s3Object,
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) performSegmentation(inputPath string, document *api_client.DocumentProps, opts api_client.PipelineRunOptions) error {
	if err := p.splitPages(inputPath, opts); err != nil {
		return err
	}
	if err := p.splitThumbnails(inputPath, opts); err != nil {
		return err
	}
	document.Thumbnails = &api_client.ThumbnailsProps{
		Extension: ".jpg",
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldSegmentation},
		Segmentation: &api_client.S3Object{
			Bucket:   opts.Bucket,
			Key:      filepath.FromSlash(opts.SnapshotID + "/segmentation"),
			Document: document,
		},
	}); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) splitPages(inputPath string, opts api_client.PipelineRunOptions) error {
	pagesDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	if err := os.MkdirAll(pagesDir, 0o750); err != nil {
		return nil
	}
	defer func(path string) {
		if err := os.RemoveAll(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(pagesDir)
	if err := p.pdfProc.SplitPages(inputPath, pagesDir); err != nil {
		return err
	}
	if err := p.s3.PutFolder(filepath.FromSlash(opts.SnapshotID+"/segmentation/pages"), pagesDir, opts.Bucket); err != nil {
		return err
	}
	return nil
}

func (p *pdfPipeline) splitThumbnails(inputPath string, opts api_client.PipelineRunOptions) error {
	thumbnailsDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	if err := os.MkdirAll(thumbnailsDir, 0o750); err != nil {
		return nil
	}
	defer func(path string) {
		if err := os.RemoveAll(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(thumbnailsDir)
	if err := p.pdfProc.SplitThumbnails(inputPath, 100, 0, ".jpg", thumbnailsDir); err != nil {
		return err
	}
	if err := p.s3.PutFolder(filepath.FromSlash(opts.SnapshotID+"/segmentation/thumbnails"), thumbnailsDir, opts.Bucket); err != nil {
		return err
	}
	return nil
}
