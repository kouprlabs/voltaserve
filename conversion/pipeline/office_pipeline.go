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

	"github.com/kouprlabs/voltaserve/conversion/client/api_client"
	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type officePipeline struct {
	pdfPipeline    model.Pipeline
	officeProc     *processor.OfficeProcessor
	pdfProc        *processor.PDFProcessor
	s3             *infra.S3Manager
	config         *config.Config
	taskClient     *api_client.TaskClient
	snapshotClient *api_client.SnapshotClient
}

func NewOfficePipeline() model.Pipeline {
	return &officePipeline{
		pdfPipeline:    NewPDFPipeline(),
		officeProc:     processor.NewOfficeProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		s3:             infra.NewS3Manager(),
		config:         config.GetConfig(),
		taskClient:     api_client.NewTaskClient(),
		snapshotClient: api_client.NewSnapshotClient(),
	}
}

func (p *officePipeline) Run(opts api_client.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket, minio.GetObjectOptions{}); err != nil {
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
	return p.RunFromLocalPath(inputPath, opts)
}

func (p *officePipeline) RunFromLocalPath(inputPath string, opts api_client.PipelineRunOptions) error {
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Converting to PDF."),
	}); err != nil {
		return err
	}
	pdfPath, err := p.convertToPDF(inputPath, opts)
	if err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(*pdfPath)
	return p.pdfPipeline.RunFromLocalPath(*pdfPath, opts)
}

func (p *officePipeline) convertToPDF(inputPath string, opts api_client.PipelineRunOptions) (*string, error) {
	outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	outputPath, err := p.officeProc.PDF(inputPath, outputDir)
	if err != nil {
		return nil, err
	}
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			infra.GetLogger().Error(err)
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
	if err := p.s3.PutFile(pdfKey, pdfPath, helper.DetectMimeFromFile(pdfPath), opts.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldPreview},
		Preview: &api_client.S3Object{
			Bucket: opts.Bucket,
			Key:    pdfKey,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return nil, err
	}
	return &pdfPath, nil
}
