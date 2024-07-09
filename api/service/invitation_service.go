// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

package service

import (
	"sort"
	"strings"
	"time"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/errorpkg"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/model"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type InvitationService struct {
	orgRepo          repo.OrganizationRepo
	orgMapper        *organizationMapper
	invitationRepo   repo.InvitationRepo
	invitationMapper *invitationMapper
	orgCache         *cache.OrganizationCache
	orgGuard         *guard.OrganizationGuard
	userRepo         repo.UserRepo
	mailTmpl         *infra.MailTemplate
	config           *config.Config
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

type InvitationCreateOptions struct {
	OrganizationID string   `json:"organizationId" validate:"required"`
	Emails         []string `json:"emails"         validate:"required,dive,email"`
}

func (svc *InvitationService) Create(opts InvitationCreateOptions, userID string) error {
	for i := range opts.Emails {
		opts.Emails[i] = strings.ToLower(opts.Emails[i])
	}
	org, err := svc.orgCache.Get(opts.OrganizationID)
	if err != nil {
		return err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionOwner); err != nil {
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
			if e == i.GetEmail() && i.GetStatus() == model.InvitationStatusPending {
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
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
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

type InvitationListOptions struct {
	Page      uint
	Size      uint
	SortBy    string
	SortOrder string
}

type Invitation struct {
	ID           string        `json:"id"`
	Owner        *User         `json:"owner,omitempty"`
	Email        string        `json:"email"`
	Organization *Organization `json:"organization,omitempty"`
	Status       string        `json:"status"`
	CreateTime   string        `json:"createTime"`
	UpdateTime   *string       `json:"updateTime"`
}

type InvitationList struct {
	Data          []*Invitation `json:"data"`
	TotalPages    uint          `json:"totalPages"`
	TotalElements uint          `json:"totalElements"`
	Page          uint          `json:"page"`
	Size          uint          `json:"size"`
}

func (svc *InvitationService) GetIncoming(opts InvitationListOptions, userID string) (*InvitationList, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.doSorting(invitations, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.invitationMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &InvitationList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(mapped)),
	}, nil
}

func (svc *InvitationService) GetIncomingCount(userID string) (*int64, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	var res int64
	if res, err = svc.invitationRepo.GetIncomingCount(user.GetEmail()); err != nil {
		return nil, err
	}
	return &res, nil
}

func (svc *InvitationService) GetOutgoing(orgID string, opts InvitationListOptions, userID string) (*InvitationList, error) {
	invitations, err := svc.invitationRepo.GetOutgoing(orgID, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = SortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = SortOrderAsc
	}
	sorted := svc.doSorting(invitations, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.doPagination(sorted, opts.Page, opts.Size)
	mapped, err := svc.invitationMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &InvitationList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint(len(mapped)),
	}, nil
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

func (svc *InvitationService) doSorting(data []model.Invitation, sortBy string, sortOrder string) []model.Invitation {
	if sortBy == SortByEmail {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == SortOrderDesc {
				return data[i].GetEmail() > data[j].GetEmail()
			} else {
				return data[i].GetEmail() < data[j].GetEmail()
			}
		})
		return data
	} else if sortBy == SortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a, _ := time.Parse(time.RFC3339, data[i].GetCreateTime())
			b, _ := time.Parse(time.RFC3339, data[j].GetCreateTime())
			if sortOrder == SortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == SortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a, _ := time.Parse(time.RFC3339, *data[i].GetUpdateTime())
				b, _ := time.Parse(time.RFC3339, *data[j].GetUpdateTime())
				if sortOrder == SortOrderDesc {
					return a.UnixMilli() > b.UnixMilli()
				} else {
					return a.UnixMilli() < b.UnixMilli()
				}
			} else {
				return false
			}
		})
		return data
	}
	return data
}

func (svc *InvitationService) doPagination(data []model.Invitation, page, size uint) (pageData []model.Invitation, totalElements uint, totalPages uint) {
	totalElements = uint(len(data))
	totalPages = (totalElements + size - 1) / size
	if page > totalPages {
		return []model.Invitation{}, totalElements, totalPages
	}
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if endIndex > totalElements {
		endIndex = totalElements
	}
	return data[startIndex:endIndex], totalElements, totalPages
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
	o, err := mp.orgMapper.mapOne(org, userID)
	if err != nil {
		return nil, err
	}
	return &Invitation{
		ID:           m.GetID(),
		Owner:        mp.userMapper.mapOne(owner),
		Email:        m.GetEmail(),
		Organization: o,
		Status:       m.GetStatus(),
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}, nil
}

func (mp *invitationMapper) mapMany(invitations []model.Invitation, userID string) ([]*Invitation, error) {
	res := make([]*Invitation, 0)
	for _, invitation := range invitations {
		i, err := mp.mapOne(invitation, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}
