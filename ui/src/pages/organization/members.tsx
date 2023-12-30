import { useEffect, useState } from 'react'
import {
  useLocation,
  useNavigate,
  useParams,
  useSearchParams,
} from 'react-router-dom'
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
  Stack,
} from '@chakra-ui/react'
import {
  variables,
  IconDotsVertical,
  IconExit,
  IconUserPlus,
  SectionSpinner,
  PagePagination,
  usePagePagination,
} from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI from '@/client/api/organization'
import { geEditorPermission } from '@/client/api/permission'
import UserAPI, { SortBy, SortOrder, User } from '@/client/api/user'
import { swrConfig } from '@/client/options'
import InviteMembers from '@/components/organization/invite-members'
import RemoveMember from '@/components/organization/remove-member'
import { decodeQuery } from '@/helpers/query'
import { organizationMemberPaginationStorage } from '@/infra/pagination'

const OrganizationMembersPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const params = useParams()
  const organizationId = params.id as string
  const { data: org, error: orgError } = OrganizationAPI.useGetById(
    organizationId,
    swrConfig(),
  )
  const { page, size, handlePageChange, setSize } = usePagePagination({
    navigate,
    location,
    storage: organizationMemberPaginationStorage(),
  })
  const [searchParams] = useSearchParams()
  const invite = Boolean(searchParams.get('invite') as string)
  const query = decodeQuery(searchParams.get('q') as string)
  const {
    data: list,
    error: membersError,
    mutate,
  } = UserAPI.useList(
    {
      query,
      organizationId,
      page,
      size,
      sortBy: SortBy.FullName,
      sortOrder: SortOrder.Asc,
    },
    swrConfig(),
  )
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

  if (!list || !org) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{org.name}</title>
      </Helmet>
      {list.data.length > 0 && (
        <Stack
          direction="column"
          spacing={variables.spacing2Xl}
          pb={variables.spacing2Xl}
        >
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Full name</Th>
                <Th>Email</Th>
                <Th></Th>
              </Tr>
            </Thead>
            <Tbody>
              {list.data.map((u) => (
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
                            Remove From Organization
                          </MenuItem>
                        </MenuList>
                      </Portal>
                    </Menu>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
          {list && (
            <HStack alignSelf="end">
              <PagePagination
                totalPages={list.totalPages}
                page={page}
                size={size}
                handlePageChange={handlePageChange}
                setSize={setSize}
              />
            </HStack>
          )}
          {userToRemove && (
            <RemoveMember
              isOpen={isRemoveMemberModalOpen}
              user={userToRemove}
              organization={org}
              onCompleted={() => mutate()}
              onClose={() => setIsRemoveMemberModalOpen(false)}
            />
          )}
        </Stack>
      )}
      {list.data.length === 0 && (
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
                  Invite Members
                </Button>
              )}
            </VStack>
          </Center>
          <InviteMembers
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
