import { useCallback, useState } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import MosaicAPI from '@/client/api/mosaic'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/mosaic'

const PerformanceCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateMetadata = useAppSelector(
    (state) => state.ui.mosaic.mutateMetadata,
  )
  const [isLoading, setIsLoading] = useState(false)
  const metadata = { isOutdated: true }
  const { data: file } = FileAPI.useGet(id, swrConfig())

  const handleCreate = useCallback(async () => {
    if (id) {
      try {
        setIsLoading(true)
        await MosaicAPI.create(id)
        mutateMetadata?.()
        setIsLoading(false)
      } catch (error) {
        setIsLoading(false)
      } finally {
        setIsLoading(false)
      }
    }
  }, [id, mutateMetadata])

  if (!id || !file || !metadata) {
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
            Optimize your image for better performance by creating a mosaic.
            <br />
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
            isDisabled={isLoading}
            onClick={() => dispatch(modalDidClose())}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="solid"
            colorScheme="blue"
            isLoading={isLoading}
            onClick={handleCreate}
          >
            Optimize Image
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default PerformanceCreate
