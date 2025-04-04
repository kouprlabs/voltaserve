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
import { Button, useDisclosure } from '@chakra-ui/react'
import { AuxiliaryDrawer, SectionError, SectionSpinner } from '@koupr/ui'
import { IconClearAll, IconStacks } from '@koupr/ui'
import cx from 'classnames'
import { TaskAPI } from '@/client/api/task'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidClose, mutateCountUpdated } from '@/store/ui/tasks'
import TasksList from './task-list'

const TaskDrawer = () => {
  const dispatch = useAppDispatch()
  const { isOpen, onOpen, onClose } = useDisclosure()
  const [isDismissing, setIsDismissing] = useState(false)
  const isDrawerOpen = useAppSelector((state) => state.ui.tasks.isDrawerOpen)
  const mutateList = useAppSelector((state) => state.ui.tasks.mutateList)
  const {
    data: count,
    error: countError,
    isLoading: countIsLoading,
    mutate: mutateCount,
  } = TaskAPI.useGetCount(swrConfig())
  const countIsReady = count != null && !countError

  useEffect(() => {
    if (isDrawerOpen) {
      onOpen()
    } else {
      onClose()
    }
  }, [isDrawerOpen, onOpen, onClose])

  useEffect(() => {
    if (mutateCount) {
      dispatch(mutateCountUpdated(mutateCount))
    }
  }, [mutateCount, dispatch])

  const handleClearCompleted = useCallback(async () => {
    try {
      setIsDismissing(true)
      await TaskAPI.dismissAll()
      await mutateList?.(await TaskAPI.list())
    } finally {
      setIsDismissing(false)
    }
  }, [dispatch, mutateList])

  return (
    <AuxiliaryDrawer
      icon={<IconStacks />}
      isOpen={isOpen}
      onClose={() => {
        onClose()
        dispatch(drawerDidClose())
      }}
      onOpen={onOpen}
      hasBadge={count !== undefined && count > 0}
      header="Tasks"
      body={
        <>
          {countIsLoading ? <SectionSpinner /> : null}
          {countError ? (
            <SectionError text={errorToString(countError)} />
          ) : null}
          {countIsReady ? <TasksList /> : null}
        </>
      }
      footer={
        <>
          {countIsReady ? (
            <>
              {count && count > 0 ? (
                <Button
                  className={cx('w-full')}
                  size="sm"
                  leftIcon={<IconClearAll />}
                  isLoading={isDismissing}
                  onClick={handleClearCompleted}
                >
                  Clear Completed Items
                </Button>
              ) : null}
            </>
          ) : null}
        </>
      }
    />
  )
}

export default TaskDrawer
