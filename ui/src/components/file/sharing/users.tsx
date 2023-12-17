import { useCallback, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  Text,
  Button,
  Stack,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  IconButton,
  Badge,
  Center,
  HStack,
  Avatar,
  VStack,
} from '@chakra-ui/react'
import { Spinner, variables } from '@koupr/ui'
import { IconCheck, IconTrash, IconUserPlus } from '@koupr/ui'
import { KeyedMutator } from 'swr'
import { Select } from 'chakra-react-select'
import FileAPI, { UserPermission } from '@/client/api/file'
import { geEditorPermission } from '@/client/api/permission'
import { User } from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import IdPUserAPI from '@/client/idp/user'
import UserSelector from '@/components/common/user-selector'
import { filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import reactSelectStyles from '@/styles/react-select'

type UsersProps = {
  users: User[]
  userPermissions: UserPermission[]
  mutateUserPermissions: KeyedMutator<UserPermission[]>
}

const Users = ({
  users,
  userPermissions,
  mutateUserPermissions,
}: UsersProps) => {
  const params = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGetById(params.id as string)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [activeUser, setActiveUser] = useState<User>()
  const [activeUserPermission, setActiveUserPermission] = useState<string>()
  const { data: user } = IdPUserAPI.useGet()
  const isSingleSelection = selection.length === 1

  const handleGrantUserPermission = useCallback(async () => {
    if (!activeUser || !activeUserPermission) {
      return
    }
    try {
      setIsGrantLoading(true)
      await FileAPI.grantUserPermission({
        ids: selection,
        userId: activeUser.id,
        permission: activeUserPermission,
      })
      const result = await FileAPI.batchGet({ ids: selection })
      dispatch(filesUpdated(result))
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
    selection,
    activeUser,
    activeUserPermission,
    isSingleSelection,
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
        const result = await FileAPI.batchGet({ ids: selection })
        dispatch(filesUpdated(result))
        if (isSingleSelection) {
          await mutateUserPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [selection, isSingleSelection, dispatch, mutateUserPermissions],
  )

  return (
    <Stack direction="column" spacing={variables.spacing}>
      {users && users.length === 0 ? (
        <Center>
          <VStack spacing={variables.spacing}>
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
          </VStack>
        </Center>
      ) : null}
      {users && users.length > 0 ? (
        <Stack direction="column" spacing={variables.spacing}>
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
                setActiveUserPermission(e.value)
              }
            }}
          />
          <Button
            leftIcon={<IconCheck />}
            colorScheme="blue"
            isLoading={isGrantLoading}
            isDisabled={!activeUser || !activeUserPermission}
            onClick={() => handleGrantUserPermission()}
          >
            Apply to User
          </Button>
        </Stack>
      ) : null}
      {isSingleSelection ? (
        <>
          <hr />
          {!userPermissions ? (
            <Center>
              <Spinner />
            </Center>
          ) : null}
          {userPermissions && userPermissions.length === 0 ? (
            <Center>
              <Text>Not shared with any users.</Text>
            </Center>
          ) : null}
          {userPermissions && userPermissions.length > 0 ? (
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
                  {userPermissions.map((p) => (
                    <Tr key={p.id}>
                      <Td p={variables.spacingSm}>
                        <HStack spacing={variables.spacingSm}>
                          <Avatar
                            name={p.user.fullName}
                            src={p.user.picture}
                            size="sm"
                            width="40px"
                            height="40px"
                          />
                          <Stack spacing={variables.spacingXs}>
                            <Text noOfLines={1}>{p.user.fullName}</Text>
                            <Text color="gray">{p.user.email}</Text>
                          </Stack>
                        </HStack>
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
    </Stack>
  )
}

export default Users
