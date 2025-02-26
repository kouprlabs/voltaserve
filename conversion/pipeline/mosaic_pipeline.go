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
	"github.com/kouprlabs/voltaserve/api/client/mosaicclient"
	apiinfra "github.com/kouprlabs/voltaserve/api/infra"
	apimodel "github.com/kouprlabs/voltaserve/api/model"
	apiservice "github.com/kouprlabs/voltaserve/api/service"

	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type mosaicPipeline struct {
	videoProc      *processor.VideoProcessor
	imageProc      *processor.ImageProcessor
	fileIdent      *apiinfra.FileIdentifier
	s3             apiinfra.S3Manager
	taskClient     *apiclient.TaskClient
	snapshotClient *apiclient.SnapshotClient
	mosaicClient   *mosaicclient.MosaicClient
}

func NewMosaicPipeline() model.Pipeline {
	return &mosaicPipeline{
		videoProc:      processor.NewVideoProcessor(),
		imageProc:      processor.NewImageProcessor(),
		fileIdent:      apiinfra.NewFileIdentifier(),
		s3:             apiinfra.NewS3Manager(),
		taskClient:     apiclient.NewTaskClient(),
		snapshotClient: apiclient.NewSnapshotClient(),
		mosaicClient:   mosaicclient.NewMosaicClient(),
	}
}

func (p *mosaicPipeline) Run(opts model.PipelineRunOptions) error {
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

func (p *mosaicPipeline) RunFromLocalPath(inputPath string, opts model.PipelineRunOptions) error {
	if !p.fileIdent.IsImage(opts.Key) {
		return errors.New("unsupported file type")
	}
	if err := p.taskClient.Patch(opts.TaskID, apiservice.TaskPatchOptions{
		Fields: []string{apimodel.TaskFieldName},
		Name:   helper.ToPtr("Creating mosaic."),
	}); err != nil {
		return err
	}
	if !p.imageProc.IsSupportedByBild(inputPath) {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpg")
		if err := p.imageProc.ConvertImage(inputPath, outputPath); err != nil {
			return err
		}
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				infra.GetLogger().Error(err)
			}
		}(outputPath)
		inputPath = outputPath
	}
	if _, err := p.mosaicClient.Create(mosaicclient.MosaicCreateOptions{
		Path:     inputPath,
		S3Key:    filepath.FromSlash(opts.SnapshotID),
		S3Bucket: opts.Bucket,
	}); err != nil {
		return err
	}
	if err := p.snapshotClient.Patch(apiservice.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{apimodel.SnapshotFieldMosaic},
		Mosaic: &apimodel.S3Object{
			Key:    filepath.FromSlash(opts.SnapshotID + "/mosaic"),
			Bucket: opts.Bucket,
		},
	}); err != nil {
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
