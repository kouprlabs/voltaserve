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
import FileAPI from '@/client/api/file'
import TaskAPI, { Status, Task } from '@/client/api/task'
import {
  IconCheckCircle,
  IconError,
  IconHourglass,
} from '@/lib/components/icons'
import { useAppSelector } from '@/store/hook'

export type TaskDrawerItemProps = {
  task: Task
}

const TaskDrawerItem = ({ task }: TaskDrawerItemProps) => {
  const [isDismissing, setIsDismissing] = useState(false)
  const mutateList = useAppSelector((state) => state.ui.tasks.mutateList)
  const { data: file } = FileAPI.useGet(task.payload?.fileId)

  const handleDismiss = useCallback(async () => {
    try {
      setIsDismissing(true)
      await TaskAPI.dismiss(task.id)
      mutateList?.(await TaskAPI.list())
    } finally {
      setIsDismissing(false)
    }
  }, [task, mutateList])

  return (
    <Card variant="outline">
      <CardBody>
        <div className={cx('flex', 'flex-col', 'gap-1')}>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1.5')}>
            {task.status === Status.Waiting ? <IconHourglass /> : null}
            {task.status === Status.Running ? (
              <CircularProgress
                value={task.percentage}
                max={100}
                isIndeterminate={task.isIndeterminate}
                className={cx('text-black')}
                size="20px"
              />
            ) : null}
            {task.status === Status.Success ? (
              <IconCheckCircle
                className={cx('shrink-0', 'text-green-500')}
                filled={true}
              />
            ) : null}
            {task.status === Status.Error ? (
              <IconError filled={true} className={cx('text-red-500')} />
            ) : null}
            <div className={cx('flex', 'flex-col', 'grow')}>
              {file ? (
                <Text className={cx('font-semibold')} noOfLines={5}>
                  {file.name}
                </Text>
              ) : null}
              {task.status !== Status.Error ? (
                <Text
                  dangerouslySetInnerHTML={{ __html: task.name }}
                  noOfLines={5}
                ></Text>
              ) : null}
            </div>
          </div>
          {task.error ? (
            <>
              <Accordion allowMultiple>
                <AccordionItem className={cx('border-none')}>
                  <AccordionButton className={cx('p-0.5')}>
                    <div className={cx('flex', 'flex-row', 'w-full')}>
                      <span className={cx('text-left', 'grow')}>
                        Task failed, click to expand
                      </span>
                      <AccordionIcon />
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

export default TaskDrawerItem
