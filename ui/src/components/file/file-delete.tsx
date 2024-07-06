// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.

import { useCallback, useState } from 'react'
import { useParams } from 'react-router-dom'
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
import FileAPI from '@/client/api/file'
import { useAppSelector } from '@/store/hook'
import { deleteModalDidClose, selectionUpdated } from '@/store/ui/files'

const FileDelete = () => {
  const { fileId } = useParams()
  const dispatch = useDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isDeleteModalOpen,
  )
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const [isLoading, setIsLoading] = useState(false)

  const handleDelete = useCallback(async () => {
    try {
      setIsLoading(true)
      await FileAPI.delete({ ids: selection })
      mutateList?.()
      dispatch(selectionUpdated([]))
      dispatch(deleteModalDidClose())
    } finally {
      setIsLoading(false)
    }
  }, [selection, fileId, dispatch, mutateList])

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => dispatch(deleteModalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        {selection.length > 1 ? (
          <ModalHeader>Delete {selection.length} Item(s)</ModalHeader>
        ) : (
          <ModalHeader>Delete Item</ModalHeader>
        )}
        <ModalCloseButton />
        <ModalBody>
          {selection.length > 1 ? (
            <span>
              Are you sure you would like to delete ({selection.length})
              item(s)?
            </span>
          ) : (
            <span>Are you sure you would like to delete this item?</span>
          )}
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              colorScheme="blue"
              disabled={isLoading}
              onClick={() => dispatch(deleteModalDidClose())}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="solid"
              colorScheme="red"
              isLoading={isLoading}
              onClick={handleDelete}
            >
              Delete
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileDelete
