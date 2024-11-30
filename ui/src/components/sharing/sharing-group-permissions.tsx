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
import { useParams } from 'react-router-dom'
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
import { IconDelete, SectionPlaceholder, SectionSpinner, Text } from '@koupr/ui'
import cx from 'classnames'
import FileAPI, { GroupPermission } from '@/client/api/file'
import { geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'

const SharingGroupPermissions = () => {
  const { fileId } = useParams()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [revokedPermission, setRevokedPermission] = useState<string>()
  const { data: file } = FileAPI.useGet(selection[0], swrConfig())
  const { data: permissions, mutate: mutatePermissions } =
    FileAPI.useGetGroupPermissions(
      file && geOwnerPermission(file.permission) ? file.id : undefined,
      swrConfig(),
    )
  const isSingleSelection = selection.length === 1

  const handleRevokePermission = useCallback(
    async (permission: GroupPermission) => {
      try {
        setRevokedPermission(permission.id)
        await FileAPI.revokeGroupPermission({
          ids: selection,
          groupId: permission.group.id,
        })
        await mutateFiles?.()
        if (isSingleSelection) {
          await mutatePermissions()
        }
      } finally {
        setRevokedPermission(undefined)
      }
    },
    [fileId, selection, isSingleSelection, mutateFiles, mutatePermissions],
  )

  return (
    <>
      {!permissions ? <SectionSpinner /> : null}
      {permissions && permissions.length === 0 ? (
        <SectionPlaceholder text="Not shared with any groups." height="auto" />
      ) : null}
      {permissions && permissions.length > 0 ? (
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
