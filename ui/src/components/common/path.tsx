// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useMemo } from 'react'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Skeleton,
} from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import Text from '@/lib/components/text'

export type PathProps = {
  rootId: string
  fileId: string
  maxCharacters?: number
  onClick?: (fileId: string) => void
}

const Path = ({ rootId, fileId, maxCharacters, onClick }: PathProps) => {
  const { data: path, error, isLoading } = FileAPI.useGetPath(fileId)
  const hasMore = path && path.length > 3
  const shortPath = useMemo(() => {
    if (!path) {
      return []
    }
    return hasMore ? path.slice(1).slice(-3) : path.slice(1)
  }, [hasMore, path])

  return (
    <>
      {path && !error ? (
        <Breadcrumb className={cx('overflow-hidden', 'shrink-0')}>
          <BreadcrumbItem>
            <BreadcrumbLink
              className={cx('nowrap')}
              onClick={() => onClick?.(rootId)}
            >
              Home
            </BreadcrumbLink>
          </BreadcrumbItem>
          {hasMore ? (
            <BreadcrumbItem>
              <BreadcrumbLink>…</BreadcrumbLink>
            </BreadcrumbItem>
          ) : null}
          {shortPath.map((file) => (
            <BreadcrumbItem key={file.id}>
              <BreadcrumbLink
                isCurrentPage={fileId === file.id}
                onClick={() => onClick?.(file.id)}
              >
                {maxCharacters ? (
                  <Text maxCharacters={maxCharacters}>{file.name}</Text>
                ) : (
                  file.name
                )}
              </BreadcrumbLink>
            </BreadcrumbItem>
          ))}
        </Breadcrumb>
      ) : null}
      {isLoading ? (
        <div
          className={cx(
            'flex',
            'flex-row',
            'items-center',
            'gap-0.5',
            'flex-shrink-0',
          )}
        >
          <Skeleton className={cx('w-[100px]', 'h-[20px]', 'rounded-[20px]')} />
          <span>/</span>
          <Skeleton className={cx('w-[100px]', 'h-[20px]', 'rounded-[20px]')} />
          <span>/</span>
          <Skeleton className={cx('w-[100px]', 'h-[20px]', 'rounded-[20px]')} />
        </div>
      ) : null}
    </>
  )
}

export default Path
