// Copyright (c) 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Avatar,
  Badge,
  Divider,
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
import {
  PagePagination,
  RelativeDate,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI from '@/client/console/console'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'

const ConsolePanelOrganization = () => {
  const { id } = useParams()
  const [userPage, setUserPage] = useState(1)
  const [workspacePage, setWorkspacePage] = useState(1)
  const [groupPage, setGroupPage] = useState(1)
  const {
    data: org,
    error: orgError,
    isLoading: orgIsLoading,
  } = ConsoleAPI.useGetOrganizationById({
    id,
  })
  const {
    data: userList,
    error: userListError,
    isLoading: userListIsLoading,
  } = ConsoleAPI.useListUsersByOrganization(
    {
      id,
      page: userPage,
      size: 5,
    },
    swrConfig(),
  )
  const {
    data: workspaceList,
    isLoading: workspaceListIsLoading,
    error: workspaceListError,
  } = ConsoleAPI.useListWorkspacesByOrganization(
    {
      id,
      page: workspacePage,
      size: 5,
    },
    swrConfig(),
  )
  const {
    data: groupList,
    error: groupListError,
    isLoading: groupListIsLoading,
  } = ConsoleAPI.useListGroupsByOrganization(
    {
      id,
      page: groupPage,
      size: 5,
    },
    swrConfig(),
  )
  const orgIsReady = org && !orgError
  const userListIsEmpty =
    userList && !userListError && userList.totalElements === 0
  const userListIsReady =
    userList && !userListError && userList.totalElements > 0
  const workspaceListIsEmpty =
    workspaceList && !workspaceListError && workspaceList.totalElements === 0
  const workspaceListIsReady =
    workspaceList && !workspaceListError && workspaceList.totalElements > 0
  const groupListIsEmpty =
    groupList && !groupListError && groupList.totalElements === 0
  const groupListIsReady =
    groupList && !groupListError && groupList.totalElements > 0

  return (
    <>
      {orgIsLoading ? <SectionSpinner /> : null}
      {orgError ? <SectionError text={errorToString(orgError)} /> : null}
      {orgIsReady ? (
        <>
          <Helmet>
            <title>{org.name}</title>
          </Helmet>
          <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
            <Heading className={cx('text-heading')} noOfLines={1}>
              {org.name}
            </Heading>
            <Grid gap={4} templateColumns="repeat(9, 1fr)">
              <GridItem>
                <div className={cx('relative', 'shrink-0')}>
                  <Avatar
                    name={org.name}
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
                {userListIsLoading ? (
                  <ListSkeleton header="Users">
                    <SectionSpinner />
                  </ListSkeleton>
                ) : null}
                {userListError ? (
                  <ListSkeleton header="Users">
                    <SectionError text={errorToString(userListError)} />
                  </ListSkeleton>
                ) : null}
                {userListIsEmpty ? (
                  <ListSkeleton header="Users">
                    <SectionPlaceholder text="There are no users." />
                  </ListSkeleton>
                ) : null}
                {userListIsReady ? (
                  <>
                    <Table>
                      <Thead>
                        <Tr>
                          <Th className={cx('p-0')}>
                            <div
                              className={cx(
                                'flex',
                                'items-center',
                                'justify-between',
                                'h-[50px]',
                              )}
                            >
                              <span className={cx('font-bold')}>Users</span>
                              <Spacer />
                              {userList.totalElements > 5 ? (
                                <PagePagination
                                  totalElements={userList.totalElements}
                                  totalPages={Math.ceil(
                                    userList.totalElements / 5,
                                  )}
                                  page={userPage}
                                  size={5}
                                  steps={[]}
                                  setPage={setUserPage}
                                  setSize={() => {}}
                                  paginationSize="xs"
                                  isFirstDisabled={true}
                                  isLastDisabled={true}
                                  isRewindDisabled={true}
                                  isFastForwardDisabled={true}
                                />
                              ) : null}
                            </div>
                          </Th>
                        </Tr>
                      </Thead>
                    </Table>
                    <Divider mb={4} />
                    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                      {userList.data.map((user) => (
                        <Stack
                          direction="row"
                          alignItems="center"
                          key={user.id}
                        >
                          <Avatar
                            name={user.username}
                            src={user.picture}
                            size="sm"
                            className={cx('w-[40px]', 'h-[40px]')}
                          />
                          <div className={cx('flex', 'flex-col', 'gap-0')}>
                            <div
                              className={cx('flex', 'items-center', 'gap-0.5')}
                            >
                              <Text fontWeight="bold" noOfLines={1}>
                                {user.username}
                              </Text>
                              <Badge variant="outline">{user.permission}</Badge>
                            </div>
                            <span className={cx('text-gray-500')}>
                              <RelativeDate date={new Date(user.createTime)} />
                            </span>
                          </div>
                        </Stack>
                      ))}
                    </div>
                  </>
                ) : null}
              </GridItem>
              <GridItem colSpan={3}>
                {workspaceListIsLoading ? (
                  <ListSkeleton header="Workspaces">
                    <SectionSpinner />
                  </ListSkeleton>
                ) : null}
                {workspaceListError ? (
                  <ListSkeleton header="Workspaces">
                    <SectionError text={errorToString(workspaceListError)} />
                  </ListSkeleton>
                ) : null}
                {workspaceListIsEmpty ? (
                  <ListSkeleton header="Workspaces">
                    <SectionPlaceholder text="There are no workspaces." />
                  </ListSkeleton>
                ) : null}
                {workspaceListIsReady ? (
                  <>
                    <Table>
                      <Thead>
                        <Tr>
                          <Th className={cx('p-0')}>
                            <div
                              className={cx(
                                'flex',
                                'items-center',
                                'justify-between',
                                'h-[50px]',
                              )}
                            >
                              <span className={cx('font-bold')}>
                                Workspaces
                              </span>
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
                                  paginationSize="xs"
                                  isFirstDisabled={true}
                                  isLastDisabled={true}
                                  isRewindDisabled={true}
                                  isFastForwardDisabled={true}
                                />
                              ) : null}
                            </div>
                          </Th>
                        </Tr>
                      </Thead>
                    </Table>
                    <Divider mb={4} />
                    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                      {workspaceList.data.map((workspace) => (
                        <div
                          key={workspace.id}
                          className={cx(
                            'flex',
                            'flex-row',
                            'items-center',
                            'gap-1',
                          )}
                        >
                          <Avatar
                            name={workspace.name}
                            size="sm"
                            className={cx('w-[40px]', 'h-[40px]')}
                          />
                          <div className={cx('flex', 'flex-col', 'gap-0')}>
                            <Text fontWeight="bold" noOfLines={1}>
                              {workspace.name}
                            </Text>
                            <span className={cx('text-gray-500')}>
                              <RelativeDate
                                date={new Date(workspace.createTime)}
                              />
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </>
                ) : null}
              </GridItem>
              <GridItem colSpan={3}>
                {groupListIsLoading ? (
                  <ListSkeleton header="Groups">
                    <SectionSpinner />
                  </ListSkeleton>
                ) : null}
                {groupListError ? (
                  <ListSkeleton header="Groups">
                    <SectionError text={errorToString(groupListError)} />
                  </ListSkeleton>
                ) : null}
                {groupListIsEmpty ? (
                  <ListSkeleton header="Groups">
                    <SectionPlaceholder text="There are no groups." />
                  </ListSkeleton>
                ) : null}
                {groupListIsReady ? (
                  <>
                    <Table>
                      <Thead>
                        <Tr>
                          <Th className={cx('p-0')}>
                            <div
                              className={cx(
                                'flex',
                                'items-center',
                                'justify-between',
                                'h-[50px]',
                              )}
                            >
                              <span className={cx('font-bold')}>Groups</span>
                              <Spacer />
                              {groupList.totalElements > 5 ? (
                                <PagePagination
                                  totalElements={groupList.totalElements}
                                  totalPages={Math.ceil(
                                    groupList.totalElements / 5,
                                  )}
                                  page={groupPage}
                                  size={5}
                                  steps={[]}
                                  setPage={setGroupPage}
                                  setSize={() => {}}
                                  paginationSize="xs"
                                  isFirstDisabled={true}
                                  isLastDisabled={true}
                                  isRewindDisabled={true}
                                  isFastForwardDisabled={true}
                                />
                              ) : null}
                            </div>
                          </Th>
                        </Tr>
                      </Thead>
                    </Table>
                    <Divider mb={4} />
                    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                      {groupList.data.map((group) => (
                        <div
                          key={group.id}
                          className={cx(
                            'flex',
                            'flex-row',
                            'items-center',
                            'gap-1',
                          )}
                        >
                          <Avatar
                            name={group.name}
                            size="sm"
                            className={cx('w-[40px]', 'h-[40px]')}
                          />
                          <div className={cx('flex', 'flex-col', 'gap-0')}>
                            <Text fontWeight="bold" noOfLines={1}>
                              {group.name}
                            </Text>
                            <span className={cx('text-gray-500')}>
                              <RelativeDate date={new Date(group.createTime)} />
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </>
                ) : null}
              </GridItem>
            </Grid>
          </div>
        </>
      ) : null}
    </>
  )
}

type ListSkeletonProps = {
  header: string
  children?: React.ReactNode
}

const ListSkeleton = ({ header, children }: ListSkeletonProps) => (
  <>
    <Table>
      <Thead>
        <Tr>
          <Th className={cx('p-0')}>
            <div
              className={cx(
                'flex',
                'items-center',
                'justify-between',
                'h-[50px]',
              )}
            >
              <span className={cx('font-bold')}>{header}</span>
              <Spacer />
            </div>
          </Th>
        </Tr>
      </Thead>
    </Table>
    <Divider mb={4} />
    {children}
  </>
)

export default ConsolePanelOrganization
