// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
import {
  Pagination,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
} from '@koupr/ui'
import cx from 'classnames'
import TaskAPI, { SortOrder } from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { useAppDispatch } from '@/store/hook'
import { mutateListUpdated } from '@/store/ui/tasks'
import TaskDrawerItem from './task-item'

const TasksList = () => {
  const dispatch = useAppDispatch()
  const [page, setPage] = useState(1)
  const {
    data: list,
    error: listError,
    isLoading: isListLoading,
    mutate: mutateList,
  } = TaskAPI.useList({ page, size: 5, sortOrder: SortOrder.Asc }, swrConfig())
  const isListError = !list && listError
  const isListEmpty = list && !listError && list.totalElements === 0
  const isListReady = list && !listError && list.totalElements > 0

  useEffect(() => {
    dispatch(mutateListUpdated(mutateList))
  }, [dispatch, mutateList])

  return (
    <>
      {isListLoading ? <SectionSpinner /> : null}
      {isListError ? <SectionError text="Failed to load tasks." /> : null}
      {isListEmpty ? <SectionPlaceholder text="There are no tasks." /> : null}
      {isListReady ? (
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
          {list.totalPages > 1 ? (
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
