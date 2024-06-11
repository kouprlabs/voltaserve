package service

import (
	"bytes"
	"voltaserve/infra"
	"voltaserve/model"
)

type DownloadResult struct {
	RangeInterval *infra.RangeInterval
	Buffer        *bytes.Buffer
	File          model.File
	Snapshot      model.Snapshot
}
