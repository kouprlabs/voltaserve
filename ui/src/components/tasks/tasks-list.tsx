import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import cx from 'classnames'
import TaskAPI, { SortOrder } from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { taskPaginationStorage } from '@/infra/pagination'
import PagePagination from '@/lib/components/page-pagination'
import usePagePagination from '@/lib/hooks/page-pagination'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/tasks'
import TasksDrawerItem from './tasks-item'

const TasksList = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigate,
    location,
    storage: taskPaginationStorage(),
  })
  const { data: list, mutate } = TaskAPI.useList(
    { page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate])

  return (
    <>
      {list && list.data.length > 0 ? (
        <div className={cx('flex', 'flex-col', 'gap-1.5')}>
          {list.data.map((task) => (
            <div key={task.id} className={cx('flex', 'flex-col', 'gap-1.5')}>
              <TasksDrawerItem task={task} />
            </div>
          ))}
        </div>
      ) : (
        <span>There are no tasks.</span>
      )}
      {list ? (
        <PagePagination
          style={{ alignSelf: 'end' }}
          totalElements={list.totalElements}
          totalPages={list.totalPages}
          page={page}
          size={size}
          steps={steps}
          setPage={setPage}
          setSize={setSize}
        />
      ) : null}
    </>
  )
}

export default TasksList
