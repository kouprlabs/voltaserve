import { useCallback, useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Button } from '@chakra-ui/react'
import cx from 'classnames'
import { FcFolder } from 'react-icons/fc'
import FileAPI, { File, FileType } from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import Path from '@/components/common/path'
import { IconChevronRight, SectionSpinner, Text } from '@/lib'

export type FileBrowseProps = {
  onChange?: (id: string) => void
}

const FileBrowse = ({ onChange }: FileBrowseProps) => {
  const { id } = useParams()
  const { data: workspace } = WorkspaceAPI.useGetById(id)
  const [folders, setFolders] = useState<File[]>([])
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [loading, setLoading] = useState(false)
  const [isSpinnerVisible, setIsSpinnerVisible] = useState(false)
  const [fileId, setFileId] = useState<string>()

  useEffect(() => {
    if (workspace) {
      setFileId(workspace.rootId)
    }
  }, [workspace])

  useEffect(() => {
    ;(async () => {
      if (fileId) {
        try {
          const timeoutId = setTimeout(() => setIsSpinnerVisible(true), 250)
          const result = await FileAPI.list(fileId, {
            page: 1,
            type: FileType.Folder,
          })
          clearTimeout(timeoutId)
          setTotalPages(result.totalPages)
          setFolders(result.data)
        } finally {
          setIsSpinnerVisible(false)
        }
      }
    })()
  }, [fileId])

  useEffect(() => {
    if (fileId) {
      onChange?.(fileId)
    }
  }, [fileId, onChange])

  const handleLoadMore = useCallback(async (fileId: string, page: number) => {
    try {
      setLoading(true)
      const result = await FileAPI.list(fileId, {
        page,
        type: FileType.Folder,
      })
      setTotalPages(result.totalPages)
      setFolders(result.data)
      setPage(page + 1)
    } finally {
      setLoading(false)
    }
  }, [])

  if (isSpinnerVisible) {
    return <SectionSpinner />
  }

  return (
    <div className={cx('flex', 'flex-col', 'gap-1')}>
      {workspace && fileId ? (
        <Path
          rootId={workspace.rootId}
          fileId={fileId}
          maxCharacters={10}
          onClick={(fileId) => setFileId(fileId)}
        />
      ) : null}
      <div
        className={cx(
          'flex',
          'flex-col',
          'gap-0',
          'border-t',
          'pt-1.5',
          'h-[250px]',
          'xl:h-[400px]',
          'overflow-y-scroll',
          'border-t-gray-300',
          'dark:border-t-gray-600',
        )}
      >
        {folders.length > 0 ? (
          folders.map((f) => (
            <div
              key={f.id}
              className={cx(
                'flex',
                'flex-row',
                'gap-1.5',
                'items-center',
                'cursor-pointer',
                'p-1',
                'rounded-md',
                'hover:bg-gray-100',
                'hover:dark:bg-gray-700',
                'active:bg-gray-100',
                'active:dark:bg-gray-700',
              )}
              onClick={() => setFileId(f.id)}
            >
              <FcFolder fontSize="36px" className={cx('shrink-0')} />
              <Text noOfLines={1}>{f.name}</Text>
              <div className={cx('grow')} />
              <IconChevronRight />
            </div>
          ))
        ) : (
          <div
            className={cx('flex', 'items-center', 'justify-center', 'h-full')}
          >
            <span>There are no folders.</span>
          </div>
        )}
      </div>
      {totalPages > page && fileId ? (
        <div
          className={cx(
            'flex',
            'items-center',
            'justify-center',
            'w-full',
            'p-1.5',
          )}
        >
          <Button
            onClick={() => handleLoadMore(fileId, page)}
            isLoading={loading}
          >
            Load More
          </Button>
        </div>
      ) : null}
    </div>
  )
}

export default FileBrowse
