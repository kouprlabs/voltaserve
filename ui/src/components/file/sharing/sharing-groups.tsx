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
import { IconAdd, IconCheck, IconTrash } from '@koupr/ui'
import { KeyedMutator, useSWRConfig } from 'swr'
import { Select } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI, { GroupPermission, List } from '@/client/api/file'
import { Group } from '@/client/api/group'
import { geEditorPermission } from '@/client/api/permission'
import WorkspaceAPI from '@/client/api/workspace'
import GroupSelector from '@/components/common/group-selector'
import useFileListSearchParams from '@/hooks/use-file-list-params'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import reactSelectStyles from '@/styles/react-select'
import SharingFormSkeleton from './sharing-form-skeleton'

export type SharingGroupsProps = {
  groups?: Group[]
  permissions?: GroupPermission[]
  mutateGroupPermissions: KeyedMutator<GroupPermission[]>
}

const SharingGroups = ({
  groups,
  permissions,
  mutateGroupPermissions,
}: SharingGroupsProps) => {
  const { mutate } = useSWRConfig()
  const { id, fileId } = useParams()
  const dispatch = useAppDispatch()
  const { data: workspace } = WorkspaceAPI.useGetById(id)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const [isGrantLoading, setIsGrantLoading] = useState(false)
  const [permissionBeingRevoked, setPermissionBeingRevoked] = useState<string>()
  const [activeGroup, setActiveGroup] = useState<Group>()
  const [activePermission, setActivePermission] = useState<string>()
  const fileListSearchParams = useFileListSearchParams()
  const isSingleSelection = selection.length === 1

  const handleGrantGroupPermission = useCallback(async () => {
    if (activeGroup && activePermission) {
      try {
        setIsGrantLoading(true)
        await FileAPI.grantGroupPermission({
          ids: selection,
          groupId: activeGroup.id,
          permission: activePermission,
        })
        await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
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
    fileId,
    selection,
    activeGroup,
    activePermission,
    isSingleSelection,
    fileListSearchParams,
    mutate,
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
        await mutate<List>(`/files/${fileId}/list?${fileListSearchParams}`)
        if (isSingleSelection) {
          await mutateGroupPermissions()
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
      mutateGroupPermissions,
    ],
  )

  return (
    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
      {!groups ? <SharingFormSkeleton /> : null}
      {groups && groups.length > 0 ? (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
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
                setActivePermission(e.value)
              }
            }}
          />
          <Button
            leftIcon={<IconCheck />}
            colorScheme="blue"
            isLoading={isGrantLoading}
            isDisabled={!activeGroup || !activePermission}
            onClick={() => handleGrantGroupPermission()}
          >
            Apply to Group
          </Button>
        </div>
      ) : null}
      {groups && groups.length === 0 ? (
        <div className={cx('flex', 'items-center', 'justify-center')}>
          <div className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}>
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
          </div>
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
              <Text>Not shared with any groups.</Text>
            </div>
          ) : null}
          {permissions && permissions.length > 0 ? (
            <Table>
              <Thead>
                <Tr>
                  <Th>Group</Th>
                  <Th>Permission</Th>
                  <Th />
                </Tr>
              </Thead>
              <Tbody>
                {permissions.map((p) => (
                  <Tr key={p.id}>
                    <Td p={variables.spacingSm}>
                      <div
                        className={cx(
                          'flex',
                          'flex-row',
                          'items-center',
                          'gap-1',
                        )}
                      >
                        <Avatar
                          name={p.group.name}
                          size="sm"
                          width="40px"
                          height="40px"
                        />
                        <Text noOfLines={1}>{p.group.name}</Text>
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
    </div>
  )
}

export default SharingGroups
