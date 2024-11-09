// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useState } from 'react'
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
  IconEdit,
  PagePagination,
  RelativeDate,
  SectionError,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleAPI, { ConsoleWorkspace } from '@/client/console/console'
import { swrConfig } from '@/client/options'
import ConsoleRenameModal from '@/components/console/console-rename-modal'
import { consoleWorkspacesPaginationStorage } from '@/infra/pagination'
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
  const [isConfirmRenameOpen, setIsConfirmRenameOpen] = useState(false)
  const [currentName, setCurrentName] = useState<string>('')
  const [workspaceId, setWorkspaceId] = useState<string>()
  const {
    data: list,
    error,
    mutate,
  } = ConsoleAPI.useListOrSearchObject<ConsoleWorkspace>(
    'workspace',
    { page, size, query },
    swrConfig(),
  )

  const renameRequest = useCallback(
    async (name: string) => {
      if (workspaceId) {
        await ConsoleAPI.renameObject({ id: workspaceId, name }, 'workspace')
        await mutate()
      }
    },
    [workspaceId],
  )

  return (
    <>
      <Helmet>
        <title>Workspaces</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Workspaces</Heading>
        {!list && error ? (
          <SectionError text="Failed to load workspaces." />
        ) : null}
        {!list && !error ? <SectionSpinner /> : null}
        {list && list.totalElements > 0 && !error ? (
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
                label: 'Rename',
                icon: <IconEdit />,
                onClick: async (workspace) => {
                  setCurrentName(workspace.name)
                  setWorkspaceId(workspace.id)
                  setIsConfirmRenameOpen(true)
                },
              },
            ]}
          />
        ) : null}
        {list ? (
          <div className={cx('self-end')}>
            <PagePagination
              totalElements={list.totalElements}
              totalPages={Math.ceil(list.totalElements / size)}
              page={page}
              size={size}
              steps={steps}
              setPage={setPage}
              setSize={setSize}
            />
          </div>
        ) : null}
      </div>
      <ConsoleRenameModal
        currentName={currentName}
        isOpen={isConfirmRenameOpen}
        onClose={() => setIsConfirmRenameOpen(false)}
        onRequest={renameRequest}
      />
    </>
  )
}

export default ConsolePanelWorkspaces
