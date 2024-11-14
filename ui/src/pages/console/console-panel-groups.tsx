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
import ConsoleAPI, { ConsoleGroup } from '@/client/console/console'
import { swrConfig } from '@/client/options'
import ConsoleConfirmationModal, {
  ConsoleConfirmationModalRequest,
} from '@/components/console/console-confirmation-modal'
import { consoleGroupsPaginationStorage } from '@/infra/pagination'
import { getUserId } from '@/infra/token'
import { decodeQuery } from '@/lib/helpers/query'

const ConsolePanelGroups = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const location = useLocation()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleGroupsPaginationStorage(),
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
  } = ConsoleAPI.useListOrSearchObject<ConsoleGroup>(
    'group',
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
        <title>Groups</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Groups</Heading>
        {isListLoading ? <SectionSpinner /> : null}
        {isListError ? <SectionError text="Failed to load groups." /> : null}
        {isListEmpty ? (
          <SectionPlaceholder text="There are no groups." />
        ) : null}
        {isListReady ? (
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Name',
                renderCell: (group) => (
                  <div
                    className={cx(
                      'flex',
                      'flex-row',
                      'items-center',
                      'gap-1.5',
                    )}
                  >
                    <Avatar
                      name={group.name}
                      size="sm"
                      className={cx('w-[40px]', 'h-[40px]')}
                    />
                    <Text noOfLines={1}>{group.name}</Text>
                  </div>
                ),
              },
              {
                title: 'Organization',
                renderCell: (group) => (
                  <ChakraLink
                    as={Link}
                    to={`/console/organizations/${group.organization.id}`}
                    className={cx('no-underline')}
                  >
                    <Text noOfLines={1}>{group.organization.name}</Text>
                  </ChakraLink>
                ),
              },
              {
                title: 'Created',
                renderCell: (group) => (
                  <RelativeDate date={new Date(group.createTime)} />
                ),
              },
              {
                title: 'Updated',
                renderCell: (group) => (
                  <RelativeDate date={new Date(group.updateTime)} />
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
                      resourceId: workspace.id,
                      resourceType: 'group',
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
                      resourceId: workspace.id,
                      resourceType: 'group',
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
                  totalPages={Math.ceil(list.totalElements / size)}
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

export default ConsolePanelGroups
