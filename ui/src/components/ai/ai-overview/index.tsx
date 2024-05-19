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
import { useAppDispatch } from '@/store/hook'
import { modalDidClose } from '@/store/ui/ai'
import AiOverviewLanguage from './ai-overview-language'
import AiOverviewNamedEntities from './ai-overview-named-entities'
import AiOverviewText from './ai-overview-text'

const AiOverview = () => {
  const dispatch = useAppDispatch()
  const [activeTab, setActiveTab] = useState(0)

  return (
    <>
      <ModalBody>
        <div className={cx('w-full', 'pb-1.5')}>
          <Tabs
            variant="solid-rounded"
            colorScheme="gray"
            index={activeTab}
            className={cx('pb-2.5')}
          >
            <TabList>
              <Tab onClick={() => setActiveTab(0)}>Text</Tab>
              <Tab onClick={() => setActiveTab(1)}>Language</Tab>
              <Tab onClick={() => setActiveTab(2)}>Named Entities</Tab>
            </TabList>
          </Tabs>
          {activeTab === 0 ? <AiOverviewText /> : null}
          {activeTab === 1 ? <AiOverviewLanguage /> : null}
          {activeTab === 2 ? <AiOverviewNamedEntities /> : null}
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

export default AiOverview
