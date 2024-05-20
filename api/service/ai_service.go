package service

import (
	"os"
	"path/filepath"
	"voltaserve/client"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"

	"go.uber.org/zap"
)

type Language struct {
	ID      string `json:"id"`
	ISO6393 string `json:"iso6393"`
	Name    string `json:"name"`
}

type AIService struct {
	languages    []*Language
	snapshotRepo repo.SnapshotRepo
	s3           *infra.S3Manager
	toolClient   *client.ToolClient
	logger       *zap.SugaredLogger
}

func NewAIService() *AIService {
	return &AIService{
		languages: []*Language{
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
		snapshotRepo: repo.NewSnapshotRepo(),
		s3:           infra.NewS3Manager(),
		toolClient:   client.NewToolClient(),
	}
}

func (svc *AIService) GetAvailableLanguages() ([]*Language, error) {
	return svc.languages, nil
}

type AIUpdateLanguageOptions struct {
	SnapshotID string `json:"snapshotId" validate:"required"`
	LanguageID string `json:"languageId" validate:"required"`
}

func (svc *AIService) UpdateLanguage(opts AIUpdateLanguageOptions) error {
	s, err := svc.snapshotRepo.Find(opts.SnapshotID)
	if err != nil {
		return err
	}
	s.SetLanguage(opts.LanguageID)
	if err := svc.snapshotRepo.Save(s); err != nil {
		return err
	}
	return nil
}

type AIExtractTextOptions struct {
	FileID     string `json:"fileId" validate:"required"`
	SnapshotID string `json:"snapshotId" validate:"required"`
}

func (svc *AIService) ExtractText(opts AIExtractTextOptions) error {
	s, err := svc.snapshotRepo.Find(opts.SnapshotID)
	if err != nil {
		return err
	}
	if s.GetOriginal() == nil {
		return errorpkg.NewS3ObjectNotFoundError(err)
	}
	if s.GetLanguage() == nil {
		return errorpkg.NewSnapshotLanguageNotSet(err)
	}
	original := s.GetOriginal()

	/* Download original */
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helper.NewID() + filepath.Ext(original.Key))
	if err := svc.s3.GetFile(original.Key, inputPath, original.Bucket); err != nil {
		return err
	}
	defer func(inputPath string, logger *zap.SugaredLogger) {
		_, err := os.Stat(inputPath)
		if os.IsExist(err) {
			if err := os.Remove(inputPath); err != nil {
				logger.Error(err)
			}
		}
	}(inputPath, svc.logger)

	/* Get DPI */
	dpi, err := svc.toolClient.DPIFromImage(inputPath)
	if err != nil {
		dpi = 72
	}

	/* Convert to PDF/A */
	pdfPath, err := svc.toolClient.OCRFromPDF(inputPath, s.GetLanguage(), &dpi)
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
		Bucket: original.Bucket,
		Key:    opts.FileID + "/" + opts.SnapshotID + "/ocr.pdf",
		Size:   stat.Size(),
	}
	if err := svc.s3.PutFile(s3Object.Key, pdfPath, helper.DetectMimeFromFile(pdfPath), s3Object.Bucket); err != nil {
		return err
	}
	s.SetOCR(&s3Object)
	if err := svc.snapshotRepo.Save(s); err != nil {
		return err
	}

	/* Extract text */
	text, err := svc.toolClient.TextFromPDF(inputPath)
	if err != nil {
		svc.logger.Errorw(err.Error())
	}
	if text == "" || err != nil {
		return err
	}

	/* Set text S3 object */
	s3Object = model.S3Object{
		Bucket: original.Bucket,
		Key:    opts.FileID + "/" + opts.SnapshotID + "/text.txt",
		Size:   int64(len(text)),
	}
	if err := svc.s3.PutText(s3Object.Key, text, "text/plain", s3Object.Bucket); err != nil {
		return err
	}
	s.SetText(&s3Object)
	if err := svc.snapshotRepo.Save(s); err != nil {
		return err
	}

	return nil
}
