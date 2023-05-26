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

type FileInsertOptions struct {
	Name        string
	WorkspaceId string
	ParentId    *string
	Type        string
}

type FileRepo interface {
	New() model.CoreFile
	Insert(opts FileInsertOptions) (model.CoreFile, error)
	Find(id string) (model.CoreFile, error)
	FindChildren(id string) ([]model.CoreFile, error)
	FindPath(id string) ([]model.CoreFile, error)
	FindTree(id string) ([]model.CoreFile, error)
	GetIdsByWorkspace(workspaceId string) ([]string, error)
	AssignSnapshots(cloneId string, originalId string) error
	MoveSourceIntoTarget(targetId string, sourceId string) error
	Save(file model.CoreFile) error
	BulkInsert(values []model.CoreFile, chunkSize int) error
	BulkInsertPermissions(values []*UserPermission, chunkSize int) error
	Delete(id string) error
	GetChildrenIDs(id string) ([]string, error)
	GetItemCount(id string) (int64, error)
	IsGrandChildOf(id string, ancestorId string) (bool, error)
	GetSize(id string) (int64, error)
	GrantUserPermission(id string, userId string, permission string) error
	RevokeUserPermission(id string, userId string) error
	GrantGroupPermission(id string, groupId string, permission string) error
	RevokeGroupPermission(id string, groupId string) error
}

func NewFileRepo() FileRepo {
	return NewPostgresFileRepo()
}

type PostgresFile struct {
	ID               string                   `json:"id"`
	WorkspaceId      string                   `json:"workspaceId"`
	Name             string                   `json:"name"`
	Type             string                   `json:"type"`
	ParentId         *string                  `json:"parentId,omitempty"`
	Snapshots        []*PostgresSnapshot      `json:"snapshots,omitempty" gorm:"-"`
	UserPermissions  []*model.UserPermission  `json:"userPermissions" gorm:"-"`
	GroupPermissions []*model.GroupPermission `json:"groupPermissions" gorm:"-"`
	Text             *string                  `json:"text,omitempty" gorm:"-"`
	CreateTime       string                   `json:"createTime"`
	UpdateTime       *string                  `json:"updateTime,omitempty"`
}

func (PostgresFile) TableName() string {
	return "file"
}

func (i *PostgresFile) BeforeCreate(tx *gorm.DB) (err error) {
	i.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (i *PostgresFile) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	i.UpdateTime = &timeNow
	return nil
}

func (i PostgresFile) GetID() string {
	return i.ID
}

func (i PostgresFile) GetWorkspaceID() string {
	return i.WorkspaceId
}

func (i PostgresFile) GetName() string {
	return i.Name
}

func (i PostgresFile) GetType() string {
	return i.Type
}

func (i PostgresFile) GetParentID() *string {
	return i.ParentId
}

func (i PostgresFile) GetSnapshots() []model.CoreSnapshot {
	var res []model.CoreSnapshot
	for _, s := range i.Snapshots {
		res = append(res, s)
	}
	return res
}

func (i PostgresFile) GetUserPermissions() []model.CoreUserPermission {
	var res []model.CoreUserPermission
	for _, p := range i.UserPermissions {
		res = append(res, p)
	}
	return res
}

func (i PostgresFile) GetGroupPermissions() []model.CoreGroupPermission {
	var res []model.CoreGroupPermission
	for _, p := range i.GroupPermissions {
		res = append(res, p)
	}
	return res
}

func (i PostgresFile) GetText() *string {
	return i.Text
}

func (i PostgresFile) GetCreateTime() string {
	return i.CreateTime
}

func (i PostgresFile) GetUpdateTime() *string {
	return i.UpdateTime
}

func (i *PostgresFile) SetID(id string) {
	i.ID = id
}

func (i *PostgresFile) SetParentID(parentId *string) {
	i.ParentId = parentId
}

func (i *PostgresFile) SetWorkspaceID(workspaceId string) {
	i.WorkspaceId = workspaceId
}

func (i *PostgresFile) SetType(fileType string) {
	i.Type = fileType
}

func (i *PostgresFile) SetName(name string) {
	i.Name = name
}

func (i *PostgresFile) SetText(text *string) {
	i.Text = text
}

func (i *PostgresFile) SetCreateTime(createTime string) {
	i.CreateTime = createTime
}

func (i *PostgresFile) SetUpdateTime(updateTime *string) {
	i.UpdateTime = updateTime
}

type PostgresFileRepo struct {
	db             *gorm.DB
	snapshotRepo   *PostgresSnapshotRepo
	permissionRepo *PostgresPermissionRepo
}

func NewPostgresFileRepo() *PostgresFileRepo {
	return &PostgresFileRepo{
		db:             infra.GetDb(),
		snapshotRepo:   NewPostgresSnapshotRepo(),
		permissionRepo: NewPostgresPermissionRepo(),
	}
}

func (repo *PostgresFileRepo) New() model.CoreFile {
	return &PostgresFile{}
}

func (repo *PostgresFileRepo) Insert(opts FileInsertOptions) (model.CoreFile, error) {
	id := helpers.NewId()
	file := PostgresFile{
		ID:          id,
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
	if err := repo.populateModelFields([]*PostgresFile{res}); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *PostgresFileRepo) Find(id string) (model.CoreFile, error) {
	file, err := repo.find(id)
	if err != nil {
		return nil, err
	}
	if err := repo.populateModelFields([]*PostgresFile{file}); err != nil {
		return nil, err
	}
	return file, nil
}

func (repo *PostgresFileRepo) find(id string) (*PostgresFile, error) {
	var res = PostgresFile{}
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

func (repo *PostgresFileRepo) FindChildren(id string) ([]model.CoreFile, error) {
	var entities []*PostgresFile
	db := repo.db.Raw("SELECT * FROM file WHERE parent_id = ? ORDER BY create_time ASC", id).Scan(&entities)
	if db.Error != nil {
		return nil, db.Error
	}
	if err := repo.populateModelFields(entities); err != nil {
		return nil, err
	}
	var res []model.CoreFile
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *PostgresFileRepo) FindPath(id string) ([]model.CoreFile, error) {
	var entities []*PostgresFile
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
	var res []model.CoreFile
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *PostgresFileRepo) FindTree(id string) ([]model.CoreFile, error) {
	var entities []*PostgresFile
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
	var res []model.CoreFile
	for _, f := range entities {
		res = append(res, f)
	}
	return res, nil
}

func (repo *PostgresFileRepo) GetIdsByWorkspace(workspaceId string) ([]string, error) {
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

func (repo *PostgresFileRepo) AssignSnapshots(cloneId string, originalId string) error {
	if db := repo.db.Exec("INSERT INTO snapshot_file (snapshot_id, file_id) SELECT s.id, ? "+
		"FROM snapshot s LEFT JOIN snapshot_file map ON s.id = map.snapshot_id "+
		"WHERE map.file_id = ? ORDER BY s.version DESC LIMIT 1", cloneId, originalId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresFileRepo) MoveSourceIntoTarget(targetId string, sourceId string) error {
	if db := repo.db.Exec("UPDATE file SET parent_id = ? WHERE id = ?", targetId, sourceId); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresFileRepo) Save(file model.CoreFile) error {
	if db := repo.db.Save(file); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresFileRepo) BulkInsert(values []model.CoreFile, chunkSize int) error {
	var entities []*PostgresFile
	for _, f := range values {
		entities = append(entities, f.(*PostgresFile))
	}
	if db := repo.db.CreateInBatches(entities, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresFileRepo) BulkInsertPermissions(values []*UserPermission, chunkSize int) error {
	if db := repo.db.CreateInBatches(values, chunkSize); db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresFileRepo) Delete(id string) error {
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

func (repo *PostgresFileRepo) GetChildrenIDs(id string) ([]string, error) {
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

func (repo *PostgresFileRepo) GetItemCount(id string) (int64, error) {
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

func (repo *PostgresFileRepo) IsGrandChildOf(id string, ancestorId string) (bool, error) {
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

func (repo *PostgresFileRepo) GetSize(id string) (int64, error) {
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

func (repo *PostgresFileRepo) GrantUserPermission(id string, userId string, permission string) error {
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
			helpers.NewId(), userId, f.GetID())
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
			helpers.NewId(), userId, f.GetID(), permission, permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *PostgresFileRepo) RevokeUserPermission(id string, userId string) error {
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.Exec("DELETE FROM userpermission WHERE user_id = ? AND resource_id = ?", userId, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *PostgresFileRepo) GrantGroupPermission(id string, groupId string, permission string) error {
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
			helpers.NewId(), groupId, f.GetID())
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
			helpers.NewId(), groupId, f.GetID(), permission, permission)
		if db.Error != nil {
			return db.Error
		}
	}

	return nil
}

func (repo *PostgresFileRepo) RevokeGroupPermission(id string, groupId string) error {
	tree, err := repo.FindTree(id)
	if err != nil {
		return err
	}
	for _, f := range tree {
		db := repo.db.Exec("DELETE FROM grouppermission WHERE group_id = ? AND resource_id = ?", groupId, f.GetID())
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func (repo *PostgresFileRepo) populateModelFields(entities []*PostgresFile) error {
	for _, f := range entities {
		f.UserPermissions = make([]*model.UserPermission, 0)
		userPermissions, err := repo.permissionRepo.GetUserPermissions(f.ID)
		if err != nil {
			return err
		}
		for _, p := range userPermissions {
			f.UserPermissions = append(f.UserPermissions, &model.UserPermission{
				UserId: p.UserID,
				Value:  p.Permission,
			})
		}
		f.GroupPermissions = make([]*model.GroupPermission, 0)
		groupPermissions, err := repo.permissionRepo.GetGroupPermissions(f.ID)
		if err != nil {
			return err
		}
		for _, p := range groupPermissions {
			f.GroupPermissions = append(f.GroupPermissions, &model.GroupPermission{
				GroupID: p.GroupID,
				Value:   p.Permission,
			})
		}
		snapshots, err := repo.snapshotRepo.findAllForFile(f.ID)
		if err != nil {
			return nil
		}
		f.Snapshots = snapshots
	}
	return nil
}
