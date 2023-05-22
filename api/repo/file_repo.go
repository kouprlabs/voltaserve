package repo

import (
	"errors"
	"time"
	"voltaserve/errorpkg"
	"voltaserve/helpers"
	"voltaserve/infra"
	"voltaserve/model"

	"gorm.io/gorm"
)

type FileEntity struct {
	Id               string                   `json:"id"`
	WorkspaceId      string                   `json:"workspaceId"`
	Name             string                   `json:"name"`
	Type             string                   `json:"type"`
	ParentId         *string                  `json:"parentId,omitempty"`
	Snapshots        []*SnapshotEntity        `json:"snapshots,omitempty" gorm:"-"`
	UserPermissions  []*model.UserPermission  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*model.GroupPermission `json:"groupPermissions" gorm:"-"`
	Text             *string                  `json:"text,omitempty" gorm:"-"`
	CreateTime       string                   `json:"createTime"`
	UpdateTime       *string                  `json:"updateTime,omitempty"`
}

func (FileEntity) TableName() string {
	return "file"
}

func (i *FileEntity) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (i *FileEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	i.UpdateTime = &timeNow
	return nil
}

func (i FileEntity) GetId() string {
	return i.Id
}

func (i FileEntity) GetWorkspaceId() string {
	return i.WorkspaceId
}

func (i FileEntity) GetName() string {
	return i.Name
}

func (i FileEntity) GetType() string {
	return i.Type
}

func (i FileEntity) GetParentId() *string {
	return i.ParentId
}

func (i FileEntity) GetSnapshots() []model.SnapshotModel {
	var res []model.SnapshotModel
	for _, s := range i.Snapshots {
		res = append(res, s)
	}
	return res
}

func (i FileEntity) GetUserPermissions() []model.UserPermissionModel {
	var res []model.UserPermissionModel
	for _, p := range i.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (i FileEntity) GetGroupPermissions() []model.GroupPermissionModel {
	var res []model.GroupPermissionModel
	for _, p := range i.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (i FileEntity) GetText() *string {
	return i.Text
}

func (i FileEntity) GetCreateTime() string {
	return i.CreateTime
}

func (i FileEntity) GetUpdateTime() *string {
	return i.UpdateTime
}

func (i *FileEntity) SetId(id string) {
	i.Id = id
}

func (i *FileEntity) SetParentId(parentId *string) {
	i.ParentId = parentId
}

func (i *FileEntity) SetWorkspaceId(workspaceId string) {
	i.WorkspaceId = workspaceId
}

func (i *FileEntity) SetType(fileType string) {
	i.Type = fileType
}

func (i *FileEntity) SetName(name string) {
	i.Name = name
}

func (i *FileEntity) SetText(text *string) {
	i.Text = text
}

func (i *FileEntity) SetCreateTime(createTime string) {
	i.CreateTime = createTime
}

func (i *FileEntity) SetUpdateTime(updateTime *string) {
	i.UpdateTime = updateTime
}

type FileRepo struct {
	db             *gorm.DB
	snapshotRepo   *SnapshotRepo
	permissionRepo *PermissionRepo
}

func NewFileRepo() *FileRepo {
	return &FileRepo{
		db:             infra.GetDb(),
		snapshotRepo:   NewSnapshotRepo(),
		permissionRepo: NewPermissionRepo(),
	}
}

func (repo *FileRepo) New() model.FileModel {
	return &FileEntity{}
}

type FileInsertOptions struct {
	Name        string
	WorkspaceId string
	ParentId    *string
	Type        string
}

func (repo *FileRepo) Insert(opts FileInsertOptions) (model.FileModel, error) {
	id := helpers.NewId()
	file := FileEntity{
		Id:          id,
		WorkspaceId: opts.WorkspaceId,
		Name:        opts.Name,
		Type:        opts.Type,
		ParentId:    opts.ParentId,
	}
	if db := repo.db.Save(&file); db.Error != nil {
		return nil, db.Error
	}
	res, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*FileEntity{res}); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *FileRepo) Find(id string) (model.FileModel, error) {
	file, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*FileEntity{file}); err != nil {
		return nil, err
	}
	return file, nil
}

func (repo *FileRepo) find(id string) (*FileEntity, error) {
	var res = FileEntity{}
	db := repo.db.Raw("SELECT * FROM file WHERE id = ?", id).Scan(&res)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewFileNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	if len(res.Id) == 0 {
		return nil, errorpkg.NewFileNotFoundError(db.Error)
	}
	return &res, nil
}

func (repo *FileRepo) FindChildren(id string) ([]model.FileModel, error) {
	var entities []*FileEntity
	db := repo.db.Raw("SELECT * FROM file WHERE parent_id = ? ORDER BY create_time ASC", id).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.FileModel
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *FileRepo) FindPath(id string) ([]model.FileModel, error) {
	var entities []*FileEntity
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
	var res []model.FileModel
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *FileRepo) FindTree(id string) ([]model.FileModel, error) {
	var entities []*FileEntity
	db := repo.db.
		Raw("WITH RECURSIVE rec (id, name, type, parent_id, workspace_id, create_time, update_time) AS "+
			"(SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.create_time, f.update_time FROM file f WHERE f.id = ? "+
			"UNION SELECT f.id, f.name, f.type, f.parent_id, f.workspace_id, f.create_time, f.update_time FROM rec, file f WHERE f.parent_id = rec.id) "+
			"SELECT rec.* FROM rec ORDER BY create_time ASC", id).
		Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.FileModel
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *FileRepo) GetIdsByWorkspace(workspaceId string) ([]string, error) {
	type IdResult struct {
		Result string
	}
	var ids []IdResult
	db := repo.db.Raw("SELECT id result FROM file WHERE workspace_id = ? ORDER BY create_time ASC", workspaceId).Scan(&ids)
	if db.Error != nil {
		return nil, db.Error
	}
	res := []string{}
	for _, id := range ids {
		res = append(res, id.Result)
	}
	return res, nil
}

func (repo *FileRepo) AssignSnapshots(cloneId string, originalId string) error {
	if db := repo.db.Exec("INSERT INTO snapshot_file (snapshot_id, file_id) SELECT s.id, ? "+
		"FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id "+
		"WHERE map.file_id = ? ORDER BY s.version DESC LIMIT 1", cloneId, originalId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) MoveSourceIntoTarget(targetId string, sourceId string) error {
	if db := repo.db.Exec("UPDATE file SET parent_id = ? WHERE id = ?", targetId, sourceId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) Save(file model.FileModel) error {
	if db := repo.db.Save(file); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) BulkInsert(values []model.FileModel, chunkSize int) error {
	var entities []*FileEntity
	for _, f := range values {
		entities = append(entities, f.(*FileEntity))
	}
	if db := repo.db.CreateInBatches(entities, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *FileRepo) BulkInsertPermissions(values []*UserPermission, chunkSize int) error {
	if db := repo.db.CreateInBatches(values, chunkSize); db.Error != nil {
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

func (repo *FileRepo) GetChildrenIds(id string) ([]string, error) {
	type Result struct {
		Result string
	}
	var results []Result
	db := repo.db.Raw("SELECT id result FROM file WHERE parent_id = ? ORDER BY create_time ASC", id).Scan(&results)
	if db.Error != nil {
		return []string{}, db.Error
	}
	res := []string{}
	for _, v := range results {
		res = append(res, v.Result)
	}
	return res, nil
}

func (repo *FileRepo) GetItemCount(id string) (int64, error) {
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

func (repo *FileRepo) IsGrandChildOf(id string, ancestorId string) (bool, error) {
	type Result struct {
		Result bool
	}
	var res Result
	if db := repo.db.
		Raw("WITH RECURSIVE rec (id, parent_id) AS "+
			"(SELECT f.id, f.parent_id FROM file f WHERE f.id = ? "+
			"UNION SELECT f.id, f.parent_id FROM rec, file f WHERE f.parent_id = rec.id) "+
			"SELECT count(rec.id) > 0 as result FROM rec WHERE rec.id = ?", ancestorId, id).
		Scan(&res); db.Error != nil {
		return false, db.Error
	}
	return res.Result, nil
}

func (repo *FileRepo) GetSize(id string) (int64, error) {
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

func (repo *FileRepo) GrantUserPermission(id string, userId string, permission string) error {
	/* Grant permission to workspace */
	db := repo.db.Exec("INSERT INTO userpermission (id, user_id, resource_id, permission) "+
		"(SELECT ?, ?, w.id, 'viewer' FROM file f "+
		"INNER JOIN workspace w ON w.id = f.workspace_id AND f.id = ?) "+
		"ON CONFLICT DO NOTHING",
		helpers.NewId(), userId, id)
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
			helpers.NewId(), userId, f.GetId())
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
			helpers.NewId(), userId, f.GetId(), permission, permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *FileRepo) RevokeUserPermission(id string, userId string) error {
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userId, f.GetId())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *FileRepo) GrantGroupPermission(id string, groupId string, permission string) error {
	/* Grant permission to workspace */
	db := repo.db.Exec("INSERT INTO grouppermission (id, group_id, resource_id, permission) "+
		"(SELECT ?, ?, w.id, 'viewer' FROM file f "+
		"INNER JOIN workspace w ON w.id = f.workspace_id AND f.id = ?) "+
		"ON CONFLICT DO NOTHING",
		helpers.NewId(), groupId, id)
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
			helpers.NewId(), groupId, f.GetId())
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
			helpers.NewId(), groupId, f.GetId(), permission, permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *FileRepo) RevokeGroupPermission(id string, groupId string) error {
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.Exec("DELETE FROM grouppermission WHERE group_id = ? AND resource_id = ?", groupId, f.GetId())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *FileRepo) populateModelFields(entities []*FileEntity) error {
	for _, f := range entities {
		f.UserPermissions = make([]*model.UserPermission, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(f.Id)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			f.UserPermissions = append(f.UserPermissions, &model.UserPermission{
				UserId: p.UserId,
				Value:  p.Permission,
			})
		}
		f.GroupPermissions = make([]*model.GroupPermission, 0)
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(f.Id)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			f.GroupPermissions = append(f.GroupPermissions, &model.GroupPermission{
				GroupId: p.GroupId,
				Value:   p.Permission,
			})
		}
		snapshots, err := repo.snapshotRepo.FindAllForFile(f.Id)
		if err != nil {
			return nil
		}
		f.Snapshots = snapshots
	}
	return nil
}
