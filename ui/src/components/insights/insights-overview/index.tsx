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
import InsightsAPI from '@/client/api/insights'
import { useAppSelector } from '@/store/hook'
import InsightsOverviewArtifacts from './insights-overview-artifacts'
import InsightsOverviewChart from './insights-overview-chart'
import InsightsOverviewEntities from './insights-overview-entities'
import InsightsOverviewSettings from './insights-overview-settings'

const InsightsOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [isWarningVisible, setIsWarningVisible] = useState(true)
  const { data: info } = InsightsAPI.useGetInfo(id)

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
                  These insights are outdated, it can be updated in the
                  settings.
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
              <Tab>Chart</Tab>
              <Tab>Entities</Tab>
              <Tab>Artifacts</Tab>
              <Tab>Settings</Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <InsightsOverviewChart />
              </TabPanel>
              <TabPanel>
                <InsightsOverviewEntities />
              </TabPanel>
              <TabPanel>
                <InsightsOverviewArtifacts />
              </TabPanel>
              <TabPanel>
                <InsightsOverviewSettings />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </div>
      </ModalBody>
    </>
  )
}

export default InsightsOverview
