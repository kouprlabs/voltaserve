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
import SecurityOverviewSettings from './security-overview-settings'
import SecurityOverviewWatermark from './security-overview-watermark'

const SecurityOverview = () => {
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
              <Tab>Watermark</Tab>
              <Tab>Settings</Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <SecurityOverviewWatermark />
              </TabPanel>
              <TabPanel>
                <SecurityOverviewSettings />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </div>
      </ModalBody>
    </>
  )
}

export default SecurityOverview
