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

type InvitationEntity struct {
	Id             string  `json:"id"`
	OrganizationId string  `json:"organizationId"`
	OwnerId        string  `json:"ownerId"`
	Email          string  `json:"email"`
	Status         string  `json:"status"`
	CreateTime     string  `json:"createTime"`
	UpdateTime     *string `json:"updateTime"`
}

func (InvitationEntity) TableName() string {
	return "invitation"
}

func (o *InvitationEntity) BeforeCreate(tx *gorm.DB) (err error) {
	o.CreateTime = time.Now().UTC().Format(time.RFC3339)
	return nil
}

func (o *InvitationEntity) BeforeSave(tx *gorm.DB) (err error) {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	o.UpdateTime = &timeNow
	return nil
}

func (i InvitationEntity) GetId() string {
	return i.Id
}

func (i InvitationEntity) GetOrganizationId() string {
	return i.OrganizationId
}

func (i InvitationEntity) GetOwnerId() string {
	return i.OwnerId
}

func (i InvitationEntity) GetEmail() string {
	return i.Email
}

func (i InvitationEntity) GetStatus() string {
	return i.Status
}

func (i InvitationEntity) GetCreateTime() string {
	return i.CreateTime
}

func (i InvitationEntity) GetUpdateTime() *string {
	return i.UpdateTime
}

func (w *InvitationEntity) SetStatus(status string) {
	w.Status = status
}

func (w *InvitationEntity) SetUpdateTime(updateTime *string) {
	w.UpdateTime = updateTime
}

type InvitationRepo struct {
	db       *gorm.DB
	userRepo *UserRepo
}

func NewInvitationRepo() *InvitationRepo {
	return &InvitationRepo{
		db:       infra.GetDb(),
		userRepo: NewUserRepo(),
	}
}

type InvitationInsertOptions struct {
	UserId         string
	OrganizationId string
	Emails         []string
}

func (repo *InvitationRepo) Insert(opts InvitationInsertOptions) ([]model.InvitationModel, error) {
	var res []model.InvitationModel
	for _, e := range opts.Emails {
		invitation := InvitationEntity{
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

func (repo *InvitationRepo) Find(id string) (model.InvitationModel, error) {
	var invitation = InvitationEntity{}
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

func (repo *InvitationRepo) GetIncoming(email string) ([]model.InvitationModel, error) {
	var invitations []*InvitationEntity
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

func (repo *InvitationRepo) GetOutgoing(organizationId string, userId string) ([]model.InvitationModel, error) {
	var invitations []*InvitationEntity
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

func (repo *InvitationRepo) Save(org model.InvitationModel) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *InvitationRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM invitation WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
