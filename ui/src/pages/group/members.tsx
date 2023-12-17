import { useState } from 'react'
import { useParams, useSearchParams } from 'react-router-dom'
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
  Stack,
} from '@chakra-ui/react'
import { variables, IconExit, IconUserPlus, SectionSpinner } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import { HiDotsVertical } from 'react-icons/hi'
import GroupAPI from '@/client/api/group'
import { geEditorPermission } from '@/client/api/permission'
import UserAPI, { SortBy, SortOrder } from '@/client/api/user'
import { User as IdPUser } from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import PagePagination, {
  usePagePagination,
} from '@/components/common/page-pagination'
import AddMember from '@/components/group/add-member'
import RemoveMember from '@/components/group/remove-member'
import { decodeQuery } from '@/helpers/query'

const GroupMembersPage = () => {
  const params = useParams()
  const groupId = params.id as string
  const { data: group, error: groupError } = GroupAPI.useGetById(
    groupId,
    swrConfig(),
  )
  const { page, size, onPageChange, onSizeChange } = usePagePagination({
    localStoragePrefix: 'voltaserve',
    localStorageNamespace: 'group_member',
  })
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const {
    data: list,
    error: membersError,
    mutate,
  } = UserAPI.useList(
    {
      query,
      groupId,
      page,
      size,
      sortBy: SortBy.FullName,
      sortOrder: SortOrder.Asc,
    },
    swrConfig(),
  )
  const [userToRemove, setUserToRemove] = useState<IdPUser>()
  const [isAddMembersModalOpen, setIsAddMembersModalOpen] = useState(false)
  const [isRemoveMemberModalOpen, setIsRemoveMemberModalOpen] =
    useState<boolean>(false)

  if (groupError || membersError) {
    return null
  }

  if (!group || !list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{group.name}</title>
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
                            Remove From Group
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
                onPageChange={onPageChange}
                onSizeChange={onSizeChange}
              />
            </HStack>
          )}
          {userToRemove && (
            <RemoveMember
              isOpen={isRemoveMemberModalOpen}
              user={userToRemove}
              group={group}
              onCompleted={() => mutate()}
              onClose={() => setIsRemoveMemberModalOpen(false)}
            />
          )}
        </Stack>
      )}
      {list.data.length === 0 && (
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
