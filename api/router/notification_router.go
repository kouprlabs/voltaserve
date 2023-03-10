package router

import (
	"voltaserve/core"

	"github.com/gofiber/fiber/v2"
)

type NotificationRouter struct {
	notificationSvc *core.NotificationService
}

func NewNotificationRouter() *NotificationRouter {
	return &NotificationRouter{
		notificationSvc: core.NewNotificationService(),
	}
}

func (r *NotificationRouter) AppendRoutes(g fiber.Router) {
	g.Get("/", r.GetAll)
}

// GetAll godoc
// @Summary     Get notifications
// @Description Get notifications
// @Tags        Notifications
// @Id          notification_get_all
// @Produce     json
// @Success     200 {array} core.Notification
// @Failure     500
// @Router      /notifications [get]
func (r *NotificationRouter) GetAll(c *fiber.Ctx) error {
	userId := GetUserId(c)
	res, err := r.notificationSvc.GetAll(userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
