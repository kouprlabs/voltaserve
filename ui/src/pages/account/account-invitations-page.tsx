// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect } from 'react'
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
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, { SortBy, SortOrder } from '@/client/api/invitation'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import { incomingInvitationPaginationStorage } from '@/infra/pagination'
import { IconMoreVert } from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import prettyDate from '@/lib/helpers/pretty-date'
import userToString from '@/lib/helpers/user-to-string'
import usePagePagination from '@/lib/hooks/page-pagination'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/incoming-invitations'

const AccountInvitationsPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
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

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [dispatch, mutate])

  const handleAccept = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.accept(invitationId)
      await mutate()
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
      await InvitationAPI.decline(invitationId)
      await mutate()
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
      {list.data.length === 0 ? (
        <div
          className={cx('flex', 'items-center', 'justify-center', 'h-[300px]')}
        >
          <span>There are no invitations.</span>
        </div>
      ) : null}
      {list.data.length > 0 ? (
        <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
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
                    <Td>{i.owner ? userToString(i.owner) : ''}</Td>
                    <Td>
                      <Text noOfLines={1}>
                        {i.organization ? i.organization.name : ''}
                      </Text>
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
                            <MenuItem onClick={() => handleAccept(i.id)}>
                              Accept
                            </MenuItem>
                            <MenuItem
                              className={cx('text-red-500')}
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
          {list ? (
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
          ) : null}
        </div>
      ) : null}
    </>
  )
}

export default AccountInvitationsPage
