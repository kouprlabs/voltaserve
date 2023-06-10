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
	g.Get("/", r.GetAll)
}

// GetAll godoc
//
//	@Summary		Get notifications
//	@Description	Get notifications
//	@Tags			Notifications
//	@Id				notification_get_all
//	@Produce		json
//	@Success		200	{array}	core.Notification
//	@Failure		500
//	@Router			/notifications [get]
func (r *NotificationRouter) GetAll(c *fiber.Ctx) error {
	userID := GetUserID(c)
	res, err := r.notificationSvc.GetAll(userID)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
