import { useCallback } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import cx from 'classnames'
import WatermarkAPI from '@/client/api/watermark'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  creatingDidStop,
  modalDidClose,
  creatingDidStart,
} from '@/store/ui/watermark'

const SecurityCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateFile = useAppSelector((state) => state.ui.watermark.mutateFile)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const isCreating = useAppSelector((state) => state.ui.watermark.isCreating)

  const handleCreate = useCallback(async () => {
    if (id) {
      try {
        dispatch(creatingDidStart())
        await WatermarkAPI.create(id)
        mutateFile?.()
        mutateList?.()
      } catch (error) {
        dispatch(creatingDidStop())
      } finally {
        dispatch(creatingDidStop())
      }
    }
  }, [id, mutateFile, mutateList, dispatch])

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
            Create a watermark for your document or image to enhance its
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
            isDisabled={isCreating}
            onClick={() => dispatch(modalDidClose())}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="solid"
            colorScheme="blue"
            isLoading={isCreating}
            onClick={handleCreate}
          >
            Protect With Watermark
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default SecurityCreate
