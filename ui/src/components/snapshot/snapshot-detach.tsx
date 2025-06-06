// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useCallback, useState } from 'react'
import { useDispatch } from 'react-redux'
import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import cx from 'classnames'
import { SnapshotAPI } from '@/client/api/snapshot'
import { useAppSelector } from '@/store/hook'
import { detachModalDidClose, selectionUpdated } from '@/store/ui/snapshots'

const SnapshotDetach = () => {
  const dispatch = useDispatch()
  const id = useAppSelector((state) =>
    state.ui.snapshots.selection.length > 0
      ? state.ui.snapshots.selection[0]
      : undefined,
  )
  const mutate = useAppSelector((state) => state.ui.snapshots.snapshotMutate)
  const isModalOpen = useAppSelector(
    (state) => state.ui.snapshots.isDetachModalOpen,
  )
  const [isLoading, setIsLoading] = useState(false)

  const handleDetach = useCallback(async () => {
    if (id) {
      setIsLoading(true)
      try {
        await SnapshotAPI.detach(id)
        await mutate?.()
        dispatch(selectionUpdated([]))
        dispatch(detachModalDidClose())
      } catch {
        setIsLoading(false)
      } finally {
        setIsLoading(false)
      }
    }
  }, [id, dispatch, mutate])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(detachModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Detach Snapshot</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <span>Are you sure you want to detach this snapshot?</span>
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isLoading}
              onClick={() => dispatch(detachModalDidClose())}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="solid"
              colorScheme="red"
              isLoading={isLoading}
              onClick={handleDetach}
            >
              Detach
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default SnapshotDetach
