// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import {
  PagePagination,
  SectionError,
  SectionSpinner,
  usePageMonitor,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import Path from '@/components/common/path'
import FileCopy from '@/components/file/file-copy'
import FileCreate from '@/components/file/file-create'
import FileDelete from '@/components/file/file-delete'
import FileInfo from '@/components/file/file-info'
import FileMove from '@/components/file/file-move'
import FileRename from '@/components/file/file-rename'
import SearchFilter from '@/components/file/file-search-filter'
import FileToolbar from '@/components/file/file-toolbar'
import FileList from '@/components/file/list'
import Insights from '@/components/insights'
import Mosaic from '@/components/mosaic'
import Sharing from '@/components/sharing'
import SnapshotDetach from '@/components/snapshot/snapshot-detach'
import SnapshotList from '@/components/snapshot/snapshot-list'
import { filePaginationSteps, filesPaginationStorage } from '@/infra/pagination'
import { decodeFileQuery } from '@/lib/helpers/query'
import { listUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { mutateUpdated, selectionUpdated } from '@/store/ui/files'

const WorkspaceFilesPage = () => {
  const navigate = useNavigate()
  const { id: workspaceId, fileId } = useParams()
  const [searchParams] = useSearchParams()
  const query = decodeFileQuery(searchParams.get('q') as string)
  const dispatch = useAppDispatch()
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
  const isSnapshotListModalOpen = useAppSelector(
    (state) => state.ui.snapshots.isListModalOpen,
  )
  const isSnapshotDetachModalOpen = useAppSelector(
    (state) => state.ui.snapshots.isDetachModalOpen,
  )
  const isMoveModalOpen = useAppSelector(
    (state) => state.ui.files.isMoveModalOpen,
  )
  const isCopyModalOpen = useAppSelector(
    (state) => state.ui.files.isCopyModalOpen,
  )
  const isCreateModalOpen = useAppSelector(
    (state) => state.ui.files.isCreateModalOpen,
  )
  const isDeleteModalOpen = useAppSelector(
    (state) => state.ui.files.isDeleteModalOpen,
  )
  const isRenameModalOpen = useAppSelector(
    (state) => state.ui.files.isRenameModalOpen,
  )
  const isInfoModalOpen = useAppSelector(
    (state) => state.ui.files.isInfoModalOpen,
  )
  const isInsightsModalOpen = useAppSelector(
    (state) => state.ui.insights.isModalOpen,
  )
  const isMosaicModalOpen = useAppSelector(
    (state) => state.ui.mosaic.isModalOpen,
  )
  const isSearchFilterModalOpen = useAppSelector(
    (state) => state.ui.searchFilter.isModalOpen,
  )
  const {
    data: workspace,
    error: workspaceError,
    isLoading: isWorkspaceLoading,
  } = WorkspaceAPI.useGet(workspaceId, swrConfig())
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: filesPaginationStorage(),
    steps: filePaginationSteps(),
  })
  const {
    data: list,
    error: listError,
    isLoading: isListLoading,
    mutate,
  } = FileAPI.useList(
    fileId!,
    {
      size,
      page,
      sortBy,
      sortOrder,
      query,
    },
    swrConfig(),
  )
  const { hasPagination } = usePageMonitor({
    totalElements: list?.totalElements ?? 0,
    totalPages: list?.totalPages ?? 1,
    steps,
  })
  const isWorkspaceError = !workspace && workspaceError
  const isWorkspaceReady = workspace && !workspaceError
  const isListError = !list && listError
  const isListReady = list && !listError

  useEffect(() => {
    if (list) {
      dispatch(listUpdated(list))
    }
  }, [list, dispatch])

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate])

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
          <div
            className={cx(
              'flex',
              'flex-col',
              'w-full',
              'gap-2.5',
              'grow',
              'overflow-hidden',
            )}
          >
            {workspace && fileId ? (
              <Path
                rootId={workspace.rootId}
                fileId={fileId}
                maxCharacters={30}
                onClick={(fileId) => {
                  dispatch(selectionUpdated([]))
                  navigate(`/workspace/${workspace.id}/file/${fileId}`)
                }}
              />
            ) : null}
            <FileToolbar list={list} />
            <div
              className={cx(
                'flex',
                'flex-col',
                'gap-1.5',
                'grow',
                'overflow-y-auto',
                'overflow-x-hidden',
              )}
            >
              <div
                className={cx(
                  'w-full',
                  'overflow-y-auto',
                  'overflow-x-hidden',
                  'border-t',
                  'border-t-gray-300',
                  'dark:border-t-gray-600',
                  {
                    'border-b': hasPagination,
                    'border-b-gray-300': hasPagination,
                    'dark:border-b-gray-600': hasPagination,
                  },
                  'py-1.5',
                  'flex-grow',
                )}
                onClick={() => dispatch(selectionUpdated([]))}
              >
                {isListLoading ? <SectionSpinner /> : null}
                {isListError ? (
                  <SectionError text="Failed to load items." />
                ) : null}
                {isListReady ? (
                  <FileList list={list} scale={iconScale} />
                ) : null}
              </div>
              {list && hasPagination ? (
                <div className={cx('self-end', 'pb-1.5')}>
                  <PagePagination
                    totalElements={list.totalElements}
                    totalPages={list.totalPages}
                    page={page}
                    size={size}
                    steps={steps}
                    setPage={setPage}
                    setSize={setSize}
                  />
                </div>
              ) : null}
            </div>
          </div>
          {list ? <Sharing list={list} /> : null}
          {isSnapshotListModalOpen ? <SnapshotList /> : null}
          {isSnapshotDetachModalOpen ? <SnapshotDetach /> : null}
          {isMoveModalOpen ? <FileMove /> : null}
          {isCopyModalOpen ? <FileCopy /> : null}
          {isCreateModalOpen ? <FileCreate /> : null}
          {isDeleteModalOpen ? <FileDelete /> : null}
          {isRenameModalOpen ? <FileRename /> : null}
          {isInfoModalOpen ? <FileInfo /> : null}
          {isInsightsModalOpen ? <Insights /> : null}
          {isMosaicModalOpen ? <Mosaic /> : null}
          {isSearchFilterModalOpen ? <SearchFilter /> : null}
        </>
      ) : null}
    </>
  )
}

export default WorkspaceFilesPage
