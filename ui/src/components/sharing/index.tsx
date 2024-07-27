// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useMemo } from 'react'
import { useParams } from 'react-router-dom'
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
import FileAPI, { List } from '@/client/api/file'
import GroupAPI from '@/client/api/group'
import { geOwnerPermission } from '@/client/api/permission'
import UserAPI from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import SharingGroups from './sharing-groups'
import SharingUsers from './sharing-users'

export type FileSharingProps = {
  list: List
}

const Sharing = ({ list }: FileSharingProps) => {
  const { id } = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isShareModalOpen)
  const singleFile = useMemo(() => {
    if (selection.length === 1) {
      return list.data.find((e) => e.id === selection[0])
    } else {
      return undefined
    }
  }, [list.data, selection])
  const { data: workspace } = WorkspaceAPI.useGet(id)
  const { data: users } = UserAPI.useList({
    organizationId: workspace?.organization.id,
  })
  const { data: groups } = GroupAPI.useList({
    organizationId: workspace?.organization.id,
  })
  const { data: userPermissions, mutate: mutateUserPermissions } =
    FileAPI.useGetUserPermissions(
      singleFile && geOwnerPermission(singleFile.permission)
        ? singleFile.id
        : undefined,
    )
  const { data: groupPermissions, mutate: mutateGroupPermissions } =
    FileAPI.useGetGroupPermissions(
      singleFile && geOwnerPermission(singleFile.permission)
        ? singleFile.id
        : undefined,
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
          <ModalHeader>Sharing {selection.length} Items(s)</ModalHeader>
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
                  {singleFile &&
                  userPermissions &&
                  userPermissions.length > 0 ? (
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
                  {singleFile &&
                  groupPermissions &&
                  groupPermissions.length > 0 ? (
                    <Tag className={cx('rounded-full')}>
                      {groupPermissions.length}
                    </Tag>
                  ) : null}
                </div>
              </Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <SharingUsers
                  users={users?.data}
                  permissions={userPermissions}
                  mutateUserPermissions={mutateUserPermissions}
                />
              </TabPanel>
              <TabPanel>
                <SharingGroups
                  groups={groups?.data}
                  permissions={groupPermissions}
                  mutateGroupPermissions={mutateGroupPermissions}
                />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}

export default Sharing
