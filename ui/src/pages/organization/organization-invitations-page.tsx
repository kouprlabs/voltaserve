// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { useCallback, useEffect, useState } from 'react'
import { useLocation, useNavigate, useParams } from 'react-router-dom'
import { Button, useToast } from '@chakra-ui/react'
import {
  DataTable,
  IconDelete,
  IconPersonAdd,
  IconSend,
  Text,
  PagePagination,
  SectionSpinner,
  usePagePagination,
  RelativeDate,
  SectionError,
} from '@koupr/ui'
import cx from 'classnames'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, { SortBy, SortOrder } from '@/client/api/invitation'
import OrganizationAPI from '@/client/api/organization'
import { geEditorPermission } from '@/client/api/permission'
import { swrConfig } from '@/client/options'
import OrganizationInviteMembers from '@/components/organization/organization-invite-members'
import OrganizationStatus from '@/components/organization/organization-status'
import { outgoingInvitationPaginationStorage } from '@/infra/pagination'
import truncateMiddle from '@/lib/helpers/truncate-middle'
import { useAppDispatch } from '@/store/hook'
import { mutateUpdated } from '@/store/ui/outgoing-invitations'

const OrganizationInvitationsPage = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const location = useLocation()
  const { id } = useParams()
  const toast = useToast()
  const { data: org, error: orgError } = OrganizationAPI.useGet(id, swrConfig())
  const { page, size, steps, setPage, setSize } = usePagePagination({
    navigateFn: navigate,
    searchFn: () => location.search,
    storage: outgoingInvitationPaginationStorage(),
  })
  const {
    data: list,
    error: invitationsError,
    mutate,
  } = InvitationAPI.useGetOutgoing(
    {
      organizationId: id,
      page,
      size,
      sortBy: SortBy.DateCreated,
      sortOrder: SortOrder.Desc,
    },
    swrConfig(),
  )
  const [isInviteMembersModalOpen, setIsInviteMembersModalOpen] =
    useState(false)

  const handleResend = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.resend(invitationId)
      toast({
        title: 'Invitation resent',
        status: 'success',
        isClosable: true,
      })
    },
    [toast],
  )

  const handleDelete = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.delete(invitationId)
      await mutate()
    },
    [mutate],
  )

  useEffect(() => {
    if (mutate) {
      dispatch(mutateUpdated(mutate))
    }
  }, [mutate])

  return (
    <>
      {org ? (
        <Helmet>
          <title>{org.name}</title>
        </Helmet>
      ) : null}
      {!list && invitationsError && org && !orgError ? (
        <SectionError text="Failed to load invitations." />
      ) : null}
      {!org && orgError && list && !invitationsError ? (
        <SectionError text="Failed to load organization." />
      ) : null}
      {!list && invitationsError && !org && orgError ? (
        <SectionError text="Failed to load organization and invitations." />
      ) : null}
      {(!list && !invitationsError) || (!org && !orgError) ? (
        <SectionSpinner />
      ) : null}
      {list && list.data.length === 0 && org ? (
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
              <span>This organization has no invitations.</span>
              {geEditorPermission(org.permission) ? (
                <Button
                  leftIcon={<IconPersonAdd />}
                  onClick={() => {
                    setIsInviteMembersModalOpen(true)
                  }}
                >
                  Invite Members
                </Button>
              ) : null}
            </div>
          </div>
          <OrganizationInviteMembers
            open={isInviteMembersModalOpen}
            id={org.id}
            onClose={() => setIsInviteMembersModalOpen(false)}
          />
        </>
      ) : null}
      {list && list.data.length > 0 && org ? (
        <div className={cx('flex', 'flex-col', 'gap-3.5', 'pb-3.5')}>
          <DataTable
            items={list.data}
            columns={[
              {
                title: 'Email',
                renderCell: (i) => <Text>{truncateMiddle(i.email, 50)}</Text>,
              },
              {
                title: 'Status',
                renderCell: (i) => <OrganizationStatus value={i.status} />,
              },
              {
                title: 'Date',
                renderCell: (i) => (
                  <RelativeDate date={new Date(i.createTime)} />
                ),
              },
            ]}
            actions={[
              {
                label: 'Resend',
                icon: <IconSend />,
                onClick: (i) => handleResend(i.id),
              },
              {
                label: 'Delete',
                icon: <IconDelete />,
                isDestructive: true,
                onClick: (i) => handleDelete(i.id),
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
        </div>
      ) : null}
    </>
  )
}

export default OrganizationInvitationsPage
