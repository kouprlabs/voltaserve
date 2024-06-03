import { useCallback, useState } from 'react'
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Button,
  Card,
  CardBody,
  Text,
  CircularProgress,
} from '@chakra-ui/react'
import cx from 'classnames'
import TaskAPI, { Task } from '@/client/api/task'
import { IconError } from '@/lib/components/icons'

export type TaskDrawerItemProps = {
  task: Task
}

const TasksDrawerItem = ({ task }: TaskDrawerItemProps) => {
  const [isDismissing, setIsDismissing] = useState(false)

  const handleDismiss = useCallback(async () => {
    try {
      setIsDismissing(true)
      await TaskAPI.delete(task.id)
    } finally {
      setIsDismissing(false)
    }
  }, [task])

  return (
    <Card variant="outline">
      <CardBody>
        <div className={cx('flex', 'flex-col', 'gap-1')}>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1.5')}>
            {task.error ? (
              <IconError filled={true} className={cx('text-red-500')} />
            ) : (
              <CircularProgress
                value={task.percentage}
                max={100}
                isIndeterminate={task.isIndeterminate}
                className={cx('text-black')}
                size="20px"
              />
            )}
            <Text
              className={cx('grow')}
              dangerouslySetInnerHTML={{ __html: task.name }}
              noOfLines={5}
            ></Text>
          </div>
          {task.error ? (
            <>
              <Accordion allowMultiple>
                <AccordionItem className={cx('border-none')}>
                  <AccordionButton className={cx('p-0.5', 'hover:bg-red-50')}>
                    <div className={cx('flex', 'flex-row', 'w-full')}>
                      <span className={cx('text-red-500', 'text-left', 'grow')}>
                        Task failed, click to expand
                      </span>
                      <AccordionIcon className={cx('text-red-500')} />
                    </div>
                  </AccordionButton>
                  <AccordionPanel className={cx('p-0.5')}>
                    <Text
                      className={cx('text-red-500')}
                      dangerouslySetInnerHTML={{ __html: task.error }}
                      noOfLines={5}
                    ></Text>
                  </AccordionPanel>
                </AccordionItem>
              </Accordion>
              <div
                className={cx(
                  'flex',
                  'flex-row',
                  'gap-0.5',
                  'w-full',
                  'justify-end',
                )}
              >
                <Button
                  size="sm"
                  variant="solid"
                  isLoading={isDismissing}
                  onClick={handleDismiss}
                >
                  Dismiss
                </Button>
              </div>
            </>
          ) : null}
        </div>
      </CardBody>
    </Card>
  )
}

export default TasksDrawerItem
