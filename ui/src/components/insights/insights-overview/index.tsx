// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import {
  ModalBody,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from '@chakra-ui/react'
import { SectionError, SectionSpinner } from '@koupr/ui'
import cx from 'classnames'
import { errorToString, FileAPI, swrConfig } from '@/client'
import { useAppSelector } from '@/store/hook'
import InsightsOverviewEntities from './insights-overview-entities'
import InsightsOverviewSettings from './insights-overview-settings'
import InsightsOverviewSummary from './insights-overview-summary'

const InsightsOverview = () => {
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const {
    data: file,
    error: fileError,
    isLoading: fileIsLoading,
  } = FileAPI.useGet(id, swrConfig())
  const fileIsReady = file && !fileError

  return (
    <>
      {fileIsLoading ? <SectionSpinner /> : null}
      {fileError ? <SectionError text={errorToString(fileError)} /> : null}
      {fileIsReady ? (
        <ModalBody>
          <div className={cx('flex', 'flex-col', 'gap-1.5', 'w-full')}>
            <Tabs colorScheme="gray">
              <TabList>
                <Tab>Summary</Tab>
                {file.snapshot?.capabilities.entities ? (
                  <Tab>Entities</Tab>
                ) : null}
                <Tab>Settings</Tab>
              </TabList>
              <TabPanels>
                <TabPanel>
                  <InsightsOverviewSummary />
                </TabPanel>
                {file.snapshot?.capabilities.entities ? (
                  <TabPanel>
                    <InsightsOverviewEntities />
                  </TabPanel>
                ) : null}
                <TabPanel>
                  <InsightsOverviewSettings />
                </TabPanel>
              </TabPanels>
            </Tabs>
          </div>
        </ModalBody>
      ) : null}
    </>
  )
}

export default InsightsOverview
