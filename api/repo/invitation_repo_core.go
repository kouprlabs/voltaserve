package repo

import "voltaserve/model"

type InvitationInsertOptions struct {
	UserId         string
	OrganizationId string
	Emails         []string
}

type CoreInvitationRepo interface {
	Insert(opts InvitationInsertOptions) ([]model.CoreInvitation, error)
	Find(id string) (model.CoreInvitation, error)
	GetIncoming(email string) ([]model.CoreInvitation, error)
	GetOutgoing(organizationId string, userId string) ([]model.CoreInvitation, error)
	Save(org model.CoreInvitation) error
	Delete(id string) error
}
