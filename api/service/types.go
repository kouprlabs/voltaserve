package service

import (
	"bytes"
	"voltaserve/model"
)

type DownloadResult struct {
	Buffer      *bytes.Buffer
	PartialSize *int64
	TotalSize   *int64
	File        model.File
	Snapshot    model.Snapshot
}
