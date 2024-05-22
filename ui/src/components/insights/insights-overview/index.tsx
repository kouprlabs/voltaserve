import { useState } from 'react'
import {
  Button,
  ModalBody,
  ModalFooter,
  Tab,
  TabList,
  Tabs,
} from '@chakra-ui/react'
import cx from 'classnames'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { modalDidClose } from '@/store/ui/insights'
import InsightsOverviewEntities from './insights-overview-entities'
import InsightsOverviewLanguage from './insights-overview-language'
import InsightsOverviewSettings from './insights-overview-settings'
import InsightsOverviewText from './insights-overview-text'

const InsightsOverview = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const [activeTab, setActiveTab] = useState(0)

  if (!id) {
    return null
  }

  return (
    <>
      <ModalBody>
        <div className={cx('w-full')}>
          <Tabs
            variant="solid-rounded"
            colorScheme="gray"
            index={activeTab}
            className={cx('pb-2.5')}
          >
            <TabList>
              <Tab onClick={() => setActiveTab(0)}>Language</Tab>
              <Tab onClick={() => setActiveTab(1)}>Text</Tab>
              <Tab onClick={() => setActiveTab(2)}>Entities</Tab>
              <Tab onClick={() => setActiveTab(3)}>Settings</Tab>
            </TabList>
          </Tabs>
          {activeTab === 0 ? <InsightsOverviewLanguage /> : null}
          {activeTab === 1 ? <InsightsOverviewText /> : null}
          {activeTab === 2 ? <InsightsOverviewEntities /> : null}
          {activeTab === 3 ? <InsightsOverviewSettings /> : null}
        </div>
      </ModalBody>
      <ModalFooter>
        <Button
          type="button"
          variant="outline"
          colorScheme="blue"
          onClick={() => dispatch(modalDidClose())}
        >
          Close
        </Button>
      </ModalFooter>
    </>
  )
}

export default InsightsOverview
