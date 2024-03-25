import { useCallback, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import {
  Badge,
  Button,
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
import {
  IconDotsVertical,
  IconSend,
  IconTrash,
  IconUserPlus,
  SectionSpinner,
  PagePagination,
  usePagePagination,
} from '@koupr/ui'
import classNames from 'classnames'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, {
  InvitationStatus,
  SortBy,
  SortOrder,
} from '@/client/api/invitation'
import OrganizationAPI from '@/client/api/organization'
import { geEditorPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import OrganizationInviteMembers from '@/components/organization/organization-invite-members'
import prettyDate from '@/helpers/pretty-date'
import { outgoingInvitationPaginationStorage } from '@/infra/pagination'

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
  const navigate = useNavigate()
  const location = useLocation()
  const { id } = useParams()
  const toast = useToast()
  const { data: org, error: orgError } = OrganizationAPI.useGetById(
    id,
    swrConfig(),
  )
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: outgoingInvitationPaginationStorage(),
  })
  const {
    data: list,
    error: invitationsError,
    mutate,
  } = InvitationAPI.useGetOutgoing(
    {
      organizationId: id,
      page,
      size,
      sortBy: SortBy.DateCreated,
      sortOrder: SortOrder.Desc,
    },
    swrConfig(),
  )
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
    [toast],
  )

  const handleDelete = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.delete(invitationId)
      mutate()
    },
    [mutate],
  )

  if (invitationsError || orgError) {
    return null
  }

  if (!list || !org) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{org.name}</title>
      </Helmet>
      {list && list.data.length === 0 ? (
        <>
          <div
            className={classNames(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div
              className={classNames(
                'flex',
                'flex-col',
                'gap-1.5',
                'items-center',
              )}
            >
              <Text>This organization has no invitations.</Text>
              {geEditorPermission(org.permission) && (
                <Button
                  leftIcon={<IconUserPlus />}
                  onClick={() => {
                    setIsInviteMembersModalOpen(true)
                  }}
                >
                  Invite Members
                </Button>
              )}
            </div>
          </div>
          <OrganizationInviteMembers
            open={isInviteMembersModalOpen}
            id={org.id}
            onClose={() => setIsInviteMembersModalOpen(false)}
          />
        </>
      ) : null}
      {list && list.data.length > 0 ? (
        <div className={classNames('flex', 'flex-col', 'gap-3.5', 'py-3.5')}>
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
              {list.data.map((i) => (
                <Tr key={i.id}>
                  <Td>{i.email}</Td>
                  <Td>
                    <Status value={i.status} />
                  </Td>
                  <Td>{prettyDate(i.createTime)}</Td>
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
                          {i.status === 'pending' && (
                            <MenuItem
                              icon={<IconSend />}
                              onClick={() => handleResend(i.id)}
                            >
                              Resend
                            </MenuItem>
                          )}
                          <MenuItem
                            icon={<IconTrash />}
                            color="red"
                            onClick={() => handleDelete(i.id)}
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
          {list && (
            <PagePagination
              style={{ alignSelf: 'end' }}
              totalElements={list.totalElements}
              totalPages={list.totalPages}
              page={page}
              size={size}
              steps={steps}
              setPage={setPage}
              setSize={setSize}
            />
          )}
        </div>
      ) : null}
    </>
  )
}

export default OrganizationInvitationsPage
