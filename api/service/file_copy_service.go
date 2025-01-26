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
	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/log"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
	"github.com/kouprlabs/voltaserve/api/search"
	"time"
)

type FileCopyService struct {
	fileRepo     repo.FileRepo
	fileSearch   *search.FileSearch
	fileCache    *cache.FileCache
	fileGuard    *guard.FileGuard
	fileMapper   *FileMapper
	fileCoreSvc  *FileCoreService
	taskSvc      *TaskService
	snapshotRepo repo.SnapshotRepo
}

func NewFileCopyService() *FileCopyService {
	return &FileCopyService{
		fileRepo:     repo.NewFileRepo(),
		fileSearch:   search.NewFileSearch(),
		fileCache:    cache.NewFileCache(),
		fileGuard:    guard.NewFileGuard(),
		fileMapper:   NewFileMapper(),
		fileCoreSvc:  NewFileCoreService(),
		taskSvc:      NewTaskService(),
		snapshotRepo: repo.NewSnapshotRepo(),
	}
}

func (svc *FileCopyService) CopyOne(sourceID string, targetID string, userID string) (*File, error) {
	target, err := svc.fileCache.Get(targetID)
	if err != nil {
		return nil, err
	}
	source, err := svc.fileCache.Get(sourceID)
	if err != nil {
		return nil, err
	}
	task, err := svc.createTask(source, userID)
	if err != nil {
		return nil, err
	}
	defer func(taskID string) {
		if err := svc.taskSvc.deleteAndSync(taskID); err != nil {
			log.GetLogger().Error(err)
		}
	}(task.GetID())
	if err := svc.check(source, target, userID); err != nil {
		return nil, err
	}
	return svc.copy(source, target, userID)
}

func (svc *FileCopyService) copy(source model.File, target model.File, userID string) (*File, error) {
	tree, err := svc.getTree(source)
	if err != nil {
		return nil, err
	}
	root, clones, permissions, err := svc.cloneTree(source, target, tree, userID)
	if err != nil {
		return nil, err
	}
	if err := svc.persist(clones, permissions); err != nil {
		return nil, err
	}
	if err := svc.attachSnapshots(clones, tree); err != nil {
		return nil, err
	}
	svc.cache(clones, userID)
	go svc.index(clones)
	if err := svc.refreshUpdateTime(target); err != nil {
		return nil, err
	}
	res, err := svc.fileMapper.mapOne(root, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *FileCopyService) createTask(file model.File, userID string) (model.Task, error) {
	task, err := svc.taskSvc.insertAndSync(repo.TaskInsertOptions{
		ID:              helper.NewID(),
		Name:            "Copying.",
		UserID:          userID,
		IsIndeterminate: true,
		Status:          model.TaskStatusRunning,
		Payload:         map[string]string{repo.TaskPayloadObjectKey: file.GetName()},
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (svc *FileCopyService) check(source model.File, target model.File, userID string) error {
	if err := svc.fileGuard.Authorize(userID, target, model.PermissionEditor); err != nil {
		return err
	}
	if err := svc.fileGuard.Authorize(userID, source, model.PermissionEditor); err != nil {
		return err
	}
	if source.GetID() == target.GetID() {
		return errorpkg.NewFileCannotBeCopiedIntoItselfError(source)
	}
	if target.GetType() != model.FileTypeFolder {
		return errorpkg.NewFileIsNotAFolderError(target)
	}
	isGrandChild, err := svc.fileRepo.IsGrandChildOf(target.GetID(), source.GetID())
	if err != nil {
		return err
	}
	if isGrandChild {
		return errorpkg.NewFileCannotBeCopiedIntoOwnSubtreeError(source)
	}
	return nil
}

func (svc *FileCopyService) getTree(source model.File) ([]model.File, error) {
	var ids []string
	ids, err := svc.fileRepo.FindTreeIDs(source.GetID())
	if err != nil {
		return nil, err
	}
	var tree []model.File
	for _, id := range ids {
		sourceFile, err := svc.fileCache.Get(id)
		if err != nil {
			return nil, err
		}
		tree = append(tree, sourceFile)
	}
	return tree, nil
}

func (svc *FileCopyService) cloneTree(source model.File, target model.File, tree []model.File, userID string) (model.File, []model.File, []model.UserPermission, error) {
	var rootIndex int
	ids := make(map[string]string)
	var clones []model.File
	var permissions []model.UserPermission
	for index, leaf := range tree {
		clone := svc.newClone(leaf)
		if leaf.GetID() == source.GetID() {
			rootIndex = index
		}
		ids[leaf.GetID()] = clone.GetID()
		clones = append(clones, clone)
		permissions = append(permissions, svc.newUserPermission(clone, userID))
	}
	root := clones[rootIndex]
	for index, clone := range clones {
		id := ids[*clone.GetParentID()]
		clones[index].SetParentID(&id)
	}
	root.SetParentID(helper.ToPtr(target.GetID()))
	existing, err := svc.fileCoreSvc.GetChildWithName(target.GetID(), root.GetName())
	if err != nil {
		return nil, nil, nil, err
	}
	if existing != nil {
		root.SetName(helper.UniqueFilename(root.GetName()))
	}
	return root, clones, permissions, nil
}

func (svc *FileCopyService) newClone(file model.File) model.File {
	f := repo.NewFile()
	f.SetID(helper.NewID())
	f.SetParentID(file.GetParentID())
	f.SetWorkspaceID(file.GetWorkspaceID())
	f.SetSnapshotID(file.GetSnapshotID())
	f.SetType(file.GetType())
	f.SetName(file.GetName())
	f.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
	return f
}

func (svc *FileCopyService) newUserPermission(file model.File, userID string) model.UserPermission {
	p := repo.NewUserPermission()
	p.SetID(helper.NewID())
	p.SetUserID(userID)
	p.SetResourceID(file.GetID())
	p.SetPermission(model.PermissionOwner)
	p.SetCreateTime(time.Now().UTC().Format(time.RFC3339))
	return p
}

func (svc *FileCopyService) persist(clones []model.File, permissions []model.UserPermission) error {
	const BulkInsertChunkSize = 1000
	if err := svc.fileRepo.BulkInsert(clones, BulkInsertChunkSize); err != nil {
		return err
	}
	if err := svc.fileRepo.BulkInsertPermissions(permissions, BulkInsertChunkSize); err != nil {
		return err
	}
	return nil
}

func (svc *FileCopyService) attachSnapshots(clones []model.File, tree []model.File) error {
	const BulkInsertChunkSize = 1000
	var mappings []*repo.SnapshotFileEntity
	for index, clone := range clones {
		leaf := tree[index]
		if leaf.GetSnapshotID() != nil {
			mappings = append(mappings, &repo.SnapshotFileEntity{
				SnapshotID: *leaf.GetSnapshotID(),
				FileID:     clone.GetID(),
			})
		}
	}
	if err := svc.snapshotRepo.BulkMapWithFile(mappings, BulkInsertChunkSize); err != nil {
		return err
	}
	return nil
}

func (svc *FileCopyService) cache(clones []model.File, userID string) {
	for _, clone := range clones {
		if _, err := svc.fileCache.RefreshWithExisting(clone, userID); err != nil {
			log.GetLogger().Error(err)
		}
	}
}

func (svc *FileCopyService) index(clones []model.File) {
	if err := svc.fileSearch.Index(clones); err != nil {
		log.GetLogger().Error(err)
	}
}

func (svc *FileCopyService) refreshUpdateTime(target model.File) error {
	now := helper.NewTimestamp()
	target.SetUpdateTime(&now)
	if err := svc.fileRepo.Save(target); err != nil {
		return err
	}
	if err := svc.fileCoreSvc.Sync(target); err != nil {
		return err
	}
	return nil
}

type FileCopyManyOptions struct {
	SourceIDs []string `json:"sourceIds" validate:"required"`
	TargetID  string   `json:"targetId"  validate:"required"`
}

type FileCopyManyResult struct {
	New       []string `json:"new"`
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

func (svc *FileCopyService) CopyMany(opts FileCopyManyOptions, userID string) (*FileCopyManyResult, error) {
	res := &FileCopyManyResult{
		New:       make([]string, 0),
		Succeeded: make([]string, 0),
		Failed:    make([]string, 0),
	}
	for _, id := range opts.SourceIDs {
		file, err := svc.CopyOne(id, opts.TargetID, userID)
		if err != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.New = append(res.New, file.ID)
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}
