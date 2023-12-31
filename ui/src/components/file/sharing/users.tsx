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
import FormSkeleton from './form-skeleton'

type UsersProps = {
  users?: User[]
  permissions?: UserPermission[]
  mutateUserPermissions: KeyedMutator<UserPermission[]>
}

const Users = ({ users, permissions, mutateUserPermissions }: UsersProps) => {
  const params = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGetById(params.id as string)
  const selectedItems = useAppSelector((state) => state.ui.files.selectedItems)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [activeUser, setActiveUser] = useState<User>()
  const [activePermission, setActivePermission] = useState<string>()
  const { data: user } = IdPUserAPI.useGet()
  const isSingleSelection = selectedItems.length === 1

  const handleGrantUserPermission = useCallback(async () => {
    if (!activeUser || !activePermission) {
      return
    }
    try {
      setIsGrantLoading(true)
      await FileAPI.grantUserPermission({
        ids: selectedItems,
        userId: activeUser.id,
        permission: activePermission,
      })
      const result = await FileAPI.batchGet({ ids: selectedItems })
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
    selectedItems,
    activeUser,
    activePermission,
    isSingleSelection,
    dispatch,
    mutateUserPermissions,
  ])

  const handleRevokeUserPermission = useCallback(
    async (permission: UserPermission) => {
      try {
        setPermissionBeingRevoked(permission.id)
        await FileAPI.revokeUserPermission({
          ids: selectedItems,
          userId: permission.user.id,
        })
        const result = await FileAPI.batchGet({ ids: selectedItems })
        dispatch(filesUpdated(result))
        if (isSingleSelection) {
          await mutateUserPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [selectedItems, isSingleSelection, dispatch, mutateUserPermissions],
  )

  return (
    <Stack direction="column" spacing={variables.spacing}>
      {!users ? <FormSkeleton /> : null}
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
        </Stack>
      ) : null}
      {isSingleSelection ? (
        <>
          <hr />
          {!permissions ? (
            <Center>
              <Spinner />
            </Center>
          ) : null}
          {permissions && permissions.length === 0 ? (
            <Center>
              <Text>Not shared with any users.</Text>
            </Center>
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
