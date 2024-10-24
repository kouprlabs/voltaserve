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
import {
  Center,
  Heading,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Stack,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react'
import * as Yup from 'yup'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import ConsoleApi, { WorkspaceManagementList } from '@/client/console/console'
import ConsoleRenameModal from '@/components/console/console-rename-modal'
import { consoleWorkspacesPaginationStorage } from '@/infra/pagination'
import { IconMoreVert } from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import prettyBytes from '@/lib/helpers/pretty-bytes'
import { decodeQuery } from '@/lib/helpers/query'
import usePagePagination from '@/lib/hooks/page-pagination'

const ConsolePanelWorkspaces = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const location = useLocation()
  const query = decodeQuery(searchParams.get('q') as string)
  const [list, setList] = useState<WorkspaceManagementList>()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
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
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Workspace name</Th>
                  <Th>Organization</Th>
                  <Th>Quota</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th></Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((workspace) => (
                  <Tr key={workspace.id}>
                    <Td>
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
                    </Td>
                    <Td>
                      <ChakraLink
                        as={Link}
                        to={`/console/organizations/${workspace.organization.id}`}
                        className={cx('no-underline')}
                      >
                        <Text noOfLines={1}>{workspace.organization.name}</Text>
                      </ChakraLink>
                    </Td>
                    <Td>
                      <Text>{prettyBytes(workspace.storageCapacity)}</Text>
                    </Td>
                    <Td>
                      <Text>
                        {new Date(workspace.createTime).toLocaleDateString()}
                      </Text>
                    </Td>
                    <Td>
                      <Text>
                        {new Date(workspace.updateTime).toLocaleString()}
                      </Text>
                    </Td>
                    <Td>
                      <Center>
                        <Menu>
                          <MenuButton
                            as={IconButton}
                            icon={<IconMoreVert />}
                            variant="ghost"
                            aria-label=""
                          />
                          <MenuList>
                            <MenuItem
                              onClick={async () => {
                                await renameWorkspace(
                                  workspace.id,
                                  workspace.name,
                                  null,
                                )
                              }}
                            >
                              Rename
                            </MenuItem>
                          </MenuList>
                        </Menu>
                      </Center>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div>No workspaces found.</div>
        )}
        {list ? (
          <PagePagination
            style={{ alignSelf: 'end' }}
            totalElements={list.totalElements}
            totalPages={Math.ceil(list.totalElements / size)}
            page={page}
            size={size}
            steps={steps}
            setPage={setPage}
            setSize={setSize}
          />
        ) : null}
      </div>
    </>
  )
}

export default ConsolePanelWorkspaces
