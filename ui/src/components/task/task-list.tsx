// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect, useState } from 'react'
import {
  Pagination,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  usePageMonitor,
} from '@koupr/ui'
import cx from 'classnames'
import TaskAPI, { SortBy, SortOrder } from '@/client/api/task'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppDispatch } from '@/store/hook'
import { mutateListUpdated } from '@/store/ui/tasks'
import TaskDrawerItem from './task-item'

const TasksList = () => {
  const dispatch = useAppDispatch()
  const [page, setPage] = useState(1)
  const size = 5
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
    mutate: mutateList,
  } = TaskAPI.useList(
    { page, size, sortOrder: SortOrder.Desc, sortBy: SortBy.DateCreated },
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
    dispatch(mutateListUpdated(mutateList))
  }, [dispatch, mutateList])

  return (
    <>
      {listIsLoading ? <SectionSpinner /> : null}
      {listError ? <SectionError text={errorToString(listError)} /> : null}
      {listIsEmpty ? <SectionPlaceholder text="There are no tasks." /> : null}
      {listIsReady ? (
        <div
          className={cx(
            'flex',
            'flex-col',
            'gap-1.5',
            'justify-between',
            'items-center',
            'h-full',
          )}
        >
          <div
            className={cx(
              'flex',
              'flex-col',
              'gap-1.5',
              'w-full',
              'overflow-y-auto',
            )}
          >
            {list.data.map((task) => (
              <TaskDrawerItem key={task.id} task={task} />
            ))}
          </div>
          {hasPageSwitcher ? (
            <Pagination
              size="sm"
              maxButtons={3}
              page={page}
              totalPages={list.totalPages}
              onPageChange={(value) => setPage(value)}
            />
          ) : null}
        </div>
      ) : null}
    </>
  )
}

export default TasksList
