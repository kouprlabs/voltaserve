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
import WatermarkOverviewArtifacts from './watermark-overview-artifacts'
import WatermarkOverviewSettings from './watermark-overview-settings'

const WatermarkOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [isWarningVisible, setIsWarningVisible] = useState(true)
  const { data: info } = WatermarkAPI.useGetInfo(id, swrConfig())

  if (!info) {
    return null
  }

  return (
    <>
      <ModalBody>
        <div className={cx('flex', 'flex-col', 'gap-1.5', 'w-full')}>
          {info.metadata?.isOutdated && isWarningVisible ? (
            <Alert status="warning" className={cx('flex')}>
              <AlertIcon />
              <Box className={cx('grow')}>
                <AlertDescription>
                  The watermark is applied on an older snapshot. You can apply
                  it on the active snapshot from the settings.
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
                <WatermarkOverviewArtifacts />
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
