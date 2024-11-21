// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { ReactElement, useState } from 'react'
import {
  useLocation,
  useNavigate,
  useSearchParams,
  Link,
} from 'react-router-dom'
import { Badge, Heading, Avatar, Link as ChakraLink } from '@chakra-ui/react'
import {
  DataTable,
  IconFrontHand,
  IconHandshake,
  IconRemoveModerator,
  IconShield,
  PagePagination,
  RelativeDate,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  Text,
  usePageMonitor,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import ConsoleConfirmationModal, {
  ConsoleConfirmationModalRequest,
} from '@/components/console/console-confirmation-modal'
import { consoleUsersPaginationStorage } from '@/infra/pagination'
import { getUserId } from '@/infra/token'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { decodeQuery } from '@/lib/helpers/query'
import userToString from '@/lib/helpers/user-to-string'

const ConsolePanelUsers = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const navigate = useNavigate()
  const location = useLocation()
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false)
  const [isConfirmationDestructive, setIsConfirmationDestructive] =
    useState(false)
  const [confirmationHeader, setConfirmationHeader] = useState<ReactElement>()
  const [confirmationBody, setConfirmationBody] = useState<ReactElement>()
  const [confirmationRequest, setConfirmationRequest] =
    useState<ConsoleConfirmationModalRequest>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleUsersPaginationStorage(),
  })
  const {
    data: list,
    error: listError,
    isLoading: isListLoading,
    mutate,
  } = UserAPI.useList({ query, page, size }, swrConfig())
  const { hasPagination } = usePageMonitor({
    totalPages: list?.totalPages ?? 1,
    totalElements: list?.totalElements ?? 0,
    steps,
  })
  const isListError = !list && listError
  const isListEmpty = list && !listError && list.totalElements === 0
  const isListReady = list && !listError && list.totalElements > 0

  return (
    <>
      <Helmet>
        <title>Users</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Users</Heading>
        {isListLoading ? <SectionSpinner /> : null}
        {isListError ? <SectionError text="Failed to load users." /> : null}
        {isListEmpty ? <SectionPlaceholder text="There are no users." /> : null}
        {isListReady ? (
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Full name',
                renderCell: (user) => (
                  <div
                    className={cx(
                      'flex',
                      'flex-row',
                      'gap-1.5',
                      'items-center',
                    )}
                  >
                    <Avatar
                      name={user.fullName}
                      src={
                        user.picture
                          ? getPictureUrlById(user.id, user.picture)
                          : undefined
                      }
                      className={cx(
                        'border',
                        'border-gray-300',
                        'dark:border-gray-700',
                      )}
                    />
                    <ChakraLink
                      as={Link}
                      to={`/console/users/${user.id}`}
                      className={cx('no-underline')}
                    >
                      <Text noOfLines={1}>{user.fullName}</Text>
                    </ChakraLink>
                  </div>
                ),
              },
              {
                title: 'Email',
                renderCell: (user) => <Text noOfLines={1}>{user.email}</Text>,
              },
              {
                title: 'Email Confirmed',
                renderCell: (user) => (
                  <Badge colorScheme={user.isEmailConfirmed ? 'green' : 'red'}>
                    {user.isEmailConfirmed ? 'Confirmed' : 'Awaiting'}
                  </Badge>
                ),
              },
              {
                title: 'Created',
                renderCell: (user) => (
                  <RelativeDate date={new Date(user.createTime)} />
                ),
              },
              {
                title: 'Updated',
                renderCell: (user) => (
                  <RelativeDate date={new Date(user.updateTime)} />
                ),
              },
              {
                title: 'Props',
                renderCell: (user) => (
                  <div className={cx('flex', 'flex-row', 'gap-0.5')}>
                    {user.isAdmin ? (
                      <Badge mr="1" fontSize="0.8em" colorScheme="blue">
                        Admin
                      </Badge>
                    ) : null}
                    {user.isActive ? (
                      <Badge mr="1" fontSize="0.8em" colorScheme="green">
                        Active
                      </Badge>
                    ) : (
                      <Badge mr="1" fontSize="0.8em" colorScheme="gray">
                        Inactive
                      </Badge>
                    )}
                    {getUserId() === user.id ? (
                      <Badge colorScheme="red">It's you</Badge>
                    ) : null}
                  </div>
                ),
              },
            ]}
            actions={[
              {
                label: 'Suspend',
                icon: <IconFrontHand />,
                isDestructive: true,
                isHiddenFn: (user) =>
                  getUserId() === user.id || user.isAdmin || !user.isActive,
                onClick: async (user) => {
                  setConfirmationHeader(<>Suspend User</>)
                  setConfirmationBody(
                    <>
                      Are you sure you want to suspend{' '}
                      <span className={cx('font-bold')}>
                        {userToString(user)}
                      </span>
                      ?
                    </>,
                  )
                  setConfirmationRequest(() => async () => {
                    await UserAPI.suspend(user.id, { suspend: true })
                    await mutate()
                  })
                  setIsConfirmationDestructive(true)
                  setIsConfirmationOpen(true)
                },
              },
              {
                label: 'Unsuspend',
                icon: <IconHandshake />,
                isHiddenFn: (user) => user.isActive,
                onClick: async (user) => {
                  setConfirmationHeader(<>Unsuspend User</>)
                  setConfirmationBody(
                    <>
                      Are you sure you want to unsuspend{' '}
                      <span className={cx('font-bold')}>
                        {userToString(user)}
                      </span>
                      ?
                    </>,
                  )
                  setConfirmationRequest(() => async () => {
                    await UserAPI.suspend(user.id, { suspend: false })
                    await mutate()
                  })
                  setIsConfirmationDestructive(false)
                  setIsConfirmationOpen(true)
                },
              },
              {
                label: 'Make Admin',
                icon: <IconShield />,
                isHiddenFn: (user) => user.isAdmin,
                onClick: async (user) => {
                  setConfirmationHeader(<>Make Admin</>)
                  setConfirmationBody(
                    <>
                      Are you sure you want to make{' '}
                      <span className={cx('font-bold')}>
                        {userToString(user)}
                      </span>{' '}
                      admin?
                    </>,
                  )
                  setConfirmationRequest(() => async () => {
                    await UserAPI.makeAdmin(user.id, { makeAdmin: true })
                    await mutate()
                  })
                  setIsConfirmationDestructive(false)
                  setIsConfirmationOpen(true)
                },
              },
              {
                label: 'Demote Admin',
                icon: <IconRemoveModerator />,
                isDestructive: true,
                isHiddenFn: (user) => !user.isAdmin,
                onClick: async (user) => {
                  setConfirmationHeader(<>Demote Admin</>)
                  setConfirmationBody(
                    <>
                      Are you sure you want to demote{' '}
                      <span className={cx('font-bold')}>
                        {userToString(user)}
                      </span>
                      ?
                    </>,
                  )
                  setConfirmationRequest(() => async () => {
                    await UserAPI.makeAdmin(user.id, { makeAdmin: false })
                    await mutate()
                    if (getUserId() === user.id) {
                      navigate('/sign-out')
                    }
                  })
                  setIsConfirmationDestructive(true)
                  setIsConfirmationOpen(true)
                },
              },
            ]}
            pagination={
              hasPagination ? (
                <PagePagination
                  totalElements={list.totalElements}
                  totalPages={list.totalPages}
                  page={page}
                  size={size}
                  steps={steps}
                  setPage={setPage}
                  setSize={setSize}
                />
              ) : undefined
            }
          />
        ) : null}
      </div>
      {confirmationHeader && confirmationBody && confirmationRequest ? (
        <ConsoleConfirmationModal
          header={confirmationHeader}
          body={confirmationBody}
          isDestructive={isConfirmationDestructive}
          isOpen={isConfirmationOpen}
          onClose={() => setIsConfirmationOpen(false)}
          onRequest={confirmationRequest}
        />
      ) : null}
    </>
  )
}

export default ConsolePanelUsers
