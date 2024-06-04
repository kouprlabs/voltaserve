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
import MosaicOverviewArtifacts from './mosaic-overview-artifacts'
import MosaicOverviewSettings from './mosaic-overview-settings'

const MosaicOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [isWarningVisible, setIsWarningVisible] = useState(true)
  const { data: info } = MosaicAPI.useGetInfo(id, swrConfig())

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
                  This mosaic is outdated, it can be updated in the settings.
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
      </ModalBody>
    </>
  )
}

export default MosaicOverview
