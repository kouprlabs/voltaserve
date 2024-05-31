import { useCallback } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { ltOwnerPermission } from '@/client/api/permission'
import WatermarkAPI from '@/client/api/watermark'
import { swrConfig } from '@/client/options'
import { IconDelete, IconSync } from '@/lib/components/icons'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  updatingDidStart,
  deletingDidStart,
  deletingDidStop,
  modalDidClose,
  updatingDidStop,
} from '@/store/ui/watermark'

const WatermarkOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFile = useAppSelector((state) => state.ui.watermark.mutateFile)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const isUpdating = useAppSelector((state) => state.ui.watermark.isUpdating)
  const isDeleting = useAppSelector((state) => state.ui.watermark.isDeleting)
  const { data: metadata } = WatermarkAPI.useGetMetadata(id, swrConfig())
  const { data: file } = FileAPI.useGet(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (!id) {
      return
    }
    try {
      dispatch(updatingDidStart())
      await WatermarkAPI.create(id)
      mutateFile?.()
      mutateList?.()
    } catch {
      dispatch(updatingDidStop())
    } finally {
      dispatch(updatingDidStop())
    }
  }, [id, mutateFile, mutateList, dispatch])

  const handleDelete = useCallback(async () => {
    if (!id) {
      return
    }
    try {
      dispatch(deletingDidStart())
      await WatermarkAPI.delete(id)
      mutateFile?.()
      mutateList?.()
      dispatch(modalDidClose())
    } catch {
      dispatch(deletingDidStop())
    } finally {
      dispatch(deletingDidStop())
    }
  }, [id, mutateFile, mutateList, dispatch])

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
            isLoading={isUpdating}
            isDisabled={!metadata.isOutdated || isDeleting || isUpdating}
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
            isLoading={isDeleting}
            isDisabled={
              isDeleting ||
              isUpdating ||
              !file ||
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
