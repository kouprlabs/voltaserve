package infra

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"voltaserve/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/sse"
)

type S3Manager struct {
	config config.S3Config
	client *minio.Client
}

func NewS3Manager() *S3Manager {
	mgr := new(S3Manager)
	mgr.config = config.GetConfig().S3
	return mgr
}

func (mgr *S3Manager) GetFile(objectName string, filePath string, bucketName string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	if err := mgr.client.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) PutFile(objectName string, filePath string, contentType string, bucketName string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if _, err := mgr.client.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) PutText(objectName string, text string, contentType string, bucketName string) error {
	if contentType != "" && contentType != "text/plain" && contentType != "application/json" {
		return errors.New("invalid content type")
	}
	if contentType == "" {
		contentType = "text/plain"
	}
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	if _, err := mgr.client.PutObject(context.Background(), bucketName, objectName, strings.NewReader(text), int64(len(text)), minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) GetObject(objectName string, bucketName string) (*bytes.Buffer, error) {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return nil, err
		}
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = io.Copy(io.Writer(&buf), reader)
	if err != nil {
		return nil, nil
	}
	return &buf, nil
}

func (mgr *S3Manager) GetText(objectName string, bucketName string) (string, error) {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return "", err
		}
	}
	reader, err := mgr.client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", nil
	}
	return buf.String(), nil
}

func (mgr *S3Manager) RemoveObject(objectName string, bucketName string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	err := mgr.client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) CreateBucket(bucketName string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	found, err := mgr.client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if !found {
		if err = mgr.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
			Region: mgr.config.Region,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *S3Manager) RemoveBucket(bucketName string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	found, err := mgr.client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	objectCh := mgr.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	mgr.client.RemoveObjects(context.Background(), bucketName, objectCh, minio.RemoveObjectsOptions{})
	if err = mgr.client.RemoveBucket(context.Background(), bucketName); err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) EnableBucketEncryption(bucketName string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	err := mgr.client.SetBucketEncryption(context.Background(), bucketName, sse.NewConfigurationSSES3())
	if err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) DisableBucketEncryption(bucketName string) error {
	if mgr.client == nil {
		if err := mgr.connect(); err != nil {
			return err
		}
	}
	err := mgr.client.RemoveBucketEncryption(context.Background(), bucketName)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *S3Manager) connect() error {
	client, err := minio.New(mgr.config.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(mgr.config.AccessKey, mgr.config.SecretKey, ""),
		Secure: mgr.config.Secure,
	})
	if err != nil {
		return err
	}
	mgr.client = client
	return nil
}
