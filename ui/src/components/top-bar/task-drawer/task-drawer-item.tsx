import { Card, CardBody } from '@chakra-ui/react'
import cx from 'classnames'
import { Task } from '@/client/api/task'

export type TaskDrawerItemProps = {
  task: Task
}

const TaskDrawerItem = ({ task }: TaskDrawerItemProps) => {
  return (
    <Card variant="outline">
      <CardBody>
        <div className={cx('flex', 'flex-col', 'gap-0.5')}>{task.name}</div>
      </CardBody>
    </Card>
  )
}

export default TaskDrawerItem
