package repo

import (
	"errors"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/helper"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type FileInsertOptions struct {
	Name        string
	WorkspaceID string
	ParentID    *string
	Type        string
}

type FileRepo interface {
	Insert(opts FileInsertOptions) (model.File, error)
	Find(id string) (model.File, error)
	FindChildren(id string) ([]model.File, error)
	FindPath(id string) ([]model.File, error)
	FindTree(id string) ([]model.File, error)
	GetIDsByWorkspace(workspaceID string) ([]string, error)
	MoveSourceIntoTarget(targetID string, sourceID string) error
	Save(file model.File) error
	BulkInsert(values []model.File, chunkSize int) error
	BulkInsertPermissions(values []model.UserPermission, chunkSize int) error
	Delete(id string) error
	GetChildrenIDs(id string) ([]string, error)
	GetItemCount(id string) (int64, error)
	IsGrandChildOf(id string, ancestorID string) (bool, error)
	GetSize(id string) (int64, error)
	GrantUserPermission(id string, userID string, permission string) error
	RevokeUserPermission(tree []model.File, userID string) error
	GrantGroupPermission(id string, groupID string, permission string) error
	RevokeGroupPermission(tree []model.File, groupID string) error
}

func NewFileRepo() FileRepo {
	return newFileRepo()
}

func NewFile() model.File {
	return &fileEntity{}
}

type fileEntity struct {
	ID               string                  `json:"id" gorm:"column:id"`
	WorkspaceID      string                  `json:"workspaceId" gorm:"column:workspace_id"`
	Name             string                  `json:"name" gorm:"column:name"`
	Type             string                  `json:"type" gorm:"column:type"`
	ParentID         *string                 `json:"parentId,omitempty" gorm:"column:parent_id"`
	UserPermissions  []*UserPermissionValue  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*GroupPermissionValue `json:"groupPermissions" gorm:"-"`
	Text             *string                 `json:"text,omitempty" gorm:"-"`
	SnapshotID       *string                 `json:"snapshotId,omitempty" gorm:"column:snapshot_id"`
	CreateTime       string                  `json:"createTime" gorm:"column:create_time"`
	UpdateTime       *string                 `json:"updateTime,omitempty" gorm:"column:update_time"`
}

func (*fileEntity) TableName() string {
	return "file"
}

func (f *fileEntity) BeforeCreate(*gorm.DB) (err error) {
	f.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (f *fileEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
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

func (f *fileEntity) SetCreateTime(createTime string) {
	f.CreateTime = createTime
}

func (f *fileEntity) SetUpdateTime(updateTime *string) {
	f.UpdateTime = updateTime
}

type fileRepo struct {
	db             *gorm.DB
	permissionRepo *permissionRepo
}

func newFileRepo() *fileRepo {
	return &fileRepo{
		db:             infra.NewPostgresManager().GetDBOrPanic(),
		permissionRepo: newPermissionRepo(),
	}
}

func (repo *fileRepo) Insert(opts FileInsertOptions) (model.File, error) {
	id := helper.NewID()
	file := fileEntity{
		ID:          id,
		WorkspaceID: opts.WorkspaceID,
		Name:        opts.Name,
		Type:        opts.Type,
		ParentID:    opts.ParentID,
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

func (repo *fileRepo) Find(id string) (model.File, error) {
	file, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*fileEntity{file}); err != nil {
		return nil, err
	}
	return file, nil
}

func (repo *fileRepo) find(id string) (*fileEntity, error) {
	var res = fileEntity{}
	db := repo.db.Raw("SELECT * FROM file WHERE id = ?", id).Scan(&res)
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

func (repo *fileRepo) FindChildren(id string) ([]model.File, error) {
	var entities []*fileEntity
	db := repo.db.Raw("SELECT * FROM file WHERE parent_id = ? ORDER BY create_time ASC", id).Scan(&entities)
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

func (repo *fileRepo) FindPath(id string) ([]model.File, error) {
	var entities []*fileEntity
	if db := repo.db.
		Raw("WITH RECURSIVE rec (id, name, type, parent_id, workspace_id, create_time, update_time) AS "+
			"(SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.create_time, f.update_time FROM file f WHERE f.id = ? "+
			"UNION SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.create_time, f.update_time FROM rec, file f WHERE f.id = rec.parent_id) "+
			"SELECT * FROM rec", id).
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

func (repo *fileRepo) FindTree(id string) ([]model.File, error) {
	var entities []*fileEntity
	db := repo.db.
		Raw("WITH RECURSIVE rec (id, name, type, parent_id, workspace_id, snapshot_id, create_time, update_time) AS "+
			"(SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.snapshot_id, f.create_time, f.update_time FROM file f WHERE f.id = ? "+
			"UNION SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.snapshot_id, f.create_time, f.update_time FROM rec, file f WHERE f.parent_id = rec.id) "+
			"SELECT rec.* FROM rec ORDER BY create_time ASC", id).
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

func (repo *fileRepo) GetIDsByWorkspace(workspaceID string) ([]string, error) {
	type IDResult struct {
		Result string
	}
	var ids []IDResult
	db := repo.db.Raw("SELECT id result FROM file WHERE workspace_id = ? ORDER BY create_time ASC", workspaceID).Scan(&ids)
	if db.Error != nil {
		return nil, db.Error
	}
	res := []string{}
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *fileRepo) MoveSourceIntoTarget(targetID string, sourceID string) error {
	if db := repo.db.Exec("UPDATE file SET parent_id = ? WHERE id = ?", targetID, sourceID); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *fileRepo) Save(file model.File) error {
	if db := repo.db.Save(file); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *fileRepo) BulkInsert(values []model.File, chunkSize int) error {
	var entities []*fileEntity
	for _, f := range values {
		entities = append(entities, f.(*fileEntity))
	}
	if db := repo.db.CreateInBatches(entities, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *fileRepo) BulkInsertPermissions(values []model.UserPermission, chunkSize int) error {
	var entities []*userPermissionEntity
	for _, p := range values {
		entities = append(entities, p.(*userPermissionEntity))
	}
	if db := repo.db.CreateInBatches(entities, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *fileRepo) Delete(id string) error {
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

func (repo *fileRepo) GetChildrenIDs(id string) ([]string, error) {
	type Value struct {
		Result string
	}
	var values []Value
	db := repo.db.Raw("SELECT id result FROM file WHERE parent_id = ? ORDER BY create_time ASC", id).Scan(&values)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range values {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *fileRepo) GetItemCount(id string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw("WITH RECURSIVE rec (id, parent_id) AS "+
			"(SELECT f.id, f.parent_id FROM file f WHERE f.id = ? "+
			"UNION SELECT f.id, f.parent_id FROM rec, file f WHERE f.parent_id = rec.id) "+
			"SELECT count(rec.id) as result FROM rec", id).
		Scan(&res)
	if db.Error != nil {
		return 0, db.Error
	}
	return res.Result - 1, nil
}

func (repo *fileRepo) IsGrandChildOf(id string, ancestorID string) (bool, error) {
	type Result struct {
		Result bool
	}
	var res Result
	if db := repo.db.
		Raw("WITH RECURSIVE rec (id, parent_id) AS "+
			"(SELECT f.id, f.parent_id FROM file f WHERE f.id = ? "+
			"UNION SELECT f.id, f.parent_id FROM rec, file f WHERE f.parent_id = rec.id) "+
			"SELECT count(rec.id) > 0 as result FROM rec WHERE rec.id = ?", ancestorID, id).
		Scan(&res); db.Error != nil {
		return false, db.Error
	}
	return res.Result, nil
}

func (repo *fileRepo) GetSize(id string) (int64, error) {
	type Result struct {
		Result int64
	}
	var res Result
	db := repo.db.
		Raw("WITH RECURSIVE rec (id, parent_id) AS "+
			"(SELECT f.id, f.parent_id FROM file f WHERE f.id = ? "+
			"UNION SELECT f.id, f.parent_id FROM rec, file f WHERE f.parent_id = rec.id) "+
			"SELECT coalesce(sum((s.original->>'size')::int), 0) as result FROM snapshot s, rec "+
			"LEFT JOIN snapshot_file map ON rec.id = map.file_id WHERE map.snapshot_id = s.id", id).
		Scan(&res)
	if db.Error != nil {
		return res.Result, db.Error
	}
	return res.Result, nil
}

func (repo *fileRepo) GrantUserPermission(id string, userID string, permission string) error {
	/* Grant permission to workspace */
	db := repo.db.Exec("INSERT INTO userpermission (id, user_id, resource_id, permission) "+
		"(SELECT ?, ?, w.id, 'viewer' FROM file f "+
		"INNER JOIN workspace w ON w.id = f.workspace_id AND f.id = ?) "+
		"ON CONFLICT DO NOTHING",
		helper.NewID(), userID, id)
	if db.Error != nil {
		return db.Error
	}

	/* Grant 'viewer' permission to path files */
	path, err := repo.FindPath(id)
	if err != nil {
		return err
	}
	for _, f := range path {
		db := repo.db.Exec("INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, 'viewer') ON CONFLICT DO NOTHING",
			helper.NewID(), userID, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}

	/* Grant the requested permission to tree files */
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.Exec("INSERT INTO userpermission (id, user_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = ?",
			helper.NewID(), userID, f.GetID(), permission, permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *fileRepo) RevokeUserPermission(tree []model.File, userID string) error {
	for _, f := range tree {
		db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userID, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *fileRepo) GrantGroupPermission(id string, groupID string, permission string) error {
	/* Grant permission to workspace */
	db := repo.db.Exec("INSERT INTO grouppermission (id, group_id, resource_id, permission) "+
		"(SELECT ?, ?, w.id, 'viewer' FROM file f "+
		"INNER JOIN workspace w ON w.id = f.workspace_id AND f.id = ?) "+
		"ON CONFLICT DO NOTHING",
		helper.NewID(), groupID, id)
	if db.Error != nil {
		return db.Error
	}

	/* Grant 'viewer' permission to path files */
	path, err := repo.FindPath(id)
	if err != nil {
		return err
	}
	for _, f := range path {
		db := repo.db.Exec("INSERT INTO grouppermission (id, group_id, resource_id, permission) "+
			"VALUES (?, ?, ?, 'viewer') ON CONFLICT DO NOTHING",
			helper.NewID(), groupID, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}

	/* Grant the requested permission to tree files */
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.Exec("INSERT INTO grouppermission (id, group_id, resource_id, permission) "+
			"VALUES (?, ?, ?, ?) ON CONFLICT (group_id, resource_id) DO UPDATE SET permission = ?",
			helper.NewID(), groupID, f.GetID(), permission, permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *fileRepo) RevokeGroupPermission(tree []model.File, groupID string) error {
	for _, f := range tree {
		db := repo.db.Exec("DELETE FROM grouppermission WHERE group_id = ? AND resource_id = ?", groupID, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *fileRepo) populateModelFields(entities []*fileEntity) error {
	for _, f := range entities {
		f.UserPermissions = make([]*UserPermissionValue, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(f.ID)
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
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(f.ID)
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
