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
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom'
import {
  Badge,
  Heading,
  Stack,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Text,
  MenuButton,
  MenuList,
  MenuItem,
  Menu,
  Center,
  IconButton,
  Avatar,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import UserAPI, { ConsoleUsersResponse } from '@/client/idp/user'
import ConsoleConfirmationModal from '@/components/console/console-confirmation-modal'
import { consoleUsersPaginationStorage } from '@/infra/pagination'
import { getUserId } from '@/infra/token'
import { IconMoreVert } from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { decodeQuery } from '@/lib/helpers/query'
import usePagePagination from '@/lib/hooks/page-pagination'

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
    navigate,
    location,
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
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Full name</Th>
                  <Th>Email</Th>
                  <Th>Email confirmed</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th>Props</Th>
                  <Th></Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((user) => (
                  <Tr
                    style={{ cursor: 'pointer' }}
                    key={user.id}
                    onClick={(event) => {
                      if (
                        !(event.target instanceof HTMLButtonElement) &&
                        !(event.target instanceof HTMLSpanElement)
                      ) {
                        navigate(`/console/users/${user.id}`)
                      }
                    }}
                  >
                    <Td>
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
                        <Text noOfLines={1}>{user.fullName}</Text>
                      </div>
                    </Td>
                    <Td>
                      <Text noOfLines={1}>{user.email}</Text>
                    </Td>
                    <Td>
                      <Badge
                        colorScheme={user.isEmailConfirmed ? 'green' : 'red'}
                      >
                        {user.isEmailConfirmed ? 'Confirmed' : 'Awaiting'}
                      </Badge>
                    </Td>
                    <Td>
                      <Text>
                        {new Date(user.createTime).toLocaleDateString()}
                      </Text>
                    </Td>
                    <Td>
                      <Text>{new Date(user.updateTime).toLocaleString()}</Text>
                    </Td>
                    <Td>
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
                    </Td>
                    <Td>
                      {getUserId() === user.id ? (
                        <Badge colorScheme="red">It's you</Badge>
                      ) : (
                        <Center>
                          <Menu>
                            <MenuButton
                              as={IconButton}
                              icon={<IconMoreVert />}
                              variant="ghost"
                              aria-label=""
                            />
                            <MenuList>
                              {user.isActive ? (
                                <MenuItem
                                  onClick={async () => {
                                    setConfirmWindowAction('suspend')
                                    await suspendUser(user.id, user.email, true)
                                  }}
                                >
                                  Suspend
                                </MenuItem>
                              ) : (
                                <MenuItem
                                  onClick={async () => {
                                    setConfirmWindowAction('unsuspend')
                                    await suspendUser(
                                      user.id,
                                      user.email,
                                      false,
                                    )
                                  }}
                                >
                                  Unsuspend
                                </MenuItem>
                              )}
                              {user.isAdmin ? (
                                <MenuItem
                                  onClick={async () => {
                                    setConfirmWindowAction(
                                      'remove console rights from',
                                    )
                                    await makeAdminUser(
                                      user.id,
                                      user.email,
                                      false,
                                    )
                                  }}
                                >
                                  Deadmin
                                </MenuItem>
                              ) : (
                                <MenuItem
                                  onClick={async () => {
                                    setConfirmWindowAction(
                                      'grant console rights to',
                                    )
                                    await makeAdminUser(
                                      user.id,
                                      user.email,
                                      true,
                                    )
                                  }}
                                >
                                  Make Admin
                                </MenuItem>
                              )}
                            </MenuList>
                          </Menu>
                        </Center>
                      )}
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div>No users found.</div>
        )}
        {list ? (
          <PagePagination
            style={{ alignSelf: 'end' }}
            totalElements={list.totalElements}
            totalPages={Math.ceil(list.totalElements / size)}
            page={page}
            size={size}
            steps={steps}
            setPage={setPage}
            setSize={setSize}
          />
        ) : null}
      </div>
    </>
  )
}

export default ConsolePanelUsers
