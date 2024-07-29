// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
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

const organizations = {
  'organizations': [
    {
      'id': 'mz4GymnAnXB2O',
      'name': 'My Organization',
      'create_time': '2024-07-22T22:35:04Z',
      'update_time': '2024-07-22T22:35:04Z',
    },
    {
      'id': 'aa4Vdo6KD27y0',
      'name': 'My Organization',
      'create_time': '2024-07-23T06:21:13Z',
      'update_time': '2024-07-23T06:21:13Z',
    },
    {
      'id': 'LWl2Kl3LXr4Q7',
      'name': 'test',
      'create_time': '2024-07-24T17:07:19Z',
      'update_time': '2024-07-24T17:07:19Z',
    },
    {
      'id': 'kxlY3zYVD6WBO',
      'name': 'My organization asdasdadsdfasfda',
      'create_time': '2024-07-24T17:10:58Z',
      'update_time': '2024-07-24T17:11:50Z',
    },
  ],
}

const AdminPanelOrganizations = () => {
  return (
    <>
      <Helmet>
        <title>Organizations management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>
          Organizations management
        </Heading>
        <Stack direction="column" spacing={2}>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Organization name</Th>
                <Th>User</Th>
                <Th>Groups</Th>
                <Th>Actions</Th>
              </Tr>
            </Thead>
            <Tbody>
              {organizations.organizations.map((organization) => (
                <Tr key={organization.id}>
                  <Td>{organization.name}</Td>
                  <Td>
                    <Button>Show owner</Button>
                  </Td>
                  <Td>
                    <Button>Show groups</Button>
                  </Td>
                  <Td>
                    <Button>Manage</Button>
                    <Button colorScheme="red">Suspend</Button>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </Stack>
      </div>
    </>
  )
}

export default AdminPanelOrganizations
