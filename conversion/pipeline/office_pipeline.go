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

type officePipeline struct {
	pdfPipeline    Pipeline
	officeProc     *processor.OfficeProcessor
	pdfProc        *processor.PDFProcessor
	s3             infra.S3Manager
	config         *config.Config
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
}

func NewOfficePipeline() Pipeline {
	return &officePipeline{
		pdfPipeline:    NewPDFPipeline(),
		officeProc:     processor.NewOfficeProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		config:         config.GetConfig(),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
	}
}

func (p *officePipeline) Run(opts dto.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket, minio.GetObjectOptions{}); err != nil {
		return err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				logger.GetLogger().Error(err)
			}
		}
	}(inputPath)
	return p.RunFromLocalPath(inputPath, opts)
}

func (p *officePipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	if opts.TaskID != nil {
		if _, err := p.taskClient.Patch(*opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName},
			Name:   helper.ToPtr("Converting to PDF."),
		}); err != nil {
			return err
		}
	}
	pdfPath, err := p.patchPreviewWithPDF(inputPath, opts)
	if err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(*pdfPath)
	return p.pdfPipeline.RunFromLocalPath(*pdfPath, opts)
}

func (p *officePipeline) patchPreviewWithPDF(inputPath string, opts dto.PipelineRunOptions) (*string, error) {
	outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	outputPath, err := p.officeProc.PDF(inputPath, outputDir)
	if err != nil {
		return nil, err
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			logger.GetLogger().Error(err)
		}
	}(outputDir)
	pdfPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".pdf")
	if err := os.Rename(*outputPath, pdfPath); err != nil {
		return nil, err
	}
	stat, err := os.Stat(pdfPath)
	if err != nil {
		return nil, err
	}
	pdfKey := opts.SnapshotID + "/preview.pdf"
	if err := p.s3.PutFile(pdfKey, pdfPath, helper.DetectMIMEFromPath(pdfPath), opts.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if _, err := p.snapshotClient.Patch(opts.SnapshotID, dto.SnapshotPatchOptions{
		Fields: []string{model.SnapshotFieldPreview},
		Preview: &model.S3Object{
			Bucket: opts.Bucket,
			Key:    pdfKey,
			Size:   stat.Size(),
		},
	}); err != nil {
		return nil, err
	}
	return &pdfPath, nil
}
