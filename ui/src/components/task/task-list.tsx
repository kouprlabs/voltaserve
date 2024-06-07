import { useEffect, useState } from 'react'
import cx from 'classnames'
import TaskAPI, { List, SortOrder, Status } from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { newHashId } from '@/infra/id'
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

  // const list: List = {
  //   data: [],
  //   page: 1,
  //   size: 5,
  //   totalPages: 20,
  //   totalElements: 30,
  // }
  // for (var i = 0; i < 5; i++) {
  //   list.data.push({
  //     id: newHashId(),
  //     userId: newHashId(),
  //     name: `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`,
  //     isIndeterminate: true,
  //     status: Status.Waiting,
  //     payload: { fileId: '9bPvDEAnZ5rYm' },
  //   })
  // }

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
