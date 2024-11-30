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
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  Tabs,
  TabList,
  TabPanels,
  Tab,
  TabPanel,
  Tag,
} from '@chakra-ui/react'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import { geOwnerPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import SharingGroupOverview from './sharing-group-overview'
import SharingUserOverview from './sharing-user-overview'

const Sharing = () => {
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isShareModalOpen)
  const { data: file } = FileAPI.useGet(selection[0], swrConfig())
  const { data: userPermissions } = FileAPI.useGetUserPermissions(
    file && geOwnerPermission(file.permission) ? file.id : undefined,
    swrConfig(),
  )
  const { data: groupPermissions } = FileAPI.useGetGroupPermissions(
    file && geOwnerPermission(file.permission) ? file.id : undefined,
    swrConfig(),
  )

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => {
        dispatch(sharingModalDidClose())
      }}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        {selection.length > 1 ? (
          <ModalHeader>Sharing ({selection.length}) Items</ModalHeader>
        ) : (
          <ModalHeader>Sharing</ModalHeader>
        )}
        <ModalCloseButton />
        <ModalBody>
          <Tabs colorScheme="gray">
            <TabList className={cx('h-[40px]')}>
              <Tab>
                <div
                  className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}
                >
                  <span>Users</span>
                  {file && userPermissions && userPermissions.length > 0 ? (
                    <Tag className={cx('rounded-full')}>
                      {userPermissions.length}
                    </Tag>
                  ) : null}
                </div>
              </Tab>
              <Tab>
                <div
                  className={cx('flex', 'flex-row', 'items-center', 'gap-0.5')}
                >
                  <span>Groups</span>
                  {file && groupPermissions && groupPermissions.length > 0 ? (
                    <Tag className={cx('rounded-full')}>
                      {groupPermissions.length}
                    </Tag>
                  ) : null}
                </div>
              </Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <SharingUserOverview />
              </TabPanel>
              <TabPanel>
                <SharingGroupOverview />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}

export default Sharing
