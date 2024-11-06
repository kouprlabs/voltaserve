// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  Divider,
  IconButton,
  IconButtonProps,
  Progress,
} from '@chakra-ui/react'
import { IconDelete, IconEdit, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import { geEditorPermission } from '@/client/api/permission'
import StorageAPI from '@/client/api/storage'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import WorkspaceDelete from '@/components/workspace/workspace-delete'
import WorkspaceEditName from '@/components/workspace/workspace-edit-name'
import WorkspaceEditStorageCapacity from '@/components/workspace/workspace-edit-storage-capacity'
import prettyBytes from '@/lib/helpers/pretty-bytes'
import { truncateEnd } from '@/lib/helpers/truncate-end'

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

  const handleEditNameClose = useCallback(async () => {
    setIsNameModalOpen(false)
    await mutate()
  }, [])

  const handleEditStorageCapacityClose = useCallback(async () => {
    setIsStorageCapacityModalOpen(false)
    await mutateStorageUsage()
  }, [])

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
          {storageUsageError ? (
            <span>Failed to load storage usage.</span>
          ) : null}
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
          <div className={rowClassName}>
            <span>Capacity</span>
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
        </div>
        <Divider />
        <div className={sectionClassName}>
          <span className={cx('font-bold')}>Basics</span>
          <div className={rowClassName}>
            <span>Name</span>
            <Spacer />
            <span>{truncateEnd(workspace.name, 50)}</span>
            <EditButton
              aria-label=""
              isDisabled={!hasEditPermission}
              onClick={() => {
                setIsNameModalOpen(true)
              }}
            />
          </div>
        </div>
        <Divider />
        <div className={sectionClassName}>
          <span className={cx('font-bold')}>Advanced</span>
          <div className={rowClassName}>
            <span>Delete workspace</span>
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
        onClose={handleEditNameClose}
      />
      <WorkspaceEditStorageCapacity
        open={isStorageCapacityModalOpen}
        workspace={workspace}
        onClose={handleEditStorageCapacityClose}
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
