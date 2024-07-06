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
	"voltaserve/client"
	"voltaserve/helper"
	"voltaserve/identifier"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/processor"

	"github.com/minio/minio-go/v7"
)

type mosaicPipeline struct {
	videoProc    *processor.VideoProcessor
	fileIdent    *identifier.FileIdentifier
	s3           *infra.S3Manager
	apiClient    *client.APIClient
	mosaicClient *client.MosaicClient
}

func NewMosaicPipeline() model.Pipeline {
	return &mosaicPipeline{
		videoProc:    processor.NewVideoProcessor(),
		fileIdent:    identifier.NewFileIdentifier(),
		s3:           infra.NewS3Manager(),
		apiClient:    client.NewAPIClient(),
		mosaicClient: client.NewMosaicClient(),
	}
}

func (p *mosaicPipeline) Run(opts client.PipelineRunOptions) error {
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
		Name:   helper.ToPtr("Creating mosaic."),
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

func (p *mosaicPipeline) create(inputPath string, opts client.PipelineRunOptions) error {
	if p.fileIdent.IsImage(opts.Key) {
		if _, err := p.mosaicClient.Create(client.MosaicCreateOptions{
			Path:     inputPath,
			S3Key:    filepath.FromSlash(opts.SnapshotID),
			S3Bucket: opts.Bucket,
		}); err != nil {
			return err
		}
		if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
			Options: opts,
			Fields:  []string{client.SnapshotFieldMosaic},
			Mosaic: &client.S3Object{
				Key:    filepath.FromSlash(opts.SnapshotID + "/mosaic.json"),
				Bucket: opts.Bucket,
			},
		}); err != nil {
			return err
		}
		return nil
	}
	return errors.New("unsupported file type")
}
