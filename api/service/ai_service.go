package service

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"voltaserve/cache"
	"voltaserve/client"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"

	"go.uber.org/zap"
)

type AILanguage struct {
	ID      string `json:"id"`
	ISO6393 string `json:"iso6393"`
	Name    string `json:"name"`
}

type AIService struct {
	languages      []*AILanguage
	snapshotRepo   repo.SnapshotRepo
	userRepo       repo.UserRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	s3             *infra.S3Manager
	toolClient     *client.ToolClient
	languageClient *client.LanguageClient
	logger         *zap.SugaredLogger
}

func NewAIService() *AIService {
	return &AIService{
		languages: []*AILanguage{
			{ID: "ara", ISO6393: "ara", Name: "Arabic"},
			{ID: "chi_sim", ISO6393: "zho", Name: "Chinese Simplified"},
			{ID: "chi_tra", ISO6393: "zho", Name: "Chinese Traditional"},
			{ID: "deu", ISO6393: "deu", Name: "German"},
			{ID: "eng", ISO6393: "eng", Name: "English"},
			{ID: "fra", ISO6393: "fra", Name: "French"},
			{ID: "hin", ISO6393: "hin", Name: "Hindi"},
			{ID: "ita", ISO6393: "ita", Name: "Italian"},
			{ID: "jpn", ISO6393: "jpn", Name: "Japanese"},
			{ID: "nld", ISO6393: "nld", Name: "Dutch"},
			{ID: "por", ISO6393: "por", Name: "Portuguese"},
			{ID: "rus", ISO6393: "rus", Name: "Russian"},
			{ID: "spa", ISO6393: "spa", Name: "Spanish"},
			{ID: "swe", ISO6393: "swe", Name: "Swedish"},
		},
		snapshotRepo:   repo.NewSnapshotRepo(),
		userRepo:       repo.NewUserRepo(),
		fileCache:      cache.NewFileCache(),
		fileGuard:      guard.NewFileGuard(),
		s3:             infra.NewS3Manager(),
		toolClient:     client.NewToolClient(),
		languageClient: client.NewLanguageClient(),
	}
}

func (svc *AIService) GetAvailableLanguages() ([]*AILanguage, error) {
	return svc.languages, nil
}

type AIUpdateLanguageOptions struct {
	LanguageID string `json:"languageId" validate:"required"`
}

func (svc *AIService) UpdateLanguage(id string, opts AIUpdateLanguageOptions, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	snapshot.SetLanguage(opts.LanguageID)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *AIService) ExtractText(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	if !snapshot.HasOriginal() {
		return errorpkg.NewS3ObjectNotFoundError(err)
	}
	if snapshot.GetLanguage() == nil {
		return errorpkg.NewSnapshotLanguageNotSet(err)
	}

	/* Download original */
	originalPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(snapshot.GetOriginal().Key))
	if err := svc.s3.GetFile(snapshot.GetOriginal().Key, originalPath, snapshot.GetOriginal().Bucket); err != nil {
		return err
	}
	defer func(inputPath string, logger *zap.SugaredLogger) {
		_, err := os.Stat(inputPath)
		if os.IsExist(err) {
			if err := os.Remove(inputPath); err != nil {
				logger.Error(err)
			}
		}
	}(originalPath, svc.logger)

	/* Get DPI */
	dpi, err := svc.toolClient.DPIFromImage(originalPath)
	if err != nil {
		dpi = 72
	}

	/* Convert to PDF/A */
	pdfPath, err := svc.toolClient.OCRFromPDF(originalPath, snapshot.GetLanguage(), &dpi)
	if err != nil {
		svc.logger.Errorw(err.Error())
	}
	defer func(pdfPath string, logger *zap.SugaredLogger) {
		_, err := os.Stat(pdfPath)
		if os.IsExist(err) {
			if err := os.Remove(pdfPath); err != nil {
				logger.Error(err)
			}
		}
	}(pdfPath, svc.logger)

	/* Set OCR S3 object */
	stat, err := os.Stat(pdfPath)
	if err != nil {
		return err
	}
	s3Object := model.S3Object{
		Bucket: snapshot.GetOriginal().Bucket,
		Key:    snapshot.GetID() + "/ocr.pdf",
		Size:   stat.Size(),
	}
	if err := svc.s3.PutFile(s3Object.Key, pdfPath, helper.DetectMimeFromFile(pdfPath), s3Object.Bucket); err != nil {
		return err
	}
	snapshot.SetOCR(&s3Object)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}

	/* Extract text */
	text, err := svc.toolClient.TextFromPDF(pdfPath)
	if err != nil {
		svc.logger.Errorw(err.Error())
	}
	if text == "" || err != nil {
		return err
	}

	/* Set text S3 object */
	s3Object = model.S3Object{
		Bucket: snapshot.GetOriginal().Bucket,
		Key:    snapshot.GetID() + "/text.txt",
		Size:   int64(len(text)),
	}
	if err := svc.s3.PutText(s3Object.Key, text, "text/plain", s3Object.Bucket); err != nil {
		return err
	}
	snapshot.SetText(&s3Object)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}

	return nil
}

func (svc *AIService) ScanEntities(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionEditor); err != nil {
		return err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return err
	}
	if !snapshot.HasText() {
		return errorpkg.NewS3ObjectNotFoundError(err)
	}
	text, err := svc.s3.GetText(snapshot.GetText().Key, snapshot.GetText().Bucket)
	if err != nil {
		return err
	}
	res, err := svc.languageClient.GetEntities(client.GetEntitiesOptions{Text: text})
	if err != nil {
		return err
	}
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	entities := string(data)
	s3Object := model.S3Object{
		Bucket: snapshot.GetOriginal().Bucket,
		Key:    snapshot.GetID() + "/entities.json",
		Size:   int64(len(entities)),
	}
	if err := svc.s3.PutText(s3Object.Key, entities, "application/json", s3Object.Bucket); err != nil {
		return err
	}
	snapshot.SetEntities(&s3Object)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	return nil
}

type AISummary struct {
	HasLanguage bool
	HasOCR      bool
	HasText     bool
	HasEntities bool
}

func (svc *AIService) GetSummary(id string, userID string) (*AISummary, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	return &AISummary{
		HasLanguage: snapshot.GetLanguage() != nil,
		HasOCR:      snapshot.HasOCR(),
		HasText:     snapshot.HasText(),
		HasEntities: snapshot.HasEntities(),
	}, nil
}

func (svc *AIService) DownloadTextBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, nil, nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return nil, nil, nil, err
	}
	if snapshot == nil {
		return nil, nil, nil, errorpkg.NewSnapshotNotFoundError(nil)
	}
	if snapshot.HasText() {
		buf, err := svc.s3.GetObject(snapshot.GetText().Key, snapshot.GetText().Bucket)
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, snapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *AIService) DownloadOCRBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, nil, nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, nil, nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, nil, nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return nil, nil, nil, err
	}
	if snapshot == nil {
		return nil, nil, nil, errorpkg.NewSnapshotNotFoundError(nil)
	}
	if snapshot.HasOCR() {
		buf, err := svc.s3.GetObject(snapshot.GetOCR().Key, snapshot.GetOCR().Bucket)
		if err != nil {
			return nil, nil, nil, err
		}
		return buf, file, snapshot, nil
	} else {
		return nil, nil, nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

type AIEntitiesListOptions struct {
	Query     string
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
}

type AIEntitiesList struct {
	Data          []*model.AIEntity `json:"data"`
	TotalPages    uint              `json:"totalPages"`
	TotalElements uint              `json:"totalElements"`
	Page          uint              `json:"page"`
	Size          uint              `json:"size"`
}

func (svc *AIService) ListEntities(id string, opts AIEntitiesListOptions, userID string) (*AIEntitiesList, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	file, err := svc.fileCache.Get(id)
	if err != nil {
		return nil, err
	}
	if err = svc.fileGuard.Authorize(user, file, model.PermissionViewer); err != nil {
		return nil, err
	}
	if file.GetType() != model.FileTypeFile || file.GetSnapshotID() == nil {
		return nil, errorpkg.NewFileIsNotAFileError(file)
	}
	snapshot, err := svc.snapshotRepo.Find(*file.GetSnapshotID())
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, errorpkg.NewSnapshotNotFoundError(nil)
	}
	if snapshot.HasEntities() {
		text, err := svc.s3.GetText(snapshot.GetEntities().Key, snapshot.GetEntities().Bucket)
		if err != nil {
			return nil, err
		}
		var entities []*model.AIEntity
		if err := json.Unmarshal([]byte(text), &entities); err != nil {
			return nil, err
		}
		if opts.SortBy == "" {
			opts.SortBy = SortByName
		}
		filtered := svc.doFiltering(entities, opts.Query)
		sorted := svc.doSorting(filtered, opts.SortBy, opts.SortOrder)
		data, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
		return &AIEntitiesList{
			Data:          data,
			TotalPages:    totalPages,
			TotalElements: totalElements,
			Page:          opts.Page,
			Size:          uint(len(data)),
		}, nil
	} else {
		return nil, errorpkg.NewS3ObjectNotFoundError(nil)
	}
}

func (svc *AIService) doFiltering(data []*model.AIEntity, query string) []*model.AIEntity {
	if query == "" {
		return data
	}
	var filtered []*model.AIEntity
	for _, entity := range data {
		if strings.Contains(strings.ToLower(entity.Text), strings.ToLower(query)) {
			filtered = append(filtered, entity)
		}
	}
	return filtered
}

func (svc *AIService) doSorting(data []*model.AIEntity, sortBy string, sortOrder string) []*model.AIEntity {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].Text > data[j].Text
			} else {
				return data[i].Text < data[j].Text
			}
		})
		return data
	}
	return data
}

func (svc *AIService) doPagination(data []*model.AIEntity, page, size uint) ([]*model.AIEntity, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return nil, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}
