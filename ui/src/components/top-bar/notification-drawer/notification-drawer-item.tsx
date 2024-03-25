import { Invitation } from '@/client/api/invitation'
import { Notification } from '@/client/api/notification'
import NotificationDrawerNewInvitationItem from './notification-drawer-new-Invitation-item'

export type NotificationDrawerItemProps = {
  notification: Notification
}

const NotificationDrawerItem = ({
  notification,
}: NotificationDrawerItemProps) => {
  if (notification.type === 'new_invitation') {
    const body: Invitation = notification.body as Invitation
    return <NotificationDrawerNewInvitationItem invitation={body} />
  }
  return null
}

export default NotificationDrawerItem
