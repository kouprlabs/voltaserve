package service

import "voltaserve/repo"

type Notification struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

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

func (svc *NotificationService) GetAll(userID string) ([]*Notification, error) {
	user, err := svc.userRepo.Find(userID)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	notifications := make([]*Notification, 0)
	for _, inv := range invitations {
		v, err := svc.invitationMapper.mapInvitation(inv, userID)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &Notification{
			Type: "new_invitation",
			Body: &v,
		})
	}
	return notifications, nil
}
