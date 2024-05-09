import { useEffect, useMemo } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import FileAPI from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import { swrConfig } from '@/client/options'
import Path from '@/components/common/path'
import FileCopy from '@/components/file/file-copy'
import FileCreate from '@/components/file/file-create'
import FileMove from '@/components/file/file-move'
import FileRename from '@/components/file/file-rename'
import FileToolbar from '@/components/file/file-toolbar'
import FileDelete from '@/components/file/fle-idelete'
import FileList from '@/components/file/list'
import FileSharing from '@/components/file/sharing'
import FileSnapshots from '@/components/file/snapshots'
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
import { selectionUpdated } from '@/store/ui/files'

const WorkspaceFilesPage = () => {
  const navigate = useNavigate()
  const { id, fileId } = useParams()
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const dispatch = useAppDispatch()
  const sortBy = useAppSelector((state) => state.ui.files.sortBy)
  const sortOrder = useAppSelector((state) => state.ui.files.sortOrder)
  const iconScale = useAppSelector((state) => state.ui.files.iconScale)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const { data: workspace } = WorkspaceAPI.useGetById(id, swrConfig())
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
  const singleFile = useMemo(() => {
    if (list && selection.length === 1) {
      return list.data.find((e) => e.id === selection[0])
    } else {
      return undefined
    }
  }, [list?.data, selection])

  useEffect(() => {
    if (list) {
      dispatch(listUpdated(list))
    }
  }, [list, dispatch])

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
              'pt-1.5',
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
      {singleFile && singleFile.type === 'file' ? (
        <FileSnapshots file={singleFile} />
      ) : null}
      <FileMove />
      <FileCopy />
      <FileCreate />
      <FileDelete />
      <FileRename />
    </>
  )
}

export default WorkspaceFilesPage
