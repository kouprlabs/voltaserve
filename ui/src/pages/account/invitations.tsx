import { useCallback } from 'react'
import {
  Center,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Portal,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
  useToast,
} from '@chakra-ui/react'
import { IconDotsVertical, SectionSpinner } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, { Invitation } from '@/api/invitation'
import { swrConfig } from '@/api/options'
import UserAPI from '@/api/user'

const AccountInvitationsPage = () => {
  const toast = useToast()
  const { data: user, error: userError } = UserAPI.useGet()
  const {
    data: invitations,
    error: invitationsError,
    mutate,
  } = InvitationAPI.useGetIncoming(swrConfig())

  const handleAccept = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.accept(invitationId)
      mutate()
      toast({
        title: 'Invitation accepted',
        status: 'success',
        isClosable: true,
      })
    },
    [mutate, toast]
  )

  const handleDecline = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.delete(invitationId)
      mutate()
      toast({
        title: 'Invitation declined',
        status: 'info',
        isClosable: true,
      })
    },
    [mutate, toast]
  )

  if (userError || invitationsError) {
    return null
  }
  if (!user || !invitations) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{user.fullName}</title>
      </Helmet>
      {invitations.length === 0 && (
        <Center h="300px">
          <Text>There are no invitations.</Text>
        </Center>
      )}
      {invitations.length > 0 && (
        <Table variant="simple">
          <Thead>
            <Tr>
              <Th>From</Th>
              <Th>Organization</Th>
              <Th>Date</Th>
              <Th></Th>
            </Tr>
          </Thead>
          <Tbody>
            {invitations.length > 0 &&
              invitations.map((e: Invitation) => (
                <Tr key={e.id}>
                  <Td>{e.owner.fullName}</Td>
                  <Td>{e.organization.name}</Td>
                  <Td>{e.updateTime || e.createTime}</Td>
                  <Td textAlign="right">
                    <Menu>
                      <MenuButton
                        as={IconButton}
                        icon={<IconDotsVertical />}
                        variant="ghost"
                        aria-label=""
                      />
                      <Portal>
                        <MenuList>
                          <MenuItem onClick={() => handleAccept(e.id)}>
                            Accept
                          </MenuItem>
                          <MenuItem
                            color="red"
                            onClick={() => handleDecline(e.id)}
                          >
                            Decline
                          </MenuItem>
                        </MenuList>
                      </Portal>
                    </Menu>
                  </Td>
                </Tr>
              ))}
          </Tbody>
        </Table>
      )}
    </>
  )
}

export default AccountInvitationsPage
