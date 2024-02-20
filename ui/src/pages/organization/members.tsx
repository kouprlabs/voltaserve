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
  Button,
  Avatar,
  Portal,
} from '@chakra-ui/react'
import {
  IconDotsVertical,
  IconExit,
  IconUserPlus,
  SectionSpinner,
  PagePagination,
  usePagePagination,
} from '@koupr/ui'
import classNames from 'classnames'
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
  const { id } = useParams()
  const { data: org, error: orgError } = OrganizationAPI.useGetById(
    id,
    swrConfig(),
  )
  const { page, size, steps, setPage, setSize } = usePagePagination({
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
      organizationId: id,
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
        <div className={classNames('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
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
                    <div
                      className={classNames(
                        'flex',
                        'flex-row',
                        'gap-1.5',
                        'items-center',
                      )}
                    >
                      <Avatar name={u.fullName} src={u.picture} />
                      <Text>{u.fullName}</Text>
                    </div>
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
            <PagePagination
              style={{ alignSelf: 'end' }}
              totalElements={list.totalElements}
              totalPages={list.totalPages}
              page={page}
              size={size}
              steps={steps}
              setPage={setPage}
              setSize={setSize}
            />
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
        </div>
      )}
      {list.data.length === 0 && (
        <>
          <div
            className={classNames(
              'flex',
              'items-center',
              'justify-center',
              'h=[300px]',
            )}
          >
            <div
              className={classNames(
                'flex',
                'flex-col',
                'gap-1.5',
                'items-center',
              )}
            >
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
            </div>
          </div>
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
