package pipeline

import (
	"os"
	"path/filepath"
	"strings"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/service"
)

type StoragePipeline struct {
	s3             *infra.S3Manager
	snapshotRepo   repo.SnapshotRepo
	fileRepo       repo.FileRepo
	fileCache      *cache.FileCache
	fileMapper     *service.FileMapper
	pdfPipeline    *PDFPipeline
	imagePipeline  *ImagePipeline
	officePipeline *OfficePipeline
	videoPipeline  *VideoPipeline
	workspaceCache *cache.WorkspaceCache
	fileIdentifier *infra.FileIdentifier
	config         config.Config
}

type StoragePipelineOptions struct {
	FileId   string
	FilePath string
}

func NewStoragePipeline() *StoragePipeline {
	return &StoragePipeline{
		s3:             infra.NewS3Manager(),
		snapshotRepo:   repo.NewSnapshotRepo(),
		fileRepo:       repo.NewFileRepo(),
		fileCache:      cache.NewFileCache(),
		fileMapper:     service.NewFileMapper(),
		pdfPipeline:    NewPDFPipeline(),
		imagePipeline:  NewImagePipeline(),
		officePipeline: NewOfficePipeline(),
		videoPipeline:  NewVideoPipeline(),
		workspaceCache: cache.NewWorkspaceCache(),
		fileIdentifier: infra.NewFileIdentifier(),
		config:         config.GetConfig(),
	}
}

func (p *StoragePipeline) Run(opts StoragePipelineOptions, userId string) (*service.File, error) {
	file, err := p.fileRepo.Find(opts.FileId)
	if err != nil {
		return nil, err
	}
	if err = p.fileCache.Set(file); err != nil {
		return nil, err
	}
	workspace, err := p.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return nil, err
	}
	latestVersion, err := p.snapshotRepo.GetLatestVersionForFile(opts.FileId)
	if err != nil {
		return nil, err
	}
	snapshotId := helpers.NewId()
	snapshot := repo.NewSnapshot()
	snapshot.SetID(snapshotId)
	snapshot.SetVersion(latestVersion)
	if err = p.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	if err = p.snapshotRepo.MapWithFile(snapshotId, opts.FileId); err != nil {
		return nil, err
	}
	stat, err := os.Stat(opts.FilePath)
	if err != nil {
		return nil, err
	}
	original := model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileId + "/" + snapshotId + "/original" + strings.ToLower(filepath.Ext(opts.FilePath))),
		Size:   stat.Size(),
	}
	if err = p.s3.PutFile(original.Key, opts.FilePath, infra.DetectMimeFromFile(opts.FilePath), workspace.GetBucket()); err != nil {
		return nil, err
	}
	snapshot.SetOriginal(&original)
	if err := p.snapshotRepo.Save(snapshot); err != nil {
		return nil, err
	}
	if stat.Size() >= int64(p.config.Limits.FileProcessingMaxSizeMB*1000000) {
		v, err := p.fileMapper.MapFile(file, userId)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	if p.fileIdentifier.IsPDF(filepath.Ext(opts.FilePath)) {
		snapshot.SetPreview(&model.S3Object{
			Bucket: original.Bucket,
			Key:    original.Key,
			Size:   original.Size,
		})
		if err = p.snapshotRepo.Save(snapshot); err != nil {
			return nil, err
		}
		if file, err = p.fileRepo.Find(opts.FileId); err != nil {
			return nil, err
		}
		if err := p.fileCache.Set(file); err != nil {
			return nil, err
		}
		if err = p.pdfPipeline.Run(PDFPipelineOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if p.fileIdentifier.IsOffice(filepath.Ext(opts.FilePath)) || p.fileIdentifier.IsPlainText(filepath.Ext(opts.FilePath)) {
		if err = p.officePipeline.Run(OfficePipelineOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if p.fileIdentifier.IsImage(filepath.Ext(opts.FilePath)) {
		if err = p.imagePipeline.Run(ImagePipelineOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	} else if p.fileIdentifier.IsVideo(filepath.Ext(opts.FilePath)) {
		if err = p.videoPipeline.Run(VideoPipelineOptions{
			FileId:     opts.FileId,
			SnapshotId: snapshotId,
			S3Bucket:   workspace.GetBucket(),
			S3Key:      original.Key,
		}); err != nil {
			return nil, err
		}
	}
	file, err = p.fileCache.Refresh(file.GetID())
	if err != nil {
		return nil, err
	}
	v, err := p.fileMapper.MapFile(file, userId)
	if err != nil {
		return nil, err
	}
	return v, nil
}
