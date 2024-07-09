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

	"github.com/kouprlabs/voltaserve/conversion/client"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"

	"github.com/minio/minio-go/v7"
)

type watermarkPipeline struct {
	videoProc       *processor.VideoProcessor
	fileIdent       *identifier.FileIdentifier
	s3              *infra.S3Manager
	apiClient       *client.APIClient
	watermarkClient *client.WatermarkClient
}

func NewWatermarkPipeline() model.Pipeline {
	return &watermarkPipeline{
		videoProc:       processor.NewVideoProcessor(),
		fileIdent:       identifier.NewFileIdentifier(),
		s3:              infra.NewS3Manager(),
		apiClient:       client.NewAPIClient(),
		watermarkClient: client.NewWatermarkClient(),
	}
}

func (p *watermarkPipeline) Run(opts client.PipelineRunOptions) error {
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
	if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Fields: []string{client.TaskFieldName},
		Name:   helper.ToPtr("Applying watermark."),
	}); err != nil {
		return err
	}
	if err := p.create(inputPath, opts); err != nil {
		return err
	}
	if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
		Fields: []string{client.TaskFieldName, client.TaskFieldStatus},
		Name:   helper.ToPtr("Done."),
		Status: helper.ToPtr(client.TaskStatusSuccess),
	}); err != nil {
		return err
	}
	return nil
}

func (p *watermarkPipeline) create(inputPath string, opts client.PipelineRunOptions) error {
	var category string
	if p.fileIdent.IsImage(opts.Key) {
		category = "image"
	} else if p.fileIdent.IsPDF(opts.Key) {
		category = "document"
	} else if p.fileIdent.IsOffice(opts.Key) || p.fileIdent.IsPlainText(opts.Key) {
		category = "document"
	} else {
		return errors.New("unsupported file type")
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	key := filepath.FromSlash(opts.SnapshotID + "/watermark" + filepath.Ext(opts.Key))
	if err := p.watermarkClient.Create(client.WatermarkCreateOptions{
		Path:     inputPath,
		S3Key:    key,
		S3Bucket: opts.Bucket,
		Category: category,
		Values: []string{
			opts.Payload["workspace"],
			opts.Payload["user"],
		},
	}); err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{client.SnapshotFieldWatermark},
		Watermark: &client.S3Object{
			Key:    key,
			Bucket: opts.Bucket,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return err
	}
	return nil
}
