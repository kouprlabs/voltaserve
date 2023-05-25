package core

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
	Id           string        `json:"id"`
	Owner        *User         `json:"owner,omitempty"`
	Email        string        `json:"email"`
	Organization *Organization `json:"organization,omitempty"`
	Status       string        `json:"status"`
	CreateTime   string        `json:"createTime"`
	UpdateTime   *string       `json:"updateTime"`
}

type InvitationCreateOptions struct {
	OrganizationId string   `json:"organizationId" validate:"required"`
	Emails         []string `json:"emails" validate:"required,dive,email"`
}

type InvitationService struct {
	orgRepo          repo.CoreOrganizationRepo
	orgMapper        *organizationMapper
	invitationRepo   repo.CoreInvitationRepo
	invitationMapper *invitationMapper
	orgCache         *cache.OrganizationCache
	orgGuard         *guard.OrganizationGuard
	userRepo         repo.CoreUserRepo
	mailTmpl         *infra.MailTemplate
	config           config.Config
}

func NewInvitationService() *InvitationService {
	return &InvitationService{
		orgRepo:          repo.NewPostgresOrganizationRepo(),
		orgCache:         cache.NewOrganizationCache(),
		orgGuard:         guard.NewOrganizationGuard(),
		invitationRepo:   repo.NewPostgresInvitationRepo(),
		invitationMapper: newInvitationMapper(),
		userRepo:         repo.NewPostgresUserRepo(),
		mailTmpl:         infra.NewMailTemplate(),
		orgMapper:        newOrganizationMapper(),
		config:           config.GetConfig(),
	}
}

func (svc *InvitationService) Create(req InvitationCreateOptions, userId string) error {
	for i := range req.Emails {
		req.Emails[i] = strings.ToLower(req.Emails[i])
	}
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return err
	}
	org, err := svc.orgCache.Get(req.OrganizationId)
	if err != nil {
		return err
	}
	if err := svc.orgGuard.Authorize(user, org, model.PermissionOwner); err != nil {
		return err
	}
	orgMembers, err := svc.orgRepo.GetMembers(req.OrganizationId)
	if err != nil {
		return err
	}

	/* Collect emails of non existing members */
	var emails []string
	for _, e := range req.Emails {
		existing := false
		for _, u := range orgMembers {
			if e == u.GetEmail() {
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
		UserId:         userId,
		OrganizationId: req.OrganizationId,
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
			"UI_URL":            svc.config.UIURL,
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

func (svc *InvitationService) GetIncoming(userId string) ([]*Invitation, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	res, err := svc.invitationMapper.mapInvitations(invitations, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *InvitationService) GetOutgoing(id string, userId string) ([]*Invitation, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetOutgoing(id, user.GetID())
	if err != nil {
		return nil, err
	}
	res, err := svc.invitationMapper.mapInvitations(invitations, userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *InvitationService) Accept(id string, userId string) error {
	user, err := svc.userRepo.Find(userId)
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
		if u == userId {
			return errorpkg.NewUserAlreadyMemberOfOrganizationError(user, org)
		}
	}
	invitation.SetStatus(model.InvitationStatusAccepted)
	if err := svc.invitationRepo.Save(invitation); err != nil {
		return err
	}
	if err := svc.orgRepo.AddUser(invitation.GetOrganizationID(), userId); err != nil {
		return err
	}
	if err := svc.orgRepo.GrantUserPermission(invitation.GetOrganizationID(), userId, model.PermissionViewer); err != nil {
		return err
	}
	if _, err := svc.orgCache.Refresh(invitation.GetOrganizationID()); err != nil {
		return err
	}
	return nil
}

func (svc *InvitationService) Decline(id string, userId string) error {
	user, err := svc.userRepo.Find(userId)
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

func (svc *InvitationService) Resend(id string, userId string) error {
	user, err := svc.userRepo.Find(userId)
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
		"UI_URL":            svc.config.UIURL,
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

func (svc *InvitationService) Delete(id string, userId string) error {
	invitation, err := svc.invitationRepo.Find(id)
	if err != nil {
		return err
	}
	if userId != invitation.GetOwnerID() {
		user, err := svc.userRepo.Find(userId)
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
	userRepo   repo.CoreUserRepo
	userMapper *userMapper
	orgMapper  *organizationMapper
}

func newInvitationMapper() *invitationMapper {
	return &invitationMapper{
		orgCache:   cache.NewOrganizationCache(),
		userRepo:   repo.NewPostgresUserRepo(),
		userMapper: newUserMapper(),
		orgMapper:  newOrganizationMapper(),
	}
}

func (mp *invitationMapper) mapInvitation(m model.InvitationModel, userId string) (*Invitation, error) {
	owner, err := mp.userRepo.Find(m.GetOwnerID())
	if err != nil {
		return nil, err
	}
	org, err := mp.orgCache.Get(m.GetOrganizationID())
	if err != nil {
		return nil, err
	}
	v, err := mp.orgMapper.mapOrganization(org, userId)
	if err != nil {
		return nil, err
	}
	return &Invitation{
		Id:           m.GetID(),
		Owner:        mp.userMapper.mapUser(owner),
		Email:        m.GetEmail(),
		Organization: v,
		Status:       m.GetStatus(),
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}, nil
}

func (mp *invitationMapper) mapInvitations(invitations []model.InvitationModel, userId string) ([]*Invitation, error) {
	res := make([]*Invitation, 0)
	for _, m := range invitations {
		v, err := mp.mapInvitation(m, userId)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
