// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Avatar,
  Badge,
  Divider,
  Grid,
  GridItem,
  Heading,
  IconButton,
  IconButtonProps,
  Spacer,
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
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI from '@/client/console/console'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { truncateEnd } from '@/lib/helpers/truncate-end'
import truncateMiddle from '@/lib/helpers/truncate-middle'

const EditButton = (props: IconButtonProps) => (
  <IconButton
    icon={props.icon ? props.icon : <IconEdit />}
    className={cx('h-[40px]', 'w-[40px]')}
    {...props}
  />
)

const ConsolePanelUser = () => {
  const { id } = useParams()
  const [workspacePage, setWorkspacePage] = useState(1)
  const [groupPage, setGroupPage] = useState(1)
  const [orgsPage, setOrgPage] = useState(1)
  const { data: user, error: userError } = UserAPI.useGetById(id)
  const { data: groupList, error: groupError } = ConsoleAPI.useListGroupsByUser(
    {
      id,
      page: groupPage,
      size: 5,
    },
    swrConfig(),
  )
  const { data: orgList, error: orgError } =
    ConsoleAPI.useListOrganizationsByUser(
      {
        id,
        page: orgsPage,
        size: 5,
      },
      swrConfig(),
    )
  const { data: workspaceList, error: workspaceError } =
    ConsoleAPI.useListWorkspacesByUser(
      {
        id,
        page: workspacePage,
        size: 5,
      },
      swrConfig(),
    )
  const isUserLoading = !user && !userError
  const isUserError = !user && userError
  const isUserReady = user && !userError
  const isGroupListLoading = !groupList && !groupError
  const isGroupListError = !groupList && groupError
  const isGroupListEmpty =
    groupList && !groupError && groupList.totalElements === 0
  const isGroupListReady =
    groupList && !groupError && groupList.totalElements > 0
  const isOrgListLoading = !orgList && !orgError
  const isOrgListError = !orgList && orgError
  const isOrgListEmpty = orgList && !orgError && orgList.totalElements === 0
  const isOrgListReady = orgList && !orgError && orgList.totalElements > 0
  const isWorkspaceListLoading = !workspaceList && !workspaceError
  const isWorkspaceListError = !workspaceList && workspaceError
  const isWorkspaceListEmpty =
    workspaceList && !workspaceError && workspaceList.totalElements === 0
  const isWorkspaceListReady =
    workspaceList && !workspaceError && workspaceList.totalElements > 0

  return (
    <>
      {isUserLoading ? <SectionSpinner /> : null}
      {isUserError ? <SectionError text="Failed to load user." /> : null}
      {isUserReady ? (
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
                {user.picture ? (
                  <IconButton
                    icon={<IconClose />}
                    variant="solid"
                    colorScheme="red"
                    right="5px"
                    bottom="10px"
                    position="absolute"
                    zIndex={1000}
                    title="Delete picture"
                    aria-label="Delete picture"
                    isDisabled={true}
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
                            <EditButton
                              title="Edit full name"
                              aria-label="Edit full name"
                              isDisabled={true}
                            />
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
                                <span>
                                  {truncateMiddle(user.pendingEmail, 50)}
                                </span>
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
                              title="Edit email"
                              aria-label="Edit email"
                              isDisabled={true}
                            />
                          </>
                        ),
                      },
                      {
                        label: 'Force change password',
                        content: (
                          <EditButton
                            title="Force change password"
                            aria-label="Force change password"
                            icon={<IconSync />}
                            isDisabled={true}
                          />
                        ),
                      },
                    ],
                  },
                ]}
              />
            </GridItem>
            <GridItem colSpan={3}>
              {isOrgListLoading ? (
                <ListSekeleton header="Organizations">
                  <SectionSpinner />
                </ListSekeleton>
              ) : null}
              {isOrgListEmpty ? (
                <ListSekeleton header="Organizations">
                  <SectionPlaceholder text="There are no organizations." />
                </ListSekeleton>
              ) : null}
              {isOrgListError ? (
                <ListSekeleton header="Organizations">
                  <SectionError text="Failed to load organizations." />
                </ListSekeleton>
              ) : null}
              {isOrgListReady ? (
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
              {isWorkspaceListLoading ? (
                <ListSekeleton header="Workspaces">
                  <SectionSpinner />
                </ListSekeleton>
              ) : null}
              {isWorkspaceListError ? (
                <ListSekeleton header="Workspaces">
                  <SectionError text="Failed to load workspaces." />
                </ListSekeleton>
              ) : null}
              {isWorkspaceListEmpty ? (
                <ListSekeleton header="Workspaces">
                  <SectionPlaceholder text="There are no workspaces." />
                </ListSekeleton>
              ) : null}
              {isWorkspaceListReady ? (
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
              {isGroupListLoading ? (
                <ListSekeleton header="Groups">
                  <SectionSpinner />
                </ListSekeleton>
              ) : null}
              {isGroupListError ? (
                <ListSekeleton header="Groups">
                  <SectionError text="Failed to load groups." />
                </ListSekeleton>
              ) : null}
              {isGroupListEmpty ? (
                <ListSekeleton header="Groups">
                  <SectionPlaceholder text="There are no groups." />
                </ListSekeleton>
              ) : null}
              {isGroupListReady ? (
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
