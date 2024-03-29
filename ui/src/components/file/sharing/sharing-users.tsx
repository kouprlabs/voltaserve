import { useCallback, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  Text,
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
import { Spinner, variables } from '@koupr/ui'
import { IconCheck, IconTrash, IconUserPlus } from '@koupr/ui'
import { KeyedMutator, useSWRConfig } from 'swr'
import { Select } from 'chakra-react-select'
import classNames from 'classnames'
import FileAPI, { List, UserPermission } from '@/client/api/file'
import { geEditorPermission } from '@/client/api/permission'
import { User } from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import IdPUserAPI from '@/client/idp/user'
import UserSelector from '@/components/common/user-selector'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import reactSelectStyles from '@/styles/react-select'
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
  const { mutate } = useSWRConfig()
  const { id, fileId } = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGetById(id)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [activeUser, setActiveUser] = useState<User>()
  const [activePermission, setActivePermission] = useState<string>()
  const { data: user } = IdPUserAPI.useGet()
  const fileListSearchParams = useFileListSearchParams()
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
      await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
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
    fileListSearchParams,
    mutate,
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
        await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
        if (isSingleSelection) {
          await mutateUserPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [
      fileId,
      selection,
      isSingleSelection,
      fileListSearchParams,
      mutate,
      mutateUserPermissions,
    ],
  )

  return (
    <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
      {!users ? <SharingFormSkeleton /> : null}
      {users && users.length === 0 ? (
        <div className={classNames('flex', 'items-center', 'justify-center')}>
          <div
            className={classNames(
              'flex',
              'flex-col',
              'items-center',
              'gap-1.5',
            )}
          >
            <Text>This organization has no members.</Text>
            {workspace &&
            geEditorPermission(workspace.organization.permission) ? (
              <Button
                as={Link}
                leftIcon={<IconUserPlus />}
                to={`/organization/${workspace.organization.id}/member?invite=true`}
              >
                Invite Members
              </Button>
            ) : null}
          </div>
        </div>
      ) : null}
      {users && users.length > 0 ? (
        <div className={classNames('flex', 'flex-col', 'gap-1.5')}>
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
            chakraStyles={reactSelectStyles}
            onChange={(e) => {
              if (e) {
                setActivePermission(e.value)
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
            <div
              className={classNames('flex', 'items-center', 'justify-center')}
            >
              <Spinner />
            </div>
          ) : null}
          {permissions && permissions.length === 0 ? (
            <div
              className={classNames('flex', 'items-center', 'justify-center')}
            >
              <Text>Not shared with any users.</Text>
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
                      <Td p={variables.spacingSm}>
                        <div
                          className={classNames(
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
                            width="40px"
                            height="40px"
                          />
                          <div
                            className={classNames(
                              'flex',
                              'flex-col',
                              'gap-0.5',
                            )}
                          >
                            <Text noOfLines={1}>{p.user.fullName}</Text>
                            <Text color="gray">{p.user.email}</Text>
                          </div>
                        </div>
                      </Td>
                      <Td>
                        <Badge>{p.permission}</Badge>
                      </Td>
                      <Td textAlign="end">
                        <IconButton
                          icon={<IconTrash />}
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
