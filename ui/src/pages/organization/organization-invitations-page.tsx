import { useCallback, useEffect, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import {
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
  Th,
  Thead,
  Tr,
  useToast,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, { SortBy, SortOrder } from '@/client/api/invitation'
import OrganizationAPI from '@/client/api/organization'
import { geEditorPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import OrganizationInviteMembers from '@/components/organization/organization-invite-members'
import OrganizationStatus from '@/components/organization/organization-status'
import prettyDate from '@/helpers/pretty-date'
import { outgoingInvitationPaginationStorage } from '@/infra/pagination'
import {
  IconMoreVert,
  IconSend,
  IconDelete,
  IconPersonAdd,
  SectionSpinner,
  PagePagination,
  usePagePagination,
} from '@/lib'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/outgoing-invitations'

const OrganizationInvitationsPage = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
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

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate])

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
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div className={cx('flex', 'flex-col', 'gap-1.5', 'items-center')}>
              <span>This organization has no invitations.</span>
              {geEditorPermission(org.permission) && (
                <Button
                  leftIcon={<IconPersonAdd />}
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
        <div className={cx('flex', 'flex-col', 'gap-3.5', 'py-3.5')}>
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
                    <OrganizationStatus value={i.status} />
                  </Td>
                  <Td>{prettyDate(i.createTime)}</Td>
                  <Td className={cx('text-right')}>
                    <Menu>
                      <MenuButton
                        as={IconButton}
                        icon={<IconMoreVert />}
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
                            icon={<IconDelete />}
                            className={cx('text-red-500')}
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
