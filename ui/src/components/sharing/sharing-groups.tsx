// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

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
import { KeyedMutator } from 'swr'
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI, { GroupPermission } from '@/client/api/file'
import { Group } from '@/client/api/group'
import { geEditorPermission } from '@/client/api/permission'
import WorkspaceAPI from '@/client/api/workspace'
import GroupSelector from '@/components/common/group-selector'
import { IconAdd, IconCheck, IconDelete } from '@/lib/components/icons'
import Spinner from '@/lib/components/spinner'
import Text from '@/lib/components/text'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import { reactSelectStyles } from '@/styles/react-select'
import SharingFormSkeleton from './sharing-form-skeleton'

export type SharingGroupsProps = {
  groups?: Group[]
  permissions?: GroupPermission[]
  mutateGroupPermissions: KeyedMutator<GroupPermission[]>
}

const SharingGroups = ({
  groups,
  permissions,
  mutateGroupPermissions,
}: SharingGroupsProps) => {
  const { id: workspaceId, fileId } = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGet(workspaceId)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [activeGroup, setActiveGroup] = useState<Group>()
  const [activePermission, setActivePermission] = useState<string>()
  const isSingleSelection = selection.length === 1

  const handleGrantGroupPermission = useCallback(async () => {
    if (activeGroup && activePermission) {
      try {
        setIsGrantLoading(true)
        await FileAPI.grantGroupPermission({
          ids: selection,
          groupId: activeGroup.id,
          permission: activePermission,
        })
        mutateList?.()
        if (isSingleSelection) {
          await mutateGroupPermissions()
        }
        setActiveGroup(undefined)
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
    activeGroup,
    activePermission,
    isSingleSelection,
    mutateList,
    dispatch,
    mutateGroupPermissions,
  ])

  const handleRevokeGroupPermission = useCallback(
    async (permission: GroupPermission) => {
      try {
        setPermissionBeingRevoked(permission.id)
        await FileAPI.revokeGroupPermission({
          ids: selection,
          groupId: permission.group.id,
        })
        mutateList?.()
        if (isSingleSelection) {
          await mutateGroupPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [fileId, selection, isSingleSelection, mutateList, mutateGroupPermissions],
  )

  return (
    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
      {!groups ? <SharingFormSkeleton /> : null}
      {groups && groups.length > 0 ? (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          <GroupSelector
            value={activeGroup}
            organizationId={workspace?.organization.id}
            onConfirm={(value) => setActiveGroup(value)}
          />
          <Select
            options={[
              { value: 'viewer', label: 'Viewer' },
              { value: 'editor', label: 'Editor' },
              { value: 'owner', label: 'Owner' },
            ]}
            placeholder="Select Permission"
            selectedOptionStyle="check"
            chakraStyles={reactSelectStyles()}
            onChange={(event) => {
              if (event) {
                setActivePermission(event.value)
              }
            }}
          />
          <Button
            leftIcon={<IconCheck />}
            colorScheme="blue"
            isLoading={isGrantLoading}
            isDisabled={!activeGroup || !activePermission}
            onClick={() => handleGrantGroupPermission()}
          >
            Apply to Group
          </Button>
        </div>
      ) : null}
      {groups && groups.length === 0 ? (
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
                        aria-label=""
                        isLoading={permissionBeingRevoked === p.id}
                        onClick={() => handleRevokeGroupPermission(p)}
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
