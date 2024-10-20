// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package repo

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/helper"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
)

type InvitationRepo interface {
	Insert(opts InvitationInsertOptions) ([]model.Invitation, error)
	Find(id string) (model.Invitation, error)
	FindIncoming(email string) ([]model.Invitation, error)
	CountIncoming(email string) (int64, error)
	FindOutgoing(orgID string, userID string) ([]model.Invitation, error)
	Save(org model.Invitation) error
	Delete(id string) error
}

func NewInvitationRepo() InvitationRepo {
	return newInvitationRepo()
}

type invitationEntity struct {
	ID             string  `gorm:"column:id"              json:"id"`
	OrganizationID string  `gorm:"column:organization_id" json:"organizationId"`
	OwnerID        string  `gorm:"column:owner_id"        json:"ownerId"`
	Email          string  `gorm:"column:email"           json:"email"`
	Status         string  `gorm:"column:status"          json:"status"`
	CreateTime     string  `gorm:"column:create_time"     json:"createTime"`
	UpdateTime     *string `gorm:"column:update_time"     json:"updateTime"`
}

func (*invitationEntity) TableName() string {
	return "invitation"
}

func (i *invitationEntity) BeforeCreate(*gorm.DB) (err error) {
	i.CreateTime = helper.NewTimestamp()
	return nil
}

func (i *invitationEntity) BeforeSave(*gorm.DB) (err error) {
	timeNow := helper.NewTimestamp()
	i.UpdateTime = &timeNow
	return nil
}

func (i *invitationEntity) GetID() string {
	return i.ID
}

func (i *invitationEntity) GetOrganizationID() string {
	return i.OrganizationID
}

func (i *invitationEntity) GetOwnerID() string {
	return i.OwnerID
}

func (i *invitationEntity) GetEmail() string {
	return i.Email
}

func (i *invitationEntity) GetStatus() string {
	return i.Status
}

func (i *invitationEntity) GetCreateTime() string {
	return i.CreateTime
}

func (i *invitationEntity) GetUpdateTime() *string {
	return i.UpdateTime
}

func (i *invitationEntity) SetStatus(status string) {
	i.Status = status
}

type invitationRepo struct {
	db       *gorm.DB
	userRepo *userRepo
}

func newInvitationRepo() *invitationRepo {
	return &invitationRepo{
		db:       infra.NewPostgresManager().GetDBOrPanic(),
		userRepo: newUserRepo(),
	}
}

type InvitationInsertOptions struct {
	UserID         string
	OrganizationID string
	Emails         []string
}

func (repo *invitationRepo) Insert(opts InvitationInsertOptions) ([]model.Invitation, error) {
	var res []model.Invitation
	for _, e := range opts.Emails {
		invitation := invitationEntity{
			ID:             helper.NewID(),
			OrganizationID: opts.OrganizationID,
			OwnerID:        opts.UserID,
			Email:          e,
			Status:         model.InvitationStatusPending,
		}
		if db := repo.db.Create(&invitation); db.Error != nil {
			return nil, db.Error
		}
		i, err := repo.Find(invitation.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

func (repo *invitationRepo) Find(id string) (model.Invitation, error) {
	invitation := invitationEntity{}
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

func (repo *invitationRepo) FindIncoming(email string) ([]model.Invitation, error) {
	var invitations []*invitationEntity
	db := repo.db.
		Raw("SELECT * FROM invitation WHERE email = ? and status = 'pending' ORDER BY create_time DESC", email).
		Scan(&invitations)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Invitation
	for _, inv := range invitations {
		res = append(res, inv)
	}
	return res, nil
}

func (repo *invitationRepo) CountIncoming(email string) (int64, error) {
	var count int64
	db := repo.db.
		Model(&invitationEntity{}).
		Where("email = ?", email).
		Where("status = 'pending'").
		Count(&count)
	if db.Error != nil {
		return -1, db.Error
	}
	return count, nil
}

func (repo *invitationRepo) FindOutgoing(orgID string, userID string) ([]model.Invitation, error) {
	var invitations []*invitationEntity
	db := repo.db.
		Raw("SELECT * FROM invitation WHERE organization_id = ? and owner_id = ? ORDER BY create_time DESC", orgID, userID).
		Scan(&invitations)
	if db.Error != nil {
		return nil, db.Error
	}
	var res []model.Invitation
	for _, inv := range invitations {
		res = append(res, inv)
	}
	return res, nil
}

func (repo *invitationRepo) Save(org model.Invitation) error {
	db := repo.db.Save(org)
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (repo *invitationRepo) Delete(id string) error {
	db := repo.db.Exec("DELETE FROM invitation WHERE id = ?", id)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
