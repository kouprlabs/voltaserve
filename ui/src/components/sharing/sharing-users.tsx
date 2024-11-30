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
  IconCheck,
  IconDelete,
  IconPersonAdd,
  Select,
  Spinner,
  Text,
} from '@koupr/ui'
import { OptionBase } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI, { UserPermission } from '@/client/api/file'
import {
  geEditorPermission,
  geOwnerPermission,
  PermissionType,
} from '@/client/api/permission'
import UserAPI, { User } from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import IdPUserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import UserSelector from '@/components/common/user-selector'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import { inviteModalDidOpen } from '@/store/ui/organizations'
import SharingFormSkeleton from './sharing-form-skeleton'

interface PermissionTypeOption extends OptionBase {
  value: PermissionType
  label: string
}

const SharingUsers = () => {
  const navigate = useNavigate()
  const { id: workspaceId, fileId } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [isGranting, setIsGranting] = useState(false)
  const [revokedPermission, setRevokedPermission] = useState<string>()
  const [user, setUser] = useState<User>()
  const [permission, setPermission] = useState<string>()
  const { data: workspace } = WorkspaceAPI.useGet(workspaceId)
  const { data: file } = FileAPI.useGet(selection[0], swrConfig())
  const { data: me } = IdPUserAPI.useGet()
  const { data: users } = UserAPI.useList(
    {
      organizationId: workspace?.organization.id,
    },
    swrConfig(),
  )
  const { data: permissions, mutate: mutatePermissions } =
    FileAPI.useGetUserPermissions(
      file && geOwnerPermission(file.permission) ? file.id : undefined,
      swrConfig(),
    )
  const isSingleSelection = selection.length === 1

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
    fileId,
    selection,
    user,
    permission,
    isSingleSelection,
    dispatch,
    mutateFiles,
    mutatePermissions,
  ])

  const handleRevokePermission = useCallback(
    async (permission: UserPermission) => {
      try {
        setRevokedPermission(permission.id)
        await FileAPI.revokeUserPermission({
          ids: selection,
          userId: permission.user.id,
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

  const handleInviteMembersClick = useCallback(async () => {
    if (workspace) {
      dispatch(inviteModalDidOpen())
      dispatch(sharingModalDidClose())
      navigate(`/organization/${workspace.organization.id}/member`)
    }
  }, [workspace, navigate, dispatch])

  return (
    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
      {!users ? <SharingFormSkeleton /> : null}
      {users && users.totalElements === 0 ? (
        <div className={cx('flex', 'items-center', 'justify-center')}>
          <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
            <span>This organization has no members.</span>
            {workspace &&
            geEditorPermission(workspace.organization.permission) ? (
              <Button
                leftIcon={<IconPersonAdd />}
                onClick={handleInviteMembersClick}
              >
                Invite Members
              </Button>
            ) : null}
          </div>
        </div>
      ) : null}
      {users && users.totalElements > 0 ? (
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
          {!permissions ? (
            <div className={cx('flex', 'items-center', 'justify-center')}>
              <Spinner />
            </div>
          ) : null}
          {permissions && permissions.length === 0 ? (
            <div className={cx('flex', 'items-center', 'justify-center')}>
              <span>Not shared with any users.</span>
            </div>
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
                          name={p.user.fullName}
                          src={
                            p.user.picture
                              ? getPictureUrlById(p.user.id, p.user.picture, {
                                  organizationId: workspace?.organization.id,
                                })
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
                          <Text noOfLines={1}>{p.user.fullName}</Text>
                          <span className={cx('text-gray-500')}>
                            {p.user.email}
                          </span>
                        </div>
                      </div>
                    </Td>
                    <Td>
                      <Badge>{p.permission}</Badge>
                    </Td>
                    <Td className={cx('text-end')}>
                      <IconButton
                        icon={<IconDelete />}
                        colorScheme="red"
                        title="Revoke user permission"
                        aria-label="Revoke user permission"
                        isLoading={revokedPermission === p.id}
                        isDisabled={me?.id === p.user.id}
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

export default SharingUsers
