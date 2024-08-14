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
import AdminApi, { GroupManagementList } from '@/client/admin/admin'
import { adminGroupsPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'

const AdminPanelGroups = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<GroupManagementList | undefined>(undefined)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminGroupsPaginationStorage(),
  })

  useEffect(() => {
    AdminApi.listGroups({ page: page, size: size }).then((value) =>
      setList(value),
    )
  }, [page, size])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>Group management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Group management</Heading>
        {list && list.data.length > 0 ? (
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Group name</Th>
                  <Th>Organization</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((group) => (
                  <Tr key={group.id}>
                    <Td>{group.name}</Td>
                    <Td>
                      <Button>{group.organization.name}</Button>
                    </Td>
                    <Td>{new Date(group.createTime).toLocaleDateString()}</Td>
                    <Td>{new Date(group.updateTime).toLocaleString()}</Td>
                    <Td>
                      <Button>Rename</Button>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No groups found </div>
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

export default AdminPanelGroups
