// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import {
  Badge,
  Button,
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
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import UserAPI, { AdminUsersResponse } from '@/client/idp/user'
import { adminUsersPaginationStorage } from '@/infra/pagination'
import { getUserId } from '@/infra/token'
import { IconChevronDown, IconChevronUp } from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'
import AdminConfirmationModal from '@/pages/admin/admin-confirmation-modal'

const AdminPanelUsers = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<AdminUsersResponse | undefined>(undefined)
  const [isSubmitting, setSubmitting] = useState(false)
  const [userId, setUserId] = useState<string | undefined>(undefined)
  const [userEmail, setUserEmail] = useState<string | undefined>(undefined)
  const [actionState, setActionState] = useState<boolean | undefined>(undefined)
  const [confirmSuspendWindowOpen, setConfirmSuspendWindowOpen] =
    useState(false)
  const [confirmAdminWindowOpen, setConfirmAdminWindowOpen] = useState(false)
  const [confirmWindowAction, setConfirmWindowAction] = useState<
    string | undefined
  >(undefined)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminUsersPaginationStorage(),
  })

  const suspendUser = async (
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
  }

  const makeAdminUser = async (
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
  }

  const closeConfirmationWindow = () => {
    setUserId(undefined)
    setUserEmail(undefined)
    setActionState(undefined)
    setConfirmSuspendWindowOpen(false)
    setSubmitting(false)
    setConfirmAdminWindowOpen(false)
  }

  useEffect(() => {
    UserAPI.getAllUsers({ page: page, size: size }).then((value) => {
      setList(value)
    })
  }, [page, size, isSubmitting])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <AdminConfirmationModal
        isOpen={confirmSuspendWindowOpen}
        action={confirmWindowAction}
        target={userEmail}
        closeConfirmationWindow={closeConfirmationWindow}
        isSubmitting={isSubmitting}
        request={suspendUser}
      />
      <AdminConfirmationModal
        isOpen={confirmAdminWindowOpen}
        action={confirmWindowAction}
        target={userEmail}
        closeConfirmationWindow={closeConfirmationWindow}
        isSubmitting={isSubmitting}
        request={makeAdminUser}
      />
      {/*<Modal*/}
      {/*  isOpen={confirmSuspendWindowOpen}*/}
      {/*  onClose={() => {*/}
      {/*    closeConfirmationWindow()*/}
      {/*  }}*/}
      {/*>*/}
      {/*  <ModalOverlay />*/}
      {/*  <ModalContent>*/}
      {/*    <ModalHeader>Are You sure?</ModalHeader>*/}
      {/*    <ModalCloseButton />*/}
      {/*    <ModalBody>*/}
      {/*      You are going to {confirmWindowAction}*/}
      {/*      <br />*/}
      {/*      <Code children={userEmail} />*/}
      {/*      <br />*/}
      {/*      Please confirm this action*/}
      {/*    </ModalBody>*/}
      {/*    <ModalFooter>*/}
      {/*      <Button*/}
      {/*        type="button"*/}
      {/*        variant="outline"*/}
      {/*        colorScheme="blue"*/}
      {/*        disabled={isSubmitting}*/}
      {/*        onClick={() => {*/}
      {/*          closeConfirmationWindow()*/}
      {/*        }}*/}
      {/*      >*/}
      {/*        Cancel*/}
      {/*      </Button>*/}
      {/*      <Button*/}
      {/*        type="button"*/}
      {/*        variant="solid"*/}
      {/*        colorScheme="blue"*/}
      {/*        isLoading={isSubmitting}*/}
      {/*        onClick={async () => {*/}
      {/*          if (suspendAction !== undefined) {*/}
      {/*            await suspendUser(null, null, null, true)*/}
      {/*          } else if (adminAction !== undefined) {*/}
      {/*            await makeAdminUser(null, null, null, true)*/}
      {/*          }*/}
      {/*        }}*/}
      {/*      >*/}
      {/*        Confirm*/}
      {/*      </Button>*/}
      {/*    </ModalFooter>*/}
      {/*  </ModalContent>*/}
      {/*</Modal>*/}
      <Helmet>
        <title>Users management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Users management</Heading>
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
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((user) => (
                  <Tr
                    key={user.id}
                    onClick={(event) => {
                      if (
                        !(event.target instanceof HTMLButtonElement) &&
                        !(event.target instanceof HTMLSpanElement) &&
                        !(event.target instanceof HTMLParagraphElement)
                      ) {
                        navigate(`/admin/users/${user.id}`)
                      }
                    }}
                  >
                    <Td>
                      <Text>{user.fullName}</Text>
                    </Td>
                    <Td>
                      <Text>{user.email}</Text>
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
                        <Menu>
                          {({ isOpen }) => (
                            <>
                              <MenuButton
                                isActive={isOpen}
                                as={Button}
                                rightIcon={
                                  isOpen ? (
                                    <IconChevronUp />
                                  ) : (
                                    <IconChevronDown />
                                  )
                                }
                              >
                                Actions
                              </MenuButton>
                              <MenuList>
                                {user.isActive ? (
                                  <MenuItem
                                    onClick={async () => {
                                      setConfirmWindowAction('suspend')
                                      await suspendUser(
                                        user.id,
                                        user.email,
                                        true,
                                      )
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
                                        'remove admin rights from',
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
                                        'grant admin rights to',
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
                            </>
                          )}
                        </Menu>
                      )}
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No users found </div>
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

export default AdminPanelUsers
