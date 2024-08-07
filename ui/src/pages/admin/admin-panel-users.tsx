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
import { useLocation, useNavigate } from 'react-router-dom'
import {
  Badge,
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
import UserAPI, { AdminUsersResponse } from '@/client/idp/user'
import { adminUsersPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import SectionSpinner from '@/lib/components/section-spinner'
import usePagePagination from '@/lib/hooks/page-pagination'

const AdminPanelUsers = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const [list, setList] = useState<AdminUsersResponse | undefined>(undefined)
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: adminUsersPaginationStorage(),
  })

  useEffect(() => {
    console.log('Przed requestem')
    UserAPI.getAllUsers({ page: page, size: size }).then((value) => {
      console.log('Przed wyciagniaciem danych')
      setList(value)
      console.log(value)
    })
  }, [page, size])

  if (!list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>Users management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Users management</Heading>
        {list && list.data.length > 0 ? (
          <Stack direction="column" spacing={2}>
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>Full name</Th>
                  <Th>Email</Th>
                  <Th>Email confirmed</Th>
                  <Th>Create time</Th>
                  <Th>Update time</Th>
                  <Th>Actions</Th>
                </Tr>
              </Thead>
              <Tbody>
                {list.data.map((user) => (
                  <Tr key={user.id}>
                    <Td>{user.fullName}</Td>
                    <Td>{user.email}</Td>
                    <Td>
                      <Badge
                        colorScheme={user.isEmailConfirmed ? 'green' : 'red'}
                      >
                        {user.isEmailConfirmed ? 'Confirmed' : 'Awaiting'}
                      </Badge>
                    </Td>
                    <Td>{new Date(user.createTime).toLocaleDateString()}</Td>
                    <Td>{new Date(user.updateTime).toLocaleString()}</Td>
                    <Td>
                      <Button>Manage</Button>
                      {user.isActive ? (
                        <Button colorScheme="red">Suspend</Button>
                      ) : (
                        <Button colorScheme="yellow">Unsuspend</Button>
                      )}
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Stack>
        ) : (
          <div> No users found </div>
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

export default AdminPanelUsers
