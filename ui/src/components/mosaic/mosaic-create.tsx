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
  const mutateInfo = useAppSelector((state) => state.ui.mosaic.mutateInfo)

  const handleCreate = useCallback(async () => {
    if (id) {
      await MosaicAPI.create(id, false)
      mutateInfo?.(await MosaicAPI.getInfo(id))
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
            Create a mosaic to enhance view performance of a large image by
            splitting it into smaller, manageable tiles. This makes browsing a
            high-resolution image faster and more efficient.
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
            Create
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default MosaicCreate
