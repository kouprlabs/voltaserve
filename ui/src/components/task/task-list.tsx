import { useEffect, useState } from 'react'
import cx from 'classnames'
import TaskAPI, { SortOrder } from '@/client/api/task'
import { swrConfig } from '@/client/options'
import Pagination from '@/lib/components/pagination'
import { useAppDispatch } from '@/store/hook'
import { mutateListUpdated } from '@/store/ui/tasks'
import TaskDrawerItem from './task-item'

const TasksList = () => {
  const dispatch = useAppDispatch()
  const [page, setPage] = useState(1)
  const { data: list, mutate: mutateList } = TaskAPI.useList(
    { page, size: 5, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    dispatch(mutateListUpdated(mutateList))
  }, [dispatch, mutateList])

  return (
    <>
      {list && list.data.length > 0 ? (
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
              uiSize="sm"
              maxButtons={3}
              page={page}
              totalPages={list.totalPages}
              onPageChange={(value) => setPage(value)}
            />
          ) : null}
        </div>
      ) : (
        <span>There are no tasks.</span>
      )}
    </>
  )
}

export default TasksList
