import { useMemo } from 'react'
import { useParams } from 'react-router-dom'
import {
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import { File } from '@/client/api/file'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { snapshotsModalDidClose } from '@/store/ui/files'

export type FileSharingProps = {
  file: File
}

const FileSnapshots = ({ file }: FileSharingProps) => {
  const dispatch = useAppDispatch()
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isSnapshotModalOpen,
  )
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
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
          eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad
          minim veniam, quis nostrud exercitation ullamco laboris nisi ut
          aliquip ex ea commodo consequat.
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}

export default FileSnapshots
