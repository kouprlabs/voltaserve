import { useCallback } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import cx from 'classnames'
import TaskAPI from '@/client/api/task'
import WatermarkAPI from '@/client/api/watermark'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { drawerDidOpen as tasksDrawerDidOpen } from '@/store/ui/tasks'
import { modalDidClose } from '@/store/ui/watermark'

const WatermarkCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutate)

  const handleCreate = useCallback(async () => {
    if (id) {
      WatermarkAPI.create(id, false)
      const tasks = await TaskAPI.list()
      mutateFiles?.()
      mutateTasks?.(tasks)
      dispatch(modalDidClose())
      dispatch(tasksDrawerDidOpen())
    }
  }, [id, mutateFiles, mutateTasks, dispatch])

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
            Protect your document or image with a watermark to enhance its
            security by clearly marking it as confidential or proprietary, thus
            deterring unauthorized use or distribution.
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
            Create Watermark
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default WatermarkCreate
