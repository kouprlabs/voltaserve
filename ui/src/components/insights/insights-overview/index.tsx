import {
  ModalBody,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from '@chakra-ui/react'
import cx from 'classnames'
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

  if (!id) {
    return null
  }

  return (
    <>
      <ModalBody>
        <div className={cx('w-full')}>
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
