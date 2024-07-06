// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useCallback } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import InsightsAPI from '@/client/api/insights'
import { ltOwnerPermission } from '@/client/api/permission'
import TaskAPI from '@/client/api/task'
import { swrConfig } from '@/client/options'
import { IconDelete, IconSync } from '@/lib/components/icons'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/insights'

const InsightsOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTaskCount = useAppSelector((state) => state.ui.tasks.mutateCount)
  const { data: info, mutate: mutateInfo } = InsightsAPI.useGetInfo(
    id,
    swrConfig(),
  )
  const { data: file, mutate: mutateFile } = FileAPI.useGet(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (id) {
      await InsightsAPI.patch(id)
      mutateFile(await FileAPI.get(id))
      mutateInfo(await InsightsAPI.getInfo(id))
      mutateFiles?.()
      mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [id, mutateFiles, mutateTaskCount, mutateInfo, dispatch])

  const handleDelete = useCallback(async () => {
    if (id) {
      await InsightsAPI.delete(id)
      mutateFile(await FileAPI.get(id))
      mutateInfo(await InsightsAPI.getInfo(id))
      mutateFiles?.()
      mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [id, mutateFile, mutateFiles, mutateTaskCount, mutateInfo, dispatch])

  if (!file || !info) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Collect insights for the active snapshot.</Text>
        </CardBody>
        <CardFooter>
          <Button
            leftIcon={<IconSync />}
            isDisabled={!info.isOutdated || file.snapshot?.task?.isPending}
            onClick={handleUpdate}
          >
            Collect
          </Button>
        </CardFooter>
      </Card>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Delete insights from the active snapshot.</Text>
        </CardBody>
        <CardFooter>
          <Button
            colorScheme="red"
            leftIcon={<IconDelete />}
            isDisabled={
              !file ||
              file.snapshot?.task?.isPending ||
              info.isOutdated ||
              ltOwnerPermission(file.permission)
            }
            onClick={handleDelete}
          >
            Delete
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default InsightsOverviewSettings
