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
import cx from 'classnames'
import MosaicAPI from '@/client/api/mosaic'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'
import MosaicOverviewFile from './mosaic-overview-file'
import MosaicOverviewSettings from './mosaic-overview-settings'

const MosaicOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [isWarningVisible, setIsWarningVisible] = useState(true)
  const { data: metadata } = MosaicAPI.useGetMetadata(id, swrConfig())

  return (
    <>
      <ModalBody>
        <div className={cx('flex', 'flex-col', 'gap-1.5', 'w-full')}>
          {metadata?.isOutdated && isWarningVisible ? (
            <Alert status="warning">
              <AlertIcon />
              <Box>
                <AlertDescription>
                  This mosaic is outdated, it originates from an older snapshot.
                  Please navigate to the settings tab to update.
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
              <Tab>File</Tab>
              <Tab>Settings</Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <MosaicOverviewFile />
              </TabPanel>
              <TabPanel>
                <MosaicOverviewSettings />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </div>
      </ModalBody>
    </>
  )
}

export default MosaicOverview
