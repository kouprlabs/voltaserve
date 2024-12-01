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
import { useParams } from 'react-router-dom'
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
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { copyModalDidClose, loadingAdded, loadingRemoved, selectionUpdated } from '@/store/ui/files'
import FileBrowse from './file-browse'

const FileCopy = () => {
  const dispatch = useAppDispatch()
  const { fileId } = useParams()
  const isModalOpen = useAppSelector((state) => state.ui.files.isCopyModalOpen)
  const selection = useAppSelector((state) => state.ui.files.selection)
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const mutateTasks = useAppSelector((state) => state.ui.tasks.mutateList)
  const [targetId, setTargetId] = useState<string>()

  const handleCopy = useCallback(async () => {
    if (!targetId) {
      return
    }
    const ids = [...selection]
    for (const id of ids) {
      dispatch(loadingAdded([id]))
      FileAPI.copyOne(id, targetId)
        .then(() => {
          if (fileId === targetId) {
            mutateList?.()
          }
        })
        .finally(() => dispatch(loadingRemoved([id])))
    }
    await mutateTasks?.()
    dispatch(selectionUpdated([]))
    dispatch(copyModalDidClose())
  }, [targetId, fileId, selection, dispatch, mutateList, mutateTasks])

  return (
    <Modal isOpen={isModalOpen} onClose={() => dispatch(copyModalDidClose())} closeOnOverlayClick={false}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>{selection.length > 1 ? `Copy (${selection.length}) Items to…` : 'Copy Item to…'}</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <FileBrowse onChange={(id) => setTargetId(id)} />
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button type="button" variant="outline" colorScheme="blue" onClick={() => dispatch(copyModalDidClose())}>
              Cancel
            </Button>
            <Button variant="solid" colorScheme="blue" onClick={handleCopy}>
              Copy Here
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileCopy
