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
import { useLocation, useNavigate } from 'react-router-dom'
import {
  Button,
  Heading,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Stack,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import AdminApi, { WorkspaceManagementList } from '@/client/admin/admin'
import { adminWorkspacesPaginationStorage } from '@/infra/pagination'
import { IconChevronDown, IconChevronUp } from '@/lib/components/icons'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import prettyBytes from '@/lib/helpers/pretty-bytes'
import usePagePagination from '@/lib/hooks/page-pagination'

const AdminPanelWorkspaces = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<WorkspaceManagementList | undefined>(
    undefined,
  )
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminWorkspacesPaginationStorage(),
  })

  useEffect(() => {
    AdminApi.listWorkspaces({ page: page, size: size }).then((value) =>
      setList(value),
    )
  }, [page, size])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>Workspaces management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Workspaces management</Heading>
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
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((workspace) => (
                  <Tr key={workspace.id}>
                    <Td>{workspace.name}</Td>
                    <Td>
                      <Button>{workspace.organization.name}</Button>
                    </Td>
                    <Td>{prettyBytes(workspace.storageCapacity)}</Td>
                    <Td>
                      {new Date(workspace.createTime).toLocaleDateString()}
                    </Td>
                    <Td>{new Date(workspace.updateTime).toLocaleString()}</Td>
                    <Td>
                      <Menu>
                        {({ isOpen }) => (
                          <>
                            <MenuButton
                              isActive={isOpen}
                              as={Button}
                              rightIcon={
                                isOpen ? <IconChevronUp /> : <IconChevronDown />
                              }
                            >
                              Actions
                            </MenuButton>
                            <MenuList>
                              <MenuItem>Rename</MenuItem>
                              <MenuItem>Change quota</MenuItem>
                            </MenuList>
                          </>
                        )}
                      </Menu>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No workspaces found </div>
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

export default AdminPanelWorkspaces
