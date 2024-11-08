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
import { useParams } from 'react-router-dom'
import {
  Avatar,
  Badge,
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
import {
  Form,
  IconClose,
  IconEdit,
  IconSync,
  IconWarning,
  PagePagination,
  RelativeDate,
  SectionSpinner,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI, {
  GroupUserManagement,
  ListResponse,
  OrganizationUserManagement,
  WorkspaceUserManagement,
} from '@/client/console/console'
import UserAPI, { ConsoleUser } from '@/client/idp/user'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { truncateEnd } from '@/lib/helpers/truncate-end'
import truncateMiddle from '@/lib/helpers/truncate-middle'

const EditButton = (props: IconButtonProps) => (
  <IconButton
    disabled
    icon={props.icon ? props.icon : <IconEdit />}
    className={cx('h-[40px]', 'w-[40px]')}
    {...props}
  />
)

const ConsolePanelUser = () => {
  const [user, setUser] = useState<ConsoleUser>()
  const [orgList, setOrgList] =
    useState<ListResponse<OrganizationUserManagement>>()
  const [workspaceList, setWorkspaceList] =
    useState<ListResponse<WorkspaceUserManagement>>()
  const [groupList, setGroupList] =
    useState<ListResponse<GroupUserManagement>>()
  const { id } = useParams()
  const [workspacePage, setWorkspacePage] = useState(1)
  const [groupPage, setGroupPage] = useState(1)
  const [orgPage, setOrgPage] = useState(1)

  const userFetch = () => {
    if (id) {
      UserAPI.getUserById({ id }).then((value) => {
        setUser(value)
      })
    }
  }
  const groupsFetch = () => {
    if (id) {
      ConsoleAPI.getGroupsByUser({ id: id, page: groupPage, size: 5 }).then(
        (value) => {
          setGroupList(value)
        },
      )
    }
  }
  const organizationsFetch = () => {
    ConsoleAPI.getOrganizationsByUser({
      id: id,
      page: orgPage,
      size: 5,
    }).then((value) => {
      setOrgList(value)
    })
  }

  const workspacesFetch = () => {
    ConsoleAPI.getWorkspacesByUser({
      id: id,
      page: workspacePage,
      size: 5,
    }).then((value) => {
      setWorkspaceList(value)
    })
  }

  useEffect(() => {
    userFetch()
    groupsFetch()
    organizationsFetch()
    workspacesFetch()
  }, [])

  useEffect(() => {
    organizationsFetch()
  }, [orgPage])

  useEffect(() => {
    groupsFetch()
  }, [groupPage])

  useEffect(() => {
    workspacesFetch()
  }, [workspacePage])

  if (!user) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>User Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')} noOfLines={1}>
          {user.fullName}
        </Heading>
      </div>
      <Grid gap={4} templateColumns="repeat(9, 1fr)">
        <GridItem>
          <div className={cx('relative', 'shrink-0')}>
            <Avatar
              name={user.fullName}
              src={
                user.picture
                  ? getPictureUrlById(user.id, user.picture)
                  : undefined
              }
              size="2xl"
              className={cx(
                'w-[165px]',
                'h-[165px]',
                'border',
                'border-gray-300',
                'dark:border-gray-700',
              )}
            />
            {user.picture ? (
              <IconButton
                icon={<IconClose />}
                variant="solid"
                colorScheme="red"
                right="5px"
                bottom="10px"
                position="absolute"
                zIndex={1000}
                aria-label=""
                onClick={() => {
                  console.log('remove')
                }}
              />
            ) : null}
          </div>
        </GridItem>
        <GridItem colSpan={8}>
          <Form
            sections={[
              {
                title: 'Basics',
                rows: [
                  {
                    label: 'Full name',
                    content: (
                      <>
                        <span>{truncateEnd(user.fullName, 50)}</span>
                        <EditButton aria-label="Edit full name" />
                      </>
                    ),
                  },
                ],
              },
              {
                title: 'Credentials',
                rows: [
                  {
                    label: 'Email',
                    content: (
                      <>
                        {user.pendingEmail ? (
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
                                <IconWarning
                                  className={cx('text-yellow-400')}
                                />
                              </div>
                            </Tooltip>
                            <span>{truncateMiddle(user.pendingEmail, 50)}</span>
                          </div>
                        ) : null}
                        {!user.pendingEmail ? (
                          <span>
                            {truncateMiddle(
                              user.pendingEmail || user.email,
                              50,
                            )}
                          </span>
                        ) : null}
                        <EditButton
                          aria-label=""
                          onClick={() => {
                            console.log('edit email')
                          }}
                        />
                      </>
                    ),
                  },
                  {
                    label: 'Force change password',
                    content: (
                      <EditButton
                        aria-label="Force change password"
                        icon={<IconSync />}
                      />
                    ),
                  },
                ],
              },
            ]}
          />
        </GridItem>
        <GridItem colSpan={3}>
          {!orgList ? (
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
                        {orgList.totalElements > 5 ? (
                          <PagePagination
                            totalElements={orgList.totalElements}
                            totalPages={Math.ceil(orgList.totalElements / 5)}
                            page={orgPage}
                            size={5}
                            steps={[]}
                            setPage={setOrgPage}
                            setSize={() => {}}
                            uiSize="xs"
                            disableLastNav
                            disableMiddleNav
                          />
                        ) : null}
                      </Flex>
                    </Th>
                  </Tr>
                </Thead>
              </Table>
              <Divider mb={4} />
              <Stack>
                {orgList.data && orgList.data.length > 0 ? (
                  orgList.data.map((organization) => (
                    <Stack
                      direction="row"
                      alignItems="center"
                      key={organization.organizationId}
                    >
                      <Avatar name={organization.organizationName} />
                      <Stack direction="column">
                        <Stack direction="row" alignItems="center">
                          <Text fontWeight="bold" noOfLines={1}>
                            {organization.organizationName}
                          </Text>
                          <Badge colorScheme="green">
                            {organization.permission}
                          </Badge>
                        </Stack>
                        <Text fontSize="sm">
                          from:{' '}
                          <RelativeDate
                            date={new Date(organization.createTime)}
                          />
                        </Text>
                      </Stack>
                    </Stack>
                  ))
                ) : (
                  <Text>No organizations found.</Text>
                )}
              </Stack>
            </>
          )}
        </GridItem>
        <GridItem colSpan={3}>
          {!workspaceList ? (
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
                        {workspaceList.totalElements > 5 ? (
                          <PagePagination
                            totalElements={workspaceList.totalElements}
                            totalPages={Math.ceil(
                              workspaceList.totalElements / 5,
                            )}
                            page={workspacePage}
                            size={5}
                            steps={[]}
                            setPage={setWorkspacePage}
                            setSize={() => {}}
                            uiSize="xs"
                            disableLastNav
                            disableMiddleNav
                          />
                        ) : null}
                      </Flex>
                    </Th>
                  </Tr>
                </Thead>
              </Table>
              <Divider mb={4} />
              <Stack overflowX="auto">
                {workspaceList.data && workspaceList.data.length > 0 ? (
                  workspaceList.data.map((workspace) => (
                    <Stack
                      direction="row"
                      alignItems="center"
                      key={workspace.workspaceId}
                    >
                      <Avatar name={workspace.workspaceName} />
                      <Stack direction="column">
                        <Stack direction="row" alignItems="center">
                          <Text fontWeight="bold" noOfLines={1}>
                            {workspace.workspaceName}
                          </Text>
                          <Badge ml="1" colorScheme="green">
                            {workspace.permission}
                          </Badge>
                        </Stack>
                        <Text fontSize="sm">
                          from:{' '}
                          <RelativeDate date={new Date(workspace.createTime)} />
                        </Text>
                      </Stack>
                    </Stack>
                  ))
                ) : (
                  <Text>No workspaces found.</Text>
                )}
              </Stack>
            </>
          )}
        </GridItem>
        <GridItem colSpan={3}>
          {!groupList ? (
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
                        {groupList.totalElements > 5 ? (
                          <PagePagination
                            totalElements={groupList.totalElements}
                            totalPages={Math.ceil(groupList.totalElements / 5)}
                            page={groupPage}
                            size={5}
                            steps={[]}
                            setPage={setGroupPage}
                            setSize={() => {}}
                            uiSize="xs"
                            disableLastNav
                            disableMiddleNav
                          />
                        ) : null}
                      </Flex>
                    </Th>
                  </Tr>
                </Thead>
              </Table>
              <Divider mb={4} />
              <Stack>
                {groupList.data && groupList.data.length > 0 ? (
                  groupList.data.map((group) => (
                    <Stack
                      direction="row"
                      alignItems="center"
                      key={group.groupId}
                    >
                      <Avatar name={group.groupName} />
                      <Stack direction="column">
                        <Stack direction="row" alignItems="center">
                          <Text fontWeight="bold" noOfLines={1}>
                            {group.groupName}
                          </Text>
                          <Badge ml="1" colorScheme="green">
                            {group.permission}
                          </Badge>
                        </Stack>
                        <Text fontSize="sm">
                          from:{' '}
                          <RelativeDate date={new Date(group.createTime)} />
                        </Text>
                      </Stack>
                    </Stack>
                  ))
                ) : (
                  <Text>No groups found.</Text>
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
