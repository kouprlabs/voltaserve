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
import { useParams } from 'react-router-dom'
import {
  Avatar,
  Badge,
  Box,
  Center,
  Divider,
  Flex,
  Grid,
  GridItem,
  Heading,
  IconButton,
  IconButtonProps,
  Spacer,
  Stack,
  Table,
  Text,
  Th,
  Thead,
  Tooltip,
  Tr,
} from '@chakra-ui/react'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleApi, {
  GroupUserManagementList,
  OrganizationUserManagementList,
  WorkspaceUserManagementList,
} from '@/client/console/console'
import UserAPI, { ConsoleUser } from '@/client/idp/user'
import ConsoleConfirmationModal from '@/components/console/console-confirmation-modal'
import ConsoleRenameModal from '@/components/console/console-rename-modal'
import {
  IconClose,
  IconEdit,
  IconSync,
  IconWarning,
} from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import UserAvatar from '@/lib/components/user-avatar'

const EditButton = (props: IconButtonProps) => (
  <IconButton
    disabled
    icon={props.icon ? props.icon : <IconEdit />}
    className={cx('h-[40px]', 'w-[40px]')}
    {...props}
  />
)

enum actionChooser {
  Email = 'email',
  FullName = 'fullName',
  Password = 'password',
  Picture = 'picture',
}

const ConsolePanelUser = () => {
  const sectionClassName = cx('flex', 'flex-col', 'gap-1', 'py-1.5')
  const rowClassName = cx(
    'flex',
    'flex-row',
    'items-center',
    'gap-1',
    `h-[40px]`,
  )
  const [userData, setUserData] = useState<ConsoleUser>()
  const [organizationsData, setOrganizationsData] =
    useState<OrganizationUserManagementList>()
  const [workspacesData, setWorkspacesData] =
    useState<WorkspaceUserManagementList>()
  const [groupsData, setGroupsData] = useState<GroupUserManagementList>()
  const { id } = useParams()
  const [workspacesPage, setWorkspacesPage] = useState(1)
  const [groupsPage, setGroupsPage] = useState(1)
  const [organizationsPage, setOrganizationsPage] = useState(1)
  const [confirmWindowOpen, setConfirmWindowOpen] = useState(false)
  const [confirmRenameWindowOpen, setConfirmRenameWindowOpen] = useState(false)
  const [currentName, setCurrentName] = useState<string>('')
  const [isSubmitting, setSubmitting] = useState(false)
  const [userTarget, setUserTarget] = useState<string>('')
  const [userId, setUserId] = useState<string>()
  const [action, setAction] = useState<string>('')
  const formSchemaEmail = Yup.object().shape({
    name: Yup.string().required('Email is required').max(255),
  })
  const formSchemaFullName = Yup.object().shape({
    name: Yup.string().required('Full Name is required').max(255),
  })

  const verboseAction: { [key: string]: string } = {
    [actionChooser.Email]: 'Change email',
    [actionChooser.FullName]: 'Change Full name',
    [actionChooser.Password]: 'Force reset password',
    [actionChooser.Picture]: 'Remove the profile picture',
  }

  const verboseActionDescription: { [key: string]: string } = {
    [actionChooser.Password]: 'You are going to force reset password on ',
    [actionChooser.Picture]: 'You are going to remove the profile picture of ',
  }

  const formSchemaChooser: { [key: string]: Yup.ObjectSchema<object> } = {
    [actionChooser.Email]: formSchemaEmail,
    [actionChooser.FullName]: formSchemaFullName,
  }

  const closeConfirmationWindow = () => {
    setConfirmWindowOpen(false)
    setSubmitting(false)
    setConfirmRenameWindowOpen(false)
    setUserTarget('')
    setUserId('')
    setCurrentName('')
    setAction('')
  }

  const forceResetPassword = useCallback(
    async (
      id: string | null,
      target: string | null,
      _action: boolean | null,
      confirm: boolean = false,
    ) => {
      if (confirm && userId) {
        setSubmitting(true)
        try {
          await UserAPI.forceResetPassword({ id: userId })
        } finally {
          closeConfirmationWindow()
        }
      } else if (id && target) {
        setConfirmWindowOpen(true)
        setUserTarget(target)
        setUserId(id)
        setAction(actionChooser.Password)
      }
    },
    [userId, isSubmitting, action],
  )

  const removeUsersPicture = useCallback(
    async (
      id: string | null,
      target: string | null,
      _action: boolean | null,
      confirm: boolean = false,
    ) => {
      if (confirm && userId) {
        setSubmitting(true)
        try {
          await UserAPI.adminUpdateUserData(userId, { picture: null })
        } finally {
          closeConfirmationWindow()
        }
      } else if (id && target) {
        setConfirmWindowOpen(true)
        setUserTarget(target)
        setUserId(id)
        setAction(actionChooser.Picture)
      }
    },
    [userId, isSubmitting, action],
  )

  const renameUser = useCallback(
    async (
      id: string | null,
      currentName: string | null,
      newName: string | null,
      confirm: boolean = false,
    ) => {
      if (confirm && userId !== undefined && newName !== null) {
        try {
          setSubmitting(true)
          await UserAPI.adminUpdateUserData(userId, { fullName: newName })
        } finally {
          closeConfirmationWindow()
        }
      } else if (id !== null && currentName !== null && currentName !== '') {
        setConfirmRenameWindowOpen(true)
        setCurrentName(currentName)
        setUserId(id)
        setAction(actionChooser.FullName)
      }
    },
    [userId, isSubmitting, action],
  )

  const changeUserEmail = useCallback(
    async (
      id: string | null,
      currentName: string | null,
      newName: string | null,
      confirm: boolean = false,
    ) => {
      if (confirm && userId !== undefined && newName !== null) {
        try {
          setSubmitting(true)
          await UserAPI.adminUpdateUserData(userId, { email: newName })
        } finally {
          closeConfirmationWindow()
        }
      } else if (id !== null && currentName !== null && currentName !== '') {
        setConfirmRenameWindowOpen(true)
        setCurrentName(currentName)
        setUserId(id)
        setAction(actionChooser.Email)
      }
    },
    [userId, isSubmitting, action],
  )

  const functionConfirmChooser: {
    [key: string]: (
      id: string | null,
      target: string | null,
      action: boolean | null,
      confirm: boolean,
    ) => Promise<void>
  } = {
    [actionChooser.Password]: forceResetPassword,
    [actionChooser.Picture]: removeUsersPicture,
  }

  const functionInputChooser: {
    [key: string]: (
      id: string | null,
      currentName: string | null,
      newName: string | null,
      confirm: boolean,
    ) => Promise<void>
  } = {
    [actionChooser.Email]: changeUserEmail,
    [actionChooser.FullName]: renameUser,
  }

  const userFetch = () => {
    if (id) {
      UserAPI.getUserById({ id }).then((value) => {
        setUserData(value)
      })
    }
  }
  const groupsFetch = () => {
    if (id) {
      ConsoleApi.getGroupsByUser({ id: id, page: groupsPage, size: 5 }).then(
        (value) => {
          setGroupsData(value)
        },
      )
    }
  }
  const organizationsFetch = () => {
    ConsoleApi.getOrganizationsByUser({
      id: id,
      page: organizationsPage,
      size: 5,
    }).then((value) => {
      setOrganizationsData(value)
    })
  }

  const workspacesFetch = () => {
    ConsoleApi.getWorkspacesByUser({
      id: id,
      page: workspacesPage,
      size: 5,
    }).then((value) => {
      setWorkspacesData(value)
    })
  }

  useEffect(() => {
    userFetch()
    groupsFetch()
    organizationsFetch()
    workspacesFetch()
  }, [isSubmitting])

  useEffect(() => {
    organizationsFetch()
  }, [organizationsPage])

  useEffect(() => {
    groupsFetch()
  }, [groupsPage])

  useEffect(() => {
    workspacesFetch()
  }, [workspacesPage])

  if (!userData) {
    return <SectionSpinner />
  }

  return (
    <>
      <ConsoleConfirmationModal
        isOpen={confirmWindowOpen}
        action={verboseAction[action]}
        verbose={verboseActionDescription[action]}
        target={userTarget}
        closeConfirmationWindow={closeConfirmationWindow}
        isSubmitting={isSubmitting}
        request={functionConfirmChooser[action]}
      />
      <ConsoleRenameModal
        action={verboseAction[action]}
        target={currentName}
        closeConfirmationWindow={closeConfirmationWindow}
        isOpen={confirmRenameWindowOpen}
        isSubmitting={isSubmitting}
        previousName={currentName}
        object={'user'}
        formSchema={formSchemaChooser[action]}
        request={functionInputChooser[action]}
      />
      <Helmet>
        <title>User Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>{userData.fullName}</Heading>
      </div>
      <Grid gap={4} templateColumns="repeat(9, 1fr)">
        <GridItem>
          <div className={cx('relative', 'shrink-0')}>
            <UserAvatar
              name={userData.fullName}
              src={userData.picture}
              height={'165px'}
              size={'2xl'}
            />
            {userData.picture ? (
              <IconButton
                icon={<IconClose />}
                variant="solid"
                colorScheme="red"
                right="5px"
                bottom="10px"
                position="absolute"
                zIndex={1000}
                aria-label=""
                onClick={async () => {
                  await removeUsersPicture(
                    userData.id,
                    `${userData.fullName} (${userData.email})`,
                    true,
                  )
                }}
              />
            ) : null}
          </div>
        </GridItem>
        <GridItem colSpan={8}>
          <div className={cx('flex', 'flex-col', 'gap-0')}>
            <div className={sectionClassName}>
              <span className={cx('font-bold')}>Basics</span>
              <div className={cx(rowClassName)}>
                <span>Full name</span>
                <Spacer />
                <span>{userData.fullName}</span>
                <EditButton
                  aria-label=""
                  onClick={async () => {
                    await renameUser(userData.id, userData.fullName, null)
                  }}
                />
              </div>
            </div>
            <Divider />
            <div className={sectionClassName}>
              <span className={cx('font-bold')}>Credentials</span>
              <div className={cx(rowClassName)}>
                <span>Email</span>
                <Spacer />
                {userData.pendingEmail ? (
                  <div
                    className={cx(
                      'flex',
                      'flex-row',
                      'gap-0.5',
                      'items-center',
                    )}
                  >
                    <Tooltip label="Please check your inbox to confirm your email.">
                      <div
                        className={cx(
                          'flex',
                          'items-center',
                          'justify-center',
                          'cursor-default',
                        )}
                      >
                        <IconWarning className={cx('text-yellow-400')} />
                      </div>
                    </Tooltip>
                    <span>{userData.pendingEmail}</span>
                  </div>
                ) : null}
                {!userData.pendingEmail ? (
                  <span>{userData.pendingEmail || userData.email}</span>
                ) : null}
                <EditButton
                  aria-label=""
                  onClick={async () => {
                    await changeUserEmail(userData.id, userData.email, null)
                  }}
                />
              </div>
              <div className={cx(rowClassName)}>
                <span>Force change password</span>
                <Spacer />
                <EditButton
                  aria-label=""
                  icon={<IconSync />}
                  onClick={async () => {
                    await forceResetPassword(
                      userData.id,
                      `${userData.fullName} (${userData.email})`,
                      true,
                    )
                  }}
                />
              </div>
            </div>
          </div>
        </GridItem>
        <GridItem colSpan={3}>
          {!organizationsData ? (
            <SectionSpinner />
          ) : (
            <>
              <Table>
                <Thead>
                  <Tr>
                    <Th>
                      <Flex>
                        <span className={cx('font-bold')}>Organizations</span>
                        <Spacer />
                        {organizationsData.totalElements > 5 ? (
                          <Center>
                            <>
                              <PagePagination
                                totalElements={organizationsData.totalElements}
                                totalPages={Math.ceil(
                                  organizationsData.totalElements / 5,
                                )}
                                page={organizationsPage}
                                size={5}
                                steps={[]}
                                setPage={setOrganizationsPage}
                                setSize={() => {}}
                                uiSize="xs"
                                disableLastNav
                                disableMiddleNav
                              />
                            </>
                          </Center>
                        ) : null}
                      </Flex>
                    </Th>
                  </Tr>
                </Thead>
              </Table>
              <Divider mb={4} />
              <Stack>
                {organizationsData.data && organizationsData.data.length > 0 ? (
                  organizationsData.data.map((organization) => (
                    <Flex key={organization.organizationId}>
                      <Avatar name={organization.organizationName} />
                      <Box ml="3">
                        <Text fontWeight="bold">
                          {organization.organizationName}
                          <Badge ml="1" colorScheme="green">
                            {organization.permission}
                          </Badge>
                        </Text>
                        <Text fontSize="sm">
                          from:{' '}
                          {new Date(
                            organization.createTime,
                          ).toLocaleDateString()}
                        </Text>
                      </Box>
                    </Flex>
                  ))
                ) : (
                  <Text>No organizations found</Text>
                )}
              </Stack>
            </>
          )}
        </GridItem>
        <GridItem colSpan={3}>
          {!workspacesData ? (
            <SectionSpinner />
          ) : (
            <>
              <Table>
                <Thead>
                  <Tr>
                    <Th>
                      <Flex>
                        <span className={cx('font-bold')}>Workspaces</span>
                        <Spacer />
                        {workspacesData.totalElements > 5 ? (
                          <>
                            <PagePagination
                              totalElements={workspacesData.totalElements}
                              totalPages={Math.ceil(
                                workspacesData.totalElements / 5,
                              )}
                              page={workspacesPage}
                              size={5}
                              steps={[]}
                              setPage={setWorkspacesPage}
                              setSize={() => {}}
                              uiSize="xs"
                              disableLastNav
                              disableMiddleNav
                            />
                          </>
                        ) : null}
                      </Flex>
                    </Th>
                  </Tr>
                </Thead>
              </Table>
              <Divider mb={4} />
              <Stack overflowX="auto">
                {workspacesData.data && workspacesData.data.length > 0 ? (
                  workspacesData.data.map((workspace) => (
                    <Flex key={workspace.workspaceId}>
                      <Avatar name={workspace.workspaceName} />
                      <Box ml="3">
                        <Text fontWeight="bold">
                          {workspace.workspaceName}
                          <Badge ml="1" colorScheme="green">
                            {workspace.permission}
                          </Badge>
                        </Text>
                        <Text fontSize="sm">
                          from:{' '}
                          {new Date(workspace.createTime).toLocaleDateString()}
                        </Text>
                      </Box>
                    </Flex>
                  ))
                ) : (
                  <Text>No workspaces found</Text>
                )}
              </Stack>
            </>
          )}
        </GridItem>
        <GridItem colSpan={3}>
          {!groupsData ? (
            <SectionSpinner />
          ) : (
            <>
              <Table>
                <Thead>
                  <Tr>
                    <Th>
                      <Flex>
                        <span className={cx('font-bold')}>Groups</span>
                        <Spacer />
                        {groupsData.totalElements > 5 ? (
                          <>
                            <PagePagination
                              totalElements={groupsData.totalElements}
                              totalPages={Math.ceil(
                                groupsData.totalElements / 5,
                              )}
                              page={groupsPage}
                              size={5}
                              steps={[]}
                              setPage={setGroupsPage}
                              setSize={() => {}}
                              uiSize="xs"
                              disableLastNav
                              disableMiddleNav
                            />
                          </>
                        ) : null}
                      </Flex>
                    </Th>
                  </Tr>
                </Thead>
              </Table>
              <Divider mb={4} />
              <Stack>
                {groupsData.data && groupsData.data.length > 0 ? (
                  groupsData.data.map((group) => (
                    <Flex key={group.groupId}>
                      <Avatar name={group.groupName} />
                      <Box ml="3">
                        <Text fontWeight="bold">
                          {group.groupName}
                          <Badge ml="1" colorScheme="green">
                            {group.permission}
                          </Badge>
                        </Text>
                        <Text fontSize="sm">
                          from:{' '}
                          {new Date(group.createTime).toLocaleDateString()}
                        </Text>
                      </Box>
                    </Flex>
                  ))
                ) : (
                  <Text>No groups found</Text>
                )}
              </Stack>
            </>
          )}
        </GridItem>
      </Grid>
    </>
  )
}

export default ConsolePanelUser
