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

  return (
    <>
      {!user && userError ? <SectionError text="Failed to load user." /> : null}
      {!user && !userError ? <SectionSpinner /> : null}
      {user && !userError ? (
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
              {!orgList && orgError ? (
                <SectionError text="Failed to load organizations." />
              ) : null}
              {!orgList && !orgError ? <SectionSpinner /> : null}
              {orgList && !orgError ? (
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
                  {orgList.totalElements > 0 ? (
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
                  ) : (
                    <SectionPlaceholder text="There are no organizations." />
                  )}
                </>
              ) : null}
            </GridItem>
            <GridItem colSpan={3}>
              {!workspaceList && workspaceError ? (
                <SectionError text="Failed to load workspaces." />
              ) : null}
              {!workspaceList && !workspaceError ? <SectionSpinner /> : null}
              {workspaceList && !workspaceError ? (
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
                  {workspaceList.totalElements > 0 ? (
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
                  ) : (
                    <SectionPlaceholder text="There are no workspaces." />
                  )}
                </>
              ) : null}
            </GridItem>
            <GridItem colSpan={3}>
              {!groupList && groupError ? (
                <SectionError text="Failed to load groups." />
              ) : null}
              {!groupList && !groupError ? <SectionSpinner /> : null}
              {groupList && !groupError ? (
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
                  {groupList.totalElements > 0 ? (
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
                              <Badge variant="outline">
                                {group.permission}
                              </Badge>
                            </div>
                            <span className={cx('text-gray-500')}>
                              <RelativeDate date={new Date(group.createTime)} />
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <SectionPlaceholder text="There are no groups." />
                  )}
                </>
              ) : null}
            </GridItem>
          </Grid>
        </>
      ) : null}
    </>
  )
}

export default ConsolePanelUser
