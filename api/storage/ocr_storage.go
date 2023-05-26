package storage

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type ocrData struct {
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

type ocrImageToDataResponse struct {
	Data                []ocrData
	NegativeConfCount   int64
	NegativeConfPercent float32
	PositiveConfCount   int64
	PositiveConfPercent float32
}

type ocrStorage struct {
	minio           *infra.S3Manager
	snapshotRepo    repo.SnapshotRepo
	pdfStorage      *pdfStorage
	cmd             *infra.Command
	metadataUpdater *storageMetadataUpdater
	workspaceCache  *cache.WorkspaceCache
	fileCache       *cache.FileCache
	config          config.Config
}

type ocrOptions struct {
	FileId     string
	SnapshotId string
	S3Bucket   string
	S3Key      string
}

func newOcrStorage() *ocrStorage {
	return &ocrStorage{
		minio:           infra.NewS3Manager(),
		snapshotRepo:    repo.NewSnapshotRepo(),
		pdfStorage:      newPDFStorage(),
		cmd:             infra.NewCommand(),
		metadataUpdater: newMetadataUpdater(),
		workspaceCache:  cache.NewWorkspaceCache(),
		fileCache:       cache.NewFileCache(),
		config:          config.GetConfig(),
	}
}

func (svc *ocrStorage) store(opts ocrOptions) error {
	snapshot, err := svc.snapshotRepo.Find(opts.SnapshotId)
	if err != nil {
		return err
	}
	inputPath, err := svc.getSuitableInputPath(opts)
	if err != nil {
		return err
	}
	outputPath, _ := svc.generatePDFA(inputPath)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if err := svc.pdfStorage.store(pdfStorageOptions(opts)); err != nil {
			return err
		}
	} else {
		if err := svc.sendToPDFStorage(snapshot, opts, outputPath); err != nil {
			return err
		}
	}
	if _, err := os.Stat(inputPath); err == nil {
		if err := os.Remove(inputPath); err != nil {
			return err
		}
	}
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return err
		}
	}
	return nil
}

func (svc *ocrStorage) generatePDFA(inputPath string) (string, error) {
	outputPath := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".pdf")
	if err := svc.cmd.Exec("ocrmypdf", "--rotate-pages", "--clean", "--deskew", "--image-dpi=300", inputPath, outputPath); err != nil {
		return "", err
	}
	return outputPath, nil
}

func (svc *ocrStorage) sendToPDFStorage(snapshot model.CoreSnapshot, opts ocrOptions, outputPath string) error {
	file, err := svc.fileCache.Get(opts.FileId)
	if err != nil {
		return err
	}
	workspace, err := svc.workspaceCache.Get(file.GetWorkspaceID())
	if err != nil {
		return err
	}
	stat, err := os.Stat(outputPath)
	if err != nil {
		return err
	}
	ocrSize := stat.Size()
	snapshot.SetOCR(&model.S3Object{
		Bucket: workspace.GetBucket(),
		Key:    filepath.FromSlash(opts.FileId + "/" + opts.SnapshotId + "/ocr.pdf"),
		Size:   ocrSize,
	})
	if err := svc.minio.PutFile(snapshot.GetOCR().Key, outputPath, DetectMimeFromFile(outputPath), workspace.GetBucket()); err != nil {
		return err
	}
	if err := svc.metadataUpdater.update(snapshot, opts.FileId); err != nil {
		return err
	}
	if err := svc.pdfStorage.store(pdfStorageOptions{
		FileId:     opts.FileId,
		SnapshotId: opts.SnapshotId,
		S3Bucket:   opts.S3Bucket,
		S3Key:      snapshot.GetOCR().Key,
	}); err != nil {
		return err
	}
	return nil
}

func (svc *ocrStorage) getSuitableInputPath(opts ocrOptions) (string, error) {
	extension := filepath.Ext(opts.S3Key)
	path := filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + extension)
	if err := svc.minio.GetFile(opts.S3Key, path, opts.S3Bucket); err != nil {
		return "", err
	}

	// If an image, convert it to jpeg, because ocrmypdf supports jpeg only
	if extension == ".jpg" || extension == ".jpeg" {
		oldPath := path
		path = filepath.FromSlash(os.TempDir() + "/" + helpers.NewId() + ".jpg")
		if err := svc.cmd.Exec("gm", "convert", oldPath, path); err != nil {
			return "", err
		}
		if err := os.Remove(oldPath); err != nil {
			return "", err
		}
	}
	return path, nil
}

func (svc *ocrStorage) imageToData(inputPath string) (ocrImageToDataResponse, error) {
	outFile := helpers.NewId()
	if err := svc.cmd.Exec("tesseract", inputPath, outFile, "tsv"); err != nil {
		return ocrImageToDataResponse{}, err
	}
	var res = ocrImageToDataResponse{}
	outFile = outFile + ".tsv"
	f, err := os.Open(outFile)
	if err != nil {
		return ocrImageToDataResponse{}, err
	}
	b, err := io.ReadAll(f)
	if err != nil {
		return ocrImageToDataResponse{}, err
	}
	text := string(b)
	lines := strings.Split(text, "\n")
	lines = lines[1 : len(lines)-2]
	for _, l := range lines {
		values := strings.Split(l, "\t")
		data := ocrData{}
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
		return ocrImageToDataResponse{}, err
	}
	return res, nil
}
