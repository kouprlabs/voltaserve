package conversion

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/config"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
	"voltaserve/search"

	log "github.com/sirupsen/logrus"
)

type imagePipeline struct {
	s3              *infra.S3Manager
	snapshotRepo    repo.SnapshotRepo
	fileSearch      *search.FileSearch
	pdfPipeline     Pipeline
	cmd             *infra.Command
	imageProc       *infra.ImageProcessor
	metadataUpdater *metadataUpdater
	config          config.Config
}

type imageToDataResponse struct {
	Data                []tesseractData
	NegativeConfCount   int64
	NegativeConfPercent float32
	PositiveConfCount   int64
	PositiveConfPercent float32
}

type tesseractData struct {
	BlockNum int64
	Conf     int64
	Height   int64
	Left     int64
	Level    int64
	LineNum  int64
	PageNum  int64
	ParNum   int64
	Text     string
	Top      int64
	Width    int64
	WordNum  int64
}

func NewImagePipeline() Pipeline {
	return &imagePipeline{
		s3:              infra.NewS3Manager(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		fileSearch:      search.NewFileSearch(),
		pdfPipeline:     NewPDFPipeline(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		imageProc:       infra.NewImageProcessor(),
		config:          config.GetConfig(),
	}
}

func (p *imagePipeline) Run(opts PipelineOptions) error {
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
	imageData, err := p.imageToData(inputPath)
	if err == nil && imageData.PositiveConfCount > imageData.NegativeConfCount {
		/* We treat this as a text image, we convert it to PDF/A */
		if err := p.pdfPipeline.Run(opts); err != nil {
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

func (p *imagePipeline) imageToData(inputPath string) (imageToDataResponse, error) {
	outFile := helpers.NewId()
	if err := p.cmd.Exec("tesseract", inputPath, outFile, "tsv"); err != nil {
		return imageToDataResponse{}, err
	}
	var res = imageToDataResponse{}
	outFile = outFile + ".tsv"
	f, err := os.Open(outFile)
	if err != nil {
		return imageToDataResponse{}, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return imageToDataResponse{}, err
	}
	text := string(b)
	lines := strings.Split(text, "\n")
	lines = lines[1 : len(lines)-2]
	for _, l := range lines {
		values := strings.Split(l, "\t")
		data := tesseractData{}
		data.Level, _ = strconv.ParseInt(values[0], 10, 64)
		data.PageNum, _ = strconv.ParseInt(values[1], 10, 64)
		data.BlockNum, _ = strconv.ParseInt(values[2], 10, 64)
		data.ParNum, _ = strconv.ParseInt(values[3], 10, 64)
		data.LineNum, _ = strconv.ParseInt(values[4], 10, 64)
		data.WordNum, _ = strconv.ParseInt(values[5], 10, 64)
		data.Left, _ = strconv.ParseInt(values[6], 10, 64)
		data.Top, _ = strconv.ParseInt(values[7], 10, 64)
		data.Width, _ = strconv.ParseInt(values[8], 10, 64)
		data.Height, _ = strconv.ParseInt(values[9], 10, 64)
		data.Conf, _ = strconv.ParseInt(values[10], 10, 64)
		data.Text = values[11]
		res.Data = append(res.Data, data)
	}
	for _, v := range res.Data {
		if v.Conf < 0 {
			res.NegativeConfCount++
		} else {
			res.PositiveConfCount++
		}
	}
	if len(res.Data) > 0 {
		res.NegativeConfPercent = float32((int(res.NegativeConfCount) * 100) / len(res.Data))
		res.PositiveConfPercent = float32((int(res.PositiveConfCount) * 100) / len(res.Data))
	}
	if err := os.Remove(outFile); err != nil {
		return imageToDataResponse{}, err
	}
	return res, nil
}

func (p *imagePipeline) measureImageProps(snapshot model.Snapshot, inputPath string) error {
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

func (p *imagePipeline) generateThumbnail(snapshot model.Snapshot, inputPath string) error {
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
