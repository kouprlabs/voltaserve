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
	snapshotClient *api_client.SnapshotClient
	taskClient     *api_client.TaskClient
}

func NewOfficePipeline() model.Pipeline {
	return &officePipeline{
		pdfPipeline:    NewPDFPipeline(),
		officeProc:     processor.NewOfficeProcessor(),
		pdfProc:        processor.NewPDFProcessor(),
		s3:             infra.NewS3Manager(),
		config:         config.GetConfig(),
		snapshotClient: api_client.NewSnapshotClient(),
		taskClient:     api_client.NewTaskClient(),
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
	if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
		Fields: []string{api_client.TaskFieldName},
		Name:   helper.ToPtr("Converting to PDF."),
	}); err != nil {
		return err
	}
	pdfKey, err := p.convertToPDF(inputPath, opts)
	if err != nil {
		return err
	}
	if err := p.pdfPipeline.Run(api_client.PipelineRunOptions{
		Bucket:     opts.Bucket,
		Key:        *pdfKey,
		SnapshotID: opts.SnapshotID,
	}); err != nil {
		return err
	}
	return nil
}

func (p *officePipeline) convertToPDF(inputPath string, opts api_client.PipelineRunOptions) (*string, error) {
	outputDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
	outputPath, err := p.officeProc.PDF(inputPath, outputDir)
	if err != nil {
		return nil, err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(*outputPath)
	defer func(path string) {
		if err := os.RemoveAll(path); err != nil {
			infra.GetLogger().Error(err)
		}
	}(outputDir)
	stat, err := os.Stat(*outputPath)
	if err != nil {
		return nil, err
	}
	pdfKey := opts.SnapshotID + "/preview.pdf"
	if err := p.s3.PutFile(pdfKey, *outputPath, helper.DetectMimeFromFile(*outputPath), opts.Bucket, minio.PutObjectOptions{}); err != nil {
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
	return &pdfKey, nil
}
