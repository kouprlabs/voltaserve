// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  IconChevronRight,
  Text,
  SectionSpinner,
  SearchInput,
  usePageMonitor,
  SectionError,
  SectionPlaceholder,
  Pagination,
} from '@koupr/ui'
import cx from 'classnames'
import FileAPI, { FileType } from '@/client/api/file'
import WorkspaceAPI from '@/client/api/workspace'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import Path from '@/components/common/path'
import FolderSvg from '@/components/file/list/item/icon/icon-folder/assets/icon-folder.svg'

export type FileBrowseProps = {
  onChange?: (id: string) => void
}

const FileBrowse = ({ onChange }: FileBrowseProps) => {
  const { id } = useParams()
  const [fileId, setFileId] = useState<string>()
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState<string | undefined>(undefined)
  const { data: workspace } = WorkspaceAPI.useGet(id)
  const size = 5
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
    mutate,
  } = FileAPI.useList(
    fileId,
    {
      page,
      size,
      query: {
        text: query,
        type: FileType.Folder,
      },
    },
    swrConfig(),
  )
  const { hasPageSwitcher } = usePageMonitor({
    totalPages: list?.totalPages ?? 1,
    totalElements: list?.totalElements ?? 0,
    steps: [size],
  })
  const listIsEmpty = list && !listError && list.totalElements === 0
  const listIsReady = list && !listError && list.totalElements > 0

  useEffect(() => {
    if (workspace) {
      setFileId(workspace.rootId)
    }
  }, [workspace])

  useEffect(() => {
    if (fileId) {
      setQuery(undefined)
      onChange?.(fileId)
    }
  }, [fileId, onChange])

  useEffect(() => {
    mutate().then()
  }, [page, query, mutate])

  const handleSearchInputValue = useCallback((value: string) => {
    setPage(1)
    setQuery(value)
  }, [])

  const handleSearchInputClear = useCallback(() => {
    setPage(1)
    setQuery(undefined)
  }, [])

  return (
    <div className={cx('flex', 'flex-col', 'gap-1.5')}>
      <SearchInput
        placeholder="Search Entities"
        query={query}
        onValue={handleSearchInputValue}
        onClear={handleSearchInputClear}
      />
      {workspace && fileId ? (
        <Path rootId={workspace.rootId} fileId={fileId} maxCharacters={10} onClick={(fileId) => setFileId(fileId)} />
      ) : null}
      {listIsLoading ? <SectionSpinner /> : null}
      {listError ? <SectionError text={errorToString(listError)} /> : null}
      {listIsEmpty ? <SectionPlaceholder text="There are no items." /> : null}
      {listIsReady ? (
        <>
          <div
            className={cx(
              'flex',
              'flex-col',
              'gap-0',
              'border-t',
              'pt-1.5',
              'overflow-y-scroll',
              'border-t-gray-300',
              'dark:border-t-gray-600',
            )}
          >
            {list.data.map((f) => (
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
                <img src={FolderSvg} className={cx('shrink-0', 'w-[36px]', 'h-[28.84px]')} />
                <Text noOfLines={1}>{f.name}</Text>
                <div className={cx('grow')} />
                <IconChevronRight />
              </div>
            ))}
          </div>
          {hasPageSwitcher ? (
            <div className={cx('self-end')}>
              <Pagination
                maxButtons={3}
                size="sm"
                page={page}
                totalPages={list.totalPages}
                onPageChange={(value) => setPage(value)}
              />
            </div>
          ) : null}
        </>
      ) : null}
    </div>
  )
}

export default FileBrowse
