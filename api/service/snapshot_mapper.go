// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"path/filepath"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
)

type snapshotMapper struct {
	taskCache *cache.TaskCache
}

func newSnapshotMapper() *snapshotMapper {
	return &snapshotMapper{
		taskCache: cache.NewTaskCache(),
	}
}

func (mp *snapshotMapper) mapOne(m model.Snapshot) *Snapshot {
	s := &Snapshot{
		ID:         m.GetID(),
		Version:    m.GetVersion(),
		Status:     m.GetStatus(),
		Language:   m.GetLanguage(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	if m.HasOriginal() {
		s.Original = mp.mapS3Object(m.GetOriginal())
	}
	if m.HasPreview() {
		s.Preview = mp.mapS3Object(m.GetPreview())
	}
	if m.HasOCR() {
		s.OCR = mp.mapS3Object(m.GetOCR())
	}
	if m.HasText() {
		s.Text = mp.mapS3Object(m.GetText())
	}
	if m.HasEntities() {
		s.Entities = mp.mapS3Object(m.GetEntities())
	}
	if m.HasMosaic() {
		s.Mosaic = mp.mapS3Object(m.GetMosaic())
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.mapS3Object(m.GetThumbnail())
	}
	if m.GetTaskID() != nil {
		s.Task = &SnapshotTaskInfo{
			ID: *m.GetTaskID(),
		}
		isPending, err := isTaskPending(m, mp.taskCache)
		if err != nil {
			log.GetLogger().Error(err)
		} else {
			s.Task.IsPending = isPending
		}
	}

	return s
}

func (mp *snapshotMapper) mapMany(snapshots []model.Snapshot, activeID string) []*Snapshot {
	res := make([]*Snapshot, 0)
	for _, snapshot := range snapshots {
		s := mp.mapOne(snapshot)
		s.IsActive = activeID == snapshot.GetID()
		res = append(res, s)
	}
	return res
}

func (mp *snapshotMapper) mapS3Object(o *model.S3Object) *Download {
	download := &Download{
		Extension: filepath.Ext(o.Key),
		Size:      o.Size,
	}
	if o.Image != nil {
		download.Image = o.Image
	}
	if o.Document != nil {
		download.Document = o.Document
	}
	return download
}
