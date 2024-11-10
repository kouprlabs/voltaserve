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
import { Button, Avatar } from '@chakra-ui/react'
import {
  DataTable,
  IconLogout,
  IconPersonAdd,
  PagePagination,
  SectionError,
  SectionPlaceholder,
  SectionSpinner,
  Text,
  usePagePagination,
} from '@koupr/ui'
import cx from 'classnames'
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
    error: listError,
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

  return (
    <>
      {!org && orgError ? (
        <SectionError text="Failed to load organization." />
      ) : null}
      {!org && !orgError ? <SectionSpinner /> : null}
      {org && !orgError ? (
        <>
          {!list && listError ? (
            <SectionError text="Failed to load members." />
          ) : null}
          {!list && !listError ? <SectionSpinner /> : null}
          {list && !listError ? (
            <>
              {list.totalElements > 0 ? (
                <DataTable
                  items={list.data}
                  columns={[
                    {
                      title: 'Full name',
                      renderCell: (u) => (
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
                      ),
                    },
                    {
                      title: 'Email',
                      renderCell: (u) => (
                        <Text>{truncateMiddle(u.email, 50)}</Text>
                      ),
                    },
                  ]}
                  actions={[
                    {
                      label: 'Remove From Organization',
                      icon: <IconLogout />,
                      isDestructive: true,
                      onClick: (u) => {
                        setUserToRemove(u)
                        setIsRemoveMemberModalOpen(true)
                      },
                    },
                  ]}
                  pagination={
                    list.totalPages > 1 ? (
                      <PagePagination
                        totalElements={list.totalElements}
                        totalPages={list.totalPages}
                        page={page}
                        size={size}
                        steps={steps}
                        setPage={setPage}
                        setSize={setSize}
                      />
                    ) : undefined
                  }
                />
              ) : (
                <SectionPlaceholder
                  text="This organization has no members."
                  content={
                    geEditorPermission(org.permission) ? (
                      <Button
                        leftIcon={<IconPersonAdd />}
                        onClick={() => dispatch(inviteModalDidOpen())}
                      >
                        Invite Members
                      </Button>
                    ) : undefined
                  }
                />
              )}
            </>
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
