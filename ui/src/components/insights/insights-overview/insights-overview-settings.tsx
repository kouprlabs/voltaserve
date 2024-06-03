import { useCallback } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import InsightsAPI from '@/client/api/insights'
import { ltOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import { IconDelete, IconSync } from '@/lib/components/icons'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/insights'
import { drawerDidOpen as tasksDrawerDidOpen } from '@/store/ui/tasks'

const InsightsOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutate)
  const { data: metadata } = InsightsAPI.useGetMetadata(id, swrConfig())
  const { data: file } = FileAPI.useGet(id, swrConfig())

  const handleUpdate = useCallback(() => {
    if (id) {
      InsightsAPI.patch(id)
      mutateFiles?.()
      mutateTasks?.()
    }
  }, [id, mutateFiles, mutateTasks, dispatch])

  const handleDelete = useCallback(() => {
    if (id) {
      InsightsAPI.delete(id)
      mutateFiles?.()
      mutateTasks?.()
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
          <Text>
            Updates the insights using the active snapshot, uses the previously
            set language.
          </Text>
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
            Deletes insights from the active snapshot, can be recreated later.
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

export default InsightsOverviewSettings
