package router

import (
	"voltaserve/service"

	"github.com/gofiber/fiber/v2"
)

type NotificationRouter struct {
	notificationSvc *service.NotificationService
}

func NewNotificationRouter() *NotificationRouter {
	return &NotificationRouter{
		notificationSvc: service.NewNotificationService(),
	}
}

func (r *NotificationRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.List)
}

// List godoc
//
//	@Summary		List
//	@Description	List
//	@Tags			Notifications
//	@Id				notifications_list
//	@Produce		json
//	@Success		200	{array}	service.Notification
//	@Failure		500
//	@Router			/notifications [get]
func (r *NotificationRouter) List(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.notificationSvc.List(userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
