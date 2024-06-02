package service

import "voltaserve/repo"

type NotificationService struct {
	userRepo         repo.UserRepo
	invitationRepo   repo.InvitationRepo
	invitationMapper *invitationMapper
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		userRepo:         repo.NewUserRepo(),
		invitationRepo:   repo.NewInvitationRepo(),
		invitationMapper: newInvitationMapper(),
	}
}

type Notification struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

func (svc *NotificationService) List(userID string) ([]*Notification, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	res := make([]*Notification, 0)
	for _, invitation := range invitations {
		i, err := svc.invitationMapper.mapOne(invitation, userID)
		if err != nil {
			return nil, err
		}
		res = append(res, &Notification{
			Type: "invitation",
			Body: &i,
		})
	}
	return res, nil
}
