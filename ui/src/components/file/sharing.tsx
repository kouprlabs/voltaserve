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
  HStack,
  Tag,
  Avatar,
  VStack,
  SystemStyleObject,
} from '@chakra-ui/react'
import { Spinner, variables } from '@koupr/ui'
import { IconAdd, IconCheck, IconTrash, IconUserPlus } from '@koupr/ui'
import { Select } from 'chakra-react-select'
import FileAPI, { GroupPermission, UserPermission } from '@/client/api/file'
import GroupAPI, { Group } from '@/client/api/group'
import { geEditorPermission } from '@/client/api/permission'
import UserAPI, { User } from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import IdPUserAPI from '@/client/idp/user'
import GroupSelector from '@/components/common/group-selector'
import UserSelector from '@/components/common/user-selector'
import { filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { selectionUpdated, sharingModalDidClose } from '@/store/ui/files'

const Sharing = () => {
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
  const { data: user } = IdPUserAPI.useGet()
  const isSingleFileMode = useMemo(() => selection.length === 1, [selection])
  const selectStyles = useMemo(() => {
    return {
      dropdownIndicator: (provided: SystemStyleObject) => ({
        ...provided,
        bg: 'transparent',
        cursor: 'inherit',
        position: 'absolute',
        right: '0px',
      }),
      indicatorSeparator: (provided: SystemStyleObject) => ({
        ...provided,
        display: 'none',
      }),
      placeholder: (provided: SystemStyleObject) => ({
        ...provided,
        textAlign: 'center',
      }),
      singleValue: (provided: SystemStyleObject) => ({
        ...provided,
        textAlign: 'center',
      }),
    }
  }, [])

  const loadUsers = useCallback(async () => {
    if (workspace) {
      const { data } = await UserAPI.list({
        organizationId: workspace.organization.id,
      })
      setUsers(data)
    }
  }, [workspace])

  const loadGroups = useCallback(async () => {
    if (workspace) {
      const { data } = await GroupAPI.list({
        organizationId: workspace.organization.id,
      })
      setGroups(data)
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
    [selection, isSingleFileMode, dispatch, loadUserPermissions],
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
    [selection, isSingleFileMode, dispatch, loadGroupPermissions],
  )

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => {
        setUsers(undefined)
        setGroups(undefined)
        setActiveUserId(undefined)
        setActiveGroupId(undefined)
        setUserPermissions(undefined)
        setGroupPermissions(undefined)
        dispatch(selectionUpdated([]))
        dispatch(sharingModalDidClose())
      }}
      closeOnOverlayClick={false}
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
                          workspace.organization.permission,
                        ) ? (
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
                        organizationId={workspace?.organization.id}
                        onConfirm={(value) => setActiveUserId(value.id)}
                      />
                      <Select
                        options={[
                          { value: 'viewer', label: 'Viewer' },
                          { value: 'editor', label: 'Editor' },
                          { value: 'owner', label: 'Owner' },
                        ]}
                        placeholder="Select Permission"
                        selectedOptionStyle="check"
                        chakraStyles={selectStyles}
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
                        isDisabled={!activeUserId || !activeUserPermission}
                        onClick={() => handleGrantUserPermission()}
                      >
                        Apply to User
                      </Button>
                    </Stack>
                  ) : null}
                  {isSingleFileMode && <hr />}
                  {!userPermissions && isSingleFileMode ? (
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
                </Stack>
              </TabPanel>
              <TabPanel>
                <Stack direction="column" spacing={variables.spacing}>
                  {!groups ? <FormSkeleton /> : null}
                  {groups && groups.length > 0 ? (
                    <Stack direction="column" spacing={variables.spacing}>
                      <GroupSelector
                        organizationId={workspace?.organization.id}
                        onConfirm={(value) => setActiveGroupId(value.id)}
                      />
                      <Select
                        options={[
                          { value: 'viewer', label: 'Viewer' },
                          { value: 'editor', label: 'Editor' },
                          { value: 'owner', label: 'Owner' },
                        ]}
                        placeholder="Select Permission"
                        selectedOptionStyle="check"
                        chakraStyles={selectStyles}
                      />
                      <Button
                        leftIcon={<IconCheck />}
                        colorScheme="blue"
                        isLoading={isGrantLoading}
                        isDisabled={!activeGroupId || !activeGroupPermission}
                        onClick={() => handleGrantGroupPermission()}
                      >
                        Apply to Group
                      </Button>
                    </Stack>
                  ) : null}
                  {groups && groups.length === 0 ? (
                    <Center>
                      <VStack spacing={variables.spacing}>
                        <Text>This organization has no groups.</Text>
                        {workspace &&
                        geEditorPermission(
                          workspace.organization.permission,
                        ) ? (
                          <Button
                            as={Link}
                            leftIcon={<IconAdd />}
                            to={`/new/group?org=${workspace.organization.id}`}
                          >
                            New Group
                          </Button>
                        ) : null}
                      </VStack>
                    </Center>
                  ) : null}
                  {isSingleFileMode ? <hr /> : null}
                  {!groupPermissions && isSingleFileMode ? (
                    <Center>
                      <Spinner />
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

export default Sharing
