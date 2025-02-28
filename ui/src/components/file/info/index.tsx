// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
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
import { FileAPI } from '@/client/api/file'
import FileInfoEmbed from '@/components/file/info/file-info-embed'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { infoModalDidClose } from '@/store/ui/files'

const FileInfo = () => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector((state) => state.ui.files.isInfoModalOpen)
  const id = useAppSelector((state) => state.ui.files.selection[0])
  const { data: file } = FileAPI.useGet(id)

  if (!file) {
    return null
  }

  return (
    <Modal isOpen={isModalOpen} onClose={() => dispatch(infoModalDidClose())}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Info</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <FileInfoEmbed file={file} />
        </ModalBody>
        <ModalFooter>
          <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
            <Button
              type="button"
              variant="outline"
              onClick={() => dispatch(infoModalDidClose())}
            >
              Close
            </Button>
          </div>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default FileInfo
