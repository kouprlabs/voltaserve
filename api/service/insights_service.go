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

type InsightsService struct {
	languages      []*InsightsLanguage
	snapshotRepo   repo.SnapshotRepo
	userRepo       repo.UserRepo
	fileCache      *cache.FileCache
	fileGuard      *guard.FileGuard
	s3             *infra.S3Manager
	toolClient     *client.ToolClient
	languageClient *client.LanguageClient
	fileIdent      *infra.FileIdentifier
	logger         *zap.SugaredLogger
}

func NewInsightsService() *InsightsService {
	logger, err := infra.GetLogger()
	if err != nil {
		panic(err)
	}
	return &InsightsService{
		languages: []*InsightsLanguage{
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
		fileIdent:      infra.NewFileIdentifier(),
		logger:         logger,
	}
}

type InsightsLanguage struct {
	ID      string `json:"id"`
	ISO6393 string `json:"iso6393"`
	Name    string `json:"name"`
}

func (svc *InsightsService) GetLanguages() ([]*InsightsLanguage, error) {
	return svc.languages, nil
}

type InsightsCreateOptions struct {
	LanguageID string `json:"languageId" validate:"required"`
}

func (svc *InsightsService) Create(id string, opts InsightsCreateOptions, userID string) error {
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
	return svc.create(snapshot)
}

func (svc *InsightsService) create(snapshot model.Snapshot) error {
	if err := svc.createText(snapshot); err != nil {
		return err
	}
	if err := svc.createEntities(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *InsightsService) Patch(id string, userID string) error {
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
	previousSnapshot, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
	if err != nil {
		return err
	}
	if previousSnapshot == nil || previousSnapshot.GetLanguage() == nil {
		return errorpkg.NewSnapshotCannotBePatchedError(nil)
	}
	snapshot.SetLanguage(*previousSnapshot.GetLanguage())
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	return svc.create(snapshot)
}

func (svc *InsightsService) createText(snapshot model.Snapshot) error {
	if snapshot.HasText() {
		return nil
	}
	if !snapshot.HasOriginal() {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	if snapshot.GetLanguage() == nil {
		return errorpkg.NewSnapshotLanguageNotSetError(nil)
	}

	/* Download originalS3 */
	originalS3 := snapshot.GetOriginal()
	originalPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(originalS3.Key))
	if err := svc.s3.GetFile(originalS3.Key, originalPath, originalS3.Bucket); err != nil {
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

	var pdfPath string
	if svc.fileIdent.IsImage(originalS3.Key) {
		/* Get DPI */
		dpi, err := svc.toolClient.DPIFromImage(originalPath)
		if err != nil {
			dpi = 72
		}

		/* Remove alpha channel */
		noAlphaImagePath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(originalS3.Key))
		if err := svc.toolClient.RemoveAlphaChannel(originalPath, noAlphaImagePath); err != nil {
			return err
		}
		defer func(inputPath string, logger *zap.SugaredLogger) {
			_, err := os.Stat(inputPath)
			if os.IsExist(err) {
				if err := os.Remove(inputPath); err != nil {
					logger.Error(err)
				}
			}
		}(noAlphaImagePath, svc.logger)

		/* Convert to PDF/A */
		pdfPath, err = svc.toolClient.OCRFromPDF(noAlphaImagePath, snapshot.GetLanguage(), &dpi)
		if err != nil {
			return err
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
	} else if svc.fileIdent.IsPDF(originalS3.Key) || svc.fileIdent.IsOffice(originalS3.Key) || svc.fileIdent.IsPlainText(originalS3.Key) {
		pdfPath = originalPath
	} else {
		return errorpkg.NewUnsupportedFileTypeError(nil)
	}

	/* Extract text */
	text, err := svc.toolClient.TextFromPDF(pdfPath)
	if text == "" || err != nil {
		return err
	}

	/* Set text S3 object */
	s3Object := model.S3Object{
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

func (svc *InsightsService) createEntities(snapshot model.Snapshot) error {
	if snapshot.HasEntities() {
		return nil
	}
	if !snapshot.HasText() {
		return errorpkg.NewS3ObjectNotFoundError(nil)
	}
	if snapshot.GetLanguage() == nil {
		return errorpkg.NewSnapshotLanguageNotSetError(nil)
	}
	text, err := svc.s3.GetText(snapshot.GetText().Key, snapshot.GetText().Bucket)
	if err != nil {
		return err
	}
	if len(text) > 1000000 {
		return errorpkg.NewSnapshotTextLengthExceedsLimitError(err)
	}
	res, err := svc.languageClient.GetEntities(client.GetEntitiesOptions{
		Text:     text,
		Language: *snapshot.GetLanguage(),
	})
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

func (svc *InsightsService) Delete(id string, userID string) error {
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
	if svc.fileIdent.IsImage(snapshot.GetOriginal().Key) {
		if err := svc.deleteText(snapshot); err != nil {
			return err
		}
	}
	if err := svc.deleteEntities(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *InsightsService) deleteText(snapshot model.Snapshot) error {
	if !snapshot.HasText() {
		return nil
	}
	s3Object := snapshot.GetText()
	if err := svc.s3.RemoveObject(s3Object.Key, s3Object.Bucket); err != nil {
		return err
	}
	snapshot.SetText(nil)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	return nil
}

func (svc *InsightsService) deleteEntities(snapshot model.Snapshot) error {
	if !snapshot.HasEntities() {
		return nil
	}
	s3Object := snapshot.GetEntities()
	if err := svc.s3.RemoveObject(s3Object.Key, s3Object.Bucket); err != nil {
		return err
	}
	snapshot.SetEntities(nil)
	if err := svc.snapshotRepo.Save(snapshot); err != nil {
		return err
	}
	return nil
}

type InsightsListEntitiesOptions struct {
	Query     string `json:"query"`
	Page      uint   `json:"page"`
	Size      uint   `json:"size"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

type InsightsEntityList struct {
	Data          []*model.InsightsEntity `json:"data"`
	TotalPages    uint                    `json:"totalPages"`
	TotalElements uint                    `json:"totalElements"`
	Page          uint                    `json:"page"`
	Size          uint                    `json:"size"`
}

func (svc *InsightsService) ListEntities(id string, opts InsightsListEntitiesOptions, userID string) (*InsightsEntityList, error) {
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
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous != nil {
			snapshot = previous
		}
	}
	if snapshot.HasEntities() {
		text, err := svc.s3.GetText(snapshot.GetEntities().Key, snapshot.GetEntities().Bucket)
		if err != nil {
			return nil, err
		}
		var entities []*model.InsightsEntity
		if err := json.Unmarshal([]byte(text), &entities); err != nil {
			return nil, err
		}
		if opts.SortBy == "" {
			opts.SortBy = SortByName
		}
		filtered := svc.doFiltering(entities, opts.Query)
		sorted := svc.doSorting(filtered, opts.SortBy, opts.SortOrder)
		data, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
		return &InsightsEntityList{
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

func (svc *InsightsService) doFiltering(data []*model.InsightsEntity, query string) []*model.InsightsEntity {
	if query == "" {
		return data
	}
	filtered := []*model.InsightsEntity{}
	for _, entity := range data {
		if strings.Contains(strings.ToLower(entity.Text), strings.ToLower(query)) {
			filtered = append(filtered, entity)
		}
	}
	return filtered
}

func (svc *InsightsService) doSorting(data []*model.InsightsEntity, sortBy string, sortOrder string) []*model.InsightsEntity {
	if sortBy == SortByName {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].Text > data[j].Text
			} else {
				return data[i].Text < data[j].Text
			}
		})
		return data
	} else if sortBy == SortByFrequency {
		sort.Slice(data, func(i, j int) bool {
			return data[i].Frequency > data[j].Frequency
		})
	}
	return data
}

func (svc *InsightsService) doPagination(data []*model.InsightsEntity, page, size uint) ([]*model.InsightsEntity, uint, uint) {
	totalElements := uint(len(data))
	totalPages := (totalElements + size - 1) / size
	if page > totalPages {
		return []*model.InsightsEntity{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	pageData := data[startIndex:endIndex]
	return pageData, totalElements, totalPages
}

type InsightsSummary struct {
	IsOutdated  bool `json:"isOutdated"`
	HasLanguage bool `json:"hasLanguage"`
	HasOCR      bool `json:"hasOcr"`
	HasText     bool `json:"hasText"`
	HasEntities bool `json:"hasEntities"`
}

func (svc *InsightsService) GetSummary(id string, userID string) (*InsightsSummary, error) {
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
	isOutdated := false
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, err
		}
		if previous != nil {
			snapshot = previous
			isOutdated = true
		}
	}
	return &InsightsSummary{
		HasLanguage: snapshot.GetLanguage() != nil,
		HasOCR:      snapshot.HasOCR(),
		HasText:     snapshot.HasText(),
		HasEntities: snapshot.HasEntities(),
		IsOutdated:  isOutdated,
	}, nil
}

func (svc *InsightsService) DownloadTextBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
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
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, nil, nil, err
		}
		if previous != nil {
			snapshot = previous
		}
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

func (svc *InsightsService) DownloadOCRBuffer(id string, userID string) (*bytes.Buffer, model.File, model.Snapshot, error) {
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
	if !snapshot.HasEntities() {
		previous, err := svc.getPreviousSnapshot(file.GetID(), snapshot.GetVersion())
		if err != nil {
			return nil, nil, nil, err
		}
		if previous != nil {
			snapshot = previous
		}
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

func (svc *InsightsService) getPreviousSnapshot(fileID string, version int64) (model.Snapshot, error) {
	snaphots, err := svc.snapshotRepo.FindAllPrevious(fileID, version)
	if err != nil {
		return nil, err
	}
	for _, snapshot := range snaphots {
		if snapshot.HasEntities() {
			return snapshot, nil
		}
	}
	return nil, nil
}