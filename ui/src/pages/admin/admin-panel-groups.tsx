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

const groups = {
  'groups': [
    {
      'id': '6o7EbGD2rZ8JQ',
      'name': 'My Group',
      'organization_id': 'mz4GymnAnXB2O',
      'create_time': '2024-07-22T22:35:04Z',
      'update_time': '2024-07-22T22:35:04Z',
    },
    {
      'id': 'AdW38nLY5GlQn',
      'name': 'My Group',
      'organization_id': 'aa4Vdo6KD27y0',
      'create_time': '2024-07-23T06:21:13Z',
      'update_time': '2024-07-23T06:21:13Z',
    },
  ],
}

const AdminPanelGroups = () => {
  return (
    <>
      <Helmet>
        <title>Group management</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
        <Heading className={cx('text-heading')}>Group management</Heading>
        <Stack direction="column" spacing={2}>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Group name</Th>
                <Th>Organization</Th>
                <Th>Actions</Th>
              </Tr>
            </Thead>
            <Tbody>
              {groups.groups.map((group) => (
                <Tr key={group.id}>
                  <Td>{group.name}</Td>
                  <Td>
                    <Button>Show org</Button>
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

export default AdminPanelGroups
