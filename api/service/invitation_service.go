package service

import (
	"strings"
	"voltaserve/cache"
	"voltaserve/config"
	"voltaserve/errorpkg"
	"voltaserve/guard"
	"voltaserve/infra"
	"voltaserve/model"
	"voltaserve/repo"
)

type Invitation struct {
	ID           string        `json:"id"`
	Owner        *User         `json:"owner,omitempty"`
	Email        string        `json:"email"`
	Organization *Organization `json:"organization,omitempty"`
	Status       string        `json:"status"`
	CreateTime   string        `json:"createTime"`
	UpdateTime   *string       `json:"updateTime"`
}

type InvitationCreateOptions struct {
	OrganizationID string   `json:"organizationId" validate:"required"`
	Emails         []string `json:"emails" validate:"required,dive,email"`
}

type InvitationService struct {
	orgRepo          repo.OrganizationRepo
	orgMapper        *organizationMapper
	invitationRepo   repo.InvitationRepo
	invitationMapper *invitationMapper
	orgCache         *cache.OrganizationCache
	orgGuard         *guard.OrganizationGuard
	userRepo         repo.UserRepo
	mailTmpl         *infra.MailTemplate
	config           config.Config
}

func NewInvitationService() *InvitationService {
	return &InvitationService{
		orgRepo:          repo.NewOrganizationRepo(),
		orgCache:         cache.NewOrganizationCache(),
		orgGuard:         guard.NewOrganizationGuard(),
		invitationRepo:   repo.NewInvitationRepo(),
		invitationMapper: newInvitationMapper(),
		userRepo:         repo.NewUserRepo(),
		mailTmpl:         infra.NewMailTemplate(),
		orgMapper:        newOrganizationMapper(),
		config:           config.GetConfig(),
	}
}

func (svc *InvitationService) Create(opts InvitationCreateOptions, userID string) error {
	for i := range opts.Emails {
		opts.Emails[i] = strings.ToLower(opts.Emails[i])
	}
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	org, err := svc.orgCache.Get(opts.OrganizationID)
	if err != nil {
		return err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionOwner); err != nil {
		return err
	}
	orgMembers, err := svc.orgRepo.GetMembers(opts.OrganizationID)
	if err != nil {
		return err
	}
	outgoingInvitations, err := svc.invitationRepo.GetOutgoing(opts.OrganizationID, userID)
	if err != nil {
		return err
	}

	var emails []string

	/* Collect emails of non existing members and outgoing invitations */
	for _, e := range opts.Emails {
		existing := false
		for _, u := range orgMembers {
			if e == u.GetEmail() {
				existing = true
				break
			}
		}
		for _, i := range outgoingInvitations {
			if e == i.GetEmail() {
				existing = true
				break
			}
		}
		if !existing {
			emails = append(emails, e)
		}
	}

	/* Persist invitations */
	invitations, err := svc.invitationRepo.Insert(repo.InvitationInsertOptions{
		UserID:         userID,
		OrganizationID: opts.OrganizationID,
		Emails:         emails,
	})
	if err != nil {
		return err
	}

	/* Send emails */
	for _, inv := range invitations {
		variables := map[string]string{
			"USER_FULL_NAME":    user.GetFullName(),
			"ORGANIZATION_NAME": org.GetName(),
			"UI_URL":            svc.config.PublicUIURL,
		}
		_, err := svc.userRepo.FindByEmail(inv.GetEmail())
		var templateName string
		if err == nil {
			templateName = "join-organization"
		} else {
			templateName = "signup-and-join-organization"
		}
		if err := svc.mailTmpl.Send(templateName, inv.GetEmail(), variables); err != nil {
			return err
		}
	}
	return nil
}

func (svc *InvitationService) GetIncoming(userID string) ([]*Invitation, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	res, err := svc.invitationMapper.mapMany(invitations, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *InvitationService) GetOutgoing(id string, userID string) ([]*Invitation, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetOutgoing(id, user.GetID())
	if err != nil {
		return nil, err
	}
	res, err := svc.invitationMapper.mapMany(invitations, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *InvitationService) Accept(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	invitation, err := svc.invitationRepo.Find(id)
	if err != nil {
		return err
	}
	if invitation.GetStatus() != model.InvitationStatusPending {
		return errorpkg.NewCannotAcceptNonPendingInvitationError(invitation)
	}
	if user.GetEmail() != invitation.GetEmail() {
		return errorpkg.NewUserNotAllowedToAcceptInvitationError(user, invitation)
	}
	org, err := svc.orgCache.Get(invitation.GetOrganizationID())
	if err != nil {
		return err
	}
	for _, u := range org.GetUsers() {
		if u == userID {
			return errorpkg.NewUserAlreadyMemberOfOrganizationError(user, org)
		}
	}
	invitation.SetStatus(model.InvitationStatusAccepted)
	if err := svc.invitationRepo.Save(invitation); err != nil {
		return err
	}
	if err := svc.orgRepo.AddUser(invitation.GetOrganizationID(), userID); err != nil {
		return err
	}
	if err := svc.orgRepo.GrantUserPermission(invitation.GetOrganizationID(), userID, model.PermissionViewer); err != nil {
		return err
	}
	if _, err := svc.orgCache.Refresh(invitation.GetOrganizationID()); err != nil {
		return err
	}
	return nil
}

func (svc *InvitationService) Decline(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	invitation, err := svc.invitationRepo.Find(id)
	if err != nil {
		return err
	}
	if invitation.GetStatus() != model.InvitationStatusPending {
		return errorpkg.NewCannotDeclineNonPendingInvitationError(invitation)
	}
	if user.GetEmail() != invitation.GetEmail() {
		return errorpkg.NewUserNotAllowedToDeclineInvitationError(user, invitation)
	}
	invitation.SetStatus(model.InvitationStatusDeclined)
	if err := svc.invitationRepo.Save(invitation); err != nil {
		return err
	}
	return nil
}

func (svc *InvitationService) Resend(id string, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	invitation, err := svc.invitationRepo.Find(id)
	if err != nil {
		return err
	}
	if invitation.GetStatus() != model.InvitationStatusPending {
		return errorpkg.NewCannotResendNonPendingInvitationError(invitation)
	}
	org, err := svc.orgCache.Get(invitation.GetOrganizationID())
	if err != nil {
		return err
	}
	variables := map[string]string{
		"USER_FULL_NAME":    user.GetFullName(),
		"ORGANIZATION_NAME": org.GetName(),
		"UI_URL":            svc.config.PublicUIURL,
	}
	_, err = svc.userRepo.FindByEmail(invitation.GetEmail())
	var templateName string
	if err == nil {
		templateName = "join-organization"
	} else {
		templateName = "signup-and-join-organization"
	}
	if err := svc.mailTmpl.Send(templateName, invitation.GetEmail(), variables); err != nil {
		return err
	}
	return nil
}

func (svc *InvitationService) Delete(id string, userID string) error {
	invitation, err := svc.invitationRepo.Find(id)
	if err != nil {
		return err
	}
	if userID != invitation.GetOwnerID() {
		user, err := svc.userRepo.Find(userID)
		if err != nil {
			return err
		}
		return errorpkg.NewUserNotAllowedToDeleteInvitationError(user, invitation)
	}
	if err := svc.invitationRepo.Delete(invitation.GetID()); err != nil {
		return err
	}
	return nil
}

type invitationMapper struct {
	orgCache   *cache.OrganizationCache
	userRepo   repo.UserRepo
	userMapper *userMapper
	orgMapper  *organizationMapper
}

func newInvitationMapper() *invitationMapper {
	return &invitationMapper{
		orgCache:   cache.NewOrganizationCache(),
		userRepo:   repo.NewUserRepo(),
		userMapper: newUserMapper(),
		orgMapper:  newOrganizationMapper(),
	}
}

func (mp *invitationMapper) mapOne(m model.Invitation, userID string) (*Invitation, error) {
	owner, err := mp.userRepo.Find(m.GetOwnerID())
	if err != nil {
		return nil, err
	}
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	v, err := mp.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return &Invitation{
		ID:           m.GetID(),
		Owner:        mp.userMapper.mapOne(owner),
		Email:        m.GetEmail(),
		Organization: v,
		Status:       m.GetStatus(),
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}, nil
}

func (mp *invitationMapper) mapMany(invitations []model.Invitation, userID string) ([]*Invitation, error) {
	res := make([]*Invitation, 0)
	for _, m := range invitations {
		v, err := mp.mapOne(m, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
