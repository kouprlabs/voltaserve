import { useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Divider,
  IconButton,
  IconButtonProps,
  Progress,
} from '@chakra-ui/react'
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
import { IconDelete, IconEdit } from '@/lib/components/icons'
import SectionSpinner from '@/lib/components/section-spinner'

const EditButton = (props: IconButtonProps) => (
  <IconButton icon={<IconEdit />} {...props} />
)

const Spacer = () => <div className={cx('grow')} />

const WorkspaceSettingsPage = () => {
  const { id } = useParams()
  const {
    data: workspace,
    error: workspaceError,
    mutate,
  } = WorkspaceAPI.useGet(id, swrConfig())
  const {
    data: storageUsage,
    error: storageUsageError,
    mutate: mutateStorageUsage,
  } = StorageAPI.useGetWorkspaceUsage(id, swrConfig())
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
          <span className={cx('font-bold')}>Storage</span>
          {storageUsageError && <span>Failed to load storage usage.</span>}
          {storageUsage && !storageUsageError ? (
            <>
              <span>
                {prettyBytes(storageUsage.bytes)} of{' '}
                {prettyBytes(storageUsage.maxBytes)} used
              </span>
              <Progress value={storageUsage.percentage} hasStripe />
            </>
          ) : null}
          {!storageUsage && !storageUsageError ? (
            <>
              <span>Calculating…</span>
              <Progress value={0} hasStripe />
            </>
          ) : null}
          <Divider />
          <div className={rowClassName}>
            <span>Storage capacity</span>
            <Spacer />
            <span>{prettyBytes(workspace.storageCapacity)}</span>
            <EditButton
              aria-label=""
              isDisabled={!hasEditPermission}
              onClick={() => {
                setIsStorageCapacityModalOpen(true)
              }}
            />
          </div>
          <Divider className={cx('mb-1.5')} />
          <span className={cx('font-bold')}>Basics</span>
          <div className={rowClassName}>
            <span>Name</span>
            <Spacer />
            <span>{workspace.name}</span>
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
          <span className={cx('font-bold')}>Advanced</span>
          <div className={rowClassName}>
            <span>Delete permanently</span>
            <Spacer />
            <IconButton
              icon={<IconDelete />}
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
        onClose={() => {
          setIsNameModalOpen(false)
          mutate()
        }}
      />
      <WorkspaceEditStorageCapacity
        open={isStorageCapacityModalOpen}
        workspace={workspace}
        onClose={() => {
          setIsStorageCapacityModalOpen(false)
          mutateStorageUsage()
        }}
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
