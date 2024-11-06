// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useState } from 'react'
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
  Button,
  Avatar,
  Portal,
} from '@chakra-ui/react'
import {
  IconLogout,
  IconMoreVert,
  IconPersonAdd,
  PagePagination,
  SectionSpinner,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import OrganizationAPI from '@/client/api/organization'
import { geEditorPermission } from '@/client/api/permission'
import UserAPI, { SortBy, SortOrder, User } from '@/client/api/user'
import { swrConfig } from '@/client/options'
import OrganizationInviteMembers from '@/components/organization/organization-invite-members'
import OrganizationRemoveMember from '@/components/organization/organization-remove-member'
import { organizationMemberPaginationStorage } from '@/infra/pagination'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { decodeQuery } from '@/lib/helpers/query'
import { truncateEnd } from '@/lib/helpers/truncate-end'
import truncateMiddle from '@/lib/helpers/truncate-middle'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  inviteModalDidClose,
  inviteModalDidOpen,
} from '@/store/ui/organizations'

const OrganizationMembersPage = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const { id } = useParams()
  const { data: org, error: orgError } = OrganizationAPI.useGet(id, swrConfig())
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: organizationMemberPaginationStorage(),
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
      organizationId: id,
      page,
      size,
      sortBy: SortBy.FullName,
      sortOrder: SortOrder.Asc,
    },
    swrConfig(),
  )
  const isInviteMembersModalOpen = useAppSelector(
    (state) => state.ui.organizations.isInviteModalOpen,
  )
  const [userToRemove, setUserToRemove] = useState<User>()
  const [isRemoveMemberModalOpen, setIsRemoveMemberModalOpen] =
    useState<boolean>(false)

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
      {!list && membersError ? (
        <div
          className={cx('flex', 'items-center', 'justify-center', 'h-[300px]')}
        >
          <span>Failed to load members.</span>
        </div>
      ) : null}
      {!list && !membersError ? <SectionSpinner /> : null}
      {list.data.length > 0 ? (
        <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
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
                      className={cx(
                        'flex',
                        'flex-row',
                        'gap-1.5',
                        'items-center',
                      )}
                    >
                      <Avatar
                        name={u.fullName}
                        src={
                          u.picture
                            ? getPictureUrlById(u.id, u.picture, {
                                organizationId: org.id,
                              })
                            : undefined
                        }
                        className={cx(
                          'border',
                          'border-gray-300',
                          'dark:border-gray-700',
                        )}
                      />
                      <span>{truncateEnd(u.fullName, 50)}</span>
                    </div>
                  </Td>
                  <Td>{truncateMiddle(u.email, 50)}</Td>
                  <Td className={cx('text-right')}>
                    <Menu>
                      <MenuButton
                        as={IconButton}
                        icon={<IconMoreVert />}
                        variant="ghost"
                        aria-label=""
                      />
                      <Portal>
                        <MenuList>
                          <MenuItem
                            icon={<IconLogout />}
                            className={cx('text-red-500')}
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
          {list ? (
            <div className={cx('self-end')}>
              <PagePagination
                totalElements={list.totalElements}
                totalPages={list.totalPages}
                page={page}
                size={size}
                steps={steps}
                setPage={setPage}
                setSize={setSize}
              />
            </div>
          ) : null}
          {userToRemove ? (
            <OrganizationRemoveMember
              isOpen={isRemoveMemberModalOpen}
              user={userToRemove}
              organization={org}
              onCompleted={() => mutate()}
              onClose={() => setIsRemoveMemberModalOpen(false)}
            />
          ) : null}
        </div>
      ) : null}
      {list.data.length === 0 ? (
        <>
          <div
            className={cx(
              'flex',
              'items-center',
              'justify-center',
              'h-[300px]',
            )}
          >
            <div className={cx('flex', 'flex-col', 'gap-1.5', 'items-center')}>
              <span>This organization has no members.</span>
              {geEditorPermission(org.permission) ? (
                <Button
                  leftIcon={<IconPersonAdd />}
                  onClick={() => dispatch(inviteModalDidOpen())}
                >
                  Invite Members
                </Button>
              ) : null}
            </div>
          </div>
          <OrganizationInviteMembers
            open={isInviteMembersModalOpen}
            id={org.id}
            onClose={() => dispatch(inviteModalDidClose())}
          />
        </>
      ) : null}
    </>
  )
}

export default OrganizationMembersPage
