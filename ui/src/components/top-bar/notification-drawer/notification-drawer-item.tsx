import { Invitation } from '@/client/api/invitation'
import { Notification, NotificationType } from '@/client/api/notification'
import NotificationDrawerInvitation from './notification-drawer-Invitation'

export type NotificationDrawerItemProps = {
  notification: Notification
}

const NotificationDrawerItem = ({
  notification,
}: NotificationDrawerItemProps) => {
  if (notification.type === NotificationType.Invitation) {
    const body: Invitation = notification.body as Invitation
    return <NotificationDrawerInvitation invitation={body} />
  }
  return null
}

export default NotificationDrawerItem
