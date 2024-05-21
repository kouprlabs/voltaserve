import { useEffect } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import FileAPI from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import AnalysisModal from '@/components/analysis/analysis-modal'
import Path from '@/components/common/path'
import FileCopy from '@/components/file/file-copy'
import FileCreate from '@/components/file/file-create'
import FileDelete from '@/components/file/file-delete'
import FileMove from '@/components/file/file-move'
import FileRename from '@/components/file/file-rename'
import FileToolbar from '@/components/file/file-toolbar'
import FileList from '@/components/file/list'
import FileSharing from '@/components/file/sharing'
import FileSnapshotList from '@/components/file/snapshot/snapshot-list'
import FileSnapshotUnlink from '@/components/file/snapshot/snapshot-unlink'
import { decodeQuery } from '@/helpers/query'
import { filePaginationSteps, filesPaginationStorage } from '@/infra/pagination'
import {
  PagePagination,
  Spinner,
  usePageMonitor,
  usePagePagination,
  variables,
} from '@/lib'
import { listUpdated } from '@/store/entities/files'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { mutateUpdated, selectionUpdated } from '@/store/ui/files'

const WorkspaceFilesPage = () => {
  const navigate = useNavigate()
  const { id, fileId } = useParams()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const dispatch = useAppDispatch()
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
  const isSnapshotListModalOpen = useAppSelector(
    (state) => state.ui.snapshots.isListModalOpen,
  )
  const isSnapshotDeleteModalOpen = useAppSelector(
    (state) => state.ui.snapshots.isDeleteModalOpen,
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
  const isAnalysisModalOpen = useAppSelector(
    (state) => state.ui.analysis.isModalOpen,
  )
  const { data: workspace } = WorkspaceAPI.useGet(id, swrConfig())
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: filesPaginationStorage(),
    steps: filePaginationSteps(),
  })
  const {
    data: list,
    error,
    isLoading,
    mutate,
  } = FileAPI.useList(
    fileId!,
    {
      size,
      page,
      sortBy,
      sortOrder,
      query: query ? { text: query } : undefined,
    },
    swrConfig(),
  )
  const { hasPageSwitcher, hasSizeSelector } = usePageMonitor({
    totalElements: list?.totalElements || 0,
    totalPages: list?.totalPages || 1,
    steps,
  })
  const hasPagination = hasPageSwitcher || hasSizeSelector

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
      <Helmet>{workspace && <title>{workspace.name}</title>}</Helmet>
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
            {isLoading ? (
              <div
                className={cx(
                  'flex',
                  'items-center',
                  'justify-center',
                  'h-full',
                )}
              >
                <Spinner />
              </div>
            ) : null}
            {list && !error ? <FileList list={list} scale={iconScale} /> : null}
          </div>
          {list ? (
            <PagePagination
              style={{ alignSelf: 'end', paddingBottom: variables.spacing }}
              totalElements={list.totalElements}
              totalPages={list.totalPages}
              page={page}
              size={size}
              steps={steps}
              setPage={setPage}
              setSize={setSize}
            />
          ) : null}
        </div>
      </div>
      {list ? <FileSharing list={list} /> : null}
      {isSnapshotListModalOpen ? <FileSnapshotList /> : null}
      {isSnapshotDeleteModalOpen ? <FileSnapshotUnlink /> : null}
      {isMoveModalOpen ? <FileMove /> : null}
      {isCopyModalOpen ? <FileCopy /> : null}
      {isCreateModalOpen ? <FileCreate /> : null}
      {isDeleteModalOpen ? <FileDelete /> : null}
      {isRenameModalOpen ? <FileRename /> : null}
      {isAnalysisModalOpen ? <AnalysisModal /> : null}
    </>
  )
}

export default WorkspaceFilesPage
