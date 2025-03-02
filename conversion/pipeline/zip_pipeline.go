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

type zipPipeline struct {
	glbPipeline    Pipeline
	zipProc        *processor.ZIPProcessor
	gltfProc       *processor.GLTFProcessor
	s3             infra.S3Manager
	fileIdent      *infra.FileIdentifier
	taskClient     *client.TaskClient
	snapshotClient *client.SnapshotClient
}

func NewZIPPipeline() Pipeline {
	return &zipPipeline{
		glbPipeline:    NewGLBPipeline(),
		zipProc:        processor.NewZIPProcessor(),
		gltfProc:       processor.NewGLTFProcessor(),
		s3:             infra.NewS3Manager(config.GetConfig().S3, config.GetConfig().Environment),
		fileIdent:      infra.NewFileIdentifier(),
		taskClient:     client.NewTaskClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
		snapshotClient: client.NewSnapshotClient(config.GetConfig().APIURL, config.GetConfig().Security.APIKey),
	}
}

func (p *zipPipeline) Run(opts dto.PipelineRunOptions) error {
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(opts.Key))
	if err := p.s3.GetFile(opts.Key, inputPath, opts.Bucket, minio.GetObjectOptions{}); err != nil {
		return err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			logger.GetLogger().Error(err)
		}
	}(inputPath)
	return p.RunFromLocalPath(inputPath, opts)
}

func (p *zipPipeline) RunFromLocalPath(inputPath string, opts dto.PipelineRunOptions) error {
	isGLTF, err := p.fileIdent.IsGLTF(inputPath)
	if err != nil {
		return err
	}
	if isGLTF {
		if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName},
			Name:   helper.ToPtr("Extracting ZIP."),
		}); err != nil {
			return err
		}
		tmpDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
		defer func(path string) {
			if err := os.RemoveAll(path); err != nil {
				logger.GetLogger().Error(err)
			}
		}(tmpDir)
		if err := p.zipProc.Extract(inputPath, tmpDir); err != nil {
			return err
		}
		gltfPath, err := helper.FindFileWithExtension(tmpDir, ".gltf")
		if err != nil {
			return err
		}
		if gltfPath == nil {
			// Do nothing, treat it as a ZIP file
			return nil
		}
		if _, err := p.taskClient.Patch(opts.TaskID, dto.TaskPatchOptions{
			Fields: []string{model.TaskFieldName},
			Name:   helper.ToPtr("Converting to GLB."),
		}); err != nil {
			return err
		}
		glbPath, err := p.convertToGLB(*gltfPath, opts)
		if err != nil {
			return err
		}
		defer func(path string) {
			if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
				return
			} else if err != nil {
				logger.GetLogger().Error(err)
			}
		}(*glbPath)
		if err := p.glbPipeline.RunFromLocalPath(*glbPath, opts); err != nil {
			return err
		}
	}
	// Do nothing, treat it as a ZIP file
	return nil
}

func (p *zipPipeline) convertToGLB(inputPath string, opts dto.PipelineRunOptions) (*string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".glb")
	if err := p.gltfProc.ToGLB(inputPath, outputPath); err != nil {
		return nil, err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return nil, err
	}
	glbKey := opts.SnapshotID + "/preview.glb"
	if err := p.s3.PutFile(glbKey, outputPath, helper.DetectMIMEFromPath(outputPath), opts.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if _, err := p.snapshotClient.Patch(opts.SnapshotID, dto.SnapshotPatchOptions{
		Fields: []string{model.SnapshotFieldPreview},
		Preview: &model.S3Object{
			Bucket: opts.Bucket,
			Key:    glbKey,
			Size:   stat.Size(),
		},
	}); err != nil {
		return nil, err
	}
	return &outputPath, nil
}
