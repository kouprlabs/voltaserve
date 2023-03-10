import { useCallback, useState } from 'react'
import { Button, Stack, Text, useToast } from '@chakra-ui/react'
import { useSWRConfig } from 'swr'
import InvitationAPI, { Invitation } from '@/api/invitation'

type NewInvitationProps = {
  invitation: Invitation
}

const NewInvitationItem = ({ invitation }: NewInvitationProps) => {
  const { mutate } = useSWRConfig()
  const toast = useToast()
  const [isAcceptLoading, setIsAcceptLoading] = useState(false)
  const [isDeclineLoading, setIsDeclineLoading] = useState(false)

  const handleAccept = useCallback(
    async (invitationId: string) => {
      try {
        setIsAcceptLoading(true)
        await InvitationAPI.accept(invitationId)
        mutate('/notifications')
        mutate('/invitations/get_incoming')
        mutate('/organizations')
        toast({
          title: 'Invitation accepted',
          status: 'success',
          isClosable: true,
        })
      } finally {
        setIsAcceptLoading(false)
      }
    },
    [mutate, toast]
  )

  const handleDecline = useCallback(
    async (invitationId: string) => {
      try {
        setIsDeclineLoading(true)
        await InvitationAPI.decline(invitationId)
        mutate('/notifications')
        mutate('/invitations/get_incoming')
        toast({
          title: 'Invitation declined',
          status: 'info',
          isClosable: true,
        })
      } finally {
        setIsDeclineLoading(false)
      }
    },
    [mutate, toast]
  )

  return (
    <Stack direction="column">
      <Text>
        You have been invited by <b>{invitation.owner.fullName}</b> to join the
        organization <b>{invitation.organization.name}</b>.<br />
      </Text>
      <Stack direction="row" justifyContent="flex-end">
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
      </Stack>
    </Stack>
  )
}

export default NewInvitationItem
