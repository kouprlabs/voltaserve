package model

const (
	InvitationStatusPending  = "pending"
	InvitationStatusAccepted = "accepted"
	InvitationStatusDeclined = "declined"
)

type InvitationModel interface {
	GetId() string
	GetOrganizationId() string
	GetOwnerId() string
	GetEmail() string
	GetStatus() string
	GetCreateTime() string
	GetUpdateTime() *string
	SetStatus(string)
	SetUpdateTime(*string)
}
