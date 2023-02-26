import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import {
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
  HStack,
  Center,
  VStack,
  Button,
  Avatar,
  Portal,
} from '@chakra-ui/react'
import {
  variables,
  IconDotsVertical,
  IconExit,
  IconUserPlus,
  SectionSpinner,
} from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import { swrConfig } from '@/api/options'
import OrganizationAPI from '@/api/organization'
import { geEditorPermission } from '@/api/permission'
import { User } from '@/api/user'
import OrganizationInviteMembers from '@/components/organization/invite-members'
import OrganizationRemoveMember from '@/components/organization/remove-member'

const OrganizationMembersPage = () => {
  const params = useParams()
  const orgId = params.id as string
  const invite = Boolean(params.invite as string)
  const { data: org, error: orgError } = OrganizationAPI.useGetById(
    orgId,
    swrConfig()
  )
  const {
    data: members,
    error: membersError,
    mutate,
  } = OrganizationAPI.useGetMembers(orgId, swrConfig())
  const [userToRemove, setUserToRemove] = useState<User>()
  const [isInviteMembersModalOpen, setIsInviteMembersModalOpen] =
    useState(false)
  const [isRemoveMemberModalOpen, setIsRemoveMemberModalOpen] =
    useState<boolean>(false)

  useEffect(() => {
    if (invite) {
      setIsInviteMembersModalOpen(true)
    }
  }, [invite])

  if (membersError || orgError) {
    return null
  }

  if (!members || !org) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{org.name}</title>
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
                        icon={<IconDotsVertical />}
                        variant="ghost"
                        aria-label=""
                      />
                      <Portal>
                        <MenuList>
                          <MenuItem
                            icon={<IconExit />}
                            color="red"
                            isDisabled={!geEditorPermission(org.permission)}
                            onClick={() => {
                              setUserToRemove(u)
                              setIsRemoveMemberModalOpen(true)
                            }}
                          >
                            Remove from organization
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
            <OrganizationRemoveMember
              isOpen={isRemoveMemberModalOpen}
              user={userToRemove}
              organization={org}
              onCompleted={() => mutate()}
              onClose={() => setIsRemoveMemberModalOpen(false)}
            />
          )}
        </>
      )}
      {members.length === 0 && (
        <>
          <Center h="300px">
            <VStack spacing={variables.spacing}>
              <Text>This organization has no members.</Text>
              {geEditorPermission(org.permission) && (
                <Button
                  leftIcon={<IconUserPlus />}
                  onClick={() => {
                    setIsInviteMembersModalOpen(true)
                  }}
                >
                  Invite members
                </Button>
              )}
            </VStack>
          </Center>
          <OrganizationInviteMembers
            open={isInviteMembersModalOpen}
            id={org.id}
            onClose={() => setIsInviteMembersModalOpen(false)}
          />
        </>
      )}
    </>
  )
}

export default OrganizationMembersPage
