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
import {
  IconAdd,
  IconCheck,
  SectionError,
  SectionPlaceholder,
  Select,
} from '@koupr/ui'
import { OptionBase } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import GroupAPI, { Group } from '@/client/api/group'
import { geEditorPermission, PermissionType } from '@/client/api/permission'
import WorkspaceAPI from '@/client/api/workspace'
import { errorToString } from '@/client/error'
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

const SharingGroupForm = () => {
  const { id: workspaceId } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [isGranting, setIsGranting] = useState(false)
  const [group, setGroup] = useState<Group>()
  const [permission, setPermission] = useState<string>()
  const {
    data: workspace,
    error: workspaceError,
    isLoading: isWorkspaceLoading,
  } = WorkspaceAPI.useGet(workspaceId, swrConfig())
  const {
    data: groups,
    error: groupsError,
    isLoading: isGroupsLoading,
  } = GroupAPI.useList(
    {
      organizationId: workspace?.organization.id,
    },
    swrConfig(),
  )
  const isSingleSelection = selection.length === 1
  const { mutate: mutatePermissions } = FileAPI.useGetGroupPermissions(
    isSingleSelection ? selection[0] : undefined,
  )
  const isWorkspaceError = !workspace && workspaceError
  const isWorkspaceReady = workspace && !workspaceError
  const isGroupsError = !groups && groupsError
  const isGroupsEmpty = groups && !groupsError && groups.totalElements === 0
  const isGroupsReady = groups && !groupsError && groups.totalElements > 0

  const handleGrantPermission = useCallback(async () => {
    if (!group || !permission) {
      return
    }
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
  }, [
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
      {isWorkspaceLoading || isGroupsLoading ? <SharingFormSkeleton /> : null}
      {isWorkspaceError || isGroupsError ? (
        <SectionError
          text={errorToString(workspaceError || groupsError)}
          height="auto"
        />
      ) : null}
      {isWorkspaceReady && isGroupsEmpty ? (
        <SectionPlaceholder
          text="This organization has no groups."
          content={
            <>
              {workspace &&
              geEditorPermission(workspace.organization.permission) ? (
                <Button
                  as={Link}
                  to={`/new/group?org=${workspace.organization.id}`}
                  leftIcon={<IconAdd />}
                >
                  New Group
                </Button>
              ) : null}
            </>
          }
          height="auto"
        />
      ) : null}
      {isWorkspaceReady && isGroupsReady ? (
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
      {isSingleSelection ? (
        <>
          <hr />
          <SharingGroupPermissions />
        </>
      ) : null}
    </div>
  )
}

export default SharingGroupForm
