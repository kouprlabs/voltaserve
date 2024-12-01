// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useEffect } from 'react'
import { Modal, ModalCloseButton, ModalContent, ModalHeader, ModalOverlay } from '@chakra-ui/react'
import { SectionError, SectionSpinner } from '@koupr/ui'
import FileAPI from '@/client/api/file'
import MosaicAPI from '@/client/api/mosaic'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { mutateInfoUpdated } from '@/store/ui/mosaic'
import { modalDidClose } from '@/store/ui/mosaic'
import MosaicCreate from './mosaic-create'
import MosaicOverview from './mosaic-overview'

const Mosaic = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) => (state.ui.files.selection.length > 0 ? state.ui.files.selection[0] : undefined))
  const isModalOpen = useAppSelector((state) => state.ui.mosaic.isModalOpen)
  const {
    data: info,
    error: infoError,
    isLoading: infoIsLoading,
    mutate: mutateInfo,
  } = MosaicAPI.useGetInfo(id, swrConfig())
  const { data: file, error: fileError, isLoading: fileIsLoading } = FileAPI.useGet(id, swrConfig())
  const fileIsReady = file && !fileError
  const infoIsReady = info && !infoError

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
    <Modal size="xl" isOpen={isModalOpen} onClose={() => dispatch(modalDidClose())} closeOnOverlayClick={false}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Mosaic</ModalHeader>
        <ModalCloseButton />
        {fileIsLoading ? <SectionSpinner /> : null}
        {fileError ? <SectionError text={errorToString(fileError)} /> : null}
        {fileIsReady ? (
          <>
            {infoIsLoading ? <SectionSpinner /> : null}
            {infoError ? <SectionError text={errorToString(infoError)} /> : null}
            {infoIsReady ? <>{info.isAvailable ? <MosaicOverview /> : <MosaicCreate />}</> : null}
          </>
        ) : null}
      </ModalContent>
    </Modal>
  )
}

export default Mosaic
