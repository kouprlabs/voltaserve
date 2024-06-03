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
import { drawerDidOpen as tasksDrawerDidOpen } from '@/store/ui/tasks'
import { modalDidClose } from '@/store/ui/watermark'

const WatermarkOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutate)
  const { data: metadata } = WatermarkAPI.useGetMetadata(id, swrConfig())
  const { data: file } = FileAPI.useGet(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (id) {
      WatermarkAPI.create(id)
      const tasks = await TaskAPI.list()
      mutateFiles?.()
      mutateTasks?.(tasks)
      dispatch(modalDidClose())
      dispatch(tasksDrawerDidOpen())
    }
  }, [id, mutateFiles, mutateTasks, dispatch])

  const handleDelete = useCallback(async () => {
    if (id) {
      WatermarkAPI.delete(id)
      const tasks = await TaskAPI.list()
      mutateFiles?.()
      mutateTasks?.(tasks)
      dispatch(modalDidClose())
      dispatch(tasksDrawerDidOpen())
    }
  }, [id, mutateFiles, mutateTasks, dispatch])

  if (!id || !metadata) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Updates the watermark using the active snapshot.</Text>
        </CardBody>
        <CardFooter>
          <Button
            leftIcon={<IconSync />}
            isDisabled={!metadata.isOutdated}
            onClick={handleUpdate}
          >
            Update
          </Button>
        </CardFooter>
      </Card>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>
            Deletes the watermark from the active snapshot, can be recreated
            later.
          </Text>
        </CardBody>
        <CardFooter>
          <Button
            colorScheme="red"
            leftIcon={<IconDelete />}
            isDisabled={!file || ltOwnerPermission(file.permission)}
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
