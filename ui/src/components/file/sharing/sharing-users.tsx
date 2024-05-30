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
import { KeyedMutator } from 'swr'
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI, { UserPermission } from '@/client/api/file'
import { geEditorPermission } from '@/client/api/permission'
import { User } from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import IdPUserAPI from '@/client/idp/user'
import UserSelector from '@/components/common/user-selector'
import { IconCheck, IconDelete, IconPersonAdd } from '@/lib/components/icons'
import Spinner from '@/lib/components/spinner'
import Text from '@/lib/components/text'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import { inviteModalDidOpen } from '@/store/ui/organizations'
import { reactSelectStyles } from '@/styles/react-select'
import SharingFormSkeleton from './sharing-form-skeleton'

export type SharingUsersProps = {
  users?: User[]
  permissions?: UserPermission[]
  mutateUserPermissions: KeyedMutator<UserPermission[]>
}

const SharingUsers = ({
  users,
  permissions,
  mutateUserPermissions,
}: SharingUsersProps) => {
  const navigate = useNavigate()
  const { id: workspaceId, fileId } = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGet(workspaceId)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [activeUser, setActiveUser] = useState<User>()
  const [activePermission, setActivePermission] = useState<string>()
  const { data: user } = IdPUserAPI.useGet()
  const isSingleSelection = selection.length === 1

  const handleGrantUserPermission = useCallback(async () => {
    if (!activeUser || !activePermission) {
      return
    }
    try {
      setIsGrantLoading(true)
      await FileAPI.grantUserPermission({
        ids: selection,
        userId: activeUser.id,
        permission: activePermission,
      })
      mutateList?.()
      if (isSingleSelection) {
        await mutateUserPermissions()
      }
      setActiveUser(undefined)
      setIsGrantLoading(false)
      if (!isSingleSelection) {
        dispatch(sharingModalDidClose())
      }
    } catch {
      setIsGrantLoading(false)
    }
  }, [
    fileId,
    selection,
    activeUser,
    activePermission,
    isSingleSelection,
    mutateList,
    dispatch,
    mutateUserPermissions,
  ])

  const handleRevokeUserPermission = useCallback(
    async (permission: UserPermission) => {
      try {
        setPermissionBeingRevoked(permission.id)
        await FileAPI.revokeUserPermission({
          ids: selection,
          userId: permission.user.id,
        })
        mutateList?.()
        if (isSingleSelection) {
          await mutateUserPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [fileId, selection, isSingleSelection, mutateList, mutateUserPermissions],
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
      {users && users.length === 0 ? (
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
      {users && users.length > 0 ? (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          <UserSelector
            value={activeUser}
            organizationId={workspace?.organization.id}
            onConfirm={(value) => setActiveUser(value)}
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
            isDisabled={!activeUser || !activePermission}
            onClick={() => handleGrantUserPermission()}
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
            <>
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
                            src={p.user.picture}
                            size="sm"
                            className={cx('w-[40px]', 'h-[40px]')}
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
                          aria-label=""
                          isLoading={permissionBeingRevoked === p.id}
                          isDisabled={user?.id === p.user.id}
                          onClick={() => handleRevokeUserPermission(p)}
                        />
                      </Td>
                    </Tr>
                  ))}
                </Tbody>
              </Table>
            </>
          ) : null}
        </>
      ) : null}
    </div>
  )
}

export default SharingUsers
