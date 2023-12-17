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
import FileAPI from '@/client/api/file'
import GroupAPI from '@/client/api/group'
import UserAPI from '@/client/api/user'
import WorkspaceAPI from '@/client/api/workspace'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import { selectionUpdated, sharingModalDidClose } from '@/store/ui/files'
import FormSkeleton from './form-skeleton'
import Groups from './groups'
import Users from './users'

const Sharing = () => {
  const params = useParams()
  const dispatch = useAppDispatch()
  const selection = useAppSelector((state) => state.ui.files.selection)
  const isModalOpen = useAppSelector((state) => state.ui.files.isShareModalOpen)
  const { data: workspace } = WorkspaceAPI.useGetById(params.id as string)
  const { data: users } = UserAPI.useList({
    organizationId: workspace?.organization.id,
  })
  const { data: groups } = GroupAPI.useList({
    organizationId: workspace?.organization.id,
  })
  const { data: userPermissions, mutate: mutateUserPermissions } =
    FileAPI.useGetUserPermissions(selection[0])
  const { data: groupPermissions, mutate: mutateGroupPermissions } =
    FileAPI.useGetGroupPermissions(selection[0])

  return (
    <Modal
      size="xl"
      isOpen={isModalOpen}
      onClose={() => {
        dispatch(selectionUpdated([]))
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
                  {userPermissions && userPermissions.length > 0 ? (
                    <Tag borderRadius="full">{userPermissions.length}</Tag>
                  ) : null}
                </HStack>
              </Tab>
              <Tab>
                <HStack>
                  <Text>Groups</Text>
                  {groupPermissions && groupPermissions.length > 0 ? (
                    <Tag borderRadius="full">{groupPermissions.length}</Tag>
                  ) : null}
                </HStack>
              </Tab>
            </TabList>
            <TabPanels>
              <TabPanel>
                {users && userPermissions ? (
                  <Users
                    users={users.data}
                    userPermissions={userPermissions}
                    mutateUserPermissions={mutateUserPermissions}
                  />
                ) : (
                  <FormSkeleton />
                )}
              </TabPanel>
              <TabPanel>
                {groups && groupPermissions ? (
                  <Groups
                    groups={groups.data}
                    groupPermissions={groupPermissions}
                    mutateGroupPermissions={mutateGroupPermissions}
                  />
                ) : (
                  <FormSkeleton />
                )}
              </TabPanel>
            </TabPanels>
          </Tabs>
        </ModalBody>
      </ModalContent>
    </Modal>
  )
}

export default Sharing
