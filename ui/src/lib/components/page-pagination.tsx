// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import React, { ChangeEvent, useCallback } from 'react'
import { Select } from '@chakra-ui/react'
import cx from 'classnames'
import usePageMonitor from '../hooks/page-monitor'
import Pagination from './pagination'

type PagePaginationProps = {
  totalPages: number
  totalElements: number
  page: number
  size: number
  steps: number[]
  uiSize?: string
  style?: React.CSSProperties
  setPage: (page: number) => void
  setSize: (size: number) => void
}

const PagePagination = ({
  totalElements,
  totalPages,
  page,
  size,
  uiSize = 'md',
  steps,
  style,
  setPage,
  setSize,
}: PagePaginationProps) => {
  const { hasPageSwitcher, hasSizeSelector } = usePageMonitor({
    totalElements,
    totalPages,
    steps,
  })

  const handleSizeChange = useCallback(
    (event: ChangeEvent<HTMLSelectElement>) => {
      setSize(parseInt(event.target.value))
      setPage(1)
    },
    [setSize, setPage],
  )

  return (
    <>
      {!hasPageSwitcher && !hasSizeSelector ? null : (
        <div
          className={cx('flex', 'flex-row', 'items-center', 'gap-1.5')}
          style={style}
        >
          {hasPageSwitcher ? (
            <Pagination
              uiSize={uiSize}
              page={page}
              totalPages={totalPages}
              onPageChange={setPage}
            />
          ) : null}
          {hasSizeSelector ? (
            <Select defaultValue={size} onChange={handleSizeChange}>
              {steps.map((step, index) => (
                <option key={index} value={step.toString()}>
                  {step} items
                </option>
              ))}
            </Select>
          ) : null}
        </div>
      )}
    </>
  )
}

export default PagePagination
