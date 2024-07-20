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
import { useSWRConfig } from 'swr'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { useAppSelector } from '@/store/hook'
import { deleteModalDidClose, selectionUpdated } from '@/store/ui/files'
import { drawerDidOpen } from '@/store/ui/tasks'

const FileDelete = () => {
  const { mutate } = useSWRConfig()
  const { fileId } = useParams()
  const dispatch = useDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isDeleteModalOpen,
  )
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutateList)
  const [isLoading, setIsLoading] = useState(false)

  const handleDelete = useCallback(async () => {
    try {
      setIsLoading(true)

      // We intentionally mutate before we delete to avoid SWR
      // trying to fetch the file while the delete process is still ongoing
      await mutate(`/files/${fileId}`, null, false)

      FileAPI.delete({ ids: selection }).then(() => mutateList?.())
      await mutateTasks?.()
      dispatch(drawerDidOpen())
      dispatch(selectionUpdated([]))
      dispatch(deleteModalDidClose())
    } finally {
      setIsLoading(false)
    }
  }, [selection, fileId, dispatch, mutateList, mutateTasks])

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
