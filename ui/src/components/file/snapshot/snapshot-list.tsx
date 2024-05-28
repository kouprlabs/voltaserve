import { useCallback, useEffect, useState } from 'react'
import {
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
import SnapshotAPI, { Snapshot, SortOrder } from '@/client/api/snapshot'
import { swrConfig } from '@/client/options'
import prettyDate from '@/helpers/pretty-date'
import { Pagination, SectionSpinner } from '@/lib'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  detachModalDidOpen,
  listModalDidClose,
  mutateUpdated,
  selectionUpdated,
} from '@/store/ui/snapshots'

const SnapshotList = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.snapshots.isListModalOpen,
  )
  const fileId = useAppSelector((state) => state.ui.files.selection[0])
  const mutateFileList = useAppSelector((state) => state.ui.files.mutate)
  const [isActivating, setIsActivating] = useState(false)
  const [page, setPage] = useState(1)
  const [selected, setSelected] = useState<Snapshot>()
  const {
    data: list,
    error,
    mutate: snapshotMutate,
  } = SnapshotAPI.useList(
    { fileId, page, size: 5, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    if (snapshotMutate) {
      dispatch(mutateUpdated(snapshotMutate))
    }
  }, [snapshotMutate])

  const handleClose = useCallback(() => {
    dispatch(listModalDidClose())
    dispatch(selectionUpdated([]))
    setSelected(undefined)
  }, [dispatch])

  const handleActivate = useCallback(async () => {
    if (selected) {
      try {
        setIsActivating(true)
        await SnapshotAPI.activate(selected.id, { fileId })
        snapshotMutate()
        mutateFileList?.()
      } finally {
        setIsActivating(false)
      }
    }
  }, [selected, dispatch, snapshotMutate, mutateFileList])

  const handleDetach = useCallback(() => {
    if (selected) {
      dispatch(selectionUpdated([selected.id]))
      dispatch(detachModalDidOpen())
    }
  }, [selected, dispatch])

  const handleSelect = useCallback((snapshot: Snapshot) => {
    setSelected(snapshot)
    dispatch(selectionUpdated([snapshot.id]))
  }, [])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={handleClose}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Snapshots</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <div className={cx('flex', 'flex-col', 'gap-1.5')}>
            {!list && error ? (
              <div
                className={cx(
                  'flex',
                  'items-center',
                  'justify-center',
                  'h-[300px]',
                )}
              >
                <span>Failed to load snapshots.</span>
              </div>
            ) : null}
            {!list && !error ? <SectionSpinner /> : null}
            {list && list.data.length === 0 ? (
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
            ) : null}
            {list && list.data.length > 0 ? (
              <Table variant="simple" size="sm">
                <colgroup>
                  <col className={cx('w-[40px]')} />
                  <col className={cx('w-[auto]')} />
                </colgroup>
                <Tbody>
                  {list.data.map((s) => (
                    <Tr
                      key={s.id}
                      className={cx(
                        'cursor-pointer',
                        'h-[52px]',
                        { 'bg-gray-100': selected?.id === s.id },
                        { 'dark:bg-gray-600': selected?.id === s.id },
                        { 'bg-transparent': selected?.id !== s.id },
                      )}
                      onClick={() => handleSelect(s)}
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
            ) : null}
            {list ? (
              <div className={cx('self-end')}>
                {list.totalPages > 1 ? (
                  <Pagination
                    uiSize="md"
                    maxButtons={3}
                    page={page}
                    totalPages={list.totalPages}
                    onPageChange={(value) => setPage(value)}
                  />
                ) : null}
              </div>
            ) : null}
          </div>
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              isDisabled={isActivating}
              onClick={handleClose}
            >
              Close
            </Button>
            <Button
              variant="outline"
              colorScheme="red"
              isDisabled={!selected || selected.isActive || isActivating}
              onClick={handleDetach}
            >
              Detach Snapshot
            </Button>
            <Button
              variant="solid"
              colorScheme="blue"
              isLoading={isActivating}
              isDisabled={!selected || selected.isActive || isActivating}
              onClick={handleActivate}
            >
              Activate Snapshot
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default SnapshotList
