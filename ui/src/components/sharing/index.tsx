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
import { FileAPI } from '@/client/api/file'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import SharingGroupForm from './sharing-group-form'
import SharingUserForm from './sharing-user-form'

const Sharing = () => {
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector(
    (state) => state.ui.files.isSharingModalOpen,
  )
  const isSingleSelection = selection.length === 1
  const { data: userPermissions } = FileAPI.useGetUserPermissions(
    isSingleSelection ? selection[0] : undefined,
    swrConfig(),
  )
  const { data: groupPermissions } = FileAPI.useGetGroupPermissions(
    isSingleSelection ? selection[0] : undefined,
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
                  {userPermissions && userPermissions.length > 0 ? (
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
                  {groupPermissions && groupPermissions.length > 0 ? (
                    <Tag className={cx('rounded-full')}>
                      {groupPermissions.length}
                    </Tag>
                  ) : null}
                </div>
              </Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <SharingUserForm />
              </TabPanel>
              <TabPanel>
                <SharingGroupForm />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}

export default Sharing
