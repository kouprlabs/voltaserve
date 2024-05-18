import { useCallback, useState } from 'react'
import { Button, useToast } from '@chakra-ui/react'
import cx from 'classnames'
import InvitationAPI, { Invitation } from '@/client/api/invitation'
import userToString from '@/helpers/user-to-string'
import { useAppSelector } from '@/store/hook'

export type NewInvitationProps = {
  invitation: Invitation
}

const NotificationDrawerNewInvitationItem = ({
  invitation,
}: NewInvitationProps) => {
  const toast = useToast()
  const mutateNotifications = useAppSelector(
    (state) => state.ui.notifications.mutate,
  )
  const mutateOrganizations = useAppSelector(
    (state) => state.ui.organizations.mutate,
  )
  const mutateIncomingInvitations = useAppSelector(
    (state) => state.ui.incomingInvitations.mutate,
  )
  const [isAcceptLoading, setIsAcceptLoading] = useState(false)
  const [isDeclineLoading, setIsDeclineLoading] = useState(false)

  const handleAccept = useCallback(
    async (invitationId: string) => {
      try {
        setIsAcceptLoading(true)
        await InvitationAPI.accept(invitationId)
        mutateNotifications?.()
        mutateIncomingInvitations?.()
        mutateOrganizations?.()
        toast({
          title: 'Invitation accepted',
          status: 'success',
          isClosable: true,
        })
      } finally {
        setIsAcceptLoading(false)
      }
    },
    [
      mutateNotifications,
      mutateOrganizations,
      mutateIncomingInvitations,
      toast,
    ],
  )

  const handleDecline = useCallback(
    async (invitationId: string) => {
      try {
        setIsDeclineLoading(true)
        await InvitationAPI.decline(invitationId)
        mutateNotifications?.()
        mutateIncomingInvitations?.()
        toast({
          title: 'Invitation declined',
          status: 'info',
          isClosable: true,
        })
      } finally {
        setIsDeclineLoading(false)
      }
    },
    [mutateNotifications, mutateIncomingInvitations, toast],
  )

  return (
    <div className={cx('flex', 'flex-col', 'gap-0.5')}>
      <div>
        You have been invited by{' '}
        <span className={cx('font-bold')}>
          {userToString(invitation.owner)}
        </span>{' '}
        to join the organization{' '}
        <span className={cx('font-bold')}>{invitation.organization.name}</span>
        .<br />
      </div>
      <div className={cx('flex', 'flex-row', 'gap-0.5', 'justify-end')}>
        <Button
          size="sm"
          variant="ghost"
          colorScheme="blue"
          disabled={isAcceptLoading || isDeclineLoading}
          isLoading={isAcceptLoading}
          onClick={() => handleAccept(invitation.id)}
        >
          Accept
        </Button>
        <Button
          size="sm"
          variant="ghost"
          colorScheme="red"
          disabled={isDeclineLoading || isAcceptLoading}
          isLoading={isDeclineLoading}
          onClick={() => handleDecline(invitation.id)}
        >
          Decline
        </Button>
      </div>
    </div>
  )
}

export default NotificationDrawerNewInvitationItem
