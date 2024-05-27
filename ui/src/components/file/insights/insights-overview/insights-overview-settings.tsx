import { useCallback, useState } from 'react'
import { Button, Card, CardBody, CardFooter, Text } from '@chakra-ui/react'
import cx from 'classnames'
import InsightsAPI from '@/client/api/insights'
import { swrConfig } from '@/client/options'
import { IconDelete, IconSync } from '@/lib'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/insights'

const InsightsOverviewSettings = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateSummary = useAppSelector(
    (state) => state.ui.insights.mutateSummary,
  )
  const [isUpdating, setIsUpdating] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)
  const { data: summary } = InsightsAPI.useGetMetadata(id, swrConfig())

  const handleUpdate = useCallback(async () => {
    if (!id) {
      return
    }
    setIsUpdating(true)
    try {
      await InsightsAPI.patch(id)
      mutateSummary?.()
    } catch {
      setIsUpdating(false)
    } finally {
      setIsUpdating(false)
    }
  }, [id, mutateSummary])

  const handleDelete = useCallback(async () => {
    if (!id) {
      return
    }
    setIsDeleting(true)
    try {
      await InsightsAPI.delete(id)
      mutateSummary?.()
      dispatch(modalDidClose())
    } catch {
      setIsDeleting(false)
    } finally {
      setIsDeleting(false)
    }
  }, [id, mutateSummary, dispatch])

  if (!id || !summary) {
    return null
  }

  return (
    <div className={cx('flex', 'flex-row', 'items-stretch', 'gap-1.5')}>
      <Card size="md" variant="outline" className={cx('w-[50%]')}>
        <CardBody>
          <Text>
            Creates new insights for the active snapshot, uses the previously
            set language.
          </Text>
        </CardBody>
        <CardFooter>
          <Button
            leftIcon={<IconSync />}
            isLoading={isUpdating}
            isDisabled={!summary.isOutdated || isDeleting || isUpdating}
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
