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
import FileAPI, { UserPermission } from '@/client/api/file'
import { geOwnerPermission } from '@/client/api/permission'
import WorkspaceAPI from '@/client/api/workspace'
import IdPUserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { useAppSelector } from '@/store/hook'

const SharingUserPermissions = () => {
  const { id: workspaceId, fileId } = useParams()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [revokedPermission, setRevokedPermission] = useState<string>()
  const { data: workspace } = WorkspaceAPI.useGet(workspaceId)
  const { data: file } = FileAPI.useGet(selection[0], swrConfig())
  const { data: me } = IdPUserAPI.useGet()
  const { data: permissions, mutate: mutatePermissions } =
    FileAPI.useGetUserPermissions(
      file && geOwnerPermission(file.permission) ? file.id : undefined,
      swrConfig(),
    )

  const handleRevokePermission = useCallback(
    async (permission: UserPermission) => {
      try {
        setRevokedPermission(permission.id)
        await FileAPI.revokeUserPermission({
          ids: selection,
          userId: permission.user.id,
        })
        await mutateFiles?.()
        await mutatePermissions()
      } finally {
        setRevokedPermission(undefined)
      }
    },
    [fileId, selection, mutateFiles, mutatePermissions],
  )

  return (
    <>
      {!permissions ? <SectionSpinner /> : null}
      {permissions && permissions.length === 0 ? (
        <SectionPlaceholder text="Not shared with any users." height="auto" />
      ) : null}
      {permissions && permissions.length > 0 ? (
        <Table>
          <Thead>
            <Tr>
              <Th>User</Th>
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
                      name={permission.user.fullName}
                      src={
                        permission.user.picture
                          ? getPictureUrlById(
                              permission.user.id,
                              permission.user.picture,
                              {
                                organizationId: workspace?.organization.id,
                              },
                            )
                          : undefined
                      }
                      size="sm"
                      className={cx(
                        'w-[40px]',
                        'h-[40px]',
                        'border',
                        'border-gray-300',
                        'dark:border-gray-700',
                      )}
                    />
                    <div className={cx('flex', 'flex-col', 'gap-0.5')}>
                      <Text noOfLines={1}>{permission.user.fullName}</Text>
                      <span className={cx('text-gray-500')}>
                        {permission.user.email}
                      </span>
                    </div>
                  </div>
                </Td>
                <Td>
                  <Badge>{permission.permission}</Badge>
                </Td>
                <Td className={cx('text-end')}>
                  <IconButton
                    icon={<IconDelete />}
                    colorScheme="red"
                    title="Revoke user permission"
                    aria-label="Revoke user permission"
                    isLoading={revokedPermission === permission.id}
                    isDisabled={me?.id === permission.user.id}
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

export default SharingUserPermissions
