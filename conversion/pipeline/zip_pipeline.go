package pipeline

import (
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

type zipPipeline struct {
	glbPipeline model.Pipeline
	zipProc     *processor.ZIPProcessor
	gltfProc    *processor.GLTFProcessor
	s3          *infra.S3Manager
	fi          *identifier.FileIdentifier
	apiClient   *client.APIClient
}

func NewZIPPipeline() model.Pipeline {
	return &zipPipeline{
		glbPipeline: NewGLBPipeline(),
		zipProc:     processor.NewZIPProcessor(),
		gltfProc:    processor.NewGLTFProcessor(),
		s3:          infra.NewS3Manager(),
		fi:          identifier.NewFileIdentifier(),
		apiClient:   client.NewAPIClient(),
	}
}

func (p *zipPipeline) Run(opts client.PipelineRunOptions) error {
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
	isGLTF, err := p.fi.IsGLTF(inputPath)
	if err != nil {
		return err
	}
	if isGLTF {
		if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
			Fields: []string{client.TaskFieldName},
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
		if err := p.apiClient.PatchTask(opts.TaskID, client.TaskPatchOptions{
			Fields: []string{client.TaskFieldName},
			Name:   helper.ToPtr("Converting to GLB."),
		}); err != nil {
			return err
		}
		glbKey, err := p.convertToGLB(*gltfPath, opts)
		if err != nil {
			return err
		}
		if err := p.glbPipeline.Run(client.PipelineRunOptions{
			Bucket:     opts.Bucket,
			Key:        *glbKey,
			SnapshotID: opts.SnapshotID,
		}); err != nil {
			return err
		}
	}
	// Do nothing, treat it as a ZIP file
	return nil
}

func (p *zipPipeline) convertToGLB(inputPath string, opts client.PipelineRunOptions) (*string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + ".glb")
	if err := p.gltfProc.ToGLB(inputPath, outputPath); err != nil {
		return nil, err
	}
	defer func(path string) {
		_, err := os.Stat(path)
		if os.IsExist(err) {
			if err := os.Remove(path); err != nil {
				infra.GetLogger().Error(err)
			}
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
	if err := p.apiClient.PatchSnapshot(client.SnapshotPatchOptions{
		Options: opts,
		Fields:  []string{client.SnapshotFieldPreview},
		Preview: &client.S3Object{
			Bucket: opts.Bucket,
			Key:    glbKey,
			Size:   helper.ToPtr(stat.Size()),
		},
	}); err != nil {
		return nil, err
	}
	return &glbKey, nil
}
