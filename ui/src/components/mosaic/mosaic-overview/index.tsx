// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { useState } from 'react'
import {
  Alert,
  AlertDescription,
  AlertIcon,
  Box,
  CloseButton,
  ModalBody,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from '@chakra-ui/react'
import { SectionError, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import MosaicAPI from '@/client/api/mosaic'
import { errorToString } from '@/client/error'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'
import MosaicOverviewArtifacts from './mosaic-overview-artifacts'
import MosaicOverviewSettings from './mosaic-overview-settings'

const MosaicOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [isWarningVisible, setIsWarningVisible] = useState(true)
  const {
    data: info,
    error: infoError,
    isLoading: infoIsLoading,
  } = MosaicAPI.useGetInfo(id, swrConfig())
  const infoIsReady = info && !infoError

  return (
    <>
      <ModalBody>
        {infoIsLoading ? <SectionSpinner /> : null}
        {infoError ? <SectionError text={errorToString(infoError)} /> : null}
        {infoIsReady ? (
          <div className={cx('flex', 'flex-col', 'gap-1.5', 'w-full')}>
            {info.isOutdated && isWarningVisible ? (
              <Alert status="warning" className={cx('flex')}>
                <AlertIcon />
                <Box className={cx('grow')}>
                  <AlertDescription>
                    This mosaic comes from an older snapshot. You can create a
                    new one for the active snapshot from the settings.
                  </AlertDescription>
                </Box>
                <CloseButton
                  alignSelf="flex-start"
                  position="relative"
                  right={-1}
                  top={-1}
                  onClick={() => setIsWarningVisible(false)}
                />
              </Alert>
            ) : null}
            <Tabs colorScheme="gray">
              <TabList>
                <Tab>Artifacts</Tab>
                <Tab>Settings</Tab>
              </TabList>
              <TabPanels>
                <TabPanel>
                  <MosaicOverviewArtifacts />
                </TabPanel>
                <TabPanel>
                  <MosaicOverviewSettings />
                </TabPanel>
              </TabPanels>
            </Tabs>
          </div>
        ) : null}
      </ModalBody>
    </>
  )
}

export default MosaicOverview
