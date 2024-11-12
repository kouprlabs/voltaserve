// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect } from 'react'
import {
  Modal,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import { SectionError, SectionSpinner } from '@koupr/ui'
import FileAPI from '@/client/api/file'
import MosaicAPI from '@/client/api/mosaic'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { mutateInfoUpdated } from '@/store/ui/mosaic'
import { modalDidClose } from '@/store/ui/mosaic'
import MosaicCreate from './mosaic-create'
import MosaicOverview from './mosaic-overview'

const Mosaic = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const isModalOpen = useAppSelector((state) => state.ui.mosaic.isModalOpen)
  const {
    data: info,
    error: infoError,
    mutate: mutateInfo,
  } = MosaicAPI.useGetInfo(id, swrConfig())
  const { data: file, error: fileError } = FileAPI.useGet(id, swrConfig())
  const isFileLoading = !file && !fileError
  const isFileError = !file && fileError
  const isFileReady = file && !fileError
  const isInfoLoading = !info && !infoError
  const isInfoError = !info && infoError
  const isInfoReady = info && !infoError

  useEffect(() => {
    if (file?.snapshot?.task?.isPending) {
      dispatch(modalDidClose())
    }
  }, [file])

  useEffect(() => {
    if (mutateInfo) {
      dispatch(mutateInfoUpdated(mutateInfo))
    }
  }, [mutateInfo])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => dispatch(modalDidClose())}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Mosaic</ModalHeader>
        <ModalCloseButton />
        {isFileLoading ? <SectionSpinner /> : null}
        {isFileError ? <SectionError text="Failed to load file." /> : null}
        {isFileReady ? (
          <>
            {isInfoLoading ? <SectionSpinner /> : null}
            {isInfoError ? <SectionError text="Failed to load info." /> : null}
            {isInfoReady ? (
              <>{info.isAvailable ? <MosaicOverview /> : <MosaicCreate />}</>
            ) : null}
          </>
        ) : null}
      </ModalContent>
    </Modal>
  )
}

export default Mosaic
