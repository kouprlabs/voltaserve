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

type PostgresInvitation struct {
	Id             string  `json:"id"`
	OrganizationId string  `json:"organizationId"`
	OwnerId        string  `json:"ownerId"`
	Email          string  `json:"email"`
	Status         string  `json:"status"`
	CreateTime     string  `json:"createTime"`
	UpdateTime     *string `json:"updateTime"`
}

func (PostgresInvitation) TableName() string {
	return "invitation"
}

func (o *PostgresInvitation) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (o *PostgresInvitation) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	o.UpdateTime = &timeNow
	return nil
}

func (i PostgresInvitation) GetID() string {
	return i.Id
}

func (i PostgresInvitation) GetOrganizationID() string {
	return i.OrganizationId
}

func (i PostgresInvitation) GetOwnerID() string {
	return i.OwnerId
}

func (i PostgresInvitation) GetEmail() string {
	return i.Email
}

func (i PostgresInvitation) GetStatus() string {
	return i.Status
}

func (i PostgresInvitation) GetCreateTime() string {
	return i.CreateTime
}

func (i PostgresInvitation) GetUpdateTime() *string {
	return i.UpdateTime
}

func (w *PostgresInvitation) SetStatus(status string) {
	w.Status = status
}

func (w *PostgresInvitation) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type PostgresInvitationRepo struct {
	db       *gorm.DB
	userRepo *PostgresUserRepo
}

func NewPostgresInvitationRepo() *PostgresInvitationRepo {
	return &PostgresInvitationRepo{
		db:       infra.GetDb(),
		userRepo: NewPostgresUserRepo(),
	}
}

func (repo *PostgresInvitationRepo) Insert(opts InvitationInsertOptions) ([]model.InvitationModel, error) {
	var res []model.InvitationModel
	for _, e := range opts.Emails {
		invitation := PostgresInvitation{
			Id:             helpers.NewId(),
			OrganizationId: opts.OrganizationId,
			OwnerId:        opts.UserId,
			Email:          e,
			Status:         model.InvitationStatusPending,
		}
		if db := repo.db.Save(&invitation); db.Error != nil {
			return nil, db.Error
		}
		i, err := repo.Find(invitation.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func (repo *PostgresInvitationRepo) Find(id string) (model.InvitationModel, error) {
	var invitation = PostgresInvitation{}
	db := repo.db.Where("id = ?", id).First(&invitation)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, errorpkg.NewInvitationNotFoundError(db.Error)
		} else {
			return nil, errorpkg.NewInternalServerError(db.Error)
		}
	}
	return &invitation, nil
}

func (repo *PostgresInvitationRepo) GetIncoming(email string) ([]model.InvitationModel, error) {
	var invitations []*PostgresInvitation
	db := repo.db.
		Raw("SELECT * FROM invitation WHERE email = ? and status = 'pending' ORDER BY create_time DESC", email).
		Scan(&invitations)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.InvitationModel
	for _, inv := range invitations {
		res = append(res, inv)
	}
	return res, nil
}

func (repo *PostgresInvitationRepo) GetOutgoing(organizationId string, userId string) ([]model.InvitationModel, error) {
	var invitations []*PostgresInvitation
	db := repo.db.
		Raw("SELECT * FROM invitation WHERE organization_id = ? and owner_id = ? ORDER BY create_time DESC", organizationId, userId).
		Scan(&invitations)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.InvitationModel
	for _, inv := range invitations {
		res = append(res, inv)
	}
	return res, nil
}

func (repo *PostgresInvitationRepo) Save(org model.InvitationModel) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *PostgresInvitationRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM invitation WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
