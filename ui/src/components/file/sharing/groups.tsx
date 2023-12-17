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
import { IconAdd, IconCheck, IconTrash } from '@koupr/ui'
import { KeyedMutator } from 'swr'
import { Select } from 'chakra-react-select'
import FileAPI, { GroupPermission } from '@/client/api/file'
import { Group } from '@/client/api/group'
import { geEditorPermission } from '@/client/api/permission'
import WorkspaceAPI from '@/client/api/workspace'
import GroupSelector from '@/components/common/group-selector'
import { filesUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import reactSelectStyles from '@/styles/react-select'

type GroupsProps = {
  groups: Group[]
  groupPermissions: GroupPermission[]
  mutateGroupPermissions: KeyedMutator<GroupPermission[]>
}

const Groups = ({
  groups,
  groupPermissions,
  mutateGroupPermissions,
}: GroupsProps) => {
  const params = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGetById(params.id as string)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [activeGroup, setActiveGroup] = useState<Group>()
  const [activeGroupPermission, setActiveGroupPermission] = useState<string>()
  const isSingleSelection = selection.length === 1

  const handleGrantGroupPermission = useCallback(async () => {
    if (activeGroup && activeGroupPermission) {
      try {
        setIsGrantLoading(true)
        await FileAPI.grantGroupPermission({
          ids: selection,
          groupId: activeGroup.id,
          permission: activeGroupPermission,
        })
        const result = await FileAPI.batchGet({ ids: selection })
        dispatch(filesUpdated(result))
        if (isSingleSelection) {
          await mutateGroupPermissions()
        }
        setActiveGroup(undefined)
        setIsGrantLoading(false)
        if (!isSingleSelection) {
          dispatch(sharingModalDidClose())
        }
      } catch {
        setIsGrantLoading(false)
      }
    }
  }, [
    selection,
    activeGroup,
    activeGroupPermission,
    isSingleSelection,
    dispatch,
    mutateGroupPermissions,
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
        if (isSingleSelection) {
          await mutateGroupPermissions()
        }
      } finally {
        setPermissionBeingRevoked(undefined)
      }
    },
    [selection, isSingleSelection, dispatch, mutateGroupPermissions],
  )

  return (
    <Stack direction="column" spacing={variables.spacing}>
      {groups && groups.length > 0 ? (
        <Stack direction="column" spacing={variables.spacing}>
          <GroupSelector
            value={activeGroup}
            organizationId={workspace?.organization.id}
            onConfirm={(value) => setActiveGroup(value)}
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
                setActiveGroupPermission(e.value)
              }
            }}
          />
          <Button
            leftIcon={<IconCheck />}
            colorScheme="blue"
            isLoading={isGrantLoading}
            isDisabled={!activeGroup || !activeGroupPermission}
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
            geEditorPermission(workspace.organization.permission) ? (
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
      {isSingleSelection ? (
        <>
          <hr />
          {!groupPermissions ? (
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
        </>
      ) : null}
    </Stack>
  )
}

export default Groups
