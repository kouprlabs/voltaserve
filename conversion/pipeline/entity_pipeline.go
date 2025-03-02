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

	"github.com/minio/minio-go/v7"

	"github.com/kouprlabs/voltaserve/shared/client"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/conversion/config"
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
	if opts.Language == nil {
		return errors.New("language is undefined")
	}
	return p.RunFromLocalPath("", opts)
}

func (p *entityPipeline) RunFromLocalPath(_ string, opts dto.PipelineRunOptions) error {
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Extracting text."),
	}); err != nil {
		return err
	}
	snapshot, err := p.snapshotClient.Find(opts.SnapshotID)
	if err != nil {
		return err
	}
	if snapshot.Text == nil {
		return nil
	}
	text, err := p.s3.GetText(snapshot.Text.Key, opts.Bucket, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
		Fields: []string{model.TaskFieldName},
		Name:   helper.ToPtr("Collecting entities."),
	}); err != nil {
		return err
	}
	if err := p.patchEntities(text, opts); err != nil {
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

func (p *entityPipeline) patchEntities(text string, opts dto.PipelineRunOptions) error {
	if len(text) == 0 {
		return errors.New("text is empty")
	}
	if len(text) > 1000000 {
		return errors.New("text exceeds supported limit of 1000K characters")
	}
	res, err := p.languageClient.GetEntities(client.GetEntitiesOptions{
		Text:     text,
		Language: *opts.Language,
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
	if _, err := p.snapshotClient.Patch(opts.SnapshotID, dto.SnapshotPatchOptions{
		Fields:   []string{model.SnapshotFieldEntities},
		Entities: &s3Object,
	}); err != nil {
		return err
	}
	return nil
}
