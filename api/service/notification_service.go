package service

import "voltaserve/repo"

type Notification struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

type NotificationService struct {
	userRepo         repo.CoreUserRepo
	invitationRepo   repo.CoreInvitationRepo
	invitationMapper *invitationMapper
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		userRepo:         repo.NewPostgresUserRepo(),
		invitationRepo:   repo.NewPostgresInvitationRepo(),
		invitationMapper: newInvitationMapper(),
	}
}

func (svc *NotificationService) GetAll(userId string) ([]*Notification, error) {
	user, err := svc.userRepo.Find(userId)
	if err != nil {
		return nil, err
	}
	invitations, err := svc.invitationRepo.GetIncoming(user.GetEmail())
	if err != nil {
		return nil, err
	}
	notifications := make([]*Notification, 0)
	for _, inv := range invitations {
		v, err := svc.invitationMapper.mapInvitation(inv, userId)
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
