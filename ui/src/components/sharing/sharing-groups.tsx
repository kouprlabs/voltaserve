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
import { Link, useParams } from 'react-router-dom'
import {
  Button,
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
  IconAdd,
  IconCheck,
  IconDelete,
  Spinner,
  Text,
  Select,
} from '@koupr/ui'
import { OptionBase } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI, { GroupPermission } from '@/client/api/file'
import GroupAPI, { Group } from '@/client/api/group'
import {
  geEditorPermission,
  geOwnerPermission,
  PermissionType,
} from '@/client/api/permission'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import GroupSelector from '@/components/common/group-selector'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import SharingFormSkeleton from './sharing-form-skeleton'

interface PermissionTypeOption extends OptionBase {
  value: PermissionType
  label: string
}

const SharingGroups = () => {
  const { id: workspaceId, fileId } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [revokedPermission, setRevokedPermission] = useState<string>()
  const [group, setGroup] = useState<Group>()
  const [permission, setPermission] = useState<string>()
  const { data: workspace } = WorkspaceAPI.useGet(workspaceId)
  const { data: file } = FileAPI.useGet(selection[0], swrConfig())
  const { data: groups } = GroupAPI.useList(
    {
      organizationId: workspace?.organization.id,
    },
    swrConfig(),
  )
  const { data: permissions, mutate: mutatePermissions } =
    FileAPI.useGetGroupPermissions(
      file && geOwnerPermission(file.permission) ? file.id : undefined,
      swrConfig(),
    )
  const isSingleSelection = selection.length === 1

  const handleGrantPermission = useCallback(async () => {
    if (group && permission) {
      try {
        setIsGrantLoading(true)
        await FileAPI.grantGroupPermission({
          ids: selection,
          groupId: group.id,
          permission: permission,
        })
        await mutateFiles?.()
        if (isSingleSelection) {
          await mutatePermissions()
        }
        setGroup(undefined)
        setIsGrantLoading(false)
        if (!isSingleSelection) {
          dispatch(sharingModalDidClose())
        }
      } catch {
        setIsGrantLoading(false)
      }
    }
  }, [
    fileId,
    selection,
    group,
    permission,
    isSingleSelection,
    dispatch,
    mutateFiles,
    mutatePermissions,
  ])

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
    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
      {!groups ? <SharingFormSkeleton /> : null}
      {groups && groups.totalElements > 0 ? (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          <GroupSelector
            value={group}
            organizationId={workspace?.organization.id}
            onConfirm={(value) => setGroup(value)}
          />
          <Select<PermissionTypeOption, false>
            options={[
              { value: 'viewer', label: 'Viewer' },
              { value: 'editor', label: 'Editor' },
              { value: 'owner', label: 'Owner' },
            ]}
            placeholder="Select Permission"
            selectedOptionStyle="check"
            onChange={(newValue) => {
              if (newValue) {
                setPermission(newValue.value)
              }
            }}
          />
          <Button
            leftIcon={<IconCheck />}
            colorScheme="blue"
            isLoading={isGrantLoading}
            isDisabled={!group || !permission}
            onClick={() => handleGrantPermission()}
          >
            Apply to Group
          </Button>
        </div>
      ) : null}
      {groups && groups.totalElements === 0 ? (
        <div className={cx('flex', 'items-center', 'justify-center')}>
          <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
            <span>This organization has no groups.</span>
            {workspace &&
            geEditorPermission(workspace.organization.permission) ? (
              <Button
                as={Link}
                leftIcon={<IconAdd />}
                to={`/new/group?org=${workspace.organization.id}`}
              >
                New Group
              </Button>
            ) : null}
          </div>
        </div>
      ) : null}
      {isSingleSelection ? (
        <>
          <hr />
          {!permissions ? (
            <div className={cx('flex', 'items-center', 'justify-center')}>
              <Spinner />
            </div>
          ) : null}
          {permissions && permissions.length === 0 ? (
            <div className={cx('flex', 'items-center', 'justify-center')}>
              <span>Not shared with any groups.</span>
            </div>
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
                {permissions.map((p) => (
                  <Tr key={p.id}>
                    <Td className={cx('p-1')}>
                      <div
                        className={cx(
                          'flex',
                          'flex-row',
                          'items-center',
                          'gap-1',
                        )}
                      >
                        <Avatar
                          name={p.group.name}
                          size="sm"
                          className={cx('w-[40px]', 'h-[40px]')}
                        />
                        <Text noOfLines={1}>{p.group.name}</Text>
                      </div>
                    </Td>
                    <Td>
                      <Badge>{p.permission}</Badge>
                    </Td>
                    <Td className={cx('text-end')}>
                      <IconButton
                        icon={<IconDelete />}
                        colorScheme="red"
                        title="Revoke group permission"
                        aria-label="Revoke group permission"
                        isLoading={revokedPermission === p.id}
                        onClick={() => handleRevokePermission(p)}
                      />
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          ) : null}
        </>
      ) : null}
    </div>
  )
}

export default SharingGroups
