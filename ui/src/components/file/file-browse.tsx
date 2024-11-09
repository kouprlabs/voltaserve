// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Button } from '@chakra-ui/react'
import { IconChevronRight, Text, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import FileAPI, { File, FileType } from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import Path from '@/components/common/path'
import FolderSvg from '@/components/file/list/item/icon/icon-folder/assets/icon-folder.svg'

export type FileBrowseProps = {
  onChange?: (id: string) => void
}

const FileBrowse = ({ onChange }: FileBrowseProps) => {
  const { id } = useParams()
  const { data: workspace } = WorkspaceAPI.useGet(id)
  const [folders, setFolders] = useState<File[]>([])
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [isLoading, setIsLoading] = useState(false)
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
      setIsLoading(true)
      const result = await FileAPI.list(fileId, {
        page,
        type: FileType.Folder,
      })
      setTotalPages(result.totalPages)
      setFolders(result.data)
      setPage(page + 1)
    } finally {
      setIsLoading(false)
    }
  }, [])

  return (
    <>
      {isSpinnerVisible ? (
        <SectionSpinner />
      ) : (
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
                    'select-none',
                    'p-1',
                    'rounded-md',
                    'hover:bg-gray-100',
                    'hover:dark:bg-gray-700',
                    'active:bg-gray-100',
                    'active:dark:bg-gray-700',
                  )}
                  onClick={() => setFileId(f.id)}
                >
                  <img
                    src={FolderSvg}
                    className={cx('shrink-0', 'w-[36px]', 'h-[28.84px]')}
                  />
                  <Text noOfLines={1}>{f.name}</Text>
                  <div className={cx('grow')} />
                  <IconChevronRight />
                </div>
              ))
            ) : (
              <div
                className={cx(
                  'flex',
                  'items-center',
                  'justify-center',
                  'h-full',
                )}
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
                isLoading={isLoading}
              >
                Load More
              </Button>
            </div>
          ) : null}
        </div>
      )}
    </>
  )
}

export default FileBrowse
