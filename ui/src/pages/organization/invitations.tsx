import { useCallback, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Badge,
  Button,
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
  VStack,
} from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, { Invitation, InvitationStatus } from '@/api/invitation'
import { swrConfig } from '@/api/options'
import OrganizationAPI from '@/api/organization'
import { geEditorPermission } from '@/api/permission'
import {
  IconDotsVertical,
  IconSend,
  IconTrash,
  IconUserPlus,
} from '@/components/common/icon'
import LoadingSpinner from '@/components/common/loading-spinner'
import OrganizationInviteMembers from '@/components/organization/invite-members'
import variables from '@/theme/variables'
import prettyDate from '@/helpers/pretty-date'

type StatusProps = {
  value: InvitationStatus
}

const Status = ({ value }: StatusProps) => {
  let colorScheme
  if (value === 'accepted') {
    colorScheme = 'green'
  } else if (value === 'declined') {
    colorScheme = 'red'
  }
  return <Badge colorScheme={colorScheme}>{value}</Badge>
}

const OrganizationInvitationsPage = () => {
  const params = useParams()
  const id = params.id as string
  const toast = useToast()
  const { data: org, error: orgError } = OrganizationAPI.useGetById(
    id,
    swrConfig()
  )
  const {
    data: invitations,
    error: invitationsError,
    mutate,
  } = InvitationAPI.useGetOutgoing(id, swrConfig())
  const [isInviteMembersModalOpen, setIsInviteMembersModalOpen] =
    useState(false)

  const handleResend = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.resend(invitationId)
      toast({
        title: 'Invitation resent',
        status: 'success',
        isClosable: true,
      })
    },
    [toast]
  )

  const handleDelete = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.delete(invitationId)
      mutate()
    },
    [mutate]
  )

  if (invitationsError || orgError) {
    return null
  }

  if (!invitations || !org) {
    return <LoadingSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{org.name}</title>
      </Helmet>
      {invitations && invitations.length === 0 ? (
        <>
          <Center h="300px">
            <VStack spacing={variables.spacing}>
              <Text>This organization has no invitations.</Text>
              {geEditorPermission(org.permission) && (
                <Button
                  leftIcon={<IconUserPlus />}
                  onClick={() => {
                    setIsInviteMembersModalOpen(true)
                  }}
                >
                  Invite members
                </Button>
              )}
            </VStack>
          </Center>
          <OrganizationInviteMembers
            open={isInviteMembersModalOpen}
            id={org.id}
            onClose={() => setIsInviteMembersModalOpen(false)}
          />
        </>
      ) : null}
      {invitations && invitations.length > 0 ? (
        <Table variant="simple">
          <Thead>
            <Tr>
              <Th>Email</Th>
              <Th>Status</Th>
              <Th>Date</Th>
              <Th></Th>
            </Tr>
          </Thead>
          <Tbody>
            {invitations.map((e: Invitation) => (
              <Tr key={e.id}>
                <Td>{e.email}</Td>
                <Td>
                  <Status value={e.status} />
                </Td>
                <Td>{prettyDate(e.updateTime || e.createTime)}</Td>
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
                        {e.status === 'pending' && (
                          <MenuItem
                            icon={<IconSend />}
                            onClick={() => handleResend(e.id)}
                          >
                            Resend
                          </MenuItem>
                        )}
                        <MenuItem
                          icon={<IconTrash />}
                          color="red"
                          onClick={() => handleDelete(e.id)}
                        >
                          Delete
                        </MenuItem>
                      </MenuList>
                    </Portal>
                  </Menu>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      ) : null}
    </>
  )
}

export default OrganizationInvitationsPage
