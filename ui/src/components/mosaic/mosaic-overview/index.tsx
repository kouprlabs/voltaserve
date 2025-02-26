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
import MosaicOverviewSettings from './mosaic-overview-settings'

const MosaicOverview = () => (
  <>
    <ModalBody>
      <Tabs colorScheme="gray">
        <TabList>
          <Tab>Settings</Tab>
        </TabList>
        <TabPanels>
          <TabPanel>
            <MosaicOverviewSettings />
          </TabPanel>
        </TabPanels>
      </Tabs>
    </ModalBody>
  </>
)

export default MosaicOverview
