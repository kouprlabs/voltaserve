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
import { Button } from '@chakra-ui/react'
import { IconAdd, IconCheck, SectionPlaceholder, Select } from '@koupr/ui'
import { OptionBase } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
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
import SharingGroupPermissions from './sharing-group-permissions'

interface PermissionTypeOption extends OptionBase {
  value: PermissionType
  label: string
}

const SharingGroupOverview = () => {
  const { id: workspaceId, fileId } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [isGranting, setIsGranting] = useState(false)
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
  const { mutate: mutatePermissions } = FileAPI.useGetGroupPermissions(
    file && geOwnerPermission(file.permission) ? file.id : undefined,
  )
  const isSingleSelection = selection.length === 1

  const handleGrantPermission = useCallback(async () => {
    if (group && permission) {
      try {
        setIsGranting(true)
        await FileAPI.grantGroupPermission({
          ids: selection,
          groupId: group.id,
          permission,
        })
        await mutateFiles?.()
        if (isSingleSelection) {
          await mutatePermissions()
        }
        setGroup(undefined)
        setIsGranting(false)
        if (!isSingleSelection) {
          dispatch(sharingModalDidClose())
        }
      } catch {
        setIsGranting(false)
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
            isLoading={isGranting}
            isDisabled={!group || !permission}
            onClick={() => handleGrantPermission()}
          >
            Apply to Group
          </Button>
        </div>
      ) : null}
      {groups && groups.totalElements === 0 ? (
        <SectionPlaceholder
          text="This organization has no groups."
          content={
            <>
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
            </>
          }
          height="auto"
        />
      ) : null}
      {isSingleSelection ? (
        <>
          <hr />
          <SharingGroupPermissions />
        </>
      ) : null}
    </div>
  )
}

export default SharingGroupOverview
