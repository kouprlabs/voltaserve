import { Invitation } from '@/client/api/invitation'
import { Notification } from '@/client/api/notification'
import NewInvitationItem from './new-Invitation-item'

type NotificationItemProps = {
  notification: Notification
}

const NotificationItem = ({ notification }: NotificationItemProps) => {
  if (notification.type === 'new_invitation') {
    const body: Invitation = notification.body as Invitation
    return <NewInvitationItem invitation={body} />
  }
  return null
}

export default NotificationItem
