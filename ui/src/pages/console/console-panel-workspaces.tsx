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
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleApi, { WorkspaceManagementList } from '@/client/console/console'
import ConsoleRenameModal from '@/components/console/console-rename-modal'
import { consoleWorkspacesPaginationStorage } from '@/infra/pagination'
import prettyBytes from '@/lib/helpers/pretty-bytes'
import { decodeQuery } from '@/lib/helpers/query'

const ConsolePanelWorkspaces = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const location = useLocation()
  const query = decodeQuery(searchParams.get('q') as string)
  const [list, setList] = useState<WorkspaceManagementList>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: consoleWorkspacesPaginationStorage(),
  })
  const [confirmRenameWindowOpen, setConfirmRenameWindowOpen] = useState(false)
  const [isSubmitting, setSubmitting] = useState(false)
  const [currentName, setCurrentName] = useState<string>('')
  const [workspaceId, setWorkspaceId] = useState<string>()
  const formSchema = Yup.object().shape({
    name: Yup.string().required('Name is required').max(255),
  })

  const renameWorkspace = useCallback(
    async (
      id: string | null,
      currentName: string | null,
      newName: string | null,
      confirm: boolean = false,
    ) => {
      if (confirm && workspaceId !== undefined && newName !== null) {
        try {
          setSubmitting(true)
          await ConsoleApi.renameObject(
            { id: workspaceId, name: newName },
            'workspace',
          )
        } finally {
          closeConfirmationWindow()
        }
      } else if (id !== null && currentName !== null && currentName !== '') {
        setConfirmRenameWindowOpen(true)
        setCurrentName(currentName)
        setWorkspaceId(id)
      }
    },
    [],
  )

  const closeConfirmationWindow = () => {
    setConfirmRenameWindowOpen(false)
    setSubmitting(false)
    setCurrentName('')
    setWorkspaceId(undefined)
  }

  useEffect(() => {
    if (query && query.length >= 3) {
      ConsoleApi.searchObject('workspace', {
        page: page,
        size: size,
        query: query,
      }).then((value) => setList(value))
    } else {
      ConsoleApi.listWorkspaces({ page: page, size: size, query: query }).then(
        (value) => setList(value),
      )
    }
  }, [page, size, isSubmitting, query])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <ConsoleRenameModal
        closeConfirmationWindow={closeConfirmationWindow}
        isOpen={confirmRenameWindowOpen}
        isSubmitting={isSubmitting}
        previousName={currentName}
        object="workspace"
        formSchema={formSchema}
        request={renameWorkspace}
      />
      <Helmet>
        <title>Workspace Management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Workspace Management</Heading>
        {list && list.data.length > 0 ? (
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
                  await renameWorkspace(workspace.id, workspace.name, null)
                },
              },
            ]}
          />
        ) : (
          <div>No workspaces found.</div>
        )}
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
    </>
  )
}

export default ConsolePanelWorkspaces
