import { useCallback } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { ltOwnerPermission } from '@/client/api/permission'
import TaskAPI from '@/client/api/task'
import WatermarkAPI from '@/client/api/watermark'
import { swrConfig } from '@/client/options'
import { IconDelete, IconSync } from '@/lib/components/icons'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/watermark'

const WatermarkOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTaskCount = useAppSelector((state) => state.ui.tasks.mutateCount)
  const { data: info, mutate: mutateInfo } = WatermarkAPI.useGetInfo(
    id,
    swrConfig(),
  )
  const { data: file, mutate: mutateFile } = FileAPI.useGet(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (id) {
      await WatermarkAPI.create(id)
      mutateFile(await FileAPI.get(id))
      mutateInfo(await WatermarkAPI.getInfo(id))
      mutateFiles?.()
      mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [id, mutateFile, mutateFiles, mutateTaskCount, mutateInfo, dispatch])

  const handleDelete = useCallback(async () => {
    if (id) {
      await WatermarkAPI.delete(id)
      mutateFile(await FileAPI.get(id))
      mutateInfo(await WatermarkAPI.getInfo(id))
      mutateFiles?.()
      mutateTaskCount?.(await TaskAPI.getCount())
      dispatch(modalDidClose())
    }
  }, [id, mutateFiles, mutateTaskCount, mutateInfo, dispatch])

  if (!file || !info) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Apply watermark on the active snapshot.</Text>
        </CardBody>
        <CardFooter>
          <Button
            leftIcon={<IconSync />}
            isDisabled={!info.isOutdated || file.snapshot?.taskId !== undefined}
            onClick={handleUpdate}
          >
            Apply
          </Button>
        </CardFooter>
      </Card>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Delete watermark from the active snapshot.</Text>
        </CardBody>
        <CardFooter>
          <Button
            colorScheme="red"
            leftIcon={<IconDelete />}
            isDisabled={
              !file ||
              file.snapshot?.taskId !== undefined ||
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

export default WatermarkOverviewSettings
