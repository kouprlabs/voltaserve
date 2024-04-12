import { useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Divider,
  IconButton,
  IconButtonProps,
  Progress,
  Text,
} from '@chakra-ui/react'
import { IconEdit, IconTrash, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import { geEditorPermission } from '@/client/api/permission'
import StorageAPI from '@/client/api/storage'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import WorkspaceDelete from '@/components/workspace/workspace-delete'
import WorkspaceEditName from '@/components/workspace/workspace-edit-name'
import WorkspaceEditStorageCapacity from '@/components/workspace/workspace-edit-storage-capacity'
import prettyBytes from '@/helpers/pretty-bytes'

const EditButton = (props: IconButtonProps) => (
  <IconButton icon={<IconEdit />} {...props} />
)

const Spacer = () => <div className={cx('grow')} />

const WorkspaceSettingsPage = () => {
  const { id } = useParams()
  const { data: workspace, error: workspaceError } = WorkspaceAPI.useGetById(
    id,
    swrConfig(),
  )
  const { data: storageUsage, error: storageUsageError } =
    StorageAPI.useGetWorkspaceUsage(id, swrConfig())
  const hasEditPermission = useMemo(
    () => workspace && geEditorPermission(workspace.permission),
    [workspace],
  )
  const [isNameModalOpen, setIsNameModalOpen] = useState(false)
  const [isStorageCapacityModalOpen, setIsStorageCapacityModalOpen] =
    useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const sectionClassName = cx('flex', 'flex-col', 'gap-1', 'py-1.5')
  const rowClassName = cx(
    'flex',
    'flex-row',
    'items-center',
    'gap-1',
    `h-[40px]`,
  )

  if (workspaceError) {
    return null
  }

  if (!workspace) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{workspace.name}</title>
      </Helmet>
      <div className={cx('flex', 'flex-col', 'gap-0')}>
        <div className={sectionClassName}>
          <Text className={cx('font-bold')}>Storage</Text>
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
          <div className={rowClassName}>
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
          </div>
          <Divider className={cx('mb-1.5')} />
          <Text className={cx('font-bold')}>Basics</Text>
          <div className={rowClassName}>
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
          </div>
        </div>
        <div className={sectionClassName}>
          <Text className={cx('font-bold')}>Advanced</Text>
          <div className={rowClassName}>
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
          </div>
        </div>
      </div>
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
    </>
  )
}

export default WorkspaceSettingsPage
