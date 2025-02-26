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

type officePipeline struct {
	pdfPipeline    model.Pipeline
	officeProc     *processor.OfficeProcessor
	pdfProc        *processor.PDFProcessor
	s3             apiinfra.S3Manager
	config         *config.Config
	taskClient     *apiclient.TaskClient
	snapshotClient *apiclient.SnapshotClient
}

func NewOfficePipeline() model.Pipeline {
	return &officePipeline{
		pdfPipeline:    NewPDFPipeline(),
		officeProc:     processor.NewOfficeProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		s3:             apiinfra.NewS3Manager(),
		config:         config.GetConfig(),
		taskClient:     apiclient.NewTaskClient(),
		snapshotClient: apiclient.NewSnapshotClient(),
	}
}

func (p *officePipeline) Run(opts model.PipelineRunOptions) error {
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

func (p *officePipeline) RunFromLocalPath(inputPath string, opts model.PipelineRunOptions) error {
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
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

func (p *officePipeline) convertToPDF(inputPath string, opts model.PipelineRunOptions) (*string, error) {
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
	if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{apimodel.SnapshotFieldPreview},
		Preview: &apimodel.S3Object{
			Bucket: opts.Bucket,
			Key:    pdfKey,
			Size:   stat.Size(),
		},
	}); err != nil {
		return nil, err
	}
	return &pdfPath, nil
}
