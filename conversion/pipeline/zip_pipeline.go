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
	"github.com/kouprlabs/voltaserve/conversion/helper"
	"github.com/kouprlabs/voltaserve/conversion/identifier"
	"github.com/kouprlabs/voltaserve/conversion/infra"
	"github.com/kouprlabs/voltaserve/conversion/model"
	"github.com/kouprlabs/voltaserve/conversion/processor"
)

type zipPipeline struct {
	glbPipeline    model.Pipeline
	zipProc        *processor.ZIPProcessor
	gltfProc       *processor.GLTFProcessor
	s3             *infra.S3Manager
	fi             *identifier.FileIdentifier
	taskClient     *api_client.TaskClient
	snapshotClient *api_client.SnapshotClient
}

func NewZIPPipeline() model.Pipeline {
	return &zipPipeline{
		glbPipeline:    NewGLBPipeline(),
		zipProc:        processor.NewZIPProcessor(),
		gltfProc:       processor.NewGLTFProcessor(),
		s3:             infra.NewS3Manager(),
		fi:             identifier.NewFileIdentifier(),
		taskClient:     api_client.NewTaskClient(),
		snapshotClient: api_client.NewSnapshotClient(),
	}
}

func (p *zipPipeline) Run(opts api_client.PipelineRunOptions) error {
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

func (p *zipPipeline) RunFromLocalPath(inputPath string, opts api_client.PipelineRunOptions) error {
	isGLTF, err := p.fi.IsGLTF(inputPath)
	if err != nil {
		return err
	}
	if isGLTF {
		if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
			Fields: []string{api_client.TaskFieldName},
			Name:   helper.ToPtr("Extracting ZIP."),
		}); err != nil {
			return err
		}
		tmpDir := filepath.FromSlash(os.TempDir() + "/" + helper.NewID())
		defer func(path string) {
			if err := os.RemoveAll(path); err != nil {
				infra.GetLogger().Error(err)
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
		if err := p.taskClient.Patch(opts.TaskID, api_client.TaskPatchOptions{
			Fields: []string{api_client.TaskFieldName},
			Name:   helper.ToPtr("Converting to GLB."),
		}); err != nil {
			return err
		}
		glbPath, err := p.convertToGLB(*gltfPath, opts)
		if err != nil {
			return err
		}
		if err := p.glbPipeline.RunFromLocalPath(*glbPath, opts); err != nil {
			return err
		}
	}
	// Do nothing, treat it as a ZIP file
	return nil
}

func (p *zipPipeline) convertToGLB(inputPath string, opts api_client.PipelineRunOptions) (*string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".glb")
	if err := p.gltfProc.ToGLB(inputPath, outputPath); err != nil {
		return nil, err
	}
	defer func(path string) {
		if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
			return
		} else if err != nil {
			infra.GetLogger().Error(err)
		}
	}(outputPath)
	stat, err := os.Stat(outputPath)
	if err != nil {
		return nil, err
	}
	glbKey := opts.SnapshotID + "/preview.glb"
	if err := p.s3.PutFile(glbKey, outputPath, helper.DetectMimeFromFile(outputPath), opts.Bucket, minio.PutObjectOptions{}); err != nil {
		return nil, err
	}
	if err := p.snapshotClient.Patch(api_client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{api_client.SnapshotFieldPreview},
		Preview: &api_client.S3Object{
			Bucket: opts.Bucket,
			Key:    glbKey,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return nil, err
	}
	return &outputPath, nil
}
