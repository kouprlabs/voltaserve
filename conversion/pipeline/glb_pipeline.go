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
	"os"
	"path/filepath"

	"github.com/kouprlabs/voltaserve/conversion/client"
	"github.com/kouprlabs/voltaserve/conversion/config"
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"

	"github.com/minio/minio-go/v7"
)

type glbPipeline struct {
	glbProc   *processor.GLBProcessor
	imageProc *processor.ImageProcessor
	s3        *infra.S3Manager
	apiClient *client.APIClient
	config    *config.Config
}

func NewGLBPipeline() model.Pipeline {
	return &glbPipeline{
		glbProc:   processor.NewGLBProcessor(),
		imageProc: processor.NewImageProcessor(),
		s3:        infra.NewS3Manager(),
		apiClient: client.NewAPIClient(),
		config:    config.GetConfig(),
	}
}

func (p *glbPipeline) Run(opts client.PipelineRunOptions) error {
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
		Name:   helper.ToPtr("Creating thumbnail."),
	}); err != nil {
		return err
	}
	if err := p.createThumbnail(inputPath, opts); err != nil {
		return err
	}
	stat, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{client.SnapshotFieldPreview},
		Preview: &client.S3Object{
			Bucket: opts.Bucket,
			Key:    opts.Key,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return err
	}
	return nil
}

func (p *glbPipeline) createThumbnail(inputPath string, opts client.PipelineRunOptions) error {
	tmpPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".jpeg")
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
		}
	}(tmpPath)
	if err := p.glbProc.Thumbnail(inputPath, p.config.Limits.ImagePreviewMaxWidth, p.config.Limits.ImagePreviewMaxHeight, "rgb(255,255,255)", tmpPath); err != nil {
		return err
	}
	props, err := p.imageProc.MeasureImage(tmpPath)
	if err != nil {
		return err
	}
	stat, err := os.Stat(tmpPath)
	if err != nil {
		return err
	}
	s3Object := &client.S3Object{
		Bucket: opts.Bucket,
		Key:    opts.SnapshotID + "/thumbnail" + filepath.Ext(tmpPath),
		Image:  props,
		Size:   helper.ToPtr(stat.Size()),
	}
	if err := p.s3.PutFile(s3Object.Key, tmpPath, helper.DetectMimeFromFile(tmpPath), s3Object.Bucket, minio.PutObjectOptions{}); err != nil {
		return err
	}
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options:   opts,
		Fields:    []string{client.SnapshotFieldThumbnail},
		Thumbnail: s3Object,
	}); err != nil {
		return err
	}
	return nil
}
