import { useCallback } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import cx from 'classnames'
import TaskAPI from '@/client/api/task'
import WatermarkAPI from '@/client/api/watermark'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/watermark'

const WatermarkCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutateList)
  const mutateInfo = useAppSelector((state) => state.ui.watermark.mutateInfo)

  const handleCreate = useCallback(async () => {
    if (id) {
      await WatermarkAPI.create(id, false)
      mutateInfo?.()
      mutateFiles?.()
      mutateTasks?.(await TaskAPI.list())
      dispatch(modalDidClose())
    }
  }, [id, mutateFiles, mutateTasks, mutateInfo, dispatch])

  if (!id) {
    return null
  }

  return (
    <>
      <ModalBody>
        <div
          className={cx(
            'flex',
            'flex-col',
            'items-center',
            'justify-center',
            'gap-1.5',
          )}
        >
          <p>
            Apply a watermark on your file to label it as confidential or
            proprietary, thus deterring unauthorized use or distribution. The
            watermark will be visible by default to users with view only
            permission.
          </p>
        </div>
      </ModalBody>
      <ModalFooter>
        <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            onClick={() => dispatch(modalDidClose())}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="solid"
            colorScheme="blue"
            onClick={handleCreate}
          >
            Apply
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default WatermarkCreate
