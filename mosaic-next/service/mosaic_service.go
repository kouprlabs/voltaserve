package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime"
	"os"
	"path/filepath"
	"voltaserve/builder"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"
)

type MosaicService struct {
	s3     *infra.S3Manager
	config *config.Config
}

func NewMosaicService() *MosaicService {
	return &MosaicService{
		s3:     infra.NewS3Manager(),
		config: config.GetConfig(),
	}
}

func (svc *MosaicService) Create(path, s3Key, s3Bucket string) (*builder.Metadata, error) {
	tmpDir := filepath.Join(os.TempDir(), helper.NewID())
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			infra.GetLogger().Error(err)
		}
	}()
	metadata, err := builder.NewMosaicBuilder(builder.MosaicBuilderOptions{
		File:            path,
		OutputDirectory: tmpDir,
	}).Build()
	if err != nil {
		return nil, err
	}
	var files []string
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		contentType := mime.TypeByExtension(filepath.Ext(file))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		putOptions := minio.PutObjectOptions{ContentType: contentType}
		relativePath, err := filepath.Rel(tmpDir, file)
		if err != nil {
			return nil, err
		}
		destinationKey := filepath.Join(s3Key, "mosaic", relativePath)
		if err := svc.s3.PutFile(destinationKey, file, contentType, s3Bucket, putOptions); err != nil {
			return nil, err
		}
	}
	return metadata, nil
}

func (svc *MosaicService) Delete(s3Bucket, s3Key string) error {
	listOptions := minio.ListObjectsOptions{
		Prefix:    filepath.Join(s3Key, "mosaic"),
		Recursive: true,
	}
	objects, err := svc.s3.ListObjects(s3Bucket, listOptions)
	if err != nil {
		return errorpkg.NewResourceNotFoundError(err)
	}
	for _, object := range objects {
		if err := svc.s3.RemoveObject(object.Key, s3Bucket, minio.RemoveObjectOptions{}); err != nil {
			return errorpkg.NewResourceNotFoundError(err)
		}
	}
	return nil
}

func (svc *MosaicService) GetTileBuffer(s3Bucket, s3Key string, zoomLevel, row, col int, ext string) (*bytes.Buffer, *string, error) {
	objectName := filepath.Join(s3Key, "mosaic", fmt.Sprintf("%d/%dx%d.%s", zoomLevel, row, col, ext))
	buf := new(bytes.Buffer)
	if _, err := svc.s3.GetObjectWithBuffer(objectName, s3Bucket, buf, minio.GetObjectOptions{}); err != nil {
		return nil, nil, errorpkg.NewResourceNotFoundError(err)
	}
	contentType := mime.TypeByExtension("." + ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return buf, &contentType, nil
}

func (svc *MosaicService) GetMetadata(s3Bucket, s3Key string) (*builder.Metadata, error) {
	objectName := filepath.Join(s3Key, "mosaic", "mosaic.json")
	text, err := svc.s3.GetText(objectName, s3Bucket, minio.GetObjectOptions{})
	if err != nil {
		return nil, errorpkg.NewResourceNotFoundError(err)
	}
	var metadata builder.Metadata
	if err := json.Unmarshal([]byte(text), &metadata); err != nil {
		return nil, err
	}
	return &metadata, nil
}
