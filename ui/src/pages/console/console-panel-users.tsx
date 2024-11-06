// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect, useState } from 'react'
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
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import UserAPI, { ConsoleUsersResponse } from '@/client/idp/user'
import ConsoleConfirmationModal from '@/components/console/console-confirmation-modal'
import { consoleUsersPaginationStorage } from '@/infra/pagination'
import { getUserId } from '@/infra/token'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { decodeQuery } from '@/lib/helpers/query'
import relativeDate from '@/lib/helpers/relative-date'

const ConsolePanelUsers = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<ConsoleUsersResponse>()
  const [isSubmitting, setSubmitting] = useState(false)
  const [userId, setUserId] = useState<string>()
  const [userEmail, setUserEmail] = useState<string>()
  const [actionState, setActionState] = useState<boolean>()
  const [confirmSuspendWindowOpen, setConfirmSuspendWindowOpen] =
    useState(false)
  const [confirmAdminWindowOpen, setConfirmAdminWindowOpen] = useState(false)
  const [confirmWindowAction, setConfirmWindowAction] = useState<string>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleUsersPaginationStorage(),
  })

  const suspendUser = useCallback(
    async (
      id: string | null,
      email: string | null,
      suspend: boolean | null,
      confirm: boolean = false,
    ) => {
      if (confirm && userId && actionState !== undefined) {
        setSubmitting(true)
        try {
          await UserAPI.suspendUser({ id: userId, suspend: actionState })
        } finally {
          closeConfirmationWindow()
        }
      } else if (id && suspend !== null && email) {
        setConfirmSuspendWindowOpen(true)
        setActionState(suspend)
        setUserEmail(email)
        setUserId(id)
      }
    },
    [],
  )

  const makeAdminUser = useCallback(
    async (
      id: string | null,
      email: string | null,
      makeAdmin: boolean | null,
      confirm: boolean = false,
    ) => {
      if (confirm && userId && actionState !== undefined) {
        setSubmitting(true)
        try {
          await UserAPI.makeAdmin({ id: userId, makeAdmin: actionState })
        } finally {
          closeConfirmationWindow()
        }
      } else if (id && makeAdmin !== null && email) {
        setConfirmAdminWindowOpen(true)
        setActionState(makeAdmin)
        setUserEmail(email)
        setUserId(id)
      }
    },
    [],
  )

  const closeConfirmationWindow = () => {
    setUserId(undefined)
    setUserEmail(undefined)
    setActionState(undefined)
    setConfirmSuspendWindowOpen(false)
    setSubmitting(false)
    setConfirmAdminWindowOpen(false)
  }

  useEffect(() => {
    UserAPI.getAllUsers({ page: page, size: size, query: query }).then(
      (value) => {
        setList(value)
      },
    )
  }, [page, size, isSubmitting, query])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <ConsoleConfirmationModal
        isOpen={confirmSuspendWindowOpen}
        action={confirmWindowAction}
        target={userEmail}
        closeConfirmationWindow={closeConfirmationWindow}
        isSubmitting={isSubmitting}
        request={suspendUser}
      />
      <ConsoleConfirmationModal
        isOpen={confirmAdminWindowOpen}
        action={confirmWindowAction}
        target={userEmail}
        closeConfirmationWindow={closeConfirmationWindow}
        isSubmitting={isSubmitting}
        request={makeAdminUser}
      />
      <Helmet>
        <title>User Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>User Management</Heading>
        {list && list.data.length > 0 ? (
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
                isHiddenFn: (user) => getUserId() === user.id,
                onClick: async (user) => {
                  setConfirmWindowAction('suspend')
                  await suspendUser(user.id, user.email, true)
                },
              },
              {
                label: 'Unsuspend',
                icon: <IconHandshake />,
                isHiddenFn: (user) => user.isActive,
                onClick: async (user) => {
                  setConfirmWindowAction('unsuspend')
                  await suspendUser(user.id, user.email, false)
                },
              },
              {
                label: 'Make Admin',
                icon: <IconShield />,
                isHiddenFn: (user) => user.isAdmin,
                onClick: async (user) => {
                  setConfirmWindowAction('grant console rights to')
                  await makeAdminUser(user.id, user.email, true)
                },
              },
              {
                label: 'Demote Admin',
                icon: <IconRemoveModerator />,
                isHiddenFn: (user) => !user.isAdmin,
                onClick: async (user) => {
                  setConfirmWindowAction('remove console rights from')
                  await makeAdminUser(user.id, user.email, false)
                },
              },
            ]}
          />
        ) : (
          <div>No users found.</div>
        )}
        {list ? (
          <div className={cx('self-end')}>
            <PagePagination
              totalElements={list.totalElements}
              totalPages={Math.ceil(list.totalElements / size)}
              page={page}
              size={size}
              steps={steps}
              setPage={setPage}
              setSize={setSize}
            />
          </div>
        ) : null}
      </div>
    </>
  )
}

export default ConsolePanelUsers
