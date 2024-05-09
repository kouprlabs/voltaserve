import { useCallback, useMemo, useState } from 'react'
import {
  Avatar,
  Badge,
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Radio,
  Table,
  Tbody,
  Td,
  Tr,
} from '@chakra-ui/react'
import cx from 'classnames'
import { File } from '@/client/api/file'
import { Snapshot } from '@/client/api/snapshot'
import prettyDate from '@/helpers/pretty-date'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { snapshotsModalDidClose } from '@/store/ui/files'

export type FileSharingProps = {
  file: File
  onConfirm?: (snapshot: Snapshot) => void
}

const FileSnapshots = ({ file, onConfirm }: FileSharingProps) => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isSnapshotModalOpen,
  )
  const [selected, setSelected] = useState<Snapshot>()
  const list = file.snapshots

  const handleClose = useCallback(() => {
    dispatch(snapshotsModalDidClose())
  }, [dispatch])

  const handleConfirm = useCallback(() => {
    if (selected) {
      onConfirm?.(selected)
      handleClose()
    }
  }, [selected, onConfirm, handleClose])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => {
        dispatch(snapshotsModalDidClose())
      }}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Snapshots</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          {list && list.length === 0 && (
            <div
              className={cx(
                'flex',
                'items-center',
                'justify-center',
                'h-[300px]',
              )}
            >
              <div
                className={cx('flex', 'flex-col', 'items-center', 'gap-1.5')}
              >
                <span>There are no snapshots.</span>
              </div>
            </div>
          )}
          {list && list.length > 0 && (
            <Table variant="simple" size="sm">
              <colgroup>
                <col className={cx('w-[40px]')} />
                <col className={cx('w-[auto]')} />
              </colgroup>
              <Tbody>
                {list.map((s) => (
                  <Tr
                    key={s.id}
                    className={cx(
                      'cursor-pointer',
                      'h-[52px]',
                      { 'bg-gray-100': selected?.id === s.id },
                      { 'dark:bg-gray-600': selected?.id === s.id },
                      { 'bg-transparent': selected?.id !== s.id },
                    )}
                    onClick={() => setSelected(s)}
                  >
                    <Td className={cx('px-0.5', 'text-center')}>
                      <Radio size="md" isChecked={selected?.id === s.id} />
                    </Td>
                    <Td className={cx('px-0.5')}>
                      <div
                        className={cx(
                          'flex',
                          'flex-row',
                          'items-center',
                          'gap-1.5',
                        )}
                      >
                        <Avatar
                          name={`V ${s.version}`}
                          size="sm"
                          className={cx('w-[40px]', 'h-[40px]')}
                        />
                        <span className={cx('text-base')}>
                          {prettyDate(s.createTime)}
                        </span>
                        {s.isActive ? (
                          <Badge colorScheme="green">Active</Badge>
                        ) : null}
                      </div>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          )}
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              onClick={handleClose}
            >
              Cancel
            </Button>
            <Button
              variant="solid"
              colorScheme="blue"
              isDisabled={!selected}
              onClick={handleConfirm}
            >
              Activate
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileSnapshots
