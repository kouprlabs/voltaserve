import { useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Box,
  Divider,
  HStack,
  IconButton,
  IconButtonProps,
  Progress,
  Stack,
  Text,
} from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import { swrConfig } from '@/api/options'
import { geEditorPermission } from '@/api/permission'
import StorageAPI from '@/api/storage'
import WorkspaceAPI from '@/api/workspace'
import { IconEdit, IconTrash } from '@/components/common/icon'
import LoadingSpinner from '@/components/common/loading-spinner'
import WorkspaceDelete from '@/components/workspace/delete'
import WorkspaceEditName from '@/components/workspace/edit-name'
import WorkspaceEditStorageCapacity from '@/components/workspace/edit-storage-capacity'
import variables from '@/theme/variables'
import prettyBytes from '@/helpers/pretty-bytes'

const EditButton = (props: IconButtonProps) => (
  <IconButton icon={<IconEdit />} {...props} />
)

const Spacer = () => <Box flexGrow={1} />

const WorkspaceSettingsPage = () => {
  const params = useParams()
  const workspaceId = params.id as string
  const { data: workspace, error: workspaceError } = WorkspaceAPI.useGetById(
    workspaceId,
    swrConfig()
  )
  const { data: storageUsage, error: storageUsageError } =
    StorageAPI.useGetWorkspaceUsage(workspaceId, swrConfig())
  const hasEditPermission = useMemo(
    () => workspace && geEditorPermission(workspace.permission),
    [workspace]
  )
  const [isNameModalOpen, setIsNameModalOpen] = useState(false)
  const [isStorageCapacityModalOpen, setIsStorageCapacityModalOpen] =
    useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)

  if (workspaceError) {
    return null
  }

  if (!workspace) {
    return <LoadingSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{workspace.name}</title>
      </Helmet>
      <Stack spacing={variables.spacing} w="100%">
        <Text fontWeight="bold">Storage</Text>
        {storageUsageError && <Text>Failed to load storage usage.</Text>}
        {storageUsage && !storageUsageError && (
          <>
            <Text>
              {prettyBytes(storageUsage.bytes)} of{' '}
              {prettyBytes(storageUsage.maxBytes)} used
            </Text>
            <Progress value={storageUsage.percentage} hasStripe />
          </>
        )}
        {!storageUsage && !storageUsageError && (
          <>
            <Text>Calculating…</Text>
            <Progress value={0} hasStripe />
          </>
        )}
        <Divider />
        <HStack spacing={variables.spacing}>
          <Text>Storage capacity</Text>
          <Spacer />
          <Text>{prettyBytes(workspace.storageCapacity)}</Text>
          <EditButton
            aria-label=""
            isDisabled={!hasEditPermission}
            onClick={() => {
              setIsStorageCapacityModalOpen(true)
            }}
          />
        </HStack>
        <Divider mb={variables.spacing} />
        <Text fontWeight="bold">Basics</Text>
        <HStack spacing={variables.spacing}>
          <Text>Name</Text>
          <Spacer />
          <Text>{workspace.name}</Text>
          <EditButton
            aria-label=""
            isDisabled={!hasEditPermission}
            onClick={() => {
              setIsNameModalOpen(true)
            }}
          />
        </HStack>
        <Divider />
        <HStack spacing={variables.spacing}>
          <Text>Delete permanently</Text>
          <Spacer />
          <IconButton
            icon={<IconTrash />}
            variant="solid"
            colorScheme="red"
            isDisabled={!hasEditPermission}
            aria-label=""
            onClick={() => setIsDeleteModalOpen(true)}
          />
        </HStack>
        <WorkspaceEditName
          open={isNameModalOpen}
          workspace={workspace}
          onClose={() => setIsNameModalOpen(false)}
        />
        <WorkspaceEditStorageCapacity
          open={isStorageCapacityModalOpen}
          workspace={workspace}
          onClose={() => setIsStorageCapacityModalOpen(false)}
        />
        <WorkspaceDelete
          open={isDeleteModalOpen}
          workspace={workspace}
          onClose={() => setIsDeleteModalOpen(false)}
        />
      </Stack>
    </>
  )
}

export default WorkspaceSettingsPage
