// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package mapper

import (
	"path/filepath"

	"github.com/kouprlabs/voltaserve/shared/cache"
	"github.com/kouprlabs/voltaserve/shared/config"
	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/infra"
	"github.com/kouprlabs/voltaserve/shared/model"
)

type SnapshotMapper struct {
	taskCache  *cache.TaskCache
	taskMapper *TaskMapper
	fileIdent  *infra.FileIdentifier
}

func NewSnapshotMapper(postgres config.PostgresConfig, redis config.RedisConfig, environment config.EnvironmentConfig) *SnapshotMapper {
	return &SnapshotMapper{
		taskCache:  cache.NewTaskCache(postgres, redis, environment),
		taskMapper: NewTaskMapper(postgres, redis, environment),
		fileIdent:  infra.NewFileIdentifier(),
	}
}

func (mp *SnapshotMapper) Map(m model.Snapshot) *dto.Snapshot {
	s := &dto.Snapshot{
		ID:         m.GetID(),
		Version:    m.GetVersion(),
		Language:   m.GetLanguage(),
		Summary:    m.GetSummary(),
		Intent:     m.GetIntent(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	if m.HasOriginal() {
		s.Original = mp.MapS3Object(m.GetOriginal())
		s.Capabilities.Original = true
	}
	if m.HasPreview() {
		s.Preview = mp.MapS3Object(m.GetPreview())
		s.Capabilities.Preview = true
	}
	if m.HasOCR() {
		s.OCR = mp.MapS3Object(m.GetOCR())
		s.Capabilities.OCR = true
	}
	if m.HasText() {
		s.Text = mp.MapS3Object(m.GetText())
		s.Capabilities.Text = true
	}
	if m.HasThumbnail() {
		s.Thumbnail = mp.MapS3Object(m.GetThumbnail())
		s.Capabilities.Thumbnail = true
	}
	if m.GetSummary() != nil {
		s.Capabilities.Summary = true
	}
	if m.HasEntities() {
		s.Capabilities.Entities = true
	}
	if m.HasMosaic() {
		s.Capabilities.Mosaic = true
	}
	if m.GetTaskID() != nil {
		task, err := mp.taskCache.Get(*m.GetTaskID())
		if err == nil {
			s.Task, _ = mp.taskMapper.Map(task)
		}
	}
	if m.HasOriginal() && m.GetIntent() == nil {
		if mp.fileIdent.IsDocument(m.GetOriginal().Key) {
			s.Intent = helper.ToPtr(model.SnapshotIntentDocument)
		} else if mp.fileIdent.IsImage(m.GetOriginal().Key) {
			s.Intent = helper.ToPtr(model.SnapshotIntentImage)
		} else if mp.fileIdent.IsAudio(m.GetOriginal().Key) {
			s.Intent = helper.ToPtr(model.SnapshotIntentAudio)
		} else if mp.fileIdent.IsVideo(m.GetOriginal().Key) {
			s.Intent = helper.ToPtr(model.SnapshotIntentVideo)
		} else if mp.fileIdent.Is3D(m.GetOriginal().Key) {
			s.Intent = helper.ToPtr(model.SnapshotIntent3D)
		}
	}
	return s
}

func (mp *SnapshotMapper) MapMany(snapshots []model.Snapshot, activeID *string) []*dto.Snapshot {
	res := make([]*dto.Snapshot, 0)
	for _, snapshot := range snapshots {
		s := mp.Map(snapshot)
		s.IsActive = activeID != nil && *activeID == snapshot.GetID()
		res = append(res, s)
	}
	return res
}

func (mp *SnapshotMapper) MapS3Object(o *model.S3Object) *dto.SnapshotDownloadable {
	download := &dto.SnapshotDownloadable{
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

func (mp *SnapshotMapper) MapWithS3Objects(m model.Snapshot) *dto.SnapshotWithS3Objects {
	s := &dto.SnapshotWithS3Objects{
		ID:         m.GetID(),
		Version:    m.GetVersion(),
		Original:   m.GetOriginal(),
		Preview:    m.GetPreview(),
		OCR:        m.GetOCR(),
		Text:       m.GetText(),
		Thumbnail:  m.GetThumbnail(),
		Entities:   m.GetEntities(),
		Mosaic:     m.GetMosaic(),
		Language:   m.GetLanguage(),
		Summary:    m.GetSummary(),
		Intent:     m.GetIntent(),
		CreateTime: m.GetCreateTime(),
		UpdateTime: m.GetUpdateTime(),
	}
	if m.GetTaskID() != nil {
		task, err := mp.taskCache.Get(*m.GetTaskID())
		if err == nil {
			s.Task, _ = mp.taskMapper.Map(task)
		}
	}
	return s
}
