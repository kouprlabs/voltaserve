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
import { IconButton, IconButtonProps, Progress } from '@chakra-ui/react'
import {
  Form,
  IconDelete,
  IconEdit,
  SectionError,
  SectionSpinner,
} from '@koupr/ui'
import cx from 'classnames'
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

const WorkspaceSettingsPage = () => {
  const { id } = useParams()
  const {
    data: workspace,
    error: workspaceError,
    isLoading: isWorkspaceLoading,
    mutate,
  } = WorkspaceAPI.useGet(id, swrConfig())
  const {
    data: storageUsage,
    error: storageUsageError,
    isLoading: isStorageUsageLoading,
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
  const isWorkspaceError = !workspace && workspaceError
  const isWorkspaceReady = workspace && !workspaceError
  const isStorageUsageError = !storageUsage && storageUsageError
  const isStorageUsageReady = storageUsage && !storageUsageError

  const handleEditNameClose = useCallback(async () => {
    setIsNameModalOpen(false)
    await mutate()
  }, [])

  const handleEditStorageCapacityClose = useCallback(async () => {
    setIsStorageCapacityModalOpen(false)
    await mutateStorageUsage()
  }, [])

  return (
    <>
      {isWorkspaceLoading ? (
        <div className={cx('block')}>
          <SectionSpinner />
        </div>
      ) : null}
      {isWorkspaceError ? (
        <div className={cx('block')}>
          <SectionError text="Failed to load workspace." />
        </div>
      ) : null}
      {isWorkspaceReady ? (
        <>
          <Form
            sections={[
              {
                title: 'Storage',
                content: (
                  <>
                    {isStorageUsageError ? (
                      <SectionError
                        text="Failed to load storage usage."
                        height="auto"
                      />
                    ) : null}
                    {isStorageUsageReady ? (
                      <>
                        <span>
                          {prettyBytes(storageUsage.bytes)} of{' '}
                          {prettyBytes(storageUsage.maxBytes)} used
                        </span>
                        <Progress value={storageUsage.percentage} hasStripe />
                      </>
                    ) : null}
                    {isStorageUsageLoading ? (
                      <>
                        <span>Calculating…</span>
                        <Progress value={0} hasStripe />
                      </>
                    ) : null}
                  </>
                ),
                rows: [
                  {
                    label: 'Capacity',
                    content: (
                      <>
                        <span>{prettyBytes(workspace.storageCapacity)}</span>
                        <EditButton
                          title="Edit storage capacity"
                          aria-label="Edit storage capacity"
                          isDisabled={!hasEditPermission}
                          onClick={() => setIsStorageCapacityModalOpen(true)}
                        />
                      </>
                    ),
                  },
                ],
              },
              {
                title: 'Basics',
                rows: [
                  {
                    label: 'Name',
                    content: (
                      <>
                        <span>{truncateEnd(workspace.name, 50)}</span>
                        <EditButton
                          title="Edit name"
                          aria-label="Edit name"
                          isDisabled={!hasEditPermission}
                          onClick={() => setIsNameModalOpen(true)}
                        />
                      </>
                    ),
                  },
                ],
              },
              {
                title: 'Advanced',
                rows: [
                  {
                    label: 'Delete workspace',
                    content: (
                      <IconButton
                        icon={<IconDelete />}
                        variant="solid"
                        colorScheme="red"
                        isDisabled={!hasEditPermission}
                        title="Delete workspace"
                        aria-label="Delete workspace"
                        onClick={() => setIsDeleteModalOpen(true)}
                      />
                    ),
                  },
                ],
              },
            ]}
          />
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
      ) : null}
    </>
  )
}

export default WorkspaceSettingsPage
