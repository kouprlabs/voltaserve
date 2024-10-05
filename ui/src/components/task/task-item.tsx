// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useState } from 'react'
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Card,
  CardBody,
  CircularProgress,
  IconButton,
  Text,
} from '@chakra-ui/react'
import cx from 'classnames'
import TaskAPI, { Status, Task } from '@/client/api/task'
import {
  IconCheckCircle,
  IconClose,
  IconError,
  IconHourglass,
} from '@/lib/components/icons'
import truncateMiddle from '@/lib/helpers/truncate-middle'
import { useAppSelector } from '@/store/hook'

export type TaskDrawerItemProps = {
  task: Task
}

const TaskDrawerItem = ({ task }: TaskDrawerItemProps) => {
  const [isDismissing, setIsDismissing] = useState(false)
  const mutateList = useAppSelector((state) => state.ui.tasks.mutateList)

  const handleDismiss = useCallback(async () => {
    try {
      setIsDismissing(true)
      await TaskAPI.dismiss(task.id)
      await mutateList?.(await TaskAPI.list())
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
              {task.payload?.object ? (
                <span className={cx('font-semibold')}>
                  {truncateMiddle(task.payload.object, 40)}
                </span>
              ) : null}
              {task.status !== Status.Error ? (
                <Text noOfLines={3}>{task.name}</Text>
              ) : null}
            </div>
            {task.error ? (
              <IconButton
                icon={<IconClose />}
                size="xs"
                variant="outline"
                colorScheme="gray"
                aria-label=""
                isLoading={isDismissing}
                onClick={handleDismiss}
              />
            ) : null}
          </div>
          {task.error ? (
            <Accordion allowMultiple>
              <AccordionItem className={cx('border-none')}>
                <AccordionButton className={cx('p-0.5')}>
                  <div className={cx('flex', 'flex-row', 'w-full')}>
                    <span className={cx('text-left', 'grow')}>
                      Task failed, click to show error.
                    </span>
                    <AccordionIcon />
                  </div>
                </AccordionButton>
                <AccordionPanel className={cx('p-0.5')}>
                  <Text className={cx('text-red-500')} noOfLines={5}>
                    {task.error}
                  </Text>
                </AccordionPanel>
              </AccordionItem>
            </Accordion>
          ) : null}
        </div>
      </CardBody>
    </Card>
  )
}

export default TaskDrawerItem
