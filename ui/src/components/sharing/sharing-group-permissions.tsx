// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useState } from 'react'
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  IconButton,
  Badge,
  Avatar,
} from '@chakra-ui/react'
import {
  IconDelete,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  Text,
} from '@koupr/ui'
import cx from 'classnames'
import { FileAPI, FileGroupPermission } from '@/client/api/file'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'

const SharingGroupPermissions = () => {
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [revokedPermission, setRevokedPermission] = useState<string>()
  const {
    data: permissions,
    error: permissionsError,
    isLoading: permissionsIsLoading,
    mutate: mutatePermissions,
  } = FileAPI.useGetGroupPermissions(selection[0], swrConfig())
  // prettier-ignore
  const permissionsIsEmpty = permissions && !permissionsError && permissions.length === 0
  // prettier-ignore
  const permissionsIsReady = permissions && !permissionsError && permissions.length > 0

  const handleRevokePermission = useCallback(
    async (permission: FileGroupPermission) => {
      try {
        setRevokedPermission(permission.id)
        await FileAPI.revokeGroupPermission({
          ids: selection,
          groupId: permission.group.id,
        })
        await mutateFiles?.()
        await mutatePermissions()
      } finally {
        setRevokedPermission(undefined)
      }
    },
    [selection, mutateFiles, mutatePermissions],
  )

  return (
    <>
      {permissionsIsLoading ? <SectionSpinner height="auto" /> : null}
      {permissionsError ? (
        <SectionError text={errorToString(permissionsError)} height="auto" />
      ) : null}
      {permissionsIsEmpty ? (
        <SectionPlaceholder text="Not shared with any groups." height="auto" />
      ) : null}
      {permissionsIsReady ? (
        <Table>
          <Thead>
            <Tr>
              <Th>Group</Th>
              <Th>Permission</Th>
              <Th />
            </Tr>
          </Thead>
          <Tbody>
            {permissions.map((permission) => (
              <Tr key={permission.id}>
                <Td className={cx('p-1')}>
                  <div
                    className={cx('flex', 'flex-row', 'items-center', 'gap-1')}
                  >
                    <Avatar
                      name={permission.group.name}
                      size="sm"
                      className={cx('w-[40px]', 'h-[40px]')}
                    />
                    <Text noOfLines={1}>{permission.group.name}</Text>
                  </div>
                </Td>
                <Td>
                  <Badge>{permission.permission}</Badge>
                </Td>
                <Td className={cx('text-end')}>
                  <IconButton
                    icon={<IconDelete />}
                    colorScheme="red"
                    title="Revoke group permission"
                    aria-label="Revoke group permission"
                    isLoading={revokedPermission === permission.id}
                    onClick={() => handleRevokePermission(permission)}
                  />
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      ) : null}
    </>
  )
}

export default SharingGroupPermissions
