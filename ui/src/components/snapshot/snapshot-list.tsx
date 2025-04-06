// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useEffect, useState } from 'react'
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
import {
  Pagination,
  RelativeDate,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  usePageMonitor,
} from '@koupr/ui'
import cx from 'classnames'
import {
  SnapshotAPI,
  Snapshot,
  SnapshotSortBy,
  SnapshotSortOrder,
} from '@/client/api/snapshot'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import prettyBytes from '@/lib/helpers/pretty-bytes'
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
  const mutateFiles = useAppSelector((state) => state.ui.files.mutate)
  const [isActivating, setIsActivating] = useState(false)
  const [page, setPage] = useState(1)
  const [selected, setSelected] = useState<Snapshot>()
  const size = 5
  const {
    data: list,
    error: listError,
    isLoading: listIsLoading,
    mutate: snapshotMutate,
  } = SnapshotAPI.useList(
    {
      fileId,
      page,
      size,
      sortBy: SnapshotSortBy.Version,
      sortOrder: SnapshotSortOrder.Desc,
    },
    swrConfig(),
  )
  const { hasPageSwitcher } = usePageMonitor({
    totalPages: list?.totalPages ?? 1,
    totalElements: list?.totalElements ?? 0,
    steps: [size],
  })
  const listIsEmpty = list && !listError && list.totalElements === 0
  const listIsReady = list && !listError && list.totalElements > 0

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
        const file = await SnapshotAPI.activate(selected.id)
        if (file.snapshot) {
          handleSelect(file.snapshot)
        }
        await snapshotMutate()
        mutateFiles?.()
      } finally {
        setIsActivating(false)
      }
    }
  }, [selected, dispatch, snapshotMutate, mutateFiles])

  const handleDetach = useCallback(() => {
    if (selected) {
      dispatch(selectionUpdated([selected.id]))
      dispatch(detachModalDidOpen())
      setSelected(undefined)
    }
  }, [selected, dispatch])

  const handleSelect = useCallback((snapshot: Snapshot) => {
    setSelected(snapshot)
    dispatch(selectionUpdated([snapshot.id]))
  }, [])

  const isSelected = useCallback(
    (snapshot: Snapshot) => {
      return selected?.id === snapshot.id || (snapshot.isActive && !selected)
    },
    [selected],
  )

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
          {listIsLoading ? <SectionSpinner /> : null}
          {listError ? <SectionError text={errorToString(listError)} /> : null}
          {listIsEmpty ? (
            <SectionPlaceholder text="There are no items." />
          ) : null}
          {listIsReady ? (
            <div className={cx('flex', 'flex-col', 'gap-1.5')}>
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
                        { 'bg-gray-100': isSelected(s) },
                        { 'dark:bg-gray-600': isSelected(s) },
                        { 'bg-transparent': !isSelected(s) },
                      )}
                      onClick={() => handleSelect(s)}
                    >
                      <Td className={cx('px-0.5', 'text-center')}>
                        <Radio size="md" isChecked={isSelected(s)} />
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
                          <div className={cx('flex', 'flex-col', 'gap-0.5')}>
                            <span className={cx('text-base')}>
                              <RelativeDate date={new Date(s.createTime)} />
                            </span>
                            <div className={cx('flex', 'flex-row', 'gap-0.5')}>
                              <Badge variant="outline">
                                Version {s.version}
                              </Badge>
                              {s.original.size ? (
                                <Badge variant="outline">
                                  {prettyBytes(s.original.size)}
                                </Badge>
                              ) : null}
                              {s.capabilities.summary ? (
                                <Badge variant="outline">Summary</Badge>
                              ) : null}
                              {s.capabilities.text ? (
                                <Badge variant="outline">Text</Badge>
                              ) : null}
                              {s.capabilities.ocr ? (
                                <Badge variant="outline">OCR</Badge>
                              ) : null}
                              {s.capabilities.entities ? (
                                <Badge variant="outline">Entities</Badge>
                              ) : null}
                              {s.capabilities.mosaic ? (
                                <Badge variant="outline">Mosaic</Badge>
                              ) : null}
                              {s.isActive ? (
                                <Badge colorScheme="green">Active</Badge>
                              ) : null}
                            </div>
                          </div>
                        </div>
                      </Td>
                    </Tr>
                  ))}
                </Tbody>
              </Table>
              {hasPageSwitcher ? (
                <div className={cx('self-end')}>
                  <Pagination
                    maxButtons={3}
                    page={page}
                    totalPages={list.totalPages}
                    onPageChange={(value) => setPage(value)}
                  />
                </div>
              ) : null}
            </div>
          ) : null}
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
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
              Detach
            </Button>
            <Button
              variant="solid"
              colorScheme="blue"
              isLoading={isActivating}
              isDisabled={!selected || selected.isActive || isActivating}
              onClick={handleActivate}
            >
              Activate
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default SnapshotList
