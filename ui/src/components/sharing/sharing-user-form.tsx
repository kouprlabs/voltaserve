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
import { useNavigate, useParams } from 'react-router-dom'
import { Button } from '@chakra-ui/react'
import {
  IconCheck,
  IconPersonAdd,
  SectionError,
  SectionPlaceholder,
  Select,
} from '@koupr/ui'
import { OptionBase } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { geEditorPermission, PermissionType } from '@/client/api/permission'
import UserAPI, { User } from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import UserSelector from '@/components/common/user-selector'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import { inviteModalDidOpen } from '@/store/ui/organizations'
import SharingFormSkeleton from './sharing-form-skeleton'
import SharingUserPermissions from './sharing-user-permissions'

interface PermissionTypeOption extends OptionBase {
  value: PermissionType
  label: string
}

const SharingUserForm = () => {
  const { id: workspaceId } = useParams()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [isGranting, setIsGranting] = useState(false)
  const [user, setUser] = useState<User>()
  const [permission, setPermission] = useState<string>()
  const {
    data: workspace,
    error: workspaceError,
    isLoading: isWorkspaceLoading,
  } = WorkspaceAPI.useGet(workspaceId, swrConfig())
  const {
    data: users,
    error: usersError,
    isLoading: isUsersLoading,
  } = UserAPI.useList(
    {
      organizationId: workspace?.organization.id,
      excludeMe: true,
    },
    swrConfig(),
  )
  const isSingleSelection = selection.length === 1
  const { mutate: mutatePermissions } = FileAPI.useGetGroupPermissions(
    isSingleSelection ? selection[0] : undefined,
  )
  const isWorkspaceError = !workspace && workspaceError
  const isWorkspaceReady = workspace && !workspaceError
  const isUsersError = !users && usersError
  const isUsersEmpty = users && !usersError && users.totalElements === 0
  const isUsersReady = users && !usersError && users.totalElements > 0

  const handleGrantPermission = useCallback(async () => {
    if (!user || !permission) {
      return
    }
    try {
      setIsGranting(true)
      await FileAPI.grantUserPermission({
        ids: selection,
        userId: user.id,
        permission,
      })
      await mutateFiles?.()
      if (isSingleSelection) {
        await mutatePermissions()
      }
      setUser(undefined)
      setIsGranting(false)
      if (!isSingleSelection) {
        dispatch(sharingModalDidClose())
      }
    } catch {
      setIsGranting(false)
    }
  }, [
    selection,
    user,
    permission,
    isSingleSelection,
    dispatch,
    mutateFiles,
    mutatePermissions,
  ])

  const handleInviteMembersClick = useCallback(async () => {
    if (workspace) {
      dispatch(inviteModalDidOpen())
      dispatch(sharingModalDidClose())
      navigate(`/organization/${workspace.organization.id}/member`)
    }
  }, [workspace, navigate, dispatch])

  return (
    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
      {isWorkspaceLoading || isUsersLoading ? <SharingFormSkeleton /> : null}
      {isWorkspaceError || isUsersError ? (
        <SectionError
          text={errorToString(workspaceError || isUsersError)}
          height="auto"
        />
      ) : null}
      {isWorkspaceReady && isUsersEmpty ? (
        <SectionPlaceholder
          text="This organization has no members."
          content={
            <>
              {workspace &&
              geEditorPermission(workspace.organization.permission) ? (
                <Button
                  leftIcon={<IconPersonAdd />}
                  onClick={handleInviteMembersClick}
                >
                  Invite Members
                </Button>
              ) : null}
            </>
          }
          height="auto"
        />
      ) : null}
      {isWorkspaceReady && isUsersReady ? (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          <UserSelector
            value={user}
            organizationId={workspace?.organization.id}
            onConfirm={(value) => setUser(value)}
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
            isDisabled={!user || !permission}
            onClick={() => handleGrantPermission()}
          >
            Apply to User
          </Button>
        </div>
      ) : null}
      {isSingleSelection ? (
        <>
          <hr />
          <SharingUserPermissions />
        </>
      ) : null}
    </div>
  )
}

export default SharingUserForm
