import { useState } from 'react'
import {
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
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
import InsightsOverviewChart from './insights-overview-chart'
import InsightsOverviewEntities from './insights-overview-entities'
import InsightsOverviewSettings from './insights-overview-settings'
import InsightsOverviewText from './insights-overview-text'

const InsightsOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [isWarningVisible, setIsWarningVisible] = useState(true)
  const { data: summary } = InsightsAPI.useGetMetadata(id)

  return (
    <>
      <ModalBody>
        <div className={cx('flex', 'flex-col', 'gap-1.5', 'w-full')}>
          {summary?.isOutdated && isWarningVisible ? (
            <Alert status="warning">
              <AlertIcon />
              <Box>
                <AlertDescription>
                  These insights are outdated, they originate from an older
                  snapshot. Please navigate to the settings tab to update.
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
              <Tab>Text</Tab>
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
                <InsightsOverviewText />
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
