import { useState } from 'react'
import { useParams } from 'react-router-dom'
import {
  HStack,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  Text,
  VStack,
  Button,
  Center,
  Avatar,
  Portal,
} from '@chakra-ui/react'
import { variables, IconExit, IconUserPlus, SectionSpinner } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import { HiDotsVertical } from 'react-icons/hi'
import GroupAPI from '@/api/group'
import { swrConfig } from '@/api/options'
import { geEditorPermission } from '@/api/permission'
import { User } from '@/api/user'
import AddMember from '@/components/group/add-member'
import RemoveMember from '@/components/group/remove-member'

const GroupMembersPage = () => {
  const params = useParams()
  const groupId = params.id as string
  const { data: group, error: groupError } = GroupAPI.useGetById(
    groupId,
    swrConfig()
  )
  const {
    data: members,
    error: membersError,
    mutate,
  } = GroupAPI.useGetMembers(groupId, swrConfig())
  const [userToRemove, setUserToRemove] = useState<User>()
  const [isAddMembersModalOpen, setIsAddMembersModalOpen] = useState(false)
  const [isRemoveMemberModalOpen, setIsRemoveMemberModalOpen] =
    useState<boolean>(false)

  if (groupError || membersError) {
    return null
  }

  if (!group || !members) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{group.name}</title>
      </Helmet>
      {members.length > 0 && (
        <>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Full name</Th>
                <Th>Email</Th>
                <Th></Th>
              </Tr>
            </Thead>
            <Tbody>
              {members.map((u: User) => (
                <Tr key={u.id}>
                  <Td>
                    <HStack direction="row" spacing={variables.spacing}>
                      <Avatar name={u.fullName} src={u.picture} />
                      <Text>{u.fullName}</Text>
                    </HStack>
                  </Td>
                  <Td>{u.email}</Td>
                  <Td textAlign="right">
                    <Menu>
                      <MenuButton
                        as={IconButton}
                        icon={<HiDotsVertical />}
                        fontSize="18px"
                        variant="ghost"
                        aria-label=""
                      />
                      <Portal>
                        <MenuList>
                          <MenuItem
                            icon={<IconExit />}
                            color="red"
                            isDisabled={!geEditorPermission(group.permission)}
                            onClick={() => {
                              setUserToRemove(u)
                              setIsRemoveMemberModalOpen(true)
                            }}
                          >
                            Remove from group
                          </MenuItem>
                        </MenuList>
                      </Portal>
                    </Menu>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
          {userToRemove && (
            <RemoveMember
              isOpen={isRemoveMemberModalOpen}
              user={userToRemove}
              group={group}
              onCompleted={() => mutate()}
              onClose={() => setIsRemoveMemberModalOpen(false)}
            />
          )}
        </>
      )}
      {members.length === 0 && (
        <>
          <Helmet>
            <title>{group.name}</title>
          </Helmet>
          <Center h="300px">
            <VStack spacing={variables.spacing}>
              <Text>This group has no members.</Text>
              {geEditorPermission(group.permission) && (
                <Button
                  leftIcon={<IconUserPlus />}
                  onClick={() => {
                    setIsAddMembersModalOpen(true)
                  }}
                >
                  Add Members
                </Button>
              )}
            </VStack>
          </Center>
          <AddMember
            open={isAddMembersModalOpen}
            group={group}
            onClose={() => setIsAddMembersModalOpen(false)}
          />
        </>
      )}
    </>
  )
}

export default GroupMembersPage
