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

const workspaces = {
  'workspaces': [
    {
      'id': '3v3L8OQ32VZy1',
      'name': 'My Workspace',
      'organization_id': 'mz4GymnAnXB2O',
      'storage_capacity': 100000000000,
      'root_id': 'nxKEbdA5Z3ezP',
      'bucket': '785db6bb56514527b6d9b4b284ca7d19',
      'create_time': '2024-07-22T22:35:05Z',
      'update_time': '2024-07-22T22:35:05.391560Z',
    },
    {
      'id': 'y6gW1nQo57Bxy',
      'name': 'My Workspace',
      'organization_id': 'aa4Vdo6KD27y0',
      'storage_capacity': 100000000000,
      'root_id': 'v3W0RV6AelNgg',
      'bucket': '0ef7bc71b75c45418ee90e39a9f12d3c',
      'create_time': '2024-07-23T06:21:14Z',
      'update_time': '2024-07-23T06:21:14.413159Z',
    },
    {
      'id': '8bmWRBRLPyLYO',
      'name': 'x',
      'organization_id': 'aa4Vdo6KD27y0',
      'storage_capacity': 1000000000,
      'root_id': 'oqAE1D16jJzEk',
      'bucket': 'bdb440af006140ceb71c57897e13defc',
      'create_time': '2024-07-24T16:56:34Z',
      'update_time': '2024-07-24T16:56:34.867934Z',
    },
  ],
}

const AdminPanelWorkspaces = () => {
  return (
    <>
      <Helmet>
        <title>Workspaces management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Workspaces management</Heading>
        <Stack direction="column" spacing={2}>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Workspace name</Th>
                <Th>Organisation</Th>
                <Th>Storage</Th>
                <Th>Actions</Th>
              </Tr>
            </Thead>
            <Tbody>
              {workspaces.workspaces.map((workspace) => (
                <Tr key={workspace.id}>
                  <Td>{workspace.name}</Td>
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

export default AdminPanelWorkspaces
