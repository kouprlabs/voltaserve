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
import WatermarkAPI from '@/client/api/watermark'
import { swrConfig } from '@/client/options'
import { useAppSelector } from '@/store/hook'
import WatermarkOverviewFile from './watermark-overview-file'
import WatermarkOverviewSettings from './watermark-overview-settings'

const WatermarkOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [isWarningVisible, setIsWarningVisible] = useState(true)
  const { data: metadata } = WatermarkAPI.useGetMetadata(id, swrConfig())

  return (
    <>
      <ModalBody>
        <div className={cx('flex', 'flex-col', 'gap-1.5', 'w-full')}>
          {metadata?.isOutdated && isWarningVisible ? (
            <Alert status="warning">
              <AlertIcon />
              <Box>
                <AlertDescription>
                  This watermark protected file is outdated, it originates from
                  an older snapshot. Please navigate to the settings tab to
                  update.
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
                <WatermarkOverviewFile />
              </TabPanel>
              <TabPanel>
                <WatermarkOverviewSettings />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </div>
      </ModalBody>
    </>
  )
}

export default WatermarkOverview
