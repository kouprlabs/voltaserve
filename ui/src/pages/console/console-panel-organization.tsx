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
  Box,
  Divider,
  Flex,
  Grid,
  GridItem,
  Heading,
  Spacer,
  Stack,
  Table,
  Text,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react'
import { PagePagination, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleApi, {
  GroupManagementList,
  OrganizationManagement,
  UserOrganizationManagementList,
  WorkspaceManagementList,
} from '@/client/console/console'

const ConsolePanelOrganization = () => {
  const [organizationData, setOrganizationData] =
    useState<OrganizationManagement>()
  const [usersData, setUsersData] = useState<UserOrganizationManagementList>()
  const [workspacesData, setWorkspacesData] =
    useState<WorkspaceManagementList>()
  const [groupsData, setGroupsData] = useState<GroupManagementList>()
  const { id } = useParams()
  const [usersPage, setUsersPage] = useState(1)
  const [workspacesPage, setWorkspacesPage] = useState(1)
  const [groupsPage, setGroupsPage] = useState(1)

  const organizationFetch = () => {
    if (id) {
      ConsoleApi.getOrganizationById({ id }).then((value) => {
        setOrganizationData(value)
      })
    }
  }

  const usersFetch = () => {
    if (id) {
      ConsoleApi.getUsersByOrganization({
        id: id,
        page: usersPage,
        size: 5,
      }).then((value) => {
        setUsersData(value)
      })
    }
  }

  const workspacesFetch = () => {
    if (id) {
      ConsoleApi.getWorkspacesByOrganization({
        id: id,
        page: workspacesPage,
        size: 5,
      }).then((value) => {
        setWorkspacesData(value)
      })
    }
  }

  const groupsFetch = () => {
    if (id) {
      ConsoleApi.getGroupsByOrganization({
        id: id,
        page: workspacesPage,
        size: 5,
      }).then((value) => {
        setGroupsData(value)
      })
    }
  }

  useEffect(() => {
    organizationFetch()
    usersFetch()
    workspacesFetch()
    groupsFetch()
  }, [])

  useEffect(() => {
    usersFetch()
  }, [usersPage])

  useEffect(() => {
    groupsFetch()
  }, [groupsPage])

  useEffect(() => {
    workspacesFetch()
  }, [workspacesPage])

  if (!organizationData) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>Organization Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')} noOfLines={1}>
          {organizationData.name}
        </Heading>
        <Grid gap={4} templateColumns="repeat(9, 1fr)">
          <GridItem>
            <div className={cx('relative', 'shrink-0')}>
              <Avatar
                name={organizationData.name}
                size="2xl"
                className={cx(
                  'w-[165px]',
                  'h-[165px]',
                  'border',
                  'border-gray-300',
                  'dark:border-gray-700',
                )}
              />
            </div>
          </GridItem>
          <GridItem colSpan={8}></GridItem>
          <GridItem colSpan={3}>
            {!usersData ? (
              <SectionSpinner />
            ) : (
              <>
                <Table>
                  <Thead>
                    <Tr>
                      <Th>
                        <Flex>
                          <span className={cx('font-bold')}>Users</span>
                          <Spacer />
                          {usersData.totalElements > 5 ? (
                            <PagePagination
                              totalElements={usersData.totalElements}
                              totalPages={Math.ceil(
                                usersData.totalElements / 5,
                              )}
                              page={usersPage}
                              size={5}
                              steps={[]}
                              setPage={setUsersPage}
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
                  {usersData.data && usersData.data.length > 0 ? (
                    usersData.data.map((user) => (
                      <Flex key={user.id}>
                        <Avatar name={user.username} src={user.picture} />
                        <Box ml="3">
                          <Text fontWeight="bold" noOfLines={1}>
                            {user.username}
                            <Badge ml="1" colorScheme="green">
                              {user.permission}
                            </Badge>
                          </Text>
                          <Text fontSize="sm">
                            created:{' '}
                            {new Date(user.createTime).toLocaleDateString()}
                          </Text>
                        </Box>
                      </Flex>
                    ))
                  ) : (
                    <Text>No organizations found.</Text>
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
                      <Flex key={workspace.id}>
                        <Avatar name={workspace.name} />
                        <Box ml="3">
                          <Text fontWeight="bold" noOfLines={1}>
                            {workspace.name}
                          </Text>
                          <Text fontSize="sm">
                            created:{' '}
                            {new Date(
                              workspace.createTime,
                            ).toLocaleDateString()}
                          </Text>
                        </Box>
                      </Flex>
                    ))
                  ) : (
                    <Text>No workspaces found.</Text>
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
                      <Flex key={group.id}>
                        <Avatar name={group.name} />
                        <Box ml="3">
                          <Text fontWeight="bold" noOfLines={1}>
                            {group.name}
                          </Text>
                          <Text fontSize="sm">
                            created:{' '}
                            {new Date(group.createTime).toLocaleDateString()}
                          </Text>
                        </Box>
                      </Flex>
                    ))
                  ) : (
                    <Text>No groups found.</Text>
                  )}
                </Stack>
              </>
            )}
          </GridItem>
        </Grid>
      </div>
    </>
  )
}

export default ConsolePanelOrganization
