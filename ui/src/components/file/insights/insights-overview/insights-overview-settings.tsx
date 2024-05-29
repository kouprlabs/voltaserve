import { useCallback } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import InsightsAPI from '@/client/api/insights'
import { swrConfig } from '@/client/options'
import { IconDelete, IconSync } from '@/lib/components/icons'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  deletingDidStart,
  deletingDidStop,
  modalDidClose,
  updatingDidStop,
  updatingDidStart,
} from '@/store/ui/insights'

const InsightsOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateMetadata = useAppSelector(
    (state) => state.ui.insights.mutateMetadata,
  )
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const isUpdating = useAppSelector((state) => state.ui.insights.isUpdating)
  const isDeleting = useAppSelector((state) => state.ui.insights.isDeleting)
  const { data: metadata } = InsightsAPI.useGetMetadata(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (!id) {
      return
    }
    try {
      dispatch(updatingDidStart())
      await InsightsAPI.patch(id)
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
      await InsightsAPI.delete(id)
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
          <Text>
            Updates the insights using the active snapshot, uses the previously
            set language.
          </Text>
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
            Deletes insights from the active snapshot, can be recreated later.
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

export default InsightsOverviewSettings
