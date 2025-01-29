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
  IconCheck,
  IconPersonAdd,
  SectionError,
  SectionPlaceholder,
  Select,
} from '@koupr/ui'
import { OptionBase } from 'chakra-react-select'
import cx from 'classnames'
import { FileAPI } from '@/client/api/file'
import { geEditorPermission, PermissionType } from '@/client/api/permission'
import { UserAPI, User } from '@/client/api/user'
import { WorkspaceAPI } from '@/client/api/workspace'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import UserSelector from '@/components/common/user-selector'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import SharingFormSkeleton from './sharing-form-skeleton'
import SharingUserPermissions from './sharing-user-permissions'

interface PermissionTypeOption extends OptionBase {
  value: PermissionType
  label: string
}

const SharingUserForm = () => {
  const { id: workspaceId } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [isGranting, setIsGranting] = useState(false)
  const [user, setUser] = useState<User>()
  const [permission, setPermission] = useState<string>()
  const {
    data: workspace,
    error: workspaceError,
    isLoading: workspaceIsLoading,
  } = WorkspaceAPI.useGet(workspaceId, swrConfig())
  const {
    data: userList,
    error: userListError,
    isLoading: userListIsLoading,
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
  const workspaceIsReady = workspace && !workspaceError
  // prettier-ignore
  const userListIsEmpty = userList && !userListError && userList.totalElements === 0
  // prettier-ignore
  const userListIsReady = userList && !userListError && userList.totalElements > 0

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

  return (
    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
      {workspaceIsLoading || userListIsLoading ? <SharingFormSkeleton /> : null}
      {workspaceError || userListError ? (
        <SectionError
          text={errorToString(workspaceError || userListError)}
          height="auto"
        />
      ) : null}
      {workspaceIsReady && userListIsEmpty ? (
        <SectionPlaceholder
          text="This organization has no members."
          content={
            <>
              {workspace &&
              geEditorPermission(workspace.organization.permission) ? (
                <Button
                  as={Link}
                  to={`/organization/${workspace.organization.id}/member?invite=true`}
                  leftIcon={<IconPersonAdd />}
                >
                  Invite Members
                </Button>
              ) : null}
            </>
          }
          height="auto"
        />
      ) : null}
      {workspaceIsReady && userListIsReady ? (
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
