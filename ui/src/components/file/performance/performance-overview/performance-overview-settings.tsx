import { useCallback } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import MosaicAPI from '@/client/api/mosaic'
import { swrConfig } from '@/client/options'
import { IconDelete, IconSync } from '@/lib/components/icons'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { updatingDidStart } from '@/store/ui/insights'
import {
  deletingDidStart,
  deletingDidStop,
  modalDidClose,
  updatingDidStop,
} from '@/store/ui/mosaic'

const PeformanceOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateMetadata = useAppSelector(
    (state) => state.ui.mosaic.mutateMetadata,
  )
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const isUpdating = useAppSelector((state) => state.ui.mosaic.isUpdating)
  const isDeleting = useAppSelector((state) => state.ui.mosaic.isDeleting)
  const { data: metadata } = MosaicAPI.useGetMetadata(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (!id) {
      return
    }
    try {
      dispatch(updatingDidStart())
      await MosaicAPI.create(id)
      mutateMetadata?.()
      mutateList?.()
    } catch {
      dispatch(updatingDidStop())
    } finally {
      dispatch(updatingDidStop())
    }
  }, [id, mutateMetadata, mutateList, dispatch])

  const handleDelete = useCallback(async () => {
    if (!id) {
      return
    }
    try {
      dispatch(deletingDidStart())
      await MosaicAPI.delete(id)
      mutateMetadata?.()
      mutateList?.()
      dispatch(modalDidClose())
    } catch {
      dispatch(deletingDidStop())
    } finally {
      dispatch(deletingDidStop())
    }
  }, [id, mutateMetadata, mutateList, dispatch])

  if (!id || !metadata) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>Updates to a new mosaic using the active snapshot.</Text>
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
            Deletes the mosaic from the active snapshot, can be recreated later.
          </Text>
        </CardBody>
        <CardFooter>
          <Button
            colorScheme="red"
            leftIcon={<IconDelete />}
            isLoading={isDeleting}
            isDisabled={isDeleting || isUpdating}
            onClick={handleDelete}
          >
            Delete
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}

export default PeformanceOverviewSettings
