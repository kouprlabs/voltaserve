// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package repo

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type fileEntity struct {
	ID               string                  `gorm:"column:id"           json:"id"`
	WorkspaceID      string                  `gorm:"column:workspace_id" json:"workspaceId"`
	Name             string                  `gorm:"column:name"         json:"name"`
	Type             string                  `gorm:"column:type"         json:"type"`
	ParentID         *string                 `gorm:"column:parent_id"    json:"parentId,omitempty"`
	UserPermissions  []*UserPermissionValue  `gorm:"-"                   json:"userPermissions"`
	GroupPermissions []*GroupPermissionValue `gorm:"-"                   json:"groupPermissions"`
	Text             *string                 `gorm:"-"                   json:"text,omitempty"`
	SnapshotID       *string                 `gorm:"column:snapshot_id"  json:"snapshotId,omitempty"`
	CreateTime       string                  `gorm:"column:create_time"  json:"createTime"`
	UpdateTime       *string                 `gorm:"column:update_time"  json:"updateTime,omitempty"`
}

func (*fileEntity) TableName() string {
	return "file"
}

func (f *fileEntity) BeforeCreate(*gorm.DB) (err error) {
	f.CreateTime = helper.NewTimestamp()
	return nil
}

func (f *fileEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := helper.NewTimestamp()
	f.UpdateTime = &timeNow
	return nil
}

func (f *fileEntity) GetID() string {
	return f.ID
}

func (f *fileEntity) GetWorkspaceID() string {
	return f.WorkspaceID
}

func (f *fileEntity) GetName() string {
	return f.Name
}

func (f *fileEntity) GetType() string {
	return f.Type
}

func (f *fileEntity) GetParentID() *string {
	return f.ParentID
}

func (f *fileEntity) GetUserPermissions() []model.CoreUserPermission {
	var res []model.CoreUserPermission
	for _, p := range f.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (f *fileEntity) GetGroupPermissions() []model.CoreGroupPermission {
	var res []model.CoreGroupPermission
	for _, p := range f.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (f *fileEntity) GetText() *string {
	return f.Text
}

func (f *fileEntity) GetSnapshotID() *string {
	return f.SnapshotID
}

func (f *fileEntity) GetCreateTime() string {
	return f.CreateTime
}

func (f *fileEntity) GetUpdateTime() *string {
	return f.UpdateTime
}

func (f *fileEntity) SetID(id string) {
	f.ID = id
}

func (f *fileEntity) SetParentID(parentID *string) {
	f.ParentID = parentID
}

func (f *fileEntity) SetWorkspaceID(workspaceID string) {
	f.WorkspaceID = workspaceID
}

func (f *fileEntity) SetType(fileType string) {
	f.Type = fileType
}

func (f *fileEntity) SetName(name string) {
	f.Name = name
}

func (f *fileEntity) SetText(text *string) {
	f.Text = text
}

func (f *fileEntity) SetSnapshotID(snapshotID *string) {
	f.SnapshotID = snapshotID
}

func (f *fileEntity) SetUserPermissions(permissions []model.CoreUserPermission) {
	f.UserPermissions = make([]*UserPermissionValue, len(permissions))
	for i, p := range permissions {
		f.UserPermissions[i] = p.(*UserPermissionValue)
	}
}

func (f *fileEntity) SetGroupPermissions(permissions []model.CoreGroupPermission) {
	f.GroupPermissions = make([]*GroupPermissionValue, len(permissions))
	for i, p := range permissions {
		f.GroupPermissions[i] = p.(*GroupPermissionValue)
	}
}

func (f *fileEntity) SetCreateTime(createTime string) {
	f.CreateTime = createTime
}

func (f *fileEntity) SetUpdateTime(updateTime *string) {
	f.UpdateTime = updateTime
}

func NewFile() model.File {
	return &fileEntity{}
}

type NewFileOptions struct {
	ID               string
	WorkspaceID      string
	ParentID         *string
	Type             string
	Name             string
	Text             *string
	SnapshotID       *string
	UserPermissions  []model.CoreUserPermission
	GroupPermissions []model.CoreGroupPermission
	CreateTime       string
	UpdateTime       *string
}

func NewFileWithOptions(opts NewFileOptions) model.File {
	res := &fileEntity{
		ID:          opts.ID,
		WorkspaceID: opts.WorkspaceID,
		ParentID:    opts.ParentID,
		Type:        opts.Type,
		Name:        opts.Name,
		Text:        opts.Text,
		SnapshotID:  opts.SnapshotID,
		CreateTime:  opts.CreateTime,
		UpdateTime:  opts.UpdateTime,
	}
	res.SetUserPermissions(opts.UserPermissions)
	res.SetGroupPermissions(opts.GroupPermissions)
	return res
}

type FileRepo struct {
	db             *gorm.DB
	permissionRepo *PermissionRepo
}

func NewFileRepo() *FileRepo {
	return &FileRepo{
		db:             infra.NewPostgresManager().GetDBOrPanic(),
		permissionRepo: NewPermissionRepo(),
	}
}

type FileInsertOptions struct {
	Name        string
	WorkspaceID string
	ParentID    string
	Type        string
}

func (repo *FileRepo) Insert(opts FileInsertOptions) (model.File, error) {
	id := helper.NewID()
	var parentID *string
	if opts.ParentID != "" {
		parentID = &opts.ParentID
	}
	file := fileEntity{
		ID:          id,
		WorkspaceID: opts.WorkspaceID,
		Name:        opts.Name,
		Type:        opts.Type,
		ParentID:    parentID,
	}
	if db := repo.db.Create(&file); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*fileEntity{res}); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *FileRepo) Find(id string) (model.File, error) {
	file, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*fileEntity{file}); err != nil {
		return nil, err
	}
	return file, nil
}

func (repo *FileRepo) find(id string) (*fileEntity, error) {
	res := fileEntity{}
	db := repo.db.
		Raw("SELECT * FROM file WHERE id = ?", id).
		Scan(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewFileNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	if len(res.ID) == 0 {
		return nil, errorpkg.NewFileNotFoundError(db.Error)
	}
	return &res, nil
}

func (repo *FileRepo) FindChildren(id string) ([]model.File, error) {
	var entities []*fileEntity
	db := repo.db.
		Raw("SELECT * FROM file WHERE parent_id = ? ORDER BY create_time", id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.File
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *FileRepo) FindPath(id string) ([]model.File, error) {
	var entities []*fileEntity
	if db := repo.db.
		Raw(`WITH RECURSIVE rec (id, name, type, parent_id, workspace_id, create_time, update_time) AS
             (SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.create_time, f.update_time FROM file f WHERE f.id = ?
             UNION SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.create_time, f.update_time FROM rec, file f WHERE f.id = rec.parent_id)
             SELECT * FROM rec`,
			id).
		Scan(&entities); db.Error != nil {
		return nil, db.Error
	}
	if err := repo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.File
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *FileRepo) FindTree(id string) ([]model.File, error) {
	var entities []*fileEntity
	db := repo.db.
		Raw(`WITH RECURSIVE rec (id, name, type, parent_id, workspace_id, snapshot_id, create_time, update_time) AS
             (SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.snapshot_id, f.create_time, f.update_time FROM file f WHERE f.id = ?
             UNION SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.snapshot_id, f.create_time, f.update_time FROM rec, file f WHERE f.parent_id = rec.id)
             SELECT rec.* FROM rec ORDER BY create_time ASC`,
			id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.File
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *FileRepo) FindTreeIDs(id string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw(`WITH RECURSIVE rec (id, parent_id, create_time) AS
             (SELECT f.id, f.parent_id, f.create_time FROM file f WHERE f.id = ?
             UNION SELECT f.id, f.parent_id, f.create_time FROM rec, file f WHERE f.parent_id = rec.id)
             SELECT rec.id as result FROM rec ORDER BY create_time ASC`,
			id).
		Scan(&values)
	if db.Error != nil {
		return nil, db.Error
	}
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *FileRepo) DeleteChunk(ids []string) error {
	if db := repo.db.Delete(&fileEntity{}, ids); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) Count() (int64, error) {
	var count int64
	db := repo.db.Model(&fileEntity{}).Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *FileRepo) FindIDsByWorkspace(workspaceID string) ([]string, error) {
	type IDResult struct {
		Result string
	}
	var ids []IDResult
	db := repo.db.
		Raw("SELECT id result FROM file WHERE workspace_id = ? ORDER BY create_time", workspaceID).
		Scan(&ids)
	if db.Error != nil {
		return nil, db.Error
	}
	res := make([]string, 0)
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *FileRepo) FindIDsBySnapshot(snapshotID string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw("SELECT file_id result FROM snapshot_file WHERE snapshot_id = ?", snapshotID).
		Scan(&values)
	if db.Error != nil {
		return nil, db.Error
	}
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *FileRepo) MoveSourceIntoTarget(targetID string, sourceID string) error {
	if db := repo.db.Exec("UPDATE file SET parent_id = ? WHERE id = ?", targetID, sourceID); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) Save(file model.File) error {
	if db := repo.db.Save(file); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) BulkInsert(values []model.File, chunkSize int) error {
	var entities []*fileEntity
	for _, f := range values {
		entities = append(entities, f.(*fileEntity))
	}
	if db := repo.db.CreateInBatches(entities, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) BulkInsertPermissions(values []model.UserPermission, chunkSize int) error {
	var entities []*userPermissionEntity
	for _, p := range values {
		entities = append(entities, p.(*userPermissionEntity))
	}
	if db := repo.db.CreateInBatches(entities, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM file WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	db = repo.db.Exec("DELETE FROM userpermission WHERE resource_id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	db = repo.db.Exec("DELETE FROM grouppermission WHERE resource_id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) FindChildrenIDs(id string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.
		Raw("SELECT id result FROM file WHERE parent_id = ? ORDER BY create_time", id).
		Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := make([]string, 0)
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *FileRepo) CountChildren(id string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw("SELECT count(*) as result FROM file WHERE parent_id = ?", id).
		Scan(&res)
	if db.Error != nil {
		return -1, db.Error
	}
	return res.Result, nil
}

func (repo *FileRepo) CountItems(id string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw(`WITH RECURSIVE rec (id, parent_id) AS
             (SELECT f.id, f.parent_id FROM file f WHERE f.id = ?
             UNION SELECT f.id, f.parent_id FROM rec, file f WHERE f.parent_id = rec.id)
             SELECT count(rec.id) as result FROM rec`,
			id).
		Scan(&res)
	if db.Error != nil {
		return -1, db.Error
	}
	return res.Result - 1, nil
}

func (repo *FileRepo) IsGrandChildOf(id string, ancestorID string) (bool, error) {
	type Result struct {
		Result bool
	}
	var res Result
	if db := repo.db.
		Raw(`WITH RECURSIVE rec (id, parent_id) AS
             (SELECT f.id, f.parent_id FROM file f WHERE f.id = ?
             UNION SELECT f.id, f.parent_id FROM rec JOIN file f ON f.id = rec.parent_id)
             SELECT count(rec.id) > 0 as result FROM rec WHERE rec.id = ?`,
			id, ancestorID).
		Scan(&res); db.Error != nil {
		return false, db.Error
	}
	return res.Result, nil
}

func (repo *FileRepo) ComputeSize(id string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw(`WITH RECURSIVE rec (id, parent_id) AS
             (SELECT f.id, f.parent_id FROM file f WHERE f.id = ?
             UNION SELECT f.id, f.parent_id FROM rec, file f WHERE f.parent_id = rec.id)
             SELECT coalesce(sum((s.original->>'size')::bigint), 0) as result FROM snapshot s, rec
             LEFT JOIN snapshot_file map ON rec.id = map.file_id WHERE map.snapshot_id = s.id`,
			id).
		Scan(&res)
	if db.Error != nil {
		return res.Result, db.Error
	}
	return res.Result, nil
}

func (repo *FileRepo) ClearSnapshotID(id string) error {
	if db := repo.db.Exec("UPDATE file SET snapshot_id = NULL WHERE id = ?", id); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) GrantUserPermission(id string, userID string, permission string) error {
	// Grant 'viewer' permission to workspace
	db := repo.db.
		Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
              (SELECT ?, ?, w.id, 'viewer', ? FROM file f
              INNER JOIN workspace w ON w.id = f.workspace_id AND f.id = ?)
              ON CONFLICT DO NOTHING`,
			helper.NewID(), userID, helper.NewTimestamp(), id)
	if db.Error != nil {
		return db.Error
	}

	// Grant 'viewer' permission to path files
	path, err := repo.FindPath(id)
	if err != nil {
		return err
	}
	for _, f := range path {
		db := repo.db.
			Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
                  VALUES (?, ?, ?, 'viewer', ?) ON CONFLICT DO NOTHING`,
				helper.NewID(), userID, f.GetID(), helper.NewTimestamp())
		if db.Error != nil {
			return db.Error
		}
	}

	// Grant the requested permission to tree files
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.
			Exec(`INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
                  VALUES (?, ?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?`,
				helper.NewID(), userID, f.GetID(), permission, helper.NewTimestamp(), permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *FileRepo) RevokeUserPermission(tree []model.File, userID string) error {
	for _, f := range tree {
		db := repo.db.
			Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?",
				userID, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *FileRepo) GrantGroupPermission(id string, groupID string, permission string) error {
	// Grant permission to workspace
	db := repo.db.
		Exec(`INSERT INTO grouppermission (id, group_id, resource_id, permission, create_time)
              (SELECT ?, ?, w.id, 'viewer', ? FROM file f
              INNER JOIN workspace w ON w.id = f.workspace_id AND f.id = ?)
              ON CONFLICT DO NOTHING`,
			helper.NewID(), groupID, helper.NewTimestamp(), id)
	if db.Error != nil {
		return db.Error
	}

	// Grant 'viewer' permission to path files
	path, err := repo.FindPath(id)
	if err != nil {
		return err
	}
	for _, f := range path {
		db := repo.db.
			Exec(`INSERT INTO grouppermission (id, group_id, resource_id, permission, create_time)
                  VALUES (?, ?, ?, 'viewer', ?) ON CONFLICT DO NOTHING`,
				helper.NewID(), groupID, f.GetID(), helper.NewTimestamp())
		if db.Error != nil {
			return db.Error
		}
	}

	// Grant the requested permission to tree files
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.
			Exec(`INSERT INTO grouppermission (id, group_id, resource_id, permission, create_time)
                  VALUES (?, ?, ?, ?, ?) ON CONFLICT (group_id, resource_id) DO UPDATE SET permission = ?`,
				helper.NewID(), groupID, f.GetID(), permission, helper.NewTimestamp(), permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *FileRepo) RevokeGroupPermission(tree []model.File, groupID string) error {
	for _, f := range tree {
		db := repo.db.
			Exec("DELETE FROM grouppermission WHERE group_id = ? AND resource_id = ?",
				groupID, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *FileRepo) PopulateModelFieldsForUser(files []model.File, userID string) error {
	for _, f := range files {
		userPermissions := make([]model.CoreUserPermission, 0)
		userPermissions = append(userPermissions, &UserPermissionValue{
			UserID: userID,
			Value:  model.PermissionOwner,
		})
		f.SetUserPermissions(userPermissions)
		f.SetGroupPermissions(make([]model.CoreGroupPermission, 0))
	}
	return nil
}

func (repo *FileRepo) populateModelFields(entities []*fileEntity) error {
	for _, f := range entities {
		f.UserPermissions = make([]*UserPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.FindUserPermissions(f.ID)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			f.UserPermissions = append(f.UserPermissions, &UserPermissionValue{
				UserID: p.GetUserID(),
				Value:  p.GetPermission(),
			})
		}
		f.GroupPermissions = make([]*GroupPermissionValue, 0)
		groupPermissions, err := repo.permissionRepo.FindGroupPermissions(f.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			f.GroupPermissions = append(f.GroupPermissions, &GroupPermissionValue{
				GroupID: p.GetGroupID(),
				Value:   p.GetPermission(),
			})
		}
	}
	return nil
}
