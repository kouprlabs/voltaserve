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
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import { getPictureUrlById } from '@/lib/helpers/picture'

const ConsolePanelUser = () => {
  const { id } = useParams()
  const [workspacePage, setWorkspacePage] = useState(1)
  const [groupPage, setGroupPage] = useState(1)
  const [orgsPage, setOrgPage] = useState(1)
  const {
    data: user,
    error: userError,
    isLoading: userIsLoading,
  } = UserAPI.useGetById(id)
  const {
    data: groupList,
    error: groupListError,
    isLoading: groupIsListLoading,
  } = ConsoleAPI.useListGroupsByUser(
    {
      id,
      page: groupPage,
      size: 5,
    },
    swrConfig(),
  )
  const {
    data: orgList,
    error: orgListError,
    isLoading: orgListIsLoading,
  } = ConsoleAPI.useListOrganizationsByUser(
    {
      id,
      page: orgsPage,
      size: 5,
    },
    swrConfig(),
  )
  const {
    data: workspaceList,
    error: workspaceListError,
    isLoading: workspaceListIsLoading,
  } = ConsoleAPI.useListWorkspacesByUser(
    {
      id,
      page: workspacePage,
      size: 5,
    },
    swrConfig(),
  )
  const userIsReady = user && !userError
  // prettier-ignore
  const groupListIsEmpty = groupList && !groupListError && groupList.totalElements === 0
  // prettier-ignore
  const groupListIsReady = groupList && !groupListError && groupList.totalElements > 0
  const orgListIsEmpty = orgList && !orgListError && orgList.totalElements === 0
  const orgListIsReady = orgList && !orgListError && orgList.totalElements > 0
  // prettier-ignore
  const workspaceListIsEmpty = workspaceList && !workspaceListError && workspaceList.totalElements === 0
  // prettier-ignore
  const workspaceListIsReady = workspaceList && !workspaceListError && workspaceList.totalElements > 0

  return (
    <>
      {userIsLoading ? <SectionSpinner /> : null}
      {userError ? <SectionError text={errorToString(userError)} /> : null}
      {userIsReady ? (
        <>
          <Helmet>
            <title>{user.fullName}</title>
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
              </div>
            </GridItem>
            <GridItem colSpan={8}></GridItem>
            <GridItem colSpan={3}>
              {orgListIsLoading ? (
                <ListSekeleton header="Organizations">
                  <SectionSpinner />
                </ListSekeleton>
              ) : null}
              {orgListIsEmpty ? (
                <ListSekeleton header="Organizations">
                  <SectionPlaceholder text="There are no organizations." />
                </ListSekeleton>
              ) : null}
              {orgListError ? (
                <ListSekeleton header="Organizations">
                  <SectionError text={errorToString(orgListError)} />
                </ListSekeleton>
              ) : null}
              {orgListIsReady ? (
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
                              Organizations
                            </span>
                            <Spacer />
                            {orgList.totalElements > 5 ? (
                              <PagePagination
                                totalElements={orgList.totalElements}
                                totalPages={Math.ceil(
                                  orgList.totalElements / 5,
                                )}
                                page={orgsPage}
                                size={5}
                                steps={[]}
                                setPage={setOrgPage}
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
                    {orgList.data.map((organization) => (
                      <div
                        key={organization.organizationId}
                        className={cx(
                          'flex',
                          'flex-row',
                          'items-center',
                          'gap-1',
                        )}
                      >
                        <Avatar
                          name={organization.organizationName}
                          size="sm"
                          className={cx('w-[40px]', 'h-[40px]')}
                        />
                        <div className={cx('flex', 'flex-col', 'gap-0')}>
                          <div
                            className={cx('flex', 'items-center', 'gap-0.5')}
                          >
                            <Text fontWeight="bold" noOfLines={1}>
                              {organization.organizationName}
                            </Text>
                            <Badge variant="outline">
                              {organization.permission}
                            </Badge>
                          </div>
                          <span className={cx('text-gray-500')}>
                            <RelativeDate
                              date={new Date(organization.createTime)}
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
              {workspaceListIsLoading ? (
                <ListSekeleton header="Workspaces">
                  <SectionSpinner />
                </ListSekeleton>
              ) : null}
              {workspaceListError ? (
                <ListSekeleton header="Workspaces">
                  <SectionError text={errorToString(workspaceListError)} />
                </ListSekeleton>
              ) : null}
              {workspaceListIsEmpty ? (
                <ListSekeleton header="Workspaces">
                  <SectionPlaceholder text="There are no workspaces." />
                </ListSekeleton>
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
                        key={workspace.workspaceId}
                        className={cx(
                          'flex',
                          'flex-row',
                          'items-center',
                          'gap-1',
                        )}
                      >
                        <Avatar
                          name={workspace.workspaceName}
                          size="sm"
                          className={cx('w-[40px]', 'h-[40px]')}
                        />
                        <div className={cx('flex', 'flex-col', 'gap-0')}>
                          <div
                            className={cx('flex', 'items-center', 'gap-0.5')}
                          >
                            <Text fontWeight="bold" noOfLines={1}>
                              {workspace.workspaceName}
                            </Text>
                            <Badge variant="outline">
                              {workspace.permission}
                            </Badge>
                          </div>
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
              {groupIsListLoading ? (
                <ListSekeleton header="Groups">
                  <SectionSpinner />
                </ListSekeleton>
              ) : null}
              {groupListError ? (
                <ListSekeleton header="Groups">
                  <SectionError text={errorToString(groupListError)} />
                </ListSekeleton>
              ) : null}
              {groupListIsEmpty ? (
                <ListSekeleton header="Groups">
                  <SectionPlaceholder text="There are no groups." />
                </ListSekeleton>
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
                        key={group.groupId}
                        className={cx(
                          'flex',
                          'flex-row',
                          'items-center',
                          'gap-1',
                        )}
                      >
                        <Avatar
                          name={group.groupName}
                          size="sm"
                          className={cx('w-[40px]', 'h-[40px]')}
                        />
                        <div className={cx('flex', 'flex-col', 'gap-0')}>
                          <div
                            className={cx(
                              'flex',
                              'flex-row',
                              'items-center',
                              'gap-0.5',
                            )}
                          >
                            <Text fontWeight="bold" noOfLines={1}>
                              {group.groupName}
                            </Text>
                            <Badge variant="outline">{group.permission}</Badge>
                          </div>
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
        </>
      ) : null}
    </>
  )
}

type ListSekeletonProps = {
  header: string
  children?: React.ReactNode
}

const ListSekeleton = ({ header, children }: ListSekeletonProps) => (
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

export default ConsolePanelUser
