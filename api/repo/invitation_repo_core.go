package repo

import "voltaserve/model"

type InvitationInsertOptions struct {
	UserId         string
	OrganizationId string
	Emails         []string
}

type CoreInvitationRepo interface {
	Insert(opts InvitationInsertOptions) ([]model.InvitationModel, error)
	Find(id string) (model.InvitationModel, error)
	GetIncoming(email string) ([]model.InvitationModel, error)
	GetOutgoing(organizationId string, userId string) ([]model.InvitationModel, error)
	Save(org model.InvitationModel) error
	Delete(id string) error
}
