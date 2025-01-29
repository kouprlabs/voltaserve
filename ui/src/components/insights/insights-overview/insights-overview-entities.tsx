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
import { Badge, Table, Tbody, Td, Tr } from '@chakra-ui/react'
import {
  Pagination,
  SearchInput,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  usePageMonitor,
} from '@koupr/ui'
import cx from 'classnames'
import {
  InsightsAPI,
  InsightsSortBy,
  InsightsSortOrder,
} from '@/client/api/insights'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'

const InsightsOverviewEntities = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [page, setPage] = useState(1)
  const [query, setQuery] = useState<string | undefined>(undefined)
  const { data: metadata } = InsightsAPI.useGetInfo(id, swrConfig())
  const size = 5
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
    mutate,
  } = InsightsAPI.useListEntities(
    metadata ? id : undefined,
    {
      query,
      page,
      size,
      sortBy: InsightsSortBy.Frequency,
      sortOrder: InsightsSortOrder.Desc,
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
      {listIsLoading ? <SectionSpinner /> : null}
      {listError ? <SectionError text={errorToString(listError)} /> : null}
      {listIsEmpty ? (
        <SectionPlaceholder text="There are no entities." />
      ) : null}
      {listIsReady ? (
        <div
          className={cx(
            'flex',
            'flex-col',
            'justify-between',
            'gap-1.5',
            'h-[320px]',
          )}
        >
          <Table variant="simple" size="sm">
            <colgroup>
              <col className={cx('w-[40px]')} />
              <col className={cx('w-[auto]')} />
            </colgroup>
            <Tbody>
              {list.data.map((entity, index) => (
                <Tr key={index} className={cx('h-[52px]')}>
                  <Td className={cx('px-0.5')}>
                    <div
                      className={cx(
                        'flex',
                        'flex-row',
                        'items-center',
                        'gap-1.5',
                      )}
                    >
                      <span className={cx('text-base')}>{entity.text}</span>
                      <Badge>{entity.frequency}</Badge>
                    </div>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
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
        </div>
      ) : null}
    </div>
  )
}

export default InsightsOverviewEntities
