// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useEffect, useState } from 'react'
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
  Text,
  IconPersonAdd,
  PagePagination,
  SectionSpinner,
  usePagePagination,
  SectionError,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import GroupAPI from '@/client/api/group'
import { geEditorPermission } from '@/client/api/permission'
import UserAPI, { SortBy, SortOrder } from '@/client/api/user'
import { User as IdPUser } from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import GroupAddMember from '@/components/group/group-add-member'
import GroupRemoveMember from '@/components/group/group-remove-member'
import { groupMemberPaginationStorage } from '@/infra/pagination'
import { getPictureUrlById } from '@/lib/helpers/picture'
import { decodeQuery } from '@/lib/helpers/query'
import { truncateEnd } from '@/lib/helpers/truncate-end'
import truncateMiddle from '@/lib/helpers/truncate-middle'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/group-members'

const GroupMembersPage = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const location = useLocation()
  const { id } = useParams()
  const { data: group, error: groupError } = GroupAPI.useGet(id, swrConfig())
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: groupMemberPaginationStorage(),
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
      groupId: id,
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

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate])

  return (
    <>
      {group ? (
        <Helmet>
          <title>{group.name}</title>
        </Helmet>
      ) : null}
      {!list && membersError && group && !groupError ? (
        <SectionError text="Failed to load members." />
      ) : null}
      {!group && groupError && list && !membersError ? (
        <SectionError text="Failed to load group." />
      ) : null}
      {!list && membersError && !group && groupError ? (
        <SectionError text="Failed to load group and members." />
      ) : null}
      {(!list && !membersError) || (!group && !groupError) ? (
        <SectionSpinner />
      ) : null}
      {list && list.data.length > 0 && group ? (
        <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
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
                              groupId: group.id,
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
                renderCell: (u) => <Text>{truncateMiddle(u.email, 50)}</Text>,
              },
            ]}
            actions={[
              {
                label: 'Remove From Group',
                icon: <IconLogout />,
                isDestructive: true,
                onClick: (u) => {
                  setUserToRemove(u)
                  setIsRemoveMemberModalOpen(true)
                },
              },
            ]}
          />
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
            <GroupRemoveMember
              isOpen={isRemoveMemberModalOpen}
              user={userToRemove}
              group={group}
              onCompleted={() => mutate()}
              onClose={() => setIsRemoveMemberModalOpen(false)}
            />
          ) : null}
        </div>
      ) : null}
      {list && list.data.length === 0 && group ? (
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
              <span>This group has no members.</span>
              {geEditorPermission(group.permission) ? (
                <Button
                  leftIcon={<IconPersonAdd />}
                  onClick={() => {
                    setIsAddMembersModalOpen(true)
                  }}
                >
                  Add Members
                </Button>
              ) : null}
            </div>
          </div>
          <GroupAddMember
            open={isAddMembersModalOpen}
            group={group}
            onClose={() => setIsAddMembersModalOpen(false)}
          />
        </>
      ) : null}
    </>
  )
}

export default GroupMembersPage
