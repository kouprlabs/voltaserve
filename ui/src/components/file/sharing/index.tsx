import { useMemo } from 'react'
import { useParams } from 'react-router-dom'
import {
  Text,
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
  HStack,
  Tag,
} from '@chakra-ui/react'
import FileAPI, { List } from '@/client/api/file'
import GroupAPI from '@/client/api/group'
import { geOwnerPermission } from '@/client/api/permission'
import UserAPI from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { sharingModalDidClose } from '@/store/ui/files'
import Groups from './groups'
import Users from './users'

type SharingProps = {
  list: List
}

const Sharing = ({ list }: SharingProps) => {
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
  const { data: workspace } = WorkspaceAPI.useGetById(id)
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
        <ModalHeader>Sharing</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Tabs>
            <TabList h="40px">
              <Tab>
                <HStack>
                  <Text>People</Text>
                  {singleFile &&
                  userPermissions &&
                  userPermissions.length > 0 ? (
                    <Tag borderRadius="full">{userPermissions.length}</Tag>
                  ) : null}
                </HStack>
              </Tab>
              <Tab>
                <HStack>
                  <Text>Groups</Text>
                  {singleFile &&
                  groupPermissions &&
                  groupPermissions.length > 0 ? (
                    <Tag borderRadius="full">{groupPermissions.length}</Tag>
                  ) : null}
                </HStack>
              </Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                <Users
                  users={users?.data}
                  permissions={userPermissions}
                  mutateUserPermissions={mutateUserPermissions}
                />
              </TabPanel>
              <TabPanel>
                <Groups
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
