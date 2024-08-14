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
import AdminApi, { OrganizationManagementList } from '@/client/admin/admin'
import { adminOrganizationsPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'

const AdminPanelOrganizations = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<OrganizationManagementList | undefined>(
    undefined,
  )
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminOrganizationsPaginationStorage(),
  })

  useEffect(() => {
    AdminApi.listOrganizations({ page: page, size: size }).then((value) =>
      setList(value),
    )
  }, [page, size])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>Organizations management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>
          Organizations management
        </Heading>
        {list && list.data.length > 0 ? (
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Workspace name</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((organization) => (
                  <Tr key={organization.id}>
                    <Td>{organization.name}</Td>
                    <Td>
                      {new Date(organization.createTime).toLocaleDateString()}
                    </Td>
                    <Td>
                      {new Date(organization.updateTime).toLocaleString()}
                    </Td>
                    <Td>
                      <Button>Rename</Button>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No organizations found </div>
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

export default AdminPanelOrganizations
