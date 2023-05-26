package pipeline

import (
	"os"
	"path/filepath"
	"voltaserve/config"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	log "github.com/sirupsen/logrus"
)

type ImagePipeline struct {
	s3              *infra.S3Manager
	snapshotRepo    repo.SnapshotRepo
	fileSearch      *search.FileSearch
	ocrPipeline     *OCRPipeline
	cmd             *infra.Command
	imageProc       *infra.ImageProcessor
	metadataUpdater *metadataUpdater
	config          config.Config
}

type ImagePipelineOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

func NewImagePipeline() *ImagePipeline {
	return &ImagePipeline{
		s3:              infra.NewS3Manager(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		fileSearch:      search.NewFileSearch(),
		ocrPipeline:     NewOCRPipeline(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		imageProc:       infra.NewImageProcessor(),
		config:          config.GetConfig(),
	}
}

func (p *ImagePipeline) Run(opts ImagePipelineOptions) error {
	snapshot, err := p.snapshotRepo.Find(opts.SnapshotId)
	if err != nil {
		return err
	}
	inputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + filepath.Ext(opts.S3Key))
	if err := p.s3.GetFile(opts.S3Key, inputPath, opts.S3Bucket); err != nil {
		return err
	}
	if filepath.Ext(opts.S3Key) == ".tiff" {
		newInputFile := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".jpg")
		if err := p.imageProc.Convert(inputPath, newInputFile); err != nil {
			return err
		}
		if err := os.Remove(inputPath); err != nil {
			return err
		}
		inputPath = newInputFile
	}
	if err := p.measureImageProps(snapshot, inputPath); err != nil {
		return err
	}
	if err := p.generateThumbnail(snapshot, inputPath); err != nil {
		return err
	}
	if err := p.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	ocrData, err := p.ocrPipeline.imageToData(inputPath)
	if err == nil && ocrData.PositiveConfCount > ocrData.NegativeConfCount {
		if err := p.ocrPipeline.Run(OCRPipelineOptions(opts)); err != nil {
			/*
				Here we intentionally ignore the error, here is the explanation why:
				The reason we came here to begin with is because of
				this condition: 'ocrData.PositiveConfCount > ocrData.NegativeConfCount',
				but it turned out that the OCR failed, that means probably the image
				does not contain text after all ¯\_(ツ)_/¯
				So we log the error and move on...
			*/
			log.Error(err)
		}
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	return nil
}

func (p *ImagePipeline) measureImageProps(snapshot model.Snapshot, inputPath string) error {
	width, height, err := p.imageProc.Measure(inputPath)
	if err != nil {
		return err
	}
	original := snapshot.GetOriginal()
	original.Image = &model.ImageProps{
		Width:  width,
		Height: height,
	}
	snapshot.SetOriginal(original)
	return nil
}

func (p *ImagePipeline) generateThumbnail(snapshot model.Snapshot, inputPath string) error {
	width := snapshot.GetOriginal().Image.Width
	height := snapshot.GetOriginal().Image.Height
	if width > p.config.Limits.ImagePreviewMaxWidth || height > p.config.Limits.ImagePreviewMaxHeight {
		outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + filepath.Ext(inputPath))
		if width > height {
			if err := p.imageProc.Resize(inputPath, p.config.Limits.ImagePreviewMaxWidth, 0, outputPath); err != nil {
				return err
			}
		} else {
			if err := p.imageProc.Resize(inputPath, 0, p.config.Limits.ImagePreviewMaxHeight, outputPath); err != nil {
				return err
			}
		}
		b64, err := infra.ImageToBase64(outputPath)
		if err != nil {
			return err
		}
		thumbnailWidth, thumbnailHeight, err := p.imageProc.Measure(outputPath)
		if err != nil {
			return err
		}
		snapshot.SetThumbnail(&model.Thumbnail{
			Base64: b64,
			Width:  thumbnailWidth,
			Height: thumbnailHeight,
		})
		if _, err := os.Stat(outputPath); err == nil {
			if err := os.Remove(outputPath); err != nil {
				return err
			}
		}
	} else {
		b64, err := infra.ImageToBase64(inputPath)
		if err != nil {
			return err
		}
		thumbnailWidth, thumbnailHeight, err := p.imageProc.Measure(inputPath)
		if err != nil {
			return err
		}
		snapshot.SetThumbnail(&model.Thumbnail{
			Base64: b64,
			Width:  thumbnailWidth,
			Height: thumbnailHeight,
		})
	}
	return nil
}
