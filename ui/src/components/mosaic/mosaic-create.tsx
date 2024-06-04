import { useCallback } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import cx from 'classnames'
import MosaicAPI from '@/client/api/mosaic'
import TaskAPI from '@/client/api/task'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/mosaic'

const MosaicCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutateList)
  const mutateMetadata = useAppSelector((state) => state.ui.mosaic.mutateInfo)

  const handleCreate = useCallback(async () => {
    if (id) {
      await MosaicAPI.create(id, false)
      mutateMetadata?.()
      mutateFiles?.()
      mutateTasks?.(await TaskAPI.list())
      dispatch(modalDidClose())
    }
  }, [id, mutateFiles, mutateTasks, mutateMetadata, dispatch])

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
            The mosaic enhances view performance of large images by splitting
            them into smaller, manageable tiles. This makes browsing
            high-resolution images faster and more efficient.
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
            Enable
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default MosaicCreate
