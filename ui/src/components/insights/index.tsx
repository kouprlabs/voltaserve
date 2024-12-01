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
import InsightsAPI from '@/client/api/insights'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose, mutateInfoUpdated } from '@/store/ui/insights'
import InsightsCreate from './insights-create'
import InsightsOverview from './insights-overview'

const Insights = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) => (state.ui.files.selection.length > 0 ? state.ui.files.selection[0] : undefined))
  const isModalOpen = useAppSelector((state) => state.ui.insights.isModalOpen)
  const {
    data: info,
    error: infoError,
    isLoading: infoIsLoading,
    mutate: mutateInfo,
  } = InsightsAPI.useGetInfo(id, swrConfig())
  const { data: file, error: fileError, isLoading: fileIsLoading } = FileAPI.useGet(id, swrConfig())
  const infoIsReady = info && !infoError
  const fileIsReady = file && !fileError

  useEffect(() => {
    if (mutateInfo) {
      dispatch(mutateInfoUpdated(mutateInfo))
    }
  }, [mutateInfo])

  useEffect(() => {
    if (file?.snapshot?.task?.isPending) {
      dispatch(modalDidClose())
    }
  }, [file])

  return (
    <Modal size="xl" isOpen={isModalOpen} onClose={() => dispatch(modalDidClose())} closeOnOverlayClick={false}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Insights</ModalHeader>
        <ModalCloseButton />
        {fileIsLoading ? <SectionSpinner /> : null}
        {fileError ? <SectionError text={errorToString(fileError)} /> : null}
        {fileIsReady ? (
          <>
            {infoIsLoading ? <SectionSpinner /> : null}
            {infoError ? <SectionError text={errorToString(infoError)} /> : null}
            {infoIsReady ? <>{info?.isAvailable ? <InsightsOverview /> : <InsightsCreate />}</> : null}
          </>
        ) : null}
      </ModalContent>
    </Modal>
  )
}

export default Insights
