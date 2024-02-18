import { useCallback } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import {
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
  SectionSpinner,
  PagePagination,
  usePagePagination,
} from '@koupr/ui'
import classNames from 'classnames'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, { SortBy, SortOrder } from '@/client/api/invitation'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import prettyDate from '@/helpers/pretty-date'
import userToString from '@/helpers/user-to-string'
import { incomingInvitationPaginationStorage } from '@/infra/pagination'

const AccountInvitationsPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const toast = useToast()
  const { data: user, error: userError } = UserAPI.useGet()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: incomingInvitationPaginationStorage(),
  })
  const {
    data: list,
    error: invitationsError,
    mutate,
  } = InvitationAPI.useGetIncoming(
    { page, size, sortBy: SortBy.DateCreated, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

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
    [mutate, toast],
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
    [mutate, toast],
  )

  if (userError || invitationsError) {
    return null
  }
  if (!user || !list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{user.fullName}</title>
      </Helmet>
      {list.data.length === 0 && (
        <div
          className={classNames(
            'flex',
            'items-center',
            'justify-center',
            'h=[300px]',
          )}
        >
          <Text>There are no invitations.</Text>
        </div>
      )}
      {list.data.length > 0 && (
        <div className={classNames('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
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
              {list.data.length > 0 &&
                list.data.map((i) => (
                  <Tr key={i.id}>
                    <Td>{userToString(i.owner)}</Td>
                    <Td>{i.organization.name}</Td>
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
                            <MenuItem onClick={() => handleAccept(i.id)}>
                              Accept
                            </MenuItem>
                            <MenuItem
                              color="red"
                              onClick={() => handleDecline(i.id)}
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
      )}
    </>
  )
}

export default AccountInvitationsPage
