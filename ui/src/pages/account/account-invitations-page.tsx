// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useEffect } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { Avatar, useToast } from '@chakra-ui/react'
import {
  DataTable,
  IconThumbDown,
  IconThumbUp,
  PagePagination,
  RelativeDate,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import InvitationAPI, { SortBy, SortOrder } from '@/client/api/invitation'
import { errorToString } from '@/client/error'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import { incomingInvitationPaginationStorage } from '@/infra/pagination'
import { getPictureUrlById } from '@/lib/helpers/picture'
import userToString from '@/lib/helpers/user-to-string'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/incoming-invitations'

const AccountInvitationsPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const toast = useToast()
  const {
    data: user,
    error: userError,
    isLoading: userIsLoading,
  } = UserAPI.useGet()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: incomingInvitationPaginationStorage(),
  })
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
    mutate,
  } = InvitationAPI.useGetIncoming(
    { page, size, sortBy: SortBy.DateCreated, sortOrder: SortOrder.Desc },
    swrConfig(),
  )
  const userIsReady = user && !userError
  const listIsEmpty = list && !listError && list.totalElements === 0
  const listIsReady = list && !listError && list.totalElements > 0

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

  return (
    <>
      {userIsLoading ? <SectionSpinner /> : null}
      {userError ? <SectionError text={errorToString(userError)} /> : null}
      {userIsReady ? (
        <>
          {listIsLoading ? <SectionSpinner /> : null}
          {listError ? <SectionError text={errorToString(listError)} /> : null}
          {listIsEmpty ? (
            <SectionPlaceholder text="There are no invitations." />
          ) : null}
          {listIsReady ? (
            <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
              <DataTable
                items={list.data}
                columns={[
                  {
                    title: 'From',
                    renderCell: (i) => (
                      <div
                        className={cx(
                          'flex',
                          'flex-row',
                          'gap-1.5',
                          'items-center',
                        )}
                      >
                        {i.owner && i.organization ? (
                          <>
                            <Avatar
                              name={i.owner.fullName}
                              src={
                                i.owner.picture
                                  ? getPictureUrlById(
                                      i.owner.id,
                                      i.owner.picture,
                                      {
                                        invitationId: i.id,
                                      },
                                    )
                                  : undefined
                              }
                              className={cx(
                                'border',
                                'border-gray-300',
                                'dark:border-gray-700',
                              )}
                            />
                            {i.owner ? userToString(i.owner) : ''}
                          </>
                        ) : null}
                      </div>
                    ),
                  },
                  {
                    title: 'Organization',
                    renderCell: (i) => (
                      <Text noOfLines={1}>
                        {i.organization ? i.organization.name : ''}
                      </Text>
                    ),
                  },
                  {
                    title: 'Date',
                    renderCell: (i) => (
                      <RelativeDate date={new Date(i.createTime)} />
                    ),
                  },
                ]}
                actions={[
                  {
                    label: 'Accept',
                    icon: <IconThumbUp />,
                    onClick: (i) => handleAccept(i.id),
                  },
                  {
                    label: 'Decline',
                    icon: <IconThumbDown />,
                    isDestructive: true,
                    onClick: (i) => handleDecline(i.id),
                  },
                ]}
              />
              <div className={cx('self-end')}>
                <PagePagination
                  totalElements={list.totalElements}
                  totalPages={list.totalPages}
                  page={page}
                  size={size}
                  steps={steps}
                  setPage={setPage}
                  setSize={setSize}
                />
              </div>
            </div>
          ) : null}
        </>
      ) : null}
    </>
  )
}

export default AccountInvitationsPage
