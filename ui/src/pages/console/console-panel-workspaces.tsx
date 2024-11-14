// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ReactElement, useState } from 'react'
import {
  Link,
  useLocation,
  useNavigate,
  useSearchParams,
} from 'react-router-dom'
import { Avatar, Link as ChakraLink } from '@chakra-ui/react'
import { Heading } from '@chakra-ui/react'
import {
  DataTable,
  IconRemoveModerator,
  IconShield,
  PagePagination,
  RelativeDate,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  Text,
  usePageMonitor,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI, { ConsoleWorkspace } from '@/client/console/console'
import { swrConfig } from '@/client/options'
import ConsoleConfirmationModal, {
  ConsoleConfirmationModalRequest,
} from '@/components/console/console-confirmation-modal'
import { consoleWorkspacesPaginationStorage } from '@/infra/pagination'
import { getUserId } from '@/infra/token'
import prettyBytes from '@/lib/helpers/pretty-bytes'
import { decodeQuery } from '@/lib/helpers/query'

const ConsolePanelWorkspaces = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const location = useLocation()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleWorkspacesPaginationStorage(),
  })
  const [isConfirmationOpen, setIsConfirmationOpen] = useState(false)
  const [isConfirmationDestructive, setIsConfirmationDestructive] =
    useState(false)
  const [confirmationHeader, setConfirmationHeader] = useState<ReactElement>()
  const [confirmationBody, setConfirmationBody] = useState<ReactElement>()
  const [confirmationRequest, setConfirmationRequest] =
    useState<ConsoleConfirmationModalRequest>()
  const {
    data: list,
    error: listError,
    isLoading: isListLoading,
    mutate,
  } = ConsoleAPI.useListOrSearchObject<ConsoleWorkspace>(
    'workspace',
    { page, size, query },
    swrConfig(),
  )
  const { hasPagination } = usePageMonitor({
    totalPages: list?.totalPages ?? 1,
    totalElements: list?.totalElements ?? 0,
    steps,
  })
  const isListError = !list && listError
  const isListEmpty = list && !listError && list.totalElements === 0
  const isListReady = list && !listError && list.totalElements > 0

  return (
    <>
      <Helmet>
        <title>Workspaces</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Workspaces</Heading>
        {isListLoading ? <SectionSpinner /> : null}
        {isListError ? (
          <SectionError text="Failed to load workspaces." />
        ) : null}
        {isListEmpty ? (
          <SectionPlaceholder text="There are no workspaces." />
        ) : null}
        {isListReady ? (
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Name',
                renderCell: (workspace) => (
                  <div
                    className={cx(
                      'flex',
                      'flex-row',
                      'gap-1.5',
                      'items-center',
                    )}
                  >
                    <Avatar
                      name={workspace.name}
                      size="sm"
                      className={cx('w-[40px]', 'h-[40px]')}
                    />

                    <Text noOfLines={1}>{workspace.name}</Text>
                  </div>
                ),
              },
              {
                title: 'Organization',
                renderCell: (workspace) => (
                  <ChakraLink
                    as={Link}
                    to={`/console/organizations/${workspace.organization.id}`}
                    className={cx('no-underline')}
                  >
                    <Text noOfLines={1}>{workspace.organization.name}</Text>
                  </ChakraLink>
                ),
              },
              {
                title: 'Quota',
                renderCell: (workspace) => (
                  <Text>{prettyBytes(workspace.storageCapacity)}</Text>
                ),
              },
              {
                title: 'Created',
                renderCell: (workspace) => (
                  <RelativeDate date={new Date(workspace.createTime)} />
                ),
              },
              {
                title: 'Updated',
                renderCell: (workspace) => (
                  <RelativeDate date={new Date(workspace.updateTime)} />
                ),
              },
            ]}
            actions={[
              {
                label: 'Grant Owner Permission',
                icon: <IconShield />,
                onClick: async (workspace) => {
                  setConfirmationHeader(<>Grant Owner Permission</>)
                  setConfirmationBody(
                    <>
                      Are you sure you want to grant yourself owner permission
                      on{' '}
                      <span className={cx('font-bold')}>{workspace.name}</span>?
                    </>,
                  )
                  setConfirmationRequest(() => async () => {
                    await ConsoleAPI.grantUserPermission({
                      userId: getUserId(),
                      resourceId: workspace.organization.id,
                      resourceType: 'organization',
                      permission: 'owner',
                    })
                    await ConsoleAPI.grantUserPermission({
                      userId: getUserId(),
                      resourceId: workspace.id,
                      resourceType: 'workspace',
                      permission: 'owner',
                    })
                    await ConsoleAPI.grantUserPermission({
                      userId: getUserId(),
                      resourceId: workspace.rootId,
                      resourceType: 'file',
                      permission: 'owner',
                    })
                    await mutate()
                  })
                  setIsConfirmationDestructive(false)
                  setIsConfirmationOpen(true)
                },
              },
              {
                label: 'Revoke Permission',
                icon: <IconRemoveModerator />,
                isDestructive: true,
                onClick: async (workspace) => {
                  setConfirmationHeader(<>Revoke Permission</>)
                  setConfirmationBody(
                    <>
                      Are you sure you want to revoke your permission on{' '}
                      <span className={cx('font-bold')}>{workspace.name}</span>?
                    </>,
                  )
                  setConfirmationRequest(() => async () => {
                    await ConsoleAPI.revokeUserPermission({
                      userId: getUserId(),
                      resourceId: workspace.rootId,
                      resourceType: 'file',
                    })
                    await ConsoleAPI.revokeUserPermission({
                      userId: getUserId(),
                      resourceId: workspace.id,
                      resourceType: 'workspace',
                    })
                    await mutate()
                  })
                  setIsConfirmationDestructive(true)
                  setIsConfirmationOpen(true)
                },
              },
            ]}
            pagination={
              hasPagination ? (
                <PagePagination
                  totalElements={list.totalElements}
                  totalPages={list.totalPages}
                  page={page}
                  size={size}
                  steps={steps}
                  setPage={setPage}
                  setSize={setSize}
                />
              ) : undefined
            }
          />
        ) : null}
      </div>
      {confirmationHeader && confirmationBody && confirmationRequest ? (
        <ConsoleConfirmationModal
          header={confirmationHeader}
          body={confirmationBody}
          isDestructive={isConfirmationDestructive}
          isOpen={isConfirmationOpen}
          onClose={() => setIsConfirmationOpen(false)}
          onRequest={confirmationRequest}
        />
      ) : null}
    </>
  )
}

export default ConsolePanelWorkspaces
