// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

package service

import (
	"sort"
	"strings"

	"github.com/kouprlabs/voltaserve/shared/dto"
	"github.com/kouprlabs/voltaserve/shared/errorpkg"
	"github.com/kouprlabs/voltaserve/shared/helper"
	"github.com/kouprlabs/voltaserve/shared/model"

	"github.com/kouprlabs/voltaserve/api/cache"
	"github.com/kouprlabs/voltaserve/api/config"
	"github.com/kouprlabs/voltaserve/api/guard"
	"github.com/kouprlabs/voltaserve/api/infra"
	"github.com/kouprlabs/voltaserve/api/repo"
)

type InvitationService struct {
	orgRepo          *repo.OrganizationRepo
	orgMapper        *organizationMapper
	invitationRepo   *repo.InvitationRepo
	invitationMapper *invitationMapper
	orgCache         *cache.OrganizationCache
	orgGuard         *guard.OrganizationGuard
	orgSvc           *OrganizationService
	userRepo         *repo.UserRepo
	mailTmpl         infra.MailTemplate
	config           *config.Config
}

func NewInvitationService() *InvitationService {
	return &InvitationService{
		orgRepo:          repo.NewOrganizationRepo(),
		orgCache:         cache.NewOrganizationCache(),
		orgGuard:         guard.NewOrganizationGuard(),
		orgSvc:           NewOrganizationService(),
		invitationRepo:   repo.NewInvitationRepo(),
		invitationMapper: newInvitationMapper(),
		userRepo:         repo.NewUserRepo(),
		mailTmpl:         infra.NewMailTemplate(config.GetConfig().SMTP),
		orgMapper:        newOrganizationMapper(),
		config:           config.GetConfig(),
	}
}

func (svc *InvitationService) Create(opts dto.InvitationCreateOptions, userID string) ([]*dto.Invitation, error) {
	for i := range opts.Emails {
		opts.Emails[i] = strings.ToLower(opts.Emails[i])
	}
	org, err := svc.orgCache.Get(opts.OrganizationID)
	if err != nil {
		return nil, err
	}
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionOwner); err != nil {
		return nil, err
	}
	orgMembers, err := svc.orgRepo.FindMembers(opts.OrganizationID)
	if err != nil {
		return nil, err
	}
	outgoing, err := svc.invitationRepo.FindOutgoing(opts.OrganizationID, userID)
	if err != nil {
		return nil, err
	}
	emails := svc.getValidOutboundEmails(opts.Emails, user.GetEmail(), orgMembers, outgoing)
	invitations, err := svc.invitationRepo.Insert(repo.InvitationInsertOptions{
		UserID:         userID,
		OrganizationID: opts.OrganizationID,
		Emails:         emails,
	})
	if err != nil {
		return nil, err
	}
	if err := svc.sendEmails(invitations, org, userID); err != nil {
		return nil, err
	}
	res, err := svc.invitationMapper.mapMany(invitations, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type InvitationListOptions struct {
	Page      uint64
	Size      uint64
	SortBy    string
	SortOrder string
}

func (svc *InvitationService) ListIncoming(opts InvitationListOptions, userID string) (*dto.InvitationList, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.FindIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.InvitationSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.InvitationSortOrderAsc
	}
	sorted := svc.sort(invitations, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.invitationMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &dto.InvitationList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *InvitationService) ProbeIncoming(opts InvitationListOptions, userID string) (*dto.InvitationProbe, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	totalElements, err := svc.invitationRepo.CountIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	return &dto.InvitationProbe{
		TotalElements: uint64(totalElements),                               // #nosec G115 integer overflow conversion
		TotalPages:    (uint64(totalElements) + opts.Size - 1) / opts.Size, // #nosec G115 integer overflow conversion
	}, nil
}

func (svc *InvitationService) GetCountIncoming(userID string) (*int64, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	var res int64
	if res, err = svc.invitationRepo.CountIncoming(user.GetEmail()); err != nil {
		return nil, err
	}
	return &res, nil
}

func (svc *InvitationService) ListOutgoing(orgID string, opts InvitationListOptions, userID string) (*dto.InvitationList, error) {
	org, err := svc.orgCache.Get(orgID)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionOwner); err != nil {
		return nil, err
	}
	all, err := svc.invitationRepo.FindOutgoing(orgID, userID)
	if err != nil {
		return nil, err
	}
	if opts.SortBy == "" {
		opts.SortBy = dto.InvitationSortByDateCreated
	}
	if opts.SortOrder == "" {
		opts.SortOrder = dto.InvitationSortOrderAsc
	}
	sorted := svc.sort(all, opts.SortBy, opts.SortOrder)
	paged, totalElements, totalPages := svc.paginate(sorted, opts.Page, opts.Size)
	mapped, err := svc.invitationMapper.mapMany(paged, userID)
	if err != nil {
		return nil, err
	}
	return &dto.InvitationList{
		Data:          mapped,
		TotalPages:    totalPages,
		TotalElements: totalElements,
		Page:          opts.Page,
		Size:          uint64(len(mapped)),
	}, nil
}

func (svc *InvitationService) ProbeOutgoing(orgID string, opts InvitationListOptions, userID string) (*dto.InvitationProbe, error) {
	org, err := svc.orgCache.Get(orgID)
	if err != nil {
		return nil, err
	}
	if err := svc.orgGuard.Authorize(userID, org, model.PermissionOwner); err != nil {
		return nil, err
	}
	all, err := svc.invitationRepo.FindOutgoing(orgID, userID)
	totalElements := uint64(len(all))
	if err != nil {
		return nil, err
	}
	return &dto.InvitationProbe{
		TotalElements: totalElements,
		TotalPages:    (totalElements + opts.Size - 1) / opts.Size,
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
	if !strings.EqualFold(user.GetEmail(), invitation.GetEmail()) {
		return errorpkg.NewUserNotAllowedToAcceptInvitationError(user, invitation)
	}
	org, err := svc.orgCache.Get(invitation.GetOrganizationID())
	if err != nil {
		return err
	}
	for _, u := range org.GetMembers() {
		if u == userID {
			return errorpkg.NewUserAlreadyMemberOfOrganizationError(user, org)
		}
	}
	invitation.SetStatus(model.InvitationStatusAccepted)
	if err := svc.invitationRepo.Save(invitation); err != nil {
		return err
	}
	if err := svc.orgRepo.GrantUserPermission(invitation.GetOrganizationID(), userID, model.PermissionViewer); err != nil {
		return err
	}
	org, err = svc.orgRepo.Find(org.GetID())
	if err != nil {
		return err
	}
	if err := svc.orgSvc.sync(org); err != nil {
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
	if !strings.EqualFold(user.GetEmail(), invitation.GetEmail()) {
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

func (svc *InvitationService) IsValidSortBy(value string) bool {
	return value == "" ||
		value == dto.InvitationSortByEmail ||
		value == dto.InvitationSortByDateCreated ||
		value == dto.InvitationSortByDateModified
}

func (svc *InvitationService) IsValidSortOrder(value string) bool {
	return value == "" || value == dto.InvitationSortOrderAsc || value == dto.InvitationSortOrderDesc
}

func (svc *InvitationService) getValidOutboundEmails(emails []string, ownerEmail string, orgMembers []model.User, outgoing []model.Invitation) []string {
	var res []string
	for _, email := range emails {
		existing := false
		for _, u := range orgMembers {
			if email == u.GetEmail() {
				existing = true
				break
			}
		}
		for _, i := range outgoing {
			if email == i.GetEmail() && i.GetStatus() == model.InvitationStatusPending {
				existing = true
				break
			}
		}
		if !existing && !strings.EqualFold(email, ownerEmail) {
			res = append(res, email)
		}
	}
	return res
}

func (svc *InvitationService) sendEmails(invitations []model.Invitation, org model.Organization, userID string) error {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return err
	}
	for _, i := range invitations {
		variables := map[string]string{
			"USER_FULL_NAME":    user.GetFullName(),
			"ORGANIZATION_NAME": org.GetName(),
			"UI_URL":            svc.config.PublicUIURL,
		}
		_, err := svc.userRepo.FindByEmail(i.GetEmail())
		var templateName string
		if err == nil {
			templateName = "join-organization"
		} else {
			templateName = "signup-and-join-organization"
		}
		if err := svc.mailTmpl.Send(templateName, i.GetEmail(), variables); err != nil {
			return err
		}
	}
	return nil
}

func (svc *InvitationService) sort(data []model.Invitation, sortBy string, sortOrder string) []model.Invitation {
	if sortBy == dto.InvitationSortByEmail {
		sort.Slice(data, func(i, j int) bool {
			if sortOrder == dto.InvitationSortOrderDesc {
				return data[i].GetEmail() > data[j].GetEmail()
			} else {
				return data[i].GetEmail() < data[j].GetEmail()
			}
		})
		return data
	} else if sortBy == dto.InvitationSortByDateCreated {
		sort.Slice(data, func(i, j int) bool {
			a := helper.StringToTime(data[i].GetCreateTime())
			b := helper.StringToTime(data[j].GetCreateTime())
			if sortOrder == dto.InvitationSortOrderDesc {
				return a.UnixMilli() > b.UnixMilli()
			} else {
				return a.UnixMilli() < b.UnixMilli()
			}
		})
		return data
	} else if sortBy == dto.InvitationSortByDateModified {
		sort.Slice(data, func(i, j int) bool {
			if data[i].GetUpdateTime() != nil && data[j].GetUpdateTime() != nil {
				a := helper.StringToTime(*data[i].GetUpdateTime())
				b := helper.StringToTime(*data[j].GetUpdateTime())
				if sortOrder == dto.InvitationSortOrderDesc {
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

func (svc *InvitationService) paginate(data []model.Invitation, page, size uint64) (pageData []model.Invitation, totalElements uint64, totalPages uint64) {
	totalElements = uint64(len(data))
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
	userRepo   *repo.UserRepo
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

func (mp *invitationMapper) mapOne(m model.Invitation, userID string) (*dto.Invitation, error) {
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
	return &dto.Invitation{
		ID:           m.GetID(),
		Owner:        mp.userMapper.mapOne(owner),
		Email:        m.GetEmail(),
		Organization: o,
		Status:       m.GetStatus(),
		CreateTime:   m.GetCreateTime(),
		UpdateTime:   m.GetUpdateTime(),
	}, nil
}

func (mp *invitationMapper) mapMany(invitations []model.Invitation, userID string) ([]*dto.Invitation, error) {
	res := make([]*dto.Invitation, 0)
	for _, invitation := range invitations {
		i, err := mp.mapOne(invitation, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}
