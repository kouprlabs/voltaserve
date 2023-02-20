import { useCallback, useEffect, useMemo, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  Text,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  Button,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  Stack,
  Select,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  IconButton,
  Badge,
  Center,
  Skeleton,
  Spinner,
  HStack,
  Tag,
  Avatar,
  VStack,
} from '@chakra-ui/react'
import FileAPI, { GroupPermission, UserPermission } from '@/api/file'
import { Group } from '@/api/group'
import OrganizationAPI from '@/api/organization'
import { geEditorPermission } from '@/api/permission'
import UserAPI, { User } from '@/api/user'
import WorkspaceAPI from '@/api/workspace'
import { filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { selectionUpdated, sharingModalDidClose } from '@/store/ui/files'
import {
  IconAdd,
  IconCheck,
  IconTrash,
  IconUserPlus,
} from '@/components/common/icon'
import variables from '@/theme/variables'
import userToString from '@/helpers/user-to-string'

const FileSharing = () => {
  const params = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGetById(params.id as string)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isShareModalOpen)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [users, setUsers] = useState<User[]>()
  const [groups, setGroups] = useState<Group[]>()
  const [userPermissions, setUserPermissions] = useState<UserPermission[]>()
  const [groupPermissions, setGroupPermissions] = useState<GroupPermission[]>()
  const [activeUserId, setActiveUserId] = useState<string>()
  const [activeUserPermission, setActiveUserPermission] = useState<string>()
  const [activeGroupId, setActiveGroupId] = useState<string>()
  const [activeGroupPermission, setActiveGroupPermission] = useState<string>()
  const { data: user } = UserAPI.useGet()
  const isSingleFileMode = useMemo(() => selection.length === 1, [selection])

  const loadUsers = useCallback(async () => {
    if (workspace) {
      setUsers(await OrganizationAPI.getMembers(workspace.organization.id))
    }
  }, [workspace])

  const loadGroups = useCallback(async () => {
    if (workspace) {
      setGroups(await OrganizationAPI.getGroups(workspace.organization.id))
    }
  }, [workspace])

  const loadUserPermissions = useCallback(async () => {
    if (isSingleFileMode) {
      setUserPermissions(await FileAPI.getUserPermissions(selection[0]))
    }
  }, [selection, isSingleFileMode])

  const loadGroupPermissions = useCallback(async () => {
    if (isSingleFileMode) {
      setGroupPermissions(await FileAPI.getGroupPermissions(selection[0]))
    }
  }, [selection, isSingleFileMode])

  useEffect(() => {
    if (selection.length === 0) {
      dispatch(sharingModalDidClose())
    }
  }, [selection, dispatch])

  useEffect(() => {
    if (isModalOpen) {
      ;(async () => {
        await loadUsers()
        await loadGroups()
        await loadUserPermissions()
        await loadGroupPermissions()
      })()
    }
  }, [
    isModalOpen,
    loadUsers,
    loadGroups,
    loadUserPermissions,
    loadGroupPermissions,
  ])

  const handleGrantUserPermission = useCallback(async () => {
    if (!activeUserId || !activeUserPermission) {
      return
    }
    try {
      setIsGrantLoading(true)
      await FileAPI.grantUserPermission({
        ids: selection,
        userId: activeUserId,
        permission: activeUserPermission,
      })
      const result = await FileAPI.batchGet({ ids: selection })
      dispatch(filesUpdated(result))
      if (isSingleFileMode) {
        await loadUserPermissions()
      }
      setActiveUserId('')
      setActiveUserPermission('')
      setIsGrantLoading(false)
      if (!isSingleFileMode) {
        dispatch(sharingModalDidClose())
      }
    } catch {
      setIsGrantLoading(false)
    }
  }, [
    selection,
    activeUserId,
    activeUserPermission,
    isSingleFileMode,
    dispatch,
    loadUserPermissions,
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
        if (isSingleFileMode) {
          await loadUserPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [selection, isSingleFileMode, dispatch, loadUserPermissions]
  )

  const handleGrantGroupPermission = useCallback(async () => {
    if (activeGroupId && activeGroupPermission) {
      try {
        setIsGrantLoading(true)
        await FileAPI.grantGroupPermission({
          ids: selection,
          groupId: activeGroupId,
          permission: activeGroupPermission,
        })
        const result = await FileAPI.batchGet({ ids: selection })
        dispatch(filesUpdated(result))
        if (isSingleFileMode) {
          await loadGroupPermissions()
        }
        setActiveGroupId('')
        setActiveGroupPermission('')
        setIsGrantLoading(false)
        if (!isSingleFileMode) {
          dispatch(sharingModalDidClose())
        }
      } catch {
        setIsGrantLoading(false)
      }
    }
  }, [
    selection,
    activeGroupId,
    activeGroupPermission,
    isSingleFileMode,
    dispatch,
    loadGroupPermissions,
  ])

  const handleRevokeGroupPermission = useCallback(
    async (permission: GroupPermission) => {
      try {
        setPermissionBeingRevoked(permission.id)
        await FileAPI.revokeGroupPermission({
          ids: selection,
          groupId: permission.group.id,
        })
        const result = await FileAPI.batchGet({ ids: selection })
        dispatch(filesUpdated(result))
        if (isSingleFileMode) {
          await loadGroupPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [selection, isSingleFileMode, dispatch, loadGroupPermissions]
  )

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => {
        setUsers(undefined)
        setGroups(undefined)
        setUserPermissions(undefined)
        setGroupPermissions(undefined)
        dispatch(selectionUpdated([]))
        dispatch(sharingModalDidClose())
      }}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Sharing</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Tabs>
            <TabList h="40px">
              <Tab>
                <HStack>
                  <Text>People</Text>
                  {userPermissions && userPermissions.length > 0 ? (
                    <Tag borderRadius="full">{userPermissions.length}</Tag>
                  ) : null}
                </HStack>
              </Tab>
              <Tab>
                <HStack>
                  <Text>Groups</Text>
                  {groupPermissions && groupPermissions.length > 0 ? (
                    <Tag borderRadius="full">{groupPermissions.length}</Tag>
                  ) : null}
                </HStack>
              </Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <Stack direction="column" spacing={variables.spacing}>
                  {!users ? <FormSkeleton /> : null}
                  {users && users.length === 0 ? (
                    <Center>
                      <VStack spacing={variables.spacing}>
                        <Text>This organization has no members.</Text>
                        {workspace &&
                        geEditorPermission(
                          workspace.organization.permission
                        ) ? (
                          <Button
                            as={Link}
                            leftIcon={<IconUserPlus />}
                            to={`/organization/${workspace.organization.id}/member?invite=true`}
                          >
                            Invite members
                          </Button>
                        ) : null}
                      </VStack>
                    </Center>
                  ) : null}
                  {users && users.length > 0 ? (
                    <Stack direction="column" spacing={variables.spacing}>
                      <Select
                        placeholder="Select user"
                        value={activeUserId}
                        isDisabled={isGrantLoading}
                        onChange={(e) => setActiveUserId(e.target.value)}
                      >
                        {users.map((u) => (
                          <option key={u.id} value={u.id}>
                            {userToString(u)}
                          </option>
                        ))}
                      </Select>
                      <Select
                        placeholder="Select permission"
                        value={activeUserPermission}
                        isDisabled={isGrantLoading}
                        onChange={(e) =>
                          setActiveUserPermission(e.target.value)
                        }
                      >
                        <option value="viewer">Viewer</option>
                        <option value="editor">Editor</option>
                        <option value="owner">Owner</option>
                      </Select>
                      <Button
                        leftIcon={<IconCheck />}
                        colorScheme="blue"
                        isLoading={isGrantLoading}
                        isDisabled={!activeUserId || !activeUserPermission}
                        onClick={() => handleGrantUserPermission()}
                      >
                        Apply to user
                      </Button>
                    </Stack>
                  ) : null}
                  {isSingleFileMode && <hr />}
                  {!userPermissions && isSingleFileMode ? (
                    <Center>
                      <Spinner size="sm" thickness="4px" />
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
                </Stack>
              </TabPanel>
              <TabPanel>
                <Stack direction="column" spacing={variables.spacing}>
                  {!groups ? <FormSkeleton /> : null}
                  {groups && groups.length > 0 ? (
                    <Stack direction="column" spacing={variables.spacing}>
                      <Select
                        placeholder="Select group"
                        value={activeGroupId}
                        isDisabled={isGrantLoading}
                        onChange={(e) => setActiveGroupId(e.target.value)}
                      >
                        {groups.map((g) => (
                          <option key={g.id} value={g.id}>
                            {g.name}
                          </option>
                        ))}
                      </Select>
                      <Select
                        placeholder="Select permission"
                        value={activeGroupPermission}
                        isDisabled={isGrantLoading}
                        onChange={(e) =>
                          setActiveGroupPermission(e.target.value)
                        }
                      >
                        <option value="viewer">Viewer</option>
                        <option value="editor">Editor</option>
                        <option value="owner">Owner</option>
                      </Select>
                      <Button
                        leftIcon={<IconCheck />}
                        colorScheme="blue"
                        isLoading={isGrantLoading}
                        isDisabled={!activeGroupId || !activeGroupPermission}
                        onClick={() => handleGrantGroupPermission()}
                      >
                        Apply to group
                      </Button>
                    </Stack>
                  ) : null}
                  {groups && groups.length === 0 ? (
                    <Center>
                      <VStack spacing={variables.spacing}>
                        <Text>This organization has no groups.</Text>
                        {workspace &&
                        geEditorPermission(
                          workspace.organization.permission
                        ) ? (
                          <Button
                            as={Link}
                            leftIcon={<IconAdd />}
                            to={`/new/group?org=${workspace.organization.id}`}
                          >
                            New group
                          </Button>
                        ) : null}
                      </VStack>
                    </Center>
                  ) : null}
                  {isSingleFileMode ? <hr /> : null}
                  {!groupPermissions && isSingleFileMode ? (
                    <Center>
                      <Spinner size="sm" thickness="4px" />
                    </Center>
                  ) : null}
                  {groupPermissions && groupPermissions.length === 0 ? (
                    <Center>
                      <Text>Not shared with any groups.</Text>
                    </Center>
                  ) : null}
                  {groupPermissions && groupPermissions.length > 0 ? (
                    <Table>
                      <Thead>
                        <Tr>
                          <Th>Group</Th>
                          <Th>Permission</Th>
                          <Th />
                        </Tr>
                      </Thead>
                      <Tbody>
                        {groupPermissions.map((p) => (
                          <Tr key={p.id}>
                            <Td p={variables.spacingSm}>
                              <HStack spacing={variables.spacingSm}>
                                <Avatar
                                  name={p.group.name}
                                  size="sm"
                                  width="40px"
                                  height="40px"
                                />
                                <Text noOfLines={1}>{p.group.name}</Text>
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
                                onClick={() => handleRevokeGroupPermission(p)}
                              />
                            </Td>
                          </Tr>
                        ))}
                      </Tbody>
                    </Table>
                  ) : null}
                </Stack>
              </TabPanel>
            </TabPanels>
          </Tabs>
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}

const FormSkeleton = () => (
  <Stack spacing={variables.spacing}>
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
    <Skeleton height="40px" borderRadius={variables.borderRadiusMd} />
  </Stack>
)

export default FileSharing
